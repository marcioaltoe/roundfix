---
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
created: 2026-11-18
status: promoted
leverage: compounding
---

# The Idea and the PRD question the rules that block the goal

## Overview

When the author of an Idea or a PRD finds that a recorded decision stops the solution from reaching its goal, the author asks the user one question: this rule imposes this, it prevents these things, do you keep it, change it or retire it? A recorded decision is an accepted ADR, a Normative Clause of a guide under `docs/agents/`, a definition in `CONTEXT.md`, or a document under `docs/references/` that the repository treats as binding. This Idea calls each of them a rule, as the maintainer does. The rule binds until the user answers, the answer is recorded in the artifact, and the Spec that asked carries the change.

It is for the person who decides what a repository builds, and for the Agent that authors for that person. In this repository and in the projects that adopt the Baseline, that person is the maintainer.

It is valuable because the stages that shape a solution receive rules today only as obligations. The cheapest answer to a conflict is then to shrink the solution and say nothing. The maintainer sees the result "com frequência": something is not built because a rule exists, nobody proposes changing the rule, and the reason surfaces one, two or three Specs later.

V1 is a Normative Clause in the Baseline, a step in the authoring skills and a record in fixed form, with no mechanical check. The effort is small, because the question contract, the ADR lifecycle, the Waiting Declaration and the message to another project already exist. The value grows with use: each question leaves the rules better than it found them, with an out-of-date rule changed or retired, or a kept rule that now says why it stays.

## Problem

The maintainer's request of 2026-11-18, quoted from [the adopted Backlog Entry](references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md): "Tenho visto com frequência esse problema na construção de soluções, porque tem uma regra, então, nem é proposto a mudança da regra, deprecar, regra, simplesmente não faz, porque existe uma regra. Em algum momento, talvez depois de uma, duas, três specs, é falado, não, não foi feito, porque existe essa regra, sim, lá do início do projeto." He names the cause in the same request: "As ADRs, regras do specific repository e outras referências que possam existir dentro do projeto, elas são para guiar a implementação, são guard-rails de guiação da implementação. No momento que a gente está trabalhando no IDEA e no PRD, a gente está trabalhando com criatividade."

The instructions an author reads today explain the silence. Line numbers are of this repository at `5d3237bd5f6db3fde9cfeaf74de65b23d8f42f65`.

- The root instruction file says a guide "cannot weaken a universal Normative Clause or confirmed project decision" (`AGENTS.md:5`). It makes no difference between the stage that shapes a solution and the stage that builds it.
- The only clause about a conflict forbids one of two silences: "flag conflicts instead of silently overriding repository decisions" (`docs/agents/domain.md:7`). No clause forbids the other one, which is giving up a goal without saying so.
- `write-idea` names `CONTEXT.md` as its only ground-rule read (`.agents/skills/write-idea/SKILL.md:32`), mentions the active ADRs once as something already read (`:92`) and has no step about a rule. `write-prd` reads `docs/adr/` first (`.agents/skills/write-prd/SKILL.md:44`) and records what it finds in the row "Active ADR obligations" (`.agents/skills/write-prd/references/prd-template.md:26`). `brainstorming` reads `CONTEXT.md` and the Specs before it asks, and names no existing ADR or rule (`.agents/skills/brainstorming/SKILL.md:29`).
- Only `write-techspec` may contest one: "respect prior decisions or challenge them explicitly — never contradict one silently" (`.agents/skills/write-techspec/SKILL.md:27`). It does so after the PRD fixed the scope, and without asking the user.

A rule therefore reaches the Idea and the PRD as an obligation to list. An author that meets a rule in the way of the best solution has two routes the instructions allow: drop the solution, or keep it and leave the rule to the TechSpec. Both are silent. The user learns that a rule decided the scope when someone notices the missing part.

The plans of this repository show which record is missing. In the 41 Specs planned from 2026-11-06 to 2026-11-18, 15 PRDs and 2 Ideas draw a scope exclusion at an ADR or a guide. Twelve PRDs amend or supersede an ADR: 11 times by the author's own statement, and twice inside an answer the maintainer gave to another question. None of the 34 lines counted records a question to the maintainer about whether a recorded decision is kept, changed or retired ([measurement](references/2026-11-18-rules-named-in-recent-plans.md)). The count cannot tell a boundary the author chose from a goal the author gave up, because the file does not say which it was.

