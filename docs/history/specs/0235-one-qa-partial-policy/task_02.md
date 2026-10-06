---
task: task_02
spec: 0235-one-qa-partial-policy
status: completed
type: backend
complexity: medium
---

# Task 02: Settlement, archive, settle and the queue agree on every partial, and settle names the report it refused

## Overview

With task_01's policy in place, this Task proves that the four QA settlement
callers reach the same verdict on the same report. It also removes the two
places where the outcome still diverged in practice. The Delivery Queue
stops classifying with a predicate of its own, and `roundfix settle` names
the report it judged. The Fluxus case of the Backlog Entry of 2026-10-06
([settlement refuses a partial that archive accepts](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md))
could not be diagnosed, because settle never said which report it had
refused.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-06 named in the Overview by
   making `qaEnvironmentPartial` in `internal/cli/deliver_workflow.go`
   classify per Invariant 8: `partial`, `rows_blocked_finding` 0 and
   `EnvironmentRowsNeedingOverride() > 0`, read from the newest report on
   the Run Branch with the shared reader, as today.
2. MUST make `roundfix settle`'s QA eligibility refusal print API Contract 4,
   `roundfix: settle QA Report is ineligible: <reason> (report <path>)`, with
   the path of the report it read relative to the settle surface's working
   tree. Everything else about settle (exit 1, Task file unchanged, the
   `Settle surface:` line) stays as it is.
3. MUST NOT add any condition to Daemon settlement, `roundfix archive` or
   `roundfix qa-report accept`. They keep calling `spec.QAReportEligibility`
   (Invariant 10).
4. MUST add `internal/cli/qa_partial_policy_test.go` with the tests below,
   using temporary repositories and homes, the existing archive, settle and
   delivery recovery harnesses, and no network:
   - `TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape`. For each
     shape of task_01's Requirement 4, `qa-report accept`, `archive` and the
     QA Task's `settle` all accept or all refuse. A refusal names the same
     reason text in each command's stderr. The shapes include Surface
     Transcript 1's (exit 0, empty stdout and stderr) and Surface
     Transcript 2's (exit 1, the quoted stderr line).
   - `TestSettleNamesTheReportItRefused`. It matches Surface Transcript 3,
     and the named path is the newest report of the settle surface.
   - `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly`. A report
     with only Pull Request rows, and one with a Pull Request row and a
     network-denied outside-evidence row, are not classified
     `qa-environment-partial`. Reports with another environment row, or
     with a network-denied row without outside-evidence provenance, still
     are.
5. MUST add `internal/daemon/qa_partial_policy_test.go` with
   `TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo`. It drives
   `Engine.TaskCycle` with the existing `taskFakeRunner` over the same
   shapes and asserts that the QA Task settles `completed` exactly for the
   shapes the commands accept, with the refusal reason
   `QA verdict partial not accepted: <reason>` otherwise.
6. MUST make the declared changes to existing tests, and no others. In
   `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch`, the case
   "pre-PR only" now expects `false`.
   `TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly`
   is removed, because its premise is reversed: its cases move into
   `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly` with the new
   expectations. Settle tests that match the refusal by substring keep
   passing unchanged.
7. MUST describe the settle refusal and the report it names in
   `docs/user-guide/commands/settle.md`, quoting
   `(report <path>)` with its placeholder.
8. MUST NOT change any skill, Baseline asset, `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Classify `qa-environment-partial` through the shared override count.
- [ ] Name the report in settle's eligibility refusal and in the settle guide.
- [ ] Add the agreement tests across the commands and the Daemon.
- [ ] Update the two reversed queue expectations.

## Acceptance Criteria

- [ ] For every shape, `qa-report accept`, `archive`, `settle` and Daemon
      settlement reach the same verdict, and every refusal carries the same
      reason.
- [ ] A qualifying partial is never classified `qa-environment-partial`,
      and a partial with an environment row that needs the override still
      is.
- [ ] A refused settle names the report it judged.

## Context

- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`
- interface: `internal/spec/qa.go`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_operator_archive_test.go`
- interface: `internal/cli/settle_test.go`
- interface: `internal/cli/archive_test.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `docs/user-guide/commands/settle.md`
- creates: `internal/cli/qa_partial_policy_test.go`
- creates: `internal/daemon/qa_partial_policy_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape|TestSettleNamesTheReportItRefused|TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly|TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch|TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo|TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport)$' ./internal/cli ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape TestSettleNamesTheReportItRefused TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; ! grep -q 'func TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly' internal/cli/deliver_operator_archive_test.go || { printf 'reversed test still present\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/settle.md | grep -qF -- '(report <path>)' || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/settle.md '(report <path>)' >&2; exit 1; }` — expected: exit 0; before this Task the four new tests do not exist, the reversed test is still present and the settle guide does not name the report, so the command fails; after it the new and changed tests pass.

