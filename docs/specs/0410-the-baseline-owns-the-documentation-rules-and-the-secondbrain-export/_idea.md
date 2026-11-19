---
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
created: 2026-11-18
status: promoted
leverage: compounding
---

# The Baseline owns the documentation rules and the Secondbrain export

## Overview

The Baseline says, for every repository that adopts it, what each documentation folder holds, which of that is knowledge, and what the repository publishes to the Secondbrain. Four things are new:

- a standard export list, which the Baseline writes into `.secondbrain-export`, and which publishes knowledge and keeps process at home;
- a Lesson Learned, a short record of something tried, what went wrong, why, and what solved it, written only when there is one worth reusing;
- a mark on each ADR that says whether the knowledge base wants it, and a line at the top of an ADR that a later one changed;
- checks in `roundfix doctor` that report a repository that departs from these rules, while `roundfix baseline update` applies what a command can apply.

It is for the maintainer of the projects that adopt the Baseline, who searches all of them through one knowledge base: the Indexus, which reads each project's mirror in the Secondbrain. It is also for the Agent that looks there for what another project already learned, and for the sessions that keep the Indexus and the Secondbrain, which today remove on their own side what the projects should not have sent.

It is valuable because nothing in the Baseline says today what a project publishes. Each project exports all of `docs/`, so 1,509 of the 3,534 documentation files in the mirrors are process: copies of the Baseline's own guides, records of delivered Specs, backlog and transient investigations. What a project learned the hard way is mostly absent, because the Finding that held the cause is deleted when it closes.

V1 is text in the Baseline, one block the Baseline keeps in one file, new checks in the Doctor Command and steps in two owned skills. No folder is renamed and no file is moved. The value grows with use: each lesson and each marked ADR adds to what every project can find, and each update keeps process out without anyone filtering by hand.

## Problem

The Indexus relayed the maintainer's decision on 2026-11-18. In the words of [the adopted Backlog Entry](references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md), Roundfix defines through the Baseline "the content of documentation files, the folder layout, the required files, what is kept and what is not, and what each project exports to the Secondbrain", and "no other project creates these rules on its own". On lessons the maintainer's own words are: "Pensar em como fazer as lições aprendidas serem algo de valor para ir para o secondbrain para serem ingeridas no indexus", and "Só seria útil se fosse algo como: tentamos usar x lib e gerou problemas tal e resolvemos assim."

Four gaps stand between that decision and the repository today. Line numbers are of this repository at `d1c43f9ef26c64bf75000c71699e8bdaae784cf6`, and the counts are in [the measurement](references/2026-11-18-mirrored-documentation-measurement.md).

The export list has no owner in the Baseline. No Baseline module names `.secondbrain-export`, and no code of this repository reads or writes it apart from one test of this repository's own file (`internal/docscontract/secondbrain_export_test.go:16`). Each project copies a template the Secondbrain keeps, which lists `docs/` whole and `AGENTS.md`. On 2026-11-18 the Secondbrain began to remove process files when it sends the mirrors to the Indexus. That is a second list, kept in another project, which no repository sees and which `roundfix doctor` cannot check.

A lesson has a rule and no record. The only rule is one sentence: "When a Finding closes, state whether its change teaches a lesson that its ADR, Archive Record or reference does not already state. Only then send one message with its source pinned" (`docs/agents/agent-instructions.md:86`). The lesson has no form and no file, and it reaches the knowledge base only if a session remembers to send it and another session files it. The archive of a Spec asks nothing about lessons; its skill only copies the files a later reader will need (`.agents/skills/archive-spec/SKILL.md:337`). Meanwhile the Finding that held the proven cause leaves the tree when it closes (`docs/agents/docs-layout.md:181`). The one time anyone counted, 7 of 383 retired files held a lesson that no ADR states (`docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md:11`), and a session rescued them by hand.

An expired decision can be read as current in two ways. The first is a superseded ADR that stays in the tree. The rule deletes it once its successor names it (`docs/agents/docs-layout.md:181`) and `roundfix baseline update` removes it (`internal/spec/history_entries.go:125`), yet two projects keep four of them, the Indexus among them, and nothing reports it. The second is larger: an ADR that a later one changed in part stays `accepted`, and the notice is a line in its body. In this repository 64 of 364 ADRs carry such a line, and in 63 of them it comes after the first twelve lines of the file.

The folders are named in a few words each (`docs/agents/docs-layout.md:123`), and [the maintainer listed on 2026-11-16](references/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md) where that leaves a repository guessing: `docs/design/` takes an API design in one repository and a front-end design in another, and reference files such as `DESIGN.md` and `DATABASE.md` sit at the repository root. The Baseline adds to the confusion, because its own frontend module places `DESIGN.md` at the root (`internal/baseline/assets/modules/frontend.json:297`). `roundfix doctor` checks none of this: its Baseline checks are the recorded version, a repository rule outside its file, and the `CLAUDE.md` link (`internal/cli/doctor.go:104`, `:135`, `:136`).

