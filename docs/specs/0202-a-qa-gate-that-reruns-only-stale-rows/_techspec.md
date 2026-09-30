---
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
prd: _prd.md
created: 2026-09-30
---

# A QA gate that reruns only stale rows — Technical Spec

## Executive Summary

The carry machinery of ADR-0097 already exists in the mechanical stage:
`Carriable`, `resolveCarriedRows`, the `evidence_snapshots` reader and the
`carried (…)` row the seeded report materializes. This Spec supplies the three
pieces it lacked. The Daemon writes the snapshot when a QA pass closes. The
Daemon imports a failed pass's QA files from the Run Branch that holds them.
The mechanical stage accepts an establishing head proven by the QA Report
commit that recorded it. It also refuses always-observed rows and records a
disposition for every prior row. The trade-off this design accepts is trust in
the Agent's input declarations, which ADR-0097 already accepted, in exchange
for skipping whole rows: in Spec 0179, 24 of 41 re-executed passing rows. It
also accepts that a carry across an unintegrated pass depends on the failed
Run's branch still holding its QA Report commit.

## Project Constraints

- Identifier strategy: not applicable — no new entity identifier; carried rows
  keep their establishing row identifiers and imported reports keep their
  names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes, Git and the Spec
  tree only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0194 and ADR-0195 (this Spec)
  refine ADR-0097; ADR-0053, ADR-0057, ADR-0080, ADR-0093, ADR-0096,
  ADR-0117, ADR-0156, ADR-0166, ADR-0168, ADR-0170 and ADR-0176 hold; its gate is bound by
  ADR-0088, ADR-0091, ADR-0104, ADR-0155 and ADR-0167; ADR-0160, ADR-0161,
  ADR-0182, ADR-0183 and ADR-0184 do not apply, as the PRD records. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer's authorization of
  2026-09-30 for every skill, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Evidence Snapshot recorder | New `internal/speccheck/evidence_record.go` | Parses a closed report, snapshots each qualifying passing row at a head, rewrites only the `evidence_snapshots` key |
| Always-observed predicate | `internal/speccheck/evidence_record.go`, `internal/speccheck/report.go` | One function both the recorder and the carry resolver call; `EvidenceCommitRange` joins the input kinds |
| Carry resolver | `resolveCarriedRows` in `internal/speccheck/mechanical.go` | Recorded-at proof beside ancestry, always-observed refusal, one disposition per prior row, establishing provenance |
| Seeded report | `WriteMechanicalResult` in `internal/speccheck/report.go` | Carried rows keep provenance; a `## Row carry-forward` section when dispositions exist |
| Prior pass import | New `internal/daemon/qa_prior_pass.go`; `ReportShapeFindings` in `internal/speccheck/mechanical.go` | Finds the newest QA Report commit of the Spec outside `HEAD`'s history, copies its QA files, and keeps them only when the mechanical stage finds nothing in them |
| QA stage wiring | `runQAGate` in `internal/daemon/task_engine.go`, `buildQAPromptContext` in `internal/daemon/task_context.go` | Import before the prompt context, record after the Agent, new event phases and counts |
| Gate guidance | `qaGateContract` in `internal/agent/spec_prompt.go`, the qa-gate skill, `docs/user-guide/context-driven-development.md` | Keep carried rows, declare inputs on every executed row, never write the snapshot |

`Carriable`'s signature and its refusals stay as they are. Task Carry-Forward,
the Reconcile Command and Delivery Retry do not change.

## Implementation Design

### Interfaces

```go
// internal/speccheck/report.go — every emitted word is a constant here.
const EvidenceCommitRange EvidenceInputKind = "commit_range"
const CarryReasonNotPass = "not pass" // and one constant per reason

// CarryDisposition is one prior row's fate in the next pass.
type CarryDisposition struct {
	ID      string
	Carried bool
	Reason  string // empty when Carried; see the Vocabulary Contract
}
// MechanicalResult gains Dispositions []CarryDisposition.
// CarriedRow gains Provenance string.

// internal/speccheck/evidence_record.go
type EvidenceRecord struct {
	Head string
	Rows []string // recorded row IDs, in Results order
}

// RecordEvidenceSnapshots rewrites the report's evidence_snapshots key from
// the rows that can carry at head. Every other byte is kept.
func RecordEvidenceSnapshots(ctx context.Context, repoRoot, reportPath, head string) (EvidenceRecord, error)

// AlwaysObserved names why a row is never carried, or reports false.
func AlwaysObserved(provenance string, inputs []EvidenceInput) (reason string, observed bool)
```

