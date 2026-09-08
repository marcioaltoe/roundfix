---
spec: 0121-baseline-decisions-and-complete-regeneration
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# A Baseline change preserves the decision it edits

Changing an HTTP Contract mode can discard its recorded exceptions and source. A second, inert default disagrees with the value the interactive command uses. Greenfield adoption can demand a classification its own interface never collects, and skill regeneration lacks the output ownership that authorization auditing expects. The maintainer needs one explainable transition that retains valid decisions and reaches a verifiable result.

This Spec is **in authoring**. Its protected scope and open decisions are
pending approval. The accompanying authorization is a proposal with no grant;
there is no TechSpec, Task Graph, or authority to start implementation.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity or identifier strategy is proposed; preserve existing Baseline and decision identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the proposal changes how the CLI edits a repository-owned HTTP Contract, preserving its confirmed mode, typed exceptions, and source; it introduces no application endpoint or authentication policy. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 keeps regenerated pins as sanctioned fallout. ADR-0149 gives the tree ownership of regeneration outputs. ADR-0130 keeps the audit's governed path set honest against authorization history. Retain semantic preservation rather than treating stale source as disposable. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/contract-v1.json`, `internal/cli/baseline_human_test.go`, `internal/baseline/derived_ownership_test.go`, `skills/_ownership.yml`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `internal/baseline/assets/decisions.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/go-cli-tui.json`, `internal/baseline/assets/profiles/rust-cli.json`, `internal/baseline/assets/templates/index.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`, `internal/baseline/plan_test.go`, `.agents/skills/setup-context-driven/SKILL.md`, `skills/setup-context-driven/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A mode-only decision change retains every untouched exception and source field.
- The displayed HTTP default has one effective owner and no inert conflicting declaration.
- Every supported adoption route either completes with its required evidence or refuses before requesting an unreachable step.
- The audit resolves skill regeneration outputs from the command ownership record.

## Core Features

1. Change the selected HTTP mode without silently replacing the rest of the typed decision. Deliberately editing exceptions remains distinguishable from retaining them.
2. Remove or reconcile the unused HTTP default and its unused diagnostic contract; keep catalog-level defaults and explicit repository decisions distinguishable.
3. Proposed: when Greenfield encounters retained managed source requiring classification, refuse early with the Preservation route named. Do not discard source or invent approval; a different design requires a recorded decision.
4. Declare skill regeneration ownership and prove that sanctioned generation covers the declared outputs while manual derived edits remain refused.
5. Update shipped command guidance and regenerate only the approved derivation closure.

6. Represent the local incremental command separately from the complete verification.gate decision and expose both through Profile output, planning, generated guidance and the Setup Manifest. Existing single-gate manifests need a reviewed migration; never invent a command or treat absence as compliance. Replay the Fiscus capture in an isolated repository without changing Fiscus tooling.
7. Reconcile obsolete external skill-lock entries only against positively established absence at an explicitly selected immutable upstream revision and a confirmed mutation plan. An unreachable source, a merely unneeded skill or local absence cannot authorize removal. Required-but-removed skills block; preserve installed files and unknown lock fields, and keep Doctor offline/read-only.

## Non-Goals / Out of Scope

- No new API mode, authentication provider, framework, or repository-wide policy default.
- No automatic adoption of unreviewed source or rewriting of archived observations.
- No unrelated skill update, dependency upgrade, or new generator framework.

## Acceptance evidence

Outside-evidence row: Retain a multi-exception HTTP decision taken from the pre-existing archived Finding while changing only its mode; prove that removal of an exception is detected. That outside fixture must keep its original provenance.

The later Task Graph and QA must prove each Core Feature with observed
positive and negative cases. Source inspection and proposed checks do not
constitute implementation or terminal QA evidence.

## Open Questions

- Approve retaining untouched HTTP fields and the early Greenfield refusal, or specify the alternative adoption behavior.
- Approve removal of the inert Profile default rather than introducing a new per-Profile override.
- Approve the exact protected paths and the new skill ownership declaration below.

Until answered, all proposed limits and protected mutations remain unapproved.

## Source ownership

The maintainer selected this intent for implementation. Ordinary sources now
have one primary owner and one copy under that owner's `references/` directory.
Active Rollups remain as shared archive-license roots; their dated addenda map
every remaining family to its consuming Spec. Adoption is not execution approval.

- [2026-08-06-rollup-baseline-and-derived-tooling.md](../../findings/2026-08-06-rollup-baseline-and-derived-tooling.md)
- [2026-09-08-skill-regeneration-declares-its-owned-outputs.md](references/2026-09-08-skill-regeneration-declares-its-owned-outputs.md)
- [2026-08-07-changing-the-http-contract-discards-its-exceptions.md](../../history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md)
- [2026-08-07-two-http-contract-defaults-and-only-one-is-read.md](../../history/findings/2026-08-07-two-http-contract-defaults-and-only-one-is-read.md)
- [2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md](../../history/findings/2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md)

## Research basis

Secondbrain `inbox/roundfix/_triaged/2026-08-30-skills-nao-declara-ownership-e-forca-todo-grant-a-enumerar-espelhos.md` established the repeated ownership gap; `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` informed negative controls for retention and regeneration. Exa consultation read [RFC 6902](https://www.rfc-editor.org/rfc/rfc6902), which makes document mutation explicit; it supports testing a scoped field change, not adopting JSON Patch as a new dependency. Local observations establish the defects; the external source does not validate a proposed Roundfix implementation.

The local Secondbrain index was read and its query workflow used before
authoring. The sources above affected the stated requirements and limits;
they do not approve this Spec or substitute for live capability/behavior proof.

## Authoring checkpoint

Spec 0119 owns the future approval-reader convention. This Spec supplies the skill ownership prerequisite used by later skill-changing Specs; no implementation starts from the proposal alone.

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
