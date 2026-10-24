---
status: approved
granted: 2026-10-24
action: make the internal/agent package adapters copies of the test binary instead of hard links, keep the re-executed fixture children from making them, name the RunPrompt error in the fake acpx helper, add the Hard Link rule test in internal/testfixture, and write the rule into docs/agents/specific-repository.md and the Single-Name Test Binary term into the glossary, without changing the Makefile or the CI workflows
consuming: 0311-the-suite-audit-and-no-test-hard-links-its-binary
paths:
  - docs/agents/specific-repository.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0311

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
  2 a 3 Specs ordenadas."); this Spec is the audit and the elimination of
  the two flakes that cost Runs on 2026-10-23;
- authorization: "Sim, só Makefile e workflows (Recomendado)" (option text:
  "Autoriza o Makefile, .github/workflows e flags de teste, limitado ao que a
  Spec declarar (ex.: rodar o docscontract uma vez por job, limitar o
  GOCACHE, paralelismo). go.mod e go.sum continuam proibidos."). This Spec
  declares no Makefile, workflow or test-flag change.

On 2026-10-24, at about 01:50 BRT, the maintainer answered the authoring
questions of this batch through AskUserQuestion; the answer this Spec rests
on is, for the leaked package adapter directories under `$TMPDIR`: "Apagar
depois da 0311 (Recomendado)", the operator deletes them after this Spec
ships, so no Task of this Spec removes them; this Spec only stops new ones.
The same answers order the batch 0311, 0312, 0313, 0314, all before 0307.

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
  ADR-0324 and the audit reference.
- `docs/agents/specific-repository.md` is governed, and task_02 adds one
  rule bullet to it.

## What is not governed

- `internal/agent/acpx_runner_test.go` and the new
  `internal/agent/package_adapter_test.go` (task_01).
- The new `internal/testfixture/hard_link_rule_test.go` and `CONTEXT.md`
  (task_02).
- `docs/adr/0324-a-test-never-gives-a-running-executable-a-second-name.md`,
  `docs/references/2026-10-24-test-suite-audit.md`, every file the Tasks read
  as an instruction, and every planning file of this Spec.

## Sanctioned regeneration

None. No skill, guide under `docs/user-guide/`, Baseline module, record or
derived file changes, so no regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any skill, any
  production `.go` file, or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
