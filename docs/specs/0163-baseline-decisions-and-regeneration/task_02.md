---
task: task_02
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: backend
complexity: medium
---

# Task 02: Greenfield refuses before an unreachable step

## Overview

With stale managed carriers, `planRootPreservationWithCatalog` passes the Greenfield early return and asks for a classification that `promptBaselineClassification` never offers outside Preservation, so adoption fails at a gate it cannot pass.

## Requirements

1. MUST make Greenfield planning with a non-empty `staleManagedSources` return `PreservationStateBlocked` with the finding `baseline.preservation.greenfield.managed-source-retained` and a `NextAction` naming `preservation.mode=preservation`, with no decision skeleton.
2. MUST make `promptBaselineClassification` plan Greenfield and return that action before any classification or approval prompt.
3. MUST surface the same finding and next action from `BuildPlan` in JSON planning.
4. MUST keep Greenfield without stale managed source planning as today, and MUST NOT discard source or invent a classification.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Greenfield with a stale managed carrier refuses with the finding and next action in the preservation plan and in `BuildPlan`.
- [ ] The human command refuses before prompting for classification.
- [ ] Greenfield without a stale managed carrier still plans.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/preservation.go`
- interface: `internal/cli/baseline_human.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestGreenfieldWithStaleManagedSourceRefusesNamingPreservation|TestBuildPlanGreenfieldStaleManagedSourceNamesPreservation|TestGreenfieldWithoutStaleManagedSourceStillPlans|TestHumanGreenfieldRefusesBeforeClassification)$" ./internal/baseline ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestGreenfieldWithStaleManagedSourceRefusesNamingPreservation TestBuildPlanGreenfieldStaleManagedSourceNamesPreservation TestGreenfieldWithoutStaleManagedSourceStillPlans TestHumanGreenfieldRefusesBeforeClassification; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The Greenfield refusal
- [2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md](../../history/findings/2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md)
