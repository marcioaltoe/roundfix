---
task: task_17
spec: 0046-public-context-driven-baseline-command
status: pending
type: test
complexity: high
---

# Task 17: Prove release-ready Baseline journeys

## Overview

Run the final product-level evidence gate over the public CLI, all maintained
profiles, every finding correction, and external composition. This Task is
AFK for hermetic journeys but HITL for the live Fluxus adoption or update,
which requires separate explicit authorization before touching that repository.

## Requirements

1. MUST create the dated Spec-local QA report and map every PRD user story,
   Core Feature, success metric, and seven Fluxus finding categories to fresh
   evidence.
2. MUST run real binary journeys for greenfield, preservation, update, profile
   change, manual ACP fallback, preferred/fallback proposal, rejected-plan
   revision, automation plan/apply, stale plan, cross-clone apply, unsafe
   carrier, rollback, recovery, and empty reapply.
3. MUST prove every maintained profile composes with its formatter,
   repository Verification, final Baseline audit, and empty reapply without a
   managed-file delta; Baseline itself must execute none of those repository
   commands.
4. MUST prove the setup skill contains no executable engine and every
   documented example matches the public parser.
5. MUST perform one separately authorized live Fluxus adoption or update and
   record repository identity, approved digest, observable results, recovery
   state, and absence of the seven original operator-friction outcomes.
6. MUST leave the Task pending when live authorization or any required
   criterion is unavailable or failing.

## Subtasks

- [ ] Create the QA matrix with every criterion pending.
- [ ] Execute hermetic real-binary profile and failure journeys.
- [ ] Execute formatter and Verification composition outside Baseline.
- [ ] Validate docs, skill cutover, schemas, and finding regressions.
- [ ] Request and execute the separately authorized Fluxus journey.
- [ ] Record exact evidence and settle every QA row.

## Acceptance Criteria

- [ ] All 10 User Stories, 22 Core Features, and success metrics have fresh passing evidence.
- [ ] All seven Fluxus finding categories have named regression evidence.
- [ ] Every maintained profile completes plan, apply, verification, formatter composition, final audit, and empty reapply.
- [ ] Failure journeys prove zero unauthorized or unrecoverable mutation.
- [ ] Interactive and non-interactive equivalent inputs produce identical Plan Digests.
- [ ] The live Fluxus journey completes only after separate authorization and needs no manual schema repair, route digest calculation, capability reinterpretation, or command correction.
- [ ] The full repository Verification passes after all QA evidence is recorded.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `docs/findings/2026-07-23-setup-context-driven-adoption-process-improvements.md`
- interface: `docs/specs/0046-public-context-driven-baseline-command/qa`
- interface: `internal/cli/implement_test.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run 'TestBaselineMacroJourneys|TestBaselineFindingRegressions|TestBaselineFormatterComposition|TestBaselineDocumentationContract'` — expected: real-binary journeys, all seven finding regressions, composition, and docs contracts pass hermetically.
- `rtk grep -R -n 'status: passed' docs/specs/0046-public-context-driven-baseline-command/qa` — expected: the dated final QA report records a passing result after the separately authorized Fluxus journey.
- `rtk make verify` — expected: the full release-ready repository gate passes.

## References

- `_prd.md` → Goals 1–6; User Stories 1–10; Core Features 1–22; Success Metrics; Non-Goals / Out of Scope.
- `_techspec.md` → Coverage Map; Testing Approach; Build Order 11; Risks & Considerations.
- ADR-0066 through ADR-0073 → complete authority, safety, portability, analysis, parity, and transaction contracts.
