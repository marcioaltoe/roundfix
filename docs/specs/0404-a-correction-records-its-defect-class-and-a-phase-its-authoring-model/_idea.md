---
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
created: 2026-11-17
status: proposed
leverage: compounding
---

# A Correction records its Defect Class and a phase entry its authoring model

## Overview

Roundfix starts recording two facts it does not hold today, and one read-only view counts them.

The first fact is why each correction of a delivery happened. A Correction is one cause repaired while a Spec is being delivered: a pre-PR review finding that was fixed, a QA finding that was repaired, or any other defect that needed a change to the delivery candidate or to the plan. Whoever makes the correction declares its Defect Class, one of `plan`, `implementation`, `environment`, `operator` or `unclear`, with a one-line reason. A `plan` defect also names the artifact that was wrong: the PRD, the TechSpec or the Tasks.

The second fact is which model wrote each part of a plan. A phase entry of the Phase Record gains the model and the reasoning effort that ran the phase.

The view joins the two: Corrections per Spec by Defect Class, and plan defects counted against the model that authored the artifact they name.

It is for the person who chooses which model authors each phase and pays for that choice. In this repository that is the maintainer. An adopter gets the same record and view once it turns the requirement on in its Project Config; the maintainer will turn it on in the Gesttione projects.

It is valuable because the question it answers cannot be answered from what Roundfix records. On 2026-11-17 the maintainer asked whether Fable is effective for Ideas, PRDs and TechSpecs compared with Opus. The operator answered with a one-off study that read harness transcripts, commit subjects, review artifacts and memory notes for Specs 0370 to 0403 ([`references/2026-11-17-fable-and-opus-authoring-study.md`](references/2026-11-17-fable-and-opus-authoring-study.md)). It found 4 delivered Specs authored on Fable and 14 on Opus, no measurable difference in plan defects, and about five times the output tokens for Fable. The cause of every correction was rebuilt by hand, and three of the 23 delivered Specs still carry a question mark.

V1 is one Spec. It records from its release on and rebuilds nothing from the past. The view is honest about its sample: every rate is a count beside the type and size of the Specs behind it, and a comparison between two models appears only above a stated minimum. The maintainer chose to compare the Specs as they come, without steering which model authors which Spec, so the view may read "insufficient sample" for weeks, and it says so.

## Problem

The request, as the operator recorded it on 2026-11-17: "each correction made during a Spec delivery (a fix commit on the item branch, a review disposition, a QA repair) carries a defect class — plan, implementation, environment or operator — with a one-line reason, recorded where Roundfix can read it", "the Phase Record entry stores the model and reasoning effort of the session or Author Subagent that ran the phase", and a read-only view "counts corrections per Spec by defect class and by the authoring model of its TechSpec and Tasks, so a `git log`-free report answers 'plan-defect rate per authoring model'". The full text is in [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md).

**Why a correction happened is recorded nowhere.** After a Run ends, a delivery candidate is corrected in several ways: a pre-PR review finding is fixed, a QA finding is repaired, a check that failed is made to pass, a parked item is corrected before its Delivery Retry. The only one Roundfix records is the Review Finding Disposition, which says that a finding was fixed by a commit or dismissed with evidence. It has no field for why the finding existed. `roundfix runs causes` reads none of these: it lists failed Verification attempts and Corrective Tasks inside Runs. For the Runs created since 2026-11-09 it lists 10 items in 31 Runs, all of them failed Verification attempts, while the study counted 24 corrections in the 23 Specs delivered in about the same period. Its Cause Class is assigned by signature and never by judgment (ADR-0214), so it also leaves most items without a class: 100 of 132 in the measurement of 2026-09-17 to 2026-10-02 (`docs/references/corrective-causes-measurement.md`).

**The Phase Record does not say which model ran a phase.** A phase entry holds the phase, the Spec, the repository, the harness and the session. On the Mac mini on 2026-11-17 the record held 15 phase entries, for Specs 0399 to 0404, and every one carries the same session: an Author Subagent's command reads its parent's session, so the record cannot tell a phase the session wrote from one it handed to a subagent on another model (ADR-0407 says so in its consequences). The delivery side already records the runtime, model and reasoning effort of each Run. The authoring side records neither. The measurements are in [`references/2026-11-17-phase-record-and-causes-view-measurement.md`](references/2026-11-17-phase-record-and-causes-view-measurement.md).

