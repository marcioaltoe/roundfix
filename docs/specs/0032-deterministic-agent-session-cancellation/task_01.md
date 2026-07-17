---
task: task_01
spec: 0032-deterministic-agent-session-cancellation
status: pending
type: test
complexity: medium
---

# Task 01: Create deterministic cancellation controls

## Overview

Create the test-only controls that let Agent Session cancellation tests observe completed subprocess milestones and advance cancellation timers explicitly. This prefactor leaves production timing unchanged while making the fixture contract independently verifiable before lifecycle behavior depends on it.

## Requirements

1. MUST provide a controllable cancellation clock fixture that records timer creation order, lets each timer fire at most once, and exposes whether a timer was stopped.
2. MUST expose completed prompt-start, cancel, and close milestones through the helper-process boundary described by the TechSpec.
3. MUST isolate every invocation log and milestone set in the test's temporary directory.
4. MUST use an outer deadline only as a deadlock guard that identifies the missing milestone; it MUST NOT advance cancellation behavior from elapsed wall-clock time.
5. MUST retain the existing helper-process coverage of real subprocess creation and acpx argument recording.

## Subtasks

- [ ] Add the controllable cancellation clock and timer fixtures.
- [ ] Cover timer creation order, single fire, and stop behavior.
- [ ] Record completed cancellation milestones from the helper process.
- [ ] Add milestone waits with actionable deadlock failures.
- [ ] Keep the existing cancellation scenarios passing through the prefactored harness.

## Acceptance Criteria

- [ ] Tests can wait for a created timer without firing it or depending on a short sleep.
- [ ] A timer cannot emit more than one event, and a stopped timer cannot be fired as active.
- [ ] Prompt start, cancel completion, and close completion are distinguishable observable milestones.
- [ ] Cancellation fixtures use unique temporary paths and retain the append-only invocation record.
- [ ] Existing Agent package cancellation tests pass after the test-harness prefactor.

## Context

- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `rtk go test ./internal/agent -run 'TestFakeCancellationClock' -count=1` — expected: the controllable clock's creation order, single-fire, and stop contract pass.
- `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=1` — expected: both existing cancellation scenarios remain green through the prefactored helper process.

## References

- `_prd.md` → Goals 1-2; User Stories 1-2; Core Features 2 and 6; Non-Goals / Out of Scope; Decisions.
- `_techspec.md` → Interfaces; Data Models; Integration Points; Testing Approach; Build Order 1; Risks & Considerations.
