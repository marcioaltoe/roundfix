---
task: task_06
spec: 0023-agent-model-catalog-and-isolation
status: completed
type: docs
complexity: medium
---

# Task 06: Ship Agent selection guidance

## Overview

Ship the configuration and operational guidance in the same change as the new
selection behavior. The slice is verifiable through copy-paste examples, the
canonical-to-embedded Skill sync check, and removal of guidance that still
depends on runtime-owned model defaults.

## Requirements

1. MUST update the canonical Roundfix Skill for the new config keys, defaults, flags, Interactive Input, preflight failures, and inspection behavior.
2. MUST synchronize the embedded Skill using the repository generator rather than editing it directly.
3. MUST document precedence, both Model Catalogs, custom values, deprecated-key migration, and unsupported-selection recovery.
4. MUST provide copy-paste Project Config and one-Run examples for Codex and Claude plus OpenCode's explicit-configuration requirement.
5. MUST update repository dogfood configuration and autonomous-work guidance to pin Roundfix-owned model and reasoning values.
6. MUST preserve canonical glossary terms and remove guidance that tells callers to rely on `~/.codex/config.toml` or another runtime default.

## Subtasks

- [x] Update user-facing configuration and command documentation.
- [x] Migrate repository dogfood configuration and autonomous runtime guidance.
- [x] Update the canonical Roundfix Skill contract and examples.
- [x] Regenerate the embedded Skill bundle.
- [x] Audit glossary vocabulary and deprecated guidance.
- [x] Verify all examples and Skill governance checks.

## Acceptance Criteria

- [x] Documentation contains valid Codex and Claude configuration plus one-Run override examples using model and reasoning together.
- [x] OpenCode documentation states that both values are required and provides no invented catalog/default.
- [x] Recovery guidance names runtime update and supported-value selection without suggesting fallback or home-config mutation.
- [x] `.roundfixrc.yml` pins the repository's effective Codex model and reasoning under the per-runtime structure.
- [x] Autonomous-work guidance uses Roundfix-owned selection and the canonical Supervisor role.
- [x] Canonical and embedded Roundfix Skills have zero sync drift and describe the shipped CLI exactly.

## Verification

- `rtk make skills-sync-check` - expected: canonical and embedded Skill bundles have zero drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` - expected: every shipped Skill passes validation.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/tech-writer/SKILL.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `docs/agents/skill-governance.md`
- interface: `docs/agents/autonomous-work.md`
- interface: `.roundfixrc.yml`
- interface: `README.md`

## References

`_prd.md` -> User Story 7; Core Feature 11; Success Metrics; Decisions. `_techspec.md` -> API Contracts; Integration Points; Build Order 6. ADR-0037; ADR-0039.

## Result

- Updated `README.md` with per-runtime `runtimes.<runtime>.model` and `runtimes.<runtime>.reasoning_effort` guidance, built-in Codex and Claude defaults, ordered Codex and Claude Model Catalogs, custom-value behavior, deprecated `defaults.model` migration, unsupported-selection recovery, and stored-selection inspection behavior.
- Added Codex and Claude one-Run examples that pass `--model` and `--reasoning-effort` together, plus Project Config examples for Codex and Claude. OpenCode guidance states that both values are required and that Roundfix ships no OpenCode Model Catalog or default.
- Updated `.roundfixrc.yml` to remove `defaults.model` and pin this repository's Codex selection under `runtimes.codex.model: gpt-5.5` and `runtimes.codex.reasoning_effort: xhigh`.
- Updated `docs/agents/autonomous-work.md` to use the canonical Supervisor role and Roundfix-owned selection instead of `~/.codex/config.toml` or runtime-owned defaults.
- Updated the canonical `.agents/skills/roundfix/SKILL.md` for new config keys, defaults, flags, Interactive Input ordering, preflight rejection recovery, inspection behavior, and examples; regenerated `skills/roundfix/SKILL.md` with `rtk make skills-sync`.

Verification:

- `rtk make skills-sync-check` passed with no drift output.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed: every shipped Skill validated.
- `rtk make verify` passed: `go test ./...` reported 1039 tests in 19 packages, skill checks passed, and `go build` completed.
