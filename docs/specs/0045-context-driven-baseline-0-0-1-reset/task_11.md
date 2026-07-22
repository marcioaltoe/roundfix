---
task: task_11
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
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

- [ ] Reset application and npm distribution versions.
- [ ] Reset all fourteen owned skill metadata versions.
- [ ] Reset setup generation and Release Plan schema identity.
- [ ] Replace the changelog with the 0.0.1 history.
- [ ] Align packaging and owned-version validation.
- [ ] Add protected operational/upstream version assertions.
- [ ] Verify no release-history mutation surface was introduced.

## Acceptance Criteria

- [ ] One owned-version check reports `0.0.1` for every authoritative Roundfix
      distribution surface and fails if any one surface is changed.
- [ ] CLI version output reports semantic version 0.0.1 while retaining Build
      Commit and Build Time fields.
- [ ] Release Plan JSON uses only `roundfix.release-plan/0.0.1`.
- [ ] The changelog begins at 0.0.1 and contains no obsolete release sections.
- [ ] Fixtures prove Run Database, external lock, third-party, and upstream
      metadata versions are unchanged.
- [ ] Tests and workflow inspection prove the implementation performs no tag
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
