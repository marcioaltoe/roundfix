---
spec: 0030-context-driven-agent-instructions
status: active
created: 2026-07-15
surfaces: [cli, docs]
---

# Context-driven agent instructions

Repository maintainers need a fast, repeatable way to establish and maintain effective agent instructions without replacing repository-specific knowledge or repeatedly answering settled setup questions. The setup workflow must provide concise, modular guidance, validate both the root agent instructions and their supporting guides, and preserve enough durable configuration to make future updates safe and mostly automatic.

## Goals

- Keep root agent instructions concise by routing detailed, conditional guidance to canonical supporting documents.
- Let maintainers compose portable instruction modules for their project type and technology stack while preserving repository-authored content.
- Detect missing, stale, duplicated, conflicting, or unmanaged setup guidance and safely repair only setup-owned content.
- Reuse confirmed setup decisions across runs and ask only when a new or changed rule requires human judgment.
- Keep documented skill requirements aligned with the selected canonical skill setup and make installation drift visible.

## User Stories

1. As a repository maintainer, I want to select modular instruction profiles for my project, so that the generated guidance fits the repository without copying a large monolithic sample.
2. As a repository maintainer, I want the setup workflow to remember confirmed decisions, so that rerunning it does not repeat settled questions.
3. As a repository maintainer, I want to audit the root agent instructions and supporting agent guides, so that missing, duplicated, stale, and conflicting guidance is reported consistently.
4. As a repository maintainer, I want safe corrections to affect only setup-owned content, so that repository-specific instructions remain intact.
5. As a repository maintainer, I want the selected profile's required skills checked against the canonical skill setup and the repository's installed skills, so that agents are not instructed to use unavailable capabilities.
6. As a repository maintainer, I want an optional list of installed skills outside the selected setup, so that I can evaluate cleanup candidates without the workflow removing anything.
7. As a repository maintainer with a local Secondbrain, I want to opt into concise usage guidance, so that agents know when and how to retrieve cross-project knowledge without exposing secrets or writing to the knowledge workspace.
8. As a setup-context-driven maintainer, I want portable samples and setup snapshots to expose upstream drift, so that the distributed workflow stays aligned with the canonical skill catalog.

## Core Features

1. The workflow provides a small universal instruction core and combinable modules for project type, language, runtime, repository shape, application surface, workflow, and optional integrations.
2. The initial module catalog prioritizes TypeScript, Bun, monorepos, backend and frontend application surfaces, while also supporting Go CLI/TUI and Rust CLI repositories.
3. A durable setup record identifies the selected modules, confirmed choices, setup-owned documents, and the template versions used to produce them.
4. Setup-owned instruction blocks and supporting guides have stable identities that allow the workflow to update them in place without relying on prose matching.
5. Auditing reports missing required sections, duplicate managed sections, stale template versions, invalid references, conflicting managed rules, unrecognized managed content, and missing supporting guides.
6. Safe correction adds or refreshes only setup-owned canonical content, preserves repository-authored sections, and remains idempotent across repeated runs.
7. The workflow validates the root agent instructions and every setup-owned guide under the repository's agent documentation area; nested agent-instruction files remain outside the initial release.
8. The workflow asks for confirmation only when a required decision has no durable answer, when a stored answer is no longer valid, or when a newly introduced module or rule requires human judgment.
9. Each project profile identifies one canonical skill setup and the skills required by its generated rules. Missing required skills or rules that reference skills outside that setup are blocking findings.
10. Installed skills outside the selected setup are available as an optional, non-blocking cleanup report. The workflow never removes skills.
11. Portable setup data records the canonical setup revision used to generate each bundled skill list and can detect drift against an available canonical setup source without requiring that source during normal use.
12. Secondbrain integration is opt-in. When selected, the workflow creates a detailed read-only usage guide and keeps only a concise mandatory pointer in the root agent instructions.
13. The audit can produce concise human-readable results and stable machine-readable results suitable for agent automation and repository verification.
14. Findings distinguish blocking errors, safe corrections, decisions requiring confirmation, and informational cleanup candidates.

## User Experience

On first use, the maintainer sees detected project characteristics, the proposed module composition, and only the decisions that cannot be derived safely. The workflow previews all managed changes before applying them and clearly distinguishes setup-owned content from repository-authored content.

On later runs, the workflow loads the durable setup record, validates the current repository against it, and proceeds without questions when the stored decisions remain valid. When the module catalog evolves, the maintainer is asked only about the new decision and can see why it is required.

Audit results are concise by default and group findings by severity and affected document. Detailed or machine-readable output is available for automation. Optional skill-cleanup output lists candidates without modifying the repository.

## Non-Goals / Out of Scope

- Rewriting or normalizing repository-authored instruction sections.
- Automatically removing installed skills.
- Validating or generating nested agent-instruction files in the initial release.
- Requiring access to a developer-specific skills checkout or Secondbrain path during normal use.
- Writing to the Secondbrain or editing its mirrors and raw sources.
- Replacing general Markdown linting, repository verification, or domain-specific documentation review.
- Automatically adopting rules whose addition requires a new product, architecture, security, or workflow decision.

## Success Metrics

- A second run with unchanged choices completes without asking a previously answered question.
- Applying safe corrections twice produces no changes on the second application.
- All setup-owned blocks and guides can be identified without matching their prose.
- Every skill referenced by a selected module is present in its bundled canonical setup snapshot.
- Missing required skills produce blocking findings, while extra installed skills never block validation and are never removed.
- The TypeScript/Bun monorepo, Go CLI/TUI, and Rust CLI profiles all pass the same validation contract.
- Repository-authored content outside managed boundaries remains byte-for-byte unchanged after safe correction.
- The optional Secondbrain module produces a concise root pointer and a complete read-only usage guide.

## Decisions

- Use audit plus safe correction: report all findings and update only setup-owned canonical content.
- Keep samples portable inside `setup-context-driven` and compose them from reusable modules rather than monolithic project copies.
- Prioritize TypeScript, Bun, and monorepo repositories; include Go CLI/TUI and Rust CLI profiles in the initial release.
- Store confirmed decisions in one canonical manifest and identify managed content with stable ownership markers. See ADR-0046.
- Validate the root agent instructions and setup-owned supporting guides together.
- Make Secondbrain integration an explicit opt-in module.
- Treat missing required skills as blocking and extra installed skills as optional, non-blocking cleanup information.
- Never remove installed skills.
- Preserve the established autonomous implementation defaults: Codex `gpt-5.5` with `xhigh` for backend Tasks and Claude Opus 4.8 with `xhigh` for frontend, UI, UX, and design Tasks.
- Write all generated repository content in English.

## Open Questions

None.
