---
task: task_08
spec: 0001-implement-command
status: completed
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

- [x] Spec field in the input collector and its field set for the Implement Command
- [x] Active-Spec listing wired into the picker
- [x] Merge and precedence of collected values into the command request
- [x] Empty-list and `--no-input` failure paths

## Acceptance Criteria

- [x] Driving the collector synchronously with a scripted selection returns the chosen slug and agent; the resulting Run targets that Spec.
- [x] `--no-input` without `--spec` exits 2 naming the missing flag; `--interactive` with both flags still opens the flow.
- [x] With no active Specs, the command exits 2 with a message saying there is nothing to implement and where Specs live.
- [x] Agent selection is remembered across invocations; the spec slug is not.
- [x] The full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Feature 12. `_techspec.md` → System Architecture (tui), Coverage Map (Story 7), Build Order 9.

## Result

The Implement Command now opens Interactive Input under the shared rules — forced by `--interactive`, suppressed by `--no-input`, and otherwise when `--spec` is missing or `--agent` is explicitly cleared (the built-in `defaults.agent` makes a missing `--spec` the primary trigger). The flow presents a Spec picker: a numbered `Active Specs:` list of `spec.ListActive` slugs, selectable by 1-based number or by typing a slug, followed by the existing Agent field with the config/remembered suggestion. Collected values merge through the established `buildInteractiveInputRequest`/`applyInteractiveValues` precedence (flag-provided values pre-fill their fields; Enter accepts them). The Agent selection is remembered through the existing interactive defaults after Run creation; the spec slug is never stored. With `--no-input` and no `--spec`, the command fails with the shared `missing required --spec because --no-input disables Interactive Input` shape, exit 2. With no active Specs, the flow fails before any Run Database access with one actionable message naming `docs/specs/<slug>/_prd.md` `status: active`, exit 2 — never an empty picker. The two placeholder "not available yet" validation messages and the help-text caveat are gone.

Commands run:

- `rtk go test ./internal/tui/ ./internal/cli/` — pass (159 tests, 2 packages).
- `rtk go test ./...` — pass (401 tests, 16 packages).
- `make verify` — pass (fmt-check, tests, `roundfix skills check`, build).
- `rtk go run ./cmd/roundfix implement --help` — help text matches shipped behavior.

Evidence per acceptance criterion:

- Scripted selection → chosen slug and agent, Run targets that Spec: `TestCollectInputSpecPickerSelectsListedSpec` (tui, scripted stdin, by-number/by-slug/out-of-range) and `TestRunImplementInteractiveInputPicksSpecThroughCollector` (CLI seam driving the real `tui.CollectInput` with `"1\nclaude\n"`; asserts `run.SpecSlug`).
- `--no-input` without `--spec` exits 2 naming the flag: `TestRunImplementValidationFailures/missing spec with no-input`. `--interactive` with both flags still opens the flow: `TestRunImplementInteractiveForcedWithFlagsProvidedStillOpensFlow`.
- No active Specs → exit 2 with the actionable message and no Run Database side effects: `TestRunImplementValidationFailures/missing spec without active Specs` and `/interactive without active Specs`.
- Agent remembered, spec slug not: `TestRunImplementInteractiveInputRemembersAgentButNotSpec` (second invocation surfaces the remembered Agent suggestion; `Values.Spec` arrives empty).
- Full suite passes: `rtk go test ./...` green. Three implement validation table cases were updated because they encoded the pre-picker placeholder behavior this task replaces; the review-path suite passed untouched.

Follow-up: the `roundfix` skill does not document the Implement Command yet — Build Order 11 (skill + docs task) must describe the picker (numbered active-Spec list, number-or-slug selection) when it adds implement.
