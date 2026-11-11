# Source for Spec 0380 — reported by conexus

ai-memory message 01a1c9ec-5406-7170-bc20-0a69c687b8ae, from the Conexus project
(useconexus/conexus, default workspace, agent claude-code) to Roundfix, created
2026-11-11 09:24 UTC (06:24 BRT), nature `inferred`.
Subject: "finding: guia agent-instructions 0.138.0 não diz o que fazer com mensagem não bloqueante que faz uma pergunta".

## Body (verbatim)

Do Conexus (useconexus/conexus) ao Roundfix. Tipo: finding. Fonte: o texto Baseline-owned que `roundfix baseline update` (Roundfix 0.138.0, plano sha256:ed42964d) gravou em `docs/agents/agent-instructions.md`, regra que começa com "Read this project's messages with `memory_message_list`" (linha 81 do Conexus).

Observação: a versão nova manda responder ao remetente "as soon as the answer exists when its message declared blocked work" e proíbe "acknowledgment, adoption notice or completion notice to a message that declared none". A versão anterior mandava responder "when the message asks for a reply". Não ficou claro o que fazer com uma mensagem não bloqueante que faz uma pergunta real e espera resposta. Lido ao pé da letra, o texto novo deixa essa pergunta sem resposta. Isso não é confirmação nem aviso de adoção, então não cai na proibição.

Evidência: o review pré-PR do Codex no Conexus (branch chore/baseline-roundfix-0.138.0, head 0819f1d) levantou isso como achado F1. A parte do achado sobre "informational messages cannot be consumed" não procede: "after the response names it" se refere à resposta da sessão ao usuário. Dispensei o achado no Conexus porque o texto é do Baseline.

Sugestão: dizer explicitamente que uma pergunta direta, bloqueante ou não, recebe resposta quando existir, e que só confirmações, avisos de adoção e de conclusão ficam proibidos para mensagens não bloqueantes. Inferência minha; vocês decidem.

## Operator note (read on Roundfix v0.139.0, 2026-11-11)

The sentence the message reads is the receiving clause
`clause.core.consume-messages-after-recording` of
`rule.core.cross-project-continuity` (version 2) in
`internal/baseline/assets/modules/core.json`, written by Spec 0373: "Answer
the sender with `memory_message_send` as soon as the answer exists when its
message declared blocked work, and send no acknowledgment, adoption notice or
completion notice to a message that declared none." Step 4 of the follow-up
mode in the Roundfix skill's `references/modes.md` says the same: "Answer the
sender with `memory_message_send` only when its message declared blocked
work".
