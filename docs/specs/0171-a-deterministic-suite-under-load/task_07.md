---
task: task_07
spec: 0171-a-deterministic-suite-under-load
status: pending
type: backend
complexity: medium
---

# Task 07: Worktree administration is serialized per repository

## Overview

Corrective Task from QA finding F-001 of 2026-09-25. Under a concurrent `go test ./...`, one of 20 stress iterations of `TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork` failed because Implement could not create `task_01`'s worktree: Git reported `fatal: failed to read .git/worktrees/<task_02>/commondir: Result too large`. `internal/worktree/worktree.go` runs `git worktree add`, `remove`, `prune` and `list` for sibling Task worktrees, Run worktrees and delivery item worktrees concurrently with no serialization, so one `git worktree add` reads another worktree's administrative files while Git is still writing them. The race is in production too: the Daemon creates Task worktrees in parallel, and a delivery owner and a separate Implement Run can administer worktrees of the same repository at the same time.

## Requirements

1. MUST serialize every `git worktree add`, `git worktree remove`, `git worktree prune` and `git worktree move` that `internal/worktree` runs, per repository (keyed by the resolved Git common directory), both inside one process and across Roundfix processes (an advisory file lock, for example `flock` on a lock file under the Git common directory or under Roundfix Home keyed by that directory).
2. MUST hold the lock only for the Git administrative command itself, never across Worktree Bootstrap, copy, or Agent work, and MUST release it on every error path and on context cancellation.
3. MUST keep waiting for the lock bounded by the caller's context, returning the context error when it ends first.
4. MUST NOT change any public output, exit code, or worktree path layout.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a named test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Many concurrent Task worktree creations in one repository all succeed, and at no time do two worktree administrative commands for that repository run at once (observed through a recording command runner).
- [ ] Two independent processes (a helper subprocess) administering worktrees of the same repository never overlap.
- [ ] A caller whose context ends while waiting for the lock returns the context error without running the command.
- [ ] Commands for two different repositories are not serialized against each other.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/worktree_test.go`
- creates: `internal/worktree/adminlock.go`
- creates: `internal/worktree/adminlock_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestConcurrentTaskWorktreeCreationIsSerialized|TestWorktreeAdministrationIsSerializedAcrossProcesses|TestWorktreeAdministrationWaitEndsWithTheContext|TestWorktreeAdministrationOfDifferentRepositoriesRunsConcurrently)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestConcurrentTaskWorktreeCreationIsSerialized TestWorktreeAdministrationIsSerializedAcrossProcesses TestWorktreeAdministrationWaitEndsWithTheContext TestWorktreeAdministrationOfDifferentRepositoriesRunsConcurrently; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && out2="$(go test -count=10 -cpu 1,4 -run "^TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out2"; exit 1; }` — expected: exit 0; before this Task the four new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
- [qa/qa-report-2026-09-25.md](qa/qa-report-2026-09-25.md) — F-001
