---
task: task_02
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
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