**The reconstruction is slow and leaves holes.** The study needed a session of its own. For Specs 0365 to 0379, authored on another machine, it had no transcripts, so its table holds no authoring cost for them. Its scan of file writes missed the plans Fable wrote through a shell. For Specs 0377, 0383 and 0386 nobody could say afterwards whether a correction came from the plan or from the implementation.

**A count per Spec hides the mix, and a bare rate misleads.** The review of Spec 0398 returned a plan defect (the TechSpec prescribed reading Archive Records from the working tree) and an implementation defect. Spec 0399 had three corrections and five causes. Counted per Spec, the Fable Specs had 1.5 corrections each against 1.1 for the Opus Specs; counted per Task, 0.40 against 0.48. The order flips with the denominator, and four Specs settle neither.

**Model and Spec type move together.** The four Fable Specs were all features, all with an Idea, and their TechSpecs ran from 5,369 to 10,443 words against 1,925 to 5,507 for the Opus Specs, 11 of which were fixes. The phases of a feature run on the declared model, and a fix is authored in the session on the session's model. Spec 0402 then declares Fable for the PRD and the TechSpec and Opus for the Tasks, so every feature Spec will mix both models, and a class recorded per Spec could not say which of them erred.

### Market Data

- The classic work on defect classification warns about exactly this kind of label. Orthogonal Defect Classification says of schemes that ask where a defect was injected: "Such direct classification schemes, by the nature of their assumptions, qualify as good opinion surveys, but do not constitute a measurement on the process." It also says who labels in practice, "A programmer making the correction usually chooses the defect type", and that "If the number of classes is small, there is a greater chance that the human mind can accurately resolve between them" (Chillarege et al., 1992).
- Agreement between two people who classify the same defects ranges widely. Published values of kappa run from 0.66 and 0.82 for trained inspectors, rising when two confusable classes are merged (El Emam and Wieczorek, 1998), to 0.44 to 0.53 with developers of mixed experience (Matalonga et al., 2009) and about 0.2 for defect categories in infrastructure code (NCSU, 2019).
- An industrial scheme described by Freimut (2001) has developers pick three descriptors for each defect: its origin, where it was injected, its type and its mode. An origin label is established practice; its reliability is the open point.
- A sample this small cannot separate two rates. With no event in n trials the 95% upper bound of the rate is about 3/n (Jovanovic and Levy, 1997). Telling 5% from 30% with 80% power needs 39 per group, 10% from 30% needs 69, and 10% from 50% needs 23 (Kim, 2016, Table 1).
- The trial that measured AI assistance on real work assigned each of 246 tasks at random and recorded the developer's forecast of its difficulty before the assignment (METR, 2025). It controls what this repository's routing does not: which arm gets which kind of work.
- A commit trailer is a proven carrier of machine-read facts: the Linux kernel's `Fixes:` tag feeds its stable tooling. A squash merge collapses a branch's commits into one, and reports of trailers lost that way are easy to find. This repository squash-merges each delivery: the nine commits of pull request #892, two of them `fix` commits, reached the default branch as one.
- Conventions that record the AI behind a change record a model and no reasoning effort: the kernel's `Assisted-by: AGENT_NAME:MODEL_VERSION [TOOL1] [TOOL2]`, the Agent Trace draft's `model_id`, git-ai's notes, Aider's co-author line.
- In the comparisons read of spec-driven tools (Spec Kit, Kiro, OpenSpec, BMAD), none records why a piece of rework happened or which model wrote a planning artifact. DORA's rework rate counts deployments, not causes. This is an absence in the articles read, not proof.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | A Correction is one cause | Critical | From the moment a Spec's plan is on the default branch until its delivery merge, each cause that needed a repair is one Correction: a fixed review finding, a repaired QA finding, or another defect that needed a change to the delivery candidate or to the plan. A commit that repairs two causes holds two Corrections, and two commits for one cause are one. |
| F2  | A Defect Class and a reason | Critical | A Correction carries one class from a closed list and a one-line reason. `plan`: the plan said something wrong, missing or impossible. `implementation`: the plan was right and the delivered change was not. `environment`: the cause lies outside the Spec, such as the machine, a tool or a failure already on the default branch. `operator`: an action of the Supervisor or of a person during the delivery caused it. `unclear`: the declarer cannot tell, and the reason says what is known. |
| F3  | A plan defect names its artifact | Critical | A Correction of class `plan` names the PRD, the TechSpec or the Tasks. The view counts it against the model that authored that phase. |
| F4  | Whoever corrects declares | Critical | The Supervisor or person who makes a correction declares its class and reason at that moment: together with the disposition of a fixed review finding, and in a declaration of its own for a repaired QA finding or any other cause. The record keeps who declared and when. The Roundfix skill tells a Supervisor when and how to declare. |
| F5  | Required only where a repository turns it on | Critical | A Project Config setting, off by default, makes the declaration required. With it on, a delivery does not reach its merge while one of its Corrections has no class, and the refusal names the Correction, the accepted values and the command. With it off, nothing is required: a class that is given is recorded, and the view counts the others as `undeclared`. This repository turns it on. |
| F6  | A later declaration replaces an earlier one | High | Anyone may declare a Correction's class and reason again, the maintainer above all. The view counts the last declaration, shows that the Correction was declared more than once, and the earlier declarations stay in the record. |
| F7  | A phase entry holds the model and the effort | Critical | An entry of an Authoring Phase records the model that ran the phase, as the harness names it, its reasoning effort, and how Roundfix learned them: reported by the session that ran the phase, or `unknown`. It also tells a phase an Author Subagent ran from one the session ran itself. The authoring skills report both when they record the phase. |
| F8  | The view counts Corrections and plan defects | Critical | A read-only view, in text and in JSON, shows for each Spec delivered in its window the type, the number of Tasks, the model and effort of each Authoring Phase, its Corrections by Defect Class and its plan defects by artifact. It then shows, for each Authoring Phase, model, effort and Spec type, the Specs and Tasks authored and the plan defects that name that phase's artifact. |
| F9  | The view shows its sample and withholds a comparison | Critical | Every rate is a count over a count, beside the Spec type and the size in Tasks, never a bare percentage. Two models are set side by side only for the same phase and Spec type and only when each has at least 20 delivered Specs; below that the line reads `insufficient sample` with both counts. The view states its window and its machine, counts the Specs whose model is `unknown` and the Corrections that are `undeclared`, and never names a better model. |
| F10 | Shown beside, never counted | High | A dismissed review finding is not a Correction, because nothing was corrected. A failed Verification attempt that the Agent repaired inside its Run is not one either and keeps its Cause Class. The view shows both beside the Corrections of each Spec. Cause Class, its signatures and the present output of `roundfix runs causes` do not change. |

