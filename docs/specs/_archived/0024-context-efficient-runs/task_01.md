---
task: task_01
spec: 0024-context-efficient-runs
status: completed
type: backend
complexity: high
---

# Task 01: Capture failed Verification diagnostics

## Overview

Move authoritative Verification output from the caller stream into bounded Run
artifacts and expose a typed command-failure contract for the repair cycle. The
slice is verifiable for successful, failed, cancelled, and infrastructure-error
commands before any Agent repair behavior is added.

## Requirements

1. MUST capture combined Verification stdout/stderr in an attempt-specific temporary artifact instead of streaming it to the Agent or caller.
2. MUST remove successful output and atomically retain failed command output at the TechSpec path.
3. MUST distinguish an exited Verification command from cancellation, process-start failure, and artifact filesystem failure.
4. MUST emit attempt-aware `daemon.verification` events with command phases and exactly one aggregate verdict per attempt.
5. MUST include only the diagnostic path and wrapped failure in events, never the command output body.
6. MUST apply the capture contract to Task command sequences and review Batch Verification without changing their current final settlement policy.

## Subtasks

- [x] Add the artifact-backed Verifier result and typed command error.
- [x] Create deterministic attempt paths with atomic failure retention.
- [x] Remove successful artifacts and preserve infrastructure error identity.
- [x] Add command-phase and aggregate-verdict Run Events.
- [x] Route Task and review Verification through the capture contract.
- [x] Cover real shell behavior and engine integration paths.

## Acceptance Criteria

- [x] A successful command leaves no Verification output artifact and emits a passed aggregate verdict.
- [x] A non-zero command retains combined output at the expected attempt path and returns the typed command failure.
- [x] Cancellation, process-start failure, and artifact write failure remain infrastructure errors and are not typed as repairable command failures.
- [x] Run Events contain attempt, phase, verdict, and diagnostic path metadata without output bytes.
- [x] Multiple Task Verification commands stop on the first failure and produce one aggregate verdict for that attempt.
- [x] Caller progress contains verdict summaries but no dots, passing-package banners, timings, or raw command output.

## Verification

- `rtk go test ./internal/daemon ./internal/runevent` - expected: real-shell capture, artifact lifecycle, error classification, and event-shape tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/runevent/event.go`

## References

`_prd.md` -> User Story 1; Core Features 1, 3; Success Metrics. `_techspec.md` -> Interfaces: Verification capture; Data Models: diagnostic artifacts; Build Order 1. ADR-0014; ADR-0038.

## Result

Implemented artifact-backed Verification capture for Daemon-owned Verification. `ExecVerifier` now writes combined stdout/stderr to a sibling temp file, removes successful output, atomically retains failed command output at `runs/<run-id>/verification/batch-<NNN>-attempt-<N>.log`, and returns `VerificationCommandError` only for commands that actually exit unsuccessfully. Cancellation, process-start failure, and artifact filesystem failures remain infrastructure errors.

Task and review Batch Verification now use a shared attempt runner. It emits `daemon.verification` command phases (`started`, `command-passed`, `failed`) plus exactly one aggregate `verdict` event per attempt, includes attempt metadata and diagnostic paths, and keeps command output bytes out of Run Events and caller progress. `settle` was updated for the verifier interface and now prints only a diagnostic path for failed Verification.

Evidence by acceptance criterion:

- Successful command cleanup and passed aggregate verdict: covered by `TestExecVerifierRemovesSuccessfulOutputArtifact`, `TestTaskCycleSettlesForgottenAgentStatus/completed_on_passing_verification`, and the affected package run.
- Non-zero command retention and typed failure: covered by `TestExecVerifierRetainsFailedOutputAsTypedCommandError` and `TestResolveCycleVerificationFailureRetainsDiagnosticsWithoutStreamingOutput`.
- Infrastructure error classification: covered by `TestExecVerifierClassifiesCancellationAsInfrastructureError`, `TestExecVerifierClassifiesProcessStartFailureAsInfrastructureError`, `TestExecVerifierClassifiesArtifactRetentionFailureAsInfrastructureError`, `TestResolveCycleVerificationInfrastructureErrorHaltsWithoutFailedSettlement`, and `TestTaskCycleVerificationInfrastructureErrorHaltsWithoutTaskSettlement`.
- Event metadata without output bytes: covered by `TestVerificationEventVocabulary`, `TestResolveCycleVerificationFailureRetainsDiagnosticsWithoutStreamingOutput`, and `TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdict`.
- Task command sequence stop/aggregate verdict: covered by `TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdict`.
- Caller progress summaries without raw command output: covered by `TestResolveCycleVerificationFailureRetainsDiagnosticsWithoutStreamingOutput`, the no-agent-console CLI tests, and `TestRunSettleVerificationFailureLeavesTaskAndTreeUntouched`.

Verification:

- `rtk go test ./internal/daemon ./internal/runevent`: passed, 86 tests.
- `rtk go test ./internal/daemon ./internal/runevent ./internal/cli`: passed, 455 tests.
- `rtk make verify`: passed; `rtk go test ./...` reported 1054 passed in 19 packages, `roundfix skills check` passed, and `go build` completed.

Follow-up notes:

- The one-repair Agent prompt behavior remains for the later repair-cycle task; this slice only exposes the typed command-failure contract and attempt artifacts.
