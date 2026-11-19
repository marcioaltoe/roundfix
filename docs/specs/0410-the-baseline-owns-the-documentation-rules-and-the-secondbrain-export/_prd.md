---
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: active
created: 2026-11-18
surfaces: [cli, docs]
type: feat
---

# The Baseline owns the documentation rules and the Secondbrain export

Nothing in the Baseline says what a repository publishes to the Secondbrain, so each project exports all of `docs/`, and 1,509 of the 3,534 documentation files in the project mirrors are process. What a project learned the hard way is mostly absent, because the Finding that held the cause is deleted when it closes, and an ADR that a later one changed still reads as current. After this Spec the Baseline states the job of each documentation folder and keeps, in each repository that has the Secondbrain turned on, a standard list of what is published: knowledge, never process. An ADR says in its header whether the knowledge base wants it and which ADRs amend it, a Lesson Learned is written only when there is a lesson, `roundfix baseline update` applies what a command can apply, and `roundfix doctor` reports what departs without failing.

## Sources

1. [`references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md`](references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md), the adopted Backlog Entry that holds the two ai-memory messages of 2026-11-18 from gesttione-solutions-indexus, `01a1eec1-4be0-72e0-8e9f-4bcf7211f8a9` and `01a1eeca-397e-77c2-8429-3668a5a0b8f8`, and its Notice Request: the sender asked to be told when the work is done, with the final export list, the form and path of the Lesson Learned, the marking of ADRs and the projects already updated.
2. [`references/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md`](references/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md), the adopted Backlog Entry with the maintainer's list about the folders.
3. [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md), the routing, the six answers of the Idea stage, the approval of the Idea with its three defaults, and the answer about open questions, recorded by the operator.
4. [`references/2026-11-18-mirrored-documentation-measurement.md`](references/2026-11-18-mirrored-documentation-measurement.md), the count of the mirrors, of the ADR statuses and of what a rename would touch, recorded by the author of the Idea.
5. [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), recorded by the author of this PRD: the prototype of the Baseline text, the tests that pin the text it changes, the behavior of the sync with an exclusion, and the Governed Paths.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from: its features F1 to F10 are this Spec, and its ten open questions are settled in `## Decisions`. This project's ai-memory inbox was listed on 2026-11-19 and held no pending message. `docs/findings/` holds no Finding. Six Backlog Entries are open at the authoring base. One is the entry the Idea wrote for the rename of `docs/agents/` and the merge of backlog and findings, which this Spec leaves out. The other five are about three unreproduced symptoms of another team, the context an Authoring Phase starts with, an abort of the consistency check, a fixed deadline of the Baseline analysis and dead probe code; none shares a component with this Spec and none is adopted.

## Prerequisites

No behavior of this Spec depends on another Spec's code. It was authored against `0ea36244`, where Specs 0407 and 0408 are delivered and Specs 0409 and 0411 are planned, and it is delivered after Specs 0411, 0414 and 0409 and before Specs 0412 and 0413, in the order the maintainer set. These files are shared, and each Task regenerates the shared records on the tree it runs on, with no version number written by hand:

- Spec 0409 adds a rule to the Baseline module `spec-workflow`, which this Spec does not touch, and changes neither `core`, `context-workflow`, `secondbrain`, the docs-layout guide nor the agent-instructions guide. Both Specs move the module version record, the derived Baseline files, the Standard TypeScript profile's pin, this repository's Setup Manifest, the Roundfix skill's version and `modes` reference, the owned-skill version record, `CONTEXT.md`, the closed list of clause forces, one user-guide page and the Behavior Surface Record.
- Spec 0411 changes the Roundfix skill's `profiles`, `setup` and `modes` references, that skill's version, `CONTEXT.md` and the Behavior Surface Record. This Spec adds sentences of its own to `setup` and `modes`.
- Spec 0414 changes the Release Plan's maintenance boundary, the release sentence that four documents repeat, the Roundfix skill's `release` reference and `CONTEXT.md`. This Spec adds two paths to the same boundary and the same sentence, for every repository, and leaves the Adopter Maintenance Path of Spec 0414 as it is: the export list is counted by the general rule, so no adopter test is needed for it.
- Specs 0412 and 0413 are delivered after this one. Their decision records arrive with their plans and carry no Knowledge Mark; the Task that marks this repository's ADRs marks every ADR in the tree when it runs.

The plan of Spec 0414 was read on its plan branch, which had not merged at the authoring base.

