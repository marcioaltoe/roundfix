---
task: task_06
spec: 0157-reconciliation-and-one-repository-identity
status: pending
type: backend
complexity: low
---

# Task 06: Same report means same content

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `supersedingQAReport` now treats a Run-side report and a target-side report with the same file name as the same report without comparing their content, so a Run-side report with a different verdict can be declared superseded and its Run Branch released.

## Requirements

1. MUST treat a Run-side report and a target-side report with the same file name as the same report only when their blobs are identical.
2. MUST keep recognising an archived copy whose content is identical to the Run-side report.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Same name with different content is not proven superseded.
- [ ] Same name with identical content across the active and archived roots is recognised.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSupersedingQAReportRequiresTheSameContent|TestSupersedingQAReportRecognisesAnArchivedCopy)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSupersedingQAReportRequiresTheSameContent TestSupersedingQAReportRecognisesAnArchivedCopy; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Report identity
