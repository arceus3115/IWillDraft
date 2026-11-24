package main

import (
	"fmt"
	"math"
	"slices"
	"time"
)

const minProTrainingWeight = 8

func enoughProHistory(training []proWeightedGame) bool {
	var mass float64
	for _, row := range training {
		mass += row.weight
	}
	return len(training) >= minProTrain && mass >= minProTrainingWeight
}

type proPrediction struct {
	game                 proGame
	probs                [len(proModelSpecs)]float64
	ready                [len(proModelSpecs)]bool
	priorLocalPatchGames int
}

func backtestPro(dataset proDataset, year int, resets []proPatch) []proPrediction {
	type fits struct {
		models               [len(proModelSpecs)]proModel
		ready                [len(proModelSpecs)]bool
		priorLocalPatchGames int
	}
	var predictions []proPrediction
	var week time.Time
	var models map[proContext]fits
	for _, game := range dataset.games {
		if game.year != year {
			continue
		}
		// The game's series shares one cutoff, including series that cross midnight.
		cutoff := proBacktestCutoff(game.startOfSeries)
		if !week.Equal(cutoff) {
			week, models = cutoff, map[proContext]fits{}
		}
		target := proContext{game.league, game.patch, cutoff}
		fit, exists := models[target]
		if !exists {
			for i, spec := range proModelSpecs {
				training := proTraining(dataset, target, resets, spec)
				if spec.samePatchOnly {
					for _, row := range training {
						if row.game.league == target.league {
							fit.priorLocalPatchGames++
						}
					}
				}
				if enoughProHistory(training) {
					fit.models[i], fit.ready[i] = fitProModel(training, spec.draft), true
				}
			}
			models[target] = fit
		}
		prediction := proPrediction{game: game, ready: fit.ready, priorLocalPatchGames: fit.priorLocalPatchGames}
		for i, model := range fit.models {
			if fit.ready[i] {
				prediction.probs[i] = logistic(model.score(game))
			}
		}
		predictions = append(predictions, prediction)
	}
	return predictions
}

func proBacktestCutoff(start time.Time) time.Time {
	day := start.UTC().Truncate(24 * time.Hour)
	return day.AddDate(0, 0, -(int(day.Weekday())+6)%7)
}

type proCalibration struct {
	Band        int     `json:"band"`
	Games       int     `json:"games"`
	Mean        float64 `json:"mean_predicted"`
	BlueWinRate float64 `json:"observed_blue_win_rate"`
}

type proAccuracy struct {
	Games       int              `json:"games"`
	Brier       float64          `json:"brier"`
	LogLoss     float64          `json:"log_loss"`
	Calibration []proCalibration `json:"calibration"`
}

func measureProAccuracy(rows []proPrediction, index int) *proAccuracy {
	out := proAccuracy{}
	var bins [5]proCalibration
	for _, row := range rows {
		if !row.ready[index] {
			continue
		}
		p, y := row.probs[index], row.game.blueWon
		out.Brier += (p - y) * (p - y)
		pClipped := min(max(p, 1e-12), 1-1e-12)
		out.LogLoss -= y*math.Log(pClipped) + (1-y)*math.Log1p(-pClipped)
		out.Games++
		band := min(int(p*5), len(bins)-1)
		bin := &bins[band]
		bin.Band = band
		bin.Games++
		bin.Mean += p
		bin.BlueWinRate += y
	}
	if out.Games == 0 {
		return nil
	}
	out.Brier /= float64(out.Games)
	out.LogLoss /= float64(out.Games)
	for _, bin := range bins {
		if bin.Games > 0 {
			bin.Mean /= float64(bin.Games)
			bin.BlueWinRate /= float64(bin.Games)
			out.Calibration = append(out.Calibration, bin)
		}
	}
	return &out
}