A plan of this cycle that asks every Idea and PRD for a `## Rule challenges` section is delivered before this Spec. This PRD carries that section.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity. The two new Baseline rules and their clauses take ids in the form their modules already use, the new checks take codes with the `DR-` prefix the Doctor Command already uses, and a Lesson Learned is named by its date and a slug, as a Backlog Entry is. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. The update and the Doctor Command read and write files of the repository they run in; the sync that publishes the list belongs to another project and is not called. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — three new ADRs record the decisions: the list (ADR-0438: "The Baseline keeps the list, and only the list."), the two header fields (ADR-0439: "Only a marked ADR is published.") and the lesson (ADR-0440: "A lesson leaves when a rule takes it in."). The retirement rule stays (ADR-0278: "A retirement deletes."), and its paragraph on lessons changes (ADR-0278: "Lessons leave before the text."), so a Task writes the amendment notice of ADR-0440 into it. ADR-0438 changes what an earlier decision says of the export (ADR-0248: "A documentation contract holds both directions."), and both new release rules extend the planning-evidence decision (ADR-0412: "Planning references and the release-preparation commit set no release impact"); a Task writes those notices too. The history cleanup of the update is used as it is (ADR-0282: "The history cleanup names a superseded ADR in its successor"; ADR-0254: "Baseline update sanitizes pending history in the change it plans"), and a proposed ADR is never reported as retired (ADR-0122: "A proposed decision is pending, not retired"). Every new check warns and none fails (ADR-0220: "A definite answer fails; an unanswered question warns."), and a departure is named, never repaired by moving a file (ADR-0411: "A repository rule outside the specific-repository guide is named, never moved"). The new clauses are written for every adopter (ADR-0186: "Baseline guidance states what the product does, in adopter-neutral words"), the rule on the mark lives in the guide of the module it depends on (ADR-0222: "A Baseline guide says only what holds for the repository that reads it"), and each rule is stated once (ADR-0257: "One statement per rule."). The evidence comes from outside Roundfix: the request of another project and the count of fifteen mirrors (ADR-0203: "The Baseline grows only with evidence from outside Roundfix."). The standard list is not a Managed Region and is not recorded in the Setup Manifest (ADR-0145: "A managed region is trusted by its marker"; ADR-0287: "every settlement checks the Setup Manifest"), and the preservation proof of the instruction carriers is untouched (ADR-0100: "A managed refresh proves preservation instead of backing it up"). The lesson stops being a message (ADR-0330: "The Secondbrain is read-only and ai-memory carries handoffs and messages"), and the sender's Notice Request is answered after the release (ADR-0396: "receives one completion notice"). The moment of a lesson in a delivery is the declaration of a Correction (ADR-0417: "A Correction carries a declared Defect Class"; ADR-0420: "The corrections record is one append-only file per Spec under Roundfix Home"). The skills and the guides describe the changed commands in the same delivery (ADR-0256: "A release waits for the skills that describe a changed surface"), and the authorization record carries the Release Ceiling (ADR-0349: "A Spec's Release Ceiling answers the release plan's approval"). ADR-0092, ADR-0120, ADR-0247, ADR-0251, ADR-0258, ADR-0260, ADR-0261, ADR-0262, ADR-0277, ADR-0281, ADR-0290, ADR-0295, ADR-0299, ADR-0308, ADR-0342, ADR-0352, ADR-0366, ADR-0368, ADR-0391, ADR-0421, ADR-0422, ADR-0423, ADR-0424, ADR-0426 and ADR-0432 cite one of the ADRs above and are listed only to record why they do not apply. The TechSpec relates each ADR. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-18 for this Spec, quoted in the authorization record, which keeps the grant of every owned skill and of the Governed Paths each Spec declares, Baseline modules and records included. It grants no Makefile, workflow, Go module, lint or formatter change except for a correction that requires one, and the design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export/_authorization.md`; bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/corrections.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/release.md`, `.agents/skills/roundfix/references/setup.md`, `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, `docs/references/coverage-record.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/secondbrain.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

- A repository that applied the update publishes knowledge and no process: its mirror holds no Baseline guide, Archive Record, Backlog Entry, Spec folder, Finding or handoff.
- The list of what a repository publishes lives in the repository, where a person reads it and a check compares it, and no other project filters on its behalf.
- A reader of the first lines of an ADR knows whether a later decision changed it, and no superseded ADR stays in a tree unreported.
- A failure that taught something leaves a record that another project can find, and work that taught nothing leaves none.
- Each documentation folder has a stated job, and a repository that departs from it is told so without being stopped.

## User Stories

1. As the maintainer of several projects, I want each one to publish only knowledge, so that a search across them stops answering with copies of the same guide and with records of delivered work.
2. As the maintainer, I want the standard list to arrive with the Baseline update and to stay current without anyone editing it, so that no project keeps export rules of its own.
3. As the maintainer of one repository, I want to add a file of knowledge or drop a standard item with a reason, so that the standard never forces me to publish or hide the wrong thing.
4. As the maintainer, I want the update to show every path that stops being published before it writes anything, so that knowledge never leaves a mirror unnoticed.
5. As the maintainer, I want each ADR to say whether the knowledge base wants it, so that rules of the business are published and decisions about tooling are not.
6. As the maintainer, I want an Agent to propose the mark of every existing ADR in one change I review, so that I correct marks and do not write a thousand of them.
7. As an Agent that opens an ADR, I want its header to name the ADRs that amend it, so that I do not act on a decision a later one changed.
8. As an Agent in another project that meets a problem, I want to find what a project tried, why it failed and what solved it, with when that applies, so that the failure is paid for once.
9. As the maintainer, I want a lesson written only when there is one, in the change that taught it, where I see it in the Pull Request and can delete it, so that the lessons folder does not become another archive of process.
10. As a person or an Agent new to a repository, I want each documentation folder to hold one stated kind of document, so that I know where a design, a reference or a lesson goes.
11. As the maintainer of a repository, I want `roundfix doctor` to tell me what departs from these rules and what repairs it, without failing, so that a repository that has not caught up keeps working.
12. As the maintainer of an adopter, I want the commit of a Baseline update and a commit that only adds a lesson to ask for no release, so that documentation never forces a version.
13. As the project that asked for this, I want one notice when it is released, with the final list, the form of the lesson and the marking of ADRs, so that I can prepare the search index.

## Core Features

