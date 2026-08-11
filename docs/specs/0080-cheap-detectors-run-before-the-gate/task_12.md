---
task: task_12
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: backend
complexity: medium
---

# Task 12: Let the mechanical report satisfy the contract it enforces

## Overview

The mechanical stage refuses its own output. Report `-04` raised two
`QA-REPORT-SHAPE` findings against report `-03`, and `-03` was written by the
mechanical stage itself: it carries `## Mechanical findings`, an empty
`## Mechanical rows` table, and no `## Results` table at all, while declaring
`rows_blocked_finding: 1`.

Both refusals are correct readings of a real defect:

- *Results table has no report rows* — the writer puts its rows under
  `## Mechanical rows`, and every reader of a QA Report looks under
  `## Results`.
- *rows_blocked_finding is 1 but the Results table contains 0 matching rows* —
  the count is declared from the findings while the rows live somewhere the
  counter cannot see.

The consequence is a loop that cannot end on its own: each mechanical refusal
writes a report that guarantees the next mechanical stage refuses again. Expect
one more refusal after this Task lands, because the next stage inspects the last
malformed report before any well-formed one exists; that is convergence, not
failure.

A mechanical report is not a lesser artifact. It becomes the Spec's newest QA
Report, and `archive-spec` and `internal/spec.QAVerdict` read it exactly like an
Agent-written one.

## Requirements

1. MUST write every mechanical row into the `## Results` table the report-shape
   detector and every downstream reader already expect.
2. MUST make each declared `rows_blocked_*` count equal the number of matching
   rows in that table, so the report's own header agrees with its body.
3. MUST keep the mechanical provenance visible: a reader must still be able to
   tell a row the mechanical stage blocked from one an Agent settled.
4. MUST NOT relax the report-shape detector to accept the current output. The
   detector is right; the writer is wrong. Changing the reader to match a
   malformed writer would remove the only check that would have caught this.
5. MUST NOT change verdict semantics or the report naming contract.

## Subtasks

- [ ] Write mechanical rows into `## Results`.
- [ ] Derive the blocked counts from those rows.
- [ ] Keep mechanical provenance readable.

## Acceptance Criteria

- [ ] A mechanical report carries its rows in `## Results`.
- [ ] Every declared blocked count equals its matching row count.
- [ ] The report-shape detector raises nothing against a mechanical report.
- [ ] A mechanically blocked row is still distinguishable from an Agent row.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/task_engine.go`
- `internal/daemon/task_engine_test.go`
- `internal/speccheck/mechanical.go`
- `internal/speccheck/mechanical_test.go`
- `docs/specs/0080-cheap-detectors-run-before-the-gate/task_12.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestMechanicalReportSatisfiesTheReportShapeContract$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestMechanicalReportSatisfiesTheReportShapeContract'` — expected: exits 0. The case does not exist before this Task.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestMechanicalReportShape' -count=1 -v 2>&1 | tee /dev/stderr | grep -q -- '--- PASS'` — expected: exits 0, proving the detector still refuses a malformed report rather than being loosened.

## References

- `_prd.md` → Goal 1.
- `qa/qa-report-2026-08-11-04.md` → the two `QA-REPORT-SHAPE` findings.
- ADR-0080 → the blocked-cause counts this report declares.
