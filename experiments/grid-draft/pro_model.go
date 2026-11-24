package main

import (
	"math"
	"strings"
	"time"
)

const (
	proHistoryWindow = 180 * 24 * time.Hour
	minProTrain      = 30
)

type proModelSpec struct {
	name                            string
	draft, localOnly, samePatchOnly bool
}

var proModelSpecs = [...]proModelSpec{
	{name: "team_player_side"},
	{name: "league_only", draft: true, localOnly: true},
	{name: "same_patch_only", draft: true, samePatchOnly: true},
	{name: "cross_league_cross_patch", draft: true},
}

type proContext struct {
	league string
	patch  proPatch
	cutoff time.Time
}

type proWeightedGame struct {
	game   proGame
	weight float64
}

func proTraining(games []proGame, target proContext, resets []proPatch, spec proModelSpec) []proWeightedGame {
	oldest := target.cutoff.Add(-proHistoryWindow)
	minPatch := 1
	for _, reset := range resets {
		if reset.version == target.patch.version && reset.update <= target.patch.update {
			minPatch = max(minPatch, reset.update)
		}
	}
	var training []proWeightedGame
	mass := map[string]float64{}
	for _, game := range games {
		if !game.available.Before(target.cutoff) || game.start.Before(oldest) ||
			game.patch.version != target.patch.version || game.patch.update < minPatch ||
			game.patch.update > target.patch.update ||
			(spec.localOnly && game.league != target.league) ||
			(spec.samePatchOnly && game.patch != target.patch) {
			continue
		}
		days := target.cutoff.Sub(game.start).Hours() / 24
		patchDistance := float64(target.patch.update - game.patch.update)
		weight := math.Exp2(-days/60 - patchDistance/2)
		training = append(training, proWeightedGame{game, weight})
		mass[game.league] += weight
	}
	// Limit the influence of high-volume source leagues without inflating sparse leagues.
	cap := 0.5 * max(12, mass[target.league])
	for i, row := range training {
		if league := row.game.league; league != target.league && mass[league] > cap {
			training[i].weight *= cap / mass[league]
		}
	}
	return training
}

type proFeatureKind int

const (
	sideEffect proFeatureKind = iota
	leagueSideEffect
	firstPickEffect
	teamEffect
	playerEffect
	championEffect
	leagueChampionEffect
	patchChampionEffect
	playerChampionEffect
	synergyEffect
	matchupEffect
)

var proFeatureNames = [...]string{
	"map_side", "league_map_side", "first_pick", "team", "player",
	"champion_role", "league_champion_role", "patch_champion_role",
	"player_champion", "allied_pair", "lane_matchup",
}

var proRidge = [...]float64{2, 8, 12, 7, 24, 26, 55, 55, 80, 100, 100}

type proFeature struct {
	kind    proFeatureKind
	role    int
	partner int
	a, b    string
}

type proTerm struct {
	feature proFeature
	sign    float64
}

func proTerms(game proGame, draft bool) []proTerm {
	terms := make([]proTerm, 0, 80)
	add := func(kind proFeatureKind, role, partner int, a, b string, sign float64) {
		terms = append(terms, proTerm{proFeature{kind, role, partner, a, b}, sign})
	}
	add(sideEffect, 0, 0, "", "", 1)
	add(leagueSideEffect, 0, 0, game.league, "", 1)
	if game.firstPick != 0 {
		add(firstPickEffect, 0, 0, "", "", float64(game.firstPick))
	}
	for side, lineup := range game.sides {
		sign := 1 - 2*float64(side)
		add(teamEffect, 0, 0, game.league, lineup.Team, sign)
		for i, pick := range lineup.Picks {
			add(playerEffect, 0, 0, game.league, pick.Player, sign)
			if !draft {
				continue
			}
			add(championEffect, i, 0, pick.Champion, "", sign)
			add(leagueChampionEffect, i, 0, pick.Champion, game.league, sign)
			add(patchChampionEffect, i, 0, pick.Champion, game.patch.String(), sign)
			add(playerChampionEffect, i, 0, pick.Champion, pick.Player, sign)
			for j := i + 1; j < len(lineup.Picks); j++ {
				add(synergyEffect, i, j, pick.Champion, lineup.Picks[j].Champion, sign)
			}
		}
	}
	if draft {
		for i, blue := range game.sides[0].Picks {
			red := game.sides[1].Picks[i]
			a, b := blue.Champion, red.Champion
			sign := 1.0
			if a > b {
				a, b, sign = b, a, -1
			}
			add(matchupEffect, i, 0, a, b, sign)
		}
	}
	return terms
}

type proCoefficient struct {
	value  float64
	games  int
	weight float64
}

type proModel struct {
	draft bool
	coef  map[proFeature]proCoefficient
}

func logistic(score float64) float64 {
	return 1 / (1 + math.Exp(-score))
}

func (m proModel) score(game proGame) float64 {
	var score float64
	for _, term := range proTerms(game, m.draft) {
		score += term.sign * m.coef[term.feature].value
	}
	return score
}

func fitProModel(training []proWeightedGame, draft bool) proModel {
	type observation struct {
		row  int
		sign float64
	}
	type column struct {
		feature proFeature
		rows    []observation
		weight  float64
		value   float64
	}
	indices := map[proFeature]int{}
	var columns []column
	for i, row := range training {
		for _, term := range proTerms(row.game, draft) {
			j, ok := indices[term.feature]
			if !ok {
				j = len(columns)
				indices[term.feature] = j
				columns = append(columns, column{feature: term.feature})
			}
			columns[j].rows = append(columns[j].rows, observation{i, term.sign})
			columns[j].weight += row.weight
		}
	}
	logits := make([]float64, len(training))
	for pass := 0; pass < 60; pass++ {
		maxChange := 0.0
		for i := range columns {
			col := &columns[i]
			ridge := proRidge[col.feature.kind]
			gradient := ridge * col.value
			for _, obs := range col.rows {
				row := training[obs.row]
				gradient += row.weight * obs.sign * (logistic(logits[obs.row]) - row.game.blueWon)
			}
			// .25 is the upper bound on the logistic loss curvature for a signed unit feature.
			change := -gradient / (ridge + 0.25*col.weight)
			col.value += change
			maxChange = max(maxChange, math.Abs(change))
			for _, obs := range col.rows {
				logits[obs.row] += obs.sign * change
			}
		}
		if maxChange < 1e-5 {
			break
		}
	}
	model := proModel{draft: draft, coef: make(map[proFeature]proCoefficient, len(columns))}
	for _, col := range columns {
		model.coef[col.feature] = proCoefficient{col.value, len(col.rows), col.weight}
	}
	return model
}

func parsePatchResets(value string) ([]proPatch, error) {
	var resets []proPatch
	if value == "" {
		return resets, nil
	}
	for _, raw := range strings.Split(value, ",") {
		patch, err := parseProPatch(raw)
		if err != nil {
			return nil, err
		}
		resets = append(resets, patch)
	}
	return resets, nil
}
