---
status: accepted
created_at: 2026-11-17T13:30:00Z
updated_at: 2026-11-17T13:30:00Z
deprecated_at: null
superseded_by: null
---

# Doctor and the Update Check recommend an Authoring Model Declaration when none is declared

When a repository has installed authoring skills it does not publish and
declares no `skills.authoring.claude`, the Doctor Command and the Update Check
warn and print the Recommended Authoring Declaration, a block to add to
`.roundfixrc.yml`. The binary carries the recommendation as one table in
`internal/config`, beside the Recommended Profile, and a repository contract
test holds this repository's own declaration and the configuration guide's
block equal to it. Nothing fails and nothing is written. This amends
ADR-0401, which shipped no model name.

## Context

ADR-0401 made the declaration optional and silent: "Nothing is declared by
default, and a repository that declares nothing keeps every byte Roundfix
writes today." The checks report only a declaration that does not match the
installed files. Measured on 2026-11-17, Argus has all four authoring skills
installed, no `skills:` key and no generated subagent, and
`roundfix update check --format json` there carries no `authoringModels`: its
authoring phases run on whatever model the session uses, and no report says
so. The maintainer chose on 2026-11-17 "Avisar e sugerir (Recommended)":
Doctor and the Update Check warn in that case and show the block to paste,
with the models Roundfix itself declares; they never fail and never write.
The same day he set the values (verbatim): "Bloco recomendado: - Fable xhigh para Ideia - Fable high para PRD e TechSpec - Opus high para Tasks - Sonnet medium para pesquisa".
This repository's own declaration then said `fable`/`xhigh` for `prd` and
`techspec` and `fable`/`high` for `tasks`, and changes to those values.

ADR-0401 rejected "Model names shipped as defaults" because "a name ages and
assumes one vendor". The Recommended Profile already ships dated model names
that never select anything by themselves; this recommendation takes the same
shape. The binary cannot read this repository's `.roundfixrc.yml` in an
adopter, and no Go package sits at the repository root to embed it.

## Decision

**When.** The report applies when Project Config was read, the repository
root is known, no `skills.authoring.claude` entry is declared, the
inspection of ADR-0401 finds no problem, and at least one of `write-idea`,
`write-prd`, `write-techspec` and `write-tasks` has
`.agents/skills/<name>/SKILL.md` while the repository does not publish
`skills/<name>/SKILL.md`. A repository that publishes every installed
authoring skill stays exempt (ADR-0408). A problem the inspection finds,
such as a leftover generated subagent, keeps today's report unchanged.

**What.** Doctor prints `authoring-models: warn` naming the installed
skills, with the next action to add the block and run
`roundfix skills install --target project`, followed by the block indented
by two spaces. The Update Check prints the same in its text form, and in JSON
sets `authoringModels.state` to `undeclared` with `installed` and
`recommendation`, the block as text. Neither changes an exit code or the
Update Check's `next`.

**The source.** `internal/config` holds the Recommended Authoring
Declaration as one table: the model and reasoning effort of `idea`, `prd`,
`techspec`, `tasks` and `research` (`fable`/`xhigh`, `fable`/`high`,
`fable`/`high`, `opus`/`high` and `sonnet`/`medium`), and the function that
renders the block.
The block holds no `research.tools`: reaching the network or an MCP server
stays each repository's choice (ADR-0402). A repository contract test fails
when this repository's `.roundfixrc.yml` declares another model or effort for
any key, and when the block in `docs/user-guide/configuration.md` differs from
the rendering. Changing this repository's declaration therefore changes the
recommendation in the same change, or the contract run stops it.

## Considered Options

- **Read this repository's `.roundfixrc.yml` at build time.** Rejected: no Go
  package sits at the repository root, so `go:embed` cannot reach it, and a
  build step that copies it would add a generated file and its regeneration
  rule to the repository's tooling.
- **A generated record with a record flag.** Rejected: one more derived record
  and a `derived_paths` entry in `.roundfixrc.yml` for five pairs of values a
  contract test can hold equal.
- **The Model Catalog.** Rejected: it lists what a runtime's picker offers,
  without a phase or an effort.
- **A hard-coded block in the command code.** Rejected: a second copy that
  drifts from the repository's declaration without a test that notices.
- **Writing the block into `.roundfixrc.yml`.** Rejected by the maintainer:
  the checks never write.
- **Failing the check.** Rejected, as in ADR-0401: the skills stay usable, and
  `make verify` and the release plan read those checks.
- **The repository's `research.tools` in the block.** Rejected: this
  repository lists MCP tools of its own machine, and ADR-0402 leaves the
  network to each repository.
- **The same warning in `roundfix skills check`.** Rejected: the maintainer
  named Doctor and the Update Check, and `skills check` runs inside
  installation flows, where a block of text would repeat on every install.

## Consequences

- Roundfix now ships model names in a recommendation, as it does in the
  Recommended Profile; it still applies none by itself.
- When the maintainer changes this repository's authoring models, the
  recommendation, the configuration guide and the release that carries them
  change together.
- A repository that wants its authoring phases on the session's model
  declares nothing and keeps seeing the warning; this decision adds no key to
  silence it.
- An adopter that pastes the block and installs gets the five generated
  subagents of ADR-0401 and a Research Subagent with no web tool until it
  lists one.
