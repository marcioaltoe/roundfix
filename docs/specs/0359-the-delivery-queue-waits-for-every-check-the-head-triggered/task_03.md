---
task: task_03
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: completed
type: backend
complexity: high
---

# Task 03: The checking stage waits for the Head Check Set and parks naming each failed check

## Overview

The checking stage moved Fluxus pull requests #466, #469 and #477 to
`merging` on one passing check and the merge state `CLEAN` while two
workflow runs of each head were in progress, and parked a failure, when it
saw one, as the bare `checks-failed` (reported by fluxus on 2026-11-05,
[references](references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md)).
After this Task `checkCandidate` evaluates the Head Check Set task_02 reads,
waits while any member is pending or any Expected Workflow has no run, parks
`checks-failed` naming every failed member with its link, writes the log
lines, and names the pending members in the `checks-timeout` blocker, as
`_techspec.md` → API Contracts 3 to 5 state.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestTheCheckWait.*|TestEachCheckWaitStartsItsOwnDeadline|TestChecksTimeoutParksTheItem|TestABlockedMergeStateParksAtTheChecksTimeout)$' ./internal/delivery`
   and record in the Result that it passes.
2. MUST create `internal/delivery/head_check_wait.go` with `headCheckSet`,
   `evaluateHeadCheckSet` and `checksFailedBlocker` exactly as `_techspec.md` →
   Interfaces and API Contracts 3 and 4 state, and change `checkCandidate`
   and `checksTimeoutBlocker` in `internal/delivery/engine.go` as API
   Contracts 4 and 5 state, writing exactly these log lines:
   - `roundfix: checks: Delivery Queue item <slug>: waiting for <n> check(s): <member>; <member>`
   - `roundfix: checks: Delivery Queue item <slug>: <m> check(s) completed for <head>`
   - `roundfix: checks: Delivery Queue item <slug>: read failed: <error>`

   and keeping the `merge state` line, the re-run rules of ADR-0236, every
   poll interval and every deadline as they are.
3. MUST create `internal/delivery/head_check_wait_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, each
   top-level test calling `t.Parallel()` first, built on
   `newCheckWaitFixture` with the fake clock and a log buffer, with links of
   the form `https://github.com/example/repo/actions/runs/<run>/job/<job>`:
   - `TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered`, the
     Fluxus #466 replay, with three scripted polls, merge state `CLEAN`:
     (1) the check `Validate PR title` of `CI | Conventions` passing (run 1),
     and the runs 1 `CI | Conventions` `completed` `success`, 2 `CI |
     Frontend` `in_progress` and 3 `CI | Backend` `in_progress`; (2) the same
     plus `Lint, Type Check, Test & Build` of `CI | Frontend` (run 2) and of
     `CI | Backend` (run 3, job 30), both `pending`; (3) the Frontend check
     `pass`, run 2 `completed` `success`, the Backend check `fail` with state
     `FAILURE` and run 3 `completed` `failure`. With a check recovery that
     reports the failure in a changed package, the item parks exactly
     `checks-failed: CI | Backend / Lint, Type Check, Test & Build (failure) https://github.com/example/repo/actions/runs/3/job/30`,
     and the log holds, in order, the `merge state CLEAN` line, `waiting for
     2 check(s): CI | Frontend (in_progress); CI | Backend (in_progress)` and
     `waiting for 2 check(s): CI | Frontend / Lint, Type Check, Test & Build;
     CI | Backend / Lint, Type Check, Test & Build`. Subtest `every member
     passes`: in poll (3) both checks pass and both runs succeed; the item
     reaches `merging` and the log ends with `3 check(s) completed for
     candidate`.
   - `TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting`: a report with
     one passing commit status `coverage` (no workflow, no link), no run and
     the Expected Workflow `.github/workflows/ci-verify.yml`, with a
     `1m` Check Wait, parks exactly
     `checks-timeout: current-head checks not settled within delivery.check_timeout 1m0s; pending: .github/workflows/ci-verify.yml (no run yet)`.
     Subtest `the run arrives`: from the third poll the run of that path is
     `completed` `success`; the item reaches `merging`.
   - `TestAFailedWorkflowRunWithoutAFailingCheckParks`: a table over the
     conclusions `failure`, `timed_out`, `cancelled`, `action_required`,
     `startup_failure` and `stale` of run 7 `CI | Verify` with `html_url`
     `https://github.com/example/repo/actions/runs/7`, beside one passing
     check of another run, each parking exactly
     `checks-failed: CI | Verify (<conclusion>) https://github.com/example/repo/actions/runs/7`;
     and `neutral` and `skipped`, each reaching `merging`.
   - `TestTheChecksTimeoutNamesWhatItWaitedFor`: subtests `no check
     reported` (an empty report: `; pending: no check reported`), `merge
     state` (every member passed, merge state `BLOCKED`: `; merge state
     BLOCKED`) and `pending check` (`recoveryReport("pending")`: `; pending:
     CI / verify`).
   - `TestAReadErrorIsLoggedOnce`: two reads failing with the same error,
     then passing reports; the log holds exactly one `read failed:` line.
