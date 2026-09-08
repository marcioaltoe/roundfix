---
spec: 0123-runtime-readiness-and-model-capabilities
prd: _prd.md
created: 2026-09-08
---

# Runtime readiness and advertised model support — Technical candidate

## Executive Summary

Extend the existing Agent Selection and ACP readiness paths instead of selecting models from a price table. Preserve opaque runtime identifiers and prove the effective tuple and access controls at the real adapter boundary. The performance trade-off proposed for fallbacks changes when proof is required, never whether an activated selection is proved.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: applicable — model IDs are runtime-advertised opaque values; keep API/provider identifiers distinct from ACP Agent Selection IDs and preserve existing Run/Session identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no new authentication or HTTP layer; use current runtime credentials and do not transmit secrets in probes or research. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0107 requires every configured Work Category to be accounted for; ADR-0147 preserves honest advertised capability evidence and the adapter refusal. The current eager Fallback Chain contract is an explicit proposed decision to revisit, not permission to omit readiness evidence. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0049 applies: each present Agent Selection Profile replaces the lower-precedence profile atomically and preserves its explicit fallback chain.
  ADR-0140 applies: prove the exact advertised runtime/model/effort tuple through the installed adapter. Any narrower change to fallback validation timing remains a proposed revision, not an accepted exception.
  ADR-0093 applies to authored capability claims: the consistency checker follows explicit citations and must not infer approval or model support from missing evidence.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Capability evidence | `internal/agent ACP discovery/selection code and CLI profile_preflight.go` | Record the adapter-advertised runtime/model/effort/access tuple and its source/time. |
| Model identity reference | `docs/references/model-selection.md and current runtime selection records` | Keep API, provider-qualified and ACP IDs distinct without changing profiles from advisory pricing. |
| Profile readiness | `internal/config/profiles.go and daemon/agent_session_owner.go` | Preserve atomic profile precedence and declared fallback reasons, proving actual activated choices. |
| Runtime state ownership | `internal/agent ACPXRunner implementation and lifecycle callers` | Eliminate value copies of a state owner containing a mutex while preserving cancellation/session identity. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Exact capabilities and model identities

The requested API model IDs are `gpt-6-astra` and `claude-fable-5-1`. Provider
catalog IDs and adapter-advertised values are separate evidence: do not rewrite
one spelling into another by string convention. The dated pricing reference
remains advisory; it cannot activate a new profile. Record installed runtime
and adapter versions, advertised controls, attempted assignment and effective
selection. Treat model-managed reasoning as an explicit supported state, not
an omitted proof. Characterize access policy alongside model/effort so a
requested permission mode fails before work when unsupported.

A characterization reaches only its explicitly named local runtime/account
under the approved consumption policy. An API catalog GET is not a paid
inference or ACP selection proof. A live probe records its actual result and
cleans up its disposable session; missing access is an explicit unsupported or
blocked outcome, never a passing surrogate based on documentation.

### Atomic profiles and fallback proof

Use the existing built-in→User Config→Project Config resolver and retain each
present profile as a complete preferred/fallback unit. The confirmed reviewer
policy reuses this machinery; do not change the current project tuple.
ADR-0049 and ADR-0140 require advertised tuple proof. The proposed refinement
would prove preferred selection before Run creation and a declared fallback
immediately before activation, recording provenance, reason and outcome.
Approval of this timing refinement is required before replacing the current
eager contract. If the refinement is declined, keep eager proof and still
implement the other readiness fixes. Never use fallback to hide an invalid
configuration or an unproved override.

### One runtime state owner

Make ACPXRunner's mutable session state have one stable owner, passed by pointer
through constructors and call sites that use its mutex/map. Audit interface
satisfaction and copying at every call boundary rather than changing one method
receiver. Retain immutable configuration as value data where useful. Session
creation, cancellation and teardown must converge on the same synchronized
state; no copied-lock analyzer suppression is permitted. The runtime fix must
land before0124 adds vet to the continuous gate.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Separate exact model/provider/runtime identities, advertised capabilities, observed access and pricing provenance. Preserve profile source and one pointer-owned runtime state object; catalog presence alone is not access evidence.

### API Contracts

Doctor/profile preflight reports requested/effective selection and source, capability/access evidence and named incompatibility. A configured unsupported runtime/model refuses; it does not silently substitute an advertised model.

## Coverage Map

- PRD Goal 1 → Capability evidence, Model identity reference.
- PRD Goal 2 → Capability evidence, Profile readiness.
- PRD Goal 3 → Profile readiness.
- PRD Goal 4 → Runtime state ownership.
- Core Feature 1 → Capability evidence.
- Core Feature 2 → Model identity reference.
- Core Feature 3 → Capability evidence and Profile readiness.
- Core Feature 4 → Profile readiness.
- Core Feature 5 → Runtime state ownership.

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

1. Advertised and unadvertised tuples, model-managed effort, unsupported access mode and a requested override absent from the adapter.
2. Distinct API/provider/ACP identifiers and unchanged effective profiles after an advisory model-reference update.
3. Preferred selection succeeds while an unused optional fallback is unavailable under the approved policy; any activated fallback is proved and its reason recorded.
4. Concurrent session start/cancel/teardown uses one state owner; copied-lock diagnostics disappear because the ownership defect is removed.
5. Replay the installed OpenCode catalog omission and record actual Codex/Claude model selection outcomes only when approved access permits a live probe.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Coherent runtime state ownership and lifecycle regressions (depends on: none).
2. Access-policy and exact advertised-selection evidence (depends on: 1).
3. Approved fallback proof timing and provenance (depends on: 1, 2).
4. Bounded real-model characterization and identifier documentation (depends on: 2, 3).
5. Readiness diagnostics/owned guidance and terminal QA (depends on: 1, 2, 3, 4).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Live capability proof may consume account quota and therefore waits for the recorded limit. Do not substitute the parent assistant's available models for ACP support. A fallback timing change remains a proposed ADR refinement; preserving eager proof is the conservative alternative.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
- Lazy fallback readiness changes the existing eager-validation decision. ADR-0049/ADR-0140 remain operative until the narrower timing revision is explicitly accepted; preserve rejection reasons and profile provenance in either case.
