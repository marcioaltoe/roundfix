---
task: task_02
spec: 0069-review-run-targets-its-pull-request
status: completed
type: backend
complexity: high
---

# Task 02: Keep the target while the Run writes

## Overview

task_01 stops a Run that starts on the wrong branch. This slice stops one that
*becomes* wrong: the checkout moves while the Run is Active, and every write
after that lands somewhere the Pull Request never named.

The PRD assumed a late check already existed to be moved. Authoring disproved
it — `checkout branch mismatch` is not in the tree, and the nearest live code
is Round artifact reuse in `rounds.go`. This slice builds the guard.

The distinction it must preserve is the expensive one: a Review Issue that
failed on its merits and a Run stopped by its environment cost different
things. One deserves attention, the other deserves a re-run unchanged. The
session this Spec came from failed two legitimate security findings for an
environmental reason, and they had to be redone from scratch.

## Requirements

1. MUST record the Run's target branch and its revision at Preflight, so the
   mid-Run comparison has an anchor that does not depend on re-querying the
   forge.
2. MUST re-read the checkout and compare it against that anchor before every
   write boundary: each Batch commit, the review artifact commit, and Final
   Push.
3. MUST reach a terminal outcome distinct from `Failed` when the checkout no
   longer matches, so an interruption never reads as a Review Issue defect.
4. MUST leave the affected Review Issues unsettled rather than failed, so a
   re-run starts from their real state.
5. MUST make every review artifact commit and push target the Pull Request's
   head branch, asserted from Git rather than from the log.
6. MUST settle through the normal outcome path under ADR-0052; an interruption
   that cannot reach a terminal state is worse than the defect it replaces.
7. MUST NOT change what a Run does when the checkout never moves, asserted over
   the existing review tests unchanged.
8. MUST NOT move, check out, or restore the working tree.

## Subtasks

- [ ] Record the target branch and revision on the Run.
- [ ] Add the write-boundary re-read and the comparison.
- [ ] Add the terminal interruption outcome and its report line.
- [ ] Assert issues stay unsettled and artifacts land on the head branch.

## Acceptance Criteria

- [ ] A checkout moved before a Batch commit reaches the interruption outcome
      and commits nothing.
- [ ] A checkout moved before the review artifact commit reaches it too.
- [ ] A checkout moved before Final Push reaches it and pushes nothing.
- [ ] The interruption is a distinct terminal outcome, not `Failed`.
- [ ] Review Issues affected by an interruption are left unsettled, not failed.
- [ ] Every review artifact commit lands on the Pull Request's head branch,
      asserted from Git.
- [ ] A Run whose checkout never moves behaves exactly as it does today.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/store`
- instruction: `docs/adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/watch -count=1 -run 'Interrupt|Moved|Target|Boundary' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the interruption tests ran and passed.
- `go test ./internal/watch ./internal/cli ./internal/store -count=1`
  — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the Skill is task_03's bounded scope.

## References

- `_prd.md` → Core Features 3 and 4; Success Metrics 2, 3 and 4.
- `_techspec.md` → Disproven premise; Build Order 2.
- ADR-0052, ADR-0036.

## Result

Implemented a Git-backed target guard anchored by the Run's recorded PR Head
Branch and Preflight revision. The guard re-reads the checkout without moving
it before each Batch starts, immediately before every Batch commit, before the
ADR-0036 review artifact commit, and before Final Push. Successful Roundfix
commits advance the guard's expected revision while retaining the recorded PR
Head Branch.

Added `CheckoutMoved` as a distinct terminal Run outcome. Resolve and watch
settle it through `Store.CompleteRun`, release the Active Run lock, journal the
normal terminal outcome, report the mismatch and recovery action, and avoid
classifying the interruption as `Failed`.

### Focused checks

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/watch -count=1` — passed, 51 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/daemon -count=1` — passed, 175 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/store -count=1` — passed, 185 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/watch -count=1 -run 'TestRunCheckoutMoved|TestTargetGuard'` — passed, 6 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/daemon -count=1 -run 'TestFinalPushCheckoutMoved|TestResolveCycleCheckoutMoved'` — passed, 3 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/store -count=1 -run 'TestCompleteRunAcceptsCheckoutMoved|TestCreateReviewRunRecordsTarget'` — passed, 2 tests.
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/cli -count=1 -run 'TestRunWatchCheckoutMovedBeforeFinalPush|TestRunWatchArtifactEvidenceInherited|TestRunResolveCheckoutMoved|TestRunOperationalCommandAcceptsMVPFlags|TestRunResolveCommitsOnUserBranchWithoutRunBranch|TestRunResolvePushRunsFromUserCheckoutWithoutRunWorktree'` — passed, 9 tests.
- `rtk git diff --check` — passed with no whitespace errors after the Result update.

### Verification feedback repair

The Daemon's first combined watch/CLI/store verification exposed
`TestBranchIntegrityPreflightWatchDisregardsOnlySupersededFailedQACycles`.
The test composed the real Git Preflight and checkout reader with a synthetic
Branch Integrity refresh left behind by `withSuccessfulPreflight`; after the
fast-forward integration, the Run therefore recorded `integrated-head` while
the write guard correctly observed the actual Git revision.

`withRealReviewPreflight` now restores the production
`defaultRefreshBranchIntegrityHead` dependency along with the real Git checkout
reader, keeping the Preflight anchor and every later guard read on the same Git
source of truth.

- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/cli -count=1 -run '^TestBranchIntegrityPreflightWatchDisregardsOnlySupersededFailedQACycles$'` — reproduced the reported failure before the repair, then passed after the repair (1 test).
- `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/cli -count=1 -run '^(TestBranchIntegrityPreflight|TestRunResolveCheckoutMoved|TestRunWatchCheckoutMoved|TestRunWatchArtifactEvidenceInherited)'` — passed, 21 tests covering the repaired Branch Integrity fixture and the Task's CLI interruption/artifact paths.

### Acceptance evidence

- Checkout moved before a Batch commit: `TestResolveCycleCheckoutMovedAtBatchCommitCommitsNothing` observes the second boundary check fail and asserts zero committer calls.
- Checkout moved before the review artifact commit: `TestRunCheckoutMovedBeforeReviewArtifactCommitInterrupts` reaches `CheckoutMoved` and asserts the artifact publisher was never called.
- Checkout moved before Final Push: `TestRunWatchCheckoutMovedBeforeFinalPushInterruptsAndPushesNothing` reaches the persisted `CheckoutMoved` outcome and asserts zero pusher calls.
- Distinct terminal outcome: `TestCompleteRunAcceptsCheckoutMovedAsTerminal` proves compare-and-set completion, timestamping, and Active Run lock release; the watch tests assert the outcome is not `Failed`.
- Review Issues remain unsettled: `TestResolveCycleCheckoutMovedBeforeBatchLeavesReviewIssuesUnsettled` and `TestRunResolveCheckoutMovedBeforeBatchInterruptsWithIssueUnsettled` assert no failed, resolved, or invalid status is written.
- Artifact commit targets the PR Head Branch: `TestRunWatchArtifactEvidenceInheritedWithoutCurrentHeadPolling` reads `feature/review` with Git and asserts its ref equals the created review artifact commit.
- Unmoved checkout behavior: the focused CLI run includes the existing operational resolve/watch flow, user-branch commit flow, user-checkout push flow, and artifact Evidence inheritance flow unchanged.

The commands under `## Verification` were not run; they remain Daemon-owned.
