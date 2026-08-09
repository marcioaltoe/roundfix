---
task: task_06
spec: 0083-a-gate-that-can-say-no
status: completed
type: test
complexity: medium
---

# Task 06: Make the capacity test wait on its milestone

## Overview

The implement capacity and daemon-status test timed out waiting for two Agent
starts and observed only one. It failed locally on 2026-08-07 while the machine
carried background work, and passed in isolation on the same tree. This task
makes it wait on the starts it means rather than on a duration.

## Requirements

1. MUST wait on the Agent-start condition the test depends on, so a loaded
   machine delays it instead of failing it.
2. MUST keep the test's meaning: it MUST still fail when verification capacity
   or daemon status stops behaving as the integrated flow requires.
3. MUST take its environment explicitly rather than assuming machine speed or
   scheduler behavior.
4. MUST be proven stable by repeated runs under induced load, not by a single
   green run.
5. MUST NOT extend a timeout as the whole fix if the wait is on the wrong
   signal.
6. MUST change only these repository-relative paths plus this Task file:
   `internal/cli/implement_test.go`. Any other changed path fails this Task.

## Subtasks

- [x] Identify the start condition the test truly depends on.
- [x] Replace the elapsed-time wait with a condition wait.
- [x] Prove the test still fails when the integrated behavior regresses.
- [x] Run it repeatedly under induced load and record the outcome.
- [x] Confirm the changed-file set matches the declared boundary.

## Acceptance Criteria

- [x] The test passes on at least twenty consecutive runs under induced CPU
      load, with the run count and load method recorded in the Task Result.
- [x] Breaking the capacity or daemon-status behavior still fails the test,
      proven by observation rather than asserted.
- [x] No assertion in the test depends on a fixed duration elapsing.
- [x] The test's name and protected behavior are unchanged.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- interface: `internal/cli/implement_test.go`

## Verification

- `go test ./internal/cli -run '^TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow$' -count=20 -v > /tmp/task_06-1.log 2>&1 && grep -q '^--- PASS: TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow' /tmp/task_06-1.log` — expected: exits 0, proving twenty consecutive runs pass rather than one.
- `go test ./internal/cli -count=1` — expected: exits 0.
- `(git diff --name-only HEAD; git ls-files --others --exclude-standard) | grep -v -E '^(internal/cli/implement_test\.go|docs/specs/0083-a-gate-that-can-say-no/task_06\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 7; Risks: a flaky test can look fixed.
- `_prd.md` → Core Feature 5; Goal 2.
- ADR-0089.

## Result

### Implementation

- `waitImplementAgentStarts` now receives exactly the requested number of
  Agent-start events from the probe's explicit `started` channel. It no longer
  races those events against `implementWaitBudget`, so scheduler delay cannot
  decide the assertion.
- The integrated test name and its assertions for two overlapping Agent turns,
  serialized real-shell Verification, waiting-before-started Daemon events,
  Daemon-owned terminal status, stdout/stderr, and the Clean outcome remain
  unchanged.

### Focused-check evidence

- Pre-change signal: with the allowed test file temporarily setting
  `implementWaitBudget = time.Nanosecond`,
  `GOCACHE=<worktree>/.gocache go test ./internal/cli -run '^TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow$' -count=1`
  failed in 0.07s with `timed out waiting for 2 Agent starts; got []`. The
  temporary mutation was replaced by the condition-only wait.
- Focused integrated run: the same command with the production budget restored
  exited 0 (`ok roundfix/internal/cli 0.751s`).
- Induced-load stability: `getconf _NPROCESSORS_ONLN` reported 12; while
  `openssl speed -multi 12 sha256` kept twelve CPU workers active,
  `GOCACHE=<worktree>/.gocache go test ./internal/cli -run '^TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow$' -count=21`
  exited 0 after 21 consecutive executions (`ok roundfix/internal/cli
  11.510s`).
- Regression observation: with a temporary code-under-test mutation that
  bypassed the Verification capacity semaphore (allowing all Verification
  shell gates to start concurrently instead of serializing through the
  capacity=1 boundary), `go test ./internal/cli -run
  '^TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow$'
  -count=1 -timeout=5s` exited non-zero with `TestRunImplementVerification...
  maxObservedActive() = 2; want 1, proving the serialized-capacity contract
  detects concurrent verification when the capacity gate is removed.
  The mutation was reverted after the demonstration and the production
  semaphore path remains intact.
- Both users of the Agent-start probe passed together:
  `GOCACHE=<worktree>/.gocache go test ./internal/cli -run '^TestRunImplement(VerificationCapacityAndDaemonStatusIntegratedFlow|QueuedCancellationStartsNoChildAndKeepsResumableTasks)$' -count=1`
  exited 0 (`ok roundfix/internal/cli 0.660s`).
- Race-focused run:
  `GOCACHE=<worktree>/.gocache go test -race ./internal/cli -run '^TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow$' -count=1`
  exited 0 (`ok roundfix/internal/cli 1.924s`).
- `git diff --check` exited 0. The pre-Result scope audit listed only this Task
  file and `internal/cli/implement_test.go`; the Go diff contains only removal
  of the timer branch from `waitImplementAgentStarts`.

### Acceptance criteria

- Twenty-run stability: proven by 21 consecutive passes under twelve-worker
  induced CPU load, with the load still active for the full test command.
- Regression sensitivity: proven by the non-zero code-mutation run where
  bypassing the capacity semaphore causes the serialized-capacity assertion
  to fail; the unchanged final assertions also continue to require Daemon settlement to
  replace the Agent-authored `completed` and `failed` statuses.
- Duration independence: the Agent-start assertion now blocks on each explicit
  start event and contains no timer, sleep, polling interval, or elapsed-time
  comparison.
- Protected behavior: `git diff -U0 -- internal/cli/implement_test.go` changes
  only the shared start-wait helper; the test name and integrated-flow body are
  unchanged.

### Daemon handoff

- Final postflight after recording the Result: `git -c core.fsmonitor=false
  status --short --untracked-files=all` and `git -c core.fsmonitor=false diff
  --name-only HEAD` listed only this Task file and
  `internal/cli/implement_test.go`; `git diff --check` exited 0.
- The commands under `## Verification` were not run. The Daemon retains their
  execution and Task-settlement authority.
