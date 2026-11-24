package main

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"time"
)

const (
	minute15          = 15 * 60
	maxSampleAge      = 15
	dragonObjectiveID = "slayDragon"
)

type goldFrame struct {
	second int
	gold   map[string]int
}

type dragonCapture struct {
	second int
	team   string
}

type gameTimeline struct {
	frames    map[int]goldFrame
	dragons   []dragonCapture
	seen      map[string]bool
	badDragon bool
}

type performanceGame struct {
	series                          matchSeries
	number                          int
	blue, red, blueName, redName    string
	winner                          string
	at15, beforeDragon, afterDragon goldFrame
	firstDragon                     dragonCapture
	dragonVerified                  bool
}

func newTimelines(series matchSeries) map[string]*gameTimeline {
	timelines := make(map[string]*gameTimeline, len(series.Games))
	for _, game := range series.Games {
		timelines[game.ID] = &gameTimeline{frames: map[int]goldFrame{}, seen: map[string]bool{}}
	}
	return timelines
}

func observePerformance(timelines map[string]*gameTimeline, event archiveEvent) error {
	var state struct {
		Games []struct {
			ID       string `json:"id"`
			Finished bool   `json:"finished"`
			Clock    *struct {
				Second int `json:"currentSeconds"`
			} `json:"clock"`
			Teams []struct {
				ID   string `json:"id"`
				Gold *int   `json:"totalMoneyEarned"`
			} `json:"teams"`
		} `json:"games"`
	}
	if err := json.Unmarshal(event.State, &state); err != nil {
		return fmt.Errorf("decode gameplay state: %w", err)
	}
	for _, game := range state.Games {
		timeline := timelines[game.ID]
		if timeline == nil || game.Finished || game.Clock == nil || game.Clock.Second < 1 {
			continue
		}
		second := game.Clock.Second
		if second <= minute15 && len(game.Teams) == 2 {
			first, secondTeam := game.Teams[0], game.Teams[1]
			if first.ID != "" && secondTeam.ID != "" && first.ID != secondTeam.ID && first.Gold != nil && secondTeam.Gold != nil {
				timeline.frames[second] = goldFrame{second, map[string]int{first.ID: *first.Gold, secondTeam.ID: *secondTeam.Gold}}
			}
		}
		if len(state.Games) == 1 {
			timeline.observeDragon(event, second)
		}
	}
	return nil
}

func (t *gameTimeline) observeDragon(event archiveEvent, second int) {
	if event.Kind != "player-completed-slayDragon" && event.Kind != "player-!completed-slayDragon" {
		return
	}
	if event.ID == "" || event.Actor.State.Team == "" {
		t.badDragon = true
		return
	}
	if t.seen[event.ID] {
		return
	}
	t.seen[event.ID] = true
	team := event.Actor.State.Team
	if event.Kind == "player-completed-slayDragon" {
		t.dragons = append(t.dragons, dragonCapture{second, team})
		return
	}
	for i := len(t.dragons) - 1; i >= 0; i-- {
		if t.dragons[i].team == team {
			t.dragons = slices.Delete(t.dragons, i, i+1)
			return
		}
	}
	t.badDragon = true
}

func (t *gameTimeline) frameAt(second int, blue, red string) goldFrame {
	for clock := second; clock >= second-maxSampleAge; clock-- {
		frame := t.frames[clock]
		_, hasBlue := frame.gold[blue]
		_, hasRed := frame.gold[red]
		if hasBlue && hasRed {
			return frame
		}
	}
	return goldFrame{}
}

type dragonSummary struct {
	count int
	first bool
}

func (team draftTeam) dragons() (dragonSummary, error) {
	for _, objective := range team.Objectives {
		if objective.ID == dragonObjectiveID {
			count, err := objective.Count.Int64()
			if err != nil || count < 0 {
				return dragonSummary{}, fmt.Errorf("invalid dragon count for team %s", team.ID)
			}
			return dragonSummary{int(count), objective.First}, nil
		}
	}
	return dragonSummary{}, nil
}

