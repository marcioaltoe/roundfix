---
spec: 0079-one-door-for-fleet-knowledge
created: 2026-08-06
status: promoted
leverage: compounding
---

# One door for fleet knowledge

## Overview

Todo conhecimento da frota passa a nascer por uma única porta — o inbox do
secondbrain — e a fluir por um ciclo com fim: triagem no projeto de destino,
Finding ou Backlog Entry no repo certo, rollup que consolida, arquivamento que
licencia, mirror que devolve ao brain. Resolve o acúmulo medido (63 findings no
roundfix, 111 na frota), o transporte cross-project que hoje é humano, e a
janela de perda entre observar e versionar. Usuários: as sessões-agente de
todos os projetos e o mantenedor como curador. V1 ambicioso-híbrido: o ciclo
completo em regras e skills, a captura automática de fim de sessão, e a
aquisição de conhecimento de pesquisa entrando pela mesma porta — nenhuma
ferramenta nova de CLI.

## Problem

Findings acumulam mais rápido do que viram Specs: 63 arquivos/500K no
roundfix (28 em julho inteiro, 35 nos primeiros 6 dias de agosto), 111 em 8
projetos. Nenhum dos 63 carrega o lifecycle que o contrato ganhou em
2026-08-06 (Spec 0075). Sem rollup, sem arquivamento, sem enforcement — o
ciclo vicioso que o SRE workbook nomeia: postmortems escritos e nunca
fechados.

Feedback entre projetos viaja por humano (vortex→roundfix e fiscus→roundfix
nesta mesma sessão, por paste) e o ride-along colide com Runs ativos por
construção — interferência medida três vezes numa noite: edições de findings
esperando checkout livre e um resolve recusado por arquivo sujo.

A captura não tem durabilidade até a triagem: uma observação vive só na
conversa até alguém commitá-la — e o commit é justamente o que o Run ativo
impede. O brain, que deveria ser o destino do conhecimento, recebe hoje a
pilha crua via mirror, sem camada de captura nem curadoria de entrada.

### Market Data

- Prior art interno: o roteiro do vortex (2026-07-27) é um rollup inventado
  ad hoc — 30+ findings sequenciados em fases, com status.
- Convergência externa: Capture→Curation→Retention com "um bom sistema
  rejeita a maioria das entradas" (fitgap); repositório central + rollups
  periódicos + action items no tracker (Google SRE); fleeting→permanent→
  archived com revisão advisória e "arquivado ≠ apagado" (open-zk-kb,
  memento-vault); Products over Projects (Communication Patterns, cap. 10).
