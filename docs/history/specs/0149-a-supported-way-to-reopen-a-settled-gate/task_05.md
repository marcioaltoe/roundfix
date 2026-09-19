---
task: task_05
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: completed
type: backend
complexity: medium
---

# Task 05: Make the reopen safe to run

## Overview

Pre-PR review found three ways the command can do damage the Spec never intended.
All three were reproduced before this Task was written.

1. **It does not check for an active Run.** `settle` guards with
   `ensureNoSettleActiveRun` because a Run that owns the checked-out tree can
   snapshot, commit or push whatever it finds there. `reopen` writes the QA Task
   without that guard, so an active review Run can carry the status change into
   unrelated work.
2. **Its two writes are not atomic.** `SetStatus` persists `status: pending`,
   then `AppendGateInvalidation` records why. If the second fails, the gate is
   reopened with no record of the invalidation — and a later `reopen` refuses,
   because the gate is no longer completed. The recovery path loses exactly the
   audit trail it exists to keep.
3. **It accepts a path-like slug.** `--spec ../other` resolves to
   `docs/other`, outside the configured Spec Root. Measured: the command
   reported `no _prd.md under "/Users/marcio/dev/roundfix/docs/other"`, which is
   a refusal by accident — a traversal reaching a directory that does hold a
   Spec would be rewritten instead.

## Requirements

1. MUST refuse when a Run is active for the Spec or the git root, reusing the
   guard `settle` already uses rather than adding a second one.
2. MUST leave the QA Task unchanged when the invalidation record cannot be
   written: either both writes land or neither does.
3. MUST reject a `--spec` value that is not a single Spec directory name, and
   MUST confirm the resolved QA Task path stays inside the configured Spec Root.
4. MUST keep every behavior Tasks 01 through 03 delivered.

## Subtasks

- [ ] Add the active-Run guard.
- [ ] Make the status write and the invalidation record one replacement.
- [ ] Validate the slug and the resolved path.
- [ ] Add a test for each of the three.

## Acceptance Criteria

- [ ] With an active Run recorded for the Spec, the command refuses and writes
      nothing.
- [ ] When the invalidation record cannot be written, the QA Task file is
      byte-identical to what it was before the command ran.
- [ ] `--spec ../other`, `--spec a/b` and an absolute value are each refused as
      invalid input, naming that condition.
- [ ] The existing reopen tests still pass unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reopen.go`
- interface: `internal/cli/settle.go`

## Verification

- `out="$(go test -count=1 -run "^TestReopenRefusesWhileARunIsActive" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestReopenRejectsPathLikeSpecValues" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestReopenLeavesTheTaskUnchangedWhenTheRecordCannotBeWritten" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

- [_techspec.md](_techspec.md) — Risks & Considerations

## Result

Implemented the reopen safety boundary without changing its stale-gate
predicate or its existing success and refusal contracts:

- `reopen` now calls the same Active Run guard as `settle` before it reads or
  writes gate evidence.
- `--spec` now accepts one directory name only. Before mutation, the command
  also resolves the QA Task and confirms that its canonical path remains under
  the configured Spec Root.
- The status change and invalidation record are now assembled in memory and
  installed with one temporary-file rename. A failed replacement leaves the
  original QA Task bytes in place.

Acceptance evidence:

- Active Run refusal: `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go
  test -count=1 -run
  'TestReopen(RefusesWhileARunIsActive|RejectsPathLikeSpecValues|LeavesTheTaskUnchangedWhenTheRecordCannotBeWritten)$'
  ./internal/cli` passed. `TestReopenRefusesWhileARunIsActive` creates a live
  Spec Run, observes the preflight refusal, and compares the QA Task bytes
  before and after.
- Atomic failure: the same focused command passed.
  `TestReopenLeavesTheTaskUnchangedWhenTheRecordCannotBeWritten` makes the QA
  Task directory read-only, observes the write failure, and compares the file
  byte-for-byte with its original contents.
- Invalid input and confinement: the same focused command passed for
  `../other`, `a/b`, and an absolute value, with each refusal naming an invalid
  Spec slug. `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test
  -count=1 -run '^TestReopen' ./internal/cli` also passed the valid-looking
  symlink escape case, which names the configured Spec Root boundary.
- Existing reopen behavior: `rtk env
  GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run
  '^TestReopen' ./internal/cli` passed all reopen command tests, including the
  previously existing success, pending-gate, healthy-gate, missing-QA, and
  unknown-flag cases.

Additional focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run
  'Test(ReopenGate|AppendGateInvalidation|SetStatus)' ./internal/spec` passed.
- `rtk git diff --check` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache make
  verify-incremental` passed with process-table access enabled. The first
  sandboxed attempt denied process inspection in two unrelated force-stop
  integration tests; rerunning those tests with that access passed before the
  complete incremental gate rerun.

No follow-up work was discovered inside this Task's slice.
