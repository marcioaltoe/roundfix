---
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: feat
---

# A Correction records its Defect Class and a phase entry its authoring model

When a delivery is corrected, nothing records why. Roundfix knows that a review finding was fixed and never whether the plan or the implementation was at fault, and the Phase Record does not say which model wrote each part of a plan. On 2026-11-17 the maintainer asked whether Fable is effective for Ideas, PRDs and TechSpecs compared with Opus, and the answer took a study that rebuilt every cause by hand from transcripts, commits and notes and still left three of 23 delivered Specs with a question mark. After this Spec, whoever corrects a delivery declares the Defect Class of each cause with a one-line reason, a phase entry records the model and reasoning effort that ran its Authoring Phase, and one read-only report counts both: Corrections per Spec by class, and plan defects against the model that authored the artifact they name. The report shows every count beside its sample and withholds a comparison it cannot support. A repository requires the declaration only when its own Project Config says so.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md), the request, the Idea-stage judgment and the maintainer's answers, recorded by the operator.
2. [`references/2026-11-17-fable-and-opus-authoring-study.md`](references/2026-11-17-fable-and-opus-authoring-study.md), the operator's study of Specs 0370 to 0403.
3. [`references/2026-11-17-phase-record-and-causes-view-measurement.md`](references/2026-11-17-phase-record-and-causes-view-measurement.md), recorded by the author of the Idea: what the Phase Record and `roundfix runs causes` held, and the commits of one delivery.
4. [`references/2026-11-17-item-branch-commits-and-author-session-measurement.md`](references/2026-11-17-item-branch-commits-and-author-session-measurement.md), recorded by the author of this PRD: what marks each commit of an Item Branch, what a command run by an Author Subagent reads, and the Governed Paths the design is expected to change.
5. [`references/2026-11-17-techspec-design-measurements.md`](references/2026-11-17-techspec-design-measurements.md), recorded by the author of the TechSpec: a third reading of what an Author Subagent's command reads, the identity of a repository by its root commits, where the plan of a delivered Spec is readable, the Governed Paths the design changes, the existing assertions it moves and the sources read outside the repository.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from. This project's ai-memory inbox was listed on 2026-11-17 and held no pending message, and `docs/findings/` holds no Finding. None of the four open Backlog Entries shares this Spec's context: a Conexus report that waits for reproduction steps, the start cost of a session, which waits for a week of per-phase numbers and is about tokens, the Delivery Queue under strict branch rules, and the documentation folders of adopters.

## Prerequisites

Specs 0401, 0402 and 0403 are planned on the default branch and are delivered before this one, as the maintainer ordered ("autoria e entrega depois da 0403"); the Tasks stage names them in `requires`. Authored against `74d0e3a9` and rebased on `a8236a87`, the delivery of Spec 0400; none of the three was delivered then.

- Spec 0401 changes the writer of the Phase Record, the phase command, the phase section of the four authoring skills and the Roundfix skill's phase reference, which this Spec also changes.
- Specs 0402 and 0403 change the Roundfix skill's version fields and references, the glossary, the surface record and the configuration guide, which this Spec also changes, each inside its own sentences. They change this repository's Project Config; no Task of this Spec does.
- Spec 0402 decides what the comparison can show: it declares Fable for the Idea, the PRD and the TechSpec and Opus for the Tasks, so within one Delivery Type a phase may have a single model.

No version number is written by hand: the record commands raise it on the tree they run on.

## Project Constraints

