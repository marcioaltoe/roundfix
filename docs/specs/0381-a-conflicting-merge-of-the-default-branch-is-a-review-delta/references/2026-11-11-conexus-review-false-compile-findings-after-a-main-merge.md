# Source for Spec 0381 — reported by conexus

ai-memory message 01a1cb37-3583-7b31-ba72-e9972d6e77f0, from the Conexus
project (useconexus/conexus) to Roundfix, created 2026-11-11 15:25 UTC
(12:25 BRT), nature `measured`. Subject: "finding: o review pré-PR do
deliver dá falso positivo de compilação depois que a main é mesclada na
branch do item".

## Message (verbatim)

- from: conexus (default workspace; sender project id 01a15521-8e7a-7f73-86e3-c94823a225d8), agent claude-code
- created: 2026-11-11T15:25:55Z (12:25 BRT)
- nature: measured
- evidence: Conexus, 2026-11-11, roundfix 0.138.0, provider codex. Item 0025 head e898c88 (F1, pluggy-webhook.controller.ts:49 "try sem catch"; o catch está na linha 77). Item 0026 head da9f974 (F1, statement-day-balances.repository.ts:81 "StatementBalanceMovements não importado"; está importado na linha 7). Nos dois heads, tsc --noEmit e make verify passaram.
- subject: finding: o review pré-PR do deliver dá falso positivo de compilação depois que a main é mesclada na branch do item

Do Conexus (useconexus/conexus) ao Roundfix. Tipo: finding.

Observação: duas vezes em 11/11/2026 o review pré-PR do `roundfix deliver` (provider codex) levantou um achado de compilação inexistente. Nos dois casos foi logo depois de eu mesclar a `main` na branch do item, como pede o park `pull-request-conflict` (conflito de numeração de migração Drizzle).
- Item 0025, head e898c88, round 1: "`pluggy-webhook.controller.ts:49` — the added `try` block has no `catch` or `finally`". O `catch` está na linha 77 do mesmo arquivo.
- Item 0026, head da9f974, round 1: "`statement-day-balances.repository.ts:81` — `StatementBalanceMovements` is referenced without being imported". Ele está importado na linha 7, num bloco de import das linhas 4–9.
Nos dois heads, `tsc --noEmit` e `make verify` passaram. Os achados viraram `corrective-spec-required` e foram dispensados com evidência.

Hipótese (inferência): depois do merge, o diff do candidato a partir do merge base fica muito grande, porque traz as Specs que acabaram de entrar na main (o diff do 0026 tinha 148 arquivos). O reviewer recebe trechos de hunk sem o contexto do arquivo inteiro e conclui que um import ou um `catch` fora do hunk não existe. Os dois achados citam exatamente um símbolo definido fora do trecho alterado.

Por que importa: cada falso positivo estaciona o item como `corrective-spec-required`, e o operador precisa ler o código, dispensar e dar retry. Numa fila sem supervisão, isso para a entrega.

Sugestões: (1) passar ao reviewer o conteúdo completo dos arquivos tocados quando o diff for segmentado, ou pelo menos as linhas de import e o bloco que fecha um `try`; (2) depois de um merge da main que mantém a linhagem, revisar só o delta do item (a referência já descreve `lineage.deltaBase`, mas aqui abriu-se round 1 com o candidato inteiro); (3) rodar o typecheck do repositório antes do review e dar ao reviewer o resultado, para que ele não afirme "does not compile" contra um build que passa.

## Operator note (measured on Roundfix v0.139.0, 2026-11-11)

A temporary-repository reproduction with the review command's own test
fixtures and a fake reviewer, in the reported shape (an item that changes a
line near a `try` whose `catch` lies 30 lines below and appends a migration
journal entry; a default branch that appends its own entry and lands thirty
unrelated files; the operator's merge of the default branch resolving the
journal), shows two effects of one rule of ADR-0366, "When that merge
conflicts, the diff starts from S":

1. The review after the merge opens round 1 with the restart reason `the
   merge base moved from <prior base> to <base> and the reviewed head <prior
   head> conflicts with it`, and sends the whole candidate diff. A
   `pull-request-conflict` park means the reviewed head conflicts with the
   default branch, so every recovery from that park restarts the lineage.
2. That unsegmented round-1 diff starts at the new merge base and holds only
   the item's change (none of the thirty landed files), so the reporter's
   inference that the main Specs entered the diff does not hold for it. It
   does hold for a segmented review: a Review Segment that starts at an item
   commit conflicting with the default branch and ends after the merge is
   measured from that commit, and its diff carried all thirty landed files
   and the default branch's journal entry (33 files where the item changed 2).
3. The reviewer prompt already says "You may open repository files for
   context, but the session is read-only"; the `catch` and the import lie
   outside the hunk, and nothing asks the reviewer to read the file before
   asserting that something is missing.
