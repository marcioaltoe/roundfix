---
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
status: active
created: 2026-10-22
surfaces: [backend, docs]
---

# The Roundfix skill's follow-up, status and unattended modes, and up to three questions per call

A maintainer who wants a repository's agent to pick up everything pending,
ask him every decision at once and then keep working alone has no way to say
so today: the Roundfix skill has no modes, and the Baseline forbids asking
more than one question per call. After this Spec the Roundfix skill reads a
Skill Mode from the request in Claude Code, Codex and OpenCode alike.
`--follow-up` inventories the pending work through Roundfix's read-only
follow-up command, plans it, gathers every decision, triages, authors and
queues the Specs, and arms its own continuation. `--status` gives a short
read-only report, and `--unattended` gathers decisions and authorizations in
calls of up to three independent questions and records each answer where the
work reads it. The Baseline clause that governs questions allows up to three
independent questions per call for every adopter.

## Sources

This Spec is a secondary Spec of one source. Spec 0306, which delivers the
read-only follow-up command, owns the source and adopted it as a promoted
Backlog Entry; this Spec links the owner's copy and adopts nothing:

- [`../0306-the-follow-up-command-lists-pending-work-in-priority-order/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md`](../../references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md),
  type `feat`: the maintainer's request of 2026-10-22 and the design settled
  with him, the operator's proposal for both follow-up Specs. This Spec
  delivers its skill and Baseline half (Spec B): the modes, decision 3 on
  questions per call and decision 4 on where answers go. The maintainer's
  decisions are quoted in `_authorization.md`.

At the authoring base `docs/backlog/` held one open entry, the QA refusal
report's audited commit, which Spec 0302 adopts, and `docs/findings/` held
none; neither shares this context.

## Prerequisites

None in `requires`: no Verification command of this Spec runs code of another
active Spec. This Spec is authored against v0.109.0 with Specs 0299 and 0300
active, and the operator delivers it after the cycle AP Specs 0301 to 0306.
Spec 0306 delivers the read-only follow-up command that the `--follow-up` and
`--status` modes run; the Tasks that describe those modes take the command's
name and options from the tree they start on and never name a command or flag
the tree lacks; Spec 0306 plans it as `roundfix follow-up [--status]
[--messages <path>] [--format <text|json>]`. The Sources link resolves to
Spec 0306's planning change, which merges with or before this one. Several earlier Specs raise the Roundfix skill's version and
add a reference to it, so this Spec's skill Tasks raise the version and the
reference count over what the tree holds when they start.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Baseline clause keeps its identity
  `clause.core.ask-user-answerable-decisions` and the Roundfix skill keeps its
  name. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; the skill text describes optional sources that the
  session or the follow-up command already reaches, and every test reads
  repository files or temporary repositories. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0320, "A question call carries up
  to three independent questions", which this Spec records and which also
  places each gathered answer; and ADR-0095, "The secondbrain inbox is the
  single capture door", whose capture and triage contract the `--follow-up`
  mode follows. ADR-0095: "committed there at the moment of capture".
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."), and the batch 2 decisions of the
  same day. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/modes.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter or Project Config change.

## Goals

- Every adopter's agent may ask up to three independent questions in one call,
  never more than the tool accepts, with every other safeguard of the clause
  kept, and an adopter's refresh records the clause `retained`.
- The Roundfix skill reads `--follow-up`, `--status` or `--unattended` from
  the request in all three harnesses, and a request without one behaves as
  today.
- `--status` reports the repository's state without writing, asking or
  consuming anything.
- `--follow-up` runs its steps in a fixed order, uses each optional source
  only when it is present, says which were missing, and never ends a turn
  without arming its continuation or saying that work stops.
- `--unattended` gathers every decision and authorization before work
  continues and records each answer where the work reads it, never an
  authorization outside a Spec's `_authorization.md`.

## User Stories

1. As a maintainer, I want to type `/roundfix --follow-up`, `$roundfix
   --follow-up` or `@roundfix --follow-up` and have the agent pick up all
   pending work of the repository, so that I do not have to list it myself.
2. As a maintainer, I want every decision the pending work needs asked in as
   few calls as possible, so that I answer once and the agent then works alone.
3. As a maintainer, I want `--status` to tell me what is running, what is
   queued, what the next release brings and what waits on me, without
   changing anything.
4. As a maintainer of several repositories, I want the modes to work whether
   or not a repository has a Secondbrain, ai-memory, a Jev key or a release
   announcement, and to tell me what they could not read.
