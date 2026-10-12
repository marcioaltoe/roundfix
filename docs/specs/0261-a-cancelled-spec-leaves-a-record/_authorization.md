---
status: approved
granted: 2026-10-11
action: add the roundfix cancel command, the cancelled Archive Record disposition and Source Return, make the Supersede Command, deliver start and deliver retry refuse a cancelled Spec, and describe cancellation in the Baseline docs layout guide, the Roundfix and archive-spec skills, the user guides and the glossary
consuming: 0261-a-cancelled-spec-leaves-a-record
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/deliver.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0261

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois".

On 2026-10-11 the maintainer decided cycle AE. For the cancellation item the
decision was "Comando com registro": `roundfix cancel <slug> --reason <text>`
writes an Archive Record with the disposition `cancelled` that keeps the
reason, the date and the completed and pending Tasks, returns the adopted
sources to the backlog under the adoption contract, refuses while a Run or a
Delivery Queue item is active for the Spec, and is described by the
archive-spec skill and the docs guide. On 2026-10-12 the maintainer split item
B, giving this Spec the cancellation and Spec 0256 the squash survival. To the
request for grants, "Grants dos caminhos governados", the maintainer answered
"Concedo" for the Governed Paths each Spec declares.

No live provider call is authorized by this record. Tests, Verification and QA
use temporary repositories and homes, the embedded catalog and this
repository's own files. No adopter repository is read or written.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a8b7f70`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against the 44 files the Tasks
declare, and 11 of them are governed, all in task_02. The derived Baseline
files among them were measured by an authoring prototype in a scratch clone,
which applied the clause and the version raises and ran the sanctioned
commands.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/spec-workflow.json` holds the docs layout
  clause the adopted entry asks to describe the cancellation. The Module
  Version Record step writes its version line.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and
  the Standard TypeScript formatter golden `docs/agents/docs-layout.md` are
  rewritten by `make baseline-digests` from the module edit. The prototype
  measured exactly these.
- `docs/agents/docs-layout.md` and `docs/agents/setup-context.json` are
  rewritten by this repository's Managed Refresh.
- `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/deliver.md` and
  `.agents/skills/archive-spec/SKILL.md`, with the mirrors
  `skills/roundfix/SKILL.md` and `skills/archive-spec/SKILL.md`: the routing
  row that indexes the new command for the Skill Coverage Map, the Cancel
  Command section, the delivery refusals and the cancelling section the
  adopted entry asks for. An owned skill's content changes only with its
  version, and the record command raises both version fields.

## What is not governed

These paths are ordinary: `CONTEXT.md`, every file under `docs/user-guide/`,
`docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json`,
`internal/baseline/module-versions.json`, `internal/baseline/testdata/**`,
`skills/roundfix/references/archive.md`, `skills/roundfix/references/deliver.md`,
`skills/testdata/owned-skill-versions.json`, every production and test Go
file the Tasks declare, and the new tests.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
skill mirrors. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes the changed module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make skills-sync
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden or generated guide is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No change to the Source Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, or any module other than
  `spec-workflow`.
- No change to the `### QA settlement` section of any skill, to archived
  Specs, existing Archive Records, existing history entries or `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
