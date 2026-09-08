---
spec: 0129-spec-authoring-and-gate-recovery
prd: _prd.md
created: 2026-09-08
---

# Explicit Spec promises and supported gate recovery — Technical candidate

## Executive Summary

Extend the authoring evidence graph beyond citation presence and make premise changes and stale QA recoverable through supported owners. Preserve completed Results, terminal QA and existing correction ceilings.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

## Project Constraints

- Identifier strategy: applicable — preserve Spec/Task identities and exact Git ref identity; any new recovery record follows the existing Run identity boundary. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential or HTTP API change is proposed; use existing read boundaries only. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014 and ADR-0057 reserve Verification and Task status to the Daemon; ADR-0091 keeps QA terminal; ADR-0117 places each detector at the stage that can establish it; ADR-0148 preserves non-vacuous Verification probing. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: preserve the distinction between environment-blocked QA rows and product failure or success.
  ADR-0020 applies at the inherited Agent result boundary: a parsed prompt result remains usable despite a later teardown exit, with the anomaly recorded before Verification.
  ADR-0038 applies: the same Agent Session has only the accepted single Verification Feedback repair; recovery cannot enlarge it.
  ADR-0096 applies: prove mechanical QA facts before spending an Agent turn and preserve the existing verdict semantics.
  ADR-0127 applies: process residue remains a readiness observation, not a synthetic Run or authority to settle work.
  ADR-0056 applies: authoring/recovery preserves distinct Task and Verification capacities; neither grants whole-queue concurrency.
  ADR-0093 applies: consistency uses explicit citations and does not infer behavior from their presence.
  ADR-0097 applies: carry forward QA only from declared, unchanged evidence; stale input requires supported revalidation.
  ADR-0104 applies: acceptance uses independent evidence with its actual origin; missing external evidence remains visible under the declared policy.
  ADR-0154 applies to archive disposition: an explicitly user-authorized QA Archive Override preserves actual QA/Task evidence and does not satisfy independent delivery gates.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `.agents/skills/implement-task/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-techspec/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-tasks/references/task-template.md`, `skills/implement-task/SKILL.md`, `internal/speccheck/coherence.go`, `internal/baseline/assets/modules/spec-workflow.json`, `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Requirement evidence map | `owned write-prd/write-techspec/write-tasks skills and templates` | Connect API contracts, metrics and accepted ADR obligations to Tasks and behavioral evidence. |
| Mechanical authoring checks | `internal/speccheck/coherence.go and stage-specific detectors` | Check reference completeness and declared graph structure without claiming semantic proof. |
| Optional semantic independent review | `Spec 0126 reviewer contract and authored review checklist` | When review is enabled, compare implementation behavior and rejected alternatives to operative decisions; explicit none records omission without a semantic-review claim. |
| Amendment and QA recovery | `existing Spec/Task loaders and Daemon lifecycle/store` | Preserve prior Results/reports and authorize supported supersession/invalidation. |
| Temporal prerequisites | `Task Graph authoring and queue preflight` | Represent future observations, ordinal-generator edges and changed prerequisite assumptions. |
| Verification contract | `canonical spec-routing, Task authoring and existing rehearsal/CI declarations` | Require property-shaped acceptance and declared feasible test surfaces without broadening corrective budgets. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Promises and evidence

Extend authored traceability to API Contracts, Success Metrics and each
applicable accepted-ADR obligation. Name the consuming Task and the observable
behavior/evidence that will settle it. Citation presence remains a mechanical
fact, not proof that a decision was obeyed. When enabled by the pre-PR policy, independent semantic review checks
the chosen design and its rejected alternatives against actual behavior;
unsupported conclusions remain findings. A documentation assertion or mock
that reproduces the implementation is insufficient product evidence.

Property-shaped acceptance states what remains true over meaningful inputs
and failure boundaries. Use the narrowest seam that exercises the actual
contract, including public use-case/persistence behavior when needed. Update
the affected existing contract instead of adding parallel tests that cannot
fail for its regression. Declare newly required test classes and check that
the actual CI command executes them. Keep Spec commits narrow and preserve
current rehearsal requirements; this proposal does not lower CI feasibility.

### Amendment and stale gate

A falsified premise produces an explicit amendment/supersession record with
its reason, changed assumptions, affected consumers and approval impact.
Retain every previous Result and report; do not erase failed or obsolete
evidence. The next authoring pass revalidates dependents after prerequisite
changes. A broader protected scope, future release or irreversible operation
blocks until its own authority exists.

Provide a public supported stale-QA recovery operation owned by the Daemon.
It identifies the exact Spec revision/report and reason, verifies no competing
active owner, records invalidation history and schedules eligible new QA work
under the accepted lifecycle. It cannot rewrite a completed status from the
Supervisor, delete the old report, reopen an immutable archive or fabricate
pass evidence. The detailed transition must be characterized against current
Task settlement before its command/schema is finalized; absence of a safe
transition remains an explicit implementation blocker, not a manual-edit
fallback. Archived cases use a new corrective Spec with sufficient authority.

### Dependencies and ceilings

Temporal prerequisites are explicit observations or operations, not ordinary
Task completion invented in advance. Record the required external fact and
its evidence; a future release requirement does not authorize creating that
release. Generated ordinal identifiers require explicit ordering edges before
parallel authoring/execution. Whole-queue consumers revalidate these facts at
start and after earlier merges.

Preserve Task Type, the authored terminal QA Task, ADR-governed corrective
ceilings and the one permitted Verification Feedback repair. When the ceiling
is reached, stop and reauthor/split through the supported contract. Neither a
new queue record nor a fresh agent session resets the existing Task's budget.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Trace each API/metric/ADR obligation to its consuming Task and behavioral evidence. An amendment retains prior Results/report identities, reason, affected consumers and approval impact; temporal evidence is separate from Task completion.

### API Contracts

Stage-specific Spec Check reports missing traceability separately from semantic review. A supported stale-QA operation must identify its exact revision/report, preserve history and use the Daemon transition; finalize that transition before authoring executable Tasks.

## Coverage Map

- PRD Goal 1 → Requirement evidence map, Mechanical authoring checks, Optional semantic independent review.
- PRD Goal 2 → Amendment and QA recovery.
- PRD Goal 3 → Temporal prerequisites, Verification contract.
- Core Feature 1 → Requirement evidence map and Mechanical authoring checks and Optional semantic independent review.
- Core Feature 2 → Amendment and QA recovery.
- Core Feature 3 → Amendment and QA recovery.
- Core Feature 4 → Temporal prerequisites.
- Core Feature 5 → Verification contract.
- Core Feature 6 → Verification contract.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
The [owned source index](references/_index.md) records each primary source.
The prerequisite Specs are listed below. Secondary consumers reference the
primary owner's adopted source instead of duplicating it. External research was read through
Exa and compared with local Secondbrain history. The
[historical research record](https://github.com/marcioaltoe/roundfix/blob/6b8ea48725cbca13974eee0b400b3482202874f6/docs/workflow/2026-09-08-pending-work-plan.md)
retains the consulted sources, their influence and limitations after the plan
was removed from the current tree. Published interfaces support feasibility,
not a claim that the proposed runtime or behavior already exists.

## Testing Approach

Use focused tests at the named package seams for deterministic rules, real
Git/store/process boundaries for integration behavior, and the authored public
QA Task for user-visible acceptance. Do not infer a terminal pass from source
inspection or a focused fixture. Required observations:

1. A missing API/metric/ADR-to-Task/evidence reference is found at its authoring stage; complete citations with contradictory behavior fail enabled semantic review. Explicit none records no such review while preserving mechanical checks and QA.
2. A premise amendment preserves prior Results and forces affected consumers to revalidate; broader scope or missing temporal evidence blocks.
3. Stale QA recovery keeps the earlier report and audit trail, respects active ownership and cannot reopen archived/completed history by manual editing.
4. Ordinal-generator dependencies prevent concurrent collision and a future-release prerequisite does not execute publication.
5. Required test classes absent from the actual CI command fail feasibility; property-shaped cases expose a regression rather than mirror implementation.
6. At the accepted corrective ceiling the workflow reauthors/splits instead of manufacturing another same-Task repair.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Requirement/evidence and temporal-prerequisite authoring contract (depends on: none).
2. Stage-specific mechanical checks and semantic review integration (depends on: 1).
3. Supported amendment and Daemon-owned stale-QA transition (depends on: 1).
4. Canonical/owned guidance, CI feasibility and preserved ceilings (depends on: 1, 2, 3).
5. Public authoring/recovery journeys and terminal QA (depends on: 1, 2, 3, 4).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Do not pretend mechanical traceability proves ADR obedience. The safe stale-QA transition and final public operation require characterization before Task authoring; shared authoring files overlap 0119/0121/0122 and must be serialized and revalidated.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Preserve configured reviewer selection; this Spec introduces no separate reviewer override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md), [0121](../0121-baseline-decisions-and-complete-regeneration/_techspec.md), [0122](../0122-verified-content-and-terminal-settlement/_techspec.md), [0126](../0126-agent-review-before-pull-request/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.


## Authorized archive disposition

Consume the archive policy from Spec 0122 and ADR-0154. An applicable explicit
user authorization can archive the covered Spec with unmet QA, recording the
override and preserving original evidence. This does not reopen or complete
Tasks, declare QA passed, imply review approval or authorize publication/merge.
A durable workflow records the overridden archive and evaluates subsequent
actions against their own approval and gates; it does not retry the waived
archive prerequisite or ask again for the same applicable archive approval.
The absence of authority for a later action remains a separate visible blocker.
