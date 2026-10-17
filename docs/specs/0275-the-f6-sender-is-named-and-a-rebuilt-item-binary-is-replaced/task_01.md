---
task: task_01
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
status: completed
type: test
complexity: medium
---

# Task 01: The item binary tests report how the probe ended, and the diagnosis records the F6 sender

## Overview

F6's failures printed only `migrate --check exited -1: `, so nobody could tell
a kernel kill from a signal another process sent (`_techspec.md` → Current
behavior). This Task makes the delivery item fixture record its process,
parent and group IDs, adds a probe report that the F6 test prints on failure,
and then runs up to ten full suites at this tree, where the build still copies
in place, to catch a recurrence and attribute it through the macOS diagnostic
report. It is verifiable on its own through the new tests and the diagnosis
record.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [trace the sender of the F6 signal](references/2026-10-13-trace-the-sender-of-the-f6-signal.md),
   for its first half: the test records the pid, parent and signal of the
   `migrate --check` child, and the suites run until F6 recurs or ten have
   run.
2. MUST extend the fixture source and `deliveryItemStart` with `Pid`, `Ppid`
   and `Pgid`, written with every start from `os.Getpid`, `os.Getppid` and
   `syscall.Getpgrp` (Data Models). `assertStarts` keeps comparing only
   `Args`, `Dir` and `Home`.
3. MUST add the probe report of API Contract 1 in the new
   `internal/cli/deliver_item_probe_report_test.go`: `deliveryProbeReport`,
   its `String` form exactly as API Contract 1 writes it, with the forms
   `absent` (inode after is 0), `created`, `rewritten in place` and
   `replaced`, and the fixture method `step` that collects it around one
   `stepExecutable` call. Collecting the report MUST NOT start a process,
   sleep or wait (Invariant 1), because a delay between two steps could hide
   the race this Task diagnoses.
4. MUST make `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` run its
   step through `step` and append the report to its existing `log=%q want=%q`
   failure line (API Contract 2), changing no assertion and no build.
5. MUST add, each calling `t.Parallel()` first:
   `TestDeliveryProbeReportFormsEachEnd`, a table over the four forms, a
   recorded and an unrecorded start, and a live and a canceled context,
   asserting each `String` exactly; and
   `TestDeliveryProbeReportRecordsTheProbeStart`, which runs one step on a
   fresh fixture and asserts the form `created`, an inode after equal to the
   item binary's, a start whose `Ppid` is the test binary's PID and whose
   `Pgid` is its process group, the context `live`, and then a step with a
   canceled context whose report reads `probe start not recorded` and
   `context canceled`. Neither test rebuilds an item binary in place.
6. MUST then run the diagnosis of API Contract 3 in the Agent turn, at this
   Task's tree before the fix: `go test -count=1 -parallel 16 ./...` up to ten
   consecutive times with `HOME` set to a temporary directory, writing raw
   output only to a temporary directory outside the checkout and creating
   nothing under the checkout while a suite runs. It stops at the first suite
   in which F6 fails; it never reruns a failed suite to obtain a pass
   (Invariant 5). A suite that fails on another test is recorded and counted,
   and the run continues.
7. MUST, for a recurrence, record the failure line with its probe report and
   look under `~/Library/Logs/DiagnosticReports` for a `roundfix-*.ips` report
   written during that suite whose parent process is `cli.test`, recording
   its termination namespace and indicator, and set `sender` by API Contract
   3. When the directory cannot be read, the record says
   `diagnostic report: unreadable: <reason>` and `sender: process signal` is
   not inferred from it; it says `sender: not reproduced` only when no suite
   showed F6. No test and no Verification command reads that directory.
8. MUST write the result to the new
   `docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/diagnosis/f6-sender.md`
   with the four lines and the per-suite table of API Contract 3, plus the
   commit the suites ran on and that this Task's changes were applied. It MUST
   NOT copy an environment value, a home path or a credential into it; paths
   are written relative to the temporary root.
9. MUST keep every existing test of `deliver_item_binary_test.go` and
   `deliver_item_binary_signal_unix_test.go` passing with its assertions
   unchanged, and MUST NOT change any production file.

## Subtasks

- [ ] Record the process IDs in the fixture's start line.
- [ ] Add the probe report, `step` and the two tests.
- [ ] Append the report to the F6 failure line.
- [ ] Run up to ten full suites and read any diagnostic report.
- [ ] Write the diagnosis record.

## Acceptance Criteria

- [ ] A probe start records its pid, parent and group, and the report names them.
- [ ] The report tells a created, an in-place, a replaced and an absent item binary apart.
- [ ] The F6 test's failure line carries the report, and the test passes as before.
- [ ] The diagnosis record states how many suites ran, how many showed F6, and the sender.