- Identifier strategy: applicable — a Correction is named by a short identifier that is unique within its Spec and that a later declaration cites; it is public CLI contract, and the TechSpec fixes its form. No other identifier is minted: a phase entry and a Correction Sign carry identifiers other components recorded. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the report starts no network request, and no judge is asked. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0417 (new) records the Correction, its declared class and the setting (ADR-0417: "one Correction per cause, with a Defect Class from a closed list of five and a one-line reason"). ADR-0418 (new) records the model and effort of a phase entry (ADR-0418: "The session reports the model, because nothing a command can read names it") and amends ADR-0407, whose entry held neither (ADR-0407: "It never holds prompt text, a token count, a price or an environment value"). ADR-0419 (new) records how the report shows its sample (ADR-0419: "prints no percentage"). Cause Class keeps its rule (ADR-0214: "never by a model judgment"), and no report merges it with a Defect Class. The phase command keeps what Spec 0401 gives it (ADR-0414: "records each phase as before and reads no earlier entry"), the status view keeps reading the same entries (ADR-0410: "it reads phase entries only, and it reports no usage"), and the Authoring Model Declaration is unchanged and is not the source of the recorded model (ADR-0401: "Roundfix ships no model name"). A new Project Config key cannot reach this repository in the change that teaches the binary that key (ADR-0408: "refuses a key it does not know"), and an older binary refuses a repository a newer one wrote (ADR-0392: "every command that loads the repository's configuration"). The record follows the Run Retention window and leaves the queue's own tables alone (ADR-0255: "the Delivery Queue tables are never touched"). A review finding's disposition keeps its record and its refusals (ADR-0272: "The record is keyed by the resolved checkout"; ADR-0279: "The disposition stays tied to the reviewed head"). The new park is the item's own work (ADR-0395: "A queue item stops only for its own work") and has its way back (ADR-0398: "A parked item always has a way back through a Delivery Retry"). Each changed surface ships with its skill text (ADR-0256: "The Roundfix-owned skills must describe the behavior Roundfix ships"). The Release Ceiling in the authorization record answers the release plan (ADR-0349: "the Spec's authorization record carries the highest impact its release grant answers"). ADR-0415 sets the models this repository declares for each phase, which decide what the comparison can show, and is unchanged. The TechSpec added seven decisions and took up two older ones. The record is a file per Spec (ADR-0420: "one JSON Lines file per Spec under Roundfix Home"), keyed by an identity that needs no remote (ADR-0421: "computed from Git objects alone"), which is the identity a Baseline Plan already carries (ADR-0071: "the stable repository identity, plan digest, and all bounded preimages match"). A correction by hand is told by what its commit carries (ADR-0422: "never by who made it"). The setting is `delivery.require_defect_class`, and its park is a `finding` (ADR-0423: "read from the default branch and never from the item's tree"), as the Delivery Queue already reads one delivery setting (ADR-0288: "so enabling it takes effect after the change that enables it merges"). The commands are `roundfix corrections declare` and `roundfix corrections report` (ADR-0424: "a command group of its own"). The session passes its model as a flag (ADR-0425: "only the exact identifier its harness states"). The report takes a Spec's size from the default branch (ADR-0426: "so it never fetches"). Not applicable: ADR-0294, ADR-0297 and ADR-0303, which cite ADR-0288 for the Linux Lane, which this Spec leaves alone; ADR-0173, which cites ADR-0071 for the citations a Baseline Plan reports; ADR-0286, whose Correction Review is a review round with a near name and is unchanged, with ADR-0197, ADR-0366 and ADR-0391, which cite it; ADR-0413, which the study names in a note and which ADR-0414 superseded before any Run implemented it; ADR-0353, which cites ADR-0410 and holds Waiting Items back from the follow-up list; ADR-0361, which cites ADR-0214 for the secret-scan item of `roundfix runs causes`; ADR-0373, which cites ADR-0279 for the liveness of a Queue Owner; ADR-0405, which cites ADR-0401 for a user-level skill copy; ADR-0412, which cites ADR-0349 for the release impact of planning references; and ADR-0290, ADR-0296 and ADR-0301, which cite ADR-0255 for storage and collection: the TechSpec keeps the record outside every Artifact Root and outside the Run Database, where neither the storage report nor GC sanitation reads. The TechSpec relates each one, and its not-applicable lines close the account of the records that cite only those. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-17, quoted in the authorization record, which keeps the grant of the owned skills. It grants no Makefile, workflow, Go module, lint or formatter change. The one-key change of this repository's Project Config that turns the setting on is the operator's step after the release, under the maintainer's answer "Ligar e manter unclear (Recommended)" of 2026-11-17, and no Task makes it. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/corrections.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/roundfix/references/review.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `.agents/skills/roundfix/references/storage.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- Every correction of a delivery has a recorded cause: a class from a closed list, a one-line reason, and the artifact at fault when the plan was.
- Every Authoring Phase of a Spec has a recorded model and reasoning effort, or says that they are unknown.
- The person who chooses the authoring models answers "plan defects per authoring model" from one report, with no transcript, commit log or memory note.
- The report never says more than its sample supports: counts beside their sample, a comparison only above a stated minimum, and no verdict on a model.
- A repository that does not turn the requirement on delivers exactly as it does today.

## User Stories

1. As the maintainer, I want each correction of a delivery to carry a class and a reason, so that I stop rebuilding the cause from commits and notes.
2. As the maintainer, I want a plan defect to name the PRD, the TechSpec or the Tasks, so that it counts against the model that wrote that artifact.
3. As a Supervisor that fixes a review finding, I want to declare the cause in the command that records the fix, so that declaring costs no extra step.
4. As a Supervisor that repairs a QA finding or any other defect, I want one command that declares the cause, so that a correction without a review finding is recorded too.
5. As the maintainer, I want to declare a Correction again when I disagree, so that my reading counts and the earlier one stays visible.
6. As the maintainer of a repository that turned the requirement on, I want a delivery to stop before its merge while a correction has no class, so that the record is complete.
7. As an adopter that did not turn it on, I want my deliveries, commands and outputs unchanged, so that an update costs me nothing.
8. As the maintainer, I want each phase entry to record the model and effort that ran the phase, so that I know which model wrote which artifact.
9. As the maintainer, I want one report of Corrections per Spec and of plan defects per phase and model, so that I answer the question of 2026-11-17 without a study.
10. As a reader of that report, I want every rate as a count over a count beside the Delivery Type and size, and no comparison below a minimum, so that I do not read a difference four Specs cannot show.
11. As a reader of that report, I want dismissed findings and failed Verification attempts shown beside the Corrections and not counted with them, so that I see the whole delivery without mixing three different facts.
12. As an Agent that reads the report, I want the reason of a Correction marked as data, so that I never follow it as an instruction.

## Core Features

