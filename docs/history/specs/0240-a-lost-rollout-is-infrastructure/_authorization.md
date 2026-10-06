---
status: approved
granted: 2026-10-06
action: treat a lost Codex rollout as runtime infrastructure, recover the Task in a new Agent Session or on its Fallback Selection without repair or a counted queue retry, and describe it in the glossary, the guides and the Roundfix Skill's profiles and deliver references
consuming: 0240-a-lost-rollout-is-infrastructure
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/profiles.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0240

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-06 the maintainer answered the cycle's
questions through AskUserQuestion:

- scope: "Tudo, na ordem sugerida", of which this Spec is second, for v0.51.0;
- the QA fallback to Claude when Codex loses the rollout before the report:
  "Sim, automaticamente". It may use the Claude subscription without spending a
  retry, recorded in the report;
- for the Governed Paths this Spec declares: "Concedo".

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use fake runners, a fake acpx and a fake ACP adapter under
a temporary home.

The governed set was measured with `GovernedPath` on the authoring branch at
`53d9b2d2`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/profiles.md` — it says that once Agent
  work started no session-loss failure starts a replacement session, which this
  Spec changes, and the skill-sync rule requires the skill to ship with the
  behavior.
- `.agents/skills/roundfix/references/deliver.md` — it lists the Park Classes
  and their blockers, and this Spec adds `runtime-infrastructure` with an
  uncounted retry.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields by the record command.

## What is not governed

The Go sources and tests under `internal/agent`, `internal/daemon`,
`internal/delivery`, `internal/store` and `internal/cli`, `CONTEXT.md`,
`docs/user-guide/commands/deliver.md`, `docs/user-guide/configuration.md`,
`docs/user-guide/usage.md`, the skill reference mirrors under
`skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json` and
`docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md` are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to acpx arguments, TTLs, `CODEX_HOME`, the review command or the
  pre-PR review, and no fallback after a First Handoff.
- No test, Verification command or QA row reaches a provider, starts a real
  Codex or Claude session, reads a credential, or reads or writes the real
  `~/.roundfix`, `~/.acpx` or `~/.codex`.
- No upstream issue is filed; the report text stays in the Spec's references.
- No change to archived Specs, existing QA Reports, `CHANGELOG.md` or the
  `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
