# GRID draft experiment

## Tier-one draft decision support

Supply one local Oracle’s Elixir yearly CSV per season. The importer takes domestic LCS, LCK, LPL, and LEC games and Worlds from 2025 onward, reconciles the two team rows with the ten player and role rows, and joins games into series. The JSON output shows the dates and patches covered in each competition and counts discarded games by reason. GRID remains the source for the ordered replay and gameplay analysis below. [Oracle’s Elixir yearly downloads](https://drive.google.com/drive/folders/1gLSw0RLjBbtaNy0dgnGQDAZOHIgCe-HH)

```sh
GOWORK=off go run . -pro-csv "$OE_2026_CSV" -pro-year 2026
GOWORK=off go run . -pro-csv "$OE_2026_CSV" -pro-scenario /path/to/scenario.json
```

The backtest refits before each held-out series. A prior series enters training four hours after its last observed game start. Series identity uses the source series key when embedded in the game ID and otherwise uses the competition, season, split, team pair, ordered game numbers, and timestamps. The report shows Brier loss, log loss, and calibration by competition, patch, and whether the competition had fewer than 15 preceding games on the target patch. It compares the patch-linked model with independent patch effects, team/player/side, target-competition-only, and same-patch models. Each comparison uses the games scored by both models and resamples whole series for intervals.

Evidence weight halves every two **competitive** patch steps within a client major. For each historical game, the model counts distinct intervening patches played or scheduled to have started in its competition or the target competition before the cutoff. It includes the target patch once even on a scheduled debut. Same-patch games receive the same recency weight throughout the patch. Dates determine patch assignment and which series have finished. Games on newer builds are excluded when projecting an older tournament build. A donor competition contributes at most half the greater of the target competition's weight and 12. `-pro-patch-resets 26.13,26.19` illustrates analyst-specified breaks for changes that make earlier observations unsuitable.

### Competitive patch calendar

Oracle's Elixir supplies a game patch where available. The GRID archive also has a game-level patch in `meta/games.csv`. The tier-one importer keeps the recorded patch when present and compiles observed windows per league from runs of games with the same tag. For a feed with missing tags or no `patch` column, pass one dated range per league and competitive patch.

```sh
GOWORK=off go run . -pro-csv "$OE_2026_CSV" -pro-patch-calendar "$PATCH_CALENDAR_CSV"
```

```csv
league,patch,start,end
LCK,26.6,2026-04-01T08:00:00Z,2026-04-08T08:00:00Z
LCK,26.7,2026-04-08T08:00:00Z,2026-04-22T08:00:00Z
```

The dates above illustrate the format. Use verified competitive dates with explicit UTC offsets. Starts are inclusive, ends are exclusive, and windows may overlap across different leagues. Patch labels such as `26.6` and `16.6` resolve to the same version. A recorded game that contradicts a scheduled window is counted in coverage and disables automatic assignment from that window. An untagged game can also inherit a patch when bracketed by matching recorded games in its own league. In a historical backtest, that game enters training only after the patch window is corroborated and its series is available. Games in uncovered gaps are reported as rejected. JSON includes the resolved windows and counts by patch source.

Use competition `WORLDS` for the tournament. Oracle’s Elixir labels those rows `WLDs`. The importer uses one locked window for each event, with Riot's published tournament dates and patches built in for 2025 and 2026. Supply a calendar window for other seasons. The importer attributes regional qualifiers with verified LPL match URLs to LPL and reports other out-of-window rows and patch conflicts. Across competitions, team and player terms share stable source IDs so a Worlds entrant can carry domestic history into the event. Domestic results alone give limited information about relative team strength between regions. [2025 schedule](https://lolesports.com/en-GB/news/worlds-2025-heads-to-china) · [2025 patch notes](https://www.leagueoflegends.com/en-us/news/game-updates/patch-25-20-notes/) · [2026 schedule](https://lolesports.com/en-GB/news/worlds-2026-primer) · [2026 patch notes](https://www.leagueoflegends.com/en-us/news/game-updates/league-of-legends-patch-26-20-notes/)

The logistic model jointly estimates map side, first-pick side, and team and player strength alongside champion-by-role effects. The patch-linked model penalizes changes between adjacent champion-role patch effects. An unseen patch inherits the learned adjustment from the most recent supported patch, and incoming stage games can move that estimate. Sparse chains allow larger changes when several competitive patch steps separate the observations. The independent-patch comparison tests whether this structure improves held-out predictions. Player-champion, allied-pair, and lane-matchup terms remain regularized. These are fixed experiment settings and the patch effects describe observed outcomes under the fitted model.

### Compare a draft choice

Set nine champion assignments and leave the focus role's champion empty. Use team and player IDs from the CSV, a timestamp with a UTC offset, and the current draft's bans, earlier Fearless locks, and any patch disables in `unavailable`. Supply a legal reference and alternatives. `as_of` controls evidence availability. Optional `match_at` identifies the scheduled match time and defaults to `as_of`. The competition's patch window supplies the version at `match_at` or verifies the explicitly supplied `patch`. Supply the competitive patch for a planned matchup outside the covered dates. Map side and first-pick side are independent under First Selection. [Riot’s First Selection rules](https://lolesports.com/en-GB/news/season-start-2026-lol-esports)

```json
{
  "league": "LCK",
  "patch": "26.9",
  "as_of": "2026-05-08T08:00:00Z",
  "first_pick_side": "red",
  "focus_side": "blue",
  "focus_role": "mid",
  "blue": {
    "team_id": "blue-team-id",
    "roles": {
      "top": {"player_id": "b-top", "champion": "Gnar"},
      "jng": {"player_id": "b-jng", "champion": "Vi"},
      "mid": {"player_id": "b-mid"},
      "bot": {"player_id": "b-bot", "champion": "Kai'Sa"},
      "sup": {"player_id": "b-sup", "champion": "Leona"}
    }
  },
  "red": {
    "team_id": "red-team-id",
    "roles": {
      "top": {"player_id": "r-top", "champion": "Renekton"},
      "jng": {"player_id": "r-jng", "champion": "Xin Zhao"},
      "mid": {"player_id": "r-mid", "champion": "Azir"},
      "bot": {"player_id": "r-bot", "champion": "Jinx"},
      "sup": {"player_id": "r-sup", "champion": "Nautilus"}
    }
  },
  "unavailable": ["Orianna"],
  "reference": "Ahri",
  "candidates": ["Ahri", "Taliyah", "Orianna"]
}
```

The decision reports probability differences versus the reference for each supported model, with whole-series refits for a conditional sampling interval on the transfer estimate. Feature contributions explain the difference in fitted log odds. Counts include stage reps for that player and role, the opposing lane matchup, the four relevant allied pairs, recorded bans, and how much champion-role evidence comes from the target patch and league. Unknown teams, players, and thin supporting roles are surfaced. A model needs 30 prior games and eight weighted games to score. Each candidate needs three role appearances and at least one weighted appearance. Fewer than five role appearances in the target league and patch marks the estimate as limited.

The comparison holds the other nine assignments fixed. Run it on different plausible completions when opponents can still change their draft. Observed stage selections and outcomes support conditional estimates. Coaching and player assignment, player proficiency, and which options teams chose can still confound the estimated differences.

## Measuring balance changes

Choose an outcome before calculating a percentage. A damage change in a specified trade, an objective clear-time change, a change in first-item purchase frequency, and a change in expected win probability answer different questions. Each needs its own comparison and denominator.

### Controlled item mechanics

`-patch-mechanics` compares two item definitions on a fixed set of fight profiles. Supply a fixed number of hits or a fixed attack window in seconds. The fixed-hit comparison isolates per-hit changes. The time-window comparison computes attack speed from the base rate, attack-speed ratio, bonuses, and cap. It assumes uninterrupted attacks and the first hit after one attack interval. Spellblade inputs specify the hit indices on which a proc is available. Stacks from previous hits modify base AD, armor reduces physical damage, and fixed other physical damage can be included in the total.

```sh
GOWORK=off go run . -patch-mechanics /path/to/trinity-force.json
```

```json
{
  "item": "Trinity Force",
  "before_patch": "13.13",
  "after_patch": "13.14",
  "before": {"flat_ad": 40, "bonus_attack_speed": 0.35, "base_ad_per_stack": 0.04, "max_stacks": 5},
  "after": {"flat_ad": 40, "bonus_attack_speed": 0.33, "base_ad_per_stack": 0.06, "max_stacks": 3},
  "profiles": [
    {"name": "short trade", "weight": 1, "base_ad": 100, "target_armor": 100, "hits": 3},
    {"name": "extended trade", "weight": 1, "base_ad": 100, "target_armor": 100, "hits": 16}
  ]
}
```

This example isolates attack damage and the Threefold Strike ramp. Under these assumptions, the same item update changes total physical damage by +1.39% over three hits and −0.32% over 16. To model attack-speed breakpoints, replace `hits` with `seconds` and provide `base_attack_speed`, `attack_speed_ratio`, and `attack_speed_cap`. To model Spellblade in a planned sequence, specify `spellblade_base_ad_ratio` in each item definition and supply `spellblade_proc_hits` in the profile. Patch notes and versioned static data supply the mechanics. Riot documents that Data Dragon versions can differ from the regional game client. [Riot’s 13.14 changes](https://www.leagueoflegends.com/en-sg/news/game-updates/patch-13-14-notes/) · [Riot’s Data Dragon documentation](https://developer.riotgames.com/docs/lol)

The tool also reports damage added above the same profile without the item. Percentage changes use `100 × (after / before − 1)`. When combining profiles, the tool applies that ratio to the weighted means using the same weights before and after. Weights can represent a fixed pre-patch mix of trade types. The example weights are illustrative. Use the actual rules on either end of a competitive patch jump. That captures changes from skipped public releases if those changes carry into the next competitive build. Build cost, purchase delay, health, movement, target access, proc uptime, and other item choices each require a corresponding outcome or scenario input.

### Estimating the stage impact

| Method | What it can estimate | Main evidence requirement |
|---|---|---|
| Linked patch effects, implemented in the draft model | Forecasts for champion-role performance that can move as pro games arrive | Series-level holdouts across patch debuts, established patches, and Worlds. Compare Brier loss and calibration with independent patch effects. |
| Standardized before/after comparison | The change in an observed outcome under the same mix of champions, roles, opponents, and game states | For item effectiveness, align telemetry at a consistent purchase or inventory landmark and compare pre-landmark covariates with adequate overlap. Keep adoption and in-game performance as separate outcomes. |
| Staggered event study | How adoption of the complete patch changed outcomes for cohorts with heavy pre-patch reliance on Trinity Force, using competitions that have yet to adopt it as comparators | Group by actual competitive adoption, use patch steps for event time, and check prior trends, concurrent changes, and within-series dependence. Patch adoption changes many mechanics at once. [Callaway and Sant’Anna](https://arxiv.org/abs/1803.09015) |
| Synthetic control | A counterfactual trajectory for an affected champion-role or build, assembled from a weighted set of comparable pre-patch trajectories | Needs a long pre-patch span, strong pre-period fit, and donors with limited exposure to the change. Run placebo patches and sensitivity analyses for substitutes and spillovers. [Abadie, Diamond, and Hainmueller](https://doi.org/10.1198/jasa.2009.ap08746) |

For an item study, define exposure from pre-patch build usage so players who drop the nerfed item still count in the affected cohort. Item purchase is influenced by gold, game state, matchup, and player decisions. The annual CSV imports draft and outcome data. Estimating an item-specific stage effect requires time-aligned inventory, gold, damage, and objective telemetry, plus credible comparisons for the simultaneous changes on that patch.

A mechanics-informed forecast can extend the implemented zero-mean patch drift. For each champion-role, multiply the item’s controlled damage change by its pre-patch build share and add separate inputs for direct champion or objective changes. Learn how those inputs predict the next patch adjustment from earlier transitions. The coefficient on item exposure translates a mechanical percentage into an estimated change in draft log odds. Hold out entire patch transitions and compare with the current patch-linked model before using that translation.

## LCS archive replay

On a completed Cloud9–FlyQuest best-of-five, the GRID Series State API and Cloud9's GRID archive agreed on all 100 ordered picks and bans and on first-dragon credit and counts for all 10 team-games. The experiment uses `grid-framework` v0.1.0 for the completed state and the archived game clock for when the lead moved. It joins both to the same Games 1–5 replay.

## Gameplay evidence for the playbook

In Game 4, Cloud9 secured first dragon at 8:18. Its gold differential moved from −234 at 8:16 to +180 at 10:16, then to −3,188 by 14:58. The Game 5 huddle can point the coach to the 10:16–14:58 stretch when reviewing whether to repeat the early-game plan. The observation identifies a review window. Attribution requires reviewing the plays in that window.

Across 189 games from the 2026 LCS splits, teams with at least a 1,000 gold lead around 15:00 won 97 of 121 games. All 189 had a usable gold sample and verified dragon counts. In 182 games, there were usable samples just before first dragon and around two minutes afterward.

| Inputs available by 15:00 | Brier loss on 165 chronological holdouts | Change from prior row, paired 95% interval |
|---|---:|---:|
| League blue-side win frequency | 0.250 | — |
| Add gold differential at 15 | 0.210 | −0.058 to −0.023 |
| Add first-dragon credit by 15 | 0.205 | −0.015 to +0.007 |

Lower Brier loss means closer probability estimates to observed results. Gold differential was the stronger signal in this cohort. The first-dragon estimate was inconclusive once gold was included. Use the objective and the nearby gold swing to give the huddle a timestamped question to investigate. There were three Game 5s in the cohort, so the pooled result offers little evidence about Game 5 specifically.

Each gold sample uses earned gold from both teams in the same event state within 15 seconds at or before its target. The two-minute swing uses the last sample before the capture and the sample at or before 120 seconds after it. The experiment reconciles event credits and counts with completed per-team objective records before scoring dragon-dependent rows. The holdouts use only prior series in the last 120 days, with at least 20 eligible training games. Lead bands split at ±1,000 and ±2,000 gold. Each estimate is shrunk toward the preceding model using fixed prior weights of 10, 8, and 8 observations. The model measures association after observing early gameplay. It supplies evidence for reviewing a played line and for scouting prior matches.

## What to put in the playbook

- **Huddle and stage sheet.** Use confirmed picks as the source of truth for legality and the strongest gameplay observation to focus the discussion. After Game 3, Twisted Fate's lock left APA with Orianna as his sole open champion with at least two appearances in his preceding 20 stage games. That is a specific contingency to recheck.
- **Fearless planner.** Show the remaining champions and their stage reps. Rank the shortlist by recency if space is tight. Recency improved top-three hit rate by about one point overall. The full last-20 window covered only 44.1% of actual Game 5 picks, which is reason to surface thin stage evidence prominently.
- **Opening tree.** Show the most relevant observed replies and evidence counts at the selected junction. Adding the actual opener improved the shortlist, but both response champions landed in the top three in just 10.6% of test drafts. Keep map side and pick order as separate fields. Red map side opened 125 of the 189 verified drafts under a ruleset that allows map side and pick order to differ. [Riot's 2026 First Selection rules](https://lolesports.com/en-GB/news/season-start-2026-lol-esports)

## Comparison

The top-three test checks whether replayed picks appeared in the three highest-ranked candidates. Rates measure how often a method recovered picks that teams actually made.

| Player model | Top-three hit rate | Difference from raw reps, paired 95% interval |
|---|---:|---:|
| Counts in the prior 20 stage games | 39.9% | — |
| Same games, weighted with a 30-day half-life | 41.0% | +0.4 to +1.7 pp |
| Same-patch counts shrunk toward last-20 frequencies | 40.1% | −0.4 to +0.7 pp |

The pool test scored 15,388 player appearances from 725 eligible series with at least five prior stage appearances per scored player. It excludes confirmed prior-game locks. The roster indexes lack current-game bans. The full last-20 pool contained 77.1% of observed Game 1 picks and 44.1% of Game 5 picks. The Game 5 cohort contained 297 scored player appearances across 30 games. Game cohorts differ because only some series reach later games.

| Opening model | At least one actual reply champion in top three | Both in top three |
|---|---:|---:|
| LCS reply frequency | 36.7% | 5.6% |
| Responding team frequency with a league prior | 36.1% | 4.4% |
| Same model with league-wide counts for this exact opener | 43.3% | 10.6% |

The opening test scored 180 of 189 verified drafts in 68 series across the three 2026 LCS splits. It required earlier responding-team history and used a 120-day window. Of the 180, 149 had an earlier match with the exact opener and three were Game 5 drafts. Current exclusions are earlier Fearless locks, the six initial bans, and the opening pick. The gain for getting both replies in the top three was 6.1 percentage points over the team model, with a paired 95% bootstrap interval of +2.6 to +9.9 pp.

Both tests froze training history before the replayed series. An earlier series qualified after four hours had elapsed from its last recorded game start. The patch model used eight games' worth of prior frequencies. The opening model used eight observations for the league prior and twelve for the team prior when adding opener counts. These were fixed experiment settings. Intervals resampled series to keep picks from the same series together. The comparisons are an exploratory backtest on this archive cohort.

## Reproduce

Use Go 1.25 or newer and `zstd`. Point `STL_DATA_DIR` at a local checkout of the official GRID starter dataset kept outside this project. The default opening cohort reads the League of Legends archives for the three 2026 LCS splits along with the roster and game indexes. [Cloud9 GRID starter data and license](https://github.com/cloud9esports/stl-hackathon-2026-grid-starter)

From this directory, run:

```sh
GOWORK=off go test -race ./...
GOWORK=off go run . -data "$STL_DATA_DIR" -replay 2953533
```

With `GRID_API_KEY` in the environment, compare the live API's historical draft and dragon records with the local archive:

```sh
GOWORK=off go run . -api-series 2953533 -data "$STL_DATA_DIR"
```
