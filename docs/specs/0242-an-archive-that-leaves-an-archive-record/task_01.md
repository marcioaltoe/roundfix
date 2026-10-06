---
task: task_01
spec: 0242-an-archive-that-leaves-an-archive-record
status: pending
type: backend
complexity: high
---

# Task 01: An archive writes the Archive Record and removes the Spec folder, and the Delivery Queue accepts it

## Overview

Today `roundfix archive` stamps `_prd.md` and moves the Spec folder under the
History Root. This Task makes it write one Archive Record,
`<archive-root>/<slug>.md`, and remove the folder in the same change. The
folder's bytes stay in Git at the recorded `source_revision`. The Task adds
the resolver every later reader uses. The Delivery Queue's archive stage and
supersede ship here too, because their existing tests run the real Archive
Command. It answers the Backlog Entry "History keeps only what the
Secondbrain needs" of 2026-10-06.

## Requirements

1. MUST add `internal/spec/archive_record.go` with the shapes of
   `_techspec.md` → Interfaces: `ArchiveRecord`, `QAArchiveOverrideRecord`,
   `ArchiveRegeneration`, `ArchiveDisposition`, `ArchiveRecordPath`,
   `BuildArchiveRecord`, `RenderArchiveRecord`, `ParseArchiveRecord`,
   `ArchivedSpec`, `ReadArchivedSpec`, `ArchivedSpecSlugs` and
   `ErrNotArchived`, per Invariants 1, 2, 3 and 7.
2. MUST change `spec.Archive` per Invariants 4 and 5. It keeps every
   eligibility rule and refusal of today in today's order. A normal archive,
   a QA Archive Override and a superseded Spec without a Task Graph each
   write the record and remove the folder. No archive stamps the PRD or
   rewrites a link. `ArchiveRequest` gains `SourceRevision`, `Promote` and
   `RepositoryRoot`, and `ArchiveResult` gains `RecordPath`, `RemovedFiles`,
   `RemovedBytes` and `Promoted`.
3. MUST make the Archive Command pass `HEAD` as `SourceRevision`. It refuses
   a Spec folder that differs from `HEAD`, per Invariant 6 and
   `_techspec.md` → Surface Transcript 5. It prints the confirmation of
   Surface Transcript 1, and the override form of Surface Transcript 2
   without its promotion suffix, which task_04 adds with `--promote`. It
   keeps every other refusal and exit code.
4. MUST make the Delivery Queue's archive stage and its reconcile accept the
   exact retirement of Invariant 8: `Archive`, `archivePaths`,
   `archiveDiffIsExact`, `archiveCommitIsExact` and `reconcileArchiveCommit`
   in `internal/cli/deliver_workflow.go`. The already-archived check at the
   reviewed head finds a record. A commit that moves the folder keeps
   today's proof.
5. MUST make supersede's `knownDelivererSpec` accept a Spec whose
   `ReadArchivedSpec` succeeds, as it accepts an archived `_prd.md` today.
6. MUST update the existing tests that ran the real Archive Command and
   asserted a moved folder or a stamped PRD. Each asserts the record and the
   removed folder instead and keeps its other assertions: the archive, link,
   active-path, delivery, archived-retry, item-binary, QA partial and
   supersede tests named in Context. A test that proves today's link pass on
   a legacy folder keeps proving it through `ArchiveLinksMatch` and the
   legacy exact move.
7. MUST add the tests named in Verification:
   - in `internal/spec/archive_record_test.go`:
     - a passing Spec with adopted sources, ADRs in Decisions and a
       sanctioned regeneration;
     - an override with a 300-byte reason that keeps every override field;
     - a superseded Spec;
     - a 5,000-byte outcome that leaves the record at 2,048 bytes or less;
     - the render and parse round-trip;
     - each refusal of Invariant 4, with the tree byte-identical afterwards;
     - `ReadArchivedSpec` over a record, a legacy folder, both and neither;
   - in `internal/cli/archive_record_test.go`, Surface Transcripts 1 and 5,
     the override confirmation, and the removed bytes equal to `git show <source_revision>:<path>`
     for every removed file;
   - in `internal/cli/deliver_archive_record_test.go`:
     - the exact retirement;
     - a commit that keeps one file;
     - one that adds another path;
     - one whose record names another revision;
     - one whose promoted copy differs from its source;
     - a legacy move;
     - the archive stage committing the record end to end through the real
       command.
