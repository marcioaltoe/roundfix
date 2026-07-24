---
spec: 0040-mandatory-context-driven-adrs
status: archived
created: 2026-07-17
archived: 2026-07-24
qa_override: true
archive_reason: superseded
superseded_at: 2026-07-24
superseded_by:
  - 0047-context-driven-guidance-composition
  - 0048-context-driven-project-decisions-and-spec-constraints
surfaces: [docs, tooling]
---

# Mandatory Context-Driven ADRs

> Superseded without implementation by
> [Spec 0047](../../0047-context-driven-guidance-composition/_prd.md) and
> [Spec 0048](../../0048-context-driven-project-decisions-and-spec-constraints/_prd.md).
> The replacement contracts keep operative guidance independent of ADR
> existence and require explicit confirmation of identifier policy.

The Context-Driven Baseline currently generates instructions and supporting
guides, but it does not install the architectural decisions that make its
domain-documentation contract portable across repositories. Projects can
therefore select the same setup profile while disagreeing about where
ubiquitous language lives, which language ADRs use, or what format a
project-controlled Internal Identifier must have.

This feature makes two project-agnostic Baseline ADRs mandatory in every
`setup-context-driven` profile. The first defines colocated ubiquitous-language
documents and requires ADRs to be written in English. The second requires
UUIDv7 for Internal Identifiers while preserving external identifiers, natural
keys, business codes, and protocol-defined identities. Audit proves the
Baseline ADRs are present and current; apply installs or refreshes them and can
transactionally renumber existing root ADRs to reserve their canonical
identities.

## Goals

- Install the same two Baseline ADRs in every supported Context-Driven profile.
- Reserve these canonical root paths:
  - `docs/adr/0001-ubiquitous-language-in-colocated-context-md.md`;
  - `docs/adr/0002-uuidv7-required-for-internal-identifiers.md`.
- Keep the Baseline ADR bodies project-agnostic, in English, and versioned as
  setup-owned managed content.
- Preserve repository-authored notes outside the managed ADR blocks.
- Safely shift existing root ADRs by two positions when the reserved paths are
  not already exact Baseline ADRs.
- Update repository-owned references to every renamed ADR as part of the same
  all-or-nothing apply operation.
- Keep context-local ADR namespaces unchanged in multi-context repositories.
- Make missing, stale, ambiguous, or unsafe Baseline ADR state a deterministic
  blocking audit result.
- Preserve the separation between repository setup compliance and Doctor Skill
  Readiness.

## User Stories

1. As a repository maintainer, I want every Context-Driven repository to carry
   the same foundational domain-documentation decisions, so agents do not have
   to infer those rules from project-specific prose.
2. As an agent authoring or changing an ADR, I want the repository to state
   unambiguously that ADR filenames, titles, and prose are in English.
3. As a domain expert, I want canonical glossary terms to remain in the
   language the business actually uses, even when code identifiers use English
   aliases.
4. As a developer creating an Internal Identifier, I want one portable UUIDv7
   invariant across domain, persistence, adapters, seeds, and fixtures.
5. As a maintainer adopting the setup in a mature repository, I want existing
   ADRs and their references preserved while the two reserved identities are
   introduced.
6. As a maintainer rerunning setup, I want exact unmarked Baseline ADRs adopted
   automatically and managed ADR blocks refreshed without another decision
   prompt.
7. As a maintainer with context-local ADRs, I want the migration limited to the
   root `docs/adr/` namespace so bounded-context history is not renumbered.

## Core Features

1. **Mandatory Baseline ADR artifacts.** The `context-workflow` module owns two
   first-class ADR artifacts, and every supported profile already activates
   that module. Neither ADR is optional or controlled by a new user decision.
2. **Portable ubiquitous-language decision.** ADR `0001` defines a root
   `CONTEXT.md` for a single-context repository, a root `CONTEXT-MAP.md` plus
   colocated context `CONTEXT.md` files for a multi-context repository, domain
   terms in the language used by domain experts, English code aliases when the
   vocabularies differ, and English for every new or modified ADR.
3. **Portable identifier decision.** ADR `0002` requires RFC 9562 UUIDv7 for
   technical entity or resource identities generated and controlled by the
   project. It excludes external identifiers, natural keys, business codes,
   protocol-defined identities, and values that are not persistent identity.
4. **Managed ownership.** Each Baseline ADR contains a stable, versioned
   setup-owned block. Apply may create or refresh that block while preserving
   project-specific notes outside it.
5. **Exact automatic adoption.** An unmarked canonical file is adopted without
   a question only when its required body is byte-equal to the bundled
   template. Divergent content is historical repository content and is never
   overwritten.
