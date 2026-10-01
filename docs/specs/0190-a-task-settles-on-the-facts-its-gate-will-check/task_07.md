---
task: task_07
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: pending
type: backend
complexity: low
---

# Task 07: The public event stream names the two Settlement Checks

## Overview

The QA gate of Run `run_20261001T024730Z_a33946bbe6ed3a5d` found F-01. `roundfix events <run> --filter verification` prints the Settlement Check events without their `command`, although the Run Event Journal holds both labels. The cause is the public projection in `internal/runevent/stream.go`, which keeps `command` only for an `unknown` classification. The stream redacts verification commands on purpose: `TestProjectStreamEventCoversStableCategoriesAndRedactsPayload` and `TestEventsReplayDefaultAndFilterJSONLRecordsOnly` require that a repository command such as `make verify` never appears. API Contract 2 asks only for the two labels Roundfix itself defines.

## Requirements

1. MUST make the verification projection in `internal/runevent/stream.go` set `command` when the payload's `command` is exactly `settlement check: spec consistency` or `settlement check: authorization`. Every other command MUST stay out of the public record, as today.
2. MUST add `TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot` in `internal/runevent/settlement_label_projection_test.go`. It covers `started`, `command-passed` and `failed` events for both labels, and a `make verify` event whose command stays absent.
3. MUST keep `TestProjectStreamEventCoversStableCategoriesAndRedactsPayload` and `TestEventsReplayDefaultAndFilterJSONLRecordsOnly` passing unchanged.
4. MUST NOT edit any other file.

## Subtasks

- [ ] Project the two labels.
- [ ] Prove the labels appear and repository commands stay redacted.

## Acceptance Criteria

- [ ] `roundfix events <run> --filter verification` shows both Settlement Check labels for a gated Task, and no repository command.

## Context

- interface: `internal/runevent/stream.go`
- creates: `internal/runevent/settlement_label_projection_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot|TestProjectStreamEventCoversStableCategoriesAndRedactsPayload)$' ./internal/runevent 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot TestProjectStreamEventCoversStableCategoriesAndRedactsPayload; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestEventsReplayDefaultAndFilterJSONLRecordsOnly$' ./internal/cli` — expected: exit 0; before this Task the new test does not exist.

## References

- task_02
- `_techspec.md` → API Contracts (API Contract 2)
