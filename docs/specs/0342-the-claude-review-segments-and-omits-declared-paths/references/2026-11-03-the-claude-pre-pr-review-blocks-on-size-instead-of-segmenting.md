---
type: feat
status: promoted
created: 2026-11-03
spec: 0342-the-claude-review-segments-and-omits-declared-paths
reason: null
---

# The Claude pre-PR review segments a large candidate and omits declared paths instead of blocking on size

## Opportunity

An adopter whose pre-PR review runs on Claude loses a delivery every time a
candidate carries one large generated file: the review refuses the prompt by
size, the Delivery Queue parks the item as `review-blocked`, and the only
remedy is to declare the file a derived path with a `regenerate` command that
does nothing. Reported by vortex (gesttione-solutions/vortex) on roundfix
0.129.x, ai-memory message `01a1a0ea-ec31-7512-a83c-d329647cad29` from
`default/vortex`, received 2026-11-03.

## Value

A Spec that adds an ORM migration snapshot of about 670 KB passes the Claude
review without manual work, and the Pre-PR Review Record says what was not
reviewed and why. A candidate that still does not fit names its largest files
and the key that would leave them out, so nobody has to open the record and
run `git diff --numstat` by hand.

## Shape

Non-binding. The maintainer chose (AskUserQuestion, 2026-11-03): Claude also
splits round 1 into at most four Review Segments under its token budget; a new
`pre_pr_review.omit_paths` key (globs) keeps files out of the review without
`delivery.derived_paths`; when the review still blocks, the message names the
three largest files and the key that would remove them.

## Report (verbatim)

The message as received, frontmatter included:

```markdown
---
source: ai-memory message 01a1a0ea-ec31-7512-a83c-d329647cad29
from: default/vortex (gesttione-solutions/vortex)
received: 2026-11-03
subject: "Sugestão: revisão pré-PR do Claude não devia travar só pelo tamanho do diff"
---


Sugestão de melhoria do Roundfix, pedida por Marcio (mantenedor) em 03/11/2026, vinda do repositório gesttione-solutions/vortex. Já não é a primeira vez que uma entrega para porque o pedido de revisão passou do limite.

## O que aconteceu
- Roundfix 0.129.x, `pre_pr_review.provider: claude`. A Delivery Queue estacionou o item `0168-a-auditoria-das-fichas-tecnicas` como `review-blocked`. Motivo no `pre-pr-review.json`: "review prompt too large: 541455 estimated tokens after omitting 60 path(s) exceeds the claude review budget of 500000 tokens, half of its 1000000-token context window".
- O diff tinha 1.049.480 bytes, e 669 KB eram um único arquivo gerado, o snapshot do Drizzle (`packages/backend/drizzle/migrations/meta/0092_snapshot.json`). Não havia nenhum achado de código.
- Casos anteriores no mesmo repositório (01/10/2026, provider codex): um diff de ~1,45 MB falhou como "agent/protocol error"; noutro, specs foram puladas (`specContextTruncated`). O contorno manual de sempre é revisar com `--base` depois do commit grande, fora do trilho automático.

## Contorno adotado agora
- O repositório declarou `packages/backend/drizzle/migrations/meta/*_snapshot.json` em `delivery.derived_paths`, com um `regenerate` que sai com 1 de propósito. A única intenção era omitir o arquivo da revisão. Funciona, mas mistura dois conceitos: o "derivado para resolver conflito de merge" passou a significar também "fora da revisão", e o repositório teve de inventar um `regenerate` que só recusa.

## O que sugerimos (escolham o desenho)
1. **Segmentar também com Claude.** Hoje só o codex divide em até 4 Review Segments, e o Claude "is never segmented". Estender a segmentação por commit ao orçamento de tokens do Claude resolveria a maioria dos casos.
2. **Omitir da revisão sem `regenerate`.** Uma chave própria, por exemplo `pre_pr_review.omit_paths` com globs, ou uma detecção de arquivo gerado (snapshots de ORM, lockfiles, `linguist-generated` do `.gitattributes`). Ela deixaria o arquivo fora da revisão sem tocar no tratamento de conflito.
3. **Degradar em vez de bloquear.** Acima do orçamento, omitir primeiro os maiores arquivos gerados ou de dados, listá-los em `omittedPaths` com um motivo (`oversize`), e só bloquear se ainda não couber.
4. **Mensagem acionável.** Quando bloquear, a mensagem do park poderia nomear os 3 maiores arquivos do diff e a configuração que os tiraria. Hoje é preciso abrir o JSON e rodar `git diff --numstat` à mão.

## Resultado esperado
Uma Spec com uma migration nova (snapshot de ~670 KB) passa pela revisão pré-PR do Claude sem intervenção manual. O registro deixa claro o que não foi revisado e por quê.

Não precisa responder a este projeto. Se mudar o formato de configuração, o vortex adota no próximo `baseline update`.
```
