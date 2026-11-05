---
task: task_03
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: pending
type: backend
complexity: medium
---

# Task 03: deliver status and deliver retry name a review that restarts at round 1

## Overview

When an item's new head will open a new Reviewer Lineage at round 1, nothing
says so until the full review has run and raised new findings. The
maintainer's answer asks that the queue warn before. This Task adds
`reviewLineageNotice`, which reads the item worktree's Pre-PR Review Record
and asks `decideReviewLineage` from task_02 for its `Restart`, and prints its
`Review:` line from `deliver status` and `deliver retry`, as ADR-0366 records.
It answers point 2 of the Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md).

## Requirements

1. MUST add `internal/cli/review_lineage_notice.go` with
   `reviewLineageNotice`, with the signature of `_techspec.md` → Interfaces
   and the behavior of API Contract 5, returning exactly:

   ```text
   Review: <slug> head <head> opens a new Reviewer Lineage at round 1, a full review of the candidate: <reason>
   ```

   It writes nothing, creates no directory (use
   `roundconfig.ResolveArtifactDirectory`, never `ValidateArtifactDirectory`),
   and returns `""` on any error.
2. MUST make `runDeliverStatus` (`internal/cli/deliver.go`) print that line,
   when it is not empty, for each `parked` or `reviewing` item, in queue
   order, after the loop that prints the `Warning:` lines and before the
   `Owner:` line; and make `runDeliverRetry` compute it from the item's
   recorded worktree, read from the queue before the store closes, when the
   retried stage is `reviewing`, and print it right after each
   `printDeliverRetryResult` call. Locate both edits by the code they follow,
   not by line number, because Spec 0355 may change the same functions first.
3. MUST keep every other line, their order and every exit code of `deliver
   status` and `deliver retry` byte for byte (`_techspec.md` → Invariant 5).
4. MUST NOT change `decideReviewLineage`, the review, `internal/delivery`,
   any document, or any existing test.
5. MUST add `internal/cli/deliver_review_lineage_notice_test.go`, opening
   with the `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment
   lines, with:
   - `TestDeliverStatusNamesAReviewThatOpensANewLineage`: subtests
     `rewritten head` and `conflict-free merge of the default branch`. Each
     uses `newReviewCommandFixture` with `newLineageRunner` answering round-1
     findings, runs `runDefaultBranchReview`, changes the head (`git commit
     --amend`, or `mergeAdvancedDefaultBranch`), creates
     `docs/specs` in the repository, opens the store under the fixture's
     home, creates a queue of `0356-example` for the repository, records the
     repository as the item's worktree on branch `feature/review`, parks it
     `review-findings`, and runs `deliver status` through `runCLI`. The first
     subtest asserts exit `0` and Surface Transcript 1 under its matching
     conventions; the second asserts no line starting `Review: `.
   - `TestDeliverRetryNamesAReviewThatOpensANewLineage`: the first subtest's
     setup, then `deliver retry 0356-example` through the `newDeliveryEngine`
     and `startDeliveryOwner` seams of `deliver_retry_test.go`, with a fake
     engine returning stage `reviewing` and an owner that prints `Delivery
     Owner: test-owner`, asserting exit `0` and the whole stdout of Surface
     Transcript 2.
   Each top-level test calls `t.Parallel()` first and uses only `t.TempDir()`
   directories.

## Subtasks

- [ ] Add the notice.
- [ ] Print it from deliver status and deliver retry.
- [ ] Add the status and retry tests.

## Acceptance Criteria

- [ ] For an item whose reviewed head was rewritten, `deliver status` and `deliver retry` print the `Review:` line of Surface Transcripts 1 and 2 before the review runs.
- [ ] After a conflict-free merge of the default branch, neither prints a `Review:` line.
- [ ] Every other line and exit code of both commands is unchanged.

## Context

- instruction: `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- instruction: `internal/cli/review_lineage.go`
- creates: `internal/cli/review_lineage_notice.go`
- interface: `internal/cli/deliver.go`
- creates: `internal/cli/deliver_review_lineage_notice_test.go`
- instruction: `internal/cli/deliver_retry_test.go`
- instruction: `internal/cli/deliver_park_status_test.go`

## Consumers

- The stdout of `deliver status` and `deliver retry` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDeliverStatusNamesAReviewThatOpensANewLineage|TestDeliverRetryNamesAReviewThatOpensANewLineage)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStatusNamesAReviewThatOpensANewLineage TestDeliverRetryNamesAReviewThatOpensANewLineage; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards both pass and the whole CLI suite, whose status and retry tests pin the other lines, passes.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 4
- `_techspec.md` → Current behavior 5; API Contract 5; Surface Transcript 1; Surface Transcript 2; Invariant 5; Testing Approach; Build Order 3
