---
task: task_03
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: completed
type: backend
complexity: high
---

# Task 03: The next pass reads a failed pass's QA Report from its Run Branch

## Overview

A failed gate's Run is never integrated, and Task Carry-Forward re-commits its
Tasks without its QA Report commit. The next pass therefore finds no previous
report and executes every row again. This Task makes the Daemon's QA stage
bring in the newest QA Report commit of the same Spec that the Run's history
does not contain, before the prompt context and the mechanical stage read the
previous report. Its report and evidence are copied byte for byte and then
ride in the next QA Report commit. With task_01 and task_02, the next pass
carries every unmoved row of the failed pass.

## Requirements

1. MUST add `findPriorQAPass` and `importPriorQAPass` in the new file
   `internal/daemon/qa_prior_pass.go`, with the selection rule of the
   TechSpec's "Importing the failed pass" section: only the QA settlement commit the Run Event Journal records for a recorded Run of this Spec, read from that Run's Run Branch, is a candidate, and `TestAPriorPassIsImportedOnlyFromARecordedRunCommit` MUST prove that a fabricated commit with the same subject and trailers on another ref is ignored. The subject prefix is
   `docs: qa report for <slug> (`, the trailers name `Roundfix-Spec: <slug>`
   and no `Roundfix-Task`, and the commit is not an ancestor of `HEAD`. The
   search reads every ref and takes the newest commit in date order.
2. MUST copy every QA path the commit added or modified, from its Git blobs,
   byte-identical with mode `0644`. A path whose existing bytes are identical
   is kept.
3. MUST refuse the import with the four reasons of that section:
   `path differs: <path>`, `not newer than <report>`, `dated after today` and
   `report shape: <code>`. A refused import MUST leave no file it created.
4. MUST add `speccheck.ReportShapeFindings(repoRoot, reportPath string)
   ([]MechanicalFinding, error)`. It runs the existing report loader and the
   report-shape and evidence-path detectors, and adds no detector of its own.
5. MUST call the import in `runQAGate` after the `before` snapshot and before
   `buildQAPromptContext`. The imported files MUST be staged in the QA Report
   commit. `buildQAPromptContext` MUST report the imported pass's first parent
   as `PreviousReportHead` when an import happened.
6. MUST publish `daemon.qa` phase `prior_report` with outcome `imported`,
   `none` or `refused` and the payload of the TechSpec's API Contract 4. Git
   errors are infrastructure errors of the QA step, and a cancelled context
   publishes the stop.
7. MUST NOT change Task Carry-Forward, the Reconcile Command, Delivery Retry,
   Run integration, report naming or any existing test, with one exception:
   `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` in
   `internal/daemon/task_engine_test.go` MUST be updated to expect the new
   `prior_report` `daemon.qa` event of Requirement 6 in its QA event list and
   positions, changing nothing else in that test. Its five subtests MUST pass.
8. MUST put the new tests in the two test files this Task creates, over
   temporary Git repositories with a side branch named like a Run Branch. The
   end-to-end test drives `TaskCycle` through the existing task-cycle fixture.

## Subtasks

- [ ] Find the newest unintegrated QA Report commit of the Spec.
- [ ] Import its QA files, or refuse with a reason and leave nothing behind.
- [ ] Wire the import and its event into the QA stage and the prompt context.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] The newest QA Report commit of the Spec on a side branch, outside
      `HEAD`'s history, is imported: its report and evidence are on disk with
      the blob bytes.
- [ ] A QA Report commit already in `HEAD`'s history imports nothing, and the
      event outcome is `none`.
- [ ] A Task commit touching the QA directory, and another Spec's QA Report
      commit, are never imported.
- [ ] A path that already holds different bytes, a report not newer than the
      tree's newest, and a report the shape detector refuses are each refused
      with their reason, and the QA directory is left as it was.
- [ ] `ReportShapeFindings` returns the report-shape finding for a report with
      a pending row and none for a closed report.
- [ ] End to end: a first pass whose QA Report commit with a recorded snapshot
      lies only on a side branch is followed by a head that re-commits its
      Task. The second pass's seeded report carries the unmoved passing row,
      records `re-run: not pass` for the failed row, and its QA Report commit
      includes the imported report.
