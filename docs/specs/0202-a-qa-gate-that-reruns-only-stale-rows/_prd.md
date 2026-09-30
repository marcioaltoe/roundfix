---
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: active
created: 2026-09-30
surfaces: [backend, docs]
---

# A QA gate that reruns only stale rows

After a corrective Task, the authored QA gate runs again, and its Agent
Session executes every row of the matrix. Most of those rows read inputs the
correction never touched. ADR-0097 already allows a passing row to carry
forward when its declared repository inputs are unmoved. In practice no row
ever carries, for two reasons.

- **No snapshot.** Carrying needs the evidence snapshot of the report that
  established the row, and nothing writes one. Spec 0179's seven QA passes
  declared repository inputs on every row, and no report recorded a snapshot.
  Across its six later passes, 41 passing rows were executed again. For 24 of
  them, every declared input lay outside the files that changed between the
  two audited heads.
- **No previous report.** A failed gate's Run is never integrated, and Task
  Carry-Forward re-commits its Tasks with new identities. The next pass
  therefore holds neither the failed report nor a head that descends from it.
  Spec 0192's first gate failed at 12:34 on 2026-09-30 after 22 minutes. The
  correction changed only its PRD and Task files. The second pass started from
  a head the first did not precede, found no previous report, and ran all 13
  rows again in about 24 more minutes.

This Spec makes the Daemon record what each passing row observed, hand a
failed pass to the next, and carry every row whose recorded inputs are
byte-identical. Only stale, failed, blocked and always-observed rows are
executed again.

## Project Constraints

- Identifier strategy: not applicable — no new entity identifier. Carried rows
  keep the row identifiers their establishing report gave them, and imported
  reports keep their names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes, Git and the Spec
  tree only; no credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0194 (this Spec) makes the Daemon
  record each passing row's evidence snapshot and hand a failed pass's report
  to the next pass. It lets an establishing head be proven by the QA Report
  commit that recorded it. ADR-0195 (this Spec) keeps rows that read the
  repository Verification, the Pull Request row or the Task commits observed
  on every pass, and records a carry disposition for every prior row. Both
  refine ADR-0097, whose carry rule otherwise holds unchanged: no carried row
  may make a verdict more permissive than a fresh observation would. ADR-0080
  keeps QA verdicts distinguishing environment-blocked rows. ADR-0096 keeps
  the mechanical stage computing facts before the Agent turn. ADR-0053 and ADR-0170 keep Task Carry-Forward
  proof-based and untouched. ADR-0057 keeps Task status Daemon-written.
  ADR-0167 keeps the pre-PR Pull Request row out of a qualifying partial.
  ADR-0093 checks Spec consistency by citation, ADR-0156 makes a declared
  promise name a consuming Task, and ADR-0117 places a check with the stage
  that can produce the defect; all three hold for this Spec's artifacts and
  for the mechanical stage it extends, as do ADR-0168, which narrows the
  related-ADR check, and ADR-0176, which reads citations only from authored
  text. ADR-0183 cites ADR-0093 but requires a receipt on an attributed claim,
  and ADR-0184 cites ADR-0156 but states a command surface as a transcript;
  both bind a Spec authored from their guide commit on, which Spec 0191 has
  not landed, so they do not bind this Spec. ADR-0166 records the paths a Task
  changed without declaring them, and the gate's scope row reads them. This
  Spec's gate is bound by ADR-0088, ADR-0091, ADR-0104, ADR-0155 and ADR-0167.
  ADR-0160 cites ADR-0057 but opens a red repository gate only through the
  frozen authorization. ADR-0161 cites ADR-0053 but releases a merged Spec's
  Runs on the merged head. ADR-0182 cites ADR-0096 but runs Settlement Checks
  when a Task settles. This Spec changes none of these three, so they do not
  apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized adjusting every
  skill on 2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário"). The authorization is recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A QA pass after a correction executes only the rows whose recorded inputs
  moved, the rows that did not pass, and the rows that are always observed.
- A failed pass whose Run was never integrated still serves as the next pass's
  previous report.
- Every prior row's fate is visible in the new report: carried with the report
  and head that established it, or re-run with a named reason.