1. **A Correction is one cause.** A **Correction** is one cause that needed a repair while a Spec was being delivered, from the moment its plan is on the default branch until its delivery merge: a pre-PR review finding that was fixed, a QA finding that was repaired, or another defect that needed a change to the delivery candidate or to the plan. A commit that repairs two causes holds two Corrections, and two commits for one cause are one. A dismissed finding, a failed Verification attempt the Agent repaired inside its Run, a stop that needed no repair and a plan change that records a new decision of the maintainer are not Corrections.
2. **A Correction carries a class and a reason.** A Correction carries one **Defect Class** from a closed list and a reason. `plan`: the plan said something wrong, missing or impossible. `implementation`: the plan was right and the delivered change was not. `environment`: the cause lies outside the Spec, such as the machine, a tool or a failure already on the default branch. `operator`: an action of the Supervisor or of a person during the delivery caused it. `unclear`: the declarer cannot tell, and the reason says what is known. The reason is one line of at most 200 characters. Hazard: the reason is text a declarer wrote, and an Agent reads it later in a report; it is recorded and printed as data on one line, a reason that is blank, longer, or holds a line break or a control character is refused, and the skill text says that a reason is never an instruction.
3. **A plan defect names its artifact.** A Correction of class `plan` names the PRD, the TechSpec or the Tasks, and a declaration of `plan` without one is refused; any other class names none. A defect that comes from the Idea names the PRD, which adopts it.
4. **Whoever corrects declares.** One command declares a Correction for a Spec: its class, its reason, its artifact when the class is `plan`, and the Correction Signs it covers. It gives the new Correction an identifier that is unique within the Spec and prints it. The record keeps when it was declared, the harness session that declared it when the environment names one, and the name the declarer gives when it declares on a person's word. A declaration is accepted for a Spec the repository holds, active or archived, and refused for a slug it does not hold. Hazard: the Supervisor corrects in the item's worktree and reads the report in the user's checkout, so a declaration made in either one is written to the one record of the repository. A refusal names the accepted values, writes nothing and exits as a usage error.
5. **A fixed review finding is declared with its disposition.** The command that records a finding as fixed also takes the class, the reason and the artifact, or the identifier of a Correction already declared for the same cause, and then records the **Review Finding Disposition** as today and the Correction with it. Without them it records the disposition exactly as today. A dismissal takes none: nothing was corrected.
6. **A later declaration replaces an earlier one.** Anyone may declare a Correction again by its identifier, the maintainer above all, or withdraw one declared by mistake with a reason. The last declaration counts, the earlier ones stay in the record, and the report shows that a Correction was declared more than once and whether its class changed. A withdrawn Correction is counted apart and covers no sign.
7. **Roundfix knows three signs of a Correction.** A **Correction Sign** is a recorded fact from which Roundfix knows, without a declaration, that a delivery of the Delivery Queue was corrected: a review finding with a fixed disposition; a commit on the Item Branch that neither a Run nor the Delivery Queue made and that is not a merge of the default branch; and a Corrective Task. A declaration names the signs its Correction covers, one sign may be named by several Corrections, and a sign that no standing Correction names counts as one `undeclared` Correction. Hazard: a new Delivery Queue replaces a finished one with its items, and an Item Branch is released after the merge, so Roundfix records the signs of a delivery no later than its merge, in a record the report still reads afterwards.
8. **The requirement is a Project Config setting, off by default.** A boolean setting of Project Config makes the declaration required in that repository. It is `false` when absent, and a value in User Config is ignored with a warning, as the Authoring Model Declaration is there. The configuration guide says what the setting does and that a Roundfix older than the release that carries it refuses a Project Config that holds it, so a repository turns it on after every machine that runs its Roundfix is upgraded.
9. **With the setting on, the Delivery Queue does not merge an undeclared Correction.** Before it merges an item, the Queue Owner reads the item's Correction Signs. While one is named by no standing Correction, the item parks with a blocker of its own instead of merging. Its Park line names each such sign (a finding by its identifier, a commit by its short hash and its subject cut at 72 characters, a Corrective Task by its identifier), the five classes and the command that declares, followed by the Delivery Retry. From the moment such a sign exists, `roundfix deliver status` shows it on a line for the item, so a Supervisor can declare before the item reaches its merge. After the declarations, a Delivery Retry returns the item to the stage its unchanged head supports, and it merges. A declaration made while the Queue Owner runs is read at the merge without restarting the owner. A merge made outside the Delivery Queue is never refused.
10. **With the setting off, nothing is required and nothing changes.** A delivery with corrections and no declaration merges as today. `roundfix deliver status`, the fixed disposition without a class and every other command print the bytes they print today, and no command asks for a class. A class that is given is recorded.
11. **A phase entry holds the model and the effort.** An entry of the `idea`, `prd`, `techspec` or `tasks` phase in the **Phase Record** records the model that ran the phase as the harness names it, which the session reports; the reasoning effort, read from the session's environment where the harness exposes it and otherwise as the session reports it; how each of the two was learned (`reported`, `environment` or `unknown`); and whether the phase's Author Subagent or the session itself ran the phase, as reported. The values are recorded when the environment names a Claude Code session; any other entry holds `unknown`. Hazard: a reported value is text from a session; a value that is missing, too long or malformed is recorded as `unknown` with one warning, the entry is still written, and the command's refusals and exit codes stay as they are. The result line names what was recorded and how it was learned. The keys are added to the record's present schema, and an entry written before them reads as `unknown`.
12. **A phase entry stays attributable to its repository.** An author works in a clone that is removed once the plan merges, so an entry of an Authoring Phase also records, when it is written, an identity of the repository that outlives the working tree. Hazard: a remote address can carry a credential; the entry never records one.
13. **The authoring skills report them.** The four authoring skills tell the session, and the Author Subagent, to report the exact model identifier the harness states for it when it records the phase, never an alias and never a guess, and to say when an Author Subagent runs the phase. When the installed Roundfix does not accept what the skill passes, the skill records the phase as before and says so in its report.
14. **The Corrections Report shows each delivered Spec.** The **Corrections Report** is a read-only view, in text and in JSON. For each Spec delivered in its window, in Spec number order, it shows the Delivery Type, the number of Tasks without the QA gate, the model and effort of each Authoring Phase, the Corrections by Defect Class with `undeclared` and withdrawn counted apart, and the plan defects by artifact. A Spec is delivered when the default branch holds its Archive Record and the record is not a cancellation; the window is every such Spec this machine's records still hold, and two flags narrow it by delivery date. Asked for one Spec, in delivery or delivered, the report lists each Correction with its identifier, class, artifact, reason, declarer, time and number of declarations, and each Correction Sign that no Correction names.
15. **The report counts plan defects by Authoring Phase, model and effort.** For each Authoring Phase, Delivery Type, model and effort, the report shows the Specs and the Tasks that group authored and the plan defects that name that phase's artifact, and beside them each model's total across its efforts. The model of a phase is that of the last entry of that phase recorded before the Spec's first Implement Run began, or before its delivery merge when this machine holds no Run of the Spec; when the entries before that moment disagree, the Spec's row says that the phase was recorded more than once. A defect of a PRD that the TechSpec phase wrote, as it writes the minimal PRD of a fix or a refactor, counts for the TechSpec phase's entry. The `idea` phase shows its Specs and Tasks and no defect count.
16. **The report shows its sample and withholds a comparison.** A rate is printed as a count over a count, for Specs and for Tasks, beside the Delivery Type; the report prints no percentage. For one Authoring Phase and one Delivery Type it sets side by side the groups, and the model totals, that each hold at least 20 delivered Specs; with fewer than two of them the line reads `insufficient sample`, with each group's count and the 20 needed. The report ranks nothing and prints no word that names a better or a worse model.
17. **The report says what it left out.** It states its window, the Roundfix Home it read, the Specs whose model is `unknown` for a phase, the Corrections that are `undeclared`, and the Specs it left out of a count with the reason: a Spec whose delivery this machine did not record shows what was declared for it and enters no count of signs.
18. **Shown beside, never counted.** For each Spec the report shows, beside its Corrections, the number of dismissed review findings and the failed Verification attempts by Cause Class. Neither enters a count of Corrections. A Corrective Task is a Correction Sign here and stays an item of `roundfix runs causes` there. Cause Class, the Classification Signatures and the text and JSON of `roundfix runs causes` do not change.
19. **Read-only, offline and bounded.** The report writes nothing under the repository or under Roundfix Home, starts no network request, asks no judge and reads no harness transcript. A source it cannot read is named as unavailable with its reason, the rows it still has are shown, and the report exits `0`; invalid usage and an unreadable record follow the exit codes of `roundfix runs causes`.
20. **The record is kept for the Run Retention window.** The declarations and the signs are kept under Roundfix Home for the Run Retention window, and the Run Retention Sweep removes what is older. Nothing is added to the Archive Record, and no text leaves the machine.
21. **The Supervisor is told when and how to declare.** The Roundfix skill tells a Supervisor to declare each Correction when it makes it: with the disposition of a fixed finding, and with the declaring command for a repaired QA finding, a correction commit, a Corrective Task or a plan change that repairs a defect. It defines the five classes with one example each, says that the declarer may be the session that wrote the plan and must still answer `plan` when the plan was at fault, and says that `unclear` with what is known is better than a guess.
22. **Documentation and glossary.** The user guide and the Roundfix skill's references describe the declaring command, the fixed disposition with a class, the setting, the new park and status line, the phase entry's new values and the report with its JSON. `CONTEXT.md` defines the four new terms and revises **Phase Record** and **Review Finding Disposition**, and ADR-0407 carries a notice naming ADR-0418.

