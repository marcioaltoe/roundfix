---
task: task_02
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: pending
type: backend
complexity: high
---

# Task 02: A Claude round-1 candidate over its budget is reviewed in Review Segments

## Overview

Today only a `codex` round-1 candidate above the byte bound is split into
Review Segments; a `claude` prompt above its 500,000-token budget is refused.
After this Task the segment planner takes a fit rule: `codex` keeps the byte
bound, byte for byte, and `claude` measures each segment's prompt as it will
be sent with the existing token estimate. Round 2 and the Correction Review
stay unsegmented. This answers the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md)),
reported by vortex, whose first suggestion asks for Claude segmentation.

## Requirements

1. MUST give `planReviewSegments` (`internal/cli/review_segments.go`) a fit
   rule (`reviewSegmentFit` in `_techspec.md` → Interfaces) used wherever it
   compares a commit's own diff or a grown segment with the bound today. The
   `codex` rule is `len(diff) <= reviewDiffBound`, so the Codex walk, its
   prompts, its refusals and its record keep their bytes (Invariant 2).
2. MUST make the `claude` rule `estimateReviewPromptTokens(prompt) <=
   claudeReviewPromptBudget`, where `prompt` is the segment prompt as it will
   be sent: the segment header written for segment `<i> of 4`,
   `buildReviewPrompt` of the segment and, for the first segment, the
   rendered Spec contexts (Invariant 3).
3. MUST, in `runConfiguredReviewSession` (`internal/cli/review.go`), for
   provider `claude` in round 1 when the full prompt with Spec contexts is
   estimated above the budget, plan segments with the Claude rule after the
   Spec contexts are read, append the Spec contexts to the first segment
   prompt only, and run the segments as Codex segments run, in one prepared
   session (API Contract 4). Below the budget nothing changes.
4. MUST refuse a `claude` candidate whose single commit's segment prompt does
   not fit with today's message `review prompt too large: <n> estimated tokens
   after omitting <k> path(s) exceeds the claude review budget of 500000
   tokens, half of its 1000000-token context window`, `<n>` being the whole
   unsegmented prompt's estimate, and a fifth segment with
   `review prompt too large: <n> estimated tokens after omitting <k> path(s)
   needs <s> review segments; a review allows at most 4`, both before any
   readiness probe or session preparation.
5. MUST keep round 2 and the Correction Review unsegmented for both providers
   with today's refusals (API Contract 5).
6. MUST add `EstimatedPromptTokens` (`estimatedPromptTokens`, omitted when
   zero) to `reviewSegment`, set for each Claude segment to its sent prompt's
   estimate and never for Codex; a segmented Claude record's
   `estimatedPromptTokens` is the sum of its segments' and its `diffBytes`
   the sum of theirs.
7. MUST create `internal/cli/review_claude_segment_test.go` with
   `TestReviewSegmentsAClaudeCandidateOverItsBudget`,
   `TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments`,
   `TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented` and
   `TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates`, with the inputs
   `_techspec.md` → Testing Approach names: three commits of 475,000 bytes
   must give exactly two segments, `<base>..<c2>` and `<c2>..<c3>`, because a
   byte-bounded planner would give three. Refusal reasons are asserted as a
   prefix of the record's reason and of the standard-error line
   `roundfix: review blocked: <reason>`, so the suffix task_03 appends does
   not change these tests. Every new top-level test and subtest calls
   `t.Parallel()` first.
8. MUST NOT change any existing test, guide, skill or ADR, the budget, the
   byte bound or the four-segment limit.

## Subtasks

- [ ] Give the planner a fit rule and keep the Codex rule byte-identical.
- [ ] Plan, run and record Claude segments in round 1 after the Spec contexts are read.
- [ ] Keep round 2, the Correction Review and an oversized single commit refused as today.
- [ ] Add the Claude segment tests.

## Acceptance Criteria

- [ ] A Claude candidate whose commits fit separately is reviewed in the fewest segments the token rule allows, in one session, each segment estimate within the budget.
- [ ] Five Claude segments, an oversized Claude commit and an over-budget Claude round 2 are refused before any provider call.
- [ ] The Codex segment, bound and Claude prompt-bound tests pass unchanged.

## Context

- instruction: `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`
- instruction: `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- instruction: `internal/cli/review_segment_test.go`
- instruction: `internal/cli/review_prompt_bound_test.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_segments.go`
- creates: `internal/cli/review_claude_segment_test.go`

## Consumers

- The Review Segment planner and the Pre-PR Review Record's `segments` — consumers: `./internal/cli`

## Verification

- `grep -q 'func TestReviewSegmentsAClaudeCandidateOverItsBudget(' internal/cli/review_claude_segment_test.go && grep -q 'func TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments(' internal/cli/review_claude_segment_test.go && go test -count=1 -timeout 30m ./internal/cli -run '^(TestReviewSegmentsAClaudeCandidateOverItsBudget|TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments|TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented|TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates|TestReviewSegmentsACandidateAboveTheBound|TestReviewBlocksACandidateNeedingMoreThanFourSegments|TestReviewSegmentsStopAtTheFirstBlockedAnswer|TestReviewGreedilyPacksFourSegments|TestReviewBlocksAnOversizedCommitHiddenByItsRange|TestReviewCommandRefusesAClaudePromptOverHalfItsWindow|TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound|TestReviewCommandKeepsTheCodexDiffByteBound|TestReviewOmitsConfiguredPathsFromTheBaseCommit)$'` — expected: exit 0. Before this Task the Claude segment test file does not exist, so the first check fails; afterwards the Claude segment tests pass and the Codex and Claude bound tests pass unchanged.

## References

- `_prd.md` → Goals; Core Feature 2; Success Metric 3; Success Metric 4; Decisions
- `_techspec.md` → Interfaces; Data Models; API Contracts 4-5; Invariants 2 and 3; Testing Approach; Build Order 2
- ADR-0352; ADR-0227; ADR-0282
