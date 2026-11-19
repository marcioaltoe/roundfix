---
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
prd: _prd.md
created: 2026-11-19
---

# The Idea and the PRD question the rules that block the goal — Technical Spec

## Executive Summary

The design is text in three places and no product code. The Baseline's `spec-workflow` module gains one rule, `rule.spec.rule-challenge`, with two Normative Clauses that the Spec routing guide renders: the first tells an author to ask, the second says what follows the answer. Four owned skills gain one section with the same core text and a paragraph for their own stage, the Idea and PRD templates gain the `## Rule challenges` section, and the Roundfix skill's `modes` reference says that its modes ask the question and never answer it. The glossary, the user guide, ADR-0436 and ADR-0437 describe it. No command, flag, configuration key, exit code or check changes, and the only Go files that change are tests.

The primary trade-off is the maintainer's: nothing mechanical holds the step. An author follows the clause and the skill or does not, and `roundfix spec check` cannot tell which. In exchange the first version is text that every adopter receives with one Baseline update, it can refuse no plan, and it leaves the records a later check would read. The design answers the risk it can answer: contract tests pin every guarantee of the clauses and of the skill section by its own phrase, so a later rewording cannot drop one unnoticed.

A second trade-off decides where the clauses live. The route rule of the Spec workflow is their natural neighbour, and the Standard TypeScript profile requires that rule: the catalog refuses a new clause in it until a Source Baseline row is written by hand, which no command creates (`references/2026-11-19-prototype-and-governed-paths.md`, section 2). A rule that no profile requires needs no row. The clauses therefore live in a rule of their own, which the Spec routing guide lists right after the route rule.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity; the rule id `rule.spec.rule-challenge` and the clause ids `clause.spec.rule-challenge-01-ask` and `clause.spec.rule-challenge-02-answer` follow the form the module already uses. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added, and no test or Verification asks a person, a provider or ai-memory anything. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands, with ADR-0436 and ADR-0437 as the decisions this design implements, ADR-0320, ADR-0403, ADR-0353, ADR-0186, ADR-0222, ADR-0257, ADR-0411, ADR-0256 and ADR-0349 as the ones it obeys, and ADR-0203, ADR-0278, ADR-0258, ADR-0338, ADR-0342, ADR-0383, ADR-0390, ADR-0404, ADR-0410, ADR-0412 and ADR-0432 listed to record why they do not apply. This stage adds no ADR. The clause texts name no ADR and no Spec of this repository (ADR-0186: "Baseline guidance states what the product does, in adopter-neutral words"), and the new clauses take new identities in a new rule, so no clause is reworded and no Source Baseline row is needed (ADR-0257: "a new identity needs a Source Baseline row that the regenerator"). Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18 for this Spec, quoted in the authorization record, which keeps the grant of every owned skill and of the Governed Paths the Spec declares, Baseline modules included. The design changes no Makefile, workflow, Go module, lint or formatter file. It changes one Baseline module, the Standard TypeScript profile's pin and formatter golden of the Spec routing guide, this repository's Spec routing guide and Setup Manifest, and no other file under `docs/agents/`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0409-the-idea-and-prd-question-rules-that-block-the-goal/_authorization.md`; bounded files: `.agents/skills/brainstorming/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-idea/references/idea-template.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/brainstorming/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`.

## System Architecture

Every component extends something that exists. Nothing is added to `internal/` but two test files and two lines of a test's closed list.

### Rule Challenge rule

`internal/baseline/assets/modules/spec-workflow.json` gains the rule `rule.spec.rule-challenge`, version 1, with the module's coverage id and the two clauses of API Contracts 1 and 2. The guide `guide.spec-routing` lists it between `rule.spec.routing` and `rule.spec.project-constraints`, and its guide version rises by one. The guide renders a rule's clauses in the order of the guide's list, so both clauses appear in `docs/agents/spec-routing.md` after the route clauses and before the Project Constraints clauses. No profile lists the rule under `requiredRules`; `rule.spec.project-constraints`, which the same guide renders, is in the same position. The module version rises through its record command, never by hand.

### Derived Baseline files and Managed Refresh

The module edit moves the files the prototype measured: the module version record, the Standard TypeScript profile's formatter pin and its golden Spec routing guide, three catalog files and four plan goldens. `make baseline-digests` writes them. This repository's own `docs/agents/spec-routing.md` and Setup Manifest follow from the Managed Refresh. `internal/baseline/clause_characterization_test.go` holds the closed list of every clause's force and gains the two ids.

### Skill section and stage paragraphs

`brainstorming`, `write-idea`, `write-prd` and `write-techspec` each gain a `## Rule Challenge` section: the core text of API Contract 3, the same bytes in the four, followed by the paragraph of API Contract 4 for that skill. Three existing sentences gain a pointer to it (API Contract 5). The canonical copies under `.agents/skills/` are edited, `make skills-sync` writes the mirrors, and the version record command raises each skill's version.

### Record section

