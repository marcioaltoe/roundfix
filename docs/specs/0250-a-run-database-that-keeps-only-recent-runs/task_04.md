---
task: task_04
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: completed
type: docs
complexity: low
---

# Task 04: The glossary and the lifecycle policy describe Run Retention

## Overview

Writes ADR-0255 into the durable documents. `CONTEXT.md` gains **Run
Retention** and **Run Retention Sweep** and revises **GC Command** and
**Journal Retention**, which today say that nothing deletes a `runs` row. The
prose of the Run Database lifecycle policy stops saying that Roundfix bounds
only the Run Event Journal. It is verifiable on its own through phrase checks
and the markdown contracts.

## Requirements

1. MUST, through the `domain-modeling` skill, add **Run Retention** to
   `CONTEXT.md` directly after **Journal Retention**, with this definition
   and `_Avoid_` line:

   > The User Config window, `store.run_retention_days` of 7, 15 or 30 days
   > and 30 by default, after which a terminal Run leaves the Run Database
   > whole: its row, its Run Event Journal, its Agent Selection records, its
   > token usage and its artifact directory. A Run that is not terminal, that a
   > Delivery Queue references, whose Run Worktree still exists, or whose
   > artifact directory sits under an unproven root is kept (ADR-0255).
   > _Avoid_: Journal Retention, TTL, purge

2. MUST add **Run Retention Sweep** directly after **Run Retention**:

   > The step that applies Run Retention and then returns freed pages to the
   > filesystem by incremental compaction. It runs at most once a day at the
   > start of a Run or a Delivery Queue under a two-second budget, continuing
   > at the next start when the budget runs out, and without a budget in the
   > GC Command (ADR-0255).
   > _Avoid_: cleanup job, cron, vacuum

3. MUST revise **GC Command** to this definition, keeping its `_Avoid_` line:

   > The support command that reclaims Run storage: it prunes the Run Event
   > Journal and artifact directory of terminal Runs older than the Journal
   > Retention window, removes orphaned run artifact directories, runs the Run
   > Retention Sweep without a budget, and reports what it freed. A Run it
   > already emptied is not counted again. It never touches Active Runs or
   > active-run locks, and it removes a `runs` row only through Run Retention
   > (ADR-0255).

4. MUST revise **Journal Retention** so that its definition keeps its two
   existing sentences and its ADR-0033 reference and ends with:

   > It is the inner window: Run Retention removes the whole Run past the
   > outer window, so a Journal Retention longer than the Run Retention has no
   > effect past it (ADR-0255).

5. MUST, in `docs/user-guide/run-database-lifecycle.md`, replace the opening
   sentence that says Roundfix bounds only the Run Event Journal, and the
   paragraphs after the table that say the GC Command never deletes `runs`
   and the Run index needs no bound, with prose that starts its first
   paragraph with "Run Retention is the outer bound" and names Journal
   Retention as the inner one, the tables Run
   Retention never touches, and the 2026-10-08 measurement of ADR-0255. The
   table between the `durable-table-lifecycle` markers, which task_01
   changed, stays byte-identical, and the paragraph naming the lifecycle
   policy test stays.
6. MUST NOT reference any Spec or Finding from `CONTEXT.md`; it cites ADRs
   only.

## Subtasks

- [ ] Add the two new glossary entries.
- [ ] Revise the two existing entries.
- [ ] Rewrite the lifecycle policy prose around the table.
- [ ] Run the markdown contracts.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Run Retention** and **Run Retention Sweep** and no longer says the GC Command never touches `runs` rows.
- [ ] The lifecycle policy names Run Retention as the outer bound, and its table is unchanged from task_01.

## Context

- interface: `CONTEXT.md`
- interface: `docs/user-guide/run-database-lifecycle.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Verification

- `t="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; for term in '**Run Retention**:' '**Run Retention Sweep**:' '**GC Command**:' '**Journal Retention**:'; do printf '%s\n' "$t" | grep -qF -- "$term" || { printf 'missing term in CONTEXT.md: %s\n' "$term" >&2; exit 1; }; done; for p in 'The User Config window, .store.run_retention_days. of 7, 15 or 30 days' 'The step that applies Run Retention and then returns freed pages' 'runs the Run Retention Sweep without a budget' 'It is the inner window: Run Retention removes the whole Run past the outer window'; do printf '%s\n' "$t" | grep -q -- "$p" || { printf 'missing phrase in CONTEXT.md: %s\n' "$p" >&2; exit 1; }; done; ! printf '%s\n' "$t" | grep -q -- 'Never touches Active Runs, .runs. rows, or active-run locks' || { printf 'CONTEXT.md still says the GC Command never touches runs rows\n' >&2; exit 1; }; l="$(tr -s '[:space:]' ' ' < docs/user-guide/run-database-lifecycle.md)"; printf '%s\n' "$l" | grep -qF -- 'Run Retention is the outer bound' || { printf 'missing the outer bound in the lifecycle policy\n' >&2; exit 1; }; ! printf '%s\n' "$l" | grep -qF -- 'bounds only the Run Event Journal' || { printf 'the lifecycle policy still says only the journal is bounded\n' >&2; exit 1; }; make verify-docs` — expected: exit 0. Before this Task `CONTEXT.md` has no **Run Retention** entry, so the first check fails. After it, every glossary phrase is present, the old claims are gone from both files, and the markdown contracts and `roundfix spec check` pass.

## References

- `_prd.md` → Core Feature 6; Glossary
- `_techspec.md` → Glossary; Data Models; Build Order 4
- ADR-0255; ADR-0033

## Result

Implemented the glossary and lifecycle-policy prose for Daemon Verification.
`CONTEXT.md` now defines Run Retention and Run Retention Sweep, revises GC
Command and Journal Retention, and cites only ADRs in the new definitions.
The lifecycle guide now describes Run Retention as the outer bound and Journal
Retention as the inner bound, preserves the durable-table-lifecycle table, and
records the 2026-10-08 ADR-0255 measurement. Focused checks after the edits
confirmed the protected table hash is unchanged and the required glossary and
lifecycle phrases are present; the declared Verification commands were not
run, and Task status remains Daemon-owned.
