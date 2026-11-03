---
task: task_03
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: completed
type: backend
complexity: medium
---

# Task 03: A size refusal names the largest files and the key, and the park carries it

## Overview

A review that still does not fit names, after its unchanged refusal text, the
three largest files of the diff it measured with their bytes and the key
`pre_pr_review.omit_paths` that would leave them out. The Delivery Queue
parks a blocked review with the first line of the record's reason, so the
queue's status carries the refusal without opening the record. This answers
the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md)),
reported by vortex, whose fourth suggestion asks for this message.

## Requirements

1. MUST create `internal/cli/review_size_refusal.go` with `reviewDiffFile`,
   `reviewLargestDiffFiles` and `reviewSizeRefusal` as `_techspec.md` →
   Interfaces and API Contract 6 state. Sections are split at lines starting
   `diff --git ` and paired, by position, with the entries of
   `git diff --name-only -z` run with the measured diff's own arguments
   (the same base, head, flags and `:(exclude,literal)` pathspecs), never
   parsed from the header. Files rank by section bytes, descending, ties by
   path ascending; at most three are named.
2. MUST route every size refusal through `reviewSizeRefusal`: the Codex byte
   bound (`checkReviewDiffBound` and its callers, every round), the Codex
   segment limit, the single-commit refusals of both providers, the Claude
   budget refusal of round 1 and round 2, and the Claude segment limit. Each
   keeps its text and appends
   `; largest files in the diff: <path> (<bytes> bytes), …; pre_pr_review.omit_paths in the base commit's Project Config leaves generated files out of the review`
   (Invariant 6). The files are those of the diff the refusal measured:
   round 1's whole scoped candidate diff, or the round-2 or Correction Review
   delta.
3. MUST, in `Engine.reviewCandidate` (`internal/delivery/engine.go`), park a
   `ReviewOutcomeBlocked` as `review-blocked: <first line>` where the first
   line is `ReviewResult.Reason` up to its first newline, trimmed, and as the
   bare `review-blocked` when that line is empty (API Contract 7). No other
   park changes.
4. MUST create `internal/cli/review_size_refusal_test.go` with
   `TestReviewSizeRefusalNamesTheLargestFiles`, one subtest per refusal of
   requirement 2, using the inputs `_techspec.md` → Testing Approach names:
   files of distinct sizes committed out of alphabetical order, two of equal
   size, the path `dados gerados/ação.json`, a larger omitted file that is
   never listed, and a two-file case. Expected bytes come from
   `git diff <base> <head> -- <path>` of each path, and each subtest asserts
   the whole standard-error line and exit `2`.
5. MUST create `internal/delivery/review_blocked_park_test.go` with
   `TestReviewBlockedParkCarriesTheReviewReason`: a two-line reason parks with
   its first line, an empty reason parks bare, and both classify `review`.
6. MUST update the declared breaks of `_techspec.md` → Declared breaks and
   nothing else in existing tests: the exact reasons of
   `TestReviewCommandRefusesAClaudePromptOverHalfItsWindow`,
   `TestReviewCommandKeepsTheCodexDiffByteBound`,
   `TestReviewBlocksACandidateNeedingMoreThanFourSegments` and
   `TestReviewBlocksAboveTheBoundWithoutAProviderCall` become the old text
   plus the suffix, and
   `TestDeliveryEngineParksDeclaredBlockersAndContinues` expects
   `review-blocked: review unavailable`. Every new top-level test and
   subtest calls `t.Parallel()` first.
7. MUST NOT change any guide, skill or ADR, the bound, the budget or the
   segment limit.

## Subtasks

- [ ] Add the largest-file helper and the refusal suffix.
- [ ] Route every size refusal through it.
- [ ] Park a blocked review with its reason line.
- [ ] Add the two new tests and update the declared pins.

## Acceptance Criteria

- [ ] Every size refusal keeps its text and ends with the three largest files of its diff, ranked by bytes then path, written as Git lists them, and the key.
- [ ] A blocked review parks with its reason line, or bare without one, and stays in the `review` Park Class.
- [ ] Only the declared tests change, and the Claude segment tests of task_02 still pass.

## Context

- instruction: `docs/adr/0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_segments.go`
- creates: `internal/cli/review_size_refusal.go`
- creates: `internal/cli/review_size_refusal_test.go`
- interface: `internal/cli/review_prompt_bound_test.go`
- interface: `internal/cli/review_segment_test.go`
- interface: `internal/cli/review_scope_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/engine_test.go`
- creates: `internal/delivery/review_blocked_park_test.go`

## Consumers

- The size refusal reason of the Pre-PR Review Record and the `review-blocked` park blocker — consumers: `./internal/cli`, `./internal/delivery`

## Verification

- `grep -q 'func TestReviewSizeRefusalNamesTheLargestFiles(' internal/cli/review_size_refusal_test.go && grep -q 'func TestReviewBlockedParkCarriesTheReviewReason(' internal/delivery/review_blocked_park_test.go && go test -count=1 -timeout 30m ./internal/cli -run '^(TestReviewSizeRefusalNamesTheLargestFiles|TestReviewCommandRefusesAClaudePromptOverHalfItsWindow|TestReviewCommandKeepsTheCodexDiffByteBound|TestReviewBlocksACandidateNeedingMoreThanFourSegments|TestReviewBlocksAboveTheBoundWithoutAProviderCall|TestReviewBlocksAnOversizedCommitHiddenByItsRange|TestReviewSegmentsAClaudeCandidateOverItsBudget|TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments|TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented|TestReviewSegmentsACandidateAboveTheBound|TestCorrectionReviewBlocksAnOversizedDeltaBeforeTheReviewer)$' && go test -count=1 -timeout 20m ./internal/delivery -run '^(TestReviewBlockedParkCarriesTheReviewReason|TestDeliveryEngineParksDeclaredBlockersAndContinues|TestParkClassesAndNextActions|TestParkClassRequiresAnExactBlockerOrColonPrefix)$'` — expected: exit 0. Before this Task the new test files do not exist, so the first check fails; afterwards the refusals and the park carry the new detail and every listed existing test passes.

