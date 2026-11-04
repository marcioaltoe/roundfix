---
task: task_01
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
status: pending
type: backend
complexity: medium
---

# Task 01: The QA Report Import refuses a Precondition Refusal, and its report's format path is pinned

## Overview

The QA Report Import (`copyPriorQAPass`, `internal/daemon/qa_prior_pass.go`)
copies the newest unintegrated Daemon QA Report commit of the Spec into a new
pass whatever its report records, so a retry after a precondition refusal
imports the refusal and, in a repository whose Verification checks the
formatting of `qa/` files, refuses again on its bytes without an Agent
session. After this Task the import refuses a candidate commit that adds or
modifies a report recording a Precondition Refusal, and the pass runs in full
with no previous report. A new test pins that the refusal report the Daemon
writes goes through the configured Format Command before its commit. This
answers the report reported by pantheon on 2026-11-03, "deliver retry
reimporta QA Report 'fail' de precondição e não roda o QA gate de novo"
([references](references/2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md)).

## Requirements

1. MUST, before changing code, characterize the starting tree with a
   temporary test that is not committed: a recorded Run Branch commit adding
   a report written by `spec.WritePreconditionRefusalReport`, then
   `TaskCycle` with the repository Verification
   `! grep -rqF -- '| - | --- |' docs/specs/<slug>/qa` and `ExecVerifier`,
   and record in the Result that the import publishes `imported`, no Agent
   prompt is sent, the verdict is `fail` and the QA Report commit holds the
   imported refusal.
2. MUST change `copyPriorQAPass` exactly as `_techspec.md` → API Contract 1
   states: after the copy loop and before `speccheck.ReportShapeFindings`,
   the first `qa-report-*.md` path of `pass.Files` that
   `spec.ReadQAReportFile` reads with `RowsBlockedPrecondition` greater than
   zero returns the reason `precondition refusal: <path>`, with the path as
   `pass.Files` lists it; the existing deferred cleanup removes what the
   import created. No other function changes, and `findRecordedPriorQAPass`
   searches no older candidate (API Contract 2).
3. MUST add `TestPriorQAPassRefusesAPreconditionRefusal` to
   `internal/daemon/qa_prior_pass_test.go`, a table with the three cases of
   `_techspec.md` → Testing Approach (`repository Verification refusal`,
   `strict Spec check refusal`, `older refusal beside a matrix report`), each
   asserting `imported` false, the last `prior_report` event `refused` with
   reason `precondition refusal: <refusal path>`, and no `qa/` directory
   under the Spec.
4. MUST add `TestARetryAfterAPreconditionRefusalRunsTheWholeGate` to
   `internal/daemon/qa_prior_pass_test.go`: the Pantheon shape of requirement
   1 with a fake Agent whose report passes; the import is `refused`, the
   repository Verification passes, the Agent is prompted once, the result is
   `pass` and accepted, the QA Task settles `completed`, and the QA Report
   commit at `HEAD` lists none of the refusal commit's files.
5. MUST add `TestQAFormatFormatsAPreconditionRefusalReport` to
   `internal/daemon/qa_format_test.go` with the subtests `formatter runs` and
   `formatter fails` of `_techspec.md` → Testing Approach, through `TaskCycle`
   with repository Verification `exit 3`, `ExecVerifier` and `GitCommitter`:
   the committed newest report has `rows_blocked_precondition` 1 and holds
   the formatter's bytes with event `commit` `formatted`, or the Daemon's
   bytes with event `commit` `failed` and the `roundfix: QA format commit
   failed:` diagnostic.
6. MUST make each new top-level test call `t.Parallel()` first, keep every
   test file compiling for Windows (`GOOS=windows go vet ./internal/daemon`),
   and keep every existing test, helper and fixture byte-identical.
7. MUST run, on a temporary copy, each sabotage `_techspec.md` → Testing
   Approach lists and record in the Result which named test fails and on
   which case.
8. MUST NOT touch `internal/daemon/task_engine.go`,
   `internal/daemon/qa_format.go`, `internal/spec`, `internal/speccheck`,
   `internal/cli`, any document, any skill, the Makefile, CI workflows,
   `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the import of a refusal on the starting tree.
- [ ] Refuse a candidate whose reports record a Precondition Refusal.
- [ ] Add the import, retry and format tests.
- [ ] Run the sabotages on a temporary copy.

## Acceptance Criteria

- [ ] A recorded refusal commit is refused with `precondition refusal: <path>` and leaves no file, and a retry after it reaches the Agent and commits none of its files.
- [ ] A failed pass that built its matrix is still imported byte for byte, and the refusal report is formatted before its commit.

## Context

- interface: `internal/daemon/qa_prior_pass.go`
- interface: `internal/daemon/qa_prior_pass_test.go`
- interface: `internal/daemon/qa_format_test.go`
- interface: `internal/daemon/qa_prior_pass_skip_test.go`
- interface: `internal/daemon/qa_two_pass_carry_test.go`
- interface: `internal/daemon/task_engine_test.go`

## Consumers

- The QA Report Import's `refused` outcome and the files a pass imports — consumers: `./internal/daemon`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestPriorQAPassRefusesAPreconditionRefusal|TestARetryAfterAPreconditionRefusalRunsTheWholeGate|TestQAFormatFormatsAPreconditionRefusalReport|TestPriorQAPass.*|TestAPriorPassIsImportedOnlyFromARecordedRunCommit|TestQAGateCarriesARowFromAnUnintegratedFailedPass|TestQAPromptNamesTheImportedPassHead|TestTwoGatePasses.*|TestQAFormat.*|TestQAGateRefusesOnFailedRepositoryVerification)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPriorQAPassRefusesAPreconditionRefusal TestARetryAfterAPreconditionRefusalRunsTheWholeGate TestQAFormatFormatsAPreconditionRefusalReport; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the three new tests exists, so the first pass line is missing; afterwards a recorded refusal is refused, a retry after it runs the whole gate, the refusal report is formatted, and the existing import, carry and format tests still pass.

## References

- `_prd.md` → Sources 1; Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 1-5; API Contract 1; API Contract 2; API Contract 3; Invariants 1-5; Declared breaks; Testing Approach; Build Order 1
