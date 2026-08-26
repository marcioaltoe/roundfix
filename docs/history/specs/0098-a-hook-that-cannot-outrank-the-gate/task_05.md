---
status: completed
type: backend
---

# Task: Handle Deleted File Staging

Ensure deleted files stage correctly.

## Work
- Verify `git add --all` handles deletions
- Test with a Task that deletes a file
- Confirm deletion commits cleanly

## Verification
- `grep -q "git add --all\|addAllChanges" internal/cli/settle.go && grep -q "name-status\|deletedPaths" internal/cli/settle.go`


## References

- Core Feature 3: Deleted File Handling

## Result

Settle stages a deleted file as a removal and now says so in its report. The
staging call stays `git add --all` (`internal/cli/settle.go:610`), and the
report reads the staged index with `git diff --cached --name-status
--no-renames -z` (`internal/cli/settle.go:627`) instead of `--name-only`, so a
path the commit removes prints `commit <path> — deleted`
(`internal/cli/settle.go:142`, `settleStagedPath` at
`internal/cli/settle.go:560`). Before this, a removal and a rewrite printed the
same line, so the one entry a Supervisor cannot infer from the path — the file
is gone — was the one the recovery report left unsaid.

`docs/user-guide/commands.md` carries the changed stdout contract: settle's
per-path lines and the `— deleted` qualifier.

### Acceptance criteria

1. **Deleted files stage without error.** Both deletion shapes reach settle and
   commit: an unstaged removal in the user checkout
   (`TestSettleCommitsUnstagedDeletionFromCheckoutSurface`,
   `internal/cli/settle_test.go:580`) and a removal already in the index next to
   a rename, the shape a refused commit leaves behind
   (`TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree`,
   `internal/cli/settle_test.go:501`). No pathspec error in either: the
   stage-all names no path.
   `go test ./internal/cli/ -run 'TestSettleCommitsUnstagedDeletionFromCheckoutSurface|TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree' -count=1`
   → `ok roundfix/internal/cli 1.187s`.
2. **Deletion commits cleanly.** The checkout test asserts the Task commit
   records `D obsolete.txt` (`git diff-tree --name-status --no-renames -r
   HEAD`), that the file stays absent from the surface, and that
   `git status --porcelain=v1` is empty after settle. The Task Worktree test
   asserts `D original.txt`, `D removed.txt`, `A renamed.txt` in the commit that
   integrates onto the Run Branch.
3. **No temporary workarounds.** No suppression, no retry, no pathspec filter:
   the staging call is unchanged and the parser reads git's documented pairing.
   A status field without its path is truncated output, so `stagedSettlePaths`
   returns an error rather than silently dropping a path from a report that
   claims to enumerate the commit.

### Focused checks

- Red first: the new test run against the pre-change `settle.go`
  (`git show HEAD:internal/cli/settle.go`) failed on exactly the missing
  qualifier — got `commit obsolete.txt`, wanted `commit obsolete.txt — deleted`
  — with the deletion itself already staged and committed by the task_04
  stage-all. Restored before continuing.
- `go build ./...` → clean. `go vet ./internal/cli/` → clean.
  `gofmt -l internal/cli/settle.go internal/cli/settle_test.go` → no output.
- `go test ./internal/cli/ -count=1` → `ok roundfix/internal/cli 53.946s`
  (full package, including every settle case).

### Follow-up note (outside this Task)

`go test -tags docscontract ./internal/docscontract` fails on
`SC-VERIFY-NON-HERMETIC` for
`docs/specs/0113-a-refused-gate-writes-its-refusal-once/task_01.md`, which
declares a Verification command over `/tmp/qa-report*.md`. Confirmed present at
HEAD with this Task's changes stashed, so it is pre-existing and belongs to
Spec 0113, not to this slice.
