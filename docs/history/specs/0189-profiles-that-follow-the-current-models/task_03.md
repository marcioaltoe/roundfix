---
task: task_03
spec: 0189-profiles-that-follow-the-current-models
status: completed
type: backend
complexity: high
---

# Task 03: Built-in selections are the Recommended Profile, and the replaced models leave production code

## Overview

`builtinProfiles` in `internal/config/profiles.go`, the generated config in `internal/config/config.go`, the legacy Codex runtime default and the Baseline semantic analysis in `internal/baselineacp/analyzer.go` each name their own models. All four name `gpt-5.5`, which leaves Codex on 2026-10-14. Preflight proves every configured tuple, so from that day a repository on built-in profiles could not start a Run. This Task derives all of them from the Recommended Profile of task_02 (ADR-0180). It also removes the two transitional entries task_01 left in the Model Catalog, `gpt-5.5` and `claude-fable-5`, which nothing names once the built-ins change. It is verifiable on its own: with no configuration, `roundfix profiles show` reports the Recommended Profile as `built-in` for each required category.

## Requirements

1. MUST make `builtinProfiles` build each of the five required categories from `RecommendedProfile`, with source `built-in`, and define no optional category.
2. MUST make `defaultConfigYAML` render its `profiles:` block from `builtinProfiles`, in required-category order, for both the user and the project scope. The rendered block MUST parse back to the same five profiles.
3. MUST make the built-in Codex runtime default the Codex selection of the `general` Recommended Profile, its model and its effort, and remove the constants `defaultCodexModel` and `defaultCodexReasoningEffort`. The built-in Claude runtime default MUST not change, and `applyLegacyRuntimeProfiles` MUST keep its behavior.
4. MUST set `PreferredModel` to `gpt-6.1-sol` and `FallbackModel` to `gpt-5.6-sol` in `internal/baselineacp/analyzer.go`, and keep `RequiredReasoningEffort` at `xhigh` and every other part of the analysis unchanged.
5. MUST remove `gpt-5.5` and `claude-fable-5` from the Model Catalog, so that each catalog is exactly the table of `_techspec.md` → Reference data → Model Catalog, and MUST leave no occurrence of `gpt-5.5` in any non-test Go file under `internal` or `cmd`.
6. MUST edit `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md` so that its two model names are `gpt-6.1-sol` and `gpt-5.6-sol`, that it says ADR-0180 sets them, and that its `updated_at` is `2026-09-30T00:00:00Z`. Nothing else in that ADR may change.
7. MUST update only the existing tests this change invalidates, and name each in the Result. The TechSpec's Risks section lists the twenty measured on a scratch copy. A test that means "the built-in value" MUST read `RecommendedProfile` or `Builtin()` instead of repeating a literal. The two Doctor characterization tests change because `_prd.md` → Declared breaks declares the break; their new expected counts MUST be stated in the Result. It MUST NOT edit a Baseline fixture or golden file: their `codex gpt-5.5 xhigh` is a decision value, not a built-in.
8. MUST put the new tests in `internal/config/built_in_profiles_test.go`, `internal/cli/built_in_review_provider_test.go`, `internal/baselineacp/analyzer_models_test.go` and `internal/agent/catalog_replaced_test.go`. The catalog test MUST require that neither catalog offers `gpt-5.5` or `claude-fable-5`, and the existing exact-list test in `internal/agent/agent_test.go` MUST be updated to the final table. The review test MUST prove that the built-in `review` profile passes `validateReviewProfileProvider` for the provider `codex`.
9. MUST prove each new gate can fail. The Result MUST record one sabotage for the built-in equality test (a built-in edited away from the Recommended Profile) and one for the review provider test (the built-in `review` derived from `general`), each with the test that failed, and that the code was restored.
10. MUST update, in place, every sentence and every example configuration of `docs/user-guide/configuration.md`, `docs/user-guide/usage.md` and `.agents/skills/roundfix/SKILL.md` that states the built-in profiles, the generated config, the legacy runtime default or the catalog identifiers, so that none names `gpt-5.5` or offers `claude-fable-5`. The example configurations MUST equal what `roundfix setup` now generates. It MUST add no text inside the `### QA settlement` section of the skill, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
11. MUST remove the two `gpt-5.5` pins from `TestProfilesDocumentationContractMatchesPublicGuidance` in `internal/docscontract/publicdocs_test.go`, change its `claude-fable-5` pin to `claude-fable-5-1`, and pin the built-in Preferred Selection through `RecommendedProfile` instead.

