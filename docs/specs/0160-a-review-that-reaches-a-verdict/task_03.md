---
task: task_03
spec: 0160-a-review-that-reaches-a-verdict
status: completed
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

- [x] Implement the requirements above.
- [x] Add a test for each acceptance criterion.

## Acceptance Criteria

- [x] A candidate adding a Spec folder yields a prompt with its Decisions and a record naming it.
- [x] A candidate without one yields the prior prompt.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewPromptCarriesSpecDecisions|TestReviewPromptWithoutSpecIsUnchanged)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewPromptCarriesSpecDecisions TestReviewPromptWithoutSpecIsUnchanged; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Spec-aware prompt

## Result

- The review command now derives added or changed Spec folders from the candidate under the configured Spec Root, reads each PRD Decisions section and complete TechSpec from the candidate commit, appends labelled untrusted-context blocks to the reviewer prompt, and records the sorted consulted slugs in `specs`.
- A candidate without a changed Spec retains the prior prompt byte-for-byte and writes an empty `specs` array.
- Acceptance criterion 1: `TestReviewPromptCarriesSpecDecisions` passed in the focused `review.go` sweep. It uses a non-default configured Spec Root and proves the prompt carries only the PRD Decisions section, the complete TechSpec, the contradiction/rejected-alternative instruction, and the consulted slug.
- Acceptance criterion 2: `TestReviewPromptWithoutSpecIsUnchanged` passed in the focused `review.go` sweep. It compares the executed prompt byte-for-byte with the prior prompt builder and proves the record has no consulted Specs.
- Focused checks: `go test -count=1 -run '^(TestReviewRecord.*|TestReviewPrompt.*|TestReviewSession.*|TestReviewClassifies.*|TestReviewBlocks.*|TestReviewKeeps.*|TestReviewRunsTheClaudeProvider|TestReviewRefusesCodeRabbitNamingTheMissingSurface|TestReviewCommand(Exits|Defaults|Blocks|Does|Keeps|Classifies|None|RefusesProvider|UsesFallback|Proves|RefusesUnknown).*)$' ./internal/cli` passed; `make verify-incremental` passed with its existing GitHub integration tests granted network access after the sandboxed attempt was blocked at `api.github.com`.
- The Task's declared Verification command was not run; Daemon Verification remains the settlement authority.
