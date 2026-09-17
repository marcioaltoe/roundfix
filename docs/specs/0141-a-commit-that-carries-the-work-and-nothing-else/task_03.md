---
task: task_03
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: completed
type: backend
complexity: medium
---

# Task 03: Keep the Verification's writes out of the QA commit

## Overview

The QA step snapshots the worktree before it begins, and the Daemon runs the
repository Verification after that snapshot, so whatever the verifier writes is
a candidate for the QA Report commit. This slice measures the Verification's own
window and subtracts what appeared inside it.

## Requirements

1. MUST take a snapshot immediately before the repository Verification starts
   and another immediately after it ends.
2. MUST exclude from the QA Report commit every path that appeared between those
   two snapshots.
3. MUST keep carrying the report the mechanical stage seeded, its evidence, and
   everything the QA Agent writes after the Verification.
4. MUST leave a Spec whose Verification writes nothing with exactly the commit it
   produces today.
5. MUST NOT change the mechanical stage, the withholding rule, the verdict
   settlement, or the QA Report's shape.
6. MUST name the excluded set `verificationWindowPaths`, so the Task's own check
   and later readers agree.

## Subtasks

- [ ] Take the two snapshots around the repository Verification.
- [ ] Subtract the window's paths from the report commit's path set.
- [ ] Cover a Verification that writes a tracked file, and one that writes
      nothing, at the QA-step seam.

## Acceptance Criteria

- [ ] A fixture whose repository Verification writes a tracked file produces a QA
      Report commit that does not carry it.
- [ ] A fixture whose Verification writes nothing produces the same commit as
      before this Task.
- [ ] The seeded report and the Agent's own writes stay in the commit.
- [ ] The verdict, the report shape and the withholding rule are untouched.

## Context

- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q "verificationWindowPaths" internal/daemon/task_engine.go` — expected: exit 0; the QA step measures the Verification's window. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestQAReportCommitExcludesVerificationWrites$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 3; User Story 4; Goal 3; Success Metric 3;
`_techspec.md` → Implementation Design: The QA Report commit excludes the
Verification's writes; API Contract 3; Build Order 3; ADR-0096; ADR-0117.

## Result

Implementation:

- The QA repository Verification now takes worktree snapshots immediately
  before and after its execution and derives `verificationWindowPaths` from
  that bounded window, including failed and unobserved command outcomes that
  still materialize a refusal report.
- The QA Report commit removes those paths from the QA step's full snapshot
  diff before applying the existing report and QA Task file guarantees. The
  mechanical stage, withholding decision, verdict settlement and report shape
  were not changed.
- `TestQAReportCommitExcludesVerificationWrites` exercises the QA-step seam with
  a tracked file written by Verification and with a Verification that writes
  nothing. Both cases require the same report, evidence, QA Agent write and QA
  Task file commit paths.

Focused checks:

- `rtk proxy go test -count=1 -run '^TestQAReportCommitExcludesVerificationWrites/tracked_file$' ./internal/daemon`
  failed before the production change because the scripted commit snapshot was
  consumed without a Verification window and the QA Agent evidence paths were
  absent from the resulting commit.
- `rtk proxy go test -count=1 -run '^TestQAReportCommitExcludesVerificationWrites/(tracked_file|no_writes)$' ./internal/daemon`
  passed after the production change.
- `rtk proxy go test -count=1 -run '^(TestQAGate|TestMechanicalStageSeedsReportBeforeAgentSession|TestQAReportCommitExcludesVerificationWrites)' ./internal/daemon`
  passed the focused QA gate regression suite.
- `rtk proxy go test -count=1 ./internal/daemon` passed the complete daemon
  package.
- `rtk git diff --check` passed.
- The Task's authored `## Verification` commands were not run; the Daemon owns
  those checks.

Acceptance evidence:

- Verification writes are excluded: the `tracked file` case modifies a
  Git-tracked file during the verifier call and asserts that the path is absent
  from the QA Report commit.
- A no-write Verification keeps the prior commit contents: the `no writes`
  case asserts the same exact commit path set without a Verification-window
  exclusion.
- The seeded report, its evidence and QA Agent writes stay included: both cases
  assert the report, evidence file and Agent source write in the commit, and
  assert that the Agent received non-empty mechanical report seed content.
- Verdict, report shape and withholding remain unchanged: both cases settle the
  existing `pass` verdict, the focused QA gate regression suite passed, and the
  production diff changes only snapshot capture and commit-path filtering.

## Carry-forward provenance

- Source Run: `run_20260917T192120Z_0e305cd0796d7680`
- Source commit: `eff761a19f168ec447d8a2062cd347e160110299`
