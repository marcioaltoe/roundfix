---
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
created: 2026-11-16
status: promoted
leverage: compounding
---

# Each phase of a Spec is recorded and starts in a fresh session

## Overview

Roundfix records which Spec and which phase each Supervisor session is working on, outside a Run as well as inside one, in a single documented place that a usage tool can read. The skills and the autonomous-work guide tell the Supervisor to end its session at each phase boundary with a Session Handoff, and Roundfix says so when one session goes on into a second phase. This repository starts running its own Authoring Phases on declared models through the generated subagents, and a fix or a refactor stays in the session unless the user asks otherwise.

It is for the person who supervises Spec work from an Agent harness sold with a usage limit, in any adopting repository. The first reader of the record is the operator's usage script, which already reads the real limit percentage from the harness.

It is valuable because the limit is spent where nothing is recorded today. On the one host measured, 94.2% of a week's Claude usage ran in sessions opened in a repository checkout, where nothing says which Spec or phase was being worked on, and three quarters of what a main session cost was the re-reading of its own context.

V1 is one Spec and stays advisory: Roundfix records and tells. It never ends a session, refuses a command or reports a limit percentage. It succeeds only if the limit falls with no loss of speed and with better plans for implementations, the two conditions the maintainer added on 2026-11-16.

## Problem

Reported by jarvis, message 01a1da5d, 2026-11-14: one of the maintainer's two subscriptions reached 100% of its weekly Claude limit between 09/11 and 10/11, and the other stood at 90% on 14/11, one day before its reset. The message places the weight in long orchestration sessions that resend a large context on every call, and asks Roundfix to write the Run and the phase in a fixed place so that the usage script can attribute the limit by phase.

**The record the request asked for would see half a percent.** The usage script names a Run from the worktree directory a call ran in. Measured by the author on the Mac mini from 2026-11-08 to 2026-11-15, 0.5% of the Claude cost ran in a Roundfix worktree and 94.2% in sessions opened in a repository checkout: the Supervisor and the subagents it starts. Runs spend the other subscription: 98.1% of the Codex tokens already carry a Run id in the directory name, and 2.5% of those also carry the Task.

**A session pays for its whole context on every call.** In the main sessions of that host, 75% of the cost was cache read, 59.2% of the calls carried 200 thousand tokens or more, and two sessions of another repository of the fleet ran for 25 hours each and reached 861 and 961 thousand tokens. The `--follow-up` Skill Mode authors, queues and follows a Spec to its release in one session, and the scheduled wake-up the autonomous-work guide offers returns to the same session with all of its context. The maintainer's rule of one Supervisor session per Spec lives in his personal instructions: no adopter receives it, and nothing shows when it is not followed.

**In this repository the cost is in the authors.** Here the subagents held 74.3% of the week's Claude cost, a subagent call carried a median of 272 thousand tokens, and Sonnet answered for 0.5% of the cost: the authors ran their own sweeps on Opus and Fable. Spec 0392 shipped the declaration that puts each Authoring Phase on its model and sends sweeps to a Research Subagent, and ADR-0401 left this repository out "until a later decision says how", because the skills here are the source of the shipped bundle. Once a phase is declared, its skill hands every TechSpec the agent starts to the declared author, a small fix included, while the maintainer keeps a small fix on the session's model.

### Market Data