## KPIs

Each is read from the view over the first 20 Specs this repository delivers after the release, unless the row says otherwise.

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Corrections with a declared class and reason, and an artifact when the class is `plan` | 100% in this repository, where the setting is on (today: 0%, no class is recorded) | The view's `undeclared` count is 0 and no `plan` Correction lacks its artifact. |
| Corrections declared `unclear` | <= 10% (the study could not classify the corrections of 3 of 23 delivered Specs, 13%) | The view's count of `unclear` over all Corrections. |
| Authoring phase entries that hold a model and an effort | >= 95% of the `idea`, `prd`, `techspec` and `tasks` entries written after the release (today: 0 of 15) | Count the entries of the Phase Record in the window; the rest read `unknown`. |
| Declared classes that survive the maintainer's reading | >= 16 of the first 20 Corrections he reads keep their class | The view's count of Corrections declared again with a different class. |
| Sources outside the view needed to answer "plan defects per authoring model" | 0 (the study read four: harness transcripts, commit subjects, review artifacts and memory notes) | The operator answers the study's question for those 20 Specs from the view alone and records what was missing. |
| New refusals or required input in a repository that leaves the setting off | 0 | A fixture with the setting off delivers a Spec with corrections and no declaration, and one adopter's first delivery after its update is read. |

