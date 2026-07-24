---
task: task_04
spec: 0047-context-driven-guidance-composition
status: pending
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

- [ ] Admit active semantic-guide destinations.
- [ ] Parse and render repository-owned rule blocks.
- [ ] Migrate recognized repository-rule carriers.
- [ ] Remove empty residual files and root pointers.
- [ ] Add ownership, retention, edit, and rollback tests.

## Acceptance Criteria

- [ ] Every accepted segment has one managed, semantic, residual, or reasoned
  rejection disposition.
- [ ] Semantic block bodies match the source segment bytes.
- [ ] Empty residual and legacy generic carrier paths are absent after apply.
- [ ] Arbitrary nested carriers remain byte-identical and warning-only.
- [ ] Empty reapply produces zero managed-file delta.

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