- Claude Code's documentation names the cause and prices the remedies: "a one-line question in a session that has been open all day still draws usage for the whole conversation"; "compacting a large context is itself a large request. When you want a fresh start instead of continuity, `/clear` costs nothing"; and a scheduled task fires "even while the session is idle, sending your full context each time".
- Anthropic's guidance on long tasks supports refusing automatic compaction: "overly aggressive compaction can result in the loss of subtle but critical context whose importance only becomes apparent later". It recommends notes kept outside the context and subagents that return "a condensed, distilled summary".
- A fresh context for each phase, with a written artifact carried between phases, is established practice: HumanLayer keeps "utilization in the 40%-60% range", and Cline proposes a new task with a carried summary at a threshold the user sets and starts it when the user approves. Cline's reason is the quality of the work: "after exceeding 50% of the context window, performance can dip".
- The practice is contested. Amp replaced compaction with a handoff in October 2025 and removed the handoff in May 2026: "we don't think it warrants the complexity introduced by many small, connected threads". Its reasons are what current models can do and the product's complexity, not the cost of a plan limit.
- Claude Code's telemetry attributes usage by labels set where the session starts, `session.id` and the user's own attributes. The phase record is the same idea for a Spec and its phase.
- No tool was found that records the phase of Spec work so that plan-limit usage is attributed to it. This is an absence in the search, not proof.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The phase record outside a Run | Critical | Each Supervisor session's work on a Spec is recorded with its harness session, the Spec and the phase: follow-up, one of the four Authoring Phases, delivery or release. |
| F2  | The Run's current Work Item | High | Inside a Run the same record names the Run and the Task or QA gate in progress, with when each began and ended. It also names the pre-PR review of a Delivery Queue item, which runs outside a Run. |
| F3  | One documented place and format | Critical | The record lives in one place that does not depend on a repository's configuration, in a versioned format the user guide documents. Roundfix reads no harness transcript and computes no price or limit percentage. |
| F4  | The session boundary rule | Critical | The Roundfix skill's modes, the authoring skills and the autonomous-work guide tell the Supervisor, at each phase boundary, to write the Session Handoff with the Spec, the phase that ended and the next one, to say what stays armed, and to end the session. The next session starts with `/roundfix follow-up`, which claims the handoff. |
| F5  | The boundary notice | High | Roundfix tells the user, and refuses nothing, when one harness session records a second phase, of the same Spec or of another. |
| F6  | This repository declares its authoring models | High | Roundfix's own Project Config declares the model and effort of each Authoring Phase and of research. Installation generates the Author Subagents and the Research Subagent here and writes no Declared Skill Header into this repository's skills, so a skill the user types runs on the session's model. |
| F7  | A fix or a refactor stays in the session | High | An authoring skill the agent activated on the fix or refactor route writes in the session and hands the phase to its Author Subagent only when the user asks; a feature's phases are always handed over. The rule ships in the skills, so every adopter receives it. |

## KPIs

The maintainer's success criterion has three parts: less of the limit spent, no loss of speed and better artifacts, in his words "com queda de performance nula ou quase nula e com aumento de qualidade de artefatos para implementações". Each part has its own measures, and the two measurements this Spec adopts hold the baselines.

### The limit

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Claude cost of this repository's Supervisor work that carries a Spec and a phase | >= 95% (0% today) | The operator's usage report joined to the record, over the first full week after the release. |
| Median context per call in the Supervisor sessions of this repository | < 150 thousand tokens (198 thousand measured) | The usage script's collected calls of main sessions: input plus cache read plus cache write. |
| Median context per call in the authoring subagents of this repository | < 200 thousand tokens (272 thousand) | The same collection, for subagents named for an Authoring Phase or for research. |
| Authoring cost carried by the generated subagents | >= 90% of the subagent cost of feature Specs (0%: every subagent is `general-purpose` today) | The agent type of each subagent transcript, which the usage script already reads. |

### No loss of speed

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Time from a plan reaching the default branch to its delivery merge | median <= 7.5 hours over the first 20 Specs delivered after the release (6.0 hours over the 58 Specs created since 2026-11-01) | Commit times in the Git history: the commit that adds the Spec's PRD and the commit that adds its Archive Record. |
| Duration of a Clean Implement Run | median <= 71 minutes over the same 20 Specs (61 minutes over 47 Runs since 2026-11-01) | `roundfix runs list`. |
| Time of the strict check with Verification on one fixed plan | the release with this Spec within 5% of the release before it, median of three runs each | The same command on the same tree with the two binaries. |

Each observed target is its baseline plus one standard error of a 20-Spec median, the smallest loss that can be told from chance; the paired check is where "almost none" is tested closely. The work a phase boundary adds is counted, not timed: one Session Handoff written, one session opened and one `/roundfix follow-up`, with no new question to the user. The wait between a phase's end and the next phase's start has no baseline today; the record gives it from its first week.

### Better artifacts for implementations

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Specs whose first Implement Run ends Clean | >= 80% of the first 20 Specs delivered after the release (68%: 32 of 47 since 2026-11-01); under 68% is a loss | `roundfix runs list`: the state of each Spec's earliest Run. |
| Implement Runs per delivered Spec | <= 1.25 over the same 20 Specs (1.38: 65 Runs for 47 Specs) | The same listing. |
| Task verdicts that fail at settlement | < 5% of a week's verdicts (6.5%: 3 of 46 in the week from 2026-11-09) | `roundfix runs causes` over the week, which also names each failure's class. |