No KPI promises that the comparison between models becomes conclusive. Under the maintainer's choice to compare the Specs as they come, that depends on which Specs each model happens to author.

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Strong: the choice of an authoring model that costs about five times the output tokens rests today on a hand study with four Specs in one arm, and the cause of a correction, which also shows where the authoring skills fail, is lost when the session ends. |
| Reach           | What % of users would this affect?                  | Maybe: the record of Corrections serves every repository that delivers through the Delivery Queue, but only where the setting is on, and the model of a phase is recorded only for Claude Code sessions, the one harness `roundfix phase start` identifies. |
| Frequency       | How often would users encounter this value?         | Strong: the study counts about one correction per delivered Spec (24 in 23), and this repository delivered 58 Specs in the 14 days to 2026-11-14. The view itself is read at decision points, a few times a month. |
| Differentiation | Does this set us apart or just match competitors?   | Strong: no tool read records the cause of rework or the model of a planning artifact. Attribution conventions record a model per commit and stop there. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong: each delivered Spec adds to a record that cannot be rebuilt afterwards. The study shows what rebuilding costs and what it loses. |
| Feasibility     | Can we actually build this?                         | Strong for the declaration, the setting and the view, which extend what Roundfix already records. Maybe for the model of a phase: a command run by a session can read the session's reasoning effort and finds nothing that names its model, so the model is what the session reports. |

Verdict: proceed in the shape the maintainer's six answers gave. Three results would change it. If `unclear` and `undeclared` together pass a quarter of the Corrections in the first 20 Specs, the list of classes or the moment of declaration is wrong and is corrected before more data is collected. If sessions cannot report their model reliably, the entry records the declared model with a mark that says so, and the view shows the mark. If every comparison line still reads `insufficient sample` after the 30 days the record holds, the view has still answered why corrections happen, and the maintainer decides whether to author some Specs of one type on the other model or to adopt the stretch goal.

## Trade-off Insights

- **Recommended approach:** record the cause where the correction is already recorded and the model where the phase is already recorded, declared by those who have the context at that moment; keep every declaration open to a later one; and let the view report counts with their sample and refuse a comparison it cannot support.
- **Key trade-offs:**
  - A declared class against a class by signature. A signature can be shown and repeated, and it left 100 of 132 items unclassified. A declaration covers every Correction and is an opinion. The two stay separate facts under separate names, and the view never merges them.
  - The corrector declares against the maintainer or a judge. The corrector has the context and costs the maintainer no time. The corrector is also often the session that wrote the plan, and it is asked whether the plan was at fault.
  - One Correction per cause against one per commit or per round. Per cause is more to declare and is the only unit that keeps a plan defect and an implementation defect of the same round apart.
  - Required only where turned on against required everywhere. No adopter's delivery changes by surprise, and a repository that leaves the setting off collects a thin record.
  - Comparing the Specs as they come against steering the model. Nothing changes in how Specs are authored, and the model stays tied to the Spec type.
  - A record on this machine against facts that travel with the Spec. V1 adds nothing to the Archive Record, and the view reads only what this machine holds, 30 days by default: a Spec authored elsewhere is absent from it.
- **Risks identified:**
  - The declarer judges its own plan. The reason, the declarer and the re-declaration stay visible, and the fourth KPI measures how often the maintainer disagrees.
  - The comparison stays inconclusive. After Spec 0402 a feature's PRD and TechSpec are on Fable and every Tasks phase is on Opus, so within one Spec type there may be a single model per phase.
  - A session reports a wrong model or none. The entry says how the value was learned, and `unknown` is counted.
  - The label becomes a target. A rate read as a verdict on a model invites declarations that protect it. The view never names a better model, for that reason and because of its sample.
  - The reason is repository text printed in a report an Agent reads. The PRD states it as data, bounded in length, never as an instruction.
  - Three planned Specs change the same files first (see Prerequisite and Shared Files).
  - Run Retention removes the record after its window, so a comparison cannot reach further back than 30 days unless the stretch goal is built.
