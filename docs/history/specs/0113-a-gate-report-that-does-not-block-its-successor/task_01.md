---
task: task_01
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 01: Record the refusal as the report's terminal row

## Overview

A gate that stops at a precondition writes a Results table with no rows, which its
own contract calls malformed. The writer already knows the row count is zero, so
it records what actually happened: one terminal row naming the refusal. Two
adjacent defects in the same function are corrected while it is open — the verdict
line is written only when the result blocks, and two of the three typed counts are
hard-coded to zero.

## Requirements

1. MUST write one terminal row when the mechanical result carries no rows, naming
   what stopped the gate.
2. MUST write the verdict line unconditionally, so a non-blocking refusal does not
   produce a report with no verdict.
3. MUST compute all three typed blocked-cause counts from the rows written, rather
   than hard-coding two of them to zero.
4. MUST leave a report with rows byte-identical to what the writer produces today.
5. MUST NOT add a frontmatter field or change what any existing count means.

## Subtasks

- [ ] Write the terminal refusal row when there are no rows.
- [ ] Make the verdict line unconditional.
- [ ] Compute the three typed counts.
- [ ] Cover the empty, non-empty, and non-blocking cases.

## Acceptance Criteria

- [ ] A result with no rows produces a report with exactly one terminal row that
      names the refusal cause.
- [ ] Every report carries a verdict line, including a non-blocking one.
- [ ] The three typed counts agree with the rows written, in each case.
- [ ] A result with rows produces the same bytes as before this Task.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestWriteMechanicalQAReport' -v > /tmp/0113-t01.log 2>&1; s=$?; grep -q '^--- PASS: TestWriteMechanicalQAReport' /tmp/0113-t01.log || { cat /tmp/0113-t01.log; exit 1; }; grep -rq 'TestWriteMechanicalQAReportRecordsTheRefusal' internal/daemon || { echo 'the existing writer tests pass, but the refusal case does not exist'; exit 1; }; exit $s` — expected: exits 0, proving the writer's existing behaviour survives *and* that the refusal case was added. The existing tests pass on an untouched tree, so the regression half is anchored to the work rather than left standing.
- `go test -count=1 ./internal/daemon -run 'TestWriteMechanicalQAReportRecordsTheRefusal' -v > /tmp/0113-t01-refusal.log 2>&1; s=$?; grep -q '^--- PASS: TestWriteMechanicalQAReportRecordsTheRefusal' /tmp/0113-t01-refusal.log || { cat /tmp/0113-t01-refusal.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing refusal case; fails today, where no such test exists.
- `test -s /tmp/0113-t01-refusal.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestWriteMechanicalQAReportRecordsTheRefusal' /tmp/0113-t01-refusal.log || { echo 'the refusal case did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t01-refusal.log > /tmp/0113-t01-n.txt; test "$(cat /tmp/0113-t01-n.txt)" -ge 3 || { echo "expected the empty, non-empty, and non-blocking cases, got $(cat /tmp/0113-t01-n.txt)"; cat /tmp/0113-t01-refusal.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each case runs on its own.
- `grep -n 'rows_blocked_environment: 0' internal/daemon/task_engine.go && { echo 'the environment count is still hard-coded to zero'; exit 1; }; grep -n 'rows_blocked_declared: 0' internal/daemon/task_engine.go && { echo 'the declared count is still hard-coded to zero'; exit 1; }; grep -q 'verdict:' internal/daemon/task_engine.go || { echo 'the writer no longer emits a verdict'; exit 1; }; exit 0` — expected: exits 0, proving both counts are computed and the verdict survives. It prints the offending line on failure. Fails today on both hard-coded counts.

## Context

- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 1; System Architecture, the report writer; Risks &
Considerations, the hard-coded counts. `_prd.md` → Core Feature 1; Goal 1; User
Story 1. ADR-0132, ADR-0080.

## Result

### Implementation

- When the rendered mechanical result has no carried or blocked rows, the writer
  now adds one `QA-PRECONDITION` terminal row. Its provenance names the mechanical
  finding that stopped the gate and records that no other checks executed.
- The writer now emits `verdict: fail` for blocking results and `verdict: pass`
  for non-blocking results.
- All three typed blocked-cause counts now come from the statuses in the rendered
  Results table. No frontmatter field or count meaning changed.
- The populated blocking case has an exact-byte assertion for the pre-Task report
  body.

### Focused checks

- Red baseline: `rtk go test ./internal/daemon -run
  '^TestWriteMechanicalQAReportRecordsTheRefusal$/empty_result_records_its_blocking_cause$'
  -count=1` failed because the Results table contained no refusal row.
- `rtk go test ./internal/daemon -run
  '^TestWriteMechanicalQAReportRecordsTheRefusal$' -count=1` passed with four
  tests.
- `rtk go test ./internal/daemon -run
  '^(TestWriteMechanicalQAReportRecordsTheRefusal|TestMechanicalReportSatisfiesTheReportShapeContract|TestMechanicalStageWithholdsAgentSession|TestMechanicalStageSeedsReportBeforeAgentSession|TestWriteMechanicalQAReportPreservesSameDayNamingAndPriorReport)$'
  -count=1` passed with eight tests.
- The Task's declared Verification commands were not run; the Daemon owns that
  gate.

### Evidence per acceptance criterion

1. `empty result records its blocking cause` asserts the report contains exactly
   one `QA-PRECONDITION` row with status `fail`, the finding code and detail, and
   the fact that no other checks executed.
2. `non-blocking result carries a verdict` asserts a non-blocking report contains
   `verdict: pass`; the populated blocking byte assertion includes
   `verdict: fail`.
3. The writer counts `blocked (environment: ...)`, `blocked (finding: ...)`, and
   `blocked (declared: ...)` statuses from the rendered Results table. The empty,
   populated, and non-blocking cases assert the resulting frontmatter counts.
4. `populated blocking result preserves its bytes` compares the complete report
   against the pre-Task output, including frontmatter, findings, Results, and
   skips.
