---
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
prd: _prd.md
created: 2026-10-16
---

# The F6 sender is named, and a rebuilt item binary is replaced — Technical Spec

## Executive Summary

F6 lives in the delivery item binary tests of `internal/cli`, and every change
here is test code. The fixture that stands in for an item binary records its
process, parent and group IDs; a test-side probe report says, on failure, how
the item binary probe ended and whether its build rewrote the file in place.
A diagnosis then runs up to ten full suites at that tree and attributes any
recurrence through the macOS diagnostic report. Last, the two delivery test
builds copy the fixture to a new file and rename it over `bin/roundfix`, and a
new test rebuilds while a copy of the previous build still runs.

The trade-off accepted is that the fix does not wait for the diagnosis to
prove which process still mapped F6's item binary. The rename removes the
kill for any holder, and the code reading below rules out the other senders a
test can reach, so a diagnosis that never sees a recurrence still leaves the
fix justified by ADR-0281's measurement. If a recurrence carries no
code-signing report, the PRD's Unreachable Acceptance applies instead.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or persisted
  format changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network, provider or
  credential; temporary repositories and homes only. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0281: "A test that rebuilds an
  executable replaces it." ADR-0266: "A Proven Signal is the only signal
  Roundfix code or tests send", with its acceptance of ten clean full suites.
  ADR-0125: "so no test executes a file that something may still hold open".
  ADR-0213: "must end when the test binary that started it ends". ADR-0259:
  "A test runs in parallel unless it names why it cannot". ADR-0271: "A test
  bounds by real time only what it asserts". Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/_authorization.md`;
  bounded files: none (`paths: []`).

## Current behavior

Characterized on 2026-10-16 at `9e723b82` on macOS 27.0.1 (arm64), Go 1.26.5,
without writing to the checkout (`git status` compared before and after).

- `newDeliveryItemFixture` declares the item binary build
  `mkdir -p bin && cp '<fixture>' bin/roundfix`, and
  `newDeliverySignalItemFixture` declares the same shape for its own fixture.
  `cp` onto an existing file keeps its inode (measured: two copies, one
  inode). In `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` the
  stderr subtest creates `bin/roundfix` and runs its probe; the stdout subtest
  copies over that file and runs the probe again. In
  `TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary` three subtests
  rebuild one fixture's `bin/roundfix` in place.
  `TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments` removes the file
  before each build and has never failed.
- The fixture records `Args`, `Dir` and `Home` with each start and no process
  IDs. The F6 failure message prints only the log; the four recorded failures
  read `migrate --check exited -1: ` with an empty detail, so the probe wrote
  nothing (two of them are in the Run artifacts of 2026-10-13 and 2026-10-14).
- An overlay probe test reusing the fixture rebuilt `bin/roundfix` with `cp`
  while a held copy of it ran: the probe's notice read "was killed by signal
  SIGKILL (9), no core dump: " 5 times in 5, the inode did not change, and each
  kill left `roundfix-<time>.ips` with the termination namespace
  `CODESIGNING`, indicator "Taskgated Invalid Signature" and parent process
  `cli.test`. The same probe with copy and rename passed 5 of 5 with a new
  inode each time. Outside Go, copying over a binary that no process maps and
  running it gave no kill in 300 identical and 100 alternating copies, and 3,000
  copy-then-run cycles of a Go harness under eight spawning goroutines gave no
  kill.
- The kernel log kept no entry from 2026-10-14, and the diagnostic reports of
  the earlier occurrences are gone, so no past occurrence can be attributed
  now.

Senders ruled out by reading:

- Context cancellation: `stepExecutable` and `runRoundfix` use
  `commandContextForTest(t, t.Context())` of the subtest, which is cancelled
  only after the subtest function returns; `Run` is synchronous.
- The suite guard: `markTestBinary` sets a mark unique to the test binary's PID
  and nonce, children inherit it, and `checkHelperResidue` runs after `m.Run`
  returns; no other package holds the `internal/cli` mark.
- A stale PID: the Signal rule of ADR-0266 allows only `SignalProvenProcess`,
  `SignalProvenProcessGroup`, `signalOwnerProcess`, `killResidue` and
  `interruptSelf`; the probe is a child of the test binary in its process
  group, and F6 recurred on 2026-10-14 after that rule merged.

## System Architecture

All changes are test code in `internal/cli`:

- `deliver_item_binary_test.go` owns the fixture, its start record and the F6
  test. task_01 extends the record and the F6 failure message; task_02 routes
  the build through one builder and adds a hold mode to the fixture source.
- `deliver_item_probe_report_test.go` (new, task_01) owns the probe report and
  its tests.