5. As a later session, I want an answer given today to be found in the Spec's
   `_authorization.md` or in the Backlog Entry or Finding it concerns, so that
   I do not ask again what was decided, and so that I do ask again for an
   authorization nobody granted to a Spec.

## Core Features

1. The question clause. The Baseline clause on cheap user decisions allows at
   most three questions per call or message, never more than the tool in use
   accepts, keeps one call pending, joins only independent questions, numbers
   the questions in the chat fallback, and forbids merging two decisions into
   one question. It keeps its identity and stop-and-ask enforcement.
2. Mode detection. The skill's entry file gains a modes section naming the
   three modes and the three invocation forms, its metadata exposes the skill
   as an OpenCode slash command, its body carries no argument placeholder that
   Claude Code would substitute, and a new reference describes how a mode is
   chosen and how each optional source is detected.
3. The status mode. A read-only report of work in progress, pending Specs in
   delivery order, the next release and what it brings, blockers, and what
   depends on the user.
4. The follow-up mode. Inventory, plan, gather decisions, triage, author,
   check, queue and arm, in that order, with other projects' findings
   captured in their own inbox and a closing report.
5. The unattended mode. Gather information, decisions and authorizations;
   ask in calls of up to three independent questions with the recommended
   option first, or numbered in text when no tool exists; record each answer
   in the consuming Spec's `_authorization.md` and PRD, or in a Backlog Entry
   or Finding without authorizations; arm the continuation.
6. Glossary. **Skill Mode** is defined in the glossary.

## User Experience

The maintainer types the skill with a mode. A `--status` request answers in
one short report. A `--follow-up` request shows a short staged plan, then one
to a few question calls of up to three questions each, then works and reports
what it did, what is queued or armed, which sources were missing and what is
still waiting on him. Without a structured question tool the same questions
come numbered in the chat and the agent waits.

## Non-Goals / Out of Scope

- The inventory, priority and status computation: Spec 0306's read-only
  follow-up command owns them, and this Spec only runs it.
- Changing the owned interview skills (`brainstorming`, `write-idea`,
  `write-prd`, `write-techspec`, `business-analyst`) that ask one question per
  message: their questions depend on the previous answer, which the new
  clause still keeps sequential.
- Sending messages to people or to other projects' agents, or naming any
  personal assistant or announcement channel in product text.
- A user-guide page for the modes, a new CLI flag, or any change to the
  Delivery Queue, the release commands or the follow-up command.
- Changing the Source Baseline corpus, its manifest or its retention
  transitions.

## Success Metrics

1. Success Metric: the embedded catalog's clause
   `clause.core.ask-user-answerable-decisions` carries the three-question
   limit and the independence sentence once, keeps `stop-and-ask` and no
   `replaces`, and no longer says "Ask exactly one question per call or
   message"; the Standard TypeScript golden and this repository's
   `docs/agents/agent-instructions.md` render it.
2. Success Metric: a Standard TypeScript adopter refreshed with the new
   catalog records the clause `retained` with itself as target, and no clause
   is unaccounted.
3. Success Metric: the canonical and embedded Roundfix skill name the three
   modes and the three invocation forms in a modes section, link the modes
   reference, set the OpenCode slash metadata, carry no argument placeholder,
   and their modes reference holds the sections on choosing a mode, optional
   sources and `--status`.
4. Success Metric: the modes reference's `--follow-up` section lists its steps
   in the order inventory, `--unattended`, triage, authoring, strict check,
   queue, and its `--unattended` section states the three-question limit the
   Baseline clause carries and the two answer-placement sentences.
