---
task: task_04
spec: 0261-a-cancelled-spec-leaves-a-record
status: pending
type: backend
complexity: medium
---

# Task 04: A cancelled Spec meets no prerequisite, supersedes nothing and is never retried

## Overview

Three readers treat an archived Spec as delivered work: the Supersede
Command's deliverer check, `deliver start`'s prerequisite check and
`deliver retry`. This Task makes each refuse a cancelled Spec, and proves
that the review, the Run causes report, the Spec check and the History
Sanitize Command read a repository holding a cancelled record without an
error and without rewriting it (`_techspec.md` → Reader inventory, ADR-0262).
It answers
[the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)
of 2026-10-12 on the reader side. Its fixtures write cancelled records with
task_01's renderer, so it is verifiable without the Cancel Command.

## Requirements

1. MUST refuse `roundfix supersede --by <slug>` for a Spec whose Archive
   Record is cancelled, with the message of `_techspec.md` → API Contract 4
   and exit 2, and keep accepting every other archived deliverer.
2. MUST refuse `roundfix deliver start` when a queued Spec's `requires` names
   a Spec whose Archive Record is cancelled, with the message of API
   Contract 4 and exit 2, before any queue is recorded.
3. MUST refuse `roundfix deliver retry <slug>` when the checkout's Archive
   Record for the Spec is cancelled, with the message of API Contract 4 and
   exit 2, before the engine reads or writes the queue. A Spec without a
   record keeps today's retry path.
4. MUST NOT change `UnmetPrerequisites`, the review, the Run causes report,
   the Spec check or the History Sanitize Command; each is characterized
   instead.
5. MUST create `internal/cli/cancelled_record_readers_test.go` with
   `TestSupersedeRefusesACancelledDeliverer`,
   `TestDeliverStartRefusesARequiredCancelledSpec`,
   `TestDeliverRetryRefusesACancelledSpec`,
   `TestReviewReadsACancelledArchiveRecord`,
   `TestRunCausesReportACancelledRecord` and
   `TestHistorySanitizeLeavesACancelledRecord`, and
   `internal/speccheck/cancelled_record_test.go` with
   `TestSpecCheckReadsACancelledRecordAndItsReturnedEntry`. The review test
   MUST show a diff that removes a Spec folder and adds its cancelled record
   reaching the fake provider with the folder omitted as an archived Spec
   source. The Run causes test MUST show the disposition `cancelled`. The
   sanitize test MUST show a plan that lists no unit for the record, and the
   record byte-identical after `--apply` of a batch over another unit. The
   Spec check test MUST show no finding for a returned `open` Backlog Entry
   beside the cancelled record. The retry test MUST show the queue rows
   unchanged. Every test MUST use temporary repositories and homes, and every
   new top-level test MUST call `t.Parallel()` first or carry a
   `// Sequential: <reason>` comment.
6. MUST keep `TestSupersedeAcceptsAnArchiveRecord`,
   `TestPrerequisitesCountAnArchiveRecord` and
   `TestReviewDiffOmitsTheRemovedSpecFolder` passing unchanged.
7. MUST NOT change any command's help, the Cancel Command, `internal/spec` or
   any skill or guide.

## Subtasks

- [ ] Refuse a cancelled deliverer in the Supersede Command.
- [ ] Refuse a cancelled prerequisite in `deliver start`.
- [ ] Refuse a cancelled Spec in `deliver retry` before the engine.
- [ ] Characterize review, Run causes, Spec check and History Sanitize.

## Acceptance Criteria

- [ ] Each of the three refusals exits 2 with its message and writes nothing.
- [ ] Every other reader reads a cancelled record as it reads any record.
- [ ] The kept tests pass unchanged.

## Context

- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md`
- instruction: `internal/cli/archive_record_readers_test.go`
- instruction: `internal/cli/archive_record_test.go`
- instruction: `internal/cli/review_archived_spec_test.go`
- instruction: `internal/cli/history_test.go`
- instruction: `internal/speccheck/backlog_test.go`
- instruction: `internal/runcause/report.go`
- interface: `internal/cli/supersede.go`
- interface: `internal/cli/deliver.go`
- creates: `internal/cli/cancelled_record_readers_test.go`
- creates: `internal/speccheck/cancelled_record_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSupersedeRefusesACancelledDeliverer|TestDeliverStartRefusesARequiredCancelledSpec|TestDeliverRetryRefusesACancelledSpec|TestReviewReadsACancelledArchiveRecord|TestRunCausesReportACancelledRecord|TestHistorySanitizeLeavesACancelledRecord|TestSupersedeAcceptsAnArchiveRecord|TestPrerequisitesCountAnArchiveRecord|TestReviewDiffOmitsTheRemovedSpecFolder)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSupersedeRefusesACancelledDeliverer TestDeliverStartRefusesARequiredCancelledSpec TestDeliverRetryRefusesACancelledSpec TestReviewReadsACancelledArchiveRecord TestRunCausesReportACancelledRecord TestHistorySanitizeLeavesACancelledRecord TestSupersedeAcceptsAnArchiveRecord TestPrerequisitesCountAnArchiveRecord TestReviewDiffOmitsTheRemovedSpecFolder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^TestSpecCheckReadsACancelledRecordAndItsReturnedEntry$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestSpecCheckReadsACancelledRecordAndItsReturnedEntry" || { printf 'missing pass: TestSpecCheckReadsACancelledRecordAndItsReturnedEntry\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; User Story 5; Core Feature 5; Success Metric 3; Success Metric 4
- [_techspec.md](_techspec.md) — Reader inventory; API Contract 4; Testing Approach; Build Order 4
- ADR-0262; ADR-0193; ADR-0223; ADR-0229; ADR-0237; ADR-0165; ADR-0169; ADR-0248; ADR-0251; ADR-0259
