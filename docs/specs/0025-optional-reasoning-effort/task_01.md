---
task: task_01
spec: 0025-optional-reasoning-effort
status: pending
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

## Subtasks

- [ ] Relax the empty-effort rejection in selection resolution while keeping
      the missing-model failure.
- [ ] Relax runtime selection validation and gate the config-key resolution on
      a non-empty value.
- [ ] Add the empty-effort skip to the disposable preflight selection path and
      the live session selection path.
- [ ] Extend the selection preflight error recovery copy.
- [ ] Update the agent-layer and selection unit tests for the empty-effort
      behavior, including the recorded-acpx-args assertions.

## Acceptance Criteria

- [ ] A RuntimeSpec with a model and an empty reasoning effort passes
      selection validation, and its preflight probe issues a model assignment
      but no reasoning set call.
- [ ] A live Agent Session for an empty-effort selection issues no reasoning
      set call.
- [ ] A non-empty reasoning effort still issues the set call, and a runtime
      rejection still produces the selection preflight failure whose recovery
      text names the model-managed remediation.
- [ ] A RuntimeSpec without a model still fails selection validation.

## Verification

- `rtk go test ./internal/agent ./internal/cli` - expected: selection,
  probe-skip, live-session-skip, and recovery-copy tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Goals; Core Features 1-2.
- `_techspec.md` → System Architecture; Interfaces; Build Order 1.
- ADR-0037, ADR-0039, ADR-0040.
