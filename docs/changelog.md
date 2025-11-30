# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

*This changelog was automatically generated on 2025-11-27 22:47:18*

## 2025-11-27 (Thursday)

### ✨ Features

- fix: remove redundant newline in fmt.Println to pass linting ([`ebeeee28`](../../commit/ebeeee2896733d1c8ff3da684e052b451f195bab)) by Kiran Shrestha (1 files, +2/-2 lines)
- feat: add TOON (Token-Oriented Object Notation) for LLM context optimization ([`a91eb5f4`](../../commit/a91eb5f414c56f9e0d5b92a5f189f94e374bc452)) by Kiran Shrestha (4 files, +1255/-0 lines)
  New optional feature for efficient LLM interaction and context management:
- feat: add Docker/OCI support with GitHub Actions CI/CD ([`f630fdb0`](../../commit/f630fdb0b0859e3d7aeaac3840ea615d1594538c)) by Kiran Shrestha (5 files, +625/-0 lines)
  - Add Dockerfile with multi-stage build (Go builder + Alpine runtime)
- feat: add server.json with binary distribution for v0.1.0 ([`de9282b5`](../../commit/de9282b5e136a8827154c3ddee84d6b932e68492)) by Kiran Shrestha (1 files, +71/-0 lines)

### 🐛 Bug Fixes

- fix: clean up TOON package formatting and remove duplicate package declaration ([`f4f41db9`](../../commit/f4f41db9fffc0d43c67153b80d68480905f7c3d3)) by Kiran Shrestha (2 files, +259/-266 lines)
- fix: update go-ci.yml to use Go 1.25 to match go.mod requirement ([`8b97b173`](../../commit/8b97b173d7a8b7b5521f03e6ef9f901b28603541)) by Kiran Shrestha (1 files, +1/-1 lines)
- fix: update Dockerfile Go version to 1.25 ([`bd39bb10`](../../commit/bd39bb1023fe6886d712a6b2900b5537f108940a)) by Kiran Shrestha (1 files, +1/-1 lines)
  go.mod requires Go >= 1.25.0, but Dockerfile was using 1.21-alpine.
- fix: use OCI (Docker) registry type with correct transport ([`4367eddf`](../../commit/4367eddfcf5d9910127d7515ae8acade5457a8ed)) by Kiran Shrestha (1 files, +3/-4 lines)
- fix: use http registry type for remote server ([`bfcf1030`](../../commit/bfcf103011e866585d659e39bf4082bbd6434610)) by Kiran Shrestha (1 files, +1/-3 lines)
- fix: use github-release registry type ([`2c94c727`](../../commit/2c94c72771a11a764c2b4a10e2a279283e6a18eb)) by Kiran Shrestha (1 files, +3/-2 lines)
- fix: simplify to docker registry type for initial publication ([`7f3bdd60`](../../commit/7f3bdd605ac766afcaa258dca0ea80a040330ce4)) by Kiran Shrestha (1 files, +1/-45 lines)
- fix: use universal registry type for binary distribution ([`41aacac9`](../../commit/41aacac98c381ecf747ac4b71c56d6444e5170f6)) by Kiran Shrestha (1 files, +1/-2 lines)

### 📚 Documentation

- docs: add comprehensive OCI publishing guide ([`2fb46b6f`](../../commit/2fb46b6fb230e1d1c9c7a9f547f887e0b54ce839)) by Kiran Shrestha (1 files, +402/-0 lines)
  Complete guide for publishing AionMCP to GitHub Container Registry:
- docs: add visibility strategy, roadmap, comparison, and contributing guide ([`3bce366d`](../../commit/3bce366df5d8caf8869273f5f5a18e2f8f95d4bc)) by Kiran Shrestha (7 files, +1915/-0 lines)
  - Add VISIBILITY_STRATEGY.md for community growth planning

### 🔧 Chores

- chore: remove packages - listing as metadata-only entry ([`5e873b19`](../../commit/5e873b19921bf8c8775e4b97be5824c33f1ffcbc)) by Kiran Shrestha (1 files, +3/-11 lines)

## 2025-11-26 (Wednesday)

### ✨ Features

