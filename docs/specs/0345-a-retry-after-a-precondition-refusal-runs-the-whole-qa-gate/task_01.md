---
task: task_01
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
status: completed
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


## Result

Skills loaded: implement-task, coding-guidelines, systematic-debugging, no-workarounds, golang-testing, testing-boss, tech-writer, evidence-gate.

Implementation: `copyPriorQAPass` now reads every `qa-report-*.md` in
`pass.Files` order after copying and before shape checking. The first readable
report with `RowsBlockedPrecondition > 0` returns
`precondition refusal: <path>` using the listed path. Existing deferred cleanup
removes the imported files and directories. No other production function changed;
candidate selection and the format path retain their bytes. Existing tests,
helpers and fixtures retain their original bytes; the three new top-level tests
are appended, each beginning with `t.Parallel()`.

Starting-tree characterization, before changing production code or permanent
tests:

- A temporary overlay test, `TestTask0345StartingTreeCharacterization`, used
  `priorFixture`, a recorded Run Branch commit whose report was written by
  `spec.WritePreconditionRefusalReport`, `TaskCycle`, `ExecVerifier`,
  `GitCommitter`, and repository Verification
  `! grep -rqF -- '| - | --- |' docs/specs/<slug>/qa`.
- Command: `go test -overlay <temporary overlay.json> -count=1 -timeout 3m -v -run '^TestTask0345StartingTreeCharacterization$' ./internal/daemon`.
  Exit 0: the characterization asserted `prior_report` outcome `imported`,
  zero Agent prompts, verdict `fail`, accepted false, and the imported refusal
  byte-identical in the QA Report commit at `HEAD`, beside the new refusal.
  The overlay and test exist only outside the repository and are not part of
  this diff.
- The permanent new tests were then run before adding the guard: all three
  import-refusal cases failed with `refusal import = true <nil>`; the retry
  failed with `Agent prompts = 0, want 1; verdict = fail`; both format cases
  passed, confirming that the existing format path needed no production edit.

Focused evidence for the acceptance criteria:

- Command: `go test -count=1 -timeout 3m -v -run '^(TestPriorQAPassRefusesAPreconditionRefusal|TestARetryAfterAPreconditionRefusalRunsTheWholeGate|TestQAFormatFormatsAPreconditionRefusalReport)$' ./internal/daemon`.
  Exit 0 after the guard. The refusal table covers repository Verification,
  strict Spec check and an older refusal beside a newer matrix report; each
  asserts the exact refusal reason, null imported files and no remaining
  `qa/` directory. The retry asserts a refused import, a repository Verification
  pass, exactly one Agent prompt, verdict `pass` and accepted, the fixture's QA
  Task settled `completed`, and no refusal file in the worktree or the QA Report
  commit. The format cases use `exit 3`, `ExecVerifier` and `GitCommitter`;
  they assert one precondition row and exact committed bytes against the
  Daemon's report captured before the formatter: appended bytes on success,
  restored original bytes on failure, plus the `commit` format outcome and
  failure diagnostic.
- Command: `go test -count=1 -timeout 3m -run '^(TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte|TestQAGateCarriesARowFromAnUnintegratedFailedPass|TestPriorQAPassRefusesAReportTheShapeDetectorRefuses|TestPriorQAPassKeepsIdenticalExistingEvidence|TestQAFormatFormatsAnImportedPassBeforeThePrecondition|TestQAFormatRevertsAFailingFormatter)$' ./internal/daemon`.
  Exit 0: the existing failed-matrix import remains byte-identical, its row
  carries, existing cleanup and evidence preservation hold, and imported
  formatting and failed-formatter restoration retain their behavior.
- `GOOS=windows go vet ./internal/daemon`: exit 0.
- `git diff --check`: exit 0. A byte-prefix comparison with `HEAD` confirms
  both changed test files preserve every existing byte. The changed-path audit
  found only the three declared Go files and this assigned Task file.

Required sabotages ran sequentially on an independent temporary copy with its
own initialized/staged Git metadata. Each used
`go test -count=1 -timeout 3m -v -run <named tests> ./internal/daemon` and exited 1
with the following inspected assertions; the sabotage driver exited 0 only
when every expected named failure appeared:

| Sabotage | Named failing test and case | Observed failure |
| --- | --- | --- |
| Remove refusal check | `TestPriorQAPassRefusesAPreconditionRefusal`: `repository Verification refusal`, `strict Spec check refusal`, `older refusal beside a matrix report`; `TestARetryAfterAPreconditionRefusalRunsTheWholeGate` | All three import cases returned imported true; retry received zero Agent prompts and verdict fail. |
| Check only newest report | `TestPriorQAPassRefusesAPreconditionRefusal/older refusal beside a matrix report` | Imported true; the two single-refusal cases still passed. |
| Refuse every fail verdict | `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte`; `TestQAGateCarriesARowFromAnUnintegratedFailedPass` | Matrix import returned false; the next seed lacked the required carry dispositions. |
| Skip format when newest report records refusal | `TestQAFormatFormatsAPreconditionRefusalReport/formatter runs` | Committed refusal lacked the formatter's bytes. The `formatter fails` case also failed because the formatter never captured the original report. |

Full sabotage transcripts are temporary execution artifacts, outside the diff,
in `roundfix-0345-sabotage-hw3gfdro/{removed,newest-only,all-fail,skip-refusal-format}.log`
under the execution host's temporary directory. No sabotage changed this Run's
production format code or any other excluded path.

Declared Verification and repository Verification were not run; the Daemon owns
them and Task settlement. Task status, `_tasks.md`, other Tasks and all excluded
paths were left untouched by this Agent. No commit, push or pull request was made.
No follow-up outside this Task's slice was identified.

## Carry-forward provenance

- Source Run: `run_20261104T103445Z_d9889df3a60eb797`
- Source commit: `420537ce292627472bb807afeef1b597e7a93f17`
