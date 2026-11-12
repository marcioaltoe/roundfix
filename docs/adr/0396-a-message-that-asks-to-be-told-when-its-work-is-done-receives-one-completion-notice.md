---
status: accepted
created_at: 2026-11-12T21:40:00Z
updated_at: 2026-11-12T21:40:00Z
deprecated_at: null
superseded_by: null
---

# A message that asks to be told when its work is done receives one completion notice

A Cross-project Message that asks to be told when the work it asks for is
done receives one completion notice once that work is delivered, even when it
declared no blocked work; the request is kept in the record of the message's
content until then, and it still earns no acknowledgment or adoption notice.

## Context

[ADR-0383](0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md)
made a Cross-project Message one-way unless its sender's work is blocked, and
[ADR-0390](0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md)
answered every direct question while keeping one prohibition: the receiving
project "still sends no acknowledgment, adoption notice or completion notice
to a message that declared none". Roundfix 0.139.0 carried that receiving
clause to the adopters.

On 2026-11-12 two messages that declare no blocked work asked to be told when
their work ships: Pantheon's report on the version-2 skills lock ("Avisar o
Pantheon por mensagem quando sair a versão que aceita a versão 2", so that its
next routine removes a workaround) and the gesttione-drive report on the Go
build cache ("Mandem uma linha para default/gesttione-drive com o que foi
feito, a versão, e quanto cache Go ficou no disco interno"). The operator asked
the maintainer whether to send both notices or keep the rule, and the
maintainer answered: "Adicione exceção a regra para quando houver uma
solicitação de aviso."

The receiving clause already makes the project record a message's content in a
Spec, a Backlog Entry or a Finding before it consumes the message, but the
notice falls due later, often in another session: after the release that
carries the work, which in this repository the operator cuts after the merge.
No Roundfix command reaches ai-memory, so only a session can send the notice.

## Decision

**A request for a notice earns one completion notice.** The receiving clause
of the Baseline's cross-project continuity rule keeps its identity and says
that a message that asks to be told when its work is done receives one
completion notice and no acknowledgment or adoption notice. Every other
message that declared no blocked work still receives none of the three.

**The notice follows the delivery.** The work recorded from the message is
delivered when the release that carries it is published or, when that work
releases nothing, when it merges. The notice names what changed and the
version or merge that carries it; when the work is closed without being
delivered, the one notice says so instead. The session that follows that
delivery sends it with `memory_message_send`; a Daemon-assigned turn never
sends one, as no Agent inside a Daemon-assigned Run announces a Spec Event.

**The request lives in the record of the message.** The receiving clause keeps
the request in the record of the message's content: the Backlog Entry or
Finding that records it, or the `## Sources` entry of the PRD that adopts it.
The Roundfix skill's follow-up mode writes it there and, after the release,
its release reference reads the Sources of each Spec the release archived, at
the Archive Record's `source_revision`, and sends the notices due. The request
is prose the session reads, in the sender's language; no command carries it.

**The sender may ask, and still neither waits nor watches.** The sender's
clause keeps its identity and adds that a message may ask to be told when the
work it asks for is done, only when the sending project will act on the
notice, and that the sender does not wait or watch for it either.

This amends ADR-0390, whose receiving clause sent no completion notice to a
message that declared no blocked work; ADR-0383, whose receiving clause and
its ADR-0390 notice say the same; and ADR-0330, whose ADR-0390 notice says the
same of the receiving project. It keeps ADR-0390's precedence and its answer to
every direct question, and the follow-up messages file of ADR-0383 and
[ADR-0319](0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md),
unchanged.

## Considered Options

- **Keep the rule and send no notice.** Rejected by the maintainer's answer:
  two senders asked for one, and each acts on it.
- **Send a completion notice to every message whose work was delivered.**
  Rejected: it brings back the notices nobody reads that ADR-0383 removed, and
  a release note to every adopter, which the Update Check already replaces.
- **Send the notice at merge.** Rejected as the rule: both requests ask for the
  version, which exists only once the release is published. A merge that
  releases nothing is the delivery itself, so it is the exception.
- **A field the Follow-up Command or the Archive Record carries.** Rejected: no
  Roundfix command can send the notice, because none reaches ai-memory, and a
  field would ask the session to read a body for a fixed marker in a language
  it does not choose, which ADR-0383 rejects. The PRD's Sources, read at the
  Archive Record's `source_revision`, already outlive the Spec folder.
- **Send the notice from the Run that implements the work.** Rejected: a Run
  ends before the merge and the release, and an Agent inside a Daemon-assigned
  Run sends no message to another project.

## Consequences

- Adopters receive the reworded sender and receiving clauses at their next
  `roundfix baseline update`; no clause is added or removed, so no Source
  Baseline row changes.
- A sender that asks for a notice gets exactly one, when the version it waits
  for ships or the work is closed, and can act on it without polling.
- A request recorded only in the message body and not in its record can be
  missed by a later session; the follow-up mode writes it into the record, and
  the release reference reads it from there.
- A notice that falls due before this rule ships is the operator's to send on
  the maintainer's answer; this decision does not reach back into Specs
  already planned.