The Idea template and the PRD template gain the `## Rule challenges` section of API Contract 6. The TechSpec template is unchanged: a question the TechSpec author asks is recorded in the PRD.

### Modes reference

The Roundfix skill's `references/modes.md` gains one paragraph under `--unattended` and one sentence at the end of step 5 of `--follow-up` (API Contract 7). The step count stays eight.

### Glossary, user guide and decision records

`CONTEXT.md` gains the two entries of API Contract 8 after the entry of Idea Stage Judgment. `docs/user-guide/context-driven-development.md` gains the subsection of API Contract 9 before `### Spec Event Announcements`, and the Behavior Surface Record is recorded again after that edit. ADR-0436 and ADR-0437 arrive with the planning change; no Task writes them.

### Contract tests

Two new test files hold the words: `internal/baseline/rule_challenge_clauses_test.go` for the clauses and their rendering, and `skills/rule_challenge_authoring_test.go` for the skill section, the stage paragraphs, the record section and the `modes` paragraphs.

## Implementation Design

### Interfaces

No Go interface changes. The module gains this object in its `rules` array, before `rule.spec.artifacts`, and one line in the guide's list:

```json
{
  "id": "rule.spec.rule-challenge",
  "version": 1,
  "coverage": ["coverage.spec-workflow"],
  "clauses": [
    {"id": "clause.spec.rule-challenge-01-ask", "enforcement": "stop-and-ask", "guidance": "<API Contract 1>"},
    {"id": "clause.spec.rule-challenge-02-answer", "enforcement": "mandatory", "guidance": "<API Contract 2>"}
  ]
}
```

```json
"rules": ["rule.spec.routing", "rule.spec.rule-challenge", "rule.spec.project-constraints"]
```

Invariants:

```text
1. rule.spec.rule-challenge holds exactly the two clauses above, with those ids and forces, and no clause declares replaces.
2. No profile lists the rule under requiredRules, and no file under internal/baseline/assets/source-baselines or internal/baseline/assets/setups changes.
3. A plan for any profile that includes the spec-workflow module renders each clause exactly once in docs/agents/spec-routing.md, after the route clauses and before the Project Constraints clauses.
4. The clause texts name no ADR number and no Spec number.
5. Every clause that exists at the delivery base keeps its id, its force and its text, and the sentence of the Instruction Hierarchy is unchanged.
6. The core text of the skill section is the same bytes in the four canonical skills, and every mirror under skills/ equals its canonical file.
7. The record section is the same in both templates except the phase word of its example: idea in the Idea template, prd in the PRD template.
8. No text is added inside a "### QA settlement" section, the --follow-up section keeps eight numbered steps, and no sentence an existing test pins is changed.
9. No version number, digest, golden, mirror or record is written by hand, and no test reads a file under the Spec root.
10. The regeneration runs in this order: module record, make baseline-digests, Managed Refresh until a run reports no file change, make skills-sync, skill version record, and the Behavior Surface Record after the last guide edit.
```

### Data Models

**The record.** One record per Rule Challenge under `## Rule challenges`, each a list of fielded lines, with a blank line between two records. The section's single entry is `None.` when no rule blocked the goal; text may follow it.

| field | holds |
| --- | --- |
| `rule` | the ADR, clause, definition or document, and its path |
| `imposes` | what the rule requires |
| `prevents` | up to three things, for this goal, separated by semicolons |
| `asked` | the date as `YYYY-MM-DD`, a comma, and the phase: `brainstorming`, `idea`, `prd` or `techspec` |
| `answer` | `keep`, `change`, `retire` or `waiting`; for the first three, a dash, the user's words in quotes and the date of the answer |
| `effect` | what the Spec does with the answer; for a rule the Baseline owns, the message sent and the scope left out |
| `new fact` | only when the rule already carried a Kept Line: what this question knows that the earlier one did not |

No code parses the record. `roundfix spec check` holds no closed list of an artifact's sections, and its coverage units are the User Stories, the Core Features and the Success Metrics, so the section moves no finding.

**The Kept Line.** One line after the rule's text: `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.` The date is the date of the user's answer. An accepted ADR takes it after its Decision text and its `updated_at` rises; a Glossary Definition takes it before its `_Avoid_` line; a rule inside a Managed Region takes it in `docs/agents/specific-repository.md`, under a heading of its own, with the rule named in the line.

**The waiting entry.** An `open` Backlog Entry in the form the docs-layout guide gives, with `waiting_on` naming the Rule Challenge and the rule, and `since`. Its body names the Spec folder when one exists and holds the record with `answer: waiting` and the three options. The Follow-up Command already reads it as a Waiting Item and the status view already lists it; neither changes.

### API Contracts

Each block below is the exact text a Task writes. A block is a fence in this file only; the fence is not part of the text.

1. API Contract: the `guidance` of `clause.spec.rule-challenge-01-ask`, force `stop-and-ask`, one JSON string with no line break.