## References

- `_prd.md` → Goals; Core Features 4-5; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Invariants 8 to 10; API Contract 1; API Contract 4; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Testing Approach; Build Order 2
- ADR-0240; ADR-0229; ADR-0154

## Result

Implemented this Task's caller changes. The Delivery Queue reads the newest
Run Branch report through the existing shared reader and classifies
`qa-environment-partial` only when the partial has no finding-blocked row and
`EnvironmentRowsNeedingOverride() > 0`. Settle selects the newest report once,
reads that file through the shared reader, and appends its path relative to
the selected working tree to the shared eligibility refusal. The settle guide
quotes `(report <path>)` and documents exit 1 and the unchanged Task file.
Daemon settlement, archive and `qa-report accept` retain their existing policy
calls with no added conditions.

Acceptance evidence:

- Agreement: `TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape`
  exercises 16 report shapes independently through all three commands, and
  `TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo` exercises the same
  shapes through `Engine.TaskCycle` and `taskFakeRunner`. Both use explicit
  expected verdicts and refusal reasons. The command test matches empty
  stdout/stderr on accepted `qa-report accept` reports and the exact refusal
  line on rejected reports. The Daemon test matches completed/failed QA Task
  status and `QA verdict partial not accepted: <reason>` for refusals.
- Queue classification: `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly`
  reads committed newest reports from the Run Branch while the checkout
  remains at the item head. Pull Request rows only and Pull Request plus
  network-denied outside evidence are not environment-only parks; another
  environment row and network denial without outside-evidence provenance
  remain parks. The removed test's three original fixtures are retained in
  this replacement. The existing environment-only matrix changes only the
  specified pre-PR-only expectation.
- Refused report identity: `TestSettleNamesTheReportItRefused` selects a kept
  Task Worktree with two reports while the checkout holds different eligible
  evidence. It matches exit 1 and the complete stderr transcript, including
  `Settle surface:` and the newest report's relative path, and proves the
  Task file is unchanged. The agreement test checks this refusal shape for
  every rejected partial. Existing settle eligibility tests remain unchanged.

Focused checks:

- Before the caller edits,
  `rtk proxy env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/cli -run '^(TestSettleNamesTheReportItRefused|TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly)$'`
  exited 1: settle omitted the report path, and both exempt-only queue cases
  incorrectly classified as environment partials.
- After the caller edits and fixture correction,
  `rtk proxy env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/cli ./internal/daemon -run 'TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape|TestSettleNamesTheReportItRefused|TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly|TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch|TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo|TestSettleAppliesEligibilityToAQATask'`
  exited 0 for both packages. An earlier agreement run exposed that the test
  harness's selected working directory does not change the process directory
  for `qa-report accept`; its fixture now passes the temporary absolute path.
- `rtk proxy env GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  initially exited 2 because the sandbox denied process-table access in
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` and
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`.
  After preserving the moved queue fixtures, the same command with elevated
  process access exited 0: formatting, vet, the Go suite, skill checks and
  build passed. Other unchanged package results reused the incremental cache;
  the final CLI package run passed in 237.651 seconds.

The authored Verification command was not run. Task status remains
Daemon-owned; no Task Graph, other Task file, skill, Baseline asset,
`CONTEXT.md` or `CHANGELOG.md` was edited. No commit, push or Pull Request was
created. No follow-up implementation outside this slice was added.

## Carry-forward provenance

- Source Run: `run_20261006T120012Z_3cae1c6ab9ca3edb`
- Source commit: `e3922de6a1294c42e235327fec4e081c4bf8fac2`
