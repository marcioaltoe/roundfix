---
task: task_02
spec: 0169-override-and-ordinal-follow-ups
status: pending
type: backend
complexity: low
---

# Task 02: Ordinals unique within a Spec, and a visible skip

## Overview

Two Tasks of one Spec creating different paths with the same ADR number are not reported until one file exists, and a Spec without `_tasks.md` does not list `SC-ORDINAL-CLAIMED` among the detectors it skipped.

## Requirements

1. MUST report `SC-ORDINAL-CLAIMED` for two claims of one Spec that share a number with different paths.
2. MUST list `SC-ORDINAL-CLAIMED` as skipped when a Spec has no Task Graph.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A same-Spec duplicate is reported; distinct numbers pass.
- [ ] A Spec without `_tasks.md` lists the skip.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/ordinal.go`
- interface: `internal/speccheck/citations.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSameSpecDuplicateOrdinalIsClaimed|TestMissingTaskGraphListsTheOrdinalSkip)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSameSpecDuplicateOrdinalIsClaimed TestMissingTaskGraphListsTheOrdinalSkip; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ordinal check
