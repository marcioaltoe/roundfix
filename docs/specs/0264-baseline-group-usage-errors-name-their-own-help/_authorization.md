---
status: approved
granted: 2026-10-12
action: make the baseline dispatcher's usage errors name the group once and point at the help of the typed command path, guard every group's usage error in a repository test, and describe the bare-group usage error in the roundfix skill's baseline reference with the recorded version raise
consuming: 0264-baseline-group-usage-errors-name-their-own-help
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0264

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, and on skills: "considere autorizado a ajustar todas as skills
se necessário". On 2026-10-11, answering "Grants dos caminhos governados" for
cycle AE, the maintainer wrote "Concedo" for the Governed Paths each Spec
declares. On 2026-10-12 the maintainer approved cycle AF, answering "Sim, como
no AE" to "As regras de operação do AE valem para o AF?", which carries the
same operations and the same grant for the Governed Paths each AF Spec
declares. Cycle AF item C is this Spec: the `baseline` group usage error.

The governed set was measured on the authoring branch at `3820bdd4` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository. It ran against
every path task_01 and task_02 declare, and three are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/baseline.md` covers the `baseline`
  Behavior Surfaces whose source `internal/cli/baseline_profile.go` changes;
  ADR-0256's authoring rule and `docs/agents/cli.md` require the skill to
  change with the command, and it gains the bare-group usage-error sentence.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  owned skill's version fields, which the record command raises with any
  content change.

## What is not governed

These paths are ordinary: `internal/cli/baseline_profile.go`,
`internal/cli/baseline_group_usage_test.go`,
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
  `docs/references/skill-coverage.json`, `internal/cli/cli_test.go`,
  archived Specs or Archive Records. `internal/cli/command_help_test.go`
  changes only in the expected standard error of the bare and unknown
  `capabilities`, `skills` and `assets` rows that task_02 Requirement 5
  names.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
