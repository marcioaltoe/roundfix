---
task: task_03
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
status: completed
type: backend
complexity: medium
---

# Task 03: Every Roundfix push refuses a history that adds a secret

## Overview

GitHub push protection refused the Delivery Queue's push of Spec 0259 because
a commit on the item branch carried an API key; the queue learned it only at
the remote. This Task makes both pushes Roundfix makes, the Delivery Queue's
`GitHubCLI.PushBranch` and the Daemon's Final Push `GitPusher.Push`, read every
commit the push would send and refuse on a Secret Match before `git push`
runs, and parks a Delivery Queue item as `secret-found`. The hazard: the
commits come from Agents and QA in the item's Runs, and the destination is the
GitHub remote, which publishes them. It is the Integration Checkpoint for
task_01's history API, proving its log shape against real Git.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [QA evidence records environment values, secrets included](references/2026-10-13-qa-evidence-records-environment-values.md),
   for its push half: nothing between the QA commit and the push scanned it.
2. MUST make `GitHubCLI.PushBranch` and `GitPusher.Push` follow
   `_techspec.md` → API Contracts 3 and 4: run `git` with
   `secretscan.PushedHistoryArgs` through their own runner, pass the output to
   `ScanPushedHistory`, and on a Match or a failed read return the wrapped
   `*secretscan.MatchesError` or error without running `git push`.
3. MUST add `BlockerSecretFound = "secret-found"` to
   `internal/delivery/engine.go`, make `Engine.push` park the item with
   `secret-found: <match>, <match>, …` on a `*secretscan.MatchesError`, and
   classify it in `ClassifyPark` as `finding` with API Contract 4's next
   action. `TestEveryBlockerHasAParkClass` must pass unchanged.
4. MUST change, in `internal/delivery/github_test.go`, only the scripts of
   `TestPushPublishesTheRecordedHeadNotTheCheckout` and
   `TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead`, which gain
   the history-read step before `git push`; and add the `secret-found` row to
   `TestParkClassesAndNextActions` in `internal/delivery/park_class_test.go`,
   asserting its class and its next action.
5. MUST add `internal/delivery/push_secret_scan_test.go` with
   `TestPushRefusesAHistoryThatAddsASecret`, through `NewGitHubCLI` over a
   temporary repository with a bare remote: one commit adds a runtime-built
   `github-token` sample on line 3 of a file and the next commit removes it.
   It requires the `MatchesError` to name the first commit, the path, line 3
   and `github-token`, the remote branch to be absent, the sample absent from
   the error, and a clean history to push. Add `TestEnginePushParksSecretFound`
   through the engine's existing fake boundary, requiring the park blocker and
   no Pull Request creation.
6. MUST add `internal/daemon/push_secret_scan_test.go` with
   `TestFinalPushRefusesAHistoryThatAddsASecret`, the same history through
   `GitPusher.Push`, requiring the refusal, the absent remote branch and a
   clean push afterwards. Every new test follows ADR-0259 and ADR-0271 and
   builds its sample at runtime.
7. MUST NOT change `internal/cli/review.go`, the push arguments, the Pull
   Request step, or any other blocker's class or next action.

## Subtasks

- [ ] Scan the pushed history in both push sites before `git push`.
- [ ] Add the `secret-found` blocker, its park and its class.
- [ ] Update the two scripted push tests and the class table.
- [ ] Add the real-Git refusal tests for both pushes.

## Acceptance Criteria

- [ ] A history that adds and then removes a sample is refused in both push
      paths, and the remote branch stays absent.
- [ ] A refused Delivery Queue push parks `secret-found` with class `finding`
      and its next action, and opens no Pull Request.
- [ ] No error, blocker or test output carries the sample.

## Consumers

- `GitHubCLI.PushBranch`, `GitPusher.Push` and the `secret-found` blocker — consumers: `./internal/delivery`, `./internal/daemon`

## Context

- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `docs/adr/0138-spec-run-git-policy-commit-per-task-push-only-at-clean.md`
- interface: `internal/delivery/github.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/github_test.go`
- interface: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/push_secret_scan_test.go`
- interface: `internal/daemon/daemon.go`
- creates: `internal/daemon/push_secret_scan_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPushRefusesAHistoryThatAddsASecret|TestEnginePushParksSecretFound|TestPushPublishesTheRecordedHeadNotTheCheckout|TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions|TestFinalPushRefusesAHistoryThatAddsASecret)$' ./internal/delivery ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPushRefusesAHistoryThatAddsASecret TestEnginePushParksSecretFound TestPushPublishesTheRecordedHeadNotTheCheckout TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead TestEveryBlockerHasAParkClass TestParkClassesAndNextActions TestFinalPushRefusesAHistoryThatAddsASecret; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q 'BlockerSecretFound' internal/delivery/park_class.go || { printf 'secret-found has no park class\n' >&2; exit 1; }` — expected: exit 0; before this Task the three new tests and the blocker do not exist.

## References

- [_prd.md](_prd.md) — Goals; User Story 3; Core Feature 4; Success Metric 3
- [_techspec.md](_techspec.md) — API Contract 3; API Contract 4; Invariants 1-3; Testing Approach; Build Order 3
- [2026-10-13-qa-evidence-records-environment-values.md](references/2026-10-13-qa-evidence-records-environment-values.md)
- ADR-0273; ADR-0138; ADR-0259; ADR-0271


## Result

### Implementation

