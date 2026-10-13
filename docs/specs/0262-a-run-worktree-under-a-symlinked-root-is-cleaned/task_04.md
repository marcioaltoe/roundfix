---
task: task_04
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: pending
type: test
complexity: medium
---

# Task 04: The disposal lock tests wait on a synchronization point

## Overview

`TestRemoveRegisteredWorktreeWaitsForTheAdminLock` and
`TestAddDetachedWorktreeWaitsForTheAdminLock` prove waiting by letting a
250 ms deadline expire, and that deadline also bounds a real
`git rev-parse` that runs before the lock; under load the process is killed
and the test fails (`_techspec.md` → Current behavior). This Task moves the
exported disposal functions onto the package's runner seam and rewrites both
tests to cancel at a synchronization point. It is verifiable on its own:
the tests contain no deadline and pass 50 consecutive times under the race
detector.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [the worktree admin-lock test fails under load](references/2026-10-12-the-worktree-admin-lock-test-is-timing-dependent.md),
   by implementing `_techspec.md` → API Contract 4 for both disposal tests;
   the second test has the same defect and is the rest of its class.
2. MUST make `RemoveRegisteredWorktree` and `AddDetachedWorktree` delegate to
   unexported `removeRegisteredWorktree` and `addDetachedWorktree`, which take
   a `gitRunner`, keeping the exported signatures and behavior.
   `RemoveWorktreeForce` may follow the same pattern.
3. MUST rewrite the two tests, keeping their names and `t.Parallel()` first:
   each holds the lock with `holdDisposalAdminLock`, runs the operation in a
   goroutine with a cancellable context and a test runner that delegates to
   `execGitRunner{}`, waits for the runner's report that the Git common
   directory is resolved, cancels, and requires `context.Canceled`, no
   `worktree` command started while the lock was held, and the worktree
   unchanged (still registered for removal, still absent for add). After
   releasing the holder it runs the same operation with the test's context
   and requires success: the worktree is removed, or added at its path.
4. MUST NOT pass a deadline to the operation under test, sleep, or decide any
   assertion on elapsed time; a bounded wait may only fail a test that would
   otherwise hang.
5. MUST leave `TestRemoveRegisteredWorktreeKeepsADirtyWorktree` and the
   admin-lock tests in `internal/worktree/adminlock_test.go` unchanged.

## Subtasks

- [ ] Add the unexported runner-taking disposal functions.
- [ ] Rewrite the two lock tests around the synchronization point.

## Acceptance Criteria

- [ ] Neither disposal test passes a deadline to the operation it proves.
- [ ] Both pass 50 consecutive times with `-race -parallel 16`.

## Consumers

None. The exported disposal functions keep their signatures and behavior; only their bodies delegate to unexported functions in the same file, and the change is otherwise confined to one test file.

## Context

- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- instruction: `internal/worktree/adminlock.go`
- instruction: `internal/worktree/adminlock_test.go`
- interface: `internal/worktree/disposal.go`
- interface: `internal/worktree/disposal_test.go`

## Verification

- `if grep -nE 'context[.]With(Timeout|Deadline)' internal/worktree/disposal_test.go; then printf 'a disposal test still proves waiting by a deadline\n' >&2; exit 1; fi; grep -q 'func removeRegisteredWorktree(' internal/worktree/disposal.go || { printf 'the disposal runner seam is missing\n' >&2; exit 1; }; out="$(go test -count=50 -race -parallel 16 -v -run '^(TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; test "$(printf '%s\n' "$out" | grep -c -- '--- PASS: TestRemoveRegisteredWorktreeWaitsForTheAdminLock ')" -eq 50 || { printf 'the removal lock test did not pass 50 times\n' >&2; exit 1; }; test "$(printf '%s\n' "$out" | grep -c -- '--- PASS: TestAddDetachedWorktreeWaitsForTheAdminLock ')" -eq 50 || { printf 'the add lock test did not pass 50 times\n' >&2; exit 1; }; go test -count=1 -run 'Admin|Disposal|RemoveRegistered|AddDetached|Lock' ./internal/worktree || exit 1` — expected: exit 0. Before this Task both tests use `context.WithTimeout`, so the command fails; afterwards they pass 50 times each under the race detector and the other lock tests still pass.

## References

- `_prd.md` → Goals; User Story 5; Core Feature 6; Success Metric 4
- `_techspec.md` → Current behavior; API Contract 4; Invariant 5; Testing Approach; Build Order 4
- ADR-0268; ADR-0089; ADR-0125
