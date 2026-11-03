---
type: feat
status: promoted
created: 2026-11-02
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
reason: null
---

# The follow-up names what waits on someone else and where each item went

## Opportunity

A `/roundfix --follow-up` run lists the repository's open Backlog Entries and
unresolved Findings, and its agent decides which become Specs. When an item
waits on something outside the session (a vendor's answer, real data from a
supplier, a maintainer decision), nothing records that wait: the agent leaves
the item out of the plan and the report says why in free text, if at all. The
`repository-inbox` source also lists `docs/_inbox/.gitkeep` as a `note` item.
Reported by fluxus on roundfix 0.129.0.

## Value

A recorded wait survives the session that decided it: the next follow-up
knows the item is not work to plan, and the maintainer can see why a Backlog
Entry or a Finding did not become a Spec. A report that names where every
listed item went can be checked, as the consumed messages already are. A
folder marker stops costing a judge call and a line of noise.

## Shape

Optional frontmatter fields on a Finding and a Backlog Entry, for example
`waiting_on: <who or what>` and `since: <date>`, read by the follow-up into a
class of their own; the follow-up mode's report names each listed item's
destination; the inbox source skips hidden files. The maintainer decided on
2026-11-03 that waiting items stay out of the list by default and that a flag
such as `--all` lists everything.

## Evidence

The Fluxus message, ai-memory message
`01a19bf3-cf55-7c62-866c-3e14b678f9cb` from `default/fluxus`, received on
2026-11-02 (11:10 UTC, 08:10 BRT), reproduced verbatim below with its
recorded evidence line:

> evidence: "fluxus, 02/11/2026 ~08:05 BRT: roundfix 0.129.0 follow-up --format json listou finding x3, backlog x1 e note '.gitkeep' (repository-inbox); nenhum virou Spec e o relatório não registrou por quê. O mantenedor perguntou se backlog e findings não deviam virar Specs."
>
> Pedido do Fluxus ao Roundfix (0.129.0), a partir de um `/roundfix --follow-up` de 02/11/2026.
>
> O que aconteceu: `roundfix follow-up --format json` listou 3 findings, 1 Backlog Entry e um item `note` chamado `.gitkeep` (fonte `repository-inbox`). O agente não transformou nenhum deles em Spec, porque três esperam algo de fora: a resposta do suporte de um ERP, um dado real de um fornecedor e uma decisão do mantenedor. O relatório final só disse isso em texto livre. O mantenedor perguntou por que backlog e findings não viraram Specs, e nada no comando nem no modo deixava a razão registrada ou verificável.
>
> Pedidos:
> 1. Um campo declarado em Finding e Backlog Entry para "aguarda algo externo", por exemplo `waiting_on: <quem/o quê>` e `since: <data>` no frontmatter. O follow-up mostraria esses itens numa classe própria (`waiting`), fora do que vira Spec, com o motivo e a data. Hoje a decisão de não autorar fica só no julgamento do agente.
> 2. No modo `--follow-up`, exigir que o relatório diga, para cada item `finding` ou `backlog` listado, o que aconteceu com ele: virou Spec X, foi agrupado na Spec Y, ou está à espera de Z (o campo do item 1). Item sem destino seria lacuna do relatório, como já acontece com as mensagens consumidas.
> 3. A fonte `repository-inbox` conta `.gitkeep` como item `note`, com prioridade 4. Arquivos ocultos ou de marcação de pasta deveriam ser ignorados.
>
> Nenhuma urgência: o Fluxus contorna registrando a espera no próprio finding. Se aceitar o item 1, diga o nome do campo para usarmos.