## Context

- instruction: `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- creates: `internal/cli/deliver_item_probe_report_test.go`
- creates: `docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/diagnosis/f6-sender.md`

## Consumers

- delivery item fixture start record — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryProbeReportFormsEachEnd|TestDeliveryProbeReportRecordsTheProbeStart|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate|TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments|TestDeliveryStepRunsTheItemBinary|TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryProbeReportFormsEachEnd TestDeliveryProbeReportRecordsTheProbeStart TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments TestDeliveryStepRunsTheItemBinary TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the two report tests do not exist, so their PASS lines are missing.
- `f=docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/diagnosis/f6-sender.md; test -f "$f" || { printf 'missing %s\n' "$f" >&2; exit 1; }; grep -Eq '^suites:.([1-9]|10)$' "$f" || { printf 'missing suites line in %s\n' "$f" >&2; exit 1; }; grep -Eq '^recurrences:.[0-9]+$' "$f" || { printf 'missing recurrences line in %s\n' "$f" >&2; exit 1; }; grep -Eq '^sender:.(kernel.code.signing|process.signal|not.reproduced)$' "$f" || { printf 'missing sender line in %s\n' "$f" >&2; exit 1; }; grep -Eq '^diagnostic.report:..+' "$f" || { printf 'missing diagnostic report line in %s\n' "$f" >&2; exit 1; }` — expected: exit 0. Before this Task the record does not exist.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 6; Invariants 1 and 5; Current behavior; Testing Approach; Build Order 1
- ADR-0281; ADR-0266; ADR-0259; ADR-0271

## Result

Implementation and diagnosis evidence for Daemon settlement:

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Probe process, parent and group are recorded and named | The fixture writes `Pid`, `Ppid` and `Pgid` on each start. `TestDeliveryProbeReportRecordsTheProbeStart` asserts a positive PID, the test binary as parent, its process group, the recorded inode and the exact live report. |
| Created, rewritten in place, replaced and absent forms are distinguished | `TestDeliveryProbeReportFormsEachEnd` checks exact strings for all four forms crossed with recorded/unrecorded starts and live/canceled contexts (16 cases). The real-fixture test also checks a canceled call excludes the previous start. |
| F6's failure carries the report and its passing assertions are preserved | F6 uses `step` and appends the report to the existing `log=%q want=%q` failure line. Its assertions, `assertStarts` comparisons and both in-place build declarations are unchanged. The focused delivery selection below exercised F6 and the existing item binary and signal tests. |
| Diagnosis records suite and recurrence counts and sender | [diagnosis/f6-sender.md](diagnosis/f6-sender.md) records ten suites, zero F6 recurrences and `sender: not reproduced`, with base commit, applied changes, start times, durations, exit statuses, signal-line counts and relative raw-log paths. |

Focused checks:

- `go test -count=1 -run '^TestDelivery(ProbeReport|Step|Steps|ItemBinary)' ./internal/cli`: exit 0 (`roundfix/internal/cli`, 1.746 s) after the final code edit. An earlier run caught an incorrect expectation in the new canceled-context test; it was corrected to API Contract 1's literal `context context canceled` (the prefix plus the error text).
- Read the report helper and fixture diff: collection uses only `Lstat`, recorded JSON lines and `ctx.Err()` around one `stepExecutable` call; no additional process, sleep or wait was introduced. No production file or existing signal test was changed. Both new top-level tests call `t.Parallel()` first; neither rebuilds an item binary in place.

Required diagnosis (research, separate from Daemon Verification):

- Ran `go test -count=1 -parallel 16 ./...` ten consecutive times at `e5da154a7bf5d991d0b796f187e1abeb25c00312` with this Task's code changes applied, isolated temporary homes and raw logs outside the checkout. Each suite preserved the checkout fingerprint, recorded zero signal-related lines and did not reproduce F6. No diagnostic-report lookup was triggered.
- Each suite exited 1 on the unrelated `TestGuidanceCompositionJourney/rust-cli`: rustup could not select `rustfmt` because the temporary home had no default toolchain. All ten attempts were counted; no failed suite was rerun to obtain a pass. This environment limitation remains a follow-up outside this Task's slice.
- The diagnosis record is intentionally not a claim of ten clean acceptance suites or of an identified signal sender. Task 02 retains ownership of replacing both in-place builds.

The authored `## Verification` commands and the repository Verification commands were not run. Task status, terminal verdict and settlement remain Daemon-owned; no commit, push or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261017T121021Z_d7caac86f32133f1`
- Source commit: `3bf5497bf4b8bee039c3b4f2c52fd1a3266c1f40`
