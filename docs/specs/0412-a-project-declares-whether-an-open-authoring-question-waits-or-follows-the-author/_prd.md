---
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
status: active
created: 2026-11-19
surfaces: [cli, docs]
type: feat
---

# A project declares whether an open authoring question waits or follows the author

When the author of a PRD, a TechSpec or a Task Graph ends its phase with a question that no source settles and nobody is there to answer, the instructions give the session no route: it may not take the option the author recommends, so the Spec stops, and nothing records that it stopped or why. After this Spec a project says what happens. By default the Spec waits, and the wait has a record. A project can declare in its Project Config that the session follows the author's recommendation instead, and each adopted option is then recorded in the PRD as a decision taken by delegation, listed in the session's report and in the status view, and open to the user's later answer. The questions that are the user's alone always wait.

## Sources

1. [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md), the maintainer's request of 2026-11-19 in his own words, the answers of the Idea stage, the approval of the Idea and the authorization of this Spec, recorded by the operator. The request exists only in the session, so this record is the source; no Backlog Entry, Finding or ai-memory message is adopted and nobody asked for a notice.
2. [`references/2026-11-19-how-an-open-authoring-question-is-handled-today.md`](references/2026-11-19-how-an-open-authoring-question-is-handled-today.md), recorded by the author of the Idea: what the instructions say about a question an author cannot settle, the count of PRDs authored with nobody present, the count of answers that took the recommended option, and the passages of the published sources.
3. [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), recorded by the author of this PRD: the characterization of the surfaces this Spec changes, the prototype of the Baseline rule and the Governed Paths.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from: its features F1 to F10 are this Spec. This project's ai-memory inbox was listed on 2026-11-19 and held no pending message, and `docs/findings/` holds no Finding. Seven Backlog Entries are open at the authoring base. Two are the sources of Spec 0410. One reports three unreproduced symptoms from another team, among them a command that waits for an authorization in silence; the Idea of the Spec about permission prompts, which is delivered after this one, names it. One is about the size of the context an Authoring Phase starts with and waits on a measurement; this Spec adds text to what an author reads (`## Risks`). Three record defects found on 2026-11-19: in the Verification reader of the Spec Consistency Check, in a deadline of the Baseline's analysis and in two unused probes. None shares this Spec's context, and none is adopted.

## Prerequisites

This Spec's text names the Rule Challenge, which Spec 0409 adds to the glossary and to the Baseline, so its Verification depends on Spec 0409 being delivered. No other behavior depends on another Spec's code. It was authored against `894710c3`, where Spec 0411 is delivered and Specs 0409 and 0414 are planned, and it is delivered after Specs 0414, 0409 and 0410 and before Spec 0413, in the order the maintainer set. These files are shared, and each Task regenerates the shared records on the tree it runs on, with no version number written by hand:

