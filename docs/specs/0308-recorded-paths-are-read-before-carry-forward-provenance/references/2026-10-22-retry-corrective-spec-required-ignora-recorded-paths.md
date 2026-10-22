---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-22
capture: manual
---

# O retry de `corrective-spec-required` recusa um caminho que está nos Recorded paths da revisão registrada

## Observação

Roundfix 0.108.0 (`0af39b63`), oraculum, Spec 0092-as-vendas-do-erp-pelo-caixa, em 22/10/2026 às 17h47 BRT.

- O candidato parado é o `d21b36d` ("docs: archive 0092-…"). O registro arquivado tem `source_revision: 01cf38411679c1138f68ef8ebbca203b8bc07cf2`.
- A revisão pré-PR, na rodada 1, achou F1 (`contextos-ferramentas.ts`) e F2 (`comparacao-resumo-vendas.use-case.ts`), e o item parou em `corrective-spec-required`.
- A correção seguiu a receita:
  1. revert do archive;
  2. Task corretiva 07;
  3. Run no worktree do item;
  4. carry-forward;
  5. archive novo (`b010828`, `source_revision b12d75b`);
  6. `review dispose F1/F2 --fixed-by 4419ac3`, as duas como `fixed`.
- `roundfix deliver retry` recusou:

> changed path packages/backend/src/application/use-cases/vendas/comparacao-resumo-vendas.use-case.ts is outside the named archived Specs and the paths their Tasks declared; author a corrective Spec with its own authorization and QA gate

## Por que parece defeito

O `references/deliver.md` diz que o caminho mudado pode estar entre os paths que as Tasks declararam, incluindo "`creates:`, `interface:` and recorded paths read at the recorded revision". Em `01cf384`, a `task_02.md` da 0092 lista esse arquivo em `## Recorded paths`, escrito pelo Daemon na Run `run_20261022T104854Z_1d8a777858a2f895`. No head do item, a `task_07.md` o declara em `interface:`. O outro arquivo, `contextos-ferramentas.ts`, está em `interface:` da task_01 e passou.

Pelo texto, a recusa não deveria acontecer. A suspeita é que a leitura não considera os Recorded paths. Outra possibilidade é que ela leia a Spec numa revisão diferente da `source_revision` do registro do candidato parado.

## Pedido

1. Confirmar se é defeito e corrigir.
2. Se não for, dizer qual caminho a fila aceita para uma correção que mexe num arquivo que só aparece em Recorded paths, sem uma Spec corretiva nova. O oraculum tem prazo em 25/10 para a release v0.26.0 dessa Spec.

O item segue parado e sem mudança. O worktree do item está em `~/.roundfix/worktrees/oraculum-d1255e75/roundfix-deliver-0092-as-vendas-do-erp-pelo-caixa-91f2c7934ce6a99a-8b1ada85`.