```text
In brainstorming, the Idea and the PRD, never drop or narrow a solution in silence because a rule forbids it, and never contradict a rule in silence: put a Rule Challenge to the user. A rule is a decision the repository recorded: an accepted ADR, a Normative Clause of a guide under `docs/agents/`, a definition in the domain context, or a document under `docs/references/` that the repository treats as binding. A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem; extra work, taste and friction are not a block. A Rule Challenge is one question about one rule, asked as the stop-and-ask clause of `docs/agents/agent-instructions.md` says: name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why it exists and whether it looks out of date or badly written, then offer keep, change and retire, each with its consequence and the recommended one first. The rule binds until the user answers: asking changes nothing, and a recommendation, silence, a dismissal or a timeout is never an answer. Ask no Rule Challenge about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records; name that request or answer where the artifact records the change. Ask none about a rule that carries a Kept Line unless the question names a fact the earlier one did not have. Ask at most three in one Authoring Phase or brainstorming session: when more rules block, ask about the three that prevent the most and list the others, which keep binding, in the artifact's open questions. The TechSpec author asks the same question before a design takes what an accepted ADR rejected, or otherwise contradicts its decision, unless the user already decided it. In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge.
```

2. API Contract: the `guidance` of `clause.spec.rule-challenge-02-answer`, force `mandatory`, one JSON string with no line break.

```text
Record every Rule Challenge in the `## Rule challenges` section of the Idea or the PRD, one fielded record per question: `rule`, `imposes`, `prevents`, `asked` with the date and the phase, `answer` with keep, change, retire or waiting and the user's words, and `effect`, plus `new fact` when the rule already carried a Kept Line. An Idea and a feature PRD always carry the section, which reads `None.` when no rule blocked the goal; the PRD carries the Idea's records, and a question the TechSpec author asks is recorded in the PRD. The Spec that asked carries the answer. On change or retire, the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it: an ADR through the lifecycle of `docs/agents/docs-layout.md`, any other rule in its own file. On keep, a Task of the Spec adds a Kept Line after the rule's text, `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.`, which names no Spec; a rule keeps every Kept Line, the newest last, and when no Spec follows, the session that asked adds the line. A rule inside `setup-context-driven` markers changes only in the Baseline, and its Kept Line goes in `docs/agents/specific-repository.md`, under a heading of its own, naming the rule; on change or retire, a repository that does not own the Baseline sends the user's decision as one message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out. When no person is present to answer, write no artifact for that phase and never answer for the user: only that Spec waits, as a Backlog Entry that holds the whole question, names the Spec folder when one exists and carries `waiting_on` and `since` in its frontmatter, while other work goes on. An Author Subagent returns the question to the session that started it.
```

3. API Contract: the core text of the `## Rule Challenge` section, the same in `brainstorming`, `write-idea`, `write-prd` and `write-techspec`. In `write-idea`, `write-prd` and `write-techspec` the section goes immediately before `## Declared model`; in `brainstorming`, immediately before `## Idea-stage judgment`.

```markdown
## Rule Challenge

A rule guides the solution and never bounds the goal in silence. A rule is a decision the repository recorded: an accepted ADR, a Normative Clause of a guide under `docs/agents/`, a definition in `CONTEXT.md`, or a document under `docs/references/` that the repository treats as binding. A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem. Extra work, taste and friction are not a block.

When a rule blocks the goal, neither drop or narrow the solution nor contradict the rule: put a Rule Challenge to the user, one question per rule, through the structured question tool the session exposes. Name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why the rule exists and whether it looks out of date or badly written. Offer keep, change and retire, each with its consequence, the recommended one first and labeled `(Recommended)`. The rule binds until the user answers: asking changes nothing, and a recommendation, silence, a dismissal or a timeout is never an answer.

- Ask no Rule Challenge about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records. Name that request or answer where the artifact records the change.
- Ask none about a rule that carries a Kept Line, a line `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.` after its text, unless the question names a fact the earlier one did not have. Record that fact as `new fact`.
- Ask at most three in one Authoring Phase or brainstorming session. When more rules block, ask about the three that prevent the most and list the others, which keep binding, in the artifact's open questions.
- Never answer for the user. An Author Subagent returns the question, with its options and the one it recommends, to the session that started it. When no person is present, write no artifact for this phase: the Spec waits as a Backlog Entry that holds the whole question in the record form with `answer: waiting`, names the Spec folder when one exists, and carries `waiting_on`, naming the Rule Challenge, and `since` in its frontmatter.

What follows the answer:

- Change or retire: the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it. An ADR changes through the lifecycle of `docs/agents/docs-layout.md`.
- Keep: the solution stays inside the rule, and the PRD requires a Task to add the Kept Line after the rule's text. The line names no Spec.
- A rule inside `setup-context-driven` markers changes only in the Baseline, and its Kept Line goes in `docs/agents/specific-repository.md`, under a heading of its own, naming the rule. On change or retire, a repository that does not own the Baseline sends the user's decision as one ai-memory message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out.
```

