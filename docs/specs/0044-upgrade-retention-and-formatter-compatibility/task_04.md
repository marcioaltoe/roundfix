---
task: task_04
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
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

- [ ] Normalize active dispatch contracts by skill and stable trigger ID.
- [ ] Reassign shared workflow and generic skill ownership.
- [ ] Add capability-proven technology and surface triggers.
- [ ] Simplify rendering to consume the validated normalized contracts.
- [ ] Add duplicate, near-duplicate, missing, and all-profile fixtures.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] Every bundled profile renders exactly one top-level entry for every
      active required skill.
- [ ] Genuinely distinct triggers for one skill appear within that one entry
      and retain stable ordering.
- [ ] A dependent module's duplicate skill owner or reused trigger ID fails
      asset validation even when its wording differs.
- [ ] A required skill absent from the selected setup snapshot or normalized
      dispatch map fails validation.
- [ ] Framework-specific dispatch remains absent unless a selected profile
      explicitly proves that capability.
- [ ] Repeated rendering is byte-identical and performs no heuristic text
      similarity analysis.
- [ ] Canonical and distributed setup skill trees are byte-identical.

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
