---
status: pending
created_at: 2026-08-10
updated_at: 2026-08-10
---

# A head the loop did not push is a head nobody reviews

`roundfix watch --source coderabbit --pr 153 --until-clean` timed out twice on
2026-08-10 — thirty minutes each, no Round, no Agent Session — against a Pull
Request that had twelve unresolved Review Issues waiting to be fixed.

## What happened

`.coderabbit.yaml` sets `auto_review.enabled: false` in this repository by
deliberate decision, so a review happens only when someone asks for it.
`.roundfixrc.yml` answers that with `review_source.request_review: true`, and
its own comment states the contract: *Roundfix publishes the review request
after a Round's Final Push.*

That contract holds only for heads the loop itself creates.
`requestReviewForResolvedHead` in `internal/watch/watch.go` returns early when
`resolvedHeadSHA == currentHeadSHA`, which is exactly the state at the start of
a Run whose head was pushed by hand. The head `d68e5a2a` was a Supervisor
commit fixing a Major from the previous review. No Round pushed it, so no
request was published for it, so CodeRabbit's commit status for that head stayed
`success` with the description `Review skipped: automatic reviews are disabled`.

The loop then waited for Review Source Evidence on a head nobody had asked
anyone to review, and `resultForTimedOut` reported the honest but late
diagnosis:

    Review Source signal was not recognised: CodeRabbit commit status
    "CodeRabbit" is success for the expected head: Review skipped: automatic
    reviews are disabled

The workaround is one comment: `gh pr comment 153 --body "@coderabbitai review"`,
then start the watch. Roundfix already owns every piece of that — the command
string is `defaultReviewRequestCommand`, the publisher is
`CodeRabbitClient.RequestReview`, and the idempotence marker keyed on the head
SHA is `ReviewRequestMarker`.

## The loop's own closure commit reproduces it

Once the request was published by hand, the loop converged: four Rounds, sixteen
Review Issues resolved, `Clean` in twenty-three minutes, and every head it
pushed itself arrived at `success :: Review completed` with no help. The
mechanism works exactly where the early return allows it to run.

Then the Daemon pushed `cb0d3a37 docs: review round 004 for pr 153` — the
separate review-artifact commit its own contract requires — and the CodeRabbit
check for that head went back to `Review skipped: automatic reviews are
disabled`. The closure commit is a head the loop pushed *after* deciding it was
done, so nothing asks for a review of it. GitHub still reports the Pull Request
`APPROVED`, so nothing breaks today; the next Run that watches this branch would
wait thirty minutes on it.

## Why it is worth fixing

The cost is not the waiting. It is that the Run ends `TimedOut` with zero work
done on twelve issues that were fetched successfully the whole time — a manual
`roundfix fetch --source coderabbit --pr 153` captured all twelve against the
same head at 20:52, with correct `head_sha` frontmatter, while the watch loop
sat blocked beside it. Sixty minutes of wall clock bought nothing.

It also misreads as a product defect. The first of the two Runs died with
`fetch CodeRabbit pull request reviews: gh command failed` after twenty-four
minutes; both `gh api` calls it makes reproduce clean, core rate limit at
4443/5000, so that one was a transient GitHub failure. Two different causes, one
symptom — a Run that waits half an hour and delivers nothing — and the second
cause is structural.

## Shape

Two directions, both non-binding.

**Request on the head the Run starts from, not only on heads it pushes.** The
early return exists to avoid re-requesting a review for a head that already has
one, and `ReviewRequestMarker(headSHA)` plus the existing `IssueComments` scan
already provide that idempotence independently of whether the head changed. The
guard could compare against published markers instead of against
`currentHeadSHA`. Whether the request belongs at Run start unconditionally, or
only when the configured Review Source has no accepted evidence for the head, is
the design question.

**Fail fast on a refusal Roundfix can already read.** `unrecognisedEvidenceReason`
had the whole answer at the first poll: the status is `success`, the description
says automatic reviews are disabled, and `request_review` is true. That
combination is not an unknown signal — it is a known state with a known remedy,
and thirty minutes of polling adds nothing to it. Reaching it in the first
minute with `ManualReviewCommand` already populated would turn a timeout into an
instruction.

Worth settling in the same work: a transient `gh` failure during
`WaitingForReview` currently discards the entire wait. `ghCommandError` already
carries `temporary`, and nothing consumes it on that path.
