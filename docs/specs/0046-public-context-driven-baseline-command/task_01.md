---
task: task_01
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Inventory the maintained Python contract and test cases.
- [x] Define the parity matrix and deterministic fixture schema.
- [x] Generate representative preimages, inputs, outputs, and refusal fixtures.
- [x] Add characterization tests for fixture completeness and stability.
- [x] Synchronize the transition corpus across the canonical and distributed skill suites.

## Acceptance Criteria

- [x] Every maintained Python behavior has one classified matrix row and a Go destination.
- [x] All supported profiles and adoption states have deterministic fixtures.
- [x] Fixture regeneration is byte-stable on an unchanged source tree.
- [x] No retired or designed-delta row lacks a written rationale.
- [x] The corpus records exact Plan Digests and rollback outcomes where the contract requires exact parity.
- [x] Canonical and distributed transition suites validate the same corpus.

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

## Result

Implemented a standalone, content-addressed transition corpus that future Go
tests can consume without importing the Python runtime or test suite. The
matrix freezes 240 Python tests plus three explicit product deltas: 163 exact,
63 semantic, 16 designed-delta, one ancillary, and zero retired rows.

The 14 deterministic fixtures cover all three maintained Baseline Profiles and
greenfield, update, preservation, profile change, stale input, unsafe carrier,
capability refusal, Formatter-Stable Output composition, skill restoration,
asset synchronization, and rollback states. Each fixture records normalized
input and output, repository file identities, referenced content blobs,
managed-entry evidence, manifests when applicable, planned byte order, exact
Plan Digests when parity requires them, and post-state or rollback evidence.

Acceptance evidence:

- Every maintained Python behavior has one classified row and Go destination:
  `test_matrix_covers_every_python_test_once_with_go_destinations` matched all
  240 pre-cutover test methods exactly once.
- Supported profiles and adoption states have deterministic fixtures:
  `test_fixtures_cover_profiles_states_digests_and_rollback` validated the
  three profiles and every required state.
- Regeneration is byte-stable:
  `test_regeneration_is_byte_identical_and_manifest_hashes_every_artifact`
  regenerated all 17 corpus files and matched every checked-in byte and
  artifact digest.
- Every non-exact row has a rationale, including all designed deltas; the
  matrix records no retired behavior.
- Exact Plan Digests are recorded for greenfield, Readoption, stale-plan,
  restoration, and rollback fixtures. The `atomic-rollback` and
  `skill-restoration-rollback` fixtures prove that paths, kinds, targets, and
  exact bytes return to their preimage.
- `test_canonical_and_distributed_suites_validate_identical_corpus` and
  `rtk make skills-sync-check` proved the canonical and distributed skill
  surfaces are byte-identical.

Verification:

- Pre-change
  `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_parity_corpus.py'`:
  failed with exit 5 because zero tests existed.
- Post-change, the same focused command: passed, 4 tests.
- `rtk make setup-context-check`: passed across canonical and distributed
  setup suites.
- `rtk make verify`: the sandboxed attempt could not read the host Go build
  cache; the unchanged rerun with cache access passed 1,727 Go tests, 256
  canonical setup tests, 256 distributed setup tests, asset loading, Roundfix
  skill checks, and the build.
