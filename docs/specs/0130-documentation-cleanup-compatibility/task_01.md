---
task: task_01
spec: 0130-documentation-cleanup-compatibility
status: pending
type: test
complexity: high
---

# Task 01: Preserve authorization and regeneration after cleanup

## Overview

Repair the approved five-file boundary described by the TechSpec. Demonstrate
that current grants and historical evidence work without a workflow directory,
while invalid authority and unowned writes remain rejected.

## Requirements

1. MUST implement all four Core Features through the TechSpec's existing seams.
2. MUST preserve deep equality and negative coverage; no suppressions or private
   exemption list. Use filesystem-only tests in suiteguardcontract.
3. MUST change only the following bounded paths and this Task's own evidence:
   - `internal/baseline/derived_ownership_test.go`
   - `internal/speccheck/mechanical_test.go`
   - `internal/speccheck/governed_repocontract_test.go`
   - `internal/suiteguardcontract/regeneration.go`
   - `internal/suiteguardcontract/regeneration_test.go`

## Subtasks

- [ ] Add positive and negative Spec-grant discovery and ownership tests.
- [ ] Implement bounded discovery and command-only ownership resolution.
- [ ] Adapt regeneration fixtures and assert parity with Baseline output ownership.
- [ ] Recover historical grant evidence with explicit unavailable-object outcomes.
- [ ] Record criterion-level evidence and limits in this Task's Result.

## Acceptance Criteria

- [ ] CF-1: approved current and legacy inputs work; proposed, null, malformed,
  unrelated, mismatched-consumer and symlinked inputs grant no authority.
- [ ] CF-2: command-only declarations match Baseline ownership and concrete
  positive/negative fixture expectations, including frozen/dedicated paths,
  sidecars, nested exceptions and invalid/conflicting declarations.
- [ ] CF-3: real historical grant bytes and non-empty bounded-path coverage remain
  tested; unavailable historical Task objects are explicit skips, with controlled
  accepted/refused Git-change replay providing non-vacuous coverage.
- [ ] CF-4/OE-1: regeneration fixtures exercise the missing-directory failure from
  PR #180 using Spec grants, preserving frozen boundaries and strict assertions.
- [ ] All implementation/test changes stay within the five approved paths.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- instruction: `docs/agents/go.md`
- instruction: `docs/agents/specific-repository.md`
- interface: `internal/baseline/derived_ownership.go`
- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguardcontract/contract.go`

## Verification

- `test -s internal/suiteguardcontract/regeneration_test.go && grep -q 'func TestCleanupRegenerationDiscovery' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestCleanupRegenerationDiscovery$'` — the new positive/negative reader contract executes; absence of its test cannot pass.
- `grep -q 'func TestCleanupRegenerationOwnershipParity' internal/baseline/derived_ownership_test.go && go test -count=1 -tags repocontract ./internal/baseline -run '^(TestCleanupRegenerationOwnershipParity|TestOutputsForCommand|TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries)$'` — actual ownership and cleanup-compatible regeneration fixtures agree.
- `grep -q 'func TestCleanupHistoricalGrantEvidence' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^(TestCleanupHistoricalGrantEvidence|TestAuditJudgesTheGrant|TestEveryBoundedPathIsGoverned)$'` — historical and controlled grant evidence is exercised, with unavailable original Task objects explicitly reported.

## References

- `_prd.md` → Goals 1–3, Core Features 1–4, Outside Evidence OE-1.
- `_techspec.md` → all Design sections, Coverage Map, Research and decision evidence.
- ADR-0149 — command authority and tree-owned outputs.
