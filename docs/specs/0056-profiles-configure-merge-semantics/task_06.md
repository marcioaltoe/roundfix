---
task: task_06
spec: 0056-profiles-configure-merge-semantics
status: completed
type: docs
complexity: low
---

# Task 06: Document the merge contract

## Overview

Two observable behaviors change in ways a reader cannot infer from the command
itself: a fragment no longer deletes by omission, and a declined write no longer
exits zero. This Task records the merge contract, the removal declaration, the
per-category summary, and the exit-code table in the command reference, so an
operator and an Agent both find the new semantics where they already look.

## Requirements

1. MUST document that a fragment merges by category: the categories it names are
   replaced atomically and every other configured category is preserved.
2. MUST document the removal declaration as the only way to delete a category,
   and state that omission no longer deletes.
3. MUST document the per-category summary shown before a write, naming the three
   classifications.
4. MUST document the exit-code table for this command, distinguishing an applied
   write, a no-op, a dry run, a refusal, and a validation failure.
5. MUST state that proof covers the categories the operation writes.
6. MUST NOT change any Go source, workflow, or tooling file.

## Subtasks

- [x] Document merge-by-category and the preservation guarantee.
- [x] Document the removal declaration and the end of deletion by omission.
- [x] Document the summary and the exit-code table.
- [x] State the proof scope.

## Acceptance Criteria

- [x] The command reference describes merge-by-category and states that
      unnamed categories are preserved.
- [x] It documents the removal declaration and says omission does not delete.
- [x] It lists the three summary classifications.
- [x] It gives an exit code for each of: applied, no-op, dry run, refusal, and
      validation failure.
- [x] It states that proof covers the written categories.
- [x] `git status --porcelain` shows no path outside `docs/user-guide/` and this
      task file.

## Context

- instruction: `docs/agents/docs-layout.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `grep -qi "merge" docs/user-guide/commands.md` — expected: exit 0.
- `grep -q "profiles configure" docs/user-guide/commands.md` — expected: exit 0.
- `grep -qi "remove" docs/user-guide/commands.md` — expected: exit 0; the
  removal declaration is documented.
- `grep -qi "replaced" docs/user-guide/commands.md` — expected: exit 0; the
  summary classifications are documented.
- `git diff --name-only HEAD -- internal/ Makefile .github/ | grep -q . && exit 1 || exit 0`
  — expected: exit 0; this task changed no code, workflow, or tooling.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0; the
  behavior this documents still holds.

## References

- `_prd.md` → User Experience; Decisions.
- `_techspec.md` → API Contracts; Build Order 6.
- ADR-0086.

## Result

### Implementation

- The command reference now defines category-level merge semantics: named
  categories replace complete profiles atomically, while unnamed categories
  remain configured.
- It documents repeatable `--remove <category>` as the only removal declaration,
  states that omission does not delete, and identifies a fragment/removal
  conflict as validation failure.
- It records the pre-write `added`, `replaced`, and `removed` classifications,
  scopes Exact Agent Selection Proof to categories the operation adds or
  replaces, and gives separate exit-code rows for applied writes, no-ops, dry
  runs, refusals, and validation failures.

### Focused checks

- `rtk git diff --check` — exit `0`; the documentation diff has no whitespace
  errors.
- `rtk proxy sed -n '182,238p' docs/user-guide/commands.md` — exit `0`; manual
  inspection found the usage line, merge and removal contract, three summary
  classifications, written-category proof scope, and all five exit situations.
- `rtk git status --short` — exit `0`; it listed only
  `docs/user-guide/commands.md` and this task file.

### Acceptance-criterion evidence

1. The first `profiles configure` paragraph says a fragment merges by Agent
   Work Category, named categories replace atomically, and every other
   configured category is preserved.
2. The same paragraph says omission does not delete and `--remove <category>`
   is the only removal mechanism.
3. The summary paragraph names `added`, `replaced`, and `removed` and says the
   summary appears before a write.
4. The exit-code table gives `0` for an applied write, already-satisfied no-op,
   and dry run; `1` for refusal; and `2` for validation failure.
5. The proof paragraph says Exact Agent Selection Proof covers categories the
   operation adds or replaces—the categories it writes—and excludes untouched
   categories.
6. The focused status inspection listed no path outside `docs/user-guide/` and
   this task file.

The Daemon-owned commands in `## Verification` were not run during this Agent
turn.
