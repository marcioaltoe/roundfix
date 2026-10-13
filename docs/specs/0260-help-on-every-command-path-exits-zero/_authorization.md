---
status: approved
granted: 2026-10-11
action: make --help and -h exit 0 with usage on qa-report accept and the baseline capabilities, skills and assets groups, sweep every command path for help in a repository test, fingerprint each Behavior Surface from its own help, and update the roundfix skill's settle and baseline references with the recorded version raise
consuming: 0260-help-on-every-command-path-exits-zero
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/settle.md
  - internal/cli/cli_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0260

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, and on skills: "considere autorizado a ajustar todas as skills
se necessário". On 2026-10-11, answering "Grants dos caminhos governados" for
cycle AE, the maintainer wrote "Concedo" for the Governed Paths each Spec
declares. Cycle AE item F is this Spec: `roundfix qa-report accept --help`
exits 0 with its usage, the help-exit test covers it, and the Behavior Surface
Record reads the subcommand's own help, within one implementation Task and
the QA gate.

The governed set was measured on the authoring branch at `7a8b7f70` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository. It ran against
every path task_01 declares, and five are governed.

## Why each governed path is unavoidable

- `internal/cli/cli_test.go` holds `TestRunCommandHelp`, the help-exit test
  the adopted Backlog Entry asks to cover the fixed command paths.
- `.agents/skills/roundfix/references/settle.md` says `accept` has no
  `--help`, which this Spec makes false; ADR-0256's authoring rule and
  `docs/agents/cli.md` require the skill to change with the command.
- `.agents/skills/roundfix/references/baseline.md` covers the `baseline`
  Behavior Surfaces whose source `internal/cli/baseline_profile.go` changes,
  and gains the group-help sentence.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  owned skill's version fields, which the record command raises with any
  content change.

## What is not governed

These paths are ordinary: `internal/cli/qa_report.go`,
`internal/cli/baseline_profile.go`, `internal/cli/command_help_test.go`,
`internal/docscontract/skill_coverage_test.go`,
`skills/roundfix/references/settle.md`,
`skills/roundfix/references/baseline.md` and
`skills/testdata/owned-skill-versions.json`.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors after the approved
skill edit and adds no source path.

```yaml
command: make skills-sync
```

task_01 also runs the owned-skill record command, whose package installs no
suite guard, so it needs no declared output list here:
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
rewrites the version lines of `.agents/skills/roundfix/SKILL.md` and its
mirror and `skills/testdata/owned-skill-versions.json`.

No Baseline module changes, so neither `make baseline-digests` nor the module
version record command is declared, and the Behavior Surface Record keeps its
bytes.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, or lint, formatter or test-runner configuration.
- No change to the user guide, `docs/references/behavior-surfaces.json`,
  `docs/references/skill-coverage.json`, archived Specs or Archive Records.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
