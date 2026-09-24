---
task: task_04
spec: 0162-a-durable-repository-key-per-run
status: pending
type: backend
complexity: medium
---

# Task 04: Every lookup, and only the recorded key

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `LatestKeptSpecRun`, which backs `settle`, still matches only `git_root` aliases, so a kept Run Worktree of a removed linked worktree is not found; the migration backfills the dead checkout path of a removed worktree as a key, because `RepositoryRoot` returns a non-existent path unchanged; and `sameRepository` consults the checkout path even when the Run has a recorded key, so a reused path makes reconcile accept another repository's Run.

## Requirements

1. MUST make `LatestKeptSpecRun` match `repository_root`, falling back to the alias match only for rows without a key, like `ListRuns`.
2. MUST backfill a key only when `git_root` exists and is a Git worktree, leaving the others empty.
3. MUST make `sameRepository` consult the checkout path only when the Run has no recorded key.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, `settle`'s lookup finds a kept Run from a removed linked worktree.
- [ ] A migrated row whose worktree was removed keeps an empty key.
- [ ] A Run with a recorded key of another repository is refused even when its old checkout path now belongs to this one.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/cli/reconcile.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSettleFindsARemovedWorktreeRunByKey|TestMigrationLeavesARemovedWorktreeKeyEmpty|TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath)$" ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettleFindsARemovedWorktreeRunByKey TestMigrationLeavesARemovedWorktreeKeyEmpty TestReconcileTrustsTheRecordedKeyOverTheCheckoutPath; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Consumers
