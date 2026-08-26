---
status: completed
type: backend
---

# Task: Classify a hook that does not name itself

QA finding F-2 (`Trust-Damage`): `ClassifyCommitHookRefusal`
(`internal/daemon/commit_hook.go`) decides only from the hook's own output,
matching hook names or a runner banner. Git prints nothing identifying when a
hook fails, so a hook that emits only its finding is not classified — and that
is the shape of the repository the Spec's three cases were measured in
(`.husky/pre-commit`, plain `set -eu` into lint-staged, no banner). The Run
then ends without the `hook_refused` record and without naming the recovery,
which is the silent failure the source finding asked to eliminate.

## Work

- Give the classifier a second, structural signal: when a commit fails after
  the authoritative Verification passed, resolve the repository's hooks
  directory with `git rev-parse --path-format=absolute --git-path hooks`
  (honouring `core.hooksPath`, and correct from inside a worktree or a
  subdirectory) and check whether an executable `pre-commit` or `commit-msg`
  exists there. An existing executable hook plus a failed commit is a hook
  refusal regardless of what the output names.
- Keep the existing output-based matching as the first signal — when the hook
  names itself the diagnostics stay as they are today, including the runner
  banner in the reported objection.
- A repository with no executable hook keeps returning the raw git error, so a
  genuine git failure (index lock, disk full) is not relabelled as a hook
  refusal in a repository that has no hook at all. Record in the doc comment
  that a genuine git failure in a repository that *does* have a hook is the
  accepted false-positive boundary of this decision, and that the recovery path
  it names is harmless in that case because settle re-runs the Verification
  before committing.
- Both commit boundaries reach the same classifier: the Daemon's Task commit
  and settle's own commit.

## References

- User Story 1: Run meets a refusing hook and keeps verified work
- Core Feature 1: Hook Refusal Detection

## Verification
- `grep -q "path-format=absolute" internal/daemon/commit_hook.go && grep -q "core.hooksPath\|hooksPath" internal/daemon/commit_hook.go && go test -count=1 ./internal/daemon -run 'TestClassifyCommitHookRefusal|TestHookRefusalRecovery' 2>&1 | grep -q "^ok"`

## Result

A hook that prints only its finding is now classified as `hook_refused`, so the
Run publishes the Run Event and names `roundfix settle` in the environment where
the three cases were measured.

### The failure this reproduces first

The three measured cases were fixtured with a runner banner (`husky - pre-commit
script failed (code 1)`) that git never prints, which is why the suite stayed
green while the real Runs died. Restoring the measured shape — a plain `set -eu`
hook that prints its finding and no banner — reproduced F-2 exactly against the
old output-only classifier:

```
go test -count=1 ./internal/daemon -run 'TestHookRefusalRecovery'
--- FAIL: TestHookRefusalRecovery/function_over_the_80-line_limit
    expected the hook refusal classified, got git commit ... failed: exit status 1:
    src/parser.go:3: function exceeds the 80-line limit
--- FAIL: TestHookRefusalRecovery/generated_file_over_the_500-line_limit
--- FAIL: TestHookRefusalRecovery/sort()_where_the_rule_requires_toSorted()
```

All three ended with a raw git error: no `hook_refused` record, no recovery
named. That is the silent failure the source finding asked to eliminate.

### Evidence per work item

**A second, structural signal.** `ClassifyCommitHookRefusal` now takes the
commit's working directory and, when the output names nothing, asks git for the
hooks directory with `rev-parse --path-format=absolute --git-path hooks` and
checks for an executable `pre-commit` or `commit-msg` there. Resolution is git's
own, so it honours `core.hooksPath`, a subdirectory, and a linked worktree.
Covered by `TestClassifyCommitHookRefusalReadsTheRepositoryWhenOutputNamesNothing`
(executable `pre-commit` and `commit-msg` classify; a non-executable hook, a
`.sample` hook, and a `pre-push` hook do not) and by
`TestClassifyCommitHookRefusalResolvesTheHooksDirectoryGitWouldRun`, whose three
subtests pin the husky layout (`core.hooksPath = .husky`), a call from a
subdirectory, and a call from a Task Worktree — the last two would fail against a
guessed `<workDir>/.git/hooks`.

**Output matching stays the first signal.** A hook that names itself keeps
today's diagnostics: `TestClassifyCommitHookRefusal` runs its whole table against
a repository with no hook installed, so only the output can decide, and every
banner and hook-name case still returns the hook it named.
`TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged` and
`TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery` keep the named-hook
path covered end to end, including the runner banner in the reported objection.

**The hookless repository keeps the raw git error.** Asserted by
`TestGitCommitterKeepsPlainGitFailureUnclassified` (nothing staged, real
committer) and by the `a repository with no hook keeps the raw git error` case.
The accepted false-positive boundary — a genuine git failure in a repository that
*does* have a hook — is recorded in the `ClassifyCommitHookRefusal` doc comment,
together with why the recovery it names is harmless there: settle re-runs the
authoritative Verification before it commits.

**Both commit boundaries reach the same classifier.** Classification happens
inside `GitCommitter.Commit`, and `daemon.GitCommitter{}` is the single committer
behind the Daemon's Task commit (`internal/daemon/task_engine.go:1603`) and
settle's own commit (`internal/cli/settle.go:585`, wired at
`internal/cli/cli.go:4293`). The settle boundary is now exercised on the measured
shape too: `TestSettleRecoversMeasuredHookRefusedWork`'s hook lost its banner, and
the test bites — with the structural signal disabled it fails with settle
reporting a raw `git commit ... failed: exit status 1` and no refusal at all.

### Focused checks

| Check | Outcome |
| --- | --- |
| `go build ./...` | pass |
| `make fmt-check` | pass |
| `go vet ./internal/daemon ./internal/cli` | pass |
| `go test -count=1 ./internal/daemon` | `ok roundfix/internal/daemon 3.099s` |
| `go test -count=1 ./internal/cli` | `ok roundfix/internal/cli 58.388s` |
| `go test -count=1 ./internal/daemon -run 'TestHookRefusalRecovery' -v` | all three measured cases PASS |
| `go test -count=1 ./internal/cli -run 'TestSettleRecoversMeasuredHookRefusedWork'` | pass |

Pre-existing `go vet ./...` findings in `internal/agent/acpx_runner.go` (lock
passed by value) are untouched by this Task.

### Notes

- Changes land in `internal/daemon/commit_hook.go`, `internal/daemon/daemon.go`,
  their tests, and `internal/cli/settle_test.go`. The last is a fixture-fidelity
  fix only — no production change in `internal/cli` — made because that test
  carried the same banner git never prints, and without it the second commit
  boundary would still be verified against a shape the measured repository never
  produced.
- The refusal of an unnamed hook reports the existing generic label: `Hook` is
  empty, `HookName()` returns `commit`, and the Run Event payload carries
  `"hook": "commit"` with the recovery command unchanged.
