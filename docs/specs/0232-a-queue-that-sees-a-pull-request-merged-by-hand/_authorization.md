---
status: approved
granted: 2026-09-30
action: let a Delivery Retry record a parked item as merged when its recorded Pull Request was merged from the item branch or when the default branch already archives its Spec through a delivery commit outside the item branch, refuse a retry whose recorded Pull Request was closed without merging, and describe the rule in the deliver guide and the Roundfix Skill
consuming: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0232

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-05 the maintainer answered the cycle's
question through AskUserQuestion with the scope "U + medição do tier"; this
Spec is U. The advisory `roundfix spec judge` was allowed for this
repository's Spec artifacts only, run once on the finished Spec, its
suggestions reported and never used as a gate.

The governed set was measured with `GovernedPath` on the authoring branch at
`69462a0c`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` — the repository's
  skill-sync rule requires a Pull Request that changes CLI behavior to ship
  the skill update, and this reference describes `deliver retry`, its
  re-entry stages and its refusals.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/delivery`, `internal/store`,
`internal/worktree` and `internal/cli` that the Tasks declare,
`docs/user-guide/commands/deliver.md`,
`skills/roundfix/references/deliver.md` and
`skills/testdata/owned-skill-versions.json` are ordinary.

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
- No change to `deliver status`, `deliver resume`, the owner's pass or the
  post-merge cleanup's proofs, and no change to the Run Database schema.
- No test, Verification command or QA row reaches GitHub, and no test reads
  or writes the real `~/.roundfix`.
- No change to archived Specs, existing QA Reports, `CONTEXT.md`,
  `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
