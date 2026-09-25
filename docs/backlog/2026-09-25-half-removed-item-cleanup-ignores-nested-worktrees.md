---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Cleaning up a half-removed item directory ignores worktrees registered inside it

## Symptom

`roundfix deliver` cleanup deletes a truly half-removed item directory whole,
even when a different registered worktree was created inside it; that nested
worktree's uncommitted files are lost and its registration dangles.

## Where

`internal/worktree/worktree.go` `removeUnregisteredItemWorktree`
(`os.RemoveAll(ref.Path)`).

## Expected

Refuse the removal, with the cleanup warning, when any registered worktree
resolves by canonical path to a location under the item path.

## Evidence

Recorded limit of Spec 0168; probe `TestAtkNestedRegistered` in the focused
verification of Task 08 on 2026-09-25.
