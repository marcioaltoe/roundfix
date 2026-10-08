---
spec: 0250-a-run-database-that-keeps-only-recent-runs
prd: _prd.md
created: 2026-10-08
---

# A Run Database that keeps only recent Runs — Technical Spec

## Executive Summary

Run Retention adds an outer window to the Run Database. A terminal Run that
completed more than `store.run_retention_days` days ago (7, 15 or 30, default
30) is removed in one write transaction per Run, and SQLite's existing
`ON DELETE CASCADE` foreign keys remove its journal, Agent Selections, token
usage and lock. A Run Retention Sweep applies the window and then returns
freed pages with `PRAGMA incremental_vacuum`. It runs at the start of a Run
or a Delivery Queue at most once a day under a two-second budget, and
`roundfix gc` runs it without a budget. The trade-off accepted: an existing
database shrinks only after one guarded full compaction converts it to
incremental mode, which took 4.9 to 5.9 s on a copy of this machine's 1.1 GB
database. That cost is paid in an explicit command, never at a Run start
(ADR-0255).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. Run IDs
  keep their shape, and the hint of API Contract 6 only parses the creation
  time they already embed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added. Every test uses `t.TempDir()` homes, fixture Run Databases
  and fake runtimes. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0255 governs every decision here.
  ADR-0033 keeps Journal Retention as the inner window, and its rule that
  retention never deletes a Run row now holds only inside the outer one:
  ADR-0033: "are never touched by retention". ADR-0171 keeps a prune
  reporting only what it reclaimed, and the Run Retention report follows it:
  ADR-0171: "A second prune after a complete one reclaims nothing and reports zero". ADR-0172 keeps the Doctor storage check unchanged and `runs list` without a
  notice: ADR-0172: "It is the Agents' hot path". ADR-0053 keeps reconcile
  proof-based, so a Run whose Run Worktree exists is kept:
  ADR-0053: "Ambiguous, dirty, and unintegrated work stays intact". ADR-0004
  makes the central run database machine-wide, so the window is User Config:
  ADR-0004: "Roundfix stores Run state in a global SQLite database". ADR-0225
  makes a queue item start `roundfix implement`, so the item's own start runs
  the sweep's due check: ADR-0225: "the `roundfix implement` it started in the
  item worktree". ADR-0184 asks for transcripts:
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_03 changes the Roundfix skill, a
  protected path. Express maintainer authorization: "considere autorizado a
  ajustar todas as skills se necessário" (2026-09-30) and the grant of the
  Governed Paths this Spec declares (2026-10-08). No Makefile, lint,
  formatter, test-runner, workflow or `go.mod` change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0250-a-run-database-that-keeps-only-recent-runs/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/runs.md`,
  `.agents/skills/roundfix/references/storage.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured read-only on 2026-10-08. The live `~/.roundfix/roundfix.db` was
1,180,217,344 bytes, 288,139 pages of 4,096 bytes, `auto_vacuum = 0`, WAL,
and 0 free pages. It held 313 Runs (149 Clean, 158 Unresolved, 6
BudgetExceeded, 0 Active), 832,085 events, 1,571 Agent Selections, 323 token
usage rows, 0 locks, 0 Run Windows, 1 Delivery Queue with 1 item, 1 Delivery
Queue Run link, 3 intents and 3 receipts. Every Run records an empty
`artifact_dir` and a `work_dir`; `~/.roundfix/worktrees` was empty and
`~/.roundfix/artifacts` held 13 MB. The User Config sets no `store` key, so
Journal Retention is the `336h` default, and the oldest event is from
2026-09-24.

On a `.backup` copy, through `modernc.org/sqlite` v1.57.0 (the repository's
driver), with the windows measured from 2026-10-08T12:00:00Z:

