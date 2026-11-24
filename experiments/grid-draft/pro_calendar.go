package main

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

type proPatchWindow struct {
	Patch      proPatch   `json:"patch"`
	Start      time.Time  `json:"start"`
	End        time.Time  `json:"end"`
	ConflictAt *time.Time `json:"first_conflicting_game,omitempty"`
}

type proPatchCalendar struct {
	Scheduled map[string][]proPatchWindow `json:"scheduled,omitempty"`
	Observed  map[string][]proPatchWindow `json:"observed,omitempty"`
}

func readProPatchWindows(path string) (map[string][]proPatchWindow, error) {
	if path == "" {
		return nil, nil
	}
	windows := map[string][]proPatchWindow{}
	err := readCSV(path, func(row map[string]string) error {
		league := row["league"]
		if !slices.Contains(tierOneLeagues, league) {
			return fmt.Errorf("unknown tier-one league %q", league)
		}
		patch, err := parseProPatch(row["patch"])
		if err != nil {
			return err
		}
		start, startErr := time.Parse(time.RFC3339, row["start"])
		end, endErr := time.Parse(time.RFC3339, row["end"])
		if startErr != nil || endErr != nil || !start.Before(end) {
			return fmt.Errorf("expected start and end timestamps with offsets, with start before end")
		}
		windows[league] = append(windows[league], proPatchWindow{Patch: patch, Start: start.UTC(), End: end.UTC()})
		return nil
	}, "league", "patch", "start", "end")
	if err != nil {
		return nil, err
	}
	for _, league := range tierOneLeagues {
		rows := windows[league]
		slices.SortFunc(rows, func(a, b proPatchWindow) int { return a.Start.Compare(b.Start) })
		for i := 1; i < len(rows); i++ {
			if rows[i].Start.Before(rows[i-1].End) {
				return nil, fmt.Errorf("%s: overlapping %s windows at %s", path, league, rows[i].Start.Format(time.RFC3339))
			}
		}
	}
	return windows, nil
}

// Observed windows cover runs of the same tagged patch. Gaps at patch changes
// remain uncovered because the match feed establishes only when games were played.
func buildProPatchCalendar(scheduled map[string][]proPatchWindow, games []proGame, coverage map[string]*proCoverage) proPatchCalendar {
	calendar := proPatchCalendar{Scheduled: map[string][]proPatchWindow{}, Observed: map[string][]proPatchWindow{}}
	for league, windows := range scheduled {
		calendar.Scheduled[league] = slices.Clone(windows)
	}
	slices.SortFunc(games, compareProGames)
	for _, game := range games {
		windows := calendar.Scheduled[game.league]
		if i := proWindowIndex(windows, game.start); i >= 0 && windows[i].Patch != game.patch {
			if windows[i].ConflictAt == nil {
				windows[i].ConflictAt = &game.start
			}
			coverage[game.league].ScheduleConflicts++
		}
		windows = calendar.Observed[game.league]
		end := game.start.Add(time.Second)
		if n := len(windows); n > 0 && windows[n-1].Patch == game.patch {
			windows[n-1].End = end
		} else {
			calendar.Observed[game.league] = append(windows, proPatchWindow{Patch: game.patch, Start: game.start, End: end})
		}
	}
	return calendar
}

func proWindowIndex(windows []proPatchWindow, at time.Time) int {
	i := sort.Search(len(windows), func(i int) bool { return windows[i].End.After(at) })
	if i == len(windows) || at.Before(windows[i].Start) ||
		(i+1 < len(windows) && !at.Before(windows[i+1].Start)) {
		return -1
	}
	return i
}

func (c proPatchCalendar) patchAt(league string, at time.Time) (proPatchWindow, string) {
	if rows := c.Scheduled[league]; len(rows) > 0 {
		if i := proWindowIndex(rows, at); i >= 0 && rows[i].ConflictAt == nil {
			return rows[i], "scheduled"
		}
	}
	if rows := c.Observed[league]; len(rows) > 0 {
		if i := proWindowIndex(rows, at); i >= 0 {
			return rows[i], "observed"
		}
	}
	return proPatchWindow{}, ""
}

func (c proPatchCalendar) competitiveUpdates(target proContext) map[string][]int {
	addLeague := func(updates map[int]bool, league string) {
		for _, source := range [][]proPatchWindow{c.Observed[league], c.Scheduled[league]} {
			for _, window := range source {
				p := window.Patch
				if (window.ConflictAt == nil || !window.ConflictAt.Before(target.cutoff)) && window.Start.Before(target.cutoff) &&
					p.version == target.patch.version && p.update <= target.patch.update {
					updates[p.update] = true
				}
			}
		}
	}
	out := make(map[string][]int, len(tierOneLeagues))
	for _, league := range tierOneLeagues {
		updates := map[int]bool{target.patch.update: true}
		addLeague(updates, target.league)
		addLeague(updates, league)
		for update := range updates {
			out[league] = append(out[league], update)
		}
		slices.Sort(out[league])
	}
	return out
}
