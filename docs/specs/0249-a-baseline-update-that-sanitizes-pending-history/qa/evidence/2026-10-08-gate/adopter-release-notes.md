---
origin: roundfix
destination: conexus
type-hint: fix
created_at: 2026-10-07
capture: manual
---

# Roundfix v0.57.0: o que muda neste repositório

- **Pastas arquivadas antigas.** O `roundfix history sanitize` passa a converter pastas arquivadas com regras antigas do Task Graph. São dois casos: uma linha da tabela de Tasks que nomeia uma Task fora do grafo, e um tipo de Task que não é mais permitido (por exemplo `refactor`). As duas situações são toleradas e aparecem no plano como `tolerates ...`. Linhas malformadas e duplicadas continuam recusadas. Specs ativas e o `roundfix archive` mantêm as regras atuais.
- **Disposição `failed-qa`.** É para uma Spec que foi arquivada com QA reprovado e sem override. O registro guarda `qa_verdict: fail` e o nome do relatório, e nunca vira pass.
- **Recusas listadas.** O plano lista cada unidade recusada com o motivo e continua com as outras. O `--apply --batch <n>` deixa as recusadas intactas e não as conta no n.

## Como aplicar (opcional)

1. Atualize para a 0.57.0 e rode `roundfix baseline update` até ver `current`.
2. Rode `roundfix history sanitize` e leia o plano.
3. Crie a tag anotada `history-full` no main e faça o push.
4. Para cada lote:
   1. Crie a branch `chore/history-sanitize-<k>` e rode `roundfix history sanitize --apply --batch 40`.
   2. Faça `git add -A docs`.
   3. Rode a verificação do repositório.
   4. Abra um PR.
