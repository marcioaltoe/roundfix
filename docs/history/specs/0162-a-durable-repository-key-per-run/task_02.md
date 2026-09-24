---
task: task_02
spec: 0162-a-durable-repository-key-per-run
status: completed
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

## Result

Reconcile now compares the current repository key with both the Run's recorded
key and the key resolved from its checkout. A linked-worktree Run is therefore
accepted from the main checkout, while a Run whose recorded and resolved keys
belong to another repository is still refused.

Artifact Root discovery now supplies `gc --sanitize` with the recorded
repository key, falling back to `git_root` only for legacy rows whose key is
empty. Run Window reads first use the repository key and then the exact
checkout, while clear removes both possible rows; replacing a legacy window
moves it to the repository key.

Focused evidence:

- Before the production edits,
  `TestReconcileAcceptsALinkedWorktreeRunFromTheMainCheckout` exited non-zero
  because reconcile reported that the linked checkout did not belong to the
  main checkout. After the edits, its focused run exited 0.
- `TestReconcileRefusesAnotherRepositorysRun` exited 0 and observed the
  preflight refusal naming the Run, its unrelated repository, and the current
  repository.
- Before the production edits,
  `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval` observed
  `Classification: overridden`. After the edits, its focused run exited 0 and
  observed the shared root as `orphaned` after the real Git worktree was
  removed.
- Before the production edits, `TestRunWindowKeyedOnAWorktreeIsFound` reported
  no Run Window. After the edits, its focused run exited 0 after showing the
  legacy worktree-keyed cutoff and clearing that row.
- Focused regression groups for the existing sanitation matrix, Run Window
  commands, reconciliation and carry-forward paths, and store artifact/window
  behavior each exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1
  ./internal/cli ./internal/store` exited 0 with host access
  (`internal/cli` 96.511s, `internal/store` 16.663s).
- `rtk make verify-incremental
  GOCACHE=/private/tmp/roundfix-task02-gocache` exited 0 with host access.

The Daemon-owned `## Verification` command was not run.
