---
task: task_03
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: completed
type: backend
complexity: high
---

# Task 03: A built-in profile for a Go CLI with a TypeScript monorepo

## Overview

This Task ships the built-in `go-cli-typescript-monorepo` Baseline Profile on the composed setup task_02 created, and gives the Go guide the scope sentence ADR-0190 requires, so its rules do not read as the TypeScript workspace's. It proves that the profile passes every catalog check, that a plan renders every guide with no clause repeated, and that the plan converges. Because the composed profile selects the modules of the Standard TypeScript Monorepo and Go CLI/TUI profiles, a repository-owned Profile draft cut from either now adapts two built-in profiles; this Task binds such a draft to the closest one, as ADR-0219 decides, and moves the test pins the new profile and the Go guide's new bytes reach.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST create `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` exactly as the TechSpec's Data Models describe the composed profile, including the `partOfGate` markers, which task_04 validates and reads. Every list copied from the Standard TypeScript Monorepo Profile MUST equal that profile's on this Task's base.
2. MUST insert the Go guide's scope paragraph the TechSpec's Fixed texts give into `internal/baseline/assets/templates/guides/go.md`, and raise the versions of `template.guide.go` in `templates/index.json` and of `guide.go` and the module in `internal/baseline/assets/modules/go.json` by one from this Task's base. No Go rule changes.
3. MUST add `go-cli-typescript-monorepo` to the built-in profile list `internal/baseline/catalog_test.go` pins, and give it a case in `TestProjectDecisionAssets` that expects `identifier.strategy` and `auth.provider`.
4. MUST create `internal/baseline/composed_profile_test.go` with the five tests the TechSpec's Testing Approach 3 names. The render test MUST plan the profile in a temporary repository built from `newAlignedTypeScriptRepository` plus a `go.mod`, collect every rendered clause line across the postimages, and fail on a repeat; the negative test MUST feed the same check a guide set with one clause line duplicated and expect that line back; the convergence test MUST apply the plan and see a second plan report no file change.
5. MUST update the Profiles paragraph and the `--profile-file` paragraph of `docs/user-guide/context-driven-development.md` as Fixed texts gives for task_03. Every scripted interview in `internal/cli/baseline_human_test.go` that selects a profile by its position in the list MUST select the same profile after the new one is listed: `TestBaselineHumanProfileAdaptation` (Standard TypeScript Monorepo, `3` to `4`) and `TestBaselineHumanProfileChangeRemainsReachable` (rust-cli, the sixth answer `2` to `3`). Nothing else in that file changes.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. This repository's Go guide gains the scope sentence through that refresh and is never hand-edited.
7. MUST NOT change the composed setup, any other profile, or any module other than `go.json`, and MUST NOT rename or remove a top-level test or an exported function.
8. MUST change `ProfileDraftInputFromDocument` in `internal/baseline/custom_profile.go` so that it collects every compatible built-in profile as today and then keeps only those `closestProfileDraftSources` returns, as the TechSpec's Interfaces give: one closest source binds the draft; several equally close sources keep `custom.profile.draft.source.ambiguous` and name only those sources; no compatible source keeps `custom.profile.draft.source.unresolved`. No other function in that file changes.
9. MUST add to `internal/baseline/custom_profile_test.go` the two tests the TechSpec's Testing Approach 3 names for the draft source. `TestADraftThatAdaptsSeveralBuiltInProfilesBindsToTheClosest` MUST build each draft with `NewProfileAdaptationDraft` on the embedded catalog and require it to bind to the profile it was cut from: the Standard TypeScript Monorepo Profile without `frontend` and `autonomous-work` and without `capability.workspace.frontend`; the Go CLI/TUI profile without `tui-surface` and `autonomous-work`; and the composed profile without `cli-surface` and `autonomous-work`. `TestEquallyCloseProfileDraftSourcesStayAmbiguous` MUST give `closestProfileDraftSources` two sources at equal distance and one farther, and require exactly the two.
10. MUST move only the pins the new profile and the Go guide's new bytes reach: add `docs/agents/go.md` to `evolvedPastFrozenCorpus` in `TestPlanDeterminismMatchesMaintainedManagedEntryFixture` (`internal/baseline/plan_test.go`); append `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` to the additions `TestOutputsForCommand` makes to the 2026-08-06 enumeration (`internal/baseline/derived_ownership_test.go`); and list `go-cli-typescript-monorepo` in the maintained profiles of `TestGuidanceCompositionJourney` (`internal/cli/baseline_release_gate_test.go`), giving its journey the Standard TypeScript Monorepo fixture plus `go.mod` and `main.go`, the same three TypeScript decisions, and `gofmt -w main.go` as its formatter step. The frozen parity corpus, the historical grant and every other assertion stay byte-identical.

