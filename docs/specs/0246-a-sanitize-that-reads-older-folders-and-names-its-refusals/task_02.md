---
task: task_02
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: completed
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

## Result

Implemented the task_02 slice for Daemon Verification. `ReadLegacyCauseGraph`
shares all manifest front-matter validation with the strict reader and reports
the two permitted projection tolerances in table order. A row outside the
graph with a retired type is named once, as outside the graph. Malformed and
duplicate rows retain the strict reader's messages.

`ArchiveRecordInput.Legacy` selects lenient QA Task reading only for legacy
conversion and the folder branch of `ReadArchivedSpec`. The Archive Command
continues to use the default strict input. `LegacyConversion.Tolerated` holds
the reader's tolerances without adding them to the Archive Record. Disposition
logic and CLI behavior are unchanged.

Acceptance evidence:

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Legacy folders convert with named tolerances | `TestLegacyConversionToleratesAProjectionRowOutsideTheGraph` and `TestLegacyConversionToleratesARetiredTaskType` exercise plan, apply, and record reading; assert exact tolerance wording/order, QA identity, and absence of tolerances from the rendered record. `TestReadArchivedSpecReadsALegacyFolderLeniently` covers both tolerated shapes together and verifies graph membership and unchanged folder bytes. |
| Active Specs refuse exactly as before | `TestActiveSpecStillRefusesWhatLegacyReadingTolerates` asserts exact errors from `Load`, `ReadCauseGraph`, `ReadCauseGraphAt`, and `BuildArchiveRecord` without `Legacy`. `TestLegacyConversionStillRefusesAMalformedProjectionRow` verifies the exact refusal and unchanged folder. `TestLegacyCauseGraphPreservesManifestRefusals` compares strict and legacy errors for 13 front-matter and projection cases, including duplicate rows outside the graph. |
| Existing spec tests pass unchanged | `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec -count=1` exited 0 (`ok roundfix/internal/spec`), including the new tests. No existing test or expectation was edited. |

Focused checks:

- Before production changes, `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec -run '^TestLegacyConversionTolerates' -count=1`
  exited 1: `LegacyConversion.Tolerated` and `ReadLegacyCauseGraph` were absent.
- After the implementation, `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test ./internal/spec -run 'Test(LegacyConversion|ActiveSpecStill|ReadArchivedSpecReadsALegacy)' -count=1`
  exited 0.
- After adding refusal-parity coverage, the complete spec package check above
  exited 0.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- The first `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  exited 2. Formatting and `go vet ./...` passed. Two CLI process-owner tests
  could not enumerate the process table in the sandbox. Suiteguard also
  detected this Agent's concurrent Result edit; that edit was made while the
  suite was active and invalidated the repository fingerprint checks. The
  rerun uses process-table permission and keeps the worktree unchanged until
  the check exits.
- The rerun of `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  with host process-table access exited 0: formatting, vet, the repository Go
  tests, skill synchronization/checks, and build passed. No worktree edit was
  made during that run. The Result was updated after the check exited.

The authored `## Verification` command was not run; the Daemon owns that
command and Task settlement. Status and task checkboxes were preserved. No
Task Graph, other Task file, disposition logic, CLI, or existing expectation
was edited. No commit, push, or Pull Request was made. No follow-up work was
identified within this slice.
