package main

import (
	"fmt"
	"math"
	"strings"
)

type patchItemStats struct {
	FlatAD                float64 `json:"flat_ad"`
	BonusAttackSpeed      float64 `json:"bonus_attack_speed"`
	BaseADPerStack        float64 `json:"base_ad_per_stack"`
	MaxStacks             int     `json:"max_stacks"`
	SpellbladeBaseADRatio float64 `json:"spellblade_base_ad_ratio"`
}

type patchDamageProfile struct {
	Name             string  `json:"name"`
	Weight           float64 `json:"weight"`
	BaseAD           float64 `json:"base_ad"`
	OtherBonusAD     float64 `json:"other_bonus_ad"`
	TargetArmor      float64 `json:"target_armor"`
	OtherRawPhysical float64 `json:"other_raw_physical_damage"`
	Hits             int     `json:"hits"`
	Seconds          float64 `json:"seconds"`
	BaseAttackSpeed  float64 `json:"base_attack_speed"`
	AttackSpeedRatio float64 `json:"attack_speed_ratio"`
	OtherBonusSpeed  float64 `json:"other_bonus_attack_speed"`
	AttackSpeedCap   float64 `json:"attack_speed_cap"`
	ProcHits         []int   `json:"spellblade_proc_hits"`
}

type patchMechanicsInput struct {
	Item        string               `json:"item"`
	BeforePatch string               `json:"before_patch"`
	AfterPatch  string               `json:"after_patch"`
	Before      patchItemStats       `json:"before"`
	After       patchItemStats       `json:"after"`
	Profiles    []patchDamageProfile `json:"profiles"`
}

type patchProfileResult struct {
	Name      string           `json:"name"`
	Weight    float64          `json:"weight"`
	Hits      [2]int           `json:"hits_before_after"`
	Total     patchDamageDelta `json:"total_physical_damage"`
	ItemAdded patchDamageDelta `json:"item_added_physical_damage"`
}

type patchDamageDelta struct {
	Before        float64  `json:"before"`
	After         float64  `json:"after"`
	PercentChange *float64 `json:"percent_change,omitempty"`
}

type patchMechanicsReport struct {
	Item        string               `json:"item"`
	BeforePatch string               `json:"before_patch"`
	AfterPatch  string               `json:"after_patch"`
	Profiles    []patchProfileResult `json:"profiles"`
	Total       patchDamageDelta     `json:"weighted_mean_total_physical_damage"`
	ItemAdded   patchDamageDelta     `json:"weighted_mean_item_added_physical_damage"`
}

func runPatchMechanics(path string) error {
	input, err := readJSONInput[patchMechanicsInput](path)
	if err != nil {
		return err
	}
	result, err := measurePatchMechanics(input)
	if err != nil {
		return err
	}
	return writeJSONResult(result)
}