| window | Runs removed | events | Agent Selections | event bytes (payload and summary) |
| --- | --- | --- | --- | --- |
| 30 days | 12 | 0 | 80 | 0 |
| 15 days | 68 | 0 | 330 | 0 |
| 7 days | 200 | 439,456 | 1,016 | 675,112,811 |

- 7-day removal: one transaction for all Runs 2.75 s; one transaction per Run
  4.29 s, slowest Run 96 ms. A full `VACUUM` afterwards: 2.33 s, file
  391,499,776 bytes.
- Full `VACUUM` of the whole copy: 6.9 s. `PRAGMA auto_vacuum = INCREMENTAL`
  then `VACUUM`: 4.9 s and 5.9 s on two runs, file 1,176,240,128 bytes.
- On the converted copy, 7-day removal per Run: 2.09 s, slowest Run 713 ms;
  incremental vacuum in 2,560-page slices: 75 slices, 14.3 s, slowest slice
  453 ms, file 392,019,968 bytes.
- One heavy day (cutoff 2026-09-25T12:00:00Z, 102 Runs, 107,746 events): on
  the converted copy, removal 1.10 s (slowest Run 719 ms) and incremental
  vacuum 0.78 s in 22 slices (slowest 50 ms), file 950,824,960 bytes; on the
  unconverted copy, removal 1.33 s and a full `VACUUM` 4.97 s.
- With the `sqlite3` 3.54.0 shell: on a new file, `auto_vacuum = INCREMENTAL`
  set after `journal_mode = WAL` is silently ignored (the database reports 0),
  and set before it holds (2). A read-only connection may set the pragma and
  `VACUUM INTO` writes the snapshot in incremental mode while the source
  stays at 0.

The payload estimate of the 7-day window was 675 MB against 789 MB measured
after removal and `VACUUM`, so the estimate understates by about 15 percent.

## System Architecture

No new package. The work extends three existing seams:

- `internal/store` gains the removal, the sweep record and the incremental
  compaction, next to `PruneTerminalRuns` and `Compact`.
- `internal/config` gains `store.run_retention_days` beside
  `store.journal_retention`.
- `internal/cli` gains the Run Retention Sweep, called from the GC Command,
  from the start paths that already run the Journal Retention prune, and from
  `deliver start`.

```text
implement / resolve / watch start      deliver start (after the queue is recorded)
  sweepRunRetention (Journal Retention, unchanged)
  runRetentionAtStart  -> due? (run_retention_sweeps) -> sweep(budget 2s) -> stderr line
roundfix gc [--dry-run]
  Journal Retention section (unchanged except the zero-window wording)
  sweep(no budget) -> Run Retention section -> compaction (incremental, or guarded full)
sweep:
  store.RunRetentionCandidates(cutoff) -> keep: queue / worktree / unproven root
  per Run: remove runs/<id> dir -> store.RemoveRetainedRun (one write tx, cascade)
  store.CompactIncrementally(2560) per slice -> store.RecordRunRetentionSweep when complete
```

## Implementation Design

### Interfaces

```go
// internal/store
type RunRetentionRows struct{ Runs, RunEvents, AgentSelections, TokenUsage, ActiveRunLocks int64 }
type RunRetentionCandidate struct {
	RunID, Repository, GitRoot, ArtifactDir, WorkDir string
	CompletedAt     time.Time
	QueueReferenced bool
	Rows            RunRetentionRows
	EstimatedBytes  int64
}
type RunRetentionKeptError struct{ RunID, Reason string }
type RunRetentionSweep struct {
	CompletedAt   time.Time
	RetentionDays int
}
type IncrementalCompaction struct {
	Incremental              bool
	PagesReleased, FreePages int64
}
const IncrementalCompactionSlicePages = 2560
func (store *Store) RunRetentionCandidates(ctx context.Context, cutoff time.Time) ([]RunRetentionCandidate, error)
func (store *Store) RemoveRetainedRun(ctx context.Context, runID string, cutoff time.Time) (RunRetentionRows, error)
func (store *Store) LastRunRetentionSweep(ctx context.Context) (RunRetentionSweep, bool, error)
func (store *Store) RecordRunRetentionSweep(ctx context.Context, sweep RunRetentionSweep) error
func (store *Store) CompactIncrementally(ctx context.Context, maxPages int64) (IncrementalCompaction, error)
```

