---
status: approved
granted: 2026-09-30
action: retire the Jev Router (its OpenCode selection, inline provider, relay, spend and credit gate, credit floor and router-prompt Judge Log lines), refuse any Agent Selection that reaches an OpenAI or Anthropic model through OpenRouter, and state the subscription rule in the Roundfix Skill, the guides and the Project Config comment
consuming: 0230-retire-the-jev-router
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - skills/roundfix/SKILL.md
  - .roundfixrc.yml
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0230

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário", and authorized the Baseline source, the guides and
`.roundfixrc.yml` with "Autorizar os dois". On 2026-10-05 the maintainer set
the rule "os modelos da openai e anthropic devem ser utilizado exclusivamente
pela assinatura que o codex e claude fornecem", decided "Aposentar o router"
and confirmed "Vamos seguir com a remoção".

The governed set was measured with `GovernedPath` on the authoring branch at
`a1b8b1d7`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/runtime.md` — its Jev Router section
  describes the retired selection, gate, floor and relay, and gives way to
  the subscription rule.
- `.roundfixrc.yml` — a comment only: the profile commentary already says
  never to route a `gpt-*-sol` model through OpenRouter, and gains the rule
  and its refusal. No configuration value changes.

## What is not governed

The Go sources and tests under `internal/jevrouter`, `internal/agent`,
`internal/daemon`, `internal/cli` and `internal/config` that the Tasks
declare, `skills/roundfix/references/runtime.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/configuration.md` and `docs/references/model-selection.md`
are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirror. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No configuration value in `.roundfixrc.yml` changes; no built-in or
  Recommended Profile changes.
- No change to the judge (`internal/judge`, `roundfix spec judge`), its keys,
  its ceiling or its Judge Log schema.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md`, `CHANGELOG.md` or the `### QA settlement` section of any
  skill.
- No test, Verification command or QA row reaches OpenRouter or TypeSafe, and
  no test reads or writes the real `~/.roundfix`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