- feat: make AionMCP model-independent with universal REST API ([`d1f85839`](../../commit/d1f85839888e0cd9778d71983aab784b1333a427)) by Kiran Shrestha (37 files, +3583/-1914 lines)
  - Remove Claude Desktop specific configuration (no longer model-dependent)

### 📚 Documentation

- docs: update root README with model-independent focus ([`d3cfa422`](../../commit/d3cfa4223a62de44f6ebaea60a8d7709411bc390)) by Kiran Shrestha (1 files, +223/-125 lines)
  - Emphasize universal REST API architecture

## 2025-11-19 (Wednesday)

### 🐛 Bug Fixes

- fix: Improve BoltDB recovery from corrupted/locked state ([`3f709e4a`](../../commit/3f709e4a2519aeb1bf1b6b34676806ec63c5be46)) by Kiran Shrestha (1 files, +29/-7 lines)
  - Reduce timeout from 5s to 2s for faster failure detection
- fix: Fix BoltDB initialization and server startup ([`69d1d840`](../../commit/69d1d840e5c412d15de9865e09ac9e1bf8352752)) by Kiran Shrestha (2720 files, +313784/-372206 lines)
  - Move data directory creation before logger initialization

### 📚 Documentation

- docs: Add comprehensive testing guides and status documentation ([`d159c9bb`](../../commit/d159c9bb26cafcf9277da2fe63a8934f14630b24)) by Kiran Shrestha (2 files, +755/-0 lines)
  - Add TESTING_STATUS.md: Current implementation status, quick start, and troubleshooting

## 2025-11-09 (Sunday)

### 🐛 Bug Fixes

- fix: Resolve merge conflict and add missing imports ([`78b79c8d`](../../commit/78b79c8d841708faa0b8015c84751160ed894200)) by copilot-swe-agent[bot] (11 files, +152/-91 lines)
  - Fixed duplicate setupHTTPRoutes call from merge conflict

### 📦 Other

