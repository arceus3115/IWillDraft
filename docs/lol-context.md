# League of Legends, Pro Play, & Data Context

Reference documentation covering League of Legends fundamentals, professional play and Hard Fearless Draft structure, patch/meta evolution, external "universe" data sources, and the official **Sky's The Limit 2026 (Cloud9 x JetBrains)** GRID dataset.

---

## 1. What League of Legends Is (Concise Core Primer)

* **Core Objective**: A 5v5 multiplayer online battle arena (`MOBA`) played on **Summoner's Rift**. Two teams (**Blue Side** bottom-left, **Red Side** top-right) compete to destroy the enemy `Nexus` by pushing minion waves through three lanes (`Top`, `Mid`, `Bot`) defended by tiered `Turrets` and `Inhibitors`.
* **The Five Roles & Resource Allocation**:
  * **`Top`**: Isolated 1v1 side-laner (bruisers, tanks, carries) focused on wave management, split-pushing, and `Teleport` flanks. *(In Season 2026 / `26.1+`, completing the Top `Role Quest` upgrades `Teleport` and raises the level cap from `18` to `20`.)*
  * **`Jungle`**: Map-wide roamer carrying `Smite` who farms neutral monster camps (spawning at `0:55`), tracks the enemy jungler, ganks lanes, and secures epic river objectives.
  * **`Mid`**: Central-lane mage, assassin, or marksman controlling the shortest lane and fastest access to both rivers and side lanes. *(Season 2026 `Role Quest` grants upgraded boots and empowered recalls.)*
  * **`Bot` (`ADC`)**: Gold-reliant ranged sustained-DPS carry paired with a support early to scale into mid/late-game teamfights and objective burn. *(Season 2026 `Role Quest` unlocks an extra item slot — `6` items plus boots.)*
  * **`Support`**: Zero-farm facilitator (engage tanks, enchanters, catchers) driving vision control, lane trading, roaming tempo, and peel/initiation.
