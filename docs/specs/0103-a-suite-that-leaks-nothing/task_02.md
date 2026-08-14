---
task: task_02
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 02: Fail in milliseconds when a fixture is already dead

## Overview

A test that waits for a spawned fixture waits out the whole agent budget when the
fixture never started, and then reports a timeout that names nothing. The wait
should observe the child: a process that exited cannot produce the file being
waited for, and saying so takes milliseconds.

## Requirements

1. MUST make every wait for a spawned fixture observe the child process, so an
   exited child ends the wait immediately.
2. MUST name the fixture and its exit status in the failure.
3. MUST keep the existing budget as the upper bound for a child that is alive but
   slow.
4. MUST NOT extend any deadline or add a retry.

## Subtasks

- [ ] Make the wait observe the child.
- [ ] Name the fixture and its exit in the failure.
- [ ] Cover a fixture that exits immediately.

## Acceptance Criteria

- [ ] A fixture that exits immediately fails its test in under one second.
- [ ] The failure names the fixture and its exit status.
- [ ] A fixture that is alive but slow still gets the full existing budget.
- [ ] No deadline was extended and no retry was added.

## Verification

- `go test -count=1 ./internal/agent -run 'TestSpawnedFixtureDeathFailsFast' -v > /tmp/0103-t02.log 2>&1; s=$?; grep -q '^--- PASS: TestSpawnedFixtureDeathFailsFast' /tmp/0103-t02.log || { cat /tmp/0103-t02.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t02.log` — expected: exits 0, refusing a vacuous run.
- `grep -E '^--- PASS: TestSpawnedFixtureDeathFailsFast \(0\.[0-9]+s\)' /tmp/0103-t02.log || { echo 'the fast-death case did not complete in under one second:'; grep '^--- PASS: TestSpawnedFixtureDeathFailsFast' /tmp/0103-t02.log; exit 1; }` — expected: exits 0, proving the wait ended in milliseconds rather than after the agent budget. It prints the observed duration on failure.
- `d=$(grep -rn 'agentWaitBudget *=' internal/agent/ | head -1); test -n "$d" || { echo 'the agent budget declaration vanished'; exit 1; }; echo "$d" | grep -qE '=[[:space:]]*90[[:space:]]*\*[[:space:]]*time.Second' || { echo "the budget was changed rather than the wait: $d"; exit 1; }; grep -rq 'TestSpawnedFixtureDeathFailsFast' internal/agent/ || { echo 'the budget is intact but the fast-death case does not exist'; exit 1; }` — expected: exits 0, proving the repair observed the child instead of moving the deadline. The forbidden-shortcut guard is anchored to the case it guards, because on an untouched tree the budget is unchanged and the guard alone would pass.

## Context

- interface: `internal/agent/acpx_integration_test.go`

## References

`_techspec.md` → Build Order 2; System Architecture, the bounded wait; Testing
Approach. `_prd.md` → Core Feature 3; Goal 3; User Story 3; Non-Goals.

## Result

### Implementation

- Every file, ACPX milestone, and prompt-start wait that owns a spawned-fixture
  result channel now selects on that channel as well as readiness and the
  unchanged wait budget. If readiness and process exit race, the wait preserves
  the result for the test's existing terminal assertion.
- An early exit reports the fixture's descriptive name and derives its exit
  status from `exec.ExitError.ProcessState`. The regression fixture exits with
  status 23 before creating its awaited file.
- `TestSpawnedFixtureDeathFailsFast` runs the real compiled test binary as the
  fixture and asserts that the resulting failure arrives in under one second
  with both the fixture name and `exit status 23`.

### Focused checks

- Pre-change signal: `rtk go test ./internal/agent -run
  '^TestSpawnedFixtureDeathFailsFast$'` failed because the old wait ignored the
  exited fixture and reached the regression's two-second containment timeout.
- `rtk proxy go test ./internal/agent -run
  '^TestSpawnedFixtureDeathFailsFast$' -v` exited 0 and reported
  `TestSpawnedFixtureDeathFailsFast (0.01s)`.
- The focused nine-test command covering all changed wait call sites exited 0;
  Go reported 11 passes including subtests.
- The same focused command with `-race` exited 0; Go reported 11 passes and no
  race.
- Source inspection after the last edit found
  `const agentWaitBudget = 90 * time.Second`; each changed wait still creates
  its deadline from that constant, and no retry was added.

### Acceptance evidence

- Immediate death under one second: the focused verbose regression completed in
  0.01 seconds.
- Named failure with exit status: the regression passes only when its child-test
  output contains `immediate-exit acpx fixture` and `exit status 23`.
- Live but slow fixture budget: the upper-bound branch remains
  `time.After(agentWaitBudget)`, and the constant remains 90 seconds; the focused
  call-site tests exercise live fixtures reaching their milestones normally.
- No deadline extension or retry: diff inspection found no change to
  `agentWaitBudget`, no replacement deadline, and no new retry path.

Daemon-owned Verification commands were not run in this Agent turn. No
follow-up work was discovered within this Task's slice.
