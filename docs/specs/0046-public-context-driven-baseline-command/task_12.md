---
task: task_12
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 12: Recalculate rejected plans from scoped proposals

## Overview

Let maintainers reject a final plan, revisit one structured decision area, and
optionally translate free-form feedback into an in-scope proposal. Every
accepted revision returns through deterministic validation and a new approval.

## Requirements

1. MUST let the maintainer select profile, repository rules, divergences, or
   files as the decision area to revisit after rejecting a plan.
2. MUST accept direct structured changes without ACP and optional free-form
   suggestions through a separate strict revision snapshot and proposal
   schema.
3. MUST reject suggestions outside Baseline scope, unknown decisions,
   unauthorized destinations, incomplete changes, extra prose, or snapshot
   digest mismatch.
4. MUST recalculate the complete plan from the immutable repository snapshot
   and normalized accepted decisions rather than patching the prior artifact.
5. MUST show the resulting file projection and ledger and require a new exact
   Plan Digest approval before apply.
6. MUST preserve manual revision when ACP is unavailable.

## Subtasks

- [ ] Add rejected-plan decision-area transitions.
- [ ] Define strict revision snapshot and proposal contracts.
- [ ] Supervise scoped ACP revision with the sealed analyzer.
- [ ] Revalidate decisions and recompute the complete plan.
- [ ] Add repeated-revision, out-of-scope, fallback, and approval tests.

## Acceptance Criteria

- [ ] A direct structured revision requires no ACP session.
- [ ] An accepted in-scope suggestion changes only permitted Baseline decisions.
- [ ] An out-of-scope or invalid proposal changes no decision and exposes a manual next action.
- [ ] Every accepted revision produces a newly computed plan and approval digest.
- [ ] A previously approved digest cannot authorize a revised plan.
- [ ] Repeated rejection remains deterministic and does not restart repository adoption.

## Context

- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- interface: `internal/baselineacp`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli -run 'TestRejectedPlanRevision|TestScopedRevisionProposal|TestRevisionOutOfScope|TestRevisionManualFallback|TestRevisionRequiresNewApproval'` — expected: structured and semantic revisions, scope rejection, fallback, recomputation, and approval invalidation pass.
- `rtk go test -count=1 ./internal/cli -run TestRepeatedPlanRevisionDeterminism` — expected: repeated review cycles preserve state and produce deterministic plans.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Story 6; Core Feature 11; User Experience.
- `_techspec.md` → Interfaces: ProposalAnalyzer; API Contracts: ACP revision; Build Order 6–7.
- ADR-0069 → proposal-only scoped semantic analysis.
