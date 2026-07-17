---
task: task_09
spec: 0035-agent-selection-profiles
status: pending
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

- [ ] Update profile configuration and CLI user guidance.
- [ ] Document official ids, fallback activation, and migration boundaries.
- [ ] Add recommendation limitations and source-date language.
- [ ] Update generated Project/User Config examples.
- [ ] Update the canonical Roundfix skill command recipes.
- [ ] Regenerate embedded owned skills and pin contract tests.

## Acceptance Criteria

- [ ] User guidance provides complete copy-paste profile examples with at least one fallback per required category.
- [ ] Optional categories, atomic precedence, one-Run overrides, and legacy migration are unambiguous.
- [ ] Recommendation documentation states the snapshot date, cost/result caveat, `category_specific: false`, and non-routing boundary.
- [ ] Fallback documentation states notification-before-activation and no fallback after work starts.
- [ ] The Roundfix skill matches every shipped public command and synchronized copies are byte-identical.
- [ ] No `write-tasks` skill file contains runtime ids, model ids, rankings, or profile configuration instructions.

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
