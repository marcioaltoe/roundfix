---
task: task_03
spec: 0026-model-fallback-guardrail
status: pending
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

- [ ] Add the confirmation prompt to the guardrail path behind the
      interactivity checks.
- [ ] Proceed with the confirmed Fallback Selection into Run creation and
      the existing effective-selection surfaces.
- [ ] Fail on decline through the existing preflight failure path.
- [ ] Cover confirm, decline, empty answer, and prompt-suppression contexts
      with scripted-stdin CLI tests.

## Acceptance Criteria

- [ ] An interactive command with a failing configured selection prompts
      once; answering yes creates the Run whose header and Run record carry
      the fallback model and effort.
- [ ] Answering no or nothing creates no Run and exits with the existing
      Preflight Validation failure.
- [ ] The same command with no-input or detach never prompts and keeps the
      task_02 report contract.
- [ ] No configuration file changes after a confirmed fallback Run.

## Verification

- `rtk go test ./internal/cli` - expected: interactive confirm, decline,
  effective-selection persistence, and prompt-suppression tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → User Stories 1, 4; Core Features 2, 4.
- `_techspec.md` → API Contracts; Build Order 3.
- ADR-0041.
