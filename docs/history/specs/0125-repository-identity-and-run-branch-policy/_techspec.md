---
spec: 0125-repository-identity-and-run-branch-policy
prd: _prd.md
created: 2026-09-08
---

# Shared repository identity and evidence-based reconciliation — Technical candidate

## Executive Summary

Derive repository identity from the shared Git repository boundary while retaining aliases for existing checkout-derived records. Preserve Roundfix Run/Task names and the delivered purpose-based work policy. The cleanup trade-off remains conservative: unproved content stays recoverable instead of being deleted because a ref disappeared.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

## Project Constraints

- Identifier strategy: applicable — canonical repository identity must group linked worktrees while preserving access to existing Run records; keep Run/Task IDs stable and recognize old branch names during migration. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local repository/ref identity and reconciliation introduce no authentication or HTTP endpoints. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0138 preserves one commit per verified Task and the existing Clean-only opt-in push boundary. ADR-0135 requires absent diagnostic output to be reported as an explicit state, not an empty message; it does not define whether a Git reference exists. Keep branch disposition evidence-based. ADR-0150 confirms purpose-based work branches and the existing Roundfix Run namespace. ADR-0118 remains relevant for compatibility of the saved decision key, with its personal-prefix aspect replaced. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Delivered naming policy | `ADR0150 and existing canonical branch guidance` | Preserve the already merged policy as a regression boundary; do not implement it again. |
| Repository identity | `internal/store repository-key resolution and internal/worktree discovery` | Group linked worktrees by their actual shared Git repository without conflating unrelated clones. |
| Legacy identity aliases | `Run Database lookup and repository discovery callers` | Keep previous checkout-derived Run records discoverable during migration. |
| Content reconciliation | `existing Reconcile, Run integration and Git evidence readers` | Retain deleted-target proof and add bounded present-target squash evidence. |
| Public recovery output | `owned Roundfix skill and CLI diagnostic paths` | Name missing refs, preserved content and supported next actions truthfully. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Identity boundary

Use Git's common-directory relationship to recognize linked worktrees of one
repository. Resolve the actual shared metadata boundary, normalize its identity
and retain existing Run/Task IDs. A matching remote URL is not enough to merge
two unrelated clones. Symlinked paths and moved/removed linked checkouts need
explicit normalization and lookup evidence. Store aliases from prior
checkout-derived keys so old Runs remain discoverable rather than moving or
recreating their records speculatively.

Migration is read/preview first, with collision and ambiguous-alias refusal.
Never drop the only path to a Run or change its ownership while Active. A
repository-level lookup must use the same resolver across start, browser,
reconcile and storage cleanup; fixing only the creation path leaves old records
stranded. Keep the existing one-owner and clean-checkout constraints.

### Reconciliation proof

Preserve0097's positive and negative deleted-target content tests and do not
repeat that completed implementation. Characterize a present-target squash,
a truly absent ref and retained Run-only content separately. Compare the
relevant delivered content and recorded merge evidence, preserving exact ref
identity, rather than treating an archive marker or mere ancestor test as
sufficient. When an accepted proof is unavailable or ambiguous, report the
missing proof and keep the branch/worktree intact.

The two PRs delivered on2026-09-08 supply an outside example of why squash
creates different commit ancestry even when candidate trees match. A recorded
GitHub merge receipt plus the actual resulting tree was used for those manual
delivery cleanups; it is not proof that current automatic Reconcile supports
every equivalent case. The new tests must exercise the public command and its
failure preservation behavior.

### Naming and guidance

The naming subset is already shipped in PR177 and has its own dated approved
record. New work branches express purpose, while Roundfix owns its existing
Run/Task namespace. No mass rename, new namespace migration or bootstrap
exception is part of remaining implementation. Update only recovery examples
and identity explanations made stale by the actual change.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Use Git common-directory identity for linked worktrees and retain explicit legacy aliases. An unrelated clone with the same remote is not silently the same local repository; content reconciliation keeps its proof and uncertain outcomes.

### API Contracts

Existing status/reconcile commands resolve repository identity and report observed refs and delivery proof. Missing or uncertain state is named, never treated as permission to delete content.

## Coverage Map

- PRD Goal 1 → Delivered naming policy.
- PRD Goal 2 → Repository identity.
- PRD Goal 3 → Legacy identity aliases.
- PRD Goal 4 → Content reconciliation, Public recovery output.
- Core Feature 1 → Delivered naming policy.
- Core Feature 2 → Repository identity and Legacy identity aliases.
- Core Feature 3 → Content reconciliation.
- Core Feature 4 → Content reconciliation.
- Core Feature 5 → Public recovery output.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
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

1. Main checkout plus linked worktrees share identity; unrelated clones with the same remote remain distinct.
2. Legacy aliases preserve discovery after a linked path is moved/removed; collisions and Active Run ownership conflicts refuse without mutation.
3. Deleted-target proof remains correct; present-target squash and generic missing refs have independent positive/negative cases.
4. A Run-only file, changed executable or differing content stays preserved even when the Spec is archived or a PR is reported merged.
5. Public diagnostics distinguish missing, unavailable, unintegrated and proved-superseded evidence without an unsafe cleanup prescription.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Shared Git repository identity and lookup normalization (depends on: none).
2. Legacy alias migration and discovery compatibility (depends on: 1).
3. Present-target squash characterization and bounded proof (depends on: 1, 2).
4. Public reconciliation diagnostics and owned guidance (depends on: 2, 3).
5. Real linked-worktree/Git negative journeys and terminal QA (depends on: 1, 2, 3, 4).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Identity migration must not combine independent clones or orphan old records. Cleanup stays conservative when delivery proof is incomplete. Existing naming and deleted-target behavior are regression obligations, not new implementation.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Preserve configured reviewer selection; this Spec introduces no separate reviewer override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md), [0122](../0122-verified-content-and-terminal-settlement/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.
