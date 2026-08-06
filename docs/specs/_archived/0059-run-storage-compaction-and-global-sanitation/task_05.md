---
task: task_05
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill

## Overview

This Spec adds compaction, global sanitation, and a storage report to the CLI,
so the Skill must teach them before the Spec can close. This is the authorized
tooling Task.

## Requirements

1. MUST document the compaction preview and its refusals — Active Run, other
   writer, insufficient temporary capacity — and state that compaction is
   explicit and never an automatic side effect of a retention sweep.
2. MUST document global sanitation: dry-run first, discovery from durable Run
   metadata only, the six classifications, and that Review Artifacts and
   ambiguous paths are always preserved.
3. MUST document the read-only storage report and that it requires no flags and
   is safe anywhere.
4. MUST state that per-repository GC and Journal Retention are unchanged.
5. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document compaction, sanitation, and the report.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill documents the compaction preview and all three refusals.
- [ ] The Skill states compaction is explicit, never automatic.
- [ ] The Skill documents sanitation's dry-run default and six classifications.
- [ ] The Skill states Review Artifacts and ambiguous paths are preserved.
- [ ] The Skill documents the read-only storage report.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0059-run-storage-compaction-and-global-sanitation/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 2 and 4.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

Implementation stopped before changing the Roundfix Skill because the shipped
CLI does not expose the compaction command this Task requires the Skill to
teach. `internal/store` contains `PreviewCompaction` and `Compact`, but
`internal/cli/gc.go` routes only per-repository GC and `gc sanitize`; its usage
text likewise contains no compaction surface. Documenting
`roundfix gc compact [--apply]` would therefore contradict Requirement 8 by
describing behaviour that has not shipped, while adding that CLI route would
change behaviour outside this documentation-only Task.

Focused checks and inspections:

- `rtk proxy rg -n 'PreviewCompaction|\\.Compact\\(|gc compact|compact \\[--apply\\]' . --glob '!docs/specs/**' --glob '!skills/**' --glob '!.agents/**'`
  — found compaction only in `internal/store/journal.go` and its tests; no CLI
  caller or command route exists.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go run -buildvcs=false ./cmd/roundfix gc --help`
  — exited `0` and listed only `roundfix gc [--dry-run]` and
  `roundfix gc sanitize [--apply]`.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go run -buildvcs=false ./cmd/roundfix gc compact`
  — the built CLI returned Preflight Validation reason
  `unexpected argument "compact"` with its documented no-side-effects report.

Acceptance evidence:

- Compaction preview and three refusals: the store layer implements and tests
  Active Run, other-writer, and insufficient-temporary-capacity refusals, but
  no shipped CLI surface can be truthfully documented in the Skill.
- Explicit-only compaction: the store test
  `TestPruneTerminalRunsNeverVacuumsAllocatedPages` covers the invariant, but
  the missing CLI surface blocks the requested operator workflow text.
- Sanitation dry-run, six classifications, Review Artifact preservation, and
  ambiguous-path preservation: shipped source and focused Task 03 evidence
  support these claims, but no partial Skill edit was made while the compaction
  contract is unresolved.
- Read-only storage report: shipped source and focused Task 01 evidence support
  `roundfix storage report`, including no flags and use outside Git, but no
  partial Skill edit was made.
- Mirror identity: not regenerated because changing the canonical Skill would
  create a false CLI contract.
- Repository Verification: not run; the Daemon owns every command under this
  Task's `## Verification` section.
- No Go source changed by this Task turn.

Follow-up: add a bounded implementation Task for the promised explicit
compaction CLI preview/apply route, or revise the Spec and this Task to state
that compaction is a store-only API. Then rerun this Skill synchronization and
its required regeneration chain.
