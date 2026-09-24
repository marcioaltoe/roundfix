---
task: task_03
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: medium
---

# Task 03: Review a Spec's delivery against its decisions

## Overview

The reviewer judges the diff for correctness but never sees the decisions and rejected alternatives of the Spec being delivered.

## Requirements

1. MUST collect the Spec folders under the configured Spec Root that the candidate adds or changes.
2. MUST append, for each, the PRD `## Decisions` section and the TechSpec under a labelled block, asking the reviewer to report choices that contradict a decision or adopt a rejected alternative.
3. MUST record the consulted slugs in the record's `specs`.
4. MUST leave the prompt unchanged and `specs` empty for a candidate without a Spec folder.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A candidate adding a Spec folder yields a prompt with its Decisions and a record naming it.
- [ ] A candidate without one yields the prior prompt.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewPromptCarriesSpecDecisions|TestReviewPromptWithoutSpecIsUnchanged)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewPromptCarriesSpecDecisions TestReviewPromptWithoutSpecIsUnchanged; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Spec-aware prompt
