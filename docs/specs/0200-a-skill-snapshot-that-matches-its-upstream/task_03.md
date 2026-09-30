---
task: task_03
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: pending
type: backend
complexity: medium
---

# Task 03: A profile's guides name only skills its setup lists

## Overview

The catalog checks only a module's required skills against the profile's setup, so a dispatched skill outside it went unnoticed. This Task adds the dispatch check, makes `core` require and dispatch `exa-web-search`, which a universal required capability already demands, and moves the Context7 capability to `context7-cli` while an installed `context7` still satisfies it.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST make `validateProfiles` report `catalog.profile.skill.dispatch-outside-setup` with the profile as subject and the skill as detail for every skill that `profileNamedSkills` returns and the profile's setup does not list, as the TechSpec's Interfaces define it: the `skill` (else `id`) of every `skillDispatch` entry of a selected module, and every skill of an activation bundle whose owner is a selected module.
2. MUST add `exa-web-search` to `core.json` as the TechSpec's "Fixed texts" give it for task_03, and raise the `core` module version by one from its value on this Task's base. Replace strings and numbers in place.
3. MUST change `capability.context7` as the TechSpec's Interfaces and "Fixed texts" give it: probe `{"skill": "context7-cli", "priorSkills": ["context7"]}`, the new NextAction, and a package-level `universalCapabilityRestoreSkills` that `applyUniversalCapabilityRemediation` reads. `collectInstalledSkillEvidence` checks the probe's skill, then each prior skill in order, and reports the first present `SKILL.md`; an unsafe prior name is invalid evidence.
4. MUST change both `--skill context7` commands in the "Profiles" section of `docs/user-guide/context-driven-development.md` to `--skill context7-cli`, and add after the paragraph that introduces them: "An installed `context7` skill, the name earlier snapshots used, still satisfies the Context7 capability."
5. MUST add `internal/baseline/catalog_dispatch_outside_setup_test.go` and `internal/baseline/context7_capability_test.go` with the eight tests of the TechSpec's Testing Approach 3. The dispatch negatives build catalogs with `cloneEmbeddedAssets`: one removes a skill the `go` module dispatches from `go-tui.json`, one adds an unknown skill to `bundle.qa`, and one confirms that `go-cli-tui`, whose setup does not list `vitest` and which selects no TypeScript module, gets no dispatch diagnostic for it. `TestEveryRestorableCapabilitySkillIsRequiredAndInEverySetup` reads `universalCapabilityRestoreSkills` and requires each skill to equal its capability's probe skill, to be in `core`'s `requiredSkills`, and to be a `github` member of every built-in profile's setup.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`; a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a catalog snapshot or a rendered guide.
7. MUST NOT change the capability identifier, its Title or Explanation, any other capability, or any existing test.

## Subtasks

- [ ] Add the dispatch check to the profile validation.
- [ ] Require and dispatch `exa-web-search` in `core`.
- [ ] Move the Context7 probe and remediation to `context7-cli` with the prior name.
- [ ] Update the user guide's restoration example.
- [ ] Regenerate, refresh this repository's guides, and add the tests.

## Acceptance Criteria

- [ ] The embedded catalog loads, and a profile whose modules or owned bundles name a skill outside its setup is reported.
- [ ] A repository with `context7-cli` or with `context7` satisfies `capability.context7`; one with neither is blocked with a next action naming `--skill context7-cli`.
- [ ] Every restorable capability skill is required by `core` and listed by every built-in profile's setup.
- [ ] `docs/agents/skill-dispatch.md` carries `trigger.core.exa-web-search`, and a second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/catalog_dispatch_outside_setup_test.go`
- creates: `internal/baseline/context7_capability_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestADispatchedSkillOutsideTheProfileSetupIsReported|TestABundledSkillOutsideTheProfileSetupIsReported|TestAModuleTheProfileDoesNotSelectIsNotChecked|TestTheContext7CapabilityIsSatisfiedByTheCurrentSkill|TestTheContext7CapabilityIsSatisfiedByItsPriorSkill|TestAMissingContext7SkillIsRemediatedWithTheCurrentName|TestEveryRestorableCapabilitySkillIsRequiredAndInEverySetup|TestCatalogDiagnosticCharacterization|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestADispatchedSkillOutsideTheProfileSetupIsReported TestABundledSkillOutsideTheProfileSetupIsReported TestAModuleTheProfileDoesNotSelectIsNotChecked TestTheContext7CapabilityIsSatisfiedByTheCurrentSkill TestTheContext7CapabilityIsSatisfiedByItsPriorSkill TestAMissingContext7SkillIsRemediatedWithTheCurrentName TestEveryRestorableCapabilitySkillIsRequiredAndInEverySetup TestCatalogDiagnosticCharacterization TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/skill-dispatch.md|trigger.core.exa-web-search" "docs/user-guide/context-driven-development.md|--skill context7-cli --confirm-plan" "docs/user-guide/context-driven-development.md|still satisfies the Context7 capability"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the eight new tests do not exist and the guide has no `exa-web-search` trigger, so the command fails.

## References

- `_prd.md` → Goals 2-3; User Stories 2-3; Core Features 3-4; Success Metrics 3-4; Declared breaks (catalog, capability remediation, required skills)
- `_techspec.md` → Interfaces; Fixed texts; API Contracts 2, 3 and 4; Testing Approach 3 and 5; Build Order 3
- ADR-0081, ADR-0103, ADR-0149, ADR-0191

## Result
