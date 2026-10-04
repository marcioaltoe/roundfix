---
task: task_02
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: completed
type: backend
complexity: low
---

# Task 02: An environment-only QA partial parks for the operator even when only Pull Request rows block it

## Overview

Spec 0220's QA Report of 2026-10-03 was `partial` with no finding-blocked row
and two environment-blocked rows that both waited for an open Pull Request,
so the queue's test of environment rows against Pull Request rows (2 against
2) failed and the item parked `run-unresolved` (operator intervention log
entry 141 of 2026-10-04). This Task parks every zero-finding partial with at
least one environment-blocked row as `qa-environment-partial` (ADR-0229),
whose status line names the override recovery.

## Requirements

1. MUST implement "The environment-only partial" of the TechSpec in
   `commandDeliveryWorkflow.qaEnvironmentPartial`: the flag is true when the
   newest Run-Branch QA Report is `partial`, `rows_blocked_finding` is `0`
   and `rows_blocked_environment` is at least `1`, as API Contract 2 states.
2. MUST keep `run-unresolved` for a missing or unreadable report, a
   finding-blocked partial, a partial without an environment-blocked row, and
   every other verdict.
3. MUST NOT change `spec.QAReport`, `RowsBlockedPrePullRequest`, the Archive
   Command's eligibility (ADR-0167), blocker names or Park Classes.
4. MUST add `TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly`
   to `internal/cli/deliver_operator_archive_test.go`, through
   `commandDeliveryWorkflow.runResult` on a disposable Run Branch, with the
   subtests `pull request rows only` (a report written inline with two rows of
   status `blocked (environment: no open Pull Request)` and provenance
   `Pull Request row`, `rows_blocked_environment: 2`, expecting the flag),
   `declared rows only` and `finding beside pull request rows` (expecting
   none), as Testing Approach 3 describes. No test reads an archived Spec.
5. MUST change the expectation of the `pre-PR only` case of
   `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch` to `true`,
   the one declared break, and leave its other cases unchanged.

## Subtasks

- [ ] Lower the environment-partial predicate to at least one environment-blocked row.
- [ ] Add the 0220-shaped classification test with its two negative cases.
- [ ] Update the declared `pre-PR only` expectation.

## Acceptance Criteria

- [ ] A zero-finding partial blocked only by rows waiting for an open Pull
      Request sets the environment partial; a declared-only partial and a
      partial with a finding row do not.
- [ ] The engine still parks the flag as `qa-environment-partial` and its
      absence as `run-unresolved`.

## Context

- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_operator_archive_test.go`
- instruction: `internal/spec/qa.go`
- instruction: `internal/delivery/operator_archive_retry_test.go`
- instruction: `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly|TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch|TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly/pull_request_rows_only TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly/declared_rows_only TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly/finding_beside_pull_request_rows TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch/pre-PR_only TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch/finding TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new test does not exist and the pre-PR only case expects false, so the command fails.

## References

- `_prd.md` → Goals 2 and 3; Core Feature 2; Success Metric 2; Success Metric 3
- `_techspec.md` → The environment-only partial; Interfaces; API Contract 2; Testing Approach 3; Build Order 2
- ADR-0229; ADR-0167

## Result

The Run-Branch QA predicate now requires a readable newest report with verdict
`partial`, zero finding-blocked rows and at least one environment-blocked row.
It no longer subtracts pre-PR Pull Request rows. The QA reader, archive
eligibility, blocker names and Park Classes are unchanged.

Added the inline two-Pull-Request-row fixture through `runResult` in a
disposable repository: the report is committed only to its Run Branch, then
the item checkout returns to its original head before classification. The
declared-only and finding-beside-Pull-Request cases use the same boundary.
Only the existing `pre-PR only` expectation changed; its other cases retain
their fixtures and expectations. No test reads an archived Spec.

Focused evidence:

- Before the predicate change,
  `GOCACHE=/private/tmp/roundfix-0224-task02-gocache rtk proxy go test -count=1 -run '^TestRunSpecReportsAPartialBlockedOnlyByPullRequestRowsAsEnvironmentOnly$' ./internal/cli`
  exited 1: `pull request rows only` returned
  `QAEnvironmentPartial:false` instead of true. Both negative cases passed.
- After the change,
  `GOCACHE=/private/tmp/roundfix-0224-task02-gocache rtk proxy go test -count=1 -v -run '^TestRunSpecReports' ./internal/cli`
  exited 0 with all ten subtests passing. Acceptance criterion 1 is covered by
  `pull request rows only`, `declared rows only` and
  `finding beside pull request rows`; the existing cases also cover missing
  and unreadable reports, findings, the newest report, and the pass verdict.
- `GOCACHE=/private/tmp/roundfix-0224-task02-gocache rtk proxy go test -count=1 -v -run '^TestAnEnvironmentOnlyPartialParksAsQAEnvironmentPartial$' ./internal/delivery`
  exited 0. Acceptance criterion 2 is covered by `environment` and
  `unresolved`: the unchanged engine parks the flag as
  `qa-environment-partial` and its absence as `run-unresolved`. The environment
  case also checks the Park Class and override recovery text.

- `GOCACHE=/private/tmp/roundfix-0224-task02-gocache rtk make verify-incremental`
  exited 2 in the sandbox: owner-process tests could not read the process
  table and HTTP fixtures could not bind local listeners. Suite guards also
  detected this Result section being appended while the suite was running.
  These diagnostics do not replace the focused evidence.
- The permission-enabled rerun,
  `GOCACHE=/private/tmp/roundfix-0224-task02-gocache rtk proxy make verify-incremental > /private/tmp/roundfix-0224-task02-incremental.log 2>&1`,
  exited 0 with the worktree unchanged throughout the check. Formatting,
  Go vet, package tests, skill synchronization/readiness checks and the CLI
  build passed. The log is local focused-check evidence, not Daemon
  Verification evidence.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0. Final scope
  inspection shows only the two assigned CLI paths and this Task file changed;
  the pre-existing Daemon status change is preserved.

The Task's declared Verification command remains for the Daemon. Status and
checkboxes are unchanged; no commit, push or Pull Request was made.

## Carry-forward provenance

- Source Run: `run_20261004T162613Z_0bb1b015eb38156b`
- Source commit: `0d99623956d9d0a5bed32bfd6e7cc47f27039bce`