- **Dissenting view:** the author recommended that the maintainer author part of the feature Specs on Opus over the next weeks, because recording alone cannot separate the model from the Spec type. The maintainer chose "Comparar o que vier". The record is worth having either way, and the view will say that its sample is insufficient for as long as it is. A second dissent comes from the literature: a label of where a defect came from is an opinion, and counts that need no judgment (plan amendments after a Run, first Runs that end Clean, findings per review) would be harder to dispute, though they cannot tell a plan defect from an implementation defect. A third: recording both facts now and building the view only when a month of data exists would have cost less and risked less.
- **Stretch goal (V2+):** Roundfix alternates two declared models for a phase between Specs of the same type, which is the design that can answer the maintainer's question. After it, the facts of a Spec travel with it in Git, so the view outlives Run Retention and a change of machine, and the judge gives a second reading of each declared class.

## Opportunity Scan

- **Original** — a class on each correction, the model and effort in the Phase Record, and a view that counts both. Medium effort, one Spec. Impact Strong, Reach Maybe, Frequency Strong, Differentiation Strong, Defensibility Strong, Feasibility Strong for the record and Maybe for the model.
- **More ambitious** — a controlled comparison: Roundfix assigns the authoring model per Spec within a type and reports the difference with its uncertainty. Large effort, because the Authoring Model Declaration names one model per phase and a fix is authored in the session. Impact Must do for the decision; Feasibility Maybe.
- **Simpler** — record only: the class and the model are recorded now, and the study is repeated from the records in three or four weeks. Small effort. Impact Maybe: it keeps the hand count the request wants to end.
- **Adjacent** — plan quality from counts that need no declaration: plan amendments after a Run, first Runs that end Clean, findings per review, QA rounds, each per authoring model. Medium effort. Impact Maybe: nobody labels anything, and nothing tells a plan defect from an implementation defect.

**Chosen direction:** hybrid — the maintainer's pick, "Completo e honesto (Recommended)": the original three parts, with the sample, the Spec type and the size beside every rate and no comparison below a stated minimum. The more ambitious alternative is the stretch goal, the simpler one is the third dissent, and the adjacent one is the second.

## Integration with Existing Features

| Integration point | How |
| ----------------- | --- |
| Review Finding Disposition | A finding recorded as fixed gains a Defect Class and a reason. A dismissal is unchanged and is never a Correction. |
| The Delivery Queue and the Delivery Retry | With the setting on, a delivery with an undeclared Correction stops before its merge and says what to declare. Park Class is unchanged: it says why an item stopped, and a Defect Class says where a repaired fault came from. |
| The QA gate and the QA Report | A repaired QA finding is a Correction. The QA verdict and the report's table are unchanged. |
| `roundfix runs causes` and Cause Class | Cause Class and its signatures are unchanged, and `roundfix runs causes` still prints everything it prints today. The view shows failed Verification attempts by Cause Class beside the Corrections and merges neither into the other. |
| The Phase Record and `roundfix phase start` | A phase entry gains the model, the effort, how they were learned and whether an Author Subagent ran the phase. The status view of the Follow-up Command reads the same entries as today. |
| The authoring skills and the Roundfix skill | The four authoring skills report the model and effort when they record their phase. The Roundfix skill tells the Supervisor to declare each Correction. |
| The Authoring Model Declaration | Unchanged, and not the source of the recorded model: a fix is authored in the session, and a typed skill runs on the session's model, so the declared model is not always the one that ran. |
| Project Config | Gains the setting, off by default, so each repository decides for itself. |
| Run Retention | Sets the window of the view. The Archive Record is unchanged in V1. |

## Prerequisite and Shared Files

Specs 0401, 0402 and 0403 are planned on the default branch and are delivered before this one: the operator recorded the maintainer's authorization with "autoria e entrega depois da 0403". The PRD names them under its Prerequisites when a Verification of this Spec depends on what they deliver.

