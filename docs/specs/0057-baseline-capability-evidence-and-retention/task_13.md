---
task: task_13
spec: 0057-baseline-capability-evidence-and-retention
status: completed
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
- `go test ./internal/baseline -run '^TestSameIdentityDriftRequiresRetention$' -count=1 -v | grep -q -- "--- PASS: TestSameIdentityDriftRequiresRetention"`
  — expected: exit 0; the gate still fires on the path it owns.
- `go test ./internal/baseline -run '^TestReadyPlanNeverCarriesEmptyLedger$' -count=1 -v | grep -q -- "--- PASS: TestReadyPlanNeverCarriesEmptyLedger"`
  — expected: exit 0.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0; the corpus now covers the public command journeys.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 1 and 11; Non-Goals (turning working plans into
  blocked ones is out of scope).
- `_techspec.md` → Risks (the retention gate is the regression risk).
- ADR-0058.

## Result

Implementation:

- The retention resolver now enters same-identity clause accounting only when
  the existing Setup Manifest's Baseline identifier exactly matches the target
  Setup Manifest. A valid current manifest with a different identifier uses
  the pre-existing transition path.
- The plan suite now covers an applied `go-cli-tui` Baseline changing to
  `rust-cli`, including the ready result, target Profile, target Baseline
  identifier, and absence of a same-identity clause delta.
- `TestBaselinePlanCharacterizationPublicCommandJourneys` makes the existing
  real `TestBaselineMacroJourneysPublicCLI` suite part of the characterization
  gate. The bridge reuses the public command suite instead of copying its
  journey logic.

Focused checks:

- Before the fix,
  `GOCACHE="$PWD/.gocache-task13" go test ./internal/baseline -run '^TestDifferentIdentityKeepsExistingTransitionPath$' -count=1 -v`
  failed because the different-identity plan returned `action_required` with
  `baseline.go-cli-tui-0.0.1` having no unique maintained transition.
- After the fix,
  `GOCACHE="/private/tmp/roundfix-task13-go-cache" go test ./internal/baseline -run '^(TestSameIdentityDriftRequiresRetention|TestReadyPlanNeverCarriesEmptyLedger|TestDifferentIdentityKeepsExistingTransitionPath)$' -count=1`
  passed all three tests.
- `GOCACHE="/private/tmp/roundfix-task13-go-cache" go test ./internal/baseline -run '^(TestDifferentIdentityKeepsExistingTransitionPath|TestBaselinePlanCharacterizationPublicCommandJourneys)$' -count=1`
  passed both tests after the final production edit.
- With the exact different-identity fall-through defect temporarily restored,
  `GOCACHE="/private/tmp/roundfix-task13-go-cache" go test ./internal/baseline -run '^TestBaselinePlanCharacterizationPublicCommandJourneys$' -count=1 -v`
  failed at
  `TestBaselineMacroJourneysPublicCLI/update_and_profile_change`: planning
  exited `3` with `action_required`. The temporary mutation was then removed,
  and the preceding two-test check passed.
- `git diff --check` passed. `git status --porcelain` and
  `git diff --name-only` reported only `internal/baseline/` paths and this Task
  file; the Task file's pre-existing change is the Daemon-owned status update.

Acceptance evidence:

- The public Profile-change journey completed through the characterization
  bridge, whose canonical CLI assertion verifies that the Setup Manifest
  records `rust-cli`.
- The combined focused check kept the matching-identifier disappearing-clause
  case action-required with its unaccounted count.
- The same check kept the fully accounted matching-identifier case ready with
  its non-empty retention ledger and clause delta.
- The new different-identifier plan test passed with the default-branch ready
  outcome and no same-identity clause delta.
- The public command macro journeys are now reachable from the characterization
  test name, and the temporary regression proved the Profile-change journey
  turns that corpus red.
- No changed path falls outside `internal/baseline/` and this Task file, which
  is narrower than the allowed `internal/baseline/`, `internal/cli/`, and Task
  scope.

Daemon Verification commands were not run in this Agent turn.

### Verification feedback attempt 1

- Inspected the Daemon diagnostic artifact for the failed combined package
  test. It identified the baseline-level reviewed Profile adaptation and its
  real public CLI journey as the two failing paths.
- Root cause: the first repair removed the existing `SourceProfileID`
  compatibility path while narrowing the new gate to the exact target
  identifier. A reviewed Profile adaptation therefore bypassed same-identity
  accounting correctly but fell through to legacy transition lookup.
- Repair: exact target identifier equality remains the only entry to
  same-identity clause accounting. A current Setup Manifest whose identifier
  matches the reviewed draft's source Profile now returns through its existing
  compatibility path before legacy transition lookup. The valid-current-profile
  fallback remains in place for ordinary different-Profile changes.
- Before the repair,
  `GOCACHE="/private/tmp/roundfix-task13-repair-cache" go test ./internal/baseline -run '^TestProfileDraftPlanAcceptsMatchingSourceBaselineWithoutTransition$' -count=1 -v`
  reproduced the diagnostic: the reviewed adaptation returned
  `action_required` because its source Baseline had no maintained transition.
- After the repair,
  `GOCACHE="/private/tmp/roundfix-task13-repair-cache" go test ./internal/baseline -run '^(TestProfileDraftPlanAcceptsMatchingSourceBaselineWithoutTransition|TestDifferentIdentityKeepsExistingTransitionPath|TestSameIdentityDriftRequiresRetention|TestReadyPlanNeverCarriesEmptyLedger)$' -count=1`
  passed all four baseline controls.
- `GOCACHE="/private/tmp/roundfix-task13-repair-cache" go test ./internal/cli -run '^TestProfileAdaptationJourney$' -count=1`
  passed the real public Profile-adaptation journey.
- `GOCACHE="/private/tmp/roundfix-task13-repair-cache" go test ./internal/baseline -run '^TestBaselinePlanCharacterizationPublicCommandJourneys$' -count=1`
  passed the public macro characterization bridge after the final production
  edit.
- The failed declared Verification command was not rerun; the Daemon owns the
  next full Verification attempt and Task settlement.