Both push paths read `secretscan.PushedHistoryArgs` through their existing
runner before publishing. Delivery scans the recorded candidate head; Final
Push scans `HEAD`. A failed read returns an operation-wrapped error, and a
Secret Match returns the specified wrapped `*secretscan.MatchesError`.
Delivery does not include a failed read's partial stdout patch in its error.
The existing push arguments and Pull Request operations are unchanged.

`Engine.push` parks a wrapped Secret Match refusal as
`secret-found: <match>, <match>, …`, using only `Match.String()` metadata.
Publication stops at that park before Pull Request creation. The new blocker
has Park Class `finding` and the specified instruction to remove the secret
from every named commit in the item worktree before `roundfix deliver retry`.

The two existing scripted push tests gain only the history-read step.
The park table gains the new class and an exact next-action assertion;
`TestEveryBlockerHasAParkClass` remains unchanged. New samples are built at
runtime, and real Git repositories and bare remotes live under `t.TempDir()`.
The delivery real-Git test names its process-wide configuration changes as
its sequential reason. Both Final Push tests call `t.Parallel()` first and
re-execute only themselves in a child with `gittest.IsolatedEnv()` and a
per-test `GIT_TRACE` path; neither changes the parent process environment.
Other new tests call `t.Parallel()`. No new test uses a real-time budget.

### Focused checks

- Before implementation:
  `go test -count=1 -run '^(TestPushRefusesAHistoryThatAddsASecret|TestEnginePushParksSecretFound|TestFinalPushRefusesAHistoryThatAddsASecret)$' ./internal/delivery ./internal/daemon`
  exited 1. Both real push tests reported that no `MatchesError` was returned,
  and the engine test reported the absent safe `secret-found` park.
- After implementation:
  `go test -count=1 -v -run '^(TestPushRefusesAHistoryThatAddsASecret|TestEnginePushParksSecretFound|TestFinalPushRefusesAHistoryThatAddsASecret|TestPushRefusesAFailedHistoryRead|TestFinalPushRefusesAFailedHistoryRead)$' ./internal/delivery ./internal/daemon`
  exited 0, with all five tests and the three failed-read subtests passing.
- Existing contract checks:
  `go test -count=1 -run '^(TestPushPublishesTheRecordedHeadNotTheCheckout|TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions|TestExistingGenericParkAnswersAreUnchanged|TestUnmatchedPushIntentRetriesWhenTheRemoteDiffers|TestDeliveryEngineResumesWithoutDoubleEffects)$' ./internal/delivery`
  exited 0.
- `git diff --check` exited 0.

### Acceptance evidence

1. Both real push tests add a runtime-built `github-token` sample on line 3,
   commit it, then commit its removal. Each refusal identifies the adding
   commit, path, line 3 and kind through `errors.As`; the bare remote has no
   branch after refusal. Both paths subsequently publish a clean history and
   verify its remote head. Delivery also proves it scans the supplied clean
   head while the checkout remains on the refused history.
2. `TestEnginePushParksSecretFound` checks the persisted parked stage, two
   comma-separated safe Matches, class `finding`, the exact next action and
   no Pull Request or merge effect. The existing park sweep and class table
   passed with the new blocker.
3. Both real push tests and the engine test reject every eight-byte sample
   fragment in errors and blockers without printing the sample on assertion
   failure. The failed-read delivery test also supplies sample-bearing partial
   stdout and requires a safe refusal with no push invocation. Focused test
   output contains test names, pass verdicts and suite-guard cost only; no
   sample is recorded here or in repository fixtures.

Additional refusal evidence: the delivery runner's original error remains
reachable through `errors.Is`; nonzero history reads run no push. The Final
Push failed-read test uses an unborn local history and a temporary `GIT_TRACE`
file to prove `git log` ran and `git push` did not, with no remote branch.

The authored Verification and repository gates remain for the Daemon to run.
Task status, checkboxes and the Task Graph manifest were not edited by this
Agent. No follow-up outside this slice was identified.


### Verification feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261014T102741Z_dd0d6c13d87ea4f9/verification/batch-003-attempt-1.log`.
The configured gate rejected the two new sequential Final Push tests because
`internal/daemon` exceeded its fixed sequential-test ceiling. A focused run
of the repository parallel-test contract reproduced that refusal before the
repair.

The repair changes only this Task's `internal/daemon/push_secret_scan_test.go`
and this Result. Both tests now start with `t.Parallel()` and invoke the
existing compiled test binary with an exact test selector and isolated Git
environment. The child runs the original real-Git assertions, including the
added-then-removed sample refusal, absent remote branch, clean push, safe
error metadata and the failed-read trace proving no push invocation. The
parent waits for the child through its test context and captures its result.
No sequential ceiling, Verification configuration, production behavior or
other Task file changed.

Fresh focused checks after the repair:

- `go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel/repository$' ./internal/testfixture`
  exited 0; the daemon has 15 sequential tests against its ceiling of 15.
- `go test -count=1 -race -run '^(TestFinalPushRefusesAHistoryThatAddsASecret|TestFinalPushRefusesAFailedHistoryRead)$' ./internal/daemon`
  exited 0, exercising both isolated child tests under the race detector.
- `go test -count=1 -shuffle=on -run '^(TestFinalPushRefusesAHistoryThatAddsASecret|TestFinalPushRefusesAFailedHistoryRead)$' ./internal/daemon`
  exited 0, exercising the same acceptance assertions with shuffled ordering.
- `git diff --check` exited 0.

The Delivery Queue implementation and its recorded acceptance evidence are
unchanged by this repair. The Daemon retains Task status and the full
configured Verification rerun; no authored Verification command or repository
gate was rerun by this Agent.
