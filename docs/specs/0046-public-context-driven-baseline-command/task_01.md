---
task: task_01
spec: 0046-public-context-driven-baseline-command
status: pending
type: test
complexity: high
---

# Task 01: Freeze the Python compatibility corpus

## Overview

Create the executable characterization contract that every Go slice must
satisfy before the Python runtime can be removed. The corpus makes maintained
behavior, intended deltas, and retired behavior reviewable instead of relying
on an implicit reading of the existing test suite.

## Requirements

1. MUST classify every maintained Python input, state, action, refusal, digest,
   planned byte sequence, and rollback outcome as exact, semantic,
   designed-delta, ancillary, or retired.
2. MUST generate deterministic fixtures containing inputs, normalized outputs,
   file identities, managed-entry ledgers, manifests, Plan Digests, and
   post-state or rollback evidence.
3. MUST cover every maintained Baseline Profile, greenfield, preservation,
   update, Readoption, profile change, stale input, unsafe carrier, capability,
   formatter-composition, skill-restoration, and asset-sync state.
4. MUST require an explicit rationale for every designed-delta or retired row.
5. MUST leave a standalone checked-in corpus that remains usable after the
   Python scripts and tests are removed.

## Subtasks

- [ ] Inventory the maintained Python contract and test cases.
- [ ] Define the parity matrix and deterministic fixture schema.
- [ ] Generate representative preimages, inputs, outputs, and refusal fixtures.
- [ ] Add characterization tests for fixture completeness and stability.
- [ ] Synchronize the transition corpus across the canonical and distributed skill suites.

## Acceptance Criteria

- [ ] Every maintained Python behavior has one classified matrix row and a Go destination.
- [ ] All supported profiles and adoption states have deterministic fixtures.
- [ ] Fixture regeneration is byte-stable on an unchanged source tree.
- [ ] No retired or designed-delta row lacks a written rationale.
- [ ] The corpus records exact Plan Digests and rollback outcomes where the contract requires exact parity.
- [ ] Canonical and distributed transition suites validate the same corpus.

## Context

- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_baseline.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_capabilities.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_parity_corpus.py'` — expected: the matrix is complete, fixtures regenerate byte-identically, and every non-exact row has a rationale.
- `rtk make setup-context-check` — expected: canonical and distributed setup suites accept the same characterization corpus.
- `rtk make verify` — expected: the repository gate passes with the new test-only contract.

## References

- `_prd.md` → Goals 2–5; User Stories 2–4 and 7–8; Core Features 19–20; Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 1; Risks & Considerations.
- ADR-0072 → full maintained Python contract preservation.
