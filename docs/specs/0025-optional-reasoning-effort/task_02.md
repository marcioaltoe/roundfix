---
task: task_02
spec: 0025-optional-reasoning-effort
status: completed
type: backend
complexity: medium
---

# Task 02: Surface model-managed reasoning across CLI commands

## Overview

Represent the model-managed reasoning state honestly on every non-interactive
CLI surface: the explicit empty reasoning-effort flag becomes a valid
override, Run headers render the state by name, and the Doctor and Setup
Commands accept a configured empty effort. The slice is verifiable through
buffer-captured CLI runs.

## Requirements

1. MUST accept an explicitly passed empty reasoning-effort flag as a valid
   model-managed override on every command that accepts the flag; the empty
   model flag stays invalid.
2. MUST render `Default Reasoning Effort: model-managed` in the resolve,
   watch, and implement Run headers when the effective reasoning effort is
   empty.
3. MUST make the Doctor and Setup Commands treat a configured empty reasoning
   effort as a valid selection whose agent probe assigns only the Agent
   Model.
4. MUST persist the empty effective reasoning effort on the Run row unchanged,
   and keep the existing Attach and Live Run View rendering for empty stored
   values.
5. MUST NOT change any other stdout report shape, stderr contract, or exit
   code.

## Subtasks

- [x] Allow the explicit empty reasoning-effort flag through selection-flag
      validation.
- [x] Render the model-managed header line across the resolve, watch, and
      implement headers.
- [x] Rework the Doctor and Setup configured-selection tests that expect an
      empty configured effort to fail.
- [x] Cover flag acceptance, header rendering, and Run-row persistence with
      CLI tests.

## Acceptance Criteria

- [x] A run command invoked with an explicitly empty reasoning-effort flag
      passes argument validation and creates its Run with an empty persisted
      reasoning effort.
- [x] The implement, resolve, and watch headers show
      `Default Reasoning Effort: model-managed` for an empty effective
      selection.
- [x] The Doctor Command reports `agent: ok` for a configured runtime with a
      model and an empty reasoning effort when the runtime probe succeeds.
- [x] An explicitly empty model flag still fails argument validation.

## Verification

- `rtk go test ./internal/cli` - expected: flag validation, header rendering,
  doctor/setup selection, and persistence tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Core Feature 3.
- `_techspec.md` → API Contracts; Build Order 2.
- ADR-0040.

## Result

Implemented the non-interactive CLI model-managed reasoning surface:

- `--reasoning-effort=` now passes explicit selection flag validation for
  resolve, watch, and implement commands; `--model=` remains invalid.
- Resolve, watch, and implement Run headers render an empty effective
  Default Reasoning Effort as `model-managed`.
- Doctor and Setup configured-selection tests cover runtime probes with a
  model and empty reasoning effort.
- Resolve, watch, and implement CLI tests verify empty reasoning persists on
  the Run row unchanged.

Pre-change signal:

- `rtk go test ./internal/cli` failed after adding the task tests:
  `TestRunResolveAcceptsExplicitEmptyReasoningEffort`,
  `TestRunWatchRendersModelManagedReasoningHeader`, and
  `TestRunImplementAcceptsExplicitEmptyReasoningEffort` exited at Preflight
  Validation because explicit empty reasoning was still rejected.

Verification:

- `rtk go test ./internal/cli`: passed (`370 passed in 1 packages`).
- `rtk make verify`: passed (`1048 passed in 19 packages`,
  `Roundfix skill check passed`, and `go build` completed).

Acceptance evidence:

- `TestRunResolveAcceptsExplicitEmptyReasoningEffort` verifies explicit empty
  reasoning passes validation, creates a Run, persists `ReasoningEffort == ""`,
  and renders `Default Reasoning Effort: model-managed`.
- `TestRunWatchRendersModelManagedReasoningHeader` and
  `TestRunImplementAcceptsExplicitEmptyReasoningEffort` verify the watch and
  implement headers render `model-managed`, with empty reasoning persisted.
- `TestRunDoctorAcceptsConfiguredEmptyReasoningEffort` verifies Doctor reports
  `agent: ok` and probes the configured model with an empty reasoning effort.
- Existing explicit-empty model rejection cases remain in
  `TestRunReviewAgentCommandsRejectExplicitEmptySelectionOverrides` and
  `TestRunImplementRejectsExplicitEmptySelectionOverrides`.

Follow-up:

- Interactive Input model-managed behavior remains in Task 03 and was not
  changed here.
