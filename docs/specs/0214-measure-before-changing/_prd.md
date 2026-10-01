---
spec: 0214-measure-before-changing
status: active
created: 2026-10-01
surfaces: [backend, cli, docs]
---

# Measure before changing

Three open questions ask Roundfix to change something on a belief, and each
can be answered from Roundfix's own records first:

- **Why corrective Tasks happen.** The maintainer declined a memory layer for
  Agent Sessions on 2026-09-30 and chose committed guides for lessons that
  cross Tasks. Whether to reopen that question depends on how many failed
  Verification attempts and corrective Tasks come from missing repository
  knowledge rather than from plain implementation error. The Run Database
  records each Verification attempt and its verdict, but nothing classifies
  its cause.
- **Whether a Task acceptance judgment is worth adopting.** The 2026-09-30
  Jev measurement asked whether a Task's acceptance criteria are all
  observable. The signal ran in the expected direction but was weak (AUROC
  0.62 over 300 Tasks), and the label was noisy, because it was read from
  whether the Agent-written result mentioned Verification Feedback. The Run
  Database now records Verification attempts per Task, which is the cleaner
  label the measurement asked for.
- **Whether archived evidence can go.** The history root holds most of the
  repository's tracked files: 4,682 of 6,254 and 44 MB on 2026-10-01, of
  which 1,627 QA evidence files take 11 MB and four design and reference
  binaries take 7 MB. Every repository-wide search and every Secondbrain
  mirror sync reads them. Nobody has shown which of them any agent, test,
  check or tool still reads.

This Spec answers the three questions with measurements and written
recommendations. It adds one read-only command, `roundfix runs causes`, that
classifies why Verification attempts failed and why corrective Tasks were
added, using deterministic signatures; a test harness that re-scores the Task
acceptance question against the Run Database's attempt count; and a written
inventory of who reads archived evidence, with a proposal. It changes no gate,
no Verification, no judgment and no archived file. A change that a
measurement supports becomes a Backlog Entry for a later Spec.

The maintainer decided on 2026-10-01 that the archive work measures and
proposes only: this Spec delivers the measurement and a written proposal, it
deletes no archived file, and any removal waits for the maintainer's explicit
approval in a later change.

## Prerequisites

This Spec is delivered after Spec 0205, which adds `roundfix spec judge` and
the judge's transports, keys, model pin, Judge Log and monthly ceiling. Task
02 re-uses that code to ask its question and its Verification compiles against
it, so the operator queues Spec 0205 first. Spec 0209 extends the same judge
and ships with Spec 0205; this Spec re-uses whatever the two deliver and
changes neither.

## Project Constraints

