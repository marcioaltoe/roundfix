---
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
status: active
created: 2026-11-10
surfaces: [backend, cli, docs]
type: feat
---

# A message expects a reply only when it declares the work it blocks

Requested by the maintainer on 2026-11-10: projects that send each other
messages through ai-memory keep waiting for replies and checking whether the
other project acted, a back-and-forth that is "Muito grande, sem necessidade".
After this Spec a Cross-project Message is one-way unless the sending
project's own implementation, Spec or other work cannot start or continue
without the answer. Then the message says so in its subject and first lines
and names the blocked work, and the receiving project handles it next, ahead
of every other pending item and after the work already running, which it never
interrupts. Every adopter receives the rule through its next `roundfix
baseline update`, and the Follow-up Command ranks such a message first from a
field the session fills, without the judge.

## Sources

One request, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-10-a-message-expects-a-reply-only-when-the-sender-is-blocked.md`](references/2026-11-10-a-message-expects-a-reply-only-when-the-sender-is-blocked.md),
   `type: feat`, requested by the maintainer in the Roundfix Supervisor
   session of 2026-11-10 ~15:15 BRT (voice dictation, quoted verbatim), with
   the operator's notes on what the Baseline and the Follow-up Command do
   today.

At the authoring base `docs/backlog/` and `docs/findings/` are empty. This
project's ai-memory inbox, listed read-only on 2026-11-10, holds two pending
messages, the maintainer's decision for Spec 0371 relayed by jarvis and a
review-lineage report from conexus for Spec 0372; neither shares this Spec's
context, and none was consumed.

## Prerequisites

None in `requires`. Authored against `9fda0b65` (v0.137.7). Specs 0372 and
0371 deliver before this one, in that order, and may change files this Spec
also changes, each in its own section:

- The Roundfix owned skill's version fields in the `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json` (0372 changes the `review`
  reference, 0371 the QA references; this Spec the `follow-up` and `modes`
  references). Each Task raises the version with the record command from the
  tree it runs on and never writes a number.
- `CONTEXT.md` (each Spec its own terms; this Spec the message and follow-up
  terms) and `docs/references/behavior-surfaces.json` (each Task re-records it
  after its own last guide edit).
- `internal/cli/cli_test.go` is not changed by this Spec.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; the messages file keeps its schema name and gains one optional
  field. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the command still never reaches ai-memory or the
  judge for a declared message. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0383 (new) records the decision:
  "A message expects a reply only when it declares the work it blocks".
  ADR-0319 is amended: it decides "Nothing is `unblocking` without an answer",
  and a message entry that declares blocked work now ranks in `unblocking`
  without one. ADR-0330 is amended: its receiving project "reads its messages
  oldest first", and a message that declares blocked work is now triaged
  first. ADR-0338 applies unchanged: "No message the mode handed in stays
  pending when the turn ends", and a declared message is recorded and
  consumed like any other. ADR-0353 applies unchanged: "A Waiting Item is
  never sent to the judge, with or without --all", the precedent a declared
  message follows. ADR-0372, which cites ADR-0330, does not apply: it names
  where the Baseline's archive rule sends a copy and a profile's missing
  decision, and this Spec changes neither. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  BB answers of 2026-11-10 quoted in the authorization record, which grant
  the Baseline source, its fixtures and goldens, the guides under
  `docs/agents/` and the owned skills to this Spec; the Governed Paths are the
  Roundfix skill files tasks 01 and 03 change, and the core Baseline module,
  the Standard TypeScript profile and its formatter golden, and this
  repository's managed guide and setup manifest that task_02 changes. No
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`.

## Goals

- A project that sends a message it does not need answered stops waiting for
  a reply and stops watching whether the other project acted.
- A project whose work is blocked by another project's answer says so in the
  message, and that message is the receiver's next item, never an
  interruption of running work.
