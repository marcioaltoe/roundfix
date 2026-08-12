---
task: task_03
spec: 0032-deterministic-agent-session-cancellation
status: completed
type: test
complexity: medium
---

# Task 03: Prove cooperative cancellation stability

## Overview

Lock the cooperative Agent Session cancellation path to its close-free contract and exercise both cancellation outcomes under repeated and race-enabled verification. This companion slice distinguishes a real lifecycle regression from scheduler load without weakening either path's assertions.

## Requirements

1. MUST make the cooperative scenario wait for completed cancel invocation and prompt termination through observable milestones.
2. MUST prove that cooperative completion returns `StopError.Killed=false` and never invokes or records Agent Session close.
3. MUST leave the controlled grace timer unfired when cancel ends the prompt within the grace period.
4. MUST run the cooperative and forced-close scenarios 100 times with no retries, tolerated failures, sleeps, or skipped assertions.
5. MUST pass the complete Agent package race gate and the repository verification gate.

## Subtasks

- [x] Rewrite the cooperative scenario around completed milestones.
- [x] Assert the exact cancel-only invocation sequence and absent close marker.
- [x] Assert prompt completion and the cooperative `StopError` result.
- [x] Exercise both cancellation outcomes repeatedly.
- [x] Run race-enabled Agent package verification.
- [x] Run the full repository verification gate.

## Acceptance Criteria

- [x] Cooperative cancellation records cancel and no Agent Session close invocation or marker.
- [x] The cooperative result is a `StopError` with `Killed=false`, and the blocked prompt terminates.
- [x] The first controlled grace timer remains unfired when the prompt exits cooperatively.
- [x] Both cancellation scenarios pass 100 independent executions in one command.
- [x] The complete Agent package reports no race, and the repository verification gate passes without retries.

## Context

- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=100` — expected: all 200 scenario executions pass with strict cancel/close assertions.
- `rtk go test -race ./internal/agent` — expected: the complete Agent package passes with no race report.
- `rtk make verify` — expected: formatting, all tests, skill synchronization, and the Roundfix build pass without retries.

## References

- `_prd.md` → Goals 1-3; User Stories 2-3; Core Features 5-6; Success Metrics; Non-Goals / Out of Scope.
- `_techspec.md` → API Contracts; Coverage Map; Testing Approach; Build Order 3-4; Risks & Considerations; Decisions.
- `docs/adr/0018-one-agent-session-per-run.md` → Agent Session lifecycle ownership.
- `docs/adr/0022-stop-requests-travel-through-the-run-database.md` → cooperative cancel and forced-stop boundary.

## Result

- Cooperative cancellation now waits for prompt start, completed cancel invocation, controlled grace timer creation, and prompt completion milestones before asserting the result.
- The cooperative scenario asserts the exact cancel-only invocation sequence and verifies that no Agent Session close completion marker exists.
- The cooperative result is asserted as `StopError` with `Killed=false`; the controlled grace timer is asserted unfired and stopped, with no post-close timer created.
- Evidence: `rtk go test ./internal/agent -run 'TestACPXRunCancelsPromptCooperatively$' -count=1` passed with 1 scenario.
- Evidence: `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=100` passed with 200 scenario executions.
- Evidence: `rtk go test -race ./internal/agent` passed with 143 tests and no race report.
- Evidence: `rtk make verify` passed formatting, all Go tests, skill synchronization, and build.
- Evidence: `rtk git diff --check` passed.