The author met the first gap while researching this Idea. A Secondbrain query about the job of each documentation folder returned the same Baseline guide, `docs/agents/docs-layout.md`, from seven project mirrors among its first twelve results.

The workarounds do not close the gaps. The Secondbrain's filter cleans only what the Indexus receives, and the mirror, which other sessions also search, stays as it is. A lesson sent as a message depends on two sessions and leaves nothing in the repository that learned it. Reading every ADR to its last line is what the notice in the body asks of a reader, and an Agent does not do it.

### Market Data

- A store of lessons that nobody curates stops being read. Nick Milton quotes the 2012 audit of NASA's lessons system: "We found that NASA program and project managers rarely consult or contribute to LLIS even though they are directed to by NASA requirements and guidance." His conclusion: "Lessons should not be stored in an ever-accumulating database - lessons need to be embedded into design, into principles and into process", and once a lesson is embedded "it can be removed from the database" (Milton, 2017).
- Lessons without a filter make a coding agent worse. A memory of lessons distilled from past tasks "does not outperform AutoCodeRover in any configuration", and the authors name "noise introduced by distracting CTIM items or exemplar trajectories as the likely source of the performance degradation" (Lindenbauer, Groh and Schütze, 2025).
- A lesson helps only where it applies: "This memory is valuable only when the current failure is genuinely compatible with a previous one", and the author describes fixes "solved" by earlier entries "that were no longer applicable" (Iscan, 2026). It supports the two fields that say when a lesson applies and when it does not.
- An Agent reads the top of an ADR: "An agent greps the ADR folder, opens a hit, reads the first few lines and stops. If the supersede note is buried in the body, it never sees it" (a public discussion of 2026 about an ADR format with the same optional status line as `.agents/skills/domain-modeling/ADR-FORMAT.md:21`). One reply asks for the status "at the very top of the file, where an agent cannot miss it".
- A curated list for machine readers is an established convention: "Agents are best served by concise, expert-level information gathered in a single, accessible location" (Howard, the llms.txt proposal, 2026).
- A linter of repository structure that also repairs exists and is archived: "By default Repolinter will automatically execute fixes as specified by the ruleset" (Repolinter). This Idea keeps the Doctor Command read-only and leaves every change to a plan a person confirms.
- The web sweep found no controlled study that compares a list of knowledge kept at the source with a filter at the reader. That choice rests on the maintainer's decision and on the Indexus's own count.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The standard export list | Critical | The Baseline writes a standard block into `.secondbrain-export`, and `roundfix baseline update` keeps it. In: `README.md`, `CONTEXT.md`, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and the ADRs marked as knowledge. Never in: `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, `docs/history/`, `docs/backlog/`, `docs/specs/`, `docs/findings/` and `docs/handoffs/`. A repository that has the Secondbrain turned off in its Baseline gets no block. |
| F2  | A declared variation | Critical | Below the standard block the repository keeps its own lines: it may drop a standard item or add a file or a folder of knowledge, each with its reason in one line. No Baseline operation changes those lines. On the first update, a line the repository already had that is neither standard nor process stays there as it is. |
| F3  | Each ADR says whether the knowledge base wants it | Critical | A field in the ADR's header takes one of two values: the decision is knowledge of the business or of the product, or it is about the tool, the code or the process. The author sets it when writing the ADR. An ADR without the field is not exported. |
| F4  | The existing ADRs are marked once per repository | High | The update that brings this rule to a repository leaves every ADR already there with its mark, in one change that a person reviews and corrects as one Pull Request. |
| F5  | The Lesson Learned record | Critical | One file per lesson in `docs/lessons/`, named by date and slug, in a short fixed form: context, what we tried, what went wrong, the cause, what solved it, when it applies and when it does not, and links. It leaves the tree when an ADR, a rule or a skill takes the lesson in, and Git keeps the text. |
| F6  | The test that decides there is a lesson | Critical | The Baseline states the test: something was tried, it went wrong, the cause is proven, something solved it, and no ADR, reference or glossary entry already says it. The Agent that applies the test writes the file in the change it is making, so a person sees it in the Pull Request and can delete it. When the test fails, nothing is written. |
| F7  | Two moments ask the question | High | When a Finding closes, and when a Spec that needed a Correction, a failed QA or a retry is archived. The closing of a Finding stops sending the lesson as a message, because the mirror carries the file. |
| F8  | An expired decision says so, or leaves | High | The header of an ADR lists the ADRs that amend it, and `roundfix baseline update` fills the list from the dated lines existing ADRs carry in their body. The rule that deletes a superseded ADR stays: the Doctor Command reports one still in the tree, and the update removes it as it does today. |
| F9  | Each folder has a stated job | High | The directory clause says what each folder holds and what it does not, `docs/design/` included, and where a reference file belongs: under `docs/references/`, with `DESIGN.md` at the root as the one exception the Baseline itself prescribes. Nothing is renamed and nothing is moved. |
| F10 | The Doctor Command reports what departs | High | New checks, each a warning with its next action and none a failure: a standard block altered or missing, a process folder exported, a variation without its reason, an ADR without the mark, a Lesson Learned off the form, a superseded ADR in the tree, a required file missing (`README.md` and the glossary), and a reference file or a folder outside the stated layout. |

## What the three records look like

Each block is an illustration with invented content. The PRD settles the names and the TechSpec the syntax.

The export list of a product repository. The sync reads one path per line, a `!pattern` line as an exclusion and `#` as a comment, and that does not change.

