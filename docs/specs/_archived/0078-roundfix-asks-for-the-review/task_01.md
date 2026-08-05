---
task: task_01
spec: 0078-roundfix-asks-for-the-review
status: completed
type: backend
complexity: high
---

# Task 01: Publish one review request, once

## Overview

The primitive this Spec is built on: publish the Review Source's request
command on the pull request, for one pushed head, at most once.

Both halves already exist in `internal/reviewsource/coderabbit`.
`CommentOnPullRequest` publishes; `IssueComments` — added by Spec 0077 for
refusal recognition — lists. Deduplication is the marker in the body plus that
list, so a resumed Run, a replayed completion, and a second call for the same
head all converge on one comment.

This Task wires to nothing. It must be provable with no call site.

## Requirements

1. MUST publish the configured request command on the pull request for a given
   head SHA, carrying a Roundfix idempotency marker that names that head.
2. MUST report whether it published or found an equivalent request already
   present, and MUST NOT publish a second time for the same head.
3. MUST recognise only a Roundfix-authored marker for the same head as
   equivalent; a marker for another head, or another author's identical text,
   does not suppress the request.
4. MUST read the complete comment list rather than one page, because a
   truncated read re-asks and buys a review.
5. MUST record one Run Event carrying the head, the command, and whether it
   published or deduplicated.
6. MUST NOT wait for, poll for, infer, or return Review Source Evidence. Asking
   is not an answer, and ADR-0054 owns what is.
7. MUST NOT retry, back off, or wait for capacity; a failed request returns its
   error to the caller.

## Subtasks

- [ ] Add the requester with its marker format and comment-list lookup.
- [ ] Add the Run Event for published and deduplicated outcomes.
- [ ] Table-test the marker matching, including the near-miss cases.

## Acceptance Criteria

- [ ] A head with no prior request publishes exactly one comment containing the
      command and the marker for that head.
- [ ] The same head called twice publishes once and reports the second call as
      already present.
- [ ] A marker naming a different head does not suppress a new request.
- [ ] A comment carrying the command text without a Roundfix marker does not
      suppress a new request.
- [ ] A comment list spanning more than one page is read completely, asserted
      by a paginated fixture.
- [ ] A publish failure returns the error and records no success event.
- [ ] No Evidence type is read, returned, or waited on anywhere in this Task.

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/reviewsource.go`
- interface: `internal/runevent`
- instruction: `docs/adr/0054-review-source-evidence-determines-review-outcomes.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/reviewsource/... -count=1 -run 'Request|Marker|Idempot' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the request tests ran and passed.
- `go test ./internal/reviewsource/... ./internal/runevent -count=1` — expected:
  exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 2 and 7.
- `_techspec.md` → Interfaces; Data Models; Build Order 1.
- ADR-0054.

## Result

Implemented the standalone Review Source request primitive, its per-head HTML
marker, complete comment-list deduplication, and the `review_source.request`
Run Event. The primitive remains unwired; Task 02 owns its call sites.

Focused-check evidence:

- Red before implementation: `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T173603Z_2c5ff24f19316764/.gocache go test ./internal/reviewsource/coderabbit -run '^TestClientRequestReview' -count=1` failed to compile because `ReviewRequestMarker`, `ReviewRequest`, `RequestReview`, and the request event types did not exist.
- After implementation: the same focused command passed (`ok roundfix/internal/reviewsource/coderabbit`).
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T173603Z_2c5ff24f19316764/.gocache go test ./internal/reviewsource/coderabbit -count=1` passed.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260805T173603Z_2c5ff24f19316764/.gocache go test ./internal/runevent -count=1` passed.
- `rtk git -c core.fsmonitor=false diff --check` passed with no diagnostics.

Acceptance evidence:

- No prior request: `TestClientRequestReview/publishes_command_and_marker_when_absent` asserts one pull-request comment with the configured command and exact same-head marker, plus one published Run Event carrying the head and command.
- Same head twice: `TestClientRequestReview/same_head_called_twice_publishes_once` asserts outcomes `published`, then `deduplicated`, one comment total, and one Run Event per outcome.
- Different head: `TestClientRequestReview/different_head_marker_does_not_suppress` asserts that a valid marker for another head still publishes.
- Missing or foreign marker: `TestClientRequestReview/command_without_marker_does_not_suppress`, `marker_embedded_in_prose_does_not_suppress`, and `Review_Source_authored_marker_does_not_suppress` cover the near misses.
- Complete pagination: `TestClientRequestReviewReadsAllIssueCommentPages` places the same-head marker on page two of a real `GHClient.IssueComments` fixture and asserts deduplication with only the paginated list call.
- Publish failure: `TestClientRequestReviewPublishFailureReturnsErrorWithoutEvent` asserts the wrapped source error and zero success events. `TestClientRequestReviewListFailureReturnsWithoutRetry` additionally asserts one list attempt, no post, and no event.
- Evidence separation: inspection of `RequestReview` confirms it calls only `IssueComments`, `CommentOnPullRequest`, and the Run Event sink; it neither reads nor returns `Evidence`, and contains no polling, retry, backoff, or capacity wait.

The Task's declared `## Verification` commands were not run; the Daemon owns
that gate and Task settlement.