1. **The Baseline keeps a standard list of what is published.** In a repository whose Baseline has the Secondbrain turned on, `.secondbrain-export` holds the **Standard Export List** between two marker lines: the README, the glossary files of the repository's declared layout, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/`, and the path of each ADR whose mark says it is knowledge. Each standard folder is listed whether or not it exists yet. `roundfix baseline update` writes those lines and keeps them current, and the list is the same for the same files, whoever runs the update. A repository with the Secondbrain turned off gets no list and no file.
2. **Process is never published.** The root instruction files, `docs/agents/`, the Spec Root, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and the History Root are in no part of the list.
3. **A repository varies in its own lines.** Every line outside the two markers is the repository's **Export Variation**: a file or a folder of knowledge the standard lacks, or an exclusion that drops a standard item. Each such line carries its reason in a comment directly above it, and one comment stands for the lines that follow it without a blank line. No Baseline operation changes a line outside the markers.
4. **The first update converts an existing list once.** A file without the markers is read as the list a project copied. A line that names a standard item or a process path, `docs/` whole included, and an exclusion that lies under a process path are removed; every other line stays as it is, in its order, below the new block. The plan names each path that stops being published, says which are process and which are only absent from the standard list, and says how many ADRs carry no mark, before anything is written.
5. **A broken pair of markers is refused.** When one marker is missing or either is repeated, the update writes no line of the file, says which marker is wrong and applies the rest of its plan.
6. **Each ADR says whether the knowledge base wants it.** In a repository with the Secondbrain turned on, the header of an ADR carries the **Knowledge Mark**: true when the decision states a rule of the business or of the product that someone outside the repository's development can use, false when it is about how the repository is built, tested, organized or delivered, or about a tool it uses. The author of an ADR sets it and adds the ADR's path to the standard list when the mark is true. An ADR without the mark is not published.
7. **The existing ADRs are marked once, by an Agent, in a change a person reviews.** The Agent of the session that applies the Baseline update that brings this rule marks every existing ADR by the same test, in one commit, and runs the update again so the list names the ADRs marked true. The `update` Skill Mode of the Roundfix skill carries the step. No command decides a mark.
8. **An amended ADR says so in its header.** An accepted ADR that a later accepted ADR changes in part stays accepted and carries the **Amendment List**, which names the ADRs that amend it. The author of the later ADR writes it, with a dated notice in the body after the text that changed. `roundfix baseline update` fills the list of an existing ADR from the notices its body already carries: "Amended by", "Narrowed by", "Extended by" and "Superseded in part by", each followed by one or more ADR identifiers. Writing either header field changes no other line of the ADR, and an ADR without a header gains one that holds only that field and stays a legacy ADR.
9. **The rule against rewriting existing ADRs gains one exception.** The docs-layout guide keeps the sentence that forbids rewriting existing ADRs to adopt lifecycle metadata, and a header field that a guide of the repository defines for existing ADRs is the exception.
10. **A superseded ADR is still deleted.** The retirement rule and the history cleanup of the update are unchanged. The Doctor Command reports a rejected, deprecated or superseded ADR that is still in the tree; a proposed ADR is not reported.
11. **A lesson has a record.** A **Lesson Learned** is one file under `docs/lessons/`, named by its date and a slug, with the date it was created, a title of one line and eight labelled lines: the context, what was tried, what went wrong, the proven cause, what solved it, when the lesson applies, when it does not, and links to the evidence. The labels stay as the Baseline writes them, and the rest is in the language of the repository's documentation. Every repository on the Baseline has the rule, with or without the Secondbrain.
12. **A test decides that there is a lesson.** A Lesson Learned is written only when every part holds: something was tried, it went wrong, the cause is proven, something solved it, and no accepted ADR, reference document or glossary entry already states it. When any part fails, nothing is written.
13. **Two moments ask the question.** When a Finding closes, in the change that answers it or retires it, and when a Spec that needed a correction, a failed QA or a retry is archived. The Agent that applies the test writes the file in the change it is making, so the Pull Request shows it to a person, who may delete it. `archive-spec` asks before the archive and names the signs it reads: an adopted Finding, a QA Report that did not pass or more than one report, a Corrective Task, a declared Correction or an undeclared Correction Sign, and a retried Task or queue item. The Delivery Queue archives without an Agent, so the Supervisor session applies the test when it declares a Correction or retries an item after a repair, writes the lesson in the commit that corrects, and after the merge writes it in its next change. The Roundfix skill's `corrections` reference carries that step.
14. **The file replaces the message, and leaves when a rule takes it in.** A closing Finding no longer sends its lesson to another project. A Lesson Learned is deleted, under the retirement rule, once an accepted ADR, a guide or a skill states what it teaches.
15. **Each folder has a stated job.** The directory clause of the docs-layout guide says what `docs/design/` holds, names `docs/lessons/`, and says that a reference document belongs under `docs/references/` and not at the repository root, the `DESIGN.md` design contract excepted. A folder under `docs/` or a reference document at the root that the clause does not name is a variation: the repository names its path and its job in its own rules guide, or moves its content. No command renames or moves a file.
16. **The update shows and applies the documentation part behind its one confirmation.** The plan of `roundfix baseline update` shows the lines the standard list gains and loses, each path that stops being published and each ADR whose Amendment List it fills. The Plan Digest covers them, the one confirmation the update already asks for applies them, a second run reports nothing to do, and the Update Check counts them among the file changes. The `update` Skill Mode asks the user before it applies a plan that stops publishing a path.
17. **The Doctor Command reports what departs.** Two new checks, each a warning with a stable code and a next action, and neither ever a failure. The first, where the Secondbrain is turned on: no standard list, a list that differs from what the update would write, a broken marker, a process path in the repository's own lines, a line without its reason, and the number of ADRs without the mark. The second, in every adopted repository: a missing README or glossary, a folder or a loose file under `docs/` or a reference document at the root that the layout does not name and the repository's own rules guide does not name either, a retired ADR still in the tree, an amendment notice the header does not list, and a file under `docs/lessons/` that is off the form.
18. **Documentation asks for no release.** The Release Plan counts `.secondbrain-export` as maintenance-only and `docs/lessons/` as planning evidence, in every repository.
19. **This repository applies the rules to itself.** Its export list becomes the standard list with its own lines, which drop the glossary, `docs/references/`, `docs/design/` and every ADR, so that it publishes its README, its user guide and its lessons. Every ADR in its tree carries the Knowledge Mark, every ADR that carries an amendment notice lists it in its header, and the three ADRs this Spec changes carry their notices. The contract test of its export file follows the new rule.
20. **Skills, guides, glossary and decisions.** The Roundfix skill's `baseline`, `setup`, `modes`, `corrections` and `release` references and the `archive-spec` skill describe the behavior. The user guide describes the list, the mark, the Amendment List, the Lesson Learned and the two checks. `CONTEXT.md` adds **Standard Export List**, **Export Variation**, **Knowledge Mark**, **Amendment List** and **Lesson Learned**, and revises **Managed Refresh**. ADR-0438, ADR-0439 and ADR-0440 record the decisions.

## User Experience

The blocks below show the content of each surface with invented names. The TechSpec fixes the exact wording.

The export list of a product repository after its first update. The file had `CONTEXT.md`, `docs/`, `README.md`, `AGENTS.md`, `DESIGN.md` and a commented line:

```text
# roundfix:standard-export:begin
# Kept by the Baseline: roundfix baseline update rewrites the lines between these two markers.
README.md
CONTEXT.md
docs/user-guide/
docs/references/
docs/design/
docs/lessons/
docs/adr/0007-a-credit-note-never-reopens-a-closed-period.md
docs/adr/0012-a-price-change-takes-effect-at-the-next-opening.md
# roundfix:standard-export:end

# The design contract the interface follows.
DESIGN.md
# API contracts (uncomment when present)
# openapi.json
```

