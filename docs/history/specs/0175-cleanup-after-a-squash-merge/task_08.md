---
task: task_08
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: medium
---

# Task 08: Automatic release proves its merge evidence and finds every Run of the repository

## Overview

Corrective Task from the pre-PR review of 2026-09-28. `ReleaseMergedRuns` in `internal/cli/deliver_workflow.go` checks only that the recorded merge commit and candidate head are non-empty strings; it never resolves them as commits or checks that the candidate head is reachable from (or squash-equivalent to) the merge, so malformed or stale evidence can drive cleanup or silently fall back to the default branch. And the automatic release in `internal/cli/deliver_release.go` lists Runs with `store.ListRunsQuery{GitRoot: repository}`, an exact Git-root match, while Runs record a durable `RepositoryRoot`; a Run started from another worktree of the same repository is missed and its worktree stays behind after item cleanup.

## Requirements

1. MUST resolve the recorded merge commit and candidate head with Git before any release, and refuse the release (keeping every Run, with a reason naming the bad value) when either does not resolve to a commit or the merge commit is not on the default branch.
2. MUST select the Spec's Runs by the repository's durable `RepositoryRoot` (the value `roundconfig.RepositoryRoot` returns), so Runs started from any worktree or checkout of the same repository are considered; add the store query support this needs without changing any existing query's results.
3. MUST keep every existing release test and message unchanged, and put new tests in their own files.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a named test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A merge record whose merge commit or candidate head does not resolve is refused and releases nothing; a valid record still releases.
- [ ] A terminal Run of the merged Spec started from a second worktree of the same repository is released; a Run of another repository is not.

## Context

- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_release.go`
- interface: `internal/store/store.go`
- creates: `internal/cli/deliver_release_evidence_test.go`
- creates: `internal/store/list_runs_repository_root_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit|TestReleaseMergedRunsRefusesAnUnresolvedCandidateHead|TestReleaseMergedRunsStillReleasesWithValidEvidence|TestAutomaticReleaseFindsARunStartedFromAnotherWorktree|TestAutomaticReleaseIgnoresAnotherRepository)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit TestReleaseMergedRunsRefusesAnUnresolvedCandidateHead TestReleaseMergedRunsStillReleasesWithValidEvidence TestAutomaticReleaseFindsARunStartedFromAnotherWorktree TestAutomaticReleaseIgnoresAnotherRepository; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && go test -count=1 ./internal/store >/dev/null` — expected: exit 0; before this Task the five new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented merge-evidence validation before automatic release enumerates or
mutates Runs. The workflow now resolves both recorded values as commits,
requires the merge commit on the detected default branch, and requires the
candidate head to be an ancestor of the merge or have the same tree as the
squash result. Refusals name the recorded value. Automatic release now resolves
the repository's durable key and uses a new exact `RepositoryRoot` store query;
the existing `GitRoot` query path is unchanged.

Focused-check evidence:

- Before production changes,
  `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 -run 'TestReleaseMergedRunsRefusesAnUnresolvedMergeCommit|TestListRunsScopesByRepositoryRoot' ./internal/cli ./internal/store`
  failed because unresolved evidence was accepted and `ListRunsQuery` had no
  `RepositoryRoot` field.
- `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 -run 'TestReleaseMergedRuns|TestAutomaticRelease|TestListRunsScopesByRepositoryRoot' ./internal/cli ./internal/store`
  passed after implementation. This covers each invalid merge-evidence case,
  a squash-equivalent valid candidate, the second-worktree Run, the
  other-repository exclusion, and the exact store query.
- `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache go test -count=1 -run 'TestDeliverRelease|TestReleaseMergedRuns|TestAutomaticRelease' ./internal/cli`
  passed, preserving the existing delivery-release cases and messages.
- `rtk env GOCACHE=/private/tmp/roundfix-task08-gocache make verify-incremental` passed
  after rerunning with macOS process-table access. The sandboxed attempt reached
  the full test target but two unrelated force-stop integration tests could not
  enumerate the process table; no product-test failure remained on the
  permitted rerun.

Acceptance evidence:

- Invalid merge commits, invalid candidate heads, off-default merge commits,
  and unrelated candidate heads all preserve the Run Worktree and branch;
  valid squash-equivalent evidence releases them.
- A terminal Run created from a second linked worktree is released through the
  shared durable repository key, while an otherwise matching Run from another
  repository remains present.

The Daemon-owned Verification command was not run in this Agent turn.