6. **Root ADR migration.** On first adoption, existing historical ADRs in
   `docs/adr/` move from `NNNN-slug.md` to `(NNNN+2)-slug.md`, preserving order,
   gaps, and slugs. Baseline ADRs already proven exact are excluded. ADRs below
   another context's `docs/adr/` are outside the migration.
7. **Reference migration.** Exact ADR filenames, repository-relative Markdown
   links, and explicit `ADR-NNNN` / `ADR NNNN` references in versioned,
   repository-owned text files move to the new identities. Ambiguous or
   protected references block apply rather than remaining stale.
8. **Transactional preview and apply.** Preview lists every creation, refresh,
   rename, reference edit, and manual language-review notice. Apply validates
   the complete target tree before changing any file and restores original
   bytes and modes if any operation fails.
9. **Blocking audit.** Missing files, missing or stale managed blocks, duplicate
   root ADR numbers, malformed ADR filenames, unsafe target numbers,
   collisions, partial Baseline ADR adoption, and unresolved references are
   stable blocking findings.
10. **Historical language review.** Setup never translates existing ADRs or
    guesses their language. A first migration reports the preserved historical
    ADR inventory as a non-blocking manual English-language review.
11. **Instruction alignment.** Generated root guidance and the domain guide
    point agents to the Baseline ADRs. The setup skill, user guide, portable
    assets, canonical glossary, and embedded skill copy describe the same
    contract.

## User Experience

Audit remains read-only. A repository without the Baseline ADRs sees both
planned files plus the complete deterministic migration preview. If no root
ADRs exist, the plan contains only the two creations and the normal manifest
refresh. If historical ADRs exist, the preview names every old and new path,
every reference-bearing file that will change, and the historical files that
need a manual English-language review.

The existing setup confirmation authorizes the complete plan; no additional
decision question is introduced. Apply either lands the two managed ADRs,
renamed history, updated references, and manifest together or leaves the
repository byte-for-byte unchanged. A second apply is a no-op and a subsequent
audit exits successfully.

## Non-Goals / Out of Scope

- Translating, rewriting, or normalizing historical repository-authored ADRs.
- Renumbering ADRs outside the root `docs/adr/` directory.
- Detecting the natural language of arbitrary Markdown through heuristics or an
  external model.
- Scanning source code or data stores to prove a project already conforms to
  the UUIDv7 decision.
- Generating UUIDv7 adapters, changing database schemas, or migrating existing
  identifiers; affected projects plan those changes separately.
- Treating every string, slug, business code, natural key, external identifier,
  or protocol-defined identity as an Internal Identifier.
- Making `roundfix doctor` parse or validate repository ADRs. Doctor continues
  to prove the installed `setup-context-driven` skill version through the
  Repository Skill Set check in Spec 0036.
- Reordering or compacting intentional gaps in an ADR sequence.
- Modifying upstream-managed Agent Skills while rewriting repository
  references.

## Success Metrics

- Every bundled profile resolves the same two mandatory ADR artifacts at the
  canonical root paths.
- A clean new repository receives both ADRs, in English, with current managed
  blocks and no decision prompt beyond the existing apply confirmation.
- A fixture with existing ADRs `0001`, `0004`, and `0015` moves them to `0003`,
  `0006`, and `0017`, updates all unambiguous tracked references, and preserves
  file bytes outside changed numbers and managed blocks.
- Exact unmarked canonical ADR bodies are adopted automatically; divergent
  canonical files are preserved through migration rather than overwritten.
- Any collision, duplicate number, malformed filename, overflow, protected
  reference, or injected write failure leaves every file and mode unchanged.
- Context-local ADRs remain byte-for-byte unchanged.
- Applying twice produces no second diff, and a clean audit returns exit code
  `0`.
- The setup Python suite covers new, exact-adoption, divergent, mixed,
  collision, ambiguous-reference, rollback, idempotency, and multi-context
  cases; the repository verification gate passes.

## Decisions

- Create a dedicated Spec rather than expanding Doctor Skill Readiness or the
  Run-lifecycle Specs 0037–0039.
- Make both ADRs unconditional members of the Context-Driven Baseline.
- Reserve `0001` and `0002` only in the root ADR namespace.
- Use setup-owned managed blocks and preserve project-authored extensions.
- Adopt only exact unmarked canonical content automatically.
- Shift historical root ADR identities by `+2` and update references as one
  transactional operation.
- Require English for new or modified ADR filenames, titles, and prose while
  allowing ubiquitous-language terms to use the domain experts' language.
- Apply UUIDv7 only to Internal Identifiers controlled by the project.
- Keep setup audit/apply responsible for ADR compliance; Doctor validates the
  installed skill, not project documentation.

## Open Questions

None.
