package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"
)

type proQueryLineup struct {
	Team  string             `json:"team_id"`
	Roles map[string]proPick `json:"roles"`
}

type proQuery struct {
	League        string         `json:"league"`
	Patch         string         `json:"patch"`
	AsOf          time.Time      `json:"as_of"`
	FirstPickSide string         `json:"first_pick_side"`
	FocusSide     string         `json:"focus_side"`
	FocusRole     string         `json:"focus_role"`
	Blue          proQueryLineup `json:"blue"`
	Red           proQueryLineup `json:"red"`
	Unavailable   []string       `json:"unavailable"`
	Reference     string         `json:"reference"`
	Candidates    []string       `json:"candidates"`
}

func readProQuery(path string) (proQuery, error) {
	file, err := os.Open(path)
	if err != nil {
		return proQuery{}, err
	}
	defer file.Close()
	var query proQuery
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&query); err != nil {
		return proQuery{}, fmt.Errorf("read draft scenario: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return proQuery{}, fmt.Errorf("provide one JSON draft scenario")
	}
	return query, nil
}

func (q proQuery) prepare() (proGame, map[string]bool, error) {
	game := proGame{league: q.League}
	patch, err := parseProPatch(q.Patch)
	if err != nil {
		return game, nil, err
	}
	game.patch = patch
	firstSide, focusSide := slices.Index(proSideNames, q.FirstPickSide), slices.Index(proSideNames, q.FocusSide)
	focusRole := slices.Index(proRoles, q.FocusRole)
	if q.AsOf.IsZero() || !slices.Contains(tierOneLeagues, q.League) || firstSide < 0 || focusSide < 0 || focusRole < 0 {
		return game, nil, fmt.Errorf("provide a tier-one league, as_of, first_pick_side, focus_side and focus_role")
	}
	game.firstPick = 1 - 2*firstSide
	if q.Unavailable == nil || len(q.Candidates) < 2 {
		return game, nil, fmt.Errorf("provide unavailable champions (prior Fearless locks, current bans and disables) and at least two candidates")
	}
	used := map[string]bool{}
	for _, name := range q.Unavailable {
		champion := championKey(name)
		if champion == "" {
			return game, nil, fmt.Errorf("empty unavailable champion")
		}
		used[champion] = true
	}
	players := map[string]bool{}
	for side, lineup := range [proSides]proQueryLineup{q.Blue, q.Red} {
		if strings.TrimSpace(lineup.Team) == "" || len(lineup.Roles) != len(proRoles) {
			return game, nil, fmt.Errorf("%s needs a team_id and all five roles", proSideNames[side])
		}
		game.sides[side].Team = strings.TrimSpace(lineup.Team)
		for i, role := range proRoles {
			pick := lineup.Roles[role]
			pick.Player = strings.TrimSpace(pick.Player)
			pick.Champion = championKey(pick.Champion)
			if pick.Player == "" || players[pick.Player] {
				return game, nil, fmt.Errorf("%s %s needs a distinct player_id", proSideNames[side], role)
			}
			players[pick.Player] = true
			if side == focusSide && i == focusRole {
				if pick.Champion != "" {
					return game, nil, fmt.Errorf("leave the focus role's champion empty")
				}
			} else {
				if pick.Champion == "" || used[pick.Champion] {
					return game, nil, fmt.Errorf("%s %s has a missing, duplicate or unavailable champion", proSideNames[side], role)
				}
				used[pick.Champion] = true
			}
			game.sides[side].Picks[i] = pick
		}
	}
	if game.sides[0].Team == game.sides[1].Team {
		return game, nil, fmt.Errorf("team ids must differ")
	}
	candidates := map[string]bool{}
	for _, value := range q.Candidates {
		champion := championKey(value)
		if champion == "" || candidates[champion] {
			return game, nil, fmt.Errorf("provide distinct nonempty candidates")
		}
		candidates[champion] = true
	}
	if !candidates[championKey(q.Reference)] || used[championKey(q.Reference)] {
		return game, nil, fmt.Errorf("reference must be a legal candidate")
	}
	return game, used, nil
}

