---
task: task_03
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: pending
type: backend
complexity: high
---

# Task 03: A built-in profile for a Go CLI with a TypeScript monorepo

## Overview

This Task ships the built-in `go-cli-typescript-monorepo` Baseline Profile on the composed setup task_02 created, and gives the Go guide the scope sentence ADR-0190 requires, so its rules do not read as the TypeScript workspace's. It proves that the profile passes every catalog check, that a plan renders every guide with no clause repeated, and that the plan converges.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST create `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` exactly as the TechSpec's Data Models describe the composed profile, including the `partOfGate` markers, which task_04 validates and reads. Every list copied from the Standard TypeScript Monorepo Profile MUST equal that profile's on this Task's base.
2. MUST insert the Go guide's scope paragraph the TechSpec's Fixed texts give into `internal/baseline/assets/templates/guides/go.md`, and raise the versions of `template.guide.go` in `templates/index.json` and of `guide.go` and the module in `internal/baseline/assets/modules/go.json` by one from this Task's base. No Go rule changes.
3. MUST add `go-cli-typescript-monorepo` to the built-in profile list `internal/baseline/catalog_test.go` pins, and give it a case in `TestProjectDecisionAssets` that expects `identifier.strategy` and `auth.provider`.
4. MUST create `internal/baseline/composed_profile_test.go` with the five tests the TechSpec's Testing Approach 3 names. The render test MUST plan the profile in a temporary repository built from `newAlignedTypeScriptRepository` plus a `go.mod`, collect every rendered clause line across the postimages, and fail on a repeat; the negative test MUST feed the same check a guide set with one clause line duplicated and expect that line back; the convergence test MUST apply the plan and see a second plan report no file change.
5. MUST update the Profiles paragraph of `docs/user-guide/context-driven-development.md` as Fixed texts gives for task_03. If a scripted interview in `internal/cli/baseline_human_test.go` selects a profile by its position in the list, it MUST select the same profile after the new one is listed, and nothing else in that file changes.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. This repository's Go guide gains the scope sentence through that refresh and is never hand-edited.
7. MUST NOT change the composed setup, any other profile, or any module other than `go.json`, and MUST NOT rename or remove a top-level test or an exported function.

## Subtasks

- [ ] Write the composed profile.
- [ ] Add the Go guide's scope sentence and raise its versions.
- [ ] Name the profile in the catalog tests and write the five new tests.
- [ ] Update the public guide, regenerate, and refresh this repository twice.

## Acceptance Criteria

- [ ] The embedded catalog loads with `go-cli-typescript-monorepo` on `go-cli-typescript-bun`, and the check that a profile's setup lists every skill its guides name passes for it.
- [ ] A plan for the composed profile renders the Go, CLI, TypeScript and Bun, monorepo, backend and frontend guides with no clause line repeated, and a duplicated line is reported.
- [ ] The Go guide opens with its scope sentence in both this repository and the composed plan.
- [ ] The composed plan converges, and a second Managed Refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- creates: `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`
- interface: `internal/baseline/assets/templates/guides/go.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/modules/go.json`
- interface: `internal/baseline/catalog_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/agents/go.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- creates: `internal/baseline/composed_profile_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheComposedProfileTakesTheComposedSetup|TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause|TestTheGoGuideNamesWhatItGoverns|TestARepeatedRenderedClauseIsReported|TestTheComposedProfilePlanConverges|TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestProjectDecisionAssets|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheComposedProfileTakesTheComposedSetup TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause TestTheGoGuideNamesWhatItGoverns TestARepeatedRenderedClauseIsReported TestTheComposedProfilePlanConverges TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestProjectDecisionAssets TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/agents/go.md | grep -qF "These rules govern the repository's Go module: its commands, packages and tests." && go run -buildvcs=false ./cmd/roundfix baseline profile validate go-cli-typescript-monorepo --format text >/dev/null && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the five new tests, the profile and the Go guide's sentence do not exist, and on a tree without Spec 0200 the dispatch check test does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Data Models; Fixed texts; Testing Approach 3; Testing Approach 5; Build Order 3
- `_prd.md` → Goal 1; Story 1; Story 3; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 3; Success Metric 7; Prerequisites
- `_techspec.md` → API Contract 5
- ADR-0059, ADR-0061, ADR-0063, ADR-0081, ADR-0149, ADR-0186, ADR-0190, ADR-0191, ADR-0204
