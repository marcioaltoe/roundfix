---
task: task_05
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
type: backend
complexity: high
---

# Task 05: Give a superseded Run Branch a disposition

## Overview

`reconcile` classifies a stopped Run's branch `unintegrated` and preserves it;
Branch Integrity Preflight then refuses to create any Run while one exists. On
2026-08-09 three such branches refused two consecutive `roundfix resolve`
invocations, and the suggested `git merge --ff-only` could not apply because they
had diverged behind the work that superseded them. Clearing them meant
`git branch -D` by hand, at the moment the maintainer was already blocked.

## Requirements

1. MUST classify a Run Branch as superseded when every commit it holds is
   reachable from the target branch, or when a later Run covered the same Tasks
   and integrated.
2. MUST write what the branch held — its commits, its changed files, and why it
   was classified superseded — before removing anything.
3. MUST remove the branch and its worktree only under an explicit flag, never as
   a side effect of reporting.
4. MUST refuse to discard any branch it cannot prove superseded, and say which
   condition failed.
5. MUST leave `reconcile` without the flag byte-identical in behaviour: it
   reports and never disposes, per ADR-0115.

## Subtasks

- [ ] Classify superseded.
- [ ] Write the branch record.
- [ ] Remove under the explicit flag only.

## Acceptance Criteria

- [ ] A branch whose commits are all reachable from the target is classified
      superseded.
- [ ] A branch holding unreachable commits is refused, naming the condition.
- [ ] The record is written before removal and survives it.
- [ ] `reconcile` without the flag changes nothing.

## Rehearsal Cases

- Case: a Run Branch whose two commits are both reachable from the target;
  Observation: classified superseded, record written, branch and worktree
  removed under the flag.
- Case: a Run Branch holding one commit absent from the target; Observation:
  refused, naming the unreachable commit, nothing removed.
- Case: the same superseded branch without the flag; Observation: reported as
  superseded and left in place.

## Bounded scope

This Task may create or modify only:

- `internal/cli/reconcile.go`
- `internal/cli/reconcile_test.go`
- `internal/daemon/reconcile.go`
- `internal/daemon/reconcile_test.go`
- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_05.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestSupersededBranch' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSupersededBranchIsClassifiedWhenEveryCommitIsReachable'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestSupersededBranch' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSupersededBranchRefusesAnUnreachableCommit'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestReconcileDiscard' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestReconcileDiscardWritesTheRecordBeforeRemoving'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestReconcileWithoutTheFlagRemovesNothing$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestReconcileWithoutTheFlagRemovesNothing'` — expected: exits 0.

## References

- `_prd.md` → Goal 4.
- `_techspec.md` → Build Order 5; API Contracts.
- ADR-0115.

## Result

Implemented a proof-bound `--discard-superseded` disposition. The Daemon now
inventories every Run commit and changed file after the recorded starting head,
classifies the branch from target reachability or a later Clean Run's completed
Task coverage, snapshots the proof for revalidation, writes an atomic Artifact
Root record, and only then removes the clean Run Worktree and Run Branch. The
CLI reports the same proof without mutation when the flag is absent and returns
a non-zero refusal when neither supersession condition can be proven.

Acceptance evidence:

- A branch whose commits are all reachable from the target is classified
  superseded: `TestSupersededBranchIsClassifiedWhenEveryCommitIsReachable`
  passed with two inventoried commits and both changed paths;
  `TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks` also
  passed for the alternative later-Run Task-coverage proof.
- A branch holding an unreachable commit is refused with the failed condition:
  `TestSupersededBranchRefusesAnUnreachableCommit` and
  `TestReconcileDiscardRefusesAnUnreachableCommit` passed, with the exact commit
  object ID in the refusal and both Git surfaces preserved.
- The record is written before removal and survives it:
  `TestWriteBranchDispositionRecordSurvivesDiscard` and
  `TestReconcileDiscardWritesTheRecordBeforeRemoving` passed; the record kept
  the schema version, Run and branch identities, commit inventory, changed-file
  inventory, and classification reason after both Git surfaces were absent.
  `TestReconcileDiscardKeepsSurfaceWhenRecordCannotBeWritten` passed and proved
  that a record-write failure leaves both surfaces present.
- `reconcile` without the flag changes nothing:
  `TestReconcileWithoutTheFlagRemovesNothing` passed after comparing the Run
  Database bytes and exact Git-surface snapshot before and after the report;
  the existing `TestRunReconcile*` focused set also passed unchanged.

Focused checks:

- Verification-feedback attempt 1: inspected
  `batch-004-attempt-1.log`; its selected `TestReconcileDiscard*` tests all
  passed, but the selector excluded the required
  `TestReconcileWithoutTheFlagRemovesNothing` test. Focused `go test -list`
  discovery confirmed the two names belong to disjoint selector sets, so the
  no-flag Verification entry now selects that existing test exactly. No test
  assertion or production behavior changed for this repair.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestReconcileWithoutTheFlagRemovesNothing$' -count=1` — passed after the selector repair.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -list '^TestReconcileWithoutTheFlagRemovesNothing$'` — listed the exact intended test.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^Test(SupersededBranch|WriteBranchDispositionRecord)' -count=1` — passed.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^Test(ReconcileDiscard|ReconcileWithoutTheFlag)' -count=1` — passed.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunReconcile' -count=1` — passed.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/daemon ./internal/cli -count=1` — the initial sandboxed attempt was blocked by denied `api.github.com` access; the final approved rerun after the last code edit passed both packages (`internal/daemon` in 3.675s and `internal/cli` in 62.723s).
- `rtk env GOCACHE="$PWD/.gocache" go vet ./internal/daemon ./internal/cli` — passed.
- `rtk git -c core.fsmonitor=false diff --check` — passed before this Result update and is rerun in the final bounded-diff audit.

Daemon-owned Verification commands were not run. Follow-up: the new flag's
top-level and command-help copy lives in `internal/cli/cli.go`, which is outside
this Task's bounded file list, so that discoverability update remains outside
this diff.
