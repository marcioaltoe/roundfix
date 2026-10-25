---
status: approved
granted: 2026-10-24
action: move the thirteen internal/cli tests that read this repository's live documentation trees, unchanged, into two docscontract files declared verify:always, classify those files in the repository contract declarations test, and write the rule into docs/agents/specific-repository.md
consuming: 0313-internal-cli-tests-stop-reading-the-live-documentation
paths:
  - docs/agents/specific-repository.md
  - internal/cli/baseline_documentation_contract_test.go
  - internal/cli/cli_test.go
  - internal/suiteguard/suiteguard.go
  - internal/suiteguard/suiteguard_test.go
  - internal/suiteguard/residue_test.go
  - internal/suiteguard/leak_cases_test.go
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

## Corrective amendment of 2026-10-25

The first Run, `run_20261025T122115Z_fe88d4e7bd75bcb2`, completed task_01
and task_02, and its QA gate failed on F-01: a guide or skill edit still
re-ran `internal/cli`. On 2026-10-25, at about 10:30 BRT, the maintainer
answered through AskUserQuestion:

- order: "0317 e 0318 agora (Recomendado)" (option text: "A 0313 fica
  estacionada; a 0317 entra já na fila (v0.120.0) e a 0318 em seguida. A
  0313 volta depois, com a correção.");
- fix: "Fingerprint pelo git (Recomendado)" (option text: "O suiteguard
  passa a calcular a impressão do repositório por `git` (índice, diff e
  arquivos não rastreados) num subprocesso, sem abrir os arquivos pelo
  processo de teste. O guarda continua igual e o cache deixa de depender
  dos docs. Emenda da 0313 (autoriza internal/suiteguard) no main + Task
  corretiva.").

Both answers rest on F-01's attribution of the read to `internal/suiteguard`.
The corrective measurement (`_techspec.md` → Current behavior 6 and 7)
disproved it. The guard walks the repository outside `m.Run`, where Go's
test input log records nothing. The reader was
`TestUpgradeStdoutAndExitCodesAreCharacterized`, and one
`withCLIWorkspace(t)` line made the cache row print `(cached)`.

On 2026-10-25, at about 11:00 BRT, the maintainer was asked through
AskUserQuestion: "Para a 0313, seguimos com a causa real (uma linha no
teste de upgrade) em vez do fingerprint pelo git que você tinha escolhido
com base no diagnóstico da QA?" The answer was "As duas coisas" (option
text: "A linha no teste e também o fingerprint pelo git no suiteguard
(segunda corretiva, última permitida). O guarda perde a detecção de mudança
de permissão, exceto o bit de execução.").

So this Spec has two corrective Tasks, the limit:

- task_04 adds the `withCLIWorkspace(t)` line;
- task_05 makes the guard fingerprint through `git` subprocesses (listing,
  untracked directories, and object IDs through a temporary index), so the
  test process opens no repository file.

The loss of permission detection that the option text accepts does not
happen. task_05 keeps a mode-only `os.Lstat` pass over the watched paths.
That pass runs before and after `m.Run`, outside Go's test input log, so it
is not a cache input, and the permission cases of its leak table pass
before and after the change. The four `internal/suiteguard` files in
`paths` are the ones task_05 changes or creates. They are not Governed
Paths, and they are listed because the maintainer granted them by name
("autoriza internal/suiteguard").

## What is not governed

- The two new `internal/cli/*_docscontract_test.go` files, the five other
  `internal/cli` test files the tests leave, and
  `internal/verifyselect/contracts_repository_test.go`.
- `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`,
  every file the Tasks read as an instruction, every planning file of this
  Spec and its adopted reference.
- `internal/cli/upgrade_characterization_test.go`, the one file task_04
  changes, and `task_04.md`. These were measured on 2026-10-25 with the same
  `go test -overlay` probe, on the branch `docs/0313-corrective-suiteguard`
  based on `1cdd5b55`. `git status --porcelain` was identical before and
  after.
- The four `internal/suiteguard` files task_05 changes or creates, and
  `task_05.md`. The same probe reported them as not governed; they are in
  `paths` by the maintainer's grant.

## Sanctioned regeneration

None. No skill, guide under `docs/user-guide/`, Baseline module, record or
derived file changes, so no regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any skill, any
  production `.go` file other than `internal/suiteguard/suiteguard.go`, any
  file under `internal/suiteguard/` other than the four in `paths`,
  `internal/suiteguardcontract/`, any caller of `suiteguard.Main`, any moved
  test's body, `CONTEXT.md`, or the `### QA settlement` section of any
  skill.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
