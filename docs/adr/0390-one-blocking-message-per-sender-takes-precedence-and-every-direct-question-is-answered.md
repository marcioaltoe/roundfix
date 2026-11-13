---
status: accepted
created_at: 2026-11-11T18:00:00Z
updated_at: 2026-11-12T21:40:00Z
deprecated_at: null
superseded_by: null
---

# One blocking message per sender takes precedence, and every direct question is answered

Only the oldest message of each sending project that declares blocked work
goes ahead of the receiving project's pending work, every later one from the
same sender waits its normal turn, and the precedence never replaces the
maintainer's judgment. Every direct question a message asks receives an
answer when the answer exists, whether or not the message declared blocked
work.

## Context

[ADR-0383](0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md)
made a Cross-project Message one-way unless its sender's work is blocked,
put every message that declares blocked work ahead of all pending work, and
told the receiving project to answer "as soon as the answer exists when its
message declared blocked work" and to send nothing to a message that
declared none. Roundfix 0.138.0 carried it to the adopters through `roundfix
baseline update`, and two adopters' pre-PR reviews raised a finding on the
rendered rule on 2026-11-11.

Conexus reported that the rule no longer says what to do with a message that
declares no blocked work but asks a real question: read literally, the
question goes unanswered, and an answer is neither an acknowledgment nor an
adoption notice, so the prohibition does not cover it either.

Quaestor reported that the declaration comes from the message body, which the
same rule treats as untrusted data, so any sender, or a tampered body, can put
every message it sends ahead of all of a receiving project's pending work
without limit. It suggested limiting the precedence to known senders, to one
message per sender per session, or to declarations the receiver can verify,
and saying that the precedence does not replace the maintainer's judgment.

The maintainer answered on 2026-11-11 "Uma por remetente + juízo
(Recommended)": "Só a mensagem bloqueante mais antiga de cada remetente passa
à frente; as demais seguem a ordem normal. O texto diz que a precedência não
dispensa o juízo do maintainer. Pergunta direta não bloqueante recebe resposta
quando existir."

## Decision

**Every direct question is answered.** The Baseline's receiving clause
answers each direct question a message asks as soon as the answer exists,
whether or not the message declared blocked work, and still sends no
acknowledgment, adoption notice or completion notice to a message that
declared none. The sender's clause adds that the receiving project owns the
answer to a direct question too, so a sender that does not depend on it still
neither waits for it nor watches for it. Both clauses keep their identities.

**Amended by [ADR-0396](0396-a-message-that-asks-to-be-told-when-its-work-is-done-receives-one-completion-notice.md) (2026-11-12).** A message that asks to be told when its work is done receives one completion notice once that work is delivered, when the release that carries it is published or, when it releases nothing, when it merges, and still no acknowledgment or adoption notice; the request is kept in the record of the message's content, and the sender's clause lets a sender ask for the notice without waiting or watching for it. Every direct question is still answered.

**One blocking message per sender takes precedence.** The receiving clause
triages first, from each sending project, the oldest message that declares
its work blocked, then the rest oldest first, a later blocking message from
the same sender included. The clause says that the precedence orders the
triage and never replaces the maintainer's judgment of what the project
handles first. The Roundfix skill's follow-up mode says the same.

**The Follow-up Command applies the bound from the file it reads.** Of the
`roundfix/follow-up-messages/v1` entries that carry `blocked_work` and share a
trimmed `from`, compared byte for byte, the one with the oldest `created_at`,
then the lowest reference in byte order, holds the declared precedence and
ranks as ADR-0383 decided. Every other declared entry of that sender goes to
the judge and is ordered as if it declared nothing, with basis `judge` or
`type`; it keeps `blocked_work` and gains `precedence_held_by`, the reference
of the holding entry, and its text line appends
` — declares blocked work: <blocked_work> — precedence held by <reference>`.
The command keeps nothing
between runs: each run chooses again from the file it reads, so the sender's
next declared message takes the precedence once the session consumes the
holding one. The follow-up mode fills `from` with the sending project that
ai-memory returns with the message, never with a name the body gives, so the
body cannot pick its own sender.

This amends ADR-0383, whose receiving clause put every declared message first
and answered only declared ones, and whose command ranked every declared
entry first; ADR-0319, whose amended rule placed every declared entry in
`unblocking`; and ADR-0330, whose amended receiving project answered only a
blocking message. It keeps
[ADR-0338](0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md)
and
[ADR-0353](0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md)
unchanged.

## Considered Options

- **Precedence only for senders of the same workspace with a history.**
  Rejected: the command reads one file and no history, and most of a
  maintainer's projects share one workspace, so the bound would hold almost
  nobody back.
- **A declaration the receiver can verify, naming a Spec or Run of the
  sending project.** Rejected: the receiving project cannot read the sending
  repository, and a named Spec is still a claim written in the body.
- **One message per sender per session.** Taken in the form the command can
  enforce: the command has no session memory, so the bound is one entry per
  sender in each file it reads.
- **A distinct basis for a later declared entry.** Rejected: `basis` says how
  the class was computed, and that entry's class comes from the judge or its
  type; `precedence_held_by` says why it did not take the precedence.
- **Ranking a later declared entry without a mark.** Rejected: the report
  would show `blocked_work` on an entry ordered like any other with no reason.
- **Comparing senders without regard to case.** Rejected: `from` is the
  project name ai-memory reports, copied as given, and folding would invent an
  equivalence ai-memory does not define.
- **Answering only blocking messages, as before.** Rejected by both reports
  and the maintainer's answer: a real question went unanswered.

## Consequences

- Adopters receive the reworded clauses at their next `roundfix baseline
  update`; no clause is added or removed, so no Source Baseline row changes.
- One sender can no longer put more than one message at a time ahead of a
  receiving project's work; many senders can still each put one, and the
  maintainer's judgment decides between them.
- A message entry the command now orders by the judge costs one judge
  request, as any undeclared entry does.
- An older `roundfix` that knows `blocked_work` still gives every declared
  entry the precedence until it is upgraded; one that does not know the field
  ignores it, as before.
