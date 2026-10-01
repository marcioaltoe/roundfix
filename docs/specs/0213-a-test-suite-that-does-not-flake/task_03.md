---
task: task_03
spec: 0213-a-test-suite-that-does-not-flake
status: pending
type: test
complexity: high
---

# Task 03: Fixture processes end with the test binary that started them

## Overview

`TestRunImplementDetachSurvivesCallerProcessGroupKill` proves that a detached
child survives its caller. It releases that child only from a `t.Cleanup`, and
its fake ACPX polls for a release file forever. When the test binary dies,
cleanup never runs and both processes stay alive. In three of three interrupted
authoring runs they did, and on 2026-09-30 one stayed alive for five hours. The
survivor children in `detach_test.go` have the same unbounded wait. This Task
makes every fixture wait end when its test binary is gone, reaps the detached
child by its recorded Owner PID, and proves both with tests that kill their own
inner test binary (ADR-0213).

## Requirements

1. MUST add `ROUNDFIX_FAKE_ACPX_OWNER_PID`, set to the test binary's PID, to
   the environment `cliHelperEnv` builds. On every pass, the fake ACPX release
   loop MUST check that process with `kill -0` and exit 1 once it is gone.
2. MUST pass the test binary's PID to the `detach_test.go` children. Both
   `waitForDetachTestSentinel` and the `detachTestChildIgnoreSentinel` loop
   MUST return `exitRunFailed` once `store.ProcessAlive` reports that PID gone.
   The second child MUST keep ignoring its sentinel.
3. MUST make `TestRunImplementDetachSurvivesCallerProcessGroupKill` read the
   Run's `OwnerPID` after the Run ID is printed. After every `t.TempDir` it
   uses, it MUST register a cleanup that sends `SIGKILL` to the process group
   `-OwnerPID`, accepts `ESRCH`, and polls `store.ProcessAlive` to false within
   the test deadline. When `ROUNDFIX_DETACH_TEST_OWNER_RECORD` names a file,
   the test MUST write the Owner PID to it. The survival assertions MUST stay
   unchanged.
4. MUST add `TestImplementDetachChildEndsWhenItsTestBinaryDies` in the new
   `internal/cli/implement_detach_teardown_test.go`:
   - it re-executes the compiled test binary for
     `^TestRunImplementDetachSurvivesCallerProcessGroupKill$`, with a private
     `TMPDIR` under its own temporary directory, the record variable set, and
     its own process group;
   - it waits for the recorded Owner PID, then sends `SIGKILL` to the inner
     binary's process group;
   - it polls until the Owner PID and its process group are gone, bounded by
     the test deadline.
5. MUST add `TestDetachSurvivorEndsWhenItsTestBinaryDies` in the same file. It
   does the same for `^TestDetachedChildIsTerminatedAtTeardown$`, and reads
   the survivor PID file from the `roundfix-detach-test-*` directory under the
   private `TMPDIR`.
6. MUST signal only PIDs the test started or read from its own records, and
   MUST NOT change production code.

## Subtasks

- [ ] Make the fake ACPX and the detach children watch their test binary.
- [ ] Reap the Implement detach child by its Owner PID at teardown.
- [ ] Add the two tests that kill their inner test binary.
- [ ] Prove the death tests fail at the deadline on the old waits.

## Acceptance Criteria

- [ ] After its test binary is killed with `SIGKILL`, the Implement detached
      child and its fake ACPX are gone within the test deadline.
- [ ] After its test binary is killed, the `detach_test.go` survivor that
      ignores its sentinel is gone within the test deadline.
- [ ] `TestRunImplementDetachSurvivesCallerProcessGroupKill` and
      `TestDetachedChildIsTerminatedAtTeardown` still pass with their survival
      assertions unchanged.

## Context

- creates: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/detach_test.go`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `docs/adr/0028-operational-runs-can-detach-from-their-caller.md`
- instruction: `internal/store/process.go`
- instruction: `internal/cli/detach_unix.go`

## Verification

- `out="$(go test -count=1 -timeout 600s -v -run '^(TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestDetachedChildIsTerminatedAtTeardown)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestImplementDetachChildEndsWhenItsTestBinaryDies TestDetachSurvivorEndsWhenItsTestBinaryDies TestRunImplementDetachSurvivesCallerProcessGroupKill TestDetachedChildIsTerminatedAtTeardown; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two death tests do not exist, so the command fails.
- `grep -q 'ROUNDFIX_FAKE_ACPX_OWNER_PID' internal/cli/implement_test.go && grep -q 'ROUNDFIX_DETACH_TEST_OWNER_RECORD' internal/cli/implement_test.go` — expected: exit 0; before this Task neither variable exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — The owner watch and reaping; Testing Approach 3; Build Order 3
- [references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md](references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md)
- ADR-0213; ADR-0028; ADR-0126