5. Success Metric: no Roundfix skill file names a personal assistant or an
   announcement channel, the skill's version is raised and recorded, the
   mirrors match, and the glossary defines **Skill Mode**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Source Baseline corpus committed long before this Spec,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/`,
  whose manifest records `clause.core.ask-user-answerable-decisions` with
  enforcement `stop-and-ask`: the identity a fleet adopter's refresh compares
  with, which proves the `retained` disposition on real adopter evidence
  rather than on this Spec's fixture.
- OpenAI's Codex pull request "Feat: request user input tool"
  (https://github.com/openai/codex/pull/9472), whose tool schema says:
  "Questions to show the user (1-3). Prefer 1 unless multiple independent
  decisions block progress."
- Claude Code's skills documentation
  (https://code.claude.com/docs/en/skills): "If you invoke a skill with
  arguments but no placeholder in the skill's content receives one, Claude
  Code appends `ARGUMENTS: ` to the end of the skill content so Claude still
  sees what you typed."

## Unreachable Acceptance

- criterion: a live session in each harness reads the mode from `/roundfix --follow-up` in Claude Code, `$roundfix --follow-up` in Codex and `@roundfix --follow-up` in OpenCode, and runs the mode the skill describes
  reason: it needs interactive sessions of three external agent harnesses and a person to answer their questions, which a Run sandbox has neither of
  satisfied-by: the maintainer's first follow-up session after the release, the use he asked for on 2026-10-22, recorded as this Spec's unproven step

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at v0.109.0. The clause text in the core module and in
  this repository's agent instructions says "Ask exactly one question per
  call or message and keep only one question pending, even when the tool
  accepts an array of questions." and "Never bundle questions"; the Source
  Baseline manifest records the clause with `stop-and-ask`; the Roundfix skill
  has 18 references and no modes; the owned interview skills ask one question
  per message. A prototype in a scratch worktree applied the reworded clause,
  ran the module record command, `make baseline-digests` and the Managed
  Refresh, measured the derived files the TechSpec lists, passed the
  `internal/baseline`, `internal/cli` baseline and documentation contract
  suites, and showed the refresh idempotent. The TechSpec records the
  prototype's tests and the sabotages that failed them.
- Secondbrain, read only. Read `wiki/index.md` and
  `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md`, which records the
  2026-09-08 research behind the one-question rule and states that the limit
  was the maintainer's choice while the tools accept groups (Claude Code 1-4,
  Codex 1-3, OpenCode a list): this Spec reverses that choice on the
  maintainer's new answer, not on new tool limits. Ran `qmd query
  "structured user questions per call AskUserQuestion request_user_input batch
  of questions agent skill" --all --files --min-score 0.3`, which returned
  mirrors of agent instructions and the brainstorming skill, and `qmd query
  "supervisor follow-up triage inbox authorizations unattended session
  continuation wake-up" --all --files --min-score 0.3`, which returned the
  Roundfix mirror's ADR-0095 and Spec 0187's record of a queue that recovers
  without a supervisor; both confirmed the capture door and the armed
  continuation the follow-up mode follows. No Inbox Entry was written,
  because the briefing forbids writing there.
- Exa: the Codex pull request and the Claude Code skills page above, OpenCode's
  skills page (https://opencode.ai/docs/skills/: "`opencode/slash`: set to
  `"true"` to expose the skill as a `/name` command in the TUI") and tools page
  (https://opencode.ai/docs/tools/: "When there are multiple questions, users
  can navigate between them before submitting all answers."). Challenging:
  Codex's own schema says "Prefer 1", which is why the clause allows several
  questions only when they are independent. The Claude Code page also says a
  `$` followed by a digit or `ARGUMENTS` in a skill body is substituted, which
  is why the skill's entry file must carry no such text.

## Glossary

- adds: **Skill Mode**

## Decisions

- 2026-10-22: Where do the follow-up modes live? — "Skill roundfix, fontes opcionais (Recomendado)" (option text: "Os modos entram na skill do produto. Cada fonte é detectada e opcional: Secondbrain pelo guia do baseline, ai-memory se o MCP existir, Jev se houver chave, bloco do GT3 só se as instruções do usuário o definirem. Sem a fonte, o modo segue e diz o que faltou.")
- 2026-10-22: Code or skill? — "Comando CLI + skill (Recomendado)"; the command is Spec 0306's.
- 2026-10-22: How many questions per call? — "Até 3 no geral".
- 2026-10-22: Where do answers go? — "As respostas e autorizacoes devem ser parte da spec que será implementada, como o _autorizations.md e se não for virar uma spec deve virar um backlog ou finding e armazenado com as informaçoes e decisoes (sem autorizacoes) para virar uma spec no futuro."
- 2026-10-22: Wait for the assistant's reply? — "Sim, suas decisões bastam (Recomendado)".

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 2 briefing for Spec 0307 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The clause keeps its identity and enforcement, so the refresh records it
  `retained`; no `replaces` entry and no new clause are needed. See ADR-0320.
- One reference, `references/modes.md`, holds the three modes, because they
  share detection, sources and the answer rule, and it is indexed as a topic
  row, not a command family.
- The skill exposes itself as an OpenCode slash command through metadata.
- The interview skills keep one question per message.

## Open Questions

None.
