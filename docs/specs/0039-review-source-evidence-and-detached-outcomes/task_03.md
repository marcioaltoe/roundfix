---
task: task_03
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: backend
complexity: high
---

# Task 03: Settle Review Skipped without Round artifacts

## Overview

Add Review Skipped as a terminal review outcome across persistence, watch
settlement, reports, streams, and Run Browser projection. An explicit skip ends
before fetch and publishes actionable reason evidence without creating an
empty or partial Round artifact.

## Requirements

1. MUST add Review Skipped to terminal Run state validation and storage.
2. MUST advance the Run Database schema fence before the new state can be
   persisted.
3. MUST map explicit skipped Evidence to Review Skipped and existing exit code
   `3`.
4. MUST render the Review Source reason and next action without zero issue
   counts.
5. MUST project the outcome through Run Events and Run Browser state.
6. MUST skip Round fetch and create no Round directory for Review Skipped.
7. MUST keep failed pre-fetch writes temporary and unpublished.
8. MUST use the spec 0037 registered-session cleanup boundary when settlement
   occurs before Agent Start, preserving the Review Source reason ahead of any
   secondary cleanup warning.

## Subtasks

- [ ] Add the terminal state and migration fence.
- [ ] Map skipped Evidence through watch and CLI settlement.
- [ ] Render actionable text and machine outcome evidence.
- [ ] Add Run Browser and stream projection.
- [ ] Prevent fetch and artifact publication on skip.
- [ ] Exercise cleanup-before-Agent settlement through the registered-session
      boundary.
- [ ] Add migration, exit, report, and no-artifact cases.

## Acceptance Criteria

- [ ] A structured skip ends Review Skipped with exit `3`.
- [ ] Persisted Review Skipped is terminal and older schema readers refuse the
      newer database.
- [ ] The report names skip reason and next action and prints no zero-valued
      Review Issue summary.
- [ ] No fetcher call, Round directory, partial artifact, or review-artifact
      commit occurs.
- [ ] A pre-Agent skip with no active Agent Selection lifecycle performs zero
      Agent Session cleanup calls; any eligible cleanup warning remains
      secondary to the Review Source reason.
- [ ] Outcome stream and Run Browser render the same terminal state.
- [ ] Non-skip Evidence preserves existing Reviewed, Clean, Clean Unverified,
      and failure behavior.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/store/store_test.go`
- interface: `internal/store/agent_selection.go`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/runbrowser.go`
- interface: `internal/runevent/event.go`
- interface: `internal/agent/sessions.go`

## Verification

- `rtk go test ./internal/store -run 'Test.*(ReviewSkipped|Schema.*Review)' -count=1`
  — expected: migration fence and terminal state persistence pass.
- `rtk go test ./internal/watch ./internal/cli ./internal/runevent -run 'Test.*(ReviewSkipped|CleanupBeforeAgent)' -count=1`
  — expected: skipped Evidence settles, reports, projects, and creates no Round
  artifact while registered-session cleanup preserves primary reason ordering.
- `rtk go test -race ./internal/watch ./internal/cli -run 'Test.*ReviewSkipped' -count=1`
  — expected: terminal skip settlement is race-free.

## References

- `_prd.md` → Goal 1; User Story 1; Core Features 2, 7, and 12; User Experience;
  Success Metrics.
- `_techspec.md` → Data Models: schema fence; API Contracts: Review Skipped and
  Artifact persistence; Build Order 3.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  Review Skipped outcome.
