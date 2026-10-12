---
task: task_01
spec: 0261-a-cancelled-spec-leaves-a-record
status: pending
type: backend
complexity: high
---

# Task 01: A cancelled Archive Record keeps the reason and the Tasks, and Source Return hands back sources

## Overview

This Task adds the `cancelled` disposition to the Archive Record and the
`Cancel` function beside `Archive` in the Spec package. `Cancel` writes the
record with the reason, the date and the completed and pending Tasks, returns
each adopted Backlog Entry and Finding through Source Return, and removes the
Spec folder. It answers
[the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)
of 2026-10-12 on the record side (ADR-0262). It is verifiable on its own
through Spec-package tests over temporary repositories; no command exposes it
yet.

## Requirements

1. MUST add `ArchiveCancelled` (`cancelled`) and the record fields
   `CancelReason`, `CompletedTasks`, `PendingTasks` and `ReturnedSources`,
   rendered and parsed as `_techspec.md` → Data Models says. A cancelled
   record MUST round-trip byte for byte through `ParseArchiveRecord` and
   `RenderArchiveRecord`. A record of any other disposition that carries one of
   the four keys MUST be refused with
   `archive record disposition <d> cannot carry <key>`. Records of every
   existing disposition MUST render exactly as before.
2. MUST make `ArchivedTaskCompleted` return true for a cancelled record only
   for a Task listed in `completed_tasks`, and keep its result for every other
   disposition.
3. MUST implement `Cancel(CancelRequest) (CancelResult, error)` in the order
   of `_techspec.md` → Interfaces, writing nothing before step 6. Tasks MUST be
   read and never written; `completed` Tasks go to `completed_tasks` and every
   other status to `pending_tasks`, in graph order. A Spec without
   `_tasks.md` MUST cancel with both lists empty. The outcome paragraph MUST be
   the PRD's opening paragraph, and `archived` the cancellation date.
4. MUST implement Source Return as `_techspec.md` → Data Models says: a
   `backlog` row to `docs/backlog/` as `open` with `spec: null`; a `finding`
   row to `docs/findings/` as `pending` with the new `updated_at`, no
   `absorbed_by` line and the dated addendum; every other byte of the source
   kept; every other row type left in the folder. `returned_sources` MUST list
   the repository-relative destinations in index order.
5. MUST refuse, with an error and no write, each case of `_techspec.md` →
   Interfaces steps 1-5: no Spec folder or `_prd.md`, a `_supersession.md`, a
   blank reason, a revision that is not 40 hex, a manifest error, an occupied
   record or archived folder, an indexed file that is missing or a symbolic
   link, a returned source without a `status` key, and an occupied return
   destination.
6. MUST remove every file it wrote when a later write fails, before the Spec
   folder is touched, so a failed cancellation leaves the tree as it found it.
   The test MUST cause the failure through the filesystem (for example a
   read-only destination directory), never through a test-only hook.
7. MUST create `internal/spec/cancel_test.go` with
   `TestCancelRecordsCompletedAndPendingTasks`,
   `TestCancelReturnsAnAdoptedBacklogEntryAsOpen`,
   `TestCancelReturnsAnAdoptedFindingAsPending`,
   `TestCancelLeavesOtherAdoptedSourcesInTheFolder`,
   `TestCancelWithoutATaskGraph`, `TestCancelRefusalsWriteNothing` and
   `TestCancelRollsBackAFailedWrite`, and
   `internal/spec/archive_record_cancelled_test.go` with
   `TestACancelledArchiveRecordRoundTrips`,
   `TestOnlyACancelledRecordCarriesCancelFields` and
   `TestArchivedTaskCompletedForACancelledRecord`. The first test MUST use the
   adopter's shape from
   `docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md`:
   eight Tasks, five `completed` and three `pending`, and one adopted Backlog
   Entry. It MUST compare every returned byte outside the changed front matter
   lines with the adopted copy. Every new top-level test MUST call
   `t.Parallel()` first or carry a `// Sequential: <reason>` comment.
8. MUST keep `TestArchiveCommandLeavesTheArchiveRecord` and the existing
   record tests passing unchanged, and MUST NOT change `Archive`,
   `BuildArchiveRecord`'s results for archives, the History Sanitize Command or
   any CLI file.

## Subtasks

- [ ] Add the disposition, the four fields, their rendering and their parsing.
- [ ] Make Task completion follow `completed_tasks` for a cancelled record.
- [ ] Implement `Cancel` with its refusals and rollback.
- [ ] Implement Source Return for Backlog Entries and Findings.
- [ ] Write the two test files.

## Acceptance Criteria

- [ ] The 5-and-3 fixture cancels into a record that lists those Tasks, the
      reason and the date, and its Backlog Entry is back as `open`.
- [ ] A returned Finding is `pending` with the addendum and no `absorbed_by`.
- [ ] Every refusal and the injected write failure leave the tree unchanged.
- [ ] A cancelled record round-trips, and no other disposition may carry its
      fields.

## Context

- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md`
- instruction: `internal/spec/archive.go`
- instruction: `internal/spec/supersession.go`
- instruction: `internal/spec/archive_record_test.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/archive_reader.go`
- creates: `internal/spec/cancel.go`
- creates: `internal/spec/cancel_test.go`
- creates: `internal/spec/archive_record_cancelled_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCancelRecordsCompletedAndPendingTasks|TestCancelReturnsAnAdoptedBacklogEntryAsOpen|TestCancelReturnsAnAdoptedFindingAsPending|TestCancelLeavesOtherAdoptedSourcesInTheFolder|TestCancelWithoutATaskGraph|TestCancelRefusalsWriteNothing|TestCancelRollsBackAFailedWrite)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCancelRecordsCompletedAndPendingTasks TestCancelReturnsAnAdoptedBacklogEntryAsOpen TestCancelReturnsAnAdoptedFindingAsPending TestCancelLeavesOtherAdoptedSourcesInTheFolder TestCancelWithoutATaskGraph TestCancelRefusalsWriteNothing TestCancelRollsBackAFailedWrite; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of these tests exists, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestACancelledArchiveRecordRoundTrips|TestOnlyACancelledRecordCarriesCancelFields|TestArchivedTaskCompletedForACancelledRecord)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACancelledArchiveRecordRoundTrips TestOnlyACancelledRecordCarriesCancelFields TestArchivedTaskCompletedForACancelledRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestArchiveCommandLeavesTheArchiveRecord$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestArchiveCommandLeavesTheArchiveRecord" || exit 1` — expected: exit 0; before this Task the record tests do not exist, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 2; Goal 3; User Story 1; User Story 2; User Story 5; Core Feature 2; Core Feature 3; Core Feature 4; Core Feature 5; Success Metric 1; Success Metric 3; Acceptance evidence
- [_techspec.md](_techspec.md) — Reader inventory; Interfaces; Data Models; API Contract 3; Testing Approach; Build Order 1
- ADR-0262; ADR-0247; ADR-0248; ADR-0251; ADR-0057; ADR-0259
