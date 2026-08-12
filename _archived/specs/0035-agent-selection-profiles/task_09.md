---
task: task_09
spec: 0035-agent-selection-profiles
status: completed
type: docs
complexity: medium
---

# Task 09: Synchronize profile guidance and owned skills

## Overview

Document complete profile configuration, official identifiers, category inheritance, CLI management, fallback boundaries, and migration without moving model policy into skills. The canonical Roundfix skill and generated embedded copy remain synchronized with the shipped CLI behavior.

## Requirements

1. MUST document built-in required profiles, optional inheritance, atomic scope precedence, complete Fallback Chains, and invocation Preferred Selection overrides.
2. MUST use official model identifiers in built-in examples and clearly label explicit runtime aliases as custom forward-compatible values when accepted.
3. MUST document `profiles show`, `configure`, and `validate` text/JSON, interactive/file, dry-run, scope, and recovery flows.
4. MUST explain that recommendations are five-entry dated advisory snapshots, not category-specific proof, automatic routing, or configuration mutation.
5. MUST document notification-first, pre-prompt-only automatic fallback and the prohibition on fallback after Agent work begins.
6. MUST update generated config guidance and legacy migration instructions without editing runtime-owned settings or credentials.
7. MUST update the canonical Roundfix skill first, regenerate its embedded copy, and keep `write-tasks` limited to Task Type authoring.

## Subtasks

- [x] Update profile configuration and CLI user guidance.
- [x] Document official ids, fallback activation, and migration boundaries.
- [x] Add recommendation limitations and source-date language.
- [x] Update generated Project/User Config examples.
- [x] Update the canonical Roundfix skill command recipes.
- [x] Regenerate embedded owned skills and pin contract tests.

## Acceptance Criteria

- [x] User guidance provides complete copy-paste profile examples with at least one fallback per required category.
- [x] Optional categories, atomic precedence, one-Run overrides, and legacy migration are unambiguous.
- [x] Recommendation documentation states the snapshot date, cost/result caveat, `category_specific: false`, and non-routing boundary.
- [x] Fallback documentation states notification-before-activation and no fallback after work starts.
- [x] The Roundfix skill matches every shipped public command and synchronized copies are byte-identical.
- [x] No `write-tasks` skill file contains runtime ids, model ids, rankings, or profile configuration instructions.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `docs/user-guide/usage.md`
- interface: `.roundfixrc.yml`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`

## Verification

- `grep -F 'roundfix profiles show' docs/user-guide/usage.md && grep -F 'roundfix profiles configure' docs/user-guide/usage.md && grep -F 'roundfix profiles validate' docs/user-guide/usage.md && grep -F 'claude-fable-5' docs/user-guide/usage.md` — expected: public profile management and official-id guidance are present.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` — expected: canonical and embedded owned skills are byte-identical.
- `rtk make skills-sync-check` — expected: all owned skill bundles are synchronized.
- `rtk go test ./internal/config ./internal/cli -run 'Test(ProfileGeneratedConfig|ProfilesDocumentationContract)' -count=1` — expected: generated examples and public command guidance match the implemented schema.

## References

- `_prd.md` → Goals 4 and 8; User Stories 2, 4-5, and 7; Core Features 3-7 and 11; Non-Goals; Decisions.
- `_techspec.md` → Configuration schema; Profile CLI; Recommendation data; `write-tasks` contract; Skill ownership risk; Build Order 8.
- `references/model-ranking.md` → recommendation source and interpretation rules.
- `references/openclaw-skill-analysis.md` → CLI ownership and fallback guardrails.

## Result

Implemented the documentation and owned-skill synchronization slice:

- Updated `docs/user-guide/usage.md` and `docs/user-guide/configuration.md` with complete required profile YAML, official identifiers, optional-category inheritance, atomic Project/User/built-in precedence, one-Run Preferred Selection overrides, profile CLI text/JSON flows, advisory recommendation boundaries, notification-first fallback, no fallback after `agent_work_started`, and legacy migration guidance.
- Updated generated config examples in `.roundfixrc.yml` and `internal/config.DefaultConfigYAML()` to emit only the new `profiles` schema with at least one fallback for each required category.
- Updated `.agents/skills/roundfix/SKILL.md`, regenerated `skills/roundfix/SKILL.md`, and refreshed the setup-context snapshot digest for the Roundfix skill so owned skill audits remain clean.
- Added `TestProfileGeneratedConfigUsesCompleteProfilesSchema` and `TestProfilesDocumentationContractMatchesPublicGuidance` to pin generated config and public documentation contracts, including the `write-tasks` no-policy boundary.

Acceptance evidence:

- Complete copy-paste examples: `rtk grep -F 'claude-fable-5' docs/user-guide/usage.md` found official frontend examples, and `TestProfileGeneratedConfigUsesCompleteProfilesSchema` passed with required profiles and fallbacks.
- Optional categories, atomic precedence, overrides, and migration: documented in usage/configuration guides and covered by `TestProfilesDocumentationContractMatchesPublicGuidance`.
- Recommendations: usage/configuration/skill docs state `2026-07-16`, cost/result evidence, `category_specific: false`, and advisory non-routing/non-mutating boundaries; contract test passed.
- Fallback boundaries: docs and skill state notification-before-activation and no fallback after `agent_work_started`; contract test passed.
- Skill sync: `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && rtk cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` passed, and `rtk make skills-sync-check` passed.
- `write-tasks` boundary: `TestProfilesDocumentationContractMatchesPublicGuidance` passed, asserting no runtime ids, model ids, rankings, or profile configuration terms in canonical or embedded `write-tasks`.

Verification:

- `rtk grep -F 'roundfix profiles show' docs/user-guide/usage.md && rtk grep -F 'roundfix profiles configure' docs/user-guide/usage.md && rtk grep -F 'roundfix profiles validate' docs/user-guide/usage.md && rtk grep -F 'claude-fable-5' docs/user-guide/usage.md` — passed.
- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && rtk cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` — passed.
- `rtk make skills-sync-check` — passed.
- `rtk go test ./internal/config ./internal/cli -run 'Test(ProfileGeneratedConfig|ProfilesDocumentationContract)' -count=1` — passed, 2 tests in 2 packages.
- `rtk make setup-context-check` — passed, 79 tests.
- `rtk make verify` — passed, including `rtk go test ./...` with 1544 tests, setup-context checks, `roundfix skills check`, and build.
