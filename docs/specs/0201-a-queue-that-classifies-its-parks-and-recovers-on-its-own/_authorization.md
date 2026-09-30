---
status: approved
granted: 2026-09-30
action: let a Spec name prerequisite Specs the queue owner waits for, park a conflicting Pull Request at once and resolve a conflict confined to declared derived paths by regeneration, resume an item the operator archived with a QA override after an environment-only partial, re-run once a check that failed outside the item's change, and give every park a class and a next command
consuming: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
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

# Approved authority for Spec 0201

On 2026-09-30 the maintainer approved Wave 8, "fila que classifica paradas e se
recupera": a Delivery Queue that classifies its parks and recovers on its own,
with the scope this Spec implements. The same day the maintainer said
"considere autorizado a ajustar todas as skills se necessário", which covers
the Roundfix skill's delivery reference and the version change every owned
skill edit carries. The governed set was measured with `GovernedPath` on the
authoring branch at `30e8504f`, through a `go test -overlay` probe that wrote
nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` is the canonical delivery
  reference of the Roundfix skill after Spec 0194 splits it. task_04 describes
  the prerequisites, the conflict park and its derived resolution, the
  operator-archived retry, the one check re-run and the Park Classes there, as
  the rule that the skill matches the shipped CLI requires.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  skill's two version fields, which task_04 raises because the skill's content
  changes. `skills/roundfix/SKILL.md` is rewritten only by the sanctioned
  mirror copy.

## What is not governed

`internal/delivery/`, `internal/cli/deliver.go`,
`internal/cli/deliver_workflow.go`, `internal/cli/deliver_test.go`,
`internal/spec/spec.go`, `internal/store/delivery.go`,
`internal/config/config.go`, their new test files,
`docs/user-guide/commands/deliver.md`, `docs/user-guide/configuration.md`,
`skills/roundfix/references/deliver.md` and
`skills/testdata/owned-skill-versions.json` are ordinary.
`internal/cli/cli_test.go`, `internal/speccheck/mechanical_test.go`,
`docs/references/coverage-record.json` and `.roundfixrc.yml` are governed and
are not touched: no help text changes, no existing top-level test is renamed or
removed, and this repository's derived-path declaration waits for a release.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

The owned-skill version record is rewritten by
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.

## Limits

- No change to archived Specs, existing QA Reports, the Run Database schema,
  the retry limit or any command's flags or help text.
- The queue owner never rebases, force-pushes or resolves a source conflict,
  and re-runs a failed check at most once.
- No Task or gate writes to the live Run Database under `~/.roundfix`, reaches
  GitHub, or runs `gh`; tests use temporary repositories and fake command
  runners.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
