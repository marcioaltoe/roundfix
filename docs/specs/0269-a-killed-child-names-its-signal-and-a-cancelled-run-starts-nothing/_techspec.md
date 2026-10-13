---
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
prd: _prd.md
created: 2026-10-13
---

# A killed child names its signal, and a cancelled Run starts nothing — Technical Spec

## Executive Summary

Roundfix reads a child's end through `exec.ExitError.ExitCode()`, which is
`-1` for a child that a signal ended, and the test helpers go one step further
through `exitCodeFromWait`, which turns that `-1` into exit `1`. A small new
package, `internal/childexit`, reads the wait status instead and describes a
Signal Exit by name, number and core state; the Delivery Queue's item-binary
paths and a test-side helper in `internal/testfixture` use it, and every test
helper that runs a child Roundfix binary and returns its code is moved onto
that helper. Separately, `acquireVerificationCapacity` returns an immediate
`TryAcquire` grant without reading the context, which is the second start the
cancellation test saw; it now applies the same post-grant cancellation check
as the queued path. The trade-off accepted: a new package of about forty lines
instead of a private helper in `internal/cli`, because both production
(`internal/cli`) and test helpers (`internal/testfixture`, which production
must not import and which must not import `internal/cli`) need the one
description.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or persisted file
  format changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. The change reads wait statuses and
  contexts; nothing reaches a network or credential. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0056: an attempt "starts no command
  until cancellation-aware capacity acquisition succeeds", which the immediate
  grant violates today and task_04 restores. ADR-0266: "A Proven Signal is the
  only signal Roundfix code or tests send"; the two new fixture programs signal
  only their own process, whose identity is their own PID, and the test sends
  no signal. ADR-0271: "A test bounds by real time only the budget it asserts";
  no new test waits on a clock. ADR-0259: each new top-level test "calls
  `t.Parallel()` as its first statement". ADR-0256: "Spec authoring requires a
  skills Task for a change a skill describes"; task_01 changes the `deliver`
  reference that covers `internal/cli/deliver_workflow.go`. ADR-0189 and
  ADR-0233 govern the skill version raise, which task_01 records with the
  record command. ADR-0250 does not apply: no Baseline module version or
  coverage record changes. ADR-0244: "A Spec now carries a Glossary Declaration"; the
  Glossary below adds one term. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile, CI
  or `go.mod` change. The Governed Paths, measured with `speccheck.GovernedPath`
  through a `go test -overlay` probe that wrote nothing, are granted by the
  cycle AG approval of 2026-10-13. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `internal/cli/baseline_release_gate_test.go`, `skills/roundfix/SKILL.md`.

## System Architecture

- `internal/childexit` — new. `From` reads a child's `*exec.ExitError` and its
  `syscall.WaitStatus`; `Exit.String` renders `exited <n>` or a Signal Exit.
  The signal name comes from `golang.org/x/sys/unix.SignalName` on darwin and
  linux, in `signal_name_unix.go`, and from `syscall.Signal.String` elsewhere,
  in `signal_name_other.go`; on Windows no wait status reports a signal, so the
  fallback only keeps the package building. The package's tests construct wait
  statuses and start no process, so it needs no suite guard.
- `internal/cli/deliver_workflow.go` — `runRoundfix` keeps `exitCode` and adds
  the child's `childexit.Exit` to `roundfixCommandResult`; `stepExecutable`'s
  fallback notice and `roundfixCommandResult.failure` render it.
- `internal/testfixture/child_exit.go` — new. `ChildExitCode` gives a test
  helper the code and a note; seven `exitCodeFromWait` call sites and two
  helpers that return `ExitError.ExitCode()` move onto it.
- `internal/daemon/task_engine.go` — `acquireVerificationCapacity` checks the
  context after a `TryAcquire` grant, releases the capacity and returns the
  queued path's error when the Run is cancelled.

### The second start, found

`runTaskVerificationRequest` publishes Waiting for Verification, then calls
`acquireVerificationCapacity`, whose first statement grants free capacity with
`TryAcquire` and returns it. In the test, task A holds the only permit and
blocks in the fake verifier; task B publishes Waiting; the test reads that
event, asserts no start and cancels. A returns on cancellation and releases.
When B was descheduled between its Waiting publish and `TryAcquire`, B now
finds the gate free, gets the permit, publishes `started` and calls the
verifier with a cancelled context: start 2. The queued path cannot do this,
because `acquireShared` checks `ctx.Err()` under the lock and
`acquireVerificationCapacity` checks it again after `Acquire` returns.
`ExecVerifier.Verify` refuses a cancelled context before it spawns, so in
production the cost is a false `started` event and a refused verifier call
rather than a child process; but a cancellation that lands between
the grant and `exec` would spawn and immediately kill one. The same immediate
grant also serves `verifyTaskPreWork`, so a pre-work probe has the same
window. Reproduction is recorded in `_prd.md` → Research record.

### Current behavior, characterized at `9963ffd3`

