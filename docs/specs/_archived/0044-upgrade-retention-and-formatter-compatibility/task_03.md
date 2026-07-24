---
task: task_03
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Add portable verification, design, testing, decision, research, and
      delivery clauses.
- [x] Build the reviewed prior-prose and legacy-sample ledger.
- [x] Map accepted clauses to current enforcement contracts.
- [x] Record reasons for every excluded or repository-specific clause.
- [x] Add weakened, missing, and obsolete-clause regression fixtures.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] Generated guidance blocks selected warnings where the profile declares
      warnings as errors and gates intentional Verification configuration
      changes on explicit authority.
- [x] Frontend and testing guidance require the selected design contract and
      dependent interfaces to be read before governed work.
- [x] General guidance asks user-answerable decisions, routes external
      research through the declared fallback, and excludes external research
      from local-code discovery.
- [x] Commit and delivery work dispatches to the governing installed skills.
- [x] Every ledger clause has exactly one accepted mapping or one exclusion
      reason, and a weakened accepted mapping fails validation.
- [x] Fixtures prove retired or conflicting sample behavior never appears in
      generated output.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Completed the portable hard-rule corpus and two declarative Upgrade Retention
Contract ledgers: one for the prior managed v2 prose and one for the reviewed
legacy TypeScript/Bun sample. Accepted clauses now target stable current clause
IDs with equal enforcement. Every excluded sample clause records why obsolete,
conflicting, or repository-specific behavior stays outside the Context-Driven
Baseline.

Capability-gated TypeScript/Bun guidance now blocks warnings when the selected
profile treats them as errors, requires `DESIGN.md` before UI work, inspects
dependent interfaces before tests, and declares the external web-research
fallback. General guidance now stops for user-answerable decisions and
intentional Verification configuration changes, prohibits external research
for local code, and dispatches commit and pull request delivery through
`conventional-commits` and `github-pr-workflow`. Setup still preserves
repository-authored policy outside managed markers.

Verification evidence:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_legacy_rule_ledger.py'` — passed, 3 tests. Proved complete one-to-one accounting, equivalent enforcement, missing-mapping and weakened-target rejection, excluded-prose absence, and repository-policy preservation.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_decision_rendering.py` — passed, 9 tests. Proved the warning, authority, design, dependent-interface, decision, research, and delivery dispatch clauses render for the TypeScript/Bun profile.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_upgrade_contracts.py` — passed, 6 tests. Proved the canonical ledgers load and canonical/distributed trees are byte-identical.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests -p 'test_legacy_rule_ledger.py'` — passed, 3 tests in the distributed copy.
- `rtk make verify` — passed: 1,694 Go tests; 139 canonical and 139 distributed setup-context-driven tests; both asset catalogs loaded; the Roundfix skill check passed; and the CLI build completed.

Acceptance evidence:

- Warning-as-error and Verification-configuration authority are asserted by `test_typescript_bun_profile_renders_complete_portable_hard_rules`.
- Design-contract and dependent-interface preconditions are asserted in the same generated-profile test.
- User-decision, external-fallback, and local-code research boundaries are asserted in the generated general guidance.
- Commit and delivery dispatch entries for both governing skills are asserted in the rendered dispatch guide.
- `test_every_reviewed_clause_has_one_mapping_with_equivalent_enforcement` checks every ledger entry, while `test_missing_mapping_and_weakened_accepted_mapping_fail_validation` proves both negative cases fail asset validation.
- `test_retired_sample_behavior_is_excluded_and_repository_policy_survives` proves retired knowledge flow, runtime defaults, product assumptions, API policy, and project layout remain absent without overwriting repository-authored policy.
- `test_canonical_and_distributed_setup_skill_trees_are_byte_identical` and the full gate prove mirror synchronization.

Follow-ups: none for this Task slice.
