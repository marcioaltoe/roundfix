---
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
status: active
created: 2026-11-18
surfaces: [cli, docs]
type: feat
---

# The Idea and the PRD question the rules that block the goal

An author of an Idea or a PRD that finds a recorded decision in the way of the best solution has two routes today, and both are silent: drop the solution, or keep it and leave the rule for later. The maintainer learns that a rule decided the scope one, two or three Specs afterwards. After this Spec, in brainstorming, the Idea and the PRD the author asks the user one question about each rule that blocks the goal, with three answers: keep, change or retire. The rule binds until the user answers, the artifact records the question and the answer in a fixed form, the Spec that asked carries the change, and a rule the user keeps says why it stays. The Baseline delivers the rule to every repository that adopts it, and the TechSpec author asks the same question before a design contradicts an ADR.

## Sources

1. [`references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md`](references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md), the adopted Backlog Entry with the maintainer's request of 2026-11-18, verbatim. The request is the maintainer's own; no ai-memory message is adopted and nobody asked for a notice.
2. [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md), the routing, the Idea-stage judgment, the six answers of the Idea stage and the approval of the Idea, recorded by the operator.
3. [`references/2026-11-18-rules-named-in-recent-plans.md`](references/2026-11-18-rules-named-in-recent-plans.md), the count of the recorded decisions that 41 plans name, recorded by the author of the Idea.
4. [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), recorded by the author of this PRD: the prototype of the whole change, the files its regeneration rewrites, the tests it passed and the Governed Paths.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from: its features F1 to F10 are this Spec. This project's ai-memory inbox was listed on 2026-11-19 and held no pending message, and `docs/findings/` holds no Finding. Four Backlog Entries are open at the authoring base. Two are the sources of Spec 0410. One reports three unreproduced symptoms from another team. One is about the size of the context an Authoring Phase starts with; it waits on a measurement and shares no component with this Spec, although this Spec adds text to what an author reads (`## Risks`). None is adopted here.

## Prerequisites

None of this Spec's behavior depends on another Spec's code. It was authored against `348962d2`, where Specs 0407 and 0408 are delivered, and is delivered after Spec 0411 and before Spec 0410, in the order the maintainer set. These files are shared, and each Task regenerates the shared records on the tree it runs on, with no version number written by hand:

- Spec 0411, planned and not yet delivered at the authoring base, changes the Roundfix skill's `modes` reference, that skill's version fields, the owned-skill version record, `CONTEXT.md` and the Behavior Surface Record. This Spec changes the same files, inside its own sentences.
- Spec 0410, delivered after this one, changes the Baseline modules `core` and `context-workflow`, which this Spec does not touch, and shares the module version record, the derived Baseline files, this repository's Setup Manifest, the Roundfix skill's version and `modes` reference, the owned-skill version record and `CONTEXT.md`. This Spec keeps its Baseline text in the `spec-workflow` module and changes neither the docs-layout guide nor the agent-instructions guide, so Spec 0410 rebases over it without meeting its sentences.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity; the new Baseline rule and its two clauses take ids in the form the Spec workflow module already uses, and a record names its rule by the ADR number or the path the repository already gives it. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the proposal to the Baseline's owner is a message an Agent session sends through ai-memory under the existing cross-project rule, and Roundfix sends nothing. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0436 (new) records the question and what follows it (ADR-0436: "The rule binds until the user answers."), and ADR-0437 (new) the line a kept rule gains (ADR-0437: "A rule keeps every Kept Line, the newest last."). The question uses the question contract as it is (ADR-0320: "Up to three, independent, in one pending call."). The wait repeats the Idea stage's wait (ADR-0403: "no user is present, no PRD is written: the feature stays a Backlog Entry that") and uses the Waiting Declaration without changing it (ADR-0353: "makes it a Waiting Item in a fifth class"). The proposal to the Baseline's owner is an ai-memory message under the existing cross-project rule, and it expects no reply (`docs/agents/agent-instructions.md`: "A message expects no reply unless the sending project's own work depends on the answer."). The new clauses are written for every adopter (ADR-0186: "Baseline guidance states what the product does, in adopter-neutral words"; ADR-0222: "A Baseline guide says only what holds for the repository that reads it") and state the rule once (ADR-0257: "One statement per rule."). The Kept Line of a rule the Baseline owns goes in the repository's own rules guide (ADR-0411: "A repository rule outside the specific-repository guide is named, never moved"). The skills and the guide describe the behavior in the same delivery (ADR-0256: "A release waits for the skills that describe a changed surface"), and the authorization record carries the Release Ceiling (ADR-0349: "A Spec's Release Ceiling answers the release plan's approval"). ADR-0203 is listed to record why it does not apply: no hand-written rule of a repository is promoted, the request itself puts the rule in the Baseline, and its evidence from outside Roundfix is the maintainer's report from the adopters and a published measurement (ADR-0203: "The Baseline grows only with evidence from outside Roundfix."). ADR-0278 is listed to record that this Spec retires nothing: a retire answer given in a later Spec follows the retirement rule as the docs-layout guide then states it (ADR-0278: "A retirement deletes."). ADR-0258, ADR-0338, ADR-0342, ADR-0383, ADR-0390, ADR-0404, ADR-0410, ADR-0412 and ADR-0432 are listed only to record why they do not apply; the TechSpec relates each one. No accepted ADR is amended or superseded. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18 for this Spec, quoted in the authorization record, which keeps the grant of every owned skill and of the Governed Paths each Spec declares, Baseline modules included. It grants no Makefile, workflow, Go module, lint or formatter change except for a correction that requires one, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0409-the-idea-and-prd-question-rules-that-block-the-goal/_authorization.md`; bounded files: `.agents/skills/brainstorming/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-idea/references/idea-template.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/brainstorming/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A solution is narrowed by a rule only when the user chose to keep that rule: no scope is given up, and no rule is contradicted, in silence.
- An Idea and a feature PRD say whether a rule blocked the goal, so a reader can tell a boundary the author chose from a goal the author gave up.
- A rule that the user changes or retires changes in the same delivery as the work that depended on it.
- A rule that the user keeps says why it stays, and is not asked about again without a new fact.
- Every repository that adopts the Baseline receives the rule with its next update, and nothing changes in the stages that build: there a rule stays a guard-rail.

