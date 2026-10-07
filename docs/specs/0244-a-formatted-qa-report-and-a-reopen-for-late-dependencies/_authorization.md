---
status: approved
granted: 2026-10-07
action: format the files a QA step commits under the Spec's qa directory with the repository's configured Format Command, let roundfix reopen return a completed QA gate to pending over a Late Dependency, and describe both in the glossary, the guides and the Roundfix Skill's implement and settle references
consuming: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/settle.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0244

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". The maintainer's standing answer for the Governed Paths
each Spec declares is "Concedo". On 2026-10-07 the maintainer approved this
cycle's order, with this Spec first: "Pode seguir nessa ordem".

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, a temporary home, fake runners
and a shell-script formatter. Adopter repositories are read only, through the
Secondbrain mirror, and nothing from them is sent to an external judge.

The governed set was measured with `GovernedPath` on the authoring branch at
`66f2d85b`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/implement.md` — it describes the QA
  step's commit boundary, which now runs the Format Command, and the
  skill-sync rule requires the skill to ship with the behavior.
- `.agents/skills/roundfix/references/settle.md` — it states that reopen acts
  only when a dependency is no longer completed, which this Spec changes.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/config`, `internal/daemon`,
`internal/spec` and `internal/cli`, `CONTEXT.md`,
`docs/user-guide/configuration.md`,
`docs/user-guide/context-driven-development.md`,
`docs/user-guide/commands/reopen.md`, the skill reference mirrors under
`skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json` and
`docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`
are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows. This Spec adds a configuration key Roundfix reads; it configures
  no formatter for this repository.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, or reads or writes the real `~/.roundfix`.
- No change to archived Specs, existing QA Reports, `CHANGELOG.md` or the
  `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
