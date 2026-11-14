---
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
status: active
created: 2026-11-14
surfaces: [cli, docs]
type: feat
---

# A feature is judged for the Idea stage before its PRD

The Idea stage, the study that asks whether a feature should be built and in
what shape, runs today only when someone asks for it: the routing rule sends
only a large or fuzzy initiative through `write-idea`, the agent decides
alone, and nothing records why the stage was skipped. On 2026-11-14 that
stage, started by hand, changed the shape of Spec 0392 before its PRD existed.
After this Spec, three criteria about a feature request are judged before its
PRD is written, by the Jev judge through `roundfix spec judge --stage idea`,
for a recorded item or for a request of the session handed in as a file, and
by the agent itself when the judge does not answer. An indicated stage starts without being asked and the user
may dismiss it with a recorded reason, an undecided one is the user's choice,
and nothing is ever refused. Every feature PRD records the judgment and what
followed, so a later Spec can measure the thresholds on each repository's own
history. Every adopter receives the rule through its next `roundfix baseline
update`.

## Sources

The maintainer's request, adopted by this Spec's planning change with the
operator's pilot (see `references/_index.md`), and the Idea the maintainer
approved:

1. [`references/2026-11-14-maintainer-request-and-answers.md`](references/2026-11-14-maintainer-request-and-answers.md),
   the request and every answer behind `## Decisions`; requested by the
   maintainer. No ai-memory message is adopted and nobody asked for a notice.
2. [`references/2026-11-14-idea-stage-judgment-pilot.md`](references/2026-11-14-idea-stage-judgment-pilot.md),
   the operator's four-case pilot, with its questions, answers and limits.
   The maintainer's later answer of 2026-11-14 on the data boundary, relayed
   by the operator, is quoted in `## Decisions` and in the authorization
   record.
3. [`_idea.md`](_idea.md), the approved Idea, read as this PRD's exploration:
   its features F1 to F5 are this Spec, and F6 and F7 are Spec 0395.
4. [`docs/references/idea-stage-judgment-measurement.md`](../../references/idea-stage-judgment-measurement.md),
   the measurement of the questions on this repository's history, recorded by
   the author of this PRD on 2026-11-14.

At the authoring base `docs/backlog/` holds one open entry, a report of three
unreproduced symptoms from another team, and `docs/findings/` holds none;
neither shares this context. The project's ai-memory inbox held one pending
message on 2026-11-14, a fix report about the Verification probe; it does not
concern this Spec and is left to the Supervisor's triage.

## Prerequisites

None in `requires`. Authored against `251ed81e` (v0.143.3) and rebased on
`b7eda503`, which adds the plans of Specs 0391 and 0392. Both were authored
beside this one and deliver before it, in that order. Files this Spec shares
with them, each edited in its own sentences:

- The Baseline spec-workflow module. Spec 0391 owns its grouping clause; this
  Spec changes only its three route clauses. Both raise the module version
  with the record command, and the Task regenerates the module version record,
  the Standard TypeScript profile pin and formatter golden, the catalog
  snapshots, the plan goldens and this repository's Spec routing guide and
  setup manifest from the tree it runs on.
- The owned skills `write-idea` and `write-prd` (Spec 0391 rewrites their
  grouping sentences, Spec 0392 adds research-subagent text), `roundfix`
  (Spec 0391 rewrites the grouping sentence of its `spec` reference, Spec 0392
  changes its modes and setup references) and `brainstorming`. This Spec adds
  sections under their own headings and the few sentences the TechSpec names.
  Each skill's version fields and the owned-skill version record are raised
  with the record command from the tree the Task runs on, never with a number.
- `CONTEXT.md`: each Spec its own terms; this Spec adds **Idea Stage
  Judgment** and changes **Advisory Judgment**.
- The `spec` command guide (Spec 0391 edits its grouping sentence) and the
  Behavior Surface Record, re-recorded as the last edit of the Task that
  changes a guide or a help text.