## User Experience

The lines below show the content of each surface. The TechSpec fixes the command and flag names and the exact wording.

A Supervisor fixes a review finding and declares its cause in the same command, then declares a repaired QA finding:

```text
$ roundfix review dispose F1 --fixed-by 4c1d9e2 --class plan --artifact techspec --reason "The TechSpec told the reader to take Archive Records from the working tree"
<the disposition line, as today>
Correction C1 of 0412-the-report-reads-archive-records-from-git: plan (techspec); covers finding F1 and commit 4c1d9e2.

$ roundfix corrections declare --spec 0412-the-report-reads-archive-records-from-git --class implementation --commit a8b702c --reason "The note was cut in the middle of a multi-byte character"
Correction C2 of 0412-the-report-reads-archive-records-from-git: implementation; covers commit a8b702c.
```

In a repository with the setting on, an item whose correction has no class does not merge:

```text
$ roundfix deliver status
...
Park: 0412-the-report-reads-archive-records-from-git <class>: declare a Defect Class (plan, implementation, environment, operator or unclear) for commit d23cf4b "docs: correct the QA contract of task_04" with roundfix corrections declare --spec 0412-the-report-reads-archive-records-from-git, then run roundfix deliver retry 0412-the-report-reads-archive-records-from-git
```

The maintainer reads a Correction and declares it again:

```text
$ roundfix corrections declare --spec 0412-the-report-reads-archive-records-from-git --correction C2 --class plan --artifact tasks --by maintainer --reason "task_03 required a cut by bytes"
Correction C2 of 0412-the-report-reads-archive-records-from-git: plan (tasks); declared 2 times, class changed from implementation.
```