What the update says before it writes that file:

```text
Documentation: pending
Export list: update .secondbrain-export (8 line(s) in the Standard Export List)
- stops exporting AGENTS.md (process)
- stops exporting docs/agents/ (process)
- stops exporting docs/backlog/ (process)
- stops exporting docs/prompts/ (not in the standard list)
ADRs without a Knowledge Mark: 41 of 43; they are not exported
ADR headers: 1
- docs/adr/0003-prices-are-stored-in-cents.md gains amended_by: [ADR-NNNN]
```

The header of an ADR with both fields:

```markdown
---
status: accepted
knowledge: true
amended_by: [ADR-NNNN, ADR-MMMM]
created_at: 2026-03-02T00:00:00Z
updated_at: 2026-09-10T00:00:00Z
deprecated_at: null
superseded_by: null
---
```

A Lesson Learned:

```markdown
---
created: 2026-11-18
---

# A time without a zone is read as UTC

- Context: the report of sales by day, which groups each sale by its local date.
- Tried: reading the times of the sales feed with the date library's default parser.
- Went wrong: sales made after 21:00 appeared in the next day's total.
- Cause: the feed sends local time without a zone, and the default parser reads such a time as UTC.
- Solved by: naming the store's zone when parsing, with a test of a sale at 23:30.
- Applies when: a source sends a time without a zone.
- Does not apply when: the source sends an offset.
- Links: Pull Request 412; the Finding at `docs/findings/2026-11-10-sales-move-to-the-next-day.md@3f2a9c1`.
```

What the Doctor Command says of a repository that departs, and of one that does not:

```text
baseline-export: warn (2 departure(s) from the export rules; next: run roundfix baseline update, then edit the repository's own lines)
  DR-EXPORT-LIST-STALE: the Standard Export List differs from what roundfix baseline update would write: 1 line(s) to add, 0 to remove
  DR-EXPORT-REASON-MISSING: .secondbrain-export:14: openapi.json has no reason; add a comment line directly above it
baseline-docs: warn (2 departure(s) from the documentation rules; next: see each line; roundfix baseline update fills an Amendment List and removes a retired ADR)
  DR-DOCS-OUTSIDE-LAYOUT: DATABASE.md: a reference document at the root the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md
  DR-DOCS-LESSON-FORM: docs/lessons/2026-11-18-zone.md: the line "Cause:" is missing

baseline-export: ok (the Standard Export List is current; 2 line(s) of the repository's own)
baseline-docs: ok (the documentation follows the docs-layout guide)
```

### Paths that leave a mirror

A mirror holds what the list names at the last push to the default branch, so a path that leaves the list leaves the mirror at the next push.

- In every repository that applies the update: `AGENTS.md` and `CLAUDE.md` where the list had them, `docs/agents/`, `docs/history/`, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and the Spec Root, `docs/specs/` by default.
- Every ADR, until it is marked true and the list names it; an ADR marked false stays out for good.
- Every folder or file under `docs/` that the standard list does not name, until the repository adds it to its own lines. The mirrors of 2026-11-18 hold `docs/reference/`, `docs/_refacs/`, `docs/prompts/` and `docs/schema/`, one project each, `docs/_inbox/` in four projects, and three loose files.
- What stays: the README, the glossary files, `docs/user-guide/`, `docs/references/`, `docs/design/`, and every line of its own a repository already had, such as `DESIGN.md`, a changelog or an API contract.
- This repository: `AGENTS.md`, `CHANGELOG.md`, `CONTEXT.md`, `docs/adr/`, `docs/agents/`, `docs/backlog/`, `docs/design/`, `docs/findings/`, `docs/handoffs/`, `docs/history/`, `docs/references/` and `docs/specs/` leave its mirror. `README.md` and `docs/user-guide/` stay, and `docs/lessons/` joins them with its first lesson.

## Non-Goals / Out of Scope

- Renaming `docs/agents/` or merging backlog and findings: the maintainer chose to state each folder's job and rename nothing, and the measured cost is in the open Backlog Entry the Idea wrote.
- Moving a file by command: a document outside the layout is reported, and the repository moves it or names it.
- Rules for what a README or a user guide must say: this Spec covers the records that have a form and the files that must exist.
- A check that fails a Run, a merge, a settlement or the Spec Consistency Check: every new check is a warning of the Doctor Command, and `roundfix spec check` reads none of it.
- A contract in this repository that fails a plan whose ADR lacks the mark: plans of this cycle were written before the rule.
- A command or a judge that decides a Knowledge Mark, and a judge that filters lessons.
- Reading lessons during authoring: it needs lessons to exist and a way to pick the ones that apply.
- Keeping a superseded ADR with a status, and a new ADR status for an amended one.
- Changing the Secondbrain's template, sync or filter, and applying the update in another repository: each project applies its own update, and Roundfix edits no other project.
- Turning past Findings and Archive Records into lessons: the record starts with what closes after the release.
- A place for a changelog in the standard list: a repository that wants it published keeps it in its own lines.
- A new field in the Archive Record, a new command, a new flag and a new configuration key.
- Recording the standard list in the Setup Manifest or treating it as a Managed Region.
- A change to the Adopter Maintenance Path of Spec 0414.

## Success Metrics

