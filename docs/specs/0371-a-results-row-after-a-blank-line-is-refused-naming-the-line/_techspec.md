---
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
prd: _prd.md
created: 2026-11-10
---

# A Results row after a blank line is refused, naming the line — Technical Spec

## Executive Summary

`internal/spec` gains one reader, `ReadQAResults`, that returns the tables of
a QA Report's Results section and the first table line no header leads, as a
`QAResultsBreak` whose text names the report path, the row's line, the table
above it and the repair. `readQAReport` records the break and the eligibility
decision refuses it before any count, so `qa-report accept`, the Settle
Command, `roundfix archive` and the Daemon's settlement name the real cause.
The Borrowed Evidence exemption, the environment-row count, the mechanical
stage and the Daemon's seeded count read rows only through it. The Daemon
tells the gate to continue the seeded table and returns the turn once on a
break; the qa-gate skill says the same. The trade-off accepted: a report a
GFM renderer shows with a paragraph of pipes is refused, not repaired by
Roundfix, and the gate pays one repair turn for a deleted line.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; QA Reports keep their names, frontmatter and columns. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; every reader works on local files. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0381 (new) records the reader, the
  refusal and the repair turn. ADR-0351: "The Daemon gives the gate one repair
  turn", which this refusal joins. ADR-0376: "Every QA settlement path applies
  it, after the evidence-link rule of ADR-0351", whose Carried Row exemption
  now reads the shared rows. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  maintainer's decision of 2026-11-10 06:20 BRT and the cycle BB rules answer
  quoted in the authorization record; no Makefile, workflow, Go module, lint
  or formatter change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0371-a-results-row-after-a-blank-line-is-refused-naming-the-line/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `9fda0b65`; every function is named so the plan survives line moves.

1. `qaReportResultsCounts` (`internal/spec/qa.go`) reads, under each `##
   Results` heading up to the next heading of depth one or two and outside
   fences, every header line followed by a separator line, and counts rows of
   tables with a `Status` column (`Provenance` for the exempt kinds) until the
   first line without cells. A table line after a blank line that no header
   leads is skipped in silence.
2. `QAReportEligibility` then compares the frontmatter counts with those
   derived rows; on the reduced Pantheon shape it returns
   `rows_blocked_environment is 1; expected 0` (characterized 2026-11-10).
3. `qaReportHollow` (`internal/spec/qa.go`) decides whether any row exists in
   the section or in any `Status` table of the report.
4. `borrowedEvidenceCarriedLines` (`internal/spec/qa_evidence_binding.go`)
   keeps the last header's `Status` column across blank lines, so it exempts a
   `carried (` line that no header leads.
5. `parseMechanicalResultsRows` (`internal/speccheck/mechanical.go`) reads only
   the first table after `## Results`, takes its first table line as header
   with or without a separator, does not skip fences, and stops at the first
   line without cells; `detectMechanicalReportShape` then reports
   "`<field>` is `<n>` but the Results table contains `<m>` matching rows" with
   the fix "Set `<field>` to the exact number of matching typed blocked rows."
6. `mechanicalBlockedRowCount` (`internal/daemon/task_engine.go`) counts
   `| blocked (<cause>:` byte sequences between `## Results` and
   `## Mechanical skips` of the seed it rendered.
7. `runQAGate` (`internal/daemon/task_engine.go`) appends `Seeded QA Report:
   <path>` and "Complete this report in place, preserving its materialized
   mechanical rows and skips; do not create another QA Report." to the prompt
   and says nothing about where new rows go; it gives one repair turn for
   unresolved links and Borrowed Evidence only.
8. `speccheck.WriteMechanicalResult` writes the seeded table as
   `| # | Status | Provenance |`, `| - | --- | --- |`, the carried and blocked
   rows, a blank line and `## Mechanical skips`; a gate adding rows just above
   that heading lands after the blank line, which is how Pantheon's did.

## System Architecture

- `internal/spec/qa_results.go` (new): `ReadQAResults`, `QAResults`,
  `QAResultsTable`, `QAResultsRow`, `QAResultsBreak` and its `Detail`,
  `Repair` and `Error` texts. It uses the package's `markdownHeading`,
  `markdownFenceMarker`, `markdownTableCells` and `markdownTableSeparator`.
- `internal/spec/qa.go`: `readQAReport` reads rows through `ReadQAResults`
  and stores `QAReport.ResultsBreak`; `qaReportResultsCounts` counts from the
  `Status` tables it returns; `QAReportEligibility` returns the break first;
  `EnvironmentRowsNeedingOverride` returns `0` for a report with a break.