```text
# Standard list, kept by the Baseline: roundfix baseline update rewrites this block.
README.md
CONTEXT.md
docs/user-guide/
docs/references/
docs/design/
docs/lessons/
docs/adr/0007-a-credit-note-never-reopens-a-closed-period.md
docs/adr/0012-a-price-change-takes-effect-at-the-next-opening.md

# This repository's own lines: the Baseline never changes them.
# The design contract the interface follows.
DESIGN.md
# The public API contract.
openapi.json
# Dropped: copies of a supplier's manuals, under the supplier's licence.
!docs/references/vendor/
```

The header of an ADR, with the two fields this Idea adds to the lifecycle fields the Baseline already defines:

```markdown
---
status: accepted
knowledge: true # true: a rule of the business or of the product; false: about the tool, the code or the process
amended_by: [ADR-NNNN, ADR-MMMM] # the ADRs that change part of this one
created_at: 2026-03-02T00:00:00Z
updated_at: 2026-11-18T00:00:00Z
deprecated_at: null
superseded_by: null
---
```

A Lesson Learned, in a file under `docs/lessons/` named by its date and its slug:

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
- Applies when: a source sends a time without a zone. Does not apply when the source sends an offset.
- Links: the Finding at `<repository>@<commit>:docs/findings/2026-11-10-sales-move-to-the-next-day.md`, and Pull Request 412.
```

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Process files in the mirrors | 0 files under `docs/agents/`, `docs/history/`, `docs/backlog/`, `docs/specs/`, `docs/findings/` and `docs/handoffs/` in the mirror of every repository that applied the update, from 1,509 on 2026-11-18 | Repeat section 1 of [the measurement](references/2026-11-18-mirrored-documentation-measurement.md) over the mirrors, 30 days after the release. |
| This repository's mirror | Under 60 files, from 851 | Count the files of this repository's mirror after the first push to `main` that follows the delivery. |
| ADRs without the mark | 0 of 364 in this repository at delivery, and 0 in each repository once its marking change merges | The Doctor Command's check in each repository. Here the QA gate runs it. |
| Expired decisions that read as current | 0 superseded ADRs in a tree, from 4 in two projects; every ADR changed in part lists its amending ADRs in the header, from 1 of 64 that says so in its first twelve lines | The Doctor Command's check, and section 2 of the measurement repeated in this repository. |
| Lessons kept after review | At least 80% of the Lessons Learned written in the first 60 days are still in the tree or were taken in by a rule, and between 1 and 10 are written for each 100 archived Specs and closed Findings | The Git history of `docs/lessons/` in this repository: the files added, and the files deleted with the reason their commit gives. The one earlier count is 1.8 for each 100, and above 10 the test is too loose. |

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Strong — the knowledge base stops answering with process, and a failure one project paid for becomes findable by the others. It adds nothing to how Roundfix runs a Spec. |
| Reach           | What % of users would this affect?                  | Must do — every repository on the Baseline: 13 mirrors hold its guides, and 1,106 ADRs wait for a mark. |
| Frequency       | How often would users encounter this value?         | Strong — the list acts on every push to `main` and the mark on every new ADR. A lesson is rare by design: 7 in 383 files the one time it was counted. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe — each part exists elsewhere: curated lists for machine readers, status in the header of an ADR, post-incident templates, linters of repository structure. The web sweep found no tool that holds them as one rule set delivered to many repositories. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong — the text can be copied. What accumulates stays in each repository and in the knowledge base: marked decisions and lessons with their limits. |
| Feasibility     | Can we actually build this?                         | Strong — the readers of an ADR's header ignore fields they do not know, and the pair "the Doctor Command warns, the update repairs" exists for `CLAUDE.md`. Two things are new machinery: a block the Baseline keeps inside a file that has no markers and that the repository also edits, and a judgment on about 1,100 existing ADRs that people must review. |

The leverage is compounding, and the verdict is to proceed. Three results would change it. If most lessons are deleted in review, the test is too loose and V1 needs a stricter one. If the marking changes wait unreviewed, a repository's ADRs stay out of the knowledge base, and the mark costs more than it gives. If the Indexus finds its answers no better with process gone from the mirrors, the list was not the problem.

## Trade-off Insights

- **Recommended approach:** the six answers the maintainer gave, delivered as one Spec. The Baseline text carries the rules to every repository, the block in the export file makes the rule visible where the project works, and the Doctor Command makes a departure visible without stopping anyone.
- **Key trade-offs:**
  - A list at the source against a filter at the reader. At the source there is one list, in the repository, that a check can read. The cost is that the mirror itself shrinks: a path that leaves the list leaves the mirror at the next push, and nobody outside the repository can search that project's Specs or backlog any more. The filter the Secondbrain runs today keeps the mirror whole and stays invisible to the projects.
  - A mark on each ADR against exporting every accepted ADR. The mark is exact, and it is the maintainer's choice. It costs a judgment on each of about 1,100 ADRs and on every new one, and an ADR nobody marked disappears from the knowledge base without the reader knowing. Exporting them all costs nothing and sends decisions about tooling that the knowledge base does not want.
  - A lesson in a file against a lesson in a rule. A file can be searched by the next project that meets the same wall. Milton's reading of the NASA audit is that a lesson belongs in the practice it changes and that a store of lessons decays. V1 does both: the file exists until a rule takes the lesson in, and then it leaves.
  - The Agent writes against a person approves. Writing in the same Pull Request waits for nobody and still shows the lesson to a person. It also lets a weak lesson through whenever the review is hurried, and the published measurements say weak lessons hurt.
  - Deleting a superseded ADR against keeping it with its status. Common ADR practice keeps the old record in the log, and ADR-0278 says so itself. This repository deletes, the successor names what it replaced, and Git keeps the text. The maintainer kept that rule.
  - Stating each folder's job against renaming. A rename of `docs/agents/` touches 18 Go files, 104 test files, 18 Baseline files and 17 skill files here, and every adopter. Stating the job costs a paragraph, and the name goes on misleading.
- **Risks identified:**
  - The first update of a repository that exports `docs/` whole stops exporting every folder the standard does not name, such as a `docs/reference/` or a `docs/prompts/`. The Change Plan must name each path that stops being exported, or knowledge leaves the mirror unnoticed.
  - Pages of the Secondbrain cite mirror paths under `docs/history/` and `docs/findings/`. They lose their source when a project applies the update, unless the Secondbrain pins those citations first.
  - The mark of an existing ADR is a judgment made in bulk. A wrong "about the tool" hides a business rule from the knowledge base, and the review of a change with hundreds of marks is the only guard.
  - The standard block lives in a file with no markers that the repository also edits. The only precedent, the `CLAUDE.md` link, takes the whole file. A wrong merge of the two blocks changes what a project publishes.
  - The Spec changes at least three Baseline modules and two owned skills, and adds checks and a kind of managed file. It is the largest of the cycle.
- **Dissenting view:** the author recommended exporting every accepted ADR, because a decision absent from the knowledge base fails silently and a decision present costs a few hundred small files. The maintainer chose the mark, and this Idea follows him. A second dissent is Milton's: the right number of stored lessons is zero, and a repository that writes lesson files will stop embedding them. A third is the ADR practice this repository left: a reader who meets a citation of a deleted ADR finds a Git Locator instead of a file.
- **Stretch goal (V2+):** the authoring skills read `docs/lessons/` before they choose a library or a design, so a lesson reaches the session that is about to repeat the mistake. The Doctor Command offers the move of a reference file found at the root. An index generated from the headers lists each decision with its status, its mark and the ADRs that amend it. The rename of `docs/agents/` and the merge of backlog and findings stay in the backlog.

## Opportunity Scan

The original idea is bounded by how much the knowledge base is used. If it works fully, a mirror holds only what a reader outside the project can use, a failure is paid once across the projects, and no project keeps documentation rules of its own.

- **More ambitious: a manifest of knowledge.** The Baseline generates in each repository one index of what it publishes: each ADR with its status and mark, each lesson with its "applies when" line, each reference with one line. The reader selects by what a record says instead of by its path. Scale: medium to large. Impact Strong, Reach Must do, Frequency Strong, Differentiation Strong, Defensibility Strong, Feasibility Maybe: the Indexus would have to read the manifest, which is a contract between two projects that neither has today.
- **Simpler: the list alone.** The standard block and its check, with no Lesson Learned and no ADR fields. Scale: small. Impact Maybe, Reach Must do, Frequency Strong, Differentiation Pass, Defensibility Pass, Feasibility Must do. It moves into each repository what the Secondbrain's filter already does, and leaves lessons lost and expired decisions unmarked.
- **Adjacent: lessons where the author looks.** Keep the lessons for the repository's own sessions: the TechSpec and the Task read them before choosing. Scale: small to medium. Impact Strong, Reach Strong, Frequency Maybe, Differentiation Strong, Defensibility Strong, Feasibility Strong. It serves the project that learned and not the others, and the measurements say memory injected without a filter hurts.

**Chosen direction:** original — the maintainer set it in the request and in his six answers. It takes one thing from the simpler idea, that nothing is renamed or moved, and leaves the adjacent idea as the stretch goal.

## Integration with Existing Features

| Integration point | What it gains or how it is used |
| --- | --- |
| The docs-layout guide (`docs/agents/docs-layout.md`) | The directory clause states each folder's job and names `docs/lessons/`. The ADR header gains the two fields. The guide gains the form of the Lesson Learned and its test, and the retirement clause names the lesson a rule took in. The clause that says "Do not rewrite existing ADRs solely to adopt lifecycle metadata" (`:23`) gains the exception the maintainer decided: the mark and the list of amending ADRs are written into existing ADRs. |
| The cross-project clause (`docs/agents/agent-instructions.md:86`) | Its sentence about a Finding that closes writes the Lesson Learned instead of sending a message. |
| The Secondbrain guide (`docs/agents/secondbrain.md`) | It states the standard list and the declared variation. It applies only where the repository's Baseline has the Secondbrain turned on (`internal/baseline/assets/decisions.json:336`). |
| `roundfix baseline update` | Its Change Plan shows the block, each ADR header it fills and each path that stops being exported, behind the one confirmation it already asks for. The request says `roundfix update`; that command only has `check` (`internal/cli/update_check.go:86`), and it previews this same plan. |
| The `--update` Skill Mode of the Roundfix skill | It already asks before a plan deletes files (`.agents/skills/roundfix/references/modes.md:59`). It is where a session marks the existing ADRs of a repository after the plan is applied, if the first open question keeps its default. |
| The Doctor Command | It gains the checks of F10 and still changes nothing (`CONTEXT.md:1201`). Each check names the command or the edit that repairs. |
| The history cleanup inside the update | Unchanged. It already deletes a superseded ADR and writes the line that names it into its successor (`internal/spec/history_entries.go:125`). |
| `archive-spec` and `roundfix archive` | Before the archive, a Spec that needed a Correction, a failed QA or a retry gets the lesson question. Archive Advice and the promotion of files stay as they are. |
| The Corrections record | A declared Correction is one of the signs that a Spec stumbled. The Archive Record keeps no such count today (`internal/spec/archive_record.go:43`). |
| The Spec Consistency Check | Unchanged in V1. It reads an ADR's header and ignores fields it does not know. |
| Vendored Skills | `domain-modeling` is not edited: the lifecycle header is the Baseline's own overlay on its ADR format. `knowledge-workspace` describes the export file and tells a session to add a new artifact to it, which the declared variation keeps true. |
| This repository's own files | `.secondbrain-export` becomes the standard block with this repository's variation, its 364 ADRs gain the mark, and those changed in part gain the list. The test of the export file (`internal/docscontract/secondbrain_export_test.go`) follows the new rule. |
| `CONTEXT.md` | The names the PRD settles: Lesson Learned, the standard list and the declared variation, and the two ADR fields. |
| The Secondbrain and the Indexus, which are other projects | The Secondbrain keeps the template the projects copy, the sync and the filter of 2026-11-18. This Spec changes none of them and tells both projects what it delivers. |

## Size and Shared Files

By the sizing criterion of the Spec routing guide this is one Spec. It holds one context: what a repository's documentation holds as knowledge, and what it publishes. The parts complete one another. The list names the folder of lessons and the marked ADRs, so it is not whole without them. A lesson has no reader outside its repository without the list. The Doctor Command checks all of them in one place, and three of the four gaps are closed by text in the same guide.

The work fits four or five implementation Tasks and the QA gate, each large enough to earn its verification cycle:

1. the rules in the Baseline source, with the files each edit regenerates, and the names in `CONTEXT.md`;
2. `roundfix baseline update` keeps the standard block and fills the list of amending ADRs;
3. the Doctor Command's checks, with the user guide and the record of command surfaces;
4. the steps in the owned skills `archive-spec` and `roundfix`;
5. this repository applies the rules to itself: its export file, the mark on its 364 ADRs, and the list on those changed in part.

The reference of four Tasks is a signal to look for a second context, and the look finds one context. If the PRD or the TechSpec finds the Spec too large to deliver safely, the clean cut is the Lesson Learned (F5 to F7): it shares one line of the list with the rest and could follow as a second Spec. This Idea does not propose that cut, because the request joined the two and the maintainer grouped them.

This Spec is delivered last in its cycle, after Specs 0407, 0408 and 0409, in the order the maintainer set. Nothing in it depends on their code. The four are authored in parallel, and these files can meet:

- The Baseline source under `internal/baseline/assets/`, `internal/baseline/module-versions.json`, the digest pin, the catalog snapshots, the plan goldens and `docs/agents/setup-context.json`. Spec 0409 adds a clause to the same modules, and this Spec changes at least `context-workflow`, `core` and `secondbrain`.
- `docs/agents/docs-layout.md`. Spec 0409 may describe there the dated line a kept rule gains, and it left the ADR statuses to this Spec. This Spec adds two header fields and no status.
- `docs/agents/agent-instructions.md`, where Spec 0409's clause may land and where this Spec changes one sentence.
- The owned Roundfix skill: its version fields and `references/modes.md`. Specs 0407, 0408 and 0409 also change this skill.
- `skills/testdata/owned-skill-versions.json` and the `skills/` mirror, for `archive-spec` and `roundfix`.
- `docs/references/behavior-surfaces.json`, because the Doctor Command and the update gain output. Spec 0409 leaves it alone; Specs 0407 and 0408 may change it.
- `CONTEXT.md`.

## Out of Scope (V1)

- **Renaming `docs/agents/` or merging backlog and findings** — the maintainer chose to state each folder's job and rename nothing. The measured cost is in a new Backlog Entry, `docs/backlog/2026-11-18-the-docs-agents-folder-is-misnamed-and-backlog-and-findings-stay-in-two-folders.md`.
- **Moving a file by command** — a `DATABASE.md` at the root or a folder outside the layout is reported, and the repository moves it. A command that moves a project's own files can break links the command cannot see.
- **Rules for what a README or a user guide must say** — the request names "the content of documentation files". V1 covers the records that have a form (ADR, Backlog Entry, Finding, Lesson Learned) and the files that must exist. What a README says is a subject of its own.
- **A check that fails a Run, a merge or the Spec Consistency Check** — every new check is a warning of the Doctor Command. A repository that has not yet applied the update must keep working.
- **Reading lessons during authoring** — it is the stretch goal, and it needs lessons to exist and a way to pick the ones that apply.
- **A judge that filters lessons** — the maintainer chose that the Agent writes and a person sees the Pull Request.
- **Keeping a superseded ADR with a status** — the maintainer kept the rule that deletes it.
- **Changing the Secondbrain's template, sync or filter, and applying the update in another repository** — they belong to other projects. Each repository applies its own update, and Roundfix never edits another project.
- **Turning past Findings and Archive Records into lessons** — the seven lessons of October are already in the Secondbrain, and the record starts with what closes after the release.
- **A place for `CHANGELOG.md` in the standard list** — a repository that wants it published keeps it in its own lines.

## Idea-stage judgment

- item: ai-memory message `01a1eeca-397e-77c2-8429-3668a5a0b8f8`, kept in [`references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md`](references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: unsure (0.68)
- unverified_assumption: no (0.17)
- market_question: no (0.06)
- verdict: undecided
- stage: ran
- shape: changed — the export list became a standard block with a declared variation, each ADR gained a mark that decides its export, a superseded ADR is still deleted while an amended one says so in its header, the Lesson Learned gained a test, two moments and an exit, the folders keep their names, and the Doctor Command reports while `roundfix baseline update` applies.

The verdict was undecided, and the maintainer chose "Spec com etapa de Ideia (Recommended)" on 2026-11-18. The `shape` line is the author's and waits for the approval of this draft.

## Decisions

Every answer is the maintainer's, as recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md). The text in quotes is the option he chose, verbatim. The text after the dash is that option's description as the operator recorded it.

- The Idea stage (2026-11-18, ~19:55 BRT): "Spec com etapa de Ideia (Recommended)". The folders entry (same answers): "Juntar à Spec do Indexus (Recommended)". Authorization: "Mesmas do BE (Recommended)" — push, PR, merge, release após a Spec; até 5 correções; Makefile/CI/go.mod/lint só se a correção exigir. The PRD quotes it in the authorization record. Delivery order: "Corretiva → regras → Indexus (Recommended)".
- When to author it (2026-11-18, ~22:35 BRT): "Ideia 0410 agora (Recommended)".
- The folders (2026-11-18, ~23:20 BRT): "Só a função de cada pasta (Recommended)" — o Baseline diz o que cada pasta guarda (inclusive docs/design/ e onde vão arquivos como DATABASE.md); o doctor aponta o que está fora; nada é renomeado nem movido; renomear/juntar volta ao backlog com o custo medido.
- The variation of the list (2026-11-18, ~23:20 BRT): "Padrão + variação declarada (Recommended)" — o Baseline escreve o bloco padrão e o baseline update o mantém; o repositório pode tirar itens ou acrescentar arquivos de conhecimento num bloco próprio com o motivo em uma linha; o doctor aponta pasta de processo exportada e bloco padrão alterado.
- The ADRs in the list (2026-11-18, ~23:20 BRT): "Só as marcadas" — cada ADR ganha um campo que diz se interessa à base de conhecimento. It was not the author's recommendation.
- Who marks (2026-11-18, ~23:20 BRT): "Agente marca, você vê no PR (Recommended)" — ao escrever uma ADR o autor a marca (conhecimento/regra de negócio ou ferramenta); nas existentes, o baseline update gera a marcação num PR por repositório, que o Marcio pode corrigir; ADR sem marca não é exportada.
- The ADR that no longer holds (2026-11-18, ~23:20 BRT): "Apagar substituída, emendada no topo (Recommended)" — a regra de apagar fica; doctor aponta a superseded que sobrou e o baseline update a remove; a ADR emendada ganha no cabeçalho a lista das ADRs que a emendam.
- The moment of the Lesson Learned (2026-11-18, ~23:20 BRT): "Finding fechado e archive c/ tropeço (Recommended)" — Finding que fecha e Spec arquivada que teve correção, QA reprovado ou retentativa; sem lição reutilizável, nada é gravado.
- Who decides the lesson (2026-11-18, ~23:20 BRT): "Agente grava, você vê no PR (Recommended)" — o Agente aplica o teste do Baseline (tentativa, problema, causa provada, solução, ainda não dito em ADR/referência/glossário) e grava no mesmo PR; quem revisa pode apagar.

The author's own, made in this stage:

- The questions did not offer the dismissal of the stage again, because the maintainer had chosen it twice. The operator asked them, because the Author Subagent cannot ask.
- Both Backlog Entries are adopted here: each has the status `promoted`, names this Spec and moved into `references/`. The rename of `docs/agents/` and the merge of backlog and findings went into a new Backlog Entry with the measured cost, as the maintainer's answer on the folders says.
- The standard list is the one the request suggested, with one addition: `docs/design/`, so that no design artifact published today is lost. `docs/handoffs/` and `CLAUDE.md` join the items never exported.
- The place of a Lesson Learned is `docs/lessons/`, a folder of its own, as the request asks. Its form holds the seven things the request lists.
- A Lesson Learned leaves the tree when a rule takes it in, under the retirement clause. It is the maintainer's direction of 2026-10-14 for everything retired, and it answers the published finding that a store of lessons decays.
- The file replaces the message a closing Finding sends today. Sending both would file the same lesson twice.
- The exception to the clause that forbids rewriting existing ADRs for lifecycle metadata (`docs/agents/docs-layout.md:23`) follows from the maintainer's answers on who marks and on the amended ADR. The clause changes with this Spec.
- Every new check of the Doctor Command is a warning. The checks it has for the Baseline today warn and never fail.
- The files that must exist are `README.md` and the glossary. The request names required files and no list of them.
- This repository publishes `README.md`, `docs/user-guide/` and `docs/lessons/`. The Secondbrain's filter of 2026-11-18 already sends only the first two, by the maintainer's decision recorded there. The lessons are the author's addition and the second open question.
- No ADR is written at this stage. The ADR numbers 0438 to 0440 are reserved for the PRD and the TechSpec.

## Open Questions

The first three are defaults the maintainer should see when he approves this draft.

- Who proposes the mark of each existing ADR? His answer says the update generates the marking, and a command cannot tell a business rule from a tooling decision by itself. Default: the Agent of the session that applies the update marks each ADR by the same test an author uses, in one change. The alternatives are the judge inside the command, which sends each ADR's text to an outside service and needs that authorization in every repository, and no proposal at all, where every existing ADR stays out of the knowledge base until someone marks it.
- Does this repository publish its Lessons Learned? Default: yes, with `README.md` and `docs/user-guide/`. The request describes lessons of exactly this kind, and the seven rescued in October came from here. The alternative is only `README.md` and `docs/user-guide/`, as the Secondbrain's filter does today.
- Is `docs/design/` in the standard list? Default: yes. The mirrors hold 15 such files, the filter of 2026-11-18 lets them through, and leaving the folder out would drop them at each project's first update. The alternative is that a repository adds it in its own lines.
- How do the marked ADRs become paths in the list? The sync belongs to the Secondbrain and reads paths. Default: the standard block names each marked ADR, the author adds the line with the ADR, the update rewrites the block and the Doctor Command reports a difference. The TechSpec decides.
- What of a Correction declared after the archive commit? The queue archives before it pushes (`.agents/skills/roundfix/references/deliver.md:43`), and a Correction can come later. Default: the session that declares it applies the same test and writes the lesson in the change that corrects, in the same Pull Request.
- How does the archiving Agent know the Spec stumbled? Default: from the QA Reports in the Spec folder, the Corrections declared for the Spec and the retries of its Run, as the PRD defines, with no new field in the Archive Record unless the TechSpec needs one.
- What does the first update do with a folder under `docs/` that the layout does not name? Default: the Change Plan names every path that stops being exported, the Doctor Command reports the folder, and the repository adds it in its own lines or moves its content.
- What does the block list in a repository with several domain contexts, or without a `docs/user-guide/`? Default: the glossary files the repository's declared layout has, and only the standard paths that exist, with the update adding a line when its folder appears. The TechSpec confirms how the sync treats a path that does not exist.
- Does an ADR without any header gain one for the mark? The mirrors hold 19 such ADRs. Default: yes, with the mark alone; the lifecycle fields stay optional for it.
- In which language is a Lesson Learned written? Default: the language the repository writes its documentation in, with the field labels as the Baseline gives them.
- Which names does the glossary take? Default: the PRD's Glossary settles them, starting from Lesson Learned.
- What does the Notice to the Indexus carry, and who tells the Secondbrain? The Indexus asked for the final list, the form and path of the lesson, the marking of ADRs and the projects already updated. Default: one completion notice when the release is published, as the Notice Request rule says (`CONTEXT.md:185`), naming this repository as the one updated, because each project applies its own update. The approved Idea can serve as the partial plan the Indexus asked for, and the Secondbrain should hear, before any project applies the update, that mirrors will shrink and that pages citing mirror paths need their sources pinned. The author sent nothing.
- On what outside evidence does the acceptance rest? Default: the measurement adopted here and the published passages quoted in `### Market Data`.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md`](references/2026-11-18-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export.md), the Backlog Entry that holds the two messages of the Indexus and its Notice Request.
2. [`references/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md`](references/2026-11-16-the-docs-folders-of-adopters-have-unclear-jobs.md), the Backlog Entry with the maintainer's list about the folders.
3. [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md), the answers behind `## Decisions`.
4. [`references/2026-11-18-mirrored-documentation-measurement.md`](references/2026-11-18-mirrored-documentation-measurement.md), recorded by the author: the counts used in `## Problem`, `## KPIs` and `## Trade-off Insights`.

