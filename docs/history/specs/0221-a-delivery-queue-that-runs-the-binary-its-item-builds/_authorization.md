---
status: approved
granted: 2026-09-30
action: let a repository declare its Roundfix item build in Project Config, have the queue owner build that binary in the item worktree and run the item's Implement, archive and review steps with it unless it would change the Run Database schema, add roundfix migrate --check, and describe this in the Roundfix Skill and the command guides
consuming: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/setup.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0221

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-01 the maintainer extended the program to
the next cycle, and on 2026-10-03, asked to authorize the cycle G, H and I,
answered "Continue". This Spec is item I: a delivery queue that runs the
binary its item builds.

The governed set was measured with `GovernedPath` on the authoring branch at
`305211a3`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/deliver.md` — the item build
  declaration, the item binary, the schema rule and the park.
- `.agents/skills/roundfix/references/setup.md` — `roundfix migrate --check`.

## What is not governed

The Go sources and tests under `internal/cli` and `internal/config` that the
Tasks declare, `skills/roundfix/references/deliver.md`,
`skills/roundfix/references/setup.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/migrate.md`,
`docs/user-guide/configuration.md` and ADR-0225 are ordinary.

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
  network, or writes under the real `~/.roundfix`; every Run Database a test
  reads lives in a disposable Roundfix Home.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