- `exitCodeFromWait` (`internal/cli/detach.go`) returns `exitRunFailed`, `1`,
  for an `*exec.ExitError` whose `ExitCode()` is negative. Test call sites:
  `implement_test.go` (`runCLIHelper` and `runRoundfixBinaryMacro`),
  `reopen_test.go`, `reopen_late_dependency_test.go`,
  `review_dispose_lock_test.go`, `version_freshness_isolation_test.go` and
  `script_fixture_test.go`: seven in six files. `baseline_release_gate_test.go`
  (`runBaselineReleaseCLI`) and `baseline_update_citations_test.go`
  (`baselineCitationGoRun`) return `ExitError.ExitCode()`, `-1`, with no
  signal. A Go child that sends itself `SIGKILL` gives
  `err=signal: killed code=-1 signaled=true signal=9 core=false` (measured
  2026-10-13).
- The delivery notice prints `migrate --check exited -1: ` for a killed probe;
  `failure` prints `failed with exit code -1`. The exited forms are pinned by
  `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate` and by
  `deliver_budget_park_test.go` (`roundfix implement failed with exit code 1`);
  both keep their bytes.
- The production `exitCodeFromWait` callers in `detach.go` are out of scope
  and keep their behavior.

## Implementation Design

### Interfaces

```go
package childexit

// Exit is how a child process ended.
type Exit struct {
	Code       int    // the exit status, -1 when a signal ended the child
	Signal     string // the signal's name, such as "SIGKILL"; empty when it exited
	Number     int    // the signal's number; 0 when it exited
	CoreDumped bool
}

// From reads the end of a child from the error Run or Wait returned.
// A nil error is an exit with code 0; ok is false for any other error type.
func From(err error) (exit Exit, ok bool)
func (exit Exit) Signaled() bool
func (exit Exit) String() string
```

```go
package testfixture

// ChildExitCode returns a child's exit code and, when the child did not exit
// on its own, a note for the helper's caller: a Signal Exit or a child that
// did not run. The note is empty for an exited child.
func ChildExitCode(err error) (code int, note string)
```

Invariants:

1. `Exit.String` is `exited <n>` for an exited child, and
   `was killed by signal <NAME> (<number>), core dumped` or
   `was killed by signal <NAME> (<number>), no core dump` for a Signal Exit.
2. `ChildExitCode` returns `(n, "")` for an exited child, `(0, "")` for nil,
   `(-1, "child " + Exit.String())` for a Signal Exit and
   `(-1, "child did not run: <error>")` for any other error.
3. A swept helper that returns a child's standard error appends a non-empty
   note as the last line of that standard error and returns the code
   unchanged, so a caller's existing failure message shows it.
4. The exited forms of the delivery notice and of `failure` keep their bytes.
5. After any grant of Verification Capacity, a cancelled context releases the
   capacity and starts no Verification command or pre-work probe and
   publishes no `started` phase. A Stop Request still does not stop an
   attempt that owned capacity immediately.

### Data Models

None. No stored or serialized shape changes.

### API Contracts

1. API Contract 1: when the item binary's `migrate --check` ends by a signal,
   the Delivery Queue writes
   `roundfix: notice: Delivery Queue item <slug>: <step> runs the owner's binary; the item binary's migrate --check <outcome>: <detail>`
   where `<outcome>` is `Exit.String()` (`exited <n>` as today, or the Signal
   Exit form) and `<detail>` is today's first non-empty line of standard
   error, else standard output, possibly empty.
2. API Contract 2: when an item or owner step child (`implement`, `archive`,
   `review`) ends by a signal, the step error reads
   `<operation> was killed by signal <NAME> (<number>), <core state>`,
   followed by `: <detail>` only when the child wrote a non-empty line; an
   exited child keeps `<operation> failed with exit code <n>: <detail>`.
3. API Contract 3: `ChildExitCode` as Interfaces and Invariants 2-3 state;
   no `_test.go` file calls `exitCodeFromWait`.

### Surface Transcripts

None. The changed lines appear only in a Delivery Queue log when an item
binary is killed by a signal, which no single command invocation reproduces;
API Contracts 1-2 are asserted by named tests in task_02.

## Coverage Map

Deviations from the PRD: None.

Cross-read from `_prd.md`:

- Goal 1, User Story 1, Core Feature 1 → `childexit.From`, `Exit.String`;
  Invariant 1.
- Goal 1, User Story 2, Core Feature 2 → `runRoundfix`, `stepExecutable`,
  `roundfixCommandResult.failure`; API Contracts 1-2; Invariant 4.
- Goal 2, Core Feature 3, Success Metrics 1-2 → `testfixture.ChildExitCode`
  and the swept helpers; API Contract 3; Invariants 2-3.
- Goals 3-4, User Story 3, Core Feature 4, Success Metric 3 →
  `acquireVerificationCapacity`; Invariant 5.
- Core Feature 5 → the `deliver` reference and user guide, `CONTEXT.md`, the
  new Backlog Entry (task_01).

## Absence Evidence

