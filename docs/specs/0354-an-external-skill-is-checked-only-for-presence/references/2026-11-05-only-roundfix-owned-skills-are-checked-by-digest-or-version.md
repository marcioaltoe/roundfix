---
type: fix
status: promoted
created: 2026-11-05
spec: 0354-an-external-skill-is-checked-only-for-presence
reason: null
---

# ai-memory message 01a1ab62-5139-7dc3-9b92-472690a044f9 (from fluxus, gesttione-solutions/fluxus, 2026-11-05 08:05 BRT, nature measured)

Subject: Pedido do Marcio: validar digest e versão só das skills autoradas pelo Roundfix

Evidence (sender): fluxus em 05/11/2026 com roundfix 0.132.0: `roundfix doctor` → "skills: ok (89 required: 12 Roundfix-owned, 77 external; DR-SKILL-DIFFERS-FROM-SNAPSHOT: ... systematic-debugging, ui-ux-pro-max)"; `roundfix baseline update --no-history` lista as duas em skills.differing com "restore deliberately with roundfix baseline skills restore"; fluxus PR #462

Body (verbatim):
Pedido de mudança no Roundfix, decidido por Marcio em 05/11/2026 numa sessão do Fluxus (gesttione-solutions/fluxus).

## Contexto
No Fluxus, com o Roundfix 0.132.0, o `roundfix doctor` e o `roundfix baseline update` comparam também as skills de terceiros com o Setup Snapshot. Hoje acusam `DR-SKILL-DIFFERS-FROM-SNAPSHOT` para `systematic-debugging` e `ui-ux-pro-max` e sugerem `roundfix baseline skills restore`. As cópias do repositório vêm do `marcioaltoe/skills`, que é a fonte de verdade escolhida por Marcio, e um restore as trocaria pela cópia do snapshot.

## O que Marcio pediu
1. O Roundfix confere versão e digest só das skills que ele mesmo autora (as 12 "Roundfix-owned": archive-spec, brainstorming, business-analyst, evidence-gate, implement-task, qa-gate, roundfix, setup-context-driven, write-idea, write-prd, write-tasks, write-techspec).
2. Para as outras skills de que o perfil precisa, o Roundfix confere só se estão instaladas, sem digest nem versão. Uma cópia diferente do snapshot não gera aviso nem sugestão de restore.
3. As skills são instaladas a partir do repositório `marcioaltoe/skills`. O Roundfix valida que o digest das próprias skills é igual ao do `marcioaltoe/skills`, ou que a versão é igual ou maior.

## Efeito esperado
- O `roundfix doctor` não mostra `DR-SKILL-DIFFERS-FROM-SNAPSHOT` para skills externas; falta de uma skill externa obrigatória continua sendo acusada.
- O `roundfix baseline update` não lista skills externas em `skills.differing` nem as restaura.

Quando houver release com essa mudança, mande ao default/gesttione-solutions-fluxus o aviso de sempre, "o que muda neste repositório".

---
Also reported by the maintainer the same morning from a Vortex session (verbatim, pasted): "MEsma questão das skills em vortex: · O baseline do Roundfix 0.132.0 manteve systematic-debugging e ui-ux-pro-max diferentes do Setup Snapshot (não se sabe qual é mais nova). Restauro essas duas para o snapshot do catálogo? → Manter como estão (Recommended)"