```go
// internal/config
type Store struct {
	JournalRetention  time.Duration
	RunRetentionDays  int // 7, 15 or 30; Builtin 30
}
var RunRetentionDayChoices = []int{7, 15, 30}

// internal/cli
type runRetentionOptions struct {
	dryRun bool
	budget time.Duration // 0 means no budget
}
func sweepRunDatabase(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, opts runRetentionOptions) (runRetentionReport, error)
func runRetentionAtStart(ctx context.Context, runStore *store.Store, loaded roundconfig.Loaded, stderr io.Writer)
```

### Data Models

Schema version 23 adds one table:

```text
run_retention_sweeps (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  completed_at TEXT NOT NULL,      -- RFC 3339, UTC
  retention_days INTEGER NOT NULL  -- 7, 15 or 30
)
```

The v22 to v23 migration creates it; a fresh database creates it with the
rest of the schema. A database this binary creates reports
`PRAGMA auto_vacuum = 2`, which on SQLite means the pragma is issued before
`journal_mode = WAL` on the new file. No other table changes. The Run Database
lifecycle policy gains a row for the new table, and the rows for `runs`,
`run_events`, `run_agent_selections`, `run_token_usage` and
`active_run_locks` name Run Retention as their age bound.

### API Contracts

1. API Contract: `store.run_retention_days`. Integer, one of 7, 15 or 30,
   Builtin 30. In User Config any other value, a non-integer included, fails
   config loading with `parse config "<path>": store.run_retention_days must be 7, 15 or 30`.
   In Project Config the key is removed before decoding and prints
   `config: store.run_retention_days in Project Config is ignored; set store.run_retention_days in User Config`
   on stderr once, as `runs.max_active` does. The User Config template that
   `roundfix init` writes carries the key with its comment under `store:`;
   the Project Config template does not. `store.journal_retention` is
   unchanged in meaning, scopes, default and validation.
2. API Contract: removal. `RunRetentionCandidates` reads, on any connection,
   every terminal Run whose `completed_at` is before the cutoff, with its row
   counts per table, `QueueReferenced` set when a `delivery_queue_items.run_id`
   or `delivery_queue_runs.run_id` names it, and `EstimatedBytes`, the sum of
   the byte lengths of its events' `payload` and `summary`. The sweep keeps a
   candidate that is queue-referenced, whose `WorkDir` is non-empty and exists
   (`lstat`), or whose `runs/<run-id>` directory exists under a root that is
   not proven. A root is proven when it is the recorded `artifact_dir`, or,
   when that is empty, a default derived from the Run's repository key or
   recorded checkout, and it is a clean absolute physical directory inside
   Roundfix Home, the proof `roundfix gc sanitize` uses. For a removable Run
   the sweep removes the directory first and then calls `RemoveRetainedRun`,
   which in one write transaction re-reads the Run, refuses with
   `RunRetentionKeptError` when it is no longer terminal, no longer before the
   cutoff or now queue-referenced, counts its rows, deletes the `runs` row and
   returns the counts the cascade removed. A Run absent at that point counts
   nothing and is not an error.
3. API Contract: the sweep. It runs Run Retention, then compaction, then
   records completion. In incremental mode it calls `CompactIncrementally`
   with 2,560 pages, each call one write transaction, until no free page is
   left, and finishes with a passive WAL checkpoint. In the default mode the
   budgeted sweep does not compact; the unbudgeted live `roundfix gc` runs the
   existing guarded `PreviewCompaction` and `Compact` when the database has
   free pages or is not in incremental mode, and reports a refusal as
   `skipped (<reason>)` without failing. With a budget, the sweep checks the
   clock before each Run removal and each slice and stops when the budget is
   spent; it then records nothing. It records `RunRetentionSweep` only when
   every candidate was removed or kept and no free page is left in
   incremental mode. The clock is the GC Command's existing `now` dependency.