- Spec 0409 adds a rule to the Baseline module `spec-workflow` and changes the Spec routing guide's rule list, the closed list of clause forces, the `write-prd` and `write-techspec` skills, the PRD template, the Roundfix skill's `modes` reference, `CONTEXT.md` and one user-guide page. This Spec adds a rule of its own to the same module and its own sentences to the same files, after the ones Spec 0409 adds.
- Spec 0410 changes other Baseline modules and the docs-layout guide, and shares the module version record, the derived Baseline files, the Standard TypeScript profile's pin, this repository's Setup Manifest, the Roundfix skill's version and `modes` reference, the owned-skill version record, `CONTEXT.md` and the Behavior Surface Record. It also gives every ADR two header fields; ADR-0442 arrives with this plan, and the Task that touches ADR-0280 writes the header and the notice in the form the docs-layout guide prescribes when it runs.
- Spec 0414 changes how the status view counts a delivering Spec's Tasks, one sentence of the `follow-up` guide and `CONTEXT.md`. This Spec adds two blocks of lines and two fields to the same view and leaves the queue table as Spec 0414 leaves it.
- Spec 0413 is delivered after this one and changes the Roundfix skill, the user guide, the Doctor Command and the Behavior Surface Record. It adds no Project Config key, and this Spec touches no permission file of a harness.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity; the new Baseline rule and its clauses take ids in the form the Spec workflow module already uses, and a Delegated Decision is identified by its Spec, its date and its question. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the status view reads local Git objects and local files, and no test or Task asks a person, a provider or ai-memory anything. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0442 (new) records the preference, the Reserved Questions and the record (ADR-0442: "A project that declares nothing waits."). The question contract is used as it is and keeps its sentence about what is never the user's choice (ADR-0320: "Up to three, independent, in one pending call."). The Clarification Record gains a third entry and keeps its two forms and their check, so ADR-0280 gains an amendment notice (ADR-0280: "A PRD records who decided."). The wait uses the Waiting Declaration without changing it (ADR-0353: "makes it a Waiting Item in a fifth class") and reads in the status view as every open Backlog Entry does (ADR-0410: "Every open Backlog Entry is a row marked as not yet a Spec"). The Idea stage and the Rule Challenge keep their own rules and are Reserved Questions (ADR-0403: "no user is present, no PRD is written: the feature stays a Backlog Entry that"; ADR-0436: "The rule binds until the user answers."). An Author Subagent still returns its questions (ADR-0401: "Author Subagent cannot ask the user, so the session gathers the decisions first"). The preference is a key only a project sets, read where a branch cannot change it, as an earlier key already is (ADR-0417: "off by default and ignored in User Config, makes the declaration required."; ADR-0423: "key from the Project Config at the default branch head as the local refs hold it"), and a binary older than the release that adds it is refused by the existing guard (ADR-0392: "refuses, before it reads Project Config, every command that loads the repository's configuration"). The new clauses are written for every adopter (ADR-0186: "Baseline guidance states what the product does, in adopter-neutral words"; ADR-0222: "A Baseline guide says only what holds for the repository that reads it") and state the rule once (ADR-0257: "One statement per rule."). The skills and the guides describe the behavior in the same delivery (ADR-0256: "A release waits for the skills that describe a changed surface"), and the authorization record carries the Release Ceiling (ADR-0349: "A Spec's Release Ceiling answers the release plan's approval"). ADR-0203 is listed to record why it does not apply: no hand-written rule of a repository is promoted, the clauses describe what the product's own preference does, and the evidence from outside Roundfix is the published sources of `## Acceptance evidence` (ADR-0203: "The Baseline grows only with evidence from outside Roundfix."). ADR-0437, ADR-0319, ADR-0338, ADR-0407, ADR-0434, ADR-0408, ADR-0427, ADR-0288, ADR-0258, ADR-0342, ADR-0383, ADR-0390, ADR-0404, ADR-0405, ADR-0412, ADR-0414, ADR-0415, ADR-0418, ADR-0420, ADR-0421, ADR-0422, ADR-0424, ADR-0425, ADR-0426, ADR-0429, ADR-0431 and ADR-0441 are listed only to record why they do not apply; the TechSpec relates each one. No accepted ADR is superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: his answer of 2026-11-19, "Mesmas do BE para as duas (Recommended)", which grants Specs 0412 and 0413 the cycle's authorization "Mesmas do BE (Recommended)" of 2026-11-18 and which the authorization record quotes with its provenance. It keeps the grant of every owned skill and of the Governed Paths each Spec declares, Baseline modules included. It grants no Makefile, workflow, Go module, lint or formatter change except for a correction that requires one, and the design needs none; `.roundfixrc.yml` is not granted and no Task changes it. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/clarification-record.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A Spec never stops in silence on a question its author could not settle: it either goes on with a recorded Delegated Decision or waits with a record that holds the question.
- A project that declared nothing is never built on an answer nobody gave.
- A project that follows the author loses no Spec to a question of detail, and its user reads a short list of adopted options in place of answering each one.
- A reader of a PRD can always tell an answer the user gave from an option a session adopted for him.
- The questions that are the user's alone, which are his limits, every authorization, the Idea stage and a question about a rule, wait in every project.

## User Stories

1. As the maintainer of a repository, I want to declare once, in the project's configuration, whether an open authoring question waits for me or follows the author's recommendation, so that every unattended flow of the project behaves the same way without asking me again.
2. As the maintainer of a project that declares nothing, I want to be asked once before a flow starts, with waiting as the recommended option, so that I choose for that flow while I am still there.
3. As the maintainer away from the session, I want a Spec whose question waits to stop alone and to hold the whole question where the status shows it, so that I find it when I return and the other work goes on.
4. As the maintainer of a project that follows the author, I want every adopted option recorded with its question, its options and the author's reason, and listed in the report and in the status view, so that one read tells me what was decided for me.
5. As the maintainer, I want a question about scope, spend, data sent outside, an authorization, the Idea stage or a rule to wait whatever the preference says, so that a preference about details never decides what is mine.
6. As the maintainer, I want my later answer to a delegated question to win, so that an adopted option I disagree with is corrected by the route the Spec's stage allows.
7. As a reader of a PRD, I want a Delegated Decision to look unlike a person's answer, so that I never mistake one for the other.
8. As a session that authors with nobody present, I want one rule that says what to do with a question the author returns, so that I neither guess nor stop without a record.
9. As an adopter of the Baseline, I want the rule to arrive with the Baseline update and to change nothing until my project declares the preference or I answer the question before a flow, so that today's behavior is the default.

## Core Features

