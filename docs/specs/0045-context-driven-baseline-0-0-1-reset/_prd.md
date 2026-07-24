---
spec: 0045-context-driven-baseline-0-0-1-reset
status: active
created: 2026-07-22
surfaces: [cli, infra, docs]
---

# Roundfix 0.0.1 Context-Driven Baseline reset

A real managed upgrade removed operational guidance while its generated corpus still passed internal coverage checks: 27 coarse entries stood in for a 573-line instruction source, and the generated agent-guidance set deleted 187 lines while adding 73. Roundfix needs a project-agnostic, independently auditable baseline that preserves every governed instruction and structured contract, gives the repository explicit control over non-portable policy, and restarts all Roundfix-owned version surfaces together at `0.0.1`.

## Goals

- Every source instruction has an independently provable disposition before setup can mutate a repository.
- The Standard TypeScript Monorepo Profile reproduces the confirmed opinionated stack, workflow, skill activation, architecture, and quality contracts without project-specific names or policy.
- Existing repositories can adopt the new baseline without silent deletion, while incompatible historical setup identities carry no compatibility obligation.
- The CLI, baseline, schemas, manifests, and Roundfix-owned distributed skills expose one coherent `0.0.1` version line.
- Generated root instructions remain compact while agent guides retain complete Operational Contracts, including templates, procedures, conditions, exceptions, and next actions.

## User Stories

1. As a repository maintainer, I want every prior Normative Clause, recommendation, and Operational Contract accounted before I authorize setup, so that an apparently successful upgrade cannot silently weaken my repository.
2. As a maintainer adopting the Standard TypeScript Monorepo Profile, I want required capabilities proved and missing decisions asked explicitly, so that generated instructions describe the repository I actually operate.
3. As a maintainer with project-specific policy, I want setup to recognize its typed owner or propose Repository-Specific Normative Rules, so that project rules survive without contaminating the portable baseline.
4. As an Agent, I want complete skill triggers and operational guides, so that a compact root index never hides the rules needed to execute work safely.
5. As a maintainer performing a Baseline Readoption, I want a fresh Change Plan with individual dispositions, so that incompatible historical metadata never becomes permission to overwrite content.
6. As a release maintainer, I want every Roundfix-owned version surface and release artifact to restart coherently, so that `0.0.1` identifies one contract instead of a mix of old generations.

## Core Features

1. The Roundfix CLI, Context-Driven Baseline, setup schemas, Setup Manifests, profile and module versions, and every Roundfix-owned distributed skill restart at `0.0.1`. Upstream-managed skills retain upstream ownership and versions. Existing User Config, Project Config, Runs, and Run Database state remain intact unless an independent functional change requires migration.
2. The Standard TypeScript Monorepo Profile replaces the previous TypeScript profile identity. It requires TypeScript, Bun, Turborepo, Vite, React, Hono, Drizzle, Zod, Tailwind, shadcn, TanStack Query and Router, Better Auth, PostgreSQL, LogTape, Oxlint, Oxfmt, and Vitest. Setup detects and persists existing versions, asks when a version is absent or incompatible, and blocks readiness when a required capability is missing.
3. The profile requires canonical frontend and backend workspaces. Frontend feature code uses domain systems with a public system boundary; internal modules import directly. Backend code uses domain, application, and infrastructure layers, with thin HTTP handlers, HTTP-independent use cases, Drizzle-owned persistence, and no generic modules or services layer. Inngest and Docker remain optional capabilities.
4. Every supported Source Baseline has an immutable, project-agnostic source corpus and an independent Normative Clause Manifest. The manifest inventories every Normative Clause, non-blocking recommendation, and Operational Contract with stable identity, enforcement, carrier, and source evidence. Removing an entry from both the corpus accounting and the transition cannot keep validation green.
5. The Upgrade Retention Contract defaults to semantic retention. Every entry maps individually to a current managed target, an existing typed repository document, Repository-Specific Normative Rules, or an explicit rejection with its own confirmed reason. Aggregate categories cannot replace independently triggered clauses or structured contracts.
6. Templates, ordered procedures, and decision matrices are first-class Operational Contracts. A prose summary does not satisfy a findings template, lifecycle, research sequence, skill-trigger matrix, Secondbrain protocol, Spec route matrix, or other shape whose structure carries behavior.
7. Setup recognizes existing project-specific policy without modifying its typed document. Missing non-portable clauses are previewed as exact Repository-Specific Normative Rules with a digest; after the confirmed first write, setup preserves those bytes. Disabling their destination while clauses remain unresolved blocks mutation.
8. Each repository owns an HTTP Contract Decision. Setup reuses a recognized REST or POST-only contract and its typed exceptions, or asks when no contract exists. Hono endpoint creation, modification, and review require the Hono contract, Hono, and Zod skills; persistence work additionally requires Drizzle guidance.
9. Skill activation is a declarative matrix of stable triggers and required bundles. Production-code work always activates `coding-guidelines`, `clean-code`, and `solid` plus matching domain skills. React feature work, UI and accessibility work, testing, debugging, security, delivery, and optional capabilities retain their distinct trigger scopes and multi-skill bundles.
10. Context7 and Exa are required research capabilities; their absence blocks readiness with a next action. Broad external research uses three to seven varied Exa searches, while external research never substitutes for local code search. Firecrawl, `rtk`, and `rg` are recommended capabilities whose absence produces a warning and an explanation on request.
11. The repository's selected Verification is persisted. Setup reuses a declared Verification or asks when none exists. Verification may format and fix lint before its blocking checks, but it cannot weaken configuration or hide diagnostics; Oxlint warnings block completion. Formatting, Verification, audit, and reapply must compose without a generated delta.
12. Generated root instructions remain a short mandatory index. Agent guides retain the complete findings template and lifecycle, Spec routing, Task ownership, Supervisor delegation boundary, backend and frontend contracts, language policy, dependency discipline, Git authority, test discipline, and the full optional Secondbrain protocol. External trackers are unsupported; Specs and Tasks remain the only issue-tracking system.
13. Older setup markers or manifests enter an explicit Baseline Readoption. Setup inventories the repository as a Source Baseline, resolves every disposition, and creates fresh `0.0.1` state only after confirmation. It neither silently interprets old transitions nor treats incompatibility as permission for a clean overwrite.
14. Go CLI/TUI and Rust CLI profiles remain available and restart their Roundfix-owned versions at `0.0.1`; this feature does not re-specify their content. The Standard TypeScript Monorepo Profile alone receives the new opinionated stack contract.
15. The repository changelog restarts at `0.0.1`. Previous tags and GitHub Releases are removed only after implementation and QA pass, through a read-only Release Plan followed by a separate explicit confirmation. Git history, Specs, operational state, and non-superseded ADRs remain.

