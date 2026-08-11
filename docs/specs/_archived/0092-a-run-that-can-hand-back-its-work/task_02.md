---
task: task_02
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
type: backend
complexity: high
---

# Task 02: Move the work-started boundary to the first Agent output

## Overview

The Fallback Chain exists so a Run survives a runtime that cannot serve it. It
did not activate when the Codex quota was exhausted, because the work-started
signal is published after the session is prepared and before the first prompt, so
an adapter that printed a usage limit and exited nineteen seconds in was already
past the boundary. This Task moves the boundary to the first Agent output.

## Requirements

1. MUST publish the work-started status when the Agent produces its first
   output, not when its Session is prepared or activated.
2. MUST publish a selection failure, distinct from a Batch failure, when a turn
   ends with no Agent output and an adapter-level refusal.
3. MUST leave a Fallback Selection eligible after a selection failure, so the
   configured chain is attempted.
4. MUST keep a Fallback Selection ineligible once any Agent output exists; this
   Task narrows when the chain is ineligible and widens nothing it may do.
5. MUST NOT change what a Fallback Chain contains, how it is configured, or
   ADR-0050's rule that preflight substitutes nothing.
6. MUST break the characterization case Task 01 declared for the boundary, and
   update it in the same commit.

## Subtasks

- [ ] Move the signal to the first Agent output.
- [ ] Classify a no-output adapter refusal as a selection failure.
- [ ] Keep the chain eligible for that class only.

## Acceptance Criteria

- [ ] A turn that produces no Agent output and fails with an adapter refusal
      publishes a selection failure and no work-started status.
- [ ] The configured Fallback Selection is attempted after that failure.
- [ ] A turn that produced Agent output before failing keeps the chain
      ineligible.
- [ ] The work-started status is published exactly once per Session.

## Rehearsal Cases

- Case: an adapter that exits immediately reporting a usage limit, with a
  configured `codex → claude` chain; Observation: a selection failure is
  published, no work-started status appears, and the `claude` selection is
  attempted.
- Case: an adapter that emits one output chunk and then exits non-zero;
  Observation: work-started is published, and the chain is not attempted.
- Case: an adapter whose command cannot start at all; Observation: a selection
  failure, and the chain is attempted.

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner.go`
- `internal/agent/agent.go`
- `internal/agent/acpx_runner_test.go`
- `internal/daemon/agent_session_owner.go`
- `internal/daemon/agent_session_owner_test.go`
- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_02.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestWorkStartedBoundary' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestWorkStartedBoundaryPublishesOnFirstAgentOutput'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestWorkStartedBoundary' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestWorkStartedBoundaryReportsSelectionFailureWithoutOutput'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestFallbackEligibility' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestFallbackEligibilitySurvivesASelectionFailure'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestFallbackEligibility' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestFallbackEligibilityEndsAfterAnyAgentOutput'` — expected: exits 0.

## References

- `_prd.md` → Goal 1.
- `_techspec.md` → Build Order 2; System Architecture.
- ADR-0050, ADR-0114.
- `docs/backlog/2026-08-08-a-session-that-never-opened-is-a-selection-failure.md`

## Result

Implemented the first-Agent-output boundary and kept fallback behavior bounded
to selections that produced no output:

- `ACPXRunner` now publishes `agent_work_started` immediately before the first
  recognized Agent stream event and records that publication per Agent Session.
  Later turns in the same Session do not publish it again, and the inert
  runtime-deferred setup turn remains outside the boundary.
- A no-output adapter refusal now returns the distinct `SelectionFailure` type
  and publishes `agent_selection_failed`; it no longer arrives as a Batch
  failure. Adapter commands missing at Session preparation receive the same
  selection disposition.
- `agentSessionOwner` activates the next configured Fallback Selection after a
  `SelectionFailure` only while no Agent output exists. Its event sink records
  the first message, thought, plan, tool, or raw Agent event as work started, so
  even a misclassified post-output `SelectionFailure` cannot widen fallback.
- The Task 01 characterization was deliberately changed from
  `WorkStartedPrecedesTheFirstPrompt` to
  `WorkStartedFollowsTheFirstAgentOutput`; it now proves event order at the new
  boundary.

Focused checks:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent ./internal/daemon -count=1`
  passed 519 tests.
- `GOCACHE="$PWD/.gocache" rtk go test -race ./internal/agent ./internal/daemon -run '(WorkStartedBoundary|FallbackEligibility|RunDispositionCharacterizationWorkStartedFollows)' -count=1`
  passed seven focused boundary, fallback, and characterization cases under the
  race detector.

Acceptance evidence:

- `TestWorkStartedBoundaryReportsSelectionFailureWithoutOutput` observes one
  `agent_selection_failed` status, no `agent_work_started` status, and a
  `SelectionFailure` that is distinct from `BatchFailureError` after a
  no-output protocol refusal.
- `TestFallbackEligibilitySurvivesASelectionFailure` observes the preferred
  `codex` selection fail without output, the fallback notification publish, and
  the configured `claude` selection run. The adapter-start rehearsal is covered
  by `TestFallbackEligibilitySurvivesAdapterStartFailure`.
- `TestFallbackEligibilityEndsAfterAnyAgentOutput` publishes an Agent message
  before the failure and observes no fallback attempt. The ACPX exit matrix also
  keeps an exit-one failure after an Agent message classified as a Batch
  failure.
- `TestWorkStartedBoundaryPublishesOnFirstAgentOutput` runs two turns in one
  Session and observes exactly one work-started status before the first Agent
  message. `TestWorkStartedBoundaryIgnoresInertSessionSetup` proves the inert
  setup turn does not consume that one publication.

The authored `## Verification` commands were not run; the Daemon owns them.
