---
spec: 0014-run-store-retention
prd: _prd.md
created: 2026-07-06
---

# Run Store Retention — Technical Spec

## Executive Summary

Bound the Run Event Journal without disturbing the Run state that powers
one-active-run, stop, detach, and recovery. The primary trade-off is what
retention deletes: only `run_events` rows and on-disk artifact directories of
**terminal** Runs older than a window — never `runs` rows or active-run locks —
so cleanup can never break a live or recoverable Run (ADR-0033). Pruning is a
single store operation reused by two callers: an on-demand GC Command (with
`--dry-run`) and a best-effort pass in the existing preflight sweep, so storage
self-heals. Row deletion (not file vacuuming) is the mechanism; SQLite reclaims
pages over time and a VACUUM is a deferred non-goal.

## System Architecture

- `internal/store` — a `PruneTerminalRuns(ctx, olderThan)` operation that
  deletes `run_events` for terminal Runs whose `completed_at` is older than the
  cutoff and returns what it removed (run ids, row counts). Active Runs (no
  terminal state / empty `completed_at`) are excluded by query. `runs` rows and
  `active_run_locks` are never deleted.
- `internal/config` — a new `store` section with `journal_retention` (duration;
  `0` = keep everything).
- `internal/cli` — a new `gc` command that calls the store prune plus artifact
  directory cleanup and reports freed counts; the implement/resolve/watch
  preflight sweep gains a best-effort prune call when retention is non-zero.
- Artifact cleanup — remove each pruned Run's `<artifact_dir>/runs/<run-id>`
  directory, and orphaned `runs/<id>` directories with no matching `runs` row.
- No journal format, cursor, replay, or lock changes.

## Implementation Design

### Interfaces

```go
// internal/store — prune journals of terminal Runs past the cutoff.
type PruneResult struct {
    RunIDs   []string // terminal Runs whose journal was pruned
    Events   int      // run_events rows deleted
}
// Deletes run_events for runs in a terminal state with completed_at < cutoff.
// Never deletes runs rows or active_run_locks; Active Runs are excluded.
func (s *Store) PruneTerminalRuns(ctx context.Context, cutoff time.Time) (PruneResult, error)
```

```go
// internal/cli — GC Command.
type gcRequest struct{ dryRun bool }
// Resolves retention from config; computes cutoff = now - journal_retention;
// dry-run: report the eligible set (from a read-only query) and change nothing.
// live: PruneTerminalRuns + remove each pruned Run's artifact dir + remove
// orphaned runs/<id> dirs; report Runs pruned, rows removed, bytes reclaimed.
```

### Journal Retention config

A `store` section with `journal_retention` as a Go duration string; default
**`336h` (14 days)**; `0` disables pruning. It joins the strict-decoded config
with Project > User > builtin precedence like every other key.

### Prune scope (ADR-0033)

The prune query selects only Runs in a terminal state (Clean, Unresolved,
Failed, IntegrationPending, Stopped, BudgetExceeded, TimedOut) with a
`completed_at` older than the cutoff, deletes their `run_events` rows, and
returns their ids for artifact cleanup. Active Runs (non-terminal state, or an
empty `completed_at`) are never selected. `runs` rows stay as a compact durable
index; locks are untouched.

### GC Command

`roundfix gc [--dry-run]` (support command, non-interactive): resolve the
retention window; on `--dry-run` run the read-only eligibility query and print
the report as a preview with nothing changed; otherwise prune the journal,
delete each pruned Run's `<artifact_dir>/runs/<run-id>` directory, delete
orphaned `runs/<id>` directories (no matching `runs` row), and print the freed
Runs, rows, and bytes. Writes requested output to stdout, diagnostics to stderr,
stable exit codes. A retention of `0` prunes nothing and says so.

### Self-healing prune

The existing preflight sweep (already reaping worktree debris and closing
sessions) gains a best-effort `PruneTerminalRuns` call plus pruned-Run artifact
cleanup when `journal_retention > 0`. Failures are one stderr warning, never
fatal, never block the Run. One optional summary line reports what the sweep
pruned.

## Coverage Map

- Stories 1, 4 → `PruneTerminalRuns` scope + config (ADR-0033)
- Stories 2, 3 → GC Command with `--dry-run`
- Story 5 → preflight-sweep best-effort prune
- Feature 3 → terminal-only, Active-Run-safe query

## Integration Points

SQLite via the existing store connection (row deletes) and the filesystem
(artifact directory removal). No new external systems. Reuses the preflight
sweep seam already present for worktree reaping and session close.

## Testing Approach

- Store: table tests over a seeded database with a mix of Active and terminal
  Runs at varied `completed_at`, asserting only terminal Runs past the cutoff
  lose their `run_events`, that Active Runs and all `runs` rows survive, and the
  returned counts.
- GC Command: buffer-captured CLI tests — dry-run changes nothing and lists the
  eligible set; live run prunes journals, removes the right artifact dirs,
  removes an orphaned `runs/<id>` dir, leaves an Active Run's dir, and reports
  freed counts; `journal_retention: 0` prunes nothing.
- Preflight sweep: a test asserts an eligible terminal Run is pruned during a
  new Run's preflight and that a sweep prune failure does not fail the Run.

## Build Order

1. Config `store.journal_retention` + `store.PruneTerminalRuns` with the
   terminal-only, Active-safe query and tests (no deps)
2. GC Command (`--dry-run` + live prune + artifact/orphan cleanup + report)
   (depends on: 1)
3. Best-effort retention prune in the preflight sweep (depends on: 1)
4. Docs and skill sync (depends on: 1, 2, 3)

## Risks & Considerations

- The prune query must be conservative: a wrong predicate could delete an Active
  Run's journal mid-flight — gate strictly on a terminal state AND a non-empty
  `completed_at` older than the cutoff, and cover it with tests seeding an
  Active Run at an old `created_at`.
- Artifact cleanup must only remove `runs/<id>` directories, never review
  artifacts under the spec tree (spec 0011) or unrelated paths — resolve paths
  under the resolved run artifact root only.
- `runs` rows are intentionally kept (compact index); if their growth ever
  matters, a separate policy — not this Spec — decides row deletion.
- Row deletion does not shrink the SQLite file immediately; note VACUUM as a
  deferred follow-up rather than doing it inline (it locks the database).

## Decisions

- Retention prunes Run Event Journal rows + artifact directories of terminal
  Runs past a window; `runs` rows and locks are never pruned. See ADR-0033.
- Default `journal_retention` is `336h` (14 days); `0` keeps everything.
- Cleanup is exposed as the GC Command (with `--dry-run`) and runs best-effort
  in the preflight sweep. Glossary gains GC Command and Journal Retention.
