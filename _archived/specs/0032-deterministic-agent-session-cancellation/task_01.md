---
task: task_01
spec: 0032-deterministic-agent-session-cancellation
status: completed
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

- [x] Add the controllable cancellation clock and timer fixtures.
- [x] Cover timer creation order, single fire, and stop behavior.
- [x] Record completed cancellation milestones from the helper process.
- [x] Add milestone waits with actionable deadlock failures.
- [x] Keep the existing cancellation scenarios passing through the prefactored harness.

## Acceptance Criteria

- [x] Tests can wait for a created timer without firing it or depending on a short sleep.
- [x] A timer cannot emit more than one event, and a stopped timer cannot be fired as active.
- [x] Prompt start, cancel completion, and close completion are distinguishable observable milestones.
- [x] Cancellation fixtures use unique temporary paths and retain the append-only invocation record.
- [x] Existing Agent package cancellation tests pass after the test-harness prefactor.

## Context

- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `rtk go test ./internal/agent -run 'TestFakeCancellationClock' -count=1` — expected: the controllable clock's creation order, single-fire, and stop contract pass.
- `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=1` — expected: both existing cancellation scenarios remain green through the prefactored helper process.

## References

- `_prd.md` → Goals 1-2; User Stories 1-2; Core Features 2 and 6; Non-Goals / Out of Scope; Decisions.
- `_techspec.md` → Interfaces; Data Models; Integration Points; Testing Approach; Build Order 1; Risks & Considerations.

## Result

Implemented the deterministic cancellation test controls for this slice:

- Added `fakeCancellationClock` and `fakeCancellationTimer` fixtures that record timer creation order, expose created timers without firing them, fire each timer at most once, and expose stopped state.
- Moved blocking fake-acpx cancellation markers into an isolated per-test milestone set under the harness temporary directory: prompt start, cancel completion, and close completion.
- Updated the existing cooperative and forced-close cancellation scenarios to wait on named milestone files with deadlock failures that identify the missing milestone and include invocations recorded so far.
- Kept the helper-process boundary and append-only `invocations.jsonl` record intact; cancellation scenarios still execute real subprocesses through the fake acpx process.

Evidence:

- `rtk gofmt -w internal/agent/acpx_runner_test.go` — passed.
- `rtk go test ./internal/agent -run 'TestFakeCancellationClock' -count=1` — passed; output: `Go test: 4 passed in 1 packages`.
- `rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=1` — passed; output: `Go test: 2 passed in 1 packages`.

Acceptance evidence:

- Created-timer wait without sleep: covered by `TestFakeCancellationClock/records creation order and waits without firing`.
- Single fire and stopped timer behavior: covered by `TestFakeCancellationClock/fires each timer at most once` and `TestFakeCancellationClock/stopped timer cannot be fired as active`.
- Distinguishable prompt/cancel/close milestones: the cancellation scenarios now wait for separate prompt-start, cancel-completed, and close-completed paths.
- Temporary path isolation and append-only invocation record: `assertCancellationFixturePaths` validates unique paths under the harness temp directory, while the unchanged fake acpx invocation log remains JSONL append-only.
- Existing cancellation scenarios: both existing Agent package cancellation tests passed through the prefactored harness.
