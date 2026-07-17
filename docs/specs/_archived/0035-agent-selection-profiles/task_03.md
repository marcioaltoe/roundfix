---
task: task_03
spec: 0035-agent-selection-profiles
status: completed
type: backend
complexity: high
---

# Task 03: Show advisory profile recommendations

## Overview

Deliver a read-only, source-aware `profiles show` surface backed by the versioned five-entry recommendation dataset. The configured profile remains primary, and recommendation rank is observable guidance that never mutates or routes a selection.

## Requirements

1. MUST store the 2026-07-16 recommendation snapshot as versioned Go product data with exactly five unique official Agent Models per category.
2. MUST expose `general`, canonical Task Type, `qa`, and `review` categories in stable order and reuse the General list for optional Task Types while naming that inheritance source.
3. MUST render effective profile source, inherited source, Preferred Selection, ordered fallbacks, and exactly five recommendations in deterministic text and `roundfix/profiles/v1` JSON.
4. MUST include benchmark, result, cost, source date, rationale, and `category_specific: false` caveat for every recommendation.
5. MUST reject unknown categories as usage errors and keep stdout free of warnings or diagnostics.
6. MUST never fetch benchmark data, reorder recommendations by machine availability, preselect rank one, change configuration, or use ranking as routing input.
7. SHOULD report a recommendation tuple unavailable when proof data is supplied without substituting or changing its rank.

## Subtasks

- [x] Encode the versioned recommendation schema and fixtures.
- [x] Add stable category and rank ordering.
- [x] Resolve optional-category recommendation inheritance.
- [x] Implement `profiles show` text output.
- [x] Implement `roundfix/profiles/v1` JSON output.
- [x] Add no-mutation and exactly-five contract tests.

## Acceptance Criteria

- [x] Every requested category renders exactly five unique official models in the documented order.
- [x] Optional categories label themselves while reporting `general` as the recommendation source.
- [x] Text and JSON expose the same effective profile, fallback order, snapshot evidence, and caveat.
- [x] Repeated output is byte-stable and independent of map iteration.
- [x] Showing one or all categories leaves User Config, Project Config, runtime-owned config, recommendation data, and Run state unchanged.
- [x] No product path reads ranking data to choose a Preferred Selection or Fallback Chain.

## Context

- interface: `docs/specs/0035-agent-selection-profiles/references/model-ranking.md`
- interface: `internal/agent/catalog.go`
- interface: `internal/cli/cli.go`
- interface: `internal/config/config.go`

## Verification

- `rtk go test ./internal/cli ./internal/agent -run 'Test(ProfilesShow|ModelRecommendations|ModelCatalog)' -count=1` — expected: stable text/JSON, exactly-five, inheritance, official-id, and no-mutation cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix profiles show --category backend --json | rtk python3 -c 'import json, sys; data=json.load(sys.stdin); profile=data["profiles"][0]; assert data["schema"] == "roundfix/profiles/v1"; assert profile["category"] == "backend"; assert len(profile["recommendations"]) == 5'` — expected: the public JSON schema contains one backend profile and five recommendations.

## References

- `_prd.md` → Goal 8; User Stories 4 and 7; Core Features 4-6; Non-Goals; Success Metrics.
- `_techspec.md` → Profile CLI: show; Recommendation data; JSON profile output; Build Order 3.
- `references/model-ranking.md` → versioned five-entry recommendation contract.
- `references/openclaw-skill-analysis.md` → ranking remains advisory and CLI-owned.

## Result

- Added versioned 2026-07-16 Go recommendation data with `DeepSWE v1.1` evidence, cost, source date, rationale, and `category_specific: false` for every row; `TestModelRecommendationsUseOfficialCatalogModels` verifies every exposed category returns exactly five unique official catalog models in rank order.
- Added `roundfix profiles show` text and `roundfix/profiles/v1` JSON output for `general`, canonical Task Types, `qa`, and `review`; `TestProfilesShowTextAndJSONAreByteStableAndConsistent` verifies stable category order, byte-stable repeated output, matching effective profile source, inherited source, Preferred Selection, Fallback Chain, snapshot evidence, and caveat.
- Optional Task Type recommendation inheritance is explicit: `TestProfilesShowOptionalCategoryReportsGeneralRecommendationSource` verifies optional categories label themselves while reporting `general` as `recommendation_source`.
- The configured profile remains primary: `TestProfilesShowJSONRendersProfileAndRecommendations` verifies Project Config Preferred Selection and Fallback Chain render unchanged while recommendations stay advisory.
- Read-only behavior is covered by `TestProfilesShowDoesNotMutateConfigOrRunState`, which shows one category and all categories, then verifies User Config, Project Config, and Run state remain unchanged; the command does not touch runtime-owned config or fetch benchmark data.
- Availability proof handling is covered by `TestProfilesShowReportsUnavailableRecommendationWithoutReordering`, which marks one recommendation unavailable without substitution or rank changes.
- Routing isolation evidence: `rtk rg -n "ModelRecommendations|modelRecommendationsByCategory|ModelRecommendationSnapshot|Recommendation source|recommendation_source" internal cmd` only finds recommendation usage in config data, `profiles show`, and tests.
- Verification passed:
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/cli ./internal/agent -run 'Test(ProfilesShow|ModelRecommendations|ModelCatalog)' -count=1` → `21 passed in 2 packages`.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go run -buildvcs=false ./cmd/roundfix profiles show --category backend --json | rtk python3 -c 'import json, sys; data=json.load(sys.stdin); profile=data["profiles"][0]; assert data["schema"] == "roundfix/profiles/v1"; assert profile["category"] == "backend"; assert len(profile["recommendations"]) == 5'` → passed.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/cli ./internal/agent -count=1` → `678 passed in 2 packages`.
  - `GOCACHE=/private/tmp/roundfix-gocache rtk make verify` → passed; included `rtk go test ./...` with `1492 passed in 20 packages`, skill sync check, and build.
