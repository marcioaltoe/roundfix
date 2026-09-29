---
task: task_06
spec: 0183-run-storage-that-says-what-it-holds
status: pending
type: backend
complexity: low
---

# Task 06: Only a missing checkout counts as vanished

## Overview

task_02 makes the retained-Run count treat a recorded Git root that no longer exists as a fact rather than an error. It detects that case with `errors.Is(err, fs.ErrNotExist)` on the error `recordedGitRoot` returns. That error also wraps the failure of launching Git. When the `git` executable itself is missing, `exec` reports `ENOENT`, which also matches `fs.ErrNotExist`. `runs list` would then silently drop the warning instead of reporting that the inspection failed. The pre-PR review of the 0183 candidate found this.

## Requirements

1. MUST classify a recorded Git root as vanished only when the `os.Lstat` of the recorded path itself reports that it does not exist, for example through a dedicated sentinel or error type returned by `recordedGitDirectory`. A failure to launch or run Git, including `exec` `ENOENT`, MUST still be reported as an inspection failure.
2. MUST keep every behavior task_02 established for a checkout that really is gone, a bare repository key and a symlinked root.
3. MUST change no exported function signature and rename or remove no top-level test.
4. MUST put the new test in `internal/worktree/retained_git_unavailable_test.go`. Use a runner that fails the way a missing Git executable does, over an existing recorded root.

## Subtasks

- [ ] Distinguish a vanished recorded root from a Git launch failure.
- [ ] Add a test for the Git-unavailable case.

## Acceptance Criteria

- [ ] With an existing recorded root and a Git runner that fails with a missing-executable error, the count reports an inspection failure.
- [ ] A recorded root that does not exist is still counted as vanished, with no warning.

## Context

- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/retained_git_unavailable_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCountRetainedTerminalRunsReportsAGitLaunchFailure|TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone|TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCountRetainedTerminalRunsReportsAGitLaunchFailure TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new named test does not exist, so the command fails.

## References

- `_prd.md` → the retained-Run count
- task_02

## Result
