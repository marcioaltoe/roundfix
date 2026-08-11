---
task: task_01
spec: 0030-context-driven-agent-instructions
status: completed
type: backend
complexity: high
---

# Task 01: Define portable instruction profiles and module contracts

## Overview

Create the declarative asset foundation that every later audit and apply operation consumes. The slice is complete when all initial profiles, modules, rule identifiers, decisions, document templates, and canonical skill-setup snapshots can be loaded and validated deterministically without external dependencies.

## Requirements

1. MUST define a versioned asset contract for profiles, modules, rules, decisions, managed root blocks, supporting guides, and canonical skill-setup snapshots.
2. MUST include initial compositions for TypeScript/Bun monorepos, Go CLI/TUI repositories, and Rust CLI repositories.
3. MUST keep the root managed guidance compact and route conditional rule bodies to supporting agent guides.
4. MUST assign stable identifiers and versions to every module, rule, decision, template, profile, and setup snapshot.
5. MUST reject missing dependencies, dependency cycles, conflicting modules, duplicate rule identifiers, duplicate managed block identifiers, and profile references to unknown assets.
6. MUST include pinned portable snapshots of the canonical `typescript-bun`, `go-cli`, and `rust-cli` skill setups with normalized skill paths, names, source metadata, and digests.
7. SHOULD keep reusable output templates in skill assets and reserve references for agent-readable workflow guidance.

## Subtasks

- [x] Define the asset directory layout and versioned JSON contracts.
- [x] Extract the universal, context-workflow, language/runtime, repository-shape, application-surface, and autonomous-work modules from production guidance.
- [x] Compose the TypeScript/Bun monorepo, Go CLI/TUI, and Rust CLI profiles.
- [x] Add normalized canonical skill-setup snapshots for the three profiles.
- [x] Add asset-contract tests covering valid compositions and invalid dependency/conflict cases.
- [x] Document asset selection and loading rules concisely in the skill.

## Acceptance Criteria

- [x] Each supported profile resolves to one deterministic ordered module list.
- [x] Every generated rule and managed artifact has a stable identifier and version.
- [x] Every skill referenced by a profile module exists in that profile's bundled setup snapshot.
- [x] Invalid cycles, conflicts, duplicate identifiers, and unknown references fail asset validation with precise diagnostics.
- [x] The portable assets require no access to `~/dev/skills`, the network, or third-party Python packages.
- [x] The managed root composition contains pointers and invariants rather than copied stack-specific manuals.

## Context

- instruction: `.agents/skills/setup-context-driven/SKILL.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/references/docs-layout.md`
- interface: `.agents/skills/setup-context-driven/references/domain.md`
- interface: `.agents/skills/setup-context-driven/references/autonomous-work.md`
- interface: `docs/agents/skill-governance.md`
- interface: `Makefile`

## Verification

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_assets.py'` — expected: every valid bundled profile passes and every invalid contract fixture fails for its declared reason.
- `rtk git diff --check` — expected: no whitespace errors in the new portable assets or tests.

## References

- `_prd.md` → Goals 1–2; Core Features 1–4, 9, 11; Decisions.
- `_techspec.md` → System Architecture; Data Models; Build Order 1; Risks & Considerations.
- ADR-0046.

## Result

Implemented the declarative setup-context-driven asset foundation:

- Added versioned contracts, decisions, templates, modules, profiles, and pinned setup snapshots under `assets/`.
- Added a Python standard-library asset loader and validator in `scripts/context_assets.py`.
- Added `test_assets.py` coverage for valid profile resolution and invalid asset fixtures.
- Mirrored the canonical setup-context-driven skill changes into `skills/setup-context-driven/` so `skills-sync-check` remains aligned.
- Documented bundled asset selection and loading rules in `SKILL.md`.

Verification:

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_assets.py'` — passed, 6 tests.
- `rtk git diff --check` — passed after recording this Result.
- `rtk make verify` — passed after rerunning with Go build-cache filesystem access; the first sandboxed attempt failed because Go could not open its external build cache.

Acceptance evidence:

- Deterministic profiles: `test_supported_profiles_resolve_to_deterministic_module_order` asserts exact ordered module lists for `typescript-bun-monorepo`, `go-cli-tui`, and `rust-cli`.
- Stable identifiers and versions: `test_every_managed_asset_has_stable_id_and_version` checks decisions, modules, profiles, setup snapshots, templates, rules, root blocks, and guides.
- Skill snapshot coverage: `test_profile_referenced_skills_exist_in_bundled_setup_snapshot` checks every module-required skill is present in the selected profile snapshot.
- Invalid assets: `test_invalid_contract_fixtures_fail_with_precise_diagnostics` covers unknown profile modules, unknown setup snapshots, missing dependencies, cycles, conflicts, duplicate rule IDs, duplicate root block IDs, unknown decisions/templates, and skills outside snapshots.
- Portability: `test_assets_are_portable` loads bundled assets locally and scans them for `~/dev/skills` and network URLs; implementation imports only Python standard-library modules.
- Compact root guidance: `test_root_templates_are_compact_pointers` enforces short root templates that point to supporting `docs/agents/` guides.

Follow-up for later Tasks:

- Wire `scripts/context_assets.py` into the read-only audit CLI in Task 02.
