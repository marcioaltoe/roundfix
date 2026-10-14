---
task: task_02
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: completed
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


## Result

Implemented the task_02 slice for Daemon Verification. `internal/childexit`
reads an `exec.ExitError` through its wait status and describes the signal's
name, number and core state. Darwin and Linux use `unix.SignalName`; the other
systems use `syscall.Signal.String`. A private wait-status decoder lets the
table test exercise constructed statuses without starting a process.

The Delivery Queue carries `childexit.Exit` alongside `exitCode`. The probe
notice renders the exit description, and signal step errors append only the
first non-empty stderr line, falling back to stdout. Ordinary exit failures
retain their existing wording and detail selection. The two new parallel
delivery tests use a fixture that signals only its own PID and waits on no
clock. The step-error test also checks empty output, stderr precedence and
stdout fallback.

### Focused-check evidence

Before the production edit,
`go test -count=1 -run 'TestDeliveryStep(NamesTheSignal|ErrorNamesTheSignal)' ./internal/cli`
exited 1: the probe printed `exited -1`, and each step-error case printed
`failed with exit code -1`. The new package was absent at the initial read.

After the production edit:

- `go test -count=1 -v -run 'TestExitDescribes|TestDelivery(Step|Steps|ItemBinary|RunResult)' ./internal/childexit ./internal/cli`
  exited 0. All five exit-table rows and all 17 selected CLI tests passed,
  including both new signal tests, the unchanged migration fallback test,
  owner-argument tests and Delivery Run result tests.
- `go test -count=1 -run '^TestDelivery' ./internal/cli` exited 0; the full
  delivery test family passed with the existing test files unchanged.
- `GOOS=windows go build -o /tmp/roundfix-task02-windows-childexit.a ./internal/childexit`
  and `GOOS=windows go build -o /tmp/roundfix-task02-windows-cli.a ./internal/cli`
  each exited 0, producing Windows production-package archives.
- `GOOS=linux go test -c -o /tmp/roundfix-task02-linux-cli.test ./internal/cli`
  and `GOOS=linux go test -c -o /tmp/roundfix-task02-linux-childexit.test ./internal/childexit`
  compiled successfully, including the Linux-constrained source and tests.
- `git diff --check` exited 0. Changed paths are this Task file,
  `internal/cli/deliver_workflow.go`, the four new `internal/childexit` files,
  and `internal/cli/deliver_item_binary_signal_unix_test.go`.

### Acceptance evidence

1. `TestExitDescribesHowAChildEnded` passed for exit 2, SIGKILL without a core,
   SIGABRT with a core, `From(nil)` and a non-exit error. It checks the decoded
   fields, `Signaled` and exact descriptions. Both Windows production packages
   built; Linux source and test compilation also succeeded.
2. `TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe` passed with
   the exact notice, including its trailing space and newline. It selected
   `os.Executable()`, ran that owner's `--version` successfully and recorded
   only the item's `migrate --check` invocation.
3. `TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary` passed with
   exact review-error text in all three output cases. The full delivery test
   family passed; existing delivery tests were not edited.

### Check limitation and execution boundary

An additional Windows CLI test-compilation attempt,
`GOOS=windows go test -c -o /tmp/roundfix-task02-windows-cli.test.exe ./internal/cli`,
reported existing Unix-only test code: `syscall.Kill` in
`detach_fixture_group_other_test.go`, `syscall.Wait4` in `detach_test.go`,
`Setpgid` in `implement_detach_teardown_test.go` and other Unix-only symbols.
No portability repairs were made outside this Task. Windows production builds
above succeeded; its childexit test-compilation attempt reported no test files
because the table test is constrained to Darwin and Linux.

Task status and checkboxes remain Daemon-owned. Declared Verification,
`make verify`, `make verify-changed` and the full repository suite were not
run. No other Task file or Task Graph was edited, and no commit, push or pull
request was made. No changes were made to detach exit codes,
`exitCodeFromWait`, Daemon Verification failure text or help text.
