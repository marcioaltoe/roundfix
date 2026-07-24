---
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: active
created: 2026-07-24
surfaces: [cli, infra, docs]
---

# Context-Driven project decisions and Spec constraints

The composed Context-Driven guidance still cannot reproduce the accepted
project contract unless maintainers manually add an Internal Identifier
strategy, a Better Auth route exception, and a broad tooling-authority rule.
New Specs can then omit those constraints and authorize incompatible work by
silence. This feature collects the project decisions through the public
Baseline Command, persists them in the Setup Manifest, renders self-contained
guidance, and makes every new Spec account for the applicable Project
Constraints. The
[greenfield acceptance finding](../../findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md)
is the acceptance source.

## Goals

- Human and automation callers resolve identifier and Better Auth decisions
  explicitly without silent policy inference.
- Confirmed project decisions render as self-contained domain and backend
  guidance and survive later Baseline updates.
- Every new PRD and TechSpec states how the repository's Project Constraints
  apply before implementation Tasks can begin.
- Tooling configuration remains maintainer-owned through one non-optional
  Normative Clause that no setup preference can weaken.
- Equivalent interactive and non-interactive inputs produce the same
  digest-bound plan and final repository guidance.

## User Stories

1. As a repository maintainer, I want UUID version 7 suggested for new
   project-owned Internal Identifiers, so that I can confirm or change one
   consistent identifier strategy.
2. As a maintainer using Better Auth, I want its provider-owned HTTP routes
   represented as an explicit typed exception, so that authentication protocol
   behavior is preserved without depending on an ADR.
3. As an automation author, I want every unresolved project decision exposed as
   stable machine-readable input, so that automation never accepts a suggested
   default implicitly.
4. As a Spec author, I want each PRD and TechSpec to account for applicable
   Project Constraints, so that Tasks begin with the same repository contract
   as implementation Agents.
5. As a repository maintainer, I want all tooling configuration mutations to
   require express authorization, so that an Agent cannot change quality or
   build policy as an incidental implementation step.
6. As an Agent, I want complete operative constraints in the semantic guides
   and Spec artifacts, so that absent or deprecated ADRs do not create policy
   gaps.

## Core Features

1. The Standard TypeScript Monorepo Profile includes a typed project decision
   for the identifier strategy applied to new project-owned Internal
   Identifiers.
2. When no compatible prior identifier decision exists, human setup proposes
   UUID version 7 and asks the maintainer to keep or change it. The suggestion
   is visible and cannot become authorization without an answer.
3. Non-interactive setup requires an explicit identifier-strategy value when
   the decision is unresolved. Missing input produces a structured next action
   and no repository mutation.
4. The confirmed identifier strategy applies to new Internal Identifiers.
   External provider identifiers, protocol identifiers, natural keys, and
   business codes retain their source contracts.
5. When Better Auth is a selected Repository Capability, setup proposes a
   typed HTTP Contract Decision exception owned by Better Auth for `GET` and
   `POST` under `/api/auth/*`, covering session, OAuth redirect, callback, and
   related provider protocol semantics.
6. Human setup asks the maintainer to keep or change the Better Auth proposal.
   Automation supplies the explicit exception when it is unresolved. Neither
   repository evidence nor framework detection silently authorizes it.
7. Compatible persisted decisions are reused during update and remain visible
   as keep-or-change choices. Missing, incompatible, or incomplete decisions
   block planning until the caller resolves them.
8. The Setup Manifest, Decision Plan, Change Plan, and machine-readable result
   retain the confirmed values and their deterministic identities.
9. The confirmed identifier rule renders completely in the domain guide. The
   confirmed HTTP mode and Better Auth exception render completely in the
   backend guide. Neither operative rule requires an ADR.
10. One universal Normative Clause requires express maintainer authorization
    before creating, editing, renaming, moving, or deleting any linter,
    formatter, typechecker, test-runner, architecture-checker, build-tool,
    package-manager, code-generator, or other repository-tooling
    configuration, script, ignore file, plugin declaration, or version pin.
11. The tooling-authority Normative Clause is not a setup preference, has no
    enable or disable prompt, and cannot be weakened by a profile, narrower
    guide, or project decision.
12. Every new PRD and TechSpec contains a Project Constraints section. It
    identifies each confirmed constraint as applicable or not applicable with
    a reason, including identifier strategy, HTTP and Better Auth policy, active
    ADR obligations, and tooling authority.
13. A Spec that proposes a tooling mutation cannot authorize its Tasks until
    it records the maintainer's express authorization and the bounded files the
    authorization covers. Silence, a generic implementation request, or setup
    completion is not authorization.
14. Generated Spec guidance, templates, and workflow instructions carry the
    same Project Constraint contract for human and Agent authors. Existing
    completed Specs are not rewritten.
