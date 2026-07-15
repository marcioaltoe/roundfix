---
task: task_06
spec: 0030-context-driven-agent-instructions
status: completed
type: test
complexity: medium
---

# Task 06: Gate and synchronize the complete portable skill

## Overview

Turn the completed feature into a repository-enforced contract. This slice adds the focused verification target, validates realistic profile flows, synchronizes the canonical and embedded skill trees, and proves the full Roundfix gate remains clean.

## Requirements

1. MUST add a focused repository verification target that runs the complete Python suite and bundled-asset audit.
2. MUST include the focused target in the authoritative `make verify` gate.
3. MUST synchronize the canonical `.agents/skills/setup-context-driven/` tree into the embedded `skills/setup-context-driven/` bundle through the existing repository mechanism.
4. MUST preserve `skills-sync-check` as the drift gate and avoid hand-maintaining divergent copies.
5. MUST run macro flows for TypeScript/Bun monorepo, Go CLI/TUI, and Rust CLI profiles, including audit, apply, repeated apply, required-skill validation, optional extra-skill reporting, and Secondbrain opt-in where selected.
6. MUST validate skill structure and metadata after adding scripts and assets.
7. MUST keep all existing repository tests, formatting checks, skill checks, and build checks passing with zero failures.

## Subtasks

- [x] Add the focused setup-context validation target and wire it into `make verify`.
- [x] Add or finalize macro fixtures for every supported profile and critical migration path.
- [x] Run the skill's structural validator against the canonical skill directory.
- [x] Regenerate the embedded owned-skill bundle using `make skills-sync`.
- [x] Run embedded drift checks and the complete Python suite from a clean process.
- [x] Run the full repository verification gate and record any residual risks in the Task result.

## Acceptance Criteria

- [x] The focused setup-context validation target passes independently.
- [x] All three initial profiles complete apply → clean audit → no-op apply in macro tests.
- [x] Required-skill failures and optional extras retain their documented blocking semantics.
- [x] The canonical and embedded `setup-context-driven` skill trees are identical.
- [x] Skill structure validation accepts the revised skill with its scripts and assets.
- [x] `make verify` passes completely with no skipped required check.

## Context

- instruction: `.agents/skills/evidence-gate/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `Makefile`
- interface: `docs/agents/skill-governance.md`

## Verification

- `rtk make setup-context-check` — expected: the complete Python behavior suite and bundled-asset audit pass.
- `rtk make skills-sync-check` — expected: canonical and embedded owned-skill trees have no drift.
- `rtk make verify` — expected: formatting, Go tests, setup-context validation, skill checks, and build all pass.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → all Success Metrics; User Story 8.
- `_techspec.md` → Testing Approach; Integration Points: Roundfix skill bundle; Build Order 5; Risks & Considerations.

## Result

- Added `setup-context-check` to `Makefile` and wired it into `make verify`; it runs the complete setup-context Python suite plus bundled asset validation for canonical and embedded skill trees.
- Added macro CLI coverage for TypeScript/Bun monorepo, Go CLI/TUI, and Rust CLI profile flows, including apply → clean audit → no-op apply, required-skill blocking, optional extra-skill info, and Secondbrain opt-in.
- Regenerated `skills/setup-context-driven/` with `rtk make skills-sync` and preserved `skills-sync-check` as the drift gate.
- Evidence: `rtk make setup-context-check` passed with 42 Python tests and asset validation.
- Evidence: `rtk make skills-sync-check` passed with no canonical/embedded drift.
- Evidence: `rtk make verify` passed after rerunning with Go build-cache filesystem approval; Go tests, setup-context validation, skill checks, and build all completed.
- Residual risks: none identified.
