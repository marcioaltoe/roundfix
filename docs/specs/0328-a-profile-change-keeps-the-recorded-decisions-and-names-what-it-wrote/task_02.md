---
task: task_02
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: pending
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
