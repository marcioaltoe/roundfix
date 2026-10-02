---
status: approved
granted: 2026-09-30
action: add three mandatory Baseline clauses on grouping sources that share a context into one Spec, the grouping bound and extending an open record before minting, with their Source Baseline rows; teach them in the write-idea, write-prd and write-techspec skills; and add the advisory grouping question to roundfix spec judge, described in the Roundfix Skill
consuming: 0209-sources-that-share-a-context-share-a-spec
paths:
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/source-baselines/index.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - docs/agents/spec-routing.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0209

On 2026-09-30 the maintainer approved the program whose waves cover the
Baseline and the Jev judge, and expressly authorized editing the Baseline
source and its generated guides, answering "Autorizar os dois" to a
structured question that named `internal/baseline/assets/` and the guides
under `docs/agents/`, and said of the skills "considere autorizado a ajustar
todas as skills se necessário". The same day the maintainer authorized
sending to Jev excerpts of "PRD, TechSpec, Task files, ADRs, findings e
Backlog Entries DESTE repositório", and nothing from other repositories.

On 2026-10-01 the maintainer decided this Spec's content:

> "podemos ter em uma mesma spec um ou mais inbox/backlogs ou findings não é
> obrigatório e nem desejavel uma spec para um finding com contexto similar ou
> complementar a outro. O cuidado é não ter specs desnecessáriamente complexas
> e grandes, o uso do jev pode ajudar a decidir pela reunião de inbox, backlogs
> e findings em uma só spec. Isso deve fazer parte do baseline de roundfix e
> parte das regras do roundfix."

> "podemos alterar e adicionar e atualizar um finding ou backlog já existente
> e ainda não implementado para que não seja necessário criar outro arquivo"

The governed set was measured with `GovernedPath` on `c3be3bc9`, through a
`go test -overlay` probe that wrote nothing to the repository, against the
files a disposable rehearsal of task_01 and of the skill edits rewrote and the
files task_02 and task_03 declare.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json` — the three
  clauses and the six raised versions.
- `internal/baseline/assets/source-baselines/index.json`, the Standard
  TypeScript Monorepo Source Baseline `baseline.json`, `manifest.json` and its
  corpus `spec-routing.md` and `docs-layout.md` — the catalog refuses a
  profile clause without its Source Baseline row; the regeneration rewrites
  the identity record and offsets.
- The two Standard TypeScript Monorepo goldens and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` —
  rewritten by `make baseline-digests` as the deterministic fallout of the
  module edit.
- `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json` — rewritten only by the Managed Refresh.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md` — the Roundfix Skill describes
  the grouping suggestion, and its version rises in both front-matter
  fields.
- `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md` and their mirrors — the new
  section and the raised versions.

## What is not governed

The Go sources and tests under `internal/judge`, `internal/cli/spec_judge.go`
and its test, the new and changed tests under `internal/baseline`, the
catalog snapshots and plan goldens under `internal/baseline/testdata`,
`docs/user-guide/commands/spec.md`, `skills/roundfix/references/spec.md` and
`skills/testdata/owned-skill-versions.json` are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved edits and add
no source paths.

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No change to an existing Baseline clause, archived Specs, existing QA
  Reports, `skills/_ownership.yml`, `CONTEXT.md` or the `### QA settlement`
  section of any skill.
- No test, Verification command or QA row reaches OpenRouter, TypeSafe or the
  network, reads a real Jev key, or writes under the real `~/.roundfix`.
- No data leaves the machine except, when the judge runs with a Jev key, the
  Spec artifacts ADR-0201 names and the Findings and Backlog Entries ADR-0209
  names, sent to OpenRouter or to TypeSafe.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
