---
status: approved
granted: 2026-11-14
action: add the Idea-stage judgment to the judge and to `roundfix spec judge`, reword the route clauses of the Baseline spec-workflow module so a feature is judged for the Idea stage before its PRD, with their derived files and this repository's managed refresh, change the `brainstorming`, `write-idea`, `write-prd` and Roundfix skills to run and record that judgment, and record the glossary term and the command's documentation
consuming: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
paths:
  - .agents/skills/brainstorming/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-idea/references/idea-template.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/brainstorming/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0393

The Spec adopts the maintainer's request and answers of 2026-11-14, given in
the session of `/roundfix --follow-up` between about 07:00 and 08:40 BRT;
requested by the maintainer. They are copied verbatim under `references/`.

The request, written while the follow-up was planning Specs 0391 and 0392:

> "Outro ponto é que gostaria de um julgamento (usando o JEV) sobre a
> necessidade do uso do write-idea para features/specs mesmo sem o
> acionamento do usuário. Essa skill levanta questões que podem ser muito
> valiosas e evitar trabalho inutil antes de implementações de novas features
> ou de roadmaps."

Where the judgment lives. Question: "Você pediu um julgamento do Jev sobre a
necessidade do write-idea mesmo sem o usuário acionar. Hoje a regra de
roteamento só manda rodar write-idea em "iniciativa grande ou difusa", quem
decide é o agente, e o `roundfix spec judge` não faz essa pergunta (só tem os
estágios prd, techspec e tasks). Onde esse julgamento deve morar?" Answer:
"Spec nova no produto (Recommended)" — option text: "Feature: antes do PRD de
uma feature, o Roundfix pergunta ao Jev se ela precisa da etapa de Ideia
(critérios da própria write-idea: forma da solução em aberto, suposição de
negócio não conferida, questão de mercado) e a regra de roteamento passa a
mandar rodar write-idea quando o juiz indicar. Vale para todos os adotantes."

The authorization. Question: "Se o julgamento do write-idea virar Spec no
produto (primeira pergunta), qual autorização vale para ela? É feature, então
a release dela é minor." Answer: "Mesmas do BE (Recommended)" — option text:
"Push, PR e merge; release minor depois dela; até 3 correções; Makefile, CI,
go.mod e lint só quando uma correção exigir, descrito no relatório."

The Idea stage, whose answers decide the product. Each is quoted with its
question in the adopted record and in the PRD's `## Decisions`:

- the waste to avoid: "Feature na forma errada (Recommended)";
- the authority of the judgment: "Roda, com saída registrada (Recommended)" —
  option text: "O agente inicia a write-idea sem pedir licença e mostra os
  números do julgamento. Você pode dispensar a etapa na primeira pergunta
  dela, e a dispensa fica escrita no PRD com o motivo.";
- when the judge does not answer: "O agente aplica os critérios (Recommended)"
  — option text: "O agente responde ele mesmo às mesmas perguntas (forma em
  aberto, suposição sem evidência, questão de mercado), diz que o juiz não
  estava disponível e segue o resultado; na dúvida, pergunta a você. Todo
  adotante ganha a regra, com ou sem chave.";
- success: "Menos correções depois";
- the shape: "Mais ambiciosa" — option text: "A original com registro, mais um
  relatório que remede os limiares com o histórico do próprio repositório e o
  conselho por item já no relatório do `roundfix follow-up`. Duas Specs.";
- the draft of the Idea: "Aprovado (Recommended)".

The cycle BE terms that "Mesmas do BE" names are the answers of 2026-11-13
(AskUserQuestion, about 08:20 BRT): releases, "Patch após cada Spec
(Recommended)", and corrections, "Até 3 por Spec + Makefile/CI (Recommended)"
— option text: "Até 3 correções por Spec, contando as rodadas do review. Fica
liberado mexer em Makefile, workflows do CI, go.mod e lint quando a correção
exigir; cada mudança vai descrita no relatório." Those terms keep the grants of
the earlier cycles, which release to the Runs the Baseline source
(`internal/baseline/assets/**`) with its fixtures and goldens, the guides
under `docs/agents/`, the owned skills, and changes of Baseline and skills
that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are granted only for a correction that requires one.
The design needs none of them, and none is in `paths`.

## What the judge may be sent

The maintainer's boundary for this Spec is the Idea's default: before a PRD
exists, the judge receives only the kind of text the follow-up judge already
sends, an item that is not yet a Spec, after the same secret scan. This Spec
sends less than that boundary allows: an open Backlog Entry, a note in
`docs/_inbox/` or a message handed in as a file, and never a Finding, an
active Spec, a Waiting Item or a request that exists only in a session. No
Task may widen it; a design that needs any other text stops and asks.

## The governed set

The governed set was measured on 2026-11-14 on the branch `docs/plan-0393`,
based on `251ed81e`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files afterwards, and the probe file never existed in the tree.
- It covered every file the TechSpec's file list names, every file a
  prototype of the Baseline and skill changes rewrote on a scratch worktree
  (the derived Baseline files included), every planning file of this Spec, its
  adopted references, the measurement reference and the two ADRs.
- The eighteen paths above are governed: the spec-workflow module, the
  Standard TypeScript profile and its formatter golden of the Spec routing
  guide, this repository's Spec routing guide and its setup manifest, and the
  canonical files and `SKILL.md` mirrors of the four skills, with the mirror of
  the PRD template.

## What is not governed

- Every Go source and test file under `internal/judge`, `internal/cli`,
  `internal/baseline` and `skills`, and `internal/judge/questions.json`.
- The mirrors under `skills/roundfix/references/` and
  `skills/write-idea/references/`, `skills/testdata/owned-skill-versions.json`,
  `internal/baseline/module-versions.json` and the catalog snapshots and plan
  goldens under `internal/baseline/testdata/`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`
  and `docs/references/skill-coverage.json`.
- `CONTEXT.md`, ADR-0403, ADR-0404 and
  `docs/references/idea-stage-judgment-measurement.md`.
- Every planning file of this Spec and its adopted references.

`internal/cli/cli_test.go` is a Governed Path this record does not grant. The
help of `roundfix spec judge` gains a line and keeps every line a test there
pins, so the new assertions go into the new test file.

## Sanctioned regeneration

The repository-owned commands write the module version record, the derived
Baseline files, this repository's managed guide and setup manifest, the skill
mirrors, the owned-skill version record and the Behavior Surface Record. This
declaration records the regeneration that follows the approved edits and adds
no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/spec-routing.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/brainstorming/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/brainstorming/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module other than `spec-workflow`, any clause of that module other than the
  three route clauses, any Source Baseline, Baseline template or Setup
  Snapshot, any guide under `docs/agents/` other than through the managed
  refresh, `docs/references/coverage-record.json`, any Archive Record, the
  `### QA settlement` section of any skill, the grouping sentences Spec 0391
  owns or the skill installation Spec 0392 owns.
- No clause is added, removed or renamed. No command gains a refusal: the
  Idea-stage judgment changes no other command's output or exit code.
- The follow-up report, the calibration report and any reader of the judgment
  record belong to Spec 0395.
- No test, Verification command or QA row:
  - reaches the network, a provider, the judge or a reviewer, calls ai-memory
    or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