- `deliver_item_binary_signal_unix_test.go` uses the same builder (task_02).
- `deliver_item_binary_rebuild_test.go` (new, task_02) holds the rebuild test.
- `diagnosis/f6-sender.md` in this Spec folder (new, task_01) records the
  suites and the attribution.

## Implementation Design

### Interfaces

```go
// internal/cli/deliver_item_binary_test.go
type deliveryItemStart struct {
	Args            []string
	Dir, Home       string
	Pid, Ppid, Pgid int // written by the fixture at start
}

// internal/cli/deliver_item_probe_report_test.go
type deliveryProbeReport struct {
	InodeBefore, InodeAfter uint64              // 0 when bin/roundfix was absent
	Start                   *deliveryItemStart  // the step's migrate --check start, nil when none
	ContextErr              error               // the step context's Err after the step
}

func (f *deliveryItemFixture) step(t *testing.T, ctx context.Context, step string) (string, error, deliveryProbeReport)
func (r deliveryProbeReport) String() string

// internal/cli/deliver_item_binary_test.go (task_02)
func deliveryItemBuildCommand(binary string) string
```

### Data Models

The fixture's JSON start line gains `Pid`, `Ppid` and `Pgid`, read from
`os.Getpid`, `os.Getppid` and `syscall.Getpgrp`. No persisted format of the
product changes.

### API Contracts

1. API Contract: probe report. `step` stats `bin/roundfix` without following a
   link before and after calling `stepExecutable`, takes the last start line
   with arguments `migrate --check` written during the call, and reads the
   context's `Err`. `String` returns exactly
   `probe report: item binary inode <before> -> <after> (<form>); probe start <start>; context <state>`,
   where `<form>` is `absent` when after is 0, `created` when only before is
   0, `rewritten in place` when both are equal and non-zero, and `replaced`
   otherwise; `<start>` is
   `pid <pid> ppid <ppid> pgid <pgid>` or `not recorded`; `<state>` is `live`
   or the error text. Collecting it starts no process and waits on nothing,
   so it adds no delay between two steps.
2. API Contract: F6 failure message. When
   `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` finds a log other
   than the one it wants, its failure line keeps today's `log=%q want=%q` and
   appends the step's probe report. The test's passing behavior and assertions
   do not change.
3. API Contract: diagnosis record. `diagnosis/f6-sender.md` holds the lines
   `suites: <n>`, `recurrences: <n>`, `sender: <kernel code signing|process signal|not reproduced>`
   and `diagnostic report: <text>`, then one table row per suite with its start
   time, wall seconds, exit status and the count of output lines matching
   `was killed by signal`, `signal: killed` or `exited -1`. `sender` is
   `kernel code signing` when a `roundfix-*.ips` report written during the
   failing suite, with parent process `cli.test`, carries the namespace
   `CODESIGNING`; `process signal` when a recurrence has no such report;
   `not reproduced` when no suite showed F6.
4. API Contract: replacement build. `deliveryItemBuildCommand(binary)`
   returns `mkdir -p bin && cp '<binary>' bin/roundfix.next && mv -f bin/roundfix.next bin/roundfix`,
   with the binary path shell-quoted as today. Both delivery fixtures use it;
   no test writes the item binary in place.
5. API Contract: rebuild while held. The fixture source gains a `hold`
   argument that reads its standard input until end of file and exits 0 before
   recording a start, so it is a Helper Lifeline. The rebuild test runs one
   step, starts `bin/roundfix hold` with a standard input pipe only the test
   holds, runs a second step, and asserts the selection line, a new inode and
   one more `migrate --check` start; it then closes the pipe and waits for the
   holder to exit 0.
6. API Contract: everything else stays. No production file, help text, output
   line, Signal Exit wording or existing test assertion changes; the existing
   item binary tests pass unchanged except for the fields their records gain.

### Surface Transcripts

None. No command's help, stdout, stderr or exit code changes.

### Invariants

```text
1. The probe report starts no process and waits on nothing.
2. No delivery test writes its item binary into an existing file.
3. A rebuild while a copy of the previous build runs leaves the probe exit status, never a Signal Exit.
4. Production code and its outputs are byte-identical to the delivery base.
5. The diagnosis never reruns a failed suite to obtain a pass.
```

## Coverage Map

Deviations from the PRD: None.

- User Story 1 → API Contract 3, API Contract 4, the QA gate's ten suites
- User Story 2 → API Contracts 1-2
- Core Feature 1 → `deliveryItemStart`, `deliveryProbeReport`, API Contracts 1-2
- Core Feature 2 → API Contract 3
- Core Feature 3 → `deliveryItemBuildCommand`, API Contracts 4-5
- Success Metric 1 → the QA gate's ten full suites
- Success Metric 2 → API Contract 5, Invariant 3, the QA gate's sabotage row
- Success Metric 3 → API Contract 3