## References

- `_prd.md` → Goals; Core Feature 3; Core Feature 4; Success Metric 5; Success Metric 6; Decisions
- `_techspec.md` → API Contracts 6-7; Invariant 6; Declared breaks; Testing Approach; Build Order 3
- ADR-0352

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

- Acceptance criterion 1: every local byte-bound, Claude budget and segment-limit refusal now uses `reviewSizeRefusal`. File sections are measured from the refused scoped patch and paired positionally with Git's NUL-delimited names using the same diff flags, revisions and literal exclusions. Ranking is descending bytes, then ascending path, capped at three. Round two and Correction Review pass their delta base and patch; segment-planning refusals pass the whole candidate, including when a rewrite's own patch is larger than the candidate patch.
- Acceptance criterion 1 evidence: `TestReviewSizeRefusalNamesTheLargestFiles` has eleven parallel cases covering both providers' single-commit refusals and segment limits, both round-two bounds, both Correction Review byte bounds, two files, and both providers' oversized rewrites hidden by the range. Each compares the complete stderr line, record reason and exit `2`, with bytes obtained from separate per-path Git diffs. Fixtures include unequal sizes committed out of alphabetical order, equal-size ties, `dados gerados/ação.json`, a larger omitted file and a large round-one-only file excluded from delta diagnostics. No refusal reaches the provider. The oracle preserves earlier sections' newlines and trims only the final section, matching the measured patch's existing Git runner behavior.
- Acceptance criterion 2: only `ReviewOutcomeBlocked` changes its park blocker: a nonempty trimmed first reason line follows `review-blocked: `; an empty or whitespace-only first line stays bare. Three parallel engine cases exercise two lines, an empty reason and a blank first line; all assert `ParkClassReview`.
- Acceptance criterion 3: only the five declared existing test functions changed their reason expectations. The task_02 Claude segment test file is unchanged, and its four top-level tests passed in a focused run. No guide, skill, ADR, byte bound, token budget or segment limit changed.
- Necessary mechanical production adjustment: `review_validation.go` already declared an unrelated `reviewDiffFile` for anchor ranges, colliding with the new type required here. Its four type references now use `reviewAnchorDiffFile`; fields and behavior are unchanged. Focused anchor-index tests passed.

Focused checks (all exited `0`; these are not the Daemon's declared Verification):

1. `rtk proxy go test ./internal/cli ./internal/delivery -run '^TestReview(SizeRefusalNamesTheLargestFiles|BlockedParkCarriesTheReviewReason)$' -count=1 -timeout=5m` — new refusal and park behavior passed. The refusal test subsequently gained the two hidden-rewrite cases, exercised in check 4.
2. `rtk proxy go test ./internal/cli -run '^(TestReviewCommandRefusesAClaudePromptOverHalfItsWindow|TestReviewCommandKeepsTheCodexDiffByteBound|TestReviewBlocksACandidateNeedingMoreThanFourSegments|TestReviewBlocksAboveTheBoundWithoutAProviderCall|TestReviewSegmentsAClaudeCandidateOverItsBudget|TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments|TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented|TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates)$' -count=1 -timeout=5m` — declared CLI pins and task_02 segment behavior passed.
3. `rtk proxy go test ./internal/delivery -run '^(TestDeliveryEngineParksDeclaredBlockersAndContinues|TestParkClassesAndNextActions|TestParkClassRequiresAnExactBlockerOrColonPrefix)$' -count=1 -timeout=5m` — the declared engine pin and Park Class checks passed.
4. `rtk proxy go test ./internal/cli -run '^(TestReviewSizeRefusalNamesTheLargestFiles|TestReviewDiffIndex.*|TestReviewPromptCarriesTheCandidateDiff|TestReviewSessionReadsWithoutWriting|TestReviewBlocksAnOversizedCommitHiddenByItsRange|TestCorrectionReviewBlocksAnOversizedDeltaBeforeTheReviewer)$' -count=1 -timeout=5m` — all eleven refusal cases, anchor-index checks, legacy session behavior and existing oversized-range/correction checks passed.
5. `rtk proxy git diff --check` — no whitespace errors.

Red evidence: before implementation, the new refusal cases reported the old reason without the suffix, and the two-line park case reported bare `review-blocked`. Initial fixture failures were traced to a missing nested directory and an overwritten provider profile missing its required fallback, then corrected in the new test only. A compile check exposed the existing type-name collision described above.

Execution note: the default executor could not start with the configured symlinked writable root. Task-scoped commands ran in the same checkout's physical `/Volumes/Brain-SSD/cache/.roundfix/worktrees/` location with reviewed execution permission. The only pre-existing changed path was this Task file's Daemon-owned status change; it was preserved.

Declared Verification and repository Verification were not run; the Daemon owns both and Task settlement. No commit, push, pull request, Task Graph edit or other Task edit was performed. Documentation and skill updates remain task_04's slice.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/review_validation.go`
