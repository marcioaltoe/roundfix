---
task: task_05
spec: 0083-a-gate-that-can-say-no
status: completed
type: test
complexity: medium
---

# Task 05: Make the cancellation test wait on its milestone

## Overview

The ACP cancellation test times out waiting for a fake adapter's prompt-start
milestone. On 2026-08-07 it failed in CI on a documentation-only pull request
and passed on re-run of the same commit — a timing flake, not a defect in the
code under test. This task makes it wait on the condition it means rather than
on a duration, so a slow runner delays it instead of failing it.

## Requirements

1. MUST wait on the milestone the test actually depends on, with a bound
   generous enough that a loaded runner delays rather than fails, or with a
   signal that does not depend on elapsed time at all.
2. MUST keep the test's meaning: it MUST still fail when close failure does not
   wait for prompt termination, which is the behavior it exists to protect.
3. MUST take its environment explicitly rather than assuming runner speed,
   following the repository's accepted rule for code under test.
4. MUST be proven stable by repeated runs under induced load, not by a single
   green run — a single pass is precisely the evidence that misled this
   repository already.
5. MUST NOT extend a timeout as the whole fix if the wait is on the wrong
   signal; a longer sleep is a slower flake.
6. MUST change only these repository-relative paths plus this Task file:
   `internal/agent/acpx_runner_test.go`. Any other changed path fails this Task.

## Subtasks

- [x] Identify the milestone the test truly depends on.
- [x] Replace the elapsed-time wait with a condition wait.
- [x] Prove the test still fails when the protected behavior regresses.
- [x] Run it repeatedly under induced load and record the outcome.
- [x] Confirm the changed-file set matches the declared boundary.

## Acceptance Criteria

- [x] The test passes on at least twenty consecutive runs under induced CPU
      load, with the run count and load method recorded in the Task Result.
- [x] Removing the close-failure wait from the code under test still fails the
      test, proven by observation rather than asserted.
- [x] No assertion in the test depends on a fixed duration elapsing.
- [x] The test's name and protected behavior are unchanged.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `go test ./internal/agent -run '^TestACPXRunCancellationCommandFailuresWarnAndContinue$' -count=20 -v > /tmp/task_05-1.log 2>&1 && grep -q '^--- PASS: TestACPXRunCancellationCommandFailuresWarnAndContinue' /tmp/task_05-1.log` — expected: exits 0, proving twenty consecutive runs pass rather than one.
- `go test ./internal/agent -count=1` — expected: exits 0.
- `git diff --name-only HEAD | grep -v -E '^(internal/agent/acpx_runner_test\.go|docs/specs/0083-a-gate-that-can-say-no/task_05\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 7; Risks: a flaky test can look fixed.
- `_prd.md` → Core Feature 5; Goal 2.
- ADR-0089.

## Result

The cancellation-failure test now receives a prompt-start Run Event through
the fake ACPX environment and blocks on the capture sink's event channel. The
fake emits that event only after the prompt process reaches its start
milestone, so cancellation no longer depends on the package's wall-clock wait
budget. The test name, close-failure scenario, fake cancellation timers, and
prompt-termination assertion remain unchanged.

Focused evidence:

- Twenty-five consecutive runs passed while
  `openssl speed -seconds 15 -multi 24 -evp sha256 -bytes 16` ran concurrently
  with 24 workers on a 12-logical-CPU host. The focused command was
  `go test ./internal/agent -run '^TestACPXRunCancellationCommandFailuresWarnAndContinue$' -count=25 -v`;
  it exited 0 with all 25 parent tests and all 50 subtests passing in 8.961s.
  The load process also exited 0.
- A temporary mutation returned immediately from `cancelPrompt` when
  `sessions close` failed, removing the protected close-failure wait. The
  focused close-failure subtest exited 1 under `-timeout=5s`; its timeout stack
  showed it waiting for cancellation timer 2 in `fakeCancellationClock.waitForTimer`.
  The mutation was then removed, and
  `git diff --exit-code -- internal/agent/acpx_runner.go` exited 0.
- `go test -race ./internal/agent -run '^TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod|CancellationCommandFailuresWarnAndContinue)$' -count=1`
  exited 0, covering the event handshake beside the two adjacent cancellation
  scenarios under the race detector.
- Source inspection confirms the target test waits on `<-promptStarted.done`;
  it contains no `time.Sleep`, `time.After`, or real timer assertion. Its
  existing fake-clock assertions still advance behavioral timers explicitly.
- `git -c core.fsmonitor=false status --short --untracked-files=all` listed
  only `internal/agent/acpx_runner_test.go` and this Task file. The Task's
  declared Verification commands were not run; the Daemon owns them.
