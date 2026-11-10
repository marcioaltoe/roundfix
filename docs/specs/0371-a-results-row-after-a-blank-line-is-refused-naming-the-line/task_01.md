---
task: task_01
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
status: completed
type: backend
complexity: medium
---

# Task 01: One reader of Results rows refuses a row no table header leads, on every QA settlement path

## Overview

`internal/spec` gains `ReadQAResults`, the one reader of a QA Report's Results
section, and `QAResultsBreak`, the refusal of a table line no header leads.
`readQAReport` records the break, the eligibility decision refuses it before
any other rule, and the environment-row count, the Borrowed Evidence
exemption and the mechanical stage read Results rows only through the new
reader. `roundfix qa-report accept`, the Settle Command and `roundfix archive`
then name the blank line instead of a blocked count. This answers the defect
reported by pantheon on 2026-11-09 and recorded by the operator on 2026-11-10
([references](references/2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md)).

## Requirements

1. MUST create `internal/spec/qa_results.go` with `QAResultsRow`,
   `QAResultsTable`, `QAResults`, `ReadQAResults`, `QAResults.StatusTables`
   and `QAResultsBreak` with `Detail`, `Repair` and `Error`, exactly as
   `_techspec.md` → Interfaces and API Contracts 1-3 define them, using the
   package's `markdownHeading`, `markdownFenceMarker`, `markdownTableCells`
   and `markdownTableSeparator`. No new package import.
2. MUST make `readQAReport` (`internal/spec/qa.go`) call
   `ReadQAResults(path, content)`, store the break in a new `QAReport` field
   `ResultsBreak *QAResultsBreak`, and derive the pre-PR, network-denied,
   network-prohibited, skipped and declared rows from `StatusTables()` with
   each table's own `Status` and `Provenance` columns, keeping today's rules
   for each kind (API Contract 4). `qaReportHollow` keeps its code. A break
   never makes `readQAReport` return an error.
3. MUST make `QAReportEligibility` return the break (`*report.ResultsBreak`)
   before every other rule for every verdict, and make
   `EnvironmentRowsNeedingOverride` return `0` when `ResultsBreak` is set
   (API Contract 5). Change no other function of `internal/spec/qa.go`
   except removing `qaReportResultsCounts` if the new derivation replaces it.
4. MUST make `borrowedEvidenceCarriedLines`
   (`internal/spec/qa_evidence_binding.go`) exempt only the `carried (` rows
   of `StatusTables()` and their `### <row-id>` blocks, and make
   `FindBorrowedEvidence` return `nil` items and the break as its error when
   the report has one (API Contract 6).
5. MUST make `parseMechanicalReport` (`internal/speccheck/mechanical.go`)
   build its rows from `spec.ReadQAResults(path, content).StatusTables()` with
   the `Status`, `Evidence` and `Provenance` columns of each table, keep the
   break on `mechanicalReport`, and remove `parseMechanicalResultsRows`. With
   a break, `detectMechanicalReportShape` adds exactly one `QA-REPORT-SHAPE`
   finding at the break's line with `Detail` equal to `Detail()` and `Fix`
   equal to `Repair()` with its first letter upper-cased and a final period,
   no `RowHint`, and adds no other report-shape finding after a frontmatter
   parse finding (API Contract 7).
6. MUST create `internal/spec/qa_results_test.go` with
   `TestQAResultsRefusesARowAfterABlankLine`,
   `TestQAResultsNamesEveryBlankLineOrTheRowsPlace`,
   `TestQAResultsReadsEveryStatusTable` and
   `TestQAResultsBreakStopsEveryCountingReader`, each calling `t.Parallel()`
   first, with every case `_techspec.md` → Testing Approach lists, on the
   reduced Pantheon shape written there byte for byte and its `_prd.md`
   declaration. The first test asserts the exact text
   `QA Report <path>: Results row at line 16 is outside the Results table that starts at line 12; remove the blank line at <path>:15 so the row continues the table`
   from both `ResultsBreak.Error()` and `QAReportEligibility`, asserts the
   eligibility text does not contain `rows_blocked_environment`, then removes
   line 15 and asserts no break, `RowsBlockedPrePullRequest` 1, one declared
   row with ID `Q04`, and eligibility nil.