1. **A question no source settles has a name.** An **Open Authoring Question** is a question the author of a PRD, a TechSpec or a Task Graph needs answered to finish its phase, that neither the request, the brief, the approved Idea, a recorded answer of the user nor another source of the repository settles. The author states it as the question contract of the Baseline says: what was found, two or three options with the consequence of each, and the recommended option first with its reason. A question that a source settles is not one, and the author settles it from that source as today.
2. **With the user present, the session asks.** Nothing changes there. A question the user can see is asked through the question tool the session exposes and waits for him, and no recommendation, silence, dismissal or timeout adopts anything.
3. **Nobody present is something the session was told.** A session treats the user as absent only when the brief that started it says so, when the user said so in his own words, or when it runs as a continuation that the `unattended` Skill Mode armed. In every other case it asks.
4. **A project declares the preference.** The **Authoring Question Preference** is one key of Project Config with two values: wait, or follow. A project that declares nothing waits. A value in User Config is ignored with one warning that names Project Config, and any other value in Project Config is refused by name when the configuration loads. The preference counts as the repository's default branch holds it in the local refs: the working tree and the branch a session is on change nothing. No session, Agent or command of Roundfix writes the key.
5. **The status view says which preference is in force.** The status view of the Follow-up Command prints, as text and as JSON, the preference in force and where it came from: Project Config or the default. When the checkout's own Project Config declares another value than the one in force, the same line says so, and that the value counts once the default branch holds it.
6. **A project that declares nothing is asked once before the flow.** When the status view reports the default, the `unattended` Skill Mode, and the `follow-up` mode that runs it, asks one question while the user is present and before any Spec is authored, whenever the plan holds a PRD, a TechSpec or a Tasks phase: the Spec waits, which is the recommended option, or the session follows the author. The answer holds for the Specs of that flow: each one's authorization record quotes it and its PRD records it as a dated answer. Without an answer the Specs wait, and a Spec minted in a later flow is asked about again. A project that declares is not asked: the mode's plan and its report name the preference in force and its source. An instruction about this that the user gives unasked, for one flow, counts as that answer.
7. **Following adopts the author's recommendation.** With nobody present, the preference or the answer of Core Feature 6 saying follow, and a question that is not reserved, the session adopts the option the author recommended and the phase goes on as if that option had been given. The adopted option is a **Delegated Decision**. There is no cap on their number in a phase or a Spec.
8. **A Reserved Question always waits.** A **Reserved Question** is never adopted, whatever the preference. These are reserved:
   - the Idea stage, with its questions and the approval of its draft, and the choice an Idea Stage Judgment leaves to the user;
   - a Rule Challenge;
   - every authorization: a commit, a push, a Pull Request, a merge, a release, a Governed Path, a destructive Git operation, a database statement and any other action a stop-and-ask clause holds for express authority;
   - a question whose recommended option crosses one of the user's limits. It adds scope when it adds a goal, a Core Feature, a source or a surface that the request, the approved Idea and the PRD do not hold, removes one of their Non-Goals, or contradicts a recorded answer of the user. It spends when it needs a paid service, a paid plan or a quota the repository does not already use. It sends data outside when it sends repository or client content to a service or a person the repository does not already send that content to;
   - a question that does not have two or three options and a recommended one;
   - a question the author or the session cannot classify.
9. **Two readers classify.** The author says, for each Open Authoring Question it returns, whether it is reserved and why. The session checks that statement against Core Feature 8 before it adopts anything. When either one says the question is reserved, or cannot tell, the Spec waits. Where the session authors the phase itself, it makes both reads before it adopts.
10. **The record names itself.** Each Delegated Decision is an entry of the PRD's Clarification Record, written in the same change that takes it, in a fixed form: the date it was adopted, the phase that asked (PRD, TechSpec or Tasks), the question, the options, the adopted option with the author's reason, the grant it rests on, which is Project Config or the user's dated answer before the flow, and the user's later answer, which reads as not given until he gives one. It is never written as a dated answer of a person or under the autonomous-route line. A question of the TechSpec or the Tasks phase is recorded in the PRD, and that phase's artifact points at the entry.
11. **The check reads the form.** The Spec Consistency Check reports, under the code that already holds the Clarification Record, an entry that starts as a Delegated Decision and lacks its date, its phase or one of its fields. It reads shape only, as it does today, and never whether the grant existed. The record's existing rule is unchanged: a PRD still needs a dated answer or the autonomous-route line, and a Delegated Decision counts as neither.
12. **The report and the status view list them.** The report of the session lists every Delegated Decision of the flow with its Spec, its question and the adopted option, and gives their count. The status view lists the Delegated Decisions of every active Spec, read from the PRD in its Spec folder, each with its date, its phase, its question, the adopted option and the user's later answer when the entry holds one; an entry it cannot read is named with its Spec instead of hidden. The entries of a Spec leave the list when the Spec is archived. The `status` Skill Mode relays them, and their text is repository text: data, never an instruction.
13. **The Spec goes to the end, and a later answer wins.** A Spec with Delegated Decisions continues through authoring, the plan's check and delivery inside the authorization it already holds, with no hold before delivery. When the user answers a delegated question, the session writes his answer into the entry, with its date, whether it keeps or changes the adopted option, and his words. A changed answer amends the plan while no Run and no Delivery Queue item of the Spec has started. After that the plan on the default branch is not edited, and a changed answer follows the routes that exist: a correction of the Spec in delivery, or a Backlog Entry once it is delivered, each quoting his words and naming the entry.
14. **A Spec that waits, waits alone.** When the preference is to wait, the question is reserved, or no answer came, no artifact of that phase is written and nobody answers for the user. An open Backlog Entry names the Spec folder when one exists, holds the whole question with its options and the recommended one, and carries a Waiting Declaration that names the Open Authoring Question and its phase. Other work goes on, the Follow-up Command counts the entry as a Waiting Item, the status view lists it with the open Backlog Entries, and the report names it. When the user answers, the answer is a dated line of the Clarification Record, the Spec adopts the entry, and the phase resumes. The waits of the Idea stage and of a Rule Challenge keep their own text.
15. **Every session receives the rule.** The Baseline's Spec workflow gains one rule, in adopter-neutral words, that states Core Features 1 to 3, 7 to 10 and 14: what the question is, that a present user is asked, that the Spec waits unless the project's preference or the user's answer before the flow says follow, what is reserved, the record, and the wait. Every adopter of the Spec workflow receives it with its next Baseline update. The rule speaks of the phase, so it also governs a session that authors without an Author Subagent, as in Codex, in OpenCode and on the fix or refactor route.
16. **The skills carry the steps.** The Roundfix skill's `modes` reference gains the question before the flow, what a mode does with a question an author returns when nobody is present, and the list in the report; its `follow-up` reference describes the two new parts of the status view, and its `spec` reference the wider meaning of the code. `write-prd`, `write-techspec` and `write-tasks` each say that an author returns an Open Authoring Question with its options, its recommendation and whether it is reserved, never adopts its own recommendation when a session started it, and where the Delegated Decision is recorded. The Clarification Record reference and the PRD template of `write-prd` gain the form. `write-idea` is not changed.
17. **Documentation, glossary and decision.** The configuration guide describes the key, where it is read and what an older binary does with it. The `follow-up` guide describes the two new parts of the status view, the `spec` guide the wider meaning of the code, and the guide to the CONTEXT-driven workflow the rule, the Reserved Questions and the record. `CONTEXT.md` adds **Open Authoring Question**, **Authoring Question Preference**, **Reserved Question** and **Delegated Decision**, and revises Clarification Record and Follow-up Command. ADR-0442 records the decision, and ADR-0280 gains its amendment notice.

