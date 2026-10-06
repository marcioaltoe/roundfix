---
status: approved
granted: 2026-10-05
action: let the Jev judge and implementation on open models each read their own Roundfix OpenRouter key first and the shared key second, record which variable each used, report each stage's keys in the Doctor Command, and describe the rule in the guides, the Roundfix Skill and the write-prd and write-techspec skills
consuming: 0234-an-openrouter-key-per-stage
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0234

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-05 the maintainer answered the cycle's
questions through AskUserQuestion: the scope "W e X", of which this Spec is X;
on creating one OpenRouter key per stage, each with its own monthly limit,
exposed as `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` and
`ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`, "Sim, eu crio"; and, for the Governed
Paths this Spec declares, "Concedo". The purpose the maintainer gave is to
measure each stage's cost exactly in OpenRouter's activity export, which
groups spend by API key. No live OpenRouter or TypeSafe call is authorized by
this record: authoring, tests, Verification and QA use fake transports and
read only the names and presence of environment variables.

The governed set was measured with `GovernedPath` on the authoring branch at
`9085540d`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/spec.md` — the repository's skill-sync
  rule requires a Pull Request that changes CLI behavior to ship the skill
  update, and this reference names the judge's keys and its output.
- `.agents/skills/roundfix/references/runtime.md` — it describes the Doctor
  Command and OpenCode's OpenRouter selections, which now read and report the
  stage keys.
- `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`
  and their mirrors `skills/write-prd/SKILL.md` and
  `skills/write-techspec/SKILL.md` — both forbid printing, storing or asking
  for "either Jev key" by name, and the judge now has a third.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/openrouterkey`, `internal/judge`,
`internal/cli` and the package that holds Spec 0233's implementation-key
helper, `docs/adr/0239-each-openrouter-stage-reads-its-own-roundfix-key-first.md`,
`docs/user-guide/commands/spec.md`, `docs/user-guide/commands/doctor.md`,
`docs/user-guide/configuration.md`, the skill reference mirrors
`skills/roundfix/references/spec.md` and
`skills/roundfix/references/runtime.md`, and
`skills/testdata/owned-skill-versions.json` are ordinary.

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
- No change to either ceiling's value, scope or sum, to the judge's
  recipients, or to which models an Agent Selection may name.
- No test, Verification command or QA row reaches OpenRouter or TypeSafe,
  reads a real key, or reads or writes the real `~/.roundfix`.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`,
  `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