type proComparison struct {
	Games       int         `json:"paired_games"`
	Series      int         `json:"paired_series"`
	BrierDelta  float64     `json:"cross_league_cross_patch_minus_comparator_brier"`
	Bootstrap95 *[2]float64 `json:"paired_series_bootstrap_95,omitempty"`
}

func compareProAccuracy(rows []proPrediction, comparator int) *proComparison {
	const full = len(proModelSpecs) - 1
	totals := map[string]pairedCount{}
	var order []string
	out := proComparison{}
	for _, row := range rows {
		if !row.ready[full] || !row.ready[comparator] {
			continue
		}
		key := row.game.league + "/" + row.game.series
		s, exists := totals[key]
		if !exists {
			order = append(order, key)
		}
		a, b := row.probs[full]-row.game.blueWon, row.probs[comparator]-row.game.blueWon
		delta := a*a - b*b
		s.difference += delta
		s.eligible++
		totals[key] = s
		out.BrierDelta += delta
		out.Games++
	}
	if out.Games == 0 {
		return nil
	}
	out.Series = len(totals)
	out.BrierDelta /= float64(out.Games)
	if out.Series >= 10 {
		paired := make([]pairedCount, 0, len(order))
		for _, key := range order {
			paired = append(paired, totals[key])
		}
		lo, hi := bootstrapDelta(paired)
		out.Bootstrap95 = &[2]float64{lo, hi}
	}
	return &out
}

type proBacktestReport struct {
	League      string                    `json:"league"`
	Patch       string                    `json:"patch,omitempty"`
	PatchPhase  string                    `json:"patch_phase,omitempty"`
	Games       int                       `json:"eligible_games"`
	Accuracy    map[string]*proAccuracy   `json:"model_accuracy"`
	Comparisons map[string]*proComparison `json:"paired_model_comparisons"`
}

func summarizeProBacktest(predictions []proPrediction, year int) ([]proBacktestReport, error) {
	var reports []proBacktestReport
	for _, league := range tierOneLeagues {
		all := []proPrediction{}
		var early, established []proPrediction
		patches := map[proPatch][]proPrediction{}
		var patchOrder []proPatch
		for _, row := range predictions {
			g := row.game
			if g.league != league || g.year != year {
				continue
			}
			if _, ok := patches[g.patch]; !ok {
				patchOrder = append(patchOrder, g.patch)
			}
			patches[g.patch] = append(patches[g.patch], row)
			all = append(all, row)
			if row.priorLocalPatchGames < 15 {
				early = append(early, row)
			} else {
				established = append(established, row)
			}
		}
		reports = append(reports, makeProBacktestReport(league, "", "", all))
		if len(early) > 0 {
			reports = append(reports, makeProBacktestReport(league, "", "under_15_prior_local_patch_games", early))
		}
		if len(established) > 0 {
			reports = append(reports, makeProBacktestReport(league, "", "at_least_15_prior_local_patch_games", established))
		}
		slices.SortFunc(patchOrder, func(a, b proPatch) int {
			if a.version != b.version {
				return a.version - b.version
			}
			return a.update - b.update
		})
		for _, patch := range patchOrder {
			reports = append(reports, makeProBacktestReport(league, patch.String(), "", patches[patch]))
		}
	}
	for _, report := range reports {
		if len(report.Accuracy) > 0 {
			return reports, nil
		}
	}
	return nil, fmt.Errorf("no %d holdouts had enough earlier training evidence", year)
}

func makeProBacktestReport(league, patch, phase string, rows []proPrediction) proBacktestReport {
	report := proBacktestReport{league, patch, phase, len(rows), map[string]*proAccuracy{}, map[string]*proComparison{}}
	for i, spec := range proModelSpecs {
		if accuracy := measureProAccuracy(rows, i); accuracy != nil {
			report.Accuracy[spec.name] = accuracy
		}
		if i < len(proModelSpecs)-1 {
			if comparison := compareProAccuracy(rows, i); comparison != nil {
				report.Comparisons[spec.name] = comparison
			}
		}
	}
	return report
}