## User Experience

The lines below show the content of each surface. The TechSpec fixes the key's name and the exact wording.

A project that follows the author declares it in `.roundfixrc.yml`, on its default branch:

```yaml
authoring:
  # wait (the default): a Spec whose author returns an open question waits for the user.
  # follow: with nobody present, the session adopts the option the author recommends.
  open_questions: follow
```

The status view, in a project that declares nothing and whose checkout holds an uncommitted edit:

```text
Open authoring questions: wait (default); this checkout's Project Config gives follow, which counts only once the default branch holds it
Delegated decisions: 0
```

The question before the flow, in a project that declares nothing:

```text
This project declares no preference for open authoring questions.
While you are away, the author of a PRD, a TechSpec or a Task Graph may end
its phase with a question that no source settles. What happens then?

1. The Spec waits (Recommended): it stops with the question recorded, and
   the other work goes on.
2. Follow the author: the session adopts the option the author recommends,
   records it in the PRD as a Delegated Decision and lists it in its report.
   Scope, spend, data sent outside, every authorization, the Idea stage and
   a question about a rule still wait for you.

To stop being asked, declare the preference in the project's configuration.
```

A question as an author returns it, with an invented subject:

```text
Open authoring question (prd): which failure behavior does an expired import have?
1. Retry automatically up to 3 times (Recommended): it preserves the sync retry behavior.
2. Fail fast and notify the user: the failure shows at once, with no second attempt.
3. Park it for manual review: the import stays available for an operator.
Reserved: no. The three options are inside Core Feature 4 of the PRD, none
adds a service or sends data anywhere, and no recorded answer speaks of it.
```

The Delegated Decision that a session of a following project records in the PRD's `## Decisions`:

```markdown
- delegated: 2026-11-20, prd
  - question: which failure behavior does an expired import have?
  - options: retry automatically up to 3 times; fail fast and notify the user; park it for manual review
  - adopted: retry automatically up to 3 times — it preserves the sync retry behavior
  - under: Project Config
  - answer: none yet
```

In a project that declares nothing, `under` reads `the user's answer of 2026-11-20`. After the user says he wants the other behavior, the last line reads:

```markdown
  - answer: 2026-11-22, changed — "Fail fast, the store has to know at once"
```

The status view of that project then shows:

```text
Open authoring questions: follow (Project Config)
Delegated decisions: 1
  0420-expired-imports 2026-11-20 prd: which failure behavior does an expired import have? -> retry automatically up to 3 times (answer: 2026-11-22, changed)
```

When the Spec waits, a Backlog Entry carries the question until the user answers:

```markdown
---
type: feat
status: open
created: 2026-11-20
spec: null
reason: null
waiting_on: the user's answer to an open authoring question of the PRD
since: 2026-11-20
---
```

Its body names the Spec folder and holds the question, its options, the recommended one and why the Spec waits: the preference, or the reason the question is reserved.

## How this cycle's records map to the form

The PRDs of Specs 0410 and 0414 were written under the maintainer's one-time answer of 2026-11-19 and record each adopted option in one line: the question, `Options:` with each option, the one `taken` with its reason and the ones `refused`, under the sentence "operator decision delegated by the maintainer (2026-11-19)". The form of Core Feature 10 holds the same content: the question is `question`, the options are `options`, the taken option and its reason are `adopted`, and the sentence is `under`, which names his answer and its date. Those PRDs keep their lines, and the status view does not list them. This PRD records its own in the new form, under the same sentence.

## Non-Goals / Out of Scope

