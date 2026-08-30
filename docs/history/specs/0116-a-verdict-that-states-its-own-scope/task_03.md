---
status: completed
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

## Result

Implementation:

- `VerificationCoverage` now carries whether the probe ran and the number of
  authored commands it covered into the text renderer.
- Every clean verdict line states either the executed command count or that the
  authored Verification commands were not executed. Finding reports retain
  their existing diagnostic rendering.
- The CLI passes its existing probe result to the renderer and no longer emits
  the superseded trailing `Verification: not run` note.

Acceptance evidence:

- `TestVerdictLineStatesProbeCoverage` reads only the verdict line and covers a
  probed clean result with two commands plus an unprobed clean result.
- `TestSpecCheckRunVerification` exercises the real CLI wiring in both
  directions, checks the executed count on the verdict line, and rejects the
  removed trailing note.
- `TestRenderResultTextAndJSON` exercises a result with findings after the
  renderer signature change; the finding code, severity, summary, locations,
  and fix remain observable.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck -run '^TestVerdictLineStatesProbeCoverage$'` — initially failed to compile because `VerificationCoverage` and the renderer argument did not exist, establishing the pre-change signal.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck -run '^(TestVerdictLineStatesProbeCoverage|TestRenderResultTextAndJSON)$'` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/cli -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestSpecCheckStageExitsZeroWithoutAFinding|TestRunSpecAuditPreservesSpecCheckBehavior)$'` — passed.
- `rtk git diff --check` — passed after the final code and test edits.

The Daemon-owned Verification command was not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260830T161359Z_31aaee7e42ecc4e4`
- Source commit: `e234ca6540d3494f5943e46d13bac49c723e1c68`
