# Maintainer answers for Spec 0392 (AskUserQuestion via `/roundfix --follow-up`, 2026-11-14 ~06:45–06:58 BRT)

Quote these verbatim in `_authorization.md` and in the PRD `## Decisions`.

## Scope (Indexus message 01a1d68d)

Question: "O Indexus (mensagem 01a1d68d, pedido seu de 13/11) pede que `roundfix skills install --target project` e o Baseline escrevam `model` e `effort` no cabeçalho das skills e gerem um subagente de pesquisa em `.claude/agents/`, a partir de uma seção `skills` do `.roundfixrc.yml`. É feature (PRD + TechSpec + Tasks, release minor v0.144.0). O que faço com o pedido?"

Answer: "Spec completa" — option text: "Uma Spec com a seção `skills` no Project Config, o cabeçalho model/effort na instalação, o subagente de pesquisa só leitura e as skills de autoria passando a usá-lo. Cabe em 4 Tasks + QA."

The maintainer added this free text to the answer (it is the source of Spec 0391, not of this Spec): "E também uma observação sobre a regra de 4 tasks que talvez tenha que virar uma nova redação: O objetivo da spec com limite de 4 tasks era porque tinhamos specs gigantes no inicio com varios contextos na mesma spec e tambem specs com tasks pequenas que consumiam runs com grandes tempos de verificação (rotina normal) e faziam a spec demorar muito mais do que o suficiente."

## User-level copy notice (Pantheon message 01a1d6a1, item 2)

First question: "O Pantheon (mensagem 01a1d6a1) relata a sua decisão de 13/11: skills do Roundfix só instaladas no projeto, nunca no nível do usuário (já removidas no mini e no MacBook; conferi no mini). Ele pede que o aviso de cópia de usuário antiga, em `skills check`, `doctor` e `update check`, pare de indicar só `refresh: roundfix skills install --target …`. O que faço com o texto do aviso?"

Answer: "Spec pequena (Recommended)" — option text: "Correção própria: o aviso passa a oferecer apagar a cópia de usuário, com o comando, ao lado do refresh. 1 Task + QA, release patch. A memória antiga é aposentada de qualquer forma."

Second question, asked after the maintainer's free text above: "Você respondeu "Spec pequena" para o ajuste do aviso de cópia de usuário (Pantheon). Pelo motivo que você lembrou, ela seria justamente uma Spec de 1 Task pagando um ciclo inteiro de Run, QA, review e release. O ajuste é do mesmo contexto da Spec do Indexus (instalação de skills e cópias de usuário), e ninguém está bloqueado por ele. Junto ou mantenho separado?"

Answer (this one stands): "Juntar à Spec do Indexus (Recommended)" — option text: "O ajuste do aviso entra numa Task da Spec do Indexus que já mexe em `skills check`. Um ciclo de entrega e uma release (minor v0.144.0); a Spec adota as duas mensagens e continua em 4 Tasks + QA."

## Authorization

Question: "As Specs que saírem das duas respostas acima precisam de autorização de push, merge e release. A entrega local só entra na fila com o Quaestor ocioso (regra do mini, 1 vaga de verificação). Qual autorização vale para elas?"

Answer: "Mesmas do BE (Recommended)" — option text: "Push, PR e merge; release após cada Spec (patch para correção, minor para feature); até 3 correções por Spec; Makefile, CI, go.mod e lint só quando uma correção exigir, descrito no relatório."

## The Idea stage (2026-11-14 ~07:20–08:05 BRT)

The Idea-stage judgment that sent this Spec through `write-idea` (a pilot run by the operator with `jev-1.14.0`, five yes/no questions over the Indexus message): solution shape still open 0.39, unverified assumption 0.36, market question 0.22, new capability 0.92, behavior already specified 0.88. Two answers fell between 0.3 and 0.7, so the maintainer decided.

Question: "Feature 0392 (modelo por skill e subagente de pesquisa, pedido do Indexus): o Jev ficou indeciso (forma em aberto 0,39; suposição sem evidência 0,36; comportamento já especificado 0,88). Lendo o código achei uma forma alternativa bem menor: preservar o cabeçalho que o projeto escreveu na skill, em vez de gerá-lo a partir do `.roundfixrc.yml` (a dor relatada é que a próxima instalação apaga o cabeçalho). Também há o risco de padrões de modelo embutidos envelhecerem. Ela passa pela etapa de Ideia?" → "Etapa de Ideia agora (Recommended)".

Idea question 1: "Etapa de Ideia da 0392, pergunta 1 (o problema). O pedido do Indexus junta três dores diferentes: trocar de modelo à mão, o cabeçalho posto à mão que a instalação apaga, e as varreduras de pesquisa pesando na sessão principal. Qual delas mais pesa hoje para você?" → free answer "1 > 3 > 2" (forgetting to switch the model, then the sweeps on the main session, then the erased header).

Idea question 2: "Etapa de Ideia da 0392, pergunta 2 (quem usa). A pesquisa confirmou que só o Claude Code troca modelo e esforço pelo cabeçalho da skill; Codex e OpenCode ignoram essas chaves e só aceitam modelo em arquivos de agente próprios. Hoje você contorna com a regra "Model by phase" do seu CLAUDE.md e com o `/model` à mão. Para quem a primeira versão precisa funcionar?" → "Claude Code, qualquer adotante (Recommended)" — option text: "Quem autora Specs pelo Claude Code em qualquer repositório que adota o Roundfix. Nada supõe a sua máquina: sem Fable, Exa ou ai-memory embutidos; Codex e OpenCode seguem como hoje."

Idea question 3: "Etapa de Ideia da 0392, pergunta 3 (sucesso). Você ordenou as dores: esquecer de trocar o modelo, depois as varreduras na sessão principal, depois o cabeçalho apagado. Daqui a um mês, qual resultado faria você dizer que a feature valeu?" → "Ciclo sem troca manual (Recommended)" — option text: "Num ciclo inteiro de autoria (Ideia, PRD, TechSpec, Tasks) cada fase roda no modelo e esforço declarados sem nenhum `/model` ou `/effort` digitado, inclusive quando é o agente que chama a skill. Medido pelo modelo que cada turno registra."

Idea question 4: "Etapa de Ideia da 0392, pergunta 4 (a forma). Medi que o Claude Code 2.1.325 ignora o `model` do cabeçalho da skill quando é o agente que a invoca (só aplica quando você digita `/skill`), e o seu critério de sucesso inclui esse caminho. Qual forma a ideia toma?" → "Mais ambiciosa" — option text: "A híbrida, mais um subagente autor gerado com o modelo forte e a mesma declaração compilada para arquivos de agente do OpenCode e do Codex. Duas Specs." The hybrid it builds on: "Seção `skills` no Project Config, cabeçalho model/effort na instalação, subagente de pesquisa só leitura gerado, e uma regra: quando o agente conduz, cada fase de autoria é delegada a um subagente com o modelo e o esforço declarados. Uma Spec."

Draft approval: "Rascunho da Ideia 0392 (cada fase de autoria roda no modelo declarado): duas Specs, a 0392 para o Claude Code (declaração no Project Config, subagente autor e subagente de pesquisa gerados, cabeçalho na instalação, aviso de cópia de usuário) e a 0394 para OpenCode e Codex. […] Aprova?" → "Aprovado (Recommended)".