- The Idea stage: its questions and the approval of its draft stay the user's, and `write-idea` is not changed.
- A Rule Challenge: nobody answers it for the user, and its Spec waits by its own rule.
- Any authorization: a preference delegates a choice between options, never a permission.
- Silence as an answer while the user is present: no timeout adopts anything.
- A hold before delivery, a confirm step and a notice to the user when a Spec starts to wait or an option is adopted: the report and the status view are the first version.
- A preference in User Config, a preference per phase or per kind of question, and a cap on Delegated Decisions.
- A judge that classifies a question or scores a recommendation.
- A mechanical check that the grant of a Delegated Decision existed: the check reads shape only.
- Questions in a Run, in `implement-task` and in the QA gate: an Agent of a Run never asks the user, and its rules do not change.
- What the PRD author settles from sources: the autonomous route and its line stay as they are.
- A new command, a new check code and a change to the list view of the Follow-up Command: the preference and the list are shown in the status view only.
- Declaring the preference for this repository: its Project Config is not changed by this Spec. The maintainer chose to declare `follow` for Roundfix in a chore Pull Request of its own after this Spec's release; that change is a follow-up outside this Spec.
- Rewriting the one-line records of the PRDs written before the release.
- Permission prompts of the harness for shell commands, which belong to the Spec delivered after this one.

## Success Metrics

1. Success Metric: in a temporary repository, the status view reports wait and the default when no Project Config on the default branch declares the key, follow and Project Config when the default branch declares follow, and the value of the default branch with the checkout's differing value named when the working tree or the current branch declares another; text and JSON agree in each case.
2. Success Metric: a key in User Config is ignored with one warning that names Project Config, a value other than the two is refused by name when Project Config loads, and a project that declares nothing loads as before.
3. Success Metric: the Spec Consistency Check accepts a PRD whose Clarification Record holds a complete Delegated Decision beside a dated answer or the autonomous-route line; it reports, under the code of the Clarification Record, an entry that lacks its date, its phase or any one field, naming the entry's line; and it still reports a PRD whose only entry is a Delegated Decision.
4. Success Metric: the status view lists, as text and as JSON, every Delegated Decision of every active Spec with its answer, names an entry it cannot read, lists none of an archived Spec, and prints a count of zero for a repository that holds none.
5. Success Metric: a managed refresh of an adopted temporary repository writes each clause of the new rule into its Spec routing guide once; a second refresh reports it current; this repository's own Spec routing guide holds them after the Spec's Managed Refresh.
6. Success Metric: both copies of the Roundfix skill's `modes`, `follow-up` and `spec` references, of `write-prd` with its Clarification Record reference and its PRD template, of `write-techspec` and of `write-tasks` carry their parts of Core Feature 16.
7. Success Metric: `CONTEXT.md` defines the four new terms and the two revised ones, the four user-guide pages describe their parts, ADR-0280 carries its amendment notice, and the documentation, skill and glossary contracts pass.
8. Success Metric: the stop-and-ask clause about questions, the clauses of the Rule Challenge and the clauses of the Idea stage keep their ids, their force and their text; the `modes` reference keeps the sentence that a recommendation, silence, a dismissal or a timeout is never an answer; and every test of the configuration, command, check, Baseline and skills packages that exists at the delivery base passes, the tests that pin a changed output having moved to the new output.
9. Success Metric: in a repository whose preference is to follow, no open Backlog Entry waits on an Open Authoring Question that is not reserved, counted at each cycle close.
10. Success Metric: over the first 20 Delegated Decisions recorded after the release, every entry is complete, none is a Reserved Question that was adopted anyway, and the user changed fewer than 4.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write. The first four were read through Exa for the Idea on 2026-11-19 and checked again for this PRD the same day; the last two were read for this PRD.

