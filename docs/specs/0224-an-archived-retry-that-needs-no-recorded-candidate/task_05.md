---
task: task_05
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: pending
type: backend
complexity: low
---

# Task 05: An override archive without a candidate or History refuses with the archived-head text

## Overview

Corrective Task from finding F1 of the 2026-10-04 QA Report. API Contract 1
and task_01 Requirement 2 say that an operator-archived item with the QA
Archive Override, no candidate head and no item History refuses with
`archived item head "<head>" differs from candidate head ""`. The admission
predicate task_01 implemented, `state.QAOverride && engine.history != nil`,
is false when the engine has no History, so the candidate error returns first
and the refusal reads `candidate head is missing`, the text API Contract 1
reserves for an archive without the override. The QA probe confirmed the
item stays unchanged, so nothing is silently admitted; only the refusal names
the wrong cause. API Contract 1 governs: the override, not the presence of a
History, selects the archived-head refusal, and a missing History is one of
the proof failures that refuse with it.

## Requirements

1. MUST make the operator-archive admission in the archived branch of
   `Engine.Retry` (`internal/delivery/engine.go`) depend on
   `state.QAOverride` alone, so that an item with no candidate head, an
   override archive and no `History` reaches the archived-head comparison and
   refuses with
   `retry Delivery Queue item "<slug>": archived item head "<head>" differs from candidate head ""`,
   leaving the item unchanged, whatever its blocker.
2. MUST keep every other outcome of API Contract 1 unchanged: with a
   `History`, the Run start anchor, the descent proof, the candidate recorded
   and the `reviewing` stage; without the override, `candidate head is missing`;
   with a non-descendant head or an unreadable Run start, the archived-head
   text. An item that has a candidate keeps its current anchor and refusals.
3. MUST NOT change the prerequisite, `corrective-spec-required`,
   `pull-request-conflict` or active-Spec branches of the retry, the refusal
   texts, blocker names, Park Classes or any interface signature.
4. MUST add the new file
   `internal/delivery/operator_archive_no_history_test.go` with
   `TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory`, whose
   subtests are named by blocker and cover `run-unresolved` and
   `qa-environment-partial`: each seeds a parked item with no candidate
   commits in a real queue store, inspects an archive with
   `QAOverride: true` and head `operator-head` through `fakeItemRecovery`,
   builds the engine with `newRetryDeliveryEngine` and no `History` (and
   asserts it has none), and requires the exact archived-head refusal with an
   empty candidate and an unchanged persisted item. The tests of
   `operator_archive_retry_test.go` and `operator_archive_any_park_test.go`
   pass unedited.

## Subtasks

- [ ] Drop the History condition from the operator-archive admission.
- [ ] Add the no-History engine test for both parks.

## Acceptance Criteria

- [ ] An override archive without a candidate and without History refuses
      with the archived-head text naming the item head and an empty
      candidate, and the item is unchanged.
- [ ] Every existing operator-archive retry test passes unedited.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/operator_archive_no_history_test.go`
- instruction: `internal/delivery/operator_archive_any_park_test.go`
- instruction: `internal/delivery/operator_archive_retry_test.go`
- instruction: `internal/delivery/retry_test.go`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory|TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview|TestRetryRefusesAnArchiveWithoutCandidateOrOverride|TestRetryResumesAnOperatorArchivedItemAtReview|TestRetryRefusesAnOperatorArchiveWithoutAQAOverride|TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor|TestOperatorArchiveRetryRequiresHistory|TestOperatorArchiveRetryRefusesUnreadableHistory)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory/run-unresolved TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory/qa-environment-partial TestRetryResumesAnOperatorArchivedItemWithoutACandidateAtReview TestRetryRefusesAnArchiveWithoutCandidateOrOverride TestRetryRefusesAnArchiveWithoutCandidateOrOverride/no_override TestRetryRefusesAnArchiveWithoutCandidateOrOverride/not_descended TestRetryResumesAnOperatorArchivedItemAtReview TestRetryRefusesAnOperatorArchiveWithoutAQAOverride TestRetryRefusesAnArchivedHeadThatDoesNotDescendFromTheAnchor TestOperatorArchiveRetryRequiresHistory TestOperatorArchiveRetryRefusesUnreadableHistory; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new test does not exist, and with the test but without the change both subtests fail with candidate head is missing, so the command fails.

## References

- QA Report of 2026-10-04 (Run run_20261004T162613Z_0bb1b015eb38156b) → R2; R11; finding F1
- `_prd.md` → Goals 1 and 3; Core Feature 1; Success Metric 1
- `_techspec.md` → The archived retry without a candidate, step 4; Interfaces; API Contract 1; Build Order 5
- ADR-0229
