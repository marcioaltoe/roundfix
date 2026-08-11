---
spec: 0032-deterministic-agent-session-cancellation
prd: _prd.md
created: 2026-07-16
---

# Deterministic Agent Session cancellation — Technical Spec

## Executive Summary

This bug fix separates cancellation grace timing from subprocess scheduling in the test without changing production semantics. `ACPXRunner` receives a small private timer factory used by `cancelPrompt`; production defaults to `time.NewTimer`, while tests hold and fire timers through channels after observing real fake-process milestones. The existing helper process remains the integration boundary for acpx command execution. The trade-off is one additional private dependency in the runner, accepted because it replaces a 20 ms scheduling guess with direct control of the lifecycle condition the test is meant to verify.

## System Architecture

Only `internal/agent` changes:

- `ACPXRunner` owns a private cancellation timer factory with a real-time default.
- `cancelPrompt` uses that factory for the cooperative grace and post-close wait while retaining context deadlines for actual acpx cancel and close subprocesses.
- The helper-process fixture exposes observable prompt, cancel, and close milestones through its existing temporary-file boundary.
- The two cancellation tests use a controllable timer and milestone waits. No CLI, Daemon, Run Database, Run Event, configuration, or public Agent interface changes.

The production call sequence remains context cancellation → acpx cancel → grace wait → acpx session close → bounded prompt-process wait → `StopError`.

## Implementation Design

### Interfaces

```go
type cancellationTimer interface {
    C() <-chan time.Time
    Stop() bool
}

type cancellationClock interface {
    NewTimer(time.Duration) cancellationTimer
}
```

`ACPXRunner` carries a private `cancelClock cancellationClock`. A helper returns a real implementation when the field is nil, matching the existing `eventClock` default pattern. The interface is package-private, used by production cancellation code, and exposes no test-only method on a public type.

The concrete real timer wraps `*time.Timer`. The fake clock records timers in creation order and lets a test wait for creation and fire each timer once. It has no wall-clock fallback inside the behavior assertion; the outer test deadline remains only a deadlock guard with an actionable failure.

### Data Models

No persisted or public data changes. The test fixture retains:

- prompt-started marker;
- cancel-completed marker written after the fake cancel command records its invocation;
- close-completed marker written after the fake close command records its invocation;
- append-only invocation JSON used for final order assertions.

The marker names describe completed observable milestones, not inferred timing. Every test uses a unique temporary directory.

### API Contracts

`cancelPrompt` preserves these behaviors:

1. Build and run the acpx cancel command before creating the cooperative grace timer.
2. If the prompt process exits before the timer fires, return `false` and do not invoke session close.
3. If the timer fires first, build and run session close, then wait for the prompt or the post-close timer, and return `true`.
4. Cancel and close command construction or execution failures remain warnings and do not skip the subsequent lifecycle step.
5. `RunPrompt` continues mapping the boolean to `StopError.Killed` and publishing the stopped Agent status.

The test uses a command timeout comfortably larger than local process startup. It does not wait that duration: after the cancel command completes and the fake clock reports the first timer, the test fires the timer immediately to trigger close. The close marker and invocation log prove the fallback ran. The blocked fake prompt exits only after observing the close marker, which completes the production wait channel without a guessed delay.

For cooperative cancellation, the fake prompt exits after the cancel marker. The test waits for the cancel milestone, allows the prompt to exit, and asserts that the first fake timer is never fired and no close invocation or marker exists.

The prior 20 ms fixture coupled three independent budgets: cancel subprocess startup, cooperative grace, and test runtime. Under load, the cancel subprocess could lose the scheduling race and be terminated by its context before `appendFakeACPXInvocation`, producing the observed missing-cancel assertion even though fallback close ran. The new fixture controls the grace condition after process completion and leaves the production default of 10 seconds unchanged.

## Coverage Map

- Goal 1 and Story 1 → exact invocation-order assertion over the helper-process log.
- Goal 2 and Story 2 → fake cancellation clock and completed-milestone synchronization.
- Goal 3 and Story 3 → unchanged `cancelPrompt` state machine and existing `StopError` expectations.
- Core Feature 5 → cooperative cancel companion test with close absence.
- Core Feature 6 → focused 100-run stress and race commands plus the full verification gate.

## Integration Points

- Go `os/exec` helper-process test boundary, which continues to exercise real subprocess creation, context cancellation, argument construction, and exit handling.
- acpx command grammar remains the existing cancel and `sessions close` contract.
- `context.Context` remains the owner of cancel/close subprocess deadlines; the injected clock controls only grace waiting.

## Cross-Spec dependencies

None. This Spec stabilizes the Agent Session lifecycle verification boundary and must be implemented first so later Agent and CLI work runs on a deterministic gate.

## Testing Approach

Keep the tests in the canonical `internal/agent/acpx_runner_test.go` suite. Replace `TestACPXRunClosesSessionAfterCancelGracePeriod` with a deterministic forced-close scenario and retain `TestACPXRunCancelsPromptCooperatively` as the negative companion. Assertions cover exact command order, marker order, `StopError` type and `Killed` value, prompt completion, and absence of unexpected close.

The fake clock itself gets table tests for timer creation order, single fire, and stop behavior. Those tests do not duplicate `time.Timer`; they protect only the fixture contract needed by cancellation tests.

Required verification:

```text
rtk go test ./internal/agent -run 'TestACPXRun(CancelsPromptCooperatively|ClosesSessionAfterCancelGracePeriod)$' -count=100
rtk go test -race ./internal/agent
rtk make verify
```

No command uses test retries. `-count=100` creates independent executions and any single failure fails the command.

## Build Order

1. Deterministic fake cancellation clock and helper-process milestone protocol, with a regression test that cannot advance grace before cancel completion.
2. Private real cancellation clock and `cancelPrompt` timer-factory integration (depends on: 1).
3. Rewrite forced-close and cooperative cancellation tests around milestones and exact ordering (depends on: 1, 2).
4. Focused 100-run stress, Agent package race test, and full repository verification (depends on: 3).

## Risks & Considerations

- A fake clock can make a test pass while real subprocess contexts still race. The helper process remains real, and only grace timers are controlled; cancel and close commands still execute through `os/exec` with context deadlines.
- A long command timeout could hide a hung helper. Every subprocess remains context-bounded, and the test's outer deadlock guard names the missing milestone.
- A broad clock package would add surface unrelated to this fix. The interface stays private beside `ACPXRunner`.
- Firing the fake timer before cancel completion would reproduce the same invalid test. Timer creation occurs only after the synchronous cancel command returns, and the test waits for both the marker and created timer.
- If stress or race verification reveals an actual production ordering defect, implementation must return to root-cause investigation and extend this Spec before changing behavior; assertion weakening and retries remain forbidden.

## Decisions

- The production cancellation state machine and 10-second default remain unchanged.
- Grace timers are injectable at the runner boundary; subprocess deadlines remain context-owned and real.
- Tests synchronize on completed fake-process milestones and timer creation, never fixed sleeps.
- The full Agent package race test is required because cancellation owns goroutines, channels, subprocesses, and mutable fixture state.
