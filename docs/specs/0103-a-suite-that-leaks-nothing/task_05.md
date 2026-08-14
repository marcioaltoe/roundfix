---
task: task_05
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 05: End the processes the detach tests prove survive

## Overview

The detach tests prove a process outlives its caller, which is the product
behaviour and stays. What does not stay is the process: four of them survived up
to three days and burned two hours and forty minutes of CPU. The cooperative
release sentinel also lives inside a directory the test framework deletes, so a
late waiter is stranded by construction.

## Requirements

1. MUST record every process a test proves survives and end it
   non-cooperatively at teardown.
2. MUST place the cooperative release sentinel outside any directory the test
   framework removes, so a late waiter can still observe it.
3. MUST keep proving the survival property; the assertion under test is unchanged.
4. MUST leave no process from these tests running after the package's tests end.

## Subtasks

- [ ] Record and terminate the proven survivors at teardown.
- [ ] Move the release sentinel outside the framework-deleted directory.
- [ ] Cover the teardown with an assertion that the process is gone.

## Acceptance Criteria

- [ ] Each detach test's spawned process is gone when the package's tests end.
- [ ] Termination is non-cooperative, so a process ignoring the sentinel still
      ends.
- [ ] The sentinel outlives the temporary directory.
- [ ] The survival property is still asserted.

## Verification

- `go test -count=1 ./internal/cli -run '^TestRunDetached' -v > /tmp/0103-t05.log 2>&1; s=$?; test $s -eq 0 || { cat /tmp/0103-t05.log; exit 1; }; grep -q '^--- PASS: TestRunDetached' /tmp/0103-t05.log || { cat /tmp/0103-t05.log; exit 1; }; grep -rq 'TestDetachedChildIsTerminatedAtTeardown' internal/cli/ || { echo 'the detach tests still pass, but nothing ends the processes they prove survive'; exit 1; }` — expected: exits 0, proving the survival property is still asserted by its eight existing tests *and* that a teardown now ends what they prove. The survival half alone passes on an untouched tree, so it is anchored rather than left standing.
- `go test -count=1 ./internal/cli -run 'TestDetachedChildIsTerminatedAtTeardown' -v > /tmp/0103-t05-teardown.log 2>&1; s=$?; grep -q '^--- PASS: TestDetachedChildIsTerminatedAtTeardown' /tmp/0103-t05-teardown.log || { cat /tmp/0103-t05-teardown.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing teardown assertion; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t05-teardown.log` — expected: exits 0, refusing a vacuous run.
- `grep -q 'sentinel' internal/cli/detach_test.go || { echo 'the detach tests name no release sentinel'; exit 1; }; grep -n 'sentinel' internal/cli/detach_test.go | grep 't.TempDir()' && { echo 'the release sentinel still lives in a framework-deleted directory'; exit 1; }; exit 0` — expected: exits 0, proving the sentinel exists and sits outside the deleted tree. It prints the offending line on failure. Fails today, where no sentinel is named at all.

## Context

- interface: `internal/cli/detach_test.go`

## References

`_techspec.md` → Build Order 5; Testing Approach, the detach teardown.
`_prd.md` → Core Feature 2; Goal 2; User Story 2; Non-Goals, the survival
property.
