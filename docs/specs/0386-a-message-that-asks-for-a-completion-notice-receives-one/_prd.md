---
spec: 0386-a-message-that-asks-for-a-completion-notice-receives-one
status: active
created: 2026-11-12
surfaces: [docs]
type: feat
---

# A message that asks for a completion notice receives one

The Baseline's receiving clause sends no acknowledgment, adoption notice or
completion notice to a Cross-project Message that declared no blocked work,
so a sender that asks to be told when its fix ships gets nothing. On
2026-11-12 Pantheon and gesttione-drive each asked for exactly that notice,
and the maintainer answered the operator's question with "Adicione exceção a
regra para quando houver uma solicitação de aviso." After this Spec a message
that asks to be told when its work is done receives one completion notice,
once the release that carries the work is published or, when the work
releases nothing, once it merges, naming what changed and the version; the
request is kept in the record of the message until then, and the notice is
sent by the session that follows the delivery, never by a Daemon-assigned
turn. Every adopter receives the reworded clauses through its next `roundfix
baseline update`, and the Roundfix skill's follow-up mode and release
reference carry the request to the notice.

## Sources

One maintainer request, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-12-maintainer-an-exception-for-a-message-that-asks-for-a-notice.md`](references/2026-11-12-maintainer-an-exception-for-a-message-that-asks-for-a-notice.md),
   the maintainer's instruction of 2026-11-12 in reply to the operator's
   question, both verbatim, with the two notice requests that prompted it;
   requested by the maintainer.

The two notice requests belong to the Specs that adopted their messages
(Pantheon's to Spec 0383, gesttione-drive's to Spec 0384); this Spec quotes
only the lines that ask for the notice. At the authoring base `docs/backlog/`
and `docs/findings/` hold no entry that shares this context.

## Prerequisites

None in `requires`. Authored against `73421db2` (v0.142.0 plus the plans of
Specs 0383, 0384 and 0385). Those three Specs deliver before this one, in that
order. Files this Spec shares with them, each edited in its own section:

- `CONTEXT.md`: each Spec its own terms; this Spec adds **Notice Request** and
  changes **Cross-project Message**.
- The Roundfix owned skill's version fields in the `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json` (0383 and 0385 raise them).
  task_02 raises the version with the record command from the tree it runs on
  and never writes a number. This Spec changes only the `modes` and `release`
  references, which none of the three edits.
- `docs/agents/setup-context.json` (0384 changes its Makefile digest) and the
  Baseline catalog snapshots and plan goldens that `make baseline-digests`
  rewrites: task_01 regenerates them from the tree it runs on.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted; the request is prose in the record of a message. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the notice is sent by a session through ai-memory,
  which no Roundfix command reaches. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0396 (new) records the decision: "A
  message that asks to be told when its work is done receives one completion
  notice". ADR-0390 is amended: its receiving clause "still sends no
  acknowledgment, adoption notice or completion notice to a message that
  declared none"; now a message that asks for a notice receives one. ADR-0383
  is amended: its receiving clause sends "no acknowledgment, adoption notice or
  completion notice to a message that declared no blocked work", and its
  ADR-0390 notice says the same. ADR-0330 is amended: its ADR-0390 notice says
  the receiving project "still sends no acknowledgment, adoption notice or
  completion notice to a message that declared no blocked work". ADR-0319
  applies unchanged: the four judge answers and the `blocked_work` exception
  order the Follow-up Command, and a request for a notice changes no rank.
  ADR-0338 applies unchanged: "No message the mode handed in stays pending
  when the turn ends"; a message that asks for a notice is recorded and
  consumed like any other. ADR-0353 applies unchanged: a message is never a
  Waiting Item, and a request for a notice is never watched with a Waiting
  Declaration. ADR-0095, which ADR-0330 cites, does not apply: it decides how
  the Secondbrain's knowledge base is read, and this Spec changes no
  Secondbrain access. ADR-0372, which cites ADR-0330, does not apply: it names
  where the Baseline's archive rule sends a copy and a profile's missing
  decision, and this Spec changes neither. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  instruction of 2026-11-12 quoted in the authorization record, taken on the
  cycle BD terms, which keep the grants of the Baseline source with its
  derived files, the guides under `docs/agents/` and the owned skills; the
  Governed Paths are the core Baseline module, the Standard TypeScript profile
  and its formatter golden, and this repository's managed guide and setup
  manifest that task_01 changes, and the Roundfix skill files task_02 changes.
  No Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0386-a-message-that-asks-for-a-completion-notice-receives-one/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/release.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`.

