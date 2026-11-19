---
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
created: 2026-11-19
status: promoted
leverage: quick-win
---

# A project declares whether an open authoring question waits or follows the author

## Overview

When the author of a PRD, a TechSpec or a Task Graph ends its phase with a question that no source settles, and nobody is there to answer, a project now says what happens. By default the Spec stops and waits for the user. A project can declare instead that the session follows the option the author recommends and records it as a decision taken by delegation. This Idea calls such a question an open authoring question.

It is for the person who answers a repository's questions, and for the session that authors Specs for that person while the person is away. In this repository and in the projects that adopt the Baseline, that person is the maintainer.

It is valuable because authoring with nobody present is the usual case here, and one unanswered question holds a whole Spec until the user returns. In the answer records of this cycle the maintainer took the recommended option in 29 of 33 answers ([measurement](references/2026-11-19-how-an-open-authoring-question-is-handled-today.md)). The questions that must wait are known and few: data sent outside, new scope, spend, every authorization, the Idea stage and a question about a rule.

V1 is one preference in Project Config, one question before the flow starts, a record in fixed form and a line in the status view. The effort is small, because the question contract, the Clarification Record, the Waiting Declaration and the authorization record already exist.

## Problem

The maintainer's request of 2026-11-19, from [his answers](references/2026-11-19-maintainer-answers.md): "Antes de iniciar o fluxo, perguntar se pode seguir o padrão de autorização decidido pelo autor ou se deve esperar; seria bacana no projeto ter uma configuração no .roundfixrc.yml com a preferencia por projeto e o default ser a spec para e espera o autor." He wrote it in answer to the operator's question of what to do when an author ends a phase with an open question while he is away. The word "autor" has two readings in that sentence. The operator's next question fixed them: the Spec waits for the person, and "o padrão do autor" is the option the Author Subagent recommends.

The instructions leave that case without a route. Line numbers are of this repository at `348962d2da564ab47970e4bd41c16dff0255ebfc`.

- An Author Subagent that lacks a decision returns its questions, "each with its options and the one you recommend" (`.agents/skills/write-prd/SKILL.md:89`, and the same sentence in `write-idea`, `write-techspec` and `write-tasks`).
- The session that started it must "ask the user and continue it with the answers" (`.agents/skills/write-prd/SKILL.md:91`). No sentence says what the session does when the user is away.
- The session may not take the recommended option: "A recommendation, silence, a dismissal or a timeout is never an answer" (`.agents/skills/roundfix/references/modes.md:44`).
- With nobody present, the PRD author may "settle questions from the available sources" (`.agents/skills/write-prd/SKILL.md:148`). A question that no source settles is outside that route.

So the Spec stops, and nothing records that it stopped or why. Authoring with nobody present is not rare: 41 of the 43 PRDs of Specs 0364 to 0408 name the autonomous route, and 29 hold no answer of a person ([measurement](references/2026-11-19-how-an-open-authoring-question-is-handled-today.md)). The same measurement cannot say how many of those phases met such a question, because no file records one.

The workaround is a one-time answer. On 2026-11-19 the maintainer told the operator, for Specs 0409 and 0410 only: "Seguir o padrão do autor (Recommended)", with the decision recorded in the PRD "como decisão do operador por delegação, com pergunta e opções", and his limits still stopping the Spec. That answer lives in one session's notes. The next cycle and every adopter start from the silence again, and a decision taken that way has no form that tells it apart from an answer he gave.

### Market Data

- Published tools wait, and treat silence as a refusal. LangGraph "waits indefinitely until you resume execution". AG2 says "an unanswered approval is not an approval". The web sweep found no tool that records a decision as taken on someone's behalf in a fixed form.
- The opt-in matches a named level of delegation. Appelo's sixth of seven levels, Inquire: "You let the team decide. And afterwards you inquire about their motivations, or you ask that they actively keep you informed." The seventh, Delegate, is the one this Idea does not take: "you don't even need to know which decisions they make."
- Autonomy is a setting, not a property of the Agent. Feng, McDonald and Zhang define "five levels of escalating agent autonomy, characterized by the roles a user can take when interacting with an agent: operator, collaborator, consultant, approver, and observer". A project that declares the preference moves its authoring one level, for one kind of question.
- A mistake nobody catches costs more with each step built on it. In chains of 20 code requests, agents "complete 95% of requests taken one at a time, but when requests are chained together, they complete 51%", and "A single external ``that's wrong'' after a failed request and one repair round lift chain completion from 51% to 59%" (Cuadron and Gupta, 2026). The study measures code, not plans. This Idea infers from it that a delegated decision must be easy to see and easy to overturn before the next phase builds on it.

