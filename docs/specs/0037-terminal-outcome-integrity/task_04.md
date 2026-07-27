---
task: task_04
spec: 0037-terminal-outcome-integrity
status: pending
type: backend
complexity: medium
---

# Task 04: Interrupt every Review Source wait on Stop Request

## Overview

Make the Run Database Stop Request visible at every Review Source wait
boundary. A graceful stop interrupts status, retry, quiet-period, and
Merge-Ready waits by the next polling boundary and starts no later Review
Source or repository mutation.

## Requirements

1. MUST expose Stop Request observation through a narrow context-aware source.
2. MUST supply the Store-backed source to every operational watch Run.
3. MUST check the source before each Review Source status access.
4. MUST check again after every interruptible wait, including quiet, retry, and
   Merge-Ready sleeps.
5. MUST return the existing stop classification when a request is observed.
6. MUST not perform another fetch, check, artifact write, commit, push, or
   Review Source mutation after observation.
7. MUST propagate Store observation failures with Run and operation context.

## Subtasks

- [ ] Add the Stop Request source contract to watch dependencies.
- [ ] Wire the Store-backed source into operational watch Runs.
- [ ] Cover every status and sleep boundary.
- [ ] Preserve the existing stopped classification and CLI exit behavior.
- [ ] Add fake-clock tests for each waiting phase.
- [ ] Prove no downstream operation starts after observation.

## Acceptance Criteria

- [ ] A request during status wait reaches Stopped by the next poll.
- [ ] Requests during quiet period, transient retry sleep, and Merge-Ready wait
      each interrupt their respective boundary.
- [ ] No later Review Source call or repository mutation occurs.
- [ ] Operational runs always provide the Store source; only isolated tests may
      omit it.
- [ ] Store read failure is distinguishable from a requested stop.
- [ ] Existing Run Budget and timeout behavior remains unchanged without a stop.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/store/store.go`

## Verification

- `rtk go test ./internal/watch -run 'TestRun.*StopRequest' -count=1`
  — expected: every wait boundary observes Stop Request and starts no later
  operation.
- `rtk go test ./internal/cli -run 'TestRunWatch.*StopRequest' -count=1`
  — expected: Store-backed stop observation reaches the existing Stopped CLI
  contract.
- `rtk go test -race ./internal/watch ./internal/cli -run 'Test.*StopRequest' -count=1`
  — expected: cancellation and polling remain race-free.

## References

- `_prd.md` → Goal 3; User Story 2; Core Feature 4; User Experience; Success
  Metrics.
- `_techspec.md` → Interfaces: StopRequestSource; API Contracts: Graceful Stop
  Requests; Build Order 4.
- `../../adr/0022-stop-requests-travel-through-the-run-database.md` → durable
  Stop Request transport.