15. Greenfield adoption and update use the same project-decision and Project
    Constraint contracts while retaining Spec 0047's Instruction Hierarchy and
    semantic ownership.
16. Documentation and examples show the suggested defaults, explicit
    alternatives, automation requirements, applicability rules, and refusal
    behavior accepted by the public Baseline Command.
17. The thin setup skill teaches the same project-decision and Project
    Constraint flows through the public Baseline Command and contains no
    independent decision collection, rendering, or mutation behavior.

## User Experience

On greenfield adoption, the maintainer sees UUID version 7 as the suggested
identifier strategy and, when Better Auth applies, the complete provider-owned
route exception. Each prompt offers a keep-or-change decision; pressing the
default choice remains an explicit answer.

On update, Roundfix presents compatible persisted values first. A changed or
incomplete decision must be resolved before the final Change Plan. Automation
receives stable missing-input diagnostics instead of prompts or inferred
defaults.

When a new Spec is authored, its PRD and TechSpec show a Project Constraints
section before Task decomposition. Tooling work remains blocked until the
maintainer records express authorization and the bounded file scope.

## Project Constraints

- Identifier strategy: applicable to the Baseline contract. The implementation
  models UUID version 7 as the suggested strategy but creates no application
  identifiers and performs no identifier migration.
- Authentication and HTTP: applicable to the Baseline contract. The
  implementation models Better Auth and its derived route exception but does
  not install Better Auth or mutate application routes.
- Active ADR obligations: applicable. The implementation must preserve
  repository-owned HTTP policy, CLI authority, confirmation-gated planning,
  portable Plans, and recoverable apply as defined by ADRs 0063, 0066, 0068,
  0071, and 0073.
- Tooling authority: applicable. This Spec grants no authorization to create,
  edit, rename, move, or delete tooling configuration, scripts, ignore files,
  plugin declarations, or version pins.

## Non-Goals / Out of Scope

- Migrating existing application identifiers to UUID version 7.
- Replacing an external provider identifier, protocol identifier, natural key,
  or business code with a project-owned identifier.
- Installing Better Auth, changing authentication providers, or creating
  application routes.
- Imposing the Better Auth exception on a repository that does not select
  Better Auth as a Repository Capability.
- Requiring an ADR to activate or explain an operative project rule.
- Mutating tooling configuration during Baseline adoption or treating setup
  approval as tooling authorization.
- Rewriting completed PRDs, TechSpecs, Tasks, or archived Specs.
- Modifying upstream-managed authoring or domain-modeling skills.
- Changing the semantic guide composition and migration contract owned by Spec
  0047.

## Success Metrics

- 100% of human greenfield journeys with unresolved decisions present the UUID
  version 7 suggestion and the conditional Better Auth exception for explicit
  keep-or-change confirmation.
- 100% of non-interactive journeys with unresolved decisions stop without
  mutation and name every missing stable input.
- Equivalent human and automation answers produce identical Decision Plans,
  Change Plans, and Plan Digests.
- 100% of new PRDs and TechSpecs contain a Project Constraints section; every
  confirmed constraint is marked applicable or not applicable with a reason.
- 100% of Specs proposing tooling mutation record express authorization and
  bounded files before any affected Task can execute.
- 0 generated domain or backend rules require an ADR to exist before they can
  be followed.
- 100% of public documentation examples and thin-skill instructions match the
  project-decision and Project Constraint behavior accepted by the Baseline
  Command.
- Greenfield and update fixtures retain confirmed values through apply,
  formatting, repository Verification, audit, and empty reapply with zero
  managed-file delta.
- Separately authorized Fluxus greenfield and update journeys require zero
  manual constraint repair and reproduce the accepted final guide set.
- The full Roundfix Verification passes after the Spec-local QA evidence and
  live Fluxus journeys are recorded.

## Decisions

- Spec 0048 depends on Spec 0047 reaching a passing QA verdict.
- UUID version 7 is the suggested default for new project-owned Internal
  Identifiers, not an inferred or universal replacement for external
  identities.
- Better Auth proposes an explicit `GET` and `POST` `/api/auth/*` exception
  owned by the provider; the maintainer can confirm or change it.
- Human callers confirm suggested defaults; automation supplies unresolved
  values explicitly.
- Tooling authority is a universal Normative Clause, not a configurable
  preference.
- Every new PRD and TechSpec accounts for Project Constraints before Task
  decomposition.
- Operative constraints are self-contained; ADRs preserve history but are not
  runtime instruction dependencies.
- Fluxus supplies the real acceptance journeys, but portable assets contain no
  Fluxus-specific names or policy.

## Open Questions

None.