1. `internal/childexit` — command: `ls internal/childexit; grep -rln 'WaitStatus' internal cmd`; result: no such directory at `9963ffd3`; the ten files that read a `syscall.WaitStatus` are all tests that assert a signal they sent themselves (detach, orphan, Helper Lifeline, residue and process tests), and no production file or shared helper describes one.
2. `testfixture.ChildExitCode` — command: `grep -rn 'ChildExitCode' internal; ls internal/testfixture`; result: no match at `9963ffd3`; the package's only non-test file is `fixture.go`, which holds `FixtureBinary`.
3. The post-grant check of the immediate path — command: `grep -n -A3 'TryAcquire(request.Mode); acquired' internal/daemon/task_engine.go`; result: the grant is returned with `return release, nil` and no `ctx.Err()` read at `9963ffd3`.
4. The new tests — command: `grep -rn 'TestRoundfixBinaryHelperNamesTheSignal\|ImmediateGrant\|ItemBinaryKilledBy' internal`; result: no match at `9963ffd3`.

## ADR Alternatives

- ADR-0056: no rejected alternative is taken; the cancellation-aware acquisition it decided is extended to the immediate grant.
- ADR-0266: its rejected "Naming the sender first" stays rejected; the Spec names the signal only and opens a Backlog Entry for the sender.
- ADR-0271: no rejected alternative is taken; no budget is raised and no new test waits on a clock.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first and no Sequential Test is added.
- ADR-0256: its rejected model judging stays rejected; task_01 is the skills Task its authoring rule requires.
- ADR-0189: no rejected alternative is taken; both version fields rise with the content.
- ADR-0233: no rejected alternative is taken; the record command writes the version.
- ADR-0244: its rejected model gate stays rejected; the term is declared and written by task_01.
- ADR-0250: does not apply; no Baseline module version or coverage record changes, so none of its alternatives is in play.

## Integration Points

None external. `make verify-docs` fingerprints the `deliver` user guide in the
Behavior Surface Record, which task_01 re-records with its command.

## Testing Approach

- `internal/childexit/childexit_unix_test.go` (darwin, linux): table test over
  constructed `syscall.WaitStatus` values for exit 2, `SIGKILL` without a core
  and `SIGABRT` with a core, plus `From(nil)` and a non-exit error.
- `internal/cli/deliver_item_binary_signal_unix_test.go`: a fixture item
  binary sends itself `SIGKILL` on `migrate --check` or on its step command,
  chosen by environment; the tests assert API Contracts 1-2 exactly. The
  existing exited-case tests stay unchanged.
- `internal/cli/child_signal_report_unix_test.go`:
  `TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild` runs a fixture
  that sends itself `SIGKILL` through `runRoundfixBinaryMacro` and asserts
  code `-1` and the last standard-error line
  `child was killed by signal SIGKILL (9), no core dump`.
- `internal/daemon/verification_capacity_cancellation_test.go`: two
  deterministic tests that cancel at the instant the race needs, with no
  timing. The attempt test calls `runTaskVerificationRequest` with a free
  gate and a `Publish` that cancels the context when it records the Waiting
  phase; the probe test calls `verifyTaskPreWork` with a free gate and a
  `RunStateStore` wrapper that cancels after `UpdateRunState`. Both assert
  `context.Canceled`, zero verifier calls, no `started` phase and a gate that
  grants again. Before the fix both see one start.
- The existing cancellation and Stop Request tests stay unchanged and keep
  passing; with the fix the cancellation test is deterministic, because B's
  grant can only follow A's release, which follows the cancellation.

## Build Order

1. The `deliver` reference and user guide name a Signal Exit, the glossary
   term, the Backlog Entry for the sender, the skill version record and the
   Behavior Surface Record (task_01).
2. `internal/childexit` and the delivery wording with their tests (task_02)
   (depends on: 1). It depends on step 1 because the Task that changes the
   `deliver` surface needs the Task that documents it in its ancestry.
3. `testfixture.ChildExitCode`, the helper sweep and the sabotage test
   (task_03) (depends on: 2), which imports `internal/childexit`.
4. The post-grant cancellation check and its two tests (task_04). It depends
   on no earlier step.
5. Final QA gate (task_05) (depends on: 3, 4).

## Risks & Considerations

- A swept helper that compared a code against `1` for a child that failed to
  start now sees `-1`; no current caller expects a start failure, and the note
  names the error.
- The self-killing fixtures are compiled and run on the test host; their files
  carry a darwin or linux build constraint, so the Windows build of
  non-test code is unaffected.
- F6 itself is not fixed. If it recurs in this Spec's QA, the report records
  the named signal; acceptance does not depend on it.

## Glossary

- adds: **Signal Exit** — the end of a child process that a signal caused, reported by the signal's name and number and whether a core was dumped.

## Decisions

- A shared package rather than two copies: production and `internal/testfixture` render the same text.
- The fix follows the grant, not precedes it: a check before `TryAcquire` still loses to a cancellation between the check and the grant, as the `x/sync/semaphore` documentation quoted in `_prd.md` warns.
- The cancellation test keeps its assertion and gets no new synchronization; the defect was in production.
- No ADR: see `_prd.md` → Decisions. Research is recorded in `_prd.md` → Research record.
