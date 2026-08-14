---
status: done
created_at: 2026-07-24
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-baseline-and-derived-tooling.md
---

# Context-Driven Baseline — greenfield guidance acceptance target (2026-07-24)

A greenfield Baseline run against Fluxus produced the compact root and modular
`docs/agents/` structure expected from the planned guidance work, but the live
result still needed project rules to be redistributed from
`docs/agents/specific-repository.md`. This finding records the accepted target
for the two follow-up Specs: generated guidance must be self-contained, must
not depend on an ADR being present, and must propagate applicable project
constraints into every new Spec.

## 1. Identifier strategy needs a typed, confirmable default

- **Symptom / evidence**: The generated guidance had no project identifier
  policy, so a new Spec or implementation could choose an incompatible ID
  format without surfacing a repository decision.
- **Root cause**: The current Baseline Profile has no typed identifier-strategy
  decision or render binding.
- **Action / suggestion**: Add a typed project decision whose suggested
  default is UUID version 7 for every new project-owned identifier. Human
  setup must ask the maintainer to keep or change that suggestion when no
  compatible prior decision exists; automation must require the explicit
  value. Render the confirmed rule into `docs/agents/domain.md`, and require
  every generated PRD and TechSpec to state whether the constraint applies.
  External provider and protocol identifiers retain their source contract.

## 2. Better Auth requires a self-contained typed route exception

- **Symptom / evidence**: The Fluxus HTTP guidance named an ADR as part of the
  operative rule. A greenfield repository may not have that ADR, so the guide
  was not independently sufficient.
- **Root cause**: The rendered rule treated repository evidence as an
  instruction dependency instead of rendering the complete confirmed
  decision.
- **Action / suggestion**: When the selected project uses Better Auth, propose
  a confirmable typed exception for `GET` and `POST` at `/api/auth/*`, owned by
  Better Auth so session, OAuth redirect, callback, and related provider
  protocol semantics remain intact. The maintainer must confirm or change the
  suggestion; setup must not silently infer authorization. Render the complete
  confirmed contract into `docs/agents/backend.md` and into every generated
  Spec that changes authentication or HTTP routes. ADRs may preserve decision
  history, but generated rules must operate without them.

## 3. Tooling configuration authority is a Normative Clause

- **Symptom / evidence**: The portable guide required authority only before
  intentionally changing selected verification configuration, leaving
  ambiguous whether renames, deletes, ignore files, plugins, scripts, and
  version pins were protected.
- **Root cause**: The rule was narrower than the maintainer's tooling
  ownership boundary.
- **Action / suggestion**: Ship one non-optional Normative Clause requiring
  explicit maintainer authorization before creating, editing, renaming,
  moving, or deleting any linter, formatter, typechecker, test-runner,
  architecture-checker, build-tool, package-manager, code-generator, or other
  repository-tooling configuration, script, ignore file, plugin declaration,
  or version pin. This rule is not a setup preference and must not produce a
  confirmation prompt. Every generated Spec repeats the constraint and records
  the authorization plus bounded files only when a tooling change is approved.

## 4. Repository-specific extension is unnecessary after distribution

- **Symptom / evidence**: The live Fluxus result retained a root pointer to
  `docs/agents/specific-repository.md`, but its HTTP and validation rules have
  stable owners in `backend.md`, `frontend.md`, `domain.md`,
  `agent-instructions.md`, and `spec-routing.md`.
- **Root cause**: The initial adoption preserved repository rules in a safe
  catch-all carrier before their final semantic owners were established.
- **Action / suggestion**: The generated root must contain only compact
  pointers to managed guides. Do not create or link
  `docs/agents/specific-repository.md` when every accepted rule has a managed
  semantic owner. A repository extension remains available only for a
  non-empty accepted rule that cannot be represented by a typed decision or an
  owned guide.

## Acceptance oracle

The two follow-up Specs must treat the following Fluxus layout and semantics as
the acceptance target, while keeping profile assets portable and free of
Fluxus branding:

- compact `AGENTS.md` pointers;
- `docs/agents/agent-instructions.md` for universal execution, evidence,
  tooling-authority, Git, dependency, and safety rules;
- `docs/agents/domain.md` for domain language, context layout, and UUID
  version 7 identifier policy;
- `docs/agents/backend.md` for backend boundaries and the complete confirmed
  HTTP and Better Auth contract;
- `docs/agents/frontend.md` for frontend boundaries, import safety, observable
  tests, and browser inspection;
