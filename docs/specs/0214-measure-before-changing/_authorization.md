---
status: approved
granted: 2026-10-01
action: add the read-only `roundfix runs causes` command and describe it in the Roundfix Skill, add a flag-gated test harness that re-measures the Task acceptance judgment through the judge's transports, and record three measurements with their recommendations, deleting no archived file
consuming: 0214-measure-before-changing
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runs.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0214

On 2026-09-30 the maintainer expressly authorized the skill files:
"considere autorizado a ajustar todas as skills se necessário". On
2026-10-01 the maintainer extended the unattended program to every planned
front, "Tudo, de A a F", of which this Spec is front D, "measure before
changing", and decided its scope: the history-root work measures and proposes
only, no archived file is deleted, and any removal waits for the maintainer's
explicit approval in a later change.

The same decisions bind the Task acceptance re-measurement: Jev is called
through OpenRouter first with `ROUNDFIX_OPENROUTER_API_KEY` (model
`jev-1.13`) and through TypeSafe directly with `ROUNDFIX_TYPESAFE_API_KEY`
(model `jev-1.13.0`) only when the first key is absent; the keys come from
the environment only and are never printed; the generic `OPENROUTER_API_KEY`
and `TYPESAFE_API_KEY` are never read; the data sent is only this
repository's Task files; and the monthly ceiling of US$5 that the Judge Log
enforces still holds.

The set was measured with `GovernedPath` on `5aa87d2b`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/runs.md` — the Roundfix Skill must
  describe the new `runs causes` command next to `runs list` and `runs show`.
  The text goes under a new heading, `### Why Verification failed`.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — a change
  to the Roundfix Skill's content raises both of its version fields, and the
  mirror's front matter is governed as well.

## What is not governed

Every Go source, test, embedded table and fixture under `internal/runcause`,
`internal/store`, `internal/cli` and `internal/judge` that this Spec's Tasks
create or change, the command reference `docs/user-guide/commands/runs.md`,
the mirror reference `skills/roundfix/references/runs.md`,
`skills/testdata/owned-skill-versions.json`, the reference documents and
records under `docs/references/`, and any Backlog Entry a rule calls for are
ordinary.

`internal/cli/cli_test.go`, `.roundfixrc.yml`, the skill manifests,
`Makefile`, the CI workflows and `go.mod` are governed and are not touched. No
Project Config key is added.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edit and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No new dependency in `go.mod`: the command and the harness use the
  standard library and the repository's existing SQLite driver.
- No change to the Makefile, the lint, formatter or test-runner
  configuration, or the CI workflows.
- No file under `docs/history/` is deleted, moved or rewritten, and no
  archived Spec, existing QA Report, `skills/_ownership.yml`, `CONTEXT.md` or
  `### QA settlement` section of any skill changes.
- No test, Verification command or QA row reaches OpenRouter, TypeSafe or the
  network, reads a real `ROUNDFIX_OPENROUTER_API_KEY`,
  `ROUNDFIX_TYPESAFE_API_KEY`, `OPENROUTER_API_KEY` or `TYPESAFE_API_KEY`, or
  writes under the real `~/.roundfix`. Only the harness, run explicitly with its
  flag by the measurement Task, sends requests, and only this repository's Task files.
- The `runs causes` command and the measurement steps open the real Run
  Database read-only and write nothing to it.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
