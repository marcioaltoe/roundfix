---
task: task_03
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Add immutable bundle and trigger values to the asset boundary.
- [x] Define the governed production, endpoint, persistence, and specialist
      bundles.
- [x] Render exact trigger-to-bundle instructions.
- [x] Persist normalized bundle membership in snapshots.
- [x] Add membership, ordering, duplication, and unknown-skill mutations.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] Production-code output names exactly the three required production
      skills once each.
- [x] A Hono endpoint without persistence names exactly the three endpoint
      skills, while a persistence endpoint adds `drizzle-orm` once.
- [x] Frontend, testing, security, and delivery triggers render their own
      declared bundles without inheriting unrelated skills.
- [x] Invalid or incomplete bundle documents fail before any repository write.
- [x] Snapshot output records stable bundle IDs and ordered exact membership.
- [x] Repeated rendering is byte-identical.

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

## Result

Added a declarative Skill Activation catalog with frozen bundle and trigger
values. Asset loading now validates stable ordering, unique ownership and
trigger IDs, exact non-duplicate membership, known skills, active-profile
snapshot coverage, and optional capability-condition syntax before rendering.
The Setup Snapshot digest now covers both normalized skill records and the
ordered active bundle records.

Generated skill-dispatch guidance now renders every governed trigger once with
its stable bundle ID and ordered exact membership. The production, Hono,
persistence, React, UI quality, testing, debugging, security, QA, and delivery
bundles remain separate from the existing per-skill trigger index. The
canonical and distributed setup skill trees contain the same catalog, code,
tests, snapshots, and formatter golden evidence.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skill_dispatch.py'`
  — passed 12 tests after the final synchronization. The suite covers exact
  production and endpoint membership, specialist separation, snapshot
  persistence, byte-identical repeated rendering, and missing-member,
  duplicate-member, duplicate-trigger, duplicate-owner, unknown-skill, and
  unstable-order mutations without repository writes.
- `rtk bunx oxfmt@0.59.0 --check AGENTS.md docs/agents` — passed against a
  freshly applied disposable TypeScript/Bun fixture; generated Markdown was
  unchanged by the pinned formatter.
- `rtk make skills-sync-check` and recursive tree comparison — passed; the
  canonical and distributed setup skill trees are byte-identical.
- `rtk git diff --check` — passed.
- `rtk make verify` — passed on the unchanged elevated rerun after the final
  synchronized implementation. The initial sandboxed attempt failed only
  because the host Go build cache was outside the worktree sandbox.

Acceptance evidence:

- `test_typescript_profile_renders_exact_activation_bundles_deterministically`
  asserts the production bundle is exactly `coding-guidelines`, `clean-code`,
  and `solid`; the endpoint bundle is exactly `hono-api-best-practices`,
  `hono`, and `zod`; and the persistence variant adds `drizzle-orm` once.
- `test_specialist_activation_bundles_remain_distinct_and_snapshot_backed`
  proves frontend, UI quality, testing, security, and delivery use distinct
  trigger-to-bundle identities and that the TypeScript Setup Snapshot records
  the normalized ordered membership.
- `test_invalid_activation_documents_fail_before_repository_writes` proves
  incomplete, duplicated, unknown, conflicting-owner, and reordered contracts
  fail during local asset loading while the write boundary is disabled.
- The focused test renders twice and compares bytes; the real formatter probe,
  macro-profile suite, canonical asset suite, full canonical setup suite, and
  repository Verification also passed.

Follow-ups: none for this Task slice.