- The receiving project replies only to a message that declared blocked work,
  as soon as the answer exists, and sends no acknowledgment or progress notice
  otherwise.
- The Follow-up Command puts a declared blocking message on top whether or not
  the judge answers.
- Every adopter gets the same rule through the Baseline, in any language its
  messages are written in.

## User Stories

1. As an Agent session in a project that sends a finding to another project,
   I want the rule to say that the message expects no reply, so that I do not
   wait, poll or arm a wake-up for it.
2. As an Agent session whose Spec cannot continue without another project's
   answer, I want to say so in the message and name the blocked work, so that
   the other project answers it first.
3. As an Agent session reading this project's inbox, I want to know which
   messages block their senders, so that I handle those next and answer them,
   and send nothing back to the rest.
4. As the maintainer running `/roundfix --follow-up`, I want a declared
   blocking message at the top of the list even when the judge is
   unavailable, so that another project is never left blocked behind routine
   work.

## Core Features

1. Sender rule: the clause that obliges a message to the owning project adds
   that a message expects no reply unless the sending project's own work
   depends on the answer, and that after sending one that expects none the
   sender does not wait for a reply, check whether the other project acted,
   send it again or record a wait for it.
2. Content rule: the clause on self-contained messages adds that a sender
   whose work cannot start or continue without the answer says so in the
   message's subject and first lines, in its own words and language, naming
   what is asked and which implementation, Spec or other work is blocked, and
   declares blocked work only when it is blocked.
3. Receiver rule: the clause on reading and consuming messages triages first
   each message that declares blocked work and the rest oldest first, makes it
   the next item handled (ahead of every other pending item, after the work
   already running, which is never interrupted for it), answers it as soon as
   the answer exists, and sends no acknowledgment, adoption notice or
   completion notice to a message that declared none. Recording before
   consuming stays as it is.
4. The three clauses keep their identities, force and Source Baseline
   dispositions; no clause is added or removed.
5. The Follow-up Command's messages file accepts an optional `blocked_work`
   per entry. An entry that carries it ranks first in `unblocking` with basis
   `declared`, oldest first, is never sent to the judge, and the text and
   JSON reports name the blocked project and the blocked work. A malformed
   value is refused naming the message. A file without the field gives the
   same report as before, and `--status` is unchanged.
6. The Roundfix skill's follow-up mode fills `blocked_work` from what a
   message says, puts a declared message first in its plan, replies only to a
   declared message, and sends another project's finding expecting no reply,
   without a wake-up or a Waiting Declaration to watch it unless its own work
   depends on the answer.
7. The glossary adds **Blocking Request** and changes **Cross-project
   Message** and **Follow-up Priority**; ADR-0383 records the decision and
   ADR-0319 and ADR-0330 carry its amendment notices.

## User Experience

The text list shows a declared blocking message on the first line, for
example
`1. unblocking information 2026-10-10 message:d1: Export field — from argus — blocks argus: release 1.4 waits on the answer`.
Its JSON item has `basis: "declared"`, `blocks` naming the sender and
`blocked_work` naming the work. A list whose only items are declared or
waiting reports the judge `not_asked` with reason
`only declared or waiting items`.

## Non-Goals / Out of Scope

- Detecting a declaration inside a message body by a word, a phrase or a
  pattern: the session reads the message and fills the field.
- Interrupting a running Run, Delivery Queue item or Task for a declared
  message.
- Adding, removing or renaming a Baseline clause, or a Source Baseline row.
- Changing the vendored `ai-memory-messaging` skill: the Baseline rule
  governs over its defaults, as the skill dispatch guide already says.
- Changing the Spec Event Announcement and release-announcement steps, which
  already never wait for a reply.
- Changing `--status`, the four judge questions, their thresholds, or the
  Waiting Declaration of Backlog Entries and Findings.
- A new schema version for the messages file, and any change to ai-memory or
  to how messages are sent.