An Author Subagent records its phase:

```text
$ roundfix phase start prd --spec 0413-a-ledger-export --model claude-fable-5-2 --author-subagent
Recorded prd of 0413-a-ledger-export for claude session 3f9a1c20: model claude-fable-5-2 (reported), effort high (environment), run by the Author Subagent
```

The report, in a repository three weeks after the release:

```text
$ roundfix corrections report
Corrections of 31 Specs delivered from 2026-11-19 to 2026-12-10, read from /Users/operator/.roundfix.

Spec  Type  Tasks  idea                      prd                      techspec                  tasks                     Corrections                        Plan defects  Beside
0405  feat  4      claude-fable-5-2 / xhigh  claude-fable-5-2 / high  claude-fable-5-2 / high   claude-opus-5-5 / high    3: plan 1, implementation 2        techspec 1    dismissed 1; failed attempts 2 (implementation-defect 1, repository-convention 1)
0406  fix   2      -                         -                        claude-opus-5-5 / medium  claude-opus-5-5 / medium  1: environment 1                   -             dismissed 0; failed attempts 0
0407  fix   3      -                         -                        unknown                   unknown                   2: implementation 1, undeclared 1  -             dismissed 0; failed attempts 1 (unclassified 1)
...

Phase     Type  Model / effort             Specs  Tasks  Plan defects of its artifact
idea      feat  claude-fable-5-2 / xhigh   9      31     not counted
prd       feat  claude-fable-5-2 / high    9      31     1/9 Specs, 1/31 Tasks
techspec  feat  claude-fable-5-2 / high    9      31     2/9 Specs, 2/31 Tasks
techspec  fix   claude-opus-5-5 / medium   20     44     1/20 Specs, 1/44 Tasks
tasks     feat  claude-opus-5-5 / high     9      31     0/9 Specs, 0/31 Tasks
...
Comparison, techspec, feat: insufficient sample (claude-fable-5-2 / high 9 Specs; 20 each needed).
Comparison, techspec, fix: insufficient sample (claude-opus-5-5 / medium 20 Specs and no second group; 20 each needed).

Not counted: 2 Specs with an unknown model for a phase; 1 undeclared Correction; 0 Specs delivered without a record on this machine.
```

With two groups of 20 delivered Specs or more for one phase and Delivery Type, the comparison line holds both counts and nothing else:

```text
Comparison, techspec, feat: claude-fable-5-2 / high 3/22 Specs, 3/86 Tasks; claude-opus-5-5 / high 5/24 Specs, 5/61 Tasks.
```

## Non-Goals / Out of Scope

- Roundfix choosing or alternating the authoring model: the maintainer chose to compare the Specs as they come.
- Rebuilding the past: nothing recorded which model authored a Spec before this one or why it was corrected, and the study remains the record for Specs 0370 to 0403.
- The cost of a phase: the Phase Record still holds no token count and no price, and no command reads a harness transcript.
- A class assigned by the judge or by a signature, and any second reading of a declared class by a model.
- A stop that needed no repair, such as a machine out of memory or an owner that went down: Park Class already names it.
- A defect found after the merge: a Corrective Spec repairs it, and linking it back to the plan that caused it is a later request.
- Facts that travel with the Spec in Git, and any change to the Archive Record: the report reads this machine's records within Run Retention.
- The Correction Signs of a delivery made outside the Delivery Queue, and a refusal of a merge made by hand.
- A refusal when a fix is recorded or committed: the one refusal is before the Delivery Queue's merge.
- The model of a phase run in another harness, or reported from a session the command does not identify.
- A model or an effort on a prompt entry, on the `follow-up`, `delivery` and `release` phases, or on anything a Run records: the Run side already records them.
- A percentage, a confidence interval, a statistical test, a ranking or a verdict on a model.
- The length of the TechSpec, or any size of a Spec other than its number of Tasks.
- A setting for the minimum of 20, for the length of a reason or for the list of classes.
- Changing `roundfix runs causes`, Cause Class, the Classification Signatures or Park Class.
- Turning the setting on in an adopter: each Gesttione project does it in its own Project Config, in its own session.
- A change to the Baseline guides: the Roundfix skill carries the instruction to declare.

## Success Metrics

