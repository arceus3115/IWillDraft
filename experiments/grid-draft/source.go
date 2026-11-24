package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	grid "github.com/arceus3115/grid-framework"
)

type appearance struct {
	Player, Name, Champion, ChampionName, Team string
}

type matchGame struct {
	ID, Patch, Winner string
	Number, Session   int
	Started           time.Time
	Players           []appearance
}

type matchSeries struct {
	ID, Event, File string
	Start, LastGame time.Time
	Games           []matchGame
}

func readCSV(path string, read func(map[string]string) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	for row := 2; ; row++ {
		values, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		record := make(map[string]string, len(header))
		for i, field := range header {
			record[field] = values[i]
		}
		if err := read(record); err != nil {
			return fmt.Errorf("%s row %d: %w", path, row, err)
		}
	}
}

func loadIndex(root string) ([]matchSeries, error) {
	series := map[string]*matchSeries{}
	err := readCSV(filepath.Join(root, "lol", "series.csv"), func(row map[string]string) error {
		id := row["series_id"]
		series[id] = &matchSeries{ID: id, Event: row["event"], File: row["file"]}
		return nil
	})
	if err != nil {
		return nil, err
	}

	type gameKey struct{ series, session, game string }
	rosters := map[gameKey][]appearance{}
	err = readCSV(filepath.Join(root, "meta", "rosters.csv"), func(row map[string]string) error {
		if row["title"] == "lol" && row["participation_status"] == "active" {
			key := gameKey{row["series_id"], row["session"], row["game_id"]}
			rosters[key] = append(rosters[key], appearance{
				Player: row["player_id"], Name: row["player_name"], Champion: row["character_id"],
				ChampionName: row["character_name"], Team: row["team_id"],
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = readCSV(filepath.Join(root, "meta", "games.csv"), func(row map[string]string) error {
		s := series[row["series_id"]]
		if row["title"] != "lol" || s == nil || row["finished"] != "True" || row["session_valid"] != "True" {
			return nil
		}
		number, err := strconv.Atoi(row["game_number"])
		if err != nil {
			return err
		}
		session, err := strconv.Atoi(row["session"])
		if err != nil {
			return err
		}
		started, err := time.Parse(time.RFC3339Nano, row["started_at"])
		if err != nil {
			return err
		}
		players := rosters[gameKey{s.ID, row["session"], row["game_id"]}]
		s.Games = append(s.Games, matchGame{
			ID: row["game_id"], Patch: row["patch"], Winner: row["winner_team_id"], Number: number,
			Session: session, Started: started, Players: players,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	var matches []matchSeries
	for _, s := range series {
		slices.SortFunc(s.Games, func(a, b matchGame) int { return a.Number - b.Number })
		if validHistory(*s) {
			s.Start = s.Games[0].Started
			s.LastGame = s.Games[len(s.Games)-1].Started
			matches = append(matches, *s)
		}
	}
	slices.SortFunc(matches, func(a, b matchSeries) int {
		if order := a.Start.Compare(b.Start); order != 0 {
			return order
		}
		return strings.Compare(a.ID, b.ID)
	})
	return matches, nil
}

func validHistory(s matchSeries) bool {
	if len(s.Games) == 0 || len(s.Games) > 5 {
		return false
	}
	picked := map[string]bool{}
	for i, game := range s.Games {
		if game.Number != i+1 || game.Patch == "" || game.Session != s.Games[0].Session || len(game.Players) != 10 ||
			(i > 0 && !game.Started.After(s.Games[i-1].Started)) {
			return false
		}
		players := map[string]bool{}
		teams := map[string]int{}
		for _, p := range game.Players {
			if p.Player == "" || p.Champion == "" || p.Team == "" || players[p.Player] || picked[p.Champion] {
				return false
			}
			players[p.Player], picked[p.Champion] = true, true
			teams[p.Team]++
		}
		if len(teams) != 2 {
			return false
		}
		for _, count := range teams {
			if count != 5 {
				return false
			}
		}
	}
	return true
}

type draftAction struct {
	Number   draftSequence `json:"sequenceNumber"`
	Kind     string        `json:"type"`
	Drafter  entity        `json:"drafter"`
	Champion entity        `json:"draftable"`
}

type draftSequence int

func (s *draftSequence) UnmarshalJSON(raw []byte) error {
	var number json.Number
	if err := json.Unmarshal(raw, &number); err != nil {
		return err
	}
	n, err := strconv.Atoi(string(number))
	if err != nil {
		return fmt.Errorf("invalid draft sequence: %w", err)
	}
	*s = draftSequence(n)
	return nil
}

type entity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type draftGame struct {
	ID       string        `json:"id"`
	Number   int           `json:"sequenceNumber"`
	Finished bool          `json:"finished"`
	Actions  []draftAction `json:"draftActions"`
	Teams    []draftTeam   `json:"teams"`
}

type draftTeam struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Side       string          `json:"side"`
	Objectives []teamObjective `json:"objectives"`
}

type teamObjective struct {
	ID    string      `json:"id"`
	Count json.Number `json:"completionCount"`
	First bool        `json:"completedFirst"`
}

type draftSeries struct {
	ID       string      `json:"id"`
	Finished bool        `json:"finished"`
	Valid    bool        `json:"valid"`
	Games    []draftGame `json:"games"`
}

var draftForm = grid.Transform(grid.NewForm[draftSeries](grid.Select(
	grid.Scalar("id"), grid.Scalar("finished"), grid.Scalar("valid"),
	grid.Field("games",
		grid.Scalar("id"), grid.Scalar("sequenceNumber"), grid.Scalar("finished"),
		grid.Field("teams", grid.Scalar("id"), grid.Scalar("name"), grid.Scalar("side"),
			grid.Field("objectives", grid.Scalar("id"), grid.Scalar("completionCount"), grid.Scalar("completedFirst")),
		),
		grid.Field("draftActions",
			grid.Scalar("sequenceNumber"), grid.Scalar("type"),
			grid.Field("drafter", grid.Scalar("id")),
			grid.Field("draftable", grid.Scalar("id"), grid.Scalar("name")),
		),
	),
)), func(series draftSeries) (draftSeries, error) {
	if series.ID == "" || !series.Valid || !series.Finished || len(series.Games) == 0 {
		return draftSeries{}, fmt.Errorf("series state lacks a valid completed series")
	}
	slices.SortFunc(series.Games, func(a, b draftGame) int { return a.Number - b.Number })
	for _, game := range series.Games {
		slices.SortFunc(game.Actions, func(a, b draftAction) int { return int(a.Number - b.Number) })
	}
	return series, nil
})

func loadArchive(ctx context.Context, root string, series matchSeries, observe func(archiveEvent) error) (draftSeries, error) {
	cmd := exec.CommandContext(ctx, "zstd", "-dc", filepath.Join(root, series.File))
	output, err := cmd.StdoutPipe()
	if err != nil {
		return draftSeries{}, err
	}
	if err := cmd.Start(); err != nil {
		return draftSeries{}, err
	}
	final, readErr := scanFinalState(output, series.Games[0].Session, observe)
	if readErr != nil {
		output.Close()
		cmd.Wait()
		return draftSeries{}, readErr
	}
	if err := cmd.Wait(); err != nil {
		return draftSeries{}, err
	}
	draft, err := draftForm.Decode(final)
	if err != nil {
		return draftSeries{}, err
	}
	if draft.ID != series.ID || len(draft.Games) != len(series.Games) {
		return draftSeries{}, fmt.Errorf("archive disagrees with index on series ID or game count")
	}
	return draft, nil
}

type archiveEvent struct {
	ID    string          `json:"id"`
	Kind  string          `json:"type"`
	Full  bool            `json:"includesFullState"`
	State json.RawMessage `json:"seriesState"`
	Actor struct {
		State struct {
			Team string `json:"teamId"`
		} `json:"state"`
	} `json:"actor"`
}

func scanFinalState(source io.Reader, indexSession int, observe func(archiveEvent) error) (json.RawMessage, error) {
	reader := bufio.NewScanner(source)
	reader.Buffer(make([]byte, 128<<10), 32<<20)
	var final json.RawMessage
	session, prior := 1, 0
	for reader.Scan() {
		var envelope struct {
			Sequence int            `json:"sequenceNumber"`
			Events   []archiveEvent `json:"events"`
		}
		if err := json.Unmarshal(reader.Bytes(), &envelope); err != nil {
			return nil, err
		}
		if envelope.Sequence == 1 && prior > 0 {
			session++
		}
		prior = envelope.Sequence
		if session != indexSession {
			continue
		}
		for _, event := range envelope.Events {
			if len(event.State) == 0 || string(event.State) == "null" {
				continue
			}
			if event.Full {
				final = event.State
			}
			if observe != nil {
				if err := observe(event); err != nil {
					return nil, err
				}
			}
		}
	}
	if err := reader.Err(); err != nil {
		return nil, err
	}
	if final == nil {
		return nil, fmt.Errorf("session %d lacks a full series state", indexSession)
	}
	return final, nil
}
