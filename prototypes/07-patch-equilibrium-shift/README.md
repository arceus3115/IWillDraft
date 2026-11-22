# 07: Patch Equilibrium Shift (Post-Training & Chess-Tree Model)

**Prototype:** Equilibrium shift simulator (`index.html`, powered by `../shared/engine.js`) treating League patches as KL-regularized model post-training (`π₁ ∝ π₀ · exp(ΔR / τ)`) with Hard Fearless Game 1–5 lockouts and 2-ply chess-like opening cascades.

**Problem Solved:** Patch notes alter reward weights across objectives and items, shifting the pre-patch equilibrium in semi-predictable cascades. Adjusting reward perturbations (`ΔR`) and scrim inertia (`τ`) forecasts archetype probability shifts, Ply-2 optimal draft responses, and map-objective EPA weights.

**Traverse:** Open `index.html`, adjust patch `ΔR` sliders or presets, toggle `G1–G5` Fearless lockouts, and explore `Policy Shift`, `Chess-Tree Cascade`, and `Macro Objective Weights`.