7. MUST create `internal/speccheck/mechanical_results_test.go` (package
   `speccheck_test`) with `TestMechanicalReportShapeNamesARowAfterABlankLine`
   and `TestMechanicalStageReadsEveryStatusTable`, each calling `t.Parallel()`
   first, through `runMechanical` over `newMechanicalGitRepo` with the report
   at `docs/specs/mechanical/qa/report.md`: on the reduced Pantheon shape
   exactly one `QA-REPORT-SHAPE` finding, line 16, `Detail`
   `Results row at line 16 is outside the Results table that starts at line 12`,
   `Fix`
   `Remove the blank line at docs/specs/mechanical/qa/report.md:15 so the row continues the table.`,
   and no finding whose `Detail` contains `Results table contains`; and a
   `pass` report with a second `| # | Status | Provenance |` table after a
   blank line whose row `R02` is `pending` gives the finding
   `row R02 remains pending instead of carrying a terminal status`.
8. MUST add `TestQAReportAcceptRefusesARowAfterABlankLine` and
   `TestSettleRefusesARowAfterABlankLine` to `internal/cli/qa_report_test.go`,
   each calling `t.Parallel()` first, on a workspace from
   `newImplementWorkspace` with `implementQAGateSeed`, as
   `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink` and
   `TestSettleRefusesAQAReportWithBorrowedEvidence` build theirs, with the
   report `qa/qa-report-2026-11-10.md` under the Spec directory holding
   `---\nverdict: pass\n---\n\n## Results\n\n| # | Status | Provenance |\n| - | --- | --- |\n| Q01 | pass | Requirement 1 |\n\n| Q02 | pass | Requirement 2 |\n`.
   The exact stderr is
   `roundfix: QA Report eligibility refused: QA Report <path>: Results row at line 11 is outside the Results table that starts at line 7; remove the blank line at <path>:10 so the row continues the table\n`
   for `qa-report accept`, where `<path>` is the path the test passed, and
   `Settle surface: <root>\nroundfix: settle QA Report is ineligible: QA Report <path>: Results row at line 11 is outside the Results table that starts at line 7; remove the blank line at <path>:10 so the row continues the table (report <report path relative to root>)\n`
   for settle, where `<path>` is the path `spec.NewestQAReport` returns for
   the Spec directory the command reads, and the Task file stays unchanged.
   Each test then removes line 10 and asserts exit 0 (for settle, the Task
   `completed`).
9. MUST keep these existing tests passing unchanged:
   `TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows`,
   `TestArchivedPassCorpusRemainsArchiveEligible`,
   `TestResultsTableIsTheOnlyRowSource`,
   `TestBorrowedEvidenceExemptsACarriedRow` and
   `TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory`.
10. MUST NOT touch any non-test file under `internal/cli`, `internal/daemon`,
    `internal/agent`, `internal/speccheck/report.go`, any skill, guide or
    Baseline module, `docs/agents/`, the Makefile, CI workflows, `go.mod` or
    `go.sum`.

## Subtasks

- [ ] Characterize: on the reduced Pantheon shape the eligibility decision returns `rows_blocked_environment is 1; expected 0` and the mechanical stage two count findings.
- [ ] Add the shared reader, its break texts and their tests.
- [ ] Read QA Report rows, eligibility, the environment-row count and the Borrowed Evidence exemption through it.
- [ ] Read the mechanical stage's rows through it and report the break.
- [ ] Add the `qa-report accept` and Settle Command tests.

## Acceptance Criteria

- [ ] On the reduced Pantheon shape every reader names line 16 and the blank line 15, and no reader counts or exempts a row outside every table.
- [ ] Two blank lines, a row prose separates and a header without its separator each get their own repair text.
- [ ] Reports whose every table line belongs to a table read exactly as before, including second `Status` tables, tables under deeper headings and the archived pass corpus.
- [ ] `qa-report accept` and the Settle Command refuse the broken report with the exact text and accept it once the blank line is removed.

## Context