## Absence Evidence

1. Process IDs in the fixture start record — command: `grep -n 'Getppid\|Getpgrp\|Pid' internal/cli/deliver_item_binary_test.go`; result: no match at `9e723b82`.
2. Probe report — command: `grep -rn 'probe report\|deliveryProbeReport\|Ino' internal/cli/deliver_item_binary_test.go internal/cli/deliver_item_binary_signal_unix_test.go`; result: no match at `9e723b82`.
3. Replacement build — command: `grep -rn 'mv -f\|roundfix.next\|os.Rename' internal/cli/deliver_item_binary_test.go internal/cli/deliver_item_binary_signal_unix_test.go`; result: no match at `9e723b82`; both builds end in `' bin/roundfix"`.
4. Rebuild while held — command: `grep -rn 'hold' internal/cli/deliver_item_binary_test.go internal/cli/deliver_item_binary_signal_unix_test.go`; result: no match at `9e723b82`.
5. Diagnosis record — command: `ls docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced`; result: no `diagnosis` directory.

## ADR Alternatives

- ADR-0281: this Spec's decision; it rejects removing before copying, waiting for the holder, retrying a killed probe and tracing with `eslogger` first.
- ADR-0266: no rejected alternative is taken; the sender is named without root, which leaves "Naming the sender first" rejected for its stated reason, and the acceptance is unchanged.
- ADR-0125: no rejected alternative is taken; fixtures stay compiled, and the rename extends its concern to files a process still maps.
- ADR-0213: no rejected alternative is taken; the held copy ends at end of file on a pipe only the test binary holds.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.

## Integration Points

macOS only, read without privileges: the diagnostic reports under
`~/Library/Logs/DiagnosticReports`, read by the diagnosis and the QA gate in
their Agent turns, never by a test or a Verification command.

## Testing Approach

New tests call `t.Parallel()` first, use `newDeliveryItemFixture` and
temporary directories, and wait only on the processes they start.
`TestDeliveryProbeReportFormsEachEnd` and
`TestDeliveryProbeReportRecordsTheProbeStart` (task_01) pin API Contract 1
without rebuilding an item binary in place. `TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns`
(task_02) pins API Contracts 4-5; on macOS it fails every time the build
writes in place. The existing tests `TestDeliveryStepRunsTheItemBinary`,
`TestDeliveryStepsStartTheItemBinaryWithTheOwnersArguments`,
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate`,
`TestDeliveryStepNamesTheSignalThatKilledTheItemBinaryProbe` and
`TestDeliveryStepErrorNamesTheSignalThatKilledTheItemBinary` keep their
assertions; only the record they read gains fields and the build they declare
changes form.

## Build Order

1. Probe report, fixture process IDs, the F6 failure message, and the
   diagnosis of up to ten full suites at that tree (no dependency).
2. Replacement build for both delivery fixtures, the hold mode and the rebuild
   test (depends on: 1, which edits the same fixture file and must observe the
   unfixed build).
3. QA gate (depends on: 2).

No other Spec is a prerequisite.

## Risks & Considerations

- The diagnosis runs full suites in an Agent turn, which the Baseline's tiers
  otherwise leave to the Daemon (`docs/agents/agent-instructions.md`: "run
  focused tests of the changed packages and neither selected command"). The
  maintainer chose this on 2026-10-16 ("roda as suítes até
  reproduzir"); the suites are research, bounded to ten, and never stand in
  for the Task's Verification.
- F6 may not recur in ten diagnosis suites; the record then says
  `not reproduced`, and the fix rests on ADR-0281's measurement and the code
  reading above.
- A recurrence without a code-signing report means another sender; the QA gate
  records the Unreachable Acceptance and the Supervisor sends the `eslogger`
  command.
- On Linux the in-place form was not measured: it either fails to open the
  running file for writing or runs without a kill. The run that is certain to
  fail on regression is the macOS one, which is where the Daemon and the QA
  gate run.
- A user's item binary build that writes in place can meet the same kill when
  an earlier item binary still runs, for example a Run owner started by the
  implement step; the Supervisor records it as a follow-up Backlog Entry.

## Glossary

None.

## Decisions

- The probe report collects only facts that need no process or wait. See ADR-0281.
- The diagnosis runs before the fix, at most ten suites, and stops at the first recurrence.
- Both delivery fixtures replace `bin/roundfix` by copy and rename. See ADR-0281.
