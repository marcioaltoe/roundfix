---
task: task_06
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: backend
complexity: high
---

# Task 06: Hand a stopped Run's settled Tasks back

## Overview

Implementing Spec 0089 re-executed Task 01 in four separate Runs and Task 02 in
three, against unchanged inputs, because a Run that stops leaves its settled
Tasks reading `pending` in the checkout while their commits sit in the Run
Worktree. This Task hands those Tasks back, on evidence, as an explicit act.

## Requirements

1. MUST identify a stopped Run's Tasks whose Verification passed and which the
   Run committed inside its Worktree.
2. MUST refuse to carry forward any Task whose declared inputs have changed since
   it settled; a Task reported complete against inputs it never saw is worse than
   re-executing it.
3. MUST carry forward only under an explicit flag; the act is never automatic.
4. MUST record, per carried Task, the Run it came from and the commit that
   settled it, so the checkout's `completed` status is traceable.
5. MUST refuse the whole operation rather than carrying a subset silently when
   any named Task fails its condition, and say which.

## Subtasks

- [ ] Identify settled Tasks in a stopped Run's Worktree.
- [ ] Prove their inputs are unchanged.
- [ ] Carry forward under the explicit flag with the record.

## Acceptance Criteria

- [ ] A settled Task with unchanged inputs is carried forward and reads
      `completed` with its source Run and commit recorded.
- [ ] A settled Task whose inputs moved is refused, naming the input.
- [ ] Nothing is carried forward without the flag.
- [ ] A mixed set refuses rather than carrying part of itself.

## Rehearsal Cases

- Case: a stopped Run holding two settled Tasks whose inputs are unchanged;
  Observation: both read `completed` in the checkout with their Run and commit
  recorded.
- Case: the same Run after one Task's declared input file changed; Observation:
  the operation refuses, names that Task and that input, and carries nothing.
- Case: the same Run without the flag; Observation: the settled Tasks are
  reported and the checkout is unchanged.

## Bounded scope

This Task may create or modify only:

- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `internal/cli/reconcile.go`
- `internal/cli/reconcile_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_06.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestCarryForward' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCarryForwardSettlesATaskWhoseInputsAreUnchanged'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestCarryForward' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCarryForwardRefusesATaskWhoseInputsMoved'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestCarryForward' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCarryForwardRefusesRatherThanCarryingASubset'` — expected: exits 0.
Regression across `internal/cli` and `internal/spec` is the Run-level gate's
job, not this Task's. A whole-package sweep here passes against the unchanged
tree whenever the tree is already green, which is what the pre-work probe
refused on 2026-08-11: it approves the Task before any work happens. The three
commands above name cases that do not exist yet, so each can fail.

## References

- `_prd.md` → Goal 3.
- `_techspec.md` → Build Order 6; Risks.
- `docs/backlog/2026-08-09-a-stopped-run-discards-the-tasks-it-already-proved.md`
