---
type: fix
status: open
created: 2026-10-12
spec: null
reason: null
---

# A Run Worktree under a symlinked worktree root is never cleaned

## Symptom

Since `~/.roundfix/worktrees` became a symbolic link to
`/Volumes/Brain-SSD/cache/.roundfix/worktrees` on 2026-10-10, every Delivery
Queue merge ends with
`warning: cleanup failed: keep Run "<run-id>": recorded Run Worktree path is unsafe`,
and `roundfix reconcile <run-id>` reports `evidence: recorded Run Worktree path is unsafe`
with `action: preserve`. On 2026-10-12, `roundfix reconcile` printed the reason 550
times and the worktree root held 24 GB that nothing reclaims.

## Where

`internal/worktree/worktree.go`: reconciliation calls `recordedWorktree(result.Path, worktrees)`
and returns `reconciliationReasonWorktreePath` when the state is `recordedWorktreeUnsafe`. The
recorded path keeps the symbolic-link spelling under `~/.roundfix/worktrees`, while
`git worktree list` registers the resolved `/Volumes/...` path, so the two never match.

## Expected

A recorded Run Worktree path and a registered worktree path that resolve to the
same directory are the same worktree; cleanup and reconcile proceed as they do
without the link. A path that resolves outside the worktree root stays unsafe.

## Evidence

Queue items 0261 (`run_20261012T171408Z_a9139d7510b0db63`) and 0257
(`run_20261012T194228Z_1d1cb31ee8783296`, `run_20261012T191532Z_192c8d919680b754`)
on 2026-10-12, Roundfix 0.68.0 and 0.69.0.
