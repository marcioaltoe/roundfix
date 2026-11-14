# Maintainer request and answers for Spec 0393 (`/roundfix --follow-up`, 2026-11-14 ~07:00–08:40 BRT)

Quote these verbatim in `_authorization.md` and in the PRD `## Decisions`. The source of this Spec is the maintainer himself, in this session; no ai-memory message is adopted.

## The request

The maintainer wrote, in the session, while the follow-up was planning Specs 0391 and 0392:

"Outro ponto é que gostaria de um julgamento (usando o JEV) sobre a necessidade do uso do write-idea para features/specs mesmo sem o acionamento do usuário. Essa skill levanta questões que podem ser muito valiosas e evitar trabalho inutil antes de implementações de novas features ou de roadmaps."

## Where the judgment lives

Question: "Você pediu um julgamento do Jev sobre a necessidade do write-idea mesmo sem o usuário acionar. Hoje a regra de roteamento só manda rodar write-idea em "iniciativa grande ou difusa", quem decide é o agente, e o `roundfix spec judge` não faz essa pergunta (só tem os estágios prd, techspec e tasks). Onde esse julgamento deve morar?"

Answer: "Spec nova no produto (Recommended)" — option text: "Feature: antes do PRD de uma feature, o Roundfix pergunta ao Jev se ela precisa da etapa de Ideia (critérios da própria write-idea: forma da solução em aberto, suposição de negócio não conferida, questão de mercado) e a regra de roteamento passa a mandar rodar write-idea quando o juiz indicar. Vale para todos os adotantes."

## Applying it to Spec 0392 before the product has it

Question: "A feature do Indexus (modelo por skill e subagente de pesquisa, Spec 0392) está para ser autorada direto pelo PRD. Aplico já o julgamento a ela, com uma chamada avulsa ao Jev que envia o texto do pedido do Indexus?"

Answer: "Julgar e seguir o resultado (Recommended)" — option text: "Rodo o julgamento agora. Se indicar a etapa de Ideia, faço write-idea antes do PRD e trago as perguntas dela para você; se não indicar, sigo direto ao PRD."

## Authorization

Question: "Se o julgamento do write-idea virar Spec no produto (primeira pergunta), qual autorização vale para ela? É feature, então a release dela é minor."

Answer: "Mesmas do BE (Recommended)" — option text: "Push, PR e merge; release minor depois dela; até 3 correções; Makefile, CI, go.mod e lint só quando uma correção exigir, descrito no relatório."

## The Idea stage (2026-11-14 ~08:10–08:40 BRT)

The Idea-stage judgment that sent this Spec through `write-idea` (the operator's pilot, `jev-1.14.0`, five yes/no questions over the request above): solution shape still open 0.23, unverified assumption 0.61, market question 0.46, new capability 0.81, behavior already specified 0.63. Two answers fell between 0.3 and 0.7, so the maintainer decided.

Question: "Feature 0393 (o Roundfix pergunta ao Jev se uma feature precisa da etapa de Ideia): o Jev também ficou indeciso sobre ela mesma (suposição sem evidência 0,61; questão de mercado 0,46; forma em aberto 0,23). Pontos abertos que achei: onde o julgamento roda quando ainda não existe Spec, o que acontece em adotante sem chave do Jev, e como uma etapa interativa se comporta em trabalho sem supervisão. Ela passa pela etapa de Ideia?" → "Etapa de Ideia agora (Recommended)".

Idea question 1: "Etapa de Ideia da 0393 (o Jev julga se uma feature precisa da etapa de Ideia), pergunta 1 (o problema). Você disse que a write-idea levanta questões valiosas e evita trabalho inútil. Qual desperdício você mais quer evitar com o julgamento automático?" → "Feature na forma errada (Recommended)" — option text: "A Spec vai direto ao PRD com a forma que o pedido sugeriu, e uma alternativa menor ou melhor só aparece depois de Tasks cortadas ou código escrito. Na 0392 de hoje, a etapa de Ideia mudou a forma antes do PRD."

Idea question 2: "Etapa de Ideia da 0393, pergunta 2 (a autoridade do julgamento). A pesquisa mostrou dois riscos opostos: etapa obrigatória vira cerimônia, e recomendação automática tende a ser aceita sem ser questionada. Quando o Jev indicar que uma feature precisa da etapa de Ideia, o que acontece?" → "Roda, com saída registrada (Recommended)" — option text: "O agente inicia a write-idea sem pedir licença e mostra os números do julgamento. Você pode dispensar a etapa na primeira pergunta dela, e a dispensa fica escrita no PRD com o motivo."

Idea question 3: "Etapa de Ideia da 0393, pergunta 3 (quando o juiz não responde). O Jev só responde com chave configurada, fora de um Run e abaixo do teto mensal; hoje o filtro de idioma do juiz também recusa texto que não esteja em inglês, e os pedidos da frota chegam em português. Um adotante sem chave nunca terá o julgamento. O que a regra manda o agente fazer quando o juiz não responde?" → "O agente aplica os critérios (Recommended)" — option text: "O agente responde ele mesmo às mesmas perguntas (forma em aberto, suposição sem evidência, questão de mercado), diz que o juiz não estava disponível e segue o resultado; na dúvida, pergunta a você. Todo adotante ganha a regra, com ou sem chave."

Idea question 4: "Etapa de Ideia da 0393, pergunta 4 (sucesso). O desperdício que você quer evitar é a feature construída na forma errada. Daqui a dois meses de uso, qual medida faria você dizer que o julgamento valeu?" → "Menos correções depois" — option text: "Specs de feature passam a precisar de menos correções de plano, QA e review do que antes do julgamento. Medido nos registros de arquivo das Specs."

Idea question 5: "Etapa de Ideia da 0393, pergunta 5 (a forma). Você já decidiu: o julgamento faz a etapa de Ideia rodar com dispensa registrada, o agente aplica os critérios quando o juiz não responde, e o sucesso é ter menos correções depois. Qual forma a ideia toma?" → "Mais ambiciosa" — option text: "A original com registro, mais um relatório que remede os limiares com o histórico do próprio repositório e o conselho por item já no relatório do `roundfix follow-up`. Duas Specs." The original it builds on: "Juiz antes do PRD de toda feature, critérios aplicados pelo agente sem juiz, e cada Spec de feature registra o julgamento, se a Ideia rodou e o que mudou. O registro permite medir a sua meta e remedir os limiares depois. Uma Spec."

Draft approval: "Rascunho da Ideia 0393 (uma feature é julgada para a etapa de Ideia antes do PRD): duas Specs, a 0393 com o julgamento, a regra, os critérios sem o juiz, o registro e o português, e a 0395 com o relatório de calibração e o julgamento no relatório do follow-up. […] Aprova?" → "Aprovado (Recommended)".