## User Stories

1. As the maintainer, I want the author of an Idea or a PRD to ask me whether a rule that blocks the goal is kept, changed or retired, so that I decide the scope and no rule decides it for me.
2. As the maintainer, I want the question to say what the rule imposes, what it prevents for this goal and why the author thinks it exists, so that I can answer without reading the rule's history.
3. As the maintainer, I want the rule to bind until I answer, so that a question never weakens a decision I confirmed.
4. As a reader of a PRD, I want one section that holds each question and its answer, or says that no rule blocked the goal, so that I know why the scope ends where it does.
5. As the maintainer, I want the Spec that asked to change the rule together with the work that depended on it, so that the rule and the work are verified and delivered together.
6. As a later author, I want a kept rule to carry a line with what was asked and why it stayed, so that I do not put the same question again.
7. As the maintainer away from the session, I want a Spec that needs my answer to wait while the other work goes on, so that nothing is built on an answer I did not give.
8. As the maintainer of a repository that adopts the Baseline, I want my decision about a rule the Baseline owns to reach the Baseline's owner as a proposal, so that the rule is reconsidered where it can change.
9. As the maintainer, I want the TechSpec author to ask me before a design contradicts an ADR I have not decided about, so that an ADR is never superseded on the author's word alone.
10. As an adopter, I want the rule to arrive with the Baseline update and to cost me a question only when a rule blocks a goal, so that ordinary authoring is unchanged.

## Core Features

