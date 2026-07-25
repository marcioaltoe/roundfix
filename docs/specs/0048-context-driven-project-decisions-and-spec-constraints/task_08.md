---
task: task_08
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
type: test
complexity: high
---

# Task 08: Prove project-constraint journeys

## Overview

Create hermetic macro journeys for typed decisions, rendered guidance, Spec
authoring, tooling refusal, apply, audit, and empty reapply. Separately
authorized Fluxus journeys remain fresh final QA evidence.

## Requirements

1. MUST cover human and automation decision collection for every affected
   maintained Profile.
2. MUST cover compatible reuse, unresolved input, invalid objects, derived
   exception conflict, and stable Plan Digests.
3. MUST author one new PRD and TechSpec fixture with complete Project
   Constraints and one authorized tooling fixture with exact bounded files.
4. MUST prove decomposition, execution, and QA refusal for missing or exceeded
   tooling authorization.
5. MUST complete formatter, apply, repository Verification recommendation,
   audit, and empty reapply with zero managed delta.
6. MUST define fresh Fluxus greenfield and update evidence for final
   `qa-gate`.
7. MUST implement the named journey tests in this Task's Verification command;
   a Spec-root-only status commit or a filter that executes no matching tests
   is not completion evidence.
8. MUST remove the Task 06 test-only Oxfmt version pin and provision the
   hermetic formatter from the version already owned by the maintained Profile
   or disposable repository fixture; this Spec authorizes no new protected
   tooling version pin.

## Subtasks

- [x] Add all-profile decision journey fixtures.
- [x] Add decision reuse and conflict journeys.
- [x] Add PRD and TechSpec authoring fixtures.
- [x] Add tooling authorization and refusal journeys.
- [x] Add formatter, audit, reapply, and final QA assertions.
- [x] Remove the duplicate test-only Oxfmt pin while preserving hermetic
  formatter execution.

## Acceptance Criteria

- [x] Equivalent human and automation answers produce identical Plans.
- [x] Missing or conflicting structured values stop without mutation.
- [x] New Spec fixtures contain every required Project Constraint row and
  source.
- [x] Tooling mutation outside bounded authorization fails before settlement.
- [x] Every affected Profile completes apply, audit, and empty reapply with
  zero managed delta.
- [x] The QA matrix requires fresh Fluxus greenfield and update evidence.
- [x] Each named journey test exists, executes assertions, and fails when its
  corresponding contract is removed.
- [x] The final Spec delta adds no protected tooling version pin, and the
  release journey still uses the exact maintained Profile formatter version.

## Context

- instruction: `docs/adr/0076-typed-project-decisions-render-identifier-and-authentication-policy.md`
- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/release_gate_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli ./skills -run 'TestProjectDecisionJourney|TestProjectConstraintJourney|TestToolingAuthorizationJourney|TestBaselineReleaseGate'` — expected: decision, authoring, refusal, formatter, apply, audit, and reapply journeys pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; Core Features 1–17; Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 7; Risks & Considerations.
- ADR-0076 and ADR-0077 → final project-decision and Spec-constraint contracts.

## Reopened QA repair

The 2026-07-25 QA gate found that the prior Daemon-owned commit changed only
this Task's status. It added no journey tests, left every checkbox open, and
recorded no Result. The rerun must implement the real journeys, including the
Fluxus-style keep-defaults update repaired by Task 02, before this Task may
return to `completed`. It must also remove the duplicate Task 06 Oxfmt pin so
the final implementation matches the Spec's explicit no-authorization tooling
constraint.

## Result

Implemented the named engine, public CLI, Project Constraint, and tooling
authorization journeys. The public journey proves human/automation Plan
parity, compatible update reuse including the persisted Better Auth reason,
missing-input no-mutation behavior, apply, formatter and repository
Verification composition, a zero-change fresh audit Plan, and empty reapply.
The existing all-Profile release journey uses the same lifecycle helper.

Added disposable PRD and TechSpec fixtures with all four Project Constraint
rows, reasons, operative `docs/agents/` sources, and matching express
authorization limited to `.golangci.yml` and `scripts/verify.sh`. The
tooling-authorization journey accepts those exact paths, refuses missing
authorization, refuses `scripts/release.sh`, and mutation-checks the
decomposition, execution, and QA skill clauses. Its QA matrix leaves separate
fresh Fluxus greenfield and update rows pending for the final `qa-gate`; prior
Task or QA evidence cannot satisfy them.

Removed the duplicate test-only Oxfmt version constant. The disposable
TypeScript repository now reads the formatter version from the maintained
Profile, and the hermetic installer reads that repository-owned version,
isolates Bun's temp, cache, and home directories, installs it, and verifies the
installed manifest before running the formatter.

Acceptance evidence:

- Equivalent human and automation answers:
  `TestProjectDecisionJourney/human and automation answers produce one Plan`
  compares byte-identical Plans and Plan Digests.
- Missing or conflicting values: the engine journey covers strict invalid
  objects and derived HTTP conflicts; the public CLI journey snapshots the
  repository before missing automation input and proves no mutation.
- Complete Spec fixtures: `TestProjectConstraintJourney` validates every row
  and source in both fixtures, then removes each row and requires rejection.
- Bounded tooling authority: `TestToolingAuthorizationJourney` accepts only
  the two exact fixture paths and rejects missing or exceeded authorization
  before settlement.
- Apply, audit, and reapply: the named TypeScript journey and the existing
  all-Profile release journey completed formatter, repository Verification,
  zero-change fresh planning, and verified empty reapply without managed-byte
  drift.
- Fresh final QA: the fixture matrix requires new disposable Fluxus
  greenfield and update clones, command transcripts, clone identities, and
  pending final `qa-gate` evidence.
- Named execution: `go test -list` reports
  `TestProjectDecisionJourneyEngine`,
  `TestToolingAuthorizationJourneyCoreClause`, `TestBaselineReleaseGate`,
  `TestProjectDecisionJourney`, `TestProjectConstraintJourney`, and
  `TestToolingAuthorizationJourney`.
- Formatter provenance: the release-gate test contains neither the removed
  test-only constant nor a formatter-version literal; the passing journey
  obtains and verifies the maintained Profile's formatter version.

Verification:

- `rtk env TMPDIR=/private/tmp
  GOCACHE=/private/tmp/roundfix-task08-go-cache go test -count=1
  ./internal/baseline ./internal/cli ./skills -run
  'TestProjectDecisionJourney|TestProjectConstraintJourney|TestToolingAuthorizationJourney|TestBaselineReleaseGate'`
  passed in all three packages. The environment only redirects caches around
  the child-Agent sandbox; the Daemon runs the Task command verbatim.
- `rtk env TMPDIR=/private/tmp
  GOCACHE=/private/tmp/roundfix-task08-go-cache make verify` passed: 2,344 Go
  tests, the four focused skill-runtime tests, `roundfix skills check`, and the
  final build.
- `rtk git diff --check` passed.

No follow-up implementation belongs to this Task. The live Fluxus executions
remain reserved for the separately authorized final `qa-gate`.
