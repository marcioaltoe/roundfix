---
spec: 0061-repository-derived-skill-requirements
status: active
created: 2026-07-29
surfaces: [backend, cli, docs]
---

# Repository-derived skill requirements

Repository Skill Set readiness requires a fixed list embedded in the Roundfix
binary instead of the skills the repository's own Setup Manifest selects, so
every consumer repository is told it needs Roundfix's development stack. Three
TypeScript repositories fail the check demanding `golang-cli`,
`golang-concurrency`, `bubbletea`, and `tui-design`; the only one that passes
does so because it historically carried those files. Evidence:
[Doctor demands Roundfix's own development skills](../../findings/2026-07-29-doctor-requires-roundfix-own-development-skills.md).

The machinery to do this correctly already exists and is bypassed: Baseline
modules declare their own `requiredSkills`, and the Setup Manifest records
which modules a repository selected.

## Project Constraints

- Identifier strategy: not applicable — skill names, module identifiers, and
  manifest paths keep their existing identities; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the readiness check reads local
  files only and contacts no network. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0049 and ADR-0055 keep Agent
  Selection Profile proof independent from Repository Skill Set readiness, so
  this change must not couple them; ADR-0066 keeps Baseline execution in the
  Go CLI; the Doctor Command's diagnosis-only contract must survive
  unchanged. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-29 the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout that
  `make baseline-digests` rewrites. No other protected tooling mutation is
  authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- A repository is required to carry only the skills its own selected modules
  declare, so a TypeScript project is never told it needs a Go or terminal-UI
  skill.
- A missing external skill is remediated by a command that installs exactly
  that skill, never a whole upstream catalog.
- A repository without a Setup Manifest keeps a defined, documented behavior
  rather than an accidental one.

## User Stories

1. As a maintainer of a TypeScript repository, I want readiness to require
   only what my selected modules declare, so that a green Doctor means my
   repository is ready rather than that it mirrors Roundfix's own stack.
2. As a maintainer facing a missing external skill, I want the printed next
   action to install exactly the missing skills, so that remediation does not
   import an unreviewed catalog of instructions that run with full agent
   permissions.
3. As a maintainer of a repository that has not adopted the Baseline, I want
   readiness to state that plainly instead of inventing a requirement.
4. As a Roundfix maintainer, I want this repository to keep requiring its Go
   and terminal-UI skills, because its own manifest selects those modules.

## Core Features

1. The required external skill set is derived from the repository's Setup
   Manifest: the modules it records, the `requiredSkills` those modules
   declare in the Baseline catalog, minus the Roundfix-owned skills, as a
   deduplicated deterministic set.
2. The required owned skill set is unchanged: the running binary's embedded
   bundle stays authoritative for Roundfix-owned skills.
3. A repository with no readable Setup Manifest requires no external skills
   and says so in the readiness line, rather than falling back to a list it
   cannot justify.
4. A missing-external-skill failure names every missing skill at once and
   prints a per-skill install command, not a package-wide install.
5. The Doctor Command remains diagnosis-only: it reads, reports, and mutates
   nothing.

## User Experience

- `roundfix doctor` in a TypeScript repository reports `skills: ok` with the
  count its modules justify; in this repository the count still includes the
  Go and terminal-UI skills its manifest selects.
- A failure lists each missing skill and the exact command that installs it.
- A repository without a Setup Manifest reports its external requirement as
  zero and names the Baseline adoption command as the next step.

## Non-Goals / Out of Scope

- Reconciling a consumer lock against upstream removals — the third finding in
  the source report, left for its own Spec.
- Changing the embedded owned bundle, its version contract, or how
  `roundfix skills install` works.
- Adding a network call to readiness.
- Deciding whether `autoresearch` and `exa-web-search` belong to a module.
  They are declared by none today, so they simply stop being required; adding
  them to a module is a separate Baseline decision.

## Success Metrics

- In a repository whose manifest selects the TypeScript modules, readiness
  requires no `golang-*`, `bubbletea`, or `tui-design` skill.
- In this repository, whose manifest selects `go`, `cli-surface`, and
  `tui-surface`, readiness still requires those skills.
- A missing external skill produces one failure naming every missing skill and
  a per-skill install command.
- A repository with no Setup Manifest reports zero external requirements and
  the adoption next action.
- The Doctor Command still writes nothing.

## Decisions

- The Setup Manifest is the authority for which modules a repository selected;
  the catalog is the authority for what each module requires.
- No Setup Manifest means no external requirement, which is honest, rather
  than a fallback to the embedded list that produced this defect.
- `recommended.txt` stays as advisory recommendation data and stops being the
  readiness authority.

## Open Questions

None.
