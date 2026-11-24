package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

const historyLag = 4 * time.Hour

type scoreMap map[string]float64

func top(scores scoreMap, excluded map[string]bool, limit int) []string {
	choices := make([]string, 0, len(scores))
	for id, score := range scores {
		if !excluded[id] && score > 0 {
			choices = append(choices, id)
		}
	}
	slices.SortFunc(choices, func(a, b string) int {
		if scores[a] > scores[b] {
			return -1
		}
		if scores[a] < scores[b] {
			return 1
		}
		return strings.Compare(a, b)
	})
	return choices[:min(limit, len(choices))]
}

type pastPick struct {
	champion, name, patch string
	played, last          time.Time
}

type playerHistory map[string][]pastPick

func histories(matches []matchSeries) playerHistory {
	out := playerHistory{}
	for _, series := range matches {
		for _, game := range series.Games {
			for _, p := range game.Players {
				out[p.Player] = append(out[p.Player], pastPick{p.Champion, p.ChampionName, game.Patch, game.Started, series.LastGame})
			}
		}
	}
	for _, games := range out {
		slices.SortFunc(games, func(a, b pastPick) int { return b.played.Compare(a.played) })
	}
	return out
}

func (h playerHistory) recent(player string, cutoff time.Time) []pastPick {
	var out []pastPick
	for _, pick := range h[player] {
		if !pick.last.Add(historyLag).Before(cutoff) {
			continue
		}
		out = append(out, pick)
		if len(out) == 20 {
			break
		}
	}
	return out
}

func poolScores(picks []pastPick, patch string, cutoff time.Time) [3]scoreMap {
	counts, recency, patchCounts := scoreMap{}, scoreMap{}, scoreMap{}
	var patchGames float64
	for _, pick := range picks {
		counts[pick.champion]++
		days := cutoff.Sub(pick.played).Hours() / 24
		recency[pick.champion] += math.Exp2(-days / 30)
		if pick.patch == patch {
			patchCounts[pick.champion]++
			patchGames++
		}
	}
	shrunk := scoreMap{}
	for champion, count := range counts {
		shrunk[champion] = (patchCounts[champion] + 8*count/float64(len(picks))) / (patchGames + 8)
	}
	return [3]scoreMap{counts, recency, shrunk}
}

type poolResult struct {
	total, eligible, seen int
	hits                  [3]int
}

func countPools(matches []matchSeries, h playerHistory) ([5]poolResult, []poolResult) {
	var byGame [5]poolResult
	var bySeries []poolResult
	for _, s := range matches {
		var series poolResult
		locked := map[string]bool{}
		for _, game := range s.Games {
			slot := &byGame[game.Number-1]
			for _, p := range game.Players {
				slot.total++
				series.total++
				recent := h.recent(p.Player, s.Start)
				if len(recent) < 5 {
					continue
				}
				slot.eligible++
				series.eligible++
				scores := poolScores(recent, game.Patch, s.Start)
				if scores[0][p.Champion] > 0 {
					slot.seen++
					series.seen++
				}
				for i, score := range scores {
					if slices.Contains(top(score, locked, 3), p.Champion) {
						slot.hits[i]++
						series.hits[i]++
					}
				}
			}
			for _, p := range game.Players {
				locked[p.Champion] = true
			}
		}
		if series.eligible > 0 {
			bySeries = append(bySeries, series)
		}
	}
	return byGame, bySeries
}

type opening struct {
	series                      matchSeries
	game                        int
	respondingTeam              string
	first, firstName, firstSide string
	reply                       [2]entity
	unavailable                 map[string]bool
}

