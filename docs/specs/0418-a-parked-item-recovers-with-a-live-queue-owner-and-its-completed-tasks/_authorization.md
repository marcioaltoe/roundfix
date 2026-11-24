---
status: approved
granted: 2026-11-23
action: record the latest Run on a Delivery Queue item, carry forward on a Delivery Retry the completed Tasks of a Run that ended Failed only in its QA gate, end a Delivery Retry with a proven live Queue Owner or with a failure, report in the delivery status a queue that no live owner advances and the Run and Tasks of each running item, and document them in the deliver guide, the Roundfix skill's deliver reference, the glossary and ADR-0450
consuming: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/roundfix/SKILL.md
  - docs/references/coverage-record.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0418

The messages and the maintainer's decisions this record rests on are copied whole in [`references/sources.md`](references/sources.md).

## Where the grant comes from

Three decisions of the maintainer of 2026-11-23, each delivered as a `DECISÃO DE MARCIO` block to this project's ai-memory inbox by the Praetor, the session that coordinates the software projects and puts their decision requests to the maintainer.

The Spec, its delivery and its release: ai-memory message `01a20996-258c-7fa1-af9d-a053fbce2441`, decided 15:25 BRT (18:25 UTC). The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20994-3bf2-7ed3-af99-543c9f49b0fc (itens a e b)
pergunta: Spec única de recuperação da entrega: o retry sobe dono vivo ou falha; o status mostra dono ausente, tarefa e Run em curso; um Run reprovado só na QA aproveita as tarefas concluídas no retry, e o run_id da Base é atualizado. Até 4 tasks mais a QA. Autoriza? Opções: sim, até release patch; sim, sem aproveitar tarefas; só o plano agora.
resposta: Sim, até release patch (Recomendado)
alcance: no repositório roundfix: uma Spec nova com PRD, TechSpec, tasks, implementação, PR, merge e release patch, incluindo a mudança de produto de aproveitar as tarefas completed num Run Failed só por reprovação na QA. Atende os achados 01a20714-3a6e e 01a20714-ed79 (Quaestor) e o item do Praetor (status sem a tarefa em curso, SEM AVANÇO falso no fluxus/0160). Fora do alcance: 0417, a regra de ADR, as correções do spec check, o amended_by e o cancel. Esses vêm em pedidos separados depois desta Spec.
quando: 23/11/2026 15:25 BRT
```

The three Governed Paths and their regeneration: ai-memory message `01a20a59-a318-7df1-89e7-a3c9baca4a62`, decided 18:42 BRT (21:42 UTC). The author of the PRD returned the question as a Reserved Question, the Supervisor session sent it as a decision request, and the answer came back as this block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a58-06b0-7f02-8adb-995c3ad6ebdf
pergunta: Roundfix, Spec 0418 (a recuperação da entrega autorizada às 15:25): ela muda o que o status e o retry da fila mostram. A skill do Roundfix descreve o comportamento antigo, e o próprio Roundfix recusa o release se a skill não mudar no mesmo PR. Três arquivos dela são protegidos e precisam de autorização. Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código; encolher a Spec.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/deliver.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 18:42 BRT
```

The coverage record, granted after the Run `run_20261124T001551Z_9b40dc21bfabfe76` failed `task_03` on `make verify-changed`, whose only failure was `TestCoverageEquivalence` reporting the intended removal of `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` as a coverage regression: ai-memory message `01a20b31-20fa-7650-8642-d4ee3b5d39a8`, decided 22:30 BRT (01:30 UTC on 2026-11-24). The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: mensagem de 23/11 22:35 BRT (0418 parada)
pergunta: A Spec 0418 pode regravar o registro de cobertura dos testes (docs/references/coverage-record.json, arquivo protegido), que hoje acusa como regressão o teste da regra antiga que a própria 0418 substitui (ADR-0450)? (Sim, regravar; não, parar)
resposta: "Sim, regravar (Recomendado)"
alcance: repositório roundfix, Spec 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks, até o release patch v0.153.1. Caminho: docs/references/coverage-record.json, regenerado só por: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1. Única remoção permitida: TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner. Inclui emendar a task_03 e o _authorization.md no branch do item e roundfix deliver retry. Sem Makefile, CI, go.mod ou lint.
quando: 23/11/2026 22:30 BRT
```

## What the grant covers

- The operations `implement`, `commit`, `push`, `pull_request`, `merge` and `release`.
- The product change the first decision names: a Run that ended Failed only because its QA gate failed hands its completed Tasks to the retry.
- Three Governed Paths:
  - `.agents/skills/roundfix/references/deliver.md`: the sections that describe what `deliver status` prints about the Queue Owner and about a running item, and what `deliver retry` carries, prints and refuses;
  - `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`: the version fields the record command raises, and nothing else.
- `docs/references/coverage-record.json`, written only by its own command under `## Sanctioned regeneration`; the only removal it may record is `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` (third decision).
- The regeneration by the repository's own commands, declared under `## Sanctioned regeneration`.
- A release at a `patch` Release Ceiling; the second decision names `v0.153.1`. The Spec is a `type: fix`. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

The grant names no number of corrections and no path beyond these. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are not granted ("Sem Makefile, CI, go.mod ou lint"), and this Spec's design needs none of them. Anything outside this scope needs a new decision request.

## How the governed set was measured

The author of the PRD measured it on 2026-11-23 at `fe247d4c` with `speccheck.GovernedPath`, over the files the change is expected to touch.

- The three paths above are the governed ones this Spec changes.
- These governed paths are not granted, and no Task may change them:
  - every other reference of the Roundfix skill, the `reconcile` reference included, and every other owned skill;
  - every Vendored Skill under `.agents/skills` and `skills-lock.json`'s entries;
  - every Baseline module, template, profile, Source Baseline and Setup Snapshot;
  - every guide under `docs/agents/` and `docs/agents/setup-context.json`;
  - `.roundfixrc.yml`, the `Makefile`, `.github/workflows` and `go.mod`.

The Tasks author measures again every file the Tasks declare. When that set holds a Governed Path outside the three above, no Run starts: the path needs a new decision of the maintainer, recorded here on the default branch before the consuming delivery.

## What is not governed

- The Go source and test files of the Delivery Queue, its commands and the Run Database.
- The mirror `skills/roundfix/references/deliver.md` and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`, `CONTEXT.md`, ADR-0450, ADR-0279 and ADR-0395, which gain their amendment notices, and `.secondbrain-export`.
- `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json`.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirror, the skill version fields and the generated records. This declaration records the regeneration that follows the approved edits and adds no source path. The commands are the ones the second decision names.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Vendored Skill, any Baseline module, template, profile, Setup Snapshot or plan golden;
  - any guide under `docs/agents/` or `docs/agents/setup-context.json`;
  - any Archive Record, and `docs/references/coverage-record.json` beyond its own regeneration command;
  - any skill other than `roundfix`, and any Roundfix reference other than `deliver`;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than the one this Spec adds and ADR-0279 and ADR-0395, which gain their amendment notices.
- No new column, table or stage in the Run Database, and no new command, flag or exit code.
- Outside this Spec, by the first decision: Spec 0417, the rule about ADRs, the corrections of the Spec Consistency Check, the amendment list of ADRs named by date, and the cancel of a Spec that holds only an Idea.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reads another repository, except a temporary one the test itself created;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