4. API Contract: the stage paragraph that ends the section, one per skill, after a blank line.

`write-idea`:

```markdown
In this stage, list the rules that touch the idea in the codebase sweep of step 4, and never drop an alternative of the Opportunity Scan only because a rule forbids it: keep it, name the rule beside it, and ask the Rule Challenge when it is the direction you would recommend. Record every Rule Challenge in the `## Rule challenges` section of `_idea.md`, in the fielded form the template shows; the section reads `None.` when no rule blocked the goal.
```

`write-prd`:

```markdown
In this stage, read every ADR and rule the Project Constraints snapshot will list as an obligation and as a candidate: one that blocks a goal becomes a Rule Challenge before the snapshot is written. Copy the `## Rule challenges` records of `_idea.md` into the PRD unchanged and add the ones this stage asks, in the fielded form the template shows; the section reads `None.` when no rule blocked the goal. A change or retire answer becomes a Core Feature that states the change of the rule, and a keep answer a Core Feature that requires the Kept Line.
```

`brainstorming`:

```markdown
In this conversation, when the direction the user wants meets a rule that blocks it, say so and ask the Rule Challenge instead of steering the user away from that direction. Hand the question and the user's answer to the skill the outcome routes to, which records them in the `## Rule challenges` section of its artifact; when no Spec follows a keep answer, add the Kept Line yourself.
```

`write-techspec`:

```markdown
In this stage the question has a narrower trigger: ask it before the design takes what an accepted ADR rejected, or otherwise contradicts its decision, when neither the request, the PRD's `## Rule challenges` section nor another recorded answer of the user decides it. Record it in the PRD's `## Rule challenges` section with the phase `techspec`, adding the section to a minimal PRD, which otherwise omits it, and end that ADR's line under `## ADR Alternatives` with `Rule Challenge of YYYY-MM-DD: <answer>`. An amendment or a successor the user already decided needs no question; the line names where the user decided.
```

5. API Contract: four existing sentences gain a pointer. Each pair is the whole change to its sentence.

`write-prd`, `## Inputs`:

```text
before: Copy its `## Idea-stage judgment` section into the PRD unchanged.
after:  Copy its `## Idea-stage judgment` section into the PRD unchanged, and its `## Rule challenges` records as `## Rule Challenge` below says.
```

`brainstorming`, the bullet "Explore before asking":

```text
before: Check project state first — files, docs, `CONTEXT.md`, active Specs
after:  Check project state first — files, docs, `CONTEXT.md`, the accepted ADRs under `docs/adr/`, the guides under `docs/agents/`, active Specs
```

`write-techspec`, the ground rule "Read first":

```text
before: (respect prior decisions or challenge them explicitly — never contradict one silently)
after:  (respect prior decisions or challenge them explicitly — never contradict one silently, and ask the user first as `## Rule Challenge` below says)
```

`write-techspec`, the bullet "Rejected alternatives" under `## Proof of what the Spec claims`, which is wrapped over several lines and keeps its wrapping:

```text
before: say so and supersede or amend the ADR.
after:  say so and supersede or amend the ADR, after the Rule Challenge when the user has not decided it.
```

6. API Contract: the record section of the templates. In `idea-template.md` it goes immediately before `## Decisions` and its example reads `asked: YYYY-MM-DD, idea`; in `prd-template.md` it goes immediately before `## Glossary` and reads as below.

```markdown
## Rule challenges

<!-- One record per Rule Challenge, in this fielded form, or the single entry `None.` when no rule blocked the goal.
     rule: the ADR, clause, definition or document, and its path.
     prevents: up to three things, for this goal.
     asked: the date and the phase: brainstorming, idea, prd or techspec.
     answer: keep | change | retire | waiting, then a dash, the user's words in quotes and the date of the answer.
     effect: what the Spec does with the answer; for a rule the Baseline owns, the message sent and the scope left out.
     new fact: only when the rule already carried a Kept Line; what this question knows that the earlier one did not.
     Leave a blank line between two records. -->

- rule: ADR-NNNN, `docs/adr/NNNN-<slug>.md`
- imposes: <what the rule requires>
- prevents: <what it prevents for this goal>
- asked: YYYY-MM-DD, prd
- answer: change — "<the user's words>", YYYY-MM-DD
- effect: <the requirement that states the change, and the Task that makes it>
```

7. API Contract: the Roundfix skill's `references/modes.md`. Under `## --unattended`, a new paragraph after the paragraph that begins "The Idea stage is the user's":

```markdown
A Rule Challenge is the user's too: this mode asks it and never answers it. When an Author Subagent returns one, ask it in the next question call, with keep, change and retire as its options, and continue the subagent with the answer. A Spec whose Rule Challenge gets no answer writes no artifact for that phase and waits as a Backlog Entry that holds the whole question and names the Spec folder when one exists, with a Waiting Declaration whose `waiting_on` names the Rule Challenge, and the report names it.
```

