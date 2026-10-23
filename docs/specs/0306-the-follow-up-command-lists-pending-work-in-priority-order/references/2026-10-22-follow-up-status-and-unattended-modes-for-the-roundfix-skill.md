---
type: feat
status: promoted
created: 2026-10-22
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
reason: null
---

# Follow-up, status and unattended modes for the Roundfix skill

## Opportunity

The maintainer asked, through `default/jarvis` on 2026-10-22 at 15:09 and
15:12 BRT, for one way to pick up a project's pending work from any agent
harness: list everything that waits (the Secondbrain inbox, ai-memory
messages, the open Backlog and Findings, pending Specs) next to what is
already in progress, order it, settle every open decision and authorization
in one round of questions, and carry the work through Specs, the Delivery
Queue and the releases. The operator's proposal below, sent to
`default/jarvis` and decided with the maintainer at 15:20 BRT, splits it into
a read-only CLI command (Spec A) and the skill modes plus a Baseline clause
change (Spec B).

## Value

A session that starts with `--follow-up` stops rebuilding the same inventory
by hand in each project, the order is the same in Claude Code, Codex and
OpenCode, and the maintainer answers the decisions once instead of being
asked one question at a time across the session.

## Shape

The proposal, translated faithfully from the Portuguese original quoted at
the end of this entry.

### Proposal: `--follow-up`, `--status` and `--unattended` modes in the roundfix skill

Marcio's request of 22/10/2026, received from `default/jarvis` at 15:09 and
15:12 BRT. The architecture decisions were taken with Marcio at 15:20 BRT.

#### Marcio's decisions

1. **Where it lives.** The modes go into the product's roundfix skill. Each
   source is optional and detected:
   - the Secondbrain, through the Baseline guide;
   - ai-memory, if the MCP is present;
   - Jev, if there is a key;
   - the RELEASE PARA O GT3 block, only if the user's instructions define
     it.

   Without a source, the mode continues and says what was missing.
2. **Code and skill.** A read-only CLI command gathers and orders. The skill
   talks, asks and acts.
3. **Questions.** Up to 3 per call, in the whole Baseline. The clause
   `clause.core.ask-user-answerable-decisions`, which today requires one
   question at a time, changes for every adopter. The limit of 3 is Codex's
   `request_user_input` limit.
4. **Where the answers go.** Answers and authorizations go into the
   `_authorization.md` of the Spec that will be implemented. An item that
   does not become a Spec becomes a Backlog Entry or Finding, with the
   information and decisions and without the authorizations, to become a Spec
   later.

#### How each harness calls the skill

The research was done in Exa on 22/10 and confirmed in the official docs.

| Harness | Call | Question tool |
|---|---|---|
| Claude Code | `/roundfix --follow-up`; the text arrives as `ARGUMENTS` | `AskUserQuestion` |
| Codex | `$roundfix --follow-up`; the text goes along in the prompt | `request_user_input`, up to 3 questions |
| OpenCode | `@roundfix --follow-up`, or `/roundfix` with `metadata.opencode/slash: "true"` | `question` |

The skill reads the mode from the message text. The names of the question
tools are already in the Baseline (`docs/agents/agent-instructions.md`).

#### Spec A (CLI): `roundfix follow-up`, read-only

**Project.** The project comes from the current repository. The inbox
namespace is configurable and, by default, is the repository name.

**Inventory.** The command lists, once, with origin, path or id, type, date
and the project each item blocks:
- the pending Inbox Entries of the Secondbrain, oldest first;
- the ai-memory messages the skill passes in a JSON file, because only the
  agent sees the MCP;
- the open Backlog and Findings;
- the pending Specs.

Next to the list, it shows what is in progress: active Runs, the Delivery
Queue, the open PRs and the git state.

**Stop rule.** Empty list and nothing in progress: it exits with "nothing to
do" and says what it checked.

**Priority.** Jev answers, in English, four questions per item: is it a fix?
does it block another project? is it small? does it affect production? It
uses the existing `internal/judge`, with the pinned model and the model id
written to the log. The code orders:
1. fixes and small items that block another project;
2. the good functioning of the system;
3. features with added value;
4. the rest.

Without Jev, the order falls back to the frontmatter `type` and the report
says Jev was missing.

