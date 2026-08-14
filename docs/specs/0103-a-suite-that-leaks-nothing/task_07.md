---
task: task_07
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 07: Prove the tree exited, not only its owner

## Overview

Stopping a Run proves the registered owner exited. The processes that owner
spawned are what survived for three days. This slice makes Force Stop prove the
tree, using the process-table reader the residue inventory built.

## Requirements

1. MUST prove the exit of the processes a Run spawned, not only of its registered
   owner.
2. MUST fail rather than report success when a descendant survives, naming the
   surviving process.
3. MUST leave the command's output and exit status unchanged when the tree does
   exit.
4. MUST NOT terminate a process outside the Run's own spawn lineage.

## Subtasks

- [ ] Prove the tree's exit rather than the owner's.
- [ ] Name a surviving descendant in the failure.
- [ ] Cover a surviving descendant and a clean tree.

## Acceptance Criteria

- [ ] A Run whose descendant survives makes the stop fail, naming that process.
- [ ] A Run whose tree exits reports success exactly as before.
- [ ] No process outside the Run's lineage is touched.

## Verification

- `go test -count=1 ./internal/cli -run 'TestForceStopProvesTheTree' -v > /tmp/0103-t07.log 2>&1; s=$?; grep -q '^--- PASS: TestForceStopProvesTheTree' /tmp/0103-t07.log || { cat /tmp/0103-t07.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t07.log` — expected: exits 0, refusing a vacuous run.
- `grep -c '^--- PASS' /tmp/0103-t07.log > /tmp/0103-t07-n.txt; test "$(cat /tmp/0103-t07-n.txt)" -ge 2 || { echo "expected the surviving-descendant and clean-tree cases, got $(cat /tmp/0103-t07-n.txt)"; cat /tmp/0103-t07.log; exit 1; }` — expected: exits 0, proving both directions are exercised.

## Context

- interface: `internal/cli/cli.go`

## References

`_techspec.md` → Build Order 7; API Contracts. `_prd.md` → Core Feature 6;
Goal 2; Non-Goals, processes started outside Roundfix. ADR-0127.

## Result

### Implementation

- Force Stop now calls the existing lineage-safe tree terminator after its
  read-only owner proof and Agent Session cleanup. It inspects every returned
  process outcome before completing the Run or releasing its Active Run lock.
- Any unproven tree member keeps the Run Active and produces a Run-failure
  diagnostic with that process's PID and the controller's reason. Multiple
  unproven members are reported together.
- A fully proven tree follows the existing completion path without changing
  stdout, stderr, or the success exit status.

### Focused checks

- Before the production change,
  `rtk go test ./internal/cli -run '^TestForceStopProvesTheTree/surviving_descendant_fails_and_names_the_process$' -count=1`
  failed because Force Stop invoked owner-only termination. This was the red
  regression signal.
- `rtk proxy go test -json -count=1 ./internal/cli -run '^TestForceStopProvesTheTree'`
  passed and emitted separate top-level pass events for the surviving-descendant
  and clean-tree tests after the Verification feedback repair.
- `rtk go test ./internal/cli ./internal/store -run '^(TestForceStopProvesTheTree|TestRunForceStop|TestRunStopForce|TestOwnerProcessControllerTerminateTreeLeavesUnrelatedProcessRunning)' -count=1`
  passed 31 Force Stop and lineage-safety tests across both packages after the
  repair.
- `rtk git diff --check` passed.
- `rtk make verify-incremental` reached the full suite. `internal/cli` and
  `internal/store` passed, but the target exited 2 because
  `TestRuntimeCatalogueReadsAdvertisedModelsWithoutAnOverride` in
  `internal/agent` received an ACP session that does not advertise the
  `sandbox_mode` option. Task 07 does not change that package or runtime
  contract.

### Acceptance evidence

- Surviving descendant: `TestForceStopProvesTheTree`
  returns `exitRunFailed`, prints no success output, names descendant PID
  `424243` with `process remains alive after force kill`, and observes the Run
  still Active with its lock retained.
- Clean tree: `TestForceStopProvesTheTreeWhenTreeExits`
  proves both owner and descendant outcomes, then compares the complete stdout
  byte-for-byte with the prior success contract, observes empty stderr and
  `exitOK`, and observes the Run Stopped.
- Lineage boundary: both command cases require the recorded owner PID and
  identity as the tree root.
  `TestOwnerProcessControllerTerminateTreeLeavesUnrelatedProcessRunning`
  additionally exercises the real process controller and observes that a live
  process outside that lineage remains running.

### Verification feedback repair

- Daemon attempt 1 showed that both behavioral cases passed, but Go indented
  their subtest PASS records while the authored gate counts only top-level PASS
  records.
- The two unchanged behavioral cases now live in independent top-level tests
  whose names both start with `TestForceStopProvesTheTree`. The repair changed
  no production code and removed or relaxed no assertion.

### Not run

- The commands under `## Verification` were not run; the Daemon owns those
  commands and Task settlement.

### Follow-up

- The incremental gate's ACP runtime-catalogue failure is outside this Task's
  Force Stop slice and remains for the owning runtime work.
