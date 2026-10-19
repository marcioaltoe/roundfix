---
task: task_03
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
status: pending
type: backend
complexity: high
---

# Task 03: The storage report measures one Storage Snapshot and groups every page

## Overview

On 2026-10-18 `roundfix storage report` failed while the Delivery Queue ran
(`_prd.md` → Sources 2). The difference it reported, 1,363,968 bytes, is
exactly the 332 pointer-map pages and the lock-byte page that `dbstat` never
reports. On 2026-10-19 the report failed the same way on an idle store. It
also reads the Run Database as an immutable file, which ignores the
write-ahead log while a writer appends to it.

This Task makes the report pick its reader the way `roundfix gc --dry-run`
does, measure inside one read transaction, and group pointer-map and lock-byte
pages, so the groups reconcile with the snapshot exactly. It is verifiable on
its own through store and command tests on temporary stores.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-storage-report-fails-while-the-run-database-wal-is-live.md`
   (`_prd.md` → Sources 2), by implementing API Contracts 7 to 9 and the
   Surface Transcript 1, "storage report without a Run Database", of
   `_techspec.md` (ADR-0296). The tolerance MUST NOT be widened.
2. MUST share one Storage Snapshot reader choice in `internal/cli/gc.go` (API
   Contract 7).
   - It returns the immutable storage reader when `<database>-wal` is absent
     or 32 bytes or fewer, and `store.OpenReader` otherwise. A stat error
     other than not-exist is returned.
   - `openGCStore`'s dry run and `measureStorageReport` both use it.
3. MUST measure the page groups in `StorageReport` in
   `internal/store/journal.go` (API Contract 8).
   - Every query runs on one connection inside one read transaction:
     `page_size`, `page_count`, `freelist_count`, `auto_vacuum`, `dbstat` and
     the Run rows.
   - The usable size is the page size minus the reserved-bytes field at offset
     20 of the main file's header.
   - Pointer-map positions are `2 + k·(usable/5 + 1)` when `auto_vacuum` is
     not 0, with a position that equals the lock-byte page moved one page
     later. The lock-byte page is `1073741824/page_size + 1`. Each counts only
     up to `page_count`.
   - Snapshot bytes are `page_count·page_size`. Grouped bytes are table, free,
     pointer-map and lock-byte bytes. Any difference refuses with
     `reconcile Run Database storage groups: snapshot bytes=<n> grouped
     bytes=<n> difference=<n>`.
   - `ReconciliationToleranceBytes` is 0, and
     `StorageReportReconciliationToleranceReason` becomes `one read
     transaction measures every page group, pointer-map and lock-byte pages
     included`.
   - Add `DatabaseWALBytes`, `DatabaseSharedMemoryBytes`,
     `DatabaseSnapshotBytes`, `DatabasePointerMapBytes` and
     `DatabaseLockBytePageBytes`. `DatabaseBytes` stays the main file's size;
     the `-wal` and `-shm` files are stat'ed, and an absent file is 0.
   - The report writes nothing, migrates nothing, checkpoints nothing and
     creates no sidecar file.
4. MUST print the Database section of API Contract 9 from `printStorageReport`
   in `internal/cli/gc.go`. Without a Run Database, the output is the Surface
   Transcript 1 exactly. Every other section is unchanged. `gc compact`'s
   live-writer fallback keeps reading `DatabaseBytes` and `DatabaseFreeBytes`.
5. MUST add these tests. Each calls `t.Parallel()` first, compiles for
   Windows, uses temporary stores and homes, and never reads `~/.roundfix`.
   - `internal/store/journal_test.go`:
     - `TestStoragePageStructureMatchesTheFileFormat`, a table test of the
       unexported page-structure helper. Pages of 4096 bytes: 271,903 → 332
       pointer-map pages plus the lock-byte page; 275,136 → 336 plus it; 3,013
       → 4 and none; 5,020 → 7 and none. `auto_vacuum` 0 gives no pointer-map
       page. Page size 1024 moves the pointer-map page that would fall on page
       1,048,577 to 1,048,578.
     - `TestStorageReportGroupsPointerMapPages`, a store of more than 1,700
       pages with incremental auto-vacuum. The report succeeds, at least two
       pointer-map pages are grouped, and the grouped bytes equal the snapshot
       bytes.
     - `TestStorageReportReadsALiveWALSnapshot`. A writer that stays open
       holds WAL frames not yet checkpointed, and the report reads through
       `OpenReader`. The snapshot bytes exceed the main file's bytes,
       `DatabaseWALBytes` is positive, and the database and WAL files are
       byte-identical before and after.
     - `TestStorageReportReconcilesMeasuredTotalsWithoutMutation` keeps its
       name. Its tolerance assertions change to 0 and the new reason, and its
       reconciliation compares snapshot and grouped bytes (declared break,
       ADR-0296).
   - `internal/cli/gc_test.go`:
     - `TestRunStorageReportWithoutDatabaseTranscript` asserts Surface
       Transcript 1 whole: every standard output line in order, empty standard
       error and exit 0, with `<home>` the only variation.
     - `TestRunStorageReportWhileAWriterHoldsTheWAL`. With a `store.Open`
       writer still open after its writes, `storage report` exits 0 and
       prints a positive `WAL bytes:` line and `Reconciliation difference:
       0`.
     - `TestRunStorageReportOutsideGitRepository` keeps its name. Its expected
       "one SQLite page" becomes "one read transaction measures every page
       group" (declared break).

   Every other existing test MUST keep its name and pass, including the
   `gc compact`, `gc --dry-run` and Doctor storage tests.
6. MUST document the change:
   - `docs/user-guide/commands/storage-report.md` and
     `.agents/skills/roundfix/references/storage.md` each name the
     `Pointer-map bytes`, `Lock-byte page bytes`, `WAL bytes`,
     `Shared-memory bytes` and `Snapshot bytes` lines.
   - Each says the report measures one Storage Snapshot, which reads the
     write-ahead log when it holds frames, and that the grouped bytes equal
     the snapshot bytes exactly.
   - The guide's sentence on `gc compact`'s fallback replaces "falls back to
     immutable storage-report measurements" with "falls back to the
     storage report's Storage Snapshot".

   No text inside any `### QA settlement` section may change.
