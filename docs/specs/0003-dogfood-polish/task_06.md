---
task: task_06
spec: 0003-dogfood-polish
status: completed
type: frontend
complexity: low
---

# Task 06: Offer the QA gate in implement Interactive Input

## Overview

`--qa` is flag-only knowledge today. The implement Interactive Input flow
gains a final yes/no field for the QA gate, defaulting to no, preset by the
flag. Verifiable by driving the collector synchronously.

## Requirements

1. MUST add a `QA gate [y/N]` field to the implement Interactive Input field
   set, after Spec and Agent; empty input means no; `y`/`yes` (case-insensitive)
   means yes; anything else re-prompts once then errors with the field name.
2. MUST preset the field default to yes when `--qa` was passed (Enter keeps
   it), matching the existing flag-prefill precedence.
3. MUST NOT remember the QA choice across Runs — it is a per-Run decision like
   the Spec slug.
4. MUST keep `--no-input` behavior unchanged (flag-only path).

## Subtasks

- [x] QA field in the collector and the implement field set
- [x] Flag preset and value merge into the request
- [x] Scripted-stdin tests for yes, no, default, and invalid input

## Acceptance Criteria

- [x] Collector tests: scripted `y` produces a QA Run; empty input produces a
      non-QA Run; `--qa` + Enter keeps QA on.
- [x] The QA choice is not persisted in interactive defaults (second
      invocation shows `[y/N]` again without memory).
- [x] Full suite passes; help text already documents `--qa` (unchanged).

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 6. `_techspec.md` → API Contracts,
Build Order 6. Dogfood finding 7. ADR-0015 (QA opt-in semantics).

## Result

Implemented the implement Interactive Input QA gate as the final field after
Spec and Agent. The collector now accepts empty input as no, `y`/`yes` as yes,
`n`/`no` as no, and invalid input gets one re-prompt before an error naming
`QA gate`. The `--qa` flag prefills the field as yes; the collected value is
merged into the implement request. Interactive defaults still remember only
the existing PR/Agent values, not Spec or QA.

Evidence:

- Red signal: after adding the tests, `rtk go test ./internal/tui/ ./internal/cli/`
  failed because `CommandValues.QA` did not exist.
- `TestCollectInputImplementQAGate` covers scripted `y`, empty input, and
  `--qa` + Enter. `TestCollectInputImplementQAGateInvalidInputRepromptsOnce`
  covers invalid input re-prompting once and erroring with `QA gate`.
- `TestRunImplementInteractiveInputMergesQAGateChoice` proves scripted `y`
  invokes the QA step, empty input does not, and `--qa` + Enter keeps QA on.
- `TestRunImplementInteractiveInputRemembersAgentButNotSpecOrQA` proves the
  QA choice is not persisted and the second invocation shows `QA gate [y/N]:`.
- Existing implement help assertions still pass in `internal/cli`, so the
  documented `--qa` flag remains unchanged.

Verification:

- `rtk go test ./internal/tui/ ./internal/cli/` passed: 195 tests in 2 packages.
- `rtk go test ./...` passed: 475 tests in 16 packages.
- `rtk make verify` passed: full Go suite, Roundfix skill check, and build.
