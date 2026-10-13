---
task: task_06
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: completed
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

The local echo runtime now selects `realACPXEchoModel` (`echo-model`) with no
reasoning effort. The compiled fixture advertises that same constant as its
current and only model through `configOptions` on session creation and
load/resume. Its runtime ID is `echo-custom`: adding the Model alone to
`codex-custom` exposed the production Codex adapter lineage check, which
correctly rejected the local echo binary as not proving
`@agentclientprotocol/codex-acp`. Giving the fixture its own identity lets the
unchanged production selection path inspect its real advertised model.

The stdio command remains the compiled local echo agent. No provider model,
network call, production code, budget, or assertion was added or changed.
The diff preserves the round-trip payload/log checks, cooperative cancellation
without kill fallback, crash and resume, force-stop checks, and both no-orphan
checks. Only the declared integration test and this Result were edited; the
pre-existing `status: in_progress` is preserved.

Focused evidence:

- Before the repair,
  `ROUNDFIX_REAL_ACPX=1 rtk proxy go test -race -count=1 -v -run '^TestRealACPXCommandOverrideRoundTripCancelCrashResume$' ./internal/agent`
  exited 1 with `round-trip prompt through real acpx: agent model is required`.
- With the fixture identity and advertised Model, that single-run focused
  command exited 0, reported PASS (25.75 s), and did not skip. It exercised
  every existing assertion through the compiled echo fixture.
- Required sabotage: temporarily removed only `Model: realACPXEchoModel`
  from the repaired runtime initializer and ran the same single-run focused
  command. It exited 1 with
  `round-trip prompt through real acpx: agent model is required` at
  `acpx_integration_test.go:155`, before the first prompt; it did not skip.
  A `finally` block restored the Model. The focused command then exited 0
  with PASS (24.82 s), without a skip or race diagnostic.
- `rtk proxy gofmt -l internal/agent/acpx_integration_test.go` exited 0 with
  no output; `rtk proxy git diff --check` exited 0.

Acceptance evidence: the focused runs prove the selected model reaches the
local fixture under `-race`, preserving cancellation and no-orphan coverage.
Three consecutive runs remain for the Daemon's declared Verification, which
was not executed in this Agent turn. Task status and checkboxes are unchanged.

Follow-up observation: cleanup after explicit session closure logs
`missing session`; the focused test still passes its no-orphan assertions and
suiteguard exits 0. No cleanup behavior was changed in this slice.

## Carry-forward provenance

- Source Run: `run_20261013T063154Z_69b987082422f8c7`
- Source commit: `7cf4d2ab529e50dbaa08c63a202f78a39b4960e0`
