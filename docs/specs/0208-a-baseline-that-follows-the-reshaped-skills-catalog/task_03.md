---
task: task_03
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: pending
type: backend
complexity: medium
---

# Task 03: Every managed repository requires the TypeSafe and README skills

## Overview

Every upstream setup lists `typesafe-ai` and `crafting-effective-readmes`, but no module requires `typesafe-ai` and only the `typescript` module requires the README skill. This Task makes `core` require and dispatch both, removes the README skill from `typescript`, and restores both skills in this repository with the public update. The repository Verification checks this repository's skill set when the Task settles, so the restore belongs to this Task.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file. The `.agents/skills/` files it may create are exactly the sixteen listed there.

## Requirements

1. MUST apply the task_03 Fixed texts to `core.json` and `typescript.json`, with the exact trigger identifiers and `when` texts, and raise both module versions by one from the Task's base.
2. MUST create `internal/baseline/core_skill_requirements_test.go` with the two tests of Testing Approach 3, built on `coreSkillFindings` as the TechSpec's Interfaces describe.
3. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`.
4. MUST run step 1 of the TechSpec's repository procedure with a clone of the local `~/dev/skills` at `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, and MUST stop without cloning from the network when the checkout lacks the commit. MUST NOT copy a skill tree by hand, write `~/dev/skills`, or change a `.agents/skills/` file outside the sixteen the Context lists.
5. MUST make `skills/recommended.txt` the sorted keys of `skills-lock.json`, and set `upstreamManagedSkillTreeDigest` in `skills/baseline_skill_contract_test.go` to the value `TestAuthorialSkillSync` reports. No other line of that file changes.
6. MUST replace the Doctor line in `README.md` and `docs/user-guide/usage.md` with `skills: ok (43 required: 14 Roundfix-owned, 29 external)`, after `go run -buildvcs=false ./cmd/roundfix doctor` prints it in this repository. When Doctor prints other numbers, MUST stop and report them.

## Subtasks

- [ ] Require and dispatch the two skills in `core`, and drop the README skill from `typescript`.
- [ ] Write the two tests.
- [ ] Regenerate, then refresh this repository's guides twice.
- [ ] Restore the two skills with the public update, and align the recommended list and the digest pin.
- [ ] Correct the Doctor line.

## Acceptance Criteria

- [ ] `core` requires and dispatches `typesafe-ai` and `crafting-effective-readmes` with the exact triggers, and no `typescript` entry names the README skill.
- [ ] The Standard TypeScript Monorepo golden and this repository's `docs/agents/skill-dispatch.md` name both `core` triggers, and a missing trigger is reported by the helper.
- [ ] Every built-in profile's setup lists both skills, as the existing dispatch check proves.
- [ ] This repository holds every required external skill with its lock hash, and the recommended list and the digest pin match the lock.
- [ ] The guides print the Doctor line this repository prints.

## Context

- instruction: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- creates: `internal/baseline/core_skill_requirements_test.go`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`
- interface: `skills-lock.json`
- interface: `skills/recommended.txt`
- interface: `skills/baseline_skill_contract_test.go`
- creates: `.agents/skills/typesafe-ai/LICENSE`
- creates: `.agents/skills/typesafe-ai/SKILL.md`
- creates: `.agents/skills/crafting-effective-readmes/README.md`
- creates: `.agents/skills/crafting-effective-readmes/SKILL.md`
- creates: `.agents/skills/crafting-effective-readmes/references/art-of-readme.md`
- creates: `.agents/skills/crafting-effective-readmes/references/make-a-readme.md`
- creates: `.agents/skills/crafting-effective-readmes/references/standard-readme-example-maximal.md`
- creates: `.agents/skills/crafting-effective-readmes/references/standard-readme-example-minimal.md`
- creates: `.agents/skills/crafting-effective-readmes/references/standard-readme-spec.md`
- creates: `.agents/skills/crafting-effective-readmes/section-checklist.md`
- creates: `.agents/skills/crafting-effective-readmes/style-guide.md`
- creates: `.agents/skills/crafting-effective-readmes/templates/internal.md`
- creates: `.agents/skills/crafting-effective-readmes/templates/oss.md`
- creates: `.agents/skills/crafting-effective-readmes/templates/personal.md`
- creates: `.agents/skills/crafting-effective-readmes/templates/xdg-config.md`
- creates: `.agents/skills/crafting-effective-readmes/using-references.md`
- interface: `README.md`
- interface: `docs/user-guide/usage.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCoreRequiresAndDispatchesTheTypeSafeAndReadmeSkills|TestAMissingCoreSkillTriggerIsReported|TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestEveryModuleRequiresEverySkillItDispatches|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestCatalogDiagnosticCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCoreRequiresAndDispatchesTheTypeSafeAndReadmeSkills TestAMissingCoreSkillTriggerIsReported TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestEveryModuleRequiresEverySkillItDispatches TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus TestCatalogDiagnosticCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestThisRepositoryHoldsEveryRequiredExternalSkill|TestARepositoryMissingARequiredExternalSkillIsReported)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryHoldsEveryRequiredExternalSkill TestARepositoryMissingARequiredExternalSkillIsReported; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestRecommendedSkillsMatchLock|TestAuthorialSkillSync|TestAuthoringConstraintOwnership|TestUpstreamADRFormatUnchanged)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecommendedSkillsMatchLock TestAuthorialSkillSync TestAuthoringConstraintOwnership TestUpstreamADRFormatUnchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && test -f .agents/skills/typesafe-ai/SKILL.md && test -f .agents/skills/crafting-effective-readmes/SKILL.md && for file in README.md docs/user-guide/usage.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'skills: ok (43 required: 14 Roundfix-owned, 29 external)' || { printf 'missing Doctor line in %s\n' "$file" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the two new tests do not exist and `typesafe-ai` is not installed, so the command fails.

## References

- `_prd.md` → Goals 4-5; User Story 4; Core Features 4-5; Success Metrics 4-5
- `_techspec.md` → Fixed texts (task_03); The repository procedure; API Contracts 4 and 5; Testing Approach 3, 4 and 5; Build Order 3
- ADR-0206, ADR-0191, ADR-0073, ADR-0103