- Evidência completa:
  [findings accumulate faster than they become Specs](references/2026-08-06-findings-accumulate-faster-than-they-become-specs.md).

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | Inbox no brain | Critical | `inbox/<destino>/` no secondbrain é A porta de captura da frota. Entrada nasce comitada (durável em <1 min), com `origin`, `destination`, `type-hint`, `created_at`. Destinos válidos = `projects/<p>/` existentes; zero config nova. |
| F2  | Triagem no destino | Critical | Sessão do destino lê seu inbox, converte cada entrada em Finding (evidência) ou Backlog Entry (intenção) no próprio repo — o commit é sempre do destino — e move a entrada para `_triaged/` com `resolved_to`. Re-roteamento = mover entre inboxes. |
| F3  | Regra de criação de arquivos nos agents | Critical | O fluxo vira cláusula mandatória no guia de layout e no módulo Baseline — a frota herda. Fallback `docs/_inbox/` documentado uma única vez, só para adotantes sem brain. `docs/_inbox/` da frota aposenta. |
| F4  | Rollup de primeira classe | Critical | Artefato tipado que consolida findings relacionados, supersede seus membros e licencia o arquivamento deles. É a "primeira página" humana do conjunto. |
| F5  | Arquivamento de findings | High | `docs/findings/_archived/` simétrico ao de specs; um finding arquiva quando `done`/`deferred` E absorvido por rollup ou Spec. Mirror preserva — o brain não perde história. |
| F6  | Varredura dos 63 legados | High | Uma passada única, rollup-primeiro: escreve o(s) primeiro(s) rollup(s), carimba status, arquiva o absorvido; o que nenhum rollup referencia default `deferred` — revisável no diff. Alvo: ≤15 ativos. |
| F7  | Skills roteiam pela porta | High | `write-idea`/`write-prd` começam checando o inbox pendente do repo; finding fechado spawna Backlog Entries tipadas; `knowledge-workspace` carrega templates e o ato atômico escrever+comitar. |
| F8  | Captura automática de fim de sessão | High | Ao fim de cada sessão, um capture rascunha entradas de inbox para as observações que a sessão expôs — marcadas `auto`, sempre pendentes de triagem, nunca auto-triadas. |
| F9  | Enforcement mecânico | High | As regras viram checagens na infra SC-* da 0065: lifecycle presente em finding ativo, rollup com membros resolvíveis, arquivado só com licença, entrada de inbox válida (lint no brain). Regra sem checagem apodrece — 63 provas. |
| F10 | Vocabulário no glossário | Medium | `Inbox Entry`, `Rollup` e `Triage` nascem no CONTEXT.md ao lado de Finding e Backlog Entry, com seus _Avoid_. |
| F11 | Captura de pesquisa para o brain | High | Sessão que fez pesquisa substantiva (web/livros) captura o digest com fontes em `inbox/secondbrain/` — o brain é o destino quando o conhecimento é da frota, não de um projeto. Check advisório via `qmd` antes de incluir: acerto forte → a entrada referencia e estende a página existente em vez de duplicar. A triagem desse namespace é a ingestão que o brain já possui, rodando em sessão do próprio brain. |

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Findings ativos pós-varredura | ≤ 15 | `ls docs/findings/*.md \| wc -l` |
| Findings ativos com lifecycle | 100% | regra SC-* |
| Captura → durabilidade | < 1 min | auditoria do piloto |
| Inbox pendente > 14 dias | 0 entradas | `ls` por data no brain |
| Entradas com origin+destination válidos | 100% | lint no brain |
| Interferência com Runs (preflight/Batch) | 0 ocorrências | Run journals |
| Pesquisa substantiva capturada com fontes | 100% das sessões que pesquisaram | auditoria do piloto |

## Feature Assessment

| Criteria        | Question                                            | Score   |
| --------------- | --------------------------------------------------- | ------- |
| Impact          | How much more valuable does this make the product?  | Must do |
| Reach           | What % of users would this affect?                  | Must do |
| Frequency       | How often would users encounter this value?         | Must do |
| Differentiation | Does this set us apart or just match competitors?   | Strong  |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong  |
| Feasibility     | Can we actually build this?                         | Strong  |

## Integration with Existing Features

| Integration point | How |
| --- | --- |
| Spec 0075 (backlog) | alvo do spawn da triagem; tipos e lifecycle intactos |
| Spec 0065 (SC-*) | trilho do enforcement de F9 |
| Mirror sync existente | fecha o ciclo: artefato triado → merge → mirror |
| Skill knowledge-workspace | dona da mecânica da porta |
| Módulos Baseline | levam F3 à frota (precedente: 0075/0065) |
| Convenção dogfood `docs/_inbox` | migra para a porta e aposenta |

## Council Insights

- **Recommended approach:** original refinada; o mantenedor escolheu o híbrido
  com F8.
- **Key trade-offs:** porta única vs custo do salto — só funciona se a porta
  for mais barata que o desvio (helper de uma linha, commit silencioso);
  brain vira load-bearing da captura — fallback definido uma vez, no Baseline.