- instruction: `docs/adr/0381-a-results-row-outside-every-results-table-refuses-the-qa-report.md`
- instruction: `docs/adr/0376-a-qa-report-relies-only-on-evidence-its-own-run-produced.md`
- instruction: `internal/spec/spec.go`
- instruction: `internal/spec/qa_partial_policy_test.go`
- instruction: `internal/speccheck/mechanical_test.go`
- creates: `internal/spec/qa_results.go`
- creates: `internal/spec/qa_results_test.go`
- creates: `internal/speccheck/mechanical_results_test.go`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/qa_evidence_binding.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/cli/qa_report_test.go`

## Consumers

- `spec.ReadQAResults`, `spec.QAResultsBreak`, `spec.QAReport.ResultsBreak` and `spec.QAReportEligibility` — consumers: `./internal/spec`, `./internal/speccheck`, `./internal/cli`
- checkpoint: task_02

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestQAResultsRefusesARowAfterABlankLine|TestQAResultsNamesEveryBlankLineOrTheRowsPlace|TestQAResultsReadsEveryStatusTable|TestQAResultsBreakStopsEveryCountingReader|TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows|TestArchivedPassCorpusRemainsArchiveEligible|TestBorrowedEvidenceExemptsACarriedRow)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAResultsRefusesARowAfterABlankLine TestQAResultsNamesEveryBlankLineOrTheRowsPlace TestQAResultsReadsEveryStatusTable TestQAResultsBreakStopsEveryCountingReader TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows TestArchivedPassCorpusRemainsArchiveEligible TestBorrowedEvidenceExemptsACarriedRow; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestMechanicalReportShapeNamesARowAfterABlankLine|TestMechanicalStageReadsEveryStatusTable|TestResultsTableIsTheOnlyRowSource)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMechanicalReportShapeNamesARowAfterABlankLine TestMechanicalStageReadsEveryStatusTable TestResultsTableIsTheOnlyRowSource; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestQAReportAcceptRefusesARowAfterABlankLine|TestSettleRefusesARowAfterABlankLine|TestQAReportAcceptRefusesAnUnresolvedEvidenceLink|TestSettleRefusesAQAReportWithBorrowedEvidence)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportAcceptRefusesARowAfterABlankLine TestSettleRefusesARowAfterABlankLine TestQAReportAcceptRefusesAnUnresolvedEvidenceLink TestSettleRefusesAQAReportWithBorrowedEvidence; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/spec ./internal/speccheck && go test -count=1 -timeout 20m -run 'QAReport|Settle|Archive|QAPartial|PartialShape|Deliver' ./internal/cli` — expected: exit 0. Before this Task none of the eight new tests exists, so the first pass line is missing; afterwards they pass with the existing partial-policy, archived-corpus, Borrowed Evidence, mechanical, evidence-link and settle tests, the Spec and Spec-check suites and the CLI's QA Report, settle, archive and delivery tests.
- `grep -q 'ReadQAResults(' internal/spec/qa.go || { echo 'readQAReport does not read the shared reader' >&2; exit 1; }; grep -q 'ReadQAResults(' internal/spec/qa_evidence_binding.go || { echo 'Borrowed Evidence does not read the shared reader' >&2; exit 1; }; grep -q 'ReadQAResults(' internal/speccheck/mechanical.go || { echo 'the mechanical stage does not read the shared reader' >&2; exit 1; }; ! grep -q 'func parseMechanicalResultsRows' internal/speccheck/mechanical.go || { echo 'the second reader is still present' >&2; exit 1; }` — expected: exit 0. Before this Task no file calls the shared reader, so the first check fails; afterwards the three readers call it and the mechanical stage's own reader is gone.

## References

- `_prd.md` → Goals; User Stories 1 and 3; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 7
- `_techspec.md` → Current behavior 1-5; API Contracts 1-7; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1
- ADR-0381, ADR-0376

## Result

Skills loaded: implement-task, coding-guidelines, systematic-debugging, no-workarounds, golang-structs-interfaces, golang-error-handling, golang-testing, golang-cli, agentic-cli-design, testing-boss, tech-writer, evidence-gate.

