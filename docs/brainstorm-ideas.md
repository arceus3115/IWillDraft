# Scouting & Draft Support Tool — Past Winners Audit, Pruned Brainstorm, & Idea Evaluation

> **Core Design Principle — Augment Staff, Never Autopilot:**
> Professional League of Legends coaches and analysts reject black-box AI telling them *"lock Champion X (Score: 0.84 / 53.8% Win Prob)"*—automated solvers ignore private scrim results, player confidence, comms, and stage-day intuition. Furthermore, Riot tournament rules prohibit electronic devices on stage during live pick/ban; software delivers real value during **weekly patch/match prep, opponent scouting, 5-minute between-game locker-room huddles, and generating Riot-legal stage sheets**.

---

## 1. Past Hackathon Winners Audit (Useful vs. Unhelpful Features)

Analysis based on Devpost submissions and the official winning repositories archived under `github.com/cloud9esports`:

### Winner 1: C9 StratOS (`Skys-The-Limit-Category-1` — Category 1 Assistant Coach)
* **What It Built**: Next.js + MongoDB **VALORANT-only** assistant coach (`0%` League of Legends coverage) combining match dashboards, LLM coaching reports, and an interactive chat/scenario panel.
* **What WAS Useful (Keep as Engineering Baseline)**:
  * **Deterministic `EvidenceV1` Schema Before AI (`EvidencePanel.tsx`)**: Pre-extracted hard metrics from GRID `.jsonl` files first so every insight cited exact `dataPoints`, rounds, players, and confidence levels instead of hallucinating over raw logs.
  * **Clean Match vs. Opponent Dossier Hierarchy**: Separated team macro trends from single-series post-mortems and opponent-specific views.
* **What Was NOT Useful for Pro Staff (Avoid)**:
  * **Chatbox-First "Ask the AI Coach" & Free-Text "What-If" (`InteractiveCoachPanel.tsx`, `whatIfHandler.ts`)**: Coaches under prep crunch will not type prompts into a chatbot, and LLM counterfactual guesses (*"What if we forced Round 4?"*) lack statistical validity on small pro samples.
  * **Broadcast "Infotainment" Charts (`WinProbabilityTimeline.tsx`, `PlayerRadarChart.tsx`)**: Win-probability curves tell staff *when* a game swung (which they already know) without explaining *what tactical lever* caused it.

### Winner 2: Spector (`spector` — Category 2 Automated Scouting Report Generator)
* **What It Built**: Dual-title scouting app whose flagship polish went into **Valorant** (interactive 2D map player at `20x` speed, freeze-time/post-plant setups, lurker & utility overlays), paired with a **League of Legends** tabular report (`app/analytics/*`) covering pick/ban counts, `counterPickGoldDiff` (`@15`), `botLaneDrakePrio`, `supportGrubRecall`, `adcJoinedGrubs`, `drakeGoldHolding`, `playerWorthAt15`, `classWinRate`, and `comebackStats`.
* **What WAS Useful**:
  * **Concrete Micro-Telemetry Tells Over Generic KDA**: Drilling into GRID events to measure *unspent gold held at Dragon spawn* (`drakeGoldHolding.ts`) or *recall timing before Voidgrubs* (`supportGrubRecall.ts`) proved far more actionable than `gol.gg`-style KDA/DPM tables.
  * **Spatial Map Playback (Built for Valorant, Missing for LoL)**: Visualizing player `{x, y}` coordinates made Valorant tendencies immediately intuitive—yet Spector never built a 2D Summoner's Rift view for League despite GRID LoL streams including `{x, y}` coordinates.
* **What Was NOT Useful / Broken for 2026 LoL**:
  * **Hardcoded `positions.json` Role Lookup (`getRoleFromPlayerName.ts`)**: Mapping player nicknames statically to roles breaks on roster moves, Academy subs, and lane/role swaps.
  * **Outdated `2024–2025` Meta Heuristics (`adcJoinedGrubs.ts`)**: Tracking ADCs rotating into the Voidgrub pit targeted the `2024–2025` lane-swap meta eliminated by Season `2026` (`26.1+`) `Role Quests`.
  * **Flat Season Averages Divorced from Patch, Matchup, & Fearless State**: Averaging `counterPickGoldDiff` or `drakePrio` across all games hides *which patch*, *which lane archetype*, or *which `Bo3`/`Bo5` Fearless game (`G1` vs. `G4`)* produced the number.

