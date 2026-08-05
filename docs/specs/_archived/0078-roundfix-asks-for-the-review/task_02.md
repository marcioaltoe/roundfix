---
task: task_02
spec: 0078-roundfix-asks-for-the-review
status: completed
type: backend
complexity: medium
---

# Task 02: Ask at the seam where the Round pushes

## Overview

`deps.Resolver.Resolve` performs a Round's fixes, commit, and Final Push and
returns the new `resolved.HeadSHA`. Every path after it waits for Evidence
bound to that head. This slice publishes the request in between.

It is the only Task that can produce a duplicate request, so it carries the
one-per-Round assertions — including the case that looks like two pushes: the
artifact-only docs commit, whose descendant inherits its parent's Evidence
under ADR-0036 and is not worth a review of its own.

## Requirements

1. MUST publish exactly one review request per Round, for the head
   `Resolve` reports, when `review_source.request_review` is enabled.
2. MUST publish it after the Final Push and before the Run waits for Evidence
   on that head, on both paths out of `Resolve`: the merge-ready confirmation
   and the next Round's wait.
3. MUST NOT publish for the artifact-only docs commit created after the Final
   Push.
4. MUST NOT publish when the Round pushed no new head.
5. MUST publish from `resolve` after its push, under the same configuration.
6. MUST NOT publish from `fetch` under any configuration.
7. MUST leave the Spec 0077 Evidence classification untouched: a refused
   request still resolves `skipped`, still ends the Run naming the refusal, and
   is followed by no second request.
8. MUST preserve today's control flow exactly when the configuration is
   disabled, including a nil requester.

## Subtasks

- [ ] Add the optional requester dependency and the call at the `Resolve` seam.
- [ ] Add the `resolve` command call site.
- [ ] Assert one request per Round across both post-`Resolve` paths.
- [ ] Assert the disabled default changes nothing.

## Acceptance Criteria

- [ ] A Round that pushes publishes exactly one request, for the pushed head.
- [ ] A Round whose Final Push is followed by the artifact-only docs commit
      still publishes exactly one request, for the fix head.
- [ ] A Round that pushes nothing publishes no request.
- [ ] `resolve` publishes one request after its push.
- [ ] `fetch` publishes none, asserted rather than assumed.
- [ ] A refused request ends the Run Review Skipped naming the refusal, with no
      second request in the same Run.
- [ ] With `request_review` disabled, every existing watch, resolve, and fetch
      test passes unchanged.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/cli/cli.go`
- instruction: `docs/adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/watch -count=1 -run 'Request|Round|Seam|Artifact'`
  — expected: exit 0; the seam tests ran and passed.
- `go test ./internal/watch ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./...`
  — expected: exit 0.
- `git diff --quiet HEAD -- .coderabbit.yaml .roundfixrc.yml`
  — expected: exit 0; turning the flow on is task_04's scope.

## References

- `_prd.md` → Core Features 1, 3, 5 and 6; Success Metrics 1, 2 and 4.
- `_techspec.md` → System Architecture; Build Order 2.
- ADR-0036, ADR-0054.

## Result

### Implementation

- Added the optional Review Source requester to the watch dependencies and
  configured request metadata to the watch request. A changed, non-empty head
  returned by `Resolve` now produces one request before either post-Resolve
  Evidence path; a nil requester or disabled configuration preserves the
  existing flow.
- Made Final Push report whether it actually pushed. Standalone `resolve`
  requests the pushed head only after that positive result; `fetch` has no
  requester call path.
- Kept artifact publication downstream of the fix-head request. The
  artifact-only descendant continues to inherit its parent's Evidence without
  producing another request.

### Acceptance evidence

- A Round that pushes requests exactly once for its pushed head:
  `TestRunRequestsReviewForResolvedHeadBeforeMergeReadyEvidence` passed and
  asserted the request payload plus the order `Resolve → request → Evidence`.
- An artifact-only docs commit produces no second request:
  `TestRunArtifactCommitDoesNotProduceSecondReviewRequest` passed with one
  request for `def456`, followed by inherited artifact Evidence for
  `artifact789`.
- A Round without a new pushed head requests nothing:
  `TestRunDoesNotRequestReviewWithoutNewResolvedHead` passed for empty and
  unchanged heads, and
  `TestRunResolveDoesNotRequestReviewWhenFinalPushIsSkipped` passed with zero
  pushes and zero requests.
- `resolve` requests after its Final Push:
  `TestRunResolveRequestsReviewAfterFinalPush` passed, asserted one Final Push,
  one request, the pushed `HEAD`, and rejected request-before-push ordering.
- `fetch` requests nothing:
  `TestRunFetchNeverRequestsReviewWhenEnabled` passed with asking enabled and
  zero requester calls.
- A refused request keeps Spec 0077 classification and stops further asking:
  `TestRunRequestsReviewBeforeNextRoundWaitAndStopsOnRefusal` passed with
  `Review Skipped`, the source refusal reason, and one request in the Run.
- Disabled configuration preserves existing behavior, including a nil
  requester: `TestRunRequestReviewDisabledPreservesNilRequesterControlFlow`
  passed, along with focused pre-existing watch, resolve, fetch, artifact
  inheritance, and Review Skipped tests.

### Focused checks

- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T180435Z_2151ad16c01390e5/.gocache go test ./internal/watch ./internal/cli -run 'TestRun(RequestsReview|ArtifactCommit|DoesNotRequestReview|RequestReviewDisabled|ResolveRequestsReview|FetchNeverRequestsReview)' -count=1`
  — passed.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T180435Z_2151ad16c01390e5/.gocache go test ./internal/cli -run 'TestRunResolve(RequestsReviewAfterFinalPush|DoesNotRequestReviewWhenFinalPushIsSkipped)|TestRunFetchNeverRequestsReviewWhenEnabled' -count=1`
  — passed after the final CLI test edit.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T180435Z_2151ad16c01390e5/.gocache go test ./internal/watch -run 'TestRun(WaitsFetchesResolvesToClean|ReentersFetchWhenHeadCheckFails|ReviewSkippedDuringMergeReadyPreservesTerminalEvidence|DoesNotConfirmMergeReadinessWithoutUntilClean)' -count=1`
  — passed.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T180435Z_2151ad16c01390e5/.gocache go test ./internal/cli -run 'TestRun(ResolvePushRunsFromUserCheckoutWithoutRunWorktree|WatchArtifactEvidenceInheritedWithoutCurrentHeadPolling|WatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup|FetchWritesReviewArtifactsUnderSpeclessRoot)|TestWatchSkipsFinalPushWhenAutoPushDisabled' -count=1`
  — passed.
- `rtk git diff --check` — passed.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and terminal settlement.