Asking by hand does not close the gap. The user can ask in each session whether a rule is in the way, but only when he already suspects one. A rule that is badly written, out of date or never meant for the case keeps binding, because the one person who can change it is never asked.

### Market Data

- Coding agents follow a project's instruction file even when they notice that it conflicts with what they see: "Without the sentence, 90.1% of the agents' choices followed the stale file. Claude Code pointed out the conflict in 79 of 80 runs and still followed the file." One sentence that gave the code precedence "raised the share of choices that follow the code from 4.2% to 91.4%" (Xu and Wu, 2026, 320 runs over 40 projects). The study measures code conventions, not product goals. This Idea infers from it that an Agent told a decision cannot be weakened will comply, and that one sentence in the instructions changes what it does.
- Questions cost attention, and a limit controls the cost: "limiting agents to one to three questions reduces user-attention cost by 36–73% while keeping 86–94% of the task reward they earn with unlimited questions", and "reviewing each question against the memory of users' earlier answers cuts attention cost by another 13%" (Zhang and others, 2026). It supports a cap on the question and a record of the earlier answer.
- An accepted decision changes through its record: "Once an ADR is accepted, it should never be reopened or changed - instead it should be superseded", and the record is "a good place to mention any changes in the product context that should trigger the team to reevaluate the decision" (Fowler, 2026).
- Spec Kit, a tool for spec-driven development that keeps binding principles in a constitution file, resolves a conflict in one direction: "Violations are resolved by changing the spec, plan, or tasks — not by diluting a principle", and a change to the principles needs "a PR with rationale, maintainer approval, and a version bump" (its own constitution). The web sweep found no tool that puts a keep, change or retire question to the person who decides.

## Core Features

The question below is called a Rule Challenge in this Idea. The name is a candidate for the glossary, which the PRD settles.

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The Baseline forbids both silences | Critical | A Normative Clause says that in brainstorming, the Idea and the PRD, a rule that blocks the goal is never obeyed in silence and never overridden in silence: the author asks the user. Every adopter receives the clause with the Baseline update. |
| F2  | The trigger | Critical | The author asks only when a rule blocks the goal: the solution it would recommend, or a part the goal needs, cannot exist without contradicting the rule, or the rule stops it from stating the problem. Extra work, taste and friction are not a block. A change that the request or a recorded answer of the user already decides needs no question. |
| F3  | The question | Critical | One question per rule, through the question tool the session exposes. It names the rule and where it lives, what the rule imposes, up to three things it prevents for this goal, and the author's reading of why the rule exists and whether it looks out of date or badly written. The answers are keep, change and retire, each with its consequence and the recommended one first. A recommendation, silence, a dismissal or a timeout is never an answer. |
| F4  | The rule binds until the answer | Critical | Asking changes nothing. The author plans on the user's answer only after the user gives it, and the rule's text stays as it is until a Task of the Spec changes it. |
| F5  | The record | High | The Idea and the PRD carry a section that holds one record per question in fixed form: the rule, what it imposes, what it prevents, when it was asked and in which phase, the answer in the user's words, and what the Spec does with it. The section reads `None.` when no rule blocked the goal. The PRD carries the Idea's records forward. |
| F6  | The Spec carries the change | High | On change or retire, the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it. An ADR changes through the lifecycle the docs-layout guide already defines. |
| F7  | A kept rule says why it stays | High | On keep, a Task of the Spec adds one dated line to the rule: what was asked and why it stayed. A later author reads the line and asks again only with a new fact, which the question names. |
| F8  | A Spec waits when nobody can answer | High | With no person present, only that Spec stops: the feature stays or becomes a Backlog Entry whose Waiting Declaration names the question, the entry holds the whole question, and the queue goes on with other work. An Author Subagent returns the question to the session that started it. |
| F9  | A rule the Baseline owns becomes a proposal | Medium | In an adopter, a clause inside a Managed Region cannot be edited. On change or retire, the user's decision goes to the Baseline's owner as one ai-memory message, the rule stays in force until the Baseline changes, and the record names the message. |
| F10 | The TechSpec asks before it contradicts an ADR | Medium | When a design needs what an ADR rejected and the user has not decided it, the TechSpec author asks the same question before it supersedes or amends the ADR. |

