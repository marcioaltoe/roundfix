---
task: task_10
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 10: Classify root instructions through sealed ACP proposals

## Overview

Add optional semantic classification without granting the Agent access to the
checkout or mutation authority. Preferred and fallback selections analyze the
same sealed snapshot, while invalid or unavailable analysis returns the
deterministic workflow to manual classification.

## Requirements

1. MUST eagerly prove Codex `gpt-5.6-sol`/`xhigh` and Codex
   `gpt-5.5`/`xhigh` through Exact Agent Selection Proof without reading
   configurable Agent Work Categories.
2. MUST run each attempt in a fresh session and private empty directory with
   terminal and tools denied, one turn, bounded input, output, and deadline,
   and no checkout or Agent-log path.
3. MUST send byte-identical canonical snapshot input to preferred and fallback
   attempts and reject tools, extra prose, unknown or missing IDs, duplicate
   dispositions, digest mismatch, unsupported destinations, oversized output,
   timeout, or invalid JSON.
4. MUST discard raw and invalid ACP output, persist no Run or Run Event, and
   admit only a complete normalized proposal to consolidated human review.
5. MUST use manual classification when both selections fail and must not start
   fallback after parent cancellation or unproven session cleanup.

## Subtasks

- [ ] Define the strict classification snapshot and proposal schemas.
- [ ] Build the non-Run sealed ACP execution adapter.
- [ ] Implement exact preferred/fallback supervision and strict cleanup.
- [ ] Integrate validated proposals with root-instruction planning.
- [ ] Add adapter, fallback, cancellation, privacy, and no-mutation tests.

## Acceptance Criteria

- [ ] A valid preferred proposal prevents a fallback attempt.
- [ ] Invalid preferred output is discarded and fallback receives identical snapshot bytes.
- [ ] Tool use or checkout access is impossible under the emitted ACPX arguments.
- [ ] Both unavailable or invalid selections return complete manual-classification destinations.
- [ ] Parent cancellation closes the active session and never activates fallback.
- [ ] No outcome creates repository changes, Run rows, Run Events, or Agent logs.
- [ ] Equivalent accepted classifications produce the same Plan Digest regardless of proposal source.

## Context

- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/acp_stream.go`

## Verification

- `rtk go test -count=1 ./internal/baselineacp ./internal/baseline -run 'TestSealedClassification|TestPreferredFallback|TestProposalValidation|TestManualClassificationFallback|TestAnalyzerCancellation|TestAnalyzerNoPersistence'` — expected: proof, sandbox, validation, fallback, cleanup, privacy, and deterministic-plan cases pass.
- `rtk go test -count=1 ./internal/baselineacp -run TestACPXReadOnlyArguments` — expected: arguments deny tools and terminal, bound the attempt, omit full access, and expose no checkout path.
- `rtk make verify` — expected: the full repository gate passes without requiring a live model.

## References

- `_prd.md` → User Story 5; Core Features 8–9, 17, 19; Non-Goals / Out of Scope.
- `_techspec.md` → Interfaces: ProposalAnalyzer; API Contracts: ACP classification; Integration Points: ACPX/Codex ACP; Build Order 6.
- ADR-0039 → disposable exact-selection proof.
- ADR-0069 → supervised read-only semantic analysis.
