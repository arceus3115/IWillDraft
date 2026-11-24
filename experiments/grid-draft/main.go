package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	grid "github.com/arceus3115/grid-framework"
)

func main() {
	ctx := context.Background()
	data := flag.String("data", os.Getenv("STL_DATA_DIR"), "local GRID starter data root")
	replay := flag.String("replay", "", "completed five-game GRID series ID to audit")
	event := flag.String("opening-event", "LCS - Split", "ordered-draft event prefix")
	year := flag.Int("opening-year", 2026, "ordered-draft season")
	api := flag.String("api-series", "", "fetch one completed series with GRID_API_KEY")
	var proPaths []string
	flag.Func("pro-csv", "local Oracle's Elixir CSV (repeat per season)", func(path string) error {
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("empty CSV path")
		}
		proPaths = append(proPaths, path)
		return nil
	})
	proYear := flag.Int("pro-year", 2026, "target season for the tier-one backtest")
	proScenario := flag.String("pro-scenario", "", "JSON draft comparison scenario instead of a backtest")
	proResets := flag.String("pro-patch-resets", "", "comma-separated patch boundaries that reset draft history")
	proCalendar := flag.String("pro-patch-calendar", "", "CSV of competitive patch windows (league,patch,start,end)")
	mechanics := flag.String("patch-mechanics", "", "JSON scenario for a controlled item damage comparison")
	flag.Parse()
	var err error
	switch {
	case *mechanics != "" && (len(proPaths) > 0 || *proScenario != "" || *proResets != "" || *proCalendar != "" || *api != "" || *replay != ""):
		err = fmt.Errorf("run patch mechanics separately from the stage models and replay audit")
	case *mechanics != "":
		err = runPatchMechanics(*mechanics)
	case len(proPaths) > 0 && (*api != "" || *replay != ""):
		err = fmt.Errorf("run the tier-one model and GRID replay or API audit separately")
	case len(proPaths) > 0:
		err = runProModel(proPaths, *proYear, *proScenario, *proResets, *proCalendar)
	case *proScenario != "" || *proResets != "" || *proCalendar != "" || *proYear != 2026:
		err = fmt.Errorf("provide at least one -pro-csv for the tier-one model")
	case *api != "":
		err = checkAPI(ctx, *api, *data)
	case *data != "":
		err = runArchive(ctx, *data, *replay, *event, *year)
	default:
		err = fmt.Errorf("provide -pro-csv, -patch-mechanics, -data or STL_DATA_DIR, or use -api-series")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runProModel(paths []string, year int, scenario, resetFlags, calendarPath string) error {
	resets, err := parsePatchResets(resetFlags)
	if err != nil {
		return err
	}
	scheduled, err := readProPatchWindows(calendarPath)
	if err != nil {
		return err
	}
	dataset, err := importProGames(paths, scheduled)
	if err != nil {
		return err
	}
	if len(dataset.games) == 0 {
		return fmt.Errorf("no eligible tier-one games in the supplied CSVs")
	}
	out := struct {
		Year     int                     `json:"backtest_year,omitempty"`
		Coverage map[string]*proCoverage `json:"input_coverage"`
		Calendar proPatchCalendar        `json:"competitive_patch_windows"`
		Backtest []proBacktestReport     `json:"backtest,omitempty"`
		Decision *proDecision            `json:"decision,omitempty"`
	}{Coverage: dataset.coverage, Calendar: dataset.calendar}
	if scenario != "" {
		query, err := readJSONInput[proQuery](scenario)
		if err != nil {
			return err
		}
		decision, err := compareProDrafts(dataset, resets, query)
		if err != nil {
			return err
		}
		out.Decision = &decision
	} else {
		out.Year = year
		out.Backtest, err = summarizeProBacktest(backtestPro(dataset, year, resets), year)
		if err != nil {
			return err
		}
	}
	return writeJSONResult(out)
}

func writeJSONResult(result any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func checkAPI(ctx context.Context, id, data string) error {
	key := strings.TrimSpace(os.Getenv("GRID_API_KEY"))
	if key == "" {
		return fmt.Errorf("set GRID_API_KEY in the shell")
	}
	client := grid.NewClient(key, grid.WithTimeout(15*time.Second))
	now := time.Now().UTC().Format(time.RFC3339)
	page, err := grid.AllSeriesWithForm(ctx, client, grid.AllSeriesParams{
		Page:           grid.PageParams{First: grid.Int(5)},
		OrderBy:        grid.SeriesOrderByStartTimeScheduled,
		OrderDirection: grid.OrderDirectionDesc,
		Filter: &grid.SeriesFilter{
			TitleID:            grid.String(grid.TitleIDLeagueOfLegends),
			StartTimeScheduled: &grid.DateTimeFilter{LTE: &now},
		},
	}, grid.NewForm[grid.Series](grid.Scalars("id", "startTimeScheduled")))
	if err != nil {
		return fmt.Errorf("Central Data: %w", err)
	}
	fmt.Printf("Central Data returned %d recent LoL series\n", len(page.Nodes()))
	if id == "latest" {
		if len(page.Nodes()) == 0 {
			return fmt.Errorf("Central Data returned no candidate LoL series")
		}
		id = page.Nodes()[0].ID
	}
	state, err := grid.SeriesStateWithForm(ctx, client, id, draftForm)
	if err != nil {
		return fmt.Errorf("Series State for %s: %w", id, err)
	}
	locked := map[string]bool{}
	fmt.Printf("Series State %s contains %d completed games\n", state.ID, len(state.Games))
	for _, game := range state.Games {
		var picks, bans int
		prior := len(locked)
		if !game.Finished {
			return fmt.Errorf("Game %d is unfinished", game.Number)
		}
		for _, action := range game.Actions {
			switch action.Kind {
			case "pick":
				if action.Champion.ID == "" || locked[action.Champion.ID] {
					return fmt.Errorf("Game %d has a missing or repeated pick", game.Number)
				}
				picks++
				locked[action.Champion.ID] = true
			case "ban":
				bans++
			}
		}
		if picks != 10 {
			return fmt.Errorf("Game %d has %d picks, expected 10", game.Number, picks)
		}
		fmt.Printf("G%d  entering locks=%d recorded picks=%d bans=%d\n", game.Number, prior, picks, bans)
	}
	if data != "" {
		matches, err := loadIndex(data)
		if err != nil {
			return err
		}
		for _, series := range matches {
			if series.ID != state.ID {
				continue
			}
			archived, err := loadArchive(ctx, data, series, nil)
			if err != nil {
				return err
			}
			return compareDrafts(series, state, archived)
		}
		return fmt.Errorf("series %s is absent from the archive index", state.ID)
	}
	return nil
}

func compareDrafts(series matchSeries, api, archive draftSeries) error {
	if len(api.Games) != len(archive.Games) {
		return fmt.Errorf("API and archive game counts differ")
	}
	if _, err := orderedOpenings(series, api); err != nil {
		return fmt.Errorf("API draft does not match the roster index: %w", err)
	}
	if _, err := orderedOpenings(series, archive); err != nil {
		return fmt.Errorf("archive draft does not match the roster index: %w", err)
	}
	total := 0
	teamTotals := 0
	for i, game := range api.Games {
		for j, action := range game.Actions {
			other := archive.Games[i].Actions[j]
			if action.Number != other.Number || action.Kind != other.Kind ||
				action.Drafter.ID != other.Drafter.ID || action.Champion.ID != other.Champion.ID {
				return fmt.Errorf("API and archive differ at Game %d action %d", game.Number, action.Number)
			}
			total++
		}
		archiveTeams := map[string]draftTeam{}
		for _, team := range archive.Games[i].Teams {
			archiveTeams[team.ID] = team
		}
		for _, team := range game.Teams {
			other, ok := archiveTeams[team.ID]
			if !ok {
				return fmt.Errorf("API and archive team IDs differ in Game %d", game.Number)
			}
			fromAPI, err := team.dragons()
			if err != nil {
				return err
			}
			fromArchive, err := other.dragons()
			if err != nil {
				return err
			}
			if fromAPI != fromArchive {
				return fmt.Errorf("API and archive dragon totals differ in Game %d for team %s", game.Number, team.ID)
			}
			teamTotals++
		}
	}
	fmt.Printf("API and archive agree on all %d ordered picks and bans\n", total)
	fmt.Printf("API and archive agree on dragon totals and first-dragon credit for %d team-games\n", teamTotals)
	return nil
}

func runArchive(ctx context.Context, root, replay, event string, year int) error {
	matches, err := loadIndex(root)
	if err != nil {
		return err
	}
	history := histories(matches)
	pools, seriesPools := countPools(matches, history)
	fmt.Printf("Roster index: %d eligible series before the player-history filter\n", len(matches))
	printPools(pools, seriesPools)
	var openings []opening
	var audit []opening
	var performance, replayPerformance []performanceGame
	var auditedState draftSeries
	loaded, rejected := 0, 0
	for _, series := range matches {
		inCohort := strings.HasPrefix(series.Event, event) && series.Start.Year() == year
		isReplay := series.ID == replay
		if !inCohort && !isReplay {
			continue
		}
		timelines := newTimelines(series)
		state, err := loadArchive(ctx, root, series, func(event archiveEvent) error {
			return observePerformance(timelines, event)
		})
		if err != nil {
			return fmt.Errorf("read series %s: %w", series.ID, err)
		}
		rows, err := orderedOpenings(series, state)
		if err != nil {
			if isReplay {
				return fmt.Errorf("verify replay %s: %w", series.ID, err)
			}
			fmt.Fprintf(os.Stderr, "Excluded series %s from ordered draft cohort: %v\n", series.ID, err)
			rejected++
			continue
		}
		if inCohort {
			loaded++
			openings = append(openings, rows...)
		}
		played, err := preparePerformance(series, state, timelines)
		if err != nil {
			return fmt.Errorf("read series %s gameplay metrics: %w", series.ID, err)
		}
		if inCohort {
			performance = append(performance, played...)
		}
		if isReplay {
			audit, auditedState = rows, state
			replayPerformance = played
		}
	}
	if len(openings) == 0 {
		return fmt.Errorf("no verified ordered draft games in cohort")
	}
	redFirst := 0
	for _, row := range openings {
		if row.firstSide == "red" {
			redFirst++
		}
	}
	fmt.Printf("\nOrdered archive: %d verified series, %d games, %d rejected series, %d red map side first picks\n",
		loaded, len(openings), rejected, redFirst)
	results, perSeries := countOpenings(openings)
	printOpenings(results, perSeries)
	printPerformance(countPerformance(performance))
	if replay != "" {
		if len(audit) != 5 {
			return fmt.Errorf("series %s has %d verified games, expected 5", replay, len(audit))
		}
		printReplay(audit, auditedState, history)
		printPerformanceReplay(replayPerformance, performance)
	}
	return nil
}

func printPools(byGame [5]poolResult, bySeries []poolResult) {
	fmt.Println("Player top-3 hit | eligible/total | seen in last 20 | reps | recency | patch shrinkage")
	var all poolResult
	for i, r := range byGame {
		fmt.Printf("G%d  %d/%d  %s  %s  %s  %s\n", i+1, r.eligible, r.total, pct(r.seen, r.eligible),
			pct(r.hits[0], r.eligible), pct(r.hits[1], r.eligible), pct(r.hits[2], r.eligible))
		all.eligible += r.eligible
		all.total += r.total
		all.seen += r.seen
		for j, hits := range r.hits {
			all.hits[j] += hits
		}
	}
	fmt.Printf("All %d/%d  %s  %s  %s  %s\n", all.eligible, all.total, pct(all.seen, all.eligible),
		pct(all.hits[0], all.eligible), pct(all.hits[1], all.eligible), pct(all.hits[2], all.eligible))
	for model, label := range []string{"recency", "patch shrinkage"} {
		var paired []pairedCount
		for _, r := range bySeries {
			paired = append(paired, pairedCount{float64(r.hits[model+1] - r.hits[0]), r.eligible})
		}
		if len(paired) > 0 {
			lo, hi := bootstrapDelta(paired)
			fmt.Printf("%s versus raw reps, paired series bootstrap 95%% interval: %+.1f to %+.1f pp\n", label, 100*lo, 100*hi)
		}
	}
}

func printOpenings(byGame [5]openingResult, bySeries []openingResult) {
	fmt.Println("Opening top-3: at least one / both real reply picks | league | team | team+opener")
	var all openingResult
	for i, r := range byGame {
		printOpeningRow(fmt.Sprintf("G%d", i+1), r)
		all.samples += r.samples
		all.context += r.context
		for j := range all.one {
			all.one[j] += r.one[j]
			all.both[j] += r.both[j]
		}
	}
	printOpeningRow("All", all)
	var paired []pairedCount
	for _, r := range bySeries {
		paired = append(paired, pairedCount{float64(r.both[2] - r.both[1]), r.samples})
	}
	if len(paired) > 0 {
		lo, hi := bootstrapDelta(paired)
		fmt.Printf("Opener versus team for both picks, paired series bootstrap 95%% interval: %+.1f to %+.1f pp\n", 100*lo, 100*hi)
	}
}

func printOpeningRow(label string, r openingResult) {
	fmt.Printf("%s  n=%d, prior opener n>0: %d  %s/%s  %s/%s  %s/%s\n",
		label, r.samples, r.context, pct(r.one[0], r.samples), pct(r.both[0], r.samples),
		pct(r.one[1], r.samples), pct(r.both[1], r.samples), pct(r.one[2], r.samples), pct(r.both[2], r.samples))
}

func printReplay(rows []opening, state draftSeries, h playerHistory) {
	fmt.Printf("\nFive-game replay %s: entering locks, observed first pick → reply pair, practiced options changing\n", state.ID)
	locked := map[string]bool{}
	for i, row := range rows {
		fmt.Printf("G%d  %d locked  %s (%s map side) → %s + %s\n", row.game, len(locked), row.firstName, row.firstSide, row.reply[0].Name, row.reply[1].Name)
		after := map[string]bool{}
		for id := range locked {
			after[id] = true
		}
		for _, action := range state.Games[i].Actions {
			if action.Kind == "pick" {
				after[action.Champion.ID] = true
			}
		}
		if i < len(rows)-1 {
			for _, player := range row.series.Games[i].Players {
				recent := h.recent(player.Player, row.series.Start)
				if len(recent) < 5 {
					continue
				}
				counts := poolScores(recent, "", row.series.Start)[0]
				var lost, open []string
				for _, pick := range recent {
					if locked[pick.champion] || counts[pick.champion] < 2 {
						continue
					}
					delete(counts, pick.champion)
					if after[pick.champion] {
						lost = append(lost, pick.name)
					} else {
						open = append(open, pick.name)
					}
				}
				if len(lost) > 0 && len(open) <= 1 {
					slices.Sort(lost)
					fmt.Printf("    %s: %d → %d stage-practiced options after %s locked",
						player.Name, len(lost)+len(open), len(open), strings.Join(lost, ", "))
					if len(open) == 1 {
						fmt.Printf(" (%s remains)", open[0])
					}
					fmt.Println()
				}
			}
		}
		locked = after
	}
}