All of them are expected to change files this Spec also changes, and the PRD and the TechSpec confirm the list against their Task files:

- the Phase Record's writer, which Spec 0401 changes so that a phase entry no longer marks a boundary;
- `.roundfixrc.yml`, where Spec 0402 changes the Authoring Model Declaration and this Spec turns its setting on;
- `CONTEXT.md`;
- the Roundfix skill and its references on phases and modes, the authoring skills, their mirrors under `skills/`, and `skills/testdata/owned-skill-versions.json`;
- `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json`;
- the user guide pages of the phase command and of the configuration.

Spec 0402 also decides what the comparison can show: it declares Fable for the PRD and the TechSpec and Opus for the Tasks.

## Out of Scope (V1)

- **Roundfix choosing or alternating the authoring model** — the maintainer chose to compare the Specs as they come. It is the stretch goal.
- **Rebuilding the past** — nothing recorded which model authored Specs before this one or why they were corrected. The study remains the record for Specs 0370 to 0403.
- **The cost of a phase** — the Phase Record still holds no token count and no price (ADR-0407). The study's cost figures came from harness transcripts, which no Roundfix command reads.
- **A class assigned by the judge or by a signature** — a Defect Class is declared. No text leaves the machine, and Cause Class keeps its own rule.
- **A stop that needed no repair** — a machine out of memory or an owner that went down delays a delivery and corrects nothing. Park Class already names it.
- **A defect found after the merge** — a Corrective Spec repairs it, and linking it back to the plan that caused it is a later request.
- **Facts that travel with the Spec in Git** — the view reads this machine's record within Run Retention. It is part of the stretch goal.
- **The model of a phase run in another harness** — `roundfix phase start` identifies a Claude Code session only. Any other entry reads `unknown`.
- **Turning the setting on in an adopter** — each Gesttione project turns it on in its own Project Config, in its own session.
- **A verdict on a model** — the view reports counts. Which model to declare stays the maintainer's decision.

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

## Decisions

Every answer is the maintainer's, as recorded in [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md).

- Route of the request (2026-11-17): "Spec já".
- Authorization (2026-11-17): "Mesmas do BE (Recommended)" — as the operator recorded it: "push, PR, merge, entra na release única se pronta a tempo; até 5 correções; autoria e entrega depois da 0403". The PRD quotes it in the authorization record.
- The authoring models stay as recommended for now (2026-11-17): "Então podemos manter, vai ser interessante para saber o consumo até lá."
- Shape of V1 (2026-11-17, ~14:45 BRT): "Completo e honesto (Recommended)" — the three parts, and the view shows the number of Specs, their type and size beside every rate, and "amostra insuficiente" below a minimum.
- Unit of a correction (2026-11-17, ~14:45 BRT): "Uma por causa (Recommended)".
- Who declares (2026-11-17, ~14:45 BRT): "Quem corrige (Recommended)" — whoever corrects declares the class and the reason; the record keeps who declared; the maintainer may declare again; the last declaration counts.
- Whether the class is required (2026-11-17, ~14:45 BRT), in his own words: "Só onde é ligado e por default é desligado, nos projetos da gesttione vai ser configurado para ligar".
- A plan defect names its artifact (2026-11-17, ~14:45 BRT): "Sim (Recommended)" — the PRD, the TechSpec or the Tasks, counted against the model of that phase.
- How the models get comparable Specs (2026-11-17, ~14:45 BRT): "Comparar o que vier" — no alternation; the view shows the model by Spec type and may stay inconclusive. The author had recommended varying the model by hand.
- The draft of this Idea (2026-11-17, ~15:30 BRT): "Aprovar e seguir (Recommended)".
- The minimum for a comparison (2026-11-17, ~15:30 BRT): "20 Specs (Recommended)" — 20 delivered Specs per model, for the same phase and Spec type. It is a floor for showing two counts side by side, not a test: at 20 a share moves in steps of five points, and even a gap as wide as 10% against 50% needs 23 per group (Kim, 2016).
- This repository and the fifth class (2026-11-17, ~15:30 BRT): "Ligar e manter unclear (Recommended)" — this Spec turns the setting on in this repository, and the list of classes keeps `unclear`, so that nobody is made to guess.