func preparePerformance(series matchSeries, state draftSeries, timelines map[string]*gameTimeline) ([]performanceGame, error) {
	performance := make([]performanceGame, 0, len(series.Games))
	for i, game := range state.Games {
		index := series.Games[i]
		record := performanceGame{series: series, number: game.Number, winner: index.Winner}
		summaries := map[string]dragonSummary{}
		for _, team := range game.Teams {
			stat, err := team.dragons()
			if err != nil {
				return nil, err
			}
			summaries[team.ID] = stat
			switch team.Side {
			case "blue":
				record.blue, record.blueName = team.ID, team.Name
			case "red":
				record.red, record.redName = team.ID, team.Name
			}
		}
		t := timelines[game.ID]
		record.at15 = t.frameAt(minute15, record.blue, record.red)
		counts := map[string]int{}
		for _, capture := range t.dragons {
			counts[capture.team]++
		}
		slices.SortStableFunc(t.dragons, func(a, b dragonCapture) int { return a.second - b.second })
		if len(t.dragons) > 0 {
			record.firstDragon = t.dragons[0]
		}
		blue, red := summaries[record.blue], summaries[record.red]
		record.dragonVerified = !t.badDragon && len(t.dragons) == blue.count+red.count &&
			blue.count == counts[record.blue] && red.count == counts[record.red] &&
			blue.first == (record.firstDragon.team == record.blue) && red.first == (record.firstDragon.team == record.red)
		if record.dragonVerified && record.firstDragon.team != "" && record.firstDragon.second+120 <= minute15 {
			record.beforeDragon = t.frameAt(record.firstDragon.second-1, record.blue, record.red)
			record.afterDragon = t.frameAt(record.firstDragon.second+120, record.blue, record.red)
		}
		performance = append(performance, record)
	}
	return performance, nil
}

func (g performanceGame) blueGold(frame goldFrame) int {
	return frame.gold[g.blue] - frame.gold[g.red]
}

func (g performanceGame) dragonStatus() int {
	if g.firstDragon.second <= minute15 {
		switch g.firstDragon.team {
		case g.blue:
			return 1
		case g.red:
			return -1
		}
	}
	return 0
}

func (g performanceGame) eligible() bool {
	return g.at15.gold != nil && g.dragonVerified && (g.winner == g.red || g.winner == g.blue)
}

func goldBand(gold int) int {
	switch {
	case gold <= -2000:
		return -2
	case gold <= -1000:
		return -1
	case gold >= 2000:
		return 2
	case gold >= 1000:
		return 1
	default:
		return 0
	}
}

type winSample struct {
	n, wins int
}

func (s *winSample) add(won bool) {
	s.n++
	if won {
		s.wins++
	}
}

func (s winSample) estimate(prior, weight float64) float64 {
	return (float64(s.wins) + weight*prior) / (float64(s.n) + weight)
}

func performanceScores(history []performanceGame, test performanceGame) ([3]float64, bool) {
	var league, gold, dragon winSample
	band := goldBand(test.blueGold(test.at15))
	for _, old := range history {
		if !old.eligible() || !old.series.LastGame.Add(historyLag).Before(test.series.Start) ||
			test.series.Start.Sub(old.series.Start) > 120*24*time.Hour {
			continue
		}
		won := old.winner == old.blue
		league.add(won)
		if goldBand(old.blueGold(old.at15)) == band {
			gold.add(won)
			if old.dragonStatus() == test.dragonStatus() {
				dragon.add(won)
			}
		}
	}
	if league.n < 20 {
		return [3]float64{}, false
	}
	base := league.estimate(0.5, 10)
	lead := gold.estimate(base, 8)
	return [3]float64{base, lead, dragon.estimate(lead, 8)}, true
}

type performanceResult struct {
	games, scored, timedGold, verifiedDragon, timedSwing int
	brier                                                [3]float64
	lead                                                 winSample
}

func countPerformance(rows []performanceGame) ([5]performanceResult, [2][]pairedCount) {
	var results [5]performanceResult
	series := map[string]*performanceResult{}
	var order []string
	for _, test := range rows {
		result := &results[test.number-1]
		result.games++
		if test.at15.gold != nil {
			result.timedGold++
			lead := test.blueGold(test.at15)
			if (lead >= 1000 || lead <= -1000) && (test.winner == test.red || test.winner == test.blue) {
				result.lead.add((lead > 0) == (test.winner == test.blue))
			}
		}
		if test.dragonVerified {
			result.verifiedDragon++
		}
		if test.beforeDragon.gold != nil && test.afterDragon.gold != nil {
			result.timedSwing++
		}
		if !test.eligible() {
			continue
		}
		scores, ok := performanceScores(rows, test)
		if !ok {
			continue
		}
		result.scored++
		aggregate := series[test.series.ID]
		if aggregate == nil {
			aggregate = &performanceResult{}
			series[test.series.ID] = aggregate
			order = append(order, test.series.ID)
		}
		aggregate.scored++
		outcome := 0.0
		if test.winner == test.blue {
			outcome = 1
		}
		for i, estimate := range scores {
			loss := math.Pow(estimate-outcome, 2)
			result.brier[i] += loss
			aggregate.brier[i] += loss
		}
	}
	var paired [2][]pairedCount
	for _, id := range order {
		result := series[id]
		for i := range paired {
			paired[i] = append(paired[i], pairedCount{result.brier[i+1] - result.brier[i], result.scored})
		}
	}
	return results, paired
}