4. API Contract: once a day. `runRetentionAtStart` reads
   `LastRunRetentionSweep`. The sweep is due when there is none, when it
   completed 24 hours or more before now, or when its `RetentionDays` differs
   from the configured window. When due it runs the sweep with a two-second
   budget. It prints on stderr, only when it removed a Run,
   `roundfix: Run Retention removed runs=<n> rows=<n> database_bytes_reclaimed=<n>`,
   where rows is the sum of every table's removed rows and the bytes are page
   count times page size before minus after; when the budget ran out,
   `roundfix: Run Retention paused after its 2s budget; it continues at the next Run start`;
   on any error, `roundfix: warning: Run Retention failed: <reason>`. It never
   changes the command's exit code. It runs in `implement`, `resolve` and
   `watch` right after their Journal Retention prune, and in `deliver start`
   after the Delivery Queue is recorded and before its store closes.
5. API Contract: `roundfix gc [--dry-run]`. The Journal Retention lines stay
   as they are; with `journal_retention: 0` the header is the normal
   `GC dry-run` or `GC complete` header and the Journal Retention block prints
   `Journal Retention: 0` and `No journal pruning performed.` in place of
   `GC skipped` and `No pruning performed.`. Then the Run Retention section of
   Surface Transcript 1 or 2 follows. The dry run opens the reader and writes
   nothing; its counts are those the live run would remove at the same clock.
   The live run removes, compacts and records completion. Its compaction line
   is `Compaction: incremental (<n> pages)`, `Compaction: full`,
   `Compaction: full (converted to incremental)`, `Compaction: not needed` or
   `Compaction: skipped (<reason>)`. Database bytes are page count times page
   size. `roundfix gc compact` keeps its surface; its preview and `--apply`
   set `auto_vacuum = INCREMENTAL` before `VACUUM INTO` and `VACUUM`, so the
   preview still matches the result within one page.
6. API Contract: unknown Run hint. `roundfix runs show <id>` and
   `roundfix events <id>` keep their refusal and exit 2. When the ID has the
   shape `run_<YYYYMMDD>T<HHMMSS>Z_<hex>` and that time is before the Run
   Retention cutoff, the message `Run "<id>" does not exist` continues with
   `; Run Retention may have removed it, because it removes terminal Runs that completed more than <N> days ago`.
   Any other ID keeps the existing message byte for byte.

### Invariants

```text
1. A Run that is not terminal is never removed.
2. A Run named by any Delivery Queue item or Delivery Queue Run link is never removed.
3. A Run whose recorded Run Worktree exists is never removed.
4. A Run whose artifact directory exists under an unproven root is never removed, and no directory outside a proven root is touched.
5. A removed Run leaves no row in runs, run_events, run_agent_selections, run_token_usage or active_run_locks.
6. run_windows, interactive_defaults and the Delivery Queue tables are byte-identical across a sweep.
7. The dry run changes no byte of the Run Database or of any artifact directory.
8. A budgeted sweep takes the write lock only per Run removal and per compaction slice, never a full VACUUM.
9. A sweep that stopped on its budget records no completion; the next due start resumes it.
10. Journal Retention prunes exactly what it pruned before this Spec.
```

### Surface Transcripts

1. Surface Transcript: the dry run on a fixture User Config without the key,
   with one removable Run, one queue-referenced Run and one Run with an
   existing Run Worktree past the cutoff.

   ```transcript
   $ roundfix gc --dry-run
   stdout:
   GC dry-run
   ...
     Run Retention: 30 days
     Run Retention cutoff: <time>
     Runs removable: 1
     Runs kept past the cutoff: 2 (queue-referenced 1, worktree present 1, artifact root unproven 0)
     Rows removable: runs=1 run_events=<n> run_agent_selections=<n> run_token_usage=<n> active_run_locks=0
     Database bytes reclaimable (estimated): <n>
     Removable Runs:
       <run-id>
   stderr:
   exit: 0
   ```