Under `## --follow-up`, step 5 gains this sentence at its end, after "its Waiting Declaration naming the Idea stage.":

```markdown
A Spec whose author returns a Rule Challenge goes on only after the user answered it; without an answer it waits as `--unattended` says.
```

8. API Contract: the two glossary entries, after the entry of Idea Stage Judgment.

```markdown
**Rule Challenge**:
The question an author puts to the user when a rule blocks the goal that a brainstorming session, an Idea or a PRD is shaping, and that a TechSpec author puts before a design takes what an accepted ADR rejected. A rule is an accepted ADR, a Normative Clause, a Glossary Definition or a document under `docs/references/` that the repository treats as binding. The question names the rule, what it imposes and up to three things it prevents, and offers keep, change and retire; the rule binds until the user answers, the Idea or the PRD records the question and the answer under `## Rule challenges`, and the Spec that asked carries the change ([ADR-0436](docs/adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md)).
_Avoid_: Rule waiver, exception request, rule review

**Kept Line**:
The dated line a rule gains when the user answers a Rule Challenge with keep: `Kept (YYYY-MM-DD):`, what was asked and why the rule stays, written after the rule's text and, in a Glossary Definition, before its `_Avoid_` line. It names no Spec, a rule keeps every one with the newest last, the line of a rule inside a Managed Region goes in the repository's own rules guide, and an author asks about that rule again only with a fact the line does not hold ([ADR-0437](docs/adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md)).
_Avoid_: Kept note, reaffirmation, review stamp
```

9. API Contract: the user guide subsection, in `docs/user-guide/context-driven-development.md` immediately before `### Spec Event Announcements`.