## Subtasks

- [ ] Write the composed profile.
- [ ] Add the Go guide's scope sentence and raise its versions.
- [ ] Name the profile in the catalog tests and write the five new tests.
- [ ] Bind a draft that adapts several built-in profiles to the closest one, with its two tests.
- [ ] Move the parity, enumeration, release-journey and interview pins the new profile reaches.
- [ ] Update the public guide, regenerate, and refresh this repository twice.

## Acceptance Criteria

- [ ] The embedded catalog loads with `go-cli-typescript-monorepo` on `go-cli-typescript-bun`, and the check that a profile's setup lists every skill its guides name passes for it.
- [ ] A plan for the composed profile renders the Go, CLI, TypeScript and Bun, monorepo, backend and frontend guides with no clause line repeated, and a duplicated line is reported.
- [ ] The Go guide opens with its scope sentence in both this repository and the composed plan.
- [ ] The composed plan converges, and a second Managed Refresh of this repository is a no-op.
- [ ] A Profile draft cut from the Standard TypeScript Monorepo, Go CLI/TUI or composed profile binds to the profile it was cut from, the automation path of the human adaptation reproduces the human plan, and a tie is still refused.
- [ ] The composed profile runs the maintained release journey, and `make verify` passes.

## Context

- instruction: `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- instruction: `docs/adr/0219-a-profile-draft-binds-to-the-closest-built-in-profile-it-adapts.md`
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
- interface: `internal/baseline/custom_profile.go`
- interface: `internal/baseline/custom_profile_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheComposedProfileTakesTheComposedSetup|TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause|TestTheGoGuideNamesWhatItGoverns|TestARepeatedRenderedClauseIsReported|TestTheComposedProfilePlanConverges|TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestProjectDecisionAssets|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestADraftThatAdaptsSeveralBuiltInProfilesBindsToTheClosest|TestEquallyCloseProfileDraftSourcesStayAmbiguous|TestPlanDeterminismMatchesMaintainedManagedEntryFixture|TestOutputsForCommand)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheComposedProfileTakesTheComposedSetup TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause TestTheGoGuideNamesWhatItGoverns TestARepeatedRenderedClauseIsReported TestTheComposedProfilePlanConverges TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestProjectDecisionAssets TestCatalogCompatibility TestBaselinePlanCharacterization TestADraftThatAdaptsSeveralBuiltInProfilesBindsToTheClosest TestEquallyCloseProfileDraftSourcesStayAmbiguous TestPlanDeterminismMatchesMaintainedManagedEntryFixture TestOutputsForCommand; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && cli="$(go test -count=1 -v -run "^(TestBaselineHumanProfileAdaptation|TestBaselineHumanProfileChangeRemainsReachable|TestBaselinePlanProfileFile|TestProfileAdaptationJourney|TestGuidanceCompositionJourney)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$cli"; exit 1; }; for name in TestBaselineHumanProfileAdaptation TestBaselineHumanProfileChangeRemainsReachable TestBaselinePlanProfileFile TestProfileAdaptationJourney TestGuidanceCompositionJourney/go-cli-typescript-monorepo; do printf "%s\\n" "$cli" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/agents/go.md | grep -qF "These rules govern the repository's Go module: its commands, packages and tests." && tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF "binds to the closest one" && go run -buildvcs=false ./cmd/roundfix baseline profile validate go-cli-typescript-monorepo --format text >/dev/null && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the five composed-profile tests, the two draft-source tests, the profile, the Go guide's sentence and the guide's draft-source sentence do not exist, and on a tree without Spec 0200 the dispatch check test does not exist, so the command fails; on a tree with the profile but without the closest-source rule and the moved pins, the parity, enumeration, profile-file, adaptation, interview and release-journey tests fail.

## References

- [_techspec.md](_techspec.md) — Interfaces; Data Models; Fixed texts; Testing Approach 3; Testing Approach 5; Build Order 3; Decisions
- `_prd.md` → Goal 1; Story 1; Story 3; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 3; Success Metric 7; Prerequisites
- `_techspec.md` → API Contract 5
- ADR-0059, ADR-0061, ADR-0063, ADR-0075, ADR-0081, ADR-0149, ADR-0186, ADR-0190, ADR-0191, ADR-0204, ADR-0219

## Result

Implemented task_03's bounded slice for Daemon Verification. Task status and
checkboxes remain unchanged; no commit, push, or Pull Request was made.