Implemented the shared Results reader and its typed break. QA Report eligibility,
environment override classification, Borrowed Evidence and mechanical shape
checking now use that reader. A broken report remains readable for its verdict,
but contributes no counted or exempted rows. Eligibility returns the break before
verdict, hollow, finding, environment, skipped or declared-row policy. Mechanical
shape checking preserves a frontmatter parse finding, emits the break once with
its repair, and skips every other shape finding for that broken report.

Characterization before the production changes:

- The reduced Pantheon report returned `rows_blocked_environment is 1; expected 0`.
- Its mechanical check emitted two misleading count findings, for environment
  and declared rows, and ignored a pending row in a second Status table.
- Both new CLI regression tests observed exit 0 for the broken pass report.
  The new regression assertions failed on those behaviors as expected.

Acceptance evidence from focused checks in this worktree:

- Criterion 1: `TestQAResultsRefusesARowAfterABlankLine` asserts the exact
  path-preserving refusal, row 16, table 12, blank 15 and LF/CRLF line numbering;
  removing line 15 yields one pre-PR row, declared row Q04 and eligible partial.
  `TestQAResultsBreakStopsEveryCountingReader` checks every verdict's refusal
  precedence, zero override count, nil archive actions, typed Borrowed Evidence
  errors and no exemption for an orphan carried row. The mechanical regression
  asserts the exact Detail and Fix, no RowHint and no count reconciliation, with
  and without a frontmatter parse error.
- Criterion 2: `TestQAResultsNamesEveryBlankLineOrTheRowsPlace` checks exact
  repairs for two blank lines, intervening prose and a missing separator. It
  also covers an intervening deeper heading and an empty preceding table.
- Criterion 3: `TestQAResultsReadsEveryStatusTable` covers repeated headers,
  a deeper heading, reordered columns, a case-insensitive Results heading,
  both fence kinds, non-Status tables, section bounds and carried-row/detail
  exemptions using each table's own Status column. The second-table mechanical
  test detects R02 pending. Existing override, archived-pass, Results-only,
  Borrowed Evidence, partial-policy, precondition and mechanical carry tests
  passed unchanged in the focused selections below.
- Criterion 4: both new CLI tests assert the complete specified stderr and exit
  1, then exit 0 after removing line 10. The settle test checks the fixture Task
  bytes remain unchanged on refusal and its status becomes completed on repair.

Focused commands and outcomes (all final runs exited 0):

- `rtk proxy go test ./internal/spec -run '^Test(QAResults|EnvironmentRowsNeedingOverride|ArchivedPassCorpus|BorrowedEvidence)' -count=1` — passed, including all four new reader tests and the existing Borrowed Evidence regressions.
- `rtk proxy go test ./internal/spec -run '^Test(QAReportEligibility|ReadQAReport|APartial|ANetwork|APullRequestRow)' -count=1` — passed, preserving the existing reader and partial-policy rules.
- `rtk proxy go test ./internal/speccheck -run '^Test(MechanicalReportShape|MechanicalStageReadsEveryStatusTable|ResultsTableIsTheOnlyRowSource|CountDisagreementReportsItsCause|MechanicalStageAcceptsThePreconditionRefusalRow|MechanicalEvidencePath|MechanicalStageCarriable)' -count=1` — passed. An initially incorrect line-18 assertion in the new second-table fixture was corrected to its actual line 16; the specified pending-row diagnostic is unchanged.
- `rtk proxy go test ./internal/cli -run '^Test(QAReportAcceptRefuses|SettleRefusesA(QAReport|RowAfterABlankLine))' -count=1` — passed, including both new exact-diagnostic/repair tests and existing evidence-link and Borrowed Evidence refusals.
- `rtk proxy git diff --check` — passed. Changed-path inspection contains only
  this Task file and the seven declared Go paths. The Task's in-progress status
  was present before implementation and has been preserved.

Initial sandboxed test attempts could not access the default Go build cache;
the focused tests ran with approved escalation, without changing cache or
toolchain configuration.

The declared Verification, repository Verification and final QA were not run;
they remain Daemon-owned. No commit, push, PR, Task Graph or sibling Task edit.
Daemon integration and documentation remain in their assigned downstream Tasks.

## Carry-forward provenance

- Source Run: `run_20261110T215130Z_fa9f763e0a43c500`
- Source commit: `71b2ec2872ff4a1969f068ecb4f5921a524000e0`
