---
task: task_04
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: backend
complexity: medium
---

# Task 04: Carry-forward stages commits without repository hooks

## Overview

Carry-forward proves and applies a Run's settled Tasks in a fresh staging
worktree that `createCarryForwardStaging` in `internal/cli/carryforward.go`
creates with no installed dependencies. `stageCarryForwardCandidate` runs
`cherry-pick <commit>` and `commit --amend --no-edit` there, and
`abortCarryForwardStaging` runs `cherry-pick --abort`, all through
`reconcileGitRaw`. Git 2.54.0 runs `prepare-commit-msg` and `post-commit` for
the cherry-pick, and `pre-commit`, `prepare-commit-msg`, `commit-msg` and
`post-commit` for the amend. In Fluxus, Node-based hooks failed there with
`ERR_MODULE_NOT_FOUND`, reconcile printed the raw hook output, and completed
Tasks were redone. The hooks come from the repository's shared hooks directory
or its `core.hooksPath`. The carried commits already passed the Daemon's
Verification and those hooks when the Daemon settled them in their Run
Worktree, and the staging amend changes only the carried Task file's
provenance. Three commands reach this code: `roundfix reconcile
--carry-forward`, the implement Preflight inspection and
`roundfix deliver retry`. The checkout receives the staged result only by the
fast-forward merge in `applyCarryForwards`. The bypass must therefore stay
inside the staging invocations and never reach the user's Git configuration.

## Requirements

1. MUST run the staging worktree's `cherry-pick <commit>`, `cherry-pick
   --abort` and `commit --amend --no-edit` through one helper in
   `internal/cli/carryforward.go`. The helper adds `-c core.hooksPath=<dir>`
   to that invocation only, where `<dir>` is an empty directory Roundfix
   creates and removes, and it keeps the `-c core.fsmonitor=false -c
   commit.gpgSign=false` options `reconcileGitRaw` applies.
2. MUST NOT write any Git configuration file, and MUST keep `reconcileGitRaw`
   for the fast-forward merge in `applyCarryForwards` and for every read-only
   command.
3. MUST keep the function signatures of `createCarryForwardStaging` and
   `stageCarryForwardCandidate`, and keep the conflict, operational-failure and
   provenance-failure classification they return today.
4. MUST prove through `roundfix reconcile <run-id> --carry-forward`, run by the
   package's CLI runner, that a proved Task carries in a repository whose
   `pre-commit` and `commit-msg` hooks exit `1`. The tests MUST also prove that
   the carry-forward proof reports the candidate ready, that marker-writing
   `pre-commit`, `prepare-commit-msg`, `commit-msg` and `post-commit` hooks
   write no marker during staging, that the checkout's `core.hooksPath` is
   unchanged, and that the empty hooks directory is removed.
5. MUST write each hook fixture before the test starts any Git process and MUST
   NOT call `t.Parallel` in a test that writes one, so no concurrent fork can
   hold a hook open while Git executes it (ADR-0125).
6. MUST state in the `reconcile --carry-forward` section of
   `docs/user-guide/commands.md` and the carry-forward paragraph of
   `.agents/skills/roundfix/SKILL.md` that staging commits run without
   repository hooks and why, using the phrase `run without repository hooks`,
   then regenerate `skills/roundfix/SKILL.md` with `make skills-sync`.
7. MUST put the new tests in `internal/cli/carryforward_hooks_test.go`, each
   negative case a test of its own, and MUST NOT rename or remove an existing
   top-level test.

## Subtasks

- [ ] Route the three staging commands through the hook-disabling helper.
- [ ] Add the hook fixtures and the named tests.
- [ ] Document the staging rule in the guide and the skill, then run
      `make skills-sync`.

## Acceptance Criteria

- [ ] `roundfix reconcile <run-id> --carry-forward` carries a proved Task in a
      repository whose `pre-commit` and `commit-msg` hooks refuse.
- [ ] No repository hook writes its marker during staging.
- [ ] The checkout's `core.hooksPath` is unchanged and the empty hooks
      directory is gone after carry-forward.
- [ ] A conflicting settlement commit is still that Task's refusal, and an
      operational cherry-pick failure is still not a conflict.

## Context

- interface: `internal/cli/carryforward.go`
- creates: `internal/cli/carryforward_hooks_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCarryForwardAppliesThroughRefusingCommitHooks|TestCarryForwardProofIgnoresRefusingCommitHooks|TestCarryForwardStagingRunsNoRepositoryHook|TestCarryForwardLeavesTheCheckoutHooksPathUnchanged|TestCarryForwardRemovesItsEmptyHooksDirectory|TestCarryForwardConflictRefusesInsteadOfFailingTheInspection|TestCarryForwardStagingFailureIsNotClassifiedAsAConflict|TestCarryForwardOperationalCherryPickFailureIsNotAConflict|TestCarryForwardAcceptsBudgetExceededRun)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCarryForwardAppliesThroughRefusingCommitHooks TestCarryForwardProofIgnoresRefusingCommitHooks TestCarryForwardStagingRunsNoRepositoryHook TestCarryForwardLeavesTheCheckoutHooksPathUnchanged TestCarryForwardRemovesItsEmptyHooksDirectory TestCarryForwardConflictRefusesInsteadOfFailingTheInspection TestCarryForwardStagingFailureIsNotClassifiedAsAConflict TestCarryForwardOperationalCherryPickFailureIsNotAConflict TestCarryForwardAcceptsBudgetExceededRun; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "run without repository hooks" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "run without repository hooks" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide says staging commits run without repository hooks, so the command fails.

