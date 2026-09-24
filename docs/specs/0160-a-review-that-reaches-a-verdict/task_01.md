---
task: task_01
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: medium
---

# Task 01: Classify the verdict by substance and keep the answer

## Overview

The classifier accepts exactly `No findings` or a message starting with `Findings:` and discards the answer, so a well-formed review is blocked and cannot be diagnosed.

## Requirements

1. MUST classify a verdict line regardless of letter case, trailing punctuation, surrounding Markdown emphasis, or a preamble before it.
2. MUST block, naming the ambiguity, an answer that carries both verdicts or neither.
3. MUST take the findings text as everything after the `Findings:` line.
4. MUST write the raw answer to `pre-pr-review-answer.txt` beside the record for every outcome that reached the reviewer, and name it in the record's `answerPath`.
5. MUST keep the exit-code contract unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Punctuated, emphasized, lowercase and preamble answers are classified.
- [ ] Both-verdict and no-verdict answers block.
- [ ] The raw answer file exists and the record names it.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewClassifiesVerdictVariants|TestReviewBlocksAmbiguousVerdict|TestReviewKeepsTheRawAnswer)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewClassifiesVerdictVariants TestReviewBlocksAmbiguousVerdict TestReviewKeepsTheRawAnswer; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Classification
