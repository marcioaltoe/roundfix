---
task: task_06
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: chore
complexity: high
---

# Task 06: Declare the two-tier verification contract

## Overview

Tooling Task two of two, under the same authorization. The gate's mechanical
stage runs the repository gate cold and complete at about ninety seconds when
the same gate on an unchanged tree costs under five. The fix is a declared
contract every adopting repository inherits: local verification is incremental
and fast, CI judges the complete tree from a fresh run, and the clause says
which tier answers which question.

The contract is expressed per profile in terms of the commands that profile
declares. An adopting repository may be a Bun or Rust repository with no
`make` at all, so a clause naming a tool would be unmeetable by construction.

## Requirements

1. MUST author the two-tier clause in the Baseline modules the authorization
   bounds, stating that the local tier is incremental, the CI tier is complete
   and fresh, and which question each answers.
2. MUST express the tiers per profile as declared commands, never as a named
   tool, so a repository without `make` can satisfy the contract.
3. MUST state that a profile declaring no incremental command inherits the
   clause as unmet rather than silently satisfied.
4. MUST adopt the regenerated managed-block postimages into the corresponding
   `docs/agents/` guides.
5. MUST add this repository's own incremental verification target to the
   `Makefile` under the name `verify-incremental`, and MUST NOT change what
   `make verify` means for CI. The name is fixed here because the declared
   Verification has to be able to fail: `verify:` already exists, so asserting
   that some verify-shaped target is present approves this Task before any work
   happens.
6. MUST run the module chain per the measured choreography: bootstrap the
   Source Baseline manifest rows for the new clauses, then run
   `make baseline-digests` twice, since the maintained fixture is the chain's
   first step and converges only on the second pass.
7. MUST NOT correct the maintained Source Baseline expectation this change
   invalidates. That correction is task_07's, landing as its own commit after
   this one, because a consequent fix folded into an authorized tooling commit
   fails the tooling-authority gate.
8. MUST change only the authorization's bounded files, their sanctioned
   deterministic digest fallout, and this task file.

## Subtasks

- [ ] Author the per-profile two-tier clause in the modules.
- [ ] Add this repository's incremental target.
- [ ] Run the chain and adopt both postimages.

## Acceptance Criteria

- [ ] Both guides carry the two-tier clause expressed as declared commands.
- [ ] A profile with no incremental command reads as unmet, not satisfied.
- [ ] `make verify` still means the complete gate.
- [ ] The digest chain is converged.
- [ ] The maintained expectation is left for task_07, and the diff stays
      inside the bounded files, fallout, and this task file.

## Context

- instruction: docs/workflow/authorizations/2026-08-06-proof-cost.md
- interface: internal/baseline/assets/modules/core.json
- interface: internal/baseline/assets/modules/spec-workflow.json
- interface: Makefile

## Verification

- `grep -rqi 'incremental' docs/agents/ && grep -rqi 'complete' docs/agents/`
  — expected: exit 0; the two tiers are adopted into the guides.
- `grep -q '^verify-incremental:' Makefile`
  — expected: exit 0; the incremental target exists beside an unchanged
  `verify`.
  — expected: exit 0; the digest chain is converged.
  — expected: exit 0; the consequent fixture correction was not folded into
  this commit.
  — expected: exit 0; nothing outside the bounded files and sanctioned fallout
  changed.

## References

- `_prd.md` → Core Feature 6; User Story 5; Goals 3 and 4.
- `_techspec.md` → Build Order 6; Decisions (per-profile commands).
- ADR-0081.
- `references/2026-08-03-verification-performance-contract.md` → the measured
  tier costs this clause is built on.
