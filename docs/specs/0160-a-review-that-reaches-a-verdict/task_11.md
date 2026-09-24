---
task: task_11
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: medium
---

# Task 11: A pass is the whole answer

## Overview

Third pre-PR review of 2026-09-24, blocker: the preamble check is a blocklist, and blockquoted lists, tables, setext headings, links to `path:line`, line ranges, unicode bullets, HTML, fenced code and plain prose describing a defect all pass beside `No findings` with exit 0. A findings header without an ASCII colon (`Findings`, `**Findings**`, `Findings：`) is also read as preamble.

## Requirements

1. MUST classify as reviewed only when the entire answer, after normalising case, trailing punctuation, surrounding Markdown emphasis and whitespace, is the no-findings verdict, or a `Findings:` line whose only text normalises to `none`, `n/a` or `no findings`; remove the preamble allowance.
2. MUST treat a line whose normalised text, with a trailing ASCII or fullwidth colon removed, equals `findings` as a findings header.
3. MUST block, as ambiguous, every other answer that carries no findings header.
4. MUST state the whole-answer rule in `.agents/skills/roundfix/SKILL.md` and `docs/user-guide/commands.md`, removing the preamble wording, and regenerate the mirror with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Every answer the third review reproduced blocks or reports findings; none exits 0.
- [ ] `No findings.`, `**No findings**` and `Findings: none` alone still pass.
- [ ] The skill, its mirror and the guide state the whole-answer rule.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewPassesOnlyAWholeAnswerVerdict|TestReviewBlocksEveryPreambleBesideNoFindings|TestReviewRecognisesAColonlessFindingsHeader)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewPassesOnlyAWholeAnswerVerdict TestReviewBlocksEveryPreambleBesideNoFindings TestReviewRecognisesAColonlessFindingsHeader; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Classification
