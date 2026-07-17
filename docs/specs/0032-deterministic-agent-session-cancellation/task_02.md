---
task: task_02
spec: 0032-deterministic-agent-session-cancellation
status: pending
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

- [ ] Add the private runner-owned cancellation clock boundary and real implementation.
- [ ] Route both `cancelPrompt` grace waits through the boundary.
- [ ] Synchronize forced-close setup on prompt start and completed cancel invocation.
- [ ] Fire the controlled grace timer and observe completed close invocation.
- [ ] Assert the exact invocation order, prompt termination, and forced-stop result.
- [ ] Retain coverage for cancel and close construction or execution warnings.

## Acceptance Criteria

- [ ] A runner with no injected clock uses real timers and retains the existing production grace value.
- [ ] The forced-close scenario cannot create or fire its first grace timer before the cancel subprocess completes.
- [ ] The invocation record places cooperative cancel before Agent Session close with no tolerated alternative order.
- [ ] The close milestone releases the blocked prompt, and the Run returns a `StopError` with `Killed=true`.
- [ ] Cancel or close command failures remain warnings and do not skip the next lifecycle step.
- [ ] The forced-close scenario passes 100 independent executions without retries or timing sleeps.

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
