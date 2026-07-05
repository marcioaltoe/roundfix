---
task: task_03
spec: 0007-command-lifecycle
status: pending
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

- [ ] Release resolution and platform asset selection
- [ ] Atomic self-replace with verification and fallback message
- [ ] `--check` and the three deterministic outcomes
- [ ] Daily cached freshness note with fake-clock tests
- [ ] Seam wiring and fixture tests

## Acceptance Criteria

- [ ] Fixture matrix: newer release (binary replaced — proven against a
      temp fake executable path), current version (no-op message), empty
      releases (clean outcome), download/verify failure (binary untouched,
      exit 1).
- [ ] Freshness tests: first run checks and caches; second run within 24h
      reads cache only; behind-version prints exactly one stderr line;
      network failure prints nothing and still caches the attempt time.
- [ ] `upgrade --help` truthful; no new Go module dependencies; full suite
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
