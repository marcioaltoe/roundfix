---
task: task_11
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

## Result

Replaced the preamble blocklist with a whole-answer check: `No findings` and
the three findings-none forms produce a reviewed outcome only when they are the
answer's sole nonblank content. Plain, emphasized, ASCII-colon, and
fullwidth-colon findings headers now start a findings body, so a later
verdict-shaped line remains findings text rather than creating a false pass.

Focused-check evidence:

- Acceptance criterion 1: before the production change,
  `TestReviewBlocksEveryPreambleBesideNoFindings` failed because tables,
  Markdown links to `path:line`, line ranges, and plain defect prose exited 0;
  `TestReviewRecognisesAColonlessFindingsHeader` failed all three header forms
  as blocked. After the change, both focused tests exit 0. Together they cover
  the reproduced blockquoted list, table, setext heading, linked path and line,
  line range, Unicode bullet, HTML, fenced code, plain prose, `Findings`,
  `**Findings**`, and `Findings：` answers; every command outcome is blocked or
  findings, never reviewed.
- Acceptance criterion 2:
  `TestReviewPassesOnlyAWholeAnswerVerdict` exits 0 for `No findings.`,
  `**No findings**`, and `Findings: none`. The adjacent classifier sweep for
  verdict variants, emphasized findings, contradictory and ambiguous answers,
  content beside no findings, all findings-none forms, and quoted verdicts also
  exits 0.
- Acceptance criterion 3: updated the canonical Roundfix skill and user guide
  to state the whole-answer and colonless-header rules, then ran
  `rtk make skills-sync` successfully. `rtk diff -q` confirms the canonical and
  distributed skills match, the focused `rtk rg` contract probe finds the new
  wording in all three surfaces, and
  `rtk env GOCACHE=/tmp/roundfix-task11-gocache go run -buildvcs=false ./cmd/roundfix skills check`
  passes. The first skill-check attempt used the default Go cache and was
  refused by the sandbox; the task-scoped cache rerun passed.
- `rtk env GOCACHE=/tmp/roundfix-task11-gocache go vet ./internal/cli` and
  `rtk git diff --check` exit 0. The repository-required
  `rtk env GOCACHE=/tmp/roundfix-task11-gocache make baseline-digests` also
  exits 0 and reports that derived artifacts already match their canonical
  sources.

Not run: the Task's authored `## Verification` command, which is reserved for
Daemon Verification.