- No carried row is ever less trustworthy than a fresh observation. The proof
  is byte-identical declared inputs, and it fails closed.

## User Stories

1. As a Supervisor, I want a QA pass after a corrective Task to repeat only
   the rows the correction could have changed, so that a one-file correction
   does not cost a full Agent Session of QA.
2. As a Supervisor recovering a parked Delivery Queue item, I want the next
   pass to inherit the failed pass's report, so that rows it proved are not
   proved again only because its Run was never integrated.
3. As an operator reading a QA Report, I want each prior row's carry
   disposition and reason, so that I can tell inherited evidence from fresh
   evidence and see why a row was executed again.
4. As a QA gate Agent, I want the seeded report to mark which rows are carried
   and which I must execute, so that I spend my session only on stale rows.

## Core Features

1. **The Daemon records evidence snapshots.** When a QA pass closes, and
   before the QA Report commit, the Daemon writes the report's
   `evidence_snapshots` for every row that can carry. Such a row passed,
   declared only repository inputs that all resolve at the audited head, and
   is not always observed. The snapshot names the audited head and each
   declared input's files and SHA-256 digests. Any value the Agent wrote is
   replaced. No other byte of the report changes.
2. **A failed pass reaches the next pass.** Before the next pass builds its
   matrix, the Daemon brings in the newest QA Report commit of the same Spec
   that the Run's history does not contain. It copies that commit's QA files
   byte for byte into the Spec's QA directory. It refuses the import, and says
   why, in three cases: a path already holds different bytes, the report is
   not newer than the tree's newest report, or the mechanical stage would find
   a defect in the report. A refused import leaves no file behind. Imported
   files ride in the next QA Report commit.
3. **The establishing head is proven by ancestry or by its recording
   commit.** A carried row's establishing head must be an ancestor of the
   current head, or the first parent of a Daemon QA Report commit that holds
   the establishing report with identical bytes. Declared inputs are then
   compared by content, and any moved input refuses the carry.
4. **Some rows are always observed.** A row whose provenance names the
   repository Verification or the Pull Request row is never carried. A row
   that declares a `commit_range` input, the kind for rows that read Task
   commits, their authorization or their changed-file scope, is never carried
   either.
5. **Every prior row gets a disposition.** The seeded report lists each row of
   the previous report as `carried`, or as `re-run: <reason>` from a closed
   list of reasons. A carried row keeps the provenance of the row that
   established it. A pass with no previous report writes the seeded report it
   writes today.
6. **The gate executes only what is stale.** The qa-gate skill and the QA
   prompt state that a carried row is kept as is and counts as passed. They
   state that every executed row declares its inputs, and that the Agent never
   writes `evidence_snapshots`. The Context-Driven Development guide describes
   the carry.

## User Experience

A Supervisor following a Run sees one more QA event before the mechanical
stage: `prior_report` with outcome `imported`, `none` or `refused` and its
reason. The mechanical event now counts carried rows and re-run rows. After
the Agent turn, an `evidence_snapshots` event names the audited head and the
number of rows recorded. In the seeded QA Report, carried rows appear under
`## Results` as `carried (established by: <report>; head: <sha>)` with their
original provenance. A `## Row carry-forward` section lists every prior row
with `carried` or `re-run: <reason>`.

## Non-Goals / Out of Scope

- Skipping the Agent Session when every row carries. The always-observed rows
  still need a pass, and settling them mechanically is follow-up work.
- Deriving a row's inputs automatically from the files a test or command
  read. The Agent declares them, as ADR-0097 already requires.
- Changing verdict rules, the typed blocked counts, report naming, the
  `### QA settlement` section, or archive eligibility.
- Changing Task Carry-Forward, the Reconcile Command, Delivery Retry or which
  Runs are integrated.
- Carrying a row across Specs, or from a report whose recording commit is no
  longer in the repository.
- Editing `CONTEXT.md`. New terms are reported for the glossary owner.

## Success Metrics

1. A pass after a correction that moved none of a passing row's declared
   inputs materializes that row as carried in the seeded report. The row names
   its establishing report and head and keeps its provenance.
2. A prior row is recorded `re-run` with its reason, and never carried, when
   any of these holds: it failed or was blocked, it declared no inputs or a
   non-repository input, a declared input moved, or it is always observed.
