---
task: task_13
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: high
---

# Task 13: Keep the retention gate on the path it was scoped to

## Overview

The retention gate now stops a journey that completes today. Adopting one
Profile, updating it, then changing to a different Profile exits
`action_required` on this branch and succeeds on the default branch — proven by
running the same journey against both. That is the regression task 08's own
requirement forbids: a plan legitimately ready today must stay ready.

The cause is scope. The Spec's gate targets a **matching** Baseline identifier
whose Profile or catalog digests changed. A Profile change produces a
**different** identifier and has its own pre-existing transition path. The gate
is reaching a case it was never meant to judge.

This Task restores that boundary and closes the reason it went unnoticed: the
characterization corpus covers plan outcomes only, while the journey that broke
lives in the public command surface.

## Requirements

1. MUST restrict the retention gate to a Setup Manifest whose Baseline
   identifier matches the target; a change to a different Baseline identifier
   MUST follow the path it follows today.
2. MUST restore the adopt-update-change-Profile journey to the outcome it
   produces on the default branch.
3. MUST keep every same-identity behavior task 08 established: a changed
   Profile or catalog digest under a matching identifier still requires
   accounting, and an unaccounted clause still exits action-required.
4. MUST extend the characterization corpus to cover the public command
   journeys, not only plan outcomes, so a regression in this class fails the
   corpus rather than the gate.
5. MUST NOT weaken the gate on the path it does own: the disappearing-clause
   fixture must still exit action-required.
6. MUST NOT change the fail-closed apply, digest confirmation, or preimage
   binding.

## Subtasks

- [ ] Restrict the gate to a matching Baseline identifier.
- [ ] Restore the Profile-change journey to its default-branch outcome.
- [ ] Extend the characterization corpus to the public command journeys.
- [ ] Confirm the same-identity gate and its fixture still hold.

## Acceptance Criteria

- [ ] The adopt-update-change-Profile journey completes, and its Setup Manifest
      records the new Profile.
- [ ] A matching identifier with changed digests and a disappearing clause
      still exits action-required with its unaccounted count.
- [ ] A matching identifier with changed digests and every clause accounted
      still produces a ready plan.
- [ ] A change to a different Baseline identifier produces the same outcome it
      produces on the default branch.
- [ ] The characterization corpus includes the public command journeys, and a
      deliberate regression in one of them fails the corpus.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/cli/baseline_release_gate_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -run 'TestBaselineMacroJourneysPublicCLI/update_and_profile_change' -count=1`
  — expected: exit 0; the journey that regressed completes again.
- `go test ./internal/baseline -run TestSameIdentityDriftRequiresRetention -count=1`
  — expected: exit 0; the gate still fires on the path it owns.
- `go test ./internal/baseline -run TestReadyPlanNeverCarriesEmptyLedger -count=1`
  — expected: exit 0.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0; the corpus now covers the public command journeys.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 1 and 11; Non-Goals (turning working plans into
  blocked ones is out of scope).
- `_techspec.md` → Risks (the retention gate is the regression risk).
- ADR-0058.