1. Success Metric: a managed refresh of an adopted temporary repository writes into its docs-layout guide the directory clause with `docs/lessons/` and the variation rule, the ADR clause with `amended_by`, the exception to the lifecycle rule, the retirement of a Lesson Learned and the two clauses of the lesson rule, each once; into its Secondbrain guide the three clauses of the export rule; and into its agent-instructions guide the sentence that a closing Finding sends no message. A second refresh reports it current, and this repository's own three guides hold the same text after the Spec's Managed Refresh.
2. Success Metric: on a temporary repository with the Secondbrain turned on and an export file that lists `docs/` whole, `AGENTS.md`, one file of its own and one comment, the update's plan names the standard list, every path that stops being exported with its kind and the count of unmarked ADRs; the confirmed update writes the block, removes the lines it named and leaves every other line byte for byte; a second run reports nothing to do; the same repository with the Secondbrain turned off gets no file; and a file with one marker missing is left untouched while the rest of the plan applies.
3. Success Metric: on a temporary repository the update fills the Amendment List of an ADR from each of the four notice forms, with one and with two identifiers, ignores a notice inside a code block, changes no other byte of the file, gives an ADR without frontmatter a header that holds only the list, and reports nothing to do on a second run.
4. Success Metric: on temporary repositories the Doctor Command prints each of the eleven codes for its seeded departure as a warning with a next action, exits as it did before, prints `ok` for both checks on a repository that follows the rules, skips the export check where the Secondbrain is turned off, and does not report a folder or a root document that the repository's own rules guide names.
5. Success Metric: the Release Plan counts a commit that changes only `.secondbrain-export` and a file under `docs/lessons/` as maintenance-only, and still blocks a commit that also changes another path.
6. Success Metric: in this repository at the delivery, the Doctor Command prints `ok` for both checks: every ADR in the tree carries the Knowledge Mark, every amendment notice is in a header, no retired ADR is in the tree, and the export list is the standard list with this repository's own lines, each with its reason; ADR-0248, ADR-0278 and ADR-0412 carry their dated notices.
7. Success Metric: both copies of the `archive-spec` skill carry the lesson step with its five signs, and both copies of the Roundfix skill's `corrections`, `modes`, `baseline`, `setup` and `release` references carry their sentences, each pinned by its own phrase in a contract test that fails when the phrase is removed.
8. Success Metric: `CONTEXT.md` defines the five new terms and the revised one, the user guide describes the list, the mark, the Amendment List, the Lesson Learned, the documentation part of the update and the two checks, and the documentation, skill and glossary contracts pass.
9. Success Metric: every Baseline clause that exists at the delivery base keeps its id and its force, the docs-layout region of the Standard TypeScript golden is at most 13,000 bytes, and every test of the Baseline package that exists at the delivery base passes, the three tests that pinned the changed text having been updated to the new text.
10. Success Metric: 30 days after the release, the mirror of every repository that applied the update holds no file under `docs/agents/`, `docs/history/`, `docs/backlog/`, `docs/specs/`, `docs/findings/` or `docs/handoffs/`, from 1,509 on 2026-11-18, and this repository's mirror holds fewer than 60 files, from 851.
11. Success Metric: of the Lessons Learned written in this repository in the first 60 days, at least 80% are still in the tree or were taken in by a rule, and between 1 and 10 are written for each 100 archived Specs and closed Findings.

## Acceptance evidence

The outside-evidence rows rest on a measurement this Spec did not design and on published sources it did not write. Each published passage was read for the Idea on 2026-11-18 and read again for this PRD on 2026-11-19.

