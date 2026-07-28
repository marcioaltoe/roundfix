---
task: task_03
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
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

- [x] Add the terminal state and migration fence.
- [x] Map skipped Evidence through watch and CLI settlement.
- [x] Render actionable text and machine outcome evidence.
- [x] Add Run Browser and stream projection.
- [x] Prevent fetch and artifact publication on skip.
- [x] Exercise cleanup-before-Agent settlement through the registered-session
      boundary.
- [x] Add migration, exit, report, and no-artifact cases.

## Acceptance Criteria

- [x] A structured skip ends Review Skipped with exit `3`.
- [x] Persisted Review Skipped is terminal and older schema readers refuse the
      newer database.
- [x] The report names skip reason and next action and prints no zero-valued
      Review Issue summary.
- [x] No fetcher call, Round directory, partial artifact, or review-artifact
      commit occurs.
- [x] A pre-Agent skip with no active Agent Selection lifecycle performs zero
      Agent Session cleanup calls; any eligible cleanup warning remains
      secondary to the Review Source reason.
- [x] Outcome stream and Run Browser render the same terminal state.
- [x] Non-skip Evidence preserves existing Reviewed, Clean, Clean Unverified,
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

## Result

Review Skipped is now a persisted terminal Run outcome behind Run Database
schema version 11. Explicit skipped Evidence settles before fetch, returns exit
`3`, preserves the bounded Review Source reason and next action in the terminal
Run Event, and prints a dedicated report without Review Issue counts. The Run
Browser lists the outcome as terminal, and the Supervisor outcome projection
renders the same state.

The skip path uses the registered Agent Selection lifecycle as its only Agent
Session cleanup registry. A Run with no active lifecycle makes no cleanup
calls. Eligible cleanup failures remain labeled secondary and print after the
Review Source reason. Fetch failures and skipped reviews publish no Round
directory, partial artifact, or review-artifact commit.

Verification:

- Pre-change
  `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/store ./internal/watch ./internal/cli ./internal/runevent -run 'Test.*(ReviewSkipped|Schema.*Review|CleanupBeforeAgent)' -count=1`
  failed because Review Skipped had no terminal state, watch result, event
  payload, report, or exit mapping.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/store -run 'Test.*(ReviewSkipped|Schema.*Review)' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/watch ./internal/cli ./internal/runevent -run 'Test.*(ReviewSkipped|CleanupBeforeAgent)' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -race ./internal/watch ./internal/cli -run 'Test.*ReviewSkipped' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/store -run 'Test(Schema9CreatesAgentSelectionTable|Schema8To9MigratesRunsEventsAndSelectionTable|OpenMigratesV[3-9]RunDatabase|TerminalOutcomeEveryStoredTerminalStateIsImmutable|PruneTerminalRunsDeletesOnlyEligibleJournalRows)' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/watch ./internal/cli -run 'Test(RunReviewEvidenceSharedByPreFetchAndMergeReady|RunWaitsFetchesResolvesToClean|RunWatchMissingHeadCheckEndsCleanUnverified|RunWatchReusesOneAgentSessionAcrossRoundsAndCloses|RunWatchPrintsDeterministicStdoutReport|ExitForWatchOutcome)' -count=1`
  — passed.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

The focused commands use a task-local writable `GOCACHE` because the sandbox
cannot access the host Go build cache. The Daemon owns the exact declared
Verification commands after this Agent turn.

Acceptance evidence:

- `TestRunReviewSkippedStopsBeforeFetch` and
  `TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup` prove a
  structured skip settles Review Skipped with exit `3`.
- `TestCompleteRunReviewSkippedIsTerminal`,
  `TestSchemaBeforeReviewSkippedMigrationRequiresWriter`, and
  `TestSchemaReviewSkippedReaderRejectsNewerDatabase` prove terminal
  persistence, schema migration, and the exact-version reader fence.
- `TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup` proves
  the report carries reason and next action without zero-valued issue
  summaries.
- The same CLI test proves zero fetch, Agent Session cleanup,
  review-artifact commit, and Round directory creation.
  `TestRunWatchReviewSkippedArtifactBoundaryKeepsFailedFetchUnpublished`
  proves a failed pre-fetch operation also leaves no published artifact.
- `TestRunWatchCleanupBeforeAgentWarningFollowsReviewSkippedReason` proves
  eligible cleanup warnings stay secondary to the primary Review Source
  reason.
- `TestProjectStreamEventReviewSkippedOutcome` and the Run Browser assertions
  in the CLI test prove both projections render ReviewSkipped.
- The non-skip regression checks preserve Reviewed-to-fetch, Clean,
  Clean Unverified, Agent Session close, report, exit, migration, terminal
  immutability, and journal-pruning behavior.

Follow-ups: none.
