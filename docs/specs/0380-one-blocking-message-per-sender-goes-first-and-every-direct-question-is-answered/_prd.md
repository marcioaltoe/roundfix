---
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
status: active
created: 2026-11-11
surfaces: [backend, cli, docs]
type: feat
---

# One blocking message per sender goes first, and every direct question is answered

Spec 0373 gave every Cross-project Message that declares blocked work a place
ahead of all of the receiving project's pending work, and told the receiver to
answer only such messages. Two adopters' pre-PR reviews found the gaps on
2026-11-11: a message that declares no blocked work but asks a real question
has no rule that answers it (reported by conexus), and a body any sender
writes can put every message it sends ahead of all pending work without limit
(reported by quaestor). After this Spec only the oldest blocking message of
each sending project takes that precedence, the sender's later ones wait
their normal turn, the precedence never replaces the maintainer's judgment,
and every direct question receives an answer once the answer exists. Every
adopter receives the reworded clauses through its next `roundfix baseline
update`, and the Follow-up Command applies the same bound.

## Sources

Two messages, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-11-conexus-a-question-without-declared-blocked-work-has-no-reply-rule.md`](references/2026-11-11-conexus-a-question-without-declared-blocked-work-has-no-reply-rule.md),
   the ai-memory message `01a1c9ec-5406-7170-bc20-0a69c687b8ae` of
   2026-11-11, reported by conexus, nature `inferred`, with the operator's
   note on the receiving clause.
2. [`references/2026-11-11-quaestor-a-declared-blocked-message-jumps-every-queue-without-bound.md`](references/2026-11-11-quaestor-a-declared-blocked-message-jumps-every-queue-without-bound.md),
   the ai-memory message `01a1cbb3-63ed-7be3-9cca-6e806ef76e54` of
   2026-11-11, reported by quaestor, nature `measured`, with the operator's
   measurement of the Follow-up Command at `013af3b4`.

Both are findings on the same rule, ADR-0383's, so one Spec adopts both. At
the authoring base `docs/backlog/` and `docs/findings/` hold no entry that
shares this context.

## Prerequisites

None in `requires`. Authored against `013af3b4` (v0.139.0). Specs 0377, 0376,
0378 and 0379 deliver before this one, in that order, and Specs 0381 and 0382
after it. Files this Spec shares with them, each edited in its own section:

- The Roundfix owned skill's version fields in the `SKILL.md` pair and
  `skills/testdata/owned-skill-versions.json` (0376 and 0379 raise them, 0381
  and 0382 may). Each Task raises the version with the record command from
  the tree it runs on and never writes a number.
- The `modes` reference pair: 0376 changes step 1 of `## --update`; this Spec
  changes steps 1, 2 and 4 of `## --follow-up` only.
- `CONTEXT.md` (each Spec its own terms; this Spec **Cross-project Message**,
  **Blocking Request** and **Follow-up Priority**) and
  `docs/references/behavior-surfaces.json` (each Task re-records it after its
  own last guide edit).
