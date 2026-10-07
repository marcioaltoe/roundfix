---
status: approved
granted: 2026-10-06
action: add roundfix history sanitize, which converts the existing history in reviewed batches after a history-full tag, and describe it in the command guide, the Roundfix skill, one Baseline clause and the glossary
consuming: 0243-a-history-that-holds-only-records
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0243

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources and guides the
maintainer said "Autorizar os dois".

On 2026-10-06 the maintainer decided the history clean-up:

- "Não quero nada no histórico que não seja relevante para o secondbrain".
- The clean-up "deve atuar também no que já existe no diretório". This Spec
  delivers the command that does it.
- "Aplicar em lotes via PR (Recommended)": after a `history-full` tag, the
  existing history is applied in batches of Specs, one Pull Request per
  batch, each with `make verify` green and each revertible.
- For the Governed Paths the clean-up declares, the maintainer said
  "Concedo".
- Jev may advise through the judge key, on this repository's artifacts
  only.

These decisions are the explicit approval ADR-0215 requires before existing
history is removed or rewritten, and ADR-0248 records them. They apply to
the batches the operator runs with `roundfix history sanitize --apply` after
this Spec merges. They do not authorize this Spec's Tasks to delete, move or
rewrite anything under this repository's `docs/history`, to create or push
the `history-full` tag, or to edit `.secondbrain-export`.

The governed set was measured with `GovernedPath` on the authoring branch at
`b07a080e`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. 8 of them are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` — the
  skill-sync rule requires a Pull Request that adds a command to ship the
  Roundfix skill's description of it, and the skill index names the
  command. An owned skill's content changes only with its version.
- `internal/baseline/assets/modules/context-workflow.json` — the clause that
  says where retired documentation goes gains the sentence on a history
  sanitize. `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`,
  the formatter golden `docs-layout.md` and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` are
  its sanctioned derived files.

## What is not governed

These are ordinary:

- the Go sources and tests under `internal/spec`, `internal/cli`,
  `internal/speccheck`, `internal/docscontract` and `internal/baseline`;
- `internal/baseline/testdata/**`;
- the user guides, `CONTEXT.md`, `CHANGELOG.md` and
  `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`;
- the skill reference mirror `skills/roundfix/references/archive.md` and
  `skills/testdata/owned-skill-versions.json`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors and the Baseline's
derived files. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

`make baseline-digests` regenerates the catalog snapshots, the plan goldens,
the formatter goldens and the profile digest.
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden or generated guide is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No deletion, move or rewrite of any file already under `docs/history`, no
  `history-full` tag and no change to `.secondbrain-export` inside this
  Spec's Tasks or QA gate. Those belong to the operator's batches.
- No test, Verification command or QA row reaches a provider, reads a
  credential, or reads or writes the real `~/.roundfix`. Advice is exercised
  only against a fake transport.
- No change to the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
