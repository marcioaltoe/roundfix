---
task: task_02
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: completed
type: backend
complexity: high
---

# Task 02: A profile change carries the recorded decisions and refuses once

## Overview

`roundfix baseline plan` in an adopted repository takes every recorded
decision the selected profile still selects from the Setup Manifest, names
the Carried Decisions in its text, warns about each recorded decision it does
not carry, and, when answers are missing, names all of them, their
catalog-conditional dependents and the preservation mode in one refusal.
`roundfix baseline update --decision profile=<id>` refuses with the plan
command that performs a profile change. This answers frictions (a) and (b)
of the adopted Backlog Entry (Quaestor, 2026-10-26).

## Requirements

1. MUST add `CarryRecordedDecisions` to `PlanRequest` and `CarriedDecisions`
   to `PlanOutcome` as `_techspec.md` → Interfaces sketches; only
   `runBaselinePlanCommand` sets the field, for `--profile` and
   `--profile-file` alike. `baseline update`, the interactive workflow,
   `RecalculatePlan` and every other caller keep their inputs.
2. MUST implement API Contract 2 in a new `internal/baseline/recorded_decisions.go`:
   read the manifest through `ResolveManifestInput` (its `Manifest` field,
   for a resolved or incompatible manifest), carry by the closure the
   selected profile computes over flag answers plus carried values, let a
   flag always win, and name every decision not carried in one
   `baseline.decision.not-carried` warning with the reason the contract
   lists; an unparsable manifest yields `baseline.decision.manifest-unreadable`
   and contributes nothing; no manifest changes nothing (Invariant 1).
3. MUST print `Carried decisions: <ids>` after `Baseline Profile:` in the
   text plan when any decision was carried; the JSON plan document keeps its
   schema.
4. MUST implement API Contract 3 in `buildPlanWithCatalog`: one refusal
   naming the missing decisions, one `required when <id> is <value>: <ids>`
   part per unanswered decision whose `equals` effect (through
   `requireDecisions` or an activated module's `requiredDecisions`) requires
   non-optional decisions not yet required, and the preservation mode when
   absent, with the next action the contract fixes. When only one part is
   missing, the message and next action keep today's bytes.
5. MUST implement API Contract 1: `DeclareManifestDecisions` returns a
   `ProfileChangeRefusal` for the id `profile` before any normalization, and
   `baseline update` reports its message and next action with category
   `invalid` and exit `2`, writing nothing (Invariant 2). Any other unknown
   id keeps today's bytes.
6. MUST keep every existing assertion of `internal/cli/baseline_plan_test.go`
   and of the update decision tests; MAY add cases to the plan corpus.
7. MUST add the tests `_techspec.md` → Testing Approach names for task_02:
   `TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects`,
   `TestCarryRecordedDecisionsWithoutAManifestChangesNothing`,
   `TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange`,
   `TestBaselinePlanNamesEveryMissingAnswerInOneRefusal` and
   `TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand`. The
   cross-profile test adopts `go-cli-tui` with `baseline plan` and
   `baseline apply` in a temporary repository whose installed fake skills
   satisfy both profiles' alignment (`context7`, `context7-cli`,
   `exa-web-search`), then plans `python` with only
   `preservation.mode=preservation` and asserts exit `0`, decisions equal to
   the recorded values `python` selects, and the `Carried decisions:` line.
   Each test fails when its behavior is removed.

## Subtasks

- [ ] Characterize today's refusals in the new CLI test file before changing code.
- [ ] Add the carried-decision computation and wire it into the plan command.
- [ ] Join the refusal parts and add the conditional dependents.
- [ ] Add the typed profile refusal to the update.
- [ ] Add the named tests.

## Acceptance Criteria

- [ ] A profile change from `go-cli-tui` to `python` plans with only `--profile` and `preservation.mode`, and its text lists the Carried Decisions.
- [ ] A flag answer overrides a recorded value; a recorded decision the profile does not select is warned about, not carried.
- [ ] A plan without a manifest and without decisions refuses once, naming the base decisions, `runtime.backend` and `runtime.design` with their condition, and the preservation mode.
- [ ] `baseline update --decision profile=python` exits `2` naming `roundfix baseline plan --repo . --profile python`, and writes nothing.

## Context

