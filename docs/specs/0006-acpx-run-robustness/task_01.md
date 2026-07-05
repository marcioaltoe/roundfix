---
task: task_01
spec: 0006-acpx-run-robustness
status: pending
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

- [ ] Result tracking and anomaly construction in the runner
- [ ] Classification matrix tests over the fake-acpx rig
- [ ] Engine journaling of the anomaly in both cycles
- [ ] Verification-still-gates proof tests

## Acceptance Criteria

- [ ] Matrix tests cover: result+exit0, result+exit1 (anomaly, success),
      no-result+exit1 (failure), no-result+exit2 and +exit4 (infrastructure),
      result+exit130 (stop), partial-stream+exit1 (failure).
- [ ] Engine tests prove: anomaly event lands in the Run Event Journal with
      the stderr tail; a Task whose verification fails after an anomaly still
      settles failed (the flip never bypasses the gate).
- [ ] The 0003/0005 incident shape — full completion report then buffer-error
      line then exit 1 — reproduced in a fixture, now ends with the Task
      committed after passing verification.
- [ ] Full suite passes with only deliberate classification-test updates.

## Verification

- `rtk go test ./internal/agent/ ./internal/daemon/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 4; Core Features 1, 4. `_techspec.md` →
Classification change, Interfaces, Build Order 1, Risks. ADR-0014, ADR-0020.
Round-2 dogfood findings 1–2.