```go
// internal/daemon/qa_prior_pass.go
type priorQAPass struct {
	Commit string   // the Daemon QA Report commit
	Head   string   // its first parent: the audited head of that pass
	Report string   // newest report path it holds, repository-relative
	Files  []string // every QA path it added or modified
}

func findPriorQAPass(ctx context.Context, workDir, specRelDir, slug string) (priorQAPass, bool, error)
func (engine *Engine) importPriorQAPass(ctx context.Context, plan TaskPlan, ordinal int) (priorQAPass, bool, error)
```

### Data Models

No database change. The QA Report frontmatter key `evidence_snapshots` keeps
the shape Spec 0080 defined and the mechanical stage already reads:

```yaml
evidence_snapshots:
  "3":
    head: <audited head>
    inputs:
      - ref: <declared ref, in declaration order>
        files:
          - path: <repository-relative path, sorted>
            sha256: <64 lowercase hex digits>
```

### Recording the snapshot

`RecordEvidenceSnapshots` parses the report with `parseMechanicalReport`. A
row qualifies when all of these hold:

- its status is exactly `pass`;
- it declares a non-empty `inputs:` list and every entry is a
  `repository_path`;
- `AlwaysObserved` reports false;
- `buildEvidenceSnapshots` resolves every input at `head`;
- every path its Evidence column cites is covered by an input's snapshot.

The key is removed from the frontmatter lines together with its indented
continuation lines. When at least one row qualifies, the new block is inserted
before the closing `---`, rendered through `yaml.v3` so that keys and paths
are quoted where needed. The body and every other frontmatter line keep their
bytes. A Git read error still strips an Agent-written key, and the error is
returned beside the stripped write. Only a failure to write the report is an
error the caller treats as infrastructure.

In `runQAGate`, after `settleQAVerdict` and before `settleTask`, the Daemon
calls the recorder on the settled report path. It records with the
`auditedHead` it resolved before the Agent turn, and publishes `daemon.qa`
phase `evidence_snapshots` (API Contract 4). A report path outside the Run
Worktree is recorded as outcome `skipped`. The QA Report commit then carries
the recorded key.

### Always-observed rows

`AlwaysObserved` splits the provenance on `;` and `,` the way
`qaReportProvenanceNames` does, and compares exact items:

- `repository Verification` → `always observed: repository Verification`;
- `Pull Request row` (`spec.QAPullRequestRowSource`) →
  `always observed: Pull Request row`;
- any input of kind `commit_range` → `always observed: commit_range input`.

The recorder never snapshots such a row, and the resolver never carries one,
even when an older report holds a snapshot for it.

### Carrying across an unintegrated pass

`resolveCarriedRows` keeps its order of checks and returns a disposition for
every prior row. The establishing-head step accepts either proof:

1. `mechanicalIsAncestor(head, current)`, as today;
2. `mechanicalReportRecordedAt`: a commit reachable from any ref, touching the
   establishing report's path, whose first parent is `head`. It must carry a
   `Roundfix-Spec` trailer equal to the Spec directory that holds the report,
   carry no `Roundfix-Task` trailer, and hold that path with the blob the
   establishing report's bytes hash to (`git hash-object --stdin`).

`mechanicalReport` keeps the report bytes it parsed so the second proof can
hash them. The changed-path set stays `git diff --name-only <head>..HEAD`,
which compares the two trees whether or not one descends from the other.
`Carriable` is unchanged. A new internal `carryRefusal` returns the reason
string, and `Carriable` reports `carryRefusal(...) == ""`. `ReportRow`'s
`AncestryVerified` field keeps its name, and its comment now names both
proofs. A carried row's `CarriedRow.Provenance` is the establishing row's
provenance cell.

### Importing the failed pass

