---
task: task_04
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: backend
complexity: high
---

# Task 04: Retry Temporary Verification Failure under exclusive capacity

## Overview

Turn child exit code `75` into the only typed Temporary Verification Failure
and give each Task one observable exclusive retry across its complete
Verification lifecycle. The retry must retain distinct diagnostics, execute
alone, remain separate from the existing Agent repair, and fail closed when
the temporary signal repeats.

## Requirements

1. MUST classify only child process exit code `75` as Temporary Verification
   Failure through a typed error that preserves the existing command failure
   and diagnostic chain.
2. MUST keep cancellation, process-start, artifact filesystem, and every other
   non-zero child exit in their existing infrastructure or deterministic
   failure classifications.
3. MUST release a normal shared permit before requesting the whole per-Run
   Verification Capacity for the exclusive retry.
4. MUST grant at most one exclusive retry per Task lifecycle across numbered
   attempts `1` and `2`, without adding a numbered attempt or consuming the one
   Agent Verification Feedback turn.
5. MUST retain the initial temporary diagnostic and write the retry diagnostic
   to a distinct `attempt-<N>-retry-1` artifact path.
6. MUST journal temporary classification, retry availability, exclusive wait,
   exclusive start, retry verdict, and exhaustion through bounded additive
   Verification event fields.
7. MUST let a deterministic retry failure use the existing Agent repair only
   when it has not already been used; a repeated exit `75` must settle failed
   without Agent repair or another retry.
8. MUST never classify from output text, elapsed time, framework name, port,
   timeout, package, database, container, or listener content.

## Subtasks

- [x] Add typed exit-75 classification at the real process boundary.
- [x] Add Task-scoped one-retry state without renumbering attempts.
- [x] Run the retry under whole-gate exclusive acquisition.
- [x] Preserve separate initial and retry diagnostics.
- [x] Publish retry lifecycle and exhaustion evidence.
- [x] Cover repair interaction, repeated signals, cancellation, and race behavior.

## Acceptance Criteria

- [x] A real shell command exiting `75` returns the typed temporary failure and
      retains its combined diagnostics; exit `1` remains deterministic.
- [x] Two normal attempts that could overlap are both drained before an
      exclusive retry begins, and no later shared attempt bypasses it.
- [x] A temporary failure followed by a passing exclusive retry settles the
      Task completed with zero Agent repair turns.
- [x] A temporary failure followed by a deterministic exclusive-retry failure
      receives the one existing Agent repair, then follows numbered attempt `2`.
- [x] A repeated exit `75`, including one occurring after Agent repair when the
      retry budget was already used, settles failed with no second retry.
- [x] Initial and retry diagnostic files both remain inspectable and no retry
      artifact overwrites an attempt artifact.
- [x] Event payloads distinguish numbered attempt, retry identity, mode,
      classification, availability, and final verdict without unbounded log text.
