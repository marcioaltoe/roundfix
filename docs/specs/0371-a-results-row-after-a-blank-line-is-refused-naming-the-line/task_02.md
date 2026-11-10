---
task: task_02
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
status: completed
type: backend
complexity: medium
---

# Task 02: The Daemon tells the gate to continue the seeded table and returns a broken report once

## Overview

The Daemon's prompt for the seeded QA Report tells the gate where its rows go,
the cause of the defect reported by pantheon on 2026-11-09
([references](references/2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md)).
When the QA turn still ends with a Results row outside every table, the
Daemon returns the turn once, in the same single repair turn as an unresolved
evidence link or Borrowed Evidence, naming the line and the repair; a report
still broken settles through the shared eligibility refusal. The Daemon's
count of the rows it seeds reads the same reader.

## Requirements

1. MUST append, in `runQAGate` (`internal/daemon/task_engine.go`), the line
   `Add each row of the matrix to the seeded Results table directly below its last row, with no blank line before it and the same three columns, #, Status and Provenance; a row after a blank line is outside the table, and the report is refused naming that line.`
   followed by a newline, directly after the line that begins `Build and run
   every executable a row exercises` and before any self-audit line (API
   Contract 8). The `Seeded QA Report:` line and the Run Worktree lines keep
   their bytes.
2. MUST add `qaResultsRepairPrompt(reportPath string, brk spec.QAResultsBreak) string`
   returning the four lines of API Contract 9, each ending with a newline:
   `QA Report Results rows lie outside the Results table.`, `Report: <path>`,
   `<Detail()>; <Repair()>.` and
   `Keep every row in the seeded Results table, directly below its last row, with no blank line before it and the same three columns; never change a verdict or a row status.`
3. MUST, where `runQAGate` decides the repair turn, read the newest QA Report
   with `spec.NewestQAReport` and `spec.ReadQAReportFile`; when it has a
   `ResultsBreak`, publish a `runevent.KindDaemonQA` event with summary
   `QA Report Results row outside the Results table at line <n>; returning the QA turn once.`
   and payload `phase` `results_table`, `outcome` `repair`, `report` and
   `line`, print the same summary to the progress writer, and append
   `qaResultsRepairPrompt` after the evidence-link and Borrowed Evidence
   repair texts, so exactly one repair turn runs. The other repair texts keep
   their bytes.
4. MUST make `borrowedQAEvidence` return the report path and no item when the
   newest report has a `ResultsBreak`, so the break reaches settlement as the
   eligibility refusal and never as an infrastructure error.
5. MUST make `mechanicalBlockedRowCount` count, over
   `spec.ReadQAResults("", <the rendered body>).StatusTables()`, the rows whose
   `Status` cell starts with `blocked (<cause>:`, returning for every seed
   what it returns today (API Contract 10).
6. MUST create `internal/daemon/qa_results_table_test.go` with
   `TestQAResultsRowAfterABlankLineGetsOneRepairTurn`,
   `TestQAResultsRowAfterABlankLineAfterRepairSettlesFailed` and
   `TestQAPromptTellsTheGateToContinueTheSeededTable`, each calling
   `t.Parallel()` first, with a runner wrapping `taskFakeRunner` that
   recognises the repair prompt by its first line, as
   `qa_evidence_links_test.go` does. The scripted QA writes
   `qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Provenance |\n| --- | --- | --- |\n| 1 | pass | Requirement 1 |\n\n| 2 | pass | Requirement 2 |\n"`.
   The first test's repair removes the blank line; it asserts exactly one
   repair turn whose prompt equals `qaResultsRepairPrompt` for that report,
   the `results_table` event with outcome `repair` and the row's line, and the
   QA Task `completed` with a QA Report commit. The second makes no repair and
   asserts the QA Task `failed` with the settled reason
   `QA verdict pass not accepted: ` followed by the break's `Error()` for the
   path `spec.NewestQAReport` returns for the Spec directory. The third
   asserts the first QA prompt holds the requirement 1 line directly after
   the line that ends `never copy, replay or cite evidence or an executable an earlier Run left.`.
7. MUST keep these existing tests passing unchanged:
   `TestQAEvidenceLinksAndBorrowedEvidenceShareOneRepairTurn`,
   `TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory`,
   `TestQAEvidenceLinkRepairTurnKeepsTheEvidence`,
   `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed`,
   `TestQAReportWithResolvedLinksRunsNoRepairTurn`,
   `TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence`,
   `TestMechanicalStageWithholdsAgentSession`,
   `TestMechanicalQAReportSeedStartsPending` and
   `TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal`.
8. MUST NOT touch `internal/agent/`, `internal/spec/`, `internal/speccheck/`,
   `internal/cli/`, any skill, guide or Baseline module, `docs/agents/`, the
   Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: the QA prompt says nothing about where rows go, and a broken report gets no repair turn.
- [ ] Add the prompt line and its test.
- [ ] Add the results-table repair to the single repair turn, with its event and tests.
- [ ] Keep Borrowed Evidence from turning a break into an infrastructure error, and count seeded rows through the shared reader.

## Acceptance Criteria

- [ ] The QA prompt tells the gate to add rows directly below the seeded table's last row with no blank line.
- [ ] A report with a row after a blank line gets exactly one repair turn naming the line, and settles `completed` after the repair or `failed` with the shared text without it.
- [ ] The evidence-link and Borrowed Evidence repair texts, the seeded counts and every existing QA Daemon test are unchanged.

