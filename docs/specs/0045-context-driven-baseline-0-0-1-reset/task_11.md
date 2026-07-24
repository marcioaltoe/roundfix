---
task: task_11
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
type: infra
complexity: high
---

# Task 11: Align the Roundfix 0.0.1 distribution

## Overview

Align every Roundfix-owned distribution surface with the single 0.0.1 product
identity. The reset must be complete while leaving operational database
schemas, external lock schemas, upstream metadata, and historical Git evidence
outside the owned version contract.

## Requirements

1. MUST set the application semantic version, launcher manifest, platform npm
   manifests, setup generation, and all fourteen Roundfix-owned skill versions
   to exactly `0.0.1`.
2. MUST change the Release Plan JSON schema to
   `roundfix.release-plan/0.0.1` and reject obsolete owned schema identities as
   current output.
3. MUST restart `CHANGELOG.md` as the 0.0.1 history without retaining obsolete
   release sections.
4. MUST align release packaging and validation workflows with the reset owned
   identity.
5. MUST make owned-version validation enumerate the authoritative surface and
   fail on any disagreement.
6. MUST preserve Run Database schema versions, external `skills-lock.json`
   schema, third-party protocol versions, upstream skill metadata, Git history,
   and existing Run configuration.
7. MUST keep Build Commit and Build Time as the source-state distinction for
   locally built binaries.
8. MUST NOT delete or mutate tags or GitHub Releases; that remains a separately
   approved operation after a fresh Release Plan.

## Subtasks

- [x] Reset application and npm distribution versions.
- [x] Reset all fourteen owned skill metadata versions.
- [x] Reset setup generation and Release Plan schema identity.
- [x] Replace the changelog with the 0.0.1 history.
- [x] Align packaging and owned-version validation.
- [x] Add protected operational/upstream version assertions.
- [x] Verify no release-history mutation surface was introduced.

## Acceptance Criteria

- [x] One owned-version check reports `0.0.1` for every authoritative Roundfix
      distribution surface and fails if any one surface is changed.
- [x] CLI version output reports semantic version 0.0.1 while retaining Build
      Commit and Build Time fields.
- [x] Release Plan JSON uses only `roundfix.release-plan/0.0.1`.
- [x] The changelog begins at 0.0.1 and contains no obsolete release sections.
- [x] Fixtures prove Run Database, external lock, third-party, and upstream
      metadata versions are unchanged.
- [x] Tests and workflow inspection prove the implementation performs no tag
      or GitHub Release deletion.

## Context

- instruction: `docs/adr/0062-roundfix-owned-versions-restart-at-zero.md`
- instruction: `docs/adr/0065-release-plan-exposes-a-read-only-reset-mode.md`
- interface: `internal/app/version.go`
- interface: `dist/npm/roundfix/package.json`
- interface: `skills/skills.go`
- interface: `CHANGELOG.md`
- interface: `.github/workflows/release.yml`

## Verification

- `rtk go test ./internal/app ./internal/releaseplan ./skills` — expected: application, Release Plan schema, and all owned skill versions agree on 0.0.1 while protected versions remain unchanged.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_version_contract.py'` — expected: setup assets and protected upstream/operational schemas satisfy the version matrix.
- `rtk go run -buildvcs=false ./cmd/roundfix version` — expected: output reports version 0.0.1 and the existing build provenance fields.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1 and 4; Core Features 1, 14, and 15; User Story 6.
- `_techspec.md` → Data Models; Integration Points; Testing Approach; Build
  Order 6–7.
- ADR-0062 → exhaustive owned-version reset and excluded contracts.
- ADR-0065 → no release-history mutation during implementation.

## Result

Aligned the checked-in Roundfix distribution with the single `0.0.1` product
identity.

- The application, launcher package, five platform packages, launcher
  dependencies, current setup generation, Release Plan schema, and all 14
  canonical and embedded Roundfix-owned skills now agree on `0.0.1`. The
  version-contract suite enumerates 54 owned version fields and fails with the
  changed surface named when one disagrees.
- The CLI retains Build Commit and Build Time as local source-state evidence.
  `rtk ./bin/roundfix version` reported
  `roundfix 0.0.1 (27f9398-dirty, built 2026-07-23 00:45:52 -0300)` after the
  verification build.
- Normal and reset Release Plan JSON now use only
  `roundfix.release-plan/0.0.1`; the previous
  `roundfix.release-plan/v1` identity is rejected by the exact schema
  assertion.
- `CHANGELOG.md` now contains only the `0.0.1` release section. The release
  workflow requires the pushed tag to equal the checked-in launcher version
  before it can build or publish packages.
- The protected-version fixture pins Run Database schema `9`, external
  `skills-lock.json` and lock-hash compatibility schema `1`, ACP protocol `1`,
  JSON-RPC `2.0`, and ten upstream skill versions. These values remain outside
  the Roundfix-owned reset.
- Static release-surface inspection rejects tag- or GitHub Release-deletion
  commands. The implementation adds no release-history mutation method; remote
  cleanup remains outside this Task.
- `rtk env GOCACHE=/tmp/roundfix-task11-go-cache go test ./internal/app ./internal/releaseplan ./skills`
  — PASS.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_version_contract.py'`
  — PASS, 4 tests.
- `rtk env GOCACHE=/tmp/roundfix-task11-go-cache go run -buildvcs=false ./cmd/roundfix version`
  — PASS; reported `roundfix 0.0.1`.
- `rtk make skills-sync-check` — PASS; canonical and distributed skill trees
  are byte-identical.
- `rtk env GOCACHE=/tmp/roundfix-task11-go-cache make verify` — PASS; 1,725 Go
  tests, 243 canonical setup tests, 243 distributed setup tests, asset
  validation, owned-skill validation, and the stamped CLI build completed.
