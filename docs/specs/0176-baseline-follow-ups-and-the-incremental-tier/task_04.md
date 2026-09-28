---
task: task_04
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: pending
type: backend
complexity: high
---

# Task 04: The incremental tier is a published decision

## Overview

The generated agent instructions and Spec routing guide require the incremental
verification command "declared by the active Baseline Profile", and call the
clause unmet when there is none. Only the Standard TypeScript Monorepo Profile
declares one. The template `internal/baseline/assets/templates/guides/agent-instructions.md`
renders only `{{verification.gate}}`. `resolveVerificationProjection` in
`internal/baseline/profile_alignment.go` projects only the `verification.gate`
decision into the Setup Manifest. So no adopter's guidance ever names an
incremental command, and every adopter Spec writes a waiver. This is Spec 0121
Core Feature 6.

The decision's value comes from the maintainer: an interactive answer, an
explicit `--adopt-suggested`, or a `--decision` flag. It is written into the
repository's Setup Manifest and its managed guides, and every agent that
verifies a Task reads those guides.

This is an authorized tooling Task. It may change only the files in its Context,
the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add the `verification.incremental` decision to
   `internal/baseline/assets/decisions.json`, directly after `verification.gate`,
   exactly as the TechSpec's "The incremental verification decision" gives it:
   type `string`, suggestion `rtk make verify-incremental`, no `default`, and one
   render binding to token `verification.incremental` of
   `template.guide.agent-instructions`.
2. MUST list `verification.incremental` directly after `verification.gate` in
   the `requiredDecisions` of `internal/baseline/assets/modules/core.json` and in
   the `entryDecisions` of `go-cli-tui.json`, `rust-cli.json` and
   `standard-typescript-monorepo.json`. The Standard TypeScript Monorepo
   `verification` list MUST stay byte-identical.
3. MUST add the token to `template.guide.agent-instructions` in
   `internal/baseline/assets/templates/index.json`, raising its version by one.
   MUST add the line `The selected incremental Verification is
   {{verification.incremental}}.` directly after the gate line of the template.
4. MUST reword `clause.core.verification-two-tiers` in `core.json` and
   `clause.spec.verification-two-tiers` in
   `internal/baseline/assets/modules/spec-workflow.json` exactly as the TechSpec
   gives them. Every other sentence of the core clause stays byte-identical, and
   the versions of the rules, guides and modules that contain them rise by one.
5. MUST make `resolveVerificationProjection` project a `verification.incremental`
   decision as the TechSpec's "Projection" paragraph describes. That means role
   `incremental` and classification `repository-command`, with a blocking
   `verification.command.undeclared` divergence and the message `selected
   incremental Verification command "<command>" has no matching local
   declaration` when the command is undeclared. A built-in Profile declaration
   with the same ID MUST be skipped while the decision is present, and MUST still
   project as a Profile expectation while it is absent.
6. MUST update, in the test files in the Context, only what the new required
   decision invalidates: complete decision lists, `--decision` arguments,
   Makefile fixtures that declare only `verify:`, the human prompt scripts, and
   the two decision counts of 15 in `internal/docscontract/publicdocs_test.go`,
   which become 16. `TestProfileDeclaresBothVerificationTiers` MUST keep its
   name and its independence assertions for both decisions. It MUST also assert
   that without the incremental decision the Profile expectation projects as
   `profile-expectation`. MUST NOT edit the parity corpus under
   `internal/baseline/testdata/parity-corpus/`. MUST NOT remove or rename a
   top-level test, so `docs/references/coverage-record.json` stays unchanged. A
   test that needs any other change MUST be reported in the Result instead of
   rewritten.
7. MUST add these tests to `internal/baseline/incremental_verification_test.go`,
   each negative case in its own test:
   - `TestIncrementalVerificationDecisionIsDeclaredWithoutADefault`: type,
     suggestion, no default, and the render binding.
   - `TestEveryBuiltInProfileRequiresTheIncrementalVerificationDecision`: the
     three Profiles and the core module.
   - `TestIncrementalVerificationDecisionProjectsLikeTheGate`: a declared
     command's projection fields, and exactly one `verification.incremental`
     projection for the Standard TypeScript Monorepo.
   - `TestUndeclaredIncrementalVerificationBlocksPlanning`: the blocking
     divergence, and planning that returns `action_required` with no Plan.
   - `TestPlanWithoutIncrementalVerificationNamesTheDecision`: `BuildPlan`
     without the decision returns `action_required` naming
     `verification.incremental` and no Plan.
   - `TestAgentInstructionsRenderTheSelectedIncrementalVerification`: the
     planned `docs/agents/agent-instructions.md` postimage contains the rendered
     line with the selected command.
   - `TestTwoTierClausesNameTheSelectedCommands`: both reworded clauses contain
     `selected incremental Verification`, and neither contains `declared by the
     active Baseline Profile`.
8. MUST add these tests to `internal/cli/baseline_incremental_verification_test.go`.
   Each drives `baseline update` on an adopted repository whose Setup Manifest
   has had the `verification.incremental` decision and projection removed:
   - `TestBaselineUpdateNamesTheMissingIncrementalDecision`: exit `3`, category
     `decision`, `newDecisions` naming `verification.incremental` with the
     suggested value, and a byte-identical tree.
   - `TestBaselineUpdateRefusesAnUndeclaredIncrementalSuggestion`:
     `--adopt-suggested` with a Makefile that has no `verify-incremental` target
     exits `3` and leaves the tree byte-identical.
   - `TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion`: with the target
     declared, `--yes --adopt-suggested` applies. The manifest records the
     decision and a role-`incremental` projection declared in `Makefile`, the
     agent instructions render the command, and a second update reports
     `current`.
   - `TestBaselinePlanWithoutIncrementalDecisionExitsThree`: `baseline plan`
     given only the gate decision exits `3`, names the decision, and writes
     nothing.
9. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --adopt-suggested --format text`,
   then the same refresh without `--adopt-suggested`, which MUST report
   `File changes: 0`. MUST NOT hand-edit a pin or a generated guide.
   `docs/agents/setup-context.json` MUST record `verification.incremental` =
   `rtk make verify-incremental` with a projection declared in `Makefile`.
10. MUST do all of the following:
    - add the decision to the published Decision Document in
      `docs/user-guide/context-driven-development.md`, and describe the
      single-gate migration there, naming `verification.incremental` and
      `--adopt-suggested`;
    - add an **Incremental Verification** entry to `CONTEXT.md` that names the
      decision and separates it from the complete repository Verification;
    - describe the decision and its migration in
      `.agents/skills/setup-context-driven/SKILL.md`, naming
      `verification.incremental`;
    - regenerate `skills/setup-context-driven/SKILL.md` with `make skills-sync`.

## Subtasks

- [ ] Declare and require the decision, render it, and reword the clauses.
- [ ] Project the decision and update the invalidated fixtures.
- [ ] Add the package and CLI tests.
- [ ] Regenerate the pins, goldens, guides and Setup Manifest, and update the
      docs.

## Acceptance Criteria

- [ ] Every built-in Profile requires `verification.incremental`. Planning
      without it names it, and a command the repository does not declare
      blocks.
- [ ] The rendered agent instructions and the Setup Manifest name the selected
      incremental command.
- [ ] A single-gate manifest makes `baseline update` name the decision and write
      nothing, and adopting a declared suggestion converges.
- [ ] The parity corpus and the coverage record are unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/profiles/go-cli-tui.json`
- interface: `internal/baseline/assets/profiles/rust-cli.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/baseline/profile_alignment_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/apply_test.go`
- interface: `internal/baseline/plan_characterization_test.go`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- creates: `internal/baseline/incremental_verification_test.go`
- interface: `internal/cli/baseline_plan_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `internal/cli/baseline_apply_test.go`
- interface: `internal/cli/baseline_update_test.go`
- creates: `internal/cli/baseline_incremental_verification_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestIncrementalVerificationDecisionIsDeclaredWithoutADefault|TestEveryBuiltInProfileRequiresTheIncrementalVerificationDecision|TestIncrementalVerificationDecisionProjectsLikeTheGate|TestUndeclaredIncrementalVerificationBlocksPlanning|TestPlanWithoutIncrementalVerificationNamesTheDecision|TestAgentInstructionsRenderTheSelectedIncrementalVerification|TestTwoTierClausesNameTheSelectedCommands|TestProfileDeclaresBothVerificationTiers|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestBaselineUpdateNamesTheMissingIncrementalDecision|TestBaselineUpdateRefusesAnUndeclaredIncrementalSuggestion|TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion|TestBaselinePlanWithoutIncrementalDecisionExitsThree|TestHumanBaselinePromptsOnlyForManifestMissingDecision|TestHumanBaselineAdoption|TestGuidanceCompositionJourney|TestBaselineUpdateNewDecisionRequiresActionWithoutWrites|TestBaselineUpdateAdoptsEverySuggestedDecision)$" ./internal/baseline ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestIncrementalVerificationDecisionIsDeclaredWithoutADefault TestEveryBuiltInProfileRequiresTheIncrementalVerificationDecision TestIncrementalVerificationDecisionProjectsLikeTheGate TestUndeclaredIncrementalVerificationBlocksPlanning TestPlanWithoutIncrementalVerificationNamesTheDecision TestAgentInstructionsRenderTheSelectedIncrementalVerification TestTwoTierClausesNameTheSelectedCommands TestProfileDeclaresBothVerificationTiers TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus TestBaselineUpdateNamesTheMissingIncrementalDecision TestBaselineUpdateRefusesAnUndeclaredIncrementalSuggestion TestBaselineUpdateAdoptsADeclaredIncrementalSuggestion TestBaselinePlanWithoutIncrementalDecisionExitsThree TestHumanBaselinePromptsOnlyForManifestMissingDecision TestHumanBaselineAdoption TestGuidanceCompositionJourney TestBaselineUpdateNewDecisionRequiresActionWithoutWrites TestBaselineUpdateAdoptsEverySuggestedDecision; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done; dout="$(go test -count=1 -tags docscontract -v -run "^(TestBaselineDecisionExamples|TestProjectConstraintDocumentation)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$dout"; exit 1; }; for name in TestBaselineDecisionExamples TestProjectConstraintDocumentation; do printf "%s\\n" "$dout" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "The selected incremental Verification is" docs/agents/agent-instructions.md && grep -q "selected incremental Verification" docs/agents/spec-routing.md && grep -q "verification.incremental" docs/agents/setup-context.json && grep -q "verification.incremental" docs/user-guide/context-driven-development.md && grep -q "Incremental Verification" CONTEXT.md && grep -q "verification.incremental" .agents/skills/setup-context-driven/SKILL.md && diff -r .agents/skills/setup-context-driven skills/setup-context-driven >/dev/null && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the eleven new named tests exists and no guide names an incremental command, so the command fails.

## References

- [_techspec.md](_techspec.md) — The incremental verification decision;
  Regeneration outputs