**`--status`.** The same command shows what is in progress, the pending Specs
in delivery order and the next release (through `release plan`, with what it
brings). It also shows the blockers (parked items) and what depends on Marcio
(the queue's pending questions).

**Output.** Text and `--format json`.

#### Spec B (skill and Baseline)

**Modes in the roundfix skill.** SKILL.md gains the modes section, and each
mode has its reference.

**`--follow-up`.** The steps are:
1. run the command;
2. show a short plan, by stages;
3. call `--unattended`;
4. after approval, triage the inbox with `resolved_to` or
   `discarded_reason`;
5. write the Specs, consulting the Secondbrain and Exa before the PRD and
   TechSpec;
6. run `spec check --strict`;
7. enqueue in `roundfix deliver` and follow through to the releases.

Before ending the turn, it leaves a detached Run, the queue or a wake-up
armed. It records in the inbox the findings of other projects and sends the
report.

**`--status`.** Short report, read-only.

**`--unattended`.** Gathers everything that is missing (information,
decisions and authorizations for push, merge, release and the GT3 block) and
asks at once, up to 3 questions per call, each with the recommended option.
Without a structured tool, it asks in text and waits. It records the answers
according to decision 4 and arms the continuation before ending.

**Baseline clause.** Goes from 1 to up to 3 questions per call.

#### Order and delivery

Cycle AP, after AO. AO ends with 0300 and v0.111.0, expected at ~21:50 BRT on
22/10. Spec A ships in one minor and Spec B in the next. Implementation only
starts after this proposal goes to `default/jarvis`.

### Original (Portuguese)

> # Proposta: modos --follow-up, --status e --unattended na skill roundfix
>
> Pedido de Marcio de 22/10/2026, recebido do default/jarvis às 15:09 e 15:12 BRT. As decisões de arquitetura foram tomadas com Marcio às 15:20 BRT.
>
> ## Decisões de Marcio
>
> 1. **Onde mora.** Os modos entram na skill roundfix do produto. Cada fonte é opcional e detectada:
>    - o Secondbrain, pelo guia do baseline;
>    - o ai-memory, se o MCP estiver presente;
>    - o Jev, se houver chave;
>    - o bloco RELEASE PARA O GT3, só se as instruções do usuário o definirem.
>
>    Sem uma fonte, o modo segue e diz o que faltou.
> 2. **Código e skill.** Um comando de CLI só de leitura levanta e ordena. A skill conversa, pergunta e age.
> 3. **Perguntas.** São até 3 por chamada, no baseline todo. A cláusula `clause.core.ask-user-answerable-decisions`, que hoje exige uma pergunta por vez, muda para todos os adotantes. O limite de 3 é o do `request_user_input` do Codex.
> 4. **Onde ficam as respostas.** Respostas e autorizações vão para o `_authorization.md` da Spec que será implementada. Um item que não vira Spec vira Backlog Entry ou Finding, com as informações e decisões e sem as autorizações, para virar Spec depois.
>
> ## Como cada harness chama a skill
>
> A pesquisa foi feita no Exa em 22/10 e confirmada nas docs oficiais.
>
> | Harness | Chamada | Ferramenta de perguntas |
> |---|---|---|
> | Claude Code | `/roundfix --follow-up`; o texto chega como `ARGUMENTS` | `AskUserQuestion` |
> | Codex | `$roundfix --follow-up`; o texto vai junto no prompt | `request_user_input`, até 3 perguntas |
> | OpenCode | `@roundfix --follow-up`, ou `/roundfix` com `metadata.opencode/slash: "true"` | `question` |
>
> A skill lê o modo do texto da mensagem. Os nomes das ferramentas de perguntas já estão no baseline (`docs/agents/agent-instructions.md`).
>
> ## Spec A (CLI): `roundfix follow-up`, só leitura
>
> **Projeto.** O projeto sai do repositório atual. O namespace do inbox é configurável e, por padrão, é o nome do repositório.
>
> **Levantamento.** O comando lista uma vez só, com origem, caminho ou id, tipo, data e o projeto que cada item trava:
> - as Inbox Entries pendentes do Secondbrain, mais antigas primeiro;
> - as mensagens do ai-memory que a skill passa num arquivo JSON, porque só o agente enxerga o MCP;
> - o Backlog e os Findings abertos;
> - as Specs pendentes.
>
> Ao lado da lista, mostra o que está em andamento: Runs ativas, a Delivery Queue, os PRs abertos e o estado do git.
>
> **Regra de parada.** Lista vazia e nada em andamento: sai com "nada a fazer" e diz o que conferiu.
>
> **Prioridade.** O Jev responde, em inglês, quatro perguntas por item: é correção? trava outro projeto? é pequeno? afeta produção? Usa o `internal/judge` que já existe, com o modelo fixado e o model id gravado no log. O código ordena:
> 1. correções e itens pequenos que travam outro projeto;
> 2. o bom funcionamento do sistema;
> 3. features com valor agregado;
> 4. o resto.
>
> Sem o Jev, a ordem cai para o `type` do frontmatter e o relatório diz que o Jev faltou.
>
> **`--status`.** O mesmo comando mostra o que está em andamento, as Specs pendentes na ordem de entrega e a próxima release (pelo `release plan`, com o que ela traz). Mostra também os bloqueios (itens estacionados) e o que depende de Marcio (perguntas pendentes da fila).
>
> **Saída.** Texto e `--format json`.
>
> ## Spec B (skill e baseline)
>
> **Modos na skill roundfix.** O SKILL.md ganha a seção de modos, e cada modo tem sua reference.
>
> **`--follow-up`.** Os passos são:
> 1. rodar o comando;
> 2. mostrar um plano curto, por etapas;
> 3. chamar o `--unattended`;
> 4. depois da aprovação, triar o inbox com `resolved_to` ou `discarded_reason`;
> 5. escrever as Specs, consultando Secondbrain e Exa antes de PRD e TechSpec;
> 6. rodar `spec check --strict`;
> 7. enfileirar no `roundfix deliver` e seguir até as releases.
>
> Antes de encerrar o turno, deixa uma Run destacada, a fila ou um wake-up armado. Registra no inbox os achados de outros projetos e manda o relatório.
>
> **`--status`.** Relatório curto, só leitura.
>
> **`--unattended`.** Levanta tudo o que falta (informações, decisões e autorizações de push, merge, release e do bloco do GT3) e pergunta de uma vez, até 3 perguntas por chamada, cada uma com a opção recomendada. Sem ferramenta estruturada, pergunta no texto e espera. Grava as respostas conforme a decisão 4 e arma a continuação antes de encerrar.
>
> **Cláusula do baseline.** Passa de 1 para até 3 perguntas por chamada.
>
> ## Ordem e entrega
>
> Ciclo AP, depois do AO. O AO termina com a 0300 e a v0.111.0, prevista para ~21:50 BRT de 22/10. A Spec A sai numa minor e a Spec B na seguinte. A implementação só começa depois desta proposta ir ao default/jarvis.
