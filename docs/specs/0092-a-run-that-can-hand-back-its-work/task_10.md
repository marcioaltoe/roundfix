---
task: task_10
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: docs
complexity: low
---

# Task 10: Let the two new acts be discoverable and named

## Overview

The QA gate returned `fail` twice on the same two findings, and both are the
same shape as Tasks 08 and 09: a Task delivered a public change whose
surrounding contract sat outside every bounded scope, so nobody could complete
it.

- **F-02.** `reconcile` gained `--discard-superseded` in Task 05 and
  `--carry-forward` in Task 06, and its help still lists only `--apply` and
  `--format`. Task 05's own Result recorded the gap when it happened: the help
  copy lives in `internal/cli/cli.go`, outside its bounded file list. An act a
  Run can perform and an operator cannot find is not delivered.
- **F-03.** Task 07 names `Selection Failure` and `Branch Disposition` as this
  Spec's coined terms. `CONTEXT.md` documents the first and not the second.
  `CONTEXT.md` is inside Task 07's bounded scope, but the `qa-gate` skill
  forbids the gate from writing anything but its report and evidence, so that
  scope entry can never be exercised — which is why the finding repeated
  instead of being fixed.

## Requirements

1. MUST document `--discard-superseded` and `--carry-forward` in the same help
   surfaces that already document `--apply` and `--format`, using the flag
   names the implementation accepts.
2. MUST describe each act by what it does to the Run Branch or the Task, not by
   its implementation, so the copy matches the vocabulary the rest of the help
   uses.
3. MUST add `Branch Disposition` to the `CONTEXT.md` glossary, defined so it
   reads as a sibling of the existing `Selection Failure` entry rather than a
   restatement of the flag.
4. MUST NOT change flag parsing, behaviour, or any test. This Task changes help
   copy and one glossary entry.

## Subtasks

- [ ] Document both flags in every help surface that lists `--apply`.
- [ ] Add the `Branch Disposition` glossary entry.

## Acceptance Criteria

- [ ] `reconcile --help` lists both new flags with a description each.
- [ ] `CONTEXT.md` defines `Branch Disposition`.
- [ ] No behaviour or test changes appear in the diff.

## Bounded scope

This Task may create or modify only:

- `internal/cli/cli.go`
- `CONTEXT.md`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_10.md`

## Verification

- `go run -buildvcs=false ./cmd/roundfix reconcile --help 2>&1 | grep -q -- '--discard-superseded'` — expected: exits 0. The help lists only `--apply` and `--format` before this Task, so the command fails against the unchanged tree.
- `go run -buildvcs=false ./cmd/roundfix reconcile --help 2>&1 | grep -q -- '--carry-forward'` — expected: exits 0.
- `grep -q 'Branch Disposition' CONTEXT.md` — expected: exits 0. The term is absent before this Task.
The reconcile suite is deliberately not listed here. It passes against the
unchanged tree, so asserting it would approve this Task before any work
happened — the pre-work probe refused exactly that on 2026-08-11. Proving the
help edit disturbed nothing is the Run-level gate's job.

## References

- `_prd.md` → Goals 3 and 4.
- `task_05.md` → the `--discard-superseded` act and its recorded help gap.
- `task_06.md` → the `--carry-forward` act.
- `task_07.md` → Requirement 10 and the two coined terms.
