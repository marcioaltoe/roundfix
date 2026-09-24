---
task: task_06
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: medium
---

# Task 06: No findings header escapes, no answer is invented

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `parseFindingsVerdictLine` strips emphasis only on the left, so `**Findings**:` and `_Findings_:` are not recognised and an answer carrying them beside a `No findings.` line is classified as reviewed with exit 0. The answer file is also written, empty, for outcomes that never reached the reviewer.

## Requirements

1. MUST normalise a findings line the same way as a no-findings line, removing Markdown emphasis on both sides of the word before looking for the colon.
2. MUST block any answer in which a line reads as a findings header after normalisation and a no-findings line also appears.
3. MUST write `pre-pr-review-answer.txt` and set `answerPath` only when the prompt was actually sent to a reviewer.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] `**Findings**:` and `_Findings_:` are recognised; beside `No findings.` they block.
- [ ] An outcome that never reached the reviewer has no answer file and an empty `answerPath`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewBlocksEmphasizedFindingsBesideNoFindings|TestReviewRecognisesEmphasizedFindingsHeader|TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached|TestReviewClassifiesVerdictVariants|TestReviewBlocksAmbiguousVerdict)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewBlocksEmphasizedFindingsBesideNoFindings TestReviewRecognisesEmphasizedFindingsHeader TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached TestReviewClassifiesVerdictVariants TestReviewBlocksAmbiguousVerdict; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task three of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Classification
