---
task: task_02
spec: 0169-override-and-ordinal-follow-ups
status: completed
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

- [x] Implement the requirements above.
- [x] Add a test for each acceptance criterion.

## Acceptance Criteria

- [x] A same-Spec duplicate is reported; distinct numbers pass.
- [x] A Spec without `_tasks.md` lists the skip.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/ordinal.go`
- interface: `internal/speccheck/citations.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSameSpecDuplicateOrdinalIsClaimed|TestMissingTaskGraphListsTheOrdinalSkip)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSameSpecDuplicateOrdinalIsClaimed TestMissingTaskGraphListsTheOrdinalSkip; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ordinal check

## Result

- Implementation: the ordinal detector now compares claims within the current Task Graph before checking repository and other-Spec claims, and the missing-Task-Graph path records `SC-ORDINAL-CLAIMED` among its skipped detectors.
- Red evidence: `TestSameSpecDuplicateOrdinalIsClaimed` initially reported no `SC-ORDINAL-CLAIMED` finding for two paths with ordinal `0042`; `TestMissingTaskGraphListsTheOrdinalSkip` initially reported a skip list without `SC-ORDINAL-CLAIMED`.
- Acceptance criterion 1: `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestSameSpecDuplicateOrdinalIsClaimed$' ./internal/speccheck` passed. Its duplicate-ordinal subtest requires one finding naming both paths, and its distinct-ordinals subtest requires no finding.
- Acceptance criterion 2: `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestMissingTaskGraphListsTheOrdinalSkip$' ./internal/speccheck` passed and requires the missing `_tasks.md` skip to name `SC-ORDINAL-CLAIMED`.
- Focused package check: `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/speccheck` passed.
- Incremental repository check: `rtk env GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental` passed after the sandbox-blocked first attempt was rerun with network permission; vet, the repository test sweep, skill checks, and the build exited zero.
- Daemon Verification: not run; the Daemon owns the declared `## Verification` command and Task settlement.