- Telling adopters by message: they receive the rule through their next
  Baseline update and the Update Check.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository
   writes into `docs/agents/agent-instructions.md`, once each, the phrases
   "A message expects no reply unless the sending project's own work depends
   on the answer", "say so in the message's subject and first lines" and
   "after the work already running, which is never interrupted for it"; a
   second refresh reports it current; the three clause ids are unchanged and
   their Source Baseline replacements still record `replaced`. At `9fda0b65`
   none of the three phrases is in the guide.
2. Success Metric: with a fake judge, a messages file of two judged messages
   and two messages with `blocked_work` lists the two declared messages first,
   oldest first, with `priority` `1`, `basis` `declared`, `blocks` the sender
   and null `answers`, then the judged ones; the judge receives exactly the two
   judged messages. At `9fda0b65` the field is ignored and the declared
   messages are sent to the judge.
3. Success Metric: with only declared messages the judge is `not_asked` with
   reason `only declared or waiting items` and receives no request; a
   `blocked_work` with a line break, or on an entry without `from`, exits `2`
   with empty standard output and the diagnostic naming the message id.
4. Success Metric: every existing Follow-up Command test passes unchanged, so
   a file without `blocked_work` keeps its report bytes.
5. Success Metric: both copies of the Roundfix `modes` reference say, in step
   1, when to give an entry `blocked_work`; in step 2, that a declared message
   leads the plan; in step 4, that the mode answers only a declared message and
   sends no acknowledgment, adoption notice or completion notice otherwise; in
   step 8, that another project's finding expects no reply and is watched by
   no wake-up or Waiting Declaration.
6. Success Metric: the command guide and both `follow-up` references document
   `blocked_work`, the `declared` basis and the refusals; `CONTEXT.md` defines
   **Blocking Request**; ADR-0319 and ADR-0330 carry an "Amended by ADR-0383"
   notice; and the documentation, skill and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published literature this Spec did not
write, read through Exa on 2026-11-10.

- Hohpe and Woolf, *Enterprise Integration Patterns*, "Request-Reply"
  (https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html):
  "When two applications communicate via Messaging, the communication is
  one-way. The applications may want a two-way conversation." and
  "Requestor – Sends a request message and waits for a reply message." A
  reply is a second message that only a requestor needs; a message that needs
  none is complete when sent. This is why the default is one-way and only the
  sender who needs the answer asks for it.