type proTrainingEvidence struct {
	Games          int            `json:"games"`
	Series         int            `json:"series"`
	Weight         float64        `json:"total_weight"`
	EffectiveGames float64        `json:"effective_games"`
	ByPatch        map[string]int `json:"games_by_patch"`
	LastGame       time.Time      `json:"last_game"`
}

func describeProTraining(rows []proWeightedGame) map[string]*proTrainingEvidence {
	out := map[string]*proTrainingEvidence{}
	squares := map[string]float64{}
	series := map[string]bool{}
	for _, row := range rows {
		g := row.game
		e := out[g.league]
		if e == nil {
			e = &proTrainingEvidence{ByPatch: map[string]int{}}
			out[g.league] = e
		}
		e.Games++
		e.Weight += row.weight
		e.ByPatch[g.patch.String()]++
		if g.start.After(e.LastGame) {
			e.LastGame = g.start
		}
		key := g.league + "/" + g.series
		if !series[key] {
			e.Series++
			series[key] = true
		}
		squares[g.league] += row.weight * row.weight
	}
	for league, e := range out {
		e.EffectiveGames = e.Weight * e.Weight / squares[league]
	}
	return out
}

type proOptionEvidence struct {
	RoleGamesByLeague     map[string]int `json:"role_picks_by_league"`
	RoleWeight            float64        `json:"weighted_role_picks"`
	TargetPatchRoleGames  int            `json:"target_patch_role_picks"`
	LocalPatchRoleGames   int            `json:"target_league_patch_role_picks"`
	LocalPatchBans        int            `json:"target_league_patch_recorded_bans"`
	PlayerRoleGames       int            `json:"player_champion_role_games"`
	OpponentMatchupGames  int            `json:"lane_matchup_games"`
	AlliedPairGamesByRole map[string]int `json:"allied_pair_games_by_role"`
}

type proDraftOption struct {
	Champion           string             `json:"champion"`
	Status             string             `json:"status"`
	Evidence           *proOptionEvidence `json:"evidence,omitempty"`
	Probability        map[string]float64 `json:"win_probability_by_model,omitempty"`
	Delta              map[string]float64 `json:"probability_delta_vs_reference_by_model,omitempty"`
	Bootstrap95        *[2]float64        `json:"series_bootstrap_delta_95,omitempty"`
	LogOddsVsReference map[string]float64 `json:"log_odds_delta_vs_reference_by_feature,omitempty"`
}

type proDecision struct {
	Status              string                          `json:"status"`
	League              string                          `json:"league"`
	Patch               string                          `json:"patch"`
	PatchSource         string                          `json:"patch_source"`
	AsOf                time.Time                       `json:"as_of"`
	Reference           string                          `json:"reference"`
	TeamSideProbability *float64                        `json:"team_player_side_probability,omitempty"`
	UnseenTeamIDs       []string                        `json:"unseen_team_ids,omitempty"`
	UnseenPlayerIDs     []string                        `json:"unseen_player_ids,omitempty"`
	ThinContextRoles    []string                        `json:"thin_context_roles,omitempty"`
	Training            map[string]*proTrainingEvidence `json:"transfer_training_by_league"`
	Options             []proDraftOption                `json:"options,omitempty"`
}

