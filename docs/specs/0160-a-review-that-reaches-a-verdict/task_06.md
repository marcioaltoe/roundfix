---
task: task_06
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

## Result

- Findings headers now split at the first colon and normalize the header text with the same whitespace, Markdown-emphasis, punctuation, and case handling as `No findings`; `**Findings**:` and `_Findings_:` therefore participate in the existing exactly-one-verdict rule.
- The configured review-session path now reports whether execution crossed the prompt-sent boundary. The command writes `pre-pr-review-answer.txt` and sets `answerPath` only after that boundary; preparation and other pre-prompt failures persist only the blocked review record.
- Acceptance criterion 1: `TestReviewRecognisesEmphasizedFindingsHeader` passed for bold and italic headers and preserved their findings text; `TestReviewBlocksEmphasizedFindingsBesideNoFindings` passed for both forms beside `No findings.`, proving exit `2`, outcome `blocked`, and a reason naming both verdicts.
- Acceptance criterion 2: `TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached` passed after a session-preparation failure, proving zero prompt calls, an empty `answerPath`, and no `pre-pr-review-answer.txt` file.
- Red signal: before the implementation change, the three new regression tests failed because emphasized findings were unclassified or incorrectly reviewed and the pre-prompt failure wrote an empty answer artifact.
- Focused checks: `GOCACHE=/private/tmp/roundfix-0160-task06-go-cache go test -count=1 -run '^(TestReviewBlocksEmphasizedFindingsBesideNoFindings|TestReviewRecognisesEmphasizedFindingsHeader|TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached)$' ./internal/cli` passed; `GOCACHE=/private/tmp/roundfix-0160-task06-go-cache go test -count=1 -run '^TestReview' ./internal/cli` passed.
- `make verify-incremental` passed with process-table access. The sandboxed attempt had first reached only the unrelated force-stop integration failures caused by `operation not permitted` while enumerating the process tree; the permission-enabled rerun passed all packages, skill checks, and the build.
- The Task's declared Verification command was not run; Daemon Verification remains the settlement authority.

## Carry-forward provenance

- Source Run: `run_20260924T193553Z_7f2a4b9ff0f723fb`
- Source commit: `dddfce839b0d063803efee9dfad7483163090119`