1. Success Metric: in a temporary repository and home, a Correction is declared for a Spec with each of the five classes; a `plan` declaration without an artifact, an artifact with another class, an unknown class, a blank reason, a reason of 201 characters and a reason with a line break are each refused with a message that names the accepted values, and a refusal writes nothing.
2. Success Metric: one commit named by two declarations counts as two Corrections, and two commits named by one declaration count as one, in the report's text and JSON.
3. Success Metric: recording a finding as fixed with a class and a reason appends the disposition line it appends today and records one Correction that covers the finding and its commit; recording it without a class prints and writes the same bytes as at the delivery base; a dismissal with a class is refused.
4. Success Metric: a Correction declared again by its identifier counts under its last declaration, the report says that it was declared twice and that its class changed, both declarations remain in the record, and a withdrawn Correction is counted apart and leaves its sign `undeclared`.
5. Success Metric: for a fixture Item Branch that repeats the commits of pull request #892 (three Task commits, a correction, a correction of the plan, the QA Report commit, a correction, the archive commit and the Pull Request Record Commit), the Correction Signs are exactly the three corrections; a merge of the default branch and the queue's own derived merge are not signs; a fixed finding and a Corrective Task each are one.
6. Success Metric: with the setting on, a Delivery Queue item with one Correction Sign and no Correction does not merge: it parks, its Park line names the sign, the five classes, the declaring command and the Delivery Retry, and `roundfix deliver status` showed the sign before the park; after the declaration a Delivery Retry merges the item, and an item whose signs are all covered merges without a park. A commit subject of 300 characters with a line that reads as an instruction is printed cut at 72 characters on one line.
7. Success Metric: with the setting off or absent, the same item merges, `roundfix deliver status` prints the bytes it prints at the delivery base, and the report counts the sign as one `undeclared` Correction; the setting in User Config is ignored with a warning, and a value that is not a boolean is refused by name.
8. Success Metric: `roundfix phase start` for each of the four Authoring Phases, run with a Claude Code session in its environment, a reported model and an effort in the environment, writes an entry with both values and how each was learned and with who ran the phase; without a reported model the model reads `unknown`; with no session both read `unknown`; a model of 300 characters is recorded as `unknown` with one warning; every case exits `0`, the refusals and exit codes of the delivery base are unchanged, and an entry written at the delivery base is still read, as `unknown`.
9. Success Metric: a phase entry written in a clone that is then deleted is still matched to its repository by the report, an entry of another repository with the same Spec number is not, and no entry holds the user name or the password of a remote address that carries them.
10. Success Metric: in a fixture of delivered Specs, the report's rows show for each Spec its Delivery Type, its Tasks without the QA gate, each phase's model and effort, its Corrections by class and its plan defects by artifact, in Spec number order; the phase's model is the last entry before the first Implement Run, a phase recorded with two models before it is marked, and a `prd` defect of a Spec with no `prd` entry counts for its `techspec` entry.
11. Success Metric: with groups of 20 and 20 delivered Specs for one phase and Delivery Type the comparison line holds both groups' counts; with 20 and 19, and with one group alone, it reads `insufficient sample` with each count and the 20 needed; no output of the report holds a percent sign or the words better, worse, best or worst; the JSON carries each numerator and denominator as integers.
12. Success Metric: the report states its window, the Roundfix Home it read, the Specs with an unknown model, the `undeclared` Corrections and the Specs it left out; a Spec whose delivery the machine did not record shows its declared Corrections and enters no count of signs; an unreadable Phase Record is named as unavailable and the report exits `0`.
13. Success Metric: the report shows dismissed findings and failed Verification attempts by Cause Class beside each Spec's Corrections and adds neither to a count; for the same fixture, `roundfix runs causes` prints the same text and JSON, with the same signature digest, as at the delivery base.
14. Success Metric: across a run of the report no file under the repository or the home changes and no network request, judge request or transcript read is made; a reason that reads as an instruction is printed on one line as data.
15. Success Metric: a delivery's declarations and signs are still reported after a new Delivery Queue replaced the queue that delivered it and after its Item Branch was released, and the Run Retention Sweep removes them once they are older than its cutoff.
16. Success Metric: the user guide, the Roundfix skill's references and the four authoring skills describe what Core Features 13, 21 and 22 name; `CONTEXT.md` defines **Correction**, **Defect Class**, **Correction Sign** and **Corrections Report** and the revised **Phase Record** and **Review Finding Disposition**; ADR-0407 carries an "Amended by ADR-0418" notice; the documentation, skill and glossary contracts pass.
17. Success Metric, outside evidence: the report's two rules of presentation are the ones published reporting standards state, as `## Acceptance evidence` quotes them: a result as its two counts, and a stated minimum with the reason printed when a number is withheld.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write, read through Exa on 2026-11-17.