The new built-in Profile combines the prescribed Go CLI and TypeScript modules
on `go-cli-typescript-bun`, copies the Standard TypeScript Monorepo lists from
this Task's base, and declares the two `partOfGate` entries. The Go template
opens with the exact scope paragraph; its template, guide and module versions
rose by one, with every Go rule unchanged. Profile drafts now retain only the
closest compatible sources by removed modules, then removed capabilities;
unresolved and equally close refusals keep their existing error codes.

The catalog list, project-decision case, maintained release journey, two
interview answers, parity exception and derived-output enumeration include the
new Profile. The user guide explains the composition and closest-source rule.
All derived pins were regenerated by `make baseline-digests`; this repository's
Go guide and Setup Manifest were refreshed through the public command.

### Focused check evidence

All Go commands below used `GOCACHE=/tmp/roundfix-task03-gocache` and RTK's exact
output proxy unless stated otherwise.

- Initial signal: the new draft-source tests failed to build with
  `undefined: closestProfileDraftSources`; the composed Profile and Go scope
  paragraph were absent on the base. The first default-cache attempt was
  denied by the sandbox, so subsequent checks used the temporary cache.
- `go test ./internal/baseline -count=1 -run 'Test(TheComposedProfile|TheGoGuide|ARepeatedRendered|ADraftThatAdapts|EquallyClose|EveryBuiltInProfileSetup|ProjectDecisionAssets|CatalogCompatibility|BaselinePlanCharacterization|PlanDeterminismMatchesMaintained|OutputsForCommand)'`
  exited 0 (`ok roundfix/internal/baseline`, 8.414s).
- `go test ./internal/cli -count=1 -run 'TestBaselineHumanProfile|TestBaselinePlanProfileFile|TestProfileAdaptationJourney|TestGuidanceCompositionJourney'`
  exited 0 (`ok roundfix/internal/cli`, 10.740s).
- `rtk make baseline-digests` exited 0 and reported deterministic regeneration
  of the catalog diagnostics, digest, normalized catalog and four affected plan
  characterization goldens. The unsatisfied-capabilities golden needed no change.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  ran twice successfully: first `File changes: 2`, then `File changes: 0` and
  `Idempotence: verified`. The initial sandbox attempt could not create the
  Git-private transaction directory; the authorized elevated reruns succeeded.
- A Python comparison of the new Profile against the base TypeScript Profile
  and TechSpec Data Models confirmed the copied lists, exact module/rule/stack
  composition, capability order and Verification entries.
- `git -c core.fsmonitor=false diff --check` exited 0. A postflight comparing
  tracked and untracked changes with this Task's Context found all 23 changed
  paths authorized, confirmed `status: in_progress`, and confirmed the rendered
  repository Go scope sentence after whitespace normalization.
- `rtk make verify-incremental` exited 0 on an elevated rerun with the
  worktree held unchanged throughout: formatting, `go vet ./...`,
  `go test -parallel 16 ./...`, skill sync/check and CLI build all passed.
  The initial run exited 2 because this Agent refreshed guides and wrote this
  Result while suiteguard fingerprinted the repository, and sandbox process-table
  restrictions denied two CLI force-stop integration tests. The serial elevated
  rerun removed those execution conditions without changing or bypassing tests;
  baseline, CLI, daemon and every other package passed.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Catalog and guide-skill coverage | The Profile selects `go-cli-typescript-bun`; embedded catalog loading, catalog compatibility, and every built-in Profile's guide-skill coverage check passed in the focused baseline run. |
| All composed guides and duplicate detection | The render test plans a temporary aligned TypeScript repository plus `go.mod`, requires all six guide paths, and scans bullet clause lines across every Markdown postimage. It passed; the negative test returns the duplicated line for both enforcement-labelled and plain Go-style clauses. |
| Go guide scope | The composed-plan scope test passed; the generated repository guide was inspected and its exact normalized scope sentence confirmed after public refresh. |
| Plan and refresh convergence | The convergence test applies the composed plan and requires zero changes on replan; it passed. The second repository refresh reported `File changes: 0`. |
| Closest draft source and automation parity | The three requested adaptation drafts bound to their original Profiles; the tie test retained exactly the two closest sources, excluding farther module and capability candidates. The profile-file, human adaptation and adaptation journey checks passed, including reproduction of the human plan through automation. |
| Maintained release journey and repository gate | The CLI focused run exercised the new composed journey with the TypeScript fixture, Go files, the same three TypeScript decisions and `gofmt -w main.go`. The selected incremental gate passed all its targets; the separate full `make verify` command was not run. Daemon Verification remains authoritative. |

No follow-up implementation was added. The Task's authored `## Verification`
command was not run, and no Daemon Verification verdict is claimed.
