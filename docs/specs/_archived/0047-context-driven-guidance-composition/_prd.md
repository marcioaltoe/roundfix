---
spec: 0047-context-driven-guidance-composition
status: archived
created: 2026-07-24
surfaces: [cli, infra, docs]
archived: "2026-07-25"
source_slug: 0047-context-driven-guidance-composition
---


# Context-Driven Baseline guidance composition

The public Baseline Command can generate a compact root index and modular agent
guides, but the accepted Fluxus greenfield result still required maintainers to
redistribute rules from a catch-all carrier, infer instruction precedence, and
complete the ADR and Findings guidance manually. This feature makes the
generated guidance self-contained and gives every accepted rule one semantic
owner while preserving the retention, confirmation, and rollback guarantees
delivered by Specs 0045 and 0046. The
[greenfield acceptance finding](../../findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md)
is the acceptance source.

## Goals

- Generated root instructions and modular guides expose one explicit
  Instruction Hierarchy without duplicating policy.
- Every accepted project rule resolves to a managed semantic guide when one
  exists, leaving Repository-Specific Normative Rules only for real residuals.
- Greenfield adoption and existing-repository update produce the same portable,
  self-contained guidance without manual redistribution.
- A maintainer whose repository does not match a built-in Profile can create a
  reviewed repository-owned adaptation from the interactive Baseline flow
  instead of ending at an unactionable capability list.
- ADR and Findings guidance retain their complete lifecycle and template
  Operational Contracts without modifying upstream-managed skills.
- Apply, formatting, repository Verification, audit, and empty reapply preserve
  one stable generated result.

## User Stories

1. As a repository maintainer, I want project rules placed in the guide that
   owns their meaning, so that contributors do not search a catch-all document
   for normal repository policy.
2. As a maintainer updating an existing Baseline, I want Roundfix to propose a
   complete rule redistribution with retention evidence, so that I can migrate
   without losing or silently weakening instructions.
3. As an Agent, I want an explicit Instruction Hierarchy, so that I know which
   guides apply and that narrower guidance cannot weaken universal policy.
4. As a repository maintainer, I want each generated guide to contain the
   operative rule instead of depending on an ADR, so that greenfield
   repositories remain usable before they accumulate decision history.
5. As a maintainer creating ADRs or Findings, I want complete lifecycle and
   template contracts, so that new documents are consistent and
   machine-readable.
6. As a maintainer with a genuinely unmodeled project rule, I want a residual
   carrier only when it has content, so that exceptional policy remains safe
   without becoming the default organization model.
7. As a maintainer whose repository implements only part of an opinionated
   built-in Profile, I want Roundfix to propose a repository-owned Profile
   adaptation, so that required capabilities remain truthful without forcing
   me to hand-edit a Profile before setup can continue.

## Core Features

1. The generated root index presents the Instruction Hierarchy in this order:
   universal instructions; context and documentation; Spec workflow; autonomous
   work when enabled; stack and surface guides; optional knowledge sources.
2. A narrower guide can add constraints for its concern but cannot weaken a
   universal Normative Clause or confirmed project decision.
3. The Baseline assigns each managed rule to one semantic owner. Universal
   execution and safety rules belong to the core agent guide; domain language
   and identifier policy belong to the domain guide; backend and HTTP rules
   belong to the backend guide; frontend rules belong to the frontend guide;
   Spec and Task rules belong to the Spec workflow guides; other selected
   modules retain their dedicated guides.
4. Generated guidance contains every operative rule, condition, exception, and
   next action needed to follow it. ADRs can preserve decision history but
   cannot be required to recover the rule's meaning.
5. Greenfield adoption creates only the semantic guides selected by the
   Baseline Profile and includes no repository-specific carrier or root pointer
   when no residual rule exists.
6. Update and Baseline Readoption classify every existing repository-specific
   rule, propose its semantic owner when representable, and retain ambiguous or
   unmodeled content as Repository-Specific Normative Rules. No source rule can
   disappear from the Upgrade Retention Contract.
7. The Change Plan shows each source rule, proposed semantic owner, residual
   disposition, affected root pointer, and before/after identity before the
   maintainer can confirm migration.
8. When redistribution leaves no residual rule, the confirmed Change Plan
   removes the empty repository-specific carrier and its root pointer. When
   residual rules remain, the carrier is created or retained with only those
   accepted rules.
9. The domain guide requires Agents to read the repository's domain context and
   every active ADR relevant to the work before naming concepts or changing
   behavior.
10. The ADR Operational Contract augments the externally owned ADR body format
    with repository-owned lifecycle metadata: `proposed`, `accepted`,
    `rejected`, `deprecated`, or `superseded`; RFC 3339 UTC creation and update
    timestamps; nullable deprecation timestamp; and nullable superseding ADR
    identity. Only `accepted` is active.
11. A legacy ADR without lifecycle metadata remains active unless its body
    explicitly marks it inactive. Baseline does not rewrite an existing ADR
    solely to adopt the lifecycle metadata.
12. The Findings Operational Contract includes one copyable template with
    `pending`, `partial`, `deferred`, and `done` lifecycle states; session
    context; symptom and evidence; proven or unknown root cause; action or Spec
    routing; optional retained practices; and dated append-only addenda.
13. Repository-owned lifecycle overlays compose with upstream-managed skills.
    Baseline never changes an upstream-managed skill or replaces its document
    body contract.
14. Every affected maintained Baseline Profile produces deterministic,
    Formatter-Stable Output and retains the existing digest confirmation,
    recoverable apply, and empty-reapply contracts.
