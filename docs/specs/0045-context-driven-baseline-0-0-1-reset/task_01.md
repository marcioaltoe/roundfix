---
task: task_01
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
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

- [ ] Add immutable Source Baseline value types and deterministic loaders.
- [ ] Define strict `setup-context-driven/<kind>/0.0.1` schema validation.
- [ ] Add bidirectional corpus, manifest, and index integrity checks.
- [ ] Add independent count and digest verification.
- [ ] Create valid fixtures and one-field structural mutations.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] A valid Source Baseline loads to identical normalized values on repeated
      reads.
- [ ] Missing or extra entries in any side of the corpus/manifest/index
      relationship fail with stable, actionable diagnostics.
- [ ] Removing one entry from both the corpus and manifest still fails against
      the independent index.
- [ ] Malformed schema identifiers, non-string versions, mixed versions,
      unsafe paths, duplicate IDs, wrong counts, and wrong digests are rejected.
- [ ] Loader execution performs no writes, subprocess execution, or network
      access.
- [ ] Canonical and distributed setup skill trees are byte-identical.

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