3. A failed pass committed only on an unintegrated Run Branch, whose Tasks
   were carried forward as new commits, has its report imported byte for byte
   by the next pass. That pass then carries its unmoved passing rows.
4. An `evidence_snapshots` value the Agent wrote is replaced by the Daemon's.
   A row the Daemon did not record never carries.
5. A pass with no previous report writes a seeded report byte-identical to
   today's. The existing carry-forward and report-materialization tests pass
   unedited.

## Prerequisites

Queue this Spec after Spec 0191. Both edit the qa-gate skill, so the second to
land must rebase its skill text and raise the version one step above the
first. This Spec's Verification reads none of Spec 0191's artifacts.

## Recorded limits

- An audit row is recognized only through its `commit_range` declaration. A
  row that reads Task commits but declares only repository paths could carry.
  The mechanical stage's own authorization audit still runs on every pass.
- A declaration that omits an input the row really read can carry a stale
  row. ADR-0097 already rests on declarations this way. The skill tells the
  Agent to declare conservatively, including every source a built binary
  compiles from.
- A carried row whose establishing report is imported needs its recording
  commit in the repository. Once the failed Run's branch is reclaimed, the row
  is re-run.
- Evidence the gate writes under the QA directory does not exist at the
  audited head. A row whose Evidence column names such a file is never
  recorded. The template's report puts evidence links in the row's detail
  block.

## Decisions

- **Content, not ancestry.** A failed pass never lies in the next pass's
  ancestry, so ancestry alone never carries in the Delivery Queue. Comparing
  declared inputs by digest is the proof Task Carry-Forward already uses. See
  ADR-0194.
- **The Daemon writes the snapshot.** The Agent declares inputs, and the
  Daemon records their bytes. An Agent-written snapshot could claim bytes it
  never read. See ADR-0194.
- **Import the failed pass instead of changing Task Carry-Forward.** The QA
  gate reads the report where the Daemon committed it, and the three
  carry-forward paths stay untouched. See ADR-0194.
- **Always-observed rows are decided mechanically.** The provenance names and
  the `commit_range` kind decide, not the Agent's judgement. See ADR-0195.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Spec 0179's seven archived QA Reports and this repository's history. They
  show that every report declared repository inputs, that none recorded a
  snapshot, and that consecutive audited heads were ancestors changing 2 to 6
  files, except once, when 186 files changed. They also show that 24 of the
  41 re-executed passing rows had no declared input among the changed files.
- Spec 0192's two QA Report commits, `20db691c` on
  `roundfix/run-run_20260930T144240Z_33c18e8a522f7217` and `3327baa1` on
  `roundfix/run-run_20260930T181925Z_a1967faa3d75e7bc`. Their audited heads
  `563ca5aa` and `f7e540b3` are not ancestors of each other, and their trees
  differ only in `_prd.md` and five Task files.
- Mokhov, Mitchell and Peyton Jones, "Build Systems à la Carte" (ICFP 2018,
  <https://www.microsoft.com/en-us/research/wp-content/uploads/2018/03/build-systems-a-la-carte.pdf>),
  section 4.2.2. It defines a verifying trace that records the hashes of a
  key's dependencies and rebuilds only "if something has changed".
- The Go command's test cache documentation
  (<https://pkg.go.dev/cmd/go/internal/test>). It caches only successful
  results, and files a test reads "only match future runs in which the files
  and environment variables are unchanged".

## Research basis

No pending Inbox Entry addresses this repository. The Secondbrain was
consulted through `wiki/index.md` and
`qmd query "QA gate rerun only stale rows carry forward evidence"`. The query
returned the mirrored ADR-0097 and earlier QA Reports of this repository. The
concept page `wiki/concepts/agent-workflows-e-loop-engineering.md` records
the orchestrator owning order, isolation and retry, which is where this Spec
puts the carry decision. Exa located the two published sources above. The
paper separates verifying traces from dirty bits, and this Spec takes the
verifying trace: ancestry is a dirty-bit signal that a new commit identity
breaks, while digests survive it. The Go documentation supports carrying only
passes and keying reuse on the bytes a run read.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