## Context

- instruction: `docs/adr/0381-a-results-row-outside-every-results-table-refuses-the-qa-report.md`
- instruction: `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- instruction: `internal/daemon/qa_evidence_links_test.go`
- instruction: `internal/daemon/qa_evidence_binding_test.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/qa_results_table_test.go`

## Consumers

- The QA prompt, the repair turn and the seeded counts — consumers: `./internal/daemon`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestQAResultsRowAfterABlankLineGetsOneRepairTurn|TestQAResultsRowAfterABlankLineAfterRepairSettlesFailed|TestQAPromptTellsTheGateToContinueTheSeededTable|TestQAEvidenceLinksAndBorrowedEvidenceShareOneRepairTurn|TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory|TestQAEvidenceLinkRepairTurnKeepsTheEvidence|TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed|TestQAReportWithResolvedLinksRunsNoRepairTurn|TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence|TestMechanicalStageWithholdsAgentSession|TestMechanicalQAReportSeedStartsPending|TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAResultsRowAfterABlankLineGetsOneRepairTurn TestQAResultsRowAfterABlankLineAfterRepairSettlesFailed TestQAPromptTellsTheGateToContinueTheSeededTable TestQAEvidenceLinksAndBorrowedEvidenceShareOneRepairTurn TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory TestQAEvidenceLinkRepairTurnKeepsTheEvidence TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed TestQAReportWithResolvedLinksRunsNoRepairTurn TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence TestMechanicalStageWithholdsAgentSession TestMechanicalQAReportSeedStartsPending TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; grep -q 'ReadQAResults(' internal/daemon/task_engine.go || { echo 'the seeded count does not read the shared reader' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'QA|Mechanical|Evidence|Prior|Carry' ./internal/daemon` — expected: exit 0. Before this Task the three new tests do not exist, so the first pass line is missing; afterwards they pass with the evidence-link, Borrowed Evidence, prompt and seeded-report tests, the seeded count reads the shared reader, and the Daemon's QA tests pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 4; Success Metric 6; Success Metric 7
- `_techspec.md` → Current behavior 6-8; API Contracts 8-10; Invariants 4-5; Testing Approach; Build Order 2
- ADR-0381, ADR-0351

## Result

Skills loaded: implement-task, coding-guidelines, systematic-debugging, no-workarounds, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Implemented this Task's Daemon slice in `internal/daemon/task_engine.go` and
added `internal/daemon/qa_results_table_test.go`. The pre-existing Task status
change remains Daemon-owned; no other Task or Task Graph was edited.

Acceptance evidence:

- The initial QA prompt adds the exact required table-continuation line
  immediately after the executable/evidence instruction and before self-audit.
  `TestQAPromptTellsTheGateToContinueTheSeededTable` asserts that ordering.
- The newest report's `ResultsBreak` appends the four-line
  `qaResultsRepairPrompt` after the existing repair texts. The Daemon publishes
  the required `results_table` event and prints its summary. The two
  `TestQAResultsRowAfterABlankLine...` tests assert one repair turn, the exact
  prompt, unchanged session and Batch, event payload and progress text. Removing
  the blank line accepts the QA Task and commits the repaired report; leaving
  it intact settles the fixture's QA Task failed with the shared eligibility
  refusal and the path returned by `spec.NewestQAReport`.
- `borrowedQAEvidence` leaves a Results break to eligibility instead of
  returning an infrastructure error. The unchanged evidence-link and Borrowed
  Evidence tests passed, including their existing shared repair test. The new
  `TestQAResultsAndEvidenceLinksShareOneRepairTurn` also asserts exact ordered
  concatenation of the link and Results repair prompts, one repair turn and
  committed regenerated evidence. `mechanicalBlockedRowCount` now reads the
  shared reader's Status tables and counts only matching Status cells. The
  existing mechanical report tests passed, including the exact populated seed
  bytes and frontmatter counts.

Focused checks:

- Before production edits,
  `go test -count=1 -timeout 3m -run 'TestQAResultsRowAfterABlankLine|TestQAPromptTellsTheGate' ./internal/daemon`
  failed as expected: both broken-report cases returned an evidence-binding
  infrastructure error naming row 13; the prompt test lacked the new line.
  After the implementation the same focused command exited 0.
- After adding the combined-repair regression,
  `go test -count=1 -timeout 5m -v -run 'TestQAResults|TestQAPrompt|TestQAEvidence|TestQAReportWithResolved|TestQABorrowed|TestMechanicalStage|TestMechanicalQAReport|TestQASettlement|TestWriteMechanicalQAReport' ./internal/daemon`
  exited 0 with every selected test passing. This includes all nine existing
  tests named in Requirement 7. The initial sandbox attempt could not access
  the Go build cache; the test commands above ran with approved cache access.
- A source comparison against HEAD confirmed that
  `qaEvidenceLinkRepairPrompt`, `qaBorrowedEvidenceRepairPrompt`, the seeded
  report line and the Run Worktree instruction lines remain byte-identical.
- `gofmt -l internal/daemon/task_engine.go internal/daemon/qa_results_table_test.go`
  produced no output; `git diff --check` exited 0.

The declared Verification, repository Verification and terminal settlement are
left to the Daemon. No follow-up outside this slice was needed.

## Carry-forward provenance

- Source Run: `run_20261110T215130Z_fa9f763e0a43c500`
- Source commit: `9e6983a6be83e12ccea06e0b476eb1de59389944`