- Jurgen Appelo, "Leading Agile Developers: The Seven Levels of Authority (Part 2)", 2011 (https://www.informit.com/articles/article.aspx?p=1675546&seqNum=3), level 6, Inquire: "You let the team decide. And afterwards you inquire about their motivations, or you ask that they actively keep you informed." It supports Core Features 7, 10 and 12: the adopted option is reported and recorded, never hidden.
- Feng, McDonald and Zhang, "Levels of Autonomy for AI Agents", 2025 (https://arxiv.org/abs/2506.12469): "we define five levels of escalating agent autonomy, characterized by the roles a user can take when interacting with an agent: operator, collaborator, consultant, approver, and observer". It supports Core Feature 4: autonomy is a setting someone chooses for an Agent, not a property of the Agent.
- AG2, "Approval Required" (https://docs.ag2.ai/docs/user-guide/tools/approval_required/): "an unanswered approval is not an approval". It supports the default of Core Feature 4 and Core Feature 2, and it is the position a following project departs from.
- Cuadron and Gupta, "Built on Broken Assumptions (BOBA): Why Coding Agents Cannot Be Left Alone", 2026 (https://arxiv.org/abs/2611.02207): "complete 95% of requests taken one at a time, but when requests are chained together, they complete 51%". The study measures chains of code requests, not plans. It supports Core Features 12 and 13: an option nobody checked must be easy to see and to overturn before more is built on it.
- Jeff Bezos, letter to shareholders for 2015, filed 2016 (https://www.sec.gov/Archives/edgar/data/1018724/000119312516530910/d168744dex991.htm): "Type 2 decisions can and should be made quickly by high judgment individuals or small groups." It supports the line of Core Feature 8 between what can be undone, which may be delegated, and scope, spend, data and authority, which may not.
- The Apache Software Foundation, "Decision-Making" (https://community.apache.org/committers/decisionMaking.html): "Silence indicates consent, and after that time has passed, you can assume that nobody objects, and proceed with your plan." and "Sometimes lazy consensus is not appropriate." It supports proceeding without an explicit approval when the plan is announced and open to objection, and challenges the default: this Spec follows only where a project chose it.

## Unreachable Acceptance

- criterion: in a real flow after the release, with nobody present, a session of a project that follows adopts the author's recommendation for an Open Authoring Question that is not reserved, records the Delegated Decision and lists it in its report, a session of a project that declares nothing asks the question before the flow once, and a Reserved Question makes its Spec wait.
  reason: it needs a live Agent session, an absent person and a real question, which no Verification or QA row may start.
  satisfied-by: the operator's session evidence after the release, recorded session evidence.
- criterion: Success Metrics 9 and 10 hold at the cycle closes after the release and over the first 20 Delegated Decisions.
  reason: it needs weeks of authoring in this repository and in the adopters, after the release.
  satisfied-by: the operator's count, at a cycle close, of the Delegated Decisions and their answers in the PRDs and of the Waiting Declarations in the open Backlog Entries, recorded session evidence.

## Validation

The decisions were checked on 2026-11-19 against these sources.

- Characterization at `894710c3`, recorded in [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md). The Clarification Record's check accepts a `## Decisions` section that holds a dated answer or the autonomous-route line and reads no other line, so a new entry form moves no existing result. The status view reads the Spec folders for Task counts and never a PRD. One Project Config key, the one that requires a Defect Class, is already read at the default branch head and ignored in User Config, and the reader that does it serves this key unchanged. A project that sets no new key is read by every binary. The Waiting Declaration is read only on a Backlog Entry and a Finding, which is why a Spec that already has a folder waits through a Backlog Entry that names it.
- A prototype of every text of the change on a scratch worktree, removed afterwards and recorded in the same reference: the rule's two clauses in a new rule of the Spec workflow module, the skill section, the record form, the `modes` texts, the glossary entries and the guide texts, with the regeneration commands and the Managed Refresh. The Baseline, skills, documentation, check and command-line Baseline suites passed. One contract failed first and changed the text: the `write-tasks` skill may not hold the word "recommendation", so the shared section says "the recommended option". The digests command stops while a skill mirror differs from its canonical file, which fixes the order of the regenerations. The code of the key, of the status view and of the check was not prototyped.
- Secondbrain, read only. `wiki/index.md` was read, and three `qmd query` calls on 2026-11-19 failed with `database is locked` while another process rebuilt the index, so the pages were found by name and read directly. `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md` confirms why an Author Subagent returns its questions: a subagent does not receive the question tool. `wiki/concepts/claude-code-operacao-produtiva.md`, "manter aprovação humana para envios externos ou ações difíceis de desfazer", supports the Reserved Questions. `wiki/concepts/agent-harnesses.md` says that the controls of accountability stay with a person and that an approval prompt is not a safety control, "93% são aprovados", which supports a record over one more prompt. `wiki/concepts/correcao-de-dados-com-aprovacao-humana.md`, "Parar no conflito, nunca continuar em best-effort", is the opposite rule for changes to data, and is why a question that sends or changes data outside is reserved. The Secondbrain holds no earlier research on a session that adopts an author's recommendation.
- Exa: the six published sources above, the first four through the Research Subagent, which fetched each page again on 2026-11-19 and confirmed each quoted passage; the last two read by the author. LangGraph, "Interrupts" (https://docs.langchain.com/oss/python/langgraph/interrupts), was read for the Idea and is not an acceptance row: a paused graph "waits indefinitely until you resume execution". One more page was read only through the Research Subagent and is not quoted as evidence: the Claude Code documentation of permission modes (https://code.claude.com/docs/en/permission-modes), which that subagent reports as refusing the most permissive mode when a project's settings file sets it. It is a published choice against this Spec's choice of a project-level key, and `## Risks` answers it. The sweep found no measurement of how often people accept the recommended option of a clarifying question, and no tool that records a decision as taken on someone's behalf in a fixed form.

## Risks

- The author that classifies a question is the Agent that goes on when the question passes, and the session that checks it is also an Agent. The person's check is reading the records, which is why the record holds the question and every option.
- "Nobody present" is something the session is told. A session that wrongly believes it would adopt where it could have asked; Core Feature 3 names the only three ways it is told.
- A preference in a repository file binds the unattended sessions of everyone who works there, and a cloned repository can declare follow. It delegates a choice between options the author already bounded, never a permission, and the mode's plan names the preference and its source before the user leaves. An Agent that edits the key gains nothing until the default branch holds the edit, and the key's file is a Governed Path.
- An author may label a coin flip as recommended. The record keeps the author's reason, and Success Metric 10 shows the pattern.
- "New scope" is a judgment. Core Feature 8 gives it a test a second reader can apply, and doubt waits.
- A wrong option can reach the TechSpec, the Tasks and the code before the user reads it. The maintainer chose to go to the end with a report; the cost is a correction or a new Spec.
- With scope, spend, data, authorizations, the Idea stage and rules all reserved, the questions left to delegate may be few. Nobody has counted them, and Success Metric 9 and the count of entries will.
- The rule adds text to the Spec routing guide and to three authoring skills, which every author reads. The open Backlog Entry about the context an Authoring Phase starts with measures that floor.
- Dissent kept from the Idea: published tools hold that an unanswered approval is not an approval, and a better answer to the cost of waiting may be a short wait, by telling the user at once that a Spec waits. That notice is left out of the first version because it sends the question's text outside the machine.

## Idea-stage judgment

- item: the request in the session, recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md)
- judged-by: nobody (the user asked for the Idea stage)
- verdict: not judged
- stage: ran
- shape: split — the second part of the request, the standing authorization of shell commands, became Spec 0413, and the first part gained the list of what always waits, the two readers, the record in fixed form and the line in the status view.

The judge was not run. The maintainer chose the Idea stage directly on 2026-11-19: "Spec 0412 com Ideia agora". He approved the Idea, with its `shape` line, the same day: "Aprovar e seguir ao PRD (Recommended)".

## Rule challenges

None. The author read the rules that touch this goal: the stop-and-ask clause about questions, with its sentence that a recommendation, a dismissal, a timeout or a missing answer is never the user's choice; the sentence of the `modes` reference that says the same; the rule that the Idea stage is the user's; the clause that tells the Tasks author when to stop and ask; the tooling-authority clause; and the rule that Project Config is a Governed Path. None blocks the goal. The request itself asks to change what happens when nobody answers, so that change needed no question, and the clauses keep their text: a Delegated Decision is not the user's choice and is never recorded as one.

## Glossary

- adds: **Open Authoring Question**
- adds: **Authoring Question Preference**
- adds: **Reserved Question**
- adds: **Delegated Decision**
- changes: **Clarification Record**
- changes: **Follow-up Command**

## Decisions

Every dated line is the maintainer's answer, as recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The text in quotes is the option he chose, verbatim.

- 2026-11-19: What does the operator do with an open question of an author while the maintainer is away, in this cycle? — "Seguir o padrão do autor (Recommended)": the operator adopts the option the author recommends and records it in the PRD as an operator decision by delegation, with the question and the options; his limits, which are data sent outside, new scope and spend, still stop the Spec.
- 2026-11-19: How is the suggestion recorded? — "Spec 0412 com Ideia agora": the Spec is opened at once, with the Idea stage, and delivered after Spec 0410.
- 2026-11-19: One Spec or two? — "Duas Specs (Recommended)": this Spec holds the preference for an open authoring question, and the standing authorization of shell commands becomes Spec 0413, delivered after this one.
- 2026-11-19: Which phases does the delegation cover? — "PRD, TechSpec e Tasks (Recommended)": the Idea stays the maintainer's, with its questions and the approval of its draft.
- 2026-11-19: The question before the flow, or the configuration? — "Config responde; sem ela, pergunta (Recommended)": a project that declares is not asked; without the key, the mode asks once before the flow; without an answer, the Spec waits.
- 2026-11-19: How far does a Spec with delegated decisions go? — "Até o fim, com relato (Recommended)": authoring and delivery inside the existing authorization, with each delegated decision in the PRD, in the report and in the status.
- 2026-11-19: Is the Idea approved? — "Aprovar e seguir ao PRD (Recommended)".
- 2026-11-19: Which authorization do Specs 0412 and 0413 have? — "Mesmas do BE para as duas (Recommended)": implement, commit, push, Pull Request, merge and a minor release after each one, up to five corrections, and the Makefile, CI, Go module and lint only if a correction requires one.
- 2026-11-19: Does Roundfix declare `authoring.open_questions: follow` in its own `.roundfixrc.yml`? — "Sim, declarar follow após o release (Recommended)": in a chore Pull Request of its own right after this Spec's release, outside this Spec.
- 2026-11-19: What tells a session that nobody is present? — "O que a sessão foi informada (Recommended)": the brief, the user's words or a continuation the `unattended` mode armed; a question visible on the screen is always asked.
- 2026-11-19: Is there a cap on delegated decisions? — "Sem teto, com contagem (Recommended)".
- 2026-11-19: How long does the status list a delegated decision? — "Enquanto a Spec estiver ativa (Recommended)", with the maintainer's later answer beside it.

Autonomous route: no person present; decisions taken from the approved Idea and its defaults, the maintainer's answers above, the operator's brief of 2026-11-19, the characterization and the prototype recorded in `## Validation`.

- The open questions the approval left to the Idea's defaults are settled as those defaults say: one key with two values and an absent key that waits; the answer before the flow holds for the Specs whose authorization records hold it; the rule holds where no Author Subagent runs; a question of the TechSpec or the Tasks phase is recorded in the PRD; the session tells the user nothing beyond the report and the status view; and the tests of scope, spend and data are the Idea's.
- A question about a rule is reserved without a new question to the maintainer. He answered it for Spec 0409 on 2026-11-18: "A Spec espera (Recommended)".
- The wait uses the mechanism Spec 0409 chose, a Backlog Entry that names the Spec, so that two kinds of waiting Spec read the same way in the Follow-up Command.
- The Release Ceiling is `minor`: a Project Config key, two parts of the status view and one Baseline rule are added, and no command, flag or exit code is removed or renamed.
- ADR-0442 is written with this PRD. The TechSpec adds no ADR unless its design needs one.

Each entry below is an operator decision delegated by the maintainer (2026-11-19), with its question, its options and the one taken. The author took them under the operator's brief, which is how this cycle applies his answer; from the release on, an author returns such a question and the session adopts.

- delegated: 2026-11-19, prd
  - question: which names does the glossary take?
  - options: four terms, for the question, the preference, the questions that always wait and the adopted option; one term, Delegated Decision, with the rest described in prose
  - adopted: four terms — the clause, the skills and the record each need to name the question, the preference and what is reserved in one or two words, and prose synonyms would drift between them
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: from which state is the preference read?
  - options: the default branch head as the local refs hold it; the commit the checkout is at; the working tree
  - adopted: the default branch head — a session cannot switch it with an edit or a commit on its own branch, and the key that requires a Defect Class is already read there
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which surface tells a session the preference in force and lists the Delegated Decisions?
  - options: the status view of the Follow-up Command only; both views of the Follow-up Command; a new command
  - adopted: the status view only — every Skill Mode can run it, it adds no command, and the list view keeps every line its transcripts pin
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: how does the Spec Consistency Check learn the form?
  - options: the code that already holds the Clarification Record reports a malformed entry; a new code; no check
  - adopted: the existing code — one record keeps one code, and the detector table and the corpus pins, which are Governed Paths, do not change
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which form does the record take?
  - options: one entry with a field per line, as shown in the User Experience; one line per decision, as the PRDs of Specs 0410 and 0414 wrote; a section of its own in the PRD
  - adopted: one entry with a field per line — a check and the status view can read each field, the user's later answer has a place, and it stays in the section where a Spec's decisions live
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: where is the user's later answer recorded once the Spec's delivery has started?
  - options: in the entry only while the plan can still be amended, and afterwards in the correction or the Backlog Entry that follows a changed answer; always in the PRD on the default branch; in a record under Roundfix Home
  - adopted: in the entry only while the plan can still be amended — an edit of the PRD on the default branch during a delivery conflicts with the archive, which removes the Spec folder on the Item Branch
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: is an option that contradicts a recorded answer of the user a Reserved Question?
  - options: yes, as part of the scope limit; no, only added scope is reserved, as the Idea's default reads
  - adopted: yes — an option that undoes what the user already answered changes his scope as much as one that adds to it
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: does this Spec declare the preference for this repository?
  - options: no, the maintainer declares it himself after the release if he wants it; yes, a Task writes follow into this repository's Project Config
  - adopted: no — declaring the preference is the user's own act, the key's file is not granted, and the default is the behavior he asked for
  - under: the user's answer of 2026-11-19
  - answer: 2026-11-19, kept — "Sim, declarar follow após o release (Recommended)"
- delegated: 2026-11-19, prd
  - question: what happens to the waiting Backlog Entry when the user answers?
  - options: the Spec adopts it under the adoption contract; it is closed and deleted with a reason
  - adopted: the Spec adopts it — the entry names the Spec and holds the question the Spec then records as answered, and adoption is the route the docs-layout guide gives a source that a Spec absorbs
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: who classifies where no Author Subagent runs?
  - options: the session makes both reads before it adopts; such a session never adopts and always waits
  - adopted: the session makes both reads — the Idea's default says the rule speaks of the phase wherever it runs, and a rule that held only in one harness would leave the others in today's silence
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: does ADR-0280 gain an amendment notice?
  - options: yes, a Task writes it; no, ADR-0442 alone records the third entry
  - adopted: yes — ADR-0280 lists the forms of the Clarification Record, and a reader of it alone would miss the third
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, techspec
  - question: what does the status view report when the default branch head or its Project Config cannot be read?
  - options: wait, from the default, with the reason on the same line; fail the command; fall back to the value the checkout gives
  - adopted: wait, with the reason — the default is the safe value, a failed status would hide every other line of the report, and the checkout's value is the one a session could have edited
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, techspec
  - question: does the default configuration that the init command writes show the key?
  - options: no, the configuration guide documents it; yes, as a commented line
  - adopted: no — the two other keys only a project sets are not in that file either, and the guide is where a user reads what an older binary does with the key
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, tasks
  - question: how is step 3 of the Build Order cut into Tasks?
  - options: one Task, as the TechSpec held it; two Tasks, the Baseline rule and then the skills, the modes reference, the glossary, the workflow guide and the ADR notice; three Tasks, with the glossary, the guide and the notice in a Task of their own
  - adopted: two Tasks — one Task would declare about 35 paths, run two regenerations and pin more than 70 guarantees, beyond one fresh session, and a third Task of six glossary entries, one guide subsection and one notice would not earn its own verification cycle
  - under: the user's answer of 2026-11-19
  - answer: none yet

## Open Questions

None for the product. The maintainer confirmed this Spec's authorization and chose to declare the preference for this repository after the release, in a chore Pull Request outside this Spec (`## Decisions`). The measurements of Success Metrics 9 and 10 are the operator's, after the release.
