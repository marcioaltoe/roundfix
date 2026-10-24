---
status: approved
granted: 2026-10-24
action: let make verify-docs take its documentation contract target from VERIFY_DOCS_TEST with a docs-test-skipped target, select it per event in the ci-verify Verify docs step, and prove it with tests in internal/verifyselect
consuming: 0312-the-ci-job-runs-the-documentation-contracts-once
paths:
  - .github/workflows/ci-verify.yml
  - Makefile
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0312

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
  0306 termina."); this Spec delivers right after Spec 0311;
- scope: "Auditoria completa" (option text: "Inventário de toda a suíte:
  testes lentos, redundantes, instáveis e com insumos ocultos, com o plano de
  2 a 3 Specs ordenadas.");
- authorization: "Sim, só Makefile e workflows (Recomendado)" (option text:
  "Autoriza o Makefile, .github/workflows e flags de teste, limitado ao que a
  Spec declarar (ex.: rodar o docscontract uma vez por job, limitar o
  GOCACHE, paralelismo). go.mod e go.sum continuam proibidos."). This Spec
  declares exactly the `Makefile` and `.github/workflows/ci-verify.yml`
  changes its Task names: running the documentation contracts once per CI
  job.

On 2026-10-24, at about 01:50 BRT, the maintainer answered the authoring
questions of this batch through AskUserQuestion:

- build cache: "roundfix gc recolhe" (option text: "Vira uma mudança do
  produto: o comando gc passa a recolher esse cache."); the shared build
  cache limit therefore left this Spec and is Spec 0314's, and this Spec
  changes no cache behavior;
- 0313: "Escrever agora (Recomendado)" (option text: "Autorar a 0313 já
  neste lote: os 13 testes identificados e os outros leitores dos docs vivos
  vão para a tag docscontract, que roda no limite do PR. Entrega depois da
  0312 e antes da 0307."); this Spec delivers before Spec 0313.

The cycle AP rules that answer carries were approved on 2026-10-22 (questions
of 15:25 BRT, answered by about 18:30 BRT): the Governed Paths each Spec
declares, with the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer
em Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
"Merge squash + release + push". The Makefile and CI exclusion of that grant
is lifted for this Spec's declared lines by the 2026-10-24 answer above.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0311`, based
on `56dd1a76`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and its adopted reference.
- `Makefile` and `.github/workflows/ci-verify.yml` are governed.

## What is not governed

- The new `internal/verifyselect/verify_docs_makefile_test.go` (task_01).
- Every file the Tasks read as an instruction, every planning file of this
  Spec and its adopted reference.

## Sanctioned regeneration

None. No skill, guide under `docs/user-guide/`, Baseline module, record or
derived file changes, so no regeneration command runs.

## Limits

- No change to `go.mod`, `go.sum`, the lint, formatter or test-runner
  configuration, any workflow other than `ci-verify.yml`, any Makefile
  target other than those task_01 names, `docs/agents/`, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any skill, any
  production `.go` file, `CONTEXT.md`, or the `### QA settlement` section of
  any skill.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
