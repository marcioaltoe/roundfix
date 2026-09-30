---
task: task_02
spec: 0189-profiles-that-follow-the-current-models
status: completed
type: backend
complexity: high
---

# Task 02: The recommendation is a dated Recommended Profile

## Overview

`internal/config/recommendations.go` holds a five-category ranking dated 2026-08-07, with a DeepSWE result and an average cost on every row. No model that leads today has both figures published, and optional categories only borrow the `general` ranking. This Task replaces the ranking with one Recommended Profile per Agent Work Category, dated 2026-09-30 (ADR-0180), and makes `roundfix profiles show` and the interactive configure flow print it. It is verifiable on its own through `roundfix profiles show`, whose JSON moves to schema `roundfix/profiles/v2`.

## Requirements

1. MUST set `ModelRecommendationSnapshotVersion` and `ModelRecommendationSnapshotDate` to `2026-09-30`, and hold one Recommended Profile for each of the ten Agent Work Categories exactly as `_techspec.md` → Reference data → Recommended Profile states, with the rationale given there for each selection.
2. MUST add `RecommendedProfile(category)`, `RecommendationRole` and its two values, and reshape `ModelRecommendation` and `ModelRecommendations` as `_techspec.md` → Interfaces states. `ModelRecommendations` MUST derive its rows from the same table `RecommendedProfile` reads: the Preferred Selection at rank 1 with role `preferred`, then each fallback with role `fallback`. The fields `Benchmark`, `ResultPercent`, `AverageCostUSD` and `CategorySpecific` MUST be removed.
3. MUST make `roundfix profiles show` print, for each category, the block under `_techspec.md` → Show and configure, keep the `unavailable:` line, and drop the lines `Recommendation source` and `Recommendations snapshot`. Its JSON MUST use schema `roundfix/profiles/v2` and the fields of API Contract 1, without `recommendation_source`. Its flags and exit codes MUST not change. The help text of `profiles` and `profiles show` in `internal/cli/cli.go` MUST name `roundfix/profiles/v2` and the Recommended Profile, not a top-five list.
4. MUST make the interactive `profiles configure` flow print the same rows, still marked advisory. Recommendations MUST stay advisory everywhere: they never select, route or write.
5. MUST NOT change `builtinProfiles`, the generated config or any default in this Task. task_03 owns them.
6. MUST put the new tests in `internal/config/recommended_profile_test.go` and `internal/cli/profiles_show_recommended_test.go`. The config tests MUST prove:
   - every category has a profile equal to the Reference data;
   - the first fallback changes runtime in every category except `review`, and `review` uses one runtime throughout;
   - every recommended `codex` or `claude` model is in `agent.ModelCatalog`;
   - `ModelRecommendations` equals the profile in order, with roles.
7. MUST update only the existing tests this change invalidates, and name each in the Result. Tests that mean "the snapshot" MUST read `ModelRecommendationSnapshotVersion` instead of repeating the date.
8. MUST prove each new gate can fail. The Result MUST record one sabotage for the catalog invariant (a recommended model removed from the catalog) and one for the `review` runtime rule (a Claude fallback in `review`), each with the test that failed, and that the code was restored.
9. MUST describe the Recommended Profile and the `v2` schema in `docs/user-guide/usage.md`, `docs/user-guide/configuration.md`, `docs/user-guide/commands.md` and `.agents/skills/roundfix/SKILL.md`, by replacing in place the sentences about the five-entry snapshot, its benchmark fields and `category_specific`. It MUST add no new section to the skill. It MUST add no text inside the `### QA settlement` section of the skill, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
10. MUST replace the glossary term **Model Recommendation Ranking** in `CONTEXT.md` with **Recommended Profile**: the dated Agent Selection Profile Roundfix recommends for one Agent Work Category, from which built-in profiles derive, and which never selects, routes or changes a configuration by itself. The old term MUST move to its `_Avoid_` line.
11. MUST update the pins of `TestProfilesDocumentationContractMatchesPublicGuidance` in `internal/docscontract/publicdocs_test.go` that this Task invalidates: the snapshot date (read from the constant), `roundfix/profiles/v1` and `category_specific: false`. It MUST leave the `gpt-5.5` pins for task_03.

## Subtasks

- [ ] Replace the ranking table with the ten Recommended Profiles and the derived rows.
- [ ] Print the Recommended Profile in `profiles show` and the configure flow, and move the JSON to `v2`.
- [ ] Add the invariant and output tests, and update the invalidated ones.
- [ ] Update the guides, the skill, the glossary and the documentation contract pins, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] `RecommendedProfile` returns the Reference data profile for each of the ten categories, and an unknown category returns false.
- [ ] Removing a recommended model from the Model Catalog fails `TestRecommendedModelsAreInTheModelCatalog`.
- [ ] `roundfix profiles show --category docs --json` reports schema `roundfix/profiles/v2` and two rows, `preferred` then `fallback`, with no `benchmark`, `result_percent`, `average_cost_usd`, `category_specific` or `recommendation_source` key.
- [ ] The text output prints `Recommended profile (snapshot 2026-09-30):` followed by the ranked rows and their rationales.
- [ ] `profiles show` still changes no configuration and no Run state.
- [ ] The guides, the skill, its mirror and the glossary describe the Recommended Profile, and the documentation contract passes.

