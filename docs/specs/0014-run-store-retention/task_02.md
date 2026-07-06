---
task: task_02
spec: 0014-run-store-retention
status: pending
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

- [ ] `gc` command parsing/dispatch with `--dry-run`
- [ ] Live path: prune journals + remove pruned-Run artifact dirs + orphan dirs
- [ ] Dry-run path: read-only eligibility listing, no changes
- [ ] Freed-counts report (Runs, rows, bytes); zero-retention message
- [ ] Tests: dry-run no-op, live prune + artifact removal, orphan dir, Active kept

## Acceptance Criteria

- [ ] `roundfix gc` prunes eligible terminal Runs' journals and artifact dirs and reports Runs/rows/bytes freed.
- [ ] `roundfix gc --dry-run` changes nothing and lists exactly the set a live run would prune.
- [ ] An orphaned `runs/<id>` directory (no `runs` row) is removed; an Active Run's directory is kept.
- [ ] With `journal_retention: 0`, gc prunes nothing and reports so.

## Verification

- `rtk go test ./internal/cli/` — expected: gc command tests pass.
- `rtk go run ./cmd/roundfix gc --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 2, 3; Core Feature 2. `_techspec.md` → GC Command,
Build Order 2. CONTEXT.md → GC Command. ADR-0033.
