---
task: task_03
spec: 0025-optional-reasoning-effort
status: completed
type: backend
complexity: low
---

# Task 03: Accept model-managed reasoning in Interactive Input

## Overview

Let Interactive Input express the model-managed reasoning state: the Default
Reasoning Effort field accepts an empty selection instead of raising the
required-field error, and collected-selection validation no longer rejects an
empty effort. The slice is verifiable through synchronous TUI model tests.

## Requirements

1. MUST accept an empty Default Reasoning Effort in the Interactive Input
   prompt loop when the runtime's configured default is empty, without
   raising the required-selection error.
2. MUST communicate in the field's prompt copy that an empty value means the
   Agent Model manages reasoning.
3. MUST keep the Agent Model field required.
4. MUST drop the effort emptiness check from collected-selection validation
   while keeping every other check.

## Subtasks

- [x] Allow the empty effort value through the prompt loop and
      collected-selection validation.
- [x] Add the model-managed meaning to the reasoning field's prompt copy.
- [x] Update the Interactive Input tests for empty-effort acceptance and
      unchanged model requirement.

## Acceptance Criteria

- [x] Interactive Input completes with an empty Default Reasoning Effort when
      the configured default is empty, and the resulting selection carries the
      empty value.
- [x] The reasoning field's copy names the model-managed meaning of an empty
      value.
- [x] An empty Agent Model still fails Interactive Input validation.

## Verification

- `rtk go test ./internal/tui ./internal/cli` - expected: Interactive Input
  acceptance, prompt copy, and validation tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Core Feature 3.
- `_techspec.md` → System Architecture; Build Order 3.
- ADR-0040.

## Result

Implemented Interactive Input support for model-managed reasoning:

- Empty Default Reasoning Effort no longer fails the prompt loop or collected
  selection validation.
- The reasoning field now tells users that an empty Default Reasoning Effort
  means the Agent Model manages reasoning.
- Agent Model remains required for commands that start an Agent.

Pre-change signal:

- `rtk go test ./internal/tui ./internal/cli` failed after adding the task
  tests: `TestCollectInputAcceptsModelManagedReasoningDefault` and
  `TestCollectInputOpenCodeAcceptsTypedModelWithModelManagedReasoning`
  returned the old required Default Reasoning Effort error.

Verification:

- `rtk go test ./internal/tui ./internal/cli`: passed
  (`555 passed in 2 packages`).
- `rtk make verify`: passed (`1050 passed in 19 packages`,
  `Roundfix skill check passed`, and `go build` completed).

Acceptance evidence:

- `TestCollectInputAcceptsModelManagedReasoningDefault` verifies Interactive
  Input completes with empty reasoning when the configured default is empty,
  preserves the empty value, and prints the model-managed meaning.
- `TestCollectInputOpenCodeAcceptsTypedModelWithModelManagedReasoning`
  verifies typed Agent Model plus empty reasoning passes collected selection
  validation.
- `TestValidateCollectedSelectionsKeepsAgentModelRequired` verifies empty
  Agent Model still fails validation.