The other open Backlog Entries belong to other contexts and stay where they are. The Idea of Spec 0409 was read for the files the two Specs share.

This project's ai-memory messages were listed on 2026-11-18 and none was pending. The Research Subagent did the two codebase sweeps and the web sweep. The author did the Secondbrain sweep, because it needs a command, counted the mirrors, and re-read every passage quoted in this file.

Secondbrain (read on 2026-11-18; `wiki/index.md`, then four `qmd query` runs, which answered):

- `wiki/concepts/politica-de-retirada-e-poda.md` — "Uma lição de finding pertence à página do conceito, não a uma página por finding", and capture at closing "é condicional: só quando a lição não está na ADR ou no record absorvente". It set the test of F6 and the exit of F5, and it records the open tension between an immutable ADR log and deleting superseded ADRs.
- `wiki/sources/roundfix-sete-licoes-history-2026-10-15.md` — the seven lessons rescued from this repository's history, each one a symptom, a cause and one sentence of lesson. It showed that the form the request asks for fits real cases, and that three of the seven would have been lost with the Specs that absorbed them.
- `wiki/entities/indexus.md` — what the Indexus is, that it reads the Secondbrain's mirrors, and that the page itself cites mirror paths under `docs/history/` and `docs/findings/`. It is the source of the risk about citations.
- `brain/RUNBOOK.md`, the section on sending to the Indexus — the filter the maintainer decided on 2026-11-18, which removes from what is sent nearly the same set this Idea never exports and reduces this repository to its README and user guide. It adds: "O espelho no repositório não muda." It confirmed the list and showed what changes: this Spec shrinks the mirror itself.
- `templates/repo/secondbrain-export.example` and `scripts/sync-project.sh` — the template every project copies, and the sync that replaces a mirror with what the list names. They showed where the list lives today and what happens to a path that leaves it.
- `projects/*/mirror/` — counted, not read, for the measurement.

