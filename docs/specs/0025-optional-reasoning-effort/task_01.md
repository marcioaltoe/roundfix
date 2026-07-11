---
task: task_01
spec: 0025-optional-reasoning-effort
status: completed
type: backend
complexity: medium
---

# Task 01: Select the Agent without a reasoning option

## Overview

Make an empty Default Reasoning Effort a valid Agent selection: selection
resolution and runtime validation accept it, and both the disposable preflight
session and the live Agent Session skip the acpx reasoning set call entirely.
The slice is verifiable on its own through the agent-layer unit tests that
record which acpx commands a selection issues.

## Requirements

1. MUST treat an empty resolved Default Reasoning Effort as a valid selection
   in selection resolution and in runtime selection validation; the Agent
   Model remains required.
2. MUST skip the reasoning set call on the disposable preflight Agent Session
   and on the live Agent Session when the trimmed Default Reasoning Effort is
   empty, while still assigning the Agent Model.
3. MUST keep the existing contract for non-empty values: the reasoning set
   call is issued, and a runtime rejection fails Preflight Validation without
   fallback.
4. MUST extend the selection preflight failure recovery text to name the
   model-managed remediation (setting the runtime's reasoning effort to the
   empty value) alongside choosing supported values.
5. MUST resolve the runtime-specific reasoning config key only when a
   non-empty value will be assigned.
6. MUST ship the claude runtime's built-in Default Reasoning Effort as the
   empty (model-managed) value, keeping the built-in models and the codex
   built-in effort unchanged.

## Subtasks

- [x] Relax the empty-effort rejection in selection resolution while keeping
      the missing-model failure.
- [x] Relax runtime selection validation and gate the config-key resolution on
      a non-empty value.
- [x] Add the empty-effort skip to the disposable preflight selection path and
      the live session selection path.
- [x] Extend the selection preflight error recovery copy.
- [x] Ship the empty claude built-in Default Reasoning Effort.
- [x] Update the agent-layer, config, and selection unit tests for the
      empty-effort behavior, including the recorded-acpx-args assertions.

## Acceptance Criteria

- [x] A RuntimeSpec with a model and an empty reasoning effort passes
      selection validation, and its preflight probe issues a model assignment
      but no reasoning set call.
- [x] A live Agent Session for an empty-effort selection issues no reasoning
      set call.
- [x] A non-empty reasoning effort still issues the set call, and a runtime
      rejection still produces the selection preflight failure whose recovery
      text names the model-managed remediation.
- [x] A RuntimeSpec without a model still fails selection validation.
- [x] The built-in claude runtime defaults resolve to a model with an empty
      reasoning effort, and the built-in codex defaults keep a non-empty
      effort.

## Verification

- `rtk go test ./internal/agent ./internal/cli ./internal/config` - expected:
  selection, probe-skip, live-session-skip, recovery-copy, and builtin-default
  tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Goals; Core Features 1-2.
- `_techspec.md` → System Architecture; Interfaces; Build Order 1.
- ADR-0037, ADR-0039, ADR-0040.

## Result

Implemented the model-managed Default Reasoning Effort slice:

- Selection resolution accepts an empty Default Reasoning Effort while still
  rejecting a missing Agent Model.
- Runtime selection validation accepts empty effort and resolves the
  runtime-specific reasoning config key only for non-empty efforts.
- Disposable preflight and live Agent Session setup skip the acpx reasoning
  `set` call when the trimmed effort is empty, while preserving the model
  assignment and the non-empty set path.
- Selection preflight failures for rejected non-empty efforts now name both
  supported-value remediation and `runtimes.<runtime>.reasoning_effort ""`
  for model-managed reasoning.
- Claude built-in defaults now resolve to `opus` with empty reasoning effort;
  Codex remains `gpt-5.5` with `xhigh`.

Pre-change signal:

- `rtk go test ./internal/agent ./internal/cli ./internal/config` failed after
  adding the task tests: empty-effort validation, preflight skip, live-session
  skip, recovery-copy, and Claude-default assertions failed.

Verification:

- `rtk go test ./internal/agent ./internal/cli ./internal/config`: passed
  (`552 passed in 3 packages`).
- `rtk make verify`: passed (`1046 passed in 19 packages`,
  `Roundfix skill check passed`, and `go build` completed).

Acceptance evidence:

- `TestACPXProbeSkipsEmptyReasoningEffort` verifies a model plus empty effort
  passes preflight, records model assignment, and records no reasoning set.
- `TestACPXRunSkipsEmptyReasoningEffort` verifies the live Agent Session
  records no reasoning set call for empty effort.
- `TestACPXProbeValidatesSelectionWithDisposableSession`,
  `TestACPXProbeSelectionRejectionClosesDisposableSession`, and
  `TestRunResolveSelectionPreflightRejectionReportsTupleAndCreatesNoRun`
  verify non-empty effort still sets reasoning and rejected values fail with
  model-managed recovery copy.
- `TestACPXRunRequiresModelSelection` verifies missing Agent Model still fails
  before acpx invocation.
- `TestBuiltinRuntimeDefaults` and
  `TestResolveSelectionUsesBuiltInRuntimeDefaults` verify Claude defaults to
  empty reasoning while Codex keeps `xhigh`.

Follow-up:

- Interactive Input model-managed rendering and acceptance remains in the
  TechSpec Build Order 3 slice and was not implemented here.
