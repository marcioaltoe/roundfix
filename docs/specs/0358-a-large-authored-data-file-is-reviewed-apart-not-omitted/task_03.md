---
task: task_03
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: completed
type: backend
complexity: high
---

# Task 03: The review reads the Review Data Paths apart, in its session

## Overview

task_02 takes the declared data files out of the code diff and plans their
Data Review Segments; nothing sends them yet. This Task makes every review
round send them after its code prompts, in the same prepared session, blocks
an oversized data diff before any provider call, and keeps the data findings,
record fields and answers apart, as ADR-0368 records. It answers point 5 of the
Fiscus report adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md).

## Requirements

1. MUST make `runConfiguredReviewSession` (`internal/cli/review.go`), in
   every round, after the code prompt or Review Segments are planned and
   before the readiness proof, do what `_techspec.md` → API Contract 4
   states: compute the data diff from `measuredBase` to the head with
   `:(literal)` pathspecs of the `data-path` paths, plan it with
   `planReviewDataSegments` and the provider's fit, block on a plan error
   with that reason and the size suffix of `reviewSizeRefusalForDiff` over
   the data diff, and otherwise append the data prompts after the code
   prompts, the single code prompt becoming the first of the list, sent in
   order through `RunPrepared` in the same session. Locate the edits by the
   code they follow, not by line number: Spec 0356 changes the same function
   first.
2. MUST record, as API Contract 5 states: `dataSegments` on the record, with
   `estimatedPromptTokens` only for `claude`; `omittedPaths` without any
   `data-path` path; `segments` for code Review Segments only; a data
   finding with `dataSegment` (one-based) and no `segment`, its anchor checked
   against its segment's diff; a single code prompt's finding without
   `segment`; ids `F1...` in send order. `classifyReviewSegmentAnswers`
   (`internal/cli/review_segments.go`) takes what it needs to label them.
3. MUST write the answer file with a single code answer under `--- Review ---`,
   code Review Segments under their current headers with `<k>` the code
   segment count, and each data answer under `--- Data Review Segment <i> of
   <d> ---`, word for word.
4. MUST keep byte for byte the prompts, record and answer file of a review
   whose base commit declares no `pre_pr_review.data_paths`, and every
   existing refusal (`_techspec.md` → Invariant 1); MUST NOT change
   `planReviewSegments`, the lineage decision, `internal/config`, any
   document, or any existing test.