`findPriorQAPass` never trusts a commit for its subject or trailers alone. It
walks the Runs the Run Database records for this Spec in this repository,
newest first, and considers only the QA settlement commit the Daemon recorded
for that Run in its Run Event Journal (the `daemon.commit` event of the `qa`
Task), read from that Run's Run Branch. A Run whose record or journal is
missing or pruned is skipped. It keeps the first such commit that also meets
all of these:

- its subject starts with `docs: qa report for <slug> (`;
- its trailers name `Roundfix-Spec: <slug>` and no `Roundfix-Task`;
- it is not an ancestor of `HEAD`.

A commit on any other ref, however well it imitates a QA Report commit, is
never imported.

The commit's QA paths come from `git diff-tree --no-renames --diff-filter=AM`
against its first parent, and its report is the newest `qa-report-*.md` among
them by `spec.NewestQAReportFromPaths`.

`importPriorQAPass` refuses, with a reason, when:

- a path already exists with different bytes: `path differs: <path>`;
- the imported report is not newer than the tree's newest report:
  `not newer than <report>`;
- the report is dated after the Daemon's today: `dated after today`;
- the imported report raises a finding in the mechanical stage's
  report-shape or evidence-path detectors: `report shape: <code>`.

Otherwise it writes every file from `git cat-file blob <commit>:<path>` with
mode `0644`, byte-identical. It then calls a new
`speccheck.ReportShapeFindings(repoRoot, reportPath string)
([]MechanicalFinding, error)`, which runs the existing `loadMechanicalReport`,
`detectMechanicalReportShape` and `detectMechanicalEvidencePaths` on the
written report. When that returns a finding, the import removes every file it
created and refuses. The import therefore never hands the next pass a finding
that the absent report would not have caused.

`runQAGate` calls it after the `before` snapshot and before
`buildQAPromptContext`, so the imported files ride in the QA Report commit.
`buildQAPromptContext` takes the imported pass's `Head` as
`PreviousReportHead` when an import happened. Git errors are infrastructure
errors of the QA step. A cancelled context publishes the stop, as the other
QA phases do.

### The seeded report

`WriteMechanicalResult` writes a carried row's provenance cell from
`CarriedRow.Provenance`, and keeps `report and head retained` when that is
empty. When `Dispositions` is non-empty, it appends after `## Mechanical
skips`:

```markdown

## Row carry-forward

| Prior row | Disposition |
| --- | --- |
| 2 | carried |
| 5 | re-run: input moved: internal/speccheck/report.go |
```

A result with no dispositions renders exactly the bytes it renders today, so
`task_engine_test.go`'s byte-exact report tests pass unedited. The table has
no `Status` column and lies outside `## Results`, so neither the hollow-report
reader nor the Results parser reads it as rows.

### API Contracts

1. API Contract: `evidence_snapshots` — Daemon-written frontmatter of a closed
   QA Report, in the shape under Data Models, keyed by row identifier. It is
   present only when at least one row qualifies. An Agent-written value never
   survives the Daemon's write.
2. API Contract: `commit_range` — a row input kind, declared as
   `- kind: commit_range` with a `ref` naming the range it read, such as
   `delivery base..audited head`. A row that declares it is never carried.
3. API Contract: seeded QA Report — a carried row reads
   `| <id> | carried (established by: <report>; head: <sha>) | <establishing provenance> |`.
   The `## Row carry-forward` section lists every prior row as `carried` or
   `re-run: <reason>`, with the reasons of the Vocabulary Contract.
4. API Contract: Run Event Stream — `daemon.qa` gains phase `prior_report`
   (outcome `imported`, `none` or `refused`; payload `commit`, `report`,
   `files`, `reason`) and phase `evidence_snapshots` (outcome `recorded`,
   `none`, `skipped` or `error`; payload `head`, `rows`, `report`, `error`).
   The `mechanical` phase payload gains `carried_rows` and `rerun_rows`.
5. API Contract: imported QA files — the failed pass's report and evidence
   land byte-identical under the Spec's `qa/` directory with their original
   names, and are staged in the next QA Report commit.

## Vocabulary Contract