## The question and its record

An illustration with an invented rule, not one of this repository:

```text
Rule: ADR-NNNN, "every report is produced at the close of the day" (docs/adr/NNNN-….md)
It imposes: no report is computed while the day is open.
It prevents, for this Idea's goal (a stock-out alert the store can act on):
  1. computing the alert when the sale happens;
  2. telling the store before the shift ends.
The author's reading: the ADR is from 2025 and gives a reason that the newer
sales feed removed. It looks out of date, not wrong when it was made.

1. Change the rule (Recommended): the PRD requires the amendment of ADR-NNNN,
   and this Spec delivers it with the alert.
2. Keep the rule: the Idea goes on with an alert at the close of the day, and
   ADR-NNNN gains a dated line with this question and your reason.
3. Retire the rule: ADR-NNNN leaves the tree under the retirement clause, and
   Git keeps its text.
```

The record of that question, as the Idea and the PRD hold it under a `## Rule challenges` heading:

```markdown
- rule: ADR-NNNN, `docs/adr/NNNN-….md`
- imposes: no report is computed while the day is open
- prevents: computing the alert when the sale happens; telling the store before the shift ends
- asked: 2026-11-18, idea
- answer: change — "Mudar a regra (Recommended)", 2026-11-18
- effect: the PRD requires the amendment of ADR-NNNN; a Task delivers it with the alert
```

The `answer` field holds `keep`, `change`, `retire` or `waiting`. On keep, `effect` names the scope the Spec keeps inside the rule and the dated line a Task adds to it. On a rule the Baseline owns, `effect` names the proposal and its message. A question about a rule that already carries a kept line adds a `new fact` field that says what this question knows and the earlier one did not.

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Reports that a delivered Spec left something out because of a rule nobody asked about | 0 in the 30 days after a repository applies the update | At each cycle close, count the Backlog Entries, Findings and messages that name that cause, in this repository and in the adopters' messages to it. No count exists today; the maintainer reports the problem "com frequência". |
| Ideas and PRDs that carry the record | 100% of the first 20 plans authored after the release, from 0 of 47 today | Repeat the method of [the measurement](references/2026-11-18-rules-named-in-recent-plans.md) over the new plans and count the files that hold the section, with records or `None.`. |
| Question load | At most 1 question in at least 90% of those 20 Specs, and never more than 3 in one Authoring Phase | Count the records per Spec and per phase in the same 20 plans. |
| Repeated questions | 0 questions that are about a rule with a kept line and name no new fact | For each record, read the rule it names at the plan's commit. |
| Wait when nobody is present | A median under 24 hours from the wait to the answer | Compare the `since` of the Waiting Declaration with the date in the record's `answer`. |

## Feature Assessment

| Criteria        | Question                                            | Score   |
| --------------- | --------------------------------------------------- | ------- |
| Impact          | How much more valuable does this make the product?  | Strong — it removes a failure the maintainer sees often and that costs whole Specs when it is found late. It adds no capability of its own. |
| Reach           | What % of users would this affect?                  | Must do — every repository that adopts the Baseline receives the clause, and every Idea, PRD and TechSpec authored there passes through the step. |
| Frequency       | How often would users encounter this value?         | Maybe — nobody has counted how often a rule blocks a goal. Here 15 of 41 PRDs draw their scope at a recorded decision, and the count cannot say how many gave up a goal. |
| Differentiation | Does this set us apart or just match competitors?   | Strong — the web sweep found no tool that asks the person who decides whether a rule stays. Spec Kit changes the spec, never the principle. |
| Defensibility   | Is this easy to copy or does it compound over time? | Maybe — a clause and a skill step are text anyone can copy. What compounds stays in each repository: rules with their reasons, and out-of-date rules retired. |
| Feasibility     | Can we actually build this?                         | Must do — the question contract, the ADR lifecycle, the Waiting Declaration and the message path exist. V1 is text in the Baseline source and in five owned skills, with their derived records. |

