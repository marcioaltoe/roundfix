---
status: approved
granted: 2026-09-30
action: run Settlement Checks (the repository Verification, the Spec Consistency findings a Task introduced and the authorization audit of the prospective Task commit) when a Task of a gated Task Graph settles, add the config key that turns the repository Verification at settlement off, and describe the stage in the Roundfix and write-tasks skills
consuming: 0190-a-task-settles-on-the-facts-its-gate-will-check
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0190

On 2026-09-30 the maintainer approved a three-wave program, answering "Três
ondas" to a structured question. Its first wave includes checking gate facts
when a Task settles. The same day the maintainer wrote "considere autorizado a
ajustar todas as skills se necessário", which authorizes the four skill files
below, and answered "Autorizar os dois" for the baseline source and guides and
for `.roundfixrc.yml`. This Spec uses only the skill authorization. The set was
measured with `GovernedPath` on the authoring branch at `9e439dbb`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the
  Roundfix skill describes what the Implement Command does when a Task
  settles and which config keys exist. It must state the Settlement Checks
  and `verification.repository_at_settlement` (task_04). `.agents/skills/` is
  canonical and `skills/` its mirror.
- `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md` — the
  stage creates an authoring obligation: every Task of a gated graph must
  leave the repository Verification green. The skill that authors Task Graphs
  must state it, or authored graphs fail at settlement (task_04).

## What is not governed

`internal/daemon/task_engine.go`, `internal/daemon/engine.go`,
`internal/speccheck/mechanical.go`, `internal/config/config.go`,
`internal/cli/implement.go`, the new file `internal/daemon/settlement_checks.go`,
the new test files this Spec's Tasks name under `internal/daemon/`,
`internal/speccheck/`, `internal/config/` and `internal/cli/`,
`docs/user-guide/commands.md`, `docs/user-guide/configuration.md` and
`CONTEXT.md` are ordinary. `internal/cli/cli_test.go`,
`internal/speccheck/mechanical_test.go`, `docs/references/coverage-record.json`,
`docs/agents/autonomous-work.md`,
`internal/baseline/assets/modules/autonomous-work.json` and `.roundfixrc.yml`
are governed and are not touched. No help text changes, no existing top-level
test is renamed or removed, and the new config key is not written to this
repository's Project Config.

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

- No change to the QA gate's precondition, mechanical stage or verdict rules,
  to the single Verification repair, to archived Specs or to existing QA
  Reports.
- No edit to the Makefile, to lint, formatter or test-runner configuration,
  to a CI workflow or to `go.mod`.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
