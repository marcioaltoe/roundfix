---
task: task_01
spec: 0136-a-rename-the-committer-can-stage
status: completed
type: backend
complexity: medium
---

# Task 01: Stage only paths Git can match

## Overview

`git add -f -- <path>` matches the worktree and the index together. After
`git mv`, a rename's source is in neither, because the index already records
the rename, so Git answers `pathspec did not match any files` and the commit
fails. The stageable-path filter already drops paths for reasons it records;
absence from both the worktree and the index is one more. Classification is
untouched: it reads the snapshot pair, not the staged list.

## Requirements

1. MUST drop a path absent from both the worktree and the index from the
   staged list, recording it as dropped with a reason rather than dropping it
   silently.
2. MUST keep staging a deletion, because a deleted path remains in the index
   until the deletion is recorded.
3. MUST prove a Task that renames a file with `git mv` commits, in a repository
   the Task itself creates and removes.
4. MUST prove a Task that renames a file without staging it also commits.
5. MUST keep the governed-rename refusal, so dropping the source from staging
   does not drop it from classification.
6. MUST keep every existing drop reason and the existing real-repository
   staging behavior.
7. MUST NOT change the changed-path reader, which Spec 0135 settled, and MUST
   NOT weaken or delete an existing assertion.

## Subtasks

- [ ] Add the absent-from-both drop reason to the stageable filter.
- [ ] Prove a staged rename commits.
- [ ] Prove an unstaged rename commits and records the deletion.
- [ ] Confirm deletions, the governed refusal and existing drops are unchanged.

## Acceptance Criteria

- [ ] A Task renaming a file with `git mv` commits; today the commit fails
      because Git cannot match the source pathspec.
- [ ] A Task renaming a file without staging it commits, and the deletion is
      recorded.
- [ ] A Task deleting a file still stages the deletion.
- [ ] A governed rename without the required operation is still refused.
- [ ] The existing executable-mode drop and the existing real-repository
      staging assertion still pass unchanged.

## Context

- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q 'func TestTaskCommitStagesAStagedRename' internal/daemon/daemon_test.go && go test -count=1 ./internal/daemon -run '^TestTaskCommitStagesAStagedRename$'` — a `git mv` rename commits; this fails today.
- `grep -q 'func TestTaskCommitStagesAnUnstagedRename' internal/daemon/daemon_test.go && go test -count=1 ./internal/daemon -run '^TestTaskCommitStagesAnUnstagedRename$'` — an unstaged rename commits and records the deletion.
- `grep -q 'func TestFilterStageablePathsDropsPathAbsentFromWorktreeAndIndex' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestFilterStageablePathsDropsPathAbsentFromWorktreeAndIndex$'` — the drop is recorded with its reason.
- `grep -q 'func TestTaskCommitStagesAStagedRename' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestGovernedRenameRefusesFromRealSnapshot$'` — the governed rename refusal survives the staging change.
- `grep -q 'func TestTaskCommitStagesAStagedRename' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestGovernedMutationDetectionUsesTheUnfilteredSnapshot$'` — classification still reads the unfiltered snapshot.
- `grep -q 'func TestTaskCommitStagesAStagedRename' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo$'` — real-repository staging is unchanged.
- `grep -q 'func TestTaskCommitStagesAStagedRename' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission$'` — the existing drop reason is unchanged.

## References

- `_prd.md` → Goal 1; Core Features 1-3; Declared intentional breaks 1;
  Regression locks.
- `_techspec.md` → Implementation Design: Stage only what Git can match;
  Testing Approach observations 1-4; Build Order 1.

## Result

### Implementation

- The stageable-path filter now receives the caller's context and consults
  `git ls-files --error-unmatch` only after `Lstat` proves a path is absent
  from the worktree. Git exit code `1` records the path as dropped with reason
  `absent from worktree and index`; any other probe failure leaves absence
  unproven for the existing commit boundary to report.
- A path deleted only from the worktree remains stageable because it still
  matches the index. The changed-path reader and governed-mutation classifier
  were not changed.
- Real disposable-repository tests cover staged and unstaged renames. Existing
  fake-snapshot tests now materialize the paths they claim the Agent created;
  their assertions remain unchanged.

### Focused checks

- Before the production change,
  `GOCACHE=/private/tmp/roundfix-task-0136-go-cache rtk go test ./internal/daemon -run 'Test(TaskCommitStagesA(Staged|nUnstaged)Rename|FilterStageablePathsDropsPathAbsentFromWorktreeAndIndex)$'`
  exited `1`: the unstaged rename passed, the filter kept `before.txt`, and the
  staged rename failed at `git add -f -- after.txt before.txt`.
- After the production change, the same focused command exited `0` with three
  tests reported passed.
- `GIT_CONFIG_COUNT=0 GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null GOCACHE=/private/tmp/roundfix-task-0136-go-cache rtk go test ./internal/daemon -run 'Test(TaskCommitStagesAStagedRename|TaskCommitStagesAnUnstagedRename|FilterStageablePathsDropsPathAbsentFromWorktreeAndIndex|GovernedRenameRefusesFromRealSnapshot|GovernedMutationDetectionUsesTheUnfilteredSnapshot|SnapshotDiffCommitStagesOnlyAgentChangesInRealRepo|FilterStageablePathsDropsRegularFileWithAnyExecutePermission)$'`
  exited `0` with ten tests and subtests reported passed.
- `GIT_CONFIG_COUNT=0 GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null GOCACHE=/private/tmp/roundfix-task-0136-go-cache rtk go test ./internal/daemon`
  exited `0` with 311 tests reported passed.
- Two unchanged attempts with the default Go cache did not reach the tests:
  the sandbox returned `EPERM` under `~/Library/Caches/go-build`. The writable
  temporary cache above was used for the behavioral checks.
- `GOCACHE=/private/tmp/roundfix-task-0136-go-cache rtk make verify-incremental`
  exited `2` at `fmt-check` before tests because
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. Neither file is
  part of this Task's diff.

### Acceptance evidence

- `TestTaskCommitStagesAStagedRename` creates a repository under `t.TempDir`,
  performs `git mv`, commits through the filtered snapshot paths, independently
  reads the committed addition and deletion, and confirms the repository is
  clean.
- `TestTaskCommitStagesAnUnstagedRename` performs `os.Rename`, confirms neither
  path is dropped, commits both paths, independently reads the source deletion
  and destination addition with rename detection disabled, and confirms the
  repository is clean.
- `TestFilterStageablePathsDropsPathAbsentFromWorktreeAndIndex` records the
  staged rename source with the new reason, then resets the repository, deletes
  the tracked source only from the worktree, and confirms that deletion remains
  stageable.
- The combined regression check exercised
  `TestGovernedRenameRefusesFromRealSnapshot` and
  `TestGovernedMutationDetectionUsesTheUnfilteredSnapshot`, preserving refusal
  and classification from the unfiltered snapshots.
- The same check exercised the unchanged assertions in
  `TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission` and
  `TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo`.

### Follow-up

- The two pre-existing CLI test files reported by `fmt-check` need a separate,
  authorized formatting repair outside this Task.

## Carry-forward provenance

- Source Run: `run_20260914T134813Z_e15c7ee834fbe8bf`
- Source commit: `ae49dc594d8863a4bf3f22b6cbba5040e84e9514`