The leverage is compounding. The verdict is to proceed. Two results would change it: questions that almost always end in keep, which means the trigger is too loose, or no questions at all while the maintainer still reports late discoveries, which means the step is skipped and V1 needs the check it left out.

## Trade-off Insights

- **Recommended approach:** the clause, the step and the record, as the maintainer chose. The clause is what reaches every session, including one that authors through a Vendored Skill this repository cannot edit. The step puts the question where the author already asks its questions. The record is what turns "no rule was in the way" from a guess into a statement.
- **Key trade-offs:**
  - Guard-rails that erode against mistakes that freeze. A question that is cheap to ask invites a question about every inconvenient rule, and a user in the middle of a feature tends to say change. A question that is hard to ask is today's silence. V1 sets a narrow trigger, a cap of three per Authoring Phase, a rule that binds until the answer, and a kept line that demands a new fact.
  - A record in every artifact against a clause alone. The published measurement suggests one sentence may be enough to change what an Agent does. Without the record nobody can tell "the author looked and found nothing" from "the author did not look", and a kept rule is asked about again.
  - The change travels with the Spec against a change made at once. With the Spec, the rule and the work that depended on it land together and are verified together. Until then the old text stays, so another Spec authored in the meantime can be narrowed by the same rule or ask the same question.
  - A Spec that waits against a Spec that goes on inside the rule. Waiting builds nothing the user did not choose. It also stops a Spec on the judgment of an Agent, when nobody is there to say the judgment is wrong.
- **Risks identified:**
  - The trigger is the author's judgment. An author that wants its own solution can read friction as a block. The question-load KPI and the share of keep answers show it.
  - The section fills with `None.` by habit, because nothing checks it in V1.
  - In this repository authors amend an ADR in 12 of 41 PRDs. If each amendment became a question, the load would rise sharply. F2 excludes a change the request or a recorded answer already decides, and that boundary needs exact words in the PRD.
  - Adopters send proposals about Baseline clauses, and each is a message this repository triages. Only a user's change or retire answer sends one, so the user's decisions bound the volume.
  - Kept lines accumulate on ADRs that are asked about often.
- **Dissenting view:** a skill step is an instruction, and the same research that supports the clause shows Agents notice a conflict and comply anyway. Without a check, V1 may change the wording an author reads and leave the behavior where it is. The opposite dissent also stands: the question hands every author a legitimate way to reopen settled decisions, and a rule that is asked about three times in a month has stopped being a guard-rail even if the answer is always keep.
- **Stretch goal (V2+):** every new ADR and repository rule states why it exists and what would make it worth revisiting, so a question finds its facts in the rule. The Advisory Judgment reads a PRD's exclusions and says which ones look like a goal given up to a rule. Proposals about the same Baseline clause from several adopters are grouped into one decision for the Baseline.

## Opportunity Scan

The original idea is bounded by how often a rule blocks a goal, which nobody has counted in the adopters. If it works fully, no narrowed solution ships unless the user chose the narrowing, and the rules improve as a side effect of ordinary authoring.

- **More ambitious: a check and a judge.** `roundfix spec check` reports a scope exclusion that names a rule and has no record, and the Advisory Judgment says which exclusions look like a goal given up. It closes the dissent about skipped steps. Scale: medium to large. Impact Strong, Reach Must do, Frequency Strong, Differentiation Strong, Defensibility Strong, Feasibility Maybe: 19 exclusions in 41 PRDs name a rule and most read as plain boundaries, so a structural check would ask about the wrong ones.
- **Simpler: the clause alone.** One Normative Clause and one line in each skill, with the answer in free text under `## Decisions`. Scale: small. Impact Maybe, Reach Must do, Frequency Maybe, Differentiation Maybe, Defensibility Pass, Feasibility Must do. It leaves no way to measure anything in `## KPIs` and no memory of a kept rule.
- **Adjacent: rules that state their reason.** Every ADR and repository rule records why it exists and when to revisit it. It attacks the badly written rule at its source and makes each later question cheaper. Scale: small to medium. Impact Strong, Reach Must do, Frequency Maybe, Differentiation Maybe, Defensibility Strong, Feasibility Strong. It does not end the silence on its own, and it shares the ADR layout with Spec 0410.

