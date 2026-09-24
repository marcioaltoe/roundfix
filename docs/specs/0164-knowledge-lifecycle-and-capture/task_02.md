---
task: task_02
spec: 0164-knowledge-lifecycle-and-capture
status: pending
type: backend
complexity: medium
---

# Task 02: Review retirement reads recorded evidence

## Overview

`ClassifyReview` in `internal/spec/review_liveness.go` answers from local Git, so the same orphan Review Artifact classifies three ways as refs are fetched and deleted, and a squash merge can never retire one by ancestry. Legacy `docs/specs/_reviews/` artifacts that are not finished are never renamed.

## Requirements

1. MUST make `ClassifyReview` read `outcome.md` at the Review Artifact root, with front matter `pull_request_state` (`merged`, `closed` or `open`), `merge_commit` and `recorded_at`, and run no Git command.
2. MUST answer `ReviewFinished` for `merged` with a hexadecimal `merge_commit` (a squash receipt) whether or not that commit or the recorded head exists locally, and for `closed`.
3. MUST answer `ReviewLive` for `open`, and `ReviewUndecidable` with a reason naming what is missing for an absent or malformed record or a `merged` record without a merge commit.
4. MUST make `internal/baseline/history_layout.go` relocate a live or undecidable artifact under `docs/specs/_reviews/` to `docs/specs/reviews/` and keep its retained report, while a finished one relocates to `docs/history/reviews/` as today, through the existing collision checks.
5. MUST replace the expectations of `TestClassifyReviewLocalGit` that ADR-0163 revokes, as a declared break, and keep every Review Artifact root resolver answer outside the history root.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, one artifact classifies the same with its recorded head present, absent, and reachable only from fetched refs.
- [ ] A recorded squash receipt retires an artifact whose merge commit is absent locally.
- [ ] An artifact with no recorded outcome is undecidable and retained.
- [ ] A legacy `_reviews` artifact relocates to `docs/specs/reviews/` when live or undecidable, and to history when finished.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/review_liveness.go`
- interface: `internal/baseline/history_layout.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestClassifyReviewIgnoresObjectStoreAvailability|TestClassifyReviewAcceptsARecordedSquashReceipt|TestClassifyReviewWithoutARecordedOutcomeIsUnknown|TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness)$" ./internal/spec ./internal/baseline 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestClassifyReviewIgnoresObjectStoreAvailability TestClassifyReviewAcceptsARecordedSquashReceipt TestClassifyReviewWithoutARecordedOutcomeIsUnknown TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Review retirement
- ADR-0163