This Spec coins three terms for the glossary owner: **Evidence Snapshot**,
the Daemon-recorded digests of one passing row's declared inputs at the
audited head; **Carried Row**, a row a later pass keeps instead of executing,
citing the report and head that established it; and **Carry Disposition**,
one prior row's recorded fate. The emitted words are:

- the input kind `commit_range`;
- the heading `## Row carry-forward`;
- the disposition `carried`, and `re-run: ` followed by one reason: `not pass`,
  `no inputs`, `non-repository input`, `always observed: repository Verification`,
  `always observed: Pull Request row`, `always observed: commit_range input`,
  `no evidence snapshot`, `establishing report unavailable`,
  `establishing head unproven`, `input moved: <paths>` or `evidence differs`;
- the `daemon.qa` phases `prior_report` and `evidence_snapshots` with the
  outcomes and payload fields of API Contract 4.

Every reason, the `commit_range` kind and the section heading are constants
in `internal/speccheck/report.go`, and every event word is written in
`internal/daemon/task_engine.go`, so the declarations below cover each emitted
word. task_04 documents each of them in
`docs/user-guide/context-driven-development.md` and asserts them there.

- emits: `internal/speccheck/report.go`
  pattern: `re-run: |commit_range|Row carry-forward|always observed: (?:repository Verification|Pull Request row|commit_range input)|establishing (?:head unproven|report unavailable)|no evidence snapshot|non-repository input|evidence differs|input moved: |not pass|no inputs`
  documented-in: `docs/user-guide/context-driven-development.md`
- emits: `internal/daemon/task_engine.go`
  pattern: `prior_report|evidence_snapshots|carried_rows|rerun_rows`
  documented-in: `docs/user-guide/context-driven-development.md`

## Coverage Map

- Goal 1 → Carrying across an unintegrated pass; Always-observed rows; The
  seeded report.
- Goal 2 → Importing the failed pass.
- Goal 3 → The seeded report; API Contract 3; API Contract 4.
- Goal 4 → Recording the snapshot; Carrying across an unintegrated pass.
- User Story 1 → Recording the snapshot; Carrying across an unintegrated pass.
- User Story 2 → Importing the failed pass.
- User Story 3 → The seeded report; API Contract 4.
- User Story 4 → Build Order 4; API Contract 3.
- Core Feature 1 → Recording the snapshot; API Contract 1.
- Core Feature 2 → Importing the failed pass; API Contract 5.
- Core Feature 3 → Carrying across an unintegrated pass.
- Core Feature 4 → Always-observed rows; API Contract 2.
- Core Feature 5 → The seeded report; API Contract 3.
- Core Feature 6 → Build Order 4.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1.
- Success Metric 5 → Testing Approach 2; Testing Approach 1.

## Integration Points

- **Task Carry-Forward.** Unchanged. It re-commits a failed Run's Tasks and
  skips the QA Report commit, which has no `Roundfix-Task` trailer. The import
  reads that commit from the Run Branch instead.
- **Run Branch reconciliation.** An imported report is byte-identical to the
  Run Branch's report, and the next pass's report is newer. So
  `supersedingQAReport` proves the old Run Branch superseded once the target
  holds both, which is today's rule.
- **Reopen and the QA settlement.** Unchanged. A carried row counts as a
  passed row when the Agent applies the verdict rules; the typed blocked
  counts ignore it.
- **External Spec Root.** A report outside the Run Worktree repository is
  neither recorded nor imported, as the carry already cannot resolve it.

## Testing Approach

1. **Recording.** New `internal/speccheck/evidence_record_test.go`, over
   temporary Git repositories:
   - every qualifying passing row is recorded with the head's digests;
   - fail, blocked, no-input, non-repository, `commit_range`,
     repository-Verification, Pull-Request-row and unresolved rows are not;
   - an Agent-written key is replaced;
   - every byte outside the key is unchanged, and `spec.ReadQAReportFile`
     reads the same verdict and counts;
   - a recorded report run through `RunMechanicalStage` at the same head
     carries the row, which proves the writer and the reader agree.

   New `internal/daemon/qa_evidence_snapshot_test.go`, over the task-cycle
   fixture (`newTaskCycleFixture`, `taskFakeRunner`): the QA Report commit
   carries the key at the audited head, and the `evidence_snapshots` event is
   published.
