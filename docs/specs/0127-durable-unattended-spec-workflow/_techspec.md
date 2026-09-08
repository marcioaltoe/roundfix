---
spec: 0127-durable-unattended-spec-workflow
prd: _prd.md
created: 2026-09-08
---

# Durable delivery of an approved Spec queue — Technical candidate

## Executive Summary

Extend the existing Daemon and SQLite store with durable queue ownership between Spec Runs. Preserve the one-Spec Implement executor and add receipt-driven publication recovery, explicit limits and observable blockers. This is a proposed architecture; the current detached Implement command does not provide it.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: applicable — preserve Run, Task, Spec, and Git identities and make durable delivery refer to their actual records. New work branches use purpose prefixes; tool-owned Run/Task branches use Roundfix's existing namespace. The maintainer removed the personal-prefix requirement on 2026-09-08, so naming no longer blocks Run creation. No new durable identity format is chosen by this PRD. Source: `docs/agents/domain.md`, `docs/agents/agent-instructions.md`.
- Authentication and HTTP: applicable — delivery uses the repository's existing authenticated runtimes and GitHub boundary, with no new credential storage or authentication policy. No backend HTTP guide exists; absence is not permission to create one. Publication authority and resource limits must be recorded, and a missing or denied credential parks the affected action. Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — preserve execution ownership and evidence semantics while proposing durable delivery across Runs. Source: `docs/agents/domain.md`, `docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon remains responsible for Task Verification and settlement as the durable Supervisor advances delivery.
  ADR-0019 applies to the inherited Watch meaning: Clean denotes a merge-ready Open Pull Request rather than an empty local queue; the replacement readiness contract is a prerequisite.
  ADR-0020 applies to retained acpx Batch execution: a parsed prompt response remains usable despite a subsequent teardown exit, and that anomaly remains journaled before authoritative Verification.
  ADR-0038 applies: the same Agent Session gets one Verification Feedback repair; queue-level retry proposals cannot enlarge that bound.
  ADR-0056 applies: per-Run Task Capacity and Verification Capacity remain separate and cannot be interpreted as whole-queue concurrency or a machine-wide semaphore.
  ADR-0057 applies: the Daemon is the sole writer of Implement Task status, so recovery and Supervisor narration cannot manufacture a completed Task.
  ADR-0091 applies: each authored QA gate remains a terminal Task in its Spec graph.
  ADR-0096 applies: mechanical facts are checked within the Daemon-owned QA stage before spending an Agent turn, with existing verdict semantics preserved.
  ADR-0117 applies: the workflow checks authoring defects when producing their artifacts and retains commit-dependent and user-surface checks at the gate that can establish them.
  ADR-0127 applies: process residue is a readiness observation, not a synthetic Run or a reason for the inventory to settle work after restart.
  ADR-0139 applies: durable recovery must retain one Active Run per work target and the same-checkout mutation guard rather than starting a competing owner.
  ADR-0142 applies to the existing head-bound evidence contract and historical outcomes; independent pre-PR review must be settled by Spec 0126 before the new delivery path relies on it.
  Preserve the accepted decision to keep `implement-spec` and repair its conflicting instructions rather than introducing a competing implementation loop.
  ADR-0151 applies: the queue consumes the configured review profile without overriding it with a forced Codex invocation.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/implement-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `internal/baseline/assets/modules/autonomous-work.json`, `skills/implement-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `skills/baseline_skill_contract_test.go`, `internal/docscontract/publicdocs_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Delivery queue owner | `internal/daemon and internal/store` | Persist queue, ordered items, owner identity and action intent/receipt separately from actual Runs. |
| Delivery CLI | `internal/cli/delivery.go (new) and CLI registration` | Prepare/start/status/resume/stop through one public command family. |
| Eligibility and limits | `internal/config, daemon preflight and Run lifecycle` | Validate approved Spec revision, dependencies, authority, Run Window, enforced Run deadline and distinct queue limits. |
| Action reconciliation | `existing Git/GitHub boundaries and store transactions` | Record intent before an external action and reconcile its observed result before retry. |
| Review and settlement integration | `Spec 0122 settlement and Spec 0126 candidate review` | Require complete current-candidate evidence before publication and merge. |
| Canonical delegation | `owned implement-spec/roundfix skills and autonomous guide` | Keep implementation and Verification in Roundfix; use one queue entry point and current authoring contract. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Ownership and records