The author's own, made in this stage and approved with the draft:

- The first question showed the judgment and offered the dismissal of the stage as its last option. The maintainer chose the first option.
- A dismissed finding and a failed Verification attempt inside a Run are shown beside the Corrections and are not counted as Corrections.
- The model is recorded as the harness names it, with the effort and how both were learned, and only for Claude Code sessions.
- The view covers what this machine holds and states its window.
- One Spec, of about four implementation Tasks and the QA gate. The two records and the view complete one another: neither answers the request alone.
- No ADR is written at this stage. ADR-0417 onward are reserved for the PRD and the TechSpec, which settle at least two decisions: a Defect Class is declared and stays apart from Cause Class (beside ADR-0214), and a phase entry holds the model and effort and a view reads them (an amendment of ADR-0407).

## Open Questions

- How does Roundfix know a Correction exists when nobody declared it, so that it can refuse or count it as `undeclared`? Default: a finding recorded as fixed, and a change to the delivery candidate that is not one of the commits the delivery itself produces. If the second cannot be told reliably, the author returns to the maintainer before narrowing what the setting requires.
- Is a change to the plan after it reached the default branch a Correction? Default: yes when it repairs a defect a Run or a gate met, with a class like any other; no when it records a new decision of the maintainer, as the amendment of Spec 0401 did.
- Where is the declaration kept so that it survives the squash merge of the delivery? Default: the TechSpec chooses. The request's two examples, a trailer on the fix commit and the review disposition, are candidates for where it is given, and a trailer alone does not reach the default branch.
- How does a session report its model? Default: the authoring skills pass the model the harness names when they record the phase, and the entry marks it as reported. ADR-0407 says an entry never holds "an environment value", and the effort is readable from the session's environment, so the amendment states what is allowed.
- Which entry gives a phase its model when the phase was recorded more than once for a Spec? Default: the last one before the Spec's first Run.
- Are two reasoning efforts of one model two groups? Default: yes. The view groups by model and effort together and shows the model's total beside them.
- Is a defect that comes from the Idea a plan defect of the PRD? Default: yes, because the PRD adopts the Idea.
- Is the number of Tasks enough as the size of a Spec? Default: yes, QA excluded, as in the study. The PRD may add the length of the TechSpec.
- Is the view part of `roundfix runs causes` or a command beside it? Default: the PRD decides, and the present text and JSON of `roundfix runs causes` keep every field.
- Which names does the glossary take? Default: the PRD's Glossary settles them. The candidates are Correction and Defect Class, kept apart from Cause Class, Corrective Task and Correction Review, and a name for the view.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md), the request, the Idea-stage judgment and the answers behind `## Decisions`, recorded by the operator.
2. [`references/2026-11-17-fable-and-opus-authoring-study.md`](references/2026-11-17-fable-and-opus-authoring-study.md), the operator's study: every count of Specs, Tasks, corrections, words and tokens in this file.
3. [`references/2026-11-17-phase-record-and-causes-view-measurement.md`](references/2026-11-17-phase-record-and-causes-view-measurement.md), recorded by the author: what the Phase Record and `roundfix runs causes` held on 2026-11-17, and the commits of pull request #892.

Read in this repository: `docs/references/corrective-causes-measurement.md` for the share of items without a Cause Class, ADR-0214, ADR-0401, ADR-0407 and ADR-0410, and the plans of Specs 0401, 0402 and 0403. The count of 58 Specs delivered in 14 days is from the baseline adopted by Spec 0397 (`docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/references/2026-11-16-delivery-speed-and-plan-quality-baseline.md@f4580437cd6f3c354c2524f7f2eca9c924446d64`).

