---
status: approved
granted: 2026-09-30
action: add the Profile Deviation key, the read-only recommendation comparison with `roundfix profiles check` and its `--apply`, the Doctor line and the `roundfix upgrade` notice, and describe them in the Roundfix skill
consuming: 0196-a-notice-when-profiles-fall-behind
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0196

On 2026-09-30 the maintainer approved the program this Spec belongs to
("Três ondas"), whose block D asks for a notice on `roundfix upgrade` with an
option to adjust the configuration, and expressly authorized the skill files:
"considere autorizado a ajustar todas as skills se necessário". This Spec is
the second half of that block. Spec 0189 holds the first half and is delivered
before it.

The set was measured with `GovernedPath` on `9e439dbb`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the skill
  is where an agent learns a Roundfix command. It must name the Profile
  Deviation (task_01), `roundfix profiles check` (task_02), `--apply`
  (task_03), and the Doctor line and the recommendation notice (task_04).
  `.agents/skills/` is canonical and `skills/` its mirror. All of the text
  goes under one new heading, `### Recommendation check`.

## What is not governed

`internal/config/profiles.go`, `internal/config/profile_config.go`,
`internal/config/recommendation_check.go`, `internal/cli/profiles.go`,
`internal/cli/profiles_check.go`, `internal/cli/profiles_configure.go`,
`internal/cli/doctor.go`, `internal/cli/health.go`, `internal/cli/upgrade.go`,
`internal/cli/cli.go`, `internal/cli/doctor_test.go`,
`internal/cli/upgrade_test.go`, the new test files,
`docs/user-guide/commands.md`, `docs/user-guide/configuration.md`,
`docs/user-guide/usage.md` and `CONTEXT.md` are ordinary.

`internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`,
`.roundfixrc.yml`, the skill manifests, `Makefile`, the CI workflows and
`go.mod` are governed and are not touched. The help text changes keep every
string `TestRunCommandHelp` requires, so `internal/cli/cli_test.go` and
`docs/references/coverage-record.json` do not move.

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

- No profile in `.roundfixrc.yml` changes, and no deviation is added to it.
- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, or the CI workflows.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `skills/_ownership.yml` or the `### QA settlement` section of any skill.
- No test reaches a real ACP adapter, a provider, a release host or the
  network, and the live Run Database under `~/.roundfix` is never opened for
  writing.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
