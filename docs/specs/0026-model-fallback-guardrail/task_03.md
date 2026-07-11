---
task: task_03
spec: 0026-model-fallback-guardrail
status: completed
type: backend
complexity: medium
---

# Task 03: Confirm the fallback interactively before Run creation

## Overview

Add the interactive confirmation on top of the fallback report: when the
command may prompt, Roundfix reports the failed selection, the proven
Fallback Selection, and a token-cost caveat, then asks one confirmation
question; confirming starts the Run on the fallback, declining fails exactly
as today. The slice is verifiable through buffer-captured CLI runs with
scripted input.

## Requirements

1. MUST ask one confirmation question on stderr with a stdin answer,
   following the Interactive Input style, only when prompting is allowed
   (interactive stderr, no no-input, no detach).
2. MUST include the failed selection, the proven Fallback Selection
   (model-managed rendered by name), and a caveat that a different model can
   consume tokens differently, before the question.
3. MUST, on confirmation, proceed with the fallback as the effective
   selection for that Run only: the Run is created with it, the header and
   Run record carry it, and no configuration is written.
4. MUST, on decline or empty answer, fail Preflight Validation exactly as
   the non-interactive path does.
5. MUST leave the QA prompt, Spec picker, and every other Interactive Input
   field unchanged.

## Subtasks

- [x] Add the confirmation prompt to the guardrail path behind the
      interactivity checks.
- [x] Proceed with the confirmed Fallback Selection into Run creation and
      the existing effective-selection surfaces.
- [x] Fail on decline through the existing preflight failure path.
- [x] Cover confirm, decline, empty answer, and prompt-suppression contexts
      with scripted-stdin CLI tests.

## Acceptance Criteria

- [x] An interactive command with a failing configured selection prompts
      once; answering yes creates the Run whose header and Run record carry
      the fallback model and effort.
- [x] Answering no or nothing creates no Run and exits with the existing
      Preflight Validation failure.
- [x] The same command with no-input or detach never prompts and keeps the
      task_02 report contract.
- [x] No configuration file changes after a confirmed fallback Run.

## Verification

- `rtk go test ./internal/cli` - expected: interactive confirm, decline,
  effective-selection persistence, and prompt-suppression tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → User Stories 1, 4; Core Features 2, 4.
- `_techspec.md` → API Contracts; Build Order 3.
- ADR-0041.

## Result

Interactive resolve, watch, and implement commands now present the failed
selection, proven Fallback Selection, token-cost caveat, and one confirmation
question before Run creation. A yes answer uses the fallback as that Run's
effective selection. No and empty answers retain the Preflight Validation
failure, while no-input, detached, and non-interactive stderr paths retain the
task_02 report without prompting. Interactive Input fields remain unchanged.

Verification:

- `rtk go test ./internal/cli`: passed — 385 tests.
- `rtk make verify`: passed — 1,070 tests across 19 packages, Roundfix skill
  checks, and the CLI build.

Acceptance evidence:

1. Scripted `yes` and `y` tests for resolve, watch, and implement observed one
   confirmation question. The resulting Run headers, Run records, and Agent
   requests carried the fallback model and effort, including model-managed
   reasoning.
2. Scripted `no` and empty-answer tests exited 2, printed the Preflight
   Validation failure, created no Run Database, and started no Agent work.
3. No-input and simulated Detached Run child tests observed no confirmation
   question and retained the explicit-flags re-run report from task_02.
4. Resolve, watch, and implement confirmation tests compared Project Config
   contents before and after each Run and found no changes.

Follow-ups: none.
