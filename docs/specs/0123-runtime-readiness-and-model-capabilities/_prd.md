---
spec: 0123-runtime-readiness-and-model-capabilities
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# A configured Agent Selection is proved before it receives work

The maintainer requested Claude Fable 5.1 and GPT-6 Astra support. Public model identifiers do not establish adapter capability, full-access configuration can fail only after a work Session starts, and eager proof of an unused fallback can block a working preferred selection. Fresh go vet also reports 31 copies of the ACPX runner mutex. Readiness must state what was actually proved and preserve one owner of mutable runtime state.

This Spec is **in authoring**. Its protected scope and open decisions are
pending approval. The accompanying authorization is a proposal with no grant;
there is no TechSpec, Task Graph, or authority to start implementation.

## Project Constraints

- Identifier strategy: applicable — model IDs are runtime-advertised opaque values; keep API/provider identifiers distinct from ACP Agent Selection IDs and preserve existing Run/Session identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no new authentication or HTTP layer; use current runtime credentials and do not transmit secrets in probes or research. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0107 requires every configured Work Category to be accounted for; ADR-0147 preserves honest advertised capability evidence and the adapter refusal. The current eager Fallback Chain contract is an explicit proposed decision to revisit, not permission to omit readiness evidence. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0049 applies: each present Agent Selection Profile replaces the lower-precedence profile atomically and preserves its explicit fallback chain.
  ADR-0140 applies: prove the exact advertised runtime/model/effort tuple through the installed adapter. Any narrower change to fallback validation timing remains a proposed revision, not an accepted exception.
  ADR-0093 applies to authored capability claims: the consistency checker follows explicit citations and must not infer approval or model support from missing evidence.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Offer the requested models only through a proved runtime identifier and supported controls.
- Report access-policy incompatibility before starting Agent work.
- Keep optional fallback failure distinguishable from preferred-selection readiness.
- Eliminate copied-lock diagnostics by correcting runtime-state ownership.

## Core Features

1. Discover and prove adapter-advertised model/effort choices for GPT-6 Astra and Claude Fable 5.1; preserve their source, timestamp, effective values, and unsupported outcomes.
2. Keep public API IDs gpt-6-astra and claude-fable-5-1, provider-qualified IDs, and ACP identifiers distinct. Do not automatically switch effective profiles when the advisory catalog changes.
3. Include effective access-policy capability in readiness and refuse an unsupported requested policy before work, with the failed predicate and its remedy named.
4. Proposed: prove the preferred selection before Run creation, and prove a fallback immediately before activation while retaining an exhaustive explicit profiles-validation command. No fallback may begin work unproved; this revises the eager-proof policy only if approved.
5. Remove runner mutex copying through one coherent state-ownership design and verify that session lifecycle and cancellation behavior still work. Spec 0124 owns the later Verification-configuration change that makes analyzer coverage continuous.

## Non-Goals / Out of Scope

- No automatic repository profile replacement, invented ACP alias, estimated price without a source, or paid probe without a spending limit.
- No runtime/adapter package upgrade unless a named capability gap and separate bounded authorization require one.
- No weakening of selection evidence, credential policy, or lock diagnostics.

## Acceptance evidence

Outside-evidence row: Use the pre-existing 31-diagnostic go-vet failure and captured runtime capability evidence. A public-API model that is absent from the selected ACP catalog must not become a usable choice merely because its name exists in this Spec.

The later Task Graph and QA must prove each Core Feature with observed
positive and negative cases. Source inspection and proposed checks do not
constitute implementation or terminal QA evidence.

## Open Questions

- Approve the proposed fallback admission policy or retain eager proof of the complete chain.
- Set the allowed real-adapter probe budget and access scope; until then, live paid validation remains blocked.
- Decide what user-visible state an unavailable model retains while its API exists but the configured adapter does not advertise it.

Until answered, all proposed limits and protected mutations remain unapproved.

## Source ownership

The maintainer selected this intent for implementation. Ordinary sources now
have one primary owner and one copy under that owner's `references/` directory.
Active Rollups remain as shared archive-license roots; their dated addenda map
every remaining family to its consuming Spec. Adoption is not execution approval.

- [2026-09-08-support-claude-fable-5-1-and-gpt-6-astra.md](references/2026-09-08-support-claude-fable-5-1-and-gpt-6-astra.md)
- [2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md](references/2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md)
- [2026-08-06-rollup-agent-selection-and-execution-environments.md](../../history/findings/2026-08-06-rollup-agent-selection-and-execution-environments.md)
- [model-selection.md](../../references/model-selection.md)

## Research basis

Secondbrain `inbox/roundfix/_triaged/2026-08-30-go-vet-acha-31-copylocks-que-o-gate-do-repo-nao-ve.md` and the Agent Selection Rollup provide local failure history. Exa read [Go sync documentation](https://pkg.go.dev/sync), supporting correction of copied synchronization state. The session's primary model research read [GPT-6 Astra](https://developers.openai.com/api/docs/models/gpt-6-astra), [Anthropic model overview](https://docs.anthropic.com/en/docs/about-claude/models/overview), and [Anthropic API release notes](https://docs.anthropic.com/en/release-notes/api). They establish public identifiers, not ACP availability or a successful local selection.

The local Secondbrain index was read and its query workflow used before
authoring. The sources above affected the stated requirements and limits;
they do not approve this Spec or substitute for live capability/behavior proof.

## Authoring checkpoint

Spec 0119 defines execution/probe authority and 0121 supplies skill ownership. Correct copied-lock defects before Spec 0124 adds continuous go-vet gating; no green baseline is fabricated by changing that gate first.

Record the maintainer's bounded decision in [_authorization.md](_authorization.md),
then commit that approval record separately before its consuming tooling
changes. Only after that checkpoint may the TechSpec settle the design and a
Task Graph authorize execution. PRD-stage checks will report their actual
scope and any pending authorization findings; a green partial check is not
implementation readiness.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
