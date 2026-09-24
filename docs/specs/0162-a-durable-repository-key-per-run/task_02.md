---
task: task_02
spec: 0162-a-durable-repository-key-per-run
status: pending
type: backend
complexity: medium
---

# Task 02: Every lookup uses the key

## Overview

`reconcile <run-id>` compares checkout paths and refuses a linked-worktree Run the main checkout lists; `gc --sanitize` derives the default root from the Run's checkout and marks the shared root `overridden` once the worktree is removed; a Run Window keyed on a worktree path is no longer found.

## Requirements

1. MUST make `reconcile <run-id>` accept a Run whose recorded key, or resolved checkout, equals the current repository's key, and still refuse another repository's Run.
2. MUST make `gc --sanitize` derive a Run's default artifact root from its recorded key.
3. MUST make Run Window lookups fall back to a row keyed on the exact checkout, and `window clear` remove either.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, `reconcile <run-id>` from the main checkout accepts a linked-worktree Run and refuses an unrelated repository's Run.
- [ ] The shared root is not classified `overridden` after the worktree is removed.
- [ ] A worktree-keyed Run Window is shown and cleared.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reconcile.go`
- interface: `internal/cli/gc.go`
- interface: `internal/store/store.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileAcceptsALinkedWorktreeRunFromTheMainCheckout|TestReconcileRefusesAnotherRepositorysRun|TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval|TestRunWindowKeyedOnAWorktreeIsFound)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileAcceptsALinkedWorktreeRunFromTheMainCheckout TestReconcileRefusesAnotherRepositorysRun TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval TestRunWindowKeyedOnAWorktreeIsFound; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Consumers
