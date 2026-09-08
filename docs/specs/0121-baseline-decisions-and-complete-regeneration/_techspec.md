---
spec: 0121-baseline-decisions-and-complete-regeneration
prd: _prd.md
created: 2026-09-08
---

# Lossless Baseline decisions and complete regeneration contracts — Technical candidate

## Executive Summary

Repair the existing typed decision and transactional planning paths. The design keeps current decision identities and uses the existing ownership registry rather than creating parallel generators. Legacy adoption becomes an explicit reviewed migration, while an unreachable source remains insufficient evidence for destructive reconciliation.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity or identifier strategy is proposed; preserve existing Baseline and decision identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the proposal changes how the CLI edits a repository-owned HTTP Contract, preserving its confirmed mode, typed exceptions, and source; it introduces no application endpoint or authentication policy. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 keeps regenerated pins as sanctioned fallout. ADR-0149 gives the tree ownership of regeneration outputs. ADR-0130 keeps the audit's governed path set honest against authorization history. Retain semantic preservation rather than treating stale source as disposable. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/contract-v1.json`, `internal/cli/baseline_human_test.go`, `internal/baseline/derived_ownership_test.go`, `skills/_ownership.yml`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `internal/baseline/assets/decisions.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/go-cli-tui.json`, `internal/baseline/assets/profiles/rust-cli.json`, `internal/baseline/assets/templates/index.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`, `internal/baseline/plan_test.go`, `.agents/skills/setup-context-driven/SKILL.md`, `skills/setup-context-driven/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Typed HTTP editor | `internal/cli/baseline_human.go and internal/baseline/project_decisions.go` | Change only the selected HTTP mode and preserve exceptions and provenance. |
| Profile/default compatibility | `internal/baseline/profile_alignment.go, profile assets and contract-v1.json` | Remove the inert HTTP default without losing the capability declaration. |
| Reachable adoption refusal | `internal/baseline/preservation.go and CLI classification flow` | Refuse the Greenfield/classification contradiction with a usable Preservation route. |
| Declared skill regeneration | `internal/baseline/derived_ownership.go and skills/_ownership.yml` | Resolve make skills-sync output ownership alongside the existing Baseline declaration. |
| Two Verification decisions | `decision assets, VerificationProjection, custom Profile and update readers` | Represent complete and incremental commands with their actual local declaration evidence. |
| External lock reconciliation | `internal/baseline/skills_restore.go, skills_restore_lock.go and skills_restore_git.go` | Preview and apply only positively proven obsolete lock entries through the existing transaction. |
| Public Baseline guidance | `owned Roundfix/setup skills and managed guides` | Expose the repaired decisions, refusal paths and exact planned mutations. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Preserve typed values and capability meaning

Reuse the existing object clone and `normalizeHTTPContract` validation. Editing
mode replaces that member only; unchanged exceptions and source provenance
remain present. Remove the unused Profile HTTP default and stale diagnostic
without deleting the meaningful HTTP capability object used by
`profileHasHTTPContract`. Canonical defaults and explicit repository decisions
retain separate owners.

A Greenfield plan cannot require an interaction whose classifier deliberately
runs only in Preservation. Detect that contradiction before prompting and
return the same actionable Preservation outcome in human and JSON planning.
Do not synthesize classifications or silently rewrite retained managed source.

### Verification tiers and legacy migration

Add `verification.incremental` beside the preserved `verification.gate` complete
command. Extend the existing Verification Projection role/command/provenance
shape and render both selected values. Built-in and custom Profiles use the
same representation. A missing incremental value in a legacy manifest is an
explicit decision in the public update plan. Validate the repository declaration
through the existing command-declaration reader; do not copy the complete gate
as an invented default or create a target in another project. Use the captured
Fiscus shape as outside evidence in a disposable repository.

### Generation and external locks

Make `OutputsFor` command-aware across the existing Baseline ownership tree and
the owned skill tree. `make skills-sync` owns only the embedded copies of the14
owned skills, not authorial `.agents/skills`, unrelated Go files or external
skills. The audit consumes this declaration; adding an ownership file without
extending the hardcoded resolver is incomplete. Keep the generator itself
unchanged unless a measured necessity is added to the grant.

Lock reconciliation reuses strict ordered JSON, immutable Git acquisition,
preview digest, preimage check and transactional rollback. Classify removal at
an explicitly selected immutable upstream revision. A missing required skill is
a blocking requirement conflict. Network/authentication errors, local absence
and an optional Profile omission cannot delete an entry. Initially remove only
approved obsolete lock entries, retain installed trees and unknown fields, and
report their separate disposition. Doctor remains offline/read-only. No actual
consumer lock is changed without its own exact approved plan.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Keep typed Decision values and source/exceptions together, distinguish full and incremental Verification decisions, and retain command-owned output declarations. Lock reconciliation records the selected immutable upstream revision, preview and preimage before apply.

### API Contracts

Use public Baseline plan/apply and Doctor boundaries. Plans expose typed changes and owned outputs; apply checks preimages and converges. Doctor remains offline/read-only. A failed upstream read cannot authorize lock removal.

## Coverage Map

- PRD Goal 1 → Typed HTTP editor.
- PRD Goal 2 → Profile/default compatibility.
- PRD Goal 3 → Reachable adoption refusal.
- PRD Goal 4 → Declared skill regeneration.
- Core Feature 1 → Typed HTTP editor.
- Core Feature 2 → Profile/default compatibility.
- Core Feature 3 → Reachable adoption refusal.
- Core Feature 4 → Declared skill regeneration.
- Core Feature 5 → Public Baseline guidance.
- Core Feature 6 → Two Verification decisions.
- Core Feature 7 → External lock reconciliation.

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

1. HTTP mode edits preserve two independent exceptions/provenance; deliberate exception changes are observable and validated.
2. Greenfield refuses before impossible interaction; the prescribed Preservation plan applies and replans idempotently.
3. Built-in and custom Profiles expose both verified commands; legacy manifests migrate only after the missing decision and actual declaration are supplied.
4. Ownership output sets match actual make skills-sync outputs and exclude authorial/external/Go files; a manual generated mutation still fails.
5. Known removed, required-but-removed, unreachable, existing and unrelated locked skills; stale preview and failed transaction preserve exact preimages.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Lossless HTTP editor and capability/default cleanup (depends on: none).
2. Reachable Greenfield refusal and preserved adoption route (depends on: 1).
3. Explicit incremental/complete decision representation and migration (depends on: 1, 2).
4. Owned-skill generation declaration and audit resolution (depends on: 1).
5. Immutable-source external lock reconciliation (depends on: 2, 4).
6. Public guidance and sanctioned regeneration (depends on: 3, 4, 5).
7. Public plan/apply/replan and terminal QA (depends on: 1, 2, 3, 4, 5, 6).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

This Spec changes shared canonical assets and skills, so later Specs must rebase their exact preimages before mutation. Early grants explicitly enumerate shipped paths until ownership resolution is delivered. Preserve optional installed skills and distinguish unavailable sources from proved absence.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