**Chosen direction:** original — the maintainer chose "Cláusula, passo e registro (Recommended)". The kept line takes one part of the adjacent idea: a rule that is asked about and kept records its reason, so the rules that get in the way are the first to explain themselves. The check, the judge and the reason on every new rule are the stretch goal.

## Integration with Existing Features

| Integration point | What it gains or how it is used |
| --- | --- |
| The Baseline guides under `docs/agents/` | One Normative Clause in stack-neutral words: the two silences, the trigger, the question, the binding until the answer, the wait and the proposal. The convention of the dated line on a kept rule goes where the ADR lifecycle and the repository rules are already described. |
| The instruction hierarchy in the root file | Unchanged in meaning. The clause says that a question is not a weakening: the user who confirmed a decision confirms it again or changes it. |
| `brainstorming` | It reads the recorded decisions while it explores. When the direction the user wants meets a rule that blocks it, it asks instead of steering the user away, and hands the record to the skill it routes to. |
| `write-idea` | A ground rule that a rule guides the solution and does not bound the problem. The codebase sweep lists the rules that touch the idea. The Opportunity Scan never drops an alternative only because a rule forbids it: it marks the rule. The template gains the section. |
| `write-prd` | It reads ADRs and rules as obligations and as candidates: an obligation that blocks a goal becomes a question before the Project Constraints snapshot is written. A change or retire answer becomes a requirement. The template gains the section, and the Idea's records are copied as the judgment is today. |
| `write-techspec` | Its existing permission to contest an ADR gains the question: the author asks before it supersedes or amends an ADR the user has not decided about, and the `## ADR Alternatives` line points at the record. |
| The Roundfix skill, `modes` reference | The `unattended` and `follow-up` Skill Modes say what happens to a Spec that waits on the question, next to what they say today about the Idea stage (`.agents/skills/roundfix/references/modes.md:42`). |
| The question contract (`docs/agents/agent-instructions.md:90`) | Used as it is: `AskUserQuestion` in Claude Code, `request_user_input` in Codex and `question` in OpenCode, the same question numbered in chat when there is no tool, two or three options with the recommended one first. |
| The ADR lifecycle (`docs/agents/docs-layout.md:9`, `:21`, `:181`) | Used as it is. Change is an amendment or a successor ADR. Retire is the retirement clause: the file leaves the tree and Git keeps the text. This Idea adds no status. |
| The Waiting Declaration and the Follow-up Command | Used as they are: the Backlog Entry of a waiting Spec names the question in `waiting_on`, and the follow-up shows it. |
| The message to the owning project (`docs/agents/agent-instructions.md:86`) | Used as it is for a rule the Baseline owns. |
| `grilling` and `domain-modeling` | Unchanged. They are Vendored Skills, and the Baseline clause governs a session that authors through them. |
| `CONTEXT.md` | The names the PRD settles for the question, for the dated line on a kept rule and, if one is needed, for the set of things that can be asked about. |

## Prerequisite and Shared Files

This Spec is delivered after Spec 0407 and Spec 0408 and before Spec 0410, in the order the maintainer set. Nothing in it depends on their code. The four are authored in parallel, and these files can meet:

- The owned Roundfix skill: its version fields and `references/modes.md`. Specs 0407, 0408 and 0410 also change this skill.
- `skills/testdata/owned-skill-versions.json` and the `skills/` mirror, for `write-idea`, `write-prd`, `write-techspec`, `brainstorming` and `roundfix`.
- The Baseline source under `internal/baseline/assets/`, `internal/baseline/module-versions.json`, the digest pin, the catalog snapshots, the plan goldens and `docs/agents/setup-context.json`. Spec 0410 changes the same modules.
- `docs/agents/docs-layout.md` and its module, if the dated line on a kept rule is described there. Spec 0410 owns the documentation rules and the status of a superseded ADR. This Spec adds no status and changes none: its retire answer follows whatever the lifecycle says when the Task runs.
- `CONTEXT.md`.

