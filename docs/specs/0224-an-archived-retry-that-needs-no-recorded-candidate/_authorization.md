---
status: approved
granted: 2026-10-04
action: let a Delivery Retry resume an item the operator archived with the QA Archive Override from the Implement start head of its recorded Run whatever its park, park a zero-finding QA partial with any environment-blocked row as qa-environment-partial, and describe both in the Roundfix Skill and the deliver guide
consuming: 0224-an-archived-retry-that-needs-no-recorded-candidate
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

# Approved authority for Spec 0224

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-04, asked which Spec to deliver next, the
maintainer answered "deliver M": this Spec, an archived retry that needs no
recorded candidate.

The governed set was measured with `GovernedPath` on the authoring branch at
`6ea9e1e9`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/deliver.md` — the two retry and park
  rules the Delivery Queue changes.

## What is not governed

`internal/delivery/engine.go`, `internal/delivery/operator_archive_any_park_test.go`,
`internal/cli/deliver_workflow.go`, `internal/cli/deliver_archived_retry_test.go`,
`internal/cli/deliver_operator_archive_test.go`,
`skills/roundfix/references/deliver.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/deliver.md` and ADR-0229 are ordinary.

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
  Park Class names, blocker names, retry limits, `CONTEXT.md` or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`; every Run Database a test
  reads lives in a disposable Roundfix Home.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
