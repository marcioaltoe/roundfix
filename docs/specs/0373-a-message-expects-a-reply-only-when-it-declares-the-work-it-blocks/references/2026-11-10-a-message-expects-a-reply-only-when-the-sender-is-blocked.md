---
type: feat
status: promoted
created: 2026-11-10
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
reason: null
---

# A message to another project expects a reply only when the sender is blocked — requested by the maintainer

Recorded by the operator from the Roundfix Supervisor session of 2026-11-10 ~15:15 BRT.
The maintainer invoked `/roundfix --follow-up` and dictated the request below (voice
dictation; the first words are transcription noise and are kept as received).

## Request (verbatim)

> Ana. Adicione e. Uma ideia que eu estou percebendo ao baseline: uma percepção que eu estou tendo ao baseline. Sempre que os projetos enviam algo para outro projeto usando AI memory, ele fica aguardando uma resposta. Essa resposta só é necessária se ele precisar da resposta para alguma coisa. Nova ação: se ele só está passando para ele algo que deve ser feito pelo outro projeto, isso não demanda uma resposta; isso está fazendo com que a gente tenha um fluxo de vai e volta e vigiar se foi feito ou não. Muito grande, sem necessidade. Então, a proposta seria uma regra para que só seja necessária uma resposta a uma mensagem enviada para outro projeto, se o projeto que enviou depender daquela resposta para executar ou continuar executando algo. E ele deve reforçar na mensagem que ele precisa da resposta para continuar o trabalho, e que se bloquear uma implementação, uma spec, ou o trabalho, deve ser tratado como urgência.

## What the Baseline says today (Roundfix v0.137.7, `internal/baseline/assets/modules/core.json`)

Three clauses of the same rule carry the cross-project message contract, rendered into
every adopter's `docs/agents/agent-instructions.md`:

- `clause.core.message-the-owning-project` — an observed defect, improvement or idea
  obliges one message to the project that owns the fix.
- `clause.core.messages-are-self-contained` — what each message carries: the sending
  project, a type hint, the observation, its evidence and the reasoning.
- `clause.core.consume-messages-after-recording` — read with `memory_message_list`,
  triage oldest first, consume only after recording, and "answer the sender with
  `memory_message_send` when the message asks for a reply".

So the receiving side already replies only when asked. Nothing tells the sender that
a message expects no reply by default, that it must not wait for one or watch whether
the other project acted, or how to say that it does depend on the answer. Triage is
"oldest first", with no urgency for a message that blocks the sender's work.

## What the Follow-up Command does today (`internal/cli/follow_up.go`, `followUpOrder`)

An item reaches the `unblocking` class only when the Jev judge answers
`blocks_other_project: yes` and also `is_fix: yes` or `is_small: yes`. Without a judge
answer (no key, inside a Run, monthly ceiling, secret-pattern match) nothing is
`unblocking` (ADR-0319; glossary **Follow-up Priority**). The messages file
(`roundfix/follow-up-messages/v1`) has no field for a sender's declaration.

## Maintainer answers (AskUserQuestion, 2026-11-10 ~15:35 BRT)

See the Cycle BB section of the briefing: "Texto + ferramenta (Recommended)",
"Próximo da fila (Recommended)" and "Mesmas do BA, teto minor (Recommended)".
