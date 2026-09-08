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
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

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
  ADR-0142 applies to the existing head-bound evidence contract and historical outcomes; the configurable pre-PR policy must be settled by Spec 0126 before the new delivery path relies on it.
  Preserve the accepted decision to keep `implement-spec` and repair its conflicting instructions rather than introducing a competing implementation loop.
  ADR-0151 preserves the Codex default and explicit project selection for agent review.
  ADR-0153 applies: the queue consumes codex, claude, coderabbit or explicit none, preserving QA/checks and recording configured omission without treating enabled-review failures as none.
  ADR-0154 applies to archive disposition: an explicitly user-authorized QA Archive Override preserves actual QA/Task evidence and does not satisfy independent delivery gates.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/implement-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `internal/baseline/assets/modules/autonomous-work.json`, `skills/implement-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `skills/baseline_skill_contract_test.go`, `internal/docscontract/publicdocs_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Delivery queue owner | `internal/daemon and internal/store` | Persist queue, ordered items, owner identity and action intent/receipt separately from actual Runs. |
| Delivery CLI | `internal/cli/delivery.go (new) and CLI registration` | Prepare/start/status/resume/stop through one public command family. |
| Eligibility and limits | `internal/config, daemon preflight and Run lifecycle` | Validate approved Spec revision, dependencies, authority, Run Window, enforced Run deadline and distinct queue limits. |
| Action reconciliation | `existing Git/GitHub boundaries and store transactions` | Record intent before an external action and reconcile its observed result before retry. |
| Review and settlement integration | `Spec 0122 settlement and Spec 0126 candidate review` | Require the selected policy outcome for the current candidate before publication and merge; none supplies a configured omission, not a review verdict. |
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
review/omission record and policy provenance, PR identity and observed merge commit. An action row contains
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
contract; commit the final candidate; obtain Spec 0126 configured review or explicit omission;
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

### Policy-aware recovery

Persist the effective review policy and its source with each candidate. For an
enabled provider, require its complete current-candidate receipt. For explicit
none, persist configured omission without creating a reviewer, probing reviewer
readiness or invoking any provider. Revalidate policy and candidate before
publication and after recovery. A changed policy, stale receipt or enabled
provider failure cannot silently become none. QA, required checks and external
branch protection remain binding in every mode.

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

Queue, queue item and action records are distinct from Runs. Persist approved source revisions/limits, actual Run links, candidate/policy-outcome/PR/merge identities, action intent/receipt and blockers; recovery never invents an actual Run.

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

