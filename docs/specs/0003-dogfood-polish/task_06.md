---
task: task_06
spec: 0003-dogfood-polish
status: pending
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

- [ ] QA field in the collector and the implement field set
- [ ] Flag preset and value merge into the request
- [ ] Scripted-stdin tests for yes, no, default, and invalid input

## Acceptance Criteria

- [ ] Collector tests: scripted `y` produces a QA Run; empty input produces a
      non-QA Run; `--qa` + Enter keeps QA on.
- [ ] The QA choice is not persisted in interactive defaults (second
      invocation shows `[y/N]` again without memory).
- [ ] Full suite passes; help text already documents `--qa` (unchanged).

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 6. `_techspec.md` → API Contracts,
Build Order 6. Dogfood finding 7. ADR-0015 (QA opt-in semantics).
