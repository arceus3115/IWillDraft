package main

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCompetitivePatchDistanceCountsPlayedLeagueUpdatesOnly(t *testing.T) {
	cutoff := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)
	start := cutoff.Add(-14 * 24 * time.Hour)
	row := func(league string, update int, date time.Time) proGame {
		return proGame{league: league, patch: proPatch{16, update}, start: date, available: date.Add(historyLag)}
	}
	games := []proGame{
		row("LCK", 3, start), row("LPL", 3, start), row("LPL", 4, start.AddDate(0, 0, 1)),
		row("LEC", 5, start.AddDate(0, 0, 1)), row("LPL", 5, cutoff.AddDate(0, 0, 1)),
	}
	calendar := buildProPatchCalendar(nil, games, nil)
	target := proContext{"LCK", proPatch{16, 6}, cutoff}
	training := proTraining(proDataset{games: games, calendar: calendar}, target, nil, proModelSpecs[3])
	oldWeight := math.Exp2(-14.0 / 60)
	want := map[string]float64{"LCK": oldWeight * math.Exp2(-0.5), "LPL": oldWeight * math.Exp2(-1)}
	for _, row := range training {
		if row.game.patch.update != 3 {
			continue
		}
		if diff := math.Abs(row.weight - want[row.game.league]); diff > 1e-12 {
			t.Errorf("%s patch weight = %.6f, want %.6f", row.game.league, row.weight, want[row.game.league])
		}
		delete(want, row.game.league)
	}
	if len(want) > 0 {
		t.Fatalf("missing training games: %v", want)
	}
}

func TestPatchWindowsUseLeagueAndLeaveUnverifiedTransitionGaps(t *testing.T) {
	date := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	games := []proGame{
		{league: "LCK", patch: proPatch{16, 6}, start: date},
		{league: "LCK", patch: proPatch{16, 6}, start: date.AddDate(0, 0, 3)},
		{league: "LCK", patch: proPatch{16, 7}, start: date.AddDate(0, 0, 6)},
		{league: "LPL", patch: proPatch{16, 7}, start: date},
		{league: "LPL", patch: proPatch{16, 7}, start: date.AddDate(0, 0, 3)},
	}
	calendar := buildProPatchCalendar(nil, games, nil)
	if lck, source := calendar.patchAt("LCK", date.AddDate(0, 0, 2)); lck.Patch.update != 6 || source != "observed" {
		t.Fatalf("LCK patch = %v, source = %s", lck, source)
	}
	if lpl, source := calendar.patchAt("LPL", date.AddDate(0, 0, 2)); lpl.Patch.update != 7 || source != "observed" {
		t.Fatalf("LPL patch = %v, source = %s", lpl, source)
	}
	if _, source := calendar.patchAt("LCK", date.AddDate(0, 0, 4)); source != "" {
		t.Fatalf("assigned a patch in an unverified transition: %q", source)
	}
}

func TestScheduledPatchWindowCSVUsesHalfOpenRanges(t *testing.T) {
	path := filepath.Join(t.TempDir(), "calendar.csv")
	data := "league,patch,start,end\n" +
		"LCK,26.6,2026-04-01T17:00:00+09:00,2026-04-08T17:00:00+09:00\n" +
		"LCK,26.7,2026-04-08T17:00:00+09:00,2026-04-15T17:00:00+09:00\n"
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	rows, err := readProPatchWindows(path)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 4, 8, 8, 0, 0, 0, time.UTC)
	calendar := proPatchCalendar{Scheduled: rows}
	if window, source := calendar.patchAt("LCK", at); window.Patch != (proPatch{16, 7}) || source != "scheduled" {
		t.Fatalf("patch at the boundary = %v, source = %s", window, source)
	}
	data += "LCK,26.8,2026-04-14T17:00:00+09:00,2026-04-25T17:00:00+09:00\n"
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readProPatchWindows(path); err == nil {
		t.Fatal("expected an error for overlapping league windows")
	}
}
