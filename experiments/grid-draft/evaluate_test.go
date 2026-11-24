package main

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestDraftFormSelectsIndexedSessionAndNormalizesSequence(t *testing.T) {
	feed := strings.NewReader(`{"sequenceNumber":1,"events":[{"includesFullState":true,"seriesState":{"id":"match","valid":true,"finished":true,"games":[{"sequenceNumber":1,"draftActions":[{"sequenceNumber":"8"},{"sequenceNumber":7}]}]}}]}
{"sequenceNumber":2,"events":[{"includesFullState":false,"seriesState":{"id":"partial"}}]}
{"sequenceNumber":1,"events":[{"includesFullState":true,"seriesState":{"id":"match","valid":false,"finished":true}}]}`)
	raw, err := scanFinalState(feed, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	state, err := draftForm.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	actions := state.Games[0].Actions
	if actions[0].Number != 7 || actions[1].Number != 8 {
		t.Fatalf("mixed sequence encodings were not ordered: %v", actions)
	}
}

func TestReplayPreservesDraftPositionAndScopesBansToTheirGame(t *testing.T) {
	series := matchSeries{ID: "match"}
	state := draftSeries{ID: series.ID}
	for slot := 1; slot <= 5; slot++ {
		game := draftGame{
			ID: fmt.Sprint(slot), Number: slot, Finished: true,
			Teams: []draftTeam{{ID: "a", Side: "red"}, {ID: "b", Side: "blue"}},
		}
		index := matchGame{ID: game.ID, Number: slot}
		for i := 1; i <= 20; i++ {
			action := draftAction{Number: draftSequence(i), Drafter: entity{ID: "a"}}
			switch i {
			case 7, 8, 9, 10, 11, 12, 17, 18, 19, 20:
				action.Kind = "pick"
				action.Champion.ID = fmt.Sprintf("pick-%d-%d", slot, i)
				if i == 8 || i == 9 || i == 12 || i == 17 || i == 20 {
					action.Drafter.ID = "b"
				}
				index.Players = append(index.Players, appearance{Champion: action.Champion.ID, Team: action.Drafter.ID})
			default:
				action.Kind = "ban"
				action.Champion.ID = fmt.Sprintf("ban-%d-%d", slot, i)
			}
			if slot == 1 && i == 1 {
				action.Champion.ID = "pick-2-7"
			}
			game.Actions = append(game.Actions, action)
		}
		state.Games = append(state.Games, game)
		series.Games = append(series.Games, index)
	}

	rows, err := orderedOpenings(series, state)
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range rows {
		if row.firstSide != "red" || row.respondingTeam != "b" {
			t.Fatalf("Game %d inferred draft position from map side", i+1)
		}
		// The response junction accounts for prior picks, six current bans, and the first pick.
		if len(row.unavailable) != i*10+7 {
			t.Fatalf("Game %d has %d unavailable at the junction", i+1, len(row.unavailable))
		}
	}
	if !rows[0].unavailable["pick-2-7"] || rows[1].first != "pick-2-7" {
		t.Fatal("a Game 1 ban should be available to pick in Game 2")
	}
	if rows[1].unavailable["ban-1-2"] {
		t.Fatal("Game 1 ban carried forward into Game 2")
	}
}

func TestHoldoutUsesEarlierSeriesAndLegalChoices(t *testing.T) {
	start := time.Date(2026, time.May, 23, 20, 0, 0, 0, time.UTC)
	old := start.Add(-24 * time.Hour)
	player := playerHistory{"player": {
		{champion: "later", played: start.Add(time.Hour), last: start.Add(2 * time.Hour)},
		{champion: "ongoing", played: start.Add(-3 * time.Hour), last: start.Add(-time.Hour)},
		{champion: "open", played: old, last: old},
	}}
	recent := player.recent("player", start)
	if len(recent) != 1 || recent[0].champion != "open" {
		t.Fatalf("held-out or overlapping series entered player history: %v", recent)
	}

	test := opening{
		series: matchSeries{Start: start, LastGame: start.Add(2 * time.Hour)},
		first:  "opener", respondingTeam: "team", unavailable: map[string]bool{"blocked": true},
		reply: [2]entity{{ID: "later-a"}, {ID: "later-b"}},
	}
	previous := opening{
		series: matchSeries{Start: old, LastGame: old},
		first:  "opener", respondingTeam: "team", reply: [2]entity{{ID: "open"}, {ID: "blocked"}},
	}
	scores, support, ok := openingScores([]opening{test, previous}, test)
	if !ok || support != 1 {
		t.Fatalf("opening support = %d, eligible = %v", support, ok)
	}
	for _, model := range scores {
		if got := top(model, test.unavailable, 3); !slices.Equal(got, []string{"open"}) {
			t.Fatalf("held-out picks or unavailable candidates entered the ranking: %v", got)
		}
	}
}
