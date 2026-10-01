---
status: approved
granted: 2026-09-30
action: make the Jev Router an opt-in OpenCode selection that only Project Config can name, gate every routed prompt on the shared monthly Jev ceiling and record its cost in the Judge Log, describe it in the Roundfix Skill, and measure it on replayed docs and chore Tasks
consuming: 0218-the-jev-router-on-docs-and-chore-tasks
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0218

On 2026-09-30 the maintainer expressly authorized the skill files: "considere
autorizado a ajustar todas as skills se necessário". On 2026-10-01 the
maintainer authorized the next cycle of unattended work with "Tudo, de A a
F". This Spec is the second half of item F, the experiment that tries the Jev
Router for cheap agent categories. The same day the maintainer bound its
design: agent prompts carrying this repository's code and diffs only, never
adopters' or the Secondbrain's content, may go to OpenRouter on
`ROUNDFIX_OPENROUTER_API_KEY`; `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY`
are never read; keys come from the environment only and are never printed;
spending stays under the US$5 monthly ceiling shared with all Jev use, of
which about US$0.40 was spent; and the experiment delivers measurement and an
opt-in entry, never a default change.

The set was measured with `GovernedPath` on `5aa87d2b`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/runtime.md` — the Roundfix Skill must
  describe the routed selection, its Project Config rule, its key and its
  ceiling, because the Spec changes what a configuration may name. The text
  goes under a new heading, `### Jev Router`.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — a change
  to the Roundfix Skill's content raises both of its version fields
  (ADR-0189), and the mirror's front matter is governed as well.

## What is not governed

Every Go source and test file under `internal/agent`, `internal/config`,
`internal/jevrouter` and `internal/daemon` that this Spec's Tasks create or
change, the guide `docs/user-guide/configuration.md`, the mirror reference
`skills/roundfix/references/runtime.md`,
`skills/testdata/owned-skill-versions.json`, and the measurement record under
this Spec's `measurement/` directory are ordinary.

`internal/cli/cli_test.go`, `.roundfixrc.yml`, the skill manifests, the
Makefile, the CI workflows and `go.mod` are governed and are not touched. No
Project Config key is added.

## Sanctioned regeneration

The repository-owned command resolves the generated mirror. This declaration
records the regeneration that follows the approved skill edit and adds no
source paths.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`: the key reader uses the standard library.
- No change to the Recommended Profile, the built-in profiles, this
  repository's `.roundfixrc.yml` or any User Config.
- No test, Verification command or QA row reaches OpenRouter or the network,
  reads a real `ROUNDFIX_OPENROUTER_API_KEY` or `OPENROUTER_API_KEY`, or
  writes under the real `~/.roundfix`.
- Only the measurement Task sends prompts to OpenRouter, from scratch clones
  of this repository, under the shared ceiling the gate enforces.
- No change to archived Specs, `CONTEXT.md` or the `### QA settlement`
  section of any skill.
