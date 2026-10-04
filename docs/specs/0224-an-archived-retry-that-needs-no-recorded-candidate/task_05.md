---
task: task_05
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: completed
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

## Result

The archived branch now admits a missing candidate to the archived-head
comparison based on `state.QAOverride` alone. Missing History still prevents
ancestry proof, so the retry refuses with the archived-head text and never
persists a change. The History guard on the proof itself, candidate anchors,
other retry branches, refusal texts and interfaces are unchanged.

Added `internal/delivery/operator_archive_no_history_test.go` using the real
queue store and existing recovery and engine helpers. Its blocker-named
subtests cover `run-unresolved` and `qa-environment-partial`; both explicitly
assert no History, compare the exact refusal including the empty candidate,
and compare the persisted item with its pre-retry snapshot.

Focused checks:

- `rtk proxy go test -count=1 -v -run '^TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory$' ./internal/delivery`
  could not access the sandbox-blocked default Go build cache.
- `GOCACHE=/tmp/roundfix-task05-gocache rtk proxy go test -count=1 -v -run '^TestRetryRefusesAnOperatorArchiveWithoutCandidateOrHistory$' ./internal/delivery`
  before the production change exited 1: both subtests returned
  `candidate head is missing` instead of the required archived-head refusal.
  Their unchanged-item assertions passed.
- `GOCACHE=/tmp/roundfix-task05-gocache rtk proxy go test -count=1 -v -run '^(TestRetry|TestOperatorArchiveRetry)' ./internal/delivery`
  after the production change exited 0. Acceptance criterion 1 is evidenced
  by both new subtests passing with the exact refusal and unchanged persisted
  items. Acceptance criterion 2 is evidenced by all existing operator-archive
  retry tests passing, including candidate and Run start anchors, no override,
  non-descendant heads, missing History and unreadable History. The two
  existing operator-archive test files were left unedited. The same check
  also passed prerequisite, corrective-Spec, conflict and active-Spec retry
  tests.

The declared Verification command was not run; the Daemon owns Verification
and status settlement. No commit, push or pull request was created. No
follow-up work was identified.

### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261004T192138Z_db255ee424790676/verification/batch-001-attempt-1.log`
and the related CLI detach test and process-group cleanup helper. The
Daemon's `make verify-changed` exited 2 after the CLI detach fixture cleanup
was denied permission to kill a live process group. The artifact records
`internal/delivery` passing; it does not establish a passing repository gate.

Focused feedback checks:

- `GOCACHE=/tmp/roundfix-task05-gocache rtk proxy go test -count=1 -v -run '^TestRunImplementDetachSurvivesCallerProcessGroupKill$' ./internal/cli`
  with host permissions exited 0. The isolated test passed, supporting an
  execution-permission blocker rather than a Task 05 implementation defect.
- `GOCACHE=/tmp/roundfix-task05-gocache rtk proxy go test -count=1 -run '^(TestRetry|TestOperatorArchiveRetry)' ./internal/delivery`
  exited 0; the regression cases and existing retry coverage remain passing.

No production or test changes were warranted by this feedback. The Daemon
retry needs permission to signal its disposable fixture process groups.
The full configured sequence and declared Verification remain Daemon-owned;
neither was rerun in this feedback turn. Task status is unchanged, and no
commit, push or pull request was created.

## Carry-forward provenance

- Source Run: `run_20261004T192138Z_db255ee424790676`
- Source commit: `3f97d1982ee2d4543cd23f384970e756ba32341a`