## Core Features

A recommended option that the session adopts is called a Delegated Decision in this Idea. The name is a candidate for the glossary, which the PRD settles.

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The preference | Critical | A project declares in its Project Config what happens to an open authoring question when nobody is present: the Spec waits, or the session follows the author's recommendation. A project that declares nothing waits. Only Project Config holds the preference, and only as the project committed it. |
| F2  | The question before the flow | Critical | When the project declares nothing, the `follow-up` and `unattended` Skill Modes ask once, while the user is present and before any Spec is authored: wait, which is the recommended option, or follow the author. The answer holds for the Specs of that flow and goes into each one's authorization record. Without an answer the Specs wait. A project that declares is not asked, and the mode's report says which preference is in force and where it came from. |
| F3  | What is delegated | Critical | With nobody present, a question that the author of a PRD, a TechSpec or a Task Graph returns with two or three options and a recommended one, that no source settles and that F4 does not cover. The session adopts the recommended option and the phase goes on. With the user present the session asks, as today. |
| F4  | What always waits | Critical | Whatever the preference: the Idea stage, with its questions and the approval of its draft; the choice an Idea-stage judgment leaves to the user; a question about a rule, which Spec 0409 defines; every authorization; and the user's limits, which are data sent outside the machine, new scope and spend. A question without a recommended option waits too. |
| F5  | Two readers classify | High | The author says, for each question it returns, whether F4 covers it and why. The session checks that statement against F4 before it adopts anything. When either one says the question is the user's, or cannot tell, the Spec waits. |
| F6  | The record | Critical | Each Delegated Decision is recorded in the PRD's `## Decisions` in a fixed form that cannot be read as a person's answer: the date, the phase that asked, the question, the options, the adopted option with the author's reason, and the grant it rests on. A question of the TechSpec or the Tasks phase is recorded in the PRD, and that phase's artifact points at the record. |
| F7  | The report and the status view | High | The session's report lists every Delegated Decision of the flow. The status view of the Follow-up Command lists them for each active Spec. |
| F8  | The Spec goes to the end, and a later answer wins | High | A Spec with Delegated Decisions continues through authoring and delivery inside the authorization it already holds. When the user later answers a delegated question, the dated answer joins the record. A different answer amends the plan while the Spec's Run has not started, and after that it follows the routes that exist. |
| F9  | The wait | High | When the Spec waits, only that Spec waits: no artifact of the phase is written, a Backlog Entry that names the Spec carries the whole question with its options under a Waiting Declaration, the queue goes on, and the Follow-up Command shows the wait. |
| F10 | The rule reaches every session | Medium | A Normative Clause of the Baseline's Spec workflow states F3, F4 and F6 in stack-neutral words. The `modes` reference of the Roundfix skill and the `write-prd`, `write-techspec` and `write-tasks` skills carry the steps, the configuration guide describes the preference, and `CONTEXT.md` gains the names. |

## A Delegated Decision in the record

An illustration with an invented question. The PRD stage settles the exact form.

```markdown
Delegated decision, 2026-11-20, asked by the PRD author:

- question: which failure behavior does an expired import have?
- options: retry automatically up to 3 times; fail fast and notify the user; park it for manual review
- adopted: retry automatically up to 3 times — the author's recommendation, "it preserves the sync retry behavior"
- under: the preference this project declares in Project Config
```

