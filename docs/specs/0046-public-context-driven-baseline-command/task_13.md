---
task: task_13
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 13: Restore Repository Skill Sets through Baseline

## Overview

Port the maintained skill-restoration operation to
`roundfix baseline skills restore`. Users can preview and confirm immutable
external skill recovery without retaining a Python fallback or adding states
to Baseline adoption.

## Requirements

1. MUST preserve the maintained profile, selected-skill, offline source,
   immutable provenance, lock-digest, preview, confirmation, stale-plan,
   transaction, and rollback contracts.
2. MUST expose the TechSpec command, text/JSON output, stdout/stderr, and exit
   categories under the Baseline command family.
3. MUST group restoration by exact provenance, validate source bytes and
   adapter compatibility before mutation, and update the Repository Skill Set
   plus lock evidence atomically.
4. MUST use the recoverable transaction and exact Plan Digest confirmation for
   every non-empty restoration.
5. MUST reproduce every exact or semantic restoration fixture and designed
   delta recorded in the compatibility matrix.

## Subtasks

- [ ] Port restoration planning and immutable provenance validation.
- [ ] Implement online and declared offline source acquisition boundaries.
- [ ] Integrate exact confirmation, transaction, and rollback.
- [ ] Expose stable text/JSON command results.
- [ ] Add compatibility, stale, unsafe-target, and failure-matrix tests.

## Acceptance Criteria

- [ ] Preview performs zero mutation and returns the exact confirmation digest.
- [ ] Exact offline sources restore without network access.
- [ ] Provenance, source, adapter, target, or lock mismatch fails before mutation.
- [ ] Partial write or lock-update failure restores the complete preimage.
- [ ] Empty restoration is idempotent.
- [ ] Maintained Python restoration cases have Go parity evidence.

## Context

- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `skills-lock.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineSkillsRestore|TestSkillsRestoreOffline|TestSkillsRestoreProvenance|TestSkillsRestoreStalePlan|TestSkillsRestoreRollback|TestSkillsRestoreCompatibility'` — expected: public command, proof, confirmation, idempotence, rollback, and parity cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline skills restore --help` — expected: help matches the approved repository, profile, skill, source, confirmation, and format contract.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 5; Core Features 19–20; Success Metrics.
- `_techspec.md` → API Contracts: maintenance operations; Testing Approach; Build Order 8.
- ADR-0072 → Go destination for maintained restore-skills behavior.
