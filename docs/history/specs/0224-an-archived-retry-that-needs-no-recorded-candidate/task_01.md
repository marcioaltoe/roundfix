---
task: task_01
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: completed
type: backend
complexity: low
---

# Task 01: A Delivery Retry resumes an operator-archived item from its Run start, whatever its park

## Overview

On 2026-10-04 the Delivery Retry of Spec 0220, parked `run-unresolved` with no
candidate head and archived by the operator with the QA Archive Override,
refused with `candidate head is missing` (operator intervention log entries
140 to 142 of 2026-10-04). Only a `qa-environment-partial` park may anchor at
the Implement start head of its Run today. This Task lets any park that
reaches the archived branch of the retry use that anchor when the archive
records `qa_override: true` and no candidate exists (ADR-0229).

## Requirements

1. MUST implement "The archived retry without a candidate" of the TechSpec in
   `Engine.Retry`: when the archive records `qa_override: true`, the engine
   has a `History` and the item has no candidate head, anchor at
   `History.RunStart` of the item's Run ID whatever the item's blocker, and on
   proven descent append the item head to the candidate commits and set the
   stage `reviewing`, as API Contract 1 states.
2. MUST keep every refusal of API Contract 1 byte-identical: without the
   override, `candidate head is missing`; with a non-descendant head, an
   unreadable Run start or no `History`, the existing
   `archived item head "<head>" differs from candidate head ""` text; each
   leaves the item unchanged.
3. MUST NOT change the prerequisite, `corrective-spec-required`,
   `pull-request-conflict` or active-Spec branches of the retry, the anchor
   for an item that has a candidate, blocker names, Park Classes or any
   interface signature.
4. MUST add the new file `internal/delivery/operator_archive_any_park_test.go`
   with `TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview`,
   whose subtests are named by blocker and cover `run-unresolved`,
   `qa-environment-partial` and `run-budget-exceeded`, and
   `TestRetryRefusesAnArchiveWithoutCandidateOrOverride`, with the subtests
   `no override` and `not descended`, over the fakes of
   `operator_archive_retry_test.go` and a real queue store, as Testing
   Approach 1 describes. The tests of `operator_archive_retry_test.go` pass
   unedited.
5. MUST, in `internal/cli/deliver_archived_retry_test.go`, move the body of
   `TestArchivedRetryOfAQueueStartedRunReturnsToReview` into a helper that
   takes the blocker, keep that test calling it with `qa-environment-partial`,
   and add `TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview`
   calling it with `run-unresolved`, through a real
   `roundfix archive --qa-override` in a linked item worktree and a Run in a
   disposable Run Database, as Testing Approach 2 describes. No test reads an
   archived Spec or writes under the real `~/.roundfix`.

## Subtasks

- [ ] Drop the blocker condition from the operator-archive admission in the retry.
- [ ] Add the engine tests for every park and for the refusals.
- [ ] Turn the real-archive test into a helper and add the `run-unresolved` case.

## Acceptance Criteria

- [ ] A `run-unresolved`, `qa-environment-partial` or `run-budget-exceeded`
      item with no candidate and an override archive resumes at `reviewing`
      with exactly the item head as its candidate, anchored at its Run start.
- [ ] Without the override the retry still refuses with
      `candidate head is missing`, and a non-descendant head still refuses
      with the existing text; the item is unchanged in both.