2. Surface Transcript: the live run on the same fixture.

   ```transcript
   $ roundfix gc
   stdout:
   GC complete
   ...
     Run Retention: 30 days
     Run Retention cutoff: <time>
     Runs removed: 1
     Runs kept past the cutoff: 2 (queue-referenced 1, worktree present 1, artifact root unproven 0)
     Rows removed: runs=1 run_events=<n> run_agent_selections=<n> run_token_usage=<n> active_run_locks=0
     Run artifact bytes reclaimed: <n>
     Compaction: <mode>
     Database bytes before: <n>
     Database bytes after: <n>
     Removed Runs:
       <run-id>
   stderr:
   exit: 0
   ```

3. Surface Transcript: a User Config window outside the three values.

   ```transcript
   $ roundfix gc --dry-run
   stdout:
   stderr:
   Preflight failed
   ...
     parse config "<path>": store.run_retention_days must be 7, 15 or 30
   ...
   exit: 2
   ```

4. Surface Transcript: an unknown Run whose ID predates the cutoff.

   ```transcript
   $ roundfix runs show run_20260801T120000Z_0123456789abcdef
   stdout:
   stderr:
   roundfix: runs show failed: Run "run_20260801T120000Z_0123456789abcdef" does not exist; Run Retention may have removed it, because it removes terminal Runs that completed more than 30 days ago
   Run 'roundfix runs show --help' for usage.
   exit: 2
   ```

5. Surface Transcript: the same lookup through the event stream.

   ```transcript
   $ roundfix events run_20260801T120000Z_0123456789abcdef
   stdout:
   stderr:
   roundfix events failed: Run "run_20260801T120000Z_0123456789abcdef" does not exist; Run Retention may have removed it, because it removes terminal Runs that completed more than 30 days ago
   exit: 2
   ```

## Coverage Map

- Goal "removed whole" → API Contract 2, Invariant 5, task_01, task_02.
- Goal "never removed" → API Contract 2, Invariants 1 to 4, task_01, task_02.
- Goal "on their own, bounded" → API Contracts 3 and 4, Invariants 8 and 9,
  task_03.
- Goal "dry run and gc" → API Contract 5, Surface Transcripts 1 and 2,
  Invariant 7, task_02.
- User Story 1 → API Contract 4 (task_03). User Story 2 → API Contract 1
  (task_02). User Story 3 → Invariant 2, task_03's `deliver status` test.
  User Story 4 → API Contract 5. User Story 5 → API Contract 6 (task_03).
- Core Feature 1 → API Contracts 1 and 2. Core Feature 2 → API Contract 2,
  Invariants 1 to 4. Core Feature 3 → API Contracts 3 and 4. Core Feature 4 →
  Data Models, API Contracts 3 and 5. Core Feature 5 → API Contracts 5 and 6,
  Surface Transcripts. Core Feature 6 → task_02 (GC Command and configuration
  guides), task_03 (Roundfix skill, `runs` and `events` guides), task_04
  (glossary and lifecycle policy).
- Success Metric 1 → task_01 and task_02 fixture tests.
- Success Metric 2 → Surface Transcript 1.
- Success Metric 3 → task_03 budget and due tests.
- Success Metric 4 → task_01 compaction tests, Surface Transcript 2.
- Success Metric 5 → QA measurement on a generated fixture.
- Success Metric 6 → task_03 `deliver status` test.

## Integration Points

- SQLite through `modernc.org/sqlite`: foreign-key cascades (already enabled
  on every connection), `PRAGMA auto_vacuum`, `PRAGMA incremental_vacuum(N)`,
  `PRAGMA freelist_count`, `PRAGMA wal_checkpoint(PASSIVE)`.