Use the existing Daemon package and SQLite store. A Delivery Queue is a new
orchestration record, not a fabricated Run. A queue contains an ordered set of
Spec revisions, dependencies, approved limits, owner/process identity, state
and durable blockers. Each item links its actual Run IDs, candidate commit,
review receipt, PR identity and observed merge commit. An action row contains
a stable action key, expected repository/remote/candidate, intent, attempt
history and receipt. Use the existing store transaction conventions; no new
service, cloud dependency or competing process inventory is required.

Proposed public commands are `roundfix delivery prepare`, `start`, `status`,
`resume` and `stop`, with repository conventions for JSON, refusal and detached
ownership. Final flags/output schema are authored with their public contract
before implementation. Preparation is read-only and cannot approve a Spec.
Start persists accepted authority and validates current sources. Status reads
the durable owner and receipts; it never calls an absent owner running.

### Transitions and recovery

Serialize whole-Spec delivery initially. For each eligible item: run the
existing Implement Task Graph and terminal QA; archive through its existing
contract; commit the final candidate; obtain Spec 0126 independent review;
push/open the PR; await required checks for that candidate; squash merge;
reconcile the observed merge; then revalidate the next Spec. Each step records
its intent before mutation and receipt after independent observation.

A crash between external success and local acknowledgement causes discovery
and reconciliation, not blind replay. PR lookup includes repository, exact
head branch and base; a matching title alone is insufficient identity. Merge
recovery validates the PR's observed merged state and commit. A changed head,
foreign PR, non-fast-forward branch or unverifiable external result blocks the
action. Do not rewrite remote history or use an administrator bypass. Preserve
all historical attempts and review receipts.

Retain the current archive boundary as a proposal requiring approval. A late
blocking review parks the candidate. Any corrective Spec needs its own
sufficient authority and gate; the queue cannot reopen archived Tasks or
manufacture a grant from remaining cycle capacity. It must review the final
corrected candidate again. A publication failure can resume only after the
candidate, authority and external state are revalidated.

### Limits and cancellation

Run Window controls new starts. A Run already admitted follows its own budget;
it is not killed merely because the start window closes. Apply the configured
maximum Run duration to the actual Implement execution context and propagate
cancellation to owned Agent/process boundaries. Current production Implement
only warns/renders this value; the static finding requires a real bounded
cancellation/recovery test before this queue is called unattended-ready.

Queue duration, spending, whole-Spec concurrency and correction cycles are
separate approved limits. Proposed initial policy is one Spec at a time, no
new paid API calls without a recorded ceiling, and at most two review
correction cycles; these are candidates, not supplied answers. Stop prevents
new actions and requests cancellation of owned work while preserving truthful
recoverable state. An exhausted limit parks the queue with its measured usage
and remaining obligations. Missing price/usage data cannot mean zero cost.

Only one question may be pending. A resolved decision becomes evidence attached
to the consuming Spec revision; silence, timeout and a recommended option do
not resolve it. Recheck grants and Spec assumptions after prerequisite merges,
using Spec 0129's authoring/recovery contract. Invalidated premises, broader
scope, a release requirement or a new irreversible action require explicit
resolution before the dependent step.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Queue, queue item and action records are distinct from Runs. Persist approved source revisions/limits, actual Run links, candidate/review/PR/merge identities, action intent/receipt and blockers; recovery never invents an actual Run.

### API Contracts

The proposed delivery command family prepares/starts/observes/resumes/stops an approved queue. Read-only preparation cannot start work. Mutating actions require sufficient recorded authority, current evidence and finite limits; schemas/flags must be finalized before Task dispatch.

## Coverage Map

