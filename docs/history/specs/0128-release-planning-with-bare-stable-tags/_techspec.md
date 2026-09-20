---
spec: 0128-release-planning-with-bare-stable-tags
prd: _prd.md
created: 2026-09-08
---

# Read-only release planning with exact stable tag identity — Technical candidate

## Executive Summary

Extend the existing stable-version parser and Git discovery to accept strict prefixed and bare tags. Separate semantic comparison from the exact selected ref so ambiguity is visible and ordinary/reset planning remain read-only.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

## Project Constraints

- Identifier strategy: applicable — preserve Spec/Task identities and exact Git ref identity; any new recovery record follows the existing Run identity boundary. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential or HTTP API change is proposed; use existing read boundaries only. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0143 requires read-only ordinary/reset planning; ADR-0062 preserves operational history during identity-reset planning. ADR-0146 publication is outside this Spec. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `internal/docscontract/publicdocs_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Stable version parser | `internal/releaseplan/version.go` | Parse numeric identity and retain spelling without weakening the stable-only grammar. |
| Git tag discovery | `internal/cli/releaseplan_git_source.go` | Inventory reachable refs across both spellings and preserve exact commits. |
| Base selection and proposal | `internal/releaseplan/build.go, proposal.go, model.go, errors.go` | Select the highest unique semantic base or require explicit --from; retain proposal spelling. |
| Reset identity | `internal/releaseplan/reset.go and Git source` | Preserve every exact local/remote ref, Release match and digest input. |
| Public contract | `CLI releaseplan command/tests, docs release runbook/usage and owned roundfix skill` | Expose typed ambiguity through existing error/refusal conventions without new publication authority. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Parsing and selection

The current parser requires `v`, Git discovery lists `v*`, and proposal
formatting always adds `v`. Change these three seams together. Accept only
`vMAJOR.MINOR.PATCH` and `MAJOR.MINOR.PATCH`, retaining current numeric validity
and rejection of prerelease/build metadata. Semantic comparison uses the
numeric tuple; the selected base also retains exact name, ref and commit.
Do not infer spelling later from a normalized Version.String value.

Discover both tag families, filter through the strict parser and existing
reachability boundary, and select the highest semantic version. Two spellings
of that highest version are ambiguous even if both resolve to one commit.
Return a typed exit-2 refusal listing exact names/commits and the existing
`--from` recovery. Explicit `--from` selects the actual ref and preserves
ancestry/range checks; it does not bypass them or silently normalize the name.
Lower-version aliases do not prevent selection of a unique highest version.

### Proposal and reset behavior

The next proposed version uses the selected base's prefix convention. Existing
classification and approval fields keep their meaning. Preserve each local
and remote tag name through reset inventory; Releases match exact tag names,
not semantic aliases. Include exact spelling/ref/commit in the existing reset
digest so a changed alias cannot retain an old approval identity.

Keep ordinary and reset planning read-only under ADR-0143. ADR-0062 retains
operational history. The existing success schema 0.0.1 fields can represent
these values; do not introduce a new readiness state, release command or
configuration option. ADR-0146 publication and `.github/workflows/release.yml`
are outside this Spec. No tag is created, deleted, renamed or pushed during
planning, including an ambiguous or explicitly selected case.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Separate the numeric stable version from exact tag spelling, ref and commit. Reset inventory retains semantic aliases as distinct refs and includes their identity in its digest.

### API Contracts

Reuse release plan, explicit --from, success schema 0.0.1 and typed exit-2 preflight refusal. Ambiguity reports both exact names/commits; ordinary/reset planning never creates tags or publishes a Release.

## Coverage Map

- PRD Goal 1 → Stable version parser, Git tag discovery.
- PRD Goal 2 → Base selection and proposal, Public contract.
- PRD Goal 3 → Reset identity.
- Core Feature 1 → Stable version parser.
- Core Feature 2 → Git tag discovery and Base selection and proposal.
- Core Feature 3 → Base selection and proposal.
- Core Feature 4 → Reset identity.
- Core Feature 5 → Public contract.

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

1. Parser accepts both strict stable spellings and rejects existing malformed/prerelease/build-metadata cases.
2. Real temporary Git repositories prove highest reachable selection across spellings, nonreachable exclusion and unchanged explicit --from ancestry validation.
3. Highest semantic aliases refuse with exact refs/commits even on one commit; explicit selection succeeds and lower aliases do not cause a false refusal.
4. Ordinary and reset proposals retain selected spelling; reset inventory/digest preserves local/remote aliases and exact Release matches.
5. Compare refs, worktree/index, remote operations and stored operational history before/after every planning path; all remain unchanged.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Strict parse and exact selected-tag model (depends on: none).
2. Cross-spelling discovery, selection and typed ambiguity (depends on: 1).
3. Proposal convention and reset identity preservation (depends on: 1, 2).
4. Public output/guidance and real Git terminal QA (depends on: 1, 2, 3).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Normalizing a tag too early can collapse distinct refs or change a reset digest. Maintain compatibility for existing v-prefixed consumers and obtain the bounded governed test/skill grant before implementation.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Preserve configured reviewer selection; this Spec introduces no separate reviewer override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.
