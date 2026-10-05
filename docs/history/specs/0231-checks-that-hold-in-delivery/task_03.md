---
task: task_03
spec: 0231-checks-that-hold-in-delivery
status: completed
type: backend
complexity: medium
---

# Task 03: The Delivery Queue re-runs a failed check that tested an older default branch instead of parking it

## Overview

The finding of 2026-10-05 ([a Pull Request check ran on a stale merge](references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md))
is that the queue parked Pull Request #391 as `flaky-check` and then
`checks-failed` for failures that tested the default branch before its fix.
This Task makes the queue read the failing job's `tested-base` annotation,
compare it with the delivery remote's default branch tip, and re-run a stale
failure instead of classifying it (ADR-0236), and describes the rule in the
deliver guide.

## Requirements

1. MUST answer the finding of 2026-10-05 named in the Overview by adding
   `TestedBase`, `DefaultTip` and `Stale` to `CheckFailure` and the behavior
   of TechSpec Invariants 1 to 5 to `InspectFailedCheck` and
   `RerunFailedCheck`, with the commands of API Contracts 2 and 3: the job id
   comes from the check link's `/job/<id>` segment, and the annotation read,
   the default branch fetch and the staleness decision come before the
   failed log is read.
2. MUST make `checkCandidate` follow Invariants 6 to 8: a stale failure is
   re-run once per (run, default tip), logged with API Contract 4, restarts
   the check timeout at most once in total and keeps polling, and is never
   parked as `flaky-check` or `checks-failed` unless its re-run or inspection
   errs; the first newer-attempt failure of a run re-run as stale on the
   current tip is eligible for the one outside-change re-run; and a failure
   whose attempt is not newer than the attempt re-run for it keeps polling.
3. MUST keep every existing blocker, Park Class, log line and the
   `flaky-check: <check> passed on re-run` Warning (API Contract 5); a pass
   after a stale re-run adds no Warning; a failure without a tested base is
   classified exactly as before.
4. MUST update the scripted command sequences of the existing
   `InspectFailedCheck` tests only where the annotation read is added, and
   MUST add, in a new test file, the tests named in Verification: the replay
   of #391 (a retry whose attempt-2 failure tested the old tip is re-run and
   the item reaches merging with its Warning unchanged), polling while GitHub
   still reports the re-run attempt, the outside-change re-run after a stale
   re-run ending in `flaky-check`, a refused stale re-run parking
   `checks-failed`, a real-git inspection that calls an older tip and an
   unknown commit stale and the current tip not stale, annotations that are
   absent, malformed or ambiguous leaving the inspection as before, and the
   adapter re-running a stale failure past its first attempt.
5. MUST describe in `docs/user-guide/commands/deliver.md`, beside the
   existing re-run rule, the `tested-base` annotation, the stale rule and its
   log line `roundfix: check stale: Delivery Queue item <slug>: <check> tested <sha>, default branch is at <sha>; re-run (run <id>)`,
   and qualify the sentence on runs past attempt one accordingly.
6. MUST NOT reach GitHub from any test, persist stale re-runs, close or
   reopen Pull Requests, push to an item branch, or change the workflow, the
   Run Database schema or any skill.

## Subtasks

- [ ] Read the tested base and the default tip in `InspectFailedCheck`.
- [ ] Relax `RerunFailedCheck` to stale or outside-change failures.
- [ ] Re-run stale failures in `checkCandidate`.
- [ ] Update the scripted sequences and add the new tests.
- [ ] Describe the rule in the deliver guide.

## Acceptance Criteria

- [ ] A failure that tested an older default branch is re-run and polled, and
      the item merges when the re-run passes.
- [ ] A failure that tested the current tip, or that has no `tested-base`
      annotation, is classified as before.
- [ ] The deliver guide states the annotation, the stale rule and its log
      line.

## Context