- The machine-wide advisory write lock of `withWriteTx`, shared with every
  running Daemon.
- The filesystem: `lstat` of Run Worktrees, removal of `runs/<run-id>`
  directories under proven roots.

## Testing Approach

- `internal/store/run_retention_test.go` (task_01): fixture databases in
  `t.TempDir()` homes with Runs, events, selections, token usage, locks and a
  Delivery Queue. It proves the candidate scan, the re-check inside the
  transaction, the cascade counts, the untouched tables, the sweep record, the
  incremental mode of a fresh database and the shrink after
  `CompactIncrementally`, and that `Compact` converts a default-mode database
  within the preview tolerance. The migration test proves v22 to v23. The
  lifecycle policy test already compares the policy table with every table.
- `internal/config/run_retention_test.go` (task_02): defaults, the three
  values, refusals, the Project Config warning and the User Config template.
- `internal/cli/gc_run_retention_test.go` (task_02): `runCLIContext` with the
  existing `withGCNow` clock and `seedGCFixture`-style seeding, for the dry
  run, the live run, the kept reasons, artifact proofs and compaction lines.
- `internal/cli/run_retention_start_test.go` (task_03): the due rule, the
  budget with a stepping clock through the existing GC Command `now`
  dependency, the stderr lines, `deliver start`, `deliver status` before and
  after, the reconcile-kept Run, and the unknown Run hint.
- No new test package; `internal/store`, `internal/config` and `internal/cli`
  already install `suiteguard.Main`.

## Build Order

1. Store: schema 23, incremental mode at creation, candidates, removal,
   sweep record, incremental compaction, conversion in `Compact`, and the
   lifecycle policy table (`internal/store`).
2. Config key, the GC Command's Run Retention section, the sweep engine and
   compaction, and the GC Command and configuration guides (depends on: 1).
3. Automatic sweep at the start of Runs and Delivery Queues, the unknown Run
   hint, the `runs` and `events` guides and the Roundfix skill (depends on: 2).
4. Glossary and the lifecycle policy prose (depends on: 1, because it edits
   the lifecycle policy prose after task_01 edits its table).
5. QA gate (depends on: 2, 3, 4).

## Risks & Considerations

- At the 30-day default this machine reclaims about 12 Run rows and almost no
  bytes, because Journal Retention already prunes at 14 days; the bytes fall
  only with a 7-day window or a shorter Journal Retention. The ADR records it.
- One very large Run removes in one transaction; the slowest measured was
  0.72 s, so the start may overshoot the budget by about one step.
- Older binaries refuse the schema-23 database. A queue item built from an
  older branch fails the same way it does after any migration.
- A Run Branch of a removed Run is no longer listed by reconcile; it holds no
  bytes and Git keeps it.
- Tests that seed terminal Runs completed more than 30 days before the real
  clock and then start a Run would now see those Runs removed. The
  `internal/cli` suite seeds old Runs only through the GC Command's fake clock
  or within 400 hours, so no existing start test is expected to change; each
  Task runs the full package to prove it. `internal/cli/cli_test.go` is a
  Governed Path and no Task changes it.
- The byte estimate of the dry run understates by about 15 percent; it is
  labelled an estimate and the live run reports measured bytes.

## Glossary

- adds: **Run Retention**
- adds: **Run Retention Sweep**
- changes: **GC Command**
- changes: **Journal Retention**

## Decisions

- One write transaction per Run, removal through the existing cascades; see
  ADR-0255.
- The daily record is a one-row table in the Run Database; see ADR-0255.
- A two-second budget in place of a background process; see ADR-0255.
- Incremental compaction in the sweep, conversion only by a guarded full
  compaction; see ADR-0255.
- The dry-run byte figure is an estimate from stored event bytes, and the
  live figure is measured from page counts; see ADR-0255.