## Context

- interface: `internal/config/recommendations.go`
- interface: `internal/cli/profiles.go`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- creates: `internal/config/recommended_profile_test.go`
- creates: `internal/cli/profiles_show_recommended_test.go`
- instruction: `internal/agent/catalog.go`
- instruction: `docs/adr/0180-built-in-selections-derive-from-one-dated-recommended-profile.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRecommendedProfileCoversEveryWorkCategory|TestRecommendedProfileFallbackChangesRuntimeExceptReview|TestRecommendedModelsAreInTheModelCatalog|TestModelRecommendationsListTheRecommendedProfileInOrder|TestProfilesShowPrintsTheRecommendedProfile|TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields|TestProfilesShowTextAndJSONAreByteStableAndConsistent|TestModelRecommendationsUseOfficialCatalogModels)$" ./internal/config ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRecommendedProfileCoversEveryWorkCategory TestRecommendedProfileFallbackChangesRuntimeExceptReview TestRecommendedModelsAreInTheModelCatalog TestModelRecommendationsListTheRecommendedProfileInOrder TestProfilesShowPrintsTheRecommendedProfile TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields TestProfilesShowTextAndJSONAreByteStableAndConsistent TestModelRecommendationsUseOfficialCatalogModels; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the six new named tests exists, so the command fails.
- `for pair in "docs/user-guide/usage.md|roundfix/profiles/v2" "docs/user-guide/commands.md|roundfix/profiles/v2" ".agents/skills/roundfix/SKILL.md|roundfix/profiles/v2" "skills/roundfix/SKILL.md|roundfix/profiles/v2" "docs/user-guide/usage.md|Recommended Profile" "docs/user-guide/configuration.md|Recommended Profile" ".agents/skills/roundfix/SKILL.md|Recommended Profile" "skills/roundfix/SKILL.md|Recommended Profile" "CONTEXT.md|**Recommended Profile**" "docs/user-guide/configuration.md|2026-09-30"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check && out="$(go test -count=1 -tags docscontract -v -run "^(TestProfilesDocumentationContractMatchesPublicGuidance)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfilesDocumentationContractMatchesPublicGuidance; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no document names `roundfix/profiles/v2`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; Core Feature 3; Success Metrics 2 and 3; Declared breaks
- [_techspec.md](_techspec.md) — Reference data; Interfaces; The Recommended Profile; Show and configure; API Contract 1; Testing Approach 3; Build Order 2
- ADR-0180; ADR-0049; ADR-0151; ADR-0153

## Result


Implemented this slice for Daemon Verification; status and authored Verification
are untouched. No commit, push, or Pull Request was made. The initial worktree
contained only the Daemon's change to this Task file.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Ten Reference data profiles; unknown category returns false | One table in `recommendations.go` supplies copied `RecommendedProfile` values and ordered `ModelRecommendations` rows. `TestRecommendedProfileCoversEveryWorkCategory` checks all ten reference tuples, unknown categories, and defensive copies. `TestModelRecommendationsListTheRecommendedProfileInOrder` checks order, roles, ranks, category, constant-derived source date, rationales, and independent returned rows. |
| Catalog removal fails the invariant | `TestRecommendedModelsAreInTheModelCatalog` checks every recommended selection against the runtime's Model Catalog. The catalog sabotage below produced the required failure; the restored check passes. `TestRecommendedProfileFallbackChangesRuntimeExceptReview` checks cross-runtime first fallbacks and Codex throughout review. |
| Docs JSON is v2 with two roles and no removed keys | `TestProfilesShowJSONIsSchemaV2WithoutBenchmarkFields` drives the public CLI, checks exact profile and recommendation key sets, roles, rank, source date, and the reference docs selections (`codex / gpt-5.6-luna / max`, then `claude / sonnet / high`). Removed benchmark and recommendation-source keys are absent. |
| Dated text block, ranked rows and rationales | `TestProfilesShowPrintsTheRecommendedProfile` checks all ten category blocks and the same rows/rationales in configure, with configure still labeled advisory. Help names the Recommended Profile and v2. `TestProfilesShowTextAndJSONAreByteStableAndConsistent` checks deterministic output. `TestProfilesShowReportsUnavailableRecommendationWithoutReordering` preserves availability metadata at the original rank; inspection confirms the text `unavailable:` line remains. |
| Show changes no config or Run state | `TestProfilesShowDoesNotMutateConfigOrRunState` passes. Configured effective profiles remain independent of recommendations in `TestProfilesShowJSONRendersProfileAndRecommendations`. `TestAgentSelectionProfilesMacro` passes through the built binary and checks recommendations remain unchanged during a Run. `internal/config/profiles.go` and `internal/config/config.go` are byte-identical to HEAD: built-ins, generated config and defaults stay with task_03. |
| Guides, skill, mirror, glossary and docs contract | Replaced the existing snapshot descriptions in all three guides and the skill, replaced the glossary term with Recommended Profile, and moved the old term to its Avoid line. `TestProfilesDocumentationContractMatchesPublicGuidance` passes with the snapshot constant and v2 pins; gpt-5.5 pins remain. Skill mirror bytes match, and the QA settlement section is byte-identical to HEAD. |

