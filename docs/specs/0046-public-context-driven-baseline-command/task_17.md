---
task: task_17
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Create the QA matrix with every criterion pending.
- [x] Execute hermetic real-binary profile and failure journeys.
- [x] Execute formatter and Verification composition outside Baseline.
- [x] Validate docs, skill cutover, schemas, and finding regressions.
- [ ] Request and execute the separately authorized Fluxus journey.
- [ ] Record exact evidence and settle every QA row.

## Acceptance Criteria

- [ ] All 10 User Stories, 22 Core Features, and success metrics have fresh passing evidence.
- [x] All seven Fluxus finding categories have named regression evidence.
- [x] Every maintained profile completes plan, apply, verification, formatter composition, final audit, and empty reapply.
- [x] Failure journeys prove zero unauthorized or unrecoverable mutation.
- [x] Interactive and non-interactive equivalent inputs produce identical Plan Digests.
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

## Result

Partial result — pending the separately authorized live Fluxus journey.

- Created
  `qa/qa-report-2026-07-23-release-ready.md` with every PRD User Story, Core
  Feature, Success Metric, Fluxus finding category, named journey, maintained
  profile, applicable Non-Goal, and Task acceptance criterion mapped to fresh
  evidence. The report currently records 87 passing rows and seven pending
  live/final rows.
- Added the named release-gate suites required by this Task. They exercise the
  real public binary across greenfield, preservation, update, profile change,
  proposal fallback, revision, automation, stale-plan, cross-clone, unsafe
  carrier, rollback, recovery, and empty-reapply journeys.
- Proved all three maintained Baseline Profiles compose with their formatter
  and repository Verification outside Baseline, then pass final audit and
  empty reapply without managed-file drift. Execution traps prove Baseline
  runs none of those external commands.
- Proved all seven Fluxus finding categories, public documentation examples,
  and the thin setup-skill cutover contract. The setup skill contains no
  executable engine.
- Found and fixed one blocking product defect: repository-owned reference
  paths such as `repository.design-contract` were not available to artifact
  rendering, so the TypeScript profile could not produce a Change Plan.
  Safe repository-owned `path` references now render at the same boundary as
  setup-owned `managedId` references.

Evidence:

- `rtk env GOCACHE=/private/tmp/roundfix-task17-go-cache go test -count=1
  ./internal/baseline ./internal/baselineacp ./internal/cli ./skills -run
  'TestBaselineMacroJourneys|TestBaselineFindingRegressions|TestBaselineFormatterComposition|TestBaselineDocumentationContract'`
  passed in all four packages.
- `rtk make verify` passed before the live journey: 2,077 tests in 22
  packages, four setup-skill cutover guards, Repository Skill Set checks, and
  the binary build. Per the acceptance criterion, it must run again after all
  live QA evidence is recorded.
- The detailed row-level evidence, finding reproduction, recovery assertions,
  and exact remaining rows are in the dated QA report.

Acceptance-criterion status:

- AC 2–5 pass with fresh hermetic evidence.
- AC 1 remains pending because Success Metrics 13 and 16 require the live
  journey and final Verification.
- AC 6 remains pending because this Task has not yet received separate
  authorization to touch the Fluxus repository.
- AC 7 remains pending until the post-live full repository Verification
  passes.

Requirement 6 requires this Task to remain `pending` at this authorization
boundary. No Fluxus files were read or changed, and no commit, push, or pull
request was performed.