1. **The Baseline forbids both silences.** A Normative Clause of the Spec workflow says that in brainstorming, the Idea and the PRD an author never drops or narrows a solution in silence because a rule forbids it, and never contradicts a rule in silence: it puts a **Rule Challenge** to the user. Every adopter of the Spec workflow receives the clause with its next Baseline update, and the clause governs a session that authors through a Vendored Skill.
2. **A rule is a decision the repository recorded.** An accepted ADR, a Normative Clause of a guide under `docs/agents/`, a definition in the domain context, or a document under `docs/references/` that the repository treats as binding. A user-level instruction file, an ai-memory page, a harness's own rules and a default of a Vendored Skill are not rules of the repository.
3. **The trigger is a blocked goal.** A rule blocks the goal when the solution the author would recommend, or a part the goal needs, cannot exist without contradicting the rule, or when the rule stops the author from stating the problem. Extra work, taste and friction are not a block. No Rule Challenge is asked about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records; the artifact names that request or answer where it records the change.
4. **The question.** One question per rule, through the structured question tool the session exposes, as the question contract of the Baseline says. It names the rule and where it lives, what the rule imposes, up to three things it prevents for this goal, and the author's reading of why the rule exists and whether it looks out of date or badly written. The answers are keep, change and retire, each with its consequence and the recommended one first. A recommendation, silence, a dismissal or a timeout is never an answer.
5. **At most three at a time.** An Authoring Phase or a brainstorming session asks at most three Rule Challenges. When more than three rules block the goal, the author asks about the three that prevent the most and lists the others, which keep binding, in the artifact's open questions.
6. **The rule binds until the answer.** Asking changes nothing. The author plans on the user's answer only after the user gives it, and the rule's text stays as it is until a Task of the Spec changes it.
7. **The record.** The Idea and the PRD carry a `## Rule challenges` section with one record per question in a fixed form: the rule and where it lives, what it imposes, what it prevents, the date and the phase it was asked in, the answer with the user's words and its date, and what the Spec does with it. An Idea and a feature PRD always carry the section, which reads `None.` when no rule blocked the goal; the minimal PRD of a fix or a refactor carries it only when a question was asked. The PRD copies the Idea's records, and a question the TechSpec author asks is recorded in the PRD.
8. **The Spec carries the change.** On change or retire, the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it. An ADR changes through the lifecycle the docs-layout guide defines when the Task runs: an amendment or a successor for change, the retirement rule for retire. This Spec adds no ADR status.
9. **A kept rule says why it stays.** On keep, the solution stays inside the rule, and a Task of the Spec adds a **Kept Line** after the rule's text: the date of the answer, what was asked and why the rule stays. The line names no Spec, a rule keeps every line with the newest last, and when no Spec follows the question, the session that asked adds the line. A later author asks about that rule again only with a fact the earlier question did not have, and the record names the fact.
10. **A Spec waits when nobody can answer.** With no person present, no artifact of that phase is written and nobody answers for the user. Only that Spec waits: the feature stays or becomes a Backlog Entry that holds the whole question, names the Spec folder when one exists and carries a Waiting Declaration that names the Rule Challenge, and the other work goes on. An Author Subagent returns the question to the session that started it. The `unattended` and `follow-up` Skill Modes of the Roundfix skill ask the question when the user is there and never answer it.
11. **A rule the Baseline owns becomes a proposal.** A clause inside a Managed Region changes only in the Baseline. On change or retire, a repository that does not own the Baseline sends the user's decision to the Baseline's owner as one ai-memory message, obeys the rule until the Baseline changes, goes on inside the rule and records the message and the scope it left out. On keep, the Kept Line goes in the repository's own rules guide, under a heading of its own, naming the rule.
12. **The TechSpec asks before it contradicts an ADR.** When a design takes what an accepted ADR rejected, or otherwise contradicts its decision, and neither the request nor a recorded answer of the user decides it, the TechSpec author puts the Rule Challenge before it supersedes or amends the ADR, records it in the PRD and points at the record from the TechSpec's line about that ADR. An amendment the user already decided needs no question.
13. **Implementation keeps its guard-rails.** In a Run, in a Task and in the QA gate a rule is a guard-rail: an Agent that meets a conflict there flags it, as the domain guide already says, and asks no Rule Challenge. The existing clauses keep their ids, force and text: the domain clause on conflicts, the question contract, the route clauses and the sentence of the Instruction Hierarchy.
14. **The skills carry the step.** `brainstorming`, `write-idea`, `write-prd` and `write-techspec` each gain a section that states the trigger, the question, the three answers, the cap, the binding rule, the wait and what follows each answer, with what is particular to their stage. The Idea template and the PRD template gain the section of Core Feature 7. `write-idea` lists the rules that touch the idea and never drops an alternative of its Opportunity Scan only because a rule forbids it. `write-prd` reads each rule as an obligation and as a candidate before it writes the Project Constraints. `brainstorming` asks instead of steering the user away from a direction.
15. **Documentation and glossary.** The user guide describes the question, the three answers and the record. `CONTEXT.md` adds **Rule Challenge** and **Kept Line**. ADR-0436 and ADR-0437 record the decisions.

## User Experience

The lines below show the content of each surface. The TechSpec fixes the exact wording.

The question, with an invented rule that is not one of this repository:

```text
Rule: ADR-NNNN, "every report is produced at the close of the day" (docs/adr/NNNN-….md)
It imposes: no report is computed while the day is open.
It prevents, for this PRD's goal (a stock-out alert the store can act on):
  1. computing the alert when the sale happens;
  2. telling the store before the shift ends.
The author's reading: the ADR is from 2025 and gives a reason that the newer
sales feed removed. It looks out of date, not wrong when it was made.

1. Change the rule (Recommended): the PRD requires the amendment of ADR-NNNN,
   and this Spec delivers it with the alert.
2. Keep the rule: the PRD goes on with an alert at the close of the day, and
   ADR-NNNN gains a dated line with this question and your reason.
3. Retire the rule: ADR-NNNN leaves the tree under the retirement rule, and
   Git keeps its text.
```

The record of that question, in the Idea or the PRD:

```markdown
## Rule challenges

- rule: ADR-NNNN, `docs/adr/NNNN-….md`
- imposes: no report is computed while the day is open
- prevents: computing the alert when the sale happens; telling the store before the shift ends
- asked: 2026-11-20, prd
- answer: change — "Change the rule (Recommended)", 2026-11-20
- effect: Core Feature 4 requires the amendment of ADR-NNNN; a Task delivers it with the alert
```

Had the user kept the rule, the record's `answer` would read `keep`, its `effect` would name the scope the Spec keeps inside the rule and the line a Task adds, and ADR-NNNN would gain, after its decision text:

```text
Kept (2026-11-20): asked whether the stock-out alert may be computed while the day is open; it stays because the close of the day is when the fiscal totals are final.
```

With nobody present, no PRD is written, and a Backlog Entry carries the question until the user answers:

```markdown
---
type: feat
status: open
created: 2026-11-20
spec: null
reason: null
waiting_on: the user's answer to the Rule Challenge about ADR-NNNN
since: 2026-11-20
---
```

Its body names the Spec folder and holds the record with `answer: waiting` and the three options. The Follow-up Command counts it as a Waiting Item, and the status view lists it with the open Backlog Entries.

## Non-Goals / Out of Scope

- A mechanical check of the section or of the question: the maintainer chose the first version without one, and `roundfix spec check` reads none of it.
- A judgment of a PRD's exclusions by the judge.
- A review of every rule on a schedule: the trigger is a blocked goal, never the age of a rule.
- Questions during implementation: in a Run, in `implement-task` and in the QA gate a rule is a guard-rail.
- New ADR statuses, a new ADR layout or a change to the retirement rule: the lifecycle is used as it is, and the documentation rules belong to Spec 0410.
- Rewriting existing rules to add their reasons: a rule gains a Kept Line only when it is asked about and kept.
- Rules that live outside the repository: user-level instruction files, ai-memory pages and a harness's own rules.
- Edits to Vendored Skills: `grilling` and `domain-modeling` stay as they are, and the Baseline clause governs a session that authors through them.
- Changes to `write-tasks`, `qa-gate`, `archive-spec` or `implement-task`: a rule change the PRD requires reaches a Task through the coverage those skills already hold.
- A change to the Waiting Declaration, the Follow-up Command or any other command: no command surface changes.
- A review of delivered Specs for scope given up in the past: the record starts with the plans authored after the release.
- Grouping proposals about the same Baseline clause from several adopters into one decision.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository writes the two clauses of the question into its Spec routing guide, each once, the first with the force stop-and-ask and the second mandatory; a second refresh reports it current; this repository's own Spec routing guide holds both after the Spec's Managed Refresh.
2. Success Metric: both copies of `brainstorming`, `write-idea`, `write-prd` and `write-techspec` carry the section with the trigger, the three answers, the cap of three, the binding rule, the wait and what follows each answer; both copies of the Idea template and of the PRD template carry the `## Rule challenges` section with its six fields and `new fact`; both copies of the Roundfix skill's `modes` reference say that the mode asks a Rule Challenge and never answers it.
3. Success Metric: `CONTEXT.md` defines **Rule Challenge** and **Kept Line**, the user guide describes the question, the three answers and the record, and the documentation, skill and glossary contracts pass.
4. Success Metric: the clauses that exist at the delivery base keep their ids, their force and their text, the Instruction Hierarchy keeps its sentence, and every test of the Baseline package and of the skills package that exists at the delivery base passes, the closed list of clause forces having gained the two new clauses.
5. Success Metric: no Idea or PRD authored in this repository after the release lacks the `## Rule challenges` section, read over the first 20 plans.
6. Success Metric: in those 20 plans no Authoring Phase asks more than three Rule Challenges, at least 18 Specs ask at most one, and no question is about a rule that carries a Kept Line without naming a new fact.
7. Success Metric: in the 30 days after a repository applies the update, no Backlog Entry, Finding or message of that repository reports that a delivered Spec left something out because of a rule nobody asked about.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write. The first two were read through Exa for the Idea on 2026-11-18 and read again for this PRD on 2026-11-19; the others were read for this PRD on 2026-11-19.

