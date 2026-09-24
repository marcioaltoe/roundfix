---
task: task_01
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

## Result

Implemented line-based verdict classification and raw-answer persistence without changing the review command's exit mapping. The classifier now ignores letter case, trailing punctuation, surrounding `*` or `_` Markdown emphasis, and prose before the verdict line. It blocks with a reason naming whether both verdicts or neither verdict appeared, and a findings outcome carries only the text after the first `Findings:` marker. Every configured review-session outcome writes the unmodified Agent message to `pre-pr-review-answer.txt`; the JSON record names its absolute location in `answerPath`.

Focused checks:

- `GOCACHE=/private/tmp/roundfix-0160-task01-go-cache go test -count=1 -run '^TestReviewClassifiesVerdictVariants$' ./internal/cli` — passed; punctuated, emphasized, lowercase, and preamble variants kept exit `0` for no findings and exit `1` for findings, and the findings record excluded the preamble and marker.
- `GOCACHE=/private/tmp/roundfix-0160-task01-go-cache go test -count=1 -run '^TestReviewBlocksAmbiguousVerdict$' ./internal/cli` — passed; both-verdict and no-verdict answers kept exit `2`, outcome `blocked`, and reasons naming `both` or `neither`.
- `GOCACHE=/private/tmp/roundfix-0160-task01-go-cache go test -count=1 -run '^TestReviewKeepsTheRawAnswer$' ./internal/cli` — passed; reviewed, findings, and blocked outcomes preserved the exact answer bytes beside the record and named that file in `answerPath`.
- `GOCACHE=/private/tmp/roundfix-0160-task01-go-cache go test -count=1 -run '^TestReview' ./internal/cli` — passed with GitHub access permitted; existing review record, prompt, provider, fallback, artifact, and exit-code cases remained green.
- `GOCACHE=/private/tmp/roundfix-0160-task01-go-cache make verify-incremental` — no verdict recorded. The sandboxed attempt was blocked when a repository test reached `api.github.com`; the permitted rerun passed `go vet`, entered the repository-wide Go test stage, then remained silent through extended polling and was interrupted. This is not recorded as passing evidence.

Acceptance evidence:

- Punctuated, emphasized, lowercase, and preamble answers are classified: `TestReviewClassifiesVerdictVariants` covers each form and the extracted multi-line findings text.
- Both-verdict and no-verdict answers block: `TestReviewBlocksAmbiguousVerdict` proves the blocked outcome, named ambiguity, and unchanged exit `2`.
- The raw answer file exists and the record names it: `TestReviewKeepsTheRawAnswer` reads `record.AnswerPath`, compares it with the expected sibling path, and compares the file byte-for-byte with the Agent message for all three reviewer outcomes.
