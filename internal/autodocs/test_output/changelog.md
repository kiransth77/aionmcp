# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

*This changelog was automatically generated on 2025-11-27 22:47:17*

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

## Summary

**Period:** 2025-11-20 to 2025-11-27

**Total commits:** 17

**Changes by type:**

- Chores: 1
- Features: 4
- Bug Fixes: 9
- Documentation: 3

**Contributors:** 1

- Kiran Shrestha: 17 commits

**Code changes:**
- Files changed: 67
- Lines added: +8349
- Lines removed: -2376
- Net change: +5973 lines

