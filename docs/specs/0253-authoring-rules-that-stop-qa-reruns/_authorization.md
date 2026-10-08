---
status: approved
granted: 2026-10-08
action: make outside evidence reachable before a Run, require whole-transcript tests and declared cross-Spec prerequisites, add the parallel re-check, the corrective needs edge, the wording sweep and a hermetic Go test clause, write the Roundfix-only authoring rules into the repository guide and the authoring skills, and make a retirement write reduced history and delete reviews and handoffs
consuming: 0253-authoring-rules-that-stop-qa-reruns
paths:
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/concrete-contracts.md
  - docs/agents/autonomous-work.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - docs/agents/specific-repository.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/go.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0253

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois". The standing answer
for the Governed Paths each Spec declares is "Concedo".

On 2026-10-08 the maintainer asked: "Atualizar e verificar o baseline de regras
para que tenhamos o máximo de performance e resultado com o roundfix e
congruente com as mudanças". The maintainer approved the route "Auditoria e
depois Spec (Recommended)". The audit's Specs are delivered without further
consultation within the bound of four implementation Tasks, with grants for
the governed paths each Spec declares, and a `qa_override` is approved for a
QA partial caused only by the environment. The audit is committed as
`docs/references/2026-10-08-baseline-audit.md`. This Spec is its Spec B, plus
finding B16.

One operator decision, taken within that grant on 2026-10-08, shapes task_04:
retirements write history directly in reduced form. A retired Finding or
Backlog Entry keeps its front matter, title, first paragraph and the commit
holding its full text, and retired Review Artifacts and handoffs are not moved
into `docs/history` at all. It follows ADR-0248 and the maintainer's "Não quero
nada no histórico que não seja relevante para o secondbrain".

No live provider call is authorized by this record. Tests, Verification and QA
use the embedded catalog, temporary repositories and this repository's own
files. No adopter repository is read or written.

The governed set was measured with `GovernedPath` on the authoring branch at
`b16ca020`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against every file the Tasks
declare, and 18 of them are governed. The derived files among them were
measured by an authoring prototype in a scratch clone, which applied the three
module changes and ran the sanctioned commands.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/spec-workflow.json`,
  `autonomous-work.json`, `go.json` and `context-workflow.json` hold the
  clauses, rule and guide versions this Spec extends or adds. The Module
  Version Record step writes their version lines.
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and
  the three Standard TypeScript formatter goldens are rewritten by
  `make baseline-digests` from the module edits. The prototype measured
  exactly these.
- `docs/agents/spec-routing.md`, `autonomous-work.md`, `docs-layout.md` and
  `setup-context.json` are rewritten by this repository's Managed Refresh.
- `docs/agents/specific-repository.md` gains the suite guard rule and the
  Roundfix-only authoring rules (finding B12 and the briefing section).
- `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`
  and `.agents/skills/write-techspec/references/concrete-contracts.md`, with
  the mirrors `skills/write-tasks/SKILL.md` and `skills/write-techspec/SKILL.md`:
  the whole-transcript, reachable-evidence, prerequisite and backtick rules. An
  owned skill's content changes only with its version, and the record command
  raises both version fields.

## What is not governed

These paths are ordinary: `CONTEXT.md`, `docs/agents/go.md`,
`docs/user-guide/context-driven-development.md`,
`docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`,
`internal/baseline/module-versions.json`, `internal/baseline/testdata/**`,
`skills/write-techspec/references/concrete-contracts.md`,
`skills/testdata/owned-skill-versions.json`, and the new and changed tests
under `internal/baseline`.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
skill mirrors. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes each changed module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/go.json
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
  `internal/baseline/assets/retention/**`, any production Go file, or any
  module other than the four named.
- No change to the `### QA settlement` section of any skill, to archived
  Specs, existing Archive Records, existing history entries or `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
