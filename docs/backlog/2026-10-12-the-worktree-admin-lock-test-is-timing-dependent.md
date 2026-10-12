---
type: fix
status: open
created: 2026-10-12
spec: null
reason: null
---

# The worktree admin-lock test fails under load

## Symptom

`TestRemoveRegisteredWorktreeWaitsForTheAdminLock` failed the repository gate of
the Delivery Queue item for Spec 0256 on 2026-10-12, then passed three times out
of three when run alone and passed on the retry.

## Where

`internal/worktree/disposal_test.go`: the test races the removal against a held
Git admin lock and depends on how fast the lock holder runs.

## Expected

The test proves that the removal waits for the admin lock through an explicit
synchronization point, not through elapsed time, and passes under a parallel
full suite.

## Evidence

Delivery Queue item `0256-an-archive-record-that-survives-a-squash-merge`,
repository gate on 2026-10-12; the retry merged as PR #495.