2. **Carrying.** New `internal/speccheck/qa_row_carry_test.go`:
   - a non-ancestor head recorded by a QA Report commit carries;
   - one recorded by a Task commit, or by no commit, does not;
   - always-observed rows never carry;
   - a moved input is named in the disposition;
   - a carried row keeps its provenance;
   - every prior row has one disposition;
   - a result without dispositions renders unchanged bytes.

   `TestCarriable` and the four `TestMechanicalStageCarriable…` tests in the
   governed `mechanical_test.go` run unedited in the same command.
3. **Importing.** New `internal/daemon/qa_prior_pass_test.go`, over temporary
   Git repositories with a side Run Branch:
   - the newest unintegrated QA Report commit is imported byte for byte;
   - a commit already in `HEAD`'s history, a Task commit and another Spec's
     commit are ignored;
   - a differing path, an older report and a report the shape detector
     refuses are refused with their reasons, and a refused import leaves no
     file behind;
   - end to end through `TaskCycle`, a first pass committed only on a side
     branch, followed by a head that re-commits its Task, makes the second
     pass's seeded report carry the unmoved row;
   - the prompt names the imported pass's head.
4. **Guidance.** Phrase checks on the qa-gate skill, its mirror and the guide;
   a prompt test in `internal/agent`; `make skills-sync-check`; the owned-skill
   version record; and the repository contract that pins the
   `### QA settlement` section.
5. **Real history.** The QA gate recomputes the Spec 0179 and Spec 0192
   measurements from archived reports and repository history, and reads the
   two published sources.

## Build Order

1. The Evidence Snapshot recorder, the always-observed predicate, the
   `commit_range` kind and the Daemon's recording step, task_01 (depends on:
   none).
2. The carry resolver's recorded-at proof, always-observed refusal,
   dispositions, carried provenance and the mechanical event counts, task_02
   (depends on: 1).
3. The prior pass import, `ReportShapeFindings` and the prompt head, task_03
   (depends on: 1, 2).
4. The qa-gate skill, the QA prompt contract and the Context-Driven
   Development guide, task_04 (depends on: 1, 2, 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

task_01, task_02 and task_03 all change `internal/daemon/task_engine.go`, so
they run in series.

## Risks & Considerations

- **Misdeclared inputs.** A row that under-declares can carry stale. The
  skill asks for conservative declarations, and always-observed rows cover
  the gate's own facts. This is ADR-0097's accepted risk.
- **Reclaimed Run Branches.** When the failed Run's branch and its QA Report
  commit are gone, an imported report's rows are re-run as
  `establishing head unproven`, and the import of that pass cannot happen.
  Both fail closed.
- **An imported report's findings.** The import keeps a report only when the
  mechanical stage's report-shape and evidence-path detectors find nothing in
  it, so a malformed failed pass cannot block the next one.
- **A stray QA commit.** The import takes the newest qualifying commit on any
  ref. A stale attempt can supply only rows whose inputs are still
  byte-identical, which is the carry proof itself.
- **Frontmatter rewrite.** The recorder edits a file the Agent wrote. The test
  that every byte outside the key survives, and the settle reader reading the
  same verdict and counts, pin that edit.
- **Declared break.** A seeded report that follows a previous report now ends
  with `## Row carry-forward`, and carried rows show their establishing
  provenance instead of `report and head retained`. No existing test pins
  either case.
- **Skill version.** The qa-gate version rises one step from the version on
  the Task's starting tree, so this Spec and Spec 0191 compose in either
  order.

## Decisions

- The Daemon, not the Agent, writes `evidence_snapshots`. See ADR-0194.
- An establishing head is proven by ancestry or by the QA Report commit that
  recorded it. See ADR-0194.
- The failed pass is imported by the QA stage rather than carried by Task
  Carry-Forward. See ADR-0194.
- Always-observed rows are decided by provenance names and the `commit_range`
  kind. See ADR-0195.
- Every prior row gets a disposition, and a carried row keeps its provenance.
  See ADR-0195.
- `Carriable` and the governed `mechanical_test.go` stay unchanged; new tests
  live in new files.
