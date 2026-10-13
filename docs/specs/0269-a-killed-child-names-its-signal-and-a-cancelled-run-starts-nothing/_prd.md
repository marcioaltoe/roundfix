---
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: active
created: 2026-10-13
surfaces: [cli]
---

# A killed child names its signal, and a cancelled Run starts nothing

This is an engineering-framed minimal PRD. The Spec fixes two defects and
enters through the TechSpec as a bug fix (`docs/agents/spec-routing.md`). It
adopts two Backlog Entries of 2026-10-13 that share one context: a test of the
full suite failed and nobody could say why from its report.

- [A test child is killed by a signal from an unknown sender](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md)
  (F6). `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout` failed
  three times because its `migrate --check` child "exited -1". Roundfix
  reduces a child that a signal ended to an exit code and throws the signal
  away: the Delivery Queue notice prints `exited -1`, and seven test call sites
  that run a child Roundfix binary turn the same end into exit `1`, the code a
  normal failure also uses. The Go documentation of
  `os.ProcessState.ExitCode` says why the number carries nothing: it "returns
  the exit code of the exited process, or -1 if the process hasn't exited or
  was terminated by a signal".
- [The Verification Capacity cancellation test sees a second command start](references/2026-10-13-the-verification-capacity-cancellation-test-starts-a-second-command.md).
  `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
  failed once on Linux CI with `total starts=2`. The cause is in production:
  an attempt that finds Verification Capacity free takes it through
  `TryAcquire` and starts its command without reading its context, so a Task
  that published Waiting for Verification just before the Run was cancelled
  starts a Verification command after the cancellation, against ADR-0056:
  "starts no command until cancellation-aware capacity acquisition succeeds".

The maintainer decided on 2026-10-13 that this Spec only names the signal: no
`sudo`, no trace of the sender. A new open Backlog Entry carries the trace.

## Prerequisites

None. The Spec is authored against `9963ffd3` and depends on no other active
Spec. Other cycle AG Specs also raise the `roundfix` skill's version and add
glossary terms; the operator rebases this plan on the default branch before its
plan Pull Request, and the record command picks the next free version.

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or persisted file
  format changes; the new wording is log and error text. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. Nothing reaches a network, a
  provider or a credential; the changed code reads a child's wait status and a
  context. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0056: an attempt "starts no command
  until cancellation-aware capacity acquisition succeeds", which task_04
  restores for the immediate grant. ADR-0266: "A Proven Signal is the only
  signal Roundfix code or tests send"; the new fixtures signal only their own
  process, and the rejected "Naming the sender first" stays rejected because
  `dtruss` and `dtrace` "need root, and the maintainer decided that F6 is fixed
  by fixing the class". ADR-0271: "A test bounds by real time only the budget
  it asserts"; the new tests give no child, goroutine or Git a real-time
  allowance. ADR-0259: every new top-level test "calls `t.Parallel()` as its
  first statement". ADR-0256: "Spec authoring requires a skills Task for a
  change a skill describes"; `internal/cli/deliver_workflow.go` is a source of
  the `deliver` Behavior Surfaces, so task_01 updates the `deliver` reference.
  ADR-0189 and ADR-0233 govern the owned skill's version raise, written only by
  the record command. ADR-0250 does not apply: no Baseline module
  version or coverage record changes. ADR-0244: "A Spec now carries a Glossary Declaration";
  this Spec adds one term. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile, CI
  or `go.mod` change; `golang.org/x/sys` is already a direct dependency. Four
  Governed Paths change: three in the owned `roundfix` skill and the
  release-gate test file whose child-binary helper task_03 sweeps. The
  maintainer approved cycle AG on 2026-10-13 with the same grant as cycles AE
  and AF for the Governed Paths each Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `internal/cli/baseline_release_gate_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A child process that a signal ended is reported by the signal's name and
  number and whether a core was dumped, wherever Roundfix or its tests reduce
  that child's end to an exit code today.
- The next time F6 happens, its failure report names the signal.
- A cancelled Implement Run starts no Verification command and no pre-work
  probe, whichever path granted the capacity.
- The cancellation test can no longer see a second start, on any scheduler.

## User Stories

1. As the maintainer reading a failed suite, I want the report to say
   `was killed by signal SIGKILL (9), no core dump` instead of `exited -1`, so
   that I can tell a kill from a crash or an abort and look for the right
   sender.
2. As an operator reading a Delivery Queue log, I want a fallback notice or a
   step failure caused by a signal to name it, so that I do not mistake a
   killed item binary for a migration refusal.
3. As an operator who stops a Run, I want no Verification command to start
   after the cancellation, so that a stopped Run does no new work.

## Core Features

1. One shared description of a Signal Exit: the signal's name and number and
   whether a core was dumped, with the exited case unchanged as `exited <n>`.
2. The Delivery Queue's item-binary fallback notice and its step failure name
   a Signal Exit; their exited forms keep today's bytes.
3. Every test helper that runs a child Roundfix binary and returns its exit
   code reports a Signal Exit on the last line of the standard error it
   returns; no test reduces a child's end through `exitCodeFromWait`.
4. Verification Capacity reads cancellation after every grant, the immediate
   one included, before any Verification command or pre-work probe starts.
