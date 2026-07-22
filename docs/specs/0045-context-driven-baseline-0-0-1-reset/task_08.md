---
task: task_08
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 08: Apply Baseline Readoption through one Change Plan

## Overview

Integrate the new baseline, capability, and disposition contracts into the
existing setup transaction. Preview, confirmation, writes, rollback, snapshot,
and audit must all describe the same immutable Change Plan.

## Requirements

1. MUST incorporate Source Baseline identity, complete entry inventory,
   structured dispositions, capability outcomes, exact outputs, and
   Verification entries into one deterministic Decision Plan and Change Plan.
2. MUST include every normalized decision and planned byte change in the
   timestamp-free plan digest.
3. MUST reject missing decisions, missing required capabilities, and stale
   plan confirmation before performing any write.
4. MUST verify preimages immediately before mutation and postwrite bytes after
   mutation, rolling back the full Change Plan on any failure.
5. MUST apply typed documentation moves and first-time repository-rules
   creation exactly as previewed while preserving existing repository-owned
   content.
6. MUST write only strict 0.0.1 owned schemas and semantic-version markers to
   the Setup Manifest and Setup Snapshot.
7. MUST preserve stable exit codes, stdout/stderr separation, and
   machine-readable output for audit, preview, and apply.
8. MUST make a successful reapply idempotent and make subsequent audit clean.

## Subtasks

- [ ] Extend Decision Plan and Change Plan composition with the new contracts.
- [ ] Extend canonical digest framing for decisions, capabilities, entries,
      outputs, and Verification.
- [ ] Integrate stale-plan and required-capability gates before writes.
- [ ] Apply exact dispositions and repository-owned rule creation atomically.
- [ ] Extend preimage, postwrite, rollback, manifest, and snapshot behavior.
- [ ] Add partial-failure, stale-plan, tampering, and idempotency probes.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Preview and apply report the same plan digest and exact byte changes for
      identical inputs.
- [ ] Missing decisions, missing required capabilities, changed preimages, and
      stale confirmation produce no repository writes.
- [ ] A forced failure after one planned write restores every preimage and
      leaves no partial manifest or snapshot update.
- [ ] Successful apply writes the exact previewed bytes, a strict 0.0.1
      manifest and snapshot, and preserves existing repository-owned rules.
- [ ] Postwrite tampering is detected and rolled back.
- [ ] Reapply produces no content change and audit reports a clean baseline.

## Context

- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- instruction: `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- instruction: `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`
- interface: `.agents/skills/setup-context-driven/tests/test_manifest_migration.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_readoption_apply.py'` — expected: preview/apply parity, no-write rejection, rollback, exact outputs, and idempotent reapply all pass.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_apply.py'` — expected: existing setup transaction and exit-code behavior remain correct.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Features 4–8 and 13; User Stories 1, 3, and 5; Success
  Metrics.
- `_techspec.md` → System Architecture; Implementation Design; Integration
  Points; Build Order 8.
- ADR-0046 → declarative setup-owned instruction rendering.
- ADR-0047 → immutable confirmed Decision Plan effects.
- ADR-0058 → fail-closed retention and atomic upgrade behavior.
- ADR-0064 → complete entry-level Readoption input.
