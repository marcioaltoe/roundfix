---
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: active
created: 2026-10-22
surfaces: [backend, cli, docs]
---

# The Follow-up Command lists a repository's pending work in priority order

A session that picks up a repository's pending work rebuilds the same
inventory by hand every time: it reads the Secondbrain inbox, the open
Backlog Entries and Findings, the active Specs, the Runs, the Delivery Queue,
the open Pull Requests and git, and then guesses an order. The maintainer
asked for one read-only command that does it the same way in every harness.
After this Spec, `roundfix follow-up` lists every pending item with its
origin, reference, type, date and the project it blocks, orders the list
from four Jev yes/no answers per item, says what it checked and what was
missing, and says when there is nothing to do; `roundfix follow-up --status`
shows the work in progress, the pending Specs in delivery order, the next
release and what waits on the maintainer. The skill modes that act on this
report belong to Spec 0307.

## Sources

One Backlog Entry, written by this Spec's planning change directly into
`references/` and adopted at birth (see `references/_index.md`):

1. [`references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md`](references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md),
   type `feat`: the operator's proposal for the `--follow-up`, `--status` and
   `--unattended` modes, translated from Portuguese with the original quoted.
   This Spec delivers its "Spec A"; Spec 0307 delivers its "Spec B" and links
   this copy.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` held none.

## Prerequisites

None that the Task Graph enforces. This Spec is authored against `42f523c2`
(v0.109.0) beside the active Specs 0299 and 0300 and the cycle AP Specs
0301-0305, and delivers after them. It touches none of their functions: it
adds a command, two User Config keys, one judgment kind and one Roundfix
skill reference. Each documentation Task raises the Roundfix skill version
over whatever the tree holds when it starts. Before this plan enters a Run,
the operator rebases it on the default branch and runs `roundfix spec check
0306-the-follow-up-command-lists-pending-work-in-priority-order --strict
--run-verification` again.

## Project Constraints

- Identifier strategy: not applicable — no persisted identifier, schema or
  table is minted; the report's `schema` value `roundfix/follow-up/v1` and the
  messages file's `roundfix/follow-up-messages/v1` name output formats, like
  every other `--format json` report. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the command sends Follow-up Items to
  the judge through the existing transports, keys and Judge Log, and lists
  open Pull Requests through the installed `gh`, whose authentication it never
  reads; no endpoint, credential or recipient is added, and every test uses a
  fake transport and a fake `gh`. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0319: "The order is computed in
  code, never asked.", which also records how the judge's earlier decisions
  apply: the advisory rule, the recipients, the keys, the Judge Log, the
  monthly ceiling and the Run Marker stay as they are, and only the data
  boundary widens to Follow-up Items; ADR-0209: "The data boundary grows by
  one reader.", the boundary ADR-0319 widens the same way, and ADR-0208,
  which cites it and is unaffected because the command adopts no source and
  suggests no grouping; ADR-0095: "The secondbrain inbox is the single
  capture door", whose positional status the inventory reads; ADR-0143:
  "Release planning is read-only, including its reset mode", which the status
  view calls; ADR-0256: "A release waits for the skills that describe a
  changed surface", which the new command and keys meet through the new
  reference. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, authorization "same as
  cycle AO", which carries the cycle AJ grant "Mexer em Governed Paths
  declarados" ("'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0306-the-follow-up-command-lists-pending-work-in-priority-order/_authorization.md`;
  bounded files: `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint,
  formatter, Baseline module or Project Config file change.

## Goals

- One command lists every pending item of the repository once, with what it
  checked and what was missing, in Claude Code, Codex and OpenCode alike.
- The order follows the maintainer's four classes, decided in code from Jev's
  answers, and still orders by type, saying why, when Jev is unavailable.
- An empty inventory with nothing in progress says "nothing to do" and lists
  what was checked, in a state a caller can test.
- A status view answers what is running, what is next, what the next release
  brings and what waits on the maintainer, without writing anything.

## User Stories

1. As the maintainer starting a follow-up session, I want one command that
   lists the inbox, the messages, the Backlog, the Findings and the pending
   Specs together, so that nothing waiting is missed.
2. As the maintainer, I want fixes that block another project first, then the
   system's health, then valuable features, so that the order matches how I
   would triage.
3. As an agent driving the skill, I want a JSON report with a stop state and
   each source's state, so that I can stop when there is nothing to do and
   say which source was missing.
4. As the maintainer, I want a short status of Runs, the queue, the next
   release and the decisions waiting on me, so that I know where things stand
   without reading four commands.

## Core Features

1. Inventory. The command lists the pending Inbox Entries of the repository's
   inbox namespace under the Secondbrain Root (oldest first), the messages of
   an optional messages file, the open Backlog Entries, the unresolved
   Findings, the notes in `docs/_inbox/` and the active Specs, each as one
   Follow-up Item with its kind, reference, type, date, title, the project
   that sent it and, when Jev says so, the project it blocks.
2. Sources. Each source reports `checked` with its count, `missing` (not
   configured or not provided) or `unavailable` (reading failed) with a
   detail; a missing or unavailable optional source never fails the command.
3. In progress. Beside the list, the command shows this repository's active
   Runs, its unfinished Delivery Queue items with their Park Class and next
   action, its open Pull Requests and its git state.
4. Follow-up Priority. Jev answers four yes/no questions per item, the code
   orders the items into the four classes, and without an answer the item's
   type orders it; the report names the judge's state, the model that
   answered and the reason when it did not.
5. Stop rule. With no item and nothing in progress, the report says "Nothing
   to do" with the sources it checked, and its JSON `state` is
   `nothing_to_do` instead of `work_found`.
6. Status view. `--status` shows the work in progress, the pending Specs in
   delivery order with their Task counts, the next release from the Release
   Plan with its changes or "no unreleased changes", the parked items with
   their Park Class and next action, and what waits on the maintainer: the
   queue's Pending Question and a release that needs approval or a manual
   classification.
7. Configuration. User Config `secondbrain.root` names the Secondbrain Root;
   `secondbrain.inbox_namespace` names the inbox namespace, by default the
   repository's name.
8. Documentation. A new command guide, a new Roundfix skill reference for the
   `follow-up` family, the configuration guide and the glossary describe the
   command, its JSON, its keys and its order.

## User Experience

The command prints plain text by default and one JSON document with
`--format json`. It never prompts, never writes a file, never commits and
never reaches the network except through the judge and `gh`. It exits `0`
whenever it produced a report, `2` for invalid usage, an invalid
configuration, an unreadable messages file or a directory outside a Git
repository, and `1` when the Run Database cannot be read.

## Non-Goals / Out of Scope

- The skill modes `--follow-up`, `--status` and `--unattended`, the questions
  they ask, Triage and the Baseline clause on questions per call: Spec 0307.
- Reading ai-memory directly, or writing to the Secondbrain, ai-memory, the
  repository or the Run Database.
- Fetching or pulling the Secondbrain, or reading another repository's
  checkout.
- Measuring thresholds for the four questions; they start at TypeSafe's
  suggested values.
- Changing `roundfix release plan`, `roundfix deliver status`,
  `roundfix runs list` or `roundfix spec judge`.

## Success Metrics

1. Success Metric: in a temporary repository, home and Secondbrain holding two
   pending and one triaged Inbox Entry of the namespace, one open and one
   promoted Backlog Entry, one pending Finding, one active Spec and a messages
   file with one message, `roundfix follow-up --format json` lists exactly the
   six pending items, the Inbox Entries oldest first, and no triaged or
   promoted record.
2. Success Metric: with a fake judge transport, each item is ordered into the
   class the four answers select, the Judge Log holds one row per item with
   the versioned model id, and an item whose text matches a secret pattern is
   not sent and is ordered by its type.
3. Success Metric: without a key, inside a Run or at the monthly ceiling, no
   request is sent, the items are ordered by type, the report names the
   reason, and the command exits `0`.
4. Success Metric: with no item and nothing in progress the text report
   starts "Nothing to do" and lists the checked sources, the JSON `state` is
   `nothing_to_do`, and one active Run turns it into `work_found`.
5. Success Metric: without `secondbrain.root` the inbox source is `missing`
   and every other item is still listed; a `secondbrain.root` in Project
   Config is ignored with a warning.
6. Success Metric: `--status` lists the queue's Specs in position order before
   the unqueued active Specs, reports "no unreleased changes" when the default
   branch sits on the last release tag and the proposed version with its
   changes otherwise, and lists each park with its Park Class and the Pending
   Question under what waits on the maintainer.
7. Success Metric: neither view changes the repository, the Secondbrain or the
   Run Database: their bytes and `git status` are identical before and after.
8. Success Metric: after release, the operator's first `roundfix follow-up
   --format json` on this repository with the real Secondbrain and a Jev key
   lists the pending Inbox Entries of the `roundfix` namespace and reports
   the judge `answered` with a versioned model id.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- The repository's Secondbrain guide, `docs/agents/secondbrain.md`, written
  by the Baseline: "pending entries live at the destination namespace root
  under `inbox/<destination>/`, and resolved entries live under
  `inbox/<destination>/_triaged/`" and "Triage works pending entries `oldest
  first`, ordered by `created_at`". Success Metric 1 reads the inbox that
  way.
- TypeSafe's documentation of the Noul primitive
  (https://docs.typesafe.ai/primitives/noul): "Ask one yes/no question per
  Noul." and "A good starting point is 0.7 for yes and 0.3 for no.", and of
  the primitives (https://docs.typesafe.ai/primitives): "If the judgment you
  want depends on several independent factors, ask about each factor
  separately and combine the answers with your own logic." These set the four
  questions, the thresholds and the order in code.
- Rex Black's bug triage framework
  (https://rexblack.com/resources/writing/bug-triage-framework): "Fix
  immediately. ... This action is reserved for defects that block other
  important work". The first class puts what blocks another project first.

## Unreachable Acceptance

- criterion: Success Metric 8, after release, the operator's first
  `roundfix follow-up --format json` on this repository with the real
  Secondbrain and a Jev key lists the pending Inbox Entries of the
  `roundfix` namespace and reports the judge `answered` with a versioned
  model id.
  reason: it needs the real Secondbrain, the operator's User Config and a
  request to the judge, and no Verification or QA row may read the real home,
  another repository or reach the network.
  satisfied-by: the operator's report of that first run, with its judge state
  and model id, after this Spec's release.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Code at `42f523c2`: `internal/judge` asks one question per request for a
  fixed set of judgment kinds and logs the model each answer reports;
  `releaseplan.Build` fails with `ErrInvalidReleaseRange` when the target is
  the base release tag; `delivery.ClassifyPark` and
  `delivery.PendingQuestionFor` give each park's class and next action and
  the queue's one question; `deliveryPlanIntent` lists `docs/backlog/`,
  `docs/findings/` and `docs/_inbox/`; no command reads a Secondbrain path and
  no User Config key names one. The TechSpec records the code path.
- Secondbrain, read only. Read `wiki/index.md`,
  `wiki/entities/jev-typesafe.md` (Jev decides, never ranks; "Julgamentos
  amplos que escondem vários fatores: é preciso decompor"),
  `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md` (the question
  tools of the three harnesses, which Spec 0307 uses) and
  `wiki/concepts/porta-de-captura-unica-do-fleet-inbox.md` (positional
  status; only the destination triages), and ran `qmd query` for the
  follow-up inventory, Jev triage questions, structured questions and inbox
  triage order with `--all --files --min-score 0.3`, which returned those
  pages and the Roundfix mirror's usage and command guides. The pages
  confirmed decomposing the judgment and reading the inbox positionally; none
  covers a follow-up command. No Inbox Entry was written, because the
  briefing forbids writing there.
- Exa: the TypeSafe Noul and primitives pages above; Rex Black's triage
  framework and the GitLab handbook's bug prioritization (oldest as a tie
  breaker, https://handbook.gitlab.com/handbook/engineering/infrastructure-platforms/developer-experience/bug-prioritization/);
  and CLI references on read-only status commands that exit `0` in every
  state and carry the state in JSON (https://sow.pgsty.com/docs/reference/exit-codes/:
  "`sow status` ... exits `0` in every repository state ... so a script can
  read the structured state instead of decoding an exit code"). Challenging:
  TypeSafe warns that a value between 0.3 and 0.7 is unclear, so an unsure
  answer must not promote an item; ADR-0319 counts it as not yes.

## Glossary

- adds: **Follow-up Command**
- adds: **Follow-up Item**
- adds: **Follow-up Priority**
- adds: **Secondbrain Root**

## Skills

- unchanged: config keys — the two new `secondbrain` keys are described by the new `follow-up` reference, which task_01 adds to this surface's covering files; no existing reference describes them.
- unchanged: guide docs/user-guide/configuration.md — the new section describes only the `secondbrain` keys, covered by the new `follow-up` reference that task_01 adds to this guide's covering files.
- unchanged: guide docs/user-guide/commands.md — the only change is the `follow-up` index row, covered by the new `follow-up` reference that task_01 adds to this guide's covering files.
- unchanged: guide docs/user-guide/usage.md — the only change is the `follow-up` command table row, covered by the new `follow-up` reference that task_01 adds to this guide's covering files.

## Decisions

- 2026-10-22: Where do the modes live (AskUserQuestion, about 15:20 BRT) —
  "Skill roundfix, fontes opcionais (Recomendado)" (option text: "Os modos
  entram na skill do produto. Cada fonte é detectada e opcional: Secondbrain
  pelo guia do baseline, ai-memory se o MCP existir, Jev se houver chave,
  bloco do GT3 só se as instruções do usuário o definirem. Sem a fonte, o
  modo segue e diz o que faltou.").
- 2026-10-22: Code or skill (AskUserQuestion, about 15:20 BRT) — "Comando CLI
  + skill (Recomendado)" (option text: "Novo comando só leitura (ex.:
  roundfix follow-up --format json) junta backlog, findings, specs, Runs,
  fila, PRs, git e o inbox do Secondbrain, pergunta ao Jev pelo
  internal/judge (model id no log) e ordena em código. A skill junta as
  mensagens do ai-memory (só o agente vê o MCP), conversa e pede as
  autorizações. Testável e igual nos três harnesses."). See ADR-0319.
- 2026-10-22: Wait for a reply from `default/jarvis` — "Sim, suas decisões
  bastam (Recomendado)".

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 2 briefing for Spec 0306 of 2026-10-22 and the proposal adopted under Sources.

Taken from that briefing and the proposal, with no further interview:

- The command is `roundfix follow-up`, with `--status`, `--messages <path>`
  and `--format <text|json>`; the messages file is the only way ai-memory
  messages arrive.
- `secondbrain.root` is User Config only, because it is a path on one
  machine; `secondbrain.inbox_namespace` may be set in either config, because
  it belongs to the repository. Without either, the namespace is the
  repository's name from its `origin` remote, or its main checkout's
  directory name.
- The four questions are asked together, one request per item, as Nouls; 0.7
  is yes, 0.3 is no, and between is unsure (ADR-0319).
- Open Pull Requests come from `gh pr list`, an optional source that reports
  why it is unavailable, because no existing reader lists all open Pull
  Requests.
- The JSON `state` carries the stop rule and the exit code stays `0`, as for
  every advisory report.

## Open Questions

None.