## Subtasks

- [ ] Derive the built-in profiles, the generated config and the legacy Codex default from the Recommended Profile.
- [ ] Remove the two transitional entries from the Model Catalog.
- [ ] Replace the Baseline analysis models and amend ADR-0069.
- [ ] Update the invalidated tests and add the new ones, each negative case separate.
- [ ] Update the guides, the skill and the documentation contract pins, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] With no User Config and no Project Config, each required category resolves to its Recommended Profile with source `built-in`, and an optional category resolves to `general` by inheritance.
- [ ] The config `DefaultConfigYAML` and `DefaultProjectConfigYAML` generate parses back to the five built-in profiles.
- [ ] A legacy config that names only `defaults.agent: codex` resolves to `gpt-6.1-sol` with `high`.
- [ ] The built-in `review` profile passes the review provider check for `codex`.
- [ ] Both Baseline analysis models are in the Codex Model Catalog.
- [ ] Each catalog is exactly the Reference data table, without `gpt-5.5` and without `claude-fable-5`.
- [ ] No non-test Go file under `internal` or `cmd` contains `gpt-5.5`, and neither do the two guides, the skill, its mirror or ADR-0069. None of those documents offers `claude-fable-5`.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`
- interface: `internal/baselineacp/analyzer.go`
- interface: `internal/agent/catalog.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/config/config_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/doctor_characterization_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/selection_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- creates: `internal/config/built_in_profiles_test.go`
- creates: `internal/cli/built_in_review_provider_test.go`
- creates: `internal/baselineacp/analyzer_models_test.go`
- creates: `internal/agent/catalog_replaced_test.go`
- instruction: `internal/config/recommendations.go`
- instruction: `internal/cli/review.go`
- interface: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBuiltinProfilesAreTheRecommendedProfile|TestGeneratedConfigRendersTheBuiltinProfiles|TestLegacyCodexDefaultIsTheGeneralCodexSelection|TestBuiltinProfilesDefineNoOptionalCategory|TestBuiltinReviewProfileStaysOnTheReviewProvidersRuntime|TestAnalysisModelsAreInTheCodexModelCatalog|TestModelCatalogOffersNoReplacedModel|TestModelCatalogsExposeOrderedPickerData|TestModelCatalogOpensWithTheCurrentModels)$" ./internal/config ./internal/cli ./internal/baselineacp ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBuiltinProfilesAreTheRecommendedProfile TestGeneratedConfigRendersTheBuiltinProfiles TestLegacyCodexDefaultIsTheGeneralCodexSelection TestBuiltinProfilesDefineNoOptionalCategory TestBuiltinReviewProfileStaysOnTheReviewProvidersRuntime TestAnalysisModelsAreInTheCodexModelCatalog TestModelCatalogOffersNoReplacedModel TestModelCatalogsExposeOrderedPickerData TestModelCatalogOpensWithTheCurrentModels; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && test -d internal && ! grep -rq --include='*.go' --exclude='*_test.go' -e 'gpt-5\.5' internal cmd` — expected: exit 0; before this Task none of the seven new named tests exists and production files name `gpt-5.5`, so the command fails.
- `for pair in "docs/user-guide/configuration.md|codex / gpt-6.1-sol / high" ".agents/skills/roundfix/SKILL.md|codex / gpt-6.1-sol / high" "skills/roundfix/SKILL.md|codex / gpt-6.1-sol / high" "docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md|ADR-0180"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && ! grep -q 'gpt-5\.5' docs/user-guide/configuration.md && ! grep -q 'gpt-5\.5' docs/user-guide/usage.md && ! grep -q 'gpt-5\.5' .agents/skills/roundfix/SKILL.md && ! grep -q 'gpt-5\.5' skills/roundfix/SKILL.md && ! grep -q 'gpt-5\.5' docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md && ! grep -qE 'claude-fable-5([^-]|$)' docs/user-guide/configuration.md docs/user-guide/usage.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check && out="$(go test -count=1 -tags docscontract -v -run "^(TestProfilesDocumentationContractMatchesPublicGuidance)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfilesDocumentationContractMatchesPublicGuidance; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the guides and ADR-0069 name `gpt-5.5`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Core Features 4 and 5; Success Metric 1; Declared breaks
- [_techspec.md](_techspec.md) — Reference data; Built-ins, generated config and the legacy default; Baseline semantic analysis; API Contract 3; API Contract 4; Testing Approach 4–5; Testing Approach 7; Build Order 3; Risks & Considerations
- ADR-0180; ADR-0069; ADR-0037; ADR-0049; ADR-0050; ADR-0107; ADR-0140; ADR-0151

## Result

Implemented this Task's slice for Daemon Verification. Task status, Subtasks,
Acceptance Criteria checkboxes and the authored Verification remain unchanged.
No commit, push or Pull Request was made.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Bare configuration resolves five required Recommended Profiles as built-in; optional categories inherit general | `builtinProfiles` iterates `requiredWorkCategories` and reads `RecommendedProfile`, with source `built-in`. `TestBuiltinProfilesAreTheRecommendedProfile` loads empty User and Project scopes and checks each complete profile and source. `TestBuiltinProfilesDefineNoOptionalCategory` checks the five-entry map and inheritance for every optional category. Both passed. |
| Generated User and Project Config parse to the same five profiles | `defaultConfigYAML` renders the profiles block from `builtinProfiles` in required-category order. `TestGeneratedConfigRendersTheBuiltinProfiles` checks YAML order, decoded equality and normal scope loading for both generated configurations; passed. |
| Legacy defaults.agent: codex selects gpt-6.1-sol/high | `Builtin().Runtimes.Codex` reads the general Recommended Profile's model and effort. The two old Codex constants were removed; Claude defaults and `applyLegacyRuntimeProfiles` are unchanged. `TestLegacyCodexDefaultIsTheGeneralCodexSelection` loads that minimal legacy configuration and checks the resolved tuple and runtime defaults; passed. |
| Built-in review passes the Codex provider check | Review is derived from its own Recommended Profile. `TestBuiltinReviewProfileStaysOnTheReviewProvidersRuntime` resolves the actual built-in review and calls `validateReviewProfileProvider("codex", ...)`; passed. |
| Both Baseline models are in the Codex catalog | Preferred is `gpt-6.1-sol`, fallback is `gpt-5.6-sol`; `xhigh` and all analysis logic are unchanged. `TestAnalysisModelsAreInTheCodexModelCatalog` passed. ADR-0069 changes only the two model names, the ADR-0180 attribution and the requested updated_at timestamp. |
| Catalogs equal the final Reference data table | Removed the two transitional entries. `TestModelCatalogOffersNoReplacedModel` checks both forbidden identifiers in both catalogs; `TestModelCatalogsExposeOrderedPickerData` checks the full exact tables; `TestModelCatalogOpensWithTheCurrentModels` checks their current opening order. All passed. |
| Production code and specified documents contain no replaced identifier | A local Python sweep checked all 226 non-test Go files under internal/cmd and all five specified documents for `gpt-5.5`, and the documents for the exact replaced Claude identifier. No matches. It also checked both mirrored skill files byte-for-byte and proved the QA settlement section equals HEAD. The documentation contract and a temporary overlay test comparing all four documented YAML profiles blocks to Builtin passed. |

The guides and canonical skill now document the same Preferred Selections and
Fallback Chains as setup generates. The skill was raised from 0.0.6 to 0.0.7 in
both version fields. The documentation contract removes both GPT-5.5 pins,
requires `claude-fable-5-1`, and derives its built-in Preferred Selection pin
from `RecommendedProfile(CategoryGeneral)`.

### Focused checks

All Go checks below used `GOCACHE=/private/tmp/roundfix-task03-cache` and
`rtk proxy`; the first attempt using the shared cache was refused by the
sandbox. The task-local cache rerun provided the expected initial red:
all five built-in equality cases and both replaced-catalog cases failed before
production changes.

- `go test ./internal/config -count=1` — passed after the final config-test edit.
- `go test ./internal/config ./internal/agent ./internal/baselineacp -run 'TestBuiltin|TestGeneratedConfig|TestLegacyCodex|TestModelCatalog|TestAnalysisModels' -count=1` — passed.
- `go test ./internal/cli -run 'TestRunInit|TestRunSetup|TestProfiles|TestRunDoctor|TestInvocation|TestDoctorNamesA|TestCharacterization|TestResolveSelection|TestImplementTaskContent|TestRunImplementDetach' -count=1 -timeout=90s` — passed.
- `go test ./internal/cli -run 'TestBuiltinReview|TestSetupCommandCompatibility|TestCharacterization|TestResolveSelectionUses' -count=1` — passed.
- `go test -tags docscontract ./internal/docscontract -run TestProfilesDocumentation -count=1` — passed.
- `go test -overlay /private/tmp/task03-doc-overlay.json ./internal/config -run TestTask03DocumentedProfiles -count=1` — passed; the overlay adds no repository file and checks both guides and both skill copies against generated profiles.
- `git -c core.fsmonitor=false diff --check` — passed.

An initial broad diagnostic run of the config/CLI packages was interrupted
after invalidated config expectations and detached-test startup failures were
identified; it is not full-suite evidence. The detached test adapter previously
advertised only the old Codex built-ins. Its helper now derives advertised
models from Builtin, and installs separate isolated Codex and Claude adapter
commands/configuration so both runtime lineages are proven. Both affected
detached flows passed in the focused CLI run. Test-only callers left by the
earlier failing runs were terminated; the active Task Agent was not touched.

The Task's authored Verification and the repository-wide gate were not run.
They remain Daemon-owned for this handoff.

### Existing tests updated

Every change below follows an invalidated built-in expectation, generated
configuration expectation or the declared model/catalog break. Tests meaning
the built-in value now read RecommendedProfile or Builtin; explicit configured
decision values remain unchanged.

`internal/config/config_test.go`:

- `TestBuiltinRuntimeDefaults`
- `TestBuiltinProfilesGeneratedCodexPolicy`
- `TestDefaultConfigYAMLGeneratedCodexPolicy`
- `TestAgentSelectionProfileBuiltinsResolveRequiredCategories`
- `TestProfileLegacyMigrationConvertsRuntimeDefaults`
- `TestProfileLegacyDefaultCodexKeepsDistinctBuiltInFallback`
- `TestProfileResolverPreferredOverridePreservesFallbackChain`
- `TestLoadWarnsAndIgnoresDeprecatedDefaultsModel`
- `TestInitCreatesUserConfig`
- `TestProfileGeneratedConfigUsesCompleteProfilesSchema`
- `TestInitForceOverwritesExistingConfig`

`internal/cli/cli_test.go`:

- `TestRunInitForceOverwritesExistingConfig`
- `TestProfilesValidateDeduplicatesProofsAndReportsEveryReference`
- `TestProfilesValidateTextNamesADegradedPolicy`
- `TestDoctorNamesADegradedPolicy`
- `TestInvocationProfileOverrideOmittedUsesTaskQAAndReviewProfiles`
- `TestInvocationProfileOverrideAppliesAcrossCategoriesPreservesFallbacksAndWarns`
- `TestRunSetupHealthyMachineIsIdempotent` and `TestSetupCommandCompatibility`, through their shared `assertSetupCommandHealthyMachineIsIdempotent` helper
- `TestRunSetupProfileProofsEveryDistinctTupleOnceBeforePersistence`
- `TestRunSetupProfilePersistenceMatchesSubsequentValidation`
- `TestRunSetupNoInputProfileProofCreatesNoTargets`
- `TestRunSetupProfileProofUsesProposedProfilesAndWorkDir`
- `TestRunSetupAcceptsConfiguredEmptyReasoningEffort`

Other existing tests:

- `internal/cli/doctor_test.go`: `TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts`, `TestRunDoctorProfileReadinessReportsLegacyAdapterThroughEffectiveProfile`.
- `internal/cli/doctor_characterization_test.go`: `TestCharacterizationInvariantDoctorCountsAreUnchangedWithoutOptionalCategories` now expects **4 distinct tuples / 10 category references**, including four proof requests; `TestCharacterizationInvariantInheritedCategoryAddsNoTuple` now expects **4 / 10**. The configured optional-category test still expects **5 / 12** and passed unchanged.
- `internal/cli/implement_test.go`: `TestRunImplementDetachPrintsReportAndCompletesRun`, `TestRunImplementDetachSurvivesCallerProcessGroupKill`, and their `fakeACPXCommand` helper.
- `internal/cli/selection_test.go`: `TestResolveSelectionUsesBuiltInRuntimeDefaults`.
- `internal/agent/agent_test.go`: `TestModelCatalogsExposeOrderedPickerData`.
- `internal/docscontract/publicdocs_test.go`: `TestProfilesDocumentationContractMatchesPublicGuidance`.

The scratch-copy list was used as a starting point; focused checks found the
additional generated-config and setup/detached expectations named above.
No Baseline decision fixture or golden file was edited.

### Sabotage evidence

Each sabotage was applied alone, followed by
`go test <package> -run '^<test-name>$' -count=1` with the task-local cache.
Each command exited 1 with the named test failure. Each source file was restored
in a finally block before the next case, and positive focused checks followed.

| Sabotage | Test that failed |
| --- | --- |
| Changed general's built-in preferred model to sabotaged after reading its Recommended Profile | `TestBuiltinProfilesAreTheRecommendedProfile/general` |
| Derived built-in review from general | `TestBuiltinReviewProfileStaysOnTheReviewProvidersRuntime`; Codex provider rejected fallback 1 runtime Claude |
| Changed only the generated YAML's preferred models to sabotaged | `TestGeneratedConfigRendersTheBuiltinProfiles/user` and `/project` |
| Changed only the legacy Codex runtime model to sabotaged | `TestLegacyCodexDefaultIsTheGeneralCodexSelection` |
| Added docs as a built-in entry | `TestBuiltinProfilesDefineNoOptionalCategory` |
| Set the Baseline preferred model to removed GPT-5.5 | `TestAnalysisModelsAreInTheCodexModelCatalog/gpt-5.5` |
| Reinserted GPT-5.5 into the Codex catalog | `TestModelCatalogOffersNoReplacedModel/codex/gpt-5.5` |

### Regeneration

`make skills-sync` exited 0. Its only content change was
`skills/roundfix/SKILL.md`. The target recreates every owned mirror, so it
also rewrote the following files with unchanged content:

- `skills/archive-spec/SKILL.md`
- `skills/brainstorming/SKILL.md`
- `skills/business-analyst/SKILL.md`
- `skills/council/SKILL.md`
- `skills/council/assets/synthesis-template.md`
- `skills/council/references/archetypes.md`
- `skills/council/references/debate-protocols.md`
- `skills/evidence-gate/SKILL.md`
- `skills/implement-spec/SKILL.md`
- `skills/implement-task/SKILL.md`
- `skills/qa-gate/SKILL.md`
- `skills/roundfix/agents/openai.yaml`
- `skills/setup-context-driven/SKILL.md`
- `skills/write-idea/SKILL.md`
- `skills/write-idea/references/idea-template.md`
- `skills/write-idea/references/opportunity-scan.md`
- `skills/write-prd/SKILL.md`
- `skills/write-prd/references/prd-template.md`
- `skills/write-tasks/SKILL.md`
- `skills/write-tasks/references/task-template.md`
- `skills/write-techspec/SKILL.md`
- `skills/write-techspec/references/techspec-template.md`

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
exited 0 and rewrote `skills/testdata/owned-skill-versions.json` with the
0.0.7 version and digest. This path is declared here as required by the
repository's owned-skill version rule.

`make baseline-digests` exited 0. Its final result was
`{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}`;
no derived artifact retained changed bytes. There are no hand-edited digest
pins or changes to Baseline decision fixtures/goldens.

No follow-up outside this Task's slice was implemented.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20260930T215154Z_eb58c6882d4f3968`
- Source commit: `f39f485084a7fe2b4935ea7b8f0572e01eaa09de`
