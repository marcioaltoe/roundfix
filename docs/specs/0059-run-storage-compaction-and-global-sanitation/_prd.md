---
spec: 0059-run-storage-compaction-and-global-sanitation
status: active
created: 2026-07-28
surfaces: [backend, cli, data, docs]
---

# Run storage compaction and global sanitation

The GC Command prunes Run Event Journal rows machine-wide but never returns
the deleted pages to the filesystem — the Run Database sat at 1.4 GB with
its journal already pruned — and its Artifact Directory cleanup resolves
only the current repository's Artifact Root, so a machine-wide database
coexists with repository-scoped artifact discovery and orphaned roots from
other repositories are never reclaimed. Tables added after the retention
contract (Run summaries, Agent Selection records) have no defined long-term
lifecycle at all. Evidence:
[sanitation is repository-scoped and does not compact SQLite](../../findings/2026-07-17-global-run-storage-sanitation-and-compaction.md).
The terminal Run Worktree half of that report shipped with Spec 0038; this
Spec owns the database and artifact halves.

## Project Constraints

- Identifier strategy: not applicable — Run IDs, Artifact Root paths, and
  table names keep their existing identities; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all behavior is local SQLite
  and filesystem maintenance; no authentication or HTTP surface. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0033 keeps the Run Event Journal
  pruned by its retention window (compaction reclaims pages, it never
  changes retention); ADR-0052 protects terminal completion — compaction
  must refuse while an Active Run or writer can exist; the Spec 0014
  retention contract's promise never to delete `runs` rows or Active Run
  locks is preserved unless this Spec's policy explicitly bounds a table.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Deleted Run Database pages are physically reclaimable through a guarded,
  explicit compaction with a deterministic preview.
- One global, dry-run-first sanitation pass can discover and clean every
  Roundfix-owned Artifact Root, not just the current repository's.
- Every durable table has a stated owner and retention rule, and a
  read-only storage report shows where the bytes are.
- No maintenance operation can touch an Active Run, a live writer, or an
  ambiguous path.

## User Stories

1. As an operator with a multi-gigabyte Run Database whose journal is
   already pruned, I want an explicit compaction with a preview of bytes
   before, reclaimable, and after, so that I can reclaim disk without
   hand-run SQLite operations.
2. As an operator running compaction, I want it to refuse while an Active
   Run or potential writer exists and to verify temporary disk capacity
   first, so that reclaiming space can never corrupt or block live work.
3. As a user with Runs from several repositories, I want a global
   sanitation mode that discovers every Roundfix-owned Artifact Root from
   durable Run metadata, so that orphaned artifacts outside the current
   repository stop accumulating forever.
4. As an operator deciding retention policy, I want a read-only storage
   report grouped by repository, state, table, and Artifact Root, so that
   policy decisions use measured data.
5. As a cautious operator, I want sanitation to classify roots that are
   missing, overridden, outside Roundfix Home, or unsafe, and preserve
   Review Artifacts and anything ambiguous, so that cleanup never guesses.

## Core Features

1. A guarded Run Database compaction: deterministic dry-run preview
   (bytes before, reclaimable, after), refusal while an Active Run or
   writer can be present, temporary-capacity verification, a failure-safe
   SQLite operation, and a completion report. Compaction is explicit —
   never an automatic side effect of retention sweeps — and takes no
   exclusive lock outside its own run.
2. A global sanitation mode, dry-run first, that discovers every
   Roundfix-owned Artifact Root from durable Run metadata, classifies each
   root (active, orphaned, missing, overridden, outside Roundfix Home,
   unsafe), removes only directories proven to belong to eligible or
   absent Runs, and leaves Review Artifacts and ambiguous paths unchanged.
   Per-repository GC remains available unchanged.
3. Table-by-table storage lifecycle policy: the compact Run index is
   preserved by default, Active Run locks stay governed by the Run
   lifecycle, and Agent Selection records prune only with their owning Run
   or under an explicit evidence-retention rule.
4. A read-only storage report grouped by repository, state, table, and
   Artifact Root, exposing measured bytes and row counts.

## User Experience

- The compaction preview and result read as three numbers an operator can
  compare; refusals name the live Run or writer that blocked it.
- Global sanitation lists each root with its classification and the proof
  that made it removable; ambiguity always reads `preserved`.
- The storage report is one command, no flags required, safe to run
  anywhere.

## Non-Goals / Out of Scope

- Changing Journal Retention windows or what retention deletes (ADR-0033).
- Automatic, threshold-triggered, or scheduled compaction.
- Deleting `runs` rows or Active Run locks outside the stated policy.
- Cloud or remote storage; everything stays in Roundfix Home and declared
  Artifact Roots.
- Terminal Run Worktree reconciliation, owned by Spec 0038.

## Success Metrics

- On a database with pruned journal rows, compaction's preview matches the
  measured reclaimable bytes and the post-compaction file size shrinks
  accordingly; an injected Active Run makes it refuse.
- Global sanitation on a machine with multiple repository roots removes
  only proven-eligible directories, preserves Review Artifacts and
  ambiguous paths, and is idempotent on a second run.
- The storage report's grouped totals reconcile with the database file
  size and artifact `du` totals within measurement tolerance.
- Every durable table appears in the documented lifecycle policy.

## Decisions

- Compaction is explicit and preview-gated; retention sweeps stay cheap and
  lock-free.
- Artifact Root discovery trusts only durable Run metadata — never
  filesystem guessing outside recorded roots.
- The Spec 0014 non-deletion promises hold unless the policy in this Spec
  explicitly bounds a table with measured justification.
- The `CONTEXT.md` glossary entries this Spec's behavior touches — at least
  `GC Command`, plus any new compaction, sanitation, or storage-report
  vocabulary — update in this Spec's documentation task, never ahead of
  implementation.

## Open Questions

None.
