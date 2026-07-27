---
task: task_04
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: backend
complexity: high
---

# Task 04: Retry transient Review Source failures and project waits

## Overview

Retry only positively typed transient Review Source failures within existing
timeout and Run Budget bounds while projecting every wait phase and retry
episode. Fake-clock coverage proves bounded recovery, exhaustion, and graceful
Stop Request interruption without real sleeps.

## Requirements

1. MUST retry only typed transient errors and never infer retryability from
   output text.
2. MUST reuse the configured poll interval, Review Source timeout, and Run
   Budget without a new retry setting.
3. MUST publish one retry-start event and one recovery or exhaustion event per
   episode.
4. MUST expose phase, expected head, start, deadline, Evidence, and retry status
   for both Review Source wait phases.
5. MUST publish progress only on phase, Evidence, or retry changes.
6. MUST let the Store-backed Stop Request interrupt retry sleep and win over
   another attempt.
7. MUST keep authentication, validation, and permanent failures terminal.

## Subtasks

- [ ] Add bounded transient retry episodes to watch.
- [ ] Reuse existing time and budget boundaries.
- [ ] Add wait-phase and deadline projection.
- [ ] Deduplicate unchanged progress and episode events.
- [ ] Integrate graceful Stop Request interruption.
- [ ] Add recovery, exhaustion, permanent, and cancellation matrices.

## Acceptance Criteria

- [ ] One transient failure retries and can recover without ending the Run.
- [ ] Exhaustion occurs at the existing timeout or Run Budget boundary.
- [ ] Each episode contains at most one start and one terminal episode event.
- [ ] Permanent failure performs zero retry sleeps.
- [ ] Both wait phases expose the documented head, deadline, Evidence, and retry
      fields.
- [ ] Unchanged polling appends no duplicate progress event.
- [ ] Stop Request during retry sleep starts no later Review Source call.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/watch -run 'TestRun.*(Transient|Retry|WaitPhase|StopRequest)' -count=1`
  — expected: fake-clock recovery, exhaustion, projection, deduplication, and
  stop cases pass.
- `rtk go test ./internal/cli ./internal/runevent -run 'Test.*(ReviewRetry|WaitingForReview|ReviewStatusEvent)' -count=1`
  — expected: CLI and stream expose bounded wait evidence.
- `rtk go test -race ./internal/watch ./internal/cli -run 'Test.*(Transient|Retry|StopRequest)' -count=1`
  — expected: retry and cancellation have no race or leaked waiter.

## References

- `_prd.md` → Goals 2–3; User Stories 2–3; Core Features 5–6; User Experience;
  Success Metrics.
- `_techspec.md` → API Contracts: Transient retry and Wait projection; Testing
  Approach; Build Order 4.
- `../0037-terminal-outcome-integrity/_techspec.md` → stop-aware Review Source
  waits.
