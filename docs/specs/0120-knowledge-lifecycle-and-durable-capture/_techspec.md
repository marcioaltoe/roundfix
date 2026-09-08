---
spec: 0120-knowledge-lifecycle-and-durable-capture
prd: _prd.md
created: 2026-09-08
---

# Knowledge lifecycle, publication evidence and review retirement — Technical candidate

## Executive Summary

Extend existing lifecycle readers and canonical capture rules instead of adding another Inbox or scheduler. The trade-off is retaining explicit unknown lifecycle outcomes until durable evidence exists. A Review Artifact relocation is a path migration and must not decide whether its review is active.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: applicable — retain existing Spec slugs, dated artifact basenames, Inbox destinations, and source provenance; a terminal disposition does not invent an implementation identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, authentication, or HTTP API change is proposed. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — preserve the one-way mirror and fleet Inbox boundary, the existing History Root, and accepted regeneration ownership under ADR-0149. Source: `docs/agents/docs-layout.md`.
  ADR-0123 remains operative: retirement currently uses conservative local reachability and the Review Artifact resolver never writes into history. The proposed stable-evidence change does not yet supersede it.
  ADR-0152 is a proposed revision only, recorded for review; it creates no current obligation or grant.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/secondbrain.json`, `internal/spec/archive.go`, `internal/spec/archive_test.go`, `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`, `internal/docscontract/publicdocs_test.go`, `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Capture contract | `canonical context-workflow/secondbrain modules and companion Secondbrain Inbox contract` | Name one fleet door and distinguish local capture, local commit and observed remote publication. |
| Terminal disposition | `internal/spec/retirement.go and internal/speccheck/backlog.go` | Represent evidence-backed fulfilled or superseded intent without misclassifying it as declined. |
| History and licenses | `internal/spec/archive.go and internal/speccheck/citations.go` | Resolve all retired families and replacement license chains while preserving old records. |
| Upstream citation check | `internal/docscontract/publicdocs_test.go` | Reject evidentiary citations from glossary/guides into active or archived downstream cases. |
| Review Artifact retirement | `existing Review Artifact resolver and retirement readers in internal/spec` | Separate stable lifecycle proof from legacy path relocation and local object availability. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Capture and publication

Keep `inbox/<destination>/` in Secondbrain as the fleet door and the project as
the owner of its triage. A capture receipt distinguishes a created file, a local
commit and confirmed remote visibility. Use the existing publication owner;
record a bounded session-owned fallback only when that scheduler is absent.
The installed scheduler's last exit and mere presence cannot certify a new
capture's remote durability. Companion changes belong to Secondbrain's own
branch and authorization, never its mirrors or immutable raw sources.

### Dispositions and archive ownership

Extend the existing backlog/finding readers with a reasoned terminal disposition
for directly fulfilled and superseded work. Preserve old `done` absorption
meaning: a Spec owner is not implementation proof. Validate the evidence/ref
that supports terminal disposition and reject cycles, missing replacements,
path escape or claims resting only on a status label. History resolution covers
handoffs as a separate family and preserves basename/provenance on moves.

Six active Rollups currently license82 historical members. Keep these roots in
place during this queue's authoring. A future Rollup retirement is eligible only
when every member has a resolvable durable replacement chain; do not repoint
historical evidence merely to empty a directory. Migration is previewed and
compares exact preimages before applying its bounded path changes.

### Review retirement and upstream boundaries

ADR-0123 remains operative until the proposed replacement in ADR-0152 is
approved. The new retirement result records stable provider/recorded evidence,
including a squash outcome, or an explicit unknown reason. Fetching or pruning
a local Git object must not change the lifecycle verdict. Move a legacy
`docs/specs/_reviews` path independently of that verdict, retaining its active
or unknown state. Consume the local-candidate receipts from0126 when available;
missing historical receipts do not become successful evidence.

The citation detector reads the glossary and active agent guides, distinguishes
literal layout declarations from actual evidence links, and rejects active and
history Spec/Finding references as upstream authority. Do not exempt the entire
layout guide: it can contain a real forbidden citation. Keep this detector in
the existing docscontract boundary with positive and negative fixtures.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Retain source provenance, terminal disposition and absorber/license identity. A capture durability receipt identifies the observed durable destination; local file creation alone cannot attest remote storage.

### API Contracts

Extend existing capture, lifecycle and consistency boundaries. Any retirement or migration command must name its exact source/destination and refusal conditions; scheduler installation and external automation remain separately scoped.

## Coverage Map

- PRD Goal 1 → Capture contract.
- PRD Goal 2 → Terminal disposition.
- PRD Goal 3 → History and licenses, Review Artifact retirement.
- PRD Goal 4 → Upstream citation check.
- Core Feature 1 → Capture contract.
- Core Feature 2 → Capture contract.
- Core Feature 3 → Terminal disposition.
- Core Feature 4 → History and licenses.
- Core Feature 5 → History and licenses.
- Core Feature 6 → Upstream citation check.
- Core Feature 7 → Capture contract.
- Core Feature 8 → Review Artifact retirement.

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

1. Capture states file/commit/remote-confirmed and scheduler-absent fallback; no mirror or raw writes.
2. Fulfilled, superseded, unresolved and declined intent with actual disposition evidence; merely changing status cannot close a defect.
3. Every one of the82 existing license edges survives; replacement chains accept valid history and reject missing/cyclic/escaping links.
4. The same Review Artifact classifies identically before/after Git object availability changes; legacy relocation preserves liveness and exact observations.
5. Actual downstream evidence citations in active/history trees fail; directory-layout examples and accepted ADR links pass.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Capture/publication contract and companion changes (depends on: none).
2. Evidence-backed terminal dispositions (depends on: 1).
3. History-family and replacement-license resolution (depends on: 2).
4. Stable Review Artifact retirement and separate legacy relocation (depends on: 2, 3).
5. Upstream citation detector and canonical guidance (depends on: 1, 3, 4).
6. Cross-repository examples, migration negatives and terminal QA (depends on: 1, 2, 3, 4, 5).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Retirement revises an accepted policy and is pending that decision. A source move must not erase evidence or leave a dangling archive license. Cross-repository publication must retain the actual owner and observed receipt rather than treating either checkout as the other.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
- [Proposed ADR-0152](../../adr/0152-review-artifact-retirement-uses-stable-evidence.md) challenges the object-dependent retirement rule in ADR-0123; it does not supersede that rule yet.