func printPerformance(byGame [5]performanceResult, paired [2][]pairedCount) {
	fmt.Println("\nGameplay holdout | GD@15 coverage | verified dragons | 2-minute swing coverage | 1k gold leads won")
	var all performanceResult
	for i, r := range byGame {
		fmt.Printf("G%d  %d/%d  %d/%d  %d/%d  %d/%d\n",
			i+1, r.timedGold, r.games, r.verifiedDragon, r.games, r.timedSwing, r.games, r.lead.wins, r.lead.n)
		all.games += r.games
		all.timedGold += r.timedGold
		all.verifiedDragon += r.verifiedDragon
		all.timedSwing += r.timedSwing
		all.lead.n += r.lead.n
		all.lead.wins += r.lead.wins
		all.scored += r.scored
		for j, loss := range r.brier {
			all.brier[j] += loss
		}
	}
	fmt.Printf("All %d/%d  %d/%d  %d/%d  %d/%d\n",
		all.timedGold, all.games, all.verifiedDragon, all.games, all.timedSwing, all.games, all.lead.wins, all.lead.n)
	if all.scored > 0 {
		n := float64(all.scored)
		fmt.Printf("Brier loss across %d held-out games (lower is better): league blue-side baseline %.3f | GD@15 %.3f | GD@15 + first dragon %.3f\n",
			all.scored, all.brier[0]/n, all.brier[1]/n, all.brier[2]/n)
		for i, label := range []string{"GD@15 versus league baseline", "First dragon added to GD@15"} {
			lo, hi := bootstrapDelta(paired[i])
			fmt.Printf("%s, paired series bootstrap 95%% interval for change in Brier loss: %+.3f to %+.3f\n", label, lo, hi)
		}
	}
}

func printPerformanceReplay(replay, history []performanceGame) {
	first := replay[0]
	view := first.blue
	viewName := first.blueName
	if first.redName < viewName {
		view, viewName = first.red, first.redName
	}
	fmt.Printf("\nGameplay replay from %s's perspective:\n", viewName)
	for _, game := range replay {
		side := 1
		if game.red == view {
			side = -1
		}
		fmt.Printf("G%d", game.number)
		if game.at15.gold != nil {
			fmt.Printf("  gold diff %+d at %s", side*game.blueGold(game.at15), matchClock(game.at15.second))
		}
		if game.dragonVerified && game.firstDragon.team != "" {
			taker := game.blueName
			if game.firstDragon.team == game.red {
				taker = game.redName
			}
			fmt.Printf("  | first dragon %s at %s", taker, matchClock(game.firstDragon.second))
		}
		if game.beforeDragon.gold != nil && game.afterDragon.gold != nil {
			before, after := side*game.blueGold(game.beforeDragon), side*game.blueGold(game.afterDragon)
			fmt.Printf("  | gold diff %+d → %+d (%+d) from %s to %s",
				before, after, after-before, matchClock(game.beforeDragon.second), matchClock(game.afterDragon.second))
		}
		fmt.Println()
	}
	printOpponentHistory(history, first, first.blue, first.blueName)
	printOpponentHistory(history, first, first.red, first.redName)
}

func printOpponentHistory(history []performanceGame, target performanceGame, team, name string) {
	var first, other winSample
	for _, old := range history {
		if !old.eligible() || !old.series.LastGame.Add(historyLag).Before(target.series.Start) ||
			target.series.Start.Sub(old.series.Start) > 120*24*time.Hour ||
			(old.blue != team && old.red != team) {
			continue
		}
		gap := old.blueGold(old.at15)
		if old.red == team {
			gap = -gap
		}
		if old.firstDragon.team == team && old.firstDragon.second <= minute15 {
			first.add(gap >= 1000)
		} else {
			other.add(gap >= 1000)
		}
	}
	fmt.Printf("%s reached a +1k gold lead at 15 in %d/%d prior games when they took first dragon by 15 and %d/%d otherwise\n",
		name, first.wins, first.n, other.wins, other.n)
}

func matchClock(second int) string {
	return fmt.Sprintf("%d:%02d", second/60, second%60)
}