- Microsoft Azure Architecture Center, "Priority Queue pattern"
  (https://learn.microsoft.com/en-us/azure/architecture/patterns/priority-queue):
  "The application or producer that sends a request assigns a priority value
  to the message, and consumers process the messages by priority." This is why
  the sender declares the blocked work and the receiver orders by it. The same
  page warns to "dynamically increase the priority of old messages to ensure
  that low-priority messages eventually get processed"; the oldest-first
  order of the rest and the rule that only blocked work is declared answer
  that risk.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, a real Blocking
  Request from another project reaches this repository's `/roundfix
  --follow-up`, the mode gives its entry `blocked_work`, and the command lists
  it first with basis `declared`.
  reason: it needs a real ai-memory message, the ai-memory MCP server and a
  live Agent session, which no Verification or QA row may reach.
  satisfied-by: the operator's follow-up session evidence after the release,
  recorded session evidence.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Characterization at `9fda0b65`:
  - the messages reader decodes each entry into a struct of `id`, `from`,
    `created_at`, `subject`, `text` and `type` with the standard JSON decoder,
    so an unknown key such as `blocked_work` is ignored without error; an
    older binary given a newer file therefore lists the entry as before;
  - an item reaches `unblocking` only when the judge answers
    `blocks_other_project` yes and `is_fix` or `is_small` yes; without an
    answer its type places it;
  - the item's `blocks` already names the sender when the judge says the
    message blocks another project; the documentation contract
    `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` pins the item's
    JSON field list in the command guide and both reference copies;
  - the three clauses are pinned in part by
    `TestEveryInboxClauseAnAdopterHoldsIsReplaced`,
    `TestTheGuidesSendCrossProjectWorkThroughAiMemory` and the closing-lesson
    check of `clause.core.message-the-owning-project`; the prescribed words
    keep every pinned phrase.
- A prototype of the three Tasks on a scratch copy at `9fda0b65`: the clause
  edits, the module record command, `make baseline-digests` and the managed
  refresh rewrote exactly the derived files task_02 declares, and a second
  refresh reported no change; `./internal/baseline`, the CLI's Baseline and
  Follow-up tests, `./skills` and the documentation contracts passed; the
  Behavior Surface Record moved only for the command guide.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for
  cross-project messages, replies and waiting. The top results were the
  mirrors of the `ai-memory-messaging` skill, this repository's
  `references/modes.md` and adopters' agent instructions, which hold the
  rule this Spec changes. `wiki/entities/ai-memory.md` records the decision of
  2026-10-25: "Pedidos, pesquisas e retornos usam exclusivamente ai-memory;
  ler não significa consumir, e o projeto destinatário decide a ação", which
  supports leaving what follows to the receiving project.
  `wiki/concepts/porta-de-captura-unica-do-fleet-inbox.md` confirms the inbox
  retirement and holds no reply rule.
- Exa: the two published sources above, and a third Microsoft page
  ("Prioritize critical flows", Well-Architected) that makes the same point
  about priority queues.

## Glossary

- adds: **Blocking Request**
- changes: **Cross-project Message**
- changes: **Follow-up Priority**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB brief for Spec 0373 of 2026-11-10, the maintainer's answers below and the adopted request.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~15:35 BRT), verbatim:

- Question "Spec 0373 (sua ideia): a regra de resposta entre projetos fica
  só no texto da Baseline, ou o `roundfix follow-up` também passa a ordenar
  por ela?" → "Texto + ferramenta (Recommended)" — option text: "As três
  cláusulas de mensagem da Baseline e a skill mudam, e o `roundfix follow-up`
  põe no topo a mensagem que declara trabalho bloqueado, sem depender do juiz
  Jev. Uma Task a mais; o release provavelmente sai como minor (v0.138.0)."
- Question "Spec 0373: quando chega uma mensagem que declara trabalho
  bloqueado, o que \"tratar como urgência\" obriga o projeto que recebe a
  fazer?" → "Próximo da fila (Recommended)" — option text: "Passa à frente de
  tudo que está pendente e é atendida logo depois do que já está rodando.
  Nenhuma Run ou entrega em curso é interrompida."
- Rules: question "Com que autorizações eu toco sem supervisão as três Specs
  deste ciclo: 0372 (review do Conexus), 0371 (relatório de QA, que você já
  decidiu às 06:20 pelo Jarvis) e 0373 (regra de resposta)?" → "Mesmas do
  BA, teto minor (Recommended)" — option text: "Escrever, entregar pela fila,
  merge e release de cada Spec assim que pronta (patch nas duas correções; a
  regra pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis,
  sem notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir."

Taken from those answers, the request and the characterization, with no
further interview:

- The three clauses are reworded in place and keep their ids; no clause is
  added, so no Source Baseline row is needed.
- The declaration is prose in the message; no code reads a message body for
  a marker. The session that writes the messages file fills `blocked_work`.
- A declared entry is not sent to the judge: no answer could move it from the
  top of `unblocking`, and sending it would cost a request and send its text
  out for nothing. ADR-0353's Waiting Item is the precedent.
- `blocked_work` requires `from`, so the report always names who is blocked;
  it is one line, like a Waiting Declaration's `waiting_on`.
- The messages file keeps `roundfix/follow-up-messages/v1`: the field is
  optional and an older binary ignores it.
- ADR-0383 records the decision and amends ADR-0319 and ADR-0330. Three
  implementation Tasks in series and a QA gate.

## Open Questions

None.