## Goals

- A message that asks to be told when its work is done receives exactly one
  completion notice, and no acknowledgment or adoption notice.
- The notice falls due when the release that carries the work is published,
  or at the merge when the work releases nothing, and names what changed and
  the version or merge.
- The request survives until the notice is due, in the record of the
  message, so a later session that follows the release finds it.
- Every other message that declared no blocked work still receives none of the
  three notices, and no project receives a release note it did not ask for.
- Every adopter gets the change through the Baseline, in any language its
  messages are written in.

## User Stories

1. As a sending project that removes a workaround once a fix ships, I want to
   ask for a notice when the version that carries it is published, so that I
   act on it without polling the other project.
2. As an Agent session triaging a message that asks for a notice, I want the
   rule to tell me to keep the request in the message's record, so that the
   session that releases the work weeks later still sends the notice.
3. As the session that publishes a release, I want the release reference to
   tell me how to find the requests the release fulfils, so that each sender
   gets one notice and nobody else gets one.

## Core Features

1. Receiving clause: a message that asks to be told when its work is done is
   the one exception to "no acknowledgment, adoption notice or completion
   notice to a message that declared none". It receives one completion notice
   and no acknowledgment or adoption notice; the request is kept in the record
   of the message's content; the notice is sent with `memory_message_send`
   once the work recorded from it is delivered (the release that carries it is
   published, or the merge when it releases nothing), naming what changed and
   the version or merge; a closure without delivery says so; the session that
   follows the delivery sends it, a Daemon-assigned turn never does. The
   clause keeps its id, force and Source Baseline replacements.
2. Sender's clause: a message may ask to be told when its work is done, only
   when the sending project will act on the notice, and the sender does not
   wait or watch for it. The clause keeps its id, force and replacements.
3. The Roundfix skill's follow-up mode keeps the request in the record it
   writes (the Backlog Entry or Finding, or the PRD `## Sources` entry that
   adopts the message), and its step 7 sends the notice after the delivery as
   the release reference says.
4. The release reference gains a `## Completion notices` section: after a
   release is published, list the Archive Records the release adds since the
   previous tag, read each Spec's PRD `## Sources` at the record's
   `source_revision`, and send one notice to each project whose message asked
   for one; never from a Daemon-assigned Run and only to the project that
   asked.
5. The glossary adds **Notice Request** and changes **Cross-project Message**;
   ADR-0396 records the decision; ADR-0390, ADR-0383 and ADR-0330 carry its
   amendment notices.

## User Experience

A Backlog Entry that records gesttione-drive's message notes, beside the
message's provenance, that the sender asked to be told what was done, the
version and how much Go cache stayed on the internal disk. After the release
that carries that Spec is published, the session that cut it lists the
Archive Records the tag added, reads that Spec's Sources at its
`source_revision`, and sends gesttione-drive one ai-memory message with the
change, the released version and the measured cache size.
Pantheon's message on the same release gets its own notice; Vortex, which
asked for none, gets nothing.

## Non-Goals / Out of Scope

- A field in the follow-up messages file, the Archive Record or any command
  that carries the request, and any command that sends a notice: no Roundfix
  command reaches ai-memory, and ADR-0396 keeps the request as prose the
  session reads. ADR-0396: "The request is prose the session reads, in the
  sender's language; no command carries it."
- Detecting a request inside a message body by a word, a phrase or a
  pattern: the session reads the message.
