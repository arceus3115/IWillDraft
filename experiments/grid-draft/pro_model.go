package main

import (
	"math"
	"slices"
	"sort"
	"strings"
	"time"
)

const minProTrain = 30

type proModelSpec struct {
	name                                        string
	draft, localOnly, samePatchOnly, patchDrift bool
}

var proModelSpecs = [...]proModelSpec{
	{name: "team_player_side"},
	{name: "league_only", draft: true, localOnly: true},
	{name: "same_patch_only", draft: true, samePatchOnly: true},
	{name: "independent_patch_effects", draft: true},
	{name: "cross_league_cross_patch", draft: true, patchDrift: true},
}

const fullProModel = len(proModelSpecs) - 1

type proContext struct {
	league string
	patch  proPatch
	cutoff time.Time
}

type proWeightedGame struct {
	game   proGame
	weight float64
}

func proTraining(dataset proDataset, target proContext, resets []proPatch, spec proModelSpec) []proWeightedGame {
	minPatch := 1
	for _, reset := range resets {
		if reset.version == target.patch.version && reset.update <= target.patch.update {
			minPatch = max(minPatch, reset.update)
		}
	}
	var training []proWeightedGame
	mass := map[string]float64{}
	updatesByLeague := dataset.calendar.competitiveUpdates(target)
	for _, game := range dataset.games {
		if !game.available.Before(target.cutoff) ||
			game.patch.version != target.patch.version || game.patch.update < minPatch ||
			game.patch.update > target.patch.update ||
			(spec.localOnly && game.league != target.league) ||
			(spec.samePatchOnly && game.patch != target.patch) {
			continue
		}
		updates := updatesByLeague[game.league]
		firstNewer := sort.Search(len(updates), func(i int) bool { return updates[i] > game.patch.update })
		patchDistance := float64(len(updates) - firstNewer)
		weight := math.Exp2(-patchDistance / 2)
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
	patch   proPatch
}

type proTerm struct {
	feature proFeature
	sign    float64
}

func proTerms(game proGame, draft bool) []proTerm {
	terms := make([]proTerm, 0, 80)
	add := func(kind proFeatureKind, role, partner int, a, b string, sign float64) {
		terms = append(terms, proTerm{proFeature{kind: kind, role: role, partner: partner, a: a, b: b}, sign})
	}
	add(sideEffect, 0, 0, "", "", 1)
	add(leagueSideEffect, 0, 0, game.league, "", 1)
	if game.firstPick != 0 {
		add(firstPickEffect, 0, 0, "", "", float64(game.firstPick))
	}
	for side, lineup := range game.sides {
		sign := 1 - 2*float64(side)
		add(teamEffect, 0, 0, "", lineup.Team, sign)
		for i, pick := range lineup.Picks {
			add(playerEffect, 0, 0, "", pick.Player, sign)
			if !draft {
				continue
			}
			add(championEffect, i, 0, pick.Champion, "", sign)
			add(leagueChampionEffect, i, 0, pick.Champion, game.league, sign)
			terms = append(terms, proTerm{proFeature{kind: patchChampionEffect, role: i, a: pick.Champion, patch: game.patch}, sign})
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

func fitProSpec(training []proWeightedGame, spec proModelSpec, target proPatch) proModel {
	var axis []proPatch
	if spec.patchDrift {
		axis = proPatchAxis(training, target)
	}
	return fitProModel(training, spec.draft, axis)
}

func proPatchAxis(training []proWeightedGame, target proPatch) []proPatch {
	seen := map[proPatch]bool{target: true}
	for _, row := range training {
		seen[row.game.patch] = true
	}
	axis := make([]proPatch, 0, len(seen))
	for patch := range seen {
		axis = append(axis, patch)
	}
	slices.SortFunc(axis, func(a, b proPatch) int { return a.update - b.update })
	return axis
}

func fitProModel(training []proWeightedGame, draft bool, axis []proPatch) proModel {
	type observation struct {
		row  int
		sign float64
	}
	type column struct {
		feature                     proFeature
		rows                        []observation
		weight                      float64
		value                       float64
		ridge                       float64
		before, after               int
		beforePenalty, afterPenalty float64
	}
	indices := map[proFeature]int{}
	var columns []column
	for i, row := range training {
		for _, term := range proTerms(row.game, draft) {
			j, ok := indices[term.feature]
			if !ok {
				j = len(columns)
				indices[term.feature] = j
				columns = append(columns, column{feature: term.feature, ridge: proRidge[term.feature.kind]})
			}
			columns[j].rows = append(columns[j].rows, observation{i, term.sign})
			columns[j].weight += row.weight
		}
	}
	if len(axis) > 1 {
		rank := make(map[proPatch]int, len(axis))
		for i, patch := range axis {
			rank[patch] = i
		}
		target := axis[len(axis)-1]
		chains := map[proFeature][]int{}
		var championRoles []proFeature
		for i, col := range columns {
			if col.feature.kind == patchChampionEffect {
				key := col.feature
				key.patch = proPatch{}
				if len(chains[key]) == 0 {
					championRoles = append(championRoles, key)
				}
				chains[key] = append(chains[key], i)
			}
		}
		for _, key := range championRoles {
			chain := chains[key]
			slices.SortFunc(chain, func(i, j int) int {
				return rank[columns[i].feature.patch] - rank[columns[j].feature.patch]
			})
			if columns[chain[len(chain)-1]].feature.patch != target {
				key.patch = target
				chain = append(chain, len(columns))
				columns = append(columns, column{feature: key})
			}
			for k := 1; k < len(chain); k++ {
				before, after := chain[k-1], chain[k]
				steps := rank[columns[after].feature.patch] - rank[columns[before].feature.patch]
				// Independent increments add variance over unobserved patch steps.
				precision := 80 / float64(steps)
				columns[before].after, columns[before].afterPenalty = after, precision
				columns[after].before, columns[after].beforePenalty, columns[after].ridge = before, precision, 0
			}
		}
	}
	logits := make([]float64, len(training))
	for pass := 0; pass < 60; pass++ {
		maxChange := 0.0
		for i := range columns {
			col := &columns[i]
			curvature := col.ridge + col.beforePenalty + col.afterPenalty
			gradient := curvature * col.value
			if col.beforePenalty > 0 {
				gradient -= col.beforePenalty * columns[col.before].value
			}
			if col.afterPenalty > 0 {
				gradient -= col.afterPenalty * columns[col.after].value
			}
			for _, obs := range col.rows {
				row := training[obs.row]
				gradient += row.weight * obs.sign * (logistic(logits[obs.row]) - row.game.blueWon)
			}
			// .25 is the upper bound on the logistic loss curvature for a signed unit feature.
			change := -gradient / (curvature + 0.25*col.weight)
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
