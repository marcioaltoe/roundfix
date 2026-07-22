---
task: task_09
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: docs
complexity: medium
---

# Task 09: Synchronize the shipped setup contract

## Overview

Align the shipped skill instructions, maintainer guidance, generated behavior,
and distributed copy with the completed upgrade contract. A reader must be
able to operate retention, repository extensions, delegation findings, skill
dispatch, and formatter verification without consulting the implementation.

## Requirements

1. MUST document the Upgrade Retention Contract, normative Change Plan output,
   blocking failures, and existing confirmation flow in the setup skill.
2. MUST document Repository-Owned Extension ownership and the informational
   baseline-floor finding without implying setup can manage project content.
3. MUST document one-entry skill dispatch and formatter compatibility,
   including the distinction between hermetic Verification and the real final
   QA probe.
4. MUST keep CLI examples truthful to implemented arguments, output fields,
   exit codes, stdout/stderr behavior, and offline guarantees.
5. MUST update maintainer guidance for transition ledgers, formatter
   provenance refresh, and canonical-to-distributed synchronization.
6. MUST leave both setup skill trees byte-identical and the complete repository
   Verification green.

## Subtasks

- [ ] Update the canonical setup skill workflow and examples.
- [ ] Update the Context-Driven development user guide.
- [ ] Document ledger and formatter-provenance maintenance.
- [ ] Reconcile help text and documented exit/output behavior.
- [ ] Synchronize the distributed setup skill copy.
- [ ] Run the complete setup and repository gates.

## Acceptance Criteria

- [ ] Skill and user guidance explain retention accounting, authorization, and
      blocking failure recovery with implemented command forms.
- [ ] Guidance states that Repository-Owned Extension bytes remain unmarked
      and outside setup management.
- [ ] Guidance explains informational delegation findings as a baseline-floor
      signal that does not block apply.
- [ ] Guidance states that each installed skill renders once and that
      formatter proof is pinned and profile-specific.
- [ ] Maintainer instructions identify the canonical source, sync check,
      transition-ledger review, and formatter-provenance refresh boundary.
- [ ] Canonical and distributed skill trees and test suites describe identical
      behavior.
- [ ] The complete repository gate passes without warnings or failures.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `Makefile`

## Verification

- `rtk grep -F 'Upgrade Retention Contract' .agents/skills/setup-context-driven/SKILL.md docs/user-guide/context-driven-development.md` — expected: both operator-facing documents name the canonical retention behavior.
- `rtk grep -F 'Formatter-Stable Output' .agents/skills/setup-context-driven/SKILL.md docs/user-guide/context-driven-development.md` — expected: both documents name the formatter compatibility contract.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make setup-context-check` — expected: canonical and distributed setup suites and asset catalogs pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → All Goals; Core Feature 12; User Experience; Success Metrics.
- `_techspec.md` → System Architecture; API Contracts; Integration Points;
  Testing Approach; Build Order 6.
- ADR-0058 → operator-visible retention authorization.
- ADR-0059 → operator-visible formatter compatibility.
