---
task: task_08
spec: 0001-implement-command
status: pending
type: frontend
complexity: medium
---

# Task 08: Add the Interactive Input Spec picker

## Overview

Give the Implement Command the same Interactive Input parity the review commands have: when the spec flag is omitted in interactive mode, the flow lists the repository's active Specs for selection. Verifiable by driving the input collector synchronously and through the CLI's interactive-input seam.

## Requirements

1. MUST open Interactive Input for the Implement Command under the existing rules: forced by `--interactive`, suppressed by `--no-input`, and otherwise opened when `--spec` or `--agent` is missing.
2. MUST present a Spec picker listing the repository's active Specs (via the Spec contract package's discovery) alongside the existing agent field; the collected values merge into the request with the established precedence (flags before interactive values only where flags were given).
3. MUST fail with the existing non-interactive error shape when `--no-input` is set and `--spec` is missing, naming the missing flag.
4. MUST handle the no-active-Specs case with one actionable message instead of an empty picker.
5. SHOULD remember the agent selection through the existing interactive defaults, and MUST NOT remember the spec slug — each Run's target is an explicit choice.

## Subtasks

- [ ] Spec field in the input collector and its field set for the Implement Command
- [ ] Active-Spec listing wired into the picker
- [ ] Merge and precedence of collected values into the command request
- [ ] Empty-list and `--no-input` failure paths

## Acceptance Criteria

- [ ] Driving the collector synchronously with a scripted selection returns the chosen slug and agent; the resulting Run targets that Spec.
- [ ] `--no-input` without `--spec` exits 2 naming the missing flag; `--interactive` with both flags still opens the flow.
- [ ] With no active Specs, the command exits 2 with a message saying there is nothing to implement and where Specs live.
- [ ] Agent selection is remembered across invocations; the spec slug is not.
- [ ] The full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Feature 12. `_techspec.md` → System Architecture (tui), Coverage Map (Story 7), Build Order 9.