* **Power Scaling ("Universe" Levers)**: Across the roster (`173` champions as of [Patch 26.23](https://www.leagueoflegends.com/en-us/news/game-updates/league-of-legends-patch-26-23-notes/)), strength scales continuously via **XP/Levels** (base stats + per-level growth + ability ranks/cooldowns), **Gold -> Items** (stats + passives/actives), **Runes** (`Keystone` + minor runes + stat shards), and **Neutral Buffs**.
* **Map & Neutral Objectives (Current `26.x` Ruleset)**:
  * **Laning & Structures**: First minion wave spawns at `0:30`; outer and Tier-2 turrets carry `Turret Plates` that award localized burst gold for lane dominance.
  * **Top River Pit**: `Voidgrubs` (early structure-pushing buff) -> `Rift Herald` (summoned siege ram) -> `Baron Nashor` at `20:00` (team-wide minion/stat siege buff; `Atakhan`, `Blood Roses`, and `Feats of Strength` from 2025 were removed in `26.1`).
  * **Bot River Pit**: Stacking `Elemental Drakes` (transforming map terrain after the 2nd drake and granting permanent `Dragon Soul` at `4` drakes) -> `Elder Dragon` (late-game true-damage burn and low-HP execute).
  * **Vision Layer**: `Stealth Wards`, `Control Wards`, `Oracle Lens`, `Scryer's Bloom`, and `Faelights` (Season 2026 glowing ward rings that turn placed wards into temporary wide-area superwards).

---

## 2. Basics of Professional Play & Drafting

* **Competitive Circuit**: Regional Tier-1 leagues (`LCK` Korea, `LPL` China, `LEC` EMEA, `LCP` Asia-Pacific, and Americas `LCS`/`LTA`/`CBLOL`) play across three annual splits feeding into three international tournaments: **`First Stand`** (March), **`MSI`** (Mid-Season Invitational), and **`Worlds`** (World Championship).
  * *Patch Context*: Pro play runs on fixed patches (often 1 patch behind live). [Patch 26.23](https://www.leagueoflegends.com/en-us/news/game-updates/league-of-legends-patch-26-23-notes/) (Nov 2026) is the post-`Worlds` offseason patch, where Riot typically buffs "pro-jailed" champions (`Azir`, `K'Sante`, `Gnar`, `Galio`) that had to be kept weak during the pro calendar.
* **How Pro Play Differs from Solo Queue**:
  * **Zero-Variance Coordination**: Voice comms enable synchronized recall/`Homeguard` timers, 3-man/4-man turret dives, layered crowd control (`CC`), and cross-map trades.
  * **Tempo & Lane Priority ("Prio")**: Whoever clears the minion wave first moves first to river vision, jungle invades, or objectives without losing gold/XP. *(Season 2026 `Role Quests` heavily penalize early cross-lane XP/gold sharing, effectively eliminating the 2024–2025 pro `lane swap` meta and restoring pure 1v1/2v2 lane matchups.)*
* **Standard Tournament Draft Order (Snake Format)**:
  * **Ban Phase 1 (6 bans)**: `B1` - `R1` - `B2` - `R2` - `B3` - `R3` (removes meta-warping picks and opponent comfort champions).
  * **Pick Phase 1 (6 picks)**: `B1` -> `R1, R2` -> `B2, B3` -> `R3`
    * *Blue Side Advantage*: Gets `B1` (uncontested #1 power pick on the patch).
    * *Red Side Advantage*: Gets `R1+R2` strong duo pairing and `R3` mid-draft counter.
  * **Ban Phase 2 (4 bans)**: `R4` - `B4` - `R5` - `B5` ("pinch bans" targeting roles the enemy hasn't drafted yet).
  * **Pick Phase 2 (4 picks)**: `R4` -> `B4, B5` -> `R5`
    * *`R5` Counter-Pick*: Red side always gets the final pick of the draft (`R5`) with full information, usually reserved for a hard `Top` or `Mid` lane counter.
* **Key Draft Concepts**:
  * **Presence**: `(Picks + Bans) / Total Games` — the true measure of pro meta priority (a `95%`-banned champion has low pick sample size despite warping every draft).
  * **Flex Picks**: Champions playable in 2–3 roles (e.g., `Rumble`, `Jayce`, `Aurora`, `Poppy`, `Taliyah`, `Maokai`) used in Phase 1 to hide lane matchups and bait unfavorable counters.
  * **Composition Archetypes**: `Front-to-Back Teamfight` (tank + scaling `ADC`/control mage), `Dive/Wombo` (`Vi`/`Nocturne`/`Orianna`/`Rell`), `Poke & Siege` (`Jayce`/`Ezreal`/`Zoe`), `Pick/Skirmish`, and `1-3-1` / `4-1` `Split-Push`.
* **[Hard Fearless Draft](https://lolesports.com/en-US/news/fearless-draft-takes-over-2025) (Global Pro Standard Since 2025)**:
  * **Core Rule**: In all `Bo3` and `Bo5` series, any champion **picked** in a game is **banned for BOTH teams** for the remainder of the series ("Hard Fearless"). Standard per-game bans do *not* carry over—only played champions.
  * **Cumulative Champion Pool Shrinkage**:
    * `Game 1`: `0` Fearless locks + `10` bans = **`10` unavailable**
    * `Game 2`: `10` Fearless locks + `10` bans = **`20` unavailable**
    * `Game 3`: `20` Fearless locks + `10` bans = **`30` unavailable**
    * `Game 4`: `30` Fearless locks + `10` bans = **`40` unavailable**
    * `Game 5`: `40` Fearless locks + `10` bans = **`50` unavailable**
  * **Strategic Impact**: Eliminates repetitive "handshake" matchups (`Azir` vs. `Corki` 5 games straight), punishes shallow player champion pools, introduces **series-level resource management** (burning S-tier combos early vs. denying an opponent's crutch pick in Game 1), and makes Phase-2 role-pinch bans devastating in Games 4–5.

---

## 3. How the Meta Evolves Over Patches & External Resources

### Why & How the Meta Changes
* **Bi-Weekly Patch Cadence (`YY.XX`)**: Riot ships balance patches roughly every 2 weeks (`25.1–25.24` in 2025, `26.1–26.24` in 2026). `.1` patches introduce systemic map/item/objective overhauls, mid-season/pre-`MSI`/pre-`Worlds` patches reshape competitive priority, and minor patches nudge outliers.
* **Direct vs. Indirect ("Systemic") Meta Shifts**:
  * *Direct*: Base stat (`+3 Base AD`, `+2 Armor`), ability ratio, cooldown, or mana tweaks that flip early wave-clear breakpoints or level-1–6 lane trades.
  * *Indirect/Systemic*: Item reworks, rune tuning, minion spawn/speed changes, turret plate durability, or epic monster HP/damage changes. For example, tankier `Elemental Drakes` and `20:00 Baron` in `26.1` boosted sustained-DPS `ADC`s and objective-shred mages (`Ryze`, `Azir`, `Cassiopeia`) without touching their kits.
* **Major Historical Pro Meta Eras**:
  * *Seasons 3–7 (Vision & Slow Macro)*: Cheap wards and `Tracker's Knife` enabled slow, low-kill vision-choke control (`40+` min games).
  * *Seasons 8–11 (Skirmish & Turret Plate Revolution)*: Removal of `Tracker's Knife` + addition of `Turret Plates` (`2018–2019`) + `Dragon Soul` (`2020`) forced mandatory `5v5` fights every 5 minutes and rewarded aggressive early solo-laners (`IG`/`FPX`/`DWG` era).
  * *Season 12 (`Patch 12.10` Durability Update)*: Universal base `HP`/`Armor`/`MR` buffs across all champions slowed burst lethality and revived `Enchanter` + hypercarry front-to-back metas.
  * *2024 (Voidgrubs & Lane Swaps)*: `Voidgrubs` (`14.1`) sparked asymmetrical bot/top `lane swaps` to dodge losing `2v2` bot lanes.
  * *2025–2026 (Hard Fearless & High-Tempo Role Quests)*: `Hard Fearless` expanded champion diversity across `Bo3`/`Bo5`s, while `26.1` removed `Atakhan`/lane swaps via `Role Quests` and accelerated `Homeguards` + minion timings into a fast-tempo, lane-priority meta.

### Curated External Resources & Data Sources
* **Pro Esports Match, Draft, & Scouting Data**
  * **GRID Esports (`grid.gg`)**: Official Riot Esports server-level data partner covering full pick/ban sequences, timestamps, live gold/XP/item states, coordinates, objective control, and turret/kill events (see Section 4 for the Hackathon starter dataset).
  * **[Leaguepedia Cargo API](https://lol.fandom.com/wiki/Help:Leaguepedia_API)**: Free MediaWiki `Cargo` SQL-like query API (`ScoreboardGames`, `ScoreboardPlayers`, `PicksAndBansS7`, `Tournaments`) with deep historical pro & Tier-2 draft order, side selection, rosters, and exact patch strings.
  * **Oracle's Elixir (`oracleselixir.com`)**: Downloadable annual `CSV` datasets of global pro matches (`GD@15`, `CSD@15`, `XPD@15`, `DPM`, vision stats, picks/bans, side, patch) from `2014–2026`, plus pro meta commentary.
  * **Games of Legends (`gol.gg`)**: Visual scouting reference for per-patch/per-tournament champion presence, Blue vs. Red win rates, blind-pick vs. counter-pick win rates, and player champion pools.
* **Historical Patch Notes & Structured Game/Stat Files ("Universe Data")**
  * **Official League of Legends Wiki (`wiki.leagueoflegends.com`)**: Official Riot-hosted wiki with complete `Patch History` changelogs and exact math/scaling formulas for every champion, item, rune, and neutral monster (`2009–26.23`).
  * **Riot Data Dragon (`ddragon.leagueoflegends.com`)**: Official static `JSON` archives indexed by `versions.json` (`champion.json`, `<Champion>.json`, `item.json`, `runesReforged.json`) for base stats and per-level stat growth across historical patches.
  * **CommunityDragon (`raw.communitydragon.org`) & Meraki Analytics (`github.com/meraki-analytics/lolstaticdata`)**: Extracted raw client bin files (`CommunityDragon`) and normalized `JSON` (`lolstaticdata`) capturing exact ability damage numbers, AP/AD/HP ratios, cooldowns per rank, and monster stats where `Data Dragon` leaves unresolved placeholders.
  * **NerfPlz Patch History & Pro Meta (`nerfplz.com/lol-patch-history`, `nerfplz.com/lol-pro-meta`)**: Aggregates `300+` patches of numeric champion/item stat diffs alongside pro vs. solo-queue presence comparisons.
* **Meta Evolution Analysis & Leading Indicators**
  * **High-Elo Statistical Trackers (`lolalytics.com`, `dpm.lol`, `u.gg`)**: `LoLalytics` (`Master+` patch-over-patch delta tabs) and `dpm.lol` (pro player solo-queue tracking and matchup VODs) serve as 1-to-2-week leading indicators before picks appear on stage.
  * **Analytical Breakdowns**: Riot's official [Patch Notes](https://www.leagueoflegends.com/en-us/news/game-updates/league-of-legends-patch-26-23-notes/) developer commentary, **LoL Esports x Skill-Capped** seasonal meta guides, and **Sheep Esports** / coach draft reviews.

---

## 4. Hackathon GRID Data Starter (`stl-hackathon-2026-grid-starter`)

Primary match and telemetry dataset for **Sky's The Limit 2026: Cloud9 x JetBrains Hackathon**:
* **Repository**: [`https://github.com/cloud9esports/stl-hackathon-2026-grid-starter`](https://github.com/cloud9esports/stl-hackathon-2026-grid-starter)

### League of Legends Coverage (`2025–2026`)
* **Volume**: **`745` series (`1,739` games)** across `20` tournament stages (`2025-01-15` through `2026-10-12`), covering Hard Fearless `Bo3`/`Bo5` eras and both `25.x` and `26.x` seasonal rulesets:
  * **2025 (`538` series)**: `LCK` (LCK Cup, Split 2, Split 3), `LEC` (Winter, Spring, Summer), `LTA North` (Splits 1–3), `LTA South` (Splits 1–3), and `LTA Cross-Conference` (Split 1 & Regional Championship).
  * **2026 (`207` series)**: `LEC` (LEC Versus 2026, Spring 2026) and `LCS` (Split 1, Split 2, Split 3, Promotion 2026).
  * **Cloud9 Identifier**: Cloud9's LCS roster appears under **`Cloud9 Kia`** (always resolve across splits via `team_id` in `meta/teams.csv`).

### Repository Structure & Lookup Tables ("Index First")
Query the plain-CSV index files in `lol/` and `meta/` first—most draft, roster, patch, side-selection, and head-to-head queries can be answered from `meta/*.csv` alone without opening `.zst` feeds:

* **`lol/series.csv`**: One row per series (`series_id`, `event`, `stage`, `scheduled_start`, `teams`, `file`, `ability_supplement`).
* **`lol/events_{seriesId}_grid.jsonl.zst`**: Complete Zstandard-compressed UTF-8 JSON Lines event stream for one series.
* **`meta/*.csv`** *(contains both `val` and `lol`; always filter `title == "lol"` and join on string `(title, id)`)*:
  * `tournaments.csv`: `tournament_id`, `tournament_name`, `event`, `event_tournament_id`, `stage`, `series_count`, `first_scheduled`, `last_scheduled`.
  * `teams.csv`: `team_id`, `name` (latest), `other_names` (`; `-separated historical broadcast names), `series_count`, `events`, `first_seen`, `last_seen`.
  * `players.csv`: `player_id`, `name` (in-game nickname), `other_names`, `latest_team_id`, `latest_team_name`, `team_ids`, `series_count`, `first_seen`, `last_seen`.
  * `characters.csv`: `character_id` (champion ID used in feeds), `name`, `games_played`.
  * `games.csv`: `series_id`, `game_id`, `game_number`, `map`, `patch`, `started_at`, `finished`, `winner_team_id`, `team{1,2}_id/name/side/score`.
  * `rosters.csv`: `series_id`, `game_id`, `game_number`, `team_id`, `team_name`, `side`, `player_id`, `player_name`, `character_id`, `character_name`, `participation_status`.

### GRID Event Stream Format & Key LoL Events
* **Envelope Structure**: One JSON object per line ordered by strictly increasing `sequenceNumber`:
  * Envelope keys: `id`, `correlationId`, `occurredAt` (UTC millisecond timestamp), `seriesId`, `sequenceNumber`, `events[]`.
  * Event keys: `id`, `type` (`"{actor.type}-{action}-{target.type}"`), `actor`, `action`, `target`, `seriesState` (full series/game/team/player snapshot *after* the event), `seriesStateDelta` (incremental diff), `includesFullState`.
* **Key League of Legends Event Types**:
  * **Draft**: `team-picked-character`, `team-banned-character` (ordered pick/ban history is also recorded in `seriesState.games[].draftActions`).
  * **Economy, XP, & Combat**: `player-purchased-item`, `player-used-ability`, `player-completed-increaseLevel`, `player-killed-player`, `player-multikilled-player`, `game-respawned-player`.
  * **Objectives & Structures**: `player-completed-slayVoidGrub`, `player-completed-slayRiftHerald`, `player-completed-slayDragon` (plus elemental drake subtypes), `player-completed-slayBaron`, `player-destroyed-tower`, `player-destroyed-fortifier`, `player-destroyed-base`.
  * **Game & Series Lifecycle**: `series-started-game`, `team-won-game`, `series-ended-game`, `team-won-series`.
* **State & Spatial Coordinates**: `seriesState` tracks team-level stats (`side`, `score`, `won`, `kills`, `money`) and player-level stats (`character`, `kills`, `deaths`, `inventory`, `abilities`, `alive`, `currentHealth`, and `{x, y}` Summoner's Rift coordinates bounded by `game.map.bounds`). The **last event in a file carrying `seriesState`** holds the final box score.

### Data Caveats & Strict Licensing / Repo Hygiene Rules
* **Never Commit Raw or Derived GRID Data**: Pushing `.zst` feeds, decompressed `.jsonl` files, copies of `meta/` or `series.csv`, or derived extracts/databases (`.parquet`, `.sqlite`, `.duckdb`, caches) to **any** git repository (even private) constitutes prohibited redistribution under the hackathon license.
* **Sibling Checkout Pattern**: Clone the data starter **next to** this repository rather than inside it (e.g., `../stl-hackathon-2026-grid-starter`, referenced via an environment variable such as `STL_DATA_DIR`), optionally using sparse checkout for League of Legends only:
  ```bash
  git clone --filter=blob:none --sparse https://github.com/cloud9esports/stl-hackathon-2026-grid-starter.git
  cd stl-hackathon-2026-grid-starter
  git sparse-checkout set lol meta
  ```
* **Stream Rather Than Decompress in Bulk**: Uncompressed `.jsonl` feeds are `50–150x` larger than `.jsonl.zst`. Stream `.zst` files line-by-line instead of loading full files or decompressing the entire corpus to disk.
* **Handle Rollbacks & Feed Gaps Safely**:
  * Actions prefixed with `!` undo earlier events (e.g., `player-!purchased-item` for undone shop purchases, or `series-!started-game` when an aborted game was remade under a new `game_id`).
  * A small number of series have feed gaps or omit `team-won-series`; always derive game outcomes from final per-game state (`meta/games.csv` or the final `seriesState`) and join entities on string IDs rather than broadcast team/player names.
