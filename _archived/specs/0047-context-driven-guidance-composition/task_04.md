---
task: task_04
spec: 0047-context-driven-guidance-composition
status: completed
type: backend
complexity: high
---

# Task 04: Distribute rules to semantic owners

## Overview

Turn reviewed segment dispositions into literal repository-owned blocks in
active semantic guides or a non-empty residual carrier. The portable Plan
accounts for every source clause and removes obsolete carriers only when no
accepted residual remains.

## Requirements

1. MUST constrain `repository-document` destinations to the active semantic
   owner registry.
2. MUST insert stable repository-rule markers outside setup-owned blocks while
   preserving the accepted body bytes exactly.
3. MUST inventory later repository edits as repository-owned content rather
   than setup drift.
4. MUST migrate the three recognized repository-rule carriers and leave all
   other nested carriers untouched.
5. MUST remove empty residual carriers, legacy generic guides, and their root
   pointer in the same confirmed Plan.
6. MUST record every move, retention, rejection, and deletion through existing
   preimage, postimage, retention, and managed-entry ledgers.

## Subtasks

- [x] Admit active semantic-guide destinations.
- [x] Parse and render repository-owned rule blocks.
- [x] Migrate recognized repository-rule carriers.
- [x] Remove empty residual files and root pointers.
- [x] Add ownership, retention, edit, and rollback tests.

## Acceptance Criteria

- [x] Every accepted segment has one managed, semantic, residual, or reasoned
  rejection disposition.
- [x] Semantic block bodies match the source segment bytes.
- [x] Empty residual and legacy generic carrier paths are absent after apply.
- [x] Arbitrary nested carriers remain byte-identical and warning-only.
- [x] Empty reapply produces zero managed-file delta.

## Context

- instruction: `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- instruction: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- instruction: `docs/adr/0074-repository-rules-use-hybrid-semantic-ownership.md`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/preservation.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestSemanticRuleDistribution|TestRepositoryRuleBlock|TestResidualCarrier|TestNestedCarrier'` — expected: exact movement, ownership, cleanup, warning, ledger, and idempotency cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 2–3; User Stories 1–2 and 6; Core Features 5–8 and 16.
- `_techspec.md` → Implementation Design: Data Models and API Contracts; Build Order 3.
- ADR-0058 → complete Upgrade Retention accounting.
- ADR-0074 → semantic blocks before residual rules.

## Result

Planning now advertises only active semantic-guide paths and managed IDs,
validates every reviewed disposition against that registry, and accepts a
locally materialized segmented Source Baseline only when it reconstructs the
current carrier bytes exactly. Accepted semantic rules render in stable
repository-owned markers outside setup-owned blocks; later body edits remain
repository-owned content and keep their current bytes.

The preservation inventory now includes
`docs/agents/specific-repository.md`, `docs/agents/repository.md`, and
`docs/agents/repository-rules.md`. A confirmed redistribution deletes all
recognized source carriers, recreates only a non-empty residual carrier, and
removes the residual root pointer when no residual remains. Arbitrary nested
carriers stay outside the mutation set.

Acceptance evidence:

- `TestSemanticRuleDistributionAdmitsOnlyActiveSemanticOwners` proves the
  sealed classification contract advertises the active CLI semantic guide,
  omits the inactive frontend guide, admits the active destination, and
  rejects the inactive path.
- `TestSemanticRuleDistributionMovesExactBytesAndAccountsLedgers` proves
  segmented source bytes produce one exact semantic block and one exact
  residual, with source dispositions in the retention ledger and the
  repository-owned block in the managed-entry ledger but not the Setup
  Manifest.
- `TestResidualCarrierRemovesAllRecognizedEmptyResults` proves each of the
  three recognized carriers receives a missing postimage after full semantic
  distribution, apply removes the path, and no residual root pointer remains.
- `TestNestedCarrierRemainsByteIdenticalAndWarningOnly` proves an arbitrary
  nested `AGENTS.md` remains outside postimages, retains its warning, and is
  byte-identical after apply.
- `TestRepositoryRuleBlockPreservesRepositoryEditAndEmptyReapply` proves an
  edited semantic block body is retained and inventoried as repository-owned
  content, while applying the same confirmed Plan again produces no managed
  file delta.
- `TestRepositoryRuleBlockRollbackRestoresSemanticGuide` proves transaction
  rollback restores the complete visible preimage after an injected semantic
  guide verification failure.

Verification:

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run
  'TestSemanticRuleDistribution|TestRepositoryRuleBlock|TestResidualCarrier|TestNestedCarrier'`
  passed 11 tests in 2 packages.
- `rtk make verify` passed 2,153 tests in 22 packages, 4 setup-skill contract
  tests, the Roundfix skill check, and the final build.

Follow-ups: none.
