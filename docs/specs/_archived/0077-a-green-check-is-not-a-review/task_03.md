---
task: task_03
spec: 0077-a-green-check-is-not-a-review
status: completed
type: backend
complexity: low
---

# Task 03: Say what was observed instead of merging

## Overview

A stalled loop that does not say why is a stalled loop an operator will
override. This slice makes the refusal and the unclassified case legible where
the operator is already looking: the watch output and the Run Event Stream.

## Requirements

1. MUST make `watch --until-clean` name what it observed when it declines to
   proceed — the refusal and its reason, or that the signal was unrecognised.
2. MUST record the evidence state and reason in the Run Event Stream, so an
   operator can tell a refused review from an absent one without opening GitHub.
3. MUST NOT change what `watch` does on `verified` or `reviewed`.
4. MUST NOT push or merge on a `skipped` or `pending` head, asserted rather than
   assumed.

## Subtasks

- [ ] Name the observation in the watch diagnostic.
- [ ] Record state and reason in the Run Event Stream.
- [ ] Assert no push and no merge on refused or unclassified heads.

## Acceptance Criteria

- [ ] A refused head produces a diagnostic naming the refusal and its reason.
- [ ] An unclassified head produces a diagnostic saying the signal was not
      recognised.
- [ ] The Run Event Stream carries the evidence state and reason for both.
- [ ] Neither case pushes or merges, asserted by a fixture that would have
      merged before this Spec.
- [ ] A verified head behaves exactly as it does today.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/runevent`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/watch -count=1 -run 'Skipped|Pending|Diagnostic|Clean' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the watch tests ran and passed.
- `go test ./internal/watch ./internal/runevent -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; the Skill is task_04's bounded scope.

## References

- `_prd.md` → Core Features 4 and 6; Goals.
- `_techspec.md` → Build Order 3.

## Result

Implementation:

- Made a timed-out `pending` observation with a concrete Evidence kind retain
  its Evidence and explain that the Review Source signal was not recognised,
  including the existing bounded Evidence detail. Missing signals keep the
  existing generic timeout behavior.
- Kept the existing `skipped` terminal path and its source reason, and made the
  watch command print the unrecognised-signal terminal reason before its
  existing timeout guidance.
- Added `evidence_state` to the additive terminal Run Event payload and public
  `roundfix-events/v1` outcome projection. The existing outcome `reason` now
  carries the refusal reason or unrecognised-signal diagnostic alongside it.
- Added command-level assertions that refused and unrecognised green signals
  make zero Final Push calls and return non-Clean outcomes. The command has no
  merge operation; its non-Clean exit prevents the Supervisor's subsequent
  merge step.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=<worktree>/.gocache go test ./internal/watch ./internal/runevent ./internal/cli -count=1 -run '^(TestRunUnrecognisedPendingEvidenceStopsBeforeFetchWithDiagnostic|TestProjectStreamEventReviewSkippedOutcome|TestProjectStreamEventOutcomeContextProjectsReviewIssuesEvidenceAndRecovery|TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup|TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush)$'`
  exited 1: watch still returned the generic timeout reason, and the Run Event
  types had no `EvidenceState` field.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/watch -count=1 -run '^(TestRunUnrecognisedPendingEvidenceStopsBeforeFetchWithDiagnostic|TestRunReviewSkippedStopsBeforeFetch|TestRunReviewSkippedDuringMergeReadyPreservesTerminalEvidence|TestRunReviewEvidenceSharedByPreFetchAndMergeReady)$'`
  — exit 0 after the last Go edit.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/runevent -count=1 -run '^(TestProjectStreamEventReviewSkippedOutcome|TestProjectStreamEventPendingUnrecognisedOutcome|TestProjectStreamEventOutcomeContextProjectsReviewIssuesEvidenceAndRecovery)$'`
  — exit 0 after the last Go edit.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/cli -count=1 -run '^(TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup|TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush|TestRunWatchReviewIssuesKnownAfterFetchedZero|TestRunWatchArtifactEvidenceInheritedWithoutCurrentHeadPolling)$'`
  — exit 0 after the last Go edit.
- `rtk git diff --check` — exit 0 after the last Go edit.

Acceptance evidence:

1. `TestRunReviewSkippedStopsBeforeFetch` preserves the refusal as
   `ReviewSkipped` with the source reason, while
   `TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup` asserts
   the visible refusal diagnostic and reason.
2. `TestRunUnrecognisedPendingEvidenceStopsBeforeFetchWithDiagnostic` asserts
   `TimedOut`, retained `pending` Evidence, and the explicit
   `signal was not recognised` reason; the command-level companion asserts the
   same text reaches stderr.
3. `TestProjectStreamEventReviewSkippedOutcome` and
   `TestProjectStreamEventPendingUnrecognisedOutcome` assert the public outcome
   stream carries `evidence_state` plus the applicable reason for both cases.
   The two command tests also decode and assert the journaled terminal payloads.
4. `TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup` and
   `TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush` both assert zero
   Final Push calls. The latter models the formerly permissive success-conclusion
   and unknown-output fixture, asserts a non-Clean exit, and therefore does not
   authorize the Supervisor's merge step.
5. `TestRunReviewEvidenceSharedByPreFetchAndMergeReady` retains the verified
   watch path. `TestRunWatchReviewIssuesKnownAfterFetchedZero` still reaches
   `Clean`, and `TestRunWatchArtifactEvidenceInheritedWithoutCurrentHeadPolling`
   still makes exactly one Final Push for verified Evidence.

Declared `## Verification` commands were not run; the Daemon owns them.

Follow-up: task_04 owns the authorized Roundfix Skill synchronisation; no Skill
file or generated mirror is part of this diff.
