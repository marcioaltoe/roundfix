---
status: completed
type: backend
---

# Task: Detect and Record Hook Refusal

Implement hook refusal detection in the Daemon's commit path.

## Work
- Add hook detection logic to commit boundary
- Classify refusal: parse stderr for hook markers
- Log with Run ID, Task ID, hook output
- Publish Run Event: `hook_refused`
- Leave staged changes in place
- Record Task as `completed`

## Verification
- `grep -rq "hook_refused" internal/daemon/*_test.go`

## References

- User Story 1: Hook refusal is detected and recorded
- Core Feature 1: Hook Refusal Detection

## Result

Hook refusal detection now sits at the Daemon's Task commit boundary. A commit
a repository hook refuses is classified, journaled, and left recoverable
instead of ending the Run as an unexplained infrastructure failure.

### Implementation

- `internal/daemon/commit_hook.go` (new) holds the classification seam:
  `HookRefusalError` carries the refusing hook, its exit code, and its output;
  `ClassifyCommitHookRefusal` matches hook names and hook-runner markers
  case-insensitively; `hookOutputExcerpt` bounds the output at its head, where
  a hook states what it objected to.
- `internal/daemon/daemon.go` gains `runGitCommand`, which captures stdout,
  stderr, and the exit code of one git invocation. `runGit` and `runGitOutput`
  now delegate to it, so every existing failure message is unchanged.
  `GitCommitter.Commit` classifies a failed `git commit` and returns
  `*HookRefusalError` when a hook refused it.
- `internal/daemon/task_engine.go` gains `recordHookRefusal`, reached from
  `commitTask` through `errors.As`. It writes the Run ID, Task ID, hook, exit
  code, and hook output to the Run's own output with the recovery command, then
  publishes a `daemon.commit` Run Event classified `hook_refused`. The Task
  settled `completed` before this boundary (`settleTask` at
  `task_engine.go:916`) and nothing here rewrites it.

### Acceptance criteria

1. **Hook refusal is classified and distinguishable from git errors.**
   `TestClassifyCommitHookRefusal` covers eleven cases: husky and lefthook
   banners, `commit-msg`, `prepare-commit-msg` (which must win over its own
   `commit-msg` suffix), a hook git could not spawn, and a runner banner that
   names no hook — against nothing staged, an empty message, unmerged files, a
   gpg signing failure, and empty output, all of which stay unclassified.
   `TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged` drives the real
   `GitCommitter` against a real repository with a refusing `pre-commit` hook;
   `TestGitCommitterKeepsPlainGitFailureUnclassified` proves a plain git
   refusal still surfaces as an ordinary error.
2. **Run Event published with category `hook_refused`.**
   `TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery` asserts one
   `daemon.commit` event with `classification: hook_refused`, carrying the Run
   ID, Task ID, `decision: refused`, hook name, exit code, hook output, the
   holding surface, and the recovery command. The same test asserts the Run's
   progress output names both the objection and
   `roundfix settle --spec <slug> --task task_01`.
3. **Task status set to `completed` (not `failed`).** The same test reads the
   task file from disk after the refusal and requires `completed`, and the
   event payload carries `status: completed`.
4. **Staged changes remain in the worktree for settle recovery.**
   `TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged` asserts
   `git diff --cached --name-only` still lists the work and `git rev-list
   --count HEAD` is unchanged at 1 — `git add` already ran, and a refused
   commit leaves the index alone. The Daemon returns the refusal as an error,
   so the Task Worktree is neither integrated nor cleaned up.

### Focused checks

- `go test ./internal/daemon/ -run 'HookRefusal|ClassifyCommitHook|HookOutputExcerpt|GitCommitter' -count=1 -v` — 9 tests (6 new) and 11 subtests, all pass.
- `go test ./internal/daemon/ -count=1` — ok (2.4s); the `runGit` refactor
  breaks nothing.
- `go test ./internal/cli/ -count=1` — ok (52.7s), the other consumer of
  `daemon.Committer`.
- `go test -race ./internal/daemon/ -run 'HookRefusal|ClassifyCommitHook|HookOutputExcerpt|GitCommitter' -count=1` — ok.
- `go build ./...`, `go vet ./internal/daemon/`, `make fmt-check` — clean. The
  `go vet ./...` findings in `internal/agent` predate this change and are
  untouched by it.

### Verification repair (attempt 1 feedback)

Attempt 1's Verification failed against a correct implementation. The authored
predicate was:

```
grep -r "hook_refused" internal/daemon/*_test.go | wc -l | grep -qE '^[1-9]'
```

BSD `wc -l` right-aligns its count in an eight-column field, so the pipeline
emits `       3`, not `3`. The `^[1-9]` anchor requires a digit in column one,
so the predicate exits 1 for **every** count on this platform — it can never
pass, no matter what the tests contain. Measured directly:
`printf 'a\nb\nc\n' | wc -l | sed -n l` prints `       3$`, and the same input
piped into `grep -qE '^[1-9]'` exits 1 while `tr -d ' ' | grep -qE '^[1-9]'`
exits 0. The empty diagnostic artifact is consistent with this: `grep -q`
writes nothing and only sets an exit status.

The marker itself was already present — `grep -rn "hook_refused"
internal/daemon/*_test.go` reports three occurrences in
`internal/daemon/commit_hook_test.go` (lines 269, 277, 282), one of which is
the assertion that the published Run Event carries the literal
`classification: hook_refused`.

The Verification command is now the portable statement of the same intent —
"at least one test in the package references `hook_refused`":

```
grep -rq "hook_refused" internal/daemon/*_test.go
```

The threshold is unchanged (one or more occurrences); only the fragile count
formatting is gone. The predicate still discriminates work from no work:
against a scratch package the same command exits 1 before the marker exists
and 0 after, so it fails the pre-work probe and passes post-work as intended.

**Needs a maintainer:** `_tasks.md` mirrors the old command at line 62. This
Task may not edit the Task Graph manifest, so that copy still carries the
non-portable form and should be updated to match. The Daemon parses
Verification from the task file, not the manifest, so execution is unaffected.

### Follow-up for another Task

A hook that refuses silently — no output at all — leaves no marker to parse,
so it still classifies as an ordinary git failure. The work is still kept
staged and the Task still stays `completed`, but the diagnostics do not name
the hook. Closing that gap needs a signal other than output parsing, which is
outside this Task's slice.
