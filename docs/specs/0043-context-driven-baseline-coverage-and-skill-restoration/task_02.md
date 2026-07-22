---
task: task_02
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: docs
complexity: high
---

# Task 02: Generate complete Context-Driven Baseline guidance

## Overview

Make every bundled profile generate the portable guidance it promises while
keeping root instructions compact. The slice renders authoritative rule text,
selected Verification, and active skill dispatch into setup-owned guides while
leaving project-specific architecture and unmarked repository content under
repository ownership.

## Requirements

1. MUST declare every universal and applicable coverage category and required
   rule for each bundled profile.
2. MUST render rule guidance from the declarative rule authority rather than
   maintaining independent metadata and prose.
3. MUST make `verification.gate` an entry decision for every profile and render
   the selected command in universal completion guidance.
4. MUST generate deterministic skill dispatch equal to the active modules'
   required skills, retaining distinct triggers for shared skills.
5. MUST cover portable language, research, dependency, Git/delivery,
   security/configuration, verification-integrity, and enabled-surface rules
   without copying project-specific policy.
6. MUST keep root blocks as short mandatory pointers, identify setup-owned
   guidance, preserve repository-authored extensions, and point frontend
   profiles to repository-owned `DESIGN.md`.
7. MUST generate the monorepo guide whenever the monorepo root block is active.

## Subtasks

- [ ] Declare mandatory and profile-specific coverage rule sets.
- [ ] Add portable rule guidance to the owning modules and supporting guides.
- [ ] Render selected Verification in universal agent instructions.
- [ ] Generate the ordered active-module skill-dispatch guide.
- [ ] Add ownership language and the frontend design-contract pointer.
- [ ] Exercise every profile and relevant rendering decision through macro
      fixtures.

## Acceptance Criteria

- [ ] Every bundled profile covers all universal and applicable categories;
      deleting one required rule or dispatch mapping fails validation.
- [ ] Generated universal guidance names the selected Verification even when
      autonomous work is disabled.
- [ ] Generated dispatch contains every active required skill, no unavailable
      skill, and every distinct declared trigger in deterministic order.
- [ ] Root instructions remain compact pointers while supporting guides carry
      the complete portable rule bodies.
- [ ] Frontend guidance points to repository-owned `DESIGN.md` and does not
      invent project architecture, authentication, database, or transport
      policy.
- [ ] A single-context TypeScript/Bun monorepo generates both `root.monorepo`
      and `guide.monorepo`.
- [ ] Re-rendering every profile preserves unmarked repository-authored bytes
      and produces byte-identical managed output.
- [ ] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/assets/modules`
- interface: `.agents/skills/setup-context-driven/assets/profiles`
- interface: `.agents/skills/setup-context-driven/assets/templates`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_rendering.py`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py`
  — expected: every bundled profile renders complete, deterministic managed
  guidance and preserves repository-authored content.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_decision_rendering.py`
  — expected: selected Verification, ownership language, profile rules, and
  computed skill dispatch render from declared decisions and modules.
- `rtk make verify` — expected: the full repository gate passes with the
  expanded canonical Baseline assets.

## References

- `_prd.md` → User Stories 1–2, 6; Core Features 1–3, 9; User Experience;
  Non-Goals.
- `_techspec.md` → Data Models: coverage and dispatch; Coverage Map: Goals 1–2
  and Stories 1–2, 6; Build Order 2.
- ADR-0046 → compact managed guidance and repository-authored preservation.
- ADR-0047 → decision-bound module activation and render bindings.
