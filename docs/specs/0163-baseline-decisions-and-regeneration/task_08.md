---
task: task_08
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: test
complexity: low
---

# Task 08: The documented reconcile example parses

## Overview

Corrective Task from the QA gate of 2026-09-25: the repository Verification failed because `TestBaselineExamplesParse` checks every `roundfix baseline` example in the shipped skills, and `parsePublishedBaselineExample` has no case for the new `baseline skills reconcile` command, so the documented example is reported as not parsing although the command exists.

## Requirements

1. MUST add a `skills reconcile` case to `parsePublishedBaselineExample` that parses the example with the real `baseline skills reconcile` parser.
2. MUST add `TestBaselineSkillsReconcileExampleParses` proving a malformed reconcile example is still refused.
3. MUST NOT change production code or the documented example.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] `TestBaselineExamplesParse` passes with the documented reconcile example.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/baseline_documentation_contract_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBaselineExamplesParse|TestBaselineSkillsReconcileExampleParses)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBaselineExamplesParse TestBaselineSkillsReconcileExampleParses; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the documented example fails to parse, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
