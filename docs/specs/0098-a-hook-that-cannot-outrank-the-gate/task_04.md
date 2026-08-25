---
status: completed
type: backend
---

# Task: Integrate Settled Commits

Commit and integrate settled work onto Run Branch.

## Work
- Stage changes with `git add --all`
- Create standard Task commit
- Integrate onto Run Branch
- Remove Task Worktree on success

## Verification
- `grep -q "stageSharableChanges\|git add\|addAllChanges" internal/cli/settle.go && grep -q "integrateTask\|runGit.*commit" internal/cli/settle.go`


## References

- Core Feature 2: Settle Recovery

## Result

Settle now stages the whole recovered surface instead of re-staging a path list
it rebuilds from `git status`, and the standard Task commit that follows carries
what a refused commit already staged there.

### What changed

`settleTaskAndCommit` wrote the Task's settled status, took a porcelain snapshot
of the surface, and handed those paths back to the committer as an explicit
`git add -f -- <paths>` list. That list cannot describe the surface a refused
commit leaves behind. The Daemon had already staged the Task before the hook
refused, so a file the Task deleted is in neither the index nor the worktree —
re-staging it by pathspec matches nothing and the whole settle dies after a
passing Verification. The same snapshot names only the new side of a rename, so
the old side's deletion never reaches the staging call at all.

Settle now flips the Task status first, stages the surface with `git add --all`,
and reads the index back to learn exactly which paths the commit will carry.
The commit itself still goes through the Daemon's `Committer`, so a hook that
refuses settle's own commit is still classified rather than reported as a bare
git error. Integration onto the Run Branch and the Task Worktree removal are
unchanged in behavior; `integrateSettledTaskWorktree` is now
`integrateTaskWorktree` and carries a doc comment stating that a conflict
returns before the removal.

The two helpers that existed only to patch the task file into the snapshot list
(`ensureSettleCommitPath`, `settleArtifactCommitPath`) are gone — `git add --all`
picks the status flip up on its own.

### Evidence per acceptance criterion

**`git add --all` stages all changes (handles deletions)** — new test
`TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree` builds the exact shape a
hook refusal leaves: a Task Worktree where a tracked file was deleted, another
renamed, and everything staged with `git add -A` before the commit was refused.
Run against the pre-change `settle.go`, it fails with the loss this Spec exists
to stop:

```
roundfix: settle failed after verification:
git add -f -- docs/specs/0001-widget-flow/task_01.md removed.txt renamed.txt failed:
exit status 128: fatal: pathspec 'removed.txt' did not match any files
```

The same run also shows the rename's old side (`original.txt`) missing from the
path list entirely. With the change, `go test ./internal/cli/ -run
TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree -count=1` passes, and the
test asserts the commit's own name-status records `D original.txt`,
`D removed.txt`, and `A renamed.txt`.

**Standard Task commit created** — the commit is still built from
`daemon.TaskCommitMessage(slug, task)` through `collaborators.committer`.
`TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage` compares the settle
commit's message against `daemon.TaskCommitMessage` and passes unchanged.

**Integrated onto Run Branch (or updates checkout)** — the new test asserts the
deletion and the rename reach the user checkout: `removed.txt` and
`original.txt` are gone from the checkout, `renamed.txt` holds the moved
content, and `git status --porcelain=v1` is empty afterwards.
`TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration` and
`TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal` pass
unchanged, and `TestRunSettleTaskWorktreeIntegrationConflictKeepsSurfaces` still
shows a conflict leaving the Run Branch tip and both surfaces untouched.

**Task Worktree removed on success** — the new test asserts the Task Worktree
path and its branch are both gone after settle, alongside the Run Worktree, and
that no Active Run is left in the git root.

### Focused checks run

- `go build ./...` — clean.
- `go vet ./internal/cli/` — clean. (`go vet ./...` reports pre-existing
  `passes lock by value` findings in `internal/agent`, untouched here.)
- `go test ./internal/cli/ ./internal/daemon/ ./internal/worktree/ -count=1` —
  all ok (cli 52.9s, daemon 4.5s, worktree 5.7s), 20 pre-existing settle tests
  plus the new one.
- `gofmt -l internal/cli/` — empty.

### Follow-ups

- `filepathInRoot` in `internal/cli/settle.go` has no callers and predates this
  Task; left in place rather than widened into this slice.
- Settle's staging no longer runs paths through `daemon.FilterStageablePaths`,
  so the executable-file and symlink-crossing drops that filter applies to the
  Daemon's own Task commit do not apply to a settle commit. Settle's contract is
  to commit the surface whole, and `git add --all` still honours `.gitignore`,
  but if that filter should also bound settle it is a decision for the Spec
  rather than this Task.