- Release notes or "what changes" messages to adopters that did not ask; the
  Update Check tells them.
- Sending the notices already due for Specs 0383 and 0384, or editing their
  folders to record the requests: those Specs deliver before this one, and
  their notices are the operator's to send on the maintainer's answer.
- Changing `clause.core.messages-are-self-contained`, adding, removing or
  renaming a Baseline clause or a Source Baseline row, the Follow-up Command,
  its schemas, the judge, the Waiting Declaration, the Spec Event
  Announcement, the release runbook or the vendored `ai-memory-messaging`
  skill.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository
   writes into `docs/agents/agent-instructions.md`, once each, "unless it asks
   to be told when its work is done", "receives one completion notice and no
   acknowledgment or adoption notice", "keep that request in the record of its
   content", "when the release that carries that work is published or, when
   that work releases nothing, when it merges", "a Daemon-assigned turn never
   sends one" and "A message may also ask to be told when the work it asks for
   is done"; a second refresh reports it current; the three message clause
   ids, their force and their Source Baseline replacements are unchanged. At
   `73421db2` none of these phrases is in the guide.
2. Success Metric: every Baseline test that exists at `73421db2` passes
   unchanged, `TestTheMessageClausesAskForAReplyOnlyForBlockedWork` and
   `TestTheReceivingClauseAnswersQuestionsAndBoundsTheDeclaredPrecedence`
   included, so the reply rule and the bounded precedence still hold.
3. Success Metric: both copies of the `modes` reference say, in step 4, that a
   Notice Request is the exception and that its record keeps the request and
   what it asked the notice to carry, and, in step 7, that the notice is sent
   after the release or the merge; both copies of the `release` reference
   carry the `## Completion notices` section.
4. Success Metric: `CONTEXT.md` defines **Notice Request** and the changed
   **Cross-project Message**; ADR-0390, ADR-0383 and ADR-0330 carry an
   "Amended by ADR-0396" notice; the documentation, skill and glossary
   contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write,
read through Exa on 2026-11-12.

