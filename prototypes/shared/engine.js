/**
 * IWillDraft Empirical Analytics Engine (shared/engine.js)
 * Backend statistical & game-theory models:
 *  1. Hard Fearless Lockout & Value Over Replacement Pick (VORP)
 *  2. Prospect Theory (Kahneman-Tversky value & probability weighting)
 *  3. Quantized Game-Phase EPA, Success Rate & Betrayal Index
 *  4. Patch-as-Post-Training KL-Regularized Equilibrium & Chess-Tree Solver
 */
(function (global) {
  'use strict';

  // --- 1. HARD FEARLESS CORE ---
  // Canonical Bo5 Hard Fearless cumulative lockout sequence (10 champs per game)
  const FEARLESS_LOCK_SCHEDULE = {
    1: [],
    2: [
      'Azir', 'Orianna', 'Vi', 'Sejuani', 'K\'Sante',
      'Rumble', 'Kalista', 'Varus', 'Rell', 'Nautilus'
    ],
    3: [
      'Azir', 'Orianna', 'Vi', 'Sejuani', 'K\'Sante',
      'Rumble', 'Kalista', 'Varus', 'Rell', 'Nautilus',
      'Ahri', 'Corki', 'Xin Zhao', 'Maokai', 'Renekton',
      'Gnar', 'Jinx', 'Ashe', 'Alistar', 'Leona'
    ],
    4: [
      'Azir', 'Orianna', 'Vi', 'Sejuani', 'K\'Sante',
      'Rumble', 'Kalista', 'Varus', 'Rell', 'Nautilus',
      'Ahri', 'Corki', 'Xin Zhao', 'Maokai', 'Renekton',
      'Gnar', 'Jinx', 'Ashe', 'Alistar', 'Leona',
      'Galio', 'Taliyah', 'Skarner', 'Wukong', 'Jax',
      'Ambessa', 'Ezreal', 'Kai\'Sa', 'Rakan', 'Braum'
    ],
    5: [
      'Azir', 'Orianna', 'Vi', 'Sejuani', 'K\'Sante',
      'Rumble', 'Kalista', 'Varus', 'Rell', 'Nautilus',
      'Ahri', 'Corki', 'Xin Zhao', 'Maokai', 'Renekton',
      'Gnar', 'Jinx', 'Ashe', 'Alistar', 'Leona',
      'Galio', 'Taliyah', 'Skarner', 'Wukong', 'Jax',
      'Ambessa', 'Ezreal', 'Kai\'Sa', 'Rakan', 'Braum',
      'Aurora', 'Syndra', 'Poppy', 'Jarvan IV', 'Aatrox',
      'Ornn', 'Xayah', 'Lucian', 'Renata', 'Lulu'
    ]
  };

  function getLockedChampions(game, extraBans) {
    const g = Math.max(1, Math.min(5, Number(game) || 1));
    const base = FEARLESS_LOCK_SCHEDULE[g] || [];
    const set = new Set(base);
    if (Array.isArray(extraBans)) {
      extraBans.forEach(function (c) { set.add(c); });
    }
    return set;
  }

  function getLockedCount(game) {
    const g = Math.max(1, Math.min(5, Number(game) || 1));
    return (g - 1) * 10;
  }

  /**
   * Computes Fearless Value Over Replacement Pick (VORP) in WPA% points.
   * As higher-EPA picks lock out in a role, replacement baseline drops and remaining
   * picks gain scarcity VORP.
   */
  function computeFearlessVORP(pool, lockedSet) {
    const available = pool.filter(function (c) { return !lockedSet.has(c.champ); });
    if (available.length === 0) {
      return { replacementEpa: -3.5, available: [], avgEpa: -3.5, cliffDelta: -4.0 };
    }
    const sorted = available.slice().sort(function (a, b) { return b.epa - a.epa; });
    // Replacement level = 3rd best remaining pick (or worst remaining if <= 2 left)
    const repIdx = Math.min(2, sorted.length - 1);
    const replacementEpa = sorted[repIdx].epa;
    const enriched = sorted.map(function (c) {
      return Object.assign({}, c, {
        vorp: Number((c.epa - replacementEpa).toFixed(2))
      });
    });
    const avgEpa = Number((sorted.reduce(function (s, c) { return s + c.epa; }, 0) / sorted.length).toFixed(2));
    const topEpa = sorted[0].epa;
    return {
      replacementEpa: Number(replacementEpa.toFixed(2)),
      topEpa: Number(topEpa.toFixed(2)),
      avgEpa: avgEpa,
      cliffDelta: Number((topEpa - replacementEpa).toFixed(2)),
      available: enriched
    };
  }

  /**
   * Re-normalizes conditional branch probabilities after removing Fearless-locked champions.
   */
  function renormalizeBranches(branches, lockedSet, champExtractor) {
    const active = branches.map(function (b) {
      const champs = champExtractor(b);
      const locked = champs.some(function (c) { return lockedSet.has(c); });
      return Object.assign({}, b, { fearlessLocked: locked });
    });
    const legalSum = active.reduce(function (s, b) {
      return s + (b.fearlessLocked ? 0 : b.baseRate);
    }, 0);

    return active.map(function (b) {
      if (b.fearlessLocked || legalSum <= 0) {
        return Object.assign({}, b, { adjRate: 0 });
      }
      return Object.assign({}, b, {
        adjRate: Math.round((b.baseRate / legalSum) * 100)
      });
    });
  }

  // --- 2. PROSPECT THEORY ENGINE (Kahneman & Tversky) ---
  const PT_DEFAULTS = {
    alpha: 0.88,   // concavity for gains
    beta: 0.88,    // convexity for losses
    lambda: 2.25,  // loss aversion coefficient
    gamma: 0.61    // probability weighting curvature
  };

  function ptValue(x, params) {
    const p = Object.assign({}, PT_DEFAULTS, params || {});
    if (x >= 0) return Math.pow(x, p.alpha);
    return -p.lambda * Math.pow(-x, p.beta);
  }

  function ptWeight(prob, params) {
    const p = Object.assign({}, PT_DEFAULTS, params || {});
    const clamped = Math.max(0.001, Math.min(0.999, prob));
    const num = Math.pow(clamped, p.gamma);
    const den = Math.pow(num + Math.pow(1 - clamped, p.gamma), 1 / p.gamma);
    return num / den;
  }

  /**
   * Evaluates a binary decision (e.g., Contest vs. Trade, Burn Pick G1 vs. Hold G4)
   * Returns risk-neutral EV (in ΔWP points), Prospect-Theory Utility U_PT, and Risk-Bias Gap.
   */
  function evaluateProspectDecision(pSuccess, gainWP, lossWP, params) {
    const pFail = 1 - pSuccess;
    const ev = pSuccess * gainWP + pFail * lossWP; // lossWP is negative
    const wSuccess = ptWeight(pSuccess, params);
    const wFail = ptWeight(pFail, params);
    const uPT = wSuccess * ptValue(gainWP, params) + wFail * ptValue(lossWP, params);
    // Normalize uPT to same rough scale as EV for intuitive comparison
    const normalizedPT = uPT / 1.45;
    return {
      ev: Number(ev.toFixed(2)),
      uPT: Number(normalizedPT.toFixed(2)),
      wSuccess: Number((wSuccess * 100).toFixed(1)),
      wFail: Number((wFail * 100).toFixed(1)),
      biasGap: Number((normalizedPT - ev).toFixed(2))
    };
  }

  // --- 3. QUANTIZED GAME-PHASE EPA & BETRAYAL INDEX ---
  /**
   * Betrayal Index (BI): adapted from NFL QB EPA vs. Defense EPA Allowed.
   * Measures how severely a unit/phase's positive EPA is offset by team/downstream negative EPA.
   * Positive BI = unit is being betrayed by team/later phases.
   * Negative BI = unit is betraying teammates/earlier leads.
   */
  function computeBetrayalIndex(unitEpa, restOfTeamEpa) {
    // Difference weighted by whether signs oppose (amplifies true betrayal)
    const rawDiff = unitEpa - restOfTeamEpa;
    const opposing = (unitEpa > 0 && restOfTeamEpa < 0) || (unitEpa < 0 && restOfTeamEpa > 0);
    const multiplier = opposing ? 1.25 : 0.85;
    return Number((rawDiff * multiplier).toFixed(2));
  }

  /**
   * Computes composite phase metrics adjusted by Fearless Game volatility.
   * Later Fearless games (G4-G5) have higher draft variance -> higher Leverage Index & wider EPA variance.
   */
  function adjustPhaseForFearless(phase, game) {
    const g = Math.max(1, Math.min(5, Number(game) || 1));
    const fearlessFactor = 1 + (g - 1) * 0.08; // +0% G1 to +32% G5 volatility
    const adjLi = Number((phase.li * (1 + (g - 1) * 0.05)).toFixed(2));
    const adjEpa = Number((phase.epa * fearlessFactor).toFixed(2));
    const adjSr = Math.max(25, Math.min(82, Math.round(phase.sr + (phase.epa >= 0 ? (g - 1) * 0.8 : -(g - 1) * 1.4))));
    return Object.assign({}, phase, {
      adjEpa: adjEpa,
      adjLi: adjLi,
      adjSr: adjSr
    });
  }

  // --- 4. PATCH-AS-POST-TRAINING EQUILIBRIUM SOLVER ---
  /**
   * Models a League patch as KL-regularized policy post-training:
   *   pi_1(a) proportional to pi_0(a) * exp(DeltaR(a) / tau)
   * where:
   *   pi_0(a) = pre-patch equilibrium share (base policy)
   *   DeltaR(a) = dot product of patch reward perturbations and archetype exposures
   *   tau = KL inertia (scrim/muscle-memory regularization; lower tau = faster convergence)
   */
  function solvePostTrainingEquilibrium(archetypes, rewardShift, tau, fearlessGame) {
    const temperature = Math.max(0.25, Number(tau) || 1.0);
    const g = Math.max(1, Math.min(5, Number(fearlessGame) || 1));
    const lockedSet = getLockedChampions(g);

    const unnormalized = archetypes.map(function (arch) {
      // Dot product of patch reward vector with archetype sensitivity vector
      const deltaR =
        (rewardShift.dragon || 0) * (arch.sens.dragon || 0) +
        (rewardShift.grubs || 0) * (arch.sens.grubs || 0) +
        (rewardShift.durability || 0) * (arch.sens.durability || 0) +
        (rewardShift.itemCurve || 0) * (arch.sens.itemCurve || 0);

      // Fearless pool availability penalty: if key archetype enablers are locked, effective reward drops
      const lockedEnablers = (arch.enablers || []).filter(function (c) { return lockedSet.has(c); }).length;
      const totalEnablers = Math.max(1, (arch.enablers || []).length);
      const availRatio = (totalEnablers - lockedEnablers) / totalEnablers;
      const fearlessPenalty = (1 - availRatio) * 1.15;

      const netDeltaR = deltaR - fearlessPenalty;
      const weight = arch.basePi * Math.exp(netDeltaR / temperature);
      return Object.assign({}, arch, {
        deltaR: Number(netDeltaR.toFixed(2)),
        rawDeltaR: Number(deltaR.toFixed(2)),
        lockedEnablers: lockedEnablers,
        totalEnablers: totalEnablers,
        weight: weight
      });
    });

    const Z = unnormalized.reduce(function (s, a) { return s + a.weight; }, 0) || 1;
    let klDivergence = 0;

    const solved = unnormalized.map(function (a) {
      const pi1 = a.weight / Z;
      if (pi1 > 0.0001 && a.basePi > 0.0001) {
        klDivergence += pi1 * Math.log(pi1 / a.basePi);
      }
      return Object.assign({}, a, {
        pi0Pct: Number((a.basePi * 100).toFixed(1)),
        pi1Pct: Number((pi1 * 100).toFixed(1)),
        shiftPct: Number(((pi1 - a.basePi) * 100).toFixed(1))
      });
    });

    return {
      archetypes: solved,
      klDivergence: Number(klDivergence.toFixed(3))
    };
  }

  /**
   * Solves Ply-2 chess-like response distribution given a Ply-1 opening move and current patch reward shift.
   */
  function solveChessResponse(opening, rewardShift, tau, fearlessGame) {
    const temperature = Math.max(0.25, Number(tau) || 1.0);
    const lockedSet = getLockedChampions(fearlessGame);

    const scored = opening.responses.map(function (r) {
      const locked = (r.champs || []).some(function (c) { return lockedSet.has(c); });
      const patchBoost =
        (rewardShift.dragon || 0) * (r.sens.dragon || 0) +
        (rewardShift.grubs || 0) * (r.sens.grubs || 0) +
        (rewardShift.durability || 0) * (r.sens.durability || 0) +
        (rewardShift.itemCurve || 0) * (r.sens.itemCurve || 0);

      const postEv = Number((r.baseEv + patchBoost * 1.8).toFixed(2));
      const weight = locked ? 0 : r.baseShare * Math.exp((postEv - r.baseEv) / (temperature * 2.2));
      return Object.assign({}, r, {
        locked: locked,
        postEv: postEv,
        evDelta: Number((postEv - r.baseEv).toFixed(2)),
        weight: weight
      });
    });

    const Z = scored.reduce(function (s, r) { return s + r.weight; }, 0);
    return scored.map(function (r) {
      return Object.assign({}, r, {
        postSharePct: Z > 0 ? Math.round((r.weight / Z) * 100) : 0
      });
    });
  }

  global.IWillDraftEngine = {
    FEARLESS_LOCK_SCHEDULE: FEARLESS_LOCK_SCHEDULE,
    getLockedChampions: getLockedChampions,
    getLockedCount: getLockedCount,
    computeFearlessVORP: computeFearlessVORP,
    renormalizeBranches: renormalizeBranches,
    PT_DEFAULTS: PT_DEFAULTS,
    ptValue: ptValue,
    ptWeight: ptWeight,
    evaluateProspectDecision: evaluateProspectDecision,
    computeBetrayalIndex: computeBetrayalIndex,
    adjustPhaseForFearless: adjustPhaseForFearless,
    solvePostTrainingEquilibrium: solvePostTrainingEquilibrium,
    solveChessResponse: solveChessResponse
  };
})(window);
