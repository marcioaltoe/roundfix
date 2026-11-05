---
task: task_02
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: pending
type: backend
complexity: medium
---

# Task 02: Review Data Paths leave the code diff and are planned into Data Review Segments

## Overview

With the key from task_01 the review still puts a declared data file into
the code diff. This Task gives it the reason `data-path`, so it leaves every
code diff and the code prompt says it is reviewed apart, and adds the plan that
cuts the data diff into at most four Data Review Segments by file, hunk and
line, as ADR-0368 records. task_03 sends what this Task plans. It answers
point 5 of the Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md).

## Requirements

1. MUST make `reviewScopedDiff` (`internal/cli/review.go`) give a changed
   path the reason `data-path` when no earlier reason applies and the
   declarations at the declaration base answer `ReviewsApart`, after the
   `omit-path` check, so every existing reason wins (`_techspec.md` → API
   Contract 2).
2. MUST make `appendReviewOmissions` leave `data-path` paths out of the
   `Omitted from this diff (not reviewed):` counts and lines and, when at
   least one exists, write after them, word for word:

   ```text
   Reviewed apart, after this diff, in Data Review Segments: <n> path(s) of pre_pr_review.data_paths.
   <path> (data-path)
   ```

   with one `<path> (data-path)` line per path in sorted order; without one,
   the prompt keeps its bytes.
3. MUST add `internal/cli/review_data_segments.go` with
   `reviewDataSegment`, `reviewDataFit` and `planReviewDataSegments` of
   `_techspec.md` → Interfaces and the behavior of API Contract 3: file
   sections, cuts at hunks and then at lines with recomputed hunk headers,
   greedy packing measured in bytes without rebuilding a string per line, the
   two refusals and the two-line prompt header, word for word.
4. MUST NOT change `runConfiguredReviewSession`, `planReviewSegments`,
   `classifyReviewSegmentAnswers`, the lineage, `internal/config`, any
   document, or any existing test; until task_03, `runConfiguredReviewSession`
   sends no data prompt.
5. MUST add `internal/cli/review_data_segments_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines,
   with `TestReviewDataSegmentsPackCutAndBound` of `_techspec.md` → Testing
   Approach, on diffs `git diff` writes in a `t.TempDir()` repository made
   with `gittest.InitRepo`, with a byte fit of a few kilobytes and the
   subtests `small files share one segment`, `a new file is cut at lines`,
   `a modified file is cut at hunks and lines`, `a line that fits no segment
   is refused` and `a fifth segment is refused`. The two cut subtests
   assert, for every piece, that it fits, that it starts with the file's
   `diff --git` line, and that each added and context line sits, by its hunk
   header, at the line where `git show <head>:<path>` holds that text; the
   refusal subtests assert the API Contract 3 texts. The test calls
   `t.Parallel()` first.

## Subtasks

- [ ] Give declared data files the `data-path` reason and the prompt line.
- [ ] Add the Data Review Segment plan.
- [ ] Add the plan tests.

## Acceptance Criteria

- [ ] A Review Data Path leaves every code diff and is named in the code prompt as reviewed apart; an omitted path stays omitted.
- [ ] The plan packs whole files, cuts a file at hunks and a hunk at lines, and every cut line keeps the line number the head's file gives it.
- [ ] A line that fits no segment and a fifth segment are refused with their texts.

## Context

- instruction: `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
- instruction: `internal/config/pre_pr_review.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_data_segments.go`
- creates: `internal/cli/review_data_segments_test.go`
- instruction: `internal/cli/review_omit_paths_test.go`

## Consumers

- `reviewScopedDiff`, its `data-path` reason and the code prompt's omission lines — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestReviewDataSegmentsPackCutAndBound$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestReviewDataSegmentsPackCutAndBound ' || { echo 'missing pass: TestReviewDataSegmentsPackCutAndBound' >&2; exit 1; }; grep -q '"data-path"' internal/cli/review.go || { echo 'missing data-path reason' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes, the reason exists, and the whole CLI suite, whose scope and omission tests pin every existing prompt, passes unchanged.

## References

- `_prd.md` → Goals; Core Feature 2; Core Feature 3; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 2 and 3; API Contract 2; API Contract 3; Invariants 2 and 4; Testing Approach; Build Order 2
