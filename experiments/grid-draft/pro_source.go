package main

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var tierOneLeagues = []string{"LCS", "LCK", "LPL", "LEC"}
var proRoles = []string{"top", "jng", "mid", "bot", "sup"}
var proSideNames = []string{"blue", "red"}

const proSides = 2

type proPatch struct{ version, update int }

func parseProPatch(raw string) (proPatch, error) {
	version, update, found := strings.Cut(strings.TrimSpace(raw), ".")
	major, majorErr := strconv.Atoi(version)
	minor, minorErr := strconv.Atoi(update)
	if !found || majorErr != nil || minorErr != nil || major <= 0 || minor <= 0 {
		return proPatch{}, fmt.Errorf("invalid patch %q, expected two positive version numbers", raw)
	}
	// Competitive exports use the client major (16.x in 2026). Riot notes use 26.x.
	if major >= 25 && major <= 29 {
		major -= 10
	}
	return proPatch{major, minor}, nil
}

func (p proPatch) String() string { return fmt.Sprintf("%d.%d", p.version, p.update) }

func championKey(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, name)
}

type proPick struct {
	Player   string `json:"player_id"`
	Champion string `json:"champion"`
}

type proLineup struct {
	Team  string
	Picks [5]proPick
}

type proGame struct {
	id, league, split               string
	patch                           proPatch
	start, startOfSeries, available time.Time
	series                          string
	number, year                    int
	sides                           [proSides]proLineup
	bans                            map[string]bool
	firstPick                       int // +1 blue, -1 red, 0 unverified
	blueWon                         float64
}

type oracleGame struct {
	rows [proSides][6]map[string]string
	err  string
}

type proCoverage struct {
	Eligible       int            `json:"eligible_games"`
	Rejected       map[string]int `json:"rejected_games_by_reason"`
	Patches        map[string]int `json:"eligible_games_by_patch"`
	First          *time.Time     `json:"first_game,omitempty"`
	Last           *time.Time     `json:"last_game,omitempty"`
	FirstPickKnown int            `json:"first_pick_known_games"`
}

var oracleColumns = []string{"gameid", "league", "year", "split", "date", "game", "patch", "side", "position",
	"teamid", "playerid", "champion", "result", "ban1", "ban2", "ban3", "ban4", "ban5"}

func importProGames(paths []string) ([]proGame, map[string]*proCoverage, error) {
	coverage := make(map[string]*proCoverage)
	for _, league := range tierOneLeagues {
		coverage[league] = &proCoverage{Rejected: map[string]int{}, Patches: map[string]int{}}
	}
	index := map[string]*oracleGame{}
	for _, path := range paths {
		err := readCSV(path, func(row map[string]string) error {
			league := row["league"]
			if coverage[league] == nil {
				return nil
			}
			if row["gameid"] == "" {
				return fmt.Errorf("missing game id for %s", league)
			}
			key := league + "/" + row["gameid"]
			group := index[key]
			if group == nil {
				group = &oracleGame{}
				index[key] = group
			}
			side := slices.Index(proSideNames, strings.ToLower(row["side"]))
			position := slices.Index(proRoles, row["position"])
			if row["position"] == "team" {
				position = len(proRoles)
			}
			if side < 0 || position < 0 {
				group.err = "invalid side or role"
			} else if group.rows[side][position] != nil && !sameOracleRow(group.rows[side][position], row) {
				group.err = "conflicting duplicate rows"
			} else {
				keep := make(map[string]string, len(oracleColumns)+1)
				for _, field := range oracleColumns {
					keep[field] = row[field]
				}
				keep["firstPick"] = row["firstPick"]
				group.rows[side][position] = keep
			}
			return nil
		}, oracleColumns...)
		if err != nil {
			return nil, nil, err
		}
	}
	type seriesKey struct {
		league, split, teamA, teamB, nativeID string
		year                                  int
	}
	type entry struct {
		game   proGame
		reason string
	}
	pairs := map[seriesKey][]entry{}
	for key, group := range index {
		game, reason := group.game()
		if game.start.IsZero() || game.number == 0 || game.sides[0].Team == "" || game.sides[1].Team == "" {
			league, _, _ := strings.Cut(key, "/")
			coverage[league].Rejected[reason]++
			continue
		}
		blue, red := game.sides[0].Team, game.sides[1].Team
		var nativeID string
		if prefix, number, ok := strings.Cut(game.id, "_game_"); ok && number == strconv.Itoa(game.number) {
			nativeID = prefix
		}
		pair := seriesKey{game.league, game.split, min(blue, red), max(blue, red), nativeID, game.year}
		pairs[pair] = append(pairs[pair], entry{game, reason})
	}

	var games []proGame
	for pair, history := range pairs {
		slices.SortFunc(history, func(a, b entry) int {
			return cmp.Or(a.game.start.Compare(b.game.start), cmp.Compare(a.game.number, b.game.number), cmp.Compare(a.game.id, b.game.id))
		})
		for len(history) > 0 {
			end := 1
			for end < len(history) && (pair.nativeID != "" || history[end].game.number != 1) {
				end++
			}
			series := history[:end]
			history = history[end:]
			first, last := series[0].game, series[len(series)-1].game
			complete := len(series) <= 5 && last.start.Sub(first.start) <= 8*time.Hour
			for i, row := range series {
				game := row.game
				if game.number != i+1 || (i > 0 && (!game.start.After(series[i-1].game.start) ||
					game.start.Sub(series[i-1].game.start) > historyLag)) {
					complete = false
				}
			}
			for _, row := range series {
				game := row.game
				report := coverage[game.league]
				reason := row.reason
				if !complete && reason == "" {
					reason = "unverified series boundary"
				}
				if reason != "" {
					report.Rejected[reason]++
					continue
				}
				game.series = first.id
				game.startOfSeries = first.start
				game.available = last.start.Add(historyLag)
				games = append(games, game)
				report.Eligible++
				if game.firstPick != 0 {
					report.FirstPickKnown++
				}
				report.Patches[game.patch.String()]++
				if report.First == nil || game.start.Before(*report.First) {
					report.First = &game.start
				}
				if report.Last == nil || game.start.After(*report.Last) {
					report.Last = &game.start
				}
			}
		}
	}
	slices.SortFunc(games, compareProGames)
	return games, coverage, nil
}

