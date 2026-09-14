---
task: task_01
spec: 0136-a-rename-the-committer-can-stage
status: pending
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
