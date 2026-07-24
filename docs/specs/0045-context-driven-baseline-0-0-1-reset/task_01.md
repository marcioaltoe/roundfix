---
task: task_01
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
type: backend
complexity: high
---

# Task 01: Enforce the 0.0.1 Source Baseline contract

## Overview

Establish the strict, versioned Source Baseline boundary used by every later
setup and release-reset slice. The loader must prove that the corpus,
manifest, and independent index describe one complete project-agnostic source
without activating migration behavior prematurely.

## Requirements

1. MUST introduce the `context_baseline.py` domain types and loaders defined by
   the TechSpec for Source Baselines, entries, manifests, and indexes.
2. MUST accept new owned documents only when their schema identifier and
   semantic version are exactly the `0.0.1` contract for that document kind.
3. MUST validate corpus-to-manifest and manifest-to-index membership in both
   directions, including independently stored entry counts and SHA-256 digests.
4. MUST reject missing, duplicate, unknown, path-unsafe, mixed-version, or
   structurally invalid entries with stable diagnostics.
5. MUST fail when the same corpus member is removed from both the corpus and
   manifest but remains required by the independent index.
6. MUST keep legacy integer-version assets readable only as explicit Baseline
   Readoption input; new plans and outputs MUST NOT emit those versions.
7. MUST keep validation local, deterministic, read-only, network-free, and
   project-agnostic.

## Subtasks

- [x] Add immutable Source Baseline value types and deterministic loaders.
- [x] Define strict `setup-context-driven/<kind>/0.0.1` schema validation.
- [x] Add bidirectional corpus, manifest, and index integrity checks.
- [x] Add independent count and digest verification.
- [x] Create valid fixtures and one-field structural mutations.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] A valid Source Baseline loads to identical normalized values on repeated
      reads.
- [x] Missing or extra entries in any side of the corpus/manifest/index
      relationship fail with stable, actionable diagnostics.
- [x] Removing one entry from both the corpus and manifest still fails against
      the independent index.
- [x] Malformed schema identifiers, non-string versions, mixed versions,
      unsafe paths, duplicate IDs, wrong counts, and wrong digests are rejected.
- [x] Loader execution performs no writes, subprocess execution, or network
      access.
- [x] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/adr/0060-source-baselines-are-exhaustive-and-project-agnostic.md`
- instruction: `docs/adr/0062-roundfix-owned-versions-restart-at-zero.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/assets/contract-v1.json`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_source_baselines.py'` — expected: the valid 0.0.1 Source Baseline loads deterministically and every integrity mutation fails with its asserted diagnostic.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1 and 5; Core Features 1, 4, and 12; User Story 4.
- `_techspec.md` → System Architecture; Implementation Design: Interfaces and
  Data Models; Build Order 1.
- ADR-0060 → exhaustive, indexed, project-agnostic Source Baselines.
- ADR-0062 → strict 0.0.1 owned-version reset.

## Result

Implemented a read-only Source Baseline boundary with frozen normalized value
types, strict `0.0.1` identity, manifest, and independent index schemas, paired
corpus entry markers, safe relative paths, and deterministic SHA-256
verification. The loader rejects incomplete or conflicting membership in every
direction and keeps existing integer-version catalog handling separate from
the new current-generation loader.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_source_baselines.py'`
  — passed 7 tests covering the valid fixture and 23 negative mutation
  scenarios.
- `rtk make skills-sync-check` — passed after synchronizing the canonical and
  distributed setup skill trees.
- `rtk make verify` — passed on the approved rerun: 1,694 Go tests, both
  177-test setup suites, asset loading, owned-skill checks, build, and all
  other repository gates passed. The first managed-sandbox attempt could not
  access the external Go build cache and did not reach a product-code failure.

Acceptance evidence:

- Determinism: `test_valid_source_baseline_loads_deterministically` compared two
  complete loads and proved frozen normalized values.
- Three-way membership:
  `test_membership_mutations_fail_with_stable_diagnostics` covered corpus-only,
  manifest-only, index-only, and corpus-plus-manifest removal cases; each
  mutation produced the same diagnostic sequence on repeated runs.
- Independent index: removing `contract.verification` from both corpus and
  manifest still produced `source-baseline.manifest.entry.missing` because the
  index retained its identity.
- Strict structure and integrity: schema, integer and mixed versions, unsafe
  paths, duplicate IDs, invalid fields and kinds, byte ranges, entry counts,
  and corpus, manifest, and entry digests each failed with their asserted
  stable diagnostic code.
- Read-only boundary:
  `test_source_baseline_loading_has_no_write_command_or_network_side_effects`
  trapped filesystem writes, subprocess execution, and network access while
  confirming fixture bytes remained unchanged.
- Distribution parity: `rtk make skills-sync-check` proved the canonical and distributed trees are byte-identical.

Follow-ups: none for this Task. Publishing the governed production corpus remains
Task 02.
