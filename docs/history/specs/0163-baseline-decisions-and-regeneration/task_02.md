---
task: task_02
spec: 0163-baseline-decisions-and-regeneration
status: completed
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

## Result

Implemented carrier-aware Greenfield preservation planning. A stale managed
source now stays in the Source Baseline while the plan returns `blocked`, emits
`baseline.preservation.greenfield.managed-source-retained`, names
`preservation.mode=preservation`, and omits the decision skeleton. The human
classification boundary now performs the same profile-aware plan before it can
invoke semantic classification or emit a classification prompt, and it returns
the finding with the action result. `BuildPlan` exposes the same finding and
next action through its existing JSON result contract.

Focused-check evidence:

- Before the implementation, the new baseline regressions observed `ready`
  preservation and a complete `BuildPlan`, while the human regression received
  no action error. This reproduced the unreachable classification gate.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run 'Test(GreenfieldWithStaleManagedSourceRefusesNamingPreservation|GreenfieldWithoutStaleManagedSourceStillPlans|BuildPlanGreenfieldStaleManagedSourceNamesPreservation)$' ./internal/baseline`
  passed after the implementation. This covers the blocked
  preservation state, retained source, absent decision skeleton, JSON finding
  and next action, plus unchanged Greenfield planning without stale source.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run '^TestHumanGreenfieldRefusesBeforeClassification$' ./internal/cli`
  passed.
  The test observes no prompt bytes and no semantic-analyzer call before the
  human action is returned.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/baseline ./internal/cli`
  passed with network access required by
  the packages' GitHub-backed test boundary (`internal/baseline` 135.135s;
  `internal/cli` 211.639s). This also covers the updated existing carrier
  characterization under the new refusal contract.

Acceptance evidence:

1. `TestGreenfieldWithStaleManagedSourceRefusesNamingPreservation` and
   `TestBuildPlanGreenfieldStaleManagedSourceNamesPreservation` cover the
   preservation and JSON-plan refusal, including the exact finding, named
   Preservation route, retained source, and absent decision skeleton.
2. `TestHumanGreenfieldRefusesBeforeClassification` covers refusal before any
   classification prompt or analyzer call and verifies that the human action
   carries the finding and Preservation route.
3. `TestGreenfieldWithoutStaleManagedSourceStillPlans` covers the unchanged
   ready Greenfield path without a stale managed carrier.

The Task's declared `## Verification` command was not run; Verification remains
Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20260925T135958Z_cb4c08f055bb883c`
- Source commit: `99c6bdbc03928f8897622d505c238f5c8c5a9fb9`
