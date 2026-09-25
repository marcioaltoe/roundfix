---
task: task_01
spec: 0166-docs-changes-run-the-tests-that-read-them
status: completed
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

## Result

### Implementation

- Removed the documentation/Markdown `NoSet` branch from `ClassifyPath`, so
  unowned documentation reaches the existing `BothSets` fallback while the
  earlier `internal/**` classification keeps package Markdown in the core set.
- Updated the classifier table and documentation-only fixture change set, and
  added named regression tests for the two escaped paths.

### Focused checks

- Red signal before the classifier edit:
  `GOCACHE=/private/tmp/roundfix-task-0166-go-cache go test -count=1 -run '^(TestDocumentationSelectsBothSets|TestRootMarkdownSelectsBothSets|TestClassifyPath|TestSelectClassifiesFixtureChangeSets)$' ./internal/verifyselect`
  failed because the documentation and root Markdown cases selected `[]`
  instead of `[core baseline]`.
- After the classifier edit, the same focused command passed.
- `GOCACHE=/private/tmp/roundfix-task-0166-go-cache go test -count=1 ./internal/verifyselect`
  passed.
- `git diff --check` passed.
- The first focused attempt without the task-scoped `GOCACHE` did not compile
  because the sandbox denied access to the default macOS Go build cache; it
  produced no behavioral verdict.
- The daemon-owned `## Verification` command was not run.

### Acceptance evidence

- `TestDocumentationSelectsBothSets` proves
  `docs/user-guide/run-database-lifecycle.md` alone selects both sets.
- `TestRootMarkdownSelectsBothSets` proves `README.md` alone selects both sets.
- The `Markdown in core package` case in `TestClassifyPath` proves
  `internal/app/README.md` still selects core.
- The documentation and root Markdown rows in `TestClassifyPath`, plus the
  documentation-only row in `TestSelectClassifiesFixtureChangeSets`, record
  the new classifier and fixture contract.