7. MUST add **Storage Snapshot** to `CONTEXT.md` through domain-modeling, with
   ADR-0296's link. Its definition contains "measures every page group".
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of each edited owned skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`)
   by one patch over the value the tree holds when this Task starts,
   because the record command refuses changed content under an unchanged
   version. The Agent runs these commands itself before handing back: the
   Daemon only runs the Verification and never regenerates a mirror or a
   record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

## Subtasks

- [ ] Share the Storage Snapshot reader choice.
- [ ] Measure in one read transaction and add the pointer-map and lock-byte groups.
- [ ] Print the new Database lines.
- [ ] Add the store and command tests and change the two declared assertions.
- [ ] Document the report in the guide, the reference and the glossary; sync, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] The report exits 0 with `Reconciliation difference: 0` on a store with several pointer-map pages and on a store with a live WAL.
- [ ] The page-structure counts match the 2026-10-18 and 2026-10-19 measurements.
- [ ] The report changes no database or WAL byte and creates no sidecar file.

## Consumers

- `store.StorageReport` and its new fields — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0296-the-storage-report-reconciles-every-page-of-one-storage-snapshot.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0172-doctor-reports-reclaimable-run-storage-from-cheap-reads.md`
- instruction: `internal/store/store.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/gc_test.go`
- interface: `docs/user-guide/commands/storage-report.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestStoragePageStructureMatchesTheFileFormat|TestStorageReportGroupsPointerMapPages|TestStorageReportReadsALiveWALSnapshot|TestStorageReportReconcilesMeasuredTotalsWithoutMutation)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestStoragePageStructureMatchesTheFileFormat TestStorageReportGroupsPointerMapPages TestStorageReportReadsALiveWALSnapshot TestStorageReportReconcilesMeasuredTotalsWithoutMutation; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestRunStorageReportWithoutDatabaseTranscript|TestRunStorageReportWhileAWriterHoldsTheWAL|TestRunStorageReportOutsideGitRepository)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunStorageReportWithoutDatabaseTranscript TestRunStorageReportWhileAWriterHoldsTheWAL TestRunStorageReportOutsideGitRepository; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/storage-report.md .agents/skills/roundfix/references/storage.md; do for name in 'Pointer-map bytes' 'Lock-byte page bytes' 'WAL bytes' 'Shared-memory bytes' 'Snapshot bytes' 'Storage Snapshot'; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "$name" || { printf 'missing %s in %s\n' "$name" "$f" >&2; exit 1; }; done; done; ! tr -s '[:space:]' ' ' < docs/user-guide/commands/storage-report.md | grep -qF -- 'falls back to immutable storage-report measurements' || { printf 'storage guide still names the immutable fallback\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Storage Snapshot**:' || { printf 'missing glossary term Storage Snapshot\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'measures every page group' || { printf 'missing Storage Snapshot definition\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/store || exit 1; go test -count=1 -run 'Storage|GC|Doctor|Compact|Sanitize' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the page structure, pointer-map and live WAL store tests, the transcript and command tests, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 3; Core Feature 4; Core Feature 5; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 7; API Contract 8; API Contract 9; API Contract 10; Surface Transcript 1; Invariants 3-4; Current behavior; Testing Approach; Build Order 3
- ADR-0296; ADR-0255; ADR-0172; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285; ADR-0271
