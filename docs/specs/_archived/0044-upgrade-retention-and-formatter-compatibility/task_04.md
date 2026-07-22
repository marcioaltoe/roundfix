---
task: task_04
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
type: backend
complexity: high
---

# Task 04: Render one dispatch entry per installed skill

## Overview

Make skill dispatch a normalized catalog contract instead of a renderer
deduplication heuristic. Every supported profile must render one actionable
entry for each active installed skill while retaining only genuinely distinct
capability-proven triggers.

## Requirements

1. MUST assign one catalog owner to each skill in a selected profile's active
   module closure.
2. MUST identify distinct triggers with stable trigger IDs and reject
   duplicate ownership or duplicate trigger intent before rendering.
3. MUST move shared skills to the least common owning module instead of
   re-declaring them through dependent modules.
4. MUST add domain dispatch only where the selected profile and Repository
   Skill Set prove the capability.
5. MUST render exactly one top-level entry per skill without semantic merging
   in the renderer.
6. MUST reject any profile whose installed, required, and dispatch skill sets
   do not agree.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [x] Normalize active dispatch contracts by skill and stable trigger ID.
- [x] Reassign shared workflow and generic skill ownership.
- [x] Add capability-proven technology and surface triggers.
- [x] Simplify rendering to consume the validated normalized contracts.
- [x] Add duplicate, near-duplicate, missing, and all-profile fixtures.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] Every bundled profile renders exactly one top-level entry for every
      active required skill.
- [x] Genuinely distinct triggers for one skill appear within that one entry
      and retain stable ordering.
- [x] A dependent module's duplicate skill owner or reused trigger ID fails
      asset validation even when its wording differs.
- [x] A required skill absent from the selected setup snapshot or normalized
      dispatch map fails validation.
- [x] Framework-specific dispatch remains absent unless a selected profile
      explicitly proves that capability.
- [x] Repeated rendering is byte-identical and performs no heuristic text
      similarity analysis.
- [x] Canonical and distributed setup skill trees are byte-identical.

## Context

- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/assets/modules/core.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/context-workflow.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/spec-workflow.json`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skill_dispatch.py'` — expected: ownership and trigger conflicts fail, while each valid profile renders one entry per skill.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py` — expected: every bundled profile's required, installed, and rendered skill sets agree.
- `rtk make verify` — expected: the full repository gate passes with synchronized skill trees.

## References

- `_prd.md` → User Story 4; Core Feature 7; Success Metrics.
- `_techspec.md` → System Architecture; Data Models; Testing Approach; Build
  Order 2.
- ADR-0047 → one declarative source for preview, audit, and apply.

## Result

Implemented normalized dispatch ownership as an asset-loading invariant. Each
selected profile now requires exact agreement between its setup snapshot,
active modules' required skills, and normalized dispatch map. Shared skills
live in `core` or `context-workflow`; profile-specific language, backend,
frontend, CLI, and TUI capabilities remain owned only by modules selected for
profiles that install them. Rendering consumes the validated map and emits one
top-level skill entry with stable trigger-ID children.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skill_dispatch.py'` — passed; 9 tests cover normalized rendering, duplicate owners, reused trigger IDs, duplicate declarations, missing sets, framework isolation, and all bundled profiles.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py` — passed; 8 macro tests cover apply, audit, reapply, and exact installed/required/rendered set agreement for every bundled profile.
- `rtk make verify` — the restricted-sandbox attempt could not access the external Go build cache; the approved rerun passed both 148-test setup suites, 1,694 Go tests across 20 packages, canonical and distributed asset loading, Roundfix skill checks, and the build.

Acceptance evidence:

1. The all-profile focused and macro tests compare setup, required, normalized,
   and rendered top-level skill sets exactly.
2. `test_distinct_triggers_render_under_one_stably_ordered_skill_entry` proves
   distinct trigger IDs remain ordered beneath one skill entry; the bundled
   `implement-task` and `implement-spec` contracts also retain execution and
   autonomous-delegation triggers under their single owners.
3. `test_dependent_module_cannot_redeclare_a_skill_owner`,
   `test_reused_trigger_id_fails_even_when_wording_differs`, and
   `test_duplicate_skill_contract_in_one_module_is_rejected` prove ownership
   and stable-intent conflicts fail during asset loading.
4. The missing-setup, missing-dispatch, and unowned-installed-skill tests prove
   exact set disagreements fail with stable diagnostics.
5. `test_framework_dispatch_only_renders_for_profiles_that_install_it` proves
   React dispatch remains TypeScript-only and Go dispatch remains Go-only.
6. Repeated fixture rendering is byte-identical, and `render_skill_dispatch`
   now formats only the validated normalized map without wording comparison or
   similarity analysis.
7. `rtk make skills-sync` refreshed the distributed copy; the full gate's
   byte-parity and dual-tree suites passed.

Follow-ups: none.
