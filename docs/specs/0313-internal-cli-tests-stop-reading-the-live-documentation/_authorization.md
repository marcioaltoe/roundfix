---
status: approved
granted: 2026-10-24
action: move the thirteen internal/cli tests that read this repository's live documentation trees, unchanged, into two docscontract files declared verify:always, classify those files in the repository contract declarations test, and write the rule into docs/agents/specific-repository.md
consuming: 0313-internal-cli-tests-stop-reading-the-live-documentation
paths:
  - docs/agents/specific-repository.md
  - internal/cli/baseline_documentation_contract_test.go
  - internal/cli/cli_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0313

On 2026-10-23, at about 23:55 BRT, the maintainer asked: "Vejo que
frequentemente temos problemas com os testes implementados para roundfix.
Precisamos trabalhar em testes mais eficientes e rápidos. Isso deve ser uma
spec de performance para os próximos passos de evolução de roundfix,
saneamento dos testes e foco de performance para garantir a qualidade."

On 2026-10-24, at about 00:05 BRT, the maintainer answered cycle AP batch 4
through AskUserQuestion:

- order: "Logo após a 0306 (Recomendado)" (option text: "Entra antes da
  0307, da 0309 e da 0310: cada Spec que vem depois roda a mesma suíte, e os
  testes instáveis continuariam custando Runs. Eu escrevo o plano enquanto a
  0306 termina.");
- scope: "Auditoria completa" (option text: "Inventário de toda a suíte:
  testes lentos, redundantes, instáveis e com insumos ocultos, com o plano de
  2 a 3 Specs ordenadas."); this Spec cuts the largest hidden test input
  the audit found;
- authorization: "Sim, só Makefile e workflows (Recomendado)" (option text:
  "Autoriza o Makefile, .github/workflows e flags de teste, limitado ao que a
  Spec declarar (ex.: rodar o docscontract uma vez por job, limitar o
  GOCACHE, paralelismo). go.mod e go.sum continuam proibidos."). This Spec
  declares no Makefile, workflow or test-flag change; it moves tests between
  build tags inside `internal/cli`.

On 2026-10-24, at about 01:50 BRT, the maintainer answered the authoring
questions of this batch through AskUserQuestion; the answer this Spec rests
on is:

- 0313: "Escrever agora (Recomendado)" (option text: "Autorar a 0313 já
  neste lote: os 13 testes identificados e os outros leitores dos docs vivos
  vão para a tag docscontract, que roda no limite do PR. Entrega depois da
  0312 e antes da 0307."). The audit found no other reader: with the trees
  removed or unreadable only the thirteen fail, and with them skipped a
  guide edit leaves the package cached. The same answers order the batch
  0311, 0312, 0313, 0314, all before 0307.

The cycle AP rules that answer carries were approved on 2026-10-22 (questions
of 15:25 BRT, answered by about 18:30 BRT): the Governed Paths each Spec
declares, with the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer
em Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
"Merge squash + release + push".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0311`, based
on `56dd1a76`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  ADR-0326 and its adopted reference.
- `docs/agents/specific-repository.md`, `internal/cli/cli_test.go` and
  `internal/cli/baseline_documentation_contract_test.go` are governed;
  task_02 adds one rule bullet to the first and moves tests out of the other
  two without changing any other byte.

## What is not governed

- The two new `internal/cli/*_docscontract_test.go` files, the five other
  `internal/cli` test files the tests leave, and
  `internal/verifyselect/contracts_repository_test.go`.
- `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`,
  every file the Tasks read as an instruction, every planning file of this
  Spec and its adopted reference.

## Sanctioned regeneration

None. No skill, guide under `docs/user-guide/`, Baseline module, record or
derived file changes, so no regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any skill, any
  production `.go` file, any moved test's body, `CONTEXT.md`, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
