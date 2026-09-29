---
type: fix
status: open
created: 2026-09-29
spec: null
reason: null
---

# `gc` reports Runs it already emptied as pruned

## Symptom

`roundfix gc` on 2026-09-29 reported `Runs pruned: 13` with `Journal rows removed: 0`, `Artifact bytes reclaimed: 0` and `Orphan artifact dirs removed: 0`. A `gc --dry-run` right afterwards reported the same 13 Runs as eligible. Every later run reports them again, so the count never says whether anything was actually reclaimed.

## Where

`terminalRunPruneCandidates` in `internal/store/journal.go` selects every terminal Run older than the retention cutoff, whether or not any journal rows or artifacts remain. The report in `internal/cli/gc.go` prints that count.

## Expected

`gc` and `gc --dry-run` count only Runs that still hold journal rows or artifacts to reclaim. A second `gc` after a complete one reports zero.

## Evidence

`bin/roundfix gc` followed by `bin/roundfix gc --dry-run` on `a626494c`, 2026-09-29.
