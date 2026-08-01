---
task: task_06
spec: 0056-profiles-configure-merge-semantics
status: pending
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

- [ ] Document merge-by-category and the preservation guarantee.
- [ ] Document the removal declaration and the end of deletion by omission.
- [ ] Document the summary and the exit-code table.
- [ ] State the proof scope.

## Acceptance Criteria

- [ ] The command reference describes merge-by-category and states that
      unnamed categories are preserved.
- [ ] It documents the removal declaration and says omission does not delete.
- [ ] It lists the three summary classifications.
- [ ] It gives an exit code for each of: applied, no-op, dry run, refusal, and
      validation failure.
- [ ] It states that proof covers the written categories.
- [ ] `git status --porcelain` shows no path outside `docs/user-guide/` and this
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
