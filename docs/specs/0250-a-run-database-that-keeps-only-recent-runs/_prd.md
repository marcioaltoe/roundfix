---
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: active
created: 2026-10-08
surfaces: [backend, cli, data, docs]
---

# A Run Database that keeps only recent Runs

The machine-wide Run Database grows with every Run and never sheds a Run. On
2026-10-08 it was 1.1 GB, of which `run_events` held 1.08 GB for 832k events
written in the two weeks before, about 77 MB a day, across 313 Runs. Journal
Retention prunes the journal of an old terminal Run but keeps its row and the
rows that hang from it. ADR-0033: "Active Runs and their journals are never
pruned". Freed pages go back to the filesystem only when someone runs
`roundfix gc compact --apply` by hand, and `gc compact` reclaimed 5 MB because
the space was live data. The maintainer decided on 2026-10-08, "Automática +
gc (Recommended)", that Roundfix keeps only the data of Runs still running,
and of any Run a Delivery Queue still references, plus the last N days. The
operator then gets a database that stays bounded without a manual step, and
`roundfix gc` shows what it would remove before it removes it.

## Prerequisites

Spec 0249 delivers before this Spec, and Spec 0251 is authored at the same
time. This Spec's Verification reads no artifact of either. When either Spec
also raises the Roundfix skill's version, the version this Spec records is the
next version at record time, chosen by the record command.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. Run IDs
  keep their existing shape; the new configuration key and the new table name
  are lowercase words in the existing style. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added. The work reads and writes the local Run Database and local
  artifact directories, and every test uses temporary homes and fixture
  databases. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0255 (this Spec) decides the
  window, the removal, the protections, the daily record, the budget and the
  compaction. It amends ADR-0033, which keeps the row index:
  ADR-0033: "are never touched by retention". The row stays load-bearing only
  while the Run is inside the window. Under ADR-0171 a retention prune reports only what it reclaimed:
  ADR-0171: "A second prune after a complete one reclaims nothing and reports zero". ADR-0172 keeps the Doctor storage check and the quiet `runs list`:
  ADR-0172: "It is the Agents' hot path". ADR-0053 keeps reconcile proof-based,
  so a Run whose Run Worktree still exists is kept:
  ADR-0053: "Ambiguous, dirty, and unintegrated work stays intact". ADR-0004
  makes the central run database machine-wide, so the window lives in User
  Config:
  ADR-0004: "Roundfix stores Run state in a global SQLite database". ADR-0184
  asks for transcripts of the changed command surfaces:
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts". ADR-0225
  has a queue item run `roundfix implement` itself, so an item's start checks
  whether the sweep is due:
  ADR-0225: "A queue item runs the Roundfix binary its branch builds".
  ADR-0022 applies only to Active Runs, which Run Retention never removes.
  ADR-0161, ADR-0212 and ADR-0232 release the Runs of a merged Spec and remove
  their Run Worktrees; Run Retention removes such a Run only after that
  release, once its Run Worktree is gone. ADR-0170 reads the Runs that Task
  Carry-Forward inspects, and those keep their Run Worktree, so Run Retention
  keeps them. ADR-0237 records on a Delivery Queue item a merge made outside
  the queue, and Run Retention keeps every Run a queue item names. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix skill is a protected path. The
  maintainer granted the skills on 2026-09-30, "considere autorizado a ajustar
  todas as skills se necessário", and on 2026-10-08 granted the Governed Paths
  this Spec declares. No Makefile, lint, formatter, test-runner, workflow or
  `go.mod` change is proposed. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0250-a-run-database-that-keeps-only-recent-runs/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/runs.md`,
  `.agents/skills/roundfix/references/storage.md`, `skills/roundfix/SKILL.md`.

## Goals

- A terminal Run older than the configured window leaves the Run Database
  whole, with every row that depends on it and its artifact directory.
- A Run that is still running, that a Delivery Queue references, or whose Run
  Worktree still exists is never removed.
- The removal and the compaction run on their own, at most once a day, and
  never hold a Run start or a Delivery Queue start for more than a couple of
  seconds.
- `roundfix gc --dry-run` shows the Runs, rows and bytes the window would
  remove and the Runs it keeps, and `roundfix gc` removes and compacts on
  demand.

## User Stories

1. As an operator, I want old terminal Runs to leave the Run Database without
   a manual step, so that the database stays bounded.
2. As an operator, I want to choose a 30, 15 or 7 day window in User Config,
   so that I trade history for disk on this machine.
3. As an operator running a Delivery Queue, I want the Runs the queue still
   reads to stay, so that `deliver status` and its token totals do not change
   under me.
4. As an operator, I want `roundfix gc --dry-run` to tell me which Runs, how
   many rows and roughly how many bytes would go, so that I can apply a window
   with confidence.
5. As an Agent or operator looking up an old Run, I want the refusal to say
   that Run Retention may have removed it, so that I do not hunt for a bug.

## Core Features

1. **Run Retention.** A User Config window of 7, 15 or 30 days, 30 by
   default. A terminal Run that completed before the window's cutoff is
   removed whole: its row, its journal, its Agent Selection records, its
   token usage and any lock, and its artifact directory under a proven
   Artifact Root. A Project Config value is ignored with a warning, and any
   other value is a configuration error (ADR-0255).
2. **Protected Runs.** Run Retention keeps a Run that is not terminal, that a
   Delivery Queue item or Delivery Queue Run link names, whose recorded Run
   Worktree still exists, or whose artifact directory sits under a root it
   cannot prove. Each removal re-checks the Run before deleting it (ADR-0255).
3. **Run Retention Sweep.** The step that applies Run Retention and then
   compacts the database. It runs automatically at the start of every Run and
   of every Delivery Queue when it is due, at most once a day, under a
   two-second budget, and it continues at the next start when the budget runs
   out. A failure prints one warning and never blocks the start. The live
   `roundfix gc` runs it without a budget (ADR-0255).
4. **Incremental compaction.** A new Run Database returns freed pages to the
   filesystem in small slices during the sweep. An existing database is
   converted by the next guarded full compaction, from
   `roundfix gc compact --apply` or the live `roundfix gc` (ADR-0255).
5. **Report surface.** `roundfix gc --dry-run` reports the window, the cutoff,
   the Runs it would remove and the Runs it keeps by reason, the rows per
   table and an estimate of the bytes. The live `roundfix gc` reports what it
   removed, the measured file bytes before and after, and how it compacted.
   `roundfix runs show` and `roundfix events` name Run Retention when an
   unknown Run ID predates the cutoff (ADR-0255).
6. **Docs, skill and glossary.** The configuration guide, the Run Database
   lifecycle policy, the GC Command guide and the Roundfix skill describe Run
   Retention, and `CONTEXT.md` gains **Run Retention** and **Run Retention
   Sweep** and revises **GC Command** and **Journal Retention**.

## User Experience

The automatic sweep is silent when it removes nothing. When it removes Runs it
prints one stderr line naming the Runs, rows and bytes it reclaimed; when its
budget runs out it prints one line saying it continues at the next start; when
it fails it prints one warning. `roundfix gc --dry-run` prints a Run Retention
section after the Journal Retention section, with the Runs it would remove
listed by ID. The live `roundfix gc` prints the same section in the past tense
with the compaction result. A refusal for an unknown Run keeps its existing
first words and adds the Run Retention hint only for an ID that predates the
cutoff.

## Non-Goals / Out of Scope

- Applying a window to `~/.roundfix` on this machine. The operator applies the
  30-day window after this Spec ships; no Task or QA row writes there.
- Changing Journal Retention's meaning, scopes or default, or the Doctor
  storage check.
- A value that turns Run Retention off, or a window other than 7, 15 or 30.
- Tombstones for removed Runs, or any record of a removed Run ID.
- Releasing the Run Branches of removed Runs; reconcile keeps its scope.
- A full `VACUUM` during the automatic sweep, or a background process for it.
- Pruning `run_windows`, `interactive_defaults` or the Delivery Queue tables.

## Success Metrics

1. Success Metric: on a fixture Run Database, a terminal Run 31 days past
   completion is removed with every row of its journal, Agent Selections,
   token usage and lock, and its artifact directory, while an Active Run, a
   queue-referenced Run, a Run with an existing Run Worktree and a Run
   completed one day ago all stay with every row.
2. Success Metric: `roundfix gc --dry-run` on that fixture lists exactly the
   removable Run, the rows per table and a non-zero byte estimate, and changes
   no byte of the database.
3. Success Metric: on a fixture whose Runs need more than the budget, the
   automatic sweep at a Run start returns within the two-second budget plus one
   step, records no completion, and a later start finishes the work and records
   it; a second start within 24 hours does nothing.
4. Success Metric: on a fresh fixture database, the sweep after a removal
   shrinks the file, and on a database in the default mode the live
   `roundfix gc` converts it and shrinks it.
5. Success Metric: measured on a generated fixture of at least 100k events,
   one sweep's removal and compaction times are recorded beside the
   TechSpec's measurements of 2026-10-08.
6. Success Metric: `deliver status` prints the same lines and token totals
   before and after a sweep removes every other old Run.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- SQLite's PRAGMA documentation (sqlite.org/pragma.html). It states that with
  auto-vacuum off freed pages go to a freelist and are reused but the file
  never shrinks, that incremental mode needs `PRAGMA incremental_vacuum(N)`,
  which removes up to N pages from the freelist, and that switching from the
  default mode needs a `VACUUM`.
- SQLite's VACUUM documentation (sqlite.org/lang_vacuum.html). It states that
  in write-ahead log mode the auto-vacuum property can still be changed by a
  `VACUUM`.
- Jim Nelson's article "SQLite, VACUUM, and auto_vacuum" (GNOME Blogs,
  2015-01-06). It proposes running `incremental_vacuum` at each startup with a
  small page count, so the free pages are chipped away without long locks.
- The read-only measurement of the live Run Database on 2026-10-08, recorded
  in the TechSpec: its size, its table counts, and the removal and compaction
  times on a copy of it.

## Glossary

- adds: **Run Retention**
- adds: **Run Retention Sweep**
- changes: **GC Command**
- changes: **Journal Retention**
- not a term: **Protected Runs** — a heading for the protection rule of Run Retention.

## Decisions

- The window is `store.run_retention_days`, User Config only, one of 7, 15
  or 30, 30 by default, and it cannot be turned off; see ADR-0255.
- Journal Retention keeps its key, scopes and default as the inner window;
  see ADR-0255.
- Queue-referenced Runs, Runs with an existing Run Worktree and Runs whose
  artifact directory cannot be proven are kept; see ADR-0255.
- The daily record lives in the Run Database, and the automatic sweep has a
  two-second budget; see ADR-0255.
- Compaction is incremental in the sweep and full only in an explicit
  command; see ADR-0255.

## Open Questions

- At the 30-day default this machine's file stays near its size, because
  Journal Retention already prunes journals at 14 days. Default until
  answered: the operator applies 30 days and reads the dry-run; choosing 7
  days, or a shorter Journal Retention, is what shrinks the file.
