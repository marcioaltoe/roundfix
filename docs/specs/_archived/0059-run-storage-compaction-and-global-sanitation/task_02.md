---
task: task_02
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
type: backend
complexity: high
---

# Task 02: Compact only when nothing can be writing

## Overview

`internal/store` carries no `VACUUM` — zero occurrences. Retention deletes rows
and the pages stay allocated, which is how the Run Database reached 1.4 GB with
its journal already pruned.

The guard is the feature here, not the reclaimed bytes. A refusal that is wrong
costs a retry; a permission that is wrong costs the database.

## Requirements

1. MUST add an explicit compaction with a deterministic preview reporting bytes
   before, reclaimable, and projected after.
2. MUST refuse while an Active Run exists or any other writer can be present,
   naming what blocked it.
3. MUST verify sufficient temporary disk capacity before starting, and refuse
   naming the shortfall rather than failing partway.
4. MUST be failure-safe: an interrupted compaction leaves a usable database.
5. MUST be explicit — never an automatic side effect of a retention sweep — and
   MUST take no exclusive lock outside its own run, so the operational prune on
   `implement`, `resolve`, and `watch` startup stays cheap and lock-free.
6. MUST NOT change what retention deletes, which ADR-0033 owns.
7. MUST assert that the reclaimed bytes equal the preview's reclaimable bytes
   within a declared tolerance, and MUST NOT assert any recorded size.
8. MUST assert that a refusal leaves the file size unchanged, compared before
   and after rather than against a constant.

## Subtasks

- [ ] Add the preview and the compaction behind it.
- [ ] Add the Active-Run and writer guard, and the capacity check.
- [ ] Assert preview-equals-result and refusal-preserves-bytes.

## Acceptance Criteria

- [ ] A preview reports bytes before, reclaimable, and projected after.
- [ ] Compaction reclaims bytes equal to the preview's reclaimable within the
      declared tolerance, asserted as a relation.
- [ ] An injected Active Run makes compaction refuse, naming that Run.
- [ ] A refusal leaves the database file size unchanged.
- [ ] Insufficient temporary capacity refuses naming the shortfall, before any
      mutation.
- [ ] Compaction never runs as a side effect of a retention sweep, asserted.
- [ ] Retention behaviour is unchanged, asserted over the existing tests.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/journal.go`
- instruction: `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/store -count=1 -run 'Compact|Preview|Vacuum|Capacity' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the compaction tests ran and passed.
- `go test ./internal/store ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Feature 1; User Stories 1 and 2; Success Metric 1.
- `_techspec.md` → System Architecture; Build Order 2; Risks & Considerations.
- ADR-0033, ADR-0052.

## Result

Implemented explicit Run Database preview and compaction in `internal/store`.
Preview builds a temporary compact SQLite snapshot and reports measured bytes
before, reclaimable, and projected after. Apply acquires SQLite exclusive
locking mode with no lock wait, rejects any non-terminal Run, revalidates the
preview fingerprint, checks the documented two-database temporary-capacity
bound, then runs transactional `VACUUM` and checkpoints the result. Cancellation
cannot skip restoration of the connection's normal locking mode and busy
timeout.

Focused checks run after the final implementation edits:

- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/store -run '^Test(CompactionPreviewMatchesResultWithinDeclaredTolerance|CompactRefusesActiveRunAndPreservesDatabaseBytes|CompactRefusesAnotherWriterAndPreservesDatabaseBytes|CompactRefusesInsufficientTemporaryCapacityBeforeMutation|CompactRefusesStalePreviewAndPreservesDatabaseBytes|PruneTerminalRunsNeverVacuumsAllocatedPages)$' -count=1 -v`
  — passed (`6` real SQLite tests).
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/store -run '^TestPruneTerminalRuns' -count=1`
  — passed the existing retention suite plus the no-automatic-compaction case.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/store -count=1`
  — exited `0` for the complete store package.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -race ./internal/store -run '^Test(CompactionPreviewMatchesResultWithinDeclaredTolerance|CompactRefusesActiveRunAndPreservesDatabaseBytes|CompactRefusesAnotherWriterAndPreservesDatabaseBytes|CompactRefusesInsufficientTemporaryCapacityBeforeMutation|CompactRefusesStalePreviewAndPreservesDatabaseBytes|PruneTerminalRunsNeverVacuumsAllocatedPages)$' -count=1`
  — exited `0` with the focused compaction and retention cases under the race
  detector.
- `rtk env GOOS=windows GOCACHE=/private/tmp/roundfix-task02-gocache go build ./internal/store`
  and `rtk env GOOS=linux GOCACHE=/private/tmp/roundfix-task02-gocache go build ./internal/store`
  — both exited `0`, covering the platform-specific capacity measurement.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go vet ./internal/store`
  — exited `0`.
- `rtk git -c core.fsmonitor=false diff --check` — exited `0`.

Acceptance evidence:

- Deterministic preview: `TestCompactionPreviewMatchesResultWithinDeclaredTolerance`
  asserts all three byte fields are measured and that `after = before -
  reclaimable`, without a recorded-size assertion.
- Preview equals result: the same test compares measured reclaimed bytes with
  preview reclaimable bytes and bounds the difference by the preview's declared
  one-SQLite-page tolerance. It also checks the completion report against the
  resulting database file and runs `PRAGMA quick_check`, which returns `ok`.
- Active Run refusal: `TestCompactRefusesActiveRunAndPreservesDatabaseBytes`
  injects an Active Run and asserts the typed refusal and diagnostic both name
  that Run ID.
- Refusal preserves bytes: the Active Run, second-writer, insufficient-capacity,
  and stale-preview cases each compare the database file size immediately
  before and after refusal and assert equality rather than a constant.
- Capacity refusal: `TestCompactRefusesInsufficientTemporaryCapacityBeforeMutation`
  injects one byte less than the measured requirement and asserts required,
  available, and shortfall bytes in the typed error before the original
  database changes.
- No automatic compaction: `TestPruneTerminalRunsNeverVacuumsAllocatedPages`
  proves a retention sweep deletes eligible journal rows while the database
  file size stays unchanged and the freed pages remain on SQLite's freelist for
  an explicit compaction.
- Retention unchanged: every existing `TestPruneTerminalRuns*` case passes,
  preserving ADR-0033's terminal cutoff, Run-row, Active Run, and lock behavior.

The commands under `## Verification` were not run; the Daemon owns those
commands and Task settlement.
