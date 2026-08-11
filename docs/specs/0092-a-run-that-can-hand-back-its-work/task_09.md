---
task: task_09
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: test
complexity: low
---

# Task 09: Let the reconcile JSON contract know about carry-forward

## Overview

Task 06 added a `carryForwards` array to `reconcile --format json`. The contract
test that pins that payload, `TestRunReconcileJSONMatchesTextFields`, compares
the field set by exact equality and lives in `internal/cli/cli_test.go`, which
is outside Task 06's bounded scope. Task 06 could not update it, so `make
verify` fails on the assembled work and the QA gate returned `fail` on
2026-08-11.

This is the second test the Spec breaks that no Task owned; Task 08 was the
first. The pattern is worth naming in the QA report: a Task that widens a
public payload needs the test pinning that payload inside its boundary.

## Requirements

1. MUST add `carryForwards` to the expected field set so the assertion states
   the payload Task 06 actually emits.
2. MUST keep the assertion an exact-equality check; do not relax it to a subset
   or a contains check, because the point of the test is that an unannounced
   field cannot appear.
3. MUST NOT change production code or any other test. If the JSON payload turns
   out to differ from Task 06's Result in any other field, stop and report it
   rather than absorbing the difference.

## Subtasks

- [ ] Add the field to the expected set.
- [ ] Confirm the assertion remains exact equality.

## Acceptance Criteria

- [ ] `TestRunReconcileJSONMatchesTextFields` passes.
- [ ] The expected set is still compared by exact equality.
- [ ] `git diff --name-only` lists only this Task's bounded paths.

## Bounded scope

This Task may create or modify only:

- `internal/cli/cli_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_09.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunReconcileJSONMatchesTextFields$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunReconcileJSONMatchesTextFields'` — expected: exits 0. The test fails against the unchanged tree, so this command cannot pass before the work.
- `grep -q 'carryForwards' internal/cli/cli_test.go` — expected: exits 0, proving the field was named rather than the assertion loosened.

## References

- `_prd.md` → Goal 3.
- `task_06.md` → the `carryForwards` payload this contract pins.