func proOptionSupport(training []proWeightedGame, model proModel, target proGame, side, role int) *proOptionEvidence {
	pick := target.sides[side].Picks[role]
	out := &proOptionEvidence{RoleGamesByLeague: map[string]int{}, AlliedPairGamesByRole: map[string]int{}}
	for _, row := range training {
		g := row.game
		if g.league == target.league && g.patch == target.patch && g.bans[pick.Champion] {
			out.LocalPatchBans++
		}
		for _, lineup := range g.sides {
			if lineup.Picks[role].Champion != pick.Champion {
				continue
			}
			out.RoleGamesByLeague[g.league]++
			out.RoleWeight += row.weight
			if g.patch == target.patch {
				out.TargetPatchRoleGames++
				if g.league == target.league {
					out.LocalPatchRoleGames++
				}
			}
			if lineup.Picks[role].Player == pick.Player {
				out.PlayerRoleGames++
			}
		}
	}
	for _, term := range proTerms(target, true) {
		f := term.feature
		switch f.kind {
		case matchupEffect:
			if f.role == role {
				out.OpponentMatchupGames = model.coef[f].games
			}
		case synergyEffect:
			if f.role == role && f.a == pick.Champion {
				out.AlliedPairGamesByRole[proRoles[f.partner]] = model.coef[f].games
			} else if f.partner == role && f.b == pick.Champion {
				out.AlliedPairGamesByRole[proRoles[f.role]] = model.coef[f].games
			}
		}
	}
	return out
}

func compareProDrafts(dataset proDataset, resets []proPatch, q proQuery) (proDecision, error) {
	patchSource := "scenario"
	if strings.TrimSpace(q.Patch) == "" {
		window, source := dataset.calendar.patchAt(q.League, q.AsOf)
		if source == "" {
			return proDecision{}, fmt.Errorf("no verified %s patch window at %s, provide the competitive patch in the scenario", q.League, q.AsOf.Format(time.RFC3339))
		}
		q.Patch, patchSource = window.Patch.String(), source
	}
	target, unavailable, err := q.prepare()
	if err != nil {
		return proDecision{}, err
	}
	const full = len(proModelSpecs) - 1
	context := proContext{q.League, target.patch, q.AsOf}
	training := proTraining(dataset, context, resets, proModelSpecs[full])
	out := proDecision{
		Status: "insufficient_training", League: q.League, Patch: target.patch.String(), PatchSource: patchSource, AsOf: q.AsOf,
		Reference: q.Reference, Training: describeProTraining(training),
	}
	if !enoughProHistory(training) {
		return out, nil
	}
	out.Status = "scored"
	var models [len(proModelSpecs)]proModel
	models[full] = fitProModel(training, true)
	models[0] = fitProModel(training, false)
	for i := 1; i < full; i++ {
		rows := proTraining(dataset, context, resets, proModelSpecs[i])
		if enoughProHistory(rows) {
			models[i] = fitProModel(rows, true)
		}
	}

	for _, term := range proTerms(target, false) {
		if models[0].coef[term.feature].games > 0 {
			continue
		}
		switch term.feature.kind {
		case teamEffect:
			out.UnseenTeamIDs = append(out.UnseenTeamIDs, term.feature.b)
		case playerEffect:
			out.UnseenPlayerIDs = append(out.UnseenPlayerIDs, term.feature.b)
		}
	}
	side, role := slices.Index(proSideNames, q.FocusSide), slices.Index(proRoles, q.FocusRole)
	for lineupSide, lineup := range target.sides {
		for slot := range lineup.Picks {
			if (lineupSide != side || slot != role) && !proRoleSupported(models[full], target, lineupSide, slot) {
				out.ThinContextRoles = append(out.ThinContextRoles, proSideNames[lineupSide]+"/"+proRoles[slot])
			}
		}
	}
	sign := 1 - 2*float64(side)
	baseline := logistic(sign * models[0].score(target))
	out.TeamSideProbability = &baseline
	reference := target
	reference.sides[side].Picks[role].Champion = championKey(q.Reference)
	referenceSupported := proRoleSupported(models[full], reference, side, role)
	var branches []proGame
	var scoredIndices []int
	for _, name := range q.Candidates {
		champion := championKey(name)
		option := proDraftOption{Champion: name, Status: "unavailable"}
		if unavailable[champion] {
			out.Options = append(out.Options, option)
			continue
		}
		branch := target
		branch.sides[side].Picks[role].Champion = champion
		option.Evidence = proOptionSupport(training, models[full], branch, side, role)
		option.Status = "insufficient_role_evidence"
		if !proRoleSupported(models[full], branch, side, role) {
			out.Options = append(out.Options, option)
			continue
		}
		option.Status = "limited_local_patch_evidence"
		if option.Evidence.LocalPatchRoleGames >= 5 {
			option.Status = "scored"
		}
		option.Probability = map[string]float64{}
		option.Delta = map[string]float64{}
		for i := 1; i < len(models); i++ {
			model := models[i]
			if !proRoleSupported(model, branch, side, role) {
				continue
			}
			p := logistic(sign * model.score(branch))
			option.Probability[proModelSpecs[i].name] = p
			if proRoleSupported(model, reference, side, role) {
				option.Delta[proModelSpecs[i].name] = p - logistic(sign*model.score(reference))
			}
		}
		if referenceSupported {
			option.LogOddsVsReference = proScoreDelta(models[full], branch, reference, sign)
			if champion != championKey(q.Reference) {
				scoredIndices = append(scoredIndices, len(out.Options))
				branches = append(branches, branch)
			}
		}
		out.Options = append(out.Options, option)
	}
	for i, interval := range bootstrapProDrafts(training, reference, branches, sign) {
		out.Options[scoredIndices[i]].Bootstrap95 = &interval
	}
	if !referenceSupported {
		out.Status = "insufficient_reference_evidence"
	} else if len(branches) == 0 {
		out.Status = "no_supported_alternative"
	}
	return out, nil
}