8. MUST update `docs/user-guide/commands/archive.md` so its usage, its
   description of the result and its confirmation line match the command.
   Keep its `--qa-override` text.
9. MUST NOT change the eligibility rules, the QA Report readers, ADR-0223's
   active-path refusal, the review, reconcile or Run-cause readers, or any
   file under `docs/history`.

## Subtasks

- [ ] Add the record, its builder, renderer, parser and resolver.
- [ ] Replace the move with the cut in `spec.Archive`.
- [ ] Pass the revision and refuse uncommitted changes in the Archive Command.
- [ ] Prove the exact retirement in the Delivery Queue and supersede.
- [ ] Update the archive tests that asserted the move and add the new ones.

## Acceptance Criteria

- [ ] Every archive disposition leaves one record of at most 2,048 bytes,
      with only the outcome shortened, and no Spec folder.
- [ ] Every removed byte reads back from Git at `source_revision`.
- [ ] The Delivery Queue accepts exactly the retirement Invariant 8
      describes and still accepts a legacy move.

## Context

- creates: `internal/spec/archive_record.go`
- creates: `internal/spec/archive_record_test.go`
- creates: `internal/cli/archive_record_test.go`
- creates: `internal/cli/deliver_archive_record_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/supersede.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `internal/spec/archive_test.go`
- interface: `internal/spec/archive_links_test.go`
- interface: `internal/cli/archive_test.go`
- interface: `internal/cli/archive_links_test.go`
- interface: `internal/cli/archive_active_spec_path_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_archive_links_test.go`
- interface: `internal/cli/deliver_archived_retry_test.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/cli/qa_partial_policy_test.go`
- interface: `internal/cli/supersede_test.go`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`
- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `internal/suiteguardcontract/regeneration.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder|TestArchiveRecordKeepsEveryOverrideField|TestArchiveRecordOfASupersededSpec|TestArchiveRecordStaysWithinTheTargetSize|TestArchiveRecordRoundTrips|TestArchiveRefusesBeforeAnyFileChanges|TestReadArchivedSpecReadsRecordAndLegacyFolder)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder TestArchiveRecordKeepsEveryOverrideField TestArchiveRecordOfASupersededSpec TestArchiveRecordStaysWithinTheTargetSize TestArchiveRecordRoundTrips TestArchiveRefusesBeforeAnyFileChanges TestReadArchivedSpecReadsRecordAndLegacyFolder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestArchiveCommandLeavesTheArchiveRecord|TestArchiveCommandRefusesUncommittedSpecChanges|TestArchiveCommandRemovedBytesStayInGit|TestDeliveryAcceptsAnExactRetirement|TestDeliveryRefusesAnInexactRetirement|TestDeliveryKeepsTheLegacyExactMove|TestDeliveryArchiveStageCommitsTheArchiveRecord|TestSupersedeAcceptsAnArchiveRecord)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveCommandLeavesTheArchiveRecord TestArchiveCommandRefusesUncommittedSpecChanges TestArchiveCommandRemovedBytesStayInGit TestDeliveryAcceptsAnExactRetirement TestDeliveryRefusesAnInexactRetirement TestDeliveryKeepsTheLegacyExactMove TestDeliveryArchiveStageCommitsTheArchiveRecord TestSupersedeAcceptsAnArchiveRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eight tests do not exist, so the command fails.

## References

- `_prd.md` → Goal 1; Goal 2; User Story 1; User Story 3; Core Feature 1; Core Feature 2; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; Invariants 1-8; Data Models; API Contract 1; API Contract 2; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 5; Build Order 1
- ADR-0247; ADR-0230; ADR-0154; ADR-0223