4. MUST update every test that pins the `checks-failed` or `checks-timeout`
   text, or an empty log after `checking`, on the default branch when the
   Task runs, changing only that expected text to what API Contracts 4 and 5
   produce for its report (`_techspec.md` → Declared breaks lists them at
   `6ceb7ff5`); no poll count, elapsed time, stage or re-run count changes.
5. MUST run sabotages 1, 2, 3 and 6 of `_techspec.md` → Testing Approach on
   a temporary copy and record which named test fails and on which input.
6. MUST NOT change `internal/delivery/github.go`,
   `internal/delivery/head_check_set.go`,
   `internal/delivery/turn.go`, `internal/delivery/park_class.go`, a blocker
   name, a Park Class, the Run Database, any document, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the check wait tests.
- [ ] Evaluate the Head Check Set and name failed members.
- [ ] Apply the completion rule, the log lines and the timeout text in `checkCandidate`.
- [ ] Add the replay tests, update the pinned texts and run the four sabotages.

## Acceptance Criteria

- [ ] An item never moves to `merging` while a member is pending, an Expected Workflow has no run or the set is empty.
- [ ] Every `checks-failed` park names each failed member and its link, and the re-run rules are unchanged.
- [ ] The log lines and the `checks-timeout` text read exactly as API Contract 5 states.

## Context

- creates: `internal/delivery/head_check_wait.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/head_check_wait_test.go`
- interface: `internal/delivery/check_wait_test.go`
- interface: `internal/delivery/check_rerun_test.go`
- interface: `internal/delivery/stale_check_test.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/delivery/merge_state_test.go`
- interface: `internal/delivery/revalidate_test.go`
- interface: `internal/delivery/check_rerun.go`
- interface: `internal/delivery/park_class.go`

## Consumers

- The `checks-failed` and `checks-timeout` blockers and the checking log — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting TestAFailedWorkflowRunWithoutAFailingCheckParks TestTheChecksTimeoutNamesWhatItWaitedFor TestAReadErrorIsLoggedOnce; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Deliver|Delivery|Queue|FollowUp' ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/delivery && GOOS=linux go build -buildvcs=false ./internal/delivery` — expected: exit 0. Before this Task none of the five new tests exists, so the first pass line is missing; afterwards the completion rule, the blockers and the log lines behave as API Contracts 3 to 5 state, every test of the package and every `deliver` command test passes, and both platforms build.

## References

- `_prd.md` → Goals; Core Feature 3; Core Feature 4; Core Feature 5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 2; Current behavior 3; API Contract 3; API Contract 4; API Contract 5; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3; ADR-0369; ADR-0236; ADR-0311


## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, no-workarounds, systematic-debugging, evidence-gate, tech-writer.

The checking stage now evaluates the Head Check Set after the existing check
recovery loop. It waits for unregistered workflow jobs, missing Expected
Workflows and an empty set, names all failed checks and runs with their
links, logs changes in pending members and read errors, and records the last
successful pending list or awaited merge state in the timeout blocker.

The existing recovery loop retains its stale/outside-change decisions,
re-run counts, poll interval and single deadline restart. Existing test
changes replace only the expected blocker or log text. The GitHub reader,
workflow discovery, Turn Record, Park Classes and database schema were not
changed; recording pending members on the Queue Turn belongs to task_04.

### Pre-edit characterization and regression signal

Before changing code, ran:

```sh
go test -count=1 -timeout 20m -v -run '^(TestTheCheckWait.*|TestEachCheckWaitStartsItsOwnDeadline|TestChecksTimeoutParksTheItem|TestABlockedMergeStateParksAtTheChecksTimeout)$' ./internal/delivery
```

Exit 0; all five selected tests passed. The default wait still made 241
polls over its fake-clock hour, and the configured duration, per-call
deadline and blocked merge-state cases passed before the change.

Added the five required regression tests before the implementation and ran:

```sh
go test -count=1 -timeout 5m -v -run '^(TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered|TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting|TestAFailedWorkflowRunWithoutAFailingCheckParks|TestTheChecksTimeoutNamesWhatItWaitedFor|TestAReadErrorIsLoggedOnce)$' ./internal/delivery
```

Exit 1 before implementation: all five top-level tests exposed the missing
behavior. The Fluxus replay moved to `merging` on poll 1 instead of waiting
through poll 3, the missing Expected Workflow was ignored, failing standalone
runs merged, timeout details were absent and the read-error line was absent.
The same focused command then exited 0 after implementation.