- `internal/spec/qa_evidence_binding.go`: `borrowedEvidenceCarriedLines` takes
  carried rows from the shared tables; `FindBorrowedEvidence` returns the
  break as its error.
- `internal/speccheck/mechanical.go`: `parseMechanicalReport` takes rows from
  `spec.ReadQAResults`; `detectMechanicalReportShape` reports the break and
  skips count reconciliation. `parseMechanicalResultsRows` is removed.
- `internal/daemon/task_engine.go`: the seeded-report prompt line,
  `qaResultsRepairPrompt`, the break in the shared repair turn,
  `borrowedQAEvidence` leaving a broken report to settlement, and
  `mechanicalBlockedRowCount` through the shared reader.
- Unchanged by design: `internal/cli` (its QA paths reach the break through
  `QAReportEligibility`), `internal/agent`, `speccheck.WriteMechanicalResult`,
  `spec.WritePreconditionRefusalReport`, `qaReportHollow`, `settleQAVerdict`.
- Docs: the qa-gate skill; `docs/user-guide/commands/qa-report.md`,
  `settle.md`, `implement.md` and `docs/user-guide/context-driven-development.md`;
  the Roundfix `settle` and `implement` references; `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/spec/qa_results.go
type QAResultsRow struct{ Line int; Cells []string }
type QAResultsTable struct{ HeaderLine int; Header []string; Rows []QAResultsRow }
type QAResults struct{ Tables []QAResultsTable; Break *QAResultsBreak }
func ReadQAResults(path string, content []byte) QAResults
func (results QAResults) StatusTables() []QAResultsTable
type QAResultsBreak struct {
	Path                                         string
	Line, TableLine, LastRowLine, BlankFrom, BlankTo int
}
func (b QAResultsBreak) Detail() string
func (b QAResultsBreak) Repair() string
func (b QAResultsBreak) Error() string
// internal/spec/qa.go: QAReport gains ResultsBreak *QAResultsBreak
```

```go
// internal/daemon/task_engine.go
func qaResultsRepairPrompt(reportPath string, brk spec.QAResultsBreak) string
```

### Data Models

No schema change. `QAReport` gains the field `ResultsBreak`, nil when every
table line of the Results section belongs to a table.

### API Contracts

