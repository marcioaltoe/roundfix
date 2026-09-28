---
task: task_06
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: medium
---

# Task 06: Protect nested worktrees and lock every worktree administration

## Overview

`removeUnregisteredItemWorktree` in `internal/worktree/worktree.go` deletes a
half-removed item directory with `os.RemoveAll`, even when a different
registered worktree was created inside it. That nested worktree's uncommitted
files are lost and its registration dangles; Spec 0168 recorded this as a
limit. Two production callers also still run `git worktree` outside the
administration lock Spec 0171 added: `discardSupersededBranch` in
`internal/daemon/reconcile.go` and the disposable checkout in
`internal/speccheck/checkout.go`. task_05 moves the carry-forward staging, the
third caller. The item path is recorded in the Delivery Queue, and the
registered worktree list comes from the repository's Git administration. The
removal destroys a directory tree in the user's machine, so a nested live
worktree must stop it.

## Requirements

1. MUST make `removeUnregisteredItemWorktree` list the registered worktrees of
   `ref.UserRoot` before `os.RemoveAll`, and refuse, without removing
   anything, when any registered worktree's canonical path lies under
   `ref.Path`. The error MUST name the nested path, and `CleanupItem` MUST
   return it, so delivery records it as the cleanup warning.
2. MUST add `internal/worktree/disposal.go` with
   `RemoveRegisteredWorktree(ctx, gitRoot, path string) error` (`git worktree
   remove <path>`, no force), `AddDetachedWorktree(ctx, repository, path,
   revision string) error` and `RemoveWorktreeForce(ctx, repository, path
   string) error`, each through `runWorktreeCommand`.
3. MUST make `discardSupersededBranch` remove the Run Worktree through
   `RemoveRegisteredWorktree`, keeping its record-first order and every error
   message it returns today.
4. MUST make `internal/speccheck/checkout.go` create and remove its disposable
   checkout through `AddDetachedWorktree` and `RemoveWorktreeForce`, keeping
   `ErrDisposableCheckoutCreate` and its cleanup order.
5. MUST add `internal/worktree/adminlock_guard_test.go` with a guard test
   that parses every non-test Go file under `internal/` and `cmd/` outside
   `internal/worktree`. The test fails when a call passes the string literal
   `"worktree"` followed by `"add"`, `"remove"`, `"prune"` or `"move"` as
   consecutive arguments. The file also needs a separate test proving the
   detector reports a synthetic direct call, so the guard cannot pass
   vacuously.
6. MUST put the other new tests in `internal/worktree/item_nested_cleanup_test.go`
   and `internal/worktree/disposal_test.go`, and MUST NOT remove or rename an
   existing test.

## Subtasks

- [ ] Refuse the nested removal and add the exported locked operations.
- [ ] Route the Daemon disposal and the disposable checkout through them.
- [ ] Add the guard and a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A half-removed item holding a nested registered worktree is refused, and the nested worktree keeps its uncommitted file and registration.
- [ ] A half-removed item without a nested worktree is still removed.
- [ ] The exported removal and add wait while another holder has the administration lock.
- [ ] The guard reports a synthetic direct call and passes on the tree after the routing; no production file outside `internal/worktree` administers a worktree directly.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/disposal.go`
- creates: `internal/worktree/disposal_test.go`
- creates: `internal/worktree/item_nested_cleanup_test.go`
- creates: `internal/worktree/adminlock_guard_test.go`
- interface: `internal/daemon/reconcile.go`
- interface: `internal/speccheck/checkout.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCleanupItemRefusesAHalfRemovedItemHoldingARegisteredWorktree|TestCleanupItemRemovesAHalfRemovedItemWithoutANestedWorktree|TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock|TestRemoveRegisteredWorktreeKeepsADirtyWorktree|TestNoProductionWorktreeAdministrationBypassesTheAdminLock|TestAdminLockGuardReportsADirectWorktreeCall|TestCleanupPreservesAnUnregisteredPathWithoutThisRepositorysMarker|TestCleanupPreservesAnUnregisteredPathWhileItsAdminDirectoryExists|TestWriteBranchDispositionRecordSurvivesDiscard|TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks|TestDisposableCheckoutIsolatesOperatorTreeAndIndex|TestDisposableCheckoutCleanupRunsWhenCallerFailsOrPanics|TestDisposableCheckoutCreationRefusalHasNamedReason|TestReconcileDiscardWritesTheRecordBeforeRemoving|TestReconcileDiscardKeepsSurfaceWhenRecordCannotBeWritten)$" ./internal/worktree ./internal/daemon ./internal/speccheck ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCleanupItemRefusesAHalfRemovedItemHoldingARegisteredWorktree TestCleanupItemRemovesAHalfRemovedItemWithoutANestedWorktree TestRemoveRegisteredWorktreeWaitsForTheAdminLock TestAddDetachedWorktreeWaitsForTheAdminLock TestRemoveRegisteredWorktreeKeepsADirtyWorktree TestNoProductionWorktreeAdministrationBypassesTheAdminLock TestAdminLockGuardReportsADirectWorktreeCall TestCleanupPreservesAnUnregisteredPathWithoutThisRepositorysMarker TestCleanupPreservesAnUnregisteredPathWhileItsAdminDirectoryExists TestWriteBranchDispositionRecordSurvivesDiscard TestSupersededBranchIsClassifiedWhenLaterIntegratedRunCoveredTasks TestDisposableCheckoutIsolatesOperatorTreeAndIndex TestDisposableCheckoutCleanupRunsWhenCallerFailsOrPanics TestDisposableCheckoutCreationRefusalHasNamedReason TestReconcileDiscardWritesTheRecordBeforeRemoving TestReconcileDiscardKeepsSurfaceWhenRecordCannotBeWritten; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the new named tests exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Item cleanup and the lock for every caller
- `_prd.md` → Core Features 5-6; Success Metric 5
- `_techspec.md` → API Contract 6