A plan is good when an implementation runs from it without stopping: its Verification is honest, its Tasks are complete and an Agent can finish them. A first Run that ends Clean is that outcome, measured where a plan defect, a failed QA gate and a failed Task all show. The expected cause of a rise is that each Authoring Phase starts from a brief and its own reads, on the model declared for it, instead of a context of hundreds of thousands of tokens. Not chosen: corrective Tasks (none in the 30 days of Runs, so nothing to improve), the QA disposition (56 of 58 records hold the same value), strict-check findings at plan time (corrected before the commit and never recorded) and `unproven` entries (declared limits, not defects).

### How the week before and the week after are read

The week before is 2026-11-08 to 2026-11-15: 30 Specs delivered, 467.8 USD-eq of Claude cost for this repository on the Mac mini, and 13.2 points of one account's weekly limit (reported by jarvis). The week after is the first full week after the release that carries this Spec.

1. The limit: the points attributed to this repository on both accounts, divided by the Specs delivered in the week. The limit is account-wide, so the split by project and by phase is the usage script's estimate.
2. The context medians of the first table.
3. Speed and quality: a week holds too few Specs to decide, so the week's numbers only warn when one leaves its band, and the verdict is read over the first 20 Specs delivered. The failed-verdict share is read weekly, because the journal keeps one week.

The Spec succeeds when the first falls, no speed measure leaves its band and no quality measure falls under its baseline, with the quality targets met over 20 Specs. If the limit falls and a speed measure leaves its band, the boundary rule is relaxed first, because it is the part that adds waits. The two weeks differ in more than this Spec: the orchestration model changed on 2026-11-14, and the first-Run rate was already rising, from 60% to 79% across the baseline. Each measure therefore names its own cause, and the share of sessions with a boundary notice says whether the rule was followed at all.

## Feature Assessment

| Criteria        | Question                                            | Score  |
| --------------- | --------------------------------------------------- | ------ |
| Impact          | How much more valuable does this make the product?  | Strong: a cost nobody can see becomes a number per phase, and the rule acts where 38% of a main session's cost is context above 150 thousand tokens. |
| Reach           | What % of users would this affect?                  | Strong: every adopter supervised from a harness with a usage limit receives the rule at its next Baseline update and the record with the binary; F6 reaches this repository only. |
| Frequency       | How often would users encounter this value?         | Must do: every Spec crosses three to five phase boundaries, and every call of a Supervisor session pays for its context. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe: a fresh context per phase is established practice; attributing plan-limit usage to the phase of Spec work was not found in the tools surveyed. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong: the record accumulates, and each later decision is measured against it. |
| Feasibility     | Can we actually build this?                         | Strong: a command the session runs can read the harness session's identity (measured on Claude Code 2.1.325), the Daemon knows the Work Item in progress, and the Session Handoff exists. Roundfix cannot end a session, and the usage script lives outside this repository. |

Verdict: proceed in the shape the maintainer's answers gave. Three results would change it: a first week of the record showing that the Supervisor session gains little context from a phase it handed to an Author Subagent, which would call for fewer boundaries; authors whose context stays near 272 thousand tokens after adoption, which would put the cost in what a phase must read; and a speed measure outside its band or a first-Run rate under 68%, which would relax the boundary rule before anything else.

## Trade-off Insights

- **Recommended approach:** record first and tell, then let the numbers set the rule's strength. The record covers the Supervisor because that is where the Claude limit is spent; the boundary is an instruction with a notice because Roundfix cannot end a session; this repository adopts the generated subagents because its authors hold three quarters of its cost.
- **Key trade-offs:** a notice against a refusal (a refusal would stop a small fix that is rightly authored and delivered in one session); a boundary at every phase against the cost of starting a session (six of seven new sessions here started at 82 to 90 thousand tokens, and by the usage script's weights a boundary pays for itself after about 12 calls, an inference; a phase of hundreds of calls gains and a short wrap-up loses); exposing the record against reporting in Roundfix (a report would read harness transcripts, which hold other projects' and clients' text); generated subagents alone against the whole declaration here (a typed skill runs on the session's model); the fix-route rule for every adopter against a setting (an adopter who wants every TechSpec on the declared model now has to ask).
- **Risks identified:** an unattended cycle stops at a boundary, so the next phase, a release included, waits for the session the maintainer opens; the Session Handoff is the only carrier between sessions, and a poor one loses decisions; the session identity was measured on one harness and one version, and a subagent shows its parent's identity; the notice can become noise on small fixes; the report by phase exists only after the usage script's owner adapts it; the boundary rule acts on a quarter of this repository's cost, and only F6 reaches the rest; the record names Specs, sessions and times on the user's disk, so its requirement must say that it holds no prompt text and no secret and who can read it.
- **Dissenting view:** Amp's reversal says that many small connected sessions are a cost of their own. The author's dissent: once a phase runs in an Author Subagent, the Supervisor session receives only the brief it wrote and the report that came back, so ending it after every Authoring Phase may spend more on session starts and on the maintainer's attention than it saves. The scope as decided keeps the boundary at every phase, and the record will settle it. A third: the request's own shape, a record written only by a Run, was the smaller Spec and would have measured half a percent of the problem.
- **Stretch goal (V2+):** a command the repository declares, started at a boundary to open the next session with nobody present; and Authoring Phases that Roundfix starts as Agent Sessions of their own, each on its declared model with its usage counted as Turn Usage.