The `under` line names either the project's declared preference or the user's answer before the flow, with its date. When the user answers the question later, a dated line of the Clarification Record follows the entry, in the form that record already has for a person's answer. The entry itself stays, so the file shows both what the session adopted and what the user said.

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Specs of a delegating project that stop on a question F3 covers | 0, where today every such question stops its Spec until the user answers | In a repository whose preference is to follow the author, count at each cycle close the Backlog Entries whose Waiting Declaration names an authoring question that F4 does not cover. |
| Delegated Decisions with a complete record | 100% of the first 20, and 0 written as a dated answer of a person or under the autonomous-route line | Read `## Decisions` in the PRDs of the first 20 Specs that hold one, and run the Spec Consistency Check over them once the PRD stage gives it the form. |
| Decisions that F4 covered and were adopted anyway | 0 | At each cycle close the maintainer reads every record of the cycle. The record holds the question and its options, so one read is enough. Each hit is a defect of the Spec that adopted it. |
| Delegated Decisions the user later changes | Under 20% over the first 20 | Compare each entry with the dated answer that follows it. The only existing count is over questions of other kinds: 29 of 33 answers of this cycle took the option labeled recommended. |
| Questions added at the start of a flow | At most 1 in a project that declares nothing, and 0 in a project that declares | Count the question in the authorization records of the Specs each flow authors. |

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Strong — it removes a stop that today lasts until the user returns, and it gives a form to a decision that this cycle recorded in session notes. It adds no capability to a project that keeps the default. |
| Reach           | What % of users would this affect?                  | Strong — every project that authors with nobody present: 41 of 43 PRDs here. A project whose user is always present gains only the record of a wait. |
| Frequency       | How often would users encounter this value?         | Maybe — nobody has counted how often an author ends a phase with a question no source settles. The Idea of this Spec returned six, all about shape, which stay the user's. |
| Differentiation | Does this set us apart or just match competitors?   | Strong — the tools the sweep read wait or refuse. None records a decision taken by delegation, with its options, for a later look. |
| Defensibility   | Is this easy to copy or does it compound over time? | Maybe — a key and a paragraph are easy to copy. What stays in each repository is the record: which recommendations its user later changed, which is the evidence for widening or narrowing the delegation. |
| Feasibility     | Can we actually build this?                         | Must do — the pieces exist: Project Config with keys only a project may set (`internal/config/authoring.go:44`), the Clarification Record with a check of its shape, the Waiting Declaration and the status view. |

The leverage type is quick win. The verdict is to proceed. Two results would change it: a share of changed decisions above the target, which means the delegation is too wide or the recommendations are weak, or no Delegated Decision at all in a project that follows the author, which means the questions that stop its Specs are the ones F4 keeps for the user.

## Trade-off Insights

- **Recommended approach:** the maintainer's four answers. The default keeps today's rule and gives the wait a record. The opt-in is declared per project, in a file the team reviews, and a project that has not declared is asked once per flow. Every adopted option leaves a record that names itself as delegated.
- **Key trade-offs:**
  - A Spec that waits against a Spec built on an answer nobody gave. Waiting costs the time until the user returns. Following costs a wrong option that the TechSpec, the Tasks and the code then build on. The default takes the first cost, and the opt-in takes the second with a record and a status line.
  - A wide list of what always waits against a narrow one. A wide list makes the preference fire rarely. A narrow one lets a question of scope pass as a detail. V1 takes the user's three limits, every authorization, the Idea stage and the question about a rule, and makes doubt wait.
  - To the end against a hold before delivery. The maintainer chose to the end. A hold would give one more chance to catch a wrong option and would also turn every delegating Spec into a waiting Spec one step later.
  - A preference of the project against one of each user. A file in the repository is reviewed and shared, so one teammate's choice binds the others' unattended sessions. A user who wants to be asked is asked whenever present.
  - The record in the PRD against a file of its own. The PRD is where a Spec's decisions live and where the shape check reads. A question of the TechSpec or the Tasks phase then makes a later phase add to the PRD.
