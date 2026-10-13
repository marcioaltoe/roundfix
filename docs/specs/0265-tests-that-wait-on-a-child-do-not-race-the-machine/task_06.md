---
task: task_06
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: pending
type: test
complexity: low
---

# Task 06: The real ACPX round trip selects a model for its local fixture

## Overview

Corrective Task 1 of 2, after QA finding F-ACPX-MODEL of Run
`run_20261013T045342Z_a78df1f6c39a67e2`. The opt-in test
`TestRealACPXCommandOverrideRoundTripCancelCrashResume` in
`internal/agent/acpx_integration_test.go` builds
`RuntimeSpec{ID: "codex-custom", Protocol: ProtocolStdio, Command: helper}`
without a Model. `ACPXRunner.Run` calls `PrepareSession`, whose
`validateRuntimeSelection` (`internal/agent/acpx_runner.go`) refuses it with
`agent model is required`, so with `ROUNDFIX_REAL_ACPX=1` all 50 executions of
row 19's acceptance stopped before the first prompt. The initializer is the
same at the delivery base: a pre-existing defect that row 19's acceptance
exposed, not a load timeout.

## Requirements

1. MUST give the test's runtime a Model (and, only if the selection check
   needs it, a reasoning effort) that `validateRuntimeSelection` and
   `config.CheckSubscriptionRule` accept for the local echo fixture, without
   naming a provider model routed outside the subscription and without any
   network request: the fixture stays the compiled local echo agent.
2. MUST keep every assertion of the test, including task_02's cooperative
   cancel and no-orphan checks of row 19.
3. MUST NOT change production code.
4. MUST record in the Result one sabotage (removing the Model again) and the
   failure it produced.
5. MUST NOT change any file outside its Context.

## Subtasks

- [ ] Select a model the local fixture's runtime accepts.
- [ ] Prove the opt-in round trip passes with `ROUNDFIX_REAL_ACPX=1`.

## Acceptance Criteria

- [ ] With `ROUNDFIX_REAL_ACPX=1`, the test reaches the echo fixture and passes three consecutive times under `-race`.

## Context

- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `internal/agent/acpx_runner.go`
- instruction: `internal/config/config.go`
- interface: `internal/agent/acpx_integration_test.go`

## Verification

- `out="$(ROUNDFIX_REAL_ACPX=1 go test -count=3 -race -v -run '^TestRealACPXCommandOverrideRoundTripCancelCrashResume$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; ! printf '%s\n' "$out" | grep -q -- '--- SKIP: TestRealACPXCommandOverrideRoundTripCancelCrashResume' || { printf 'the real ACPX test skipped\n' >&2; exit 1; }; test "$(printf '%s\n' "$out" | grep -c -- '--- PASS: TestRealACPXCommandOverrideRoundTripCancelCrashResume ')" -eq 3 || { printf 'fewer than 3 passes\n' >&2; exit 1; }` — expected: exit 0. Before this Task every run fails with `agent model is required`.

## References

- QA report `qa/qa-report-2026-10-13.md` of Run `run_20261013T045342Z_a78df1f6c39a67e2`, F-ACPX-MODEL
- `_techspec.md` → Current behavior row 19
- ADR-0271

## Result