- [ ] The QA prompt names the imported report and its pass's head.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_context.go`
- interface: `internal/speccheck/mechanical.go`
- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/qa_prior_pass.go`
- creates: `internal/daemon/qa_prior_pass_test.go`
- creates: `internal/speccheck/report_shape_findings_test.go`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- instruction: `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAPriorPassIsImportedOnlyFromARecordedRunCommit|TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte|TestPriorQAPassIgnoresReportsAlreadyInHistory|TestPriorQAPassIgnoresTaskCommitsAndOtherSpecs|TestPriorQAPassRefusesADifferingPath|TestPriorQAPassRefusesAReportOlderThanTheTreesNewest|TestPriorQAPassRefusesAReportTheShapeDetectorRefuses|TestQAGateCarriesARowFromAnUnintegratedFailedPass|TestQAPromptNamesTheImportedPassHead)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAPriorPassIsImportedOnlyFromARecordedRunCommit TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte TestPriorQAPassIgnoresReportsAlreadyInHistory TestPriorQAPassIgnoresTaskCommitsAndOtherSpecs TestPriorQAPassRefusesADifferingPath TestPriorQAPassRefusesAReportOlderThanTheTreesNewest TestPriorQAPassRefusesAReportTheShapeDetectorRefuses TestQAGateCarriesARowFromAnUnintegratedFailedPass TestQAPromptNamesTheImportedPassHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the eight named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestReportShapeFindingsNamesAPendingRow|TestReportShapeFindingsAcceptsAClosedReport)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReportShapeFindingsNamesAPendingRow TestReportShapeFindingsAcceptsAClosedReport; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither named test exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; User Story 2; Core Feature 2; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Importing the failed pass; API Contract 4; API Contract 5; Integration Points; Testing Approach 3; Build Order 3
- ADR-0194; ADR-0170; ADR-0053; ADR-0096

## Result

Implemented the prior-pass import in the QA stage, between the before snapshot
and prompt context. Selection checks the recorded Run and its journal's exact
QA settlement SHA on its Run Branch, then orders qualifying commits across
refs by commit date. QA settlement events now record the SHA and QA Task ID
needed by that selection. Older journal entries without this provenance are
skipped. The import preserves blob bytes, creates files with mode `0644`, keeps
identical existing files, refuses conflicts, old/future reports and existing
shape/evidence-path findings, and rolls back newly created files and directories
on refusal or error. Imported paths are explicitly included in the next QA
Report commit, and prompt context uses the imported commit's first parent.

Focused evidence for the acceptance criteria:

| Criterion | Implementation and focused-check evidence |
| --- | --- |
| Newest unintegrated report and evidence imported byte for byte | `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte` checks selection, blob bytes, mode and the `imported` event; `TestPriorQAPassKeepsIdenticalExistingEvidence` checks preservation of an existing identical file. |
| Report already in HEAD imports nothing | `TestPriorQAPassIgnoresReportsAlreadyInHistory` checks the `none` event. |
| Task and other-Spec commits ignored | `TestPriorQAPassIgnoresTaskCommitsAndOtherSpecs` exercises both exclusions; `TestAPriorPassIsImportedOnlyFromARecordedRunCommit` rejects an imitating commit on another ref without journal provenance. `TestPriorQAPassSkipsAPrunedRunBranch` proves a surviving derived ref cannot replace the recorded Run Branch. |
| Conflicting, old and malformed reports refused without residual files | `TestPriorQAPassRefusesADifferingPath`, `TestPriorQAPassRefusesAReportOlderThanTheTreesNewest` and `TestPriorQAPassRefusesAReportTheShapeDetectorRefuses` check reasons and preserved/absent files. `TestPriorQAPassRefusesAReportDatedAfterToday` covers the fourth refusal. |
| Existing shape detectors distinguish pending and closed reports | `TestReportShapeFindingsNamesAPendingRow` and `TestReportShapeFindingsAcceptsAClosedReport` exercise the exported loader/detector wrapper. |
| Next TaskCycle carries passing row from a failed side-branch pass | `TestQAGateCarriesARowFromAnUnintegratedFailedPass` uses a recorded failed Run, a distinct later Run and a re-committed Task; checks the seeded carried row, `re-run: not pass`, and imported report/evidence blobs in the new QA commit. |
| Prompt names imported report and audited head | `TestQAPromptNamesTheImportedPassHead` checks context; the TaskCycle test checks the Agent prompt. |

Focused command (exit 0 after the final implementation edits):

```text
GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/daemon ./internal/speccheck -run 'TestPriorQAPass|TestAPriorPass|TestQAPromptNames|TestQAGateCarriesARow|TestReportShapeFindings|TestTaskCycleQAVerdictMatrix' -count=1
```

Both packages passed, including all five existing QA verdict-matrix subtests.
`TestPriorQAPassCancellationPublishesStop` and
`TestPriorQAPassGitErrorIsInfrastructure` also passed in that focused run.
`git -c core.fsmonitor=false diff --check` exited 0.

The first incremental check was invalidated by concurrent edits; it also
encountered sandbox process-table restrictions in CLI force-stop tests and a
budget-test timing failure under concurrent test load. The stable-tree rerun passed
with process-table access, as recorded below. The initial default-cache check was blocked by
sandbox access to the host Go cache; subsequent checks use the task-scoped cache.

Task status and the Task Graph were left under Daemon ownership. The declared
Verification commands were not run. No commit, push or Pull Request was made.

Incremental check (exit 0 on the stable implementation tree):

```text
GOCACHE=/private/tmp/roundfix-task03-gocache rtk make verify-incremental
```

This rerun used sandbox escalation for the CLI force-stop fixtures' process-table
access. Formatting, vet, the repository test suite, skill sync/checks and build
passed. The CLI force-stop tests and the budget test passed on this rerun.

## Carry-forward provenance

- Source Run: `run_20261001T124541Z_e71bc4409b20b125`
- Source commit: `b18eeacac35c32d619e513812f413b18768433ca`
