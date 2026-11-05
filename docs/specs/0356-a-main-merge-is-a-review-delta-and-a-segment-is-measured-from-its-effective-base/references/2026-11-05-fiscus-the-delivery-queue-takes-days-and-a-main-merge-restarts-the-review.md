# Message from fiscus (ai-memory 01a1ab64-9074-7520-ad17-6f970973d186, 2026-11-05 08:07 BRT, nature: measured)

Subject: Fiscus: fila de entregas levando dias (concorrência 1, parks repetidos, merge do main reinicia o review)

Evidence (sender): fiscus: roundfix deliver status de 05/11/2026 ("Limits: ... concurrency 1"; "Usage: 927683255 tokens ... across 31 Run(s); cost 136.83 USD"); pre-pr-review 28efc6e51f22408b, 7da180a083754229, 9d001a318037a243; ~/.roundfix/roundfix.db run_events

Body (verbatim):
Relato do Fiscus (gesttione-solutions/fiscus), a pedido de Marcio em 05/11/2026. Ele quer que o Roundfix avalie o problema.

**O que aconteceu.** A fila do `roundfix deliver`, com cerca de 21 Specs, começou em 03/11 e já roda há 2 dias. 7 Specs foram mergeadas, 7 estão em andamento e 7 esperam pré-requisitos. Quase todos os itens estacionaram de 2 a 3 vezes: QA reprovado (`run-unresolved`), `review-findings`, `review-blocked` por tamanho do diff e colisão do número de migration drizzle entre Specs paralelas. Uso registrado: 927 M tokens em 31 Runs, 136,83 USD em 18 das 85 sessões medidas.

**Pontos para o Roundfix avaliar:**

1. **Concorrência fixa em 1, sem opção.** Num momento, 5 itens estavam em `reviewing` sem nenhum processo de review rodando, enquanto um único `implement` (0068) ocupava o dono da fila. O `deliver status` mostra `reviewing`, mas não diz que o item espera a vez. Sugestões:
   - mostrar "aguardando capacidade" (ou a posição na vez) no status;
   - permitir concorrência maior para o review pré-PR, que só lê, separada do implement.

2. **Juntar o main no item reinicia o review.** Na recuperação de `review-findings`, a gente juntava o main antes do retry. O merge move a merge base, o pre-PR review abre uma linhagem nova na rodada 1, faz o review completo de novo e traz achados novos. Isso aconteceu com a 0066 e a 0074 no mesmo dia. Sem o merge, a rodada 2 leria só o delta. Sugestões:
   - avisar no `deliver status` ou no retry que o head novo mudou a merge base e vai abrir review completo;
   - ou permitir que um merge do main sem conflito, cujo único pai novo é o main, conte como delta.

3. **Migrations drizzle de Specs paralelas colidem sempre.** Cada Spec gera a partir do main que conhece. Na fila, cinco Specs (0068, 0070, 0074, 0079 e outra) começaram todas na 0052. Cada merge obriga as demais a renumerar à mão: SQL renomeado, snapshot, prevId, `when` no `_journal.json`, `drizzle-kit generate` e `check`, `schema-reference`. O `derived_paths` com regenerate que falha de propósito só faz o item estacionar. Sugestão: um detector na revalidação que acuse colisão de numeração com itens à frente na fila antes do Run, ou um passo de renumeração apoiado pela fila.

4. **Review bloqueado por um commit sozinho acima do limite**, mesmo com o diff total da Spec dentro dele. A segmentação first-parent pega um commit antigo do item que nasceu de um main velho, e esse trecho mostra o trabalho do main como revertido (~1,5–2 MB). O retry só aceita head descendente do candidato, então não dá para refazer o histórico. Na 0070, a saída foi reaplicar a mudança em commits lineares sobre o main e fechar com um merge de árvore idêntica. Funciona, mas é contorno. Sugestão: medir o trecho contra a base efetiva (merge-base do commit com o main atual) ou ignorar trechos que só revertem o main.

5. **Catálogos grandes de dados autorais** (`icms-es-v2.json`, ~550 KB) passavam sozinhos do limite do review. Resolvemos com `derived_paths` (PR #263 do Fiscus), mas o revisor deixa de ler as regras fiscais. Talvez valha uma categoria "dado autoral grande": revisado à parte ou resumido, em vez de omitido.

**Situação agora.** Marcio autorizou `--qa-override` e dispensa do review só para a 0070 e a 0079, entregues fora da fila por PR, CI e squash merge, com `deliver retry` registrando "Merged outside the queue". O resto continua na fila.

Não precisa responder agora. Se algum ponto virar Spec do Roundfix, avise o Fiscus.
