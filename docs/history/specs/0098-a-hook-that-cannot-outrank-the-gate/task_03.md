---
status: completed
type: backend
---

# Task: Re-run Verification in Settle

Execute Task Verification in the selected settle surface.

## Work
- Reuse Verification logic from Implement
- Run commands verbatim (no edits)
- On pass: proceed to staging
- On fail: stop, print diagnostics

## Verification
- `grep -q "executeVerification\|runVerification" internal/cli/settle.go && grep -q "verificationCommand" internal/cli/settle.go`


## References

- Core Feature 2: Settle Recovery

## Result

Settle's Verification re-run is now a named, documented step of the settle
surface contract instead of an inline loop, and its failure reporting covers
every outcome the Implement Verification can produce.

### Implementation

- `internal/cli/settle.go:160` — `runVerificationInSettleSurface(ctx, plan,
  verifier, stdout, stderr)` re-runs `plan.task.Verification` in the resolved
  settle surface through the same `daemon.Verifier` collaborator the Implement
  Run engine uses (`collaborators.verifier`, resolved from
  `commandDependenciesForContext`). The loop variable is `verificationCommand`
  and it is passed to `daemon.VerifyRequest.Command` untouched; `WorkDir` is
  `plan.workDir` (the selected surface) and diagnostics land at
  `daemon.VerificationOutputPath(plan.artifactDir, settleVerificationRunID(plan), 1, 1)`.
- `internal/cli/settle.go:109` — `runSettleCommand` returns `exitRunFailed`
  the moment the re-run does not pass, before `settleTaskAndCommit`, so
  staging, the commit, and Run Branch integration are unreachable on failure.
- `internal/cli/settle.go:182` — new branch for
  `daemon.VerificationUnknownError`: a verdict the runner never observed is
  reported apart from a failed command (`verify <cmd> — verdict unknown`,
  `<task_id> stays <status> — verification verdict unknown`) with the cause on
  stderr. Previously this case exited 1 with only a bare stderr line and no
  settle report.
- `internal/cli/settle.go:41` — exit-code help now reads `verification did not
  pass` so it covers both non-passing verdicts.
- `docs/user-guide/commands.md:565` and `docs/user-guide/usage.md:506` — the
  settle contract documents verbatim, task-file-order execution with no Agent
  session, plus the failed and verdict-unknown output shapes.

### Acceptance criteria

1. **Verification runs commands verbatim (no code changes)** —
   `TestSettleVerificationRunsSurfaceCommandsVerbatim`
   (`internal/cli/settle_test.go:781`) settles a `completed`-but-uncommitted
   Task whose three commands include shell quoting and a pipe
   (`printf 'a|b' | grep -q 'a|b'`) and a `test "$(pwd -P)" = <repo>` check
   that only passes when the command runs in the selected surface. Asserted on
   exact stdout: the three `verify <command> — ok` lines echo the commands
   byte-for-byte, in task file order, before `commit done.txt` and the settled
   line. Nothing in settle rewrites, wraps, or reorders a command.
2. **Failure stops, leaves surface unchanged** —
   `TestSettleVerificationFailureKeepsHookRefusedWorkInTaskWorktree`
   (`internal/cli/settle_test.go:829`) reproduces the hook-refusal shape: a
   `completed` task file and staged `done.txt` in a kept Task Worktree, with a
   Verification whose second command fails. Exit code is 1; stdout is exactly
   the passing line, `verify test -f missing.txt — failed (diagnostics: <path>)`
   and `task_01 stays completed — verification failed`; the third command never
   ran (`should-not-run` absent); the Task Worktree and Run Worktree still
   exist; `git status --porcelain=v1` in the surface is byte-identical to
   before with `A  done.txt` still staged; HEAD count unchanged; the task file
   unchanged; the checkout still has no `done.txt`.
   `TestRunSettleVerificationFailureLeavesTaskAndTreeUntouched` continues to
   cover the same stop for a `failed` Task in the checkout.
   `TestSettleVerificationUnknownVerdictStopsWithoutCommitting`
   (`internal/cli/settle_test.go:904`) injects a verifier returning
   `daemon.VerificationUnknownError` and asserts exit 1, the two-line unknown
   report on stdout, the cause on stderr, no commit, and an unchanged tree.
3. **No Agent session, no repair prompt** — settle constructs no Agent runner;
   both new surface tests install an `implementFakeRunner` and assert
   `runner.calls == 0` after the pass and the failure path. There is no repair
   loop in `runVerificationInSettleSurface`: a temporary failure (exit 75)
   arrives wrapped around its command failure and stops the settle like any
   other failed command, which the doc comment records against ADR-0038.

### Focused checks

- `go build ./...` — pass.
- `go vet ./internal/cli/` — pass (no output).
- `gofmt -l internal/cli/settle.go internal/cli/settle_test.go` — empty.
- `go test ./internal/cli/ -run 'TestSettleVerification' -count=1` — `ok
  roundfix/internal/cli 1.122s` (the three new tests).
- `go test ./internal/cli/ -run 'Settle|settle' -count=1` — `ok
  roundfix/internal/cli 1.716s`.
- `go test ./internal/cli/ ./internal/daemon/ -count=1` — `ok` for both
  packages (52.0s, 2.6s).
- `go test -count=1 -parallel 16 ./...` — `ok` for all 26 packages with tests
  (internal/cli 116.0s, internal/baseline 81.3s, internal/speccheck 69.6s,
  internal/spec 60.9s, internal/store 52.3s, rest under 30s).
- Task Verification identifiers are in place without running the command:
  `internal/cli/settle.go:160` defines `runVerificationInSettleSurface` (matches
  `runVerification`) and line 162 declares `verificationCommand`.

### Out of scope / notes

- `make docs-test` fails on `TestCheckActiveCorpusHasNoErrors` and
  `TestCheckCorpusGolden` with `SC-VERIFY-NON-HERMETIC` pointing at
  `docs/specs/0113-a-refused-gate-writes-its-refusal-once/task_01.md`
  (`grep -q ... /tmp/qa-report*.md`). That file belongs to Spec 0113, is
  untouched by this Task (last changed in `0cdd4d97`), and the corpus golden
  expects `0` for that code — a pre-existing failure, not a regression from
  this diff.
- Staging the settled commit and its deletion handling stay with task_04 and
  task_05; this slice stops at "verification passed, proceed to staging".