- The count of the project mirrors in [`references/2026-11-18-mirrored-documentation-measurement.md`](references/2026-11-18-mirrored-documentation-measurement.md), taken from the Secondbrain at its commit `a2b713574`: 1,509 of 3,534 documentation files are process, and no mirror holds a folder of lessons. The project that reads the mirrors counted 4,150 files of 11 projects on the same day with the same proportion. It supports Core Features 1 and 2.
- Nick Milton, "Leaving lessons in a lessons database doesn't work - an example from NASA", 2017 (http://www.nickmilton.com/2017/07/leaving-lessons-in-lessons-database.html), quoting the audit of NASA's lessons system: "We found that NASA program and project managers rarely consult or contribute to LLIS even though they are directed to by NASA requirements and guidance." His conclusion: "Once the lesson has been embedded as a practice change, or a design-principle change, or a change in a checklist, then it can be removed from the database." It supports Core Feature 14 and stands as a dissent against storing lessons at all.
- Lindenbauer, Groh and Schütze, "From Knowledge to Noise: CTIM-Rover and the Pitfalls of Episodic Memory in Software Engineering Agents", 2025 (https://aclanthology.org/2025.realm-1.30.pdf): "We find that CTIM-Rover does not outperform AutoCodeRover in any configuration", and "Our analysis indicates noise introduced by distracting CTIM items or exemplar trajectories as the likely source of the performance degradation." It supports a test that writes nothing by default (Core Feature 12).
- Mehmet Iscan, "Learning When to Remember: Risk-Sensitive Contextual Bandits for Abstention-Aware Memory Retrieval in LLM-Based Coding Agents", 2026 (https://arxiv.org/abs/2604.27283): "This memory is valuable only when the current failure is genuinely compatible with a previous one." It supports the two lines that say when a lesson applies and when it does not (Core Feature 11).
- "How to deal with conflicting ADRs? AI picks what suits it best", a public discussion of 2026 (https://github.com/mattpocock/skills/discussions/847): "An agent greps the ADR folder, opens a hit, reads the first few lines and stops. If the supersede note is buried in the body, it never sees it". It supports Core Feature 8.
- The Ansible documentation of `blockinfile` (https://docs.ansible.com/ansible/latest/collections/ansible/builtin/blockinfile_module.html): a marker that does not tell the begin line from the end line "may result in the block being repeatedly inserted on subsequent playbook runs." It supports the pair of markers and the refusal of Core Feature 5.

## Unreachable Acceptance

- criterion: after the release, this repository's mirror in the Secondbrain holds fewer than 60 files once the first push to the default branch has synced.
  reason: the sync runs in another project's workflow, over the network, after the merge.
  satisfied-by: the operator's count of the mirror after the release, recorded session evidence.
- criterion: Success Metric 10 holds for the adopters 30 days after the release.
  reason: it needs each adopter to apply its own update and the Secondbrain to sync, weeks after the release, in other repositories.
  satisfied-by: section 1 of the measurement in `references/2026-11-18-mirrored-documentation-measurement.md` repeated by the operator at a cycle close, recorded session evidence.
- criterion: in an adopter, the Agent of the update session marks the existing ADRs, runs the update again, and a person reviews the marks in one Pull Request.
  reason: it needs a live Agent session, a person and another repository, which no Verification or QA row may start.
  satisfied-by: the first adopter's update after the release, reported by its project, recorded session evidence.
- criterion: Success Metric 11 holds over the first 60 days.
  reason: it needs weeks of deliveries and closed Findings after the release.
  satisfied-by: the Git history of `docs/lessons/` in this repository read by the operator at a cycle close, recorded session evidence.
- criterion: the sync drops what this repository's own lines exclude, the line `!docs/adr/*` included, with the rsync its workflow runs.
  reason: the exclusion was measured during authoring with the rsync of the author's machine, and the workflow of another project runs a different one.
  satisfied-by: the content of this repository's mirror after the release, read by the operator, recorded session evidence.

## Validation

The decisions were checked on 2026-11-19 against these sources.

- Characterization at `0ea36244`. No Baseline module names `.secondbrain-export`, and its only reader in this repository is one contract test, which holds the history exclusion the History Root made moot. The update builds its guidance plan first and then plans the history cleanup beside it, folds that plan into the Plan Digest and applies it after the guidance; the documentation part attaches at the same place. Every settlement checks that the guidance plan of the Setup Manifest has no change, which is why the standard list stays out of that plan. The readers of an ADR's header decode the fields they know and ignore the rest, so the two new fields break none. The Doctor Command builds each check as one result with a status, a detail, a next action and lines, and only a failed status changes its exit. The Release Plan lists six `docs/` prefixes as planning evidence, and neither `.secondbrain-export` nor `docs/lessons/` is among them, so without Core Feature 18 the commit of an update would ask an adopter for a classification by hand.
- A prototype of the Baseline text on a scratch worktree of `0ea36244`, removed afterwards and recorded in [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md). The docs-layout rule is one the Standard TypeScript profile requires, so its existing clauses were reworded under their ids and the lesson rule and the export rule are new rules, which need no Source Baseline row. The module record command, `make baseline-digests` and the managed refresh rewrote three modules, the module version record, the Standard TypeScript profile and three of its formatter goldens, three catalog files, four plan goldens and this repository's three guides and Setup Manifest, and a second refresh reported no change. Three existing tests pin text this Spec changes: the test of the sentence a closing Finding sends, the test of the directory clause, and the ceiling of 11,000 bytes on the docs-layout region, which the new text takes to 12,948. A fourth pins the sentence that forbids rewriting existing ADRs, and the new clause keeps that sentence whole.
- The sync, read in the Secondbrain at its commit `4a496f4a4`: it reads one path per line, treats a `!` line as an rsync exclusion and skips a listed path that does not exist with one warning. Measured on the author's machine with its rsync: an exclusion `!docs/adr/` does not drop an ADR listed by its own path, `!docs/adr/*` does, `!docs/references/` drops that folder and `!CONTEXT.md` drops that file.
- Secondbrain, read only on 2026-11-19: `wiki/index.md`; `wiki/concepts/politica-de-retirada-e-poda.md`, which says that capture at closing is conditional, "só quando a lição não está na ADR ou no record absorvente", and which set the test and the exit of a lesson; `wiki/entities/indexus.md`, whose own sources are mirror paths under `docs/history/` and `docs/findings/`, the risk named below; `brain/RUNBOOK.md`, the section on sending to the Indexus, for the filter of 2026-11-18 ("O espelho no repositório não muda."); and `scripts/sync-project.sh` with `templates/repo/secondbrain-export.example`, for the contract of the file. `qmd query` was tried three times and refused each time with "database is locked", while an index update of another session ran; no query result informed this PRD, and the four queries of the Idea stand. The digest of this research for the project that owns the Secondbrain was not sent by the author, who sends nothing in this cycle.
- Exa: the passages of `## Acceptance evidence`, each read again, and for the design the documentation of `adr-tools`, whose `adr link` command writes the pair "Amends" and "Amended by" into the two records, and of MADR, whose only relation in a header is `superseded by`; both show that an amendment list in the header is this Baseline's own extension. The sweep found no controlled study that compares a list kept at the source with a filter at the reader.

## Risks

- The first update of a repository that exports `docs/` whole stops publishing every folder the layout does not name. The plan names each one and the `update` Skill Mode asks first, and a person who confirms without reading still loses it from the mirror until a line is added.
- Pages of the Secondbrain cite mirror paths under process folders, and lose their source when a project applies the update. The project that owns the Secondbrain is told when this plan merges, before any repository can apply the update.
- Between a repository's update and the merge of its marks, it publishes no ADR. The same session does both, and a change with hundreds of marks may still wait for review.
- A mark is a judgment made in bulk. A rule of the business marked as tooling is hidden from the knowledge base, and the review of one large change is the only guard.
- A person who edits inside the markers loses the edit at the next update, and one who deletes a marker stops the list from being kept until the pair is restored. The Doctor Command reports both.
- The docs-layout guide grows by about 2,500 bytes and the Secondbrain guide by about 2,100, which every session that works on documentation reads. The ceiling the lean Baseline set for the docs-layout region rises from 11,000 to 13,000 bytes.
- A lesson written in a hurried review can be weak, and the published measurements say weak lessons hurt.
- A correction the queue makes with no session is asked about only after the merge, so its lesson reaches a later Pull Request.
- The list names each published ADR by its path, so a new ADR costs one line in two files, and in this repository those lines are excluded by its own variation.
- The Spec changes three Baseline modules, two owned skills, two commands and about 370 decision records of this repository. It is the largest of the cycle, and its last Task alone changes every ADR.
- Dissent kept from the Idea: the author recommended publishing every accepted ADR, because a decision absent from the knowledge base fails in silence. Milton's position is that the right number of stored lessons is zero. And common ADR practice keeps a superseded record in the log.

## Idea-stage judgment

- item: ai-memory message `01a1eeca-397e-77c2-8429-3668a5a0b8f8`, kept in [`references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md`](references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: unsure (0.68)
- unverified_assumption: no (0.17)
- market_question: no (0.06)
- verdict: undecided
- stage: ran
- shape: changed — the export list became a standard block with a declared variation, each ADR gained a mark that decides its export, a superseded ADR is still deleted while an amended one says so in its header, the Lesson Learned gained a test, two moments and an exit, the folders keep their names, and the Doctor Command reports while `roundfix baseline update` applies.

The verdict was undecided, and the maintainer chose "Spec com etapa de Ideia (Recommended)" on 2026-11-18. He approved the Idea, with its `shape` line, on 2026-11-18: "Aprovar e seguir ao PRD (Recommended)".

## Rule challenges

None. Two rules stand in the way of this goal, and the maintainer decided the change of each before this PRD, so no question was asked. The docs-layout clause that forbids rewriting existing ADRs for lifecycle metadata gains an exception by his answers "Agente marca, você vê no PR (Recommended)" and "Apagar substituída, emendada no topo (Recommended)", confirmed in the approval of the Idea, whose option text names the "exceção à regra 'não reescrever ADR só para metadados'". The sentence of the cross-project clause that sends a closing Finding's lesson as a message is replaced by his answers on the Lesson Learned. The ceiling of 11,000 bytes on the docs-layout region is a third rule this Spec changes; it follows from the request, which puts these rules in that guide, and is recorded under the delegated decisions.

## Glossary

- adds: **Standard Export List**
- adds: **Export Variation**
- adds: **Knowledge Mark**
- adds: **Amendment List**
- adds: **Lesson Learned**
- changes: **Managed Refresh**

## Decisions

Every dated line is the maintainer's answer, as recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The text in quotes is the option he chose, verbatim.

- 2026-11-18: How is the request of the Indexus routed, and under which authorization? — "Spec com etapa de Ideia (Recommended)"; "Mesmas do BE (Recommended)": push, Pull Request, merge and release after the Spec, up to five corrections, and the Makefile, the CI workflows, `go.mod` and the lint configuration only for a correction that requires one.
- 2026-11-18: Where does the Backlog Entry about the folders go, and in which order is the Spec delivered? — "Juntar à Spec do Indexus (Recommended)"; "Corretiva → regras → Indexus (Recommended)".
- 2026-11-18: What does the Baseline do about the folders? — "Só a função de cada pasta (Recommended)": the Baseline says what each folder holds, `docs/design/` included, and where files such as `DATABASE.md` go; the Doctor Command points out what is outside; nothing is renamed or moved; the rename and the merge return to the backlog with the measured cost.
- 2026-11-18: May a repository's list differ from the standard? — "Padrão + variação declarada (Recommended)": the Baseline writes the standard block and the update keeps it; the repository may drop items or add files of knowledge in a block of its own, with the reason in one line; the Doctor Command points out a process folder exported and a standard block altered.
- 2026-11-18: Which ADRs are published? — "Só as marcadas": each ADR gains a field that says whether it interests the knowledge base. It was not the author's recommendation.
- 2026-11-18: Who marks? — "Agente marca, você vê no PR (Recommended)": the author marks an ADR when writing it; for the existing ones the update generates the marking in one Pull Request per repository, which the maintainer can correct; an ADR without the mark is not exported.
- 2026-11-18: What of an ADR that no longer holds? — "Apagar substituída, emendada no topo (Recommended)": the rule that deletes stays; the Doctor Command points out the superseded one that remained and the update removes it; an amended ADR gains in its header the list of the ADRs that amend it.
- 2026-11-18: When is a Lesson Learned written? — "Finding fechado e archive c/ tropeço (Recommended)": a Finding that closes and an archived Spec that had a correction, a failed QA or a retry; with no reusable lesson, nothing is written.
- 2026-11-18: Who decides the lesson? — "Agente grava, você vê no PR (Recommended)": the Agent applies the Baseline's test and writes in the same Pull Request; whoever reviews can delete.
- 2026-11-18: Is the Idea approved? — "Aprovar e seguir ao PRD (Recommended)": one Spec, four or five Tasks and the QA gate, the fifth being this repository applying the rules to itself; the standard list of the README, `CONTEXT.md`, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and only the marked ADRs; the exception to the rule against rewriting an ADR for metadata.
- 2026-11-18: Who proposes the mark of the existing ADRs? — "O Agente da sessão do update (Recommended)".
- 2026-11-18: Does this repository publish its lessons? — "Sim, mais docs/lessons/ (Recommended)": Roundfix exports `README.md`, `docs/user-guide/` and `docs/lessons/`.
- 2026-11-18: Is `docs/design/` in the standard list? — "Sim, na lista padrão (Recommended)".
- 2026-11-19: How are the open questions of the PRD and the TechSpec settled with nobody present? — "Seguir o padrão do autor (Recommended)", recorded in the operator's brief of 2026-11-19: for a question that is not a maintainer boundary, the author takes the option it recommends and records it as an operator decision delegated by the maintainer.

Autonomous route: no person present; decisions taken from the approved Idea and its defaults, the maintainer's answers above, the operator's brief of 2026-11-19, the characterization and the prototype recorded in `## Validation`.

Each line below is an operator decision delegated by the maintainer (2026-11-19), with the question, the options and the one taken.

- How do the marked ADRs become paths? Options: (a) the standard list names each ADR marked true by its path, the author adds the line, the update rewrites the block and the Doctor Command reports a difference (taken: the sync reads paths and belongs to another project, and it is the Idea's default); (b) `docs/adr/` whole with an exclusion for each ADR not marked (refused: a new ADR without a mark would be published until someone excluded it).
- Does the list name a standard folder that does not exist? Options: (a) every standard folder, always (taken: the sync skips a missing path with one warning, the list is the same in every repository, and the first lesson is published by the push that merges it); (b) only the folders that exist, as the Idea's default said (refused: `docs/lessons/` created with a first lesson would stay unpublished until someone ran the update).
- How does the Baseline find its lines in a file the repository also edits? Options: (a) two marker lines, written as comments the sync ignores (taken: it is how a tool tells its lines from a person's, and a broken pair can be refused); (b) no markers, with the standard lines matched by content (refused: the Baseline could not tell a standard line it wrote from one a person kept, and could never remove a retired standard item).
- What does the first update do with the lines a repository already had? Options: (a) remove the standard and process lines and the exclusions under a process path, keep every other line below the block, and name every path that stops being exported (taken: the Idea's default, made exact); (b) keep the old file and add the block (refused: `docs/` whole would go on publishing process).
- What does the update do when a marker is missing or repeated? Options: (a) write no line of the file, say which marker is wrong, apply the rest of the plan (taken); (b) rebuild the block from the surviving marker (refused: it would guess which lines are the repository's); (c) fail the whole update (refused: a line of documentation would stop a Baseline update).
- Where does the rule of the Knowledge Mark live, and who marks in a repository without the Secondbrain? Options: (a) in the Secondbrain guide, so a repository without it marks nothing (taken: the mark answers a question only that repository has); (b) in the ADR template of the docs-layout guide, for every repository (refused: a guide would state a rule that does not hold for its reader).
- Do the mark and the Amendment List move `updated_at`? Options: (a) no, the field is added and no other line changes (taken: the decision did not change, and about 370 files would otherwise change two lines each); (b) yes.
- Does an ADR without any header gain one? Options: (a) yes, holding only the field written, and it stays a legacy ADR whose status is read from its body (taken: the Idea's default); (b) no, it is reported and left unpublished.
- When does a Finding close, for the lesson question? Options: (a) with the change that answers it, which for a Finding a Spec adopted is that Spec's archive, or with its retirement without one (taken: only then are the cause and the solution both known); (b) when its status becomes `done`, which is when a Spec adopts it (refused: nothing has solved it yet).
- Who asks the lesson question in a Delivery Queue, which archives without an Agent? Options: (a) the Supervisor session, when it declares a Correction or retries an item after a repair, writing in the commit that corrects, and after the merge in its next change (taken: the declarer has just named the cause, and the Idea's default for a late Correction); (b) the QA gate (refused: its commit holds QA evidence only, and a Correction can come after it); (c) the queue itself (refused: it runs no Agent at the archive).
- How does the archiving Agent know that the Spec stumbled? Options: (a) from the adopted sources, the QA Reports in the Spec folder, the Task Graph, the Corrections Report and the retries of its Run, with no new field in the Archive Record (taken: the Idea's default); (b) a count in the Archive Record (refused: a new field for a question asked once).
- Is the lesson form seven lines or eight? Options: (a) eight, with "Applies when" and "Does not apply when" apart (taken: each can be checked and neither can be left out unnoticed); (b) seven, with both in one line, as the request lists them.
- In which language is a Lesson Learned written? Options: (a) the language of the repository's documentation, with the labels as the Baseline gives them (taken: the Idea's default); (b) English always.
- How does a repository keep a folder or a root document the layout does not name without a warning for good? Options: (a) it names the path in its own rules guide, and the Doctor Command then does not report it (taken: the Baseline already sends every variation of a repository there); (b) every such path is reported for good (refused: a warning that never clears teaches people to ignore the check); (c) a configuration key (refused: a new key for what a guide already records).
- Which root documents does the layout name? Options: (a) the README, the root instruction file and its link, the glossary files and their maps, the changelog, the licence, the contribution, security, conduct and support files and `DESIGN.md`, and only a document whose name is in capitals is read as a reference document (taken: it matches the files seen in the mirrors and leaves tool files alone); (b) every Markdown file at the root.
- How many checks does the Doctor Command gain? Options: (a) two, one for the export rules where the Secondbrain is turned on and one for the documentation rules, each with a code per departure (taken); (b) one check per departure, eleven more lines in every report; (c) one check for everything.
- Does the commit of an update ask an adopter for a release? Options: (a) no: the Release Plan counts `.secondbrain-export` as maintenance-only and `docs/lessons/` as planning evidence, in every repository (taken: neither is shipped by any product, and without it this Spec would undo for the export list what Spec 0414 does for the files the Baseline manages); (b) add the file to the Adopter Maintenance Path (refused: it would depend on Spec 0414's code and hold only for adopters); (c) leave it (refused).
- Does the docs-layout guide keep its ceiling of 11,000 bytes? Options: (a) the ceiling rises to 13,000 and the new text is kept short, at 12,948 bytes (taken: the request puts these rules in that guide, a second guide would need a pointer in every root instruction file, and the lean Baseline of 2026-10-16 measured 14,494 before it was trimmed); (b) a guide of its own for lessons (refused for that pointer and for a new managed file in every adopter); (c) trim existing clauses to pay for the new ones (refused: each existing sentence carries a guarantee a test pins).
- Does this repository's contract test require every ADR to carry the mark? Options: (a) no: the test holds that this repository exports no process path, and the Doctor Command reports an unmarked ADR (taken: plans written before the rule would fail a contract they could not have known); (b) yes.
- What does this repository's own variation drop, and how? Options: (a) `CONTEXT.md`, `docs/references/`, `docs/design/` and every ADR, each by an exclusion line with the maintainer's answer as its reason, and `AGENTS.md` and `CHANGELOG.md` leave the file (taken: it is exactly "README.md, docs/user-guide/ e docs/lessons/"); (b) mark every ADR of this repository false (refused: the mark says what a decision is, not what a repository publishes).
- What does the Notice to the Indexus carry, and who tells the Secondbrain? The operator's session sends one completion notice when the release is published, with the final list, the form and path of the Lesson Learned, the two header fields and this repository as the one updated, because each project applies its own update. The operator tells the project that owns the Secondbrain when this plan merges that mirrors will shrink and which paths leave them. The author sent nothing.
- Tasks stage: who runs this repository's first Managed Refresh with the documentation part? Options: (a) the Task that adds the documentation part runs it as its last step, which converts this repository's export file and fills the Amendment Lists its ADRs' notices call for, and the Task that applies the rules to this repository writes the own lines, the marks and the notices and refreshes again (taken: once the part exists, this repository's own update plans a change, and the two tests that run the update on this repository and expect `current` would fail that Task's settlement); (b) leave every change to the last Task and let those two tests fail until then (refused: every Task leaves the repository Verification green); (c) make the update skip this repository until the last Task (refused: a special case in the product).
- On what outside evidence does the acceptance rest? The measurement adopted here and the published passages of `## Acceptance evidence`.
- The Release Ceiling is `minor`: every adopter receives new clauses and a file the update writes, the update gains a part of its plan and a JSON member, the Doctor Command gains two lines, and no command, flag, configuration key or exit code changes or is removed.
- The approved Idea changed in two places only: its status became `promoted`, and the two invented ADR numbers of its header illustration, which are numbers of real ADRs of this repository, became placeholders, because the consistency check reads every ADR number a Spec file cites.
- ADR-0438, ADR-0439 and ADR-0440 are written with this PRD. The amendment notices in ADR-0248, ADR-0278 and ADR-0412 are written by a Task.

## Open Questions

None. The measurements of Success Metrics 10 and 11 are the operator's, after the release.
