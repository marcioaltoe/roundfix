---
task: task_09
spec: 0084-an-update-that-can-run
status: pending
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

- [ ] Enumerate the divergence patterns from the adopted finding.
- [ ] Build one adopted copy per pattern inside the test's temporary tree.
- [ ] Assert the expected outcome and reason for each copy.
- [ ] Assert the applied copies report current on their next run.
- [ ] Document which recorded pattern each copy reproduces.
- [ ] Demonstrate the negative for at least one pattern.

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

- [ ] The corpus contains one copy per recorded divergence pattern, each built by
      the test.
- [ ] Every copy reaches an applicable plan or a reported condition naming a
      human action.
- [ ] Every applied copy reports the repository current on its next run.
- [ ] Each copy's documentation names the recorded fleet pattern it reproduces.
- [ ] The sweep depends on no path outside its own temporary tree.
- [ ] Removing the classification makes at least one copy's assertion fail.

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