This project's ai-memory messages were listed on 2026-11-17 and none was pending. The four open Backlog Entries and the two Findings belong to other contexts and stay where they are; the nearest, the cost of starting a session, waits on per-phase numbers and is about tokens, which this Spec does not record. The Research Subagent made two sweeps of the codebase and one of the web; the author checked what it returned, corrected one statement (it read `docs/specs` as empty) and re-read the passages quoted from Chillarege, Kim and the kernel documentation.

Secondbrain (read on 2026-11-17; `wiki/index.md`, then two `qmd query` runs):

- `wiki/concepts/avaliacao-de-llms-evals.md` — a judge model agreed with people at kappa 0.32 to 0.53 where people agreed at 0.55 to 0.73, and "um avaliador só é frágil". It kept the judge out of V1 and put the maintainer's re-declaration in.
- `raw/files/books/site-reliability-engineering-2016/17-effective-troubleshooting.md` — publishing a negative result "improves our industry's data-driven culture". It supports a view that says `insufficient sample` instead of a number.
- `raw/web/2026-08-08-arxiv-csai-agent-software-reliability.md` — a capture of papers on agent reliability; nothing in it on the cause of rework or on comparing models.

A digest of the external sources went to the project that owns the Secondbrain on 2026-11-17 (ai-memory message `01a1ea8a-0bd8-7ac2-b6cd-1b33d21e9ad8`), after a duplicate check that found no page on these topics.

External sources (read through Exa on 2026-11-17):

- Chillarege et al., "Orthogonal Defect Classification — A Concept for In-Process Measurements", IEEE TSE 1992, https://chillarege.com/assets/files/ODC-Concepts-Breakthrough-Paper-8c0458e4ced84069fcf7b8b71a203532.pdf — an origin label is an opinion, the person who fixes labels, few classes; it gave the reason line, the `unclear` class and the second dissent.
- El Emam and Wieczorek, "The repeatability of code defect classifications", ISSRE 1998, https://www.ehealthinformation.ca/web/default/files/wp-files/isern-98-09.pdf; Matalonga et al., CLEI 2009, https://clei.org/proceedings_data/CLEI2009/55640.pdf; NCSU TR-2019-2, https://techrep.csc.ncsu.edu/2019/TR-2019-2.pdf — the range of agreement; they set the fourth KPI.
- Freimut, "Developing and Using Defect Classification Schemes", Fraunhofer IESE 2001, http://www.inf.fu-berlin.de/inst/ag-se/teaching/S-ERROR-2004/Fre01.pdf — a scheme in use that asks for a defect's origin, type and mode.
- Jovanovic and Levy, "A Look at the Rule of Three", The American Statistician 1997, https://doi.org/10.1080/00031305.1997.10473947, and Kim, "Sample size calculation 2. Comparison of two independent proportions", 2016, https://pmc.ncbi.nlm.nih.gov/articles/PMC4868880/ — what a small sample can tell apart; they gave the minimum of F9.
- METR, "Measuring the Impact of Early-2025 AI on Experienced Open-Source Developer Productivity", https://metr.org/Early_2025_AI_Experienced_OS_Devs_Study-paper.pdf — assignment before the work and a difficulty measure; it challenged comparing the Specs as they come.
- Linux kernel, "Submitting patches", https://docs.kernel.org/process/submitting-patches.html, and "AI Coding Assistants", https://docs.kernel.org/7.1/process/coding-assistants.html — a trailer as a machine-read fact and a model named per commit.
- GitHub, "About pull request merges", https://docs.github.com/en/pull-requests/reference/pull-request-merges — what a squash merge keeps; it ruled out a branch trailer as the only carrier.
- Agent Trace 0.1.0, https://agent-trace.dev/; git-ai standard 3.0.0, https://github.com/git-ai-project/git-ai/blob/main/specs/git_ai_standard_v3.0.0.md; Aider, https://aider.chat/docs/git.html — a model identifier and no effort.
- Comparisons of spec-driven tools, https://www.getlockstep.dev/blog/spec-kit-vs-kiro-vs-openspec-vs-bmad and https://github.com/ecodan/cicadas/blob/master/docs/sdd-comparison.md, and DORA, https://dora.dev/guides/dora-metrics/ — no recorded cause of rework.
