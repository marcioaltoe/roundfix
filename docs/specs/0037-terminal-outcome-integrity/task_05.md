---
task: task_05
spec: 0037-terminal-outcome-integrity
status: pending
type: backend
complexity: high
---

# Task 05: Publish only the winning terminal outcome

## Overview

Carry the completion transition result through all terminal publication paths
so only the winner emits outcome evidence and notification. Preserve primary
failure ordering while retaining registered-session cleanup problems as
explicit secondary diagnostics.

## Requirements

1. MUST publish a terminal outcome Run Event only when completion transitioned
   the Run.
2. MUST attempt Run Outcome Notification only for that same winning transition.
3. MUST make identical replays and conflicting losers observe stored state
   without duplicate publication.
4. MUST preserve one terminal outcome across Resolve, Watch, Implement, Stop,
   and cleanup paths.
5. MUST print and journal the primary failure before any cleanup warning.
6. MUST label cleanup failures as secondary and prevent them from replacing
   the terminal reason or exit code.
7. MUST retain best-effort notification failure without changing the Run
   outcome.

## Subtasks

- [ ] Thread the completion transition result through terminal coordinators.
- [ ] Gate outcome events and notifications on the winning transition.
- [ ] Normalize replay and conflict handling around stored state.
- [ ] Order primary failure evidence before cleanup diagnostics.
- [ ] Cover terminal paths across operational commands.
- [ ] Add deterministic owner-versus-Force-Stop publication coverage.

## Acceptance Criteria

- [ ] A completion race stores and publishes exactly one terminal outcome.
- [ ] The losing owner emits no second outcome event or notification.
- [ ] An identical replay is silent and returns the stored result.
- [ ] A conflicting loser cannot change the primary reason, exit code, or
      notification context.
- [ ] Cleanup warnings follow the primary failure and are labeled secondary.
- [ ] Notification failure remains a warning and leaves persisted outcome
      unchanged.
- [ ] Resolve, Watch, Implement, and Stop terminal regressions pass.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/notify/notify.go`
- interface: `internal/notify/notify_test.go`
- interface: `internal/runevent/event.go`

## Verification

- `rtk go test ./internal/cli ./internal/store ./internal/notify ./internal/runevent -run 'Test.*(CompletionWinner|TerminalOutcome|PrimaryFailure|OutcomeNotification)' -count=1`
  — expected: only the completion winner publishes and primary diagnostics
  precede secondary cleanup warnings.
- `rtk go test -race ./internal/cli ./internal/store -run 'Test.*(CompletionWinner|TerminalOutcome)' -count=1`
  — expected: concurrent terminal publication has one race-free winner.

## References

- `_prd.md` → Goals 1 and 4; User Stories 3–4; Core Features 6–7; User
  Experience; Success Metrics.
- `_techspec.md` → System Architecture: Run Event Journal; Testing Approach;
  Build Order 5.
- `../../adr/0052-run-completion-is-compare-and-set.md` → winner-only
  publication.
