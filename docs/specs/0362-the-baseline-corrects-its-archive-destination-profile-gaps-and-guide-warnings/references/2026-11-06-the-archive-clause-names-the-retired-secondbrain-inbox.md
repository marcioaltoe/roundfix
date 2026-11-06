---
type: fix
status: promoted
created: 2026-11-06
spec: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
reason: null
---

# ai-memory message 01a1b085-3c22-78c2-9205-5ba611e0b140 — from onioncry (marcioaltoe/onioncry), 08:01 BRT
Subject: Baseline docs-layout ainda manda copiar para o "Secondbrain inbox" aposentado
Evidence: marcioaltoe/onioncry@9c251d6:docs/agents/docs-layout.md:194 (setup-context-driven managed region, after roundfix 0.136.0 baseline update)
Lacuna de Baseline observada no OnionCry (marcioaltoe/onioncry) em 06/11/2026, depois de aplicar `roundfix baseline update` com o roundfix 0.136.0 (commit 9c251d6).
O guia gerenciado `docs/agents/docs-layout.md` (região setup-context-driven, linha 194) ainda diz, na regra de arquivamento: "so copy every file a later reader needs to `docs/references/`, an accepted ADR, the glossary, an agent guide or the Secondbrain inbox before or at archive".
Em 25/10/2026, o Marcio aposentou o inbox físico do Secondbrain (`~/dev/secondbrain/inbox/`) para enviar e receber. As trocas entre projetos passaram a usar só as mensagens do ai-memory. O próprio guia gerenciado `docs/agents/secondbrain.md` já proíbe escrever em `inbox/`, e o `follow-up.md` diz que a fonte do inbox do Secondbrain foi retirada. Os dois textos se contradizem.
Sugestão: trocar "or the Secondbrain inbox" por um destino vigente, por exemplo uma mensagem ai-memory ao `default/marcioaltoe-secondbrain` para ingestão, e reemitir o Baseline. O OnionCry não editou a região gerenciada.

