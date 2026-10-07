---
task: task_02
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: pending
type: backend
complexity: medium
---

# Task 02: A Legacy Archive Folder's Task Graph is read leniently

## Overview

Lets a Legacy Archive Folder become an Archive Record when its projection
table names a Task the author took out of the graph, or a Task type the closed
set no longer holds. The current schema stays enforced for active Specs. This
Task answers reports 1 to 3 of the Backlog Entry "History sanitize refuses
older archived Specs and aborts the whole plan" of 2026-10-07. It is
verifiable on its own through the spec tests below.

## Requirements

1. MUST add `ReadLegacyCauseGraph(root, slug string) (CauseGraph, []string, error)`.
   It reads `<root>/<slug>/_tasks.md` with every front-matter check of
   `parseManifestNodes` and the projection rules of `_techspec.md` →
   Invariant 2. It returns the tolerances in table order, worded exactly
   `projection row <id> names a Task outside the graph` and
   `projection row <id> has retired Task type "<type>"`.
2. MUST keep `parseManifestNodes`, `spec.Load`, `ReadCauseGraph` and
   `ReadCauseGraphAt` byte-identical in behavior and messages
   (`_techspec.md` → Invariant 3). A malformed or duplicate projection row
   refuses in the legacy reader too, with today's message.
3. MUST add `Legacy bool` to `ArchiveRecordInput`. `BuildArchiveRecord` reads
   the QA Task through `ReadLegacyCauseGraph` when it is set and through
   `ReadCauseGraph` otherwise (`_techspec.md` → Invariant 1).
4. MUST set `Legacy` in `PlanLegacyConversion` and in the folder branch of
   `ReadArchivedSpec`, and MUST NOT set it in the Archive Command.
5. MUST add `Tolerated []string` to `LegacyConversion`, filled from
   `ReadLegacyCauseGraph` and never written into the record (`_techspec.md` →
   Invariant 4).
6. MUST add `internal/spec/history_sanitize_legacy_test.go` with synthetic
   folders only, never adopter content:
   - `TestLegacyConversionToleratesAProjectionRowOutsideTheGraph`: graph nodes
     commented out of the front matter but still in the table, so the folder
     converts and `Tolerated` names each row.
   - `TestLegacyConversionToleratesARetiredTaskType`: a `refactor` row
     converts and is named.
   - `TestLegacyConversionStillRefusesAMalformedProjectionRow`: a row with
     fewer than five cells refuses, and the folder stays whole.
   - `TestActiveSpecStillRefusesWhatLegacyReadingTolerates`: an active Spec
     with both rows still fails `spec.Load` and `ReadCauseGraph` with today's
     messages, and fails `BuildArchiveRecord` without `Legacy`.
   - `TestReadArchivedSpecReadsALegacyFolderLeniently`: the folder branch
     returns a record for a folder with both rows.
7. MUST NOT change the disposition logic, which is task_03's, the CLI, or any
   existing test expectation.

## Subtasks

- [ ] Add the lenient manifest reader and its tolerances.
- [ ] Thread `Legacy` through the record builder and its two legacy callers.
- [ ] Record the tolerances on the conversion.
- [ ] Write the spec tests.

## Acceptance Criteria

- [ ] A Legacy Archive Folder with a row outside the graph or a retired type
      converts, and its tolerances are named.
- [ ] An active Spec with the same rows refuses exactly as before.
- [ ] Every existing spec test passes unchanged.

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/spec/cause_graph.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/history_sanitize.go`
- creates: `internal/spec/history_sanitize_legacy_test.go`
- instruction: `internal/spec/history_sanitize_test.go`
- instruction: `internal/spec/archive.go`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacyConversionToleratesAProjectionRowOutsideTheGraph|TestLegacyConversionToleratesARetiredTaskType|TestLegacyConversionStillRefusesAMalformedProjectionRow|TestActiveSpecStillRefusesWhatLegacyReadingTolerates|TestReadArchivedSpecReadsALegacyFolderLeniently)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestLegacyConversionToleratesAProjectionRowOutsideTheGraph TestLegacyConversionToleratesARetiredTaskType TestLegacyConversionStillRefusesAMalformedProjectionRow TestActiveSpecStillRefusesWhatLegacyReadingTolerates TestReadArchivedSpecReadsALegacyFolderLeniently; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the five tests exists. After it, a legacy folder converts with named tolerances, and an active Spec refuses as before.

## References

- `_prd.md` → Core Feature 1; Goals; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 1; Invariant 1; Invariant 2; Invariant 3; Invariant 4; Build Order 2
- ADR-0251
