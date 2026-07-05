---
task: task_03
spec: 0007-command-lifecycle
status: completed
type: backend
complexity: high
---

# Task 03: Build the Upgrade Command and the freshness check

## Overview

`roundfix upgrade` fetches the latest release for the current platform and
replaces the running binary atomically; operational commands gain a daily,
best-effort, never-blocking stderr note when the installed version is
behind. The release channel does not exist yet — "no releases published" is
a first-class outcome by design. Verifiable through fixture releases and a
fake clock.

## Requirements

1. MUST add `roundfix upgrade [--check]`: resolve the latest release of the
   binary's release repository through the existing GitHub CLI dependency;
   outcomes `upgraded <old> → <new>`, `already current <version>`, and
   `no releases published`, each deterministic on stdout; failures exit 1
   with the bounded stderr-tail convention; `--check` reports without
   installing.
2. MUST download the platform-matching asset to a sibling temp path, verify
   it (size, plus checksum when a checksum asset is published), and rename
   over the current executable atomically; any failure leaves the current
   binary untouched and names the manual fallback.
3. MUST add the freshness check to fetch, resolve, watch, and implement:
   cached in a Roundfix Home JSON file, refreshed at most every 24 hours
   with a short timeout, one stderr line when behind naming both versions
   and `roundfix upgrade`; silent on every failure and offline; never
   observable in the Run's outcome or timing beyond the timeout.
4. MUST route clock, release lookup, and download through seams; tests
   never reach the network.

## Subtasks

- [x] Release resolution and platform asset selection
- [x] Atomic self-replace with verification and fallback message
- [x] `--check` and the three deterministic outcomes
- [x] Daily cached freshness note with fake-clock tests
- [x] Seam wiring and fixture tests

## Acceptance Criteria

- [x] Fixture matrix: newer release (binary replaced — proven against a
      temp fake executable path), current version (no-op message), empty
      releases (clean outcome), download/verify failure (binary untouched,
      exit 1).
- [x] Freshness tests: first run checks and caches; second run within 24h
      reads cache only; behind-version prints exactly one stderr line;
      network failure prints nothing and still caches the attempt time.
- [x] `upgrade --help` truthful; no new Go module dependencies; full suite
      passes.

## Verification

- `rtk go test ./internal/cli/ ./internal/app/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix upgrade --check` — expected: a clean
  deterministic outcome line (expected today: no releases published),
  exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 2, 3; Core Feature 2; Decisions. `_techspec.md` →
Upgrade Command and freshness check, Risks (self-replace), Build Order 3.
Round-1 dogfood finding 23.

## Result

- Added `roundfix upgrade [--check]` with GitHub CLI release lookup,
  platform asset selection, size/checksum verification, sibling temp-file
  replacement, deterministic stdout outcomes, and manual fallback failures.
- Added daily cached freshness checks for `fetch`, `resolve`, `watch`, and
  `implement`, backed by `.roundfix/version-check.json`, a fakeable clock,
  silent failure behavior, and a short timeout.
- Acceptance evidence: `TestRunUpgradeFixtureMatrix` covers newer/current/no
  releases/verify-failure outcomes, proving temp fake executable replacement
  and binary preservation on failure.
- Acceptance evidence: freshness tests prove first-run lookup and cache write,
  cache reuse within 24h, exactly one behind-version stderr line, silent
  network failure, and fetch wiring that does not change the command outcome.
- Acceptance evidence: `TestRunCommandHelp` covers `upgrade --help`; `go.mod`
  and `go.sum` were not changed; Roundfix skill command docs were updated and
  `skills check` passed.
- Verification: `rtk go test ./internal/cli/ ./internal/app/` passed
  (`224 passed in 2 packages`).
- Verification: `rtk go run ./cmd/roundfix upgrade --check` passed with
  stdout `no releases published`.
- Verification: `rtk go test ./...` passed (`640 passed in 16 packages`).
- Verification: `rtk make verify` passed, including full tests, skill check,
  and build.
