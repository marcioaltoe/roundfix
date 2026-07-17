---
status: pending
created_at: 2026-07-17
updated_at: 2026-07-17
---

# Run storage — sanitation is repository-scoped and does not compact SQLite (2026-07-17)

An inspection of Roundfix Home followed the cleanup of terminal Run Worktrees and evaluated
whether Roundfix can reclaim its remaining Run storage without direct filesystem or SQLite
operations. The existing GC Command provides automated retention cleanup, but no active Spec
covers global Artifact Directory sanitation or physical Run Database compaction.

Environment:

- Roundfix Home: `/Users/marcio/.roundfix`.
- Run Database: `/Users/marcio/.roundfix/roundfix.db`.
- Default Artifact Root: `/Users/marcio/.roundfix/artifacts`.
- Journal Retention: `336h` (14 days).

## 1. The GC Command does not physically reclaim deleted SQLite pages

- **Symptom / evidence**: on 2026-07-17, `du -sh` reported a 1.4 GB Run Database. SQLite
  reported 1,208,625 `run_events` rows and attributed 1,469,747,200 bytes to the
  `run_events` table plus 67,338,240 bytes to its primary-key index. A fresh
  `roundfix gc --dry-run` reported 33 eligible Runs but zero eligible Journal rows because the
  expired rows had already been pruned. The archived Run store retention Tech Spec explicitly
  deferred `VACUUM`, and no current implementation or active Spec adds a compaction operation.
- **Root cause**: `PruneTerminalRuns` deletes eligible Run Event Journal rows but SQLite row
  deletion does not return database pages to the filesystem. The GC Command has no guarded
  database compaction phase or separate compaction mode.
- **Action / suggestion**: route a future Spec for safe Run Database compaction. It should provide
  a deterministic preview, refuse while an Active Run or writer can be present, verify temporary
  disk capacity, compact through a failure-safe SQLite operation, and report database bytes before,
  reclaimable, and after. The Spec must decide whether compaction is explicit, threshold-based, or
  an optional GC phase; it must not make every operational retention sweep take an exclusive
  database lock.

## 2. Artifact cleanup resolves only the current repository's Artifact Root

- **Symptom / evidence**: `du -sh /Users/marcio/.roundfix/artifacts` reported 250 MB across
  multiple repository-scoped directories. A fresh `roundfix gc --dry-run` from the Roundfix
  repository reported zero reclaimable Artifact bytes and only two empty orphan Artifact
  Directories, even though the Run Database is machine-wide and contains Runs from multiple
  repositories. `runGC` resolves one Artifact Root from the current Project Config and scans only
  that root's `runs/` directory.
- **Root cause**: Run Event Journal pruning is machine-wide because it queries the global Run
  Database, while Artifact Directory discovery is scoped to the repository from which the GC
  Command runs. Roundfix does not retain or derive a complete set of Artifact Roots for one global
  cleanup pass.
- **Action / suggestion**: add a global, dry-run-first sanitation mode that discovers every
  Roundfix-owned Artifact Root using durable Run metadata or another validated registry. It must
  classify roots that are missing, overridden, outside Roundfix Home, or unsafe; remove only
  directories proven to belong to eligible or absent Runs; and leave Review Artifacts and
  ambiguous paths unchanged. Per-repository GC should remain available for the narrow case.

## 3. Storage lifecycle policy does not cover every durable Run record

- **Symptom / evidence**: the Run Database currently contains 279 `runs` rows, zero
  `active_run_locks`, and zero `run_agent_selections` rows. The implemented retention primitive
  promises never to delete `runs` rows or Active Run locks and has no pruning contract for Agent
  Selection records. The current 279 Run summary rows occupy only 118,784 bytes, so immediate
  deletion is unnecessary, but their long-term lifecycle is undefined.
- **Root cause**: Spec 0014 intentionally scoped retention to Run Event Journal rows and Run
  Artifact Directories. Tables added after that contract do not automatically inherit a retention
  policy, and the durable Run index has no explicit long-term bound.
- **Action / suggestion**: the future sanitation Spec should define table-by-table ownership and
  retention. Preserve the compact Run index by default unless measured growth justifies pruning;
  keep Active Run locks governed by Run lifecycle; and prune Agent Selection records only with
  their owning Run or according to an explicit evidence-retention rule. Include a read-only storage
  report grouped by repository, state, table, and Artifact Root so policy decisions use measured
  data.

## What worked — keep

- `roundfix gc --dry-run` is deterministic, non-interactive, and reports the exact eligible set.
- Journal Retention safely excludes Active Runs and retains compact Run history.
- Artifact cleanup constrains deletion to validated `runs/<run-id>` directories.
- Operational commands perform best-effort retention pruning, so users do not need to invoke GC
  for normal Journal Retention enforcement.

## Routing

No active Spec owns the remaining database-compaction and global Artifact sanitation work as of
2026-07-17. The terminal Run Worktree cleanup gap is separate and remains owned by
[Spec 0038 — Terminal Run Worktree reconciliation](../specs/0038-terminal-run-worktree-reconciliation/_prd.md).
This finding should remain `pending` until the storage sanitation work is routed to a dedicated
Spec.
