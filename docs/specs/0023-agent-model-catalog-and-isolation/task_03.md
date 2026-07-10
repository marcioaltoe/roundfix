---
task: task_03
spec: 0023-agent-model-catalog-and-isolation
status: completed
type: infra
complexity: high
---

# Task 03: Reject unavailable selections during Preflight Validation

## Overview

Prove the effective selection against the installed ACP adapter before an
operational command creates durable Run state. The slice is verifiable by
accepting supported values and rejecting unsupported model metadata or
reasoning values with zero Run rows and zero durable Agent Sessions.

## Requirements

1. MUST validate the exact model and reasoning pair through a uniquely named disposable Agent Session in the Git root.
2. MUST run the existing acpx/runtime readiness checks before selection validation and send no Agent prompt.
3. MUST close the disposable session on success, rejection, cancellation, and partial setup failure with bounded cleanup.
4. MUST complete selection validation before `resolve`, `watch`, or `implement` creates a Run or durable Agent Session.
5. MUST fail with exit `2` and an actionable stderr diagnostic naming runtime, model, reasoning, and both recovery paths.
6. MUST never substitute a model or reasoning value when the adapter rejects the selection.

## Subtasks

- [x] Add the workdir-aware probe request and disposable-session lifecycle.
- [x] Distinguish selection rejection from runtime infrastructure failure.
- [x] Join cleanup failures without losing the original validation error.
- [x] Wire operational command preflight ahead of Run creation.
- [x] Add adapter and CLI regression coverage for every terminal path.

## Acceptance Criteria

- [x] A supported selection completes preflight, closes its disposable session, and permits normal Run creation.
- [x] A missing-model-metadata rejection creates no Run row and leaves no disposable or durable Agent Session.
- [x] A rejected reasoning value reports the exact runtime/model/reasoning tuple and no fallback attempt occurs.
- [x] The disposable session receives model then reasoning and receives no prompt.
- [x] Cancellation and cleanup failures preserve context/error identity and do not leak a session silently.
- [x] Existing Doctor/setup runtime readiness checks remain functional with the evolved probe contract.

## Verification

- `rtk go test ./internal/agent ./internal/cli` - expected: disposable preflight, cleanup, rejection, zero-side-effect, and compatibility tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/store/store.go`

## References

`_prd.md` -> User Story 6; Core Feature 4; User Experience; Success Metrics. `_techspec.md` -> System Architecture; Interfaces: Agent Runner preflight; API Contracts; Build Order 3. ADR-0037; ADR-0039.

## Result

Status: completed.

Acceptance evidence:

- Supported selection: `TestACPXProbeValidatesSelectionWithDisposableSession` proves version check, disposable `sessions ensure --model`, reasoning `set`, close, and no prompt; existing resolve/watch/implement happy-path tests still create Runs after the probe passes.
- Missing model metadata: `TestACPXProbeModelRejectionSkipsReasoningAndClosesDisposableSession` proves model rejection closes the disposable session, skips reasoning and prompt; `TestRunResolveSelectionPreflightRejectionReportsTupleAndCreatesNoRun` proves zero Run database creation and zero durable Agent run calls.
- Rejected reasoning: `TestRunResolveSelectionPreflightRejectionReportsTupleAndCreatesNoRun` and `TestRunWatchSelectionPreflightFailureCreatesNoRun` assert the runtime/model/reasoning tuple, recovery text, one preflight attempt, and no fallback Agent selection.
- Disposable command stream: `TestACPXProbeValidatesSelectionWithDisposableSession` asserts `--model` precedes `set reasoning_effort`; `containsCommandKey(..., "prompt")` stays false.
- Cancellation and cleanup: `TestACPXProbeCancellationStillClosesDisposableSession` preserves `context.Canceled` and closes the disposable session; `TestACPXProbeCleanupFailureJoinsSelectionError` preserves both `SelectionPreflightError` and `AgentSessionCleanupError`.
- Doctor/setup compatibility: `TestACPXProbePassesWhenVersionMatchesPin` keeps workdir-less probes version-only; the existing doctor/setup health tests pass under the evolved probe bridge.

Verification:

- `rtk go test ./internal/agent ./internal/cli` passed: `Go test: 456 passed in 2 packages`.
- `rtk make verify` passed: `Go test: 1013 passed in 19 packages`; `roundfix skills check` passed; `go build` completed.