1. End the initiating chat/process and observe a real detached Spec Run; restart the delivery owner between Runs and reconcile its next eligible action.
2. Interrupt before and after push, PR creation and merge acknowledgement; prove no duplicate PR/merge and retain actual remote receipts.
3. Change candidate, prerequisite or grant; verify publication/next Run blocks with a durable reason and cannot consume stale review.
4. Exercise Implement deadline and explicit cancellation across Agent/process ownership, then resume without a competing Run or fabricated completed Task.
5. Exhaust each approved queue limit, deny credentials and leave a question unanswered; dependent operations do not occur.
6. Exercise none across restart and publication with zero reviewer calls; an enabled-review failure never selects none and failed required checks still block merge.
7. Exercise the archive-first late-blocker path with a separately authorized corrective Spec and new review when enabled in a real repository under its grant.

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
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Review policy is settled by ADR-0153: codex, claude, coderabbit or explicit none. ADR-0151 retains the Codex default and explicit project precedence.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md), [0122](../0122-verified-content-and-terminal-settlement/_techspec.md), [0125](../0125-repository-identity-and-run-branch-policy/_techspec.md), [0126](../0126-agent-review-before-pull-request/_techspec.md), [0129](../0129-spec-authoring-and-gate-recovery/_techspec.md), [0123](../0123-runtime-readiness-and-model-capabilities/_techspec.md), [0124](../0124-verification-capacity-and-measured-economics/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.

## Portfolio sequence retained from the removed routing plan

Proposed serial order is 0119, 0123, 0126, 0122, 0121, 0129, 0125, 0124, 0120,
0128, then 0127. Each preceding Spec defines its own protected scope and source
ownership. This order prepares authority and configured native review early,
then integrity, authoring, identity and measured capacity before the durable
queue is enabled. Revalidate the order and shared paths after each merge.

The existing one-Spec Implement command is the bootstrap executor. Author and
validate an approved Task Graph before dispatch; the proposed durable queue
cannot implement its own prerequisites before it exists. No Task Graph or
implementation readiness is established by this sequence.

Reviewer selection and conditional delivery through squash merge are confirmed.
Consumption policy, finite queue limits, archive-first correction authority and
the named ADR revisions remain pending. New governed grants must already be
independently merged into target ancestry before a consuming squash.

The primary source indexes in Specs 0119–0129 account for 23 ordinary sources;
the six archived Rollups record 82 historical members, now linked directly to Specs. These counts describe
routing, not implementation. The removed cross-Spec plan is available at Git
revision `6b8ea48725cbca13974eee0b400b3482202874f6` for historical inspection.


## Readiness recheck — 2026-09-08

The maintainer asked whether the current skills/runtime can deliver without
supervision and whether `_authorization.md` has a canonical template. Local
inspection distinguishes available components from this proposed delivery owner:

- The CLI dispatch in `internal/cli/cli.go` exposes Implement, Run events,
  window, settle, reconcile and archive. Its current help exposes no `delivery`
  or native pre-PR `review` command. Detached one-Spec execution exists; durable
  ownership of the whole delivery queue remains proposed here.
- `.agents/skills/roundfix/SKILL.md` describes Supervisor-driven delivery through
  merge. `.agents/skills/implement-spec/SKILL.md` still describes a direct Task
  loop and forbids autonomous publication. The inspected Fluxus kickoff skill
  still permits delivery without an independent reviewer when no review source
  exists. These instructions need the alignment already scoped here. The later
  confirmed four-choice policy permits explicit none, but absence of a legacy
  review source alone cannot select none or override the Codex default.
- `.roundfixrc.yml` still names `review_source.name: coderabbit`. Its separate
  `profiles.review` selects Codex/gpt-5.6-luna/max with the declared Sol fallback;
  that selection is not proof of an implemented native review adapter.
- Specs 0119–0129 have no `_tasks.md` or `task_*.md`. Their technical candidates
  and proposed authorizations do not constitute an executable approved queue.
- Searching `internal/baseline/assets/modules/`, generated guides and the
  authoring skills found no canonical `_authorization.md` template. Existing
  PRD/TechSpec templates put tooling authority in Project Constraints. Spec
  0119 explicitly calls the new filename/schema proposed; current examples
  even differ between scalar and list forms of `consuming`. Its implementation
  must standardize the source template, generated guidance and consuming
  readers together before a Spec-local proposal is advertised as that contract.

Secondbrain was consulted through its index and a qmd query about Roundfix
kickoff, autonomous queues, recovery and independent review. The verified
entries `inbox/secondbrain/2026-09-08-configured-reviewer-and-autonomous-queue-readiness.md`
and `inbox/secondbrain/2026-09-08-revisao-local-antes-do-pr-e-supervisao-duravel.md`
already preserve the same distinction between detached Runs and queue ownership;
this recheck confirms that distinction against local sources without treating
historical mirror paths as current files.

Exa MCP returned the opening excerpt of Anthropic's
[Effective harnesses for long-running agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents)
(published 2025-11-26). It describes failures across context windows and why
compaction alone is insufficient. That limited excerpt supports the recovery
concern, not a claim about Roundfix implementation or measured reliability;
no full-article review is claimed. No new external design conclusion is added
beyond the verified existing Secondbrain captures. This was source/help
inspection only: no live Run, runtime readiness probe, paid model call or
terminal QA was performed.


## Autonomy-first implementation order — 2026-09-08

For the maintainer's objective of reaching unattended delivery, the recommended
serial order is **0119 → 0123 → 0126 → 0122 → 0121 → 0129 → 0125 → 0124 → 0120
→ 0127**, followed by **0128**. This recommendation supersedes the placement
of 0128 in the historical portfolio sequence above; it changes no declared
hard dependency, operative grant or approval state.

| Position | Spec | What it supplies before the queue is enabled |
| --- | --- | --- |
| 1 | 0119 | Canonical Spec-contained authorization, approved command sources and bounded actions. |
| 2 | 0123 | Proven runtime/model/access capabilities for execution and review. |
| 3 | 0126 | Configured codex/claude/coderabbit review of the current candidate, or explicit none; no universal provider dependency. |
| 4 | 0122 | Settlement of the exact verified content, terminal QA and supported repair entry. |
| 5 | 0121 | Consistent Baseline decisions and complete sanctioned regeneration. |
| 6 | 0129 | Authored obligations, changed-premise handling and supported gate recovery. |
| 7 | 0125 | Shared repository identity, Run-branch integrity and reconciliation. |
| 8 | 0124 | Measured Verification capacity, cost and reliability. |
| 9 | 0120 | Source adoption without retained originals, terminal lifecycle and durable knowledge capture. |
| 10 | 0127 | Persistent queue ownership, enforced limits, reviewed publication and crash-safe delivery recovery. |
| After activation | 0128 | Bare stable-tag release planning; releases are outside the approved through-merge delivery scope. |

The order was checked against every sibling TechSpec's declared predecessor
list. The transitive prerequisite set for 0127 is 0119, 0121, 0122, 0123, 0124,
0125, 0126 and 0129. Neither 0120 nor 0128 belongs to that declared set.
Putting 0120 before activation is a readiness recommendation based on the
recently reproduced source-adoption failures, not a silently added graph edge.
0128 may follow activation because tag planning is not needed to implement,
review, open a PR or squash merge.

Before dispatch, settle the current cleanup's verification/regeneration
compatibility failures, complete each consuming Spec's technical decisions and
exact authority, and author its validated Task Graph with terminal QA. None
of Specs 0119–0129 currently has `_tasks.md`; the sequence alone does not
make them runnable. Required new/widened tooling grants must be independently
merged into target ancestry before their consuming squash. Existing session
publication/reviewer decisions remain settled; only unanswered decisions stay
open.

Use existing one-Spec Implement Runs during this bootstrap, coordinated by the
supervising session. Until 0126 is delivered, that session follows the explicitly selected policy
through available supported tools before publication; it cannot invoke a future
adapter or mistake a failed/skipped provider for an explicit none choice.
The future 0127 queue cannot bootstrap itself before it exists. Revalidate the
next Spec's source assumptions after each prerequisite merge and preserve the
existing verification gates throughout preparation and delivery.

Activation requires 0127's real interruption/recovery journeys: the initiating
session exits; execution reaches terminal QA, the configured review or explicit
omission, passing required PR checks and observed squash merge; interruption
between external success and local receipt resumes without duplicate actions.
Limits and unresolved authority stop the dependent action with durable evidence.
Source inspection or a green Task Graph is not that acceptance proof.

The local dependency declarations determine this order. Secondbrain's verified
`2026-09-08-configured-reviewer-and-autonomous-queue-readiness.md` capture confirms
the historical distinction between detached Runs and durable queue ownership;
its old source-routing counts are not used as current inventory. Exa MCP again
returned the opening excerpt of the Anthropic long-running harness article
linked in the preceding recheck. Its context-window limitations support the
need for durable recovery, not this repository's particular ordering; that
ordering is an engineering recommendation from the local dependency graph.


## Authorized archive disposition

Consume the archive policy from Spec 0122 and ADR-0154. An applicable explicit
user authorization can archive the covered Spec with unmet QA, recording the
override and preserving original evidence. This does not reopen or complete
Tasks, declare QA passed, imply review approval or authorize publication/merge.
A durable workflow records the overridden archive and evaluates subsequent
actions against their own approval and gates; it does not retry the waived
archive prerequisite or ask again for the same applicable archive approval.
The absence of authority for a later action remains a separate visible blocker.