5. MUST add `internal/cli/review_data_paths_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines,
   with the five tests of `_techspec.md` → Testing Approach for task_03:
   `TestReviewSendsADataPathApart` (subtests `codex` and `claude`),
   `TestReviewSplitsALargeDataFileAtLineBoundaries`,
   `TestReviewBlocksDataNeedingMoreThanFourDataSegments`,
   `TestReviewOmissionWinsOverADataPath` and
   `TestReviewReadsADataPathApartInRoundTwo`. Each commits the Project
   Config with the key and makes that commit `main` and the review base
   before the candidate commits, runs `roundfix review` through `runCLI`
   with `segmentReviewRunner` or `newLineageRunner`, and asserts what the
   Testing Approach lists, the API Contract 2 and 3 lines and the API Contract
   5 headers word for word. Each top-level test calls `t.Parallel()` first
   and uses only `t.TempDir()` repositories.

## Subtasks

- [ ] Send the Data Review Segments after the code prompts in every round.
- [ ] Record and label them, and join their answers.
- [ ] Add the command tests.

## Acceptance Criteria

- [ ] A declared catalog is read in its own Data Review Segment, for `codex` and `claude`, and the code prompt never holds it.
- [ ] A catalog above the bound is read in up to four segments with valid anchors; a fifth blocks before any session.
- [ ] An omitted path stays omitted, round 2 reads the catalog's change apart, and a review without the key is unchanged.

## Context

- instruction: `docs/adr/0368-a-large-authored-data-file-is-reviewed-apart-in-data-review-segments.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_segments.go`
- creates: `internal/cli/review_data_paths_test.go`
- instruction: `internal/cli/review_segment_test.go`
- instruction: `internal/cli/review_claude_segment_test.go`
- instruction: `internal/cli/review_lineage_test.go`

## Consumers

- The review prompts, the record fields `dataSegments` and `dataSegment`, and the answer file — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReviewSendsADataPathApart|TestReviewSplitsALargeDataFileAtLineBoundaries|TestReviewBlocksDataNeedingMoreThanFourDataSegments|TestReviewOmissionWinsOverADataPath|TestReviewReadsADataPathApartInRoundTwo)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewSendsADataPathApart TestReviewSplitsALargeDataFileAtLineBoundaries TestReviewBlocksDataNeedingMoreThanFourDataSegments TestReviewOmissionWinsOverADataPath TestReviewReadsADataPathApartInRoundTwo; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task none of the five tests exists, so the first pass line is missing; afterwards all pass and the whole CLI suite passes unchanged.

## References

- `_prd.md` → Goals; Core Feature 4; User Story 1; User Story 2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 4; API Contract 4; API Contract 5; Invariants 1 to 3; Testing Approach; Build Order 3

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging.

The session now plans the data diff from the round's measured base using literal
pathspecs before readiness checks. It appends Data Review Segments after the code
prompts in the same prepared session. A plan refusal measures the data diff alone
for its size suffix, including when the candidate also changes code.

The record separates data segments and findings from code segments, excludes
data-path entries from omittedPaths, and records data token estimates only for
Claude. Findings retain send-order ids and validate against their own segment.
Answer files preserve the required single-code, code-segment and data-segment
headers. A blocked answer stops subsequent prompts.

Acceptance evidence from focused checks:

- Catalog separation: TestReviewSendsADataPathApart passes for codex and claude
  with 9,000-entry JSON catalogs, exact scope and data headers, persisted fields,
  single-code findings without segment labels, data findings with dataSegment,
  and exact joined answer bytes. Its code-segments cases cover two code segments
  followed by data in the same session. Its Claude budget case admits a data diff
  above the Codex byte bound within the Claude token budget.
- Bounds and anchors: TestReviewSplitsALargeDataFileAtLineBoundaries passes with
  24,000 entries in two bounded data segments. Every added line is compared with
  git show at the hunk's coordinates. Line 20,000 is dismissed in data segment 1
  and stands in data segment 2. Blocked-code and blocked-data answer cases stop
  the sequence. TestReviewBlocksDataNeedingMoreThanFourDataSegments passes for
  five 14,000-entry catalogs: exact five-segment refusal and size suffix, exit 2,
  no answer file, and no provider probe, preparation or prompt.
- Existing behavior and round 2: TestReviewOmissionWinsOverADataPath retains
  omit-path and the original unwrapped answer. TestReviewReadsADataPathApartInRoundTwo
  sends exactly the changed entry's Git delta apart, pins declarations to the
  base even after the candidate clears data_paths, and resumes the recorded
  session. Focused existing tests for unchanged prompts, code-only segment
  records and answers, Claude admission, refusals and round-two behavior pass.
  The complete CLI suite remains for Daemon Verification.

Commands and outcomes:

- Pre-change signal: rtk proxy go test -count=1 -timeout 5m -run '^TestReviewSendsADataPathApart$/^codex$' ./internal/cli
  exited 1: only one prompt was sent instead of two.
- rtk proxy go test -count=1 -timeout 5m -v -run '^TestReviewSendsADataPathApart$' ./internal/cli
  exited 0, including both providers, mixed code/data segments and Claude's data
  token-budget case.
- rtk proxy go test -count=1 -timeout 5m -v -run '^TestReview(SendsADataPathApart|SplitsALargeDataFileAtLineBoundaries)$' ./internal/cli
  exited 0, including the line-coordinate oracle and blocked-answer cases.
- rtk proxy go test -count=1 -timeout 5m -v -run '^TestReview(BlocksDataNeedingMoreThanFourDataSegments|OmissionWinsOverADataPath|ReadsADataPathApartInRoundTwo)$' ./internal/cli
  exited 0. An earlier fixture used entries long enough to require ten segments;
  shortening the fixture values produced the specified five-segment input
  without changing the refusal assertion or production bound.
- rtk proxy go test -count=1 -timeout 5m -run '^TestReview(PromptWithoutSpecIsUnchanged|SegmentsACandidateAboveTheBound|SegmentsAClaudeCandidateOverItsBudget|KeepsCodexSegmentRecordsWithoutTokenEstimates|KeepsClaudeRoundTwoAndOversizedCommitsUnsegmented|SegmentsStopAtTheFirstBlockedAnswer|RoundTwo.*|SizeRefusalNamesTheLargestFiles|CommandExitsZeroOnExplicitClean|CommandExitsOneAndRecordsFindings|GreedilyPacksFourSegments|BlocksACandidateNeedingMoreThanFourSegments)$' ./internal/cli
  exited 0.
- rtk proxy git diff --check exited 0.

Only review.go, review_segments.go, the new review_data_paths_test.go and this
Result are changed by this Task. planReviewSegments, the lineage decision,
internal/config, existing tests and other documents are untouched. The Task
status remains Daemon-owned. Declared Verification, repository Verification and
settlement were not run; no commit, push or pull request was made.
