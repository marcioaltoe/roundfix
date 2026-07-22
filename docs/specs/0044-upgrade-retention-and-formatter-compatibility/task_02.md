---
task: task_02
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
type: docs
complexity: high
---

# Task 02: Restore operational workflow clauses

## Overview

Restore the operational behavior lost when supporting guides were compressed,
using the new clause contracts as the only rendered source. The generated
workflow corpus must answer routing, lifecycle, evidence, authority, and
Secondbrain questions without consulting deleted legacy prose.

## Requirements

1. MUST express every restored instruction as a stable clause with canonical
   enforcement strength and one supporting-guide carrier.
2. MUST restore the Spec entry-point routing contract separately from Task
   tracker hygiene.
3. MUST restore the full findings lifecycle, acceptance-evidence, and slice
   discipline contracts named by the PRD.
4. MUST restore the Supervisor's implementation prohibition and named
   delegation channel without changing autonomous runtime selection.
5. MUST restore the selected Secondbrain query and safety contract without
   broadening write authority.
6. MUST keep root instructions compact and reject identical effective clause
   lists between supporting guides in one module.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [x] Migrate the affected workflow rules to clause-based rendering.
- [x] Restore distinct Spec routing and issue-tracker behavior.
- [x] Restore findings lifecycle, evidence, and slice guidance.
- [x] Restore Supervisor authority and Secondbrain operational guidance.
- [x] Add guide-specific behavioral and duplicate-carrier fixtures.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] Generated Spec routing distinguishes initiative, feature, refactor or
      bug fix, and trivial-change entry points from tracker mechanics.
- [x] Generated findings guidance includes the copyable frontmatter contract,
      every lifecycle state, recorded reasons, append-only evidence behavior,
      and timestamp updates.
- [x] General guidance records acceptance evidence and keeps follow-up work
      outside the current slice.
- [x] Generated autonomous guidance prohibits the Supervisor from writing
      feature code or tests and names the delegation channel.
- [x] Generated Secondbrain guidance carries the selected query order,
      concrete consult triggers, citation, secret, and read-only boundaries.
- [x] Two supporting guides in one module cannot render the same effective
      clause list.
- [x] Canonical and distributed setup skill trees are byte-identical.

## Context

- interface: `.agents/skills/setup-context-driven/assets/modules/context-workflow.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/spec-workflow.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/autonomous-work.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/secondbrain.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_spec_triage_decisions.py`
- interface: `.agents/skills/setup-context-driven/tests/test_secondbrain.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_operational_guidance.py'` — expected: every restored workflow clause renders through its intended guide and duplicate guide clause sets are rejected.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_spec_triage_decisions.py` — expected: Spec routing decisions generate the complete routing contract.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_secondbrain.py` — expected: the Secondbrain decision generates the complete operational contract.
- `rtk make verify` — expected: the full repository gate passes with synchronized skill trees.

## References

- `_prd.md` → User Story 3; Core Features 5–6; Success Metrics.
- `_techspec.md` → Data Models; Testing Approach; Build Order 2.
- ADR-0046 → setup-owned supporting-guide boundaries.

## Result

Restored the compressed operational workflow corpus as module-v3 clauses.
Generated guides now render explicit `mandatory`, `prohibited`, or
`stop-and-ask` labels from the clause contract; the renderer no longer relies
on free-text modal strength. Catalog loading also rejects both multiply carried
rules and supporting guides with identical effective clause lists.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_operational_guidance.py'` — passed, 2 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_spec_triage_decisions.py` — passed, 2 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_secondbrain.py` — passed, 5 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_*.py'` — passed, 135 tests.
- `rtk diff -qr .agents/skills/setup-context-driven skills/setup-context-driven` — passed with no differences after `rtk make skills-sync`.
- `rtk make verify` — passed after granting the required existing Go build-cache access: 1,694 Go tests, 135 canonical setup tests, 135 distributed setup tests, asset loading, Roundfix skill checks, and the Go build all succeeded.

Acceptance evidence:

- `test_selected_workflow_guides_render_complete_distinct_contracts` proves the generated Spec routing matrix is present only in `spec-routing.md`, while tracker mechanics remain in `issue-tracker.md`.
- The same generated-guide test proves the findings frontmatter, all four lifecycle states, recorded reasons, append-only addenda, and `updated_at` rule render from the docs-layout clause carrier.
- The generated agent-instructions assertions prove every acceptance criterion records fresh evidence and follow-up work remains outside the active slice.
- The autonomous assertions prove the Supervisor cannot write feature code or tests, delegates through a Roundfix Run, and retains the selected backend and design runtime values.
- The Secondbrain assertions prove index-first lookup, the exact `qmd query "<question>" --all --files --min-score 0.3` command, mirror ordering, concrete consult triggers, citations, secret safety, and read-only boundaries.
- `test_duplicate_effective_clause_lists_are_rejected` proves a duplicate guide carrier fails with `guide.clause-list.duplicate`; catalog validation also enforces one carrier per selected rule with `profile.rule.carrier.duplicate`.
- The recursive tree comparison and both mirrored setup suites prove the canonical and distributed setup skill trees are byte-identical.

Follow-ups: none for this Task slice.