- **Risks identified:** rot de inbox sem cadência (mitigado: leitura amarrada
  ao ritual de write-idea/write-prd); varredura em escala classificando
  errado (mitigado: rollup-primeiro, default deferred revisável); mudança no
  repo do brain é dependência externa de spec roundfix — precedente novo,
  escopar explicitamente; F8 amplia a superfície do primeiro corte.
- **Dissenting view:** devils-advocate manteria F8 fora do V1 — "a porta
  ainda não provou ser mais barata que o desvio; automatizar a captura antes
  do piloto é otimizar o que não se mediu". Registrado; o mantenedor decidiu
  incluí-la.
- **Stretch goal (V2+):** graduação automatizada para o wiki — rollup que
  arquiva membros produz síntese candidata no brain; exige contrato de
  escrita próprio.

## Opportunity Scan

- Mais ambiciosa — Knowledge OS completo (captura auto + rollup auto + wiki):
  teto máximo, Feasibility Maybe; fatiada — F8 entrou, resto V2.
- Mais simples — só regras (lifecycle+rollup+archive): Pass — falha
  durabilidade e cross-project, metade do pedido.
- Adjacente — ferramenta primeiro (`roundfix inbox`): Maybe — ordem errada;
  contrato antes de ferramenta.

**Chosen direction:** hybrid — original refinada + F8 (captura automática),
por decisão do mantenedor com o dissenso registrado.

## Out of Scope (V1)

- **Comando `roundfix inbox` no CLI** — contrato antes de ferramenta; skills e
  regras provam a forma primeiro.
- **Graduação automatizada para o wiki** — exige contrato de escrita no brain;
  é o stretch V2.
- **Transporte remoto (issue forms, outra máquina)** — o V1 é local-first;
  caso real só quando existir.
- **Mudanças no mirror/.secondbrain-export** — o export segue intocado; o
  ciclo fecha com o sync existente.
- **Re-tipagem do backlog da 0075** — tipos e boundary ficam como estão.
- **Automação de varredura/retenção no brain** (`_triaged/` sweep) — política
  declarada, execução manual até doer.
- **Mecânica de ingestão ao wiki** — F11 alimenta a fila; o pipeline de
  ingestão e suas regras são e continuam sendo contrato do repo do brain.

## Decisions

- O secondbrain é a porta única de captura da frota; projetos guardam só
  artefatos triados. See ADR-0095.
- Só o inbox vive no brain; findings e backlog ficam nos repos de destino —
  o pipeline de spec não depende de repo externo (lição conexus).
- Commit sempre do destino; origem é frontmatter, nunca autoria de commit.
- Rollup licencia arquivamento; status sozinho não arquiva.
- A regra do fluxo de criação de arquivos entra nos agents e no Baseline.
- Zero configuração nova no V1; `projects/<p>/` é o registro de destinos.
- F8 (captura automática) entra no V1 por decisão do mantenedor, contra o
  dissenso do conselho — registrado dos dois lados.
- Aquisição de conhecimento é fundamento do brain: pesquisa de sessão captura
  no `inbox/secondbrain/` com check advisório de duplicata via `qmd` antes da
  inclusão; a ingestão ao wiki permanece contrato do próprio brain. (Adição
  do mantenedor, 2026-08-06 — a pesquisa desta mesma ideia é o caso-exemplo
  e o primeiro material do piloto.)

## Open Questions

- Cadência-limite da triagem: 14 dias como alvo de KPI — endurecer para
  cláusula? Default: KPI apenas, até o piloto medir.
- A captura automática roda em toda sessão ou opt-in? Default: toda sessão,
  entradas marcadas `auto`.
- Coleção qmd para o inbox: nome/escopo. Default: seguir o contrato do repo
  do brain quando a mudança for feita lá.
- Alias origem↔nome quando basename divergir. Default: inexistente até
  precisar.
- Limiar de "pesquisa substantiva" que obriga captura F11. Default: qualquer
  sessão cujo artefato de decisão citou fontes externas captura o digest.
