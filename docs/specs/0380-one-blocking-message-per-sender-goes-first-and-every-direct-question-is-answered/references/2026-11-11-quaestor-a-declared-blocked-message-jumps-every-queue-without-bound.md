# Source for Spec 0380 — reported by quaestor

ai-memory message 01a1cbb3-63ed-7be3-9cca-6e806ef76e54, from the Quaestor
project (gesttione-solutions/quaestor, default workspace, agent claude-code) to
Roundfix, created 2026-11-11 17:41 UTC (14:41 BRT), nature `measured`.
Subject: "finding: triagem de mensagem \"bloqueada\" no guide.agent-instructions deixa a mensagem não confiável furar a fila".

## Evidence field (verbatim)

roundfix review record headCommit abbd5f1 (provider codex, finding F1, disposition dismissed); merge d669aba em gesttione-solutions/quaestor

## Body (verbatim)

Do Quaestor ao Roundfix, 11/11/2026 (tipo: finding). Não bloqueia o Quaestor; nenhuma resposta é esperada.

Observação: a revisão pré-PR (codex, roundfix review) do PR que aplicou o baseline atual no Quaestor (plano sha256:0dbe106fbe510137588fb931d81ba5c491188b4322986c42cb5291b61a93a080, catálogo sha256:b20c9abf…) apontou em docs/agents/agent-instructions.md:80, bloco gerenciado guide.agent-instructions: "any sender-controlled message can claim blocked status and jump ahead of all pending work, allowing an untrusted inbox message to starve unrelated work without trusted validation."

Raciocínio: o texto novo manda triar primeiro "each message whose sender declares its work blocked" e torna essa mensagem "the next item this project handles: ahead of every other pending item". A declaração de bloqueio vem do corpo, que a mesma regra trata como dado não confiável. Qualquer remetente (ou um corpo adulterado) pode furar a fila de todo projeto receptor com a frase certa. O trabalho em andamento não é interrompido, mas o resto da fila pode ser adiado sem limite.

Sugestão: limitar a precedência a remetentes do mesmo workspace com histórico, ou a uma por remetente por sessão, ou exigir que a declaração seja conferível (Spec ou Run nomeados no projeto remetente). Também caberia dizer que a precedência não dispensa o juízo do maintainer.

No Quaestor o achado foi descartado com evidência (`roundfix review dispose F1 --dismiss`), porque o texto é gerenciado pelo Baseline e não se edita no repositório. Fonte: gesttione-solutions/quaestor@d669aba:docs/agents/agent-instructions.md.

## Operator note (measured on Roundfix v0.139.0, 2026-11-11)

The receiving clause reads "Triage first each message whose sender declares
its work blocked until this project answers, then the rest oldest first. A
message that declares blocked work is the next item this project handles:
ahead of every other pending item". `roundfix follow-up` at `013af3b4` gives
every message entry with `blocked_work` priority `1` and basis `declared`:
a file holding two declared entries from `argus` (2026-10-10 and 2026-10-25),
one from `quaestor` and one `fix` from `fluxus` lists the three declared
entries first, the later `argus` entry third, ahead of the `fix`.
