---
task: task_02
spec: 0014-run-store-retention
status: completed
type: backend
complexity: medium
---

# Task 02: GC Command with dry-run and artifact cleanup

## Overview

Add `roundfix gc`, the support command that reclaims Run storage on demand: it
prunes eligible terminal Runs' journals (via the store primitive), removes their
artifact directories and orphaned `runs/<id>` directories, and reports what it
freed. `--dry-run` previews the same set without changing anything.

## Requirements

1. MUST add a non-interactive `gc` command that resolves the retention window,
   computes the cutoff, and prunes eligible terminal Runs' journal rows plus
   their `<artifact_dir>/runs/<run-id>` directories, and removes orphaned
   `runs/<id>` directories that have no matching `runs` row.
2. MUST support `--dry-run` that runs a read-only eligibility query and prints
   the would-prune set without deleting anything.
3. MUST report the Runs pruned, journal rows removed, and artifact bytes
   reclaimed on stdout, with diagnostics on stderr and stable exit codes; a
   `journal_retention: 0` MUST prune nothing and say so.
4. MUST only remove `runs/<id>` directories under the resolved run artifact
   root — never review artifacts under the spec tree or unrelated paths.

## Subtasks

- [x] `gc` command parsing/dispatch with `--dry-run`
- [x] Live path: prune journals + remove pruned-Run artifact dirs + orphan dirs
- [x] Dry-run path: read-only eligibility listing, no changes
- [x] Freed-counts report (Runs, rows, bytes); zero-retention message
- [x] Tests: dry-run no-op, live prune + artifact removal, orphan dir, Active kept

## Acceptance Criteria

- [x] `roundfix gc` prunes eligible terminal Runs' journals and artifact dirs and reports Runs/rows/bytes freed.
- [x] `roundfix gc --dry-run` changes nothing and lists exactly the set a live run would prune.
- [x] An orphaned `runs/<id>` directory (no `runs` row) is removed; an Active Run's directory is kept.
- [x] With `journal_retention: 0`, gc prunes nothing and reports so.

## Verification

- `rtk go test ./internal/cli/` — expected: gc command tests pass.
- `rtk go run ./cmd/roundfix gc --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 2, 3; Core Feature 2. `_techspec.md` → GC Command,
Build Order 2. CONTEXT.md → GC Command. ADR-0033.

## Result

Implemented `roundfix gc` with `--dry-run`, Journal Retention cutoff handling,
read-only terminal Run eligibility listing, live journal pruning, scoped
`<artifact_root>/runs/<id>` cleanup, orphan run-artifact cleanup, and stdout
reports for Runs, journal rows, and artifact bytes.

Acceptance evidence:

- Live GC: `TestRunGCPrunesEligibleJournalsArtifactsAndOrphans` verifies
  `roundfix gc` reports `Runs pruned: 1`, `Journal rows removed: 2`, and
  `Artifact bytes reclaimed: 16`; deletes only the eligible terminal Run's Run
  Events and run artifact directory; removes the orphaned `runs/<id>` directory;
  and preserves all `runs` rows.
- Dry run: `TestRunGCDryRunListsEligibleRunsAndChangesNothing` verifies
  `roundfix gc --dry-run` reports the same eligible Run, row count, bytes, and
  orphan ID while leaving Run Events and artifact directories unchanged.
- Active safety: the live GC test verifies an old Active Run keeps its journal
  and run artifact directory, and the review artifact path under `reviews/`
  remains present.
- Zero retention: `TestRunGCSkipsWhenJournalRetentionIsZero` verifies
  `journal_retention: 0` prints `GC skipped`, `Journal Retention: 0`, and
  `No pruning performed.`, while leaving journals and artifact directories in
  place.

Verification evidence:

- `rtk go test ./internal/cli/` passed: `Go test: 287 passed in 1 packages`.
- `rtk go run ./cmd/roundfix gc --help` passed and printed concise `roundfix gc
  [--dry-run]` usage with the Journal Retention and safety contract.
- `rtk go test ./...` passed: `Go test: 798 passed in 18 packages`.
- Repo gate `rtk make verify` passed: full Go suite, `roundfix skills check`,
  and build completed.
