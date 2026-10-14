---
task: task_02
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: pending
type: backend
complexity: medium
---

# Task 02: The Delivery Queue names the signal that ended an item binary

## Overview

When the item binary's `migrate --check` or a step child is ended by a signal,
the Delivery Queue prints `exited -1` or `failed with exit code -1`, which is
how F6 reached its test report. This Task adds the `internal/childexit`
package, which describes a Signal Exit from the child's wait status, and makes
the Delivery Queue's item-binary notice and step failures use it. It is
verifiable on its own through the named tests and a Windows build.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [a test child is killed by a signal from an unknown sender](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md):
   a child that a signal ended is reported as a Signal Exit, never as
   `exited -1` or `exit code -1`.
2. MUST add `internal/childexit` with the `Exit`, `From`, `Signaled` and
   `String` signatures of `_techspec.md` → Interfaces and Invariant 1. The
   signal name comes from `golang.org/x/sys/unix.SignalName` in a file
   constrained to darwin and linux and from `syscall.Signal.String` in a file
   for every other system. The package's tests MUST start no process.
3. MUST carry the child's `childexit.Exit` in `roundfixCommandResult`, keep
   `exitCode`, and render the item-binary fallback notice as API Contract 1
   and `roundfixCommandResult.failure` as API Contract 2. The exited forms
   keep their bytes (Invariant 4), so
   `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` and
   the `TestDeliveryRunResult` tests of `deliver_budget_park_test.go` pass
   unchanged.
4. MUST add `internal/childexit/childexit_unix_test.go` with
   `TestExitDescribesHowAChildEnded`, a table over constructed
   `syscall.WaitStatus` values for exit 2, `SIGKILL` without a core and
   `SIGABRT` with a core, and over `From(nil)` and a non-exit error.
5. MUST add `internal/cli/deliver_item_binary_signal_unix_test.go`, with a
   darwin or linux build constraint, whose fixture item binary sends itself
   `SIGKILL` (its own PID, so no other process is signaled, ADR-0266) on
   `migrate --check` or on its step command as the environment selects:
   `TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe` asserts that
   the log is exactly the line
   `roundfix: notice: Delivery Queue item widget: review runs the owner's binary; the item binary's migrate --check was killed by signal SIGKILL (9), no core dump:`
   followed by one space, for the empty detail of a probe that wrote nothing,
   and a newline, and that the owner's executable is selected;
   `TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary` asserts that
   the `review` step error contains
   `roundfix review was killed by signal SIGKILL (9), no core dump`. Each
   test calls `t.Parallel()` first (ADR-0259) and waits on no clock
   (ADR-0271).
6. MUST NOT change `exitCodeFromWait`, the detach exit codes, the Daemon's
   Verification failure text, or any help text.

## Subtasks

- [ ] Add the `internal/childexit` package and its table test.
- [ ] Carry the exit in `roundfixCommandResult` and render the notice and the failure.
- [ ] Add the two signaled delivery tests.

## Acceptance Criteria

- [ ] The table test passes for every row, and the package builds for Windows.
- [ ] A killed `migrate --check` produces the exact notice line and the owner's binary runs.
- [ ] A killed step child produces the Signal Exit error, and every existing delivery test passes unchanged.

## Consumers

- The Delivery Queue notice and step error text — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- interface: `internal/cli/deliver_workflow.go`
- instruction: `internal/cli/deliver_item_binary_test.go`
- creates: `internal/childexit/childexit.go`
- creates: `internal/childexit/signal_name_unix.go`
- creates: `internal/childexit/signal_name_other.go`
- creates: `internal/childexit/childexit_unix_test.go`
- creates: `internal/cli/deliver_item_binary_signal_unix_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestExitDescribesHowAChildEnded$' ./internal/childexit 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestExitDescribesHowAChildEnded ' || { printf 'missing pass: TestExitDescribesHowAChildEnded\n' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/childexit ./internal/cli || exit 1; GOOS=linux go build -buildvcs=false ./internal/childexit ./internal/cli` — expected: exit 0; before this Task the package does not exist, so the test selects nothing and the builds fail.
- `out="$(go test -count=1 -v -run '^(TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe|TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate|TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^TestDeliveryRunResult' ./internal/cli` — expected: exit 0; before this Task the two signaled tests do not exist, so no PASS line is printed for them.

## References

- [_prd.md](_prd.md) — Goals; User Stories 1-2; Core Features 1-2
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; Interfaces; Invariant 1; Invariant 4; Testing Approach; Build Order 2
- [2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md)
- ADR-0266; ADR-0259; ADR-0271