func proRoleSupported(model proModel, game proGame, side, role int) bool {
	pick := game.sides[side].Picks[role]
	support := model.coef[proFeature{kind: championEffect, role: role, a: pick.Champion}]
	return support.games >= 3 && support.weight >= 1
}

func proScoreDelta(model proModel, branch, reference proGame, sign float64) map[string]float64 {
	var grouped [len(proFeatureNames)]float64
	for _, term := range proTerms(branch, true) {
		grouped[term.feature.kind] += sign * term.sign * model.coef[term.feature].value
	}
	for _, term := range proTerms(reference, true) {
		grouped[term.feature.kind] -= sign * term.sign * model.coef[term.feature].value
	}
	out := map[string]float64{}
	for kind := championEffect; int(kind) < len(grouped); kind++ {
		out[proFeatureNames[kind]] = grouped[kind]
	}
	return out
}

func bootstrapProDrafts(training []proWeightedGame, reference proGame, branches []proGame, sign float64) [][2]float64 {
	if len(branches) == 0 {
		return nil
	}
	type cluster struct {
		rows []int
	}
	var clusters []cluster
	byKey := map[string]int{}
	byLeague := map[string][]int{}
	for i, row := range training {
		key := row.game.league + "/" + row.game.series
		index, exists := byKey[key]
		if !exists {
			index = len(clusters)
			byKey[key] = index
			clusters = append(clusters, cluster{})
			byLeague[row.game.league] = append(byLeague[row.game.league], index)
		}
		clusters[index].rows = append(clusters[index].rows, i)
	}
	if len(clusters) < 10 {
		return nil
	}
	const draws = 150
	rng := rand.New(rand.NewPCG(1, 2))
	samples := make([][]float64, len(branches))
	multiplicity := make([]float64, len(clusters))
	resample := make([]proWeightedGame, 0, len(training))
	for range draws {
		clear(multiplicity)
		for _, league := range tierOneLeagues {
			stratum := byLeague[league]
			for range stratum {
				multiplicity[stratum[rng.IntN(len(stratum))]]++
			}
		}
		resample = resample[:0]
		for i, count := range multiplicity {
			if count > 0 {
				for _, row := range clusters[i].rows {
					resample = append(resample, proWeightedGame{training[row].game, count * training[row].weight})
				}
			}
		}
		model := fitProModel(resample, true)
		ref := logistic(sign * model.score(reference))
		for i, branch := range branches {
			samples[i] = append(samples[i], logistic(sign*model.score(branch))-ref)
		}
	}
	out := make([][2]float64, len(branches))
	for i, values := range samples {
		slices.Sort(values)
		out[i] = [2]float64{values[int(math.Floor(0.025*float64(draws-1)))], values[int(math.Ceil(0.975*float64(draws-1)))]}
	}
	return out
}
