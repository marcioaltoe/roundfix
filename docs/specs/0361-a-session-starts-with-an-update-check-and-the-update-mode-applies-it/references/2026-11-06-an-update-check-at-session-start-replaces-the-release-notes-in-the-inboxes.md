---
type: feature
status: promoted
created: 2026-11-06
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
reason: null
---

# ai-memory message 01a1b0fb-92e1-7262-bc05-97272b9de2cb (from jarvis, default/jarvis, 2026-11-06 10:10 BRT, nature measured)
Subject: Trocar as notas de release nas caixas por verificação de update no Baseline
Evidence: ai-memory message list --workspace default --project <argus|quaestor|fiscus|vortex> em 06/11/2026 08:58 BRT: 11, 8, 6 e 5 mensagens cujo corpo começa com "# Roundfix v"

Pedido de Marcio, 06/11/2026, na sessão do Jarvis.
Problema: cada release do Roundfix manda a nota "Roundfix vX: o que muda neste repositório" à caixa de cada projeto. Em 06/11 havia 30 dessas notas paradas (Argus 11, Quaestor 8, Fiscus 6, Vortex 5), misturadas aos pedidos reais entre projetos.
Proposta:
1. Parar de mandar as notas de release às caixas do ai-memory.
2. Incluir no Baseline uma verificação de update ao iniciar a sessão: compara a versão do Roundfix e do Baseline do repositório com a última publicada e, se houver versão nova, mostra o que muda naquele repositório.
3. Comando manual pela skill: /roundfix --update, que verifica, mostra o que muda e aplica.
4. Opção por repositório entre atualização automática (aplica ao iniciar a sessão) e manual (só avisa e espera o /roundfix --update).
O conteúdo "o que muda neste repositório" continua, só muda o canal. O desenho (hook de SessionStart, chave no profile, qual modo é o padrão) fica com vocês, no fluxo normal de spec.
Notas que já estão nas caixas: quando a verificação existir, o Roundfix cancela da própria saída (memory_message_cancel) as notas de release ainda pendentes nas caixas dos projetos, e cada projeto passa a ver a mudança pela verificação.
Sem retorno formal; procurem o Jarvis (default/jarvis) só se precisar de decisão do Marcio.

---
Operator note (2026-11-06 ~12:00 BRT): the maintainer chose "Spec agora, padrão avisar (Recommended)" for this request; the answer is quoted in the Spec's `_authorization.md` and its PRD `## Decisions`. Cancelling the release notes already pending in the inboxes stays the operator's step after the release that carries this Spec.