### Winner 3: Cloud9 Draft Assistant (`Cloud9DraftAssistant-JamesLandry` — Category 3 Draft Predictor)
* **What It Built**: Python (`FastAPI` + `Streamlit`) single-game `20`-step snake-draft recommender scoring every legal champion in `<100ms` via a weighted linear blend of 5 signals: PyTorch behavioral-cloning `policy_score` (what pros usually pick/ban), static `meta_score` (`gol_pro_draft_data.xlsx`), `synergy_score` (`duos_map.py`), `counter_score`, and `comfort_score`/`denial_adapter` (`ROLE_DENIAL_BONUS`).
* **What WAS Useful**:
  * **Transparent Sub-Score Decomposition (`explanations.py`)**: Showing explicit `comfort`, `denial`, `synergy`, `counter`, and `meta` sub-scores avoided the opaque "54.2% win probability meter" trap used by other draft bots.
  * **Instant `<100ms` Deterministic State Transitions**: Snappy legality and role filtering made mock-draft stepping frictionless.
* **What Was NOT Useful / Fundamentally Flawed for Pro Coaching**:
  * **Step-by-Step "Pick This Champion" Autopilot**: Training a neural policy (`policy_score`) to imitate historical pro drafts regresses recommendations to stale consensus picks, ignores private scrim context, and tries to automate the head coach's core job.
  * **Single-Game `DraftState` Blind to Multi-Game Hard Fearless Strategy**: Models one `20`-action game at a time (`draft_state.py`) with zero concept of `Bo3`/`Bo5` series-level champion starvation, Game 1->5 player pool decay, or cross-game resource planning.
  * **100% Backward-Looking Static Meta (`LoLDraftingData.xlsx`)**: Completely blind on Day 1 of a new patch because it only counts *past* pro picks/bans with zero awareness of Riot patch notes or underlying champion/item/monster stat changes.
  * **Missing Comp "Shape" & Non-Stage-Legal UI**: As noted in its own `README`, it lacks composition win-condition/damage-profile reasoning, context-dependent counters, and stage-usable outputs (Streamlit cannot go on stage).

