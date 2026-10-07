---
task: task_01
spec: 0243-a-history-that-holds-only-records
status: completed
type: backend
complexity: high
---

# Task 01: A Legacy Archive Folder converts to an Archive Record that Git still backs

## Overview

Spec 0242's archive writes an Archive Record for a Spec it retires now. This
Task converts a Spec folder an earlier archive left under the archive root,
a Legacy Archive Folder, into the same record through the same
`BuildArchiveRecord`. The record's `source` is the folder's own path and its
`source_revision` is the revision the caller passes, which the command sets
to `HEAD`. The Task also reads the delivery commit and Pull Request from Git,
and names a folder that never had a QA Report `no-qa` instead of inventing a
verdict. It answers the Backlog Entry "History keeps only what the
Secondbrain needs" of 2026-10-06, Expected 3, and the list "What Spec 0243
must do" in Spec 0242's TechSpec.

## Requirements

1. MUST add the disposition `ArchiveNoQA` (`no-qa`) to
   `internal/spec/archive_record.go` per Invariant 3. `BuildArchiveRecord`
   sets it only when no other rule set a disposition. `ParseArchiveRecord`
   accepts it. No other disposition rule, field or key order changes, and
   every existing test of the record keeps passing unchanged.
2. MUST add `LegacyArchiveFolders`, `FindLegacyDelivery`,
   `PlanLegacyConversion` and `ApplyLegacyConversion` in
   `internal/spec/history_sanitize.go` with the signatures of
   `_techspec.md` → Interfaces and the behavior of Invariants 1 and 4-6.
   - A folder whose `_prd.md` says `status: active` converts like any other.
   - The record's `source` is the folder's repository-relative path, and
     `archived` falls back to the delivery date only when the PRD has no
     `archived` stamp.
   - The plan renders, parses back and refuses a record that does not
     round-trip. It writes nothing.
   - Promotion refusals follow Invariant 4. The promotion destination is
     `docs/references/<basename>` under `RepositoryRoot`.
3. MUST read Git only through the local `git` binary with the repository
   root passed explicitly, with rename detection off where Invariant 5 says
   so. Data comes from the repository's own first-parent history, and the
   result is written only into the record.
4. MUST name every unexported helper this Task adds with a `legacy` prefix,
   because task_02 adds files to the same package in parallel.
5. MUST add the tests named in Verification to
   `internal/spec/history_sanitize_test.go`. Each builds its own temporary
   Git repository with `internal/gittest`, never reads this repository's
   `docs/history`, and leaves the repository unchanged under the package's
   suite guard:
   - a passing folder whose record is at most 2,048 bytes, round-trips, and
     whose removed files read back byte-identical from
     `git show <source_revision>:<source>/<path>`;
   - a folder without a QA Report, override or supersession is `no-qa`;
   - an override folder keeps all six override fields;
   - a pre-stamp folder (`status: active`) converts and takes `archived`
     from the delivery date;
   - delivery: a commit that removed the folder from the Spec Root, a commit
     that added it directly under the archive root with a `(#<n>)` subject,
     and a relocation commit that yields empty fields;
   - each promotion refusal, and a promotion copied byte-identically;
   - a failure before removal leaves the folder whole and no record;
   - `LegacyArchiveFolders` skips records and directories without `_prd.md`;
   - a `no-qa` record round-trips through `RenderArchiveRecord` and
     `ParseArchiveRecord`.
6. MUST NOT change `internal/spec/archive.go`, `ReadArchivedSpec`,
   `ReadArchivedSpecAt`, any reader of archived Specs, or any file under
   `docs/history`.

## Subtasks

- [ ] Add the no-qa disposition to the record.
- [ ] Add folder listing, delivery discovery, plan and apply.
- [ ] Add the fixture-repository tests.

## Acceptance Criteria

- [ ] A Legacy Archive Folder converts to a record that round-trips and that
      Git still backs at `source_revision`.
- [ ] A folder without QA is `no-qa`, and a pre-stamp folder converts.
- [ ] The delivery commit is found or left empty, never guessed from a
      relocation.

## Context

- interface: `internal/spec/archive_record.go`
- creates: `internal/spec/history_sanitize.go`
- creates: `internal/spec/history_sanitize_test.go`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `internal/spec/archive.go`
- instruction: `internal/spec/archive_reader.go`
- instruction: `internal/gittest/gittest.go`

## Verification