- `docs/agents/setup-context.json` (0377 and 0378 change the Makefile digest
  through the managed refresh; task_02 refreshes it from the tree it runs
  on) and the Baseline catalog snapshots and plan goldens that `make
  baseline-digests` rewrites.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; the messages file keeps `roundfix/follow-up-messages/v1` and the
  report `roundfix/follow-up/v1`, which gains one optional key. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; a held-back declared entry reaches the judge through
  the existing transport, as any undeclared entry does. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0390 (new) records the decision:
  "One blocking message per sender takes precedence, and every direct
  question is answered". ADR-0383 is amended: its receiving clause
  "triages such a message first" for every declared message and answers only
  "when its message declared blocked work", and its command ranks every
  entry with `blocked_work` first; now only the oldest per sender takes the
  precedence and every direct question is answered. ADR-0319 is amended: its
  ADR-0383 notice places every message entry that carries `blocked_work` in
  `unblocking`. ADR-0330 is amended: its ADR-0383 notice says the receiving
  project "answers only it". ADR-0338 applies unchanged: "No message the mode
  handed in stays pending when the turn ends", and a held-back declared
  message is recorded and consumed like any other. ADR-0353 applies
  unchanged: "A Waiting Item is never sent to the judge, with or without
  --all"; a message is never a Waiting Item. ADR-0372, which cites ADR-0330,
  does not apply: it names where the Baseline's archive rule sends a copy and
  a profile's missing decision, and this Spec changes neither. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  BC answers of 2026-11-11 quoted in the authorization record, which keep the
  cycle BB grants of the Baseline source with its derived files, the guides
  under `docs/agents/` and the owned skills; the Governed Paths are the
  Roundfix skill files tasks 01 and 03 change, and the core Baseline module,
  the Standard TypeScript profile and its formatter golden, and this
  repository's managed guide and setup manifest that task_02 changes. No
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`.

## Goals

- A direct question in a Cross-project Message receives an answer once the
  answer exists, whether or not the message declared blocked work.
- No sender can put more than one message at a time ahead of a receiving
  project's pending work by declaring blocked work.
- The rule and the skill say that the declared precedence never replaces the
  maintainer's judgment.
- The Follow-up Command applies the same bound, and its report says which
  declared message holds the precedence and which ones wait.
- Every adopter gets the change through the Baseline, in any language its
  messages are written in.

## User Stories

1. As an Agent session reading a message that declares no blocked work but
   asks a real question, I want the rule to tell me to answer it, so that the
   sender is not left without an answer it asked for.
2. As an Agent session in a receiving project, I want only one blocking
   message per sender at a time to go ahead of my pending work, so that a
   sender that declares blocked work on every message cannot hold my work
   back without end.
3. As the maintainer running `/roundfix --follow-up`, I want the list to show
   which declared message of a sender holds the precedence and why the others
   are ordered normally, so that I can decide with the whole picture.

## Core Features

1. Reply rule: the receiving clause answers each direct question a message
   asks as soon as the answer exists, whether or not the message declared
   blocked work; acknowledgments, adoption notices and completion notices
   stay forbidden for a message that declared none. The sender's clause adds
   that the receiving project owns that answer too.
2. Bounded precedence in the Baseline: the receiving clause triages first,
   from each sending project, the oldest message that declares blocked work,
   then the rest oldest first, a later blocking message from the same sender
   included, and says the precedence never replaces the maintainer's
   judgment. The clauses keep their identities, force and Source Baseline
   dispositions; no clause is added or removed.
3. Bounded precedence in the Follow-up Command: of the message entries with
   `blocked_work` that share a trimmed `from`, only the oldest, then the
   lowest reference, ranks first with basis `declared`. Every other declared
   entry of that sender is sent to the judge and ordered as if it declared
   nothing, keeps `blocked_work`, and carries `precedence_held_by` naming the
   holding entry, in JSON and on its text line.
4. The Roundfix skill's follow-up mode fills `from` with the sending project
   ai-memory returns, gives `blocked_work` to every declaring message, keeps
   the place the command gives a held-back one, says the precedence never
   replaces the maintainer's judgment, and answers each direct question.
5. The glossary changes **Cross-project Message**, **Blocking Request** and
   **Follow-up Priority**; ADR-0390 records the decision; ADR-0383, ADR-0319
   and ADR-0330 carry its amendment notices.

## User Experience

A sender's later declared message keeps its work on its text line, for
example
`4. other message 2026-10-25 message:d3: Second question — from argus — declares blocked work: release 1.5 waits on the label — precedence held by message:d1`.
Its JSON item has `basis` `judge` or `type`, `blocked_work` its work and
`precedence_held_by` `message:d1`. The holding entry's line is unchanged from
Spec 0373.

## Non-Goals / Out of Scope

- Verifying a declaration against the sending project's repository, or
  limiting the precedence to known senders: ADR-0390 rejects both.
- Remembering across runs or sessions which declared message already took
  the precedence: the command keeps nothing between runs.
- Comparing senders without regard to case or across workspaces.
- Detecting a declaration or a question inside a message body by a word, a
  phrase or a pattern: the session reads the message.
- Adding, removing or renaming a Baseline clause or a Source Baseline row;
  changing the content clause `clause.core.messages-are-self-contained`.
- Changing `--status`, the four judge questions, their thresholds, the
  Waiting Declaration, the schema names, or the vendored
  `ai-memory-messaging` skill.
- Telling adopters by message: they receive the rule through their next
  Baseline update and the Update Check.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository
   writes into `docs/agents/agent-instructions.md`, once each, "the answer to
   a direct question the message asks included", "takes precedence for one
   message per sender at a time", "a later blocking message from the same
   sender included", "never replaces the maintainer's judgment" and "Answer
   each direct question a message asks with `memory_message_send` as soon as
   the answer exists, whether or not the message declared blocked work"; a
   second refresh reports it current; the three message clause ids, their
   force and their Source Baseline replacements are unchanged. At `013af3b4`
   none of these phrases is in the guide.
2. Success Metric: given two declared entries from one sender (the newer
   listed first), one from another sender and an undeclared one, the
   Follow-up Command lists the older entry of each sender with basis
   `declared`, sends the newer one of the first sender to the judge with the
   undeclared one, and reports it with its `blocked_work` and
   `precedence_held_by` naming the older one. At `013af3b4` all three
   declared entries rank first with basis `declared` and the judge receives
   only the undeclared one.
3. Success Metric: every Follow-up Command test that exists at `013af3b4`
   passes unchanged, so a file whose declared entries come from distinct
   senders, or that has none, keeps its report bytes.
4. Success Metric: both copies of the `modes` reference say, in step 1, that
   `from` is the sending project ai-memory returns and that every declaring
   message gets `blocked_work`; in step 2, that a held-back message keeps the
   command's place and that the precedence never replaces the maintainer's
   judgment; in step 4, that each direct question is answered whether or not
   the message declared blocked work.
5. Success Metric: the command guide and both `follow-up` references document
   the bound and `precedence_held_by`; `CONTEXT.md` carries the changed
   entries; ADR-0383, ADR-0319 and ADR-0330 carry an "Amended by ADR-0390"
   notice; and the documentation, skill and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write,
read through Exa on 2026-11-11.

- Kubernetes documentation, "API Priority and Fairness"
  (https://kubernetes.io/docs/concepts/cluster-administration/flow-control/):
  "It may be `ByUser`, in which one requesting user will not be able to
  starve other users of capacity" and "it is valuable to prevent one stream
  of requests from starving others (in particular, in the relatively common
  case of a single buggy client flooding the kube-apiserver with requests,
  that buggy client would ideally not have much measurable impact on other
  clients at all)". A priority a client claims is bounded per client so that
  one source cannot starve the rest; this is why the precedence is bounded
  per sender.
- Hohpe and Woolf, *Enterprise Integration Patterns*, "Request-Reply"
  (https://www.enterpriseintegrationpatterns.com/patterns/messaging/RequestReply.html):
  "The applications may want a two-way conversation." A message that asks a
  question opens that two-way conversation whether or not its sender is
  blocked; this is why every direct question is answered, while the sender
  that does not depend on the answer still does not wait for it.

## Unreachable Acceptance

- criterion: after the release that carries this Spec, an adopter's
  `roundfix baseline update` writes the reworded receiving clause, and a real
  `/roundfix --follow-up` session that holds two declared messages from one
  sender lists the older one first and the newer one with
  `precedence_held_by`.
  reason: it needs a released binary, real ai-memory messages, the ai-memory
  MCP server and a live Agent session, which no Verification or QA row may
  reach.
  satisfied-by: the operator's follow-up session evidence after the release,
  recorded session evidence.

## Validation

The decisions were checked on 2026-11-11 against these sources.

- Characterization at `013af3b4` with the built binary, a temporary
  repository and a messages file holding `d1` and `d3` from `argus`
  (2026-10-10 and 2026-10-25), `d2` from `quaestor` and `j1`, a `fix` from
  `fluxus`, no judge key: the text list ranked `d1`, `d2` and `d3` as
  `unblocking` with ` — blocks <from>: <blocked_work>` and `j1` fourth; the
  JSON gave all three declared entries basis `declared`.
  `followUpDeclaredOrder` sorts every entry with `blocked_work` and
  `runFollowUpCommand` keeps all of them from the judge. Every existing
  declared-work test uses distinct senders.
- The receiving clause at `013af3b4` reads "Triage first each message whose
  sender declares its work blocked until this project answers" and "Answer
  the sender with `memory_message_send` as soon as the answer exists when its
  message declared blocked work"; step 4 of the `modes` reference reads
  "only when its message declared blocked work", which
  `TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork` pins.
- A prototype of the three Tasks on a scratch copy at `013af3b4`: the
  Follow-up Command change, its three tests and every existing
  `FollowUp` test passed; the clause edits, the module record command, `make
  baseline-digests` and the managed refresh rewrote exactly the derived
  files task_02 declares, and a second refresh reported no change;
  `./internal/baseline`, `./skills`, the documentation contracts with the
  glossary test and the follow-up guide contract passed.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for
  cross-project messages, blocked work and replies. The top results were the
  mirror of the `ai-memory-messaging` skill, which says nothing on replies or
  precedence, and the mirror of this repository's agent instructions, which
  holds the rule this Spec changes. `wiki/entities/ai-memory.md`
  records the decision of 2026-10-25: "Pedidos, pesquisas e retornos usam
  exclusivamente ai-memory; ler não significa consumir, e o projeto
  destinatário decide a ação", which supports keeping the maintainer's
  judgment over the declared order.
- Exa: the two published sources above, and the Azure Architecture Center
  "Priority Queue pattern" that Spec 0373 cites, whose warning on starving
  low-priority messages this bound answers for a single sender.

## Glossary

- changes: **Cross-project Message**
- changes: **Blocking Request**
- changes: **Follow-up Priority**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BC brief for Spec 0380 of 2026-11-11, the maintainer's answers below and the two adopted messages.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-11 ~14:45–14:50 BRT), verbatim:

- Question "Specs novas D (mensagens: pergunta não bloqueante recebe
  resposta; precedência de 'bloqueado' limitada) e E (review pré-PR revisa
  só o delta e com contexto depois de merge da main): como autorizar?" →
  "Após a 0379 (Recommended)" — option text: "Mesmas autorizações do ciclo
  BB (push, merge, release, teto minor); entregues depois da fila
  0377→0379, D antes de E."
- Question "Spec D (mensagens entre projetos): como limitar a precedência de
  uma mensagem que declara trabalho bloqueado, já que o corpo não é
  confiável (achado do Quaestor)?" → "Uma por remetente + juízo
  (Recommended)" — option text: "Só a mensagem bloqueante mais antiga de
  cada remetente passa à frente; as demais seguem a ordem normal. O texto
  diz que a precedência não dispensa o juízo do maintainer. Pergunta direta
  não bloqueante recebe resposta quando existir."

Taken from those answers, the messages and the characterization, with no
further interview:

- The sender is the message entry's `from`, trimmed and compared byte for
  byte; the follow-up mode fills it with the `from_project` ai-memory returns
  with the message, never with a name the body gives.
- "Oldest" is chosen per run from the file the command reads: the oldest
  `created_at`, then the lowest reference. The command keeps nothing between
  runs, so after the holding message is consumed the sender's next declared
  message holds the precedence.
- A held-back declared entry goes to the judge and keeps the basis that
  ordered it (`judge` or `type`); a new optional key `precedence_held_by`
  names the holding entry, and its text line says so. A distinct basis would
  hide how its class was computed.
- The three clauses keep their ids; only the sender's and the receiving
  clause change words, and the rule's version rises. The content clause is
  unchanged.
- ADR-0390 records the decision and amends ADR-0383, ADR-0319 and ADR-0330.
  Three implementation Tasks in series and a QA gate.

## Open Questions

None.