- **Risks identified:**
  - The author that classifies a question is the Agent that goes on when the question passes. F5 adds the session as a second reader, and both are Agents. The third KPI is the check a person makes.
  - "Nobody present" is something the session is told. A session that wrongly believes it would adopt where it could have asked.
  - An Agent could switch the preference by editing Project Config. The file is a Governed Path, the tooling-authority clause forbids the edit, a Task commit leaves the file out unless the authorization bounds it, and F1 reads the preference only as committed. The TechSpec must prove the last one.
  - An author may label a coin flip as recommended. The record keeps the author's reason, and the fourth KPI shows the pattern.
  - "New scope" is a judgment. The PRD must give it a test a second reader can apply.
- **Dissenting view:** published tools hold the opposite line: "an unanswered approval is not an approval". A better answer to the cost of waiting may be to make the wait short, by telling the user at once that a Spec waits and on what, and never to build on an option no person chose. The opposite dissent also stands: with scope, spend, data, authorizations and rules all kept for the user, the questions left to delegate may be so few that the key changes nothing, and nobody has counted them.
- **Stretch goal (V2+):** a preference per phase and per kind of question, as a delegation board sets a level per decision area. A notice to the user when a Spec starts to wait. One question call at the next follow-up that shows the Delegated Decisions since the last one and takes a confirm or a change for each. The Advisory Judgment as a third reader of F4.

## Opportunity Scan

The original idea is bounded by how many open authoring questions fall outside F4, which nobody has counted. If it works fully, a project that follows the author loses no Spec to a question of detail, and its user reads a short list of adopted options in place of answering each one.

- **More ambitious: a delegation board and a judge.** Project Config holds a level per phase and per kind of question, and the judge classifies each returned question and scores the recommendation. Scale: medium to large. Impact Strong, Reach Strong, Frequency Strong, Differentiation Strong, Defensibility Strong, Feasibility Maybe: the judge would receive the text of the question, which is itself data sent outside, and its thresholds need a record of delegated decisions that does not exist yet.
- **Simpler: the question alone.** No key. The Skill Modes ask before every flow and record the answer in each Spec's authorization record, which makes this cycle's one-time answer a rule. Scale: small. Impact Strong, Reach Maybe, Frequency Strong, Differentiation Maybe, Defensibility Pass, Feasibility Must do. It leaves out the standing preference the maintainer asked for, and it asks again in every flow.
- **Adjacent: a short wait.** A Spec that starts to wait tells the user at once, with the question, through the notification the project already configures for Runs. Scale: small to medium. Impact Strong, Reach Strong, Frequency Strong, Differentiation Maybe, Defensibility Maybe, Feasibility Strong. It lowers the cost of the default and delegates nothing, and a notification that carries the question's text sends it outside the machine.

**Chosen direction:** original — the maintainer's four answers of 2026-11-19 describe it. It takes the question of the simpler shape for a project that has not declared. The board, the judge and the notice are the stretch goal: each needs the record V1 creates to know what to tune.

## Integration with Existing Features

| Integration point | What it gains or how it is used |
| --- | --- |
| Project Config (`docs/user-guide/configuration.md`) | One preference that only a project sets. A project without it is read by every binary; a project that sets it needs the release that knows the key, because an older binary refuses a key it does not know (`internal/config/config.go:1910`). |
| The `follow-up` and `unattended` Skill Modes (`.agents/skills/roundfix/references/modes.md`) | The question of F2 joins the decisions the `unattended` mode gathers. The modes say what the session does with a returned question when nobody is present, and what the report lists. The sentence that a recommendation "is never an answer" stays true: a Delegated Decision is not an answer of the user and is never recorded as one. |
| `write-prd`, `write-techspec`, `write-tasks` | Where each says "ask the user and continue it with the answers", it gains the branch for nobody present. The Author Subagent's return gains the statement of F5 for each question. The autonomous route of `write-prd` is unchanged: it covers what a source settles. |
| `write-idea` | Unchanged. The Idea stage stays the user's (`.agents/skills/write-idea/SKILL.md:26`, `:44`). |
| The Clarification Record (`CONTEXT.md:1691`) | A third form beside the dated answer and the autonomous-route line. The Spec Consistency Check, which reads the record's shape, learns it. |
| The authorization record | It holds the user's answer to the question before the flow, as it holds his other answers (`modes.md:38`). It gains no field that grants anything to a later Spec. |
| The Waiting Declaration and the Follow-up Command | Used as they are for the wait, in the way Spec 0409 chose for a question about a rule: a Backlog Entry that names the Spec carries the question. The status view gains the list of F7. |
| The Baseline's Spec workflow | One Normative Clause, beside the one Spec 0409 adds, so that a session of an adopter knows the default, what always waits and the record. |
| The stop-and-ask clause (`docs/agents/agent-instructions.md:90`) | Used as it is. With the user present the session asks through the question tool the harness exposes. |
| The Idea Stage Judgment and the Rule Challenge of Spec 0409 | Unchanged. Both are in F4. |
| `CONTEXT.md` | The names the PRD settles for the adopted option and for the preference, and the changed definition of the Clarification Record. |