1. API Contract: `ReadQAResults` scope. Lines are numbered from 1 over the
   whole content, CRLF read as LF; a leading frontmatter block (first line
   `---` up to the next `---`) is skipped. The section starts at the first
   depth-2 heading whose text equals `Results` ignoring case and ends at the
   next heading of depth one or two; fenced blocks (` ``` ` or `~~~`) are
   skipped. A table line is one whose trimmed text starts and ends with `|`.
   A table starts at a table line whose next line is a separator line; its
   rows are the table lines below the separator up to the first line that is
   not a table line, separator lines excluded. `StatusTables` returns the
   tables whose header holds a cell equal to `Status` ignoring case and
   surrounding spaces. Nothing after the section is read.
2. API Contract: the break. The first table line of the section that starts
   no table and belongs to none ends the reading and is the break: `Line` is
   its line; `TableLine` is the header line of the last table above it, or 0
   when none; `LastRowLine` is that table's last row line, or its header line
   when it has no row; `BlankFrom` and `BlankTo` are the first and last line
   between that table and the row when every line between them is blank, and
   0 otherwise. A heading between a table and a later table line resets
   nothing: the line is still a break.
3. API Contract: break texts. `Detail()` is
   `Results row at line <Line> is outside the Results table that starts at line <TableLine>`,
   or `Results row at line <Line> is outside every Results table` when
   `TableLine` is 0. `Repair()` is, in this order of cases:
   `add a separator line directly below the table header at <Path>:<Line>`
   when `TableLine` is 0;
   `remove the blank line at <Path>:<BlankFrom> so the row continues the table`
   when one blank line separates them;
   `remove the blank lines at <Path>:<BlankFrom>-<BlankTo> so the row continues the table`
   when several do; otherwise
   `move the row directly below the last row of the table, at <Path>:<LastRowLine>`.
   `Error()` is `QA Report <Path>: <Detail()>; <Repair()>`. `Path` is the
   path the caller gave.
4. API Contract: QA Report reading. `readQAReport` calls
   `ReadQAResults(path, content)`, sets `ResultsBreak`, and derives the
   pre-PR, network-denied, network-prohibited, skipped and declared rows from
   `StatusTables()` with each table's own `Status` and `Provenance` columns,
   exactly as `qaReportResultsCounts` decides them today. Reading never fails
   on a break; verdict-only readers (the QA formatter, the pending check, the
   QA Report Import) keep working.
5. API Contract: eligibility. `QAReportEligibility` returns
   `*report.ResultsBreak` (its `Error()` text) before every other rule, for
   every verdict. `qa-report accept` therefore prints
   `roundfix: QA Report eligibility refused: <Error()>`, the Settle Command
   `roundfix: settle QA Report is ineligible: <Error()> (report <path>)`, and
   `roundfix archive` refuses with `no passing QA verdict: <Error()>`, through
   their existing wrappers. `EnvironmentRowsNeedingOverride` returns `0` when
   `ResultsBreak` is set, so the Delivery Queue does not park such a report as
   an environment partial.
6. API Contract: Borrowed Evidence. `borrowedEvidenceCarriedLines` exempts the
   rows of `StatusTables()` whose `Status` cell starts with `carried (`, and
   their `### <row-id>` blocks, as today. `FindBorrowedEvidence` returns
   `(nil, *results.Break)` when the report has a break; `RequireOwnRunEvidence`
   returns that error.
7. API Contract: mechanical stage. `parseMechanicalReport` builds its rows from
   `spec.ReadQAResults(path, content).StatusTables()` with the `Status`,
   `Evidence` and `Provenance` columns of each table, and keeps the break.
   With a break, `detectMechanicalReportShape` adds one
   `QA-REPORT-SHAPE` finding at the break's line, `Detail` = `Detail()`,
   `Fix` = `Repair()` with its first letter upper-cased and a final period,
   no `RowHint`; it still reports a frontmatter parse error, and then adds no
   other report-shape finding: no empty-matrix, row-status or count finding.
8. API Contract: Daemon prompt. In `runQAGate`, the line
   `Add each row of the matrix to the seeded Results table directly below its last row, with no blank line before it and the same three columns, #, Status and Provenance; a row after a blank line is outside the table, and the report is refused naming that line.`
   follows the line that begins `Build and run every executable a row
   exercises`, before any self-audit line.
9. API Contract: Daemon repair turn. Where the evidence-link and Borrowed
   Evidence repairs are decided, the Daemon reads the newest report with
   `spec.ReadQAReportFile`; a `ResultsBreak` adds
   `qaResultsRepairPrompt(reportPath, brk)` after the other repair texts in
   the same single repair turn. Its lines are
   `QA Report Results rows lie outside the Results table.`, `Report: <path>`,
   `<Detail()>; <Repair()>.` and
   `Keep every row in the seeded Results table, directly below its last row, with no blank line before it and the same three columns; never change a verdict or a row status.`
   It publishes a `daemon.qa` event with payload `phase` `results_table`,
   `outcome` `repair`, `report` and `line`, and prints
   `QA Report Results row outside the Results table at line <n>; returning the QA turn once.`
   `borrowedQAEvidence` returns no item for a report with a break. A break
   left after the turn settles the QA Task through `settleQAVerdict` unchanged:
   `QA verdict <verdict> not accepted: <Error()>` for `pass` or `partial`.
10. API Contract: seeded count. `mechanicalBlockedRowCount` counts, over
    `spec.ReadQAResults("", mechanicalBody).StatusTables()`, the rows whose
    `Status` cell starts with `blocked (<cause>:`; for every seed
    `WriteMechanicalResult` renders it returns what it returns today.
11. API Contract: qa-gate skill. A section `### Results rows stay inside the
    table` after `### Evidence belongs to this Run` and before `### Declared
    rows name their criterion`, outside `### QA settlement`.

### Surface Transcripts

None. The stderr lines are asserted byte for byte in
`internal/cli/qa_report_test.go` over temporary workspaces; the Daemon texts
are asserted in `internal/daemon` with the scripted QA Agent.

### Invariants

1. A report whose every Results table line belongs to a table reads exactly as
   before: the same derived rows, eligibility, findings and exemptions.
2. No reader counts a row of a report with a break, and each one that decides
   names the break.
3. Archived reports are never rewritten or re-read by archive.
4. The evidence-link and Borrowed Evidence repair prompts keep their bytes; the
   Daemon runs at most one repair turn.
5. The seeded table and its counts keep their bytes.

### Declared breaks

- A report with a table line no header leads under `## Results` was accepted
  (the row ignored) and is now refused on every QA settlement path and by
  `roundfix archive`. No existing test pins such a report: the prototype
  measurement in `_prd.md` → Validation ran every QA test green.
- The mechanical stage now reads every `Status` table of the section, skips
  fences and requires the separator line, as the eligibility check did; a
  header without a separator is now a break. `TestResultsTableIsTheOnlyRowSource`
  keeps passing, because its comparison table has no `Status` column.
- On a break the mechanical stage reports the break instead of count
  mismatches.
- `EnvironmentRowsNeedingOverride` is `0` for a broken report, so its Delivery
  Queue park is no longer `qa-environment-partial`.
- The `## Results` heading matches ignoring case in the counting reader, as in
  the mechanical stage.

## Coverage Map

- Goal 1, Story 1, Core Features 1-2, Success Metrics 1, 2 and 4 →
  `ReadQAResults` and `QAResultsBreak` (API Contracts 1-3), `readQAReport` and
  `QAReportEligibility` (API Contracts 4-5), `FindBorrowedEvidence` (API
  Contract 6); Build Order 1.
- Goal 2 → `QAResultsBreak.Detail` and `QAResultsBreak.Repair` (API Contract
  3) and the report-shape finding of `detectMechanicalReportShape` (API
  Contract 7); Build Order 1.
- Goal 3, Story 3, Core Features 1 and 3, Success Metrics 3, 5 and 7 →
  `QAResults.StatusTables`, read by `readQAReport`, `borrowedEvidenceCarriedLines`,
  `parseMechanicalReport` and `mechanicalBlockedRowCount` (API Contracts 1, 4,
  6, 7 and 10), with `parseMechanicalResultsRows` removed; Build Orders 1-2.
- Goal 4, Story 2, Core Feature 4, Success Metric 6 → the seeded-report prompt
  line and `qaResultsRepairPrompt` in `runQAGate` (API Contracts 8-9); Build
  Order 2.
- Core Feature 5 → the qa-gate section `### Results rows stay inside the
  table` (API Contract 11), the guides, references and `CONTEXT.md`; Build
  Order 3.

Deviations from the PRD: None.

## Absence Evidence

1. `ReadQAResults` and `QAResultsBreak` (the one reader and its refusal) — command: `grep -rn 'ReadQAResults\|QAResultsBreak\|outside the Results table' internal`; result: no match; four readers carry their own table rules (Current behavior 1, 4, 5, 6).
2. The Daemon's results-table repair and prompt line — command: `grep -n 'results_table\|seeded Results table\|Results rows lie outside' internal/daemon/task_engine.go`; result: no match; the repair turn covers links and Borrowed Evidence only.
3. The qa-gate section — command: `grep -n 'Results rows stay inside the table\|no blank line' .agents/skills/qa-gate/SKILL.md`; result: no match.

## ADR Alternatives

- ADR-0381: the design takes none of its rejected options: it does not read
  through the blank line, does not narrow to the first table, does not refuse
  a second `Status` table, and keeps the repair turn.
- ADR-0351: rejected warning only and no repair turn; this refusal likewise
  stops settlement and joins its single repair turn.
- ADR-0376: rejected prose-only rules; the Borrowed Evidence exemption keeps
  its meaning and refuses a report it cannot read whole instead of exempting
  a row no table holds.

## Integration Points

None beyond local files: no network, no Git call added (the Borrowed
Evidence check keeps its Git reads).

## Testing Approach

The reduced Pantheon shape, used by every Task, is this report content
(lines: header 12, blank 15, first outside row 16):

```text
---\nverdict: partial\nrows_blocked_environment: 1\nrows_blocked_finding: 0\nrows_blocked_declared: 1\n---\n\n# QA Report\n\n## Results\n\n| #   | Status | Provenance |\n| --- | ------ | ---------- |\n| Q01 | pass | PRD Goal 1 |\n\n| Q02 | pass | Task 01 Acceptance Criterion 1 |\n| Q03 | blocked (environment: no open Pull Request) | Pull Request row |\n| Q04 | blocked (declared: the probe reads the live service) | PRD Unreachable Acceptance: the probe reads the live service |\n\n## Mechanical skips\n\nNone.\n
```

with a `_prd.md` holding
`## Unreachable Acceptance\n\n- criterion: the probe reads the live service\n  reason: needs a deploy\n  satisfied-by: the operator after release\n`.

- `internal/spec/qa_results_test.go` (new), each test calling `t.Parallel()`
  first:
  - `TestQAResultsRefusesARowAfterABlankLine`: the shape above → break line
    16, table 12, blank 15, the exact `Error()` text, eligibility returns it
    and never contains `rows_blocked_environment`; with line 15 removed, no
    break, `RowsBlockedPrePullRequest` 1, one declared row `Q04`, eligibility
    nil.
  - `TestQAResultsNamesEveryBlankLineOrTheRowsPlace`: two blank lines →
    `remove the blank lines at <path>:15-16 ...` (row line 17); `Some prose.`
    instead of the blank line → `move the row directly below the last row of
    the table, at <path>:14` (row line 16); a header with no separator as the
    first table line → `Results row at line 12 is outside every Results
    table; add a separator line directly below the table header at
    <path>:12`.
  - `TestQAResultsReadsEveryStatusTable`: a second table with the same header
    after a blank line, a `Status` table under `### Outside evidence`, a table
    without a `Status` column, a fenced table, and a table under `## Notes` →
    no break; rows of the first two counted, the others not.
  - `TestQAResultsBreakStopsEveryCountingReader`: on the shape above
    `EnvironmentRowsNeedingOverride` is 0, `archiveUnprovenActions` returns
    the break, and `FindBorrowedEvidence` returns an error matching
    `QAResultsBreak` with `errors.As`.
- `internal/cli/qa_report_test.go`: `TestQAReportAcceptRefusesARowAfterABlankLine`
  and `TestSettleRefusesARowAfterABlankLine` over a `pass` report whose row
  `Q02` follows a blank line (content in task_01), exact stderr, then
  accepted (and settled) once the blank line is removed.
- `internal/speccheck/mechanical_results_test.go` (new):
  `TestMechanicalReportShapeNamesARowAfterABlankLine` (one finding at line 16
  with the repair, no `Results table contains` finding) and
  `TestMechanicalStageReadsEveryStatusTable` (a `pending` row in a second
  `Status` table is reported).
- `internal/daemon/qa_results_table_test.go` (new), with a wrapper runner over
  `taskFakeRunner` that recognises the repair prompt by its first line, as
  `qa_evidence_links_test.go` does:
  `TestQAResultsRowAfterABlankLineGetsOneRepairTurn`,
  `TestQAResultsRowAfterABlankLineAfterRepairSettlesFailed` and
  `TestQAPromptTellsTheGateToContinueTheSeededTable`.
- Existing seams that must keep passing unchanged:
  `TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows`,
  `TestArchivedPassCorpusRemainsArchiveEligible`,
  `TestResultsTableIsTheOnlyRowSource`,
  `TestQAEvidenceLinksAndBorrowedEvidenceShareOneRepairTurn`,
  `TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory`,
  `TestQAEvidenceLinkRepairTurnKeepsTheEvidence` and
  `TestBorrowedEvidenceExemptsACarriedRow`.

## Build Order

1. The shared reader, the QA Report reading, eligibility, the environment-row
   count, the Borrowed Evidence exemption and the mechanical stage, with their
   tests in `internal/spec`, `internal/speccheck` and `internal/cli`.
2. The Daemon's prompt line, repair turn and seeded count, with their tests
   (depends on: 1).
3. The qa-gate skill section, the guides, the Roundfix settle and implement
   references, `CONTEXT.md`, the owned-skill versions and the Behavior Surface
   Record (depends on: 1).
4. Final QA gate (depends on: 2, 3).

## Risks & Considerations

- An adopter whose gate already writes rows after a blank line sees refusals
  where it saw wrong counts; the refusal names the line and the repair, and
  in a Run the repair turn fixes it.
- A table line inside a row-detail block under `### <row-id>` that no header
  leads is now a break too; the skill says where rows go, and the refusal
  names the line.
- Rebase order: 0372 merges first and may raise the Roundfix skill version and
  re-record the Behavior Surface Record; 0373 follows. Every shared record is
  regenerated by its command on the Task's own base.
- No file carries a build constraint.

## Glossary

- changes: **QA Report**

## Decisions

- Refuse, never read through; one reader with the eligibility check's rule.
  See ADR-0381.
- The break is recorded on `QAReport` and refused by the eligibility decision,
  so verdict-only readers keep working.
- One repair turn in a Run, beside the existing two.
- Three implementation Tasks and the QA gate: the reader with its consumers,
  the Daemon, and the documentation.