- Identifier strategy: not applicable — no new persistent identifier. A cause
  class, a signature and a trigger are named by fixed kebab-case words in an
  embedded table, a record by its schema name, and a corrective Task by the
  Task identifier its Task Graph already gives it. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — only the Task acceptance harness
  sends requests, and only when its own flag is passed. It sends them through
  the judge's two transports with the judge's keys: OpenRouter's System One
  API with `ROUNDFIX_OPENROUTER_API_KEY` (model `jev-1.13`) first, and
  TypeSafe directly with `ROUNDFIX_TYPESAFE_API_KEY` (model `jev-1.13.0`)
  only when the first key is absent. Each key is read only from the
  environment, sent only in the authorization header of its own endpoint and
  never printed, logged or stored; the generic `OPENROUTER_API_KEY` and
  `TYPESAFE_API_KEY` are never read. A request carries only this
  repository's Task files. The `runs causes` command and every test,
  Verification and QA row send nothing. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0214 (this Spec) decides that a
  cause is classified by deterministic signatures from a closed list, that a
  corrective Task is found by its place after the QA Task, and that the report
  is read-only. ADR-0215 (this Spec) decides that a measurement lands as a
  record with a decision rule fixed in advance, changes no gate, removes
  nothing archived, and that a re-measured judgment keeps the judge's
  boundary. ADR-0201 bounds what the judge may send, and its consequences
  already place Task files inside the authorization, ADR-0201: "Findings,
  Backlog Entries and Task files are within the authorization but not read";
  the harness reads them under ADR-0215. ADR-0200 pins Jev 1.13 and compares
  an answer with a threshold only for that version, and the harness keeps the
  pin. ADR-0004 places every Run in one Run Database, ADR-0004: "Roundfix
  stores Run state in a global SQLite database". ADR-0014 and ADR-0038 make
  the Daemon run Verification and allow one repair, ADR-0038: "allows one
  repair, and settles the Work Item from the second verdict", which is why a
  Task's failed verdicts are a label. ADR-0033 prunes the journal of old
  terminal Runs, which bounds what can be classified. ADR-0008 keeps a Run
  Event's payload raw, ADR-0008: "Daemon-owned events define their own small
  JSON payloads", and the command reads those payloads without rewriting them.
  ADR-0120 places retired documentation under one history root, ADR-0120: "The
  root is `docs/history/`", which is the tree the archive measurement
  inventories; ADR-0121 records an archive move as a ledger of identities;
  neither changes. ADR-0035 lets the Spec Root live outside the repository,
  and the command reads Task Graphs where the Spec Root and its archive
  resolve. ADR-0089 passes the environment to code under test explicitly,
  which is how every test supplies a home and the harness's flag. ADR-0184 has
  the TechSpec state the new command as Surface Transcripts, ADR-0184: "A
  TechSpec now declares numbered Surface Transcripts". ADR-0187 splits the
  Roundfix Skill and the command reference by command, and the new command is
  described in the `runs` files. ADR-0189 ties an owned skill's version to its
  content, so the edited skill raises its version. ADR-0193 names a Spec's
  prerequisites; this Spec states Spec 0205 as one in its PRD and Build Order.
  ADR-0081 and ADR-0149 sanction the regeneration that follows the skill edit.
  ADR-0182 runs Settlement Checks before each Task commit, ADR-0178 authorizes
  a Task commit by the grant it ran under, ADR-0166 records undeclared paths,
  and ADR-0130 judges Governed Paths by kind; every Task here declares its
  paths. ADR-0167 keeps the pre-PR Pull Request row from deciding a qualifying
  partial, ADR-0167: "it runs before any Pull Request exists"; this gate aims
  at `pass`. This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and
  ADR-0094 check its consistency by citation and artifact presence, with
  ADR-0168, ADR-0176 and ADR-0183 narrowing which records and receipts that
  check reads. ADR-0159 lets an attempt run past a failed command, so a failed
  attempt may name several commands and the report takes the first. ADR-0171
  makes a retention prune report only what it reclaimed, ADR-0171: "reports
  only what it reclaimed", and ADR-0098 appends Run Events in batches; the
  report reads whatever survives either. ADR-0209 extends the same judge with
  a grouping question, which the harness neither asks nor changes. The
  remaining records that cite a listed ADR decide something this Spec does not
  touch: ADR-0020 (a parsed prompt result outranks the runtime's exit code);
  ADR-0022 (stop requests); ADR-0029 (review artifacts live with the Spec);
  ADR-0030 (opt-in agent run logs); ADR-0036 (review artifacts committed in a
  separate docs commit); ADR-0056 (Task and Verification capacity); ADR-0057
  (the Daemon owns Task status); ADR-0097 (QA row carry-forward); ADR-0127
  (process residue); ADR-0142 (head-bound Review Source Evidence decides the
  watch outcome); ADR-0160 (the frozen authorization that opens a red gate);
  ADR-0165 (a blocking review after archive); ADR-0169 (the pre-PR review
  diffs the candidate from its merge base); ADR-0170 (a Task already completed
  on its target is nothing to carry); ADR-0172 (doctor reports reclaimable Run
  storage from cheap reads); ADR-0179 (an empty paths list); ADR-0192
  (conflicts on derived paths); ADR-0194 (the Daemon records what a QA row
  observed and hands a failed pass to the next); ADR-0195 (rows that read the
  gate or the commits are observed on every pass); ADR-0208 (work items that
  share a context share a Spec, and an open item is extended); ADR-0210 (an
  evidence snapshot records one digest per declared input); ADR-0196 (a pre-PR
  review finding parks only after validation); none of them applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário"), extended to this Spec by the maintainer's 2026-10-01 decision
  to deliver the whole A to F program ("Tudo, de A a F"), recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runs.md`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A maintainer can read, for any period the Run Database still holds, why
  each failed Verification attempt and each corrective Task happened, as one
  class from a closed list or as `unclassified`, with the signature that
  decided it.
- The memory question gets a numbered answer: the share of failures that
  come from repository knowledge, and a recommendation computed by a rule
  stated before the data were read.
- The Task acceptance judgment is re-scored against the Run Database's attempt
  count on the pinned model, with a random-Task baseline and a confidence
  interval, and gets an adopt, do-not-adopt or inconclusive verdict by a rule
  stated in advance.
- The maintainer gets a measured inventory of archived evidence (what it
  holds, who reads it, what a tree without it would break) and a written
  proposal, while every archived file stays in place.
- Every number is recorded in the repository in a form a test can recompute,
  and a change any number supports is written down as a Backlog Entry, not
  made.

## User Stories

1. As a maintainer deciding whether to reopen a memory layer, I want the
   failed Verification attempts and corrective Tasks of a period counted by
   cause, so that I decide on Roundfix's own data instead of published claims
   that point both ways.
2. As a maintainer auditing a cause count, I want each item to name the
   signature and the evidence that classified it, so that I can dispute a
   wrong class without re-running anything.
3. As a maintainer weighing a third advisory judgment, I want the Task
   acceptance question scored against a clean repair label with an interval
   and a baseline, so that I know whether it can flag one Task at a time.
4. As a maintainer considering pruning the history root, I want to know which
   archived evidence any agent, test, check or tool reads and what removing
   the rest would break and save, so that I can approve or refuse a removal
   knowing its cost.
5. As a maintainer, I want nothing archived removed and no gate changed by
   this work, so that a measurement never becomes a change I did not approve.

## Core Features

1. **`roundfix runs causes` reports why.** For the current repository's
   terminal Spec Runs created in an optional window (`--since`, `--until`),
   it lists each failed Verification attempt and each corrective Task with its
   failing check, its cause class and the signature that decided it, then a
   summary by class, as text or `--format json`. It reads the Run Database
   with the reader that never migrates, and writes nothing.
2. **A closed list of causes.** Each item has exactly one class:
   `scope-or-authorization`, `shared-section-contract`,
   `repository-convention`, `implementation-defect` or `environment`, or it
   is `unclassified`. The first three are counted together as repository
   knowledge. An unmatched item is never forced into a class.
3. **Deterministic signatures, one reviewable table.** An ordered table,
   embedded in the binary and reproduced byte for byte in the TechSpec, binds
   each signature to one class, one evidence source and one regular
   expression. The first match decides. Every report names the table's
   digest. No model is asked.
4. **Corrective Tasks found by position.** A corrective Task is a Task that
   is not the QA Task and whose number follows the QA Task's in its Task
   Graph, read from the active Spec Root or its archive. Its failing check is
   the trigger its overview names (pre-PR review, QA gate or Verification),
   or `unknown`.
5. **Per-Task attempt counts.** The JSON report also lists every Task that
   ran in the window with its passed and failed Verification verdicts and its
   Verification Feedback rounds, and whether it is the QA Task or a corrective
   Task. That list is the label the Task acceptance measurement reads.
6. **The Task acceptance re-measurement.** A test harness, run only when its
   flag is passed, asks Jev 1.13 one fixed question about each implementation
   Task's acceptance criteria, through the judge's transports, keys, Judge Log
   and monthly ceiling, and writes a record of every answer with the Task's
   repair label. Without its flag it sends nothing. A Task's Agent-written
   result is never part of a request.
7. **Statistics a test recomputes.** The record carries AUROC with a 95%
   interval from a bootstrap that resamples whole Specs, the AUROC of a seeded
   random score and of the criteria's length, and precision and recall at a
   fixed threshold. A test recomputes every figure and the verdict from the
   record's own answers.
8. **Rules stated in advance.** The memory recommendation, the Task
   acceptance verdict and the archive candidate set each follow a rule
   written in the TechSpec before any data are read. The recommendation states
   which rule fired and on which numbers.
9. **An inventory of archived evidence.** A reference document records, at a
   named commit, the history root's files and bytes by kind; every
   repository file outside history that reads it; what `make verify` and
   `make verify-docs` do in a disposable copy without QA evidence and
   binaries; how often Agents read archived evidence in the Run Database; and
   what the Secondbrain mirror holds of it.
10. **A proposal, not a removal.** The same document proposes what could be
    removed, what it would save, what must change first, and what it would
    not reclaim. No file under the history root is deleted, moved or
    rewritten. The removal is a later change that the maintainer approves
    explicitly.
11. **Follow-ups as Backlog Entries.** When a rule says to act, the
    measurement Task writes an open Backlog Entry for the later Spec, and
    writes none when it does not.

## User Experience

The maintainer runs `roundfix runs causes --since <date>` in the repository
and reads one line per item, then a summary line that gives the count per
class and the repository-knowledge share. `--format json` gives the full
record, including the per-Task attempt counts. The three reference documents
read top-down: what was measured, at which commit or window, the numbers, the
rule and the recommendation.

## Declared breaks

- `roundfix runs --help` and the top-level help list one more command.
- The Roundfix Skill's version rises by one patch step.

## Non-Goals / Out of Scope

- Deleting, moving or rewriting any file under the history root, or any
  other archived file.
- Changing any gate, Verification, Settlement Check, QA rule, judgment
  threshold or skill rule because of a measurement.
- Adding the Task acceptance question to `roundfix spec judge`.
- Building or reopening a memory layer for Agent Sessions.
- A model judgment of a failure's cause.
- Changing Journal Retention or the GC Command so that more evidence
  survives.
- Rewriting Git history, or shrinking a clone of the repository.
- Classifying review-only Runs, or Runs of another repository.
- Any network call from a test, a Verification command or the QA gate.

## Success Metrics

1. Success Metric: on a fixture Run Database with one failed Verification
   attempt per class, one unmatched attempt and one corrective Task, `runs
   causes` reports each with its expected class and signature, counts the
   unmatched one as `unclassified`, and leaves the Run Database and Roundfix
   Home byte-identical.
2. Success Metric: the committed cause record's summary equals the counts of
   its own items, its table digest equals the embedded table's, and its
   recommendation equals what the rule gives for those counts.
3. Success Metric: the committed Task acceptance record's AUROC, interval,
   baselines, precision, recall and verdict equal what a test recomputes from
   its answers, every answered Task's model normalizes to Jev 1.13, and no
   request state carries text from a Task's result section.
4. Success Metric: the archive document's file and byte counts equal what Git
   reports for the named commit, and from the merge base to the delivered
   head no file under the history root is deleted or renamed.
5. Success Metric: a Backlog Entry exists for each rule that fired and for
   none that did not.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence rows rest on sources this Spec did not produce:

- **A record this Spec did not design.** The Run Database on the
  maintainer's machine holds every Verification attempt of this repository's
  Spec Runs since 2026-09-17: on 2026-10-01, 2,618 Verification events, 85
  failed verdicts and 42 Verification Feedback rounds over 357 Task
  executions of 63 Specs, with all 86 failed-attempt diagnostics still on
  disk. The measurement Task reads it read-only, and the QA gate re-runs the
  command over the recorded window and compares.
- **Published failure taxonomies.** A study of 1,184 failed terminal-agent
  trajectories classified decisive errors into a closed list of nine types in
  three families, epistemic 57.9%, competence 32.8% and environment 9.4%, with
  0.6% left as other (<https://arxiv.org/pdf/2607.09510>). A study of 20,574
  real coding-agent sessions found that a single-stage model extraction of
  failure causes produced systematic false positives and needed a second
  evidence filter, and still left 26.85% of causes as "cannot determine"
  (<https://arxiv.org/html/2605.29442v2>). The first supports a closed list
  with a counted remainder; the second supports deterministic signatures over
  a model judgment. Read 2026-10-01.
- **Published evidence on repository context.** The Secondbrain records,
  from its 2026-09-30 reading, that context files lowered agents' success
  rate and raised cost by more than 20% (arXiv 2602.11988), and that a
  288-run ablation found agents failed on implementation skill, not missing
  repository knowledge (arXiv 2607.27250). Both point away from a memory
  layer, which is why Roundfix's own count decides.
- **Published statistics.** A comparison of 29 interval methods for AUC at
  small sample sizes found the discrepancy between methods grows as the
  sample shrinks (<https://pubmed.ncbi.nlm.nih.gov/26323286/>), and a study of
  correlated diagnostic data found that a bootstrap ignoring clustering
  under-covers and that a cluster bootstrap holds its coverage
  (<https://pmc.ncbi.nlm.nih.gov/articles/PMC10728486/>). Tasks of one Spec
  share an author and a design, so the interval resamples whole Specs. Read
  2026-10-01.
- **Git's own documentation.** The Git book states that a large file stays
  reachable in history and every clone downloads it "even if it was removed
  from the project in the very next commit"
  (<https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery>,
  read 2026-10-01), so the proposal states what removal would not reclaim.
- **The prior measurement.** The adopted Backlog Entry's finding, the
  2026-09-30 Jev measurement, reports Task lint at AUROC 0.62 [0.55, 0.70]
  over 300 Tasks with a 23% repaired base rate, 0.67 for backend Tasks only.
  Its scripts and question text were not committed, so this Spec writes the
  question anew and compares only the direction and size of the signal.

## Recorded limits

- The Run Database holds Verification events only from 2026-09-17, and
  Journal Retention may prune older terminal Runs, so the window is short and
  the repaired Tasks are few (34 of 357 Task executions on 2026-10-01).
- Keyword signatures can match for the wrong reason; every item names its
  signature so a reviewer can see it.
- The Task acceptance question is written anew, so its numbers are not the
  same measurement as 2026-09-30's.
- About 2% of Jev answers flip between two runs of the same states
  (ADR-0200's consequences), so a Task near the threshold can change sides.
- The ablation shows what the repository's own gates need; it cannot show
  what a future reader would have wanted from removed evidence.

## Decisions

- **Deterministic signatures, a closed list, read-only.** See ADR-0214.
- **Measure, record, recommend; remove nothing.** See ADR-0215.
- **A command for causes, a harness for Jev.** The cause report is useful
  again after every release and touches no network, so it is a command. The
  Task acceptance re-measurement is a one-off that spends money and sends
  Task files, so it stays a test harness behind a flag and is not compiled
  into the binary.
- **The cause report supplies the label.** The harness reads the cause
  report's per-Task attempt counts instead of opening the Run Database, so
  the Run Database has one reader in this Spec.
- **Clusters, not Tasks, are resampled.** Tasks of one Spec are not
  independent, so the interval resamples Specs.
- **The decision rules are written before the data are read.** A rule fitted
  after seeing the numbers would confirm whatever they show.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the queries `qmd
query "causa de tarefas corretivas e retentativas de verificação em agentes de
código" --all --files --min-score 0.3` and `qmd query "podar docs/history
evidência de QA arquivada tamanho do repositório" --all --files --min-score
0.3`. The first returned this repository's mirrored Backlog Entry and
`wiki/concepts/agent-evaluation.md`, whose AI SDK course notes put
deterministic scorers first and a model judge only where they cannot reach;
that shaped ADR-0214. The triaged inbox note
`inbox/secondbrain/_triaged/2026-09-30-ai-memory-e-graphify-avaliados-contra-o-roundfix.md`
records the published memory evidence above and that the Secondbrain mirror
syncs with deletion, so removing archived files would also remove them from
the mirror; the archive document measures that. The second query returned
archived QA reports of other projects and nothing on pruning. Exa found the
two failure taxonomies, the AUC interval comparisons and the Git book section
cited above. No Inbox Entry for this repository was pending.

## Open Questions

None. The maintainer's decisions of 2026-10-01 settle the data boundary, the
keys, the ceiling and the measure-only scope of the archive work.

## Technical candidate

The [_techspec.md](_techspec.md) records the signature table, the decision
rules, the command transcripts, the measurement procedures, coverage and
build order.
