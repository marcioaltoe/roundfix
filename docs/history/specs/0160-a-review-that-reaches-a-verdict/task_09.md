---
task: task_09
spec: 0160-a-review-that-reaches-a-verdict
status: completed
type: backend
complexity: medium
---

# Task 09: A no-findings verdict must account for the whole answer

## Overview

Second pre-PR review of 2026-09-24, blocker: the classifier accepts a `No findings` line anywhere, so `"## Findings\n1. file:line bug\n\n## Security\nNo findings"`, findings under topic headings, and findings written after the verdict are recorded as reviewed with exit 0. It also reads `Findings: none` as a finding, blocks a findings body that quotes the verdict, leaves a stale answer file behind, records `skippedSpecs` as `null`, and drops Specs past the total context bound without naming them.

## Requirements

1. MUST classify as reviewed only when, after normalisation, the no-findings verdict is the sole content line, or the last content line with only a preamble before it that carries no list item, heading or `path:line` reference; any other content beside it blocks as ambiguous.
2. MUST read a `Findings:` line whose only text normalises to `none`, `n/a` or `no findings`, with nothing after it, as no findings.
3. MUST ignore verdict-shaped lines after a `Findings:` header when detecting conflicting verdicts.
4. MUST remove an existing answer file when the prompt did not reach a reviewer.
5. MUST record `skippedSpecs` as a JSON list, never `null`.
6. MUST record the slugs of Specs dropped by the total context bound.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The three answers the review reproduced block with exit 2.
- [ ] `Findings: none` is reviewed; a findings body quoting the verdict is findings.
- [ ] No stale answer file survives an unreached outcome; `skippedSpecs` is `[]`; dropped Specs are named.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewBlocksContentBesideANoFindingsVerdict|TestReviewReadsFindingsNoneAsNoFindings|TestReviewIgnoresAQuotedVerdictInsideFindings|TestReviewRemovesAStaleAnswerFile|TestReviewRecordsEmptySkippedSpecsAsAList|TestReviewRecordsDroppedSpecs)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewBlocksContentBesideANoFindingsVerdict TestReviewReadsFindingsNoneAsNoFindings TestReviewIgnoresAQuotedVerdictInsideFindings TestReviewRemovesAStaleAnswerFile TestReviewRecordsEmptySkippedSpecsAsAList TestReviewRecordsDroppedSpecs; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Classification

## Result

Implemented structural verdict classification so a clean verdict accounts for
the complete answer, while a `Findings:` body owns every later line. A
`Findings:` value of `none`, `n/a`, or `no findings` now produces a reviewed
outcome only when it is the terminal content and any preceding content is a
plain preamble. The command also clears stale answer evidence before an
unreached review, always serializes `skippedSpecs` as a list, and records each
Spec omitted by the total context bound.

Focused-check evidence:

- Acceptance criterion 1: the pre-change run of
  `TestReviewBlocksContentBesideANoFindingsVerdict` reproduced all three false
  passes with exit 0; the post-change focused run exits 0 with all three cases
  observing command exit 2. `TestReviewClassifiesVerdictVariants`,
  `TestReviewBlocksEmphasizedFindingsBesideNoFindings`, and
  `TestReviewBlocksAmbiguousVerdict` also pass individually.
- Acceptance criterion 2: the pre-change runs of
  `TestReviewReadsFindingsNoneAsNoFindings` and
  `TestReviewIgnoresAQuotedVerdictInsideFindings` failed; their post-change
  focused runs exit 0. The first covers `none`, `n/a`, `no findings`, a safe
  preamble, structured content before the clean verdict, and text after it.
  The second proves a verdict-shaped line in the findings body remains finding
  text. `TestReviewRecognisesEmphasizedFindingsHeader` also passes individually.
- Acceptance criterion 3: the pre-change runs of
  `TestReviewRemovesAStaleAnswerFile`,
  `TestReviewRecordsEmptySkippedSpecsAsAList`, and
  `TestReviewRecordsDroppedSpecs` failed on the stale file, `null` list, and
  unnamed omissions respectively; all three post-change focused runs exit 0.
  `TestReviewKeepsTheRawAnswer`,
  `TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached`,
  `TestReviewBoundsSpecContext`, and
  `TestReviewReportsSpecReadingFailureDistinctly` also pass individually.
- `rtk env GOCACHE=/tmp/roundfix-task09-gocache go vet ./internal/cli` exits 0.
- `rtk git diff --check` exits 0.

Not run: the Task's authored `## Verification` command, which is reserved for
Daemon Verification. A broader `^TestReview` sweep was not completed because it
attempted a live GitHub request; sandbox escalation was rejected because this
Task does not authorize external API use. The local-only focused tests above
do not use that boundary.

## Carry-forward provenance

- Source Run: `run_20260924T212727Z_e4cc3cafb0fa837a`
- Source commit: `f5a15f062d6c3d5dbba41c44a8e2a567fc73e744`