External sources (found through Exa by the Research Subagent on 2026-11-18; the author re-read every passage quoted here):

- Nick Milton, "Leaving lessons in a lessons database doesn't work - an example from NASA", http://www.nickmilton.com/2017/07/leaving-lessons-in-lessons-database.html — the audit's finding and the argument that a lesson belongs in the practice it changes. It gave the exit of F5 and stands as a dissent against storing lessons at all.
- Lindenbauer, Groh and Schütze, "From Knowledge to Noise: CTIM-Rover and the Pitfalls of Episodic Memory in Software Engineering Agents", https://aclanthology.org/2025.realm-1.30.pdf — distilled lessons made a coding agent worse in every configuration, by noise. It supports a test that writes nothing by default and challenges the stretch goal.
- Mehmet Iscan, "Learning When to Remember: Risk-Sensitive Contextual Bandits for Abstention-Aware Memory Retrieval in LLM-Based Coding Agents", https://arxiv.org/abs/2604.27283 — a stored fix misleads when the new failure only looks like the old one. It supports the fields that say when a lesson applies and when it does not.
- "How to deal with conflicting ADRs? AI picks what suits it best", https://github.com/mattpocock/skills/discussions/847 — an Agent reads the first lines of an ADR and misses a notice in the body. It supports the list of amending ADRs in the header.
- Jeremy Howard, "The /llms.txt file, v2", https://llmstxt.org/ — a curated list for machine readers as a convention. It supports a list of knowledge kept at the source.
- Repolinter, https://github.com/todogroup/repolinter — a ruleset that reports and repairs the structure of a repository, now archived. It is the pattern this Idea splits in two: the Doctor Command reports, and a plan a person confirms repairs.

The Research Subagent also read, and this file quotes none of them: NASA's Office of Inspector General report on its lessons system (2012), the example postmortem of the Google SRE book, two further papers on memory for coding agents, three texts on how ADR conventions and tools mark a superseded decision, Diátaxis, and two practitioner texts on stale documents in a retrieval corpus.

The digest of this research for the project that owns the Secondbrain was not sent by the author, who sends nothing in this cycle. The operator has it in the author's reports.