func sameOracleRow(a, b map[string]string) bool {
	for _, field := range []string{"league", "year", "split", "date", "game", "patch", "teamid",
		"playerid", "champion", "result", "ban1", "ban2", "ban3", "ban4", "ban5", "firstPick"} {
		if a[field] != b[field] {
			return false
		}
	}
	return true
}

func compareProGames(a, b proGame) int {
	return cmp.Or(a.start.Compare(b.start), cmp.Compare(a.league, b.league), cmp.Compare(a.id, b.id))
}

func (group *oracleGame) game() (proGame, string) {
	blue, red := group.rows[0][len(proRoles)], group.rows[1][len(proRoles)]
	if blue == nil || red == nil {
		return proGame{}, "missing team rows"
	}
	g := proGame{id: blue["gameid"], league: blue["league"], split: blue["split"], bans: map[string]bool{}}
	g.sides[0].Team, g.sides[1].Team = blue["teamid"], red["teamid"]
	if g.sides[0].Team == "" || g.sides[1].Team == "" || g.sides[0].Team == g.sides[1].Team {
		g.sides[0].Team, g.sides[1].Team = "", ""
		return g, "missing or duplicate team ids"
	}
	var err error
	g.start, err = time.Parse("2006-01-02 15:04:05", blue["date"])
	if err != nil {
		return g, "invalid timestamp"
	}
	g.year, err = strconv.Atoi(blue["year"])
	if err != nil || g.year < 2025 {
		return g, "outside supported seasons (2025 onward)"
	}
	g.number, err = strconv.Atoi(blue["game"])
	if err != nil || g.number < 1 {
		g.number = 0
		return g, "invalid game number"
	}
	if group.err != "" {
		return g, group.err
	}
	g.patch, err = parseProPatch(blue["patch"])
	if err != nil {
		return g, "missing or invalid patch"
	}
	if blue["firstPick"] == "1" && red["firstPick"] == "0" {
		g.firstPick = 1
	}
	if blue["firstPick"] == "0" && red["firstPick"] == "1" {
		g.firstPick = -1
	}
	if blue["result"] == "1" && red["result"] == "0" {
		g.blueWon = 1
	} else if blue["result"] != "0" || red["result"] != "1" {
		return g, "invalid game outcome"
	}
	picked, players := map[string]bool{}, map[string]bool{}
	for side, rows := range group.rows {
		team := rows[len(proRoles)]
		for _, row := range rows {
			if row == nil {
				return g, "incomplete role assignments"
			}
			if row["teamid"] != team["teamid"] || row["result"] != team["result"] ||
				row["date"] != blue["date"] || row["game"] != blue["game"] ||
				row["year"] != blue["year"] || row["split"] != blue["split"] || row["patch"] != blue["patch"] {
				return g, "inconsistent game metadata"
			}
		}
		for role, row := range rows[:len(proRoles)] {
			pick := proPick{row["playerid"], championKey(row["champion"])}
			if pick.Player == "" || pick.Champion == "" || players[pick.Player] || picked[pick.Champion] {
				return g, "missing or duplicate champion or player"
			}
			g.sides[side].Picks[role] = pick
			players[pick.Player], picked[pick.Champion] = true, true
		}
		for _, field := range []string{"ban1", "ban2", "ban3", "ban4", "ban5"} {
			if ban := championKey(team[field]); ban != "" {
				g.bans[ban] = true
			}
		}
	}
	for ban := range g.bans {
		if picked[ban] {
			return g, "picked champion also recorded as banned"
		}
	}
	return g, ""
}
