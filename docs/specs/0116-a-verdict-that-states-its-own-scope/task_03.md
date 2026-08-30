---
status: pending
type: backend
---

# Task: The verdict line states the probe's coverage

An author reads `No findings.` and stops. The note saying the authored
Verification commands were never executed is printed after the verdict and
after the skipped-detector list, which is past where the reader stopped. The
scope belongs on the line that carries the verdict.

## Work

- Tell the renderer whether the probe ran and how many commands it covered. The
  caller already knows; today it keeps the answer and appends its own note.
- A clean verdict states its own coverage on the verdict line, in terms that
  name what was not covered rather than naming a flag.
- Remove the trailing note this replaces. One fact reported once: a line that
  says it twice in two voices is the same defect wearing a second coat.
- A verdict with findings is unchanged. This Task is about what a clean result
  claims.
- Cover both directions on the verdict line itself, not on the whole report: a
  probed clean run and an unprobed clean run render differently, and the
  assertion reads the verdict line so it cannot pass on text printed elsewhere.

## References

- `_prd.md` → Goal 2, User Story 2, Core Feature 2; Open Questions, resolved as
  every clean verdict
- `_techspec.md` → Build Order 3; Interfaces: `VerificationCoverage`
- ADR-0093 bounds the check to what the artifacts say, which this Task reports
  from rather than widens

## Verification
- `grep -q "VerificationCoverage" internal/speccheck/report.go && grep -q "TestVerdictLineStatesProbeCoverage" internal/speccheck/report_test.go && ! grep -q "Verification: not run (use --run-verification)" internal/cli/spec_check.go && go test -count=1 ./internal/speccheck ./internal/cli`
