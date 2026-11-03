---
type: feat
status: promoted
created: 2026-11-03
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
reason: null
---

# The Spec skills announce each Spec event to the destination the user defines

## Opportunity

A user who tracks deliveries and deadlines outside the repository, on a
board, a tracker or a personal assistant, learns that a Spec was created,
authorized, started, blocked, given a deadline, completed or cancelled only
by reading the repository. The Spec skills that perform those steps already
know each moment and its facts. Reported by jarvis, which keeps such a board
for the maintainer and asked every project to send it one block per Spec
event.

## Value

Each Spec event reaches the user's board when it happens, from the session
that made it happen, with no project writing to the board itself and no step
waiting for a reply. Adopters without such a board see no change.

## Shape

A step in `write-prd`, `write-techspec`, `write-tasks`, the Roundfix skill's
implement reference and `archive-spec` that composes an announcement at each
event and sends it where the user's own instructions say. The skills name no
destination and no format: the maintainer chose a generic step over the
assistant-specific block the report asks for ("Passo genérico nas skills
(Recommended)", 2026-11-03), and the block and its destination go to the
user's own global instructions. Non-binding.

## Evidence

The report, ai-memory message `01a1a1ec-277d-7871-901b-fe53531a2216` from
`default/jarvis`, received 2026-11-03, subject "Quadro de specs no Todoist:
bloco SPEC PARA O QUADRO e ajuste nas skills de spec", as the operator
reproduced it (verbatim below):

> Pedido de Marcio em 03/11/2026, na sessão do Jarvis.
>
> Desde 03/11/2026 o quadro "Gesttione Solutions" do Todoist é o controle de entregas e prazos das specs, gerido pelo Jarvis. Este projeto não edita o Todoist. A cada mudança de spec, mande para default/jarvis, por memory_message_send, um bloco neste formato (pode ir mais de um bloco por mensagem):
>
> ```text
> SPEC PARA O QUADRO
> projeto: <nome do produto>
> spec: <NNNN>
> título: <título da spec>
> evento: nova | autorizada | iniciada | prazo proposto | bloqueada | concluída | cancelada
> data do evento: DD/MM/AAAA
> início: DD/MM/AAAA ou -
> prazo proposto: DD/MM/AAAA ou -
> caminho: docs/specs/<slug> ou docs/history/specs/<slug>.md
> cliente: Grupo … ou nenhum
> resumo: <1 a 3 frases, só em concluída ou cancelada, escritas para o empresário, sem termos técnicos nem dados de outros clientes>
> ```
>
> Quando mandar: ao criar a spec (nova), quando a autorização for aprovada (autorizada), quando a primeira task começar (iniciada), quando houver prazo real, de cliente, de lei ou de medição (prazo proposto: o prazo só vale depois da aprovação de Marcio), se a spec travar (bloqueada) e ao arquivar (concluída, ou cancelada quando a disposição for cancelled). Não é preciso esperar retorno. Se o repositório tiver números de spec repetidos, informe o slug completo no campo spec.
>
> As specs já existentes deste projeto foram carregadas no quadro em 03/11 (ativas como cartões abertos; arquivadas desde 13/10 como cartões concluídos; anteriores num cartão-histórico). Daqui em diante basta mandar os eventos novos.
>
> Específico do Roundfix: as specs do próprio Roundfix não ganham cartão; o Jarvis registra cada evento só como comentário no cartão-histórico do Roundfix. Mande o bloco do mesmo jeito.
>
> Segundo pedido, para as skills que o Roundfix distribui (write-prd, write-techspec, write-tasks, roundfix/implement e archive-spec), seguindo a curadoria de skills (~/.vault/skills/curation.md): incluir o envio automático do bloco acima para default/jarvis nos pontos certos — write-prd/write-techspec ao criar a pasta da spec (evento: nova, início = data de criação); write-tasks com _authorization aprovada (evento: autorizada; prazo real conhecido vai como prazo proposto); implement quando a primeira task entra em andamento (evento: iniciada) e em Run bloqueado (evento: bloqueada); archive-spec depois de gravar o registro (evento: concluída, ou cancelada quando a disposição for cancelled, com resumo de 1 a 3 frases para o empresário e o caminho em docs/history/specs/). O envio não espera retorno e o projeto nunca escreve no Todoist. Avise default/jarvis quando a versão com isso estiver publicada.