- Initial plan ([`7a263d0a`](../../commit/7a263d0a020055f3f5b82c96b95497beca602a35)) by copilot-swe-agent[bot]
- Code quality verification and conflict resolution for Iteration 4 ([`4610de60`](../../commit/4610de60e68aa64b60062c9c810ccbdf2ce17dc9)) by copilot-swe-agent[bot] (30 files, +970/-829 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- Initial plan ([`2281c156`](../../commit/2281c156acc1b22062c59250e21399ac81ffe8e4)) by copilot-swe-agent[bot]

## 2025-11-07 (Friday)

### 🐛 Bug Fixes

- fix: Correct semaphore release logic with acquisition tracking ([`6a8bcb57`](../../commit/6a8bcb57df04f9b4e2c67d69c2ac723bb2a080a4)) by copilot-swe-agent[bot] (7 files, +51/-39 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Address PR review feedback - improve concurrency safety and test reliability ([`c6d73fec`](../../commit/c6d73fec2483f20bcebc9d5fd305b13e38eb9f24)) by copilot-swe-agent[bot] (8 files, +193/-114 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### 📦 Other

- Update internal/core/registry.go ([`00f652e7`](../../commit/00f652e7482387f43e4ae13d4070c8769df8c2d3)) by Kiran Shrestha (1 files, +9/-5 lines)
  Co-authored-by: Copilot <175728472+Copilot@users.noreply.github.com>

## 2025-11-06 (Thursday)

### 🐛 Bug Fixes

- fix: Apply PR review feedback - implement handler removal and fix concurrency issues ([`85d4dd58`](../../commit/85d4dd5894bced7e9e7694f6101a9aed1d29da9a)) by copilot-swe-agent[bot] (10 files, +325/-92 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Address code review feedback - error classification, metadata access, and comment accuracy ([`1580003c`](../../commit/1580003c8972e02959df07ad17caf6f2c4c3454e)) by copilot-swe-agent[bot] (3 files, +8/-6 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### 📦 Other

- Initial plan ([`0f89f83c`](../../commit/0f89f83c3f120cf09cfb4976996649bf2cb5604d)) by copilot-swe-agent[bot]
- Initial plan ([`4ab45194`](../../commit/4ab451948d7107e984e7420e30460698315f28d4)) by copilot-swe-agent[bot]

## 2025-11-02 (Sunday)

### ✨ Features

- refactor: Improve JSON parsing robustness and implement schema parsing ([`e019252f`](../../commit/e019252fa67325712a370c02f24b0f2c9a6d1a89)) by copilot-swe-agent[bot] (7 files, +93/-109 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### 🐛 Bug Fixes

- fix: Address code review feedback - error message capitalization, comment accuracy, and goroutine parameter passing ([`a8da604f`](../../commit/a8da604f5a084412f65bafa11e943c3e18c7bcf0)) by copilot-swe-agent[bot] (5 files, +9/-9 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Address all code review feedback from PR #1 review 3408186579 ([`b94355ac`](../../commit/b94355acfbdc528b68d9183df8a095c32e1e1732)) by copilot-swe-agent[bot] (5 files, +68/-35 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Address all code review feedback from PR #1 ([`fd76bcbc`](../../commit/fd76bcbcb33db5f68ff02203ad13fdbd21ca28fb)) by copilot-swe-agent[bot] (7 files, +59/-37 lines)
  - Fix BoltDB key copying issue in cleanup function (keys must be copied during cursor iteration)
- fix: Address code review feedback - error handling, race conditions, and dependency updates ([`e768c31f`](../../commit/e768c31fb0edac31b8a8ff0f54787c762317ba44)) by copilot-swe-agent[bot] (7 files, +28/-18 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Apply PR review feedback - fix JSON parsing and serialization issues ([`ef4d6c87`](../../commit/ef4d6c871b1b129e374810f6faf3b2063cfdf208)) by copilot-swe-agent[bot] (3 files, +46/-21 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### 📚 Documentation

- refactor: Consolidate constants and add comprehensive documentation ([`8bf1e1a7`](../../commit/8bf1e1a7f3d0794dacb1e9129777d96c3ad814e2)) by copilot-swe-agent[bot] (11 files, +110/-72 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Apply PR review feedback - improve API design and documentation ([`15fe8862`](../../commit/15fe88627bc54a85f2440eea8068aa9c22fa6719)) by copilot-swe-agent[bot] (9 files, +213/-74 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### ♻️ Code Refactoring

- refactor: Extract magic numbers to constants and simplify code structure ([`77248113`](../../commit/772481134d3b465bc56a0547288e8dfd5cc6d611)) by copilot-swe-agent[bot] (10 files, +86/-93 lines)
  - Add health score deduction constants to utils.go
- refactor: Extract duplicated logic and add missing constants ([`cd486cfe`](../../commit/cd486cfe4f821fe30f18bfba9d3c5aacf23ace34)) by copilot-swe-agent[bot] (12 files, +220/-166 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>

### 📦 Other

- Initial plan ([`42110703`](../../commit/42110703db96b228d6fadbaf2ccc633768e1c849)) by copilot-swe-agent[bot]
- Initial plan ([`45f8d9f7`](../../commit/45f8d9f78df25650e24e9e0346c404b5c8302e41)) by copilot-swe-agent[bot]
- Initial plan ([`fb2fb11f`](../../commit/fb2fb11f4e57925fcfcaaf5ffde136bf48882916)) by copilot-swe-agent[bot]
- Initial plan ([`56455e47`](../../commit/56455e47a982d174ab712e45c7ad42c462f406a6)) by copilot-swe-agent[bot]
- Initial plan ([`ed72de97`](../../commit/ed72de972116228d3cb6e1c3a258e2d61e087f62)) by copilot-swe-agent[bot]
- Initial plan ([`62b6699c`](../../commit/62b6699c0a64df48fd1f518c76747c71d2d80187)) by copilot-swe-agent[bot]
- Initial plan ([`e8a2a4bb`](../../commit/e8a2a4bb652a2477bac85a073e25e16821385931)) by copilot-swe-agent[bot]
- Initial plan ([`1c2292ae`](../../commit/1c2292ae07c091431e7844fc62e8e296c5512821)) by copilot-swe-agent[bot]
- Initial plan ([`04bdfaed`](../../commit/04bdfaedf1c941e48f95d3b53595af1784b7060a)) by copilot-swe-agent[bot]

## 2025-10-29 (Wednesday)

### ✨ Features

- feat: Complete Iteration 5 - VS Code Extension with Agent Management ([`474df328`](../../commit/474df32885a90ecd4752a42e766ed7d867cf8176)) by Kiran Shrestha (8 files, +503/-49 lines)
  - Implemented VS Code extension with Tree Views for agents and tools

### 🐛 Bug Fixes

- fix: Apply PR review feedback from code review comments ([`738fb1ce`](../../commit/738fb1cea5496acfc136808ebee904687b091129)) by copilot-swe-agent[bot] (11 files, +181/-64 lines)
  - Remove redundant zero-check logic in server.go (lines 527-529)
- fix: Address code review feedback for self-learning engine ([`e9e8c3ff`](../../commit/e9e8c3ff9d35559946c9e17cf3f9ca149a9adedf)) by copilot-swe-agent[bot] (7 files, +29/-9 lines)
  - Fix nil pointer dereference in collector.go filterPII method
- fix: Address code review comments - constants, date format, regex patterns ([`8d53b468`](../../commit/8d53b468520415da7ea89bfce4cc2040f1854278)) by copilot-swe-agent[bot] (11 files, +211/-70 lines)
  Co-authored-by: kiransth77 <23469105+kiransth77@users.noreply.github.com>
- fix: Address code review feedback for self-learning engine ([`a243b634`](../../commit/a243b634f93f1f9a987d0d85649a7fca357ef4fa)) by copilot-swe-agent[bot] (5 files, +17/-7 lines)
  - Fix variable shadowing bug in GetToolInsights (engine.go:110)

### ♻️ Code Refactoring

- refactor: Address PR review feedback - improve code maintainability ([`da73c6f7`](../../commit/da73c6f7b3166a9e29faf3f55422a50d58f3a762)) by copilot-swe-agent[bot] (2 files, +30/-20 lines)
  - Replace custom path parsing with filepath.Dir() for cross-platform compatibility
- refactor: Extract duplicated methods and add configuration options ([`7588b7cc`](../../commit/7588b7ccfaf34e42a5d8dc1f59c78632ebe34ebf)) by copilot-swe-agent[bot] (12 files, +156/-131 lines)
  - Extract getHealthStatus and writeToFile to shared utils.go

### 📦 Other

- Initial plan ([`d5fc4293`](../../commit/d5fc42934a34e6d0586b56b5b043a1719f67f482)) by copilot-swe-agent[bot]
- Initial plan ([`a73b75e8`](../../commit/a73b75e819d27ebd380d9299d79ce60f1562a2ea)) by copilot-swe-agent[bot]
- Initial plan ([`e0b74312`](../../commit/e0b74312ce1e58e306b46752af4e632bd7154479)) by copilot-swe-agent[bot]
- Initial plan ([`f9191460`](../../commit/f9191460ce1dbd298c49495de9c2a59b47aad58c)) by copilot-swe-agent[bot]
- Initial plan ([`ea588b75`](../../commit/ea588b75afdad2cba4fd5a0ce9cde721ee25601d)) by copilot-swe-agent[bot]
- Initial plan ([`509cc649`](../../commit/509cc649efd7bf45abf96863010bca563e86f3b5)) by copilot-swe-agent[bot]
- update gitignore to include node artifacts ([`c2b32a93`](../../commit/c2b32a93de3e0e98547dcf43ed90b958926d85e2)) by Kiran Shrestha (1 files, +19/-1 lines)

## Summary

**Period:** 2025-10-28 to 2025-11-27

**Total commits:** 64

**Changes by type:**

- Other: 22
- Documentation: 4
- Refactoring: 5
- Features: 7
- Bug Fixes: 25
- Chores: 1

**Contributors:** 2

- copilot-swe-agent[bot]: 41 commits
- Kiran Shrestha: 23 commits

**Code changes:**
- Files changed: 2993
- Lines added: +326703
- Lines removed: -376750
- Net change: -50047 lines