- `out="$(go test -count=1 -v ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLegacyArchiveFoldersSkipRecords TestLegacyConversionWritesARecordGitStillHolds TestLegacyConversionOfAFolderWithoutQAIsNoQA TestLegacyConversionKeepsEveryOverrideField TestLegacyConversionAcceptsAPreStampFolder TestLegacyDeliveryNamesTheCommitThatRetiredTheSpec TestLegacyDeliveryIgnoresARelocation TestLegacyConversionPromotionRefusals TestLegacyConversionLeavesTheFolderWholeOnFailure TestNoQARecordRoundTrips; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the ten tests exists, so the command fails; after it the whole `internal/spec` package passes with them.
- `grep -qF -- 'ArchiveNoQA' internal/spec/archive_record.go && grep -qF -- '"no-qa"' internal/spec/archive_record.go && grep -qF -- 'func PlanLegacyConversion(' internal/spec/history_sanitize.go && grep -qF -- 'func FindLegacyDelivery(' internal/spec/history_sanitize.go && go build -buildvcs=false ./internal/spec` — expected: exit 0; before this Task neither the disposition nor `internal/spec/history_sanitize.go` exists, so the command fails.

## References

- `_prd.md` → Goal 3; User Story 1; Core Feature 3; Success Metric 2
- `_techspec.md` → Measured inventory; Interfaces; Invariants 1, 3, 4, 5 and 6; Data Models; Build Order 1
- ADR-0248; ADR-0247; ADR-0154

## Result

Implemented this Task's slice for Daemon Verification; status remains
Daemon-owned.

- Added `ArchiveNoQA` and parser support. The builder uses it only after QA,
  override and supersession rules leave the disposition unset, and accepts
  a pre-stamp folder without a Task Graph.
- Added sorted legacy-folder discovery, first-parent Git delivery discovery,
  read-only conversion planning and conversion application. Plans use the
  folder's repository-relative source and the caller's revision, preserve an
  existing archive stamp, and otherwise use the delivery's author date.
- Plans validate promotion containment, regular files, core-artifact
  exclusions and destination collisions, then compare rendered and parsed
  metadata and verify stable rendering. The existing renderer still owns
  outcome shortening. Application creates outputs exclusively, copies
  promotion bytes, and rolls back written files on failure before removal.
- Added the ten required fixture-repository tests and two additional tests
  for metadata round-trip refusal and the renderer's outcome budget. Every
  added unexported helper has a `legacy` prefix.

### Acceptance evidence

| Acceptance criterion | Focused evidence |
| --- | --- |
| A Legacy Archive Folder converts to a record that round-trips and Git still backs | `TestLegacyConversionWritesARecordGitStillHolds` passed: planning leaves the folder unchanged, inventory counts match, the passing record is at most `ArchiveRecordTargetBytes`, application removes the folder, its promotion is byte-identical, and every removed file reads back byte-identically from `git show <source_revision>:<source>/<path>`. `TestLegacyConversionLeavesTheFolderWholeOnFailure` passed after a second promotion failed, proving rollback of the record and first promotion. All promotion-refusal cases passed. |
| A folder without QA is `no-qa`, and a pre-stamp folder converts | `TestLegacyConversionOfAFolderWithoutQAIsNoQA`, `TestNoQARecordRoundTrips` and `TestLegacyConversionAcceptsAPreStampFolder` passed. The pre-stamp fixture has `status: active` and receives the delivery date. `TestLegacyConversionKeepsEveryOverrideField` passed, retaining the override boolean and all five associated metadata fields. |
| Delivery is found or empty, never guessed from relocation | `TestLegacyDeliveryNamesTheCommitThatRetiredTheSpec` passed for Spec Root deletion, direct archive addition, merge retirement and the newest retirement. It checks the commit, Pull Request suffix and author date. `TestLegacyDeliveryIgnoresARelocation` passed with every delivery field empty; it also verifies context cancellation. |

### Checks run

- Initial inspection found no legacy-conversion implementation or required
  tests. The first focused run exposed an absolute-versus-relative path
  error in the new fixture adapter; the adapter was corrected.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk proxy go test ./internal/spec -run '^(TestLegacy|TestNoQA|TestArchiveRecord|TestArchiveWrites)' -count=1 -v`
  — exit 0 after the final code edit; all 12 new tests and six selected
  existing Archive Record tests passed. The suite guard reported no
  repository mutation.
- `GOCACHE=/private/tmp/roundfix-task01-gocache rtk make verify-incremental`
  — the sandboxed attempt passed `internal/spec` but exited 2 because two
  existing CLI force-stop integration tests could not read the process
  table. The same command rerun with the required process-table permission
  exited 0, including formatting, vet, package tests, skill checks and build.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

The declared Verification commands were not run. No commit, push or Pull
Request was made. No changes were made to the Task Graph, another Task,
archive execution, archived-Spec readers or `docs/history`. The initial
Daemon change from `pending` to `in_progress` was preserved. No follow-up
implementation was added to this slice.

## Carry-forward provenance

- Source Run: `run_20261007T040343Z_bb67f44effbee11e`
- Source commit: `b44eeb1c7aaf13b7c89a023fae00e8952480c5c9`
