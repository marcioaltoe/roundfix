---
task: task_11
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: test
complexity: low
---

# Task 11: Make the assembled tree pass its own gate

## Overview

`TestRunCommandHelp` pins the `reconcile` synopsis by exact string, and Task 10
changed that synopsis to name the two acts Tasks 05 and 06 added. The contract
lives in `internal/cli/cli_test.go`, outside Task 10's bounded scope, so Task 10
could not update it and the QA gate returned `fail` a third time.

This is the fourth Task this Spec has needed for the same reason: a public
change whose surrounding contract sat outside every boundary. The first three
were minted one instance at a time, and each rerun of the gate found the next
one. This Task closes the class instead of the instance — its acceptance is the
repository gate itself, so it cannot settle while any contract this Spec broke
is still pinned to the pre-Spec text.

## Requirements

1. MUST update every contract assertion that pins text this Spec changed, so
   each states what the surface now says.
2. MUST keep every such assertion exact. Do not relax a pinned string to a
   substring or a regex to make it pass; the point of pinning is that an
   unannounced change fails.
3. MUST NOT change production code, help copy, or behaviour. If an assertion
   and the surface genuinely disagree about what is correct, stop and report
   which one is wrong rather than editing the test to match.

## Subtasks

- [ ] Update the pinned `reconcile` synopsis.
- [ ] Run the repository gate and resolve whatever else it names.

## Acceptance Criteria

- [ ] `make verify` exits 0 on the assembled tree.
- [ ] Every updated assertion still compares exact text.
- [ ] `git diff --name-only` lists only test files and this Task file.

## Bounded scope

This Task may create or modify only:

- `internal/cli/cli_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_11.md`

## Verification

- `go run -buildvcs=false ./cmd/roundfix reconcile --help 2>&1 | grep -q -- '--carry-forward'` — expected: exits 0, proving the help copy Task 10 wrote is still the copy the contract pins.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunCommandHelp$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunCommandHelp'` — expected: exits 0. The assertion pins the pre-Spec synopsis and fails against the unchanged tree.
- `make verify` — expected: exits 0. It fails against the unchanged tree, which is what makes this command a proof rather than a formality: this Task is done when the assembled work passes the gate it currently fails.

## References

- `_prd.md` → Goals 3 and 4.
- `task_10.md` → the help copy this contract pins.