func orderedOpenings(series matchSeries, state draftSeries) ([]opening, error) {
	locked := map[string]bool{}
	var out []opening
	for i, game := range state.Games {
		indexGame := series.Games[i]
		if game.ID != indexGame.ID || game.Number != indexGame.Number || !game.Finished {
			return nil, fmt.Errorf("game %d differs from the index", indexGame.Number)
		}
		blue, red := "", ""
		for _, team := range game.Teams {
			switch team.Side {
			case "blue":
				blue = team.ID
			case "red":
				red = team.ID
			}
		}
		if blue == "" || red == "" || blue == red || len(game.Actions) != 20 {
			return nil, fmt.Errorf("game %d lacks sides or the 20 Phase 1 and 2 actions", game.Number)
		}
		current := opening{series: series, game: game.Number, unavailable: map[string]bool{}}
		for id := range locked {
			current.unavailable[id] = true
		}
		indexPicks := map[string]string{}
		for _, p := range indexGame.Players {
			indexPicks[p.Champion] = p.Team
		}
		picks := map[string]bool{}
		for j, action := range game.Actions {
			champion, team := action.Champion.ID, action.Drafter.ID
			if int(action.Number) != j+1 || (team != red && team != blue) {
				return nil, fmt.Errorf("game %d action %d has an invalid sequence or side", game.Number, j+1)
			}
			if action.Kind == "ban" {
				if j < 6 && champion != "" {
					current.unavailable[champion] = true
				}
				continue
			}
			if action.Kind != "pick" || indexPicks[champion] != team || locked[champion] || picks[champion] {
				return nil, fmt.Errorf("game %d action %d disagrees with completed picks", game.Number, j+1)
			}
			picks[champion] = true
		}
		for _, action := range game.Actions[:6] {
			if action.Kind != "ban" {
				return nil, fmt.Errorf("game %d lacks six opening bans", game.Number)
			}
		}
		first, reply1, reply2 := game.Actions[6], game.Actions[7], game.Actions[8]
		if len(picks) != 10 || first.Kind != "pick" || reply1.Kind != "pick" || reply2.Kind != "pick" ||
			reply1.Drafter.ID != reply2.Drafter.ID || first.Drafter.ID == reply1.Drafter.ID {
			return nil, fmt.Errorf("game %d lacks a valid first-pick → reply-pair junction", game.Number)
		}
		current.first, current.firstName = first.Champion.ID, first.Champion.Name
		current.respondingTeam = reply1.Drafter.ID
		current.reply = [2]entity{reply1.Champion, reply2.Champion}
		current.firstSide = "blue"
		if first.Drafter.ID == red {
			current.firstSide = "red"
		}
		current.unavailable[current.first] = true
		if current.unavailable[reply1.Champion.ID] || current.unavailable[reply2.Champion.ID] {
			return nil, fmt.Errorf("game %d response conflicts with a prior lock, ban or pick", game.Number)
		}
		out = append(out, current)
		for champion := range picks {
			locked[champion] = true
		}
	}
	return out, nil
}

type openingResult struct {
	samples, context int
	one, both        [3]int
}

func openingScores(train []opening, test opening) ([3]scoreMap, int, bool) {
	league, team, opener := scoreMap{}, scoreMap{}, scoreMap{}
	var leagueN, teamN, openerN int
	for _, old := range train {
		if !old.series.LastGame.Add(historyLag).Before(test.series.Start) ||
			test.series.Start.Sub(old.series.Start) > 120*24*time.Hour {
			continue
		}
		leagueN++
		for _, pick := range old.reply {
			league[pick.ID]++
		}
		if old.respondingTeam == test.respondingTeam {
			teamN++
			for _, pick := range old.reply {
				team[pick.ID]++
			}
		}
		if old.first == test.first {
			openerN++
			for _, pick := range old.reply {
				opener[pick.ID]++
			}
		}
	}
	if teamN == 0 || leagueN == 0 {
		return [3]scoreMap{}, openerN, false
	}
	local, combined := scoreMap{}, scoreMap{}
	for id, count := range league {
		local[id] = (team[id] + 8*count/float64(leagueN)) / float64(teamN+8)
		combined[id] = (opener[id] + 12*local[id]) / float64(openerN+12)
	}
	return [3]scoreMap{league, local, combined}, openerN, true
}

func countOpenings(rows []opening) ([5]openingResult, []openingResult) {
	var results [5]openingResult
	bySeries := map[string]*openingResult{}
	var order []string
	for _, test := range rows {
		scores, openerN, ok := openingScores(rows, test)
		if !ok {
			continue
		}
		series := bySeries[test.series.ID]
		if series == nil {
			series = &openingResult{}
			bySeries[test.series.ID] = series
			order = append(order, test.series.ID)
		}
		result := &results[test.game-1]
		for _, r := range []*openingResult{result, series} {
			r.samples++
			if openerN > 0 {
				r.context++
			}
		}
		for i, score := range scores {
			ranked := top(score, test.unavailable, 3)
			a := slices.Contains(ranked, test.reply[0].ID)
			b := slices.Contains(ranked, test.reply[1].ID)
			for _, r := range []*openingResult{result, series} {
				if a || b {
					r.one[i]++
				}
				if a && b {
					r.both[i]++
				}
			}
		}
	}
	perSeries := make([]openingResult, 0, len(order))
	for _, id := range order {
		perSeries = append(perSeries, *bySeries[id])
	}
	return results, perSeries
}

type pairedCount struct {
	difference float64
	eligible   int
}

func bootstrapDelta(series []pairedCount) (float64, float64) {
	rng := rand.New(rand.NewPCG(1, 2))
	const draws = 4000
	samples := make([]float64, draws)
	for i := range samples {
		difference, eligible := 0.0, 0
		for range series {
			row := series[rng.IntN(len(series))]
			difference += row.difference
			eligible += row.eligible
		}
		samples[i] = difference / float64(eligible)
	}
	slices.Sort(samples)
	return samples[draws*25/1000], samples[draws*975/1000]
}

func pct(n, d int) string {
	if d == 0 {
		return "—"
	}
	return fmt.Sprintf("%.1f%%", 100*float64(n)/float64(d))
}
