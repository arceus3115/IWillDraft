package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestPerformanceUsesTimedBilateralGoldAndVerifiedDragon(t *testing.T) {
	series := matchSeries{Games: []matchGame{{ID: "g1", Number: 1, Winner: "blue"}}}
	timelines := newTimelines(series)
	observe := func(id, kind, team string, seconds, blue, red int) {
		t.Helper()
		event := archiveEvent{ID: id, Kind: kind, State: json.RawMessage(fmt.Sprintf(
			`{"games":[{"id":"g1","clock":{"currentSeconds":%d},"teams":[{"id":"blue","totalMoneyEarned":%d},{"id":"red","totalMoneyEarned":%d}]}]}`,
			seconds, blue, red,
		))}
		event.Actor.State.Team = team
		if err := observePerformance(timelines, event); err != nil {
			t.Fatal(err)
		}
	}
	observe("1", "player-completed-slayDragon", "red", 300, 1000, 900)
	observe("2", "player-!completed-slayDragon", "red", 301, 1000, 900)
	observe("3", "state", "", 398, 2020, 2180)
	observe("4", "player-completed-slayDragon", "blue", 400, 2150, 2180)
	observe("4", "player-completed-slayDragon", "blue", 400, 2150, 2180)
	observe("5", "state", "", 516, 2900, 2690)
	observe("6", "state", "", 898, 5000, 4600)
	observe("7", "state", "", 901, 5200, 4600)

	seriesState := draftSeries{Games: []draftGame{{
		ID: "g1", Number: 1,
		Teams: []draftTeam{
			{ID: "blue", Side: "blue", Objectives: []teamObjective{{ID: dragonObjectiveID, Count: "1", First: true}}},
			{ID: "red", Side: "red"},
		},
	}}}
	rows, err := preparePerformance(series, seriesState, timelines)
	if err != nil {
		t.Fatal(err)
	}
	row := rows[0]
	if !row.dragonVerified || row.firstDragon != (dragonCapture{400, "blue"}) {
		t.Fatalf("first dragon is not reconciled: %#v", row.firstDragon)
	}
	if row.at15.second != 898 || row.blueGold(row.at15) != 400 {
		t.Fatalf("at-15 sample should use the last bilateral frame at or before 15:00: %#v", row.at15)
	}
	swing := row.blueGold(row.afterDragon) - row.blueGold(row.beforeDragon)
	if row.beforeDragon.second != 398 || row.afterDragon.second != 516 || swing != 370 {
		t.Fatalf("wrong first-dragon swing: %d at %d–%d", swing, row.beforeDragon.second, row.afterDragon.second)
	}

	seriesState.Games[0].Teams[0].Objectives[0].Count = "2"
	rows, err = preparePerformance(series, seriesState, timelines)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].dragonVerified || rows[0].eligible() || rows[0].afterDragon.gold != nil {
		t.Fatal("disagreement with the final objective count must exclude dragon-dependent metrics")
	}
}

func TestPerformanceHoldoutExcludesCurrentSeriesAndPost15Dragon(t *testing.T) {
	start := time.Date(2026, time.May, 23, 20, 0, 0, 0, time.UTC)
	row := performanceGame{
		series: matchSeries{Start: start, LastGame: start.Add(time.Hour)},
		blue:   "blue", red: "red", winner: "blue", dragonVerified: true,
		at15:        goldFrame{898, map[string]int{"blue": 12000, "red": 10000}},
		firstDragon: dragonCapture{901, "blue"},
	}
	if row.dragonStatus() != 0 {
		t.Fatal("first dragon after 15:00 leaked into the model")
	}
	var history []performanceGame
	for i := range 20 {
		old := row
		old.series.Start = start.Add(-time.Duration(i+1) * 24 * time.Hour)
		old.series.LastGame = old.series.Start.Add(time.Hour)
		if i%2 == 0 {
			old.winner = "red"
		}
		history = append(history, old)
	}
	scores, ok := performanceScores(append(history, row), row)
	if !ok || scores != [3]float64{0.5, 0.5, 0.5} {
		t.Fatalf("held-out series influenced gameplay estimates: %v, eligible %v", scores, ok)
	}
}
