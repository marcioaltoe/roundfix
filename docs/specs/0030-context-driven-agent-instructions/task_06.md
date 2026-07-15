---
task: task_06
spec: 0030-context-driven-agent-instructions
status: pending
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

- [ ] Add the focused setup-context validation target and wire it into `make verify`.
- [ ] Add or finalize macro fixtures for every supported profile and critical migration path.
- [ ] Run the skill's structural validator against the canonical skill directory.
- [ ] Regenerate the embedded owned-skill bundle using `make skills-sync`.
- [ ] Run embedded drift checks and the complete Python suite from a clean process.
- [ ] Run the full repository verification gate and record any residual risks in the Task result.

## Acceptance Criteria

- [ ] The focused setup-context validation target passes independently.
- [ ] All three initial profiles complete apply → clean audit → no-op apply in macro tests.
- [ ] Required-skill failures and optional extras retain their documented blocking semantics.
- [ ] The canonical and embedded `setup-context-driven` skill trees are identical.
- [ ] Skill structure validation accepts the revised skill with its scripts and assets.
- [ ] `make verify` passes completely with no skipped required check.

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
