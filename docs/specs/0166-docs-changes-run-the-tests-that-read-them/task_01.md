---
task: task_01
spec: 0166-docs-changes-run-the-tests-that-read-them
status: pending
type: backend
complexity: low
---

# Task 01: No path selects nothing

## Overview

`ClassifyPath` sends `docs/**` and root Markdown to no test set, yet tests read `docs/user-guide/run-database-lifecycle.md`, `README.md` and `docs/history/specs`, so a change to them passes `make verify-changed` untested.

## Requirements

1. MUST remove the case that returns `NoSet` for documentation and Markdown, so those paths select both sets.
2. MUST keep Markdown inside a package directory selecting its owning set.
3. MUST update `TestClassifyPath` and the fixture change sets to the new contract.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] `docs/user-guide/run-database-lifecycle.md` alone and `README.md` alone each select both sets.
- [ ] `internal/app/README.md` still selects core.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDocumentationSelectsBothSets|TestRootMarkdownSelectsBothSets|TestClassifyPath)$" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestDocumentationSelectsBothSets TestRootMarkdownSelectsBothSets TestClassifyPath; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first two cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The classifier
