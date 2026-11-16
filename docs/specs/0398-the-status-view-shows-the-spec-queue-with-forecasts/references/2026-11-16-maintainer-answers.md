Spec da tabela da fila com previsão no /roundfix --status (backlog docs/backlog/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md, PR #879).
Pedido original do Marcio (15/11/2026 08:52 BRT): "É possível que no comando /roundfix --status seja informados além do que já temos uma tabela com a fila de specs programada, da ordem de implementacao mais antiga para a mais nova (incluindo possíveis mudanças, ou seja, a fila corrente), com o status e breve observação e para as que ainda não iniciaram, previsão de inicio e fim?"
Roteamento (15/11): "Backlog agora, Spec depois (Recommended)"; quando autorar (15/11 ~20:15 BRT): "Logo após o plano da 0397 (Recommended)".
Julgamento da Ideia (16/11, roundfix spec judge --stage idea --request, typesafe/jev-1.13-20260917): undecided — shape_open unsure (0.37), unverified_assumption unsure (0.39), market_question no (0.07).
Respostas (AskUserQuestion, 16/11/2026 ~08:00 BRT), verbatim:
- Etapa de Ideia → "Rodar a Ideia" — "write-idea em Fable xhigh com perguntas antes do PRD, como na 0397."
- Autorização → "Mesmas do BE (Recommended)" — "Push, PR, merge, release minor depois da Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir; com o Quaestor ocioso, entregue depois da 0397."

Respostas às perguntas da Ideia (AskUserQuestion, 16/11/2026 ~09:15 BRT), opção escolhida verbatim:
1. Tamanho da V1 → "Tabela com previsão (Recommended)" — "O pedido original, neste repositório: status, observação e previsão pelo histórico de Runs do próprio repositório."
2. Linhas além da fila → "Autoria e backlog (Recommended)" — "Entram Specs em Ideia, PRD ou TechSpec (pelo registro de fase da 0397) e backlog aberto, marcados 'ainda não é Spec', com a espera declarada e sem previsão de fim."
3. Fim previsto → "Merge no main (Recommended)" — "Fim é a entrega integrada (Run, revisão, checks e merge); a release segue na linha 'próxima release' do status."
4. Formato → "Provável e limite (Recommended)" — "Dois horários por célula, ex.: '~10:30, até 12:10' (mediana e percentil 85), com a base no rodapé."
5. Esperas → "Condição e duração (Recommended)" — "Início 'depois de <condição>' e fim como duração após o início; as de trás encadeiam igual, sem horário inventado."
6. Validação → "Backtest antes (Recommended)" — "No PRD, o método é repetido sobre ~100 Specs dos últimos 30 dias e só é aceito se o fim real cair dentro do limite em pelo menos 80%."
- Spec em autoria mostra previsão? (AskUserQuestion, 16/11/2026 ~09:50 BRT) → "Mostra previsão (Recommended)" — "Uma Spec em autoria recebe o provável e o limite pelo histórico de plano-até-merge; só o backlog fica sem fim."
- Aprovação da Ideia 0398 → "Aprovar e seguir (Recommended)".