func measurePatchMechanics(input patchMechanicsInput) (patchMechanicsReport, error) {
	out := patchMechanicsReport{Item: input.Item, BeforePatch: input.BeforePatch, AfterPatch: input.AfterPatch}
	beforePatch, beforeErr := parseProPatch(input.BeforePatch)
	afterPatch, afterErr := parseProPatch(input.AfterPatch)
	if strings.TrimSpace(input.Item) == "" || beforeErr != nil || afterErr != nil || beforePatch.version > afterPatch.version ||
		(beforePatch.version == afterPatch.version && beforePatch.update >= afterPatch.update) || len(input.Profiles) == 0 {
		return out, fmt.Errorf("provide an item, two ascending patch versions and at least one fight profile")
	}
	for _, item := range []patchItemStats{input.Before, input.After} {
		if item.FlatAD < 0 || item.BonusAttackSpeed < 0 || item.BaseADPerStack < 0 || item.MaxStacks < 0 || item.SpellbladeBaseADRatio < 0 {
			return out, fmt.Errorf("item damage and attack-speed parameters must be nonnegative")
		}
	}
	var weightSum float64
	totalDamage, itemDamage := [2]float64{}, [2]float64{}
	for _, profile := range input.Profiles {
		if strings.TrimSpace(profile.Name) == "" || profile.Weight <= 0 || profile.BaseAD <= 0 ||
			profile.OtherBonusAD < 0 || profile.TargetArmor < 0 || profile.OtherRawPhysical < 0 ||
			profile.Hits < 0 || profile.Seconds < 0 || (profile.Hits == 0) == (profile.Seconds == 0) {
			return out, fmt.Errorf("profile %q needs a positive weight and base AD, nonnegative damage and armor, and exactly one of hits or seconds", profile.Name)
		}
		if profile.Seconds > 0 && (profile.BaseAttackSpeed <= 0 || profile.AttackSpeedRatio <= 0 ||
			profile.AttackSpeedCap <= 0 || profile.OtherBonusSpeed < 0 || profile.Seconds*profile.AttackSpeedCap > 10000) {
			return out, fmt.Errorf("profile %q needs positive speed inputs, nonnegative other bonus speed and at most 10000 possible hits", profile.Name)
		}
		if profile.Hits > 10000 {
			return out, fmt.Errorf("profile %q exceeds 10000 hits", profile.Name)
		}
		procs := map[int]bool{}
		for _, hit := range profile.ProcHits {
			if hit <= 0 || (profile.Hits > 0 && hit > profile.Hits) || procs[hit] {
				return out, fmt.Errorf("profile %q has an invalid or duplicate Spellblade proc hit", profile.Name)
			}
			procs[hit] = true
		}
		baseline, _ := patchPhysicalDamage(profile, patchItemStats{}, procs)
		before, beforeHits := patchPhysicalDamage(profile, input.Before, procs)
		after, afterHits := patchPhysicalDamage(profile, input.After, procs)
		out.Profiles = append(out.Profiles, patchProfileResult{
			Name: profile.Name, Weight: profile.Weight, Hits: [2]int{beforeHits, afterHits},
			Total: damageDelta(before, after), ItemAdded: damageDelta(before-baseline, after-baseline),
		})
		weightSum += profile.Weight
		totalDamage[0] += profile.Weight * before
		totalDamage[1] += profile.Weight * after
		itemDamage[0] += profile.Weight * (before - baseline)
		itemDamage[1] += profile.Weight * (after - baseline)
	}
	for _, value := range []float64{weightSum, totalDamage[0], totalDamage[1], itemDamage[0], itemDamage[1]} {
		if math.IsInf(value, 0) || math.IsNaN(value) {
			return out, fmt.Errorf("damage calculation exceeds numeric range")
		}
	}
	out.Total = damageDelta(totalDamage[0]/weightSum, totalDamage[1]/weightSum)
	out.ItemAdded = damageDelta(itemDamage[0]/weightSum, itemDamage[1]/weightSum)
	return out, nil
}

func patchPhysicalDamage(profile patchDamageProfile, item patchItemStats, procs map[int]bool) (float64, int) {
	hits := profile.Hits
	if profile.Seconds > 0 {
		speed := min(profile.AttackSpeedCap, profile.BaseAttackSpeed+
			profile.AttackSpeedRatio*(profile.OtherBonusSpeed+item.BonusAttackSpeed))
		hits = int(math.Floor(profile.Seconds * speed))
	}
	raw := profile.OtherRawPhysical
	for hit := 1; hit <= hits; hit++ {
		baseAD := profile.BaseAD * (1 + item.BaseADPerStack*float64(min(hit-1, item.MaxStacks)))
		raw += baseAD + profile.OtherBonusAD + item.FlatAD
		if procs[hit] {
			raw += item.SpellbladeBaseADRatio * baseAD
		}
	}
	return raw * (100 / (100 + profile.TargetArmor)), hits
}

func damageDelta(before, after float64) patchDamageDelta {
	out := patchDamageDelta{Before: before, After: after}
	if before > 0 {
		change := 100 * (after/before - 1)
		if !math.IsInf(change, 0) && !math.IsNaN(change) {
			out.PercentChange = &change
		}
	}
	return out
}
