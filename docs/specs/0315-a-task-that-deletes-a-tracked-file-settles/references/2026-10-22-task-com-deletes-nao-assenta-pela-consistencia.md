---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-22
capture: manual
---

# Uma Task com `deletes:` não assenta: a consistência da Spec lê o arquivo apagado pelo índice

## Observação

Roundfix 0.108.0 (`0af39b63`), oraculum, Spec 0092, Task 08 (`deletes: packages/backend/test/caixa-sem-venda-e-fluxus-meio-configurado.test.ts`). A Task foi escrita a pedido da resposta de vocês sobre o retry da 0092.

Duas Runs falharam no assentamento com a mesma mensagem:

- `run_20261022T225619Z_446a7c0f4282a311`
- `run_20261022T230908Z_46d197b7316725c1`

> Settlement check failed: settlement check: spec consistency: read refusing Spec Consistency findings: run "spec check --strict" for Spec 0092-as-vendas-do-erp-pelo-caixa: inspect citing file packages/backend/test/caixa-sem-venda-e-fluxus-meio-configurado.test.ts: lstat …: no such file or directory

O diagnóstico está em `~/.roundfix/artifacts/d1255e752e9d9f49/runs/run_20261022T230908Z_46d197b7316725c1/verification/batch-001-attempt-2-settlement-spec-consistency.log`.

## Causa provável

O Agente apaga o arquivo só do disco, sem commit e sem mexer no índice, porque o commit é do Daemon. O arquivo continua em `git ls-files`. A consistência parece enumerar os arquivos rastreados e fazer `lstat` em cada um como possível citador. Num worktree descartável:

- com a remoção preparada no índice (`git rm`), o `spec check --strict` passa;
- com a remoção só no disco, a leitura falha.

O efeito é que nenhuma Task com `deletes:` consegue assentar, e o `deliver.md` prevê esse tipo de Task.

## O que o oraculum fez

Para contornar, a Task 08 passou a pedir `git rm` ao Agente (commit `0d42da5` no worktree do item), e uma Run nova está rodando: `run_20261023T004828Z_d579279a4a2193eb`.

## Pedido

Confirmar se é defeito e corrigir. Por exemplo, ignorar no inventário de citadores os caminhos rastreados que não existem no disco, ou ler o inventário da árvore de trabalho. Se a correção couber na 0.110.0, melhor ainda.
