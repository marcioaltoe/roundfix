---
status: approved
granted: 2026-09-30
action: add `cursor` as an opt-in ACP Runtime whose selections name Cursor's advertised model value, check the maintainer's Cursor login without performing it, describe the runtime in the Roundfix Skill, and measure Grok through it
consuming: 0217-a-cursor-runtime-to-measure-grok-on
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

# Approved authority for Spec 0217

On 2026-09-30 the maintainer expressly authorized the skill files: "considere
autorizado a ajustar todas as skills se necessário". On 2026-10-01 the
maintainer authorized the next cycle of unattended work with "Tudo, de A a
F". This Spec is the first half of item F, the experiment that runs Grok
through the Cursor runtime. The same day the maintainer bound its design:
Cursor runs on the maintainer's existing local login; a login that needs a
person stops the work with a named blocker; Roundfix never automates a login
or handles a credential; and the experiment delivers measurement and an
opt-in runtime entry, never a default change.

The set was measured with `GovernedPath` on `5aa87d2b`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/runtime.md` — the Roundfix Skill must
  describe the `cursor` runtime, its model-value rule and its login rule,
  because the Spec changes CLI behavior (the repository's skill sync rule).
  The text goes under a new heading, `### Cursor`.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — a change
  to the Roundfix Skill's content raises both of its version fields
  (ADR-0189), and the mirror's front matter is governed as well.

## What is not governed

Every Go source, test and fixture file under `internal/config`,
`internal/agent` and `internal/cli` that this Spec's Tasks create or change,
the guides under `docs/user-guide/`, the mirror reference
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

- No new dependency in `go.mod`.
- No change to the Recommended Profile, the built-in profiles or this
  repository's Project Config.
- No test, Verification command or QA row reaches Cursor or the network, or
  runs a real `cursor-agent`.
- Roundfix never runs `cursor-agent login` or `logout`, never passes a Cursor
  key or token, never reads `CURSOR_API_KEY` or `CURSOR_AUTH_TOKEN`, and never
  prints or stores the account.
- Only the measurement Task sends prompts to Cursor, from scratch clones of
  this repository, on the maintainer's machine and existing login.
- No change to archived Specs, `CONTEXT.md` or the `### QA settlement`
  section of any skill.