- instruction: `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- interface: `internal/delivery/check_rerun.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/check_rerun_test.go`
- interface: `internal/delivery/check_rerun_remote_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- creates: `internal/delivery/stale_check_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestARetryAfterTheDefaultBranchMovedRerunsTheStaleCheckAndMerges|TestAStaleCheckKeepsPollingWhileGitHubReportsTheAttemptItReran|TestAStaleRerunThatFailsOnTheCurrentTipGetsTheOutsideChangeRerun|TestAStaleCheckWhoseRerunIsRefusedParksChecksFailed|TestGitHubCLIJudgesAFailedCheckByTheTipItTested|TestGitHubCLIWithoutATestedBaseInspectsAsBefore|TestGitHubCLIRerunsAStaleFailurePastItsFirstAttempt|TestGitHubCLIAttributionUsesTheRefreshedMergeBaseWithRealGit|TestGitHubCLIInspectsAgainstTheDeliveryRemote|TestARerunCheckThatFailsAgainParksAsFlakyCheck|TestACheckPastItsFirstAttemptIsNotRerun)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARetryAfterTheDefaultBranchMovedRerunsTheStaleCheckAndMerges TestAStaleCheckKeepsPollingWhileGitHubReportsTheAttemptItReran TestAStaleRerunThatFailsOnTheCurrentTipGetsTheOutsideChangeRerun TestAStaleCheckWhoseRerunIsRefusedParksChecksFailed TestGitHubCLIJudgesAFailedCheckByTheTipItTested TestGitHubCLIWithoutATestedBaseInspectsAsBefore TestGitHubCLIRerunsAStaleFailurePastItsFirstAttempt TestGitHubCLIAttributionUsesTheRefreshedMergeBaseWithRealGit TestGitHubCLIInspectsAgainstTheDeliveryRemote TestARerunCheckThatFailsAgainParksAsFlakyCheck TestACheckPastItsFirstAttemptIsNotRerun; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven new tests do not exist, so their pass lines are missing and the command fails.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'roundfix: check stale: Delivery Queue item <slug>: <check> tested <sha>, default branch is at <sha>; re-run (run <id>)' && grep -qF -- 'tested-base' docs/user-guide/commands/deliver.md && grep -qF -- 'check stale: Delivery Queue item' internal/delivery/engine.go && grep -qF -- '"tested-base"' internal/delivery/check_rerun.go` — expected: exit 0; before this Task neither the guide nor the code names the stale rule, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 3; Success Metric 4; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 1 to 8; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Vocabulary Contract; Testing Approach; Build Order 3
- ADR-0236

## Result

Implemented the Task 03 slice for Daemon Verification; status and the authored
Verification commands remain Daemon-owned.

- `InspectFailedCheck` now extracts the job id from the check link and reads
  its annotations after the attempt. Exactly one `tested-base` annotation with
  a 40-character lowercase hexadecimal message is accepted. With that evidence,
  inspection fetches the delivery remote's default branch, resolves its tip,
  checks the tested commit and its ancestry, and returns a stale failure before
  reading the failed log. Without usable evidence, attribution retains its
  previous command order and behavior. A current or descendant tested base
  proceeds through the existing package attribution using the refreshed ref.
- `RerunFailedCheck` accepts stale or outside-change failures with valid run
  ids; the engine owns attempt limits. Check polling remembers stale (run, tip)
  pairs, runs re-run as stale, and the attempt re-run per run only in memory.
  Stale failures bypass classification and add the documented owner log line.
  Delayed attempt reports keep polling; a newer current-tip failure can receive
  the existing one outside-change re-run. Both recovery paths share the one
  timeout restart. Existing blockers, Park Classes, recovery logs and Warning
  behavior remain in place.
- The deliver guide explains the annotation, stale rule, exact log line,
  timeout bound, in-memory lifetime and eligibility past attempt one.
- Existing scripted inspection sequences changed only by adding the empty
  annotation response. The new `stale_check_test.go` uses the existing queue
  and command-runner seams; GitHub commands are always scripted and local Git
  fixtures use disposable repositories.

Acceptance evidence:

1. Older-base failures re-run and keep polling: the #391 replay
   `TestARetryAfterTheDefaultBranchMovedRerunsTheStaleCheckAndMerges` reaches
   `merging` from an attempt-2 failure, preserving the existing Warning and
   asserting the exact stale log. The delayed-attempt test proves both stale
   and non-stale reports of the re-run attempt keep polling. Additional tests
   prove once-per-tip re-runs, only one timeout restart, refusal parking as
   `checks-failed`, and the newer-attempt outside-change path ending in
   `flaky-check` or preserving its existing pass Warning.
2. Current-base and unannotated failures retain classification: real Git
   inspection calls older and unknown commits stale, current and descendant
   commits current, and reads the failed log only for current evidence.
   Named subtests cover absent, unrelated, malformed, uppercase and ambiguous
   annotations, including a valid annotation alongside a malformed duplicate.
   Existing attribution and engine re-run tests pass in the focused selection.
   Additional tests cover links without a job and inspection command failures.
3. Guide contract: a focused Python assertion confirmed `tested-base`, the
   exact stale log line, and both later-attempt qualifications are present.

Focused checks (after the last code edit):

- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 ./internal/delivery -run 'Stale|TestGitHubCLI|TestA.*Check|TestCheck|TestSecondCheck|TestMultipleFailed|TestSkippedRerun|TestARetryAfter|TestActionsLink'`
  — exit 0, `ok roundfix/internal/delivery` (3.292s).
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.
- Focused Python guide and Task-status assertions — exit 0; the annotation,
  log line and eligibility prose are present, and status remains `in_progress`.

Initial inspection showed no tested-base fields or annotation read, and
classification allowed only attempt-one recovery. An initial focused test run
caught an omitted existing re-run log line in the implementation edit; the line
was restored and the final focused selection above passed.

The declared Verification commands and repository-wide gate were not run in
this child turn. No commits, pushes, Pull Request operations, workflow edits,
Database schema changes, skill edits, or edits to other Task files or the Task
Graph were made. No follow-up scope was discovered.
