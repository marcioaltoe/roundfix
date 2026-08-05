---
task: task_02
spec: 0077-a-green-check-is-not-a-review
status: completed
type: backend
complexity: medium
---

# Task 02: Name the refusal so the stall is legible

## Overview

With the default inverted, a refused head already stalls — but it stalls as
`pending`, which reads as "no signal yet" rather than "the reviewer declined".
An operator cannot tell a refusal from an absence.

This slice recognises the documented refusal shapes and resolves them `skipped`
with the reason preserved. It explains the stall; it does not create it.

## Requirements

1. MUST recognise the documented shapes in which the Review Source declines to
   review the expected head, including the rate-limit refusal and the
   path-filter skip, and resolve them as `skipped`.
2. MUST recognise by class rather than by one literal: title casing and the
   documented variants all resolve the same way.
3. MUST read the authoritative signal the vendor documents — the comment as well
   as the check — because the check conclusion is green by design.
4. MUST preserve the refusal reason verbatim in the evidence detail, bounded by
   the existing helper.
5. MUST NOT let recognition widen what reaches `verified`. If refusal
   recognition and the closed default ever disagree, the default wins.
6. MUST NOT retry, re-request, or wait for capacity. That is deferred to the
   follow-on Spec and is out of scope here.

## Subtasks

- [ ] Recognise the documented refusal shapes and their reasons.
- [ ] Resolve them `skipped` ahead of every other classification.
- [ ] Add the class table and the #107 replay.
- [ ] Assert stale-head isolation for refusals.

## Acceptance Criteria

- [ ] The Pull Request #107 payload resolves `skipped` with a reason naming the
      rate limit.
- [ ] A path-filter skip resolves `skipped`, unchanged from today.
- [ ] Title-case variants of each documented refusal resolve identically,
      asserted by a table rather than one case.
- [ ] The refusal reason appears in the evidence detail.
- [ ] A refusal recorded against an earlier commit does not settle the current
      head.
- [ ] No payload reaches `verified` that did not reach it after task_01.
- [ ] No retry, re-request, or capacity wait exists in this Task's changes.

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/reviewsource/... -count=1 -run 'Refusal|Skip|RateLimit|Class' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the refusal tests ran and passed.
- `go test ./internal/reviewsource/... ./internal/watch -count=1` — expected:
  exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 3 and 5; Success Metrics 1, 2 and 4.
- `_techspec.md` → Interfaces; Build Order 2.
- ADR-0054.

## Result

Implementation:

- Added refusal-first classification for the documented CodeRabbit shapes. A
  structured `Review skipped` check preserves its bounded output reason, while a
  rate-limit status resolves only when CodeRabbit's authoritative rate-limit
  comment marker supplies the refusal reason.
- Added bounded GitHub issue-comment ingestion only for current-head
  rate-limit-shaped signals. The green status without that authoritative comment
  remains `pending`, preserving task_01's closed default.
- Made skipped evidence detail equal the bounded source reason and retained the
  same reason in `Evidence.Reason` for both check-run and commit-status refusals.
- Added the refusal class table, the recorded Pull Request #107 status/comment
  replay, the negative missing-comment case, and stale-head isolation at the
  existing CodeRabbit classifier seam.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=<worktree>/.gocache go test ./internal/reviewsource/coderabbit -count=1 -run '^TestEvidenceRefusal'`
  exited 1 because `IssueComment` and the comment-aware classifier did not exist.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/reviewsource/coderabbit -count=1`
  — exit 0; the complete CodeRabbit package suite passed after implementation.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/reviewsource/coderabbit -count=1 -run '^(TestEvidenceHierarchyPrecedence|TestEvidenceRefusalClassTable|TestEvidenceRefusalReasonIsBoundedVerbatim|TestEvidenceRateLimitWithoutAuthoritativeCommentStaysPending|TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals|TestEvidenceRefusalForStaleHeadDoesNotSettleCurrentHead|TestEvidenceRecordedCommitStatusCorpus|TestIssueCommentsMapGitHubRateLimitCommentJSON)$'`
  — exit 0 after the last implementation edit.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0 after the last
  implementation edit.
- Public Pull Request #107 inspection confirmed the recorded head
  `c6c14bece33bddf153c81c16029a97537f94d7c9`, commit status
  `CodeRabbit` / `success` / `Review rate limited`, CodeRabbit rate-limit comment
  marker, and `Review limit reached` heading used by the replay.

Acceptance evidence:

1. `TestEvidenceRecordedCommitStatusCorpus/pull_request_107_rate_limit_is_skipped`
   replays the recorded green status and authoritative comment as `skipped`; both
   `Evidence.Reason` and `Evidence.Detail` equal `Review limit reached`.
2. `TestEvidenceRefusalClassTable/path_filter_skip` preserves the existing
   structured path-filter refusal as `skipped`.
3. `TestEvidenceRefusalClassTable` covers lower/title-case variants for both the
   rate-limit and path-filter classes in one table, with identical states and
   reasons.
4. The class table asserts the source reason verbatim in `Evidence.Detail`, and
   `TestEvidenceRefusalReasonIsBoundedVerbatim` asserts the existing detail bound.
5. `TestEvidenceRefusalForStaleHeadDoesNotSettleCurrentHead` keeps an earlier-head
   refusal `pending` with no accepted evidence kind and records the old observed
   head.
6. The complete package suite passed with task_01's completion/default-deny and
   recorded corpus coverage; the new negative missing-comment case stays
   `pending`, and refusal recognition adds only `skipped` outcomes.
7. Diff inspection found no retry, re-request, backoff, capacity wait, or other
   follow-on policy; the changes are limited to signal ingestion,
   classification, and tests.

Declared `## Verification` commands were not run; the Daemon owns them.

Follow-ups: retry and re-request policy remains deferred to the follow-on Spec;
no part of it is included in this diff.
