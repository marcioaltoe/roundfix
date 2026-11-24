---
status: approved
granted: 2026-11-23
action: record in the Reviewer Lineage the model of the session each review round used, open a new reviewer session in round 2 when the review profile's selection no longer names that model, and document it in the review guide, the Roundfix skill's review reference, the glossary and ADR-0451
consuming: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
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

# Approved authority for Spec 0419

The message and the maintainer's decisions this record rests on are copied whole in [`references/sources.md`](references/sources.md).

## Where the grant comes from

Three decisions of the maintainer of 2026-11-23, each delivered as a `DECISÃO DE MARCIO` block to this project's ai-memory inbox by the Praetor, the session that coordinates the software projects and puts their decision requests to the maintainer.

The Spec, its delivery and its release: ai-memory message `01a20a6c-7dce-7bb3-9a50-b09b27461492`, decided 18:48 BRT (21:48 UTC). The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a5a-ebf0-7ba2-8ac0-95a330f45eeb
pergunta: Roundfix: a fila do Fiscus está parada na Spec 0108 porque, na segunda rodada da review antes do PR, o Roundfix retoma a sessão antiga do revisor (Codex) mesmo com o modelo configurado diferente (gpt-6-luna → gpt-6.1-sol); o Codex avisa a troca e o classificador recusa a resposta, mesmo "No findings". Abrir uma Spec curta de correção? Opções: antes da 0418; depois da 0418; não, o Fiscus contorna.
resposta: "Sim, antes da 0418 (Recomendado)"
alcance: repositório roundfix, Spec nova 0419: PRD, TechSpec, Tasks, implementação, PR, merge e release patch v0.153.1 (a 0418 passa a v0.153.2). Código em internal/cli/review*.go e testes; sem Makefile, CI, go.mod, lint nem skill protegida — se precisar de caminho protegido, novo PEDIDO DE DECISÃO. Avise o Fiscus (default/gesttione-solutions-fiscus) quando o release sair, para ele dar o retry da 0108.
quando: 23/11/2026 18:48 BRT
```

The three Governed Paths and their regeneration: ai-memory message `01a20ac2-3453-71c1-8227-b02e7cb2c905`, decided 20:33 BRT (23:33 UTC). The author of the TechSpec returned the question as a Reserved Question, the Supervisor session sent it as a decision request, and the answer came back as this block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a8f-f044-7402-ab98-56f7e59b3420
pergunta: Roundfix, Spec 0419 (a revisão abre sessão nova quando o modelo muda): a skill descreve a retomada da sessão e a regra do repositório exige mudá-la no mesmo PR. Três arquivos são protegidos (references/review.md e as duas cópias do SKILL.md). Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0419, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/references/review.md, .agents/skills/roundfix/SKILL.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 20:33 BRT
```

The order of the two Specs and the number of the release: ai-memory message `01a20ace-16f5-7c32-b81d-7e73d0340f02`, decided 20:46 BRT (23:46 UTC). The block, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a20a93-6180-76a2-b59f-20ee9895af31
pergunta: O plano da 0418 está pronto e pode rodar esta noite (6 a 10 h); a 0419 acabou de receber a autorização da skill (01a20ac2-3453) mas ainda precisa ser planejada, e o Fiscus sai com um segundo retry. Inverter a ordem e começar a 0418 hoje?
resposta: "Sim, 0418 hoje (Recomendado)"
alcance: repositório roundfix: merge do PR #946 (plano da 0418) e roundfix deliver start da 0418 hoje; a 0418 sai como v0.153.1 e a 0419 como v0.153.2, entrando na fila logo que planejada. Nada mais muda nas autorizações das duas Specs (0418: 01a20a59-a318; 0419: 01a20a6c-7dce e 01a20ac2-3453).
quando: 23/11/2026 20:46 BRT
```

The third decision changes two things the first two say and nothing else: Spec 0419 is delivered after Spec 0418, and its release is `v0.153.2`, not `v0.153.1`.

## What the grant covers

- The operations `implement`, `commit`, `push`, `pull_request`, `merge` and `release`.
- Three Governed Paths:
  - `.agents/skills/roundfix/references/review.md`: the two paragraphs that describe the reviewer session round 1 leaves open, what round 2 does with it and what each round records;
  - `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`: the version fields the record command raises, and nothing else.
- The regeneration by the repository's own commands, declared under `## Sanctioned regeneration`.
- A release at a `patch` Release Ceiling; the third decision names `v0.153.2`. The Spec is a `type: fix`. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

The grant names no number of corrections and no path beyond these. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are not granted ("Sem Makefile, CI, go.mod ou lint"), and this Spec's design needs none of them. Anything outside this scope needs a new decision request.

The first decision bounds the code to `internal/cli/review*.go` and tests. The design changes two such files, `internal/cli/review.go` and `internal/cli/review_lineage.go`, and adds test files only.

## How the governed set was measured

The author of the TechSpec measured it on 2026-11-23 at `50143e74` with `speccheck.GovernedPath`, in a `go test -overlay` probe whose test file lives outside the repository, over the fifteen files the TechSpec's Build Order names, `docs/references/skill-coverage.json` and this Spec's reference file.

- The three paths above are the governed ones this Spec changes.
- These governed paths are not granted, and no Task may change them:
  - every other reference of the Roundfix skill, and every other owned skill;
  - every Vendored Skill under `.agents/skills` and `skills-lock.json`'s entries;
  - every Baseline module, template, profile, Source Baseline and Setup Snapshot;
  - every guide under `docs/agents/` and `docs/agents/setup-context.json`;
  - `.roundfixrc.yml`, the `Makefile`, `.github/workflows` and `go.mod`.

The Tasks author measures again every file the Tasks declare. When that set holds a Governed Path outside the three above, no Run starts: the path needs a new decision of the maintainer, recorded here on the default branch before the consuming delivery.

## What is not governed

- The Go source and test files of the Pre-PR Review Command, and the new test file of `internal/docscontract`.
- The mirror `skills/roundfix/references/review.md` and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/review.md`, `CONTEXT.md`, ADR-0451, ADR-0197, which gains its amendment notice, and `.secondbrain-export`.
- `docs/references/behavior-surfaces.json`, and `docs/references/skill-coverage.json`, which this Spec does not change.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirror, the skill version fields and the generated record. This declaration records the regeneration that follows the approved edits and adds no source path. The commands are the ones the second decision names.

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

No mirror, version field, digest or record is hand-edited. Spec 0418 raises the same version fields and re-records the same two records before this Spec is delivered: every command above runs on the tree the Task finds, and the Task writes the next version, never a number.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Vendored Skill, any Baseline module, template, profile, Setup Snapshot or plan golden;
  - any guide under `docs/agents/` or `docs/agents/setup-context.json`;
  - `docs/references/skill-coverage.json`, `docs/references/coverage-record.json` and any Archive Record;
  - any skill other than `roundfix`, and any Roundfix reference other than `review`;
  - the `### QA settlement` section of any skill;
  - any accepted ADR other than the one this Spec adds and ADR-0197, which gains its amendment notice.
- No change to the classifier of the reviewer's answer: a no-findings verdict beside other content still blocks.
- No new command, flag, output line or exit code, and no new key of Project Config.
- Outside this Spec: everything of Spec 0418, and every other pending message of this project's inbox.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reads another repository, except a temporary one the test itself created;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement. The operator sends the notice to the Fiscus that the first decision asks for, after the release.
