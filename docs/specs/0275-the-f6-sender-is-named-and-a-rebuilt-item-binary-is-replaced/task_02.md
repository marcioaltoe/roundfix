---
task: task_02
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
status: pending
type: test
complexity: medium
---

# Task 02: A rebuilt item binary replaces the file the previous build left

## Overview

Both delivery fixtures rebuild `bin/roundfix` by copying onto the existing
file, and macOS kills a probe whose signed executable was rewritten in place
while another process still maps it, with exactly F6's notice (ADR-0281,
`_techspec.md` → Current behavior). This Task routes both builds through one
builder that copies to a new file and renames it over the item binary path,
and adds a test that rebuilds while a copy of the previous build still runs.
It is verifiable on its own: on macOS the new test fails every time a build
writes in place.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [trace the sender of the F6 signal](references/2026-10-13-trace-the-sender-of-the-f6-signal.md),
   for its second half, its Expected "Find and fix the sender so ten full
   suites pass": the fix removes the in-place rewrite that ADR-0281 measured,
   whichever sender the diagnosis record names. It MUST read
   `diagnosis/f6-sender.md` first and copy its `sender` line into this Task's
   Result.
2. MUST add `deliveryItemBuildCommand(binary)` returning exactly the command
   of API Contract 4, with the binary path shell-quoted as today, and use it
   in `newDeliveryItemFixture` and `newDeliverySignalItemFixture`. No delivery
   fixture may build its item binary by copying onto `bin/roundfix`
   (Invariant 2); the builds that print, exit or only create `bin/` stay as
   they are.
3. MUST add the `hold` argument of API Contract 5 to the fixture source: it
   reads standard input until end of file and exits 0, before recording a
   start. The held copy is a Helper Lifeline (ADR-0213): its standard input is
   a pipe whose write end only the test holds, so it ends when the test binary
   ends even if no cleanup runs.
4. MUST add `TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns`
   in the new `internal/cli/deliver_item_binary_rebuild_test.go`, calling
   `t.Parallel()` first: one step through the task_01 probe report; then
   `bin/roundfix hold` started with a standard input pipe; then a second step;
   asserting the second step's selection line, the report form `replaced`
   with a new inode, and starts of exactly two `migrate --check` probes; then
   closing the pipe and asserting the held copy exits 0 (Invariant 3). It
   waits only on the process it started and sets no real-time bound
   (ADR-0271).
5. MUST prove the test against the old build in the Agent turn and record it
   in the Result: with `deliveryItemBuildCommand` temporarily returning the
   in-place copy, the test fails on macOS with the probe notice "was killed by
   signal SIGKILL (9), no core dump"; with the replacement it passes. The
   temporary change is not committed.
6. MUST keep the assertions of `TestDeliveryStepRunsTheItemBinary`,
   `TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments`,
   `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate`,
   `TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe`,
   `TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary` and the two
   task_01 report tests unchanged; their declared build is the only thing
   that changes. MUST NOT change any production file (API Contract 6).

## Subtasks

- [ ] Add the replacement builder and use it in both delivery fixtures.
- [ ] Add the hold mode to the fixture source.
- [ ] Add the rebuild test.
- [ ] Prove the test fails with the in-place build and passes with the replacement.

## Acceptance Criteria

- [ ] No delivery fixture copies its item binary onto the existing file.
- [ ] A rebuild while the previous build still runs leaves the probe exit status, not a Signal Exit.
- [ ] Every existing item binary test passes with its assertions unchanged.

## Context

- instruction: `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/cli/deliver_item_binary_signal_unix_test.go`
- creates: `internal/cli/deliver_item_binary_rebuild_test.go`

## Consumers

- delivery item fixture build — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns|TestDeliveryStepRunsTheItemBinary|TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate|TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe|TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary|TestDeliveryProbeReportFormsEachEnd|TestDeliveryProbeReportRecordsTheProbeStart)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns TestDeliveryStepRunsTheItemBinary TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary TestDeliveryProbeReportFormsEachEnd TestDeliveryProbeReportRecordsTheProbeStart; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the rebuild test does not exist, so its PASS line is missing.
- `matches="$(grep -n -F "' bin/roundfix\"" internal/cli/deliver_item_binary_test.go internal/cli/deliver_item_binary_signal_unix_test.go)"; test -z "$matches" || { printf 'a delivery fixture still copies onto its item binary:\n%s\n' "$matches" >&2; exit 1; }` — expected: exit 0. Before this Task both fixture builds end by copying onto bin/roundfix and match.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 6; Invariants 2 and 3; Current behavior; Testing Approach; Build Order 2
- ADR-0281; ADR-0213; ADR-0125; ADR-0259; ADR-0271
