---
status: approved
granted: 2026-09-30
action: add the advisory `roundfix spec judge` command backed by TypeSafe Jev, describe it in the Roundfix Skill, and have the write-prd and write-techspec skills run it and answer what it raises
consuming: 0205-an-advisory-judge-for-spec-authoring
paths:
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0205

On 2026-09-30 the maintainer reopened the Jev front and decided this Spec's
content: a dedicated command that the `write-prd` and `write-techspec` skills
call, advisory only and never a gate; the pinned model `jev-1.13.0`; a JSONL
record of every call; a spending ceiling of US$5 per calendar month that stops
calling and fails open; requests carrying only this repository's Spec
artifacts; the key read only from `TYPESAFE_API_KEY`; and English artifacts
only. The same day the maintainer expressly authorized the skill files:
"considere autorizado a ajustar todas as skills se necessário".

On 2026-10-01, before delivery, the maintainer amended the transport and the
key: "Podemos consumir através do openrouter. Estou muito disposto a seguir
por esse caminho", which authorizes sending the same Spec artifacts to Jev
through OpenRouter, and "Gosto de ter chaves separadas para determinar o real
custo de cada trabalho/projeto no openrouter", which gives the OpenRouter
transport its own key, `ROUNDFIX_JEV_OPENROUTER_API_KEY`. The command calls
OpenRouter's System One API when that key is set and TypeSafe directly only
when it is absent and `TYPESAFE_API_KEY` is set; the generic
`OPENROUTER_API_KEY` is never read. The pinned version stays Jev 1.13,
requested as `jev-1.13` on OpenRouter and `jev-1.13.0` on TypeSafe. The
amendment adds no Governed Path.

The set was measured with `GovernedPath` on `5f182757`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`
  — the two authoring skills gain the heading that runs the command and tells
  the authoring model how to answer it, and each raises both version fields.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — a change
  to the Roundfix Skill's content raises both of its version fields, and the
  mirror's front matter is governed as well.
- `.agents/skills/roundfix/references/spec.md` — the Roundfix Skill must
  describe the new command. Spec 0194 creates this per-command file, and this
  Spec is delivered after it. The text goes under a new heading,
  `### Advisory judge`.

## What is not governed

Every Go source and test file under `internal/judge` and `internal/cli` that
this Spec's Tasks create or change, including the embedded
`internal/judge/questions.json`, the command reference file
`docs/user-guide/commands/spec.md`, the mirror reference
`skills/roundfix/references/spec.md` and
`skills/testdata/owned-skill-versions.json`, are ordinary.

`internal/cli/cli_test.go`, `.roundfixrc.yml`, the skill manifests,
`Makefile`, the CI workflows and `go.mod` are governed and are not touched. No
Project Config key is added.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No new dependency in `go.mod`: the client uses the standard library.
- No change to the Makefile, the lint, formatter or test-runner
  configuration, or the CI workflows.
- No change to archived Specs, existing QA Reports, `skills/_ownership.yml`,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches OpenRouter, TypeSafe or the
  network, reads a real `ROUNDFIX_JEV_OPENROUTER_API_KEY`, `TYPESAFE_API_KEY`
  or `OPENROUTER_API_KEY`, or writes under the real `~/.roundfix`.
- No data leaves the machine except, when the command runs with a Jev key, the
  Spec artifacts ADR-0201 names, sent to OpenRouter or to TypeSafe.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
