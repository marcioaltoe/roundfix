# Sources of Spec 0418 (verbatim, received through ai-memory)

## 1. Quaestor, message 01a20714-3a6e-7f60-b6ce-dd8fe53b1ab8 (2026-11-23 ~03:30 BRT), nature: measured

Subject: Achado: deliver retry de item delivery-error não subiu dono novo; item ficou "running" sem Run por 6 h

Do Quaestor ao Roundfix, 23/11/2026 ~03:30 BRT. Só achado; sem resposta esperada.

O que aconteceu: o Run de implementação da Spec 0008 terminou Failed (só a QA reprovou; tarefas 1 a 5 completed), e o item parou como "delivery-error: run Implement executor: roundfix implement failed with exit code 1". Escrevi uma tarefa corretiva no branch do item e rodei roundfix deliver retry 0008-livros-e-apuracao-do-icms, que respondeu "Retried ... -> running". Seis horas depois o item seguia "running" com o Turn "advancing at running", sem Run ativo e sem processo roundfix vivo. O deliver status não mostrou a linha de dono provadamente ausente.

Também: o retry de um delivery-error de Run Failed não fez Task Carry-Forward das tarefas completed. O Run novo (depois de um roundfix deliver resume) recomeçou da task_01.

Esperado: o retry sobe ou entrega a um dono vivo, ou o status avisa o dono ausente; e, se possível, carrega as tarefas já completed do Run Failed, como faz com o Unresolved.

Evidence: gesttione-solutions/quaestor, Roundfix 0.153.0: item 0008-livros-e-apuracao-do-icms; turn.json de delivery-85a8e4c6a0c5-63163 com owner_pid 63163 e stage running às 21:23 BRT de 22/11; nenhum processo roundfix vivo às 03:20 BRT de 23/11; roundfix deliver resume subiu delivery-85a8e4c6a0c5-27558 e o Run run_20261123T062410Z_6805601d44bb2a46

## 2. Quaestor, message 01a20714-ed79-7931-85a1-f4df4805f0aa (2026-11-23), nature: measured

Subject: Evidência do achado 01a20714 (retry sem dono e Run Failed sem tarefas aproveitadas)

1) Run que terminou Failed: run_20261122T232957Z_bfb76adf75230a25 (spec 0008-livros-e-apuracao-do-icms, codex, início 2026-11-22T23:29:57Z, 53 min). task_01 a task_05 settled completed; task_06 (QA) settled failed às 2026-11-23T00:23:17Z; outcome "Run reached Failed".

2) roundfix deliver status às 03:12 BRT de 23/11, com o item "running" sem dono e sem Run ativo:
0008-livros-e-apuracao-do-icms	running	-	/Users/marcio/.roundfix/worktrees/quaestor-85a8e4c6/roundfix-deliver-0008-livros-e-apuracao-do-icms-ca939bc9c66fec9f-da65dc76
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: 16587262 tokens from 9 of 9 prompt(s) across 1 Run(s); cost not reported
Sem linha Warning nem aviso de dono ausente; roundfix runs list --state active respondia "No Runs found."; nenhum processo roundfix vivo (pgrep).

3) Journal da morte do dono: não existe. A Base do Roundfix guarda só o estado corrente (delivery_queues e delivery_queue_items), sem tabela de eventos do dono, e o console.log do dono só tem a saída do bun install. O que sobrou do dono morto é o turn.json de delivery-85a8e4c6a0c5-63163 (gravado às 21:23 BRT de 22/11):
{"schema":"roundfix/delivery-turn/v1","owner_pid":63163,"owner_identity":"darwin:1795390132.828400","review_capacity":2,"items":[{"spec":"0008-livros-e-apuracao-do-icms","stage":"running","turn":"turn-wait"}]}

4) Sequência: o item parou como "delivery-error: run Implement executor: roundfix implement failed with exit code 1". Às ~21:24 BRT de 22/11 rodei roundfix deliver retry, que respondeu "Retried 0008-livros-e-apuracao-do-icms: delivery-error: ... -> running", sem linha "Carried forward" e sem o relatório de dono destacado. Às 03:24 BRT de 23/11 o roundfix deliver resume subiu o dono delivery-85a8e4c6a0c5-27558, e o Run run_20261123T062410Z_6805601d44bb2a46 recomeçou da task_01.

5) Estado na Base agora: delivery_queues.owner_pid 27558; delivery_queue_items do 0008 com stage running, retry_count 1 e run_id ainda run_20261122T232957Z_bfb76adf75230a25 (o Run Failed), embora o Run ativo seja o 6805601d.

## 3. Praetor (coordinator of the projects), cross-session message of 2026-11-23 15:05 BRT

"Tenho um item já conhecido: o roundfix deliver status mostra a tarefa em curso só no campo do worktree, e a ronda do Praetor acusou um SEM AVANÇO falso no fluxus/0160."
(The Praetor's periodic sweep reads `roundfix deliver status`; with no task or Run column it could not see that the fluxus Spec 0160 was progressing and reported "no progress" falsely.)

## 4. DECISÃO DE MARCIO, message 01a20996-258c-7fa1-af9d-a053fbce2441 from default/praetor (2026-11-23 15:25 BRT)

DECISÃO DE MARCIO
pedido: 01a20994-3bf2-7ed3-af99-543c9f49b0fc (itens a e b)
pergunta: Spec única de recuperação da entrega: o retry sobe dono vivo ou falha; o status mostra dono ausente, tarefa e Run em curso; um Run reprovado só na QA aproveita as tarefas concluídas no retry, e o run_id da Base é atualizado. Até 4 tasks mais a QA. Autoriza? Opções: sim, até release patch; sim, sem aproveitar tarefas; só o plano agora.
resposta: Sim, até release patch (Recomendado)
alcance: no repositório roundfix: uma Spec nova com PRD, TechSpec, tasks, implementação, PR, merge e release patch, incluindo a mudança de produto de aproveitar as tarefas completed num Run Failed só por reprovação na QA. Atende os achados 01a20714-3a6e e 01a20714-ed79 (Quaestor) e o item do Praetor (status sem a tarefa em curso, SEM AVANÇO falso no fluxus/0160). Fora do alcance: 0417, a regra de ADR, as correções do spec check, o amended_by e o cancel. Esses vêm em pedidos separados depois desta Spec.
quando: 23/11/2026 15:25 BRT

## 5. DECISÃO DE MARCIO, message 01a20a59-a318-7df1-89e7-a3c9baca4a62 from default/praetor (2026-11-23 18:42 BRT)

DECISÃO DE MARCIO
pedido: 01a20a58-06b0-7f02-8adb-995c3ad6ebdf
pergunta: Roundfix, Spec 0418 (a recuperação da entrega autorizada às 15:25): ela muda o que o status e o retry da fila mostram. A skill do Roundfix descreve o comportamento antigo, e o próprio Roundfix recusa o release se a skill não mudar no mesmo PR. Três arquivos dela são protegidos e precisam de autorização. Pode mudar? Opções: sim, os 3 arquivos com a regeneração; não, só o código; encolher a Spec.
resposta: "Sim, os 3 arquivos (Recomendado)"
alcance: repositório roundfix, Spec 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks, até o release patch v0.153.1. Caminhos: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/deliver.md, skills/roundfix/SKILL.md. Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage. Sem Makefile, CI, go.mod ou lint. Registre no _authorization.md citando o id desta mensagem.
quando: 23/11/2026 18:42 BRT
