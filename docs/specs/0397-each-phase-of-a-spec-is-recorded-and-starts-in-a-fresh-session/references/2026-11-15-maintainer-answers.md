Respostas do Marcio (AskUserQuestion, 15/11/2026 ~20:15 BRT):
- Ideia 0397 (Jev: not indicated, shape 0.26 / premissa 0.18 / mercado 0.07; judged-by typesafe/jev-1.13-20260917): "Manter a Ideia".
- Modelo do subagente research em skills.authoring.claude do Roundfix: "Sonnet, effort medium (Recommended)".
- Backlog da tabela da fila no --status: "Logo após o plano da 0397 (Recommended)".
Decisões anteriores (14/11, ~11:35 BRT): Jarvis 01a1da5d → "Medidas já, sem Spec" + "Spec com Ideia após o reset" (0397) + "Orquestração mais barata" = "Opus medium"; autocompactação RECUSADA ("passa para determinístico e pode perder informação importante"); adoção do modelo por fase da 0392 no Roundfix → "Adotar na Spec 0397"; autorização "Mesmas do BE" (push, PR, merge, release após cada Spec, até 3 correções, Makefile/CI/go.mod/lint só se a correção exigir) vale p/ 0396 e 0397.
Escopo decidido (handoff 01a1dc8f): o Run grava o id e a fase num lugar fixo para o ai_usage.py; o Supervisor fecha a sessão em cada fronteira de fase.

Respostas às perguntas da Ideia (AskUserQuestion, 15–16/11/2026 ~00:00 BRT), opção escolhida verbatim:
1. Cobertura do registro de fase → "Run + sessões Supervisor (Recommended)" — "Registra sessão, Spec e fase também fora do Run (autoria, entrega, release, triagem) e, dentro do Run, a Task ou o QA em curso."
2. Quem relata o consumo por fase → "Roundfix só expõe (Recommended)" — "Lugar e formato documentados e estáveis; o ai_usage.py cruza com os transcripts; o Roundfix não lê transcript de harness nem calcula preço."
3. Fronteira de sessão do Supervisor → "Instrução e aviso (Recommended)" — "Skills e guia mandam gravar handoff, dizer o que ficou armado e encerrar; o Roundfix avisa, sem bloquear, quando a mesma sessão registra duas fases; você abre a próxima com /roundfix follow-up."
4. Adoção dos modelos por fase no Roundfix → "Só subagentes gerados (Recommended)" — "Na 0397: autores por fase e pesquisador em Sonnet em .claude/agents/, sem cabeçalho nas skills deste repositório; uma skill digitada segue no modelo da sessão."
5. Correção pequena × declaração por fase → "Exceção pela rota (Recommended)" — "Na rota de correção ou refactor a skill escreve na sessão e só delega quando você pedir; feature sempre delega; muda o texto das skills para todos os adotantes."
6. Direção da Ideia → "Original + respostas (Recommended)" — "Registro de fase, fronteira de sessão e adoção dos modelos numa Spec; o custo de nascer de cada sessão vira Backlog Entry."

Acréscimo do Marcio ao critério de sucesso (16/11/2026, na sessão), verbatim: "Adicione ao critério de sucesso: com queda de performance nula ou quase nula e com aumento de qualidade de artefatos para implementações."
