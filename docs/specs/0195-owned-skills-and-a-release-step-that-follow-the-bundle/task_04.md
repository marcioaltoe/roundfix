---
task: task_04
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: pending
type: docs
complexity: medium
---

# Task 04: Every release checks skills and guides first

## Overview

Skills and guides drift because no step of a release re-reads them. This Task adds a mandatory step to the release runbook, run before the release Pull Request, that names the checks to run and ends with a reading pass. It appends one sentence to the Baseline's release clause so every adopter's Agent knows the step exists, and it adds a test that refuses a step naming a check the repository does not have.

The checks the step names belong to Spec 0192, Spec 0193 and this Spec's task_01 and task_02. The runbook names them; it does not restate what they assert.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add the section "Checking skills and guides before the release" to `docs/user-guide/release-runbook.md`, before "Cutting a release", with the six items the TechSpec's "Fixed texts" gives, each command written exactly as given. The section MUST say the step is mandatory and runs before the release Pull Request.
2. MUST add one step to "Cutting a release", after the release plan step and before the tag step, that points at the new section. Every existing step keeps its text.
3. MUST append to `clause.core.plan-the-release-first` in `internal/baseline/assets/modules/core.json` the sentence the TechSpec gives, keep the clause `mandatory`, keep every existing sentence, and raise by one the versions of `rule.core.git-delivery`, `guide.agent-instructions` and the `core` module. Replace the string and the numbers in place.
4. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a catalog snapshot or a generated guide.
5. MUST add `internal/docscontract/release_step_test.go`, in the package and under the `docscontract` build tag the directory's other tests use, with the three tests the TechSpec's Testing Approach 4 names. `TestEveryCheckTheReleaseStepNamesExists` MUST read the names from the runbook section, not from a list in the test, and MUST fail when the section names fewer than seven tests.
6. MUST add `internal/baseline/release_clause_test.go` with `TestTheReleaseClauseNamesTheSkillsAndGuidesCheck`, which expects the sentence in the embedded clause and in the formatter golden.
7. MUST NOT edit `internal/docscontract/publicdocs_test.go`, any skill, the Release Plan Command, or any existing test. The sentence cites no Spec number and no ADR number.
8. MUST re-record `skills/testdata/owned-skill-versions.json` with the sanctioned command `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, because Spec 0192 landed after task_02 and raised owned skill versions (for example `archive-spec` to `0.0.3`). The record MUST change only by that command, and no skill may change.

## Subtasks

- [ ] Add the runbook section and the pointer step.
- [ ] Append the clause sentence and raise the versions.
- [ ] Regenerate, refresh this repository's guides, and confirm the second refresh is a no-op.
- [ ] Add the runbook tests and the clause test, the negative case separate.

## Acceptance Criteria

- [ ] The runbook carries the section with its six items, and "Cutting a release" points at it before the tag step.
- [ ] Every test the section names exists in the repository, and a section that names a missing test is reported.
- [ ] `docs/agents/agent-instructions.md` carries the appended sentence, and the existing release-planning contract test passes.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `internal/docscontract/publicdocs_test.go`
- interface: `docs/user-guide/release-runbook.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `internal/docscontract/release_step_test.go`
- creates: `internal/baseline/release_clause_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestTheReleaseRunbookRequiresTheSkillsAndGuidesCheck|TestEveryCheckTheReleaseStepNamesExists|TestAReleaseStepThatNamesAMissingCheckIsReported|TestReleasePlanDocumentationContract|TestTheReleaseClauseNamesTheSkillsAndGuidesCheck|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/docscontract ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheReleaseRunbookRequiresTheSkillsAndGuidesCheck TestEveryCheckTheReleaseStepNamesExists TestAReleaseStepThatNamesAMissingCheckIsReported TestReleasePlanDocumentationContract TestTheReleaseClauseNamesTheSkillsAndGuidesCheck TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/user-guide/release-runbook.md|Checking skills and guides before the release" "docs/user-guide/release-runbook.md|TestEveryOwnedSkillVersionIsRecorded" "docs/agents/agent-instructions.md|describe the behavior being released"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the version record lacks the versions Spec 0192 raised, and none of the four new named tests exists and the runbook has no such section, so the command fails.

## References

- [_techspec.md](_techspec.md) — Fixed texts; Testing Approach 4
- `_prd.md` → Goal 4; Core Feature 4; Success Metrics 5 and 6
- `_techspec.md` → API Contract 6; Testing Approach 5
- [references/2026-09-30-a-release-does-not-check-skills-and-guides.md](references/2026-09-30-a-release-does-not-check-skills-and-guides.md)
- ADR-0143, ADR-0186, ADR-0189

## Result