## Size and Shared Files

By the sizing criterion of the Spec routing guide this is one Spec. It holds one context: what a session does with an open authoring question when nobody is present. The preference, the question before the flow, the record and the rule complete one another. A key without the rule changes nothing in a session a skill leads, and a rule without the record leaves an adopted option looking like an answer. The work fits three implementation Tasks and the QA gate, each large enough to earn its verification cycle:

1. the preference in Project Config, with its validation, the way a session reads it and the configuration guide;
2. the record: its form in the PRD template and the Clarification Record reference, the shape check, and the list in the status view;
3. the rule: the Normative Clause with the files each Baseline edit regenerates, the `modes` reference and the three authoring skills, and the names in `CONTEXT.md`.

The request had a second part, the standing authorization of the shell commands the flow runs. It is Spec 0413. The two share no component and no contract: this Spec changes how an Authoring Phase treats a question, and that one changes what a harness asks before it runs a command. Together they would need seven or eight Tasks. The maintainer chose "Duas Specs (Recommended)".

This Spec is delivered after Spec 0410 and before Spec 0413. It was authored at `348962d2`, before Specs 0409 and 0410 merged, so the PRD must be written against the tree that holds their text. These files can meet:

- `.agents/skills/roundfix/references/modes.md` and the version fields of the Roundfix skill. Spec 0409 adds a paragraph under `--unattended` and a sentence in step 5, and Specs 0410, 0411 and 0413 also change this skill.
- `.agents/skills/write-prd/SKILL.md`, its PRD template and `.agents/skills/write-techspec/SKILL.md`, which Spec 0409 changes, with `skills/testdata/owned-skill-versions.json` and the `skills/` mirror.
- `internal/baseline/assets/modules/spec-workflow.json`, `docs/agents/spec-routing.md` and the files each Baseline edit regenerates. Specs 0409 and 0410 change the same modules.
- `docs/references/behavior-surfaces.json`, because the status view and the configuration guide change. Specs 0410 and 0413 change it too.
- `CONTEXT.md`.

## Out of Scope (V1)

- **The Idea stage** — its questions and the approval of its draft stay the user's. The maintainer chose "PRD, TechSpec e Tasks (Recommended)", and the Idea is where scope is decided.
- **A question about a rule** — Spec 0409 decided that nobody answers it for the user and that its Spec waits. This Spec changes neither.
- **Any authorization** — a commit, a push, a Pull Request, a merge, a release, a Governed Path, a destructive Git operation and a database statement keep their stop-and-ask clauses. A preference delegates a choice between options, never a permission.
- **Silence as an answer while the user is present** — a question the user can see is asked and waits for him. No timeout adopts anything.
- **A hold before delivery** — the maintainer chose "Até o fim, com relato (Recommended)".
- **A preference in User Config** — the request asks for a preference of the project. A machine-wide switch would follow a person into repositories that never chose it.
- **Questions in a Run, in `implement-task` and in the QA gate** — an Agent of a Run never asks the user, and its rules do not change.
- **What the PRD author settles from sources** — the autonomous route stays as it is.
- **Permission prompts of the harness for shell commands** — Spec 0413.
- **A preference per phase or per kind of question, a judge, a notice of a wait and a confirm step** — the stretch goal.

## Idea-stage judgment