- Moher, Schulz and Altman, "The CONSORT statement: revised recommendations for improving the quality of reports of parallel-group randomised trials", The Lancet, 2001 (https://www.thelancet.com/journals/lancet/article/PIIS0140-6736(00)04337-3/fulltext), checklist item 16: "State the results in absolute numbers when feasible (eg, 10/20, not 50%)." This is why a rate is printed as a count over a count.
- National Center for Health Statistics, "Data Presentation Standards for Proportions", Vital and Health Statistics, Series 2, Number 175, 2017 (https://www.cdc.gov/nchs/data/series/sr_02/sr02_175.pdf): "Estimates with either a denominator sample size or an effective denominator sample size (when applicable) less than 30 should be suppressed." and "When an estimate is flagged or suppressed, a footnote indicating the reason the estimate has been flagged or suppressed should be provided in the publication." This is why a comparison below the minimum is withheld and the line says why. The standard's minimum is 30 for one proportion; the maintainer chose 20 as a floor for showing two counts side by side, and every count is still printed.
- Chillarege et al., "Orthogonal Defect Classification — A Concept for In-Process Measurements", IEEE Transactions on Software Engineering, 1992 (https://chillarege.com/assets/files/ODC-Concepts-Breakthrough-Paper-8c0458e4ced84069fcf7b8b71a203532.pdf): "A programmer making the correction usually chooses the defect type", and "If the number of classes is small, there is a greater chance that the human mind can accurately resolve between them". This is why whoever corrects declares, and why the list has five classes.

The study, the two measurements and the reading of the source that this planning change commits under `references/` are this Spec's own evidence, not outside evidence; `## Validation` reports them.

## Unreachable Acceptance

- criterion: over the first 20 Specs this repository delivers after the setting is on, the report's `undeclared` count is 0 and no `plan` Correction lacks its artifact.
  criterion: over the same 20 Specs, at most one Correction in ten is declared `unclear`.
  criterion: at least 95 of every 100 entries of an Authoring Phase written after the release hold a model and an effort.
  criterion: of the first 20 Corrections the maintainer reads, at least 16 keep their class.
  criterion: the operator answers "plan defects per authoring model" for those 20 Specs from the report alone and records what was missing.
  criterion: one adopter that leaves the setting off delivers its first Spec after its update with no new refusal and no new required input.
  criterion: this repository's Project Config holds the setting on.
  reason: each needs the released binary, real deliveries over days or weeks, the maintainer's reading, another repository, or a change of this repository's Project Config that cannot arrive before the release is the installed binary.
  satisfied-by: the operator turns the setting on after the release, reads the report after 20 delivered Specs and records the counts, the maintainer's declarations and the adopter's first delivery in the session; the Verdict of the Idea says what a miss changes.

## Validation

The decisions were checked against these sources on 2026-11-17.

- The study ([`references/2026-11-17-fable-and-opus-authoring-study.md`](references/2026-11-17-fable-and-opus-authoring-study.md)): 23 delivered Specs and 24 corrections; 4 Specs authored on Fable and 14 on Opus with transcripts; corrections per Spec 1.50 against 1.14 and per Task 0.40 against 0.48, so the order flips with the denominator; three Specs whose correction nobody could classify afterwards. The delivery of Spec 0399 had three corrections and five causes, and the review of Spec 0398 returned a defect of the plan and one of the implementation in one round, the two cases behind one Correction per cause.
- The first measurement ([`references/2026-11-17-phase-record-and-causes-view-measurement.md`](references/2026-11-17-phase-record-and-causes-view-measurement.md)): 15 phase entries with one session and no model; `roundfix runs causes` listing 10 items in 31 Runs, none of them a correction of the delivery; the nine commits of pull request #892 reaching the default branch as one.
- The second measurement ([`references/2026-11-17-item-branch-commits-and-author-session-measurement.md`](references/2026-11-17-item-branch-commits-and-author-session-measurement.md)). Each commit a Run or the Delivery Queue makes on an Item Branch carries a mark: a Task trailer, the QA Report's message, the archive's change, the Pull Request Record Commit's change or the derived merge's trailer. A correction made by hand carries none, and of the nine commits of pull request #892 exactly the Supervisor's three carry none. This settles the Idea's first open question: a correction commit can be told. In a second Author Subagent session the environment again held the reasoning effort and no model name, and the session is told its model. A queue item stores only its candidate heads, a new queue replaces a finished one, and nothing removes the review records, so the signs are recorded by the merge.
- Characterization at `74d0e3a9`, read by the Research Subagent in two sweeps and checked by the author. The command that records a disposition accepts a fix by a commit that descends from the reviewed head and a dismissal with evidence, appends one line to a ledger and stores no cause. A Delivery Retry that accepts a moved head records that head and nothing else about its commits. The phase command loads no configuration and reads the session from one environment variable. Configuration loading refuses an unknown key, and one existing delivery key is a boolean that Project Config may set. `roundfix runs causes` reads Runs in a window of Run creation dates and assigns a Cause Class by signature.
- The web, read by the Research Subagent through Exa and re-read by the author for the two passages quoted from the CONSORT statement and the NCHS standard. IEEE Std 1044-2009 gives a defect an "Insertion activity" whose sample values are phases, and none of them says that the origin is not known. In a survey experiment, respondents made to answer with no "do not know" option chose the midpoint, and about 30% said they gave an incorrect answer to go on (Gutenberg and Lim, 2017, https://doi.org/10.20429/amtp.2017.47), which supports keeping `unclear`; a randomized design found that required fields did not impair reliability (Faran and Zanbar, 2019, https://doi.org/10.1080/13645579.2019.1630899), which supports a requirement that offers it. A study of 190,558 bug reports found that about 80% had a field reassigned later (Xia et al., 2014, https://xin-xia.github.io/publication/wcre141.pdf), which supports the later declaration. The sweep found no study of a required classification at merge time; that is an absence in what was read, not proof.
- The Secondbrain, read on 2026-11-17 (`wiki/index.md`, then two `qmd query` runs). `wiki/concepts/agentic-coding-e-unknowns.md` says that when a long task comes back wrong the cause may be the operator's brief and not a weak model, which supports the `operator` class and a report that names no better model. `wiki/concepts/the-art-of-agile-development-capitulo-21-16-quality.md` lists programmer, design and requirements errors as separate things to prevent, beside the heading "Don't Play the Bug Blame Game". `wiki/concepts/avaliacao-de-llms-evals.md` records a judge model agreeing with people at kappa 0.32 to 0.53, which kept a second reading by a model out of this version. The other hits were mirrors of this repository. A digest of the external sources this phase added went to the project that owns the Secondbrain on 2026-11-17 (ai-memory message `01a1eade-d6cc-7f92-a8fd-b2701604e822`), after a duplicate check that found no page on them.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.81)
- unverified_assumption: no (0.22)
- market_question: no (0.07)
- verdict: indicated
- stage: ran
- shape: changed — a correction became one cause, the list of classes gained `unclear`, a plan defect names its artifact and counts against the model of that phase, the declaration is required only where a repository turns it on, and the view shows its sample and withholds a comparison below a minimum.

The judged text is the request as the operator recorded it: [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md).

## Glossary

- adds: **Correction**
- adds: **Defect Class**
- adds: **Correction Sign**
- adds: **Corrections Report**
- changes: **Phase Record**
- changes: **Review Finding Disposition**

## Decisions

The maintainer's answers of 2026-11-17, as the operator recorded them in [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md):

- 2026-11-17: How should the record of the cause of each correction be routed? — "Spec já".
- 2026-11-17: Which authorization applies to Spec 0404? — "Mesmas do BE (Recommended)": push, pull request, merge, in the single release if ready in time; up to five corrections; authored and delivered after Spec 0403.
- 2026-11-17: Do the authoring models change meanwhile? — "Então podemos manter, vai ser interessante para saber o consumo até lá."
- 2026-11-17: Which shape does the first version take? — "Completo e honesto (Recommended)": the three parts, with the number of Specs, their type and their size beside every rate and "amostra insuficiente" below a minimum.
- 2026-11-17: What is the unit of a correction? — "Uma por causa (Recommended)".
- 2026-11-17: Who declares the class? — "Quem corrige (Recommended)": whoever corrects declares the class and the reason; the record keeps who declared; the maintainer may declare again; the last declaration counts.
- 2026-11-17: Is the class required? — "Só onde é ligado e por default é desligado, nos projetos da gesttione vai ser configurado para ligar".
- 2026-11-17: Does a plan defect name its artifact? — "Sim (Recommended)": the PRD, the TechSpec or the Tasks, counted against the model of that phase.
- 2026-11-17: How do the models get comparable Specs? — "Comparar o que vier": no alternation; the view shows the model by Delivery Type and may stay inconclusive.
- 2026-11-17: Is the draft of the Idea approved? — "Aprovar e seguir (Recommended)".
- 2026-11-17: What is the minimum for a comparison? — "20 Specs (Recommended)".
- 2026-11-17: Does this repository turn the setting on, and does the list keep a fifth class? — "Ligar e manter unclear (Recommended)".

The author settled the following from those answers, the approved Idea and its defaults, and the research of this phase, without a further question:

- The Idea's open questions. Roundfix knows an undeclared Correction by the three Correction Signs: the Idea's default of a fixed finding and a commit the delivery did not produce holds, because every commit of a Run or of the queue carries a mark, and a Corrective Task is added because it is found by its place in the graph. A plan change is a Correction when it repairs a defect and not when it records a new decision; it is a sign only as a commit on the Item Branch. Where the record is kept is the TechSpec's, within Core Feature 20. The session reports its model and the effort is read from the environment. The phase's model is the last entry before the first Implement Run. Two efforts of one model are two groups, with the model's total beside them. A defect of the Idea names the PRD. The size of a Spec is its number of Tasks. The report is a command of its own beside `roundfix runs causes`, which keeps every byte. The glossary takes the four names of `## Glossary`.
- A merge of the default branch into the Item Branch is not a sign: it answers another Spec that merged first, not a defect. A repair made inside one is declared like any Correction without a sign.
- The one refusal is before the Delivery Queue's merge. A correction commit runs no Roundfix command that could refuse, so an earlier refusal would cover one kind of correction only; the status line shows the sign from the moment it exists.
- The setting belongs to Project Config alone, as the Authoring Model Declaration does: a requirement is the repository's rule and travels with it.
- A Correction declared by mistake is withdrawn, not left to be counted: the record is append-only and the last declaration counts, so a withdrawal is one more declaration.
- This repository's setting is turned on by the operator after the release, not by a Task. The maintainer decided that this Spec turns it on; the installed binary of the Delivery Queue refuses a Project Config key it does not know, so the key cannot arrive in the change that teaches the binary that key.
- The model and the effort are recorded only for a session the command identifies, as the Idea states, and a reported value never costs the phase entry.
- A phase entry records an identity of its repository when it is written, because the clone an author works in is gone when the report reads the entry.
- No percentage is printed. The maintainer asked for the sample beside every rate, and the two counts are the rate.
- `type: feat` and Release Ceiling `minor`, the ceiling the cycle BE rules set. ADR-0417, ADR-0418 and ADR-0419 record the decisions; ADR-0418 amends ADR-0407.

## Open Questions

The TechSpec settled four of the six questions this PRD left open, and its Decisions section records each answer.

- Settled: a merge of the default branch is not a Correction Sign, and a repair inside one is declared by hand (ADR-0422). Still open: does such a repair go unseen too often? The operator reads the first 20 deliveries for one.
- Settled: the record keeps the declaring session and the name the declarer gives, and the report shows both and verifies neither (ADR-0424).
- Settled: the new blocker is `defect-class-undeclared`, of Park Class `finding` (ADR-0423).
- Settled: a session that is told no model passes none, and the entry reads `unknown` (ADR-0425).
- Do `unclear` and `undeclared` together pass a quarter of the Corrections in the first 20 Specs? Default: the list and the moment of declaration stand; past a quarter, the maintainer revisits them before more data is collected, as the Idea's verdict says.
- Does every comparison still read `insufficient sample` after the 30 days the record holds? Default: yes is an accepted outcome; the maintainer then decides whether to author some Specs of one type on another model or to take up the alternation the Idea names as its stretch goal.
