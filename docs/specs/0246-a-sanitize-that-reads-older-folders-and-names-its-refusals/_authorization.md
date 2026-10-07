---
status: approved
granted: 2026-10-07
action: let roundfix history sanitize read Legacy Archive Folders leniently, record a folder archived with a failing QA and no override under the failed-qa disposition, and list every Refused Unit while a batch converts the others, and describe all three in the glossary, the history command reference and the Roundfix Skill's archive reference
consuming: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0246

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". The maintainer's standing answer for the Governed Paths
each Spec declares is "Concedo".

On 2026-10-07 the maintainer decided the scope of this Spec from the adopted
Backlog Entry:

- On the three repairs: "Sim, as três frentes". These are the lenient reading of
  Legacy Archive Folders (projection rows outside the graph and retired Task
  types, for legacy conversion only, with the current schema still enforced for
  active Specs), a new disposition for a Spec archived with a failing QA and no
  override, and a plan that lists every refused unit with its reason instead of
  aborting.
- On a batch that meets a unit it cannot convert: "Pular e seguir". The batch
  applies the valid units, leaves the refused ones untouched and lists each
  with its reason.
- Grants: the Governed Paths this Spec declares, and read-only access to the
  adopter folders `docs/history/specs/{0001,0004,0005,0036}-*` to build
  faithful but minimal synthetic fixtures. No adopter content is copied beyond
  the structural shape needed, and nothing from the adopter reaches Jev, so
  `roundfix spec judge` is not run for this Spec.
- A `qa_override` for an environment-only partial is standing.

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, a temporary home and
synthetic legacy folders.

The governed set was measured with `GovernedPath` on the authoring branch at
`405274c6`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/archive.md` describes the History
  Sanitize Command. The skill-sync rule requires the skill to ship with the
  new refusal, lenient reading and `failed-qa` behavior.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  are included because an owned skill's content changes only with its version,
  and the record command raises both front-matter fields.

## What is not governed

These paths are ordinary:

- The Go sources and tests under `internal/spec` and `internal/cli`.
- `CONTEXT.md` and `docs/user-guide/commands/history.md`.
- The skill reference mirror `skills/roundfix/references/archive.md`.
- `skills/testdata/owned-skill-versions.json`.
- `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, the CI workflows
  or any Baseline module.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, reads or writes the real `~/.roundfix`,
  or writes to the adopter repository.
- No change to archived Specs, existing Archive Records, `CHANGELOG.md` or the
  `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