## User Experience

Setup first reports required capabilities, recommended-capability warnings, detected versions, existing Verification, repository-owned contracts, and unresolved decisions. It asks one decision at a time when the repository cannot answer, then presents one digest-bound Decision Plan and Change Plan.

Every proposed exclusion names one source entry and requires its own reason. Interactive setup asks for those decisions individually; non-interactive setup exits without mutation and lists every unresolved identifier plus its next action. The final confirmation never serves as implicit approval for unreviewed exclusions.

During Baseline Readoption, maintainers see the exact current source identity, every managed and repository-specific destination, all structured Operational Contracts, and the bytes proposed for Repository-Specific Normative Rules. Unknown or changed input fails closed instead of falling back to fuzzy matching.

Required-capability failures block readiness. Recommended-capability gaps are warnings and reveal their rationale when requested. No flow installs the application stack, changes an HTTP policy, edits a typed repository document, removes remote releases, or mutates an existing repository before the relevant authority is explicit.

## Non-Goals / Out of Scope

- Copying a real project, its names, branding, tokens, paths, or generated artifacts into a distributed skill, agent guide, fixture, or profile.
- Copying the raw legacy sample into Roundfix; the maintained corpus is project-agnostic and records project-specific exclusions individually.
- Preserving compatibility with previous setup profile identifiers, baseline versions, schemas, manifests, or synthetic transition fixtures.
- Deleting or recreating User Config, Project Config, Runs, the Run Database, Git history, Specs, or accepted and partially superseded ADRs.
- Changing versions or content owned by upstream-managed skills.
- Installing application dependencies, provisioning PostgreSQL, generating application scaffolding, or adding missing required stack capabilities.
- Rewriting existing design, domain, architecture-decision, or package instruction documents during adoption.
- Supporting an external issue tracker or its triage-label vocabulary.
- Re-specifying the Go CLI/TUI or Rust CLI profile content.
- Applying setup to an external project as part of implementation or QA.
- Deleting tags or GitHub Releases before the post-QA Release Plan and explicit confirmation.

## Success Metrics

- 100% of entries in every supported source corpus have exactly one independently validated disposition; zero entries may be omitted, duplicated, or represented only by an aggregate category.
- 100% of Operational Contracts retain their required structure and content; compressed prose cannot satisfy a structured-contract fixture.
- 100% of Roundfix-owned version surfaces report `0.0.1`, and 0 upstream-managed skill versions change.
- The Standard TypeScript Monorepo Profile proves every required capability and blocks on one missing required capability; each absent recommended capability produces exactly one warning and no readiness failure.
- All mandatory skill-trigger fixtures activate the exact required bundles for production code, frontend, UI quality, endpoints, persistence, tests, debugging, security, QA, and delivery.
- The generated corpus contains 0 project-specific names, brands, product tokens, or external-project dependencies.
- Interactive and non-interactive fixtures prove that one unresolved exclusion, HTTP Contract Decision, version choice, Verification choice, or Repository-Specific Normative Rule blocks mutation.
- Apply, formatter, Verification, audit, and second apply produce zero generated-file delta for every maintained profile.
- Repository documentation contains 0 `done` findings, 0 findings outside the required template, and 0 ADR files explicitly marked fully superseded after the authorized cleanup.
- The full repository Verification passes with no errors or warnings before the Spec can complete.
- The post-QA Release Plan identifies 100% of previous tags and GitHub Releases before any remote deletion, and mutation still requires a separate confirmation.

## Decisions

- Semantic retention is the default; exclusion is individual, reasoned, and confirmation-gated. See ADR-0060.
- The canonical terms are Source Baseline, Baseline Readoption, Normative Clause, Normative Clause Manifest, Operational Contract, Repository-Specific Normative Rules, and HTTP Contract Decision.
- The profile identifier is `standard-typescript-monorepo`, and it is deliberately opinionated. See ADR-0061.
- The entire Roundfix-owned version line restarts at `0.0.1`; operational state and upstream-owned versions do not. See ADR-0062.
- HTTP semantics belong to each repository and cannot be inferred from Hono capability. See ADR-0063.
- Code and technical artifacts use English; domain documentation follows the language declared by the repository's domain context.
- Environment variables required by the application have safe example keys, and secret exposure remains prohibited.
- The sample's style values are suggestions when no repository contract exists; setup asks before adopting them.
- The Secondbrain protocol is complete when enabled, but the obsolete sparse-checkout knowledge-workspace and separate documentation-commit flow are excluded.
- Static Fable runtime defaults are excluded; autonomous work uses effective Agent Selection Profiles.
- Archive remains an explicit post-QA action, and workflow invocation authorizes only the commits that workflow declares.

## Open Questions

None.
