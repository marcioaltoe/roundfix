---
task: task_01
spec: 0006-acpx-run-robustness
status: completed
type: backend
complexity: medium
---

# Task 01: Classify by parsed result, journal the anomaly

## Overview

Implement ADR-0020 in the acpx runner: when the NDJSON stream already
delivered the `session/prompt` response, a nonzero acpx exit becomes a
transport anomaly riding the successful result instead of a Batch failure;
both engines journal the anomaly and proceed to verification. Twice this week
a finished Task was discarded over exactly this. Verifiable through the
fake-acpx classification matrix and engine journaling tests.

## Requirements

1. MUST track in the runner whether the prompt response line was fully
   parsed; on nonzero exit with a parsed result, return the result with a
   bounded anomaly summary (exit code + trimmed stderr tail, same bounds as
   infrastructure errors) instead of an error.
2. MUST keep every no-result classification byte-identical: nonzero exits
   without a parsed result remain Batch failures; exits 2 and 4 remain
   infrastructure errors; context cancellation keeps Stop semantics.
3. MUST journal a non-empty anomaly from both engines (resolve cycle and
   TaskCycle) through an existing daemon event kind with the anomaly text in
   the payload, before verification runs; verification and settlement flow
   exactly as on a clean exit (ADR-0014 unchanged).
4. MUST NOT change the `Runner` interface; the anomaly rides the execute
   result.

## Subtasks

- [x] Result tracking and anomaly construction in the runner
- [x] Classification matrix tests over the fake-acpx rig
- [x] Engine journaling of the anomaly in both cycles
- [x] Verification-still-gates proof tests

## Acceptance Criteria

- [x] Matrix tests cover: result+exit0, result+exit1 (anomaly, success),
      no-result+exit1 (failure), no-result+exit2 and +exit4 (infrastructure),
      result+exit130 (stop), partial-stream+exit1 (failure).
- [x] Engine tests prove: anomaly event lands in the Run Event Journal with
      the stderr tail; a Task whose verification fails after an anomaly still
      settles failed (the flip never bypasses the gate).
- [x] The 0003/0005 incident shape — full completion report then buffer-error
      line then exit 1 — reproduced in a fixture, now ends with the Task
      committed after passing verification.
- [x] Full suite passes with only deliberate classification-test updates.

## Verification

- `rtk go test ./internal/agent/ ./internal/daemon/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 4; Core Features 1, 4. `_techspec.md` →
Classification change, Interfaces, Build Order 1, Risks. ADR-0014, ADR-0020.
Round-2 dogfood findings 1–2.

## Result

Implemented result-over-exit classification for acpx prompt runs. `ExecuteResult`
now carries `TransportAnomaly`; the runner sets it when a parsed
`session/prompt` result with a stop reason is followed by a nonzero acpx exit
other than 130. The anomaly includes the exit code and the same bounded stderr
tail used by infrastructure errors. No-result exits still use the existing
Batch failure, infrastructure, and Stop Request classifications.

Both engines now journal a non-empty transport anomaly through existing daemon
event kinds before verification: resolve cycles use `daemon.batch`, and
TaskCycle uses `daemon.task`. Verification and settlement still run exactly as
before; a failing verification after an anomaly settles the Task failed and
creates no commit.

Evidence:

- `TestACPXPromptExitClassificationMatrix` covers result+exit0,
  result+exit1 anomaly success, no-result+exit1 Batch failure,
  no-result+exit2 and no-result+exit4 infrastructure errors, result+exit130
  Stop Request, partial-stream+exit1 Batch failure, and the incident shape
  with a prompt result followed by the buffer-error line and exit 1.
- `TestResolveCycleJournalsTransportAnomalyBeforeVerification` proves the
  resolve engine journals the anomaly with the stderr tail before verification
  and still runs `agent>verify>commit>source`.
- `TestTaskCycleTransportAnomalyStillLetsVerificationGateSettleFailure`
  proves a Task that sees the anomaly but then fails verification settles
  `failed` and creates no commit.
- `TestTaskCycleCommitsAfterTransportAnomalyAndPassingVerification` proves the
  incident-shaped Task path proceeds through verification and creates the Task
  commit after the gate passes.

Verification:

- `rtk go test ./internal/agent/ ./internal/daemon/` — passed:
  124 tests in 2 packages.
- `rtk go test ./...` — passed: 591 tests in 16 packages.
- `rtk make verify` — passed: full Go suite, `roundfix skills check`, and
  `go build -o bin/roundfix ./cmd/roundfix`.
