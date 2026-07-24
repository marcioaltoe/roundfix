---
task: task_02
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
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

- [x] Declare mandatory and profile-specific coverage rule sets.
- [x] Add portable rule guidance to the owning modules and supporting guides.
- [x] Render selected Verification in universal agent instructions.
- [x] Generate the ordered active-module skill-dispatch guide.
- [x] Add ownership language and the frontend design-contract pointer.
- [x] Exercise every profile and relevant rendering decision through macro
      fixtures.

## Acceptance Criteria

- [x] Every bundled profile covers all universal and applicable categories;
      deleting one required rule or dispatch mapping fails validation.
- [x] Generated universal guidance names the selected Verification even when
      autonomous work is disabled.
- [x] Generated dispatch contains every active required skill, no unavailable
      skill, and every distinct declared trigger in deterministic order.
- [x] Root instructions remain compact pointers while supporting guides carry
      the complete portable rule bodies.
- [x] Frontend guidance points to repository-owned `DESIGN.md` and does not
      invent project architecture, authentication, database, or transport
      policy.
- [x] A single-context TypeScript/Bun monorepo generates both `root.monorepo`
      and `guide.monorepo`.
- [x] Re-rendering every profile preserves unmarked repository-authored bytes
      and produces byte-identical managed output.
- [x] Canonical and embedded setup skill trees are synchronized after the
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

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py`
  — expected: every bundled profile renders complete, deterministic managed
  guidance and preserves repository-authored content.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_decision_rendering.py`
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

## Result

Implemented the complete declarative Context-Driven Baseline guidance for all
bundled profiles. Profiles now declare every rule supplied by their universal,
language, workflow, optional, and enabled-surface modules. Module rules own the
portable guidance text and coverage categories, while templates render those
rules into setup-owned guides and keep root blocks as short mandatory pointers.

`verification.gate` is now an entry decision for every profile and renders in
the universal agent-instructions guide even when autonomous work is disabled.
The new skill-dispatch guide is computed in active-module order from exact
`requiredSkills`/`skillDispatch` mappings, deduplicating identical pairs while
retaining distinct triggers for shared skills. Frontend guidance identifies the
repository-owned `DESIGN.md` contract without creating project architecture or
authentication, database, or transport policy. The monorepo guide is now
unconditional whenever the monorepo module is active.

Acceptance evidence:

- Complete coverage and deletion detection: the catalog requires each
  profile's `requiredRules` to equal its declared module rules; the mutation
  suite rejects a removed required rule and missing or extra dispatch entries.
- Universal Verification and ownership: decision-rendering tests prove the
  selected command appears without autonomous work, setup-owned and
  repository-owned language renders, and frontend points to `DESIGN.md`.
- Exact deterministic dispatch: macro tests derive the expected ordered
  `(skill, trigger)` entries from active modules and compare them byte-for-byte
  with the generated guide for all three profiles.
- Compact roots and complete guides: asset tests keep every root template under
  the compactness limit and require declarative reference tokens; generated
  guides render the rule authority through `artifact.rules`.
- Monorepo and preservation: macro tests prove the single-context
  TypeScript/Bun profile creates both monorepo artifacts and that every profile
  preserves repository-authored extension bytes across an idempotent reapply.
- Distribution parity: `rtk make skills-sync` regenerated the embedded tree and
  `rtk make skills-sync-check` passed with no drift.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py`
  — passed, 8 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_decision_rendering.py`
  — passed, 8 tests.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_assets.py`
  — passed, 11 tests, including required-rule and dispatch mutation checks.
- `rtk make verify` — the pre-feedback local run passed 1,687 Go tests and 90
  setup Python tests, canonical and embedded assets, shipped skills, and the Go
  build. Daemon attempt 1 subsequently exposed generated Python cache drift at
  the skill synchronization gate.
- `rtk make skills-sync-check` — reproduced the Daemon failure before repair,
  then passed after removing only the generated canonical `__pycache__`
  artifact; `diff -qr` confirmed byte-identical skill trees.
- Focused post-repair macro and decision-rendering checks — passed, 8 tests
  each, with bytecode generation disabled so the repaired sync state remains
  clean for the Daemon retry.
- `rtk git diff --check` — passed with no whitespace errors.

Follow-ups: none within this Task's slice. Exact Decision Plan reference audit,
change-plan authority, and external skill restoration remain owned by later
Tasks in this Spec.