- PRD Goal 1 → Delivery queue owner, Review and settlement integration.
- PRD Goal 2 → Delivery queue owner, Action reconciliation.
- PRD Goal 3 → Eligibility and limits.
- PRD Goal 4 → Delivery queue owner, Delivery CLI.
- PRD Goal 5 → Canonical delegation.
- User Story 1 → Delivery CLI, Eligibility and limits.
- User Story 2 → Delivery queue owner, Action reconciliation.
- User Story 3 → Eligibility and limits.
- User Story 4 → Eligibility and limits.
- User Story 5 → Delivery queue owner, Delivery CLI, Action reconciliation.
- Core Feature 1 → Delivery CLI and Eligibility and limits.
- Core Feature 2 → Delivery queue owner and Action reconciliation.
- Core Feature 3 → Review and settlement integration.
- Core Feature 4 → Eligibility and limits.
- Core Feature 5 → Eligibility and limits.
- Core Feature 6 → Eligibility and limits.
- Core Feature 7 → Delivery queue owner and Action reconciliation.
- Core Feature 8 → Canonical delegation.
- Core Feature 9 → Review and settlement integration.
- Core Feature 10 → Eligibility and limits.
- Core Feature 11 → Eligibility and limits and Canonical delegation.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
The [owned source index](references/_index.md) records each primary source.
The [portfolio plan](../../workflow/2026-09-08-pending-work-plan.md) records
secondary consumers and prerequisite Specs. External research was read through
Exa and compared with local Secondbrain history; the PRD and portfolio plan
retain links and describe its effect. Published interfaces support feasibility,
not a claim that the proposed runtime or behavior already exists.

## Testing Approach

Use focused tests at the named package seams for deterministic rules, real
Git/store/process boundaries for integration behavior, and the authored public
QA Task for user-visible acceptance. Do not infer a terminal pass from source
inspection or a focused fixture. Required observations:

1. End the initiating chat/process and observe a real detached Spec Run; restart the delivery owner between Runs and reconcile its next eligible action.
2. Interrupt before and after push, PR creation and merge acknowledgement; prove no duplicate PR/merge and retain actual remote receipts.
3. Change candidate, prerequisite or grant; verify publication/next Run blocks with a durable reason and cannot consume stale review.
4. Exercise Implement deadline and explicit cancellation across Agent/process ownership, then resume without a competing Run or fabricated completed Task.
5. Exhaust each approved queue limit, deny credentials and leave a question unanswered; dependent operations do not occur.
6. Exercise the archive-first late-blocker path with a separately authorized corrective Spec and independent new review in a real repository under its grant.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Enforced Implement deadline and cancellation evidence (depends on: none).
2. Queue/action records and single durable owner (depends on: 1).
3. Read-only preparation and authority/dependency preflight (depends on: 2).
4. Spec execution, settlement and review transitions (depends on: 1, 2, 3).
5. External action receipts and replay reconciliation (depends on: 2, 4).
6. Status/resume/stop, explicit limits and canonical delegation (depends on: 3, 4, 5).
7. Public interrupted-delivery journeys and terminal QA (depends on: 1, 2, 3, 4, 5, 6).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

### Squash delivery and prior authority

For a new or widened governed grant, a separate commit on the consuming branch
is insufficient for squash delivery: the final squash would collapse approval
and consumption. Independently merge the approved grant/amendment into target
ancestry before creating the consuming candidate, then prove that ancestry at
preflight and against the final merged history. A proposal-only planning PR
cannot serve as that grant. The queue parks any consumer whose operative grant
is absent from the target; it never bundles its own new approval into the
consuming squash. Exercise this with a negative same-PR grant case and a
positive independently merged grant followed by the consuming squash.

## Risks & Considerations

Final queue ceilings, archive-first correction trade-off and exact governed grants remain pending. This Spec consumes 0119, 0122, 0125, 0126 and 0129, with runtime/capacity readiness from 0123/0124; it must not bootstrap itself by pretending the new delivery owner already exists.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer policy is settled by [ADR-0151](../../adr/0151-configured-review-profiles-select-the-independent-reviewer.md): Codex by default, explicit review profile in `.roundfixrc.yml` taking precedence.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