- [ ] Every existing operator-archive and Run start test passes unedited.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/operator_archive_any_park_test.go`
- interface: `internal/cli/deliver_archived_retry_test.go`
- instruction: `internal/delivery/operator_archive_retry_test.go`
- instruction: `internal/delivery/retry_test.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview|TestRetryRefusesAnArchiveWithoutCandidateOrOverride|TestRetryResumesAnOperatorArchivedItemAtReview|TestRetryRefusesAnOperatorArchiveWithoutAQAOverride|TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor|TestOperatorArchiveRetryRequiresHistory|TestOperatorArchiveRetryRefusesUnreadableHistory)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview/run-unresolved TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview/run-budget-exceeded TestRetryRefusesAnArchiveWithoutCandidateOrOverride TestRetryResumesAnOperatorArchivedItemAtReview TestRetryRefusesAnOperatorArchiveWithoutAQAOverride TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor TestOperatorArchiveRetryRequiresHistory TestOperatorArchiveRetryRefusesUnreadableHistory; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview|TestArchivedRetryOfAQueueStartedRunReturnsToReview|TestRunStartFindsAQueueStartedRunByItsRepository|TestRunStartRefusesARunOfAnotherRepository)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview TestArchivedRetryOfAQueueStartedRunReturnsToReview TestRunStartFindsAQueueStartedRunByItsRepository TestRunStartRefusesARunOfAnotherRepository; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the run-unresolved test does not exist, and with the test but without the change it fails with candidate head is missing, so the command fails.

## References

- `_prd.md` → Goals 1 and 3; Core Feature 1; Success Metric 1; Success Metric 3
- `_techspec.md` → The archived retry without a candidate; Interfaces; API Contract 1; Testing Approach 1; Testing Approach 2; Build Order 1
- ADR-0229; ADR-0223

## Result

The archived retry admission now uses the TechSpec's predicate
`state.QAOverride && engine.history != nil`, independent of blocker. The
existing Run start lookup, ancestry proof, candidate recording and review
transition remain in place. No other retry branch or interface changed.
The real-archive CLI fixture now takes a blocker; the existing environment
park and the new unresolved park both exercise it in disposable repositories
and a disposable Run Database.

Focused evidence per acceptance criterion:

- Three parks: the new engine test covers `run-unresolved`,
  `qa-environment-partial` and `run-budget-exceeded`, asserting `reviewing`,
  cleared blocker, exactly the operator head as candidate, the item's Run ID
  and the Run start anchor. Before the engine change, the unresolved and
  budget subtests refused with `candidate head is missing`.
- Refusals: the new refusal test checks the complete error text for no
  override and non-descent and compares the persisted item before and after.
  Both cases pass in the focused delivery check.
- Existing coverage: the focused delivery check includes the existing
  operator-archive tests, unchanged, and the CLI check includes the existing
  queue-started archive and repository-scoped Run start tests. Both checks
  pass. Before the change, the new real-archive unresolved test reproduced
  `candidate head is missing`; after the change it passes.

Commands and outcomes:

- The initial focused Go test could not read the sandbox-restricted default
  Go cache. Subsequent checks use `GOCACHE=/tmp/roundfix-0224-task01-gocache`.
- `GOCACHE=/tmp/roundfix-0224-task01-gocache rtk proxy go test -count=1 -run 'TestRetry|TestOperatorArchive' ./internal/delivery`: exit 0.
- `GOCACHE=/tmp/roundfix-0224-task01-gocache rtk proxy go test -count=1 -run 'TestArchivedRetry|TestRunStart' ./internal/cli`: exit 0.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.
- `GOCACHE=/tmp/roundfix-0224-task01-gocache rtk make verify-incremental`:
  initial attempt exited 2 due to sandbox-denied process inspection and local
  test ports, plus repository guards detecting this Result edit during the
  suite. Repeated with required permissions and no concurrent worktree edits:
  exit 0, including formatting, vet, tests, skill checks and build.
- Authored Verification commands were not run; they remain Daemon-owned.

Unresolved contract conflict: requirement 2 and API Contract 1 specify the
archived-head mismatch error for an override archive without a candidate or
History. The TechSpec's exact predicate above retains the earlier
`candidate head is missing` refusal when History is absent. Existing
missing-History coverage has a candidate and is unchanged. The maintainer
was asked which contract governs; the missing-History admission guard is
left unchanged pending that decision. No terminal Task verdict is claimed.

## Carry-forward provenance

- Source Run: `run_20261004T162613Z_0bb1b015eb38156b`
- Source commit: `be108d18093416f8f4354ebb7a7bde05b52727dc`
