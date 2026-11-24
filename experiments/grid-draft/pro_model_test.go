package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOracleImportBuildsSeriesBeforeDroppingUnusableGames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oracle.csv")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := csv.NewWriter(file)
	fields := append(append([]string{}, oracleColumns...), "firstPick", "patch")
	header := append([]string{}, fields...)
	header[0] = "\ufeff" + header[0]
	if err := writer.Write(header); err != nil {
		t.Fatal(err)
	}
	recordGame := func(league, id, patch string, number int, start time.Time) {
		t.Helper()
		for side, sideName := range proSideNames {
			for role, position := range append(append([]string{}, proRoles...), "team") {
				row := map[string]string{
					"gameid": id, "league": league, "year": "2026", "split": "Split 1",
					"date": start.Format("2006-01-02 15:04:05"), "game": fmt.Sprint(number), "patch": patch,
					"side": sideName, "position": position, "teamid": fmt.Sprintf("t%d", (side+number)%2),
					"playerid": fmt.Sprintf("p%d%d", side, role), "champion": fmt.Sprintf("c%d%d", side, role),
					"result": fmt.Sprint(1 - side), "firstPick": fmt.Sprint(side), "ban1": "Dr. Mundo",
				}
				record := make([]string, len(header))
				for i, column := range fields {
					record[i] = row[column]
				}
				if err := writer.Write(record); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	start := time.Date(2026, 2, 1, 8, 0, 0, 0, time.UTC)
	recordGame("LCK", "k1", "16.02", 1, start)
	recordGame("LEC", "e1", "16.02", 1, start)
	recordGame("LCS", "s1", "16.02", 1, start)
	recordGame("LDL", "ignored", "16.02", 1, start)
	recordGame("LPL", "native_game_1", "16.02", 1, start)
	recordGame("LPL", "native_game_2", "", 2, start.Add(time.Hour))
	writer.Flush()
	if err := writer.Error(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	dataset, err := importProGames([]string{path, path}, nil)
	if err != nil {
		t.Fatal(err)
	}
	games, coverage := dataset.games, dataset.coverage
	if len(games) != 4 || coverage["LPL"].Rejected["missing patch outside verified league windows"] != 1 ||
		coverage["LCK"].Eligible != 1 || coverage["LEC"].Eligible != 1 || coverage["LCS"].Eligible != 1 {
		t.Fatalf("unexpected tier-one import: %+v, %+v", games, coverage)
	}
	for _, game := range games {
		if game.patch != (proPatch{16, 2}) || game.firstPick != -1 || !game.bans["drmundo"] {
			t.Fatalf("incorrect draft normalization: %+v", game)
		}
		if game.league == "LPL" && !game.available.Equal(start.Add(5*time.Hour)) {
			t.Fatalf("unusable game shortened the series cutoff: %s", game.available)
		}
	}
	windows := map[string][]proPatchWindow{
		"LPL": {{Patch: proPatch{16, 2}, Start: start, End: start.Add(2 * time.Hour)}},
		"LCK": {{Patch: proPatch{16, 3}, Start: start, End: start.Add(2 * time.Hour)}},
	}
	withCalendar, err := importProGames([]string{path}, windows)
	if err != nil {
		t.Fatal(err)
	}
	if len(withCalendar.games) != 5 || withCalendar.coverage["LPL"].PatchSources["scheduled"] != 1 ||
		withCalendar.coverage["LPL"].PatchSources["recorded"] != 1 || withCalendar.coverage["LCK"].ScheduleConflicts != 1 {
		t.Fatalf("unexpected scheduled patch resolution: %+v", withCalendar.coverage)
	}
	for _, game := range withCalendar.games {
		if game.patch != (proPatch{16, 2}) {
			t.Fatalf("calendar altered a recorded patch: %+v", game)
		}
	}
}

func TestProTrainingUsesPriorSeriesComparablePatchesAndLeagueBudget(t *testing.T) {
	cutoff := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	row := func(id, league string, patch proPatch, days int) proGame {
		start := cutoff.AddDate(0, 0, -days)
		return proGame{id: id, league: league, patch: patch, start: start, available: start.Add(historyLag)}
	}
	local := row("local", "LCK", proPatch{16, 9}, 1)
	older := row("older", "LCK", proPatch{16, 8}, 1)
	overlap := row("overlap", "LCK", proPatch{16, 9}, 1)
	overlap.available = cutoff
	games := []proGame{local, older, overlap, row("future", "LPL", proPatch{16, 10}, 1),
		row("prior-season", "LPL", proPatch{15, 9}, 1), row("stale", "LPL", proPatch{16, 9}, 181)}
	for i := range 50 {
		games = append(games, row(fmt.Sprintf("donor-%d", i), "LPL", proPatch{16, 9}, 1))
	}
	target := proContext{"LCK", proPatch{16, 9}, cutoff}
	reset := []proPatch{{16, 9}, {16, 10}}
	dataset := proDataset{games: games}
	training := proTraining(dataset, target, reset, proModelSpecs[3])
	if len(training) != 51 || training[0].game.id != "local" {
		t.Fatalf("selected %d games including %+v", len(training), training[0].game)
	}
	var donorWeight float64
	for _, donor := range training[1:] {
		donorWeight += donor.weight
	}
	if math.Abs(donorWeight-6) > 1e-10 {
		t.Fatalf("donor league weight = %g, want 6", donorWeight)
	}
	earlier := proTraining(dataset, target, nil, proModelSpecs[1])
	if len(earlier) != 2 || earlier[1].game.id != "older" || earlier[1].weight >= earlier[0].weight {
		t.Fatalf("league-only training did not discount the older patch: %+v", earlier)
	}
}

func TestProDraftComparisonKeepsReferenceAndLegalityExplicit(t *testing.T) {
	const focus = 2
	start := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	lineups := [2]proQueryLineup{
		{Team: "blue-team", Roles: map[string]proPick{}},
		{Team: "red-team", Roles: map[string]proPick{}},
	}
	for side := range lineups {
		for i, role := range proRoles {
			lineups[side].Roles[role] = proPick{Player: fmt.Sprintf("player%d%d", side, i), Champion: fmt.Sprintf("champ%d%d", side, i)}
		}
	}
	lineups[0].Roles["mid"] = proPick{Player: "player02"}
	query := proQuery{
		League: "LCK", Patch: "26.8", AsOf: start.AddDate(0, 0, 14), FirstPickSide: "red",
		FocusSide: "blue", FocusRole: "mid", Blue: lineups[0], Red: lineups[1],
		Unavailable: []string{"Dr. Mundo"}, Reference: "reference",
		Candidates: []string{"reference", "alternative", "dr mundo", "unseen"},
	}
	target, _, err := query.prepare()
	if err != nil {
		t.Fatal(err)
	}
	if target.firstPick != -1 {
		t.Fatal("blue focus side overrode red first pick")
	}
	var games []proGame
	for i := range 40 {
		g := target
		g.id, g.series = fmt.Sprint(i), fmt.Sprint(i/2)
		g.start, g.available = start, start.Add(historyLag)
		g.sides[0].Picks[focus].Champion = "reference"
		if i%2 == 0 {
			g.blueWon = 1
			g.sides[0].Picks[focus].Champion = "alternative"
		}
		games = append(games, g)
	}
	query.Patch = ""
	dataset := proDataset{
		games: games,
		calendar: proPatchCalendar{Scheduled: map[string][]proPatchWindow{
			"LCK": {{Patch: proPatch{16, 8}, Start: start, End: query.AsOf.Add(time.Hour)}},
		}},
	}
	decision, err := compareProDrafts(dataset, nil, query)
	if err != nil {
		t.Fatal(err)
	}
	ref, candidate, banned, unseen := decision.Options[0], decision.Options[1], decision.Options[2], decision.Options[3]
	const full = "cross_league_cross_patch"
	if decision.Status != "scored" || decision.PatchSource != "scheduled" || decision.Patch != "16.8" ||
		ref.Delta[full] != 0 || candidate.Delta[full] <= 0 || candidate.Bootstrap95 == nil {
		t.Fatalf("expected a supported contrast: %+v", decision)
	}
	if candidate.Evidence.PlayerRoleGames != 20 || candidate.Evidence.OpponentMatchupGames != 20 ||
		banned.Status != "unavailable" || banned.Probability != nil ||
		unseen.Status != "insufficient_role_evidence" || unseen.Probability != nil {
		t.Fatalf("unexpected option evidence or legality: %+v", decision.Options)
	}
	var contribution float64
	for _, delta := range candidate.LogOddsVsReference {
		contribution += delta
	}
	logit := func(p float64) float64 { return math.Log(p / (1 - p)) }
	if math.Abs(contribution-(logit(candidate.Probability[full])-logit(ref.Probability[full]))) > 1e-9 {
		t.Fatalf("explanations do not sum to the fitted contrast: %g", contribution)
	}
}
