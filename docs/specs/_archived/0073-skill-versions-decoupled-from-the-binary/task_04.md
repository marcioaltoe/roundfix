---
task: task_04
spec: 0073-skill-versions-decoupled-from-the-binary
status: completed
type: backend
complexity: medium
---

# Task 04: Make every surface use the one comparison

## Overview

Readiness is reported by the Doctor Command and gated on by every command that
requires skills. If those surfaces answer the question separately, they will
eventually answer it differently, and an operator will be told a skill is fine
by one command and blocking by another.

This slice makes them call the same comparison, and pins the boundary that
keeps third-party skills out of it.

## Requirements

1. MUST make the Doctor Command report owned-skill readiness through task_02's
   comparison.
2. MUST make every command that gates on skills use that same comparison, so
   two surfaces cannot disagree.
3. MUST name, when a skill is below the minimum, all four facts: the skill, the
   required minimum, the version found, and how to upgrade.
4. MUST report `unversioned` distinctly from both satisfying and failing, at
   every surface.
5. MUST leave third-party skills with their present treatment and MUST NOT hold
   them to a version Roundfix invented for them. Finding
   `2026-07-29-doctor-requires-roundfix-own-development-skills` records what
   happened the last time Roundfix imposed its own needs on repositories that
   had no reason to hold them.
6. MUST assert the boundary directly: a third-party skill without a version
   passes, one below an owned skill's minimum is not consulted at all.

## Subtasks

- [ ] Route Doctor through the shared comparison.
- [ ] Route every gating command through it.
- [ ] Assert the four-fact diagnostic and the third-party boundary.

## Acceptance Criteria

- [ ] Doctor reports readiness through the shared comparison.
- [ ] Every gating command reports the same state for the same skill, asserted
      by comparing two surfaces on one fixture.
- [ ] A below-minimum skill names skill, minimum, found version, and upgrade
      path.
- [ ] `unversioned` is reported distinctly at every surface.
- [ ] A third-party skill without a version passes.
- [ ] A third-party skill is never compared against an owned minimum.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/baseline/catalog_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/cli ./internal/baseline -count=1 -run 'Doctor|Skills|Readiness|ThirdParty' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the surface tests ran and passed.
- `go test ./internal/cli ./internal/baseline -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 4, 5 and 6; Success Metrics 4 and 5.
- `_techspec.md` → System Architecture; Build Order 4.

## Result

### Implementation

- Centralized the owned-skill readiness states, comparison, and four-fact
  diagnostic in `skills/readiness.go`. The Task 02 interface in
  `internal/baseline/catalog_validate.go` remains available as aliases and a
  delegating wrapper, so catalog validation, Doctor, and CLI gating all execute
  the same comparison.
- Replaced repository content equality for Roundfix-owned skills with declared
  version versus owned minimum readiness. Missing, `unversioned`, `below`, and
  `satisfies` remain separate facts; only `below` blocks on version.
- Routed `roundfix doctor` and `roundfix skills check` through the shared owned
  readiness result. Both surfaces render `unversioned` distinctly, and a
  below-minimum diagnostic names the skill, required minimum, found version,
  and `roundfix skills install --target project` upgrade action.
- Preserved third-party skill handling through the existing lock hash and
  provenance path. The owned-version comparison is invoked only while
  iterating the Roundfix-owned set.

### Focused checks

- Red signal: `rtk go test ./internal/cli -run
  '^TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness$'` initially failed
  to compile because the shared readiness result was not routed to either
  surface.
- Red signal: `rtk go test ./skills -run
  '^TestRepositoryReadinessNeverComparesThirdPartySkillVersions$'` initially
  failed to compile because the repository checker had no owned-only comparison
  boundary.
- `GOCACHE=<worktree>/.gocache rtk go test ./skills -count=1` exited 0 with 134
  tests passing.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -run
  '^(TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness|TestRunDoctorRepositorySkillReadiness|TestRunDoctorRealRepositoryCheckDoesNotMutateState|TestRunSkillsCheck)$'
  -count=1` exited 0 with 14 tests passing.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/baseline -run
  '^(TestReadinessComparesDeclaredVersionToMinimum|TestCatalogRejectsMissingOwnedSkillMinimum|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent)$'
  -count=1` exited 0 with 9 tests passing.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -count=1` exited 0
  with 994 tests passing.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.
- `rtk git -c core.fsmonitor=false status --short --untracked-files=all`
  showed only this Task file and the owned-readiness implementation and test
  files; the pre-existing Daemon status transition remains `in_progress`.

### Verification feedback repair

- Daemon attempt 1 exposed a coverage-equivalence regression: Task 04 had
  updated the behavior of the recorded
  `TestOwnedSkillContractRejectsSetAndVersionDisagreement` test and also renamed
  its top-level Go function. The repository coverage record treats that stable
  test identity as removed.
- Focused reproduction
  `GOCACHE=<worktree>/.gocache rtk go test ./internal/spec -run
  '^TestCoverageEquivalence$' -count=1` exited 1 before the repair.
- Restored the recorded test function name while retaining its new
  below-minimum assertions. No coverage record, unrelated test, or production
  behavior changed for this repair.
- `GOCACHE=<worktree>/.gocache rtk go test ./skills -run
  '^(TestOwnedSkillContractRejectsSetAndVersionDisagreement|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestRepositoryReadinessNeverComparesThirdPartySkillVersions)$'
  -count=1` exited 0 with 10 tests passing.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/spec -run
  '^TestCoverageEquivalence$' -count=1` exited 0 with the focused coverage test
  passing.

### Acceptance evidence

- Doctor readiness: `skills.CheckRepositoryWithExternal` now compares every
  installed owned skill through `skills.Readiness`; Doctor renders the returned
  state. `TestRunDoctorRepositorySkillReadiness` covers repository output.
- Surface agreement: `TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness`
  sends one owned-skill readiness fixture through Doctor and `skills check` for
  `satisfies`, `below`, and `unversioned`, and asserts corresponding states and
  exit behavior.
- Four-fact failure: the shared `OwnedSkillReadiness.Diagnostic` and the surface
  agreement test assert the skill name, minimum, found version, and project
  upgrade command for `below`.
- Distinct unversioned state: the same surface test asserts `skills:
  unversioned` and `Roundfix skill check unversioned`, distinct from both the
  pass and failure labels.
- Third-party without a version: the external fixture in
  `TestRepositoryReadinessNeverComparesThirdPartySkillVersions` has no
  top-level compatibility version and remains ready after its existing lock
  hash is refreshed.
- Third-party comparison boundary: that test records every call to the shared
  comparison, asserts exactly one call per owned skill, and asserts the
  third-party skill is absent from the call ledger even though its metadata
  carries a value below the owned minimum.

The commands under this Task's `## Verification` were not run; terminal
verification remains Daemon-owned.
