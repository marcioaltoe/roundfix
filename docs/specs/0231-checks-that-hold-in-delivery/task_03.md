---
task: task_03
spec: 0231-checks-that-hold-in-delivery
status: pending
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