- `internal/cli/cli.go`: Spec 0392 changes the skills stage of the Baseline
  update there; this Spec adds one synopsis line to the top-level usage.

Spec 0392 generates author and research subagents. The Idea-stage judgment and
the Idea stage happen in the Supervisor session, before any delegation: an
author subagent cannot ask the user, so it never runs the stage and never
decides an undecided judgment. No Verification of this Spec depends on Spec
0392, so it is not in `requires`.

Spec 0395 follows this one, once recorded cases exist: the calibration report
(F6) and the judgment in the follow-up report (F7). The record this Spec
defines is what that report reads.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The judged
  item is named by the reference the Follow-up Command already prints, and the
  record is prose lines in the PRD. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the judgment sends one request to the
  judge's existing recipients with the existing keys, read from the
  environment and never printed, under the existing monthly ceiling; no
  endpoint, key or recipient is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0403 (new) records the rule: "The
  verdict starts the stage; it refuses nothing." ADR-0404 (new) records what
  the judge reads: "Nothing else can become a request". ADR-0200 applies
  unchanged: "A raised judgment is a finding the authoring model answers", the
  command exits zero whenever it ran, and answers count only from the pinned
  model. ADR-0201 is amended: by the maintainer's answer of 2026-11-14 a
  request of the session, handed in as a file, may reach a request, and a Task
  writes the notice; "Each request appends one line to a Judge Log" stays, and
  the recipients, the keys and the ceiling are its own. ADR-0319 applies unchanged: its request state is reused
  as it is, "one item's kind, type, origin, title and an excerpt of its text",
  and its sentence "Items are sent whatever their language" holds for this
  judgment too. ADR-0400 applies: the Spec holds one context in three
  implementation Tasks that each earn their verification cycle. ADR-0209 does
  not apply: it decides which Findings and Backlog Entries the grouping
  question reads through its own reader, and this judgment reads its item
  through the follow-up inventory and sends no Finding. ADR-0231 and ADR-0239
  apply unchanged: the ceiling is the User
  Config value and the judge reads its own key first. ADR-0283 applies
  unchanged: inside a Run the judge sends nothing. ADR-0353 applies unchanged:
  a Waiting Item is never sent to the judge. ADR-0208, which ADR-0400 amends,
  does not apply: it decides how sources are grouped into Specs, and this Spec
  changes no grouping clause. ADR-0235, which cites ADR-0231, does not apply:
  it decides which models implementation work may run on, and this Spec starts
  no Agent Session. ADR-0273, which cites ADR-0283, does not apply: it decides
  what the children of a Run may hold and what a Run pushes, and this Spec adds
  no key variable, no Run child and no evidence path. ADR-0383, ADR-0390 and
  ADR-0396, which cite ADR-0319, do not apply: they decide the reply,
  precedence and notice rules of a Cross-project Message, and this Spec changes
  neither those rules nor the follow-up order. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answer of 2026-11-14 "Mesmas do BE (Recommended)", quoted in the
  authorization record with the cycle BE terms it names, which keep the grants
  of the Baseline source with its derived files, the guides under
  `docs/agents/` and the owned skills; the Governed Paths are the
  spec-workflow module, the Standard TypeScript profile and its formatter
  golden, this repository's managed Spec routing guide and setup manifest, and
  the files of the four skills this Spec changes. No Makefile, workflow, Go
  module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0393-a-feature-is-judged-for-the-idea-stage-before-its-prd/_authorization.md`;
  bounded files: `.agents/skills/brainstorming/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-idea/references/idea-template.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/brainstorming/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`.

## Goals

- Before the PRD of a feature is written, the feature has an Idea-stage
  verdict with its three answers, from the judge or from the agent, and
  whoever must act on it saw the answers and not a bare recommendation.
- An indicated Idea stage starts without being asked, and dismissing it costs
  the user one answer, recorded with its reason.
- A request that exists only in the session is judged like a recorded one,
  and a repository without a judge key or a session inside a Run follows the
  same rule through the agent's own answers to the same criteria.
- Every feature Spec carries a record of the judgment and of what followed
  that a later report can read from Git.
- A request written in Portuguese is judged as it is.
- Nothing is refused: the judgment changes no other command's output or exit
  code.

## User Stories

1. As the person who supervises Spec authoring, I want each feature request
   judged for the Idea stage before its PRD is written, so that a wrong shape
   is caught before Tasks are cut.
2. As that person, I want an indicated Idea stage to start without my asking
   and to end with one recorded answer when I dismiss it, so that the study
   happens when it is needed and costs little when it is not.
3. As the maintainer of a repository without a judge key, I want the agent to
   answer the same criteria and say that the judge did not answer, so that the
   rule works for me too and I can tell whose answers I am reading.
4. As the maintainer who will read the calibration report, I want every
   feature Spec to record who judged, the answers, whether the stage ran and
   what it changed, so that the thresholds can be measured on my own history.
5. As a maintainer whose requests arrive in Portuguese, I want them judged as
   they are written, so that the judgment does not skip my requests.
6. As a maintainer who is away while a cycle runs, I want a feature whose
   judgment needs me to wait for me, so that no agent answers the Idea stage's
   questions in my place.
7. As a maintainer who asks for a feature in the session, I want that request
   judged like a recorded one, so that the features I type are not left to the
   agent's judgment alone.

## Core Features

1. The Idea-stage judgment. `roundfix spec judge --stage idea --item
   <reference>` judges one item that is not yet a Spec, named by the reference
   the Follow-up Command prints for it: an open Backlog Entry, a note in the
   repository inbox, or an entry of a messages file handed in with
   `--messages`. With `--request <path>` it judges instead a request of the
   session that the agent wrote to a file. It asks the judge three yes or no criteria in one request:
   the right shape of the solution is still open; the request rests on an
   assumption about its users, the business or the market that it gives no
   evidence for; deciding to build it depends on a market or competitive
   question. It prints the verdict, each criterion's answer with its
   probability, the rule and the thresholds, the lines to copy into the
   record, and the judge's cost and model, in text or JSON.
2. The verdict. An answer is yes at a probability of 0.70 or more, no at 0.30
   or less and unsure between. The stage is indicated when any criterion is
   yes, not indicated when every criterion is no and undecided otherwise.
   Answers count only from the pinned judge model. The thresholds are the
   judge documentation's starting points: the measurement of 2026-11-14 was
   too small to fix them, so the undecided band stays wide and goes to the
   user.
3. What the judge is sent. The project's name and the one item's kind, type,
   origin, title and the first 4,000 characters of its text, the state the
   follow-up judge already sends, in whatever language the item is written.
   The item passes the secret scan first, and an item that matches a secret
   pattern is not sent. A request of the session is sent only from the file
   `--request` names, and of that file only the first line, as the title, and
   the first 4,000 characters, with the project's name, after the same secret
   scan; the command never keeps the file. Nothing else can be sent: a Spec
   artifact, source code and any other file never reach the judge. A Finding, an item
   typed `fix`, `perf` or `refactor`, and a Waiting Item are not sent either
   and are reported as not judged. The recipients, the keys, the Judge Log,
   the monthly ceiling and the Run Marker are the judge's existing ones.
4. When the judge does not answer. Inside a Run, without a key, at the monthly
   ceiling, with an unreadable Judge Log, on a service failure, for an item
   that matches a secret pattern, and when another model answered, the command
   says why in its first line, prints the three criteria with their yes and no
   meanings for the agent to answer, and exits zero. An unknown reference, an
   active Spec, an unreadable or empty request file and a usage error exit 2. No other command reads the result or
   changes its exit code because of it.
5. The routing rule. The three route clauses of the Baseline's Spec routing
   rule keep their ids and force and say: before the PRD of a feature that has
   no Idea and that the user did not send to `write-idea`, the three criteria
   are judged, by the command, for a recorded item or for a request of the
   session written to a file, and by the agent, saying so, when the judge did
   not answer; indicated
   starts `write-idea` without asking and shows the verdict, each answer with
   its probability and who judged, and the user may dismiss the stage at its
   first question; undecided, or any doubt, is the user's decision, seeing the
   same answers; not indicated continues to `write-prd`; the judgment informs
   and never refuses; a fix or a refactor is never judged; and the smaller
   sufficient route no longer decides whether a feature starts at
   `write-idea`.
6. The agent's own answers. When the judge did not answer, the agent answers
   the same three criteria as yes, no or unsure,
   says that the judge did not answer, and follows the same rule. Its answers
   are recorded as the agent's and are never presented as the judge's.
7. The record. Every feature PRD carries an `## Idea-stage judgment` section in
   a fielded form: the judged item, who judged (the judge with its model, the
   agent with the reason the judge did not answer, or nobody because the user
   asked for the stage), each criterion's answer with its probability when the
   judge answered, the verdict, whether the stage ran, was dismissed (with the
   user's reason) or was not run, and what the Idea did to the shape of the
   request: changed, split or unchanged. An Idea carries the same section,
   with `shelved` when the stage ended the request. `write-idea` writes it,
   `write-prd` copies it or writes it when no Idea ran, and the command prints
   the lines it knows.
8. The skills. `brainstorming` routes a feature through the judgment, and
   writes a request of the session to a file in a temporary directory outside
   the repository, which it removes after the command ran.
   `write-prd` judges before any other step when no judgment exists, starts
   `write-idea` for an indicated stage, asks the user for an undecided one and
   records the outcome. `write-idea`, when the rule starts it, opens its first
   question with the verdict and the answers and gives that question one last
   option that dismisses the stage. The Roundfix skill's follow-up mode runs
   the judgment for each planned feature, shows it in the plan and asks the
   Idea stage's questions while the user is present.
9. The user away. No agent answers the Idea stage's questions or an undecided
   judgment for the user. When the stage is indicated or undecided and no user
   is present, no PRD is written for that feature: it stays a Backlog Entry
   that records the judgment and waits, with a Waiting Declaration naming the
   Idea stage, and the cycle's report names it.
10. Documentation and vocabulary. The `spec` command guide and the Roundfix
    skill's `spec` reference describe the stage; the glossary adds **Idea
    Stage Judgment** and changes **Advisory Judgment**; ADR-0403 and ADR-0404
    record the decisions.

## User Experience

A maintainer asks for a feature that an open Backlog Entry already describes.
The agent runs the command on that entry and reads:

- indicated, with `shape_open: yes (0.98)`: it starts `write-idea`, and the
  stage's first question shows the three answers, says the judge answered with
  its model, and offers, last, to dismiss the stage. The maintainer answers
  the stage's questions, or dismisses it in one answer; the PRD then records
  `stage: dismissed` with the reason given.
- undecided: the agent asks one question, "does the Idea stage run?", showing
  the three answers, and follows the answer.
- not indicated: the agent writes the PRD, whose record shows the three
  answers and `stage: not run`.

A maintainer types a feature request in the session. The agent writes it to
a temporary file, runs the command with `--request` and removes the file; the
output and the record name the item `the request in the session`.

In a repository without a judge key the command's first line says the judge
was not asked and why, and lists the three criteria with what yes and no mean.
The agent answers them, tells the maintainer that these are its own answers
because the judge did not answer, and follows the same rule; the record says
`judged-by: agent`.

In a follow-up cycle the plan shows each planned feature with its verdict and
answers. The Idea stage's questions are asked with the cycle's other
questions, before any Spec is authored. A feature whose judgment needs the
maintainer and gets no answer stays in the backlog, waiting, and the report
says so.

## Non-Goals / Out of Scope

- A command that refuses a PRD without an Idea, and any check that fails a
  Spec for a missing judgment or record: an advisory judgment never gates.
- Sending the judge a Spec artifact, a Finding or source code, or more of a
  request of the session than its first line and its first 4,000 characters.
- Judging fixes and refactors, whose route stays `write-techspec`.
- An agent answering the Idea stage's questions, or deciding an undecided
  judgment, for the user.
- Judging inside a Run.
- The calibration report, the judgment in the follow-up report, any reader of
  the record, and any automatic change of a threshold: Spec 0395.
- A language gate, a second pinned model, and the three measured questions
  that were dropped.
- Changing the Idea stage itself: its questions, research, scoring and
  template sections, other than the judgment section and the dismissal
  option.
- Changing the grouping rule (Spec 0391), skill installation and the
  generated subagents (Spec 0392), `roundfix follow-up`, the other stages of
  `roundfix spec judge` beyond their usage text, the judge's keys, ceiling,
  log format or recipients, the Baseline's skill dispatch triggers, and any
  Source Baseline row.

## Success Metrics

1. Success Metric: with a judge that answers 0.98, 0.84 and 0.17, the command
   on an open Backlog Entry typed `feat` prints the verdict `indicated`, the
   three answers as `yes (0.98)`, `yes (0.84)` and `no (0.17)`, the rule with
   both thresholds and the record lines, exits 0 and appends one Judge Log
   line for the judgment; with 0.12, 0.17 and 0.07 the verdict is `not
   indicated`; with 0.42, 0.20 and 0.05 it is `undecided`. At `251ed81e` the
   command refuses `--stage idea`.
2. Success Metric: the one request the command sends holds the project's name,
   one item's kind, type, origin, title and text cut at 4,000 characters, and
   the three questions, and nothing else. Nothing is sent for an item that
   matches a secret pattern, for a Waiting Item, for an item typed `fix`,
   `perf`, `refactor` or `finding`, without a key, at the ceiling, or with a
   non-empty Run Marker, and each of those exits 0 naming its reason.
3. Success Metric: without a key the command lists the three criteria with
   their yes and no meanings and tells the agent to answer them and to say
   that the judge did not answer; an unknown reference and an active Spec exit
   2 naming the reference.
4. Success Metric: every test of the judge, of `roundfix spec judge` and of
   `roundfix follow-up` that exists at the delivery base passes unchanged,
   except the one that pins the unknown-stage message, which now names `idea`.
5. Success Metric: a managed refresh of an adopted temporary repository writes
   the three reworded route clauses into its Spec routing guide, each new
   sentence once; a second refresh reports it current; the clause ids and
   their force are unchanged, and the Surface Transcript sentence of the third
   clause is still there once.
6. Success Metric: both copies of `brainstorming`, `write-idea`, `write-prd`
   and the Roundfix skill's `modes` and `spec` references carry the rule, the
   record form and the waiting rule, and both templates carry the record
   section.
7. Success Metric: `CONTEXT.md` defines **Idea Stage Judgment** and the
   changed **Advisory Judgment**; the `spec` command guide documents the
   stage; the documentation, skill and glossary contracts pass.
8. Success Metric: a judgment costs less than US$0.01. The committed
   measurement recorded US$0.000054 per request with six questions.
9. Success Metric: for a request file whose first line is longer than 200
   characters and whose text is longer than 4,000, the one request the command
   sends holds the project's name, the kind and type `request`, a title of 200
   characters and a text of 4,000, and neither what follows them nor the
   file's path; the output and the record name `the request in the session`.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write.
The first was read through Exa for the Idea on 2026-11-14 and the others for
this PRD on the same day.

- Kohavi, Deng, Frasca, Walker, Xu and Pohlmann, "Online Controlled
  Experiments at Large Scale", KDD 2013
  (https://dl.acm.org/doi/pdf/10.1145/2487575.2488217): "Only one third of
  the ideas tested at Microsoft improved the metric(s) they were designed to
  improve." Ideas fail often enough that asking whether to build one, before
  it is built, has something to find; it does not show that an earlier study
  would have caught them.
- Herbei and Wegkamp, "Classification with reject option", 2005
  (https://webmath.univ-rennes1.fr/stat/GPE/COURBES/articles/hw05.pdf): "we
  consider classifiers f(X) that render three possible outputs: 0, 1 and R.
  The option R expresses doubt and is to be used for few observations that
  are hard to classify in an automatic way." This is the three-way verdict:
  two cut points on a probability and a middle band that goes to a person.
- Hanley and McNeil, "The meaning and use of the area under a receiver
  operating characteristic (ROC) curve", Radiology 1982
  (https://pubs.rsna.org/doi/10.1148/radiology.143.1.7063747): the area
  "represents the probability that a randomly chosen diseased subject is
  (correctly) rated or ranked with greater suspicion than a randomly chosen
  non-diseased subject", and its standard error is used to "guide in
  determining the size of the sample required to provide a sufficiently
  reliable estimate of this area". This is why 15 positive and 23 negative
  cases rank the questions and fix no threshold.
- TypeSafe documentation, "Noul"
  (https://docs.typesafe.ai/primitives/noul.md): "A good starting point is 0.7
  for yes and 0.3 for no." and "Avoid fitting thresholds closely to a handful
  of test cases. Values between 0.3 and 0.7 can go to a person rather than
  either code path."

## Unreachable Acceptance

- criterion: after the release that carries this Spec, `roundfix spec judge
  --stage idea --item <reference>` run with a real judge key on a real open
  Backlog Entry prints a verdict with three probabilities from the pinned
  model and appends its line to the real Judge Log.
  reason: it needs the released binary, the network, a judge key and the real
  Roundfix Home, which no Verification or QA row may reach.
  satisfied-by: the operator's session evidence after the release, recorded
  session evidence.
- criterion: after two months of use, feature Specs need 30% fewer corrections
  than the 20 feature Specs delivered before the release, at least 33% of the
  indicated features had their shape changed, split or shelved by the Idea,
  fewer than 20% of indicated stages were dismissed, every feature Spec
  authored after the release carries the record, and fewer than 10% of the
  features judged not indicated later changed scope.
  reason: it needs two months of authoring in adopting repositories and the
  report that reads the records, which Spec 0395 delivers.
  satisfied-by: the calibration report of Spec 0395, recorded session
  evidence.

## Validation

The decisions were checked on 2026-11-14 against these sources.

- Characterization at `251ed81e`: `roundfix spec judge` knows the stages
  `prd`, `techspec` and `tasks`, each needs an active Spec with a PRD, and any
  other stage is a usage error. The follow-up judge is the only judgment that
  reads an item that is not yet a Spec; it builds its request from the
  follow-up inventory, scans it for secrets and sends items whatever their
  language. No command reads an Idea or its status. The first route clause
  reads "Large or fuzzy product initiative", the second "Standard feature that
  changes product behavior" and the third "start with the smaller sufficient
  route when two routes fit"; `write-prd` treats an existing Idea as ground
  truth and never sends a request to `write-idea`.
- The measurement in `docs/references/idea-stage-judgment-measurement.md`,
  with `jev-1.13.0`: 38 Backlog Entries typed `feat`, labeled by what happened
  to each, plus controls and Portuguese pairs, 270 requests, US$0.0146. It
  kept the three criteria the maintainer named, dropped three other
  questions, found Portuguese judged as English is, found that fixes and
  Findings read as indicated too often to be judged, and was too small to fix
  a threshold.
- A prototype of the Baseline, skill and glossary changes on a scratch
  worktree at `251ed81e`: the clause edits, the module record command, `make
  baseline-digests` and the managed refresh rewrote the module version record,
  the Standard TypeScript profile and its formatter golden of the Spec routing
  guide, three catalog files, four plan goldens and this repository's Spec
  routing guide and setup manifest; the Baseline package and the CLI's
  Baseline tests passed; the skill edits, `make skills-sync` and the version
  record passed the skills package; the glossary edits passed the glossary
  contracts, and the skill edits did not move the Behavior Surface Record.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for a judgment
  that decides whether discovery runs and for measured Jev thresholds.
  `wiki/entities/jev-typesafe.md` records "manter gates humanos até calibrar
  com dados próprios" and that English works best, which is why the undecided
  band goes to the user and why Portuguese was measured instead of assumed.
  `raw/web/akita/2026-10-10-akita-qual-modelo-de-decisao-classificador-de-ia-tipo-jev-e-melhor.md`,
  an independent benchmark with half its cases in Portuguese, found that Jev
  "em português repetiu 97,6% no agregado sem perder nada na tradução" and
  warns that "os intervalos importam"; this Spec's own pairs agree.
  `wiki/sources/branas-formacao-ia-processo-de-desenvolvimento.md`, read for
  the Idea, shows a process with no such stage.
- Exa: the four published sources above, and Johnson and Goldstein, "Do
  Defaults Save Lives?", Science 2003 (https://papers.ssrn.com/abstract=1324774),
  whose abstract notes that "defaults can lead to two kinds of
  misclassification"; it is why the record keeps what the Idea did to the
  shape and not only whether the stage ran.

## Risks

- The thresholds are starting points. On the measured history 12 of 38
  feature items read indicated, four of them delivered as asked, and 18 read
  undecided, so the user is asked something for most features until the
  thresholds are measured on more cases. The record exists for that.
- A verdict near a cut point can change between two runs of the same model.
  The output always shows the probabilities, and the rule treats the verdict
  as advice.
- A "not indicated" verdict may be accepted without thought. It is never
  shown bare: the answers and their probabilities are printed and recorded.
- A default is rarely overridden, so few dismissals would not prove the stage
  was needed. The record of what the Idea changed is the evidence.
- The Idea stage needs the user. A cycle with the user away authors only the
  features judged not indicated; the others wait.
- A request of the session was not measured: the thresholds come from Backlog
  Entries, and the user's own words are often short and in Portuguese. The
  record names the judged input, so the later report can compare them.
- The judge now receives text that has no record in the repository. Only its
  start is sent, after the secret scan, and the Judge Log keeps a hash of the
  state, never the text.
- Dissent kept from the Idea: the operator would have shipped the record
  without a later report, because a report over a handful of cases measures
  noise; and the rule alone, with the agent's answers and no judge, reaches
  every adopter with no data decision. On 2026-11-14 the evidence that changed
  Spec 0392 came from the Idea stage's research, not from the judge, which was
  undecided.

## Idea-stage judgment

- item: the request in the session, recorded in `references/2026-11-14-maintainer-request-and-answers.md`
- judged-by: judge (jev-1.14.0, in the operator's pilot)
- shape_open: no (0.23)
- unverified_assumption: unsure (0.61)
- market_question: unsure (0.46)
- verdict: undecided
- stage: ran — the maintainer chose it on 2026-11-14
- shape: split — into this Spec and Spec 0395

## Glossary

- adds: **Idea Stage Judgment**
- changes: **Advisory Judgment**

## Decisions

Autonomous route: no person present; decisions taken from the approved Idea, the maintainer's answers of 2026-11-14 quoted below, the operator's cycle BF brief for Spec 0393, the measurement and the characterization above.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-14), each with its question in the adopted record:

- 2026-11-14: Onde esse julgamento deve morar? — "Spec nova no produto (Recommended)", option text: "Feature: antes do PRD de uma feature, o Roundfix pergunta ao Jev se ela precisa da etapa de Ideia (critérios da própria write-idea: forma da solução em aberto, suposição de negócio não conferida, questão de mercado) e a regra de roteamento passa a mandar rodar write-idea quando o juiz indicar. Vale para todos os adotantes."
- 2026-11-14: Qual autorização vale para ela? — "Mesmas do BE (Recommended)", option text: "Push, PR e merge; release minor depois dela; até 3 correções; Makefile, CI, go.mod e lint só quando uma correção exigir, descrito no relatório."
- 2026-11-14: Qual desperdício você mais quer evitar com o julgamento automático? — "Feature na forma errada (Recommended)".
- 2026-11-14: Quando o Jev indicar que uma feature precisa da etapa de Ideia, o que acontece? — "Roda, com saída registrada (Recommended)", option text: "O agente inicia a write-idea sem pedir licença e mostra os números do julgamento. Você pode dispensar a etapa na primeira pergunta dela, e a dispensa fica escrita no PRD com o motivo."
- 2026-11-14: O que a regra manda o agente fazer quando o juiz não responde? — "O agente aplica os critérios (Recommended)", option text: "O agente responde ele mesmo às mesmas perguntas (forma em aberto, suposição sem evidência, questão de mercado), diz que o juiz não estava disponível e segue o resultado; na dúvida, pergunta a você. Todo adotante ganha a regra, com ou sem chave."
- 2026-11-14: Qual medida faria você dizer que o julgamento valeu? — "Menos correções depois".
- 2026-11-14: Qual forma a ideia toma? — "Mais ambiciosa", option text: "A original com registro, mais um relatório que remede os limiares com o histórico do próprio repositório e o conselho por item já no relatório do `roundfix follow-up`. Duas Specs."
- 2026-11-14: Rascunho da Ideia 0393: aprova? — "Aprovado (Recommended)".
- 2026-11-14: Spec 0393, fronteira de dados do juiz. Hoje o Jev só recebe itens que já estão no repositório (Backlog Entry, nota do inbox ou mensagem de outro projeto), depois da varredura de segredos. Um pedido que você digita na sessão, como o seu pedido desta manhã que virou a 0393, não tem arquivo, e nesse caso o agente responde aos três critérios sozinho. Na medição com 38 features do histórico, o juiz acertou 8 de 12 "indicada" e deixou 18 de 38 como "indecisa". O juiz pode receber o pedido que você digita na sessão? — "Pode, com varredura (Recommended)", option text: "O agente grava o pedido num arquivo entregue como uma mensagem; vai ao juiz só o nome do projeto, a primeira linha como título e os primeiros 4.000 caracteres, depois da mesma varredura de segredos. Os pedidos que nascem na sessão passam a ser julgados como os outros." (AskUserQuestion, about 10:05 BRT, relayed by the operator.)

Settled from the Idea's open questions, with no further interview:

- The questions are the three criteria the maintainer named, worded as
  measured. The thresholds stay at the documentation's starting points, 0.7
  and 0.3, because the measured set cannot fix them; the wide undecided band
  goes to the user. The model is the judge's pinned one.
- Before a PRD exists the judge is sent an item the follow-up inventory holds
  or, by the maintainer's last answer above, a request of the session handed
  in as a file, of which only the first line and the first 4,000 characters
  are sent. The file is plain text: a request has no sender, id or date.
- A request in Portuguese is judged; the measurement found no reason for a
  language gate.
- The record lives in the PRD of every feature Spec, and in the Idea when one
  is written.
- A feature is work on either feature route of the routing rule. An item typed
  `fix`, `perf` or `refactor` and a Finding are not judged.
- With the user away, a feature whose stage is indicated or undecided waits as
  a Backlog Entry with a Waiting Declaration.
- Which records hold the corrections a Spec needed is Spec 0395's to settle.
- The judgment is a stage of `roundfix spec judge`, the command authors
  already run for advisory judgments, with the item named by a flag.
- `type: feat`, Release Ceiling `minor`. See ADR-0403 and ADR-0404.

## Open Questions

None.
