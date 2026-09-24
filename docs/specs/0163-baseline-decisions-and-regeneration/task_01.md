---
task: task_01
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: backend
complexity: medium
---

# Task 01: A mode change keeps the HTTP decision

## Overview

The `http-contract` branch of `promptBaselineDecision` returns `{"mode": ...}` on "Change", which silently drops every recorded exception and the source of the stored `http.contract` decision.

## Requirements

1. MUST, when a valid current `http.contract` value exists and the maintainer chooses "Change", return a clone of that value with only `mode` replaced, validated through `baseline.ValidateDecisionValue`.
2. MUST name each kept exception scope in the review line of that prompt.
3. MUST keep the prompt as today when no current value exists.
4. MUST keep a value supplied through `--decision` or `--decision-file` replacing the whole typed decision, so a deliberate exception edit stays distinguishable from retention.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The multi-exception decision of the archived Finding, with a source, keeps every exception and the source, compared field by field, after a mode-only change.
- [ ] The review names each kept exception scope.
- [ ] An explicit decision value with one exception fewer yields exactly that value.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/baseline_human.go`
- interface: `internal/baseline/project_decisions.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestHTTPContractModeChangeRetainsExceptionsAndSource|TestHTTPContractModeChangeReviewNamesKeptExceptions|TestHTTPContractExplicitValueStillReplacesExceptions)$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestHTTPContractModeChangeRetainsExceptionsAndSource TestHTTPContractModeChangeReviewNamesKeptExceptions TestHTTPContractExplicitValueStillReplacesExceptions; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The HTTP mode change
- [2026-08-07-changing-the-http-contract-discards-its-exceptions.md](../../history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md)
