---
task: task_02
spec: 0157-reconciliation-and-one-repository-identity
status: pending
type: backend
complexity: low
---

# Task 02: Report identity survives archiving

## Overview

Archive moves a QA Report between the active and archived Spec roots without renaming it, and the supersession comparison breaks ties on the full path, so an archived copy is not recognised as the same report.

## Requirements

1. MUST treat a QA Report and its same-named archived copy as the same report when judging supersession.
2. MUST keep genuinely newer reports winning by date and sequence as today.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A Run-side report and its archived copy are the same report.
- [ ] A newer report still supersedes an older one.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSupersedingQAReportRecognisesAnArchivedCopy$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSupersedingQAReportRecognisesAnArchivedCopy"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md)
