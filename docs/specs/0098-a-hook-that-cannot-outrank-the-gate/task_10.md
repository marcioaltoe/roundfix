---
status: pending
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
A hook that prints only its finding is classified as `hook_refused`, so the Run
publishes the Run Event and names `roundfix settle` in the environment where
the three cases were measured.
