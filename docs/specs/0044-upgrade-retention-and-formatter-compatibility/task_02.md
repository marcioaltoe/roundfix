---
task: task_02
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
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

- [ ] Migrate the affected workflow rules to clause-based rendering.
- [ ] Restore distinct Spec routing and issue-tracker behavior.
- [ ] Restore findings lifecycle, evidence, and slice guidance.
- [ ] Restore Supervisor authority and Secondbrain operational guidance.
- [ ] Add guide-specific behavioral and duplicate-carrier fixtures.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] Generated Spec routing distinguishes initiative, feature, refactor or
      bug fix, and trivial-change entry points from tracker mechanics.
- [ ] Generated findings guidance includes the copyable frontmatter contract,
      every lifecycle state, recorded reasons, append-only evidence behavior,
      and timestamp updates.
- [ ] General guidance records acceptance evidence and keeps follow-up work
      outside the current slice.
- [ ] Generated autonomous guidance prohibits the Supervisor from writing
      feature code or tests and names the delegation channel.
- [ ] Generated Secondbrain guidance carries the selected query order,
      concrete consult triggers, citation, secret, and read-only boundaries.
- [ ] Two supporting guides in one module cannot render the same effective
      clause list.
- [ ] Canonical and distributed setup skill trees are byte-identical.

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
