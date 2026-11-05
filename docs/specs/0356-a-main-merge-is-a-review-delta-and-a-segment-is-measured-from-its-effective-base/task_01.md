---
task: task_01
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: pending
type: backend
complexity: medium
---

# Task 01: Each Review Segment is measured from its Effective Base

## Overview

`planReviewSegments` diffs each first-parent commit of the candidate from the
commit before it, starting at the merge base, so an item commit cut from an
old default branch shows the default branch's later work as removed and the
item's merge of the default branch shows it as added. On 2026-11-05 that
blocked Fiscus's Spec 0070 by size although its own change fit. This Task
adds the measurement of ADR-0366, `reviewMeasuredFrom`, and starts every diff
the segment plan measures from it, for `codex` and `claude` alike. It answers
point 4 of the Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md).

## Requirements

1. MUST add `internal/cli/review_effective_base.go` with
   `reviewEffectiveBase` and `reviewMeasuredFrom`, with the signatures of
   `_techspec.md` → Interfaces and the behavior of API Contract 1: the carry
   runs `git merge-tree --write-tree <from> <E2>` and takes its first output
   line; exit status 1 returns `from, false` and no error; any other failure
   returns an error naming both commits.
2. MUST make `planReviewSegments` (`internal/cli/review_segments.go`) start
   the own diff of each commit and the diff of each growing segment from
   `reviewMeasuredFrom(base, <start>, <commit>)` instead of `<start>`, keeping
   the declaration base argument of `reviewScopedDiff` at `base`
   (`_techspec.md` → API Contract 2).
3. MUST add `MeasuredFrom string` with JSON name `measuredFrom,omitempty` to
   `reviewSegment`, set to the object a segment's diff starts from only when
   it differs from the segment's `Base`, and make `reviewSegmentHeader` append,
   when it is set, this second line, in the prompts sent and in every fit
   estimate:

   ```text
   This segment's diff starts from <measuredFrom> instead of <segment base>, so it leaves out what the default branch changed.
   ```

4. MUST keep byte for byte: the first header line, every size refusal, the
   review bound, the Claude budget, the four-segment limit, the prompt and
   record of a candidate whose commits share one Effective Base, and a
   conflicting carry's measurement from `<start>` (`_techspec.md` →
   Invariants 2 and 4).
5. MUST NOT change `decideReviewLineage`, `runReviewCommand`,
   `runConfiguredReviewSession`, the round-2 or Correction Review delta,
   `internal/cli/deliver.go`, any skill, any other document, or any existing
   test.
6. MUST add `internal/cli/review_effective_base_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, with
   the two tests of `_techspec.md` → Testing Approach for task_01:
   - `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange`, with the rows
     `same Effective Base`, `start on the base`, `merge of the default
     branch` and `conflicting carry` as described there, asserting the
     returned object (and, for the merge row, that it is a tree whose diff to
     the merge is empty) and the carried flag;
   - `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase`, subtests
     `codex` and `claude`, running `runDefaultBranchReview` with
     `segmentReviewRunner`, asserting exit `0`, outcome `reviewed`, the three
     segments' `head`, `base` and `measuredFrom` (the old base, a tree, none),
     the second header line exactly in the first two prompts and absent from
     the third, and that no prompt holds 1,000 consecutive `m`.
   The file also defines the helper `runDefaultBranchReview(t)`, which runs
   `review --base main` through `runCLI` and decodes the record from stdout;
   tasks 02 and 03 reuse it. Each top-level test calls `t.Parallel()` first
   and uses only `t.TempDir()` repositories.

7. MUST add to `docs/user-guide/commands/review.md`, after the paragraph
   ending "all segments are one round.", the Effective Base paragraph of
   `_techspec.md` → API Contract 6, word for word (line wrapping free), and
   change nothing else in the guide; task_04 writes the rest of API Contract
   6.
8. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the guide edit and the code changes, and
   never hand-edit it.

## Subtasks

- [ ] Add the Effective Base measurement.
- [ ] Measure every segment diff from it and record `measuredFrom` with its header line.
- [ ] Add the measurement and segment tests.
- [ ] Write the Effective Base paragraph into the `review` guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A commit cut from an older default branch and a merge of the default branch add none of the default branch's work to a segment, for `codex` and `claude`.
- [ ] A segment measured from elsewhere than its base records `measuredFrom` and says so in its prompt; every other segment is unchanged.
- [ ] A conflicting carry keeps today's measurement.

## Context

- interface: `docs/user-guide/commands/review.md`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- creates: `internal/cli/review_effective_base.go`
- creates: `internal/cli/review_effective_base_test.go`
- interface: `internal/cli/review_segments.go`
- instruction: `internal/cli/review_segment_test.go`
- instruction: `internal/cli/review_claude_segment_test.go`

## Consumers

- The Review Segment plan, its `segments[].measuredFrom` record field and its segment prompts — consumers: `./internal/cli`
- The `review` guide fingerprint in the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange|TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- 'starts from the Effective Base of the commit it ends at' || { echo 'missing Effective Base paragraph' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/cli && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither test exists, so the first pass line is missing; afterwards both pass, the guide states the Effective Base paragraph, and the whole CLI suite, which consumes the segment plan, and the documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 1; Core Feature 2; Success Metric 1
- `_techspec.md` → Current behavior 3 and 4; API Contract 1; API Contract 2; API Contract 6; Invariants 1, 2 and 4; Testing Approach; Build Order 1
