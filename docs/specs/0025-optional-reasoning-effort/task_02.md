---
task: task_02
spec: 0025-optional-reasoning-effort
status: pending
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

- [ ] Allow the explicit empty reasoning-effort flag through selection-flag
      validation.
- [ ] Render the model-managed header line across the resolve, watch, and
      implement headers.
- [ ] Rework the Doctor and Setup configured-selection tests that expect an
      empty configured effort to fail.
- [ ] Cover flag acceptance, header rendering, and Run-row persistence with
      CLI tests.

## Acceptance Criteria

- [ ] A run command invoked with an explicitly empty reasoning-effort flag
      passes argument validation and creates its Run with an empty persisted
      reasoning effort.
- [ ] The implement, resolve, and watch headers show
      `Default Reasoning Effort: model-managed` for an empty effective
      selection.
- [ ] The Doctor Command reports `agent: ok` for a configured runtime with a
      model and an empty reasoning effort when the runtime probe succeeds.
- [ ] An explicitly empty model flag still fails argument validation.

## Verification

- `rtk go test ./internal/cli` - expected: flag validation, header rendering,
  doctor/setup selection, and persistence tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → Core Feature 3.
- `_techspec.md` → API Contracts; Build Order 2.
- ADR-0040.