### Focused checks

All Go checks used `GOCACHE=/tmp/roundfix-task02-gocache`.

- Starting signal: `go test ./internal/config -run '^TestRecommendedProfileCoversEveryWorkCategory$'`
  exited 1 against the old implementation: missing RecommendedProfile/roles
  and the old three-result ModelRecommendations signature.
- `go test ./internal/config ./internal/cli -run 'TestRecommended|TestModelRecommendations|TestProfilesShow|TestProfilesConfigure' -count=1`
  exited 0 after the final code edits.
- `go test ./internal/config ./internal/cli -run 'TestRecommended|TestModelRecommendations|TestProfilesShow|TestProfilesConfigure|TestAgentSelectionProfilesMacro' -count=1`
  exited 0, including the binary macro.
- `go test -tags docscontract ./internal/docscontract -run ProfilesDocumentation -count=1`
  exited 0 after the final code edits.
- `rtk make verify-incremental` initially exited 2: two force-stop integration
  tests could not enumerate the process table in the sandbox, and coverage
  equivalence rejected renaming an existing test. Preserved the recorded test
  name and reran with process-table permission; the rerun exited 0, covering
  formatting, vet, all package tests, skill sync/checks, and build.
- `git -c core.fsmonitor=false diff --check` exited 0.
- Authored `## Verification` commands were not run; they remain Daemon-owned.

### Sabotage evidence

- Removed the `gpt-6.1-sol` entry from `internal/agent/catalog.go`, then ran
  `go test ./internal/config -run '^TestRecommendedModelsAreInTheModelCatalog$' -count=1`.
  It exited 1 with `--- FAIL: TestRecommendedModelsAreInTheModelCatalog`
  and reported the absent recommended model across implementation, frontend
  and review categories. Restored the catalog byte-for-byte; the subsequent
  focused checks and incremental check pass.
- Replaced review's fallback with `claude / opus / high`, then ran
  `go test ./internal/config -run '^TestRecommendedProfileFallbackChangesRuntimeExceptReview$' -count=1`.
  It exited 1 with `--- FAIL: TestRecommendedProfileFallbackChangesRuntimeExceptReview`
  and `review must stay on Codex`. Restored recommendations.go byte-for-byte;
  the subsequent focused checks and incremental check pass.

### Existing tests intentionally updated

- `TestProfilesShowJSONRendersProfileAndRecommendations`: v2 shape and
  the new preferred recommendation, while preserving configured precedence.
- `TestProfilesShowOptionalCategoryReportsGeneralRecommendationSource`:
  preserved its recorded identity; now checks the category's own recommendation
  while its effective profile still inherits general.
- `TestProfilesShowTextAndJSONAreByteStableAndConsistent`: its shared text
  assertion reads the snapshot constant and compares roles with selections.
- `TestProfilesShowReportsUnavailableRecommendationWithoutReordering`:
  marks the recommended Opus fallback unavailable without moving it.
- `TestModelRecommendationsUseOfficialCatalogModels`: two-result API,
  two rows, current evidence fields and constant-derived date.
- `TestAgentSelectionProfilesMacro`: its frontend-show helper expects two
  recommendations; all configured profile and Run assertions remain.
- `TestProfilesDocumentationContractMatchesPublicGuidance`: snapshot
  constant, Recommended Profile and v2; task_03 model pins remain.

### Regeneration and recorded paths

- Raised both canonical Roundfix skill version fields from 0.0.5 to 0.0.6.
- `make skills-sync` changed the contents of `skills/roundfix/SKILL.md`;
  every other copied skill file remained byte-identical.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and rewrote `skills/testdata/owned-skill-versions.json`.
  That repository-required generated version record is declared here as part
  of this Task's changed paths.
- `make baseline-digests` exited 0 and reported
  `changed: false`: it rewrote no artifact contents.
- Final changed paths are the Task's named interfaces, its two new test files,
  this Task's Result, the mirrored skill, and the generated owned-skill version
  record. No Task Graph, other Task, adapter catalog, built-in/default source,
  Baseline source, or QA settlement section changed.

No follow-up implementation was included from task_03 or task_04.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20260930T215154Z_eb58c6882d4f3968`
- Source commit: `3433ffea266017b4446a53940ace4bd61fbf9de7`
