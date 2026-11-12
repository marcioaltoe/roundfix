# Source for Spec 0386 — requested by the maintainer

The maintainer's instruction in chat on 2026-11-12, about 18:40 BRT
(21:40 UTC), in reply to the operator's cycle BD report. Both texts are
copied verbatim.

## The operator's question (verbatim)

"O Pantheon e o Drive pediram aviso quando a correção sair, e o Drive diz que
o pedido do retorno veio de você. A regra da 0380 proíbe aviso de conclusão a
mensagem que não declarou trabalho bloqueado, e nenhuma das duas declarou.
Mando os dois avisos quando as Specs saírem, ou mantenho a regra?"

## The maintainer's instruction (verbatim)

"Adicione exceção a regra para quando houver uma solicitação de aviso."

## The two notice requests the question names

Both are ai-memory messages that declare no blocked work and that Specs
planned on the same day adopted. The lines that ask for the notice are quoted
verbatim; the rest of each message is about its own defect.

- ai-memory message `01a1d081-347f-70c1-976e-f00851ad9a19` from pantheon
  (default/pantheon), 2026-11-12 16:04 UTC (13:04 BRT), adopted by Spec 0383,
  third request: "Avisar o Pantheon por mensagem quando sair a versão que
  aceita a versão 2. Na rotina seguinte, o Pantheon deixa o lock na versão 2 e
  remove o contorno."
- ai-memory message `01a1d0ad-dd6f-7bb0-ad09-8a1e38911657` from
  gesttione-drive (default workspace), 2026-11-12 16:53 UTC (13:53 BRT),
  adopted by Spec 0384, section "AO TERMINAR": "Mandem uma linha para
  default/gesttione-drive com o que foi feito, a versão, e quanto cache Go
  ficou no disco interno. Não precisa da aprovação de Marcio para o retorno."

## Operator note (read on Roundfix v0.142.0, `73421db2`, 2026-11-12)

The receiving clause of `rule.core.cross-project-continuity`
(`clause.core.consume-messages-after-recording`, rule version 3) ends "Answer
each direct question a message asks with `memory_message_send` as soon as the
answer exists, whether or not the message declared blocked work, and send no
acknowledgment, adoption notice or completion notice to a message that
declared none. An empty inbox is rest: continue the session's current work."
Step 4 of the Roundfix skill's `--follow-up` mode says the same. Both requests
ask for the version, so the notice falls due when the release that carries
the work is published, not at merge.