5. The `roundfix` skill, the `deliver` user guide and the glossary describe a
   Signal Exit, and a new open Backlog Entry carries the trace of F6's sender.

## Non-Goals / Out of Scope

- Finding or fixing whoever sends F6's signal, `sudo`, `eslogger`, `dtrace` or
  any other trace. The new Backlog Entry owns it.
- Ten clean full suites as acceptance: F6 stays open by decision.
- Other production places that print an exit code, among them the Verification
  failure line of the Daemon and the detach exit code, and test helpers that
  run Git, `make` or a non-Roundfix fixture.
- Stop Request handling, the queued `Acquire` path and the gate's fairness,
  which already behave.

## Success Metrics

1. Success Metric: a child Roundfix binary that kills itself with `SIGKILL`,
   run through `runRoundfixBinaryMacro`, returns `-1` and a last standard-error
   line naming `SIGKILL (9)` and the core state; at `9963ffd3` it returns `1`
   and no such line.
2. Success Metric: `exitCodeFromWait` has 0 call sites in test files, against
   7 in 6 files at `9963ffd3`.
3. Success Metric: with a 50 ms pause injected before the immediate grant in a
   disposable copy, the cancellation test fails 0 of 20 runs at the audited
   head, against 20 of 20 on the same copy with the post-grant cancellation
   check removed.

## Acceptance evidence

The outside-evidence row rests on documentation this Spec did not write, read
through Exa during authoring and checkable offline:

- Go standard library, `os.ProcessState.ExitCode`
  (https://pkg.go.dev/os#ProcessState.ExitCode, also printed offline by
  `go doc os.ProcessState.ExitCode`): "ExitCode returns the exit code of the
  exited process, or -1 if the process hasn't exited or was terminated by a
  signal." It is why an exit code cannot name a signal and the wait status
  must be read.
- `golang.org/x/sync/semaphore`, `Weighted.Acquire`, the Go team's own
  capacity primitive. An earlier revision
  (https://go.googlesource.com/sync/+/37e7f081c4d4c64e13b10787722085407fe5d15f/semaphore/semaphore.go)
  documented the same defect: "If ctx is already done, Acquire may still
  succeed without blocking." Release v0.22.0
  (https://cs.opensource.google/go/x/sync/+/refs/tags/v0.22.0:semaphore/semaphore.go),
  the version this repository's `go.mod` requires, reads the context before
  the free grant, "We prefer to fail even if we could acquire the mutex without
  blocking", and again after a queued grant: "Acquired the semaphore. Check
  that ctx isn't already done." Roundfix's queued path already does the second;
  task_04 adds the first, after the grant, because `TryAcquire` takes no
  context and a check before it would leave a window before the grant.
- The Go Programming Language Specification, Select statements
  (https://go.dev/ref/spec#Select_statements): "If one or more of the
  communications can proceed, a single one that can proceed is chosen via a
  uniform pseudo-random selection." It is why the queued path already re-reads
  the context after waking, and why a check before the grant is not enough.

## Research record

- Secondbrain: read `wiki/index.md`, then ran
  `qmd query "child process killed by signal exit code -1 test flake" --all --files --min-score 0.3`,
  which returned the Roundfix mirror of `CHANGELOG.md` (0.76) and of ADR-0266
  (0.62): prior work on the same F6, which kept its sender unknown and
  decided that tracing needs root. A second query,
  `qmd query "context cancellation race select semaphore acquire goroutine" --all --files --min-score 0.3`,
  returned `projects/skills/mirror/skills/05-implementation-loop/golang-context/references/cancellation.md`
  (0.99), whose "Checking cancellation in loops" section checks `ctx.Err()`
  before each unit of work, and the `golang-concurrency` skill mirror (0.62).
  Neither changed the design; the first confirmed the maintainer's
  name-only scope.
- Exa: fetched the `os` and `syscall` package documentation and the Go
  specification, and searched for the `x/sync/semaphore` documentation. The
  search returned the earlier revision quoted above, whose queued path still
  kept a grant that landed after cancellation ("pretend we didn't notice the
  cancelation"); the v0.22.0 source in the local module cache gives the grant
  back instead. The history challenged the first design, a check before
  `TryAcquire`, and confirmed the check after it.
- Measurement: on 2026-10-13 the cancellation test passed 200 of 200 runs on
  macOS and 1,200 of 1,200 in `golang:1.26` containers at half a CPU or less,
  and failed 1 of 360 with six concurrent copies at half a CPU and
  `GOMAXPROCS=2`, with the CI message `total starts=2`. A 50 ms pause before
  `TryAcquire` for an attempt that found the gate busy, in a scratch copy,
  failed it 20 of 20; adding the post-grant context check to the same copy
  passed it 50 of 50.

## Glossary

- adds: **Signal Exit** — the end of a child process that a signal caused, reported by the signal's name and number and whether a core was dumped.

## Decisions

- Only the signal is named; the sender is not traced (maintainer, 2026-10-13,
  "Só nomear o sinal").
- The cancellation flake is a production defect, so production is fixed and
  the test keeps its assertion; no fixture synchronization is added.
- No ADR: both fixes restore what ADR-0056 and ADR-0266 already decided, and
  the wording of a Signal Exit is reversible text.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the shared description, the delivery
wording, the test-helper sweep, the capacity fix and the build order.