This Spec changes no command, so `docs/references/behavior-surfaces.json` stays as it is. It does not change `qa-gate`, `write-tasks` or `archive-spec`.

By the sizing criterion of the Spec routing guide this is one Spec. It holds one context: how an Authoring Phase treats a recorded decision that blocks the goal. The clause, the skill steps, the record and the glossary names complete one another: a clause without the step changes nothing in a session a skill leads, and a step without the clause never reaches a session that authors through a Vendored Skill. The work is text in the Baseline source and in owned skills, with the records each edit regenerates, so it fits two or three implementation Tasks and the QA gate, each large enough to earn its verification cycle. A split would deliver half a behavior twice.

## Out of Scope (V1)

- **A mechanical check** — the maintainer chose V1 without one. A structural check cannot tell a boundary from a goal given up: 19 exclusions in 41 PRDs name a rule, and most read as boundaries.
- **A judgment of exclusions by the judge** — it is the stretch goal, and it needs the record of V1 to judge against.
- **A review of every rule on a schedule** — the trigger is a blocked goal, not the age of a rule.
- **Questions during implementation** — in a Run, in `implement-task` and in the QA gate a rule is a guard-rail, as the maintainer says. An Agent that meets a conflict there flags it, as the existing clause says (`docs/agents/domain.md:7`).
- **New ADR statuses or a new ADR layout** — the lifecycle is used as it is, and the documentation rules belong to Spec 0410.
- **Rewriting existing rules to add their reasons** — a rule gains a reason only when it is asked about and kept.
- **Rules that live outside the repository** — user-level instruction files, ai-memory pages and a harness's own rules are not the project's recorded decisions. A default of a Vendored Skill is already overridden by a repository rule.
- **Edits to Vendored Skills** — a correction to `grilling` or `domain-modeling` goes upstream.
- **A review of delivered Specs for scope given up in the past** — the record starts with the plans authored after the release.

## Idea-stage judgment

- item: [`references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md`](references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.86)
- unverified_assumption: unsure (0.64)
- market_question: no (0.12)
- verdict: indicated
- stage: ran
- shape: changed — the question gained a trigger and a cap, the Spec that asks carries the change, a kept rule records why it stays, a Spec waits when nobody can answer, a rule the Baseline owns becomes a proposal to its owner, and the TechSpec asks before it contradicts an ADR.

The verdict was indicated, and the maintainer chose "Ideia agora (Recommended)" on 2026-11-18. The judged item is the adopted Backlog Entry. The `shape` line is the author's and waits for the approval of this draft.

## Decisions

Every answer is the maintainer's, as recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The text in quotes is the option he chose, verbatim. The text after the dash is that option's description as the operator recorded it.

- When to author it (2026-11-18): "Spec após a release". Authorization (2026-11-18): "Mesmas do BE (Recommended)". The PRD quotes it in the authorization record.
- The Idea stage (2026-11-18, ~19:55 BRT): "Ideia agora (Recommended)". Delivery order (2026-11-18, ~19:55 BRT): "Corretiva → regras → Indexus (Recommended)" — o fix da fila entra logo após a corretiva.
- Size of V1 (2026-11-18, ~20:30 BRT): "Cláusula, passo e registro (Recommended)" — cláusula no Baseline, passo nas skills que faz a pergunta, registro em formato fixo da pergunta e da resposta no artefato; sem verificação mecânica.
- The TechSpec (2026-11-18, ~20:30 BRT): "Entra, só ao contradizer ADR (Recommended)" — Ideia, PRD e brainstorming questionam a regra que bloqueia o objetivo; o TechSpec pergunta quando o desenho precisa contradizer uma ADR.
- Who applies the answer (2026-11-18, ~20:30 BRT): "A própria Spec (Recommended)" — resposta registrada na Ideia/PRD; uma Task da mesma Spec altera a regra junto com o que dependia dela.
- Nobody present (2026-11-18, ~20:30 BRT): "A Spec espera (Recommended)" — só essa Spec para: Backlog Entry com waiting_on nomeando a pergunta; a fila segue.
- Reach (2026-11-18, ~20:30 BRT): "Sim, vira proposta ao Roundfix (Recommended)" — tudo gravado pode ser questionado; regra do Baseline vira proposta ao Roundfix por ai-memory e vale até o Baseline mudar.
- The keep answer (2026-11-18, ~20:30 BRT): "Na própria regra (Recommended)" — a regra/ADR mantida ganha uma linha datada com o que foi perguntado e por que ficou; só se pergunta de novo com fato novo.

