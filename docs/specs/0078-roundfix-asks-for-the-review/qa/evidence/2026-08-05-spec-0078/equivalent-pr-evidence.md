# Equivalent observed evidence for hosted Pull Request journeys

Environment cause: the QA prompt confirms that no Pull Request is open for
the Spec target branch. The per-Run branch is never pushed and has no Pull
Request of its own. Hosted comment, Final Push, Evidence, and Event Stream
readback are therefore unreachable without expanding QA authority.

Unblocking action: open a Pull Request from
`ma/0078-roundfix-asks-for-the-review`, then run the read-only hosted checks
against that Pull Request. QA must still not commit, push, create the Pull
Request, or mutate review threads.

ADR-0080 permits a pass verdict when each environment-blocked row records this
cause and equivalent observed or supervised evidence. The following commands
ran fresh against build `bdf6ff8d4d680188a97986ee1340ab9ff052a499`:

1. `rtk go test ./internal/reviewsource/coderabbit -count=1 -run
   '^TestClientRequestReview' -v` — exit 0, 11 tests passed.
2. `rtk go test ./internal/watch -count=1 -run
   'TestRun(RequestsReviewForResolvedHeadBeforeMergeReadyEvidence|RequestsReviewBeforeNextRoundWaitAndStopsOnRefusal|ArtifactCommitDoesNotProduceSecondReviewRequest|DoesNotRequestReviewWithoutNewResolvedHead|RequestReviewDisabledPreservesNilRequesterControlFlow|StopsAtMaxRoundsReached)$'
   -v` — exit 0, 8 tests passed.
3. `rtk go test ./internal/cli -count=1 -run
   'TestRun(ResolveRequestsReviewAfterFinalPush|ResolveDoesNotRequestReviewWhenFinalPushIsSkipped|FetchNeverRequestsReviewWhenEnabled|WatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup)$'
   -v` — exit 0, 4 tests passed.
4. `rtk go test ./internal/runevent -count=1` — exit 0, 45 tests passed.

Observed equivalence by journey:

- One request per pushed Round: the watch seam asserted
  `Evidence(old head) -> request(new head) -> Evidence(new head)`, one request,
  and the exact pushed head/command payload.
- Artifact-only descendant: one request targets fix head `def456`; inherited
  artifact head `artifact789` creates no second request.
- Same-head replay: two calls return `published`, then `deduplicated`, one
  hosted-comment mutation, and one Run Event for each outcome.
- Resolve: the CLI integration test observes one Final Push followed by one
  request for the new public HEAD; the request callback fails if invoked before
  the push. The skipped-Final-Push case observes zero requests.
- Fetch: the CLI integration test runs fetch with asking enabled and observes
  zero requester calls. The four public scratch-CLI fetch invocations also
  bypassed the coherence refusal in every configuration pair before the
  sandboxed Run Database stopped them.
- Refusal/no retry: the watched next-head request is followed by skipped
  Evidence, terminal Review Skipped with the source reason, and exactly one
  request. Comment-list failure makes one list call, zero posts/events, and no
  retry; publish failure records no success event.
- Round cap: the seam requests once inside each completed Round, while the
  independent Max Rounds test ends after exactly two Rounds when `MaxRounds=2`.
  The repository's Project Config fixes the same cap at 2, so the composed
  upper bound is two requests after the initial Pull Request review.
- Event Stream: published and deduplicated calls each emit
  `review_source.request` with `head_sha`, `command`, and outcome. The request
  interface neither accepts nor returns Evidence, and watch continues to read
  head-bound Evidence after the request.

These are supervised equivalents, not claims that a hosted Pull Request was
used. The report retains the environment-blocked status for every affected
row.
