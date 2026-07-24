---
task: task_09
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Update the canonical setup skill workflow and examples.
- [x] Update the Context-Driven development user guide.
- [x] Document ledger and formatter-provenance maintenance.
- [x] Reconcile help text and documented exit/output behavior.
- [x] Synchronize the distributed setup skill copy.
- [x] Run the complete setup and repository gates.

## Acceptance Criteria

- [x] Skill and user guidance explain retention accounting, authorization, and
      blocking failure recovery with implemented command forms.
- [x] Guidance states that Repository-Owned Extension bytes remain unmarked
      and outside setup management.
- [x] Guidance explains informational delegation findings as a baseline-floor
      signal that does not block apply.
- [x] Guidance states that each installed skill renders once and that
      formatter proof is pinned and profile-specific.
- [x] Maintainer instructions identify the canonical source, sync check,
      transition-ledger review, and formatter-provenance refresh boundary.
- [x] Canonical and distributed skill trees and test suites describe identical
      behavior.
- [x] The complete repository gate passes without warnings or failures.

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

## Result

Updated the operator workflow to make the Upgrade Retention Contract part of
Change Plan review and confirmation, including ordered `retentionAccounting`,
blocking failure recovery, exit categories, stdout/stderr ownership, and the
offline boundary. The same guidance now states the one-time unmarked
Repository-Owned Extension boundary, the informational
`delegation.baseline-floor` signal, one rendered entry per installed skill,
and the difference between hermetic Verification and the real pinned final-QA
formatter probe.

Maintainer guidance now identifies `.agents/skills/setup-context-driven/` as
the canonical authorial source, gives focused transition-ledger checks, defines
when formatter provenance and the complete golden corpus must be refreshed,
and requires `make skills-sync` plus `make skills-sync-check` before delivery.
The repository workflow regenerated the complete distributed skill copy from
that canonical tree.

Acceptance evidence:

- Retention accounting, authorization, and recovery: both required
  `rtk grep -F 'Upgrade Retention Contract' ...` checks found the named
  contract in the canonical skill and user guide; the documented commands and
  exit/output behavior were reconciled against fresh `context_setup.py`,
  `audit`, and `apply` help output and the retention tests.
- Repository ownership: both operator documents state that the extension is
  unmarked, excluded from `managedArtifacts`, and never compared, rewritten,
  formatted, removed, or automatically recreated by setup.
- Delegation floor: both operator documents state that
  `delegation.baseline-floor` is informational, does not affect the Change Plan
  or exit status, and does not block apply.
- Dispatch and formatter compatibility: both operator documents state that
  every installed skill renders once and that Formatter-Stable Output uses a
  pinned, profile-specific proof with hermetic Verification and a separate real
  final-QA probe. The required
  `rtk grep -F 'Formatter-Stable Output' ...` checks found the named contract in
  both documents.
- Maintainer boundary: the canonical maintenance reference names the
  transition-ledger review triggers, focused tests, formatter-provenance refresh
  boundary, canonical source, and sync commands.
- Canonical/distributed parity: `rtk make skills-sync-check` passed after
  regeneration. `rtk make setup-context-check` passed 170 canonical tests, 170
  distributed tests, and both asset-catalog loads.
- Repository Verification: `rtk make verify` passed with 1,694 Go tests, both
  170-test setup suites, both asset catalogs, the shipped Roundfix skill check,
  and the Go build. The first sandboxed attempt could not read the existing Go
  build cache; rerunning the same command with cache access exited `0` with no
  warnings or failures.

Follow-ups: none.
