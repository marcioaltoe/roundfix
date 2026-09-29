---
task: task_02
spec: 0183-run-storage-that-says-what-it-holds
status: pending
type: backend
complexity: medium
---

# Task 02: The retained-Run count ignores a vanished checkout

## Overview

`runs list` counts terminal Runs that still keep a Run Worktree or Run Branch, and it inspects each Run's recorded Git root to list branches. Spec 0175 made Runs of since-removed linked checkouts visible from the main checkout through their repository key. Their recorded roots no longer exist, so each one produces a stderr warning on every call: 21 warnings on 2026-09-29 in this repository. This Task makes an absent recorded root a fact rather than a failure. The count lists that Run's branches through its repository key when the key still exists, and otherwise counts only an existing recorded Run Worktree. The Run facts come from the Run Database. Branch listings come from Git in the key repository, read-only, and the note goes to the operator's stderr.

## Requirements

1. MUST make `countRetainedTerminalRuns` in `internal/worktree/worktree.go` treat a group whose recorded Git root fails validation with an error wrapping `fs.ErrNotExist` as having no branches at that root, without appending a failure.
2. MUST, for such a group, look up each distinct non-empty `RepositoryRoot` of its Runs that differs from the recorded root, listing each key at most once per call: an absent key is skipped silently; a key that passes the existing recorded-root validation, or is a bare repository whose `git rev-parse --absolute-git-dir` equals the key, has its Run Branches listed; any other validation failure of an existing key is appended to the failures.
3. MUST keep checking each Run's recorded Run Worktree path, so an existing path still counts as retained, and MUST keep every other recorded-root validation failure (symlink, not a directory, not the repository root, Git failure) as a warning.
4. MUST NOT change what `reconcile`, the Branch Set classification or delivery release inspect; they keep their current behavior for a Run whose checkout is gone.
5. MUST put the unit tests in `internal/worktree/retained_vanished_checkout_test.go`, using the package's fake Git runner, and a public `runs list` test in `internal/cli/runs_list_vanished_checkout_test.go`, against a real repository with a removed linked worktree checkout and a temporary Roundfix Home. MUST keep `TestCountRetainedTerminalRunsBatchesGitInspectionByRepository` green without renaming any test.

## Subtasks

- [ ] Recognize an absent recorded root as a fact, not a failure.
- [ ] Fall back to the repository key, work tree or bare, listing each key once.
- [ ] Keep every other validation failure loud.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] With the recorded checkout removed and the repository key present, no failure is reported and a Run Branch that exists in the key repository is counted.
- [ ] With both the checkout and the key absent, no failure is reported and only an existing recorded Run Worktree path is counted.
- [ ] A bare repository key has its Run Branches listed.
- [ ] A symlinked recorded root still produces a failure.
- [ ] `roundfix runs list` over a Run whose linked checkout was deleted prints no `warning:` line on stderr.

## Context

- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/retained_vanished_checkout_test.go`
- creates: `internal/cli/runs_list_vanished_checkout_test.go`
- instruction: `internal/cli/runs.go`
- instruction: `docs/history/specs/0175-cleanup-after-a-squash-merge/references/2026-09-25-legacy-run-repository-key-from-its-worktree.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone|TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone|TestCountRetainedTerminalRunsListsABareRepositoryKey|TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone TestCountRetainedTerminalRunsListsABareRepositoryKey TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted"` — expected: exit 0; before this Task the test does not exist.

## References

- [_prd.md](_prd.md) — Goal 5; Core Feature 4; Success Metric 5
- [_techspec.md](_techspec.md) — `runs list` ignores a vanished checkout; API Contract 4; Testing Approach 4; Build Order 2
- [references/2026-09-29-runs-list-warns-for-runs-whose-checkout-is-gone.md](references/2026-09-29-runs-list-warns-for-runs-whose-checkout-is-gone.md)
- ADR-0023; ADR-0053
