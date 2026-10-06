---
status: approved
granted: 2026-10-05
action: run complexity-low Tasks on an open OpenRouter model through OpenCode by default, with a User Config allow-list, a separate monthly ceiling, a Light Spend Log and one escalation to the Preferred Selection, let the advisory judge suggest a tier, and describe the light tier in the Roundfix Skill, the write-tasks skill and the guides
consuming: 0233-a-light-tier-on-open-models
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0233

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário", and authorized the Baseline source, the guides and
`.roundfixrc.yml` with "Autorizar os dois". On 2026-10-05, answering structured
questions, the maintainer chose the scope "W e X", set the light tier on for
every project ("Ligado para todos"), capped its spend with a new User Config
key `openrouter.implement_monthly_ceiling_usd` defaulting to US$10 and summed
across repositories, separate from the judge's ceiling, chose that a light
Task without an OpenRouter key falls back to the profile's default with a
report warning ("Cai para o padrão"), granted the governed paths this Spec
declares ("Concedo"), and allowed live test spend of up to US$5 on open models
only, with the cost recorded.

The governed set was measured with `GovernedPath` on the authoring branch at
`9085540d`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/runtime.md` — the runtime reference
  describes Agent Selection and the subscription rule; it gains the light
  tier, its keys, its warnings and its escalation.
- `.agents/skills/roundfix/references/spec.md` — the advisory judge's
  reference gains the `tasks` stage and its model-tier suggestion.
- `.agents/skills/write-tasks/SKILL.md` and its mirror
  `skills/write-tasks/SKILL.md` — Task authors choose `complexity`, which now
  decides where a Task runs, and run the judge's `tasks` stage; the version is
  raised with the text.

## What is not governed

The Go sources and tests under `internal/config`, `internal/lighttier`,
`internal/spec`, `internal/agent`, `internal/daemon`, `internal/judge` and
`internal/cli` that the Tasks declare, the mirrors
`skills/roundfix/references/runtime.md` and
`skills/roundfix/references/spec.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/configuration.md`, `docs/user-guide/commands/spec.md` and
`docs/references/model-selection.md` are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to `.roundfixrc.yml`, to any built-in or Recommended Profile, or
  to the Baseline source.
- No change to ADR-0235's refusal: every light model passes it, and no
  OpenAI, Anthropic or router model is ever a light model.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md`, `CHANGELOG.md` or the `### QA settlement` section of any
  skill.
- No test or Verification command reaches OpenRouter, OpenCode's providers or
  TypeSafe, and no test reads or writes the real `~/.roundfix`. Only the QA
  gate may run real light Tasks: at most two, on open models only, within
  US$5, with the cost recorded.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
