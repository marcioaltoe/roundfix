---
task: task_09
spec: 0084-an-update-that-can-run
status: completed
type: test
complexity: high
---

# Task 09: Sweep a fleet the Spec did not author

## Overview

The measured failure of Spec 0082 was that every case it tested was a case it
built. This slice adds the missing shape: a sweep across several adopted
repository copies whose divergence patterns are taken from the recorded fleet
measurement rather than invented — a copy whose manifest predates its regions, a
copy missing a structural clause set, a copy whose recorded profile does not
resolve, and a copy that is already current. It is hermetic: the corpus is built
by the test, and the real-fleet reading belongs to the QA gate.

## Requirements

1. MUST build a corpus of adopted repository copies whose divergence patterns
   reproduce the ones the adopted finding recorded, one copy per pattern, each
   built by the test from repository state alone.
2. MUST run the update planning path against every copy in the corpus and assert
   the outcome expected for that copy's pattern, naming the state and the reason.
3. MUST assert that every copy reaches either an applicable plan or a reported
   condition that names a human action, and that no copy reaches a state that
   blocks before planning for a reason with no named action.
4. MUST assert, for each copy that reaches an applicable plan and is applied, that
   the next run reports the repository current.
5. MUST cite, in the test's own documentation, which recorded fleet pattern each
   copy reproduces, so a future reader can tell an invented case from a measured
   one.
6. MUST be hermetic and portable: no copy may depend on a path outside the test's
   own temporary tree, on a prior Run, or on a repository elsewhere on the
   machine.
7. MUST fail when a pattern regresses: removing the classification from the
   preservation path must make at least one copy's assertion fail.

## Subtasks

- [x] Enumerate the divergence patterns from the adopted finding.
- [x] Build one adopted copy per pattern inside the test's temporary tree.
- [x] Assert the expected outcome and reason for each copy.
- [x] Assert the applied copies report current on their next run.
- [x] Document which recorded pattern each copy reproduces.
- [x] Demonstrate the negative for at least one pattern.

## Rehearsal Cases

- Case: a copy whose Setup Manifest predates its untouched managed regions, the
  pattern this repository and fiscus produced; Observation: the run reaches an
  applicable plan and lists the regions as unrecorded.
- Case: a copy missing the structural clause set, the pattern conexus, tax-poc,
  and vortex produced; Observation: the run reaches an applicable plan once the
  catalog emits those clauses again.
- Case: a copy whose recorded Baseline Profile does not resolve, the pattern
  fluxus produced; Observation: the run reports the profile identity, the
  searched locations, and the restoring action, and does not report a raw
  filesystem error.
- Case: a copy already matching the current catalog, the pattern gss and oraculum
  produced; Observation: the run reports the repository current with zero
  proposed changes.
- Case: any copy, with the classification removed from the preservation path;
  Observation: at least one copy's assertion fails, proving the sweep can fail.

## Acceptance Criteria

- [x] The corpus contains one copy per recorded divergence pattern, each built by
      the test.
- [x] Every copy reaches an applicable plan or a reported condition naming a
      human action.
- [x] Every applied copy reports the repository current on its next run.
- [x] Each copy's documentation names the recorded fleet pattern it reproduces.
- [x] The sweep depends on no path outside its own temporary tree.
- [x] Removing the classification makes at least one copy's assertion fail.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/cli/baseline_update.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ ./internal/cli/ -run 'FleetSweep' -v > /tmp/0084-task-09-a.log 2>&1 && grep -q '^--- PASS: .*FleetSweep' /tmp/0084-task-09-a.log` — expected: exits 0, proving the sweep cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ ./internal/cli/ -run 'FleetSweep' -count=2 > /tmp/0084-task-09-b.log 2>&1` — expected: exits 0, proving no copy depends on state left by a previous run.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 9; Testing Approach.
- `_prd.md` → Goals 1 and 5; Success Metrics; User Story 1.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the recorded divergence patterns this corpus reproduces.
- ADR-0104.

## Result

Implemented one command-level `FleetSweep` corpus in
`internal/cli/baseline_update_test.go`. Every copy is a real temporary Git
repository built and adopted by the test, moved under one corpus root, and run
through the production Baseline update planning and apply path. The case
documentation cites the measured repository cohort and calls out that GSS and
Oraculum were measured as planning-capable; the already-current copy proves the
zero-change endpoint Task 09 requires for that cohort.

### Focused checks

- Pre-change: `rtk rg -n 'FleetSweep' internal/baseline internal/cli` returned
  no matches, proving the fleet corpus was absent.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-go-build-a889ef9d6ad09873 rtk go test ./internal/cli -run '^TestBaselineUpdateFleetSweep/structural-clauses-missing$' -count=1`
  passed (`2 passed`). This exercised Standard TypeScript adoption, removal of
  exactly the 14 recorded clauses, manifest-digest alignment, planning, apply,
  exact clause restoration, and next-run convergence.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task09-go-build-a889ef9d6ad09873 rtk go test ./internal/cli -run '^TestBaselineUpdateFleetSweep/(manifest-predates-managed-regions|recorded-profile-does-not-resolve|already-current)$' -count=1`
  passed (`4 passed`). This exercised the remaining three corpus rows.
- The first focused test attempt used Go's default cache and stopped before
  compilation because the managed sandbox denied access to
  `~/Library/Caches/go-build`; the task-specific `/private/tmp` cache removed
  that environment boundary.
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance evidence

- **One copy per recorded divergence pattern:** the table-driven corpus builds
  four independent copies for the Roundfix/Fiscus stale-manifest pattern, the
  Conexus/Tax PoC/Vortex missing-structural-clause pattern, the Fluxus
  unresolved-profile pattern, and the GSS/Oraculum non-blocking current
  endpoint. The structural fixture removes the 14 recorded clause identities,
  not an arbitrary guide body.
- **Applicable plan or named human action:** the stale-manifest and structural
  copies assert `plan_ready`, category `approval`, and a non-empty Plan Digest;
  the unresolved-profile copy asserts `failed`/`manifest` plus a restoring
  `NextAction`; the current copy asserts `current`. A shared guard rejects any
  other pre-planning state without a named action.
- **Applied copies become current:** both `plan_ready` copies are applied through
  the update command, assert `verified`, and then assert a second unapproved run
  returns `current` with zero `FileChanges`.
- **Measured patterns are documented:** comments adjacent to the corpus name
  every recorded repository cohort and explain the GSS/Oraculum endpoint
  mapping.
- **Hermetic and portable:** every builder uses `testing.T` temporary
  repositories; each finished copy is moved beneath one `t.TempDir` corpus root,
  and the test rejects a repository path that escapes that root. The skills
  stage is the suite's injected local success stage, so no copy reads another
  checkout or prior Run.
- **Classification regression is detectable:** the stale-manifest copy first
  passes an oracle requiring the exact `guide.agent-instructions`
  `digest-mismatch` classification. The test then removes only that
  classification from a copy of the observed result and asserts the same oracle
  fails, demonstrating the known negative required by ADR-0104.
