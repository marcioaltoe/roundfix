---
spec: 0119-spec-contained-authorization
prd: _prd.md
created: 2026-09-08
---

# Spec-contained authority and trusted Verification — Technical candidate

## Executive Summary

Extend the existing Spec and consistency-checking packages with one role-based authorization reader. The trade-off is explicit source/scope evidence at execution boundaries instead of inferring permission from a filename, a passing checker or matching prose. Preserve historical grants and their bounded paths; do not build a new permission service.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, Task IDs, and existing diagnostic identities; authorization belongs to one consuming Spec rather than a new global identity service. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no authentication system or HTTP endpoint is introduced. Existing credentials remain in their current runtime boundary. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0014 keeps Verification with the Daemon. Source: `docs/agents/spec-routing.md`.
  ADR-0130 keeps the audit on governed paths and preserves the historical governed set; this Spec changes grant discovery without exempting those paths.
  ADR-0149 makes the grant name the regeneration command and the ownership tree name its outputs; preserve that division.
  ADR-0057 keeps the Daemon the exclusive writer of Implement Task status; authorizations do not let Agents settle themselves.
  ADR-0096 requires mechanical facts, including bounded-path authority, before the QA Agent turn; extend that check rather than bypassing it.
  ADR-0117 places a defect check at the stage that can produce it; validate authored grants at authoring and actual changed files after commits exist.
  ADR-0020 is not applicable to the implementation scope: parsed prompt result versus acpx teardown exit classification remains unchanged.
  ADR-0038 is not applicable to the implementation scope: the one Verification repair allowance remains unchanged and is not widened by a grant.
  ADR-0056 is not applicable to the implementation scope: Task Capacity, Verification Capacity, and the temporary-failure retry remain unchanged.
  ADR-0127 is not applicable to the implementation scope: reporting process residue as a readiness fact remains unchanged.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/modules/context-workflow.json`, `internal/speccheck/constraints.go`, `internal/speccheck/constraints_characterization_test.go`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`, `skills/write-tasks/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/speccheck/governed.go`, `internal/speccheck/governed_repocontract_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Authorization reader | `internal/spec/authorization.go (new), existing internal/spec package` | Parse operative record fields and normalize bounded legacy records without accepting a proposed record as a grant. |
| Authoring authority checks | `internal/speccheck/constraints.go and citations.go` | Apply the same typed semantics to an undated Spec-contained record and preserve precise refusal locations. |
| Commit authority audit | `internal/speccheck/mechanical.go and governed.go` | Prove the relevant grant existed before each consuming commit and that actual governed changes stay inside its scope. |
| Command-source approval | `internal/cli/spec_check.go, implement.go and settle.go` | Validate the committed authored-command source before probing or executing it. |
| Authoring guidance | `owned PRD, TechSpec and Task skills; canonical core/spec/context modules` | Place records inside Specs and explain approved versus proposed authority consistently. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Operative record and source identity

Use an explicitly versioned record role with `status`, `consuming`, granted date,
actions, exact governed paths, sanctioned regeneration and resource limits.
A proposal is inspectable but grants nothing. Unknown states, null approval,
wrong consuming fields, escaping/symlinked paths and contradictory records are
refusals. Matching a Spec name somewhere in prose is not a consuming relation.
Preserve the current explicit legacy formats through a bounded adapter rather
than rewriting archived records. A multi-Spec historical grant retains its
actual consuming list and scope; a new Spec-owned record has one primary owner.

The authored source identity must separate requirements from execution results:
PRD/TechSpec contracts, Task Graph structure and executable Task requirements
are approved inputs; Daemon-owned status and Result updates are execution
records, not new permission. Record the committed source and the exact commands
selected for execution. Imported external Spec roots need explicit source
approval. A disposable checkout isolates output but does not authorize the
commands copied into it. An amendment that changes authority, commands or the
approved source boundary is re-evaluated before execution; no Agent can edit its
own grant to make a change valid.

### Audit and compatibility

Read the grant from the ancestor that authorized the consuming commit, retaining
its path and digest in the audit result. A later amendment cannot authorize an
earlier change retroactively. Reject a consuming commit that creates or widens
its own grant, while accepting an earlier bounded amendment for later work.
Discover operative records in active Specs, archived Specs and the legacy
workflow directory. Keep the historical Governed Path set monotonic, including
owned shipped templates currently missed by the narrow path predicate.

Apply source approval at all authored-command entry points: non-vacuous probing,
Implement dispatch and Settle. A read-only `spec check` remains usable on
untrusted documentation. Source approval does not grant network access,
credentials, release actions or a different sandbox. These fields are a design
candidate; their exact protected readers/templates await the operative grant.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Parse a grant by role with its consuming Spec, approval evidence, exact paths/actions, source revision and limits. Preserve legacy dated records and immutable prior decisions; proposed and absent approval are distinct from granted.

### API Contracts

Spec Check and Implement preflight consume the same approval/source decision. A refused command-source approval must prevent authored shell execution; settlement audits actual changed paths against the committed operative grant.

## Coverage Map

- PRD Goal 1 → Authorization reader, Commit authority audit, Authoring guidance.
- PRD Goal 2 → Authorization reader, Authoring authority checks, Commit authority audit.
- PRD Goal 3 → Command-source approval.
- PRD Goal 4 → Authorization reader, Commit authority audit.
- User Story 1 → Authorization reader, Commit authority audit, Authoring guidance.
- User Story 2 → Authoring authority checks, Command-source approval.
- User Story 3 → Commit authority audit.
- Core Feature 1 → Authorization reader.
- Core Feature 2 → Commit authority audit.
- Core Feature 3 → Authoring guidance.
- Core Feature 4 → Authorization reader and Authoring authority checks.
- Core Feature 5 → Authorization reader and Command-source approval.
- Core Feature 6 → Command-source approval.

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

1. Valid approved record; proposed/null/withdrawn/wrong-Spec records; a prose mention that does not match the consuming field; path traversal and symlink escape.
2. Earlier grant, earlier amendment, same-commit self-approval and later retroactive approval across a real temporary Git history.
3. An archived Spec-contained grant and unchanged legacy grant resolve; discovered paths never narrow the historical governed set.
4. A supplied third-party Verification marker cannot execute before source approval; unchanged approved local commands can run in a disposable checkout.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Typed authorization reader and legacy adapter (depends on: none).
2. Authoring diagnostics and all record locations (depends on: 1).
3. Ancestor grant audit and governed template coverage (depends on: 1, 2).
4. Command-source approval at probe, Implement and Settle (depends on: 1, 3).
5. Canonical authoring guidance and sanctioned generated outputs (depends on: 2, 4).
6. Public refusal/approval journeys and terminal QA (depends on: 3, 4, 5).

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

The main risk is circular source identity or retroactive permission. Keep execution Result/status outside the authored contract identity and bind approval to the real ancestor record. Preserve legitimate legacy grants without allowing the legacy reader to accept arbitrary prose.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
