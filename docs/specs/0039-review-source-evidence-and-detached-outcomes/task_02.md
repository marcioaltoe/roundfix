---
task: task_02
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: backend
complexity: high
---

# Task 02: Unify CodeRabbit evidence classification

## Overview

Replace separate pre-fetch and Merge-Ready interpretations with one CodeRabbit
evidence hierarchy. Both watch phases receive the same result for the same
head, including explicit skip, current-head approval, unresolved-thread, stale
signal, and no-signal cases.

## Requirements

1. MUST expose one evidence operation for an Open Pull Request and expected
   head.
2. MUST apply the TechSpec evidence hierarchy in stable precedence order.
3. MUST classify explicit structured skip before successful or pending signals.
4. MUST require the expected head for every accepted check, status, or review.
5. MUST accept a current-head CodeRabbit `APPROVED` review as verified only when
   zero unresolved CodeRabbit threads remain.
6. MUST classify other current-head reviews as reviewed, never verified
   approval.
7. MUST publish changed observations once and suppress duplicate unchanged
   polling events.

## Subtasks

- [x] Implement the shared evidence operation.
- [x] Unify check-run, commit-status, review, and thread precedence.
- [x] Route pre-fetch and Merge-Ready checks through the classifier.
- [x] Enforce expected-head and approval requirements.
- [x] Publish bounded changed-observation events.
- [x] Add stale-head, unresolved-thread, and precedence tables.

## Acceptance Criteria

- [x] Pre-fetch and Merge-Ready return byte-equivalent Evidence for the same
      fixture.
- [x] Explicit skip wins over success or pending signals and preserves reason.
- [x] Current-head approval with zero unresolved threads produces verified
      `review_approval` Evidence.
- [x] Commented, stale-head, or unresolved-thread reviews cannot prove
      Merge-Ready.
- [x] No usable signal produces pending without guessing.
- [x] One changed observation creates one event; an unchanged poll creates
      none.
- [x] Existing Review Source authentication and repository boundaries remain
      unchanged.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/reviewsource/reviewsource.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/coderabbit/coderabbit_test.go`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/reviewsource/... -run 'Test.*(EvidenceHierarchy|Approval|Skipped|ExpectedHead|UnresolvedThread|Precedence)' -count=1`
  — expected: one classifier produces every accepted and refused state.
- `rtk go test ./internal/watch ./internal/runevent -run 'Test.*(ReviewEvidence|ReviewStatusEvent|UnchangedEvidence)' -count=1`
  — expected: both watch phases share Evidence and events deduplicate unchanged
  polls.

## References

- `_prd.md` → Goal 1; User Stories 1 and 7; Core Features 1–3; Success Metrics.
- `_techspec.md` → API Contracts: Review Source evidence hierarchy; Integration
  Points: GitHub through `gh`; Build Order 2.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  accepted Review Source Evidence.

## Result

Implemented one CodeRabbit Evidence operation for an Open Pull Request and
expected head. The watch command now sends both pre-fetch and Merge-Ready
observations through that operation. Classification applies the documented
skip, reviewing, successful check/status, approval, other review, stale signal,
failure, and no-signal precedence while counting unresolved CodeRabbit threads.
Changed Evidence publishes one bounded `daemon.review_status` payload; an
unchanged poll publishes nothing.

Verification:

- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test ./internal/reviewsource/... -run 'Test.*(EvidenceHierarchy|Approval|Skipped|ExpectedHead|UnresolvedThread|Precedence)' -count=1`
  — passed, 10 tests in 2 packages.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test ./internal/watch ./internal/runevent -run 'Test.*(ReviewEvidence|ReviewStatusEvent|UnchangedEvidence)' -count=1`
  — passed, 2 tests in 2 packages.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk go test ./internal/cli -run 'Test(RunOperationalCommandAcceptsMVPFlags|RunWatchNoAgentConsoleSuppressesAgentDisplayOnly|AttachRendersWatchDaemonEventsInTimeline)' -count=1`
  — passed, 6 tests.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify` — passed outside
  the sandbox: 2,598 Go tests, 4 skill tests, the Roundfix skill check, and the
  production build. The sandboxed attempt could not execute `/bin/ps`; the
  approved unsandboxed rerun passed.

Acceptance evidence:

- `TestRunReviewEvidenceSharedByPreFetchAndMergeReady` proves both phases send
  the same request, receive the same Evidence, and suppress the unchanged
  second event.
- `TestEvidenceHierarchyPrecedence` proves explicit structured skip preserves
  its reason and wins over pending and successful signals.
- The approval and unresolved-thread rows prove only a current-head
  `APPROVED` review with zero unresolved CodeRabbit threads produces verified
  `review_approval` Evidence.
- Commented and unresolved-thread table rows plus
  `TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals` prove those signals
  cannot verify Merge-Ready.
- The no-signal row proves the classifier returns pending without inferring a
  successful review.
- `TestReviewStatusEventPayloadUsesStableEvidenceFields` and the shared watch
  test prove bounded changed-observation event fields and deduplication.
- The implementation reuses the existing `GitHubClient` methods, base
  repository metadata, and `gh` transport; the full verification gate proves
  the existing authentication and repository-boundary regressions remain
  green.

Follow-ups: none.
