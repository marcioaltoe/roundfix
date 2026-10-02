---
status: approved
granted: 2026-09-30
action: add coded gh, git, remote, toolchain and environment readiness lines to Doctor and Setup, refuse deliver start when this machine cannot publish, report and restore required upstream skills that trail their Setup Snapshot, add the verification.tools Project Config key, declare this repository's tools and derived paths in its Project Config, and describe this in the Roundfix Skill and the command guides
consuming: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
paths:
  - .roundfixrc.yml
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0215

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, under the broad autonomy granted on 2026-09-29 for authoring,
corrective work and delivery through merge. The maintainer said of the skills
"considere autorizado a ajustar todas as skills se necessário", and answered
"Autorizar os dois" to a structured question that named the Baseline source,
the guides under `docs/agents/` and `.roundfixrc.yml`. On 2026-10-01 the
maintainer extended the program to the next cycle, "Tudo, de A a F", with this
Spec as wave C, and asked for it in these words: "é importante que a parte do
roundfix de setup verifique e reconheça os requerimentos para que possa rodar
sem problemas".

The governed set was measured with `GovernedPath` on the authoring branch at
`18ef15eb`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.roundfixrc.yml` — this repository's Project Config gains its
  `verification.tools` list and its `delivery.derived_paths` declarations.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/setup.md` — the new Doctor and Setup
  lines, codes and the `warn` status.
- `.agents/skills/roundfix/references/deliver.md` — the `deliver start`
  refusal.
- `.agents/skills/roundfix/references/baseline.md` — the managed refresh
  restores a trailing upstream skill.

## What is not governed

The Go sources and tests under `internal/cli`, `internal/config`,
`internal/agent`, `internal/baseline` and `skills` that the Tasks declare,
the mirrors `skills/roundfix/references/setup.md`,
`skills/roundfix/references/deliver.md` and
`skills/roundfix/references/baseline.md`,
`skills/testdata/owned-skill-versions.json`, the command guides under
`docs/user-guide/`, and ADR-0220 and ADR-0221 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirror. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to `internal/cli/cli_test.go`, the Baseline assets, archived
  Specs, existing QA Reports, `CONTEXT.md` or the `### QA settlement` section
  of any skill.
- No test, Verification command or QA row reaches GitHub or the network,
  reads a real token or key value, or writes under the real `~/.roundfix`;
  `gh` and `git` are replaced by scripted runners and fake executables.
- Roundfix never installs a tool, logs in or out of `gh`, or edits the user's
  shell environment or Git configuration.
- No restore of this repository's own `.agents/skills/` trees.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