- Xu and Wu, "Stale Instruction Files Override the Code in Coding Agents: A One-Sentence Fix", 2026 (https://arxiv.org/abs/2610.30571), 320 runs over 40 projects: "Claude Code pointed out the conflict in 79 of 80 runs and still followed the file." One sentence that gave the code precedence "raised the share of choices that follow the code from 4.2% to 91.4%". The study measures code conventions, not product goals. It supports the two premises of Core Feature 1: an Agent that sees a conflict with an instruction complies with the instruction, and one sentence in the instructions changes what it does.
- Zhang and others, "Ask Less, Ask Better: Allocating User Attention in Proactive LLM Agents with Interruption Budgets", 2026 (https://arxiv.org/abs/2610.24411): "limiting agents to one to three questions reduces user-attention cost by 36–73% while keeping 86–94% of the task reward they earn with unlimited questions", and "reviewing each question against the memory of users' earlier answers cuts attention cost by another 13%". It supports the cap of Core Feature 5 and the reading of the Kept Line before asking again.
- Fowler, "Architecture Decision Record", 2026 (https://martinfowler.com/bliki/ArchitectureDecisionRecord.html): "Once an ADR is accepted, it should never be reopened or changed - instead it should be superseded." It supports changing an ADR through a successor and challenges the Kept Line, which ADR-0437 answers.
- Shadewright, "Amending accepted ADRs", a published decision record of that project, 2026 (https://shadewright.io/adr/0021-amending-accepted-adrs/): "Every amendment gets a dated entry in an `## Amendments` section at the end of the ADR: the date, what changed, and why". It supports a dated line on an accepted ADR whose decision stands.

## Unreachable Acceptance

- criterion: in a real authoring session after the release, an author that meets a rule that blocks the goal puts the Rule Challenge to the user through the session's question tool, plans on the answer only after it is given, and writes the record.
  reason: it needs a live Agent session, a person who answers and a real rule in the way, which no Verification or QA row may start.
  satisfied-by: the operator's session evidence after the release, recorded session evidence.
- criterion: Success Metrics 5, 6 and 7 hold over the first 20 plans authored after the release and over the 30 days after a repository applies the update.
  reason: it needs weeks of authoring in this repository and in the adopters, after the release.
  satisfied-by: the count of the measurement in `references/2026-11-18-rules-named-in-recent-plans.md` repeated by the operator at a cycle close, recorded session evidence.

## Validation

The decisions were checked on 2026-11-19 against these sources.

- Characterization at `348962d2`. The Instruction Hierarchy says a narrower guide "cannot weaken a universal Normative Clause or confirmed project decision". The domain guide says "flag conflicts instead of silently overriding repository decisions", and no clause speaks of a goal given up. `write-idea` and `brainstorming` have no step about a rule; `write-prd` lists ADRs as obligations; `write-techspec` says "respect prior decisions or challenge them explicitly — never contradict one silently" and asks nobody. The Waiting Declaration is read only on a Backlog Entry and a Finding, which is why a Spec that already has a folder waits through a Backlog Entry that names it. `roundfix spec check` holds no closed list of the sections of an Idea or a PRD, and its coverage reads only the User Stories, the Core Features and the Success Metrics, so the new section moves no check.
- A prototype of the whole change on a scratch worktree of `348962d2`, removed afterwards and recorded in [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md). Two new clauses added to the route rule of the Spec workflow were refused by the catalog, because the Standard TypeScript profile requires that rule and a new clause in it needs a Source Baseline row that no command creates. The same two clauses in a new rule of the same module passed: the module record command, `make baseline-digests` and the managed refresh rewrote the module, the module version record, the Standard TypeScript profile and its formatter golden of the Spec routing guide, three catalog files, four plan goldens and this repository's Spec routing guide and Setup Manifest, and a second refresh reported no change. With the two ids added to the closed list of clause forces, the Baseline package passed, and so did the CLI's Baseline tests. The skill, template and `modes` edits with `make skills-sync` and the version record passed the skills package without changing a pinned sentence. The two glossary entries passed the glossary contracts. The user guide edit moved the Behavior Surface Record, and the documentation contracts passed after its record command.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for a question about keeping, changing or retiring a rule, for a kept decision that records its reason, and for a new Baseline clause with a skill step. The queries returned this repository's own mirrored ADRs, of which ADR-0203 and ADR-0222 bear on a new clause: they are why the clauses are adopter-neutral and why the request's own evidence is named. The pages the Idea read stand: `wiki/concepts/politica-de-retirada-e-poda.md` ("Retirar com razão, nunca em silêncio") for a trigger that is an event and an answer that is the user's, and `wiki/concepts/perguntas-estruturadas-em-agentes-cli.md` for three answers and for the Author Subagent that returns the question.
- Exa: the four published sources above. Two more were read and are not acceptance rows. AWS Prescriptive Guidance, "Architectural decision record process" (https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html): "When the team accepts an ADR, it becomes immutable." It is the position ADR-0437 departs from. LangChain, "Human-in-the-loop" (https://docs.langchain.com/oss/python/langchain/human-in-the-loop): "The graph state is saved using LangGraph's persistence layer, so execution can pause safely and resume later." It shows the pattern of Core Feature 10, one item paused on a person while the rest goes on; it pauses a tool call, not a planning decision. The sweep found no study of the cost of an Agent that decides for an absent person, and no tool that puts a keep, change or retire question to the person who decides.

## Risks

- The trigger is the author's judgment. An author that wants its own solution can read friction as a block. The cap, the rule that binds until the answer and the Kept Line bound the cost, and Success Metric 6 shows the load.
- The section can fill with `None.` by habit, because nothing checks it. Success Metric 7 is the signal, and a check is later work that needs these records to read.
- A user in the middle of a feature tends to answer change. The question shows the author's reading of why the rule exists, and the change is verified with the Spec.
- A Spec stops on an Agent's judgment when nobody is there to say the judgment is wrong. Only that Spec stops, and the status view shows it.
- Until the Spec that asked is delivered the old text of the rule stays, so another Spec authored in the meantime can be narrowed by the same rule or ask the same question.
- The two clauses add about 3,900 characters to the Spec routing guide and a section of about 3,500 to each of four skills, which every author reads. The open Backlog Entry about the context an Authoring Phase starts with measures that floor.
- Kept Lines accumulate on a rule that is asked about often.
- Dissent kept from the Idea: a skill step is an instruction, and the same research that supports the clause shows Agents notice a conflict and comply anyway, so the first version may change the words an author reads and leave the behavior where it is. The opposite dissent also stands: the question gives every author a legitimate way to reopen settled decisions.

## Idea-stage judgment

- item: [`references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md`](references/2026-11-18-idea-and-prd-question-rules-that-block-the-goal.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.86)
- unverified_assumption: unsure (0.64)
- market_question: no (0.12)
- verdict: indicated
- stage: ran
- shape: changed — the question gained a trigger and a cap, the Spec that asks carries the change, a kept rule records why it stays, a Spec waits when nobody can answer, a rule the Baseline owns becomes a proposal to its owner, and the TechSpec asks before it contradicts an ADR.

The verdict was indicated, and the maintainer chose "Ideia agora (Recommended)" on 2026-11-18. The judged item is the adopted Backlog Entry. The maintainer approved the Idea, with its `shape` line, on 2026-11-18: "Aprovar e seguir ao PRD (Recommended)".

## Rule challenges

None. The author read the rules that touch this goal: the Instruction Hierarchy, the domain clause on conflicts, the question contract, the ownership of Vendored Skills, the upstream flow of durable knowledge, and ADR-0203, ADR-0222, ADR-0257 and ADR-0186 on what a Baseline clause may say. None blocks the goal. The request itself asks for the rule in the Baseline, and the repository's preference for extending an existing clause is extra work, not a block.

## Glossary

- adds: **Rule Challenge**
- adds: **Kept Line**

## Decisions

Every dated line is the maintainer's answer, as recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The text in quotes is the option he chose, verbatim.

- 2026-11-18: When is the Spec authored, and under which authorization? — "Spec após a release"; "Mesmas do BE (Recommended)": push, Pull Request, merge and release after the Spec, up to five corrections, and the Makefile, the CI workflows, `go.mod` and the lint configuration only for a correction that requires one.
- 2026-11-18: Does the feature start at the Idea stage, and in which order is it delivered? — "Ideia agora (Recommended)"; "Corretiva → regras → Indexus (Recommended)", with the fix of the queue right after the corrective Spec.
- 2026-11-18: How large is the first version? — "Cláusula, passo e registro (Recommended)": a clause in the Baseline, a step in the skills that asks the question, a record of the question and the answer in a fixed form in the artifact, and no mechanical check.
- 2026-11-18: Does the TechSpec ask too? — "Entra, só ao contradizer ADR (Recommended)": the Idea, the PRD and brainstorming question the rule that blocks the goal, and the TechSpec asks when the design needs to contradict an ADR.
- 2026-11-18: Who applies the answer? — "A própria Spec (Recommended)": the answer is recorded in the Idea or the PRD, and a Task of the same Spec changes the rule together with what depended on it.
- 2026-11-18: What happens with nobody present? — "A Spec espera (Recommended)": only that Spec stops, as a Backlog Entry whose `waiting_on` names the question, and the queue goes on.
- 2026-11-18: Can a rule the Baseline owns be questioned? — "Sim, vira proposta ao Roundfix (Recommended)": everything recorded can be questioned, and a Baseline rule becomes a proposal to Roundfix through ai-memory and holds until the Baseline changes.
- 2026-11-18: Where does a keep answer live? — "Na própria regra (Recommended)": the kept rule or ADR gains a dated line with what was asked and why it stayed, and it is asked about again only with a new fact.
- 2026-11-18: Is the Idea approved, with its defaults? — "Aprovar e seguir ao PRD (Recommended)", with the defaults: up to three questions per phase; the keep line of a Baseline rule in an adopter goes in `specific-repository.md`; a Spec that has a folder and nobody present waits through a Backlog Entry that names the Spec and carries the question.
- 2026-11-18: Three answers or four? — "Três: manter, mudar, retirar (Recommended)".
- 2026-11-18: What does a Spec do while its proposal to the Baseline's owner has no answer? — "Segue dentro da regra (Recommended)": the adopter's Spec obeys the rule and records what it left out.

Autonomous route: no person present; decisions taken from the approved Idea and its defaults, the maintainer's answers above, the operator's brief of 2026-11-19 and the prototype recorded in `## Validation`.

- The four open questions the approval left to the Idea's defaults are settled as those defaults say. The minimal PRD of a fix or a refactor carries the section only when a question was asked. A rule keeps every Kept Line, the newest last. Acceptance rests on the words each Task adds, on the published evidence above and on the counts of Success Metrics 5 to 7. The glossary takes two names, Rule Challenge and Kept Line; the things that can be asked about are listed in the definition of the first and get no term of their own.
- The line a kept rule gains reads `Kept (YYYY-MM-DD):`, what was asked and why the rule stays. The Kept Line of a rule inside a Managed Region goes in the repository's own rules guide in every repository, the one that owns the Baseline included, because a Baseline guide holds only what every repository that reads it can use. See ADR-0437.
- When no Spec follows the question, the session that asked adds the Kept Line. The Idea did not say who does; without this a keep answer given in brainstorming would be lost.
- A question the TechSpec author asks is recorded in the PRD, so a Spec has one place for its records and the TechSpec template is unchanged.
- The Baseline text is two clauses in a new rule of the Spec workflow module: the first asks, with the force stop-and-ask, and the second says what follows the answer. The Spec routing guide renders them. The domain clause on conflicts is left as it is, because it still governs the stages that build. The prototype measured that a new clause in the existing route rule needs a hand-written Source Baseline row; the new rule needs none. See ADR-0436.
- The user guide gains a short description, because the maintainer and the adopters' people are the ones who answer the question.
- No accepted ADR is amended. The Instruction Hierarchy keeps its sentence: a question weakens nothing, because the user who confirmed a decision is the one who keeps or changes it.
- The Release Ceiling is `minor`: every adopter of the Spec workflow receives two new clauses, and no command, flag, configuration key or exit code changes.
- ADR-0436 and ADR-0437 are written with this PRD. The TechSpec adds no ADR unless its design needs one.

## Open Questions

None. The measurements of Success Metrics 5 to 7 are the operator's, after the release.