The author's own, made in this stage:

- The questions did not offer the dismissal of the stage again, because the maintainer had already chosen "Ideia agora (Recommended)". The operator asked them, because the Author Subagent cannot ask.
- The Backlog Entry is adopted here: its status is `promoted`, it names this Spec, and it moved into `references/`.
- The trigger excludes a change the request or a recorded answer already decides. Without that line, the 13 amendments counted in 12 of 41 PRDs would each become a question.
- The dated line on a kept rule names no Spec: a guide must never reference one, and an archived Spec can be deleted at any time (`docs/agents/docs-layout.md:198`).
- The author read the rules that touch this Idea: the instruction hierarchy, the question contract, the ownership of Vendored Skills and the upstream flow of knowledge. None blocks the goal, so this Idea asks no such question. The one that shapes it most is the limit of two or three options, which is the first open question.
- No ADR is written at this stage. The ADR numbers 0436 and 0437 are reserved for the PRD and the TechSpec.

## Open Questions

- Three answers or four? The maintainer's words are "manter", "mudar", "deprecar ou excluir". Default: three answers, keep, change and retire. The question contract takes two or three options, the Codex tool takes no more than three, and under the retirement clause a deprecated ADR and a deleted rule end the same way: the file leaves the tree and Git keeps the text. The user's own words go into the record. If the maintainer wants deprecate and delete apart, the question contract changes first.
- How many questions can one Authoring Phase ask? Default: three, the limit of one question call. When more rules block, the author asks about the three that prevent the most and lists the others in the artifact's open questions.
- Where does the dated line go when the kept rule is a Baseline clause in an adopter? Default: in the repository's own `docs/agents/specific-repository.md`, because the repository cannot edit inside a Managed Region.
- What does a Spec do while its proposal to the Baseline's owner has no answer? Default: it goes on inside the rule and records the scope it gave up. The user can make it wait instead.
- What carries the wait when the Spec already has a folder, as when an approved Idea reaches the PRD with nobody present? A Waiting Declaration exists only on a Backlog Entry and on a Finding. Default: no PRD is written and a Backlog Entry that names the Spec carries the question until the answer arrives. The PRD stage checks this against how the Follow-up Command lists what waits.
- Does a minimal PRD written for a fix or a refactor carry the section? Default: only when a question was asked. An Idea and a feature PRD always carry it.
- How does the dated line read, and do lines stack? Default: one line per kept question, the newest last, in the form of the "Amended by" lines ADRs already carry.
- How is the Spec accepted without a mechanical check? Default: each Task's Verification checks the words it adds to the clause, the skills and the templates, the outside evidence is the published measurement quoted in `### Market Data`, and the second KPI is read over the first 20 plans after the release.
- Which names does the glossary take? Default: the PRD's Glossary settles them, starting from Rule Challenge.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md`](references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md), the Backlog Entry that holds the maintainer's request, verbatim.
2. [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md), the answers behind `## Decisions`.
3. [`references/2026-11-18-rules-named-in-recent-plans.md`](references/2026-11-18-rules-named-in-recent-plans.md), recorded by the author: the count used in `## Problem` and `## KPIs`.

The source of Spec 0410 is not adopted here. It asks for "a clear status for superseded ADRs so no expired decision is cited as current" (`docs/backlog/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md@5d3237bd5f6db3fde9cfeaf74de65b23d8f42f65`), which is why this Idea adds no status. The other open Backlog Entries belong to other contexts and stay where they are.