## Opportunity Scan

- **Original** — the three items decided on 2026-11-14 in the shape the six answers gave. Medium effort, one Spec. Impact Strong, Reach Strong, Frequency Must do, Differentiation Maybe, Defensibility Strong, Feasibility Strong.
- **More ambitious** — authoring becomes work Roundfix starts: each Authoring Phase runs as an Agent Session on its declared model, in a fresh context, with its usage counted by phase. Massive effort. Impact Must do, Feasibility Maybe: the Idea stage is a conversation with the user, and a Run would stop targeting only a Spec's Tasks.
- **Simpler** — measure first: the record and the adoption, with the boundary left as the maintainer's own rule until a week of numbers by phase exists. Small effort. Impact Maybe: it leaves out the adopter whose sessions reached 861 and 961 thousand tokens.
- **Adjacent** — what a session carries when it starts: six of the seven sessions of this repository in the week measured began at 82 to 90 thousand tokens, and each Authoring Phase reads about 52 thousand tokens of glossary before anything about its Spec. Medium effort. Impact Strong, because every call pays it and the boundary rule adds session starts; Feasibility Maybe, because part of that floor is the user's own configuration.

**Chosen direction:** original — the maintainer's pick ("Original + respostas"). The adjacent alternative is recorded as the Backlog Entry [`2026-11-16-a-new-session-and-each-authoring-phase-start-with-a-large-context.md`](../../backlog/2026-11-16-a-new-session-and-each-authoring-phase-start-with-a-large-context.md), and the more ambitious one is the stretch goal.

## Integration with Existing Features

| Integration point | How |
| ----------------- | --- |
| The Run Marker and Turn Usage | The record's Run entries name the Run the marker names. Turn Usage is unchanged: the record counts no tokens; it says what was in progress and when. |
| The Session Handoff and the `--follow-up` Skill Mode | The boundary writes the handoff, and the mode already claims it before its inventory. |
| The autonomous-work guide | Gains the boundary rule beside its rule about arming what outlives a turn; adopters receive it at their next `roundfix baseline update`. |
| The Authoring Model Declaration and its generated subagents | This repository declares it. The generated subagents' names already say the phase in a transcript; the record adds the Spec. The skills' delegation rule gains the fix and refactor route. |
| The Delivery Queue | Goes on delivering without a session once it is armed. The record names its item's phase and its pre-PR review. |
| Spec 0392 | Its Archive Record leaves the operator's pilot evidence unproven; a full authoring cycle here can supply it. |
| The operator's usage script | Reads the record. It is outside this repository, and its owner adapts it. |

## Out of Scope (V1)

- **Automatic or proactive compaction** — refused by the maintainer: "passa para determinístico e pode perder informação importante".
- **Roundfix ending or opening a harness session, and a command that opens the next session unattended** — the harness owns its sessions; the maintainer chose instruction and notice, and the command stays a stretch goal.
- **A refusal when a session crosses a boundary** — the notice informs; a refusal would stop a small fix authored and delivered in one session.
- **Reading harness transcripts, pricing tokens or computing a limit percentage in Roundfix** — the maintainer chose that Roundfix only exposes the record, and Roundfix keeps no price table.
- **Changing the operator's usage script** — it lives outside this repository; this Spec documents what it can read.
- **The Declared Skill Header in this repository's skills** — they are the source of the shipped bundle; the maintainer chose the generated subagents alone.
- **The declaration for OpenCode and Codex** — it is Spec 0394.
- **Reducing what a session or a phase carries when it starts** — recorded as a Backlog Entry; it changes a different component and needs the record to be measured by phase.
- **Ending a session when its context passes a size** — Roundfix cannot see a session's context, and nobody asked for it.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: no (0.26)
- unverified_assumption: no (0.18)
- market_question: no (0.07)
- verdict: not indicated
- stage: ran
- shape: changed — the record covers the Supervisor's sessions outside a Run and not only the Run; this repository adopts the generated subagents without the header; and a fix or a refactor is authored in the session unless the user asks for the Author Subagent.