- item: the request in the session, recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md)
- judged-by: nobody (the user asked for the Idea stage)
- verdict: not judged
- stage: ran
- shape: split — the second part of the request, the standing authorization of shell commands, became Spec 0413, and the first part gained the list of what always waits, the two readers, the record in fixed form and the line in the status view.

The judge was not run. The maintainer chose the Idea stage directly on 2026-11-19: "Spec 0412 com Ideia agora". The `shape` line is the author's and waits for the approval of this draft.

## Decisions

Every answer is the maintainer's, as recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The text in quotes is the option he chose, verbatim. The text after the dash is that option's description as the operator recorded it.

- This cycle (2026-11-19, ~00:20 BRT): "Seguir o padrão do autor (Recommended)" — o operador adota a opção recomendada pelo autor e registra no PRD como decisão do operador por delegação, com pergunta e opções; limite do Marcio (dados para fora, escopo novo, gasto) continua parando a Spec.
- How to record the suggestion (2026-11-19, ~00:20 BRT): "Spec 0412 com Ideia agora" — abre a Spec já, com etapa de Ideia, entregue depois da 0410; a Ideia volta com perguntas ao Marcio.
- One Spec or two (2026-11-19, ~01:10 BRT): "Duas Specs (Recommended)" — 0412 fica com (A) preferência para pergunta aberta de autoria; (B) autorização permanente de comandos de shell vira a 0413, com Ideia própria partindo desta pesquisa, entregue depois da 0412.
- The phases the delegation covers (2026-11-19, ~01:10 BRT): "PRD, TechSpec e Tasks (Recommended)" — a Ideia continua do Marcio (perguntas e aprovação do rascunho esperam).
- The question before the flow and the configuration (2026-11-19, ~01:10 BRT): "Config responde; sem ela, pergunta (Recommended)" — projeto que declara não é perguntado; sem a chave, pergunta uma vez antes do fluxo; sem resposta, a Spec espera.
- How far a Spec goes (2026-11-19, ~01:10 BRT): "Até o fim, com relato (Recommended)" — autoria e entrega dentro da autorização existente; cada decisão delegada no PRD (pergunta, opções, escolha, data), no relato e no --status.

The author's own, made in this stage:

- The questions did not offer the dismissal of the stage, because the maintainer had already chosen it. The operator asked them, because the Author Subagent cannot ask.
- A question about a rule is in F4 without a new question to the maintainer. He answered it for Spec 0409 on 2026-11-18: "A Spec espera (Recommended)".
- The wait of F9 uses the mechanism Spec 0409 chose, a Backlog Entry that names the Spec, so that two kinds of waiting Spec read the same way in the Follow-up Command.
- A Delegated Decision is not written as a dated line of the Clarification Record. The glossary tells authors to avoid "assumed approval" (`CONTEXT.md:1693`), and the shape check cannot tell a person's line from one an Agent wrote in the same form.
- The preference lives only in Project Config. The request says "por projeto", and two keys already work that way (`internal/config/authoring.go:44-62`).
- No source of the repository is adopted. The request exists only in the session, so the maintainer's answers are the source. No open Backlog Entry or unresolved Finding shares this context.
- The author read the rules that touch this Idea: the stop-and-ask clauses, the sentence that a recommendation is never an answer, the rule that the Idea stage is the user's and the tooling-authority clause. The request itself asks to change what happens when nobody answers, so that change needs no further question. No other rule blocks the goal.
- No ADR is written at this stage. The numbers 0442 and 0443 are reserved for the PRD and the TechSpec.

## Open Questions

