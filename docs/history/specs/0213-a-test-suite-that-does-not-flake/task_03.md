---
task: task_03
spec: 0213-a-test-suite-that-does-not-flake
status: completed
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

## Result

Implemented this Task's fixture lifetime contract without production changes.
`cliHelperEnv` passes the test binary's PID as
`ROUNDFIX_FAKE_ACPX_OWNER_PID`; the fake ACPX checks it with `kill -0` on
every release-loop pass and exits 1 when it is gone. Both detach survivor
loops check their owning test binary through `store.ProcessAlive` and return
`exitRunFailed` on owner death. The sentinel-ignoring child still ignores its
sentinel, records its PID, and runs in its own session so killing the inner
test binary's process group cannot directly kill that child.

The Implement survival test reads `OwnerPID` from the Run Database after the
Run ID is printed. Cleanup registrations follow the workspace, prompt and
adapter temporary directories, kill only the recorded Owner PID's process
group (accepting `ESRCH`), and prove PID and group absence through
`testwait.Poll`. A caller cleanup also covers failures before the explicit
caller kill. The original survival, Attach, Clean and journal assertions
remain unchanged.

The two new death tests re-execute the compiled test binary with their own
process groups, private `TMPDIR` directories and PID records. A test-only
record rendezvous keeps each inner binary alive with its fixture waiting;
the outer test sends `SIGKILL` to the inner binary's process group, confirms
that signal in its wait status, then proves both the fixture PID and its
process group absent before sending any cleanup signal to the fixture. Group
absence includes the Implement fake ACPX and its shell children. Every wait
uses the enclosing test deadline. Logs go into the outer temporary directory,
not inherited pipes that detached processes could hold open.

Focused evidence, with host process observation permitted:

- Acceptance criterion 1: `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -timeout=90s -v -run '^Test(ImplementDetachChild|DetachSurvivor)EndsWhenItsTestBinaryDies$' ./internal/cli`
  exited 0. `TestImplementDetachChildEndsWhenItsTestBinaryDies` passed in
  1.82s, proving the detached Owner PID and complete process group absent
  after its inner test binary was killed.
- Acceptance criterion 2: the same focused command reported
  `TestDetachSurvivorEndsWhenItsTestBinaryDies` passing in 0.29s. Its survivor
  was alive before the inner binary's kill, in a separate session, and its
  PID and group were absent afterward.
- Acceptance criterion 3: `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -timeout=120s -v -run '^(TestRunDetachedCommand.*|TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent|TestDetachedChildIsTerminatedAtTeardown|TestRunImplementDetachSurvivesCallerProcessGroupKill)$' ./internal/cli`
  exited 0. Both original survival tests passed, along with the adjacent
  liveness, Run-creation, handshake and refusal tests. The Implement survival
  test passed in 3.55s; the sentinel-ignoring teardown test passed in 0.01s.
- Old-wait failure proof: before adding either owner check, compiled the new
  death-test and rendezvous infrastructure with
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -c -o /tmp/roundfix-task03.test ./internal/cli`.
  From `internal/cli`, ran
  `rtk proxy /tmp/roundfix-task03.test -test.run='^Test(ImplementDetachChild|DetachSurvivor)EndsWhenItsTestBinaryDies$' -test.timeout=25s -test.v`.
  Both death tests failed at the deadline margin, at 24.01s and 24.03s,
  with `PID alive=true; process group alive=true`. The old release-file loop
  and hour-sleep loop were still in place for this probe. Outer cleanup
  reclaimed only the recorded fixture groups. The probe output is in
  `/tmp/roundfix-task03-red.log`.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk make verify-incremental` exited 0:
  formatting, vet, the full Go test suite, skill checks and build passed.
  `internal/cli` took 205.234s. Output is in
  `/tmp/roundfix-task03-incremental.log`.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The initial compiled-binary probe was launched from the repository root,
which gave the suite guard the wrong relative scan root; its known test
processes were stopped, and the evidence above comes from the correctly
located probe. Restricted process-table access required sandbox escalation
for the process checks.

Task status, Subtask and Acceptance Criteria checkboxes remain unchanged.
Authored Verification was not run; the Daemon owns Verification and
settlement. No other Task, Task Graph or production file was edited, and no
commit, push or pull request was made. No follow-up implementation was added.

## Carry-forward provenance

- Source Run: `run_20261002T064202Z_a22e1dd4644c4dc6`
- Source commit: `5ffe5e0ff1201df10f6916238a2118c3046916c8`