- [x] Output containing known timeout/listener phrases with a non-`75` exit
      never triggers a retry.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/golang-concurrency/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/daemon_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/engine_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/stream.go`

## Verification

- `rtk go test ./internal/daemon -run 'TestExecVerifier.*(Temporary|Exit|Diagnostic|Cancel)' -count=1` — expected: real shell exits preserve the typed classification boundary and distinct artifacts.
- `rtk go test ./internal/daemon -run 'TestTaskCycle.*(TemporaryVerification|ExclusiveRetry|RetryBudget|DeterministicRetry)' -count=1` — expected: pass, deterministic repair, repeated temporary failure, and one-budget matrices settle correctly.
- `rtk go test ./internal/runevent -run 'Test.*Verification.*(Temporary|Retry|Exclusive)' -count=1` — expected: additive event projection preserves retry identity and bounded classification evidence.
- `rtk go test -race ./internal/daemon ./internal/runevent -run 'Test(ExecVerifier.*(Temporary|Exit|Diagnostic|Cancel)|TaskCycle.*(TemporaryVerification|ExclusiveRetry|RetryBudget|DeterministicRetry)|.*Verification.*(Temporary|Retry|Exclusive))' -count=1` — expected: classification, exclusive scheduling, artifact, and event paths are race-free.

## References

- `_prd.md` → Goal 5; User Stories 6–8; Core Features 6–9; User Experience; Success Metrics.
- `_techspec.md` → Implementation Design: Interfaces, Data Models, and API Contracts; Testing Approach; Build Order 4.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → exit-75 protocol, exclusive acquisition, and one-retry bound.
- `../../adr/0038-daemon-allows-one-verification-repair.md` → independent one-Agent-repair contract.

## Result

Implemented the exit-`75` protocol at the real child-process boundary. The
typed temporary error unwraps through the existing command failure to the
child exit error, while cancellation, process-start failures, diagnostic
filesystem failures, and all other non-zero exits keep their prior
classification.

Each Task execution now owns one temporary retry budget across numbered
attempts `1` and `2`. A temporary shared attempt releases its permit before
publishing and acquiring an exclusive `retry: 1`; deterministic retry failures
can use the existing Verification Feedback turn, while a repeated temporary
failure settles failed without repair or another retry. Retry diagnostics use
`attempt-<N>-retry-1`, and bounded Run Event fields plus the public stream
projection retain attempt, retry, mode, classification, availability, reason,
and verdict without projecting command output.

Acceptance evidence:

1. `TestExecVerifierTemporaryExit75PreservesDiagnosticChain` used a real shell
   exit `75` and proved the combined diagnostic plus the temporary, command,
   and `exec.ExitError` chain. `TestExecVerifierExit1WithTimeoutTextRemainsDeterministic`
   proved exit `1` stays deterministic despite timeout, listener, database,
   and port text.
2. `TestTaskCycleExclusiveRetryDrainsSharedAttemptsAndBlocksLaterShared` proved
   two active shared acquisitions must drain and a later shared request cannot
   bypass the queued exclusive request.
3. `TestTaskCycleTemporaryVerificationPassesExclusiveRetryWithoutAgentRepair`
   settled completed after one exclusive retry with the observed call sequence
   `agent>verify>verify>commit`.
4. `TestTaskCycleDeterministicRetryUsesAgentRepairThenAttemptTwo` observed one
   Verification Feedback turn after the exclusive deterministic failure and
   then numbered attempt `2`.
5. `TestTaskCycleRetryBudgetExhaustsOnRepeatedTemporaryVerification` and
   `TestTaskCycleTemporaryVerificationAfterDeterministicRetryAndRepairDoesNotRetryAgain`
   proved exhaustion both immediately and after repair, with no second retry
   or extra Agent turn.
6. `TestExecVerifierTemporaryRetryDiagnosticPathsRemainDistinct` retained and
   read both real diagnostic files; the Task-cycle tests also observed the
   separate normal, retry, and attempt-`2` paths.
7. `TestVerificationTemporaryRetryExclusiveProjection` proved the bounded
   public fields survive projection and that unbounded child error text does
   not. Task-cycle event assertions covered temporary availability, exclusive
   waiting, retry verdict, and exhaustion.
8. The real-shell exit-`1` phrase test proves classification never depends on
   output content.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test ./internal/daemon -run 'TestExecVerifier.*(Temporary|Exit|Diagnostic|Cancel)' -count=1`
  — passed, 4 tests.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test ./internal/daemon -run 'TestTaskCycle.*(TemporaryVerification|ExclusiveRetry|RetryBudget|DeterministicRetry)' -count=1`
  — passed, 6 tests.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test ./internal/runevent -run 'Test.*Verification.*(Temporary|Retry|Exclusive)' -count=1`
  — passed, 1 test.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test -race ./internal/daemon ./internal/runevent -run 'Test(ExecVerifier.*(Temporary|Exit|Diagnostic|Cancel)|TaskCycle.*(TemporaryVerification|ExclusiveRetry|RetryBudget|DeterministicRetry)|.*Verification.*(Temporary|Retry|Exclusive))' -count=1`
  — passed, 11 tests in 2 packages.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test ./internal/daemon ./internal/runevent -count=1`
  — passed, 172 tests in 2 packages.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk make verify`
  — passed outside the sandbox: 2,712 tests in 23 packages, 4 Skill baseline
  tests, Roundfix Skill check, and CLI build. The sandboxed run first isolated
  five unrelated `/bin/ps` process-identity denials; the same gate passed when
  process inspection was allowed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-build rtk go test -race ./... -count=1`
  — passed outside the sandbox, 2,712 tests in 23 packages.
- `rtk git -c core.fsmonitor=false diff --check`
  — passed.

Follow-ups: none.