### Overlaps Removed from Our Initial Brainstorm
To guarantee strong differentiation on **Quality of the Idea**, the following earlier brainstorm items were **cut or absorbed into baseline infrastructure**:
1. ~~`R5 Counter-Pick Dossier` (v1 `#10`)~~ -> **Removed** (overlaps Spector's `counterPickGoldDiff.ts` and Draft Assistant's `counter_map.py`).
2. ~~`Draft-to-Minute-15 Voidgrub/Dragon Priority Summary` (v1 `#11`)~~ -> **Removed** (overlaps Spector's `botLaneDrakePrio.ts` and `adcJoinedGrubs.ts`).
3. ~~`Early Recall & Unspent-Gold Habit Cards` (v1 `#12`)~~ -> **Removed** (directly duplicates Spector's `supportGrubRecall.ts`, `drakeGoldHolding.ts`, and `playerWorthAt15.ts`).
4. ~~`Single-Game Phase-2 Role Pinch-Ban Scorer & Static Duo-Synergy Map` (v1 `#9`, `#15`)~~ -> **Removed** (overlaps Draft Assistant's `denial_adapter.py` `ROLE_DENIAL_BONUS`, `comfort_adapter.py`, `duos_map.py`, and Spector's `banPhaseAnalysis.ts`).
5. ~~`Standalone Claim-to-Evidence Drilldown & "What-If" Scenario Tool` (v1 `#18`, `#19`)~~ -> **Removed as standalone concepts** (StratOS already established `EvidencePanel` + `whatIfHandler`; evidence linking is treated strictly as baseline engineering hygiene).

---

## 2. Refined Non-Overlapping Brainstorm (20 Cumulative Ideas, `<50` Words Each)

Every idea below is strictly `<50` words, chains progressively from your two initial seeds, and targets **untouched white space** across all three past winners. Items tagged **`[Top Staff Pick]`** are the strongest real-world coaching force-multipliers.

### Phase I: Forward-Looking "Universe" Patch Physics (Ideas 1–6)

1. **Universe-Normalized Patch Power Baseline *(Seed Idea 1)***
   Normalize historical pro champion priority against patch "universe statistics" (base stats, per-level growth, cooldowns, ratios, items, runes, monsters) to flag champions gaining hidden systemic power before stage sample sizes exist—solving past winners' Day-1 patch blindness.

2. **High-Leverage Draft & Series Junction Map *(Seed Idea 2)***
   Instead of ranking "best next picks" at all 20 draft steps, isolate specific high-friction coaching junctions—`B1` vs. `R1/R2` trades, `B2/B3` pivots, and between-game Fearless transitions—where targeted context sharpens human judgment without automating draft calls.

3. **Patch Breakpoint & Time-to-Kill Calculator**
   Building on `#1`, compute concrete patch thresholds—waveclear minion-kill breakpoints, jungle first-clear speeds, and `Dragon`/`Baron` burn rates—so analysts see mathematically *why* item, rune, or monster stat tweaks elevate specific champions before stage games occur.

4. **Season 2026 Role-Quest & Scaling Curve Filter**
   Layer `26.x` ruleset mechanics onto `#3`—modeling Top's level-`20` cap, Mid's empowered boots/recalls, and Bot's `7`-slot item ceiling—to separate solo-queue stat outliers from champions hitting pro-viable tempo and objective spikes.

5. **Player Pool x Universe Riser Matrix `[Top Staff Pick]`**
   Cross-reference `#4`'s forward-looking patch risers against pro players' historical GRID champion pools, instantly flagging when a systemically buffed champion activates an upcoming opponent's signature comfort pick—or unlocks a sleeper pick for your roster.

6. **Cross-Region Meta Propagation Radar (`LCK` -> `LEC`/`LCS`) `[Top Staff Pick]`**
   Compare `#5`'s universe risers across the `745` GRID `LCK`, `LEC`, `LTA`, and `LCS` series to spot picks surging in Korea 1–2 weeks before Western adoption, filtered to champions matching your roster's comfort profiles.

### Phase II: Conditional Draft Trees & Hard Fearless Series Chess (Ideas 7–12)

7. **Conditional Phase-1 Opening Response Trees `[Top Staff Pick]`**
   Applying `#2` to GRID `draftActions`, replace flat pick/ban tables with conditional Blue/Red decision trees (`If B1 = X -> Opponent answers R1/R2 = Y/Z`), exposing repeatable branch habits rather than context-free season averages.

8. **Flex-Pick Bluff vs. Final-Lane Resolver**
   Extend `#7`'s opening trees by inferring final lane assignments from in-game GRID data, revealing whether an opponent genuinely flexes multi-role champions (`Rumble`, `Aurora`, `Poppy`) across players or merely bluffs during Phase 1.

9. **Hard Fearless Game 1-to-5 Performance Decay Curve `[Top Staff Pick]`**
   Unlike single-game tools, track each opponent player's GRID gold, XP, and win-rate trajectory across `Bo3`/`Bo5` games as `10–40` champions lock out, pinpointing who thrives deep in series versus who collapses off top-two comfort.

10. **Series-Spanning "Take vs. Starve" Resource Planner `[Top Staff Pick]`**
    Combine `#7` and `#9` to model `Bo3`/`Bo5` Game-1/2 tradeoffs: locking an S-tier pick for immediate win equity versus deliberately picking or handshaking champions that starve an opponent's shallowest role before Game 4.

11. **Multi-Game Comp Archetype & Damage-Split Auditor**
    Layer team composition shape (AD/AP split, reliable engage, front-to-back vs. dive, Baron shred) onto `#10`, highlighting when cumulative Fearless locks leave an opponent unable to assemble balanced damage or objective-burn compositions.

12. **Game-4/5 Forced-Pocket-Pick Forecaster `[Top Staff Pick]`**
    Merge `#5`'s universe risers with `#9`'s Fearless decay curves to flag exact late-series states where `30–40` Fearless locks plus two Phase-2 bans force an opponent player onto rusty career picks or untested patch risers.

### Phase III: Summoner's Rift Spatial Geometry & Spike Windows (Ideas 13–16)

13. **Coordinate- & Item-Driven Dynamic Role Classifier**
    Eliminate brittle static player-role files by inferring true in-game roles from GRID `Smite`/support starters and `0:00–5:00` `{x, y}` Rift coordinates—cleanly handling roster swaps, Academy subs, and off-role curveballs across `#8–#12`.

14. **Summoner's Rift `0:00–8:00` Jungle & Level-1 Pathing Map `[Top Staff Pick]`**
    Bring spatial scouting to League using GRID `{x, y}` coordinates: overlay an opponent jungler's early pathing, level-1 invade routes, and gank angles conditioned on whether adjacent lanes drafted pushing or losing matchups (`#7`).

15. **`Faelight` Vision & River Setup Formation Overlay**
    Extend `#14`'s Rift map to the `60` seconds before `Dragon`, `Rift Herald`, and `20:00` `Baron`: visualize where opponents anchor Season 2026 `Faelight` superwards, choke posture, and flank angles across winning vs. trailing game states.

16. **Item & Role-Quest Power-Spike Fight Profiler `[Top Staff Pick]`**
    Join `#3`'s item/level curves with GRID `player-purchased-item` and objective timestamps to reveal whether an opponent forces skirmishes on 1-item spikes or concedes early drakes until hitting 2-item and Season 2026 Role-Quest completions.

### Phase IV: Stage-Compliant Coach Prep & Inter-Game Workflow (Ideas 17–20)

17. **Coach-Authored Contingent Series Playbook `[Top Staff Pick]`**
    Unify `#7`'s opening trees, `#10`'s Fearless planner, and `#11`'s comp auditor into a pre-match canvas where coaches map 2–3 human-authored `Bo3`/`Bo5` branches—assisted by instant pool-cliff and damage-split warnings, never autopilot pick rankings.

18. **5-Minute Inter-Game Locker-Room Pivot Diff `[Top Staff Pick]`**
    Built for the tight huddle between `Bo3`/`Bo5` games: ingest the just-finished game's picks, update `#10`'s Fearless starvation state, flag `#16` spike-window misfires, and surface the opponent's likeliest next-game side selection and Phase-1 pivot.

19. **Patch-Shock Weekly Scrim Experiment Matrix**
    Combine `#1`'s universe stat deltas, `#6`'s `LCK` adoption lag, and `#16`'s item-spike timings into a Monday prep brief surfacing 3–5 concrete champion/build hypotheses worth testing in scrim blocks before stage lock-in.

20. **Riot-Compliant Printable Stage Binder Generator `[Top Staff Pick]`**
    Compile coach-approved branches (`#17`), Fearless pool-cliff triggers (`#12`), and Rift level-1/setup maps (`#14–#15`) into a scannable, Riot-legal paper one-pager—respecting the ban on stage electronics during live 30-second draft clocks.

---

## 3. Evaluating New Product Directions (Rubric & Coaching Fit)

Against the hackathon rubric (**Macro Scouting & Draft Management**, **GRID Data Depth**, **Design Coherence**, **Potential Impact**, and **Quality of the Idea / Novelty**), here is how the four distinct non-overlapping product concepts compare:

| Product Concept | Core Ideas Combined | Novelty vs. StratOS / Spector / Draft Assistant | Real Coaching Value (Support vs. Automation) | Technical & Data Risk | Overall Fit |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **A. Hard Fearless `Bo3`/`Bo5` Series Planner & Stage Binder** | `#7` + `#9` + `#10` + `#11` + `#12` + `#17` + `#18` + `#20` | **Very High**: Draft Assistant only models 1 isolated game; Spector uses flat series-blind averages. Nobody modeled multi-game Fearless pool decay or locker-room/stage constraints. | **Highest**: Solves the #1 workflow headache introduced in `2025–2026` pro LoL (tracking 30–50 vanishing champions across branches without spreadsheets) while keeping 100% of draft choices coach-authored. | **Low–Medium**: `745` Fearless-era series in `meta/*.csv` + `draftActions` provide rich, clean ground truth without external scraping fragility. | **#1 Anchor Wedge** |
| **B. "Universe" Patch Delta x Player Pool Radar** | `#1` + `#3` + `#4` + `#5` + `#6` + `#19` | **Highest**: Every past winner is 100% backward-looking and breaks on patch day. Bridging Riot game-file stat math (`Data Dragon`/`CommunityDragon`) with pro player pools & `LCK`->`LCS` lag is genuinely new. | **High**: Directly accelerates Monday patch-transition prep (*"Which patch changes actually matter against THIS Sunday's opponent?"*) without pretending to solve live draft. | **Medium–High**: Normalizing heterogeneous champion abilities/passives across patches can turn into a rabbit hole if scoped as a full damage simulator; keep math focused on clean stat/cooldown/item deltas + `LCK` leading indicators. | **#2 Differentiator (Best Paired w/ A)** |
| **C. Summoner's Rift 2D Spatial Macro & Setup Lab** | `#7` + `#13` + `#14` + `#15` | **High for LoL** *(Medium cross-title)*: Spector built 2D coordinate playback for Valorant, **not** LoL. Bringing matchup-conditioned jungle pathing & `Faelight` setup maps to LoL fills a glaring gap. | **High**: Analysts spend hours screenshotting VOD minimaps for Level-1 invades and river setups; instant Draft-Conditioned Rift heatmaps save massive prep time. | **Medium**: Requires streaming `.jsonl.zst` coordinate snapshots across dozens of games, so pre-indexed extraction + fast caching is essential. | **#3 High-Visual-Impact Add-On** |
| **D. Item & Role-Quest Power-Spike Window Profiler** | `#3` + `#4` + `#13` + `#16` | **High**: Neither Spector nor Draft Assistant connects drafted scaling/item curves to *when* teams choose to fight vs. concede neutral objectives in GRID timelines. | **Medium–High**: Gives macro coaches crisp tempo callouts (*"Opponent concedes 1st/2nd Drake 78% of games until 2-item spike at 18:30"*), though narrower as a standalone app. | **Low–Medium**: Cleanly computable from `player-purchased-item`, level-ups, and objective timestamps in GRID `.jsonl.zst`. | **Best as a Module Inside A or C** |

### Recommended Winning Synthesis
To maximize both **novelty over all three past winners** and **practical coaching utility** without over-scoping:
* **Primary Spine — Concept A (`Bo3`/`Bo5` Hard Fearless Junction & Series Playbook)**: Frame the product explicitly as the **anti-autopilot** alternative to `Cloud9 Draft Assistant`—built for multi-game Hard Fearless series (`Game 1 -> Game 5` pool decay, conditional Phase-1 opponent trees, comp damage/engage auditing, 5-minute inter-game locker-room diff, and printable Riot-legal stage sheets).
* **Signature Differentiator — Concept B (`Universe Patch Riser x Player Comfort Overlay`)**: Feed forward-looking patch "universe stat" deltas (`Data Dragon`/`CommunityDragon` + `LCK` adoption lag) directly into the Fearless player pools (`#5` + `#12`), answering the exact question every past historical-only tool fails on: **"When 30 champions are Fearless-banned in Game 4 on a new patch, which buffed sleeper picks can this specific player pull out?"**
* **Optional Visual Wow-Factor — Focused Slice of Concept C (`#14`)**: Attach a single clean Summoner's Rift `0:00–8:00` Level-1 / Jungle Pathing minimap card conditioned on the drafted Phase-1 lanes, delivering visual polish Spector lacked on League without rebuilding a full 40-minute replay player.