The stage ran although the verdict was not indicated: the maintainer chose "Manter a Ideia" on 2026-11-15, after choosing "Spec com Ideia após o reset" on 2026-11-14, before the judge existed. The judged text is [`references/2026-11-15-request-read-by-the-idea-stage-judge.md`](references/2026-11-15-request-read-by-the-idea-stage-judge.md).

## Decisions

Every answer is the maintainer's, quoted as recorded in [`references/2026-11-15-maintainer-answers.md`](references/2026-11-15-maintainer-answers.md).

- Route for the message from jarvis: "Medidas já, sem Spec" and "Spec com Ideia após o reset" (2026-11-14, ~11:35 BRT); the Idea stage is kept: "Manter a Ideia" (2026-11-15, ~20:15 BRT).
- Orchestration model: "Orquestração mais barata" = "Opus medium" (2026-11-14); it is the user's configuration, not this Spec.
- Automatic compaction is refused: "passa para determinístico e pode perder informação importante" (2026-11-14).
- Adoption of the per-phase models in this repository: "Adotar na Spec 0397" (2026-11-14). The research subagent's model: "Sonnet, effort medium (Recommended)" (2026-11-15).
- Scope, as the operator recorded it (handoff 01a1dc8f): "o Run grava o id e a fase num lugar fixo para o ai_usage.py; o Supervisor fecha a sessão em cada fronteira de fase".
- Authorization: "Mesmas do BE" (push, PR, merge, release após cada Spec, até 3 correções, Makefile/CI/go.mod/lint só se a correção exigir) (2026-11-14); the PRD quotes it in the authorization record.
- Coverage of the phase record (15–16/11/2026, ~00:00 BRT): "Run + sessões Supervisor (Recommended)" — "Registra sessão, Spec e fase também fora do Run (autoria, entrega, release, triagem) e, dentro do Run, a Task ou o QA em curso."
- Who reports usage by phase (15–16/11/2026): "Roundfix só expõe (Recommended)" — "Lugar e formato documentados e estáveis; o ai_usage.py cruza com os transcripts; o Roundfix não lê transcript de harness nem calcula preço."
- The Supervisor's session boundary (15–16/11/2026): "Instrução e aviso (Recommended)" — "Skills e guia mandam gravar handoff, dizer o que ficou armado e encerrar; o Roundfix avisa, sem bloquear, quando a mesma sessão registra duas fases; você abre a próxima com /roundfix follow-up."
- Adoption of the per-phase models in Roundfix (15–16/11/2026): "Só subagentes gerados (Recommended)" — "Na 0397: autores por fase e pesquisador em Sonnet em .claude/agents/, sem cabeçalho nas skills deste repositório; uma skill digitada segue no modelo da sessão."
- A small fix against the per-phase declaration (15–16/11/2026): "Exceção pela rota (Recommended)" — "Na rota de correção ou refactor a skill escreve na sessão e só delega quando você pedir; feature sempre delega; muda o texto das skills para todos os adotantes."
- Direction of the Idea (15–16/11/2026): "Original + respostas (Recommended)" — "Registro de fase, fronteira de sessão e adoção dos modelos numa Spec; o custo de nascer de cada sessão vira Backlog Entry."
- Addition to the success criterion (2026-11-16, in the session): "Adicione ao critério de sucesso: com queda de performance nula ou quase nula e com aumento de qualidade de artefatos para implementações."
- The declared models are the ones the request names: Fable at xhigh effort for the Idea, the PRD and the TechSpec, and Fable at high effort for the Tasks.

## Open Questions

