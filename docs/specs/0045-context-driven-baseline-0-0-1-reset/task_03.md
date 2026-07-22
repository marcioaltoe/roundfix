---
task: task_03
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 03: Render exact skill activation bundles

## Overview

Make Skill Activation a first-class, validated setup contract instead of
loosely worded guidance. Generated instructions must name the exact bundle for
each governed trigger and snapshots must preserve that membership precisely.

## Requirements

1. MUST implement the typed Skill Activation bundle model and validation
   contract defined by the TechSpec.
2. MUST render the exact production-code bundle containing
   `coding-guidelines`, `clean-code`, and `solid`.
3. MUST render the exact Hono-endpoint bundle containing
   `hono-api-best-practices`, `hono`, and `zod`.
4. MUST extend the endpoint bundle with `drizzle-orm` when persistence behavior
   is in scope.
5. MUST keep frontend, testing, security, and delivery triggers distinct and
   explicit rather than relying on broad inference.
6. MUST reject missing members, duplicate members, duplicate trigger IDs,
   unknown skills, and unstable ordering before rendering.
7. MUST persist normalized activation bundle identity and exact membership in
   the Setup Snapshot.

## Subtasks

- [ ] Add immutable bundle and trigger values to the asset boundary.
- [ ] Define the governed production, endpoint, persistence, and specialist
      bundles.
- [ ] Render exact trigger-to-bundle instructions.
- [ ] Persist normalized bundle membership in snapshots.
- [ ] Add membership, ordering, duplication, and unknown-skill mutations.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Production-code output names exactly the three required production
      skills once each.
- [ ] A Hono endpoint without persistence names exactly the three endpoint
      skills, while a persistence endpoint adds `drizzle-orm` once.
- [ ] Frontend, testing, security, and delivery triggers render their own
      declared bundles without inheriting unrelated skills.
- [ ] Invalid or incomplete bundle documents fail before any repository write.
- [ ] Snapshot output records stable bundle IDs and ordered exact membership.
- [ ] Repeated rendering is byte-identical.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_skill_dispatch.py`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skill_dispatch.py'` — expected: each governed trigger renders its exact ordered bundle and invalid membership is rejected.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Feature 9; User Story 4; Success Metrics.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; Coverage
  Map; Build Order 3.
- ADR-0060 → activation language belongs to the governed corpus.
