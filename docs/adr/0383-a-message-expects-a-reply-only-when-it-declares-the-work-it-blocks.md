---
status: accepted
created_at: 2026-11-10T19:00:00Z
updated_at: 2026-11-11T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A message expects a reply only when it declares the work it blocks

A Cross-project Message is one-way unless the sending project's own work
cannot start or continue without the answer; then the message says so in its
subject and first lines, the receiving project handles it next without
interrupting running work, and the Follow-up Command ranks it first from the
sender's declaration instead of from the judge.

## Context

[ADR-0330](0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md)
made ai-memory the channel for messages between projects: "a project reads
its messages oldest first and consumes one only after recording its content in
a Spec, a Backlog Entry or a Finding, or naming it when it only informs". The
Baseline's receiving clause already answered only "when the message asks for a
reply", but no clause told the sender that a message expects no reply by
default, that it must not wait for one or watch whether the other project
acted, or how to say that it does depend on the answer. On 2026-11-10 the
maintainer observed that projects sending a message through ai-memory kept
waiting for a reply, which made "um fluxo de vai e volta e vigiar se foi feito
ou não. Muito grande, sem necessidade", and asked for a rule: a reply is needed
only "se o projeto que enviou depender daquela resposta para executar ou
continuar executando algo", the message says so, and "se bloquear uma
implementação, uma spec, ou o trabalho, deve ser tratado como urgência".

[ADR-0319](0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md)
computes the Follow-up Priority from four judge answers and decides that
"Nothing is `unblocking` without an answer", so a message that blocks its
sender was ranked like any other whenever the judge was unavailable: without
a key, inside a Run, at the monthly ceiling or on a secret-pattern match.

The maintainer answered on 2026-11-10 "Texto + ferramenta (Recommended)" (the
three message clauses and the skill change, and the Follow-up Command puts on
top the message that declares blocked work without depending on the judge)
and "Próximo da fila (Recommended)": "Passa à frente de tudo que está
pendente e é atendida logo depois do que já está rodando. Nenhuma Run ou
entrega em curso é interrompida."

## Decision

**A message is one-way by default.** The Baseline's three message clauses of
the cross-project continuity rule keep their identities and change their
words. The sender's clause says a message expects no reply unless the sending
project's own work depends on the answer, and that after sending one that
expects none the sender does not wait for a reply, check whether the other
project acted, send it again or record a `waiting_on` wait for it. The
content clause says that a message whose sender depends on the answer says
so in its subject and first lines, in its own words and language: what is
asked and which implementation, Spec or other work is blocked. The receiving
clause triages such a message first and the rest oldest first, handles it next
(ahead of every other pending item, after the work already running, which is
never interrupted for it), answers it as soon as the answer exists, and sends
no acknowledgment, adoption notice or completion notice to a message that
declared no blocked work.

**Amended by [ADR-0390](0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md) (2026-11-11).** The receiving clause triages first only the oldest message of each sender that declares blocked work, then the rest oldest first, a later blocking message from the same sender included, and says that the precedence never replaces the maintainer's judgment. It answers each direct question a message asks as soon as the answer exists, whether or not the message declared blocked work; acknowledgments, adoption notices and completion notices stay forbidden for a message that declared none. The sender's and the content clauses still hold.

**The declaration is prose, never a marker.** Adopters write messages in
their own language, so no code reads a message body for a fixed word or
phrase. The receiving session reads the message and decides.

**The Follow-up Command reads the declaration from a field the caller fills.**
An entry of the `roundfix/follow-up-messages/v1` file may carry
`blocked_work`, one line naming the work the sender says is blocked, which
the session that writes the file fills from the message it read. The entry
needs `from`. It ranks in `unblocking` with basis `declared`, before every
other item, oldest first, `blocks` names its sender and `blocked_work` its
work, and it is never sent to the judge, because no answer could move it
from that place. The file keeps its schema name: an older `roundfix`
ignores the unknown key and ranks the entry as before. `--status` reads no
message and does not change.

**Amended by [ADR-0390](0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md) (2026-11-11).** Only the oldest entry with `blocked_work` of each `from`, then the lowest reference, ranks in `unblocking` with basis `declared` without the judge. Every later declared entry of the same sender goes to the judge, is ordered as if it declared nothing, keeps its `blocked_work` and names the holding entry in `precedence_held_by`. The declaration still comes only from the field the caller fills, and the file keeps its schema name.

This amends ADR-0319, whose rule that nothing is `unblocking` without an
answer now has one exception, and ADR-0330, whose receiving project no longer
reads every message oldest first. It keeps
[ADR-0338](0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md)
unchanged: every handed-in message is still recorded and consumed before the
turn ends, and a `handoff:` entry never declares blocked work. It mirrors
[ADR-0353](0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md),
whose Waiting Item is placed from its own declaration and never judged; a
message is never a Waiting Item, and a message the sender does not need
answered is never watched with a Waiting Declaration.

## Considered Options

- **Keep the rule in the Baseline text alone.** Rejected by the maintainer's
  answer: without the command, a blocking message reaches `unblocking` only
  when the judge answers, and the judge is unavailable in Runs and without a
  key.
- **Detect a marker phrase in the message body.** Rejected: adopters write in
  Portuguese and other languages, and a check that waits for one fixed phrase
  misses every message written another way.
- **Ask the judge a fifth question, whether the sender says it is blocked.**
  Rejected: the order would still need the judge, and the session that reads
  the message already knows the answer.
- **Treat a blocking message as an interruption.** Rejected by the maintainer's
  answer "Próximo da fila": it waits for the work already running.
- **A new schema version for the messages file.** Rejected: the field is
  optional, an older binary ignores it, and every existing file keeps its
  report.

## Consequences

- Adopters receive the reworded clauses at their next `roundfix baseline
  update`; no clause is added or removed, so no Source Baseline row changes.
- Senders stop waiting for replies that were never needed, and receivers stop
  sending notices nobody reads.
- A sender that declares blocked work without being blocked jumps the queue;
  the clause allows the declaration only when the work is blocked, and the
  receiving session reads the message before it acts on it.
- The Roundfix skill's `modes` and `follow-up` references, the command guide
  and the glossary change, and the owned-skill version rises.
