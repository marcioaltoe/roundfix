---
task: task_05
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
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

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestSupersededBranch' -count=1 -v 2>&1 | grep -q '^--- PASS: TestSupersededBranchIsClassifiedWhenEveryCommitIsReachable'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestSupersededBranch' -count=1 -v 2>&1 | grep -q '^--- PASS: TestSupersededBranchRefusesAnUnreachableCommit'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestReconcileDiscard' -count=1 -v 2>&1 | grep -q '^--- PASS: TestReconcileDiscardWritesTheRecordBeforeRemoving'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestReconcileDiscard' -count=1 -v 2>&1 | grep -q '^--- PASS: TestReconcileWithoutTheFlagRemovesNothing'` — expected: exits 0.

## References

- `_prd.md` → Goal 4.
- `_techspec.md` → Build Order 5; API Contracts.
- ADR-0115.