- What tells a session that nobody is present? Default: what the session was told, through the brief's line on whether a person is present, the user's own words, or a continuation that the `unattended` mode armed. A question the user can see is always asked.
- Does a Spec with many open authoring questions still follow the author? Default: yes, with no cap, and the report gives the count. The alternative is to wait when one phase returns more than three, the size of one question call.
- How long does the status view list a Delegated Decision? Default: while its Spec is active, with the user's later answer shown beside it. The alternative is to list it until the user confirms it, which needs a confirm step.
- What are the key and its values? Default: one key with two values, wait and follow, and an absent key means wait. The TechSpec names them.
- From which state is the preference read? Default: as the project committed it, so that an edit made in the session grants nothing, and the session never writes the key.
- Does the answer before the flow reach a later session? Default: it holds for the Specs whose authorization records hold it. A Spec minted later is asked about again, unless the project declares.
- Does the rule hold where no Author Subagent runs, as in Codex, in OpenCode and on the fix or refactor route? Default: yes. The rule speaks of the phase, wherever it runs.
- What is the test of "new scope", "spend" and "data sent outside"? Default: an option is scope when it adds a goal, a Core Feature, a source or a surface the PRD does not hold, or removes a Non-Goal. It is spend when it needs a paid service or a quota the repository does not already use. It is data sent outside when it sends repository or client content to a service the repository does not already send it to.
- Where is a question of the TechSpec or the Tasks phase recorded? Default: in the PRD's `## Decisions`, as Spec 0409 does for a question the TechSpec author asks.
- Does the session tell the user when it adopts an option? Default: no, beyond the report and the status view.
- Which names does the glossary take? Default: the PRD's Glossary settles them, starting from Delegated Decision.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md), the maintainer's request and answers of 2026-11-19, as the operator recorded them. Spec 0413 keeps a copy of the same record.
2. [`references/2026-11-19-how-an-open-authoring-question-is-handled-today.md`](references/2026-11-19-how-an-open-authoring-question-is-handled-today.md), recorded by the author: the instructions quoted in `## Problem`, the two counts and the passages of the external sources.

The plan of Spec 0409 was read in its planning clone and is not adopted. The Backlog Entry `docs/backlog/2026-11-11-the-conexus-team-saw-a-false-status-a-silent-authorization-wait-and-a-black-screen.md` belongs to another context and is named by Spec 0413.

This project's ai-memory messages were listed on 2026-11-19 and none was pending. The Research Subagent did two codebase sweeps and two web sweeps, the second one through Exa. The author did the Secondbrain sweep, because it needs a command, and re-read the passages quoted in this file.

Secondbrain (read on 2026-11-19; `wiki/index.md`, then two `qmd query --all --files --min-score 0.3` calls):

- `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md` — what each harness's question tool takes, and that a subagent of the Claude Agent SDK does not receive the tool. It confirms why an Author Subagent returns its questions.
- `wiki/concepts/claude-code-operacao-produtiva.md` — "manter aprovação humana para envios externos ou ações difíceis de desfazer". It supports the list of what always waits.
- `wiki/concepts/agent-harnesses.md` — the controls of accountability stay with a person however capable the model, and an approval prompt is not a safety control. It supports a record over a prompt.
- The query about a session that adopts an author's recommendation returned only this repository's own mirror (`projects/roundfix/mirror/docs/agents/autonomous-work.md` and two ADRs). The Secondbrain holds no earlier research on the subject.

External sources (found through Exa by the Research Subagent on 2026-11-19; the author re-read every passage quoted here):

- LangGraph, "Interrupts", https://docs.langchain.com/oss/python/langgraph/interrupts — a paused graph waits without limit. It supports the default.
- AG2, "Approval Required", https://docs.ag2.ai/docs/user-guide/tools/approval_required/ — "an unanswered approval is not an approval". It is the position the opt-in departs from, and the source of the dissent.
- Jurgen Appelo, "Leading Agile Developers: The Seven Levels of Authority (Part 2)", https://www.informit.com/articles/article.aspx?p=1675546&seqNum=3 — the level this Idea takes, Inquire, and the one it refuses, Delegate. It shaped F6 and F7.
- Feng, McDonald and Zhang, "Levels of Autonomy for AI Agents", https://arxiv.org/abs/2506.12469 — autonomy as a setting chosen for an Agent. It supports a preference per project.
- Cuadron and Gupta, "Built on Broken Assumptions (BOBA): Why Coding Agents Cannot Be Left Alone", https://arxiv.org/abs/2611.02207 — the cost of a mistake nobody catches, in chains of code requests. It challenged the answer "to the end" and is why F7 and F8 make a decision visible and easy to overturn.

The digest of this research for the project that owns the Secondbrain was not sent by the author, who sends nothing in this cycle. The operator has it in the author's report.