- `docs/agents/spec-routing.md` and `docs/agents/issue-tracker.md` for Spec
  routing, constraint propagation, Task ownership, verification, and closure;
- `docs/agents/autonomous-work.md`, `docs/agents/docs-layout.md`,
  `docs/agents/monorepo.md`, `docs/agents/secondbrain.md`,
  `docs/agents/skill-dispatch.md`, and `docs/agents/typescript-bun.md` for
  their selected concerns;
- `docs/agents/setup-context.json` as the typed, auditable record of confirmed
  decisions; and
- no `docs/agents/specific-repository.md` or root pointer when no residual
  repository-specific rule remains.

The
[Context-Driven guidance composition Spec](../specs/0047-context-driven-guidance-composition/_prd.md)
owns profile composition, semantic rule distribution, self-contained rendered
guidance, and the empty-extension outcome. The
[Context-Driven project decisions and Spec constraints Spec](../specs/0048-context-driven-project-decisions-and-spec-constraints/_prd.md)
depends on it and owns interactive and non-interactive decision collection,
UUID version 7 and Better Auth confirmations, and propagation of applicable
constraints into generated Spec artifacts.

## 5. The modular instruction hierarchy must be explicit

- **Symptom / evidence**: The generated files had distinct semantic owners,
  but their root order did not explain which rules were universal, contextual,
  workflow-specific, runtime-specific, or selected by the technology profile.
- **Root cause**: Module activation described membership without rendering an
  instruction-precedence map.
- **Action / suggestion**: Keep the modular files and current transparent
  names. Render the root pointers and core guide in this order: universal core;
  context and documentation; Spec workflow; autonomous execution when enabled;
  stack and surface guides; optional knowledge sources. A narrower guide may
  add constraints but must not weaken a universal Normative Clause or a
  confirmed project decision. Renaming remains possible only through an
  explicit artifact-ID, path, Manifest, retention, and upgrade migration.

## 6. ADR lifecycle metadata must augment, not replace, external skills

- **Symptom / evidence**: Existing ADRs use either a short body-only format or
  an inline `Status: Accepted` line. Neither form provides consistent
  machine-readable lifecycle timestamps.
- **Root cause**: The Baseline has no repository-owned ADR lifecycle
  frontmatter, while `.agents/skills/domain-modeling/ADR-FORMAT.md` is
  externally owned and must not be modified.
- **Action / suggestion**: Generate a repository-owned frontmatter overlay in
  `docs/agents/docs-layout.md` with statuses `proposed`, `accepted`, `rejected`,
  `deprecated`, and `superseded`; RFC 3339 UTC `created_at`, `updated_at`, and
  nullable `deprecated_at`; plus nullable `superseded_by`. Only `accepted` is
  active. Preserve the external skill's body format and prepend the confirmed
  repository metadata when creating a new ADR. Treat a legacy ADR without
  lifecycle frontmatter as accepted unless its body explicitly marks it
  inactive. Do not rewrite existing ADRs solely to adopt metadata.

## 7. Findings guidance needs a complete copyable template

- **Symptom / evidence**: The generated guide explained Findings frontmatter
  and status transitions but did not show the complete evidence-first body or
  dated addendum structure.
- **Root cause**: The Baseline rendered lifecycle rules without the document
  skeleton required to apply them consistently.
- **Action / suggestion**: Render one copyable Findings template in
  `docs/agents/docs-layout.md` containing the existing
  `pending | partial | deferred | done` frontmatter, session context,
  symptom/evidence, proven-or-unknown root cause, action or Spec routing,
  optional “What worked — keep,” and dated append-only addenda. Keep Findings
  dates and routing semantics compatible with the existing contract.

## Addendum — 2026-07-24 — Oraculum greenfield profile alignment

A greenfield run against Oraculum selected
`standard-typescript-monorepo`, collected all Baseline decisions, and then
stopped with one aggregate `required profile alignment is unresolved` result.
The repository is a TypeScript/Bun/Turborepo backend monorepo without the
Profile's frontend workspace, React stack, Better Auth, or LogTape, and its
Repository Skill Set lacks the universal Context7 and Exa capabilities.

The required-capability block is correct, but the human workflow did not expose
the divergence-resolution state promised by Spec 0046. Spec 0047 now owns the
accepted correction: propose a catalog-valid repository-owned Profile
adaptation for profile-specific gaps, review every removal, include the Profile
file in the digest-bound Change Plan, and return automatically to audit.
Universal requirements remain non-waivable and receive exact remediation
operations.
