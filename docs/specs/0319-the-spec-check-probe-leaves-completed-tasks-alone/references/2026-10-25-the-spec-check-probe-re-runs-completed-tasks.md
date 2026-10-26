---
type: fix
status: promoted
created: 2026-10-25
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
reason: null
---

# The Spec Check probe re-runs the Verification of completed Tasks

## Symptom

`roundfix spec check <slug> --strict --run-verification` probes the
Verification commands of every Task in the graph, including a Task whose
status is already `completed`. On a resumed Spec a completed Task's command
either passes, because its work now exists, and is reported `vacuous`, or
cannot run in the disposable checkout without the preparation its work had,
and is reported `unknown`. Either verdict makes the check exit `1` before a
Run that would never execute that Task again.

## Where

`internal/cli/spec_check.go` `probeSpecVerifications`, which loops over every
`graph.Tasks` entry and calls `daemon.ProbeCommands` whatever the Task's
status, and the `--run-verification` report (text and JSON), which has no way
to say that a Task was left unprobed.

## Expected

The check probes the same Tasks the Daemon would still execute, names every
Task it leaves unprobed with the reason, and keeps failing on a vacuous,
unknown or malformed command of a Task that is still to run. Structure,
provenance and QA reopen checks are unchanged, an `unknown` of a pending Task
never becomes success, and nothing installs packages to hide the distinction.

## Evidence

Vortex, Roundfix 0.119.0 (`f3c39140`), Spec 0162 at
`gesttione-solutions/vortex@a961a29`: Task 01 completed (commit `f3f683fdd`,
Verification proved by the Daemon and recovered by Task Carry-Forward from Run
`run_20261025T151524Z_7d3d717be10dcbdd`), Tasks 02–04 pending. The check
probed Task 01's Bun step in a checkout without `node_modules`: exit 127,
verdict `unknown`, final exit `1`; the three pending commands were `honest`.

Provenance: ai-memory message `01a17436-c629-7622-a37d-12528b1da9f2` from
`default/marcioaltoe-secondbrain`, created 2026-10-25 14:58 BRT, migrating the
retired Secondbrain inbox entry
`inbox/roundfix/2026-10-25-spec-check-sonda-verification-de-task-ja-concluida.md`
(origin vortex). The message body, as received:

```markdown
---
origin: vortex
destination: roundfix
type-hint: finding
created_at: 2026-10-25
capture: manual
---

# A sonda de honestidade executa novamente Tasks já concluídas

Roundfix 0.119.0, build f3c39140. Na retomada da Spec 0162 do Vortex, a Task 01
estava completed, com commit f3f683fdd e prova de Verification aprovada pelo
daemon e recuperada por carry-forward. As Tasks 02–04 continuavam pending.

Em `gesttione-solutions/vortex@a961a29`,
`roundfix spec check 0162-dependencias-e-tooling-sem-diagnosticos --strict --run-verification`
sondou também o comando da Task 01. O grep da nova cobertura passou, e a etapa
Bun tentou executar testes no checkout descartável sem node_modules: exit 127,
verdict unknown, exit final 1. Os três comandos pendentes foram honest.
A execução normal por implement continuou somente pelas Tasks pendentes.

O efeito observado é um resultado vermelho de pré-execução por reavaliar a
honestidade de trabalho já entregue, sem sua preparação. Não se trata de falha
da implementação ou de prova perdida no commit. O contrato também permite que
um teste já implementado passe, tornando vacuous o que é uma pós-condição válida.

Investigar se a sonda deve limitar-se a Tasks pendentes, ou separar resultados
históricos dos que bloqueiam a retomada. Preservar checagens de estrutura,
proveniência e reabertura de QA. Não transformar unknown de Task pendente em
sucesso e não instalar pacotes indiscriminadamente como substituto da distinção.

Evidência local: `/private/tmp/vortex-tooling-spec-verification.log` nesta sessão;
Task 01 recuperada da Run run_20261025T151524Z_7d3d717be10dcbdd. O arquivo temporário
pode expirar; os fatos reproduzíveis e revisões estão descritos acima.
```
