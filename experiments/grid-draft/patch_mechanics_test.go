package main

import (
	"math"
	"testing"
)

func TestTrinityForcePatchImpactDependsOnFightLength(t *testing.T) {
	input := patchMechanicsInput{
		Item: "Trinity Force", BeforePatch: "13.13", AfterPatch: "13.14",
		Before: patchItemStats{FlatAD: 40, BonusAttackSpeed: .35, BaseADPerStack: .04, MaxStacks: 5},
		After:  patchItemStats{FlatAD: 40, BonusAttackSpeed: .33, BaseADPerStack: .06, MaxStacks: 3},
		Profiles: []patchDamageProfile{
			{Name: "short", Weight: 3, BaseAD: 100, TargetArmor: 100, Hits: 3},
			{Name: "long", Weight: 1, BaseAD: 100, TargetArmor: 100, Hits: 16},
		},
	}
	out, err := measurePatchMechanics(input)
	if err != nil {
		t.Fatal(err)
	}
	short, long := out.Profiles[0], out.Profiles[1]
	if short.Total.Before != 216 || short.Total.After != 219 ||
		long.Total.Before != 1250 || long.Total.After != 1246 {
		t.Fatalf("unexpected damage: %+v", out.Profiles)
	}
	want := 100 * ((3*219.0+1246)/(3*216.0+1250) - 1)
	if out.Total.PercentChange == nil || math.Abs(*out.Total.PercentChange-want) > 1e-10 {
		t.Fatalf("expected percent of the standardized damage means: %+v", out.Total)
	}
}

func TestAttackSpeedCanChangeProcCountWithinFixedWindow(t *testing.T) {
	profile := patchDamageProfile{
		BaseAD: 100, Seconds: 3, BaseAttackSpeed: 1, AttackSpeedRatio: 1, AttackSpeedCap: 2.5,
	}
	procs := map[int]bool{4: true}
	before, beforeHits := patchPhysicalDamage(profile, patchItemStats{BonusAttackSpeed: .35, SpellbladeBaseADRatio: 2}, procs)
	after, afterHits := patchPhysicalDamage(profile, patchItemStats{BonusAttackSpeed: .33, SpellbladeBaseADRatio: 2}, procs)
	if beforeHits != 4 || afterHits != 3 || before != 600 || after != 300 {
		t.Fatalf("damage or hit counts: before %g/%d, after %g/%d", before, beforeHits, after, afterHits)
	}
}