- RFC 8098, *Message Disposition Notification*, section 2.1
  (https://www.rfc-editor.org/rfc/rfc8098.html): "The presence of a
  Disposition-Notification-To header field in a message is merely a request
  for an MDN." and "A user agent MUST NOT issue more than one MDN on behalf of
  each particular recipient." A notification goes only to a message that
  requested it, and only once; this is why the exception covers only a
  message that asks and sends one notice.
- Hohpe and Woolf, *Enterprise Integration Patterns*, "Return Address"
  (https://www.enterpriseintegrationpatterns.com/patterns/messaging/ReturnAddress.html):
  "The request message should contain a Return Address that indicates where to
  send the reply message." The replier keeps what the request says about its
  reply with the request; this is why the receiving project keeps the request
  in the record of the message's content until the notice is due.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, an adopter's
  `roundfix baseline update` writes the reworded clauses, and a real session
  that publishes a release carrying a Spec that adopted a Notice Request sends
  that sender one completion notice through ai-memory.
  reason: it needs a released binary, real ai-memory messages, the ai-memory
  MCP server and a live Agent session, which no Verification or QA row may
  reach.
  satisfied-by: the operator's release session evidence after the release,
  recorded session evidence.

## Validation

The decisions were checked on 2026-11-12 against these sources.

- Characterization at `73421db2`: `rule.core.cross-project-continuity` is at
  version 3; its receiving clause ends "send no acknowledgment, adoption
  notice or completion notice to a message that declared none. An empty inbox
  is rest: continue the session's current work.", rendered verbatim at line
  81 of this repository's `docs/agents/agent-instructions.md` and of the
  Standard TypeScript formatter golden of that guide, and carried by the
  catalog snapshots and plan goldens under `internal/baseline/testdata/`.
  `TestTheMessageClausesAskForAReplyOnlyForBlockedWork` pins "send no
  acknowledgment, adoption notice or completion notice to a message that
  declared none" once; the new wording keeps that phrase once. Step 4 of both
  `modes` copies ends with the same prohibition for "a message that declared
  no blocked work", and `TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender`
  pins the step's answer sentence, which stays. Neither `release` copy says
  anything about notices. `CONTEXT.md` **Cross-project Message** ends "sends no
  acknowledgment, adoption notice or completion notice to a message that is
  not a Blocking Request". No other guide, user-guide page or skill copy names
  a completion notice.
- A prototype of tasks 01 and 02 on a worktree at `73421db2`: the clause
  edits, the module record command, `make baseline-digests` and the managed
  refresh rewrote exactly the derived files task_01 declares, and a second
  refresh reported `current`; `./internal/baseline` and the CLI's Baseline
  tests passed with the new clause test; the skill edits, `make skills-sync`
  and the version record passed `./skills` with the new test; the glossary and
  ADR edits passed `-tags docscontract ./internal/docscontract`, including
  `TestGlossaryPhrasesAreDefinedTerms` and `TestTheSkillCoverageMapIsCurrent`,
  so the Behavior Surface Record does not move.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for completion
  notices on cross-project messages. The top results were
  `wiki/entities/ai-memory.md` and the mirror of this repository's ADR-0330,
  neither of which names a completion notice. `wiki/entities/ai-memory.md`
  records the decision of 2026-10-25: "Pedidos, pesquisas e retornos usam
  exclusivamente ai-memory; ler não significa consumir, e o projeto
  destinatário decide a ação", which supports sending the notice through
  ai-memory and leaving to the receiving project when it is due.
- Exa: the two published sources above. RFC 8098 also says a notification
  "MUST NOT be generated in response to an MDN", which this rule matches by
  sending nothing back to a completion notice; it leaves the recipient "free
  to silently ignore such a request", which this rule does not take: the
  maintainer made the notice an obligation.

## Glossary

- adds: **Notice Request**
- changes: **Cross-project Message**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BD batch 2 brief for Spec 0386 of 2026-11-12, the maintainer's instruction below and the characterization above.

The operator's question and the maintainer's instruction (chat, 2026-11-12
~18:40 BRT), verbatim:

- Question: "O Pantheon e o Drive pediram aviso quando a correção sair, e o
  Drive diz que o pedido do retorno veio de você. A regra da 0380 proíbe
  aviso de conclusão a mensagem que não declarou trabalho bloqueado, e
  nenhuma das duas declarou. Mando os dois avisos quando as Specs saírem, ou
  mantenho a regra?"
- Instruction: "Adicione exceção a regra para quando houver uma solicitação
  de aviso."

Taken from that instruction, the two requests and the characterization, with
no further interview:

- The term is **Notice Request**, beside **Blocking Request**: a Cross-project
  Message that asks to be told when the work it asks for is done. The
  Baseline clause says it in plain words, because adopters do not share this
  glossary.
- "Delivered" means the release that carries the work is published, or the
  merge when the work releases nothing: both requests ask for the version,
  which exists only after the release.
- The notice names what changed, the version or merge, and anything else the
  request asked it to carry; a closure without delivery is the same one
  notice, saying so, so no sender keeps a request nobody will answer.
- The request survives as prose in the record of the message (the Backlog
  Entry or Finding, or the PRD `## Sources` entry that adopts it), and the
  release reference reads it at the Archive Record's `source_revision`. It is
  text-level because no Roundfix command reaches ai-memory to send the notice,
  and because a field would make the session match a marker in the sender's
  language, which ADR-0383 rejects.
- The session that follows the delivery sends it; a Daemon-assigned turn
  never does, as with the Spec Event Announcement.
- The sender's clause lets a sender ask only when it will act on the notice;
  the content clause is unchanged.
- `type: feat`, as Spec 0380, which changed the same clauses; the Release
  Ceiling is `minor`. ADR-0396 amends ADR-0390, ADR-0383 and ADR-0330. Two
  implementation Tasks in series and a QA gate.

## Open Questions

None.