This project's ai-memory messages were listed on 2026-11-18 and none was pending. The Research Subagent did the codebase sweep and the web sweep. The author did the Secondbrain sweep, because it needs a command, and re-read the passages quoted in this file.

Secondbrain (read on 2026-11-18; `wiki/index.md`, then four `qmd query` attempts that all failed with "database is locked" while another process ran `qmd update`, so the pages below were found through the index):

- `wiki/concepts/politica-de-retirada-e-poda.md` — "Retirar com razão, nunca em silêncio", pruning by an event and never by age, and "agentes varrem … a pessoa decide". It supports a trigger that is an event, an answer that is the user's, and no review on a schedule.
- `wiki/concepts/ciclo-de-vida-de-decisoes-e-raiz-docs-history.md` — the lifecycle of a decision and its change of 2026-10 from a permanent history to deletion with the text in Git. It gave the retire answer its meaning.
- `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md` — what each harness's question tool takes: Claude Code two to four options, Codex two or three, and no question tool in a subagent. It set the default of three answers and confirms that an Author Subagent returns the question.
- `wiki/concepts/fundamentals-of-software-architecture-capitulo-23-19-architecture-decisions.md` — the Groundhog Day anti-pattern: "people don't know why a decision was made, so it keeps getting discussed over and over". It supports the dated line on a kept rule.
- `wiki/concepts/lddd-capitulo-11-evoluindo-decisoes-de-design.md` — design decisions age, and its checklist asks "O que faria essa decisão mudar?". It supports the stretch goal of rules that state when to revisit them.

External sources (found through Exa by the Research Subagent on 2026-11-18; the author re-read every passage quoted here):

- Xu and Wu, "Stale Instruction Files Override the Code in Coding Agents: A One-Sentence Fix", https://arxiv.org/abs/2610.30571 — the evidence that an Agent complies with an instruction it sees a conflict with, and that one sentence changes it. It supports the clause and challenges the hope that Agents ask on their own.
- Zhang and others, "Ask Less, Ask Better: Allocating User Attention in Proactive LLM Agents with Interruption Budgets", https://arxiv.org/abs/2610.24411 — it gave the cap and supports reading the earlier answer before asking again.
- Martin Fowler, "Architecture Decision Record", https://martinfowler.com/bliki/ArchitectureDecisionRecord.html — an accepted ADR is superseded, never edited. It challenged the dated line, which this Idea keeps because it records a question and changes no decision, as the "Amended by" lines of this repository already do.
- The Spec Kit constitution, https://github.com/github/spec-kit/blob/bba473c2/.specify/memory/constitution.md — a conflict is resolved by changing the spec and never the principle. It is the position this Idea departs from.
- Mike Salari, "Disagree and commit: what engineers owe after the call", https://salari.dev/writing/disagree-and-commit-what-engineers-owe-after-the-call/ — "Record the disagreement next to the decision: the option you argued for, what you predicted, and the reopen criteria." It supports the kept line and the new fact.
- Digital.gov, "Divergent and convergent thinking", https://digital.gov/guides/hcd/design-operations/thinking — the design phase alternates the two, and the guide tells the reader to "generate project ideas, then identify constraints". It matches the maintainer's split between the creative stages and implementation.
- Mahato, Sieczkowski and Kuppusamy, InfoQ, "Agentic Fitness Functions: Extending Evolutionary Architecture Beyond Deterministic Rules", https://www.infoq.com/articles/agentic-fitness-functions-evolutionary-architecture/ — it lists "stale ADR assumptions" among the risks that need judgment, and separates "deterministic gates from agentic advisory signals" that escalate to people. It supports leaving the answer to a person and the check out of V1.
- Scripta, "Reopen a decision only when the trigger occurs", https://www.scripta.cloud/blog/reopen-a-decision-only-when-the-trigger-occurs — "A decision stays closed by default and reopens because evidence crossed an agreed boundary." It supports the new fact a repeated question must name.

The digest of this research for the project that owns the Secondbrain was not sent by the author, who sends nothing in this cycle. The operator has it in the author's report.
