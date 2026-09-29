---
task: task_02
spec: 0183-run-storage-that-says-what-it-holds
status: completed
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

## Result

### Implementation

- `countRetainedTerminalRuns` now treats only a recorded-root error matching
  `fs.ErrNotExist` as an absent checkout and falls back to the Runs' distinct
  non-empty repository keys. Per-call caches validate each fallback key once
  and list branches from any resolved repository path at most once.
- Repository keys keep the recorded-root safety checks. A bare key is accepted
  only when `git rev-parse --absolute-git-dir` resolves to the key itself;
  absent keys add no failure, while unsafe or invalid existing keys do.
- Recorded Run Worktree paths are still checked for every Run. Symlink, file,
  non-root and Git-failure cases at a recorded root remain failures. No
  reconciliation, Branch Set or delivery-release path changed.

### Focused checks

- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run 'TestCountRetainedTerminalRuns' ./internal/worktree`
  passed after the final implementation edit. This selection includes the new
  absent-checkout, absent-key, bare-key and validation-failure cases and the
  unchanged `TestCountRetainedTerminalRunsBatchesGitInspectionByRepository`.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run 'TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted' ./internal/cli`
  passed against a real repository, removed linked checkout and temporary
  Roundfix Home.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental`
  passed outside the sandbox, including `go vet`, all Go packages, skill
  checks and the build. The sandboxed attempt had first failed in two
  force-stop tests because process-table access was denied and in one
  250-millisecond daemon timing test under concurrent load; the elevated rerun
  passed those same packages.
- The Daemon-owned commands under `## Verification` were not run.

### Acceptance evidence

- `TestCountRetainedTerminalRunsListsBranchesThroughTheKeyWhenTheCheckoutIsGone`
  passed: no failure was returned, the key was listed once and the existing
  Run Branch counted.
- `TestCountRetainedTerminalRunsCountsOnlyAnExistingWorktreeWhenTheRepositoryIsGone`
  passed: absent checkout and key produced no failure or Git call, and only the
  existing recorded Run Worktree counted.
- `TestCountRetainedTerminalRunsListsABareRepositoryKey` passed: the bare key's
  absolute Git directory matched the key and its Run Branch counted.
- `TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot` passed;
  separate tests also kept a non-directory root, a non-repository root, a Git
  validation failure and an invalid existing repository key loud.
- `TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted` passed: the Run
  remained visible, retained guidance stayed on stderr and no `warning:` line
  printed.
