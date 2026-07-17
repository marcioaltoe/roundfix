---
spec: 0032-deterministic-agent-session-cancellation
status: archived
created: 2026-07-16
surfaces: [backend]
archived: "2026-07-16"
source_slug: 0032-deterministic-agent-session-cancellation
---


# Deterministic Agent Session cancellation

The Agent Session cancellation test failed once because its 20 ms Stop Request grace also bounded the fake cancel subprocess, allowing scheduler delay to kill that process before the helper recorded the invocation. The production default remains 10 seconds and the failure did not reproduce, but the test cannot currently prove cancel-before-close ordering without depending on wall-clock scheduling. The cancellation boundary needs deterministic lifecycle control so the test reports a real ordering regression instead of machine load.

## Goals

- Prove that cancellation attempts cooperative Agent Session cancel before forced session close.
- Remove arbitrary short wall-clock deadlines from the cancellation-order test.
- Preserve the production Stop Request, grace-period, fallback-close, `StopError`, and exit contracts.
- Make repeated and race-enabled verification stable enough to distinguish a test-harness defect from a product race.

## User Stories

1. As a maintainer changing Agent Session lifecycle code, I want a deterministic cancel-before-close test, so that a failure identifies a real ordering regression.
2. As a developer running the full verification gate under load, I want cancellation tests independent of scheduler speed, so that an unrelated Spec does not fail intermittently.
3. As a developer stopping a Run, I want the existing cooperative-cancel and fallback-close behavior preserved, so that test repair does not weaken production cleanup.

## Core Features

1. Cancellation timing is represented by an owned, injectable clock boundary with the current real-time implementation as the production default.
2. The helper-process test synchronizes on prompt start, completed cancel invocation, grace expiration, completed close invocation, and prompt termination through observable events rather than sleeps.
3. The forced-close test uses a realistic command timeout and advances only the injected grace timer, so process scheduling cannot consume the complete test budget.
4. The test asserts exact cancel-before-close order, `StopError.Killed=true`, and completion of the blocked prompt after close.
5. The cooperative path remains covered and proves close is absent when cancel ends the prompt within the grace period.
6. Stress and race verification run both cancellation paths repeatedly without retries, skipped assertions, enlarged sleeps, or tolerated failures.

## Non-Goals / Out of Scope

- Changing the production default Stop Request grace period.
- Retrying failed cancellation tests or weakening their assertions.
- Adding sleeps, scheduler yields, or polling as substitutes for lifecycle events.
- Changing acpx command grammar, Agent Session identity, Stop Command behavior, or terminal Run outcomes.
- Adding exported test-only methods or a general-purpose clock package.
- Treating a test-harness fix as evidence that every external acpx process race is impossible.

## Success Metrics

- The forced-close test records cooperative cancel before session close in every repeated run.
- The cooperative test records cancel and no session close in every repeated run.
- Both tests pass under `go test -race` and a 100-run focused stress command.
- The full repository verification gate passes without retries.
- Existing Stop Request, `StopError`, acpx invocation, and Agent Session lifecycle tests retain their behavioral expectations.

## Decisions

- Fix the timing source, not the assertion: cancellation ordering remains strict.
- A private cancellation clock belongs to the existing runner lifecycle boundary; no global clock abstraction or exported test hook is introduced.
- Production behavior changes only if the deterministic tests expose a reproducible runtime ordering defect.

## Open Questions

None.