- instruction: `docs/adr/0263-the-better-auth-base-path-is-a-repository-decision-a-baseline-update-records.md`
- instruction: `docs/adr/0341-a-profile-change-carries-the-recorded-decisions-and-a-restore-keeps-the-claude-skill-links.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`
- instruction: `internal/cli/baseline_human.go`
- instruction: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/update.go`
- creates: `internal/baseline/recorded_decisions.go`
- creates: `internal/baseline/recorded_decisions_test.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/cli/baseline_update.go`
- creates: `internal/cli/baseline_profile_change_test.go`
- interface: `internal/cli/baseline_plan_test.go`
- interface: `internal/cli/baseline_update_decision_test.go`

## Consumers

- `PlanRequest`, `PlanOutcome` and the plan refusal message — consumers: `./internal/baseline`, `./internal/cli`
- `DeclareManifestDecisions` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestCarryRecordedDecisions.*|TestBaselinePlan.*|TestBaselineUpdate.*Decision.*|TestBaselineUpdateDeclaresTheBetterAuthBasePath)$' ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects TestCarryRecordedDecisionsWithoutAManifestChangesNothing TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange TestBaselinePlanNamesEveryMissingAnswerInOneRefusal TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus TestBaselineUpdateRefusesAnUndeclarableDecision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the five new tests do not exist, so the first pass check fails. Afterwards the new tests and the existing plan and update decision tests that consume the changed contracts pass.

## References

- `_prd.md` → Goals; Core Feature 1; Core Feature 2; Core Feature 3; User Experience; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior 1-3; Interfaces; API Contract 1; API Contract 2; API Contract 3; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0263; ADR-0341


## Result

Implemented the task_02 slice for Daemon Verification:

- Added opt-in recorded-decision carrying to `PlanRequest` and the sorted
  `PlanOutcome.CarriedDecisions` projection. Only `runBaselinePlanCommand`
  enables it, for built-in profiles and profile-file drafts; update,
  interactive planning and recalculation retain their existing inputs.
- Read the Setup Manifest through `ResolveManifestInput.Manifest`, including
  incompatible recorded profiles. Shared the existing catalog closure with
  normalization, so explicit answers and fixed profile values determine which
  recorded answers remain selected. Invalid, unselected and fixed recorded
  values produce individual deterministic warnings; absent manifests leave
  input unchanged and unparsable manifests contribute only a warning.
- Printed the sorted carried IDs immediately after `Baseline Profile:` in
  text; the portable JSON Plan Document schema is unchanged. Carried-decision
  warnings survive both repository snapshot passes exactly once.
- Joined missing base decisions, conditional non-optional dependents from
  direct requirements or activated modules, and absent preservation mode into
  one refusal. Single-part refusals retain the prior message and next action.
- Added `ProfileChangeRefusal` before declaration normalization and mapped it
  to update category `invalid`, exit 2, with the original repository argument
  in the plan/apply next action. Other unknown-decision handling is unchanged.

Acceptance evidence:

1. `TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange` adopts
   `go-cli-tui` through CLI plan and apply with fake `context7`, `context7-cli`
   and `exa-web-search` skills. Planning `python` with only preservation mode
   exits 0, retains every selected recorded value and prints the sorted carried
   IDs. A narrowed profile-file draft also carries recorded answers.
2. The same CLI test proves a flag overrides `secondbrain.enabled` and a
   narrowed draft emits exactly one warning for its excluded decision.
   `TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects` additionally
   exercises conditional runtime decisions, flag-driven closure changes,
   fixed values, catalog-rejected values and an incompatible old profile.
   `TestCarryRecordedDecisionsWithoutAManifestChangesNothing` checks absent
   and unparsable manifests without changing explicit answers.
3. `TestBaselinePlanNamesEveryMissingAnswerInOneRefusal` asserts every base
   required ID, `required when autonomous.enabled is true: runtime.backend,
   runtime.design`, preservation mode and the exact combined next action.
4. `TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand` checks
   exit 2, category `invalid`, the recorded profile and plan/apply next action,
   JSON and text output, precedence over another unknown decision, and an
   unchanged repository tree.

Focused checks (all exited 0):

- `go test ./internal/baseline ./internal/cli -run '^(TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects|TestCarryRecordedDecisionsWithoutAManifestChangesNothing|TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange|TestBaselinePlanNamesEveryMissingAnswerInOneRefusal|TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand|TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus|TestBaselineUpdateRefusesAnUndeclarableDecision|TestBaselineUpdateDeclaresTheBetterAuthBasePath)$' -count=1 -timeout=3m`
- `go test ./internal/cli -run '^(TestBaselinePlan.*|TestBaselineUpdateRecordsADeclaredDecision|TestBaselineUpdateRefusesAnUndeclarableDecision|TestBaselineUpdateDeclaresTheBetterAuthBasePath|TestBaselineUpdateRefusesAProfileDecisionAndNamesThePlanCommand)$' -count=1 -timeout=4m`
- `go test ./internal/baseline -run '^(TestCarryRecordedDecisionsTakesWhatTheProfileStillSelects|TestCarryRecordedDecisionsWithoutAManifestChangesNothing)$' -count=1 -timeout=2m`
- `git diff --check`

Before production edits, the new update regression reproduced the original
unknown-decision message and generic next action. The initial missing-answer
fixture needed a seed file before it could reach the command boundary.
Isolated Go-overlay negative probes then removed CLI carrying, conditional
reporting, the typed update refusal, recorded values and absent-manifest
handling individually. Each of the five required new tests produced an
assertion failure when its corresponding behavior was removed; all overlays
lived outside the checkout.

No existing plan or update test assertion was changed. The initial worktree
already contained the Daemon's task-file status change; it was preserved.
Declared Verification, repository Verification and terminal settlement were
not run. No Task Graph, other Task, commit, push or Pull Request was changed.


### Verification feedback repair, attempt 1

The Daemon diagnostic identified an unintended serialized `CarriedDecisions`
field in the package's PlanOutcome characterization records. The TechSpec
requires this projection to remain non-serialized. Added `json:"-"` to
`PlanOutcome.CarriedDecisions`; existing characterization fixtures and
assertions remain unchanged.

Fresh focused checks after the repair (each exited 0):

- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -timeout=3m`
- `go test ./internal/cli -run '^TestBaselinePlanCarriesTheRecordedDecisionsAcrossAProfileChange$' -count=1 -timeout=3m`
- `git diff --check`

The package characterization now matches all five existing repository shapes,
and the CLI still carries the recorded answers and prints the carried IDs.
The declared Verification command was not rerun; Task status remains
Daemon-owned and unchanged by this repair.