## References

- `_prd.md` → Goal 4; Core Feature 4; Success Metric 4; Decisions.
- `_techspec.md` → Carry-forward without repository hooks; API Contract 6;
  Testing Approach 4; ADR-0053; ADR-0125; ADR-0158.

## Result

### Implementation

- Added one staging-only Git helper that creates an empty hooks directory,
  adds `-c core.hooksPath=<dir>` to the existing `reconcileGitRaw` invocation,
  and removes the directory on return. The staging cherry-pick, abort, and
  provenance amend use it; read-only commands, `git add`, and the checkout's
  fast-forward merge keep their existing paths.
- Added CLI-runner coverage with hook fixtures written before any Git process.
  Separate tests cover refusing hooks, proof readiness, all four marker hooks,
  unchanged checkout configuration, temporary-directory cleanup, conflicts,
  provenance failures, and operational cherry-pick failures.
- Documented why carry-forward staging commits run without repository hooks in
  the user guide and canonical Roundfix skill. `make skills-sync` regenerated
  the shipped mirror, and `make baseline-digests` reported no derived changes.

### Focused checks

- Red signal: before the production change,
  `TestCarryForwardAppliesThroughRefusingCommitHooks` failed because
  `git commit --amend --no-edit` exited `1` through the repository's refusing
  hook.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run
  '^(TestCarryForwardAppliesThroughRefusingCommitHooks|TestCarryForwardProofIgnoresRefusingCommitHooks|TestCarryForwardStagingRunsNoRepositoryHook|TestCarryForwardLeavesTheCheckoutHooksPathUnchanged|TestCarryForwardRemovesItsEmptyHooksDirectory)$'
  ./internal/cli` exited `0`.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run
  '^(TestCarryForwardConflictRefusesInsteadOfFailingTheInspection|TestCarryForwardStagingFailureIsNotClassifiedAsAConflict|TestCarryForwardOperationalCherryPickFailureIsNotAConflict)$'
  ./internal/cli` exited `0`.
- `rtk rg -n -F 'run without repository hooks'
  docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` and `rtk cmp .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` exited `0`.
- The first sandboxed `make verify-incremental` attempt was environment-blocked
  when a test tried to reach `cafe.github.com`. The network-capable rerun,
  `rtk env GOCACHE=/tmp/roundfix-task04-gocache make verify-incremental`,
  exited `0`, including repository-wide vet, tests, skill checks, and build.
- The Task's declared `## Verification` command was not run; the Daemon owns
  that gate.

### Acceptance evidence

- Refusing hooks: `TestCarryForwardAppliesThroughRefusingCommitHooks` exercised
  `roundfix reconcile <run-id> --carry-forward` through the package CLI runner
  and observed the proved Task in the checkout despite refusing `pre-commit`
  and `commit-msg` hooks. `TestCarryForwardProofIgnoresRefusingCommitHooks`
  observed the candidate action `would carry forward with --carry-forward`.
- No staging hooks: `TestCarryForwardStagingRunsNoRepositoryHook` first proved
  that `pre-commit`, `prepare-commit-msg`, `commit-msg`, and `post-commit` each
  wrote its marker for a normal commit, cleared the marker, then observed no
  marker from carry-forward staging.
- Configuration and cleanup:
  `TestCarryForwardLeavesTheCheckoutHooksPathUnchanged` observed the same local
  `core.hooksPath` before and after carry-forward, and
  `TestCarryForwardRemovesItsEmptyHooksDirectory` observed an empty temporary
  root after the command returned.
- Error classification: the three focused negative tests preserved conflict as
  a Task refusal and kept provenance and missing-object cherry-pick failures as
  operational errors rather than conflicts.
