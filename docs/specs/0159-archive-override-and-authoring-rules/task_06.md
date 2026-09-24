---
task: task_06
spec: 0159-archive-override-and-authoring-rules
status: pending
type: backend
complexity: low
---

# Task 06: Blame the latecomer, survive a broken neighbour

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `SC-ORDINAL-CLAIMED` compares a fulfilled claim with other active Specs, so the Spec that took a number first is refused because a later Spec claims it too, and the QA gate precondition runs this strict check. And one active Spec whose Task Graph fails to load aborts `spec check` for every Spec.

## Requirements

1. MUST skip the other-Spec comparison for a claim whose exact path already exists on the tree; the later claim is still reported through the tree comparison.
2. MUST skip an active Spec whose Task Graph fails to load when collecting other Specs' claims, leaving that failure to its own check.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A fulfilled claim is not reported when a later Spec claims its number; the later Spec is.
- [ ] An unloadable neighbour does not abort another Spec's check.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/ordinal.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestFulfilledOrdinalClaimIsNotBlamedForALaterClaim|TestOrdinalCheckSkipsAnUnloadableActiveSpec)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestFulfilledOrdinalClaimIsNotBlamedForALaterClaim TestOrdinalCheckSkipsAnUnloadableActiveSpec; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Claimed ordinals
