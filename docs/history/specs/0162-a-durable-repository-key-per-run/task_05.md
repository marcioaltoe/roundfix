---
task: task_05
spec: 0162-a-durable-repository-key-per-run
status: completed
type: backend
complexity: medium
---

# Task 05: Removed worktrees and bare layouts

## Overview

Corrective Task from the second pre-PR review of 2026-09-24. `ListRuns` now lists a terminal Run from a removed linked worktree, and `InspectTerminalRun` then stats the dead checkout, so bare `reconcile` and `reconcile <run-id>` from the main checkout exit 1 for the whole repository after every parallel delivery. And in a bare-repository layout `RepositoryRoot` returns each worktree's own path, so a Run's key is per worktree and sibling worktrees no longer list it, including after the backfill.

## Requirements

1. MUST inspect a terminal Run whose recorded checkout is gone, but whose recorded key matches the current repository, against the current repository instead of the dead checkout, so `reconcile` and `reconcile --apply` from the main checkout exit 0 and treat the Run's branch and records as the evidence allows.
2. MUST derive the repository key from the common Git directory when it is not `<root>/.git`, so every worktree of a bare-repository layout shares one key, for creation and backfill alike.
3. MUST keep the main checkout's existing key unchanged for non-bare repositories.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, after `git worktree remove`, bare `reconcile` and `reconcile <run-id>` from the main checkout exit 0.
- [ ] With real Git, two worktrees of a bare repository share one key, and a Run from one is listed from the other.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/cli/reconcile.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileFromMainAfterALinkedWorktreeIsRemoved|TestBareRepositoryWorktreesShareOneKey|TestBareRepositoryWorktreeRunsListFromASibling)$" ./internal/config ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileFromMainAfterALinkedWorktreeIsRemoved TestBareRepositoryWorktreesShareOneKey TestBareRepositoryWorktreeRunsListFromASibling; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Consumers

## Result

- Implementation: repository identity now uses a linked worktree's common Git
  directory as the key when that directory is not a main checkout's `.git`.
  Reconciliation uses the current live checkout for Git proof after
  repository-key-scoped listing has established that each Run belongs to it;
  the persisted checkout path remains unchanged.
- Acceptance criterion 1 evidence: `rtk env
  GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestReconcileFromMainAfterALinkedWorktreeIsRemoved|TestReconcileAcceptsALinkedWorktreeRunFromTheMainCheckout|TestReconcileRefusesAnotherRepositorysRun|TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath)$'
  ./internal/cli` exited 0. The real-Git removed-worktree case covers bare
  `reconcile`, `reconcile <run-id>`, and `reconcile --apply`; the adjacent
  cases retain linked-worktree acceptance and cross-repository refusal.
- Acceptance criterion 2 evidence: `rtk env
  GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestBareRepositoryWorktreesShareOneKey|TestRepositoryIdentityIsSharedByLinkedWorktrees|TestRepositoryRootDoesNotEnumerateWorktrees)$'
  ./internal/config` exited 0, and `rtk env
  GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run
  '^(TestBareRepositoryWorktreeRunsListFromASibling|TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout|TestARemovedWorktreeRunStaysListed)$'
  ./internal/store` exited 0. These real-Git cases prove two bare-repository
  worktrees share the common Git-directory key, a Run created in one lists
  from the other, and ordinary non-bare repository identity stays unchanged.
- Repository incremental check: `rtk env
  GOCACHE=/tmp/roundfix-task05-gocache make verify-incremental` exited 0 after
  vet, the Go suite, skill checks, and the CLI build.
- Daemon Verification was not run; the Daemon owns the declared command and
  terminal Task status.
