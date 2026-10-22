---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-21
capture: manual
---

# O portão de QA recusa SC-WAVE-COLLISION com `needs` declarado e `spec check --strict` limpo

## Observação

No oraculum (Roundfix 0.101.0, `956066a7`), a Spec
`0087-o-rateio-repartido-pelo-percentual` ganhou uma Task corretiva `task_05`
depois que a revisão pré-PR achou F1 num candidato já arquivado
(`corrective-spec-required`). O portão (`task_04`) recusou a pré-condição
quatro vezes seguidas, sem medir nada, com o mesmo texto:

`SC-WAVE-COLLISION: Task task_02 and Task task_05 can run in the same Wave and share repository files: packages/backend/test/rateio-repartido-versoes.test.ts (declared context; prior Run settlement commit)`

Nas quatro, `roundfix spec check 0087-o-rateio-repartido-pelo-percentual --strict`
sai 0, sem erro nem gap, no mesmo head que o portão auditou, rodado no worktree
da própria Run e no worktree do item.

## O que foi tentado

1. `task_05` com `needs: [task_01]` (mesma onda da `task_02`): recusa, que é
   o esperado. Run `run_20261021T055620Z_3fb17e656d321264`.
2. Edge `task_05 needs: [task_02]` commitado depois do carry-forward da
   `task_05`: recusa igual. Run `run_20261021T060151Z_58ef485f1d7de325`.
3. Branch do item refeita: edge commitado ANTES de a `task_05` rodar, e a
   `task_05` executada de novo já com o edge: recusa igual. Run
   `run_20261021T061154Z_951220e90496fbcf`.
4. `task_05` passando a declarar o arquivo em `interface:`: recusa igual.
   Run `run_20261021T062605Z_f9df1c41dbcd4303`.

Os relatórios estão em `docs/specs/0087-o-rateio-repartido-pelo-percentual/qa/qa-report-2026-10-21-01.md`
nas branches `roundfix/run-run_20261021T06*` do oraculum.

## Hipótese

O "prior Run settlement commit" sugere que o detector do portão lê o
assentamento da `task_02` na primeira Run da Spec
(`run_20261021T051927Z_681196740b6ebb70`), cujo grafo ainda não tinha a
`task_05`, e conclui que as duas "podem" cair na mesma onda. O
`spec check --strict` fora do portão não vê essa Run e passa. Uma Task
corretiva acrescentada depois de um portão completo, que precisa mudar um
arquivo que uma Task anterior criou, fica sem saída: nenhuma ordem de `needs`
satisfaz o detector.

## Efeito e contorno

O item parou em `corrective-spec-required` e depois o portão nunca passou. O
oraculum arquivou com `roundfix archive --qa-override`, citando o relatório
`partial` completo anterior (15 linhas aprovadas, zero falha, zero achado) e a
Verification verde da `task_05`, registrou F1 como corrigido com
`roundfix review dispose` e seguiu com `deliver retry`.

## Pedido

Fazer o detector do portão usar o grafo atual (o mesmo do
`spec check --strict`), ou dizer na mensagem qual edge ou declaração
resolveria, para que uma Task corretiva tenha uma saída sem override.
