---
task: task_01
spec: 0261-a-cancelled-spec-leaves-a-record
status: completed
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
9. MUST apply the Archive Relink of ADR-0260 (Spec 0255, merged first) before
   the folder is removed, exactly as `Archive` does: call
   `PlanArchiveRelink` with no manual promotions, copy each non-core file a
   Live Document still cites to `docs/references/` and rewrite the inbound
   links, list the rewritten files under the record's `relinked:`, and refuse
   with no write on the relink plan's refusals. Add
   `TestCancelRelinksLiveCitations` to `internal/spec/cancel_test.go`: an
   active Spec that links a `references/` file of the cancelled Spec still
   resolves after the cancellation, and the record lists it under `relinked:`.

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

- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
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

- `out="$(go test -count=1 -v -run '^(TestCancelRecordsCompletedAndPendingTasks|TestCancelReturnsAnAdoptedBacklogEntryAsOpen|TestCancelReturnsAnAdoptedFindingAsPending|TestCancelLeavesOtherAdoptedSourcesInTheFolder|TestCancelWithoutATaskGraph|TestCancelRefusalsWriteNothing|TestCancelRollsBackAFailedWrite|TestCancelRelinksLiveCitations)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCancelRecordsCompletedAndPendingTasks TestCancelReturnsAnAdoptedBacklogEntryAsOpen TestCancelReturnsAnAdoptedFindingAsPending TestCancelLeavesOtherAdoptedSourcesInTheFolder TestCancelWithoutATaskGraph TestCancelRefusalsWriteNothing TestCancelRollsBackAFailedWrite TestCancelRelinksLiveCitations; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of these tests exists, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestACancelledArchiveRecordRoundTrips|TestOnlyACancelledRecordCarriesCancelFields|TestArchivedTaskCompletedForACancelledRecord)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACancelledArchiveRecordRoundTrips TestOnlyACancelledRecordCarriesCancelFields TestArchivedTaskCompletedForACancelledRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestArchiveCommandLeavesTheArchiveRecord$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestArchiveCommandLeavesTheArchiveRecord" || exit 1` — expected: exit 0; before this Task the record tests do not exist, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 2; Goal 3; User Story 1; User Story 2; User Story 5; Core Feature 2; Core Feature 3; Core Feature 4; Core Feature 5; Success Metric 1; Success Metric 3; Acceptance evidence
- [_techspec.md](_techspec.md) — Reader inventory; Interfaces; Data Models; API Contract 3; Testing Approach; Build Order 1
- ADR-0262; ADR-0247; ADR-0248; ADR-0251; ADR-0057; ADR-0259

## Result

Implemented the Task 01 package slice. `Cancel` reads Task state in graph
order, records cancellation without changing Tasks, returns adopted Backlog
Entries and Findings, applies Archive Relink with no manual promotions, and
removes the Spec folder only after every write succeeds. Cancelled records
carry the four required fields; other dispositions reject those keys and keep
their existing rendering. Task completion for cancellation follows only
`completed_tasks`.

### Focused checks

- Starting signal: `rtk proxy go test ./internal/spec -run '^TestACancelledArchiveRecordRoundTrips$/true$' -count=1`
  exited 1 with undefined `ArchiveCancelled` and missing cancellation fields.
- `rtk proxy go test ./internal/spec -run '^(TestCancel|TestACancelled|TestOnlyACancelled|TestArchivedTaskCompletedForACancelled|TestArchiveRecord|TestArchiveWritesTheArchiveRecord|TestArchiveRefusesBeforeAnyFileChanges)' -count=1`
  exited 0 after correcting the new fixture's QA declaration and snapshot
  path keys. This exercises every new cancellation and cancelled-record test,
  plus existing record rendering and refusal regressions.
- `rtk proxy go test ./internal/spec -run '^(TestReadArchivedSpecReadsRecordAndLegacyFolder|TestArchiveRollsBackTheRecordWhenPromotionCannotBeWritten)$' -count=1`
  exited 0. Together with the preceding selection, all tests in the unchanged
  `archive_record_test.go` were exercised.
- `rtk git diff --check` exited 0.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Five completed and three pending Tasks, reason, date and returned open Backlog Entry | `TestCancelRecordsCompletedAndPendingTasks` uses eight Tasks and one adopted Backlog Entry; checks both Task lists, normalized reason, cancellation date, PRD opening paragraph, revision, removed paths/count/bytes and exact returned bytes with only the two front matter values replaced. `TestCancelKeepsObservedQAAndTaskGraphOrder` covers dependency order, failed/in-progress Tasks and observed failed QA metadata. |
| Finding returns pending with addendum and no absorption pointer | `TestCancelReturnsAnAdoptedFindingAsPending` compares the entire returned document with only the prescribed front matter edits and exact dated addendum. It preserves prior evidence, comments and source-index return order; also covers a missing `updated_at`. Backlog tests cover comments, inserted `spec: null` and CRLF bytes. |
| Every refusal and later filesystem write failure leave the tree unchanged | `TestCancelRefusalsWriteNothing` checks 17 refusal cases against snapshots including file bytes, directory presence, modes and symlink targets. `TestCancelRollsBackAFailedWrite` makes the Findings directory read-only: record and Backlog Entry writes happen before the Finding write fails, then all created files and directories are removed and the original tree matches. `TestCancelRelinksLiveCitations/rewrite_failure` makes the later live document read-only and proves prior relinks and copies are rolled back. No test-only production hook is used. |
| Cancelled record round-trips and other dispositions reject its fields | `TestACancelledArchiveRecordRoundTrips` checks populated and empty lists plus every required key's absence; `TestOnlyACancelledRecordCarriesCancelFields` checks each key against all six existing dispositions and the exact refusal text. `TestArchivedTaskCompletedForACancelledRecord` checks listed Tasks, pending Tasks, unknown Tasks, empty identity and an empty completed list. |

Additional evidence: `TestCancelRelinksLiveCitations` checks that a secondary
Spec's reference and core-artifact links resolve after cancellation, promoted
bytes equal the original source, and the record lists `promoted` and
`relinked`. It also exercises uncommitted live-document, occupied promotion,
duplicate basename and symlink refusals. `TestCancelWithoutATaskGraph` covers
both archive-root forms with empty Task lists.
`TestCancelLeavesOtherAdoptedSourcesInTheFolder` checks that another source
type remains named in the record and is never returned or promoted.

### Ownership and follow-up

Task status, checkboxes, the Task Graph, other Task files, `Archive`,
`BuildArchiveRecord`, History Sanitize and CLI files were not edited. The
Task's declared Verification commands, its CLI regression command and the
repository Verification remain for the Daemon. No commit, push or pull
request was made.

The PRD's Non-Goals still says Markdown-link rewriting is out of scope;
Requirement 9, ADR-0260 and the TechSpec Project Constraints explicitly
require Archive Relink for cancellation. This implementation follows that
explicit Task requirement; reconciling the PRD wording is a documentation
follow-up outside this slice.
