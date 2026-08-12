---
task: task_02
spec: 0032-deterministic-agent-session-cancellation
status: completed
type: backend
complexity: high
---

# Task 02: Enforce deterministic forced-close ordering

## Overview

Move cancellation grace waits behind the runner-owned clock boundary and prove the forced-close path through real helper subprocesses. The slice must demonstrate cooperative cancel before Agent Session close while preserving production defaults, command deadlines, warnings, and `StopError` behavior.

## Requirements

1. MUST add the private cancellation timer and clock boundary described by the TechSpec, with real timers as the nil-safe production default.
2. MUST use that clock for the cooperative grace and post-close wait without replacing context-owned cancel or close subprocess deadlines.
3. MUST complete the cancel command before creating the cooperative grace timer.
4. MUST make the forced-close test advance grace only after cancel completion and timer creation are both observable.
5. MUST prove exact cancel-before-close invocation order, completion of the blocked prompt, and `StopError.Killed=true`.
6. MUST preserve the 10-second production default, public interfaces, acpx command grammar, warning behavior, and terminal status publication.

## Subtasks

- [x] Add the private runner-owned cancellation clock boundary and real implementation.
- [x] Route both `cancelPrompt` grace waits through the boundary.
- [x] Synchronize forced-close setup on prompt start and completed cancel invocation.
- [x] Fire the controlled grace timer and observe completed close invocation.
- [x] Assert the exact invocation order, prompt termination, and forced-stop result.
- [x] Retain coverage for cancel and close construction or execution warnings.

## Acceptance Criteria

- [x] A runner with no injected clock uses real timers and retains the existing production grace value.
- [x] The forced-close scenario cannot create or fire its first grace timer before the cancel subprocess completes.
- [x] The invocation record places cooperative cancel before Agent Session close with no tolerated alternative order.
- [x] The close milestone releases the blocked prompt, and the Run returns a `StopError` with `Killed=true`.
- [x] Cancel or close command failures remain warnings and do not skip the next lifecycle step.
- [x] The forced-close scenario passes 100 independent executions without retries or timing sleeps.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/agent.go`

## Verification

- `rtk go test ./internal/agent -run 'TestACPXRunClosesSessionAfterCancelGracePeriod$' -count=100` — expected: every execution records cancel before close and returns the forced-stop result.
- `rtk go test -race ./internal/agent -run 'TestACPXRunClosesSessionAfterCancelGracePeriod$' -count=1` — expected: the forced-close lifecycle completes with no race report.

## References

- `_prd.md` → Goals 1 and 3; User Stories 1 and 3; Core Features 1, 3, and 4; Success Metrics; Decisions.
- `_techspec.md` → System Architecture; Interfaces; API Contracts 1-5; Integration Points; Testing Approach; Build Order 2-3; Risks & Considerations.
- `docs/adr/0018-one-agent-session-per-run.md` → Agent Session lifecycle ownership.
- `docs/adr/0022-stop-requests-travel-through-the-run-database.md` → cooperative cancel and forced-stop boundary.

## Result

Implemented the deterministic forced-close lifecycle slice:

- Added private `cancellationClock` and `cancellationTimer` boundaries to `ACPXRunner`, with nil-safe real-time defaults backed by `time.NewTimer`.
- Routed the cooperative grace wait and post-close wait in `cancelPrompt` through the runner-owned cancellation clock while keeping cancel and close subprocess deadlines on real `context.WithTimeout` timers.
- Updated the forced-close scenario to wait for prompt start and completed cancel, wait for the first created grace timer, fire it explicitly, then observe completed close.
- Asserted the exact helper-process invocation order: session ensure, reasoning setup, prompt, cooperative cancel, then Agent Session close.
- Added warning-path coverage proving cancel and close command failures are warnings and do not prevent the fallback close lifecycle from completing.

Evidence:

- Pre-change signal: `rtk rg "cancelClock|cancellationClock|cancellationTimer|realCancellationClock" internal/agent/acpx_runner.go` found no runner-owned cancellation clock seam, while the old `TestACPXRunClosesSessionAfterCancelGracePeriod` still passed only through elapsed wall-clock time.
- `rtk gofmt -w internal/agent/acpx_runner.go internal/agent/acpx_runner_test.go` — passed.
- `rtk go test ./internal/agent -run 'TestACPXRunnerCancellationClockDefaultsToRealTimer|TestFakeCancellationClock' -count=1` — passed; output: `Go test: 5 passed in 1 packages`.
- `rtk go test ./internal/agent -run 'TestACPXRun(ClosesSessionAfterCancelGracePeriod|CancellationCommandFailuresWarnAndContinue)$' -count=1` — passed; output: `Go test: 4 passed in 1 packages`.
- `rtk go test ./internal/agent -run 'TestACPXRunClosesSessionAfterCancelGracePeriod$' -count=100` — passed; output: `Go test: 100 passed in 1 packages`.
- `rtk go test -race ./internal/agent -run 'TestACPXRunClosesSessionAfterCancelGracePeriod$' -count=1` — passed; output: `Go test: 1 passed in 1 packages`.
- `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=1` — passed; output: `Go test: 2 passed in 1 packages`.
- `rtk git diff --check` — passed.
- `rtk make verify` — passed; output included `Go test: 1280 passed in 19 packages`, `Roundfix skill check passed`, and a successful `go build`.

Acceptance evidence:

- Nil clock/default grace: `TestACPXRunnerCancellationClockDefaultsToRealTimer` verifies no injected clock uses a real timer and `stopGrace(0)` remains `10s`.
- No early grace advancement: forced-close waits for completed cancel and then for timer creation before calling `Fire`.
- Strict cancel-before-close order: forced-close compares the full invocation JSONL sequence with no alternate accepted ordering.
- Close releases prompt and returns forced stop: the close milestone is observed, the post-close fake timer is stopped without firing, and the returned `StopError` has `Killed=true`.
- Warning behavior: `TestACPXRunCancellationCommandFailuresWarnAndContinue` covers cancel and close execution failures as warnings while still reaching fallback close and prompt termination.
- Stress/race stability: the forced-close test passed 100 independent runs and the race-enabled focused command.