15. The generated guide set keeps dedicated semantic owners for universal
    instructions, domain language, documentation layout, Spec routing, Task
    tracking, autonomous work, backend, frontend, monorepo, stack, skill
    dispatch, and optional Secondbrain guidance whenever their modules apply.
16. Baseline does not generate a generic repository guide. During update, rules
    found in a legacy generic repository guide receive the same semantic-owner,
    residual, or explicit-rejection accounting as other repository rules.
17. Roundfix user documentation and the thin setup skill explain the
    Instruction Hierarchy, greenfield composition, update redistribution,
    residual-carrier behavior, ADR lifecycle, and Findings template without
    introducing another execution engine.
18. When profile alignment finds missing required profile-specific
    capabilities, the human flow offers a digest-bound repository-owned Profile
    adaptation, re-audits it, and includes its file in the final Change Plan.
    It never converts a required capability into a waiver.
19. A Profile adaptation can remove incompatible selected modules and
    profile-specific capabilities only after explicit review and catalog
    validation. Universal required capabilities remain non-removable and name
    their exact remediation operations.
20. Automation can provide the same strict repository-owned Profile draft as
    explicit planning input. The draft becomes a repository file only through
    the confirmed portable Plan; no interactive or automation path writes it
    before Plan Digest approval.

## User Experience

During greenfield adoption, the maintainer reviews a compact root index and the
complete modular guide set in the Change Plan. No empty extension choice or
placeholder carrier appears.

During update or Baseline Readoption, Roundfix shows one consolidated
redistribution proposal. Rules with known semantic owners are grouped by their
destination; ambiguous rules remain visible as residuals. The maintainer can
revise classifications before confirming the new Plan Digest.

When a built-in Profile is broader than the repository, Roundfix shows the
blocking capabilities individually and offers to construct a repository-owned
adaptation. The maintainer reviews every removed module and capability, then
the audit runs again before instruction classification or the final Change
Plan. Missing universal skills remain blocking and receive exact restoration
commands because a custom Profile cannot weaken them.

After apply, the root index explains the Instruction Hierarchy and each guide
stands on its own. ADRs can explain why a rule exists, but removing or never
creating an ADR does not make the guide incomplete.

## Project Constraints

- Identifier strategy: not applicable. This feature creates no project-owned
  application identifiers; Spec 0048 defines the Baseline decision.
- Authentication and HTTP: not applicable. This feature neither installs an
  authentication provider nor changes application routes.
- Active ADR obligations: applicable. The implementation must preserve the
  retention, profile ownership, supervised analysis, portable Plan, and
  recoverable apply contracts in ADRs 0058, 0067, 0069, 0070, 0071, and 0073.
- Tooling authority: applicable. This Spec grants no authorization to create,
  edit, rename, move, or delete tooling configuration, scripts, ignore files,
  plugin declarations, or version pins.

## Non-Goals / Out of Scope

- Collecting identifier-strategy or Better Auth decisions; Spec 0048 owns those
  decisions and their propagation.
- Combining the modular guides into one large instruction file.
- Renaming the accepted transparent guide set without a separate migration
  decision.
- Modifying the externally owned domain-modeling skill or any other
  upstream-managed skill.
- Rewriting existing ADRs only to add lifecycle metadata.
- Automatically rewriting arbitrary nested instruction carriers that Spec 0046
  leaves under repository ownership.
- Copying Fluxus names, product rules, or branding into portable Baseline
  assets.
- Changing the public Baseline plan, apply, result, recovery, or exit-code
  contracts delivered by Spec 0046.

## Success Metrics

- 100% of accepted project rules have exactly one managed semantic owner,
  residual disposition, or explicit reasoned rejection.
- 0 generated operative rules require an ADR to exist before an Agent can
  follow them.
- 0 repository-specific carriers or root pointers are generated when the
  residual rule count is zero.
- 0 generic repository guides are generated after greenfield or confirmed
  update.
- 100% of generated root indexes present the confirmed Instruction Hierarchy,
  and 0 narrower guides weaken a universal Normative Clause or project
  decision.
- 100% of ADR and Findings fixtures contain the complete confirmed Operational
  Contract while 0 upstream-managed skill bytes change.
- 100% of public guidance examples and thin-skill instructions describe the
  same behavior accepted by the Baseline Command.
- 100% of blocking profile-specific divergence journeys either select another
  Profile or produce one catalog-valid, explicitly reviewed repository-owned
  adaptation in the confirmed Change Plan; 0 required capabilities become
  waivers.
- Greenfield and update fixtures for every affected maintained profile complete
  apply, formatting, repository Verification, audit, and empty reapply with
  zero managed-file delta.
- Separately authorized Fluxus greenfield and update journeys require zero
  manual rule redistribution and produce no empty repository-specific carrier.
- The full Roundfix Verification passes after the Spec-local QA evidence and
  live Fluxus journeys are recorded.

## Decisions

- The guidance remains modular; clearer hierarchy replaces consolidation.
- Existing transparent guide names remain the acceptance target.
- A managed semantic owner takes precedence over the residual
  repository-specific carrier.
- The residual carrier exists only for non-empty accepted rules that no typed
  decision or managed guide can represent.
- ADR lifecycle status is `proposed`, `accepted`, `rejected`, `deprecated`, or
  `superseded`; only `accepted` is active.
- ADR lifecycle metadata augments rather than replaces the externally owned
  ADR body format.
- Required profile divergence is resolved through a reviewed repository-owned
  Profile adaptation, never through a waiver on the built-in Profile.
- Fluxus supplies the real acceptance journey, but portable assets contain no
  Fluxus-specific names or policy.

## Open Questions

None.
