---
status: approved
granted: 2026-09-30
action: match an archived retry's Run by repository, wait on GitHub's merge state, return a base-branch policy merge refusal to checking, print a refused retry as a refusal, drop a missing NODE_OPTIONS preload from the agent environment with a notice, and describe this in the Roundfix Skill and the command guides
consuming: 0211-a-delivery-queue-that-finishes-without-intervention
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/runtime.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0211

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, under the broad autonomy granted on 2026-09-29 for authoring,
corrective work and delivery through merge, and said of the skills
"considere autorizado a ajustar todas as skills se necessário". On 2026-10-01
the maintainer extended the program to the next cycle, "Tudo, de A a F", one
minor release per wave, with this Spec as part of wave A: a delivery queue
that finishes without operator intervention. For the dead preload the
maintainer preferred dropping it from the agent environment with a named
notice over refusing, and ruled out editing the user's environment otherwise.

The governed set was measured with `GovernedPath` on the authoring branch at
`30cd147c`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/deliver.md` — the retry, checking and
  merge behavior and the refusal wording.
- `.agents/skills/roundfix/references/runtime.md` — the agent environment and
  its notice.

## What is not governed

The Go sources and tests under `internal/cli`, `internal/delivery` and
`internal/agent` that the Tasks declare, `skills/roundfix/references/deliver.md`,
`skills/roundfix/references/runtime.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`
and ADR-0211 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirror. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  Park Class names, retry limits, `CONTEXT.md` or the `### QA settlement`
  section of any skill.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`; GitHub is replaced by the
  scripted command runner and fake boundaries.
- The user's shell environment and configuration are never edited.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
