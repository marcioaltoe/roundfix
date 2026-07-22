---
task: task_03
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: docs
complexity: high
---

# Task 03: Complete portable hard-rule coverage

## Overview

Complete the portable rule corpus and the reviewed legacy hard-rule ledger.
This slice must retain useful safety and delivery behavior while recording why
obsolete, conflicting, or project-specific sample clauses stay outside the
Context-Driven Baseline.

## Requirements

1. MUST add the missing portable hard rules listed in PRD Core Feature 4 with
   stable clause identities and explicit enforcement strength.
2. MUST maintain a reviewed ledger for prior managed prose and the supported
   legacy TypeScript/Bun sample without copying project branding or obsolete
   assumptions into generated guidance.
3. MUST map every accepted prior clause to a current clause with equivalent
   enforcement and give every excluded clause a recorded reason.
4. MUST keep stack-specific guidance conditional on confirmed profile
   capabilities or classify it outside the baseline.
5. MUST prove obsolete knowledge flow, runtime defaults, product assumptions,
   and conflicting API guidance remain absent.
6. MUST preserve repository-authored policy outside setup markers and keep
   canonical and distributed setup skill trees synchronized.

## Subtasks

- [ ] Add portable verification, design, testing, decision, research, and
      delivery clauses.
- [ ] Build the reviewed prior-prose and legacy-sample ledger.
- [ ] Map accepted clauses to current enforcement contracts.
- [ ] Record reasons for every excluded or repository-specific clause.
- [ ] Add weakened, missing, and obsolete-clause regression fixtures.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] Generated guidance blocks selected warnings where the profile declares
      warnings as errors and gates intentional Verification configuration
      changes on explicit authority.
- [ ] Frontend and testing guidance require the selected design contract and
      dependent interfaces to be read before governed work.
- [ ] General guidance asks user-answerable decisions, routes external
      research through the declared fallback, and excludes external research
      from local-code discovery.
- [ ] Commit and delivery work dispatches to the governing installed skills.
- [ ] Every ledger clause has exactly one accepted mapping or one exclusion
      reason, and a weakened accepted mapping fails validation.
- [ ] Fixtures prove retired or conflicting sample behavior never appears in
      generated output.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `CONTEXT.md`
- interface: `.agents/skills/setup-context-driven/assets/modules/core.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/frontend.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/typescript.json`
- interface: `.agents/skills/setup-context-driven/assets/modules/bun.json`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_rendering.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_legacy_rule_ledger.py'` — expected: all prior clauses are accounted, accepted clauses retain enforcement, and excluded behavior stays absent.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_decision_rendering.py` — expected: profile decisions render the complete portable hard-rule set without unsafe repository-specific content.
- `rtk make verify` — expected: the full repository gate passes with synchronized skill trees.

## References

- `_prd.md` → User Story 3; Core Features 3–4 and 6; Non-Goals / Out of Scope;
  Success Metrics.
- `_techspec.md` → Data Models; Testing Approach; Build Order 2.
- ADR-0058 → explicit accounting for prior mandatory clauses.