### Focused checks and acceptance evidence

After the final code and test edits, ran:

```sh
go test -count=1 -timeout 5m -v -run '^(TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered|TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting|TestAFailedWorkflowRunWithoutAFailingCheckParks|TestTheChecksTimeoutNamesWhatItWaitedFor|TestAReadErrorIsLoggedOnce|TestEveryFailedHeadMemberIsNamed|TestAWorkflowRunWaitsAfterItsCheckPasses|TestARequestedRerunWaitsForANewerWorkflowAttempt|TestCheckReadErrorsKeepTheLastSuccessfulWaitAndResetTheirLog)$' ./internal/delivery
```

Exit 0; all nine top-level tests and their subtests passed.

- Acceptance criterion 1: the Fluxus replay waits for three polls (30 seconds
  on the fake clock) even with `CLEAN`; its passing variant reaches `merging`
  only after every check and run succeeds. The Expected Workflow case times
  out with its missing path, or reaches `merging` when the run arrives at poll
  3. The empty report waits through its deadline as `no check reported`.
  `TestAWorkflowRunWaitsAfterItsCheckPasses` also keeps waiting after the
  run's listed check passes, until the run itself completes.
- Acceptance criterion 2: the Fluxus failure parks with the exact Backend
  job name, lower-cased `failure` and job URL. All six failing workflow
  conclusions park with the run name and URL; `neutral` and `skipped` reach
  `merging`. `TestEveryFailedHeadMemberIsNamed` proves ordered aggregation
  across failed/cancelled checks, a commit status without a link, and a run
  named by its path; the run already named by failed jobs is not duplicated.
  `TestARequestedRerunWaitsForANewerWorkflowAttempt` keeps the rerun's old
  failed workflow attempt pending and parks a newer failed attempt.
- Acceptance criterion 3: replay tests compare the complete ordered log,
  including both waiting lines and the final `3 check(s) completed for
  candidate` line. Timeout tests pin `no check reported`, `CI / verify`,
  `merge state BLOCKED`, `merge state UNKNOWN`, and the missing workflow
  path. Read-error tests prove one line for consecutive identical errors,
  reset after a successful read, and preservation of the last successful
  pending list when later reads fail.

Also ran:

```sh
go test -count=1 -timeout 5m -v -run 'Check|Rerun|MergeState|Premise|Revalidation|ParksDeclaredBlockers' ./internal/delivery
```

Exit 0; 67 top-level tests passed. The unchanged recovery assertions include
one outside-change rerun, stale reruns bounded per default-branch tip, no
repeated deadline extension, no rerun for cancelled/unattributable failures,
and the existing flaky-check warning and park. The original characterization
cases pass with only their new expected timeout detail. Persistent read
errors with no successful report retain the existing timeout text.

### Required sabotage probes

Built a temporary repository from `git archive HEAD`, overlaid this Task's
Go files, initialized/staged its Git metadata without committing, and
changed only that copy. Ran each probe with
`go test -count=1 -timeout 5m -v -run '<selection>' ./internal/delivery`.
All four returned exit 1 with test assertion failures, not build or fixture
errors; the transcripts were inspected.

| Sabotage | Selection | Input and observed assertion failure |
| --- | --- | --- |
| 1: omit workflow runs from evaluation | `^TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered$` | Poll 1: passing Conventions job beside Frontend and Backend runs `in_progress`; both variants report `polls=1 elapsed=0s`, item at `merging`. |
| 2: omit Expected Workflows from evaluation | `^TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting$` | Passing `coverage` status and `.github/workflows/ci-verify.yml` without a run; both variants merge at poll 1 instead of waiting. |
| 3: accept `action_required` as passing | `^TestAFailedWorkflowRunWithoutAFailingCheckParks$/action_required$` | Run 7 `CI | Verify` concludes `action_required` beside a passing check of another run; the item reaches `merging` instead of the exact named failure park. |
| 6: park bare `checks-failed` | `^TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered$/backend_fails$` | Poll 3: Backend job 30 fails in a changed package; the assertion sees bare `checks-failed` instead of its required job name, conclusion and URL. |

Probe logs and script remain outside the checkout at
`/tmp/0359-task03-sabotage-pribcxw8/sabotage-{1,2,3,6}.log` and
`/tmp/0359-task03-sabotages.py`. The copy's mutated Go files were restored
after the probes. `git diff --check` exits 0.

The initial worktree change was the Daemon-owned `task_03.md` status. This
turn preserves that status and adds only this Result. No other Task, Task
Graph, document, prohibited implementation file, or tooling configuration
was changed. Declared Verification, CLI verification, platform builds and
repository Verification are left to the Daemon; no commit, push or pull
request was made.