- Which transitions are a boundary for the notice? Default: a change of Authoring Phase, the change from authoring to delivery, and a change of Spec. Delivery to release of the same Spec is not one, because a boundary there would leave a merged Spec unreleased until someone opens the next session. Neither is follow-up to the first Authoring Phase it plans.
- Does a phase handed to an Author Subagent count as a boundary for the Supervisor session? Default: yes, as the scope was decided; the PRD states the measured number that would relax the rule.
- How does Roundfix learn the Supervisor's phase? Default: the TechSpec decides between the commands the Supervisor already runs and a step the skill takes when a phase starts; a phase in which no Roundfix command runs, such as an Idea stage, must still be recorded.
- What does the record hold when the harness gives a command no session identity? Default: the entry is written without one and the notice is skipped; Codex and OpenCode were not measured.
- Where does the notice appear, and how long is the record kept? Default: on the command that records the second phase and in the `--status` report; kept at least as long as a Run's usage rows.
- How does installation know that this repository's skills are the source of the bundle? Default: the PRD decides and records it in an ADR that amends ADR-0401, which left the question open.
- How does the user ask for the Author Subagent on the fix or refactor route when nobody is present? Default: the `--unattended` Skill Mode gathers it with the other decisions, and without an answer the fix is authored in the session.
- Which web and MCP tools does this repository's Research Subagent list? Default: the PRD lists them; a sweep that needs a command, such as the Secondbrain query, stays with the author.
- Which names does the glossary take? Default: the PRD's Glossary settles them; the candidates are a Supervisor Phase (follow-up, the four Authoring Phases, delivery and release), the Phase Record and the Boundary Notice.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-14-jarvis-long-orchestration-consumes-the-weekly-claude-limit.md`](references/2026-11-14-jarvis-long-orchestration-consumes-the-weekly-claude-limit.md), reported by jarvis, message 01a1da5d, 2026-11-14. The sender asked for a reply with the decision, and the operator replied on 2026-11-14 (01a1da82). No completion notice was requested.
2. [`references/2026-11-15-maintainer-answers.md`](references/2026-11-15-maintainer-answers.md), the answers behind `## Decisions`.
3. [`references/2026-11-15-request-read-by-the-idea-stage-judge.md`](references/2026-11-15-request-read-by-the-idea-stage-judge.md), the request as the judge read it.
4. [`references/2026-11-15-usage-by-context-and-session-identity-measurement.md`](references/2026-11-15-usage-by-context-and-session-identity-measurement.md), recorded by the author: the usage numbers of this file that name no other source.
5. [`references/2026-11-16-delivery-speed-and-plan-quality-baseline.md`](references/2026-11-16-delivery-speed-and-plan-quality-baseline.md), recorded by the author: the baselines of speed and of plan quality in `## KPIs`.

This project's ai-memory messages were listed on 2026-11-15 and none was consumed: the only pending one is the message from jarvis adopted above. The session listed no Research Subagent, so the author did the sweeps.

Secondbrain (read on 2026-11-15; `wiki/index.md`, then `qmd query`):

- `wiki/concepts/compactacao-de-contexto-de-agentes.md` — calls the restart with a handoff the way that depends least on the harness and says a poor handoff loses everything, which is why F4 states what the handoff carries.
- `wiki/concepts/claude-code-operacao-produtiva.md` — keeps work inside about 150 thousand tokens and clears at phase boundaries; it gave the second KPI its number.
- `wiki/concepts/custos-e-limites-de-agentes-de-codigo.md` — subagents inherit the parent's model and multiply cost unseen; it supports declared models for the authors.
- `wiki/concepts/memoria-de-sessoes-para-agentes-de-codigo.md` — clear, compact, handoff and subagent as boundary decisions.

External sources (read through Exa on 2026-11-15):

- Claude Code documentation, "Manage costs effectively", https://code.claude.com/docs/en/costs — why a long session costs more and what a scheduled task resends; it changed the reading of the guide's scheduled wake-up.
- Claude Code documentation, "Monitoring", https://code.claude.com/docs/en/monitoring-usage/ — usage attributed by `session.id` and by labels set at the source.
- Anthropic, "Effective context engineering for AI agents", https://www.anthropic.com/research/effective-context-engineering-for-ai-agents — compaction's loss, notes outside the context, and subagents that return a summary.
- Anthropic, "Effective harnesses for long-running agents", https://www.anthropic.com/research/effective-harnesses-for-long-running-agents — a new session has no memory of the last, bridged by a progress record and the Git history.
- HumanLayer, "Advanced Context Engineering for Coding Agents", https://www.humanlayer.com/blog/advanced-context-engineering — a fresh window per phase with a written artifact.
- Cline, "Unlocking Persistent Memory", https://cline.bot/blog/unlocking-persistent-memory-how-clines-new_task-tool-eliminates-context-window-limitations — a handoff proposed at a threshold and started on the user's approval.
- Amp, "Handoff (No More Compaction)", https://ampcode.com/news/handoff, and "Amp, Rebuilt", https://ampcode.com/news/neo — the handoff adopted and then removed; the dissent in Trade-off Insights.