```markdown
### Rule Challenges

In brainstorming, the Idea and the PRD, a recorded decision guides the solution and never bounds the goal in silence. When an accepted ADR, a Normative Clause of a guide, a glossary definition or a binding document under `docs/references/` blocks the goal, the author puts a Rule Challenge to the user: it names the rule, what the rule imposes and up to three things it prevents for this goal, and offers keep, change and retire. The rule binds until the user answers, and a recommendation, silence, a dismissal or a timeout is never an answer.

| Answer | What follows |
| --- | --- |
| keep | The solution stays inside the rule, and a Task of the Spec adds a Kept Line after the rule's text: `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.` The rule is asked about again only with a new fact. |
| change | The PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it. |
| retire | The PRD states the retirement as a requirement, and a Task of the same Spec removes the rule as the retirement rule below says. |

The Idea and the PRD record each question under `## Rule challenges`, and write `None.` there when no rule blocked the goal. An author asks at most three Rule Challenges in one Authoring Phase or brainstorming session, and the TechSpec author asks one before a design takes what an accepted ADR rejected. With nobody present to answer, only that Spec waits, as a Backlog Entry with a Waiting Declaration, and other work goes on. A rule inside a Managed Region changes only in the Baseline: the user's decision goes to the Baseline's owner as one ai-memory message, and the rule stays in force until the Baseline changes. In a Run, in a Task and in the QA gate a rule is a guard-rail, and no Rule Challenge is asked there ([ADR-0436](../adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md), [ADR-0437](../adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md)).
```

### Surface Transcripts

None. No command surface changes: no command, flag, help text, output line or exit code is added or changed, and `roundfix baseline update` prints its existing lines for a plan that carries two more clauses.

## Coverage Map

- Goal 1 → Rule Challenge rule: its first clause forbids both silences, with the force stop-and-ask, in every repository that adopts the Spec workflow (API Contract 1); Skill section and stage paragraphs: each skill tells its author when to ask, how, and that the rule binds until the answer (API Contracts 3 and 4); Record section: the artifact says whether a rule blocked the goal (API Contract 6). Nothing enforces it mechanically, by the decision the PRD records, and the PRD declares the live behavior under Unreachable Acceptance.
- Goal 2 → Record section (API Contract 6); Rule Challenge rule (API Contract 2).
- Goal 3 → Rule Challenge rule (API Contract 2), Skill section: "What follows the answer".
- Goal 4 → Rule Challenge rule (API Contracts 1 and 2), the Kept Line in Data Models, Glossary (API Contract 8).
- Goal 5 → Rule Challenge rule with Derived Baseline files and Managed Refresh; API Contract 1, its last sentence, and invariant 5.
- Story 1 → Rule Challenge rule (API Contract 1), Skill section (API Contract 3).
- Story 2 → the question's content in API Contracts 1 and 3.
- Story 3 → API Contracts 1 and 3: "The rule binds until the user answers".
- Story 4 → Record section (API Contract 6), the record in Data Models.
- Story 5 → API Contract 2 and the stage paragraph of `write-prd` (API Contract 4).
- Story 6 → the Kept Line in Data Models, API Contracts 2, 3 and 8.
- Story 7 → API Contract 2, its last two sentences, Skill section, Modes reference (API Contract 7), the waiting entry in Data Models.
- Story 8 → API Contract 2 and the last bullet of API Contract 3.
- Story 9 → API Contract 1, the stage paragraph of `write-techspec` (API Contract 4), API Contract 5.
- Story 10 → Rule Challenge rule, Derived Baseline files and Managed Refresh; the trigger in API Contract 1.
- Core Feature 1 → Rule Challenge rule (API Contract 1), Derived Baseline files and Managed Refresh.
- Core Feature 2 → API Contracts 1 and 3, the sentence that says what a rule is.
- Core Feature 3 → API Contracts 1 and 3, the trigger and the change the user already decided.
- Core Feature 4 → API Contracts 1 and 3, the question.
- Core Feature 5 → API Contracts 1 and 3, the cap.
- Core Feature 6 → API Contracts 1 and 3, the binding sentence.
- Core Feature 7 → Record section (API Contract 6), API Contract 2, the stage paragraphs of `write-idea`, `write-prd` and `write-techspec`.
- Core Feature 8 → API Contract 2, Skill section: "Change or retire".
- Core Feature 9 → API Contract 2, Skill section: "Keep", the stage paragraph of `brainstorming`, the Kept Line in Data Models.
- Core Feature 10 → API Contract 2, Skill section, Modes reference (API Contract 7), the waiting entry in Data Models.
- Core Feature 11 → API Contract 2, the last bullet of API Contract 3.
- Core Feature 12 → API Contract 1, the stage paragraph of `write-techspec`, API Contract 5.
- Core Feature 13 → API Contract 1, its last sentence; invariant 5; Contract tests.
- Core Feature 14 → Skill section and stage paragraphs (API Contracts 3, 4 and 5), Record section (API Contract 6).
- Core Feature 15 → Glossary, user guide and decision records (API Contracts 8 and 9).
- Success Metric 1 → Rule Challenge rule, Derived Baseline files and Managed Refresh; Contract tests, the rendering test.
- Success Metric 2 → Skill section and stage paragraphs, Record section, Modes reference; Contract tests, the skills tests.
- Success Metric 3 → Glossary, user guide and decision records; the existing glossary and documentation contracts.
- Success Metric 4 → invariant 5; the closed list of clause forces; the existing tests of `./internal/baseline` and `./skills`.
- Success Metric 5 → Record section and the stage paragraphs, which make every Idea and feature PRD carry the section; declared unreachable in the PRD.
- Success Metric 6 → the cap and the Kept Line in API Contracts 1 and 3; declared unreachable in the PRD.
- Success Metric 7 → Rule Challenge rule and Skill section as a whole; declared unreachable in the PRD.

Deviations from the PRD: None.

## Absence Evidence

Each search ran on the planning clone at `348962d2` and excludes this Spec's own folder and its two ADRs.

1. The Rule Challenge rule and its clauses — command: `grep -rn "rule-challenge\|Rule Challenge" internal/baseline docs/agents`; result: no match.
2. The skill section, the stage paragraphs and the record section — command: `grep -rln "Rule Challenge\|Rule challenges" .agents/skills skills`; result: no file.
3. An existing question with these three answers — command: `grep -rn "keep, change" .agents/skills internal/baseline/assets docs/agents docs/user-guide`; result: no match.
4. The glossary terms — command: `grep -c "Rule Challenge\|Kept Line" CONTEXT.md`; result: 0.
5. The Kept Line form — command: `grep -rn "^Kept (" docs/adr docs/agents docs/references CONTEXT.md`; result: no match.
6. The user guide subsection — command: `grep -rn "Rule Challenge" docs/user-guide`; result: no match.
7. The two test files — command: `ls internal/baseline/rule_challenge_clauses_test.go skills/rule_challenge_authoring_test.go`; result: no such file, both.

## ADR Alternatives

- ADR-0436: the design is its decision; listing a rule as an obligation, an override by the author, a clause without a record, a mechanical check, a change outside the Spec, going on inside the rule with nobody present, four answers and questions during implementation are the options it rejected.
- ADR-0437: the design is its decision; an answer kept only in the Spec, a new ADR for each keep answer, a review date in the frontmatter, strict immutability and one place for every line are the options it rejected.
- ADR-0320: the question uses the clause as it is; three Rule Challenges about three rules are independent and fit one call, and no option it rejected is taken.
- ADR-0403: the wait with nobody present repeats its rule for the Idea stage, a Backlog Entry that waits; the Idea-stage judgment and its texts are unchanged, and no option it rejected is taken.
- ADR-0353: the Waiting Declaration stays on a Backlog Entry or a Finding; a Spec folder carries none, which is why the waiting entry names the Spec. No option it rejected is taken.
- ADR-0186: the clause texts name no ADR and no Spec of this repository; no option it rejected is taken.
- ADR-0222: the clauses hold for every repository that reads them: they name the repository that owns the Baseline and the one that does not, and the Kept Line of a rule the Baseline owns stays out of the Baseline's own text.
- ADR-0257: each rule is stated once, in the new clauses, and the domain clause on conflicts keeps its own subject; the clauses take new identities in a rule no profile requires, so no Source Baseline row is hand-written. Extending an existing clause, which that ADR prefers for a reworded rule, is not taken because nothing is reworded.
- ADR-0411: the Kept Line of a rule the Baseline owns goes in the repository's own rules guide under a heading of its own, where that ADR places what a repository writes; nothing is moved or detected.
- ADR-0256: the skills and the user guide describe the behavior in the same delivery; the covering skill of the changed guide page, `write-prd`, changes with it.
- ADR-0349: the authorization record carries the Release Ceiling `minor`.
- ADR-0203: not applicable — no hand-written rule of a repository is promoted; the maintainer asked for a new rule in the Baseline, and its evidence from outside Roundfix is recorded in the PRD.
- ADR-0278: not applicable — no Task of this Spec retires a file; a retire answer given in a later Spec follows the retirement clause as the docs-layout guide then states it, and the design adds no status and no history file.
- ADR-0258: not applicable — the QA rerun classes and the retirement of history are untouched.
- ADR-0338: not applicable — the follow-up mode still claims the handoff and leaves no handed-in message pending; its step 5 gains one sentence.
- ADR-0342: not applicable — every Task still settles on the repository Verification.
- ADR-0383: not applicable — the rule on replies is untouched; the proposal about a rule the Baseline owns is a message that declares no blocked work, so the Spec records no wait for it.
- ADR-0390: not applicable — the precedence of a blocking message and the answer to a direct question are untouched.
- ADR-0404: not applicable — the Idea-stage judge and the item it reads are untouched.
- ADR-0410: not applicable — the status view is untouched; it lists the waiting entry as it lists every open Backlog Entry.
- ADR-0412: not applicable — planning references and the release preparation keep their release impact.
- ADR-0432: not applicable — `CLAUDE.md` stays a link to `AGENTS.md`, and no root instruction file changes.

## Integration Points

- The structured question tool of each harness, through the existing question contract only. Roundfix calls no tool; the clause and the skills tell the author to use the one its session exposes.
- ai-memory, through the existing cross-project rule only: an Agent session sends the proposal about a rule the Baseline owns. No Task and no test calls it.
- An adopting repository, through `roundfix baseline update`: its next update writes the two clauses into its Spec routing guide, and `roundfix skills install` brings the four skills.
- The Follow-up Command and the status view read the waiting entry as they read every Backlog Entry with a Waiting Declaration.

## Testing Approach

- **Existing seams, no new one.** The clause tests read the module file as the route clause tests do and build a plan on a temporary repository with `newProjectDecisionPlanRepository`, `buildProjectDecisionPlan`, `standardTypeScriptDecisions` and `planPostimage`. The skill tests read the embedded skills with `assertIdeaStageSkillPhrases`. Every new test calls `t.Parallel()` first.
- **Clause tests**, in `internal/baseline/rule_challenge_clauses_test.go`:
  - `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses` reads the module and asserts invariant 1, and one phrase for each guarantee of API Contracts 1 and 2, on whitespace-normalized text. For the first clause: both silences, what a rule is, the trigger, what is not a block, the question's content, the three answers, the binding sentence, what is never an answer, the change the user already decided, the Kept Line and its new fact, the cap and the rules that keep binding, the TechSpec, and the guard-rail of implementation. For the second: the record's fields, `None.`, the Idea's records, the TechSpec's record, the requirement and the Task, the ADR lifecycle, the Kept Line's form, no Spec named, the newest last, the session that adds it, the Baseline's rule, its Kept Line, the proposal, the wait, never answering for the user, and the Author Subagent.
  - `TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce` builds a plan and asserts invariant 3 on the postimage of `docs/agents/spec-routing.md`, and the same on the Standard TypeScript formatter golden.
- **Skill tests**, in `skills/rule_challenge_authoring_test.go`:
  - `TestTheAuthoringSkillsPutARuleChallengeToTheUser` asserts, in the four skills, the heading and one phrase for each guarantee of API Contract 3, that the core text is equal in the four (invariant 6), each stage paragraph of API Contract 4 and each pointer of API Contract 5.
  - `TestTheTemplatesCarryTheRuleChallengeRecord` asserts the heading, the seven field names and `None.` in both templates, and the phrase `asked: YYYY-MM-DD, idea` or `asked: YYYY-MM-DD, prd` (invariant 7).
  - `TestTheSkillModesNeverAnswerARuleChallenge` asserts the paragraph and the sentence of API Contract 7.
- **Tests this change moves.** One: `TestAnUnlistedBaselineClauseIsReported` reads the closed list `characterizedBaselineForce`, which gains the two ids with their forces. `make baseline-digests` rewrites the goldens that the catalog, formatter and plan characterization tests compare.
- **Tests that hold what must not change.** They pass unchanged: `TestInstructionHierarchyRendersActivePointersOnce` for the hierarchy sentence, `TestTheGlossaryClausesCarryTheirForceAndText` for the domain clause on conflicts, `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall` for the question contract, the route clause tests of `internal/baseline/idea_stage_route_clauses_test.go`, `TestShippedGuidanceCitesNoRepositoryRecord`, `TestNoTwoBaselineClausesShareText`, `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage`, `TestTheFollowUpModeRunsItsStepsInOrder`, `TestTheUnattendedModeAsksWithinTheBaselineLimit`, `TestTheFollowUpModeRunsTheIdeaStageJudgment` and the skills contract tests of `skills/baseline_skill_contract_test.go`.
- **Glossary and guide.** The existing glossary contracts hold the two entries: `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryAvoidLinesArePlainLists` and `TestGlossaryLinksNameAcceptedADRs`. `TestTheSkillCoverageMapIsCurrent` holds the Behavior Surface Record after its record command. A Task's Verification names each term and a phrase of the guide subsection.
- **Measured.** The prototype with these exact texts and the six regeneration commands passed `./internal/baseline`, `./skills`, `./internal/docscontract` with its tag, `./internal/speccheck`, `./internal/spec`, the Baseline tests of `./internal/cli` and `make verify-docs` (`references/2026-11-19-prototype-and-governed-paths.md`, section 4).
- **Not tested, by decision.** Whether an author asks. No test starts an Agent session, asks a person, calls ai-memory or reaches the network, and no test depends on a Git tag, on the platform or on a file under the Spec root.

## Build Order

1. Rule Challenge rule: the module edit of API Contracts 1 and 2, the guide's list, the closed list of clause forces, the clause tests, then the module record, `make baseline-digests` and the Managed Refresh.
2. Skill section, stage paragraphs, Record section and Modes reference: the canonical edits of API Contracts 3 to 7, the skill tests, `make skills-sync` and the skill version record (depends on: 1).
3. Glossary and user guide: API Contracts 8 and 9, with the Behavior Surface Record recorded after the guide edit (depends on: 2).
4. The QA gate (depends on: 1, 2, 3).

Steps 1 and 2 share no file; step 2 follows step 1 so that the skills and the clause they repeat arrive in one order and `make baseline-digests`, which reads the skill mirrors, runs on settled skills. Step 3 is small: by the sizing criterion it may join the Task of step 2, which the Tasks stage decides. Spec 0411 is delivered before this Spec and shares the files the PRD's Prerequisites name; no step depends on its work.

## Risks & Considerations

- **A guarantee lost in a later rewording.** The clauses are long. Each guarantee has its own phrase in a contract test, as the Testing Approach lists.
- **The route rule's Source Baseline row.** A later author who moves the clauses into `rule.spec.routing` meets the catalog's refusal `catalog.sourceBaseline.required-clause.missing`. The rule stays its own.
- **Regeneration order.** The module record runs before `make baseline-digests`, the Managed Refresh after it, and the Behavior Surface Record after the last guide edit (invariant 10). A second Managed Refresh must report no file change.
- **Shared files.** Specs 0411 and 0410 change the Roundfix skill's version and `modes` reference, the version records, `CONTEXT.md` and the Behavior Surface Record. Each Task records versions and fingerprints on the tree it runs on and edits only its own sentences.
- **The Setup Manifest.** The Managed Refresh also writes the Roundfix version it runs with into `docs/agents/setup-context.json`; that is the command's own output and is not hand-edited.
- **Reading cost.** The two clauses add about 3,900 characters to the Spec routing guide, and the section about 3,500 to each of four skills. The open Backlog Entry about the context an Authoring Phase starts with measures that floor.
- **Adopter-neutral words.** The clauses name `docs/agents/specific-repository.md` and `docs/agents/docs-layout.md`, which the Baseline writes in every profile that includes their modules, and name no file of this repository alone.

## Glossary

None.

## Decisions

- The clauses live in a new rule of the `spec-workflow` module, not in the route rule and not in the domain rule: the prototype measured that the route rule needs a hand-written Source Baseline row, and the domain rule, which the Standard TypeScript profile also requires, belongs to a module Spec 0410 changes.
- Two clauses, not one: the first asks, with the force stop-and-ask, and the second records and carries the answer, with the force mandatory. One clause would hide the second force.
- The domain clause on conflicts is unchanged. It still governs the stages that build, and the new first clause ends by saying so.
- One core text in four skills, with a paragraph per stage, as the skills already repeat the sizing rule and the Declared model section. A shared reference file would not reach `brainstorming`, which has no references directory.
- The TechSpec template is unchanged, and the minimal PRD of a fix or a refactor gains the section only when a question was asked.
- No code parses the record. The form is fixed so that a later check or count can read it.
- The user guide gains one subsection in the page that describes the pipeline; its covering skill, `write-prd`, changes in the same Spec.
- No new ADR at this stage: ADR-0436 and ADR-0437, written with the PRD, hold the durable decisions, and where the clauses live is cheap to change.
