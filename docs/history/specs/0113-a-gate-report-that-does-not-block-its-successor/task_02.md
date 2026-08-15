---
task: task_02
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 02: Prove a refused report does not block its successor

## Overview

The defect is a sequence, not a state: a gate refuses, writes a report, and the
next run of the same Spec refuses on that report with a fix it cannot perform.
This slice replays the measured 2026-08-14 sequence end to end and asserts it now
completes.

## Requirements

1. MUST write a precondition-refused report through the real writer, then run the
   mechanical stage against it, and assert no `QA-REPORT-SHAPE` finding.
2. MUST assert the same for a refusal whose cause is environmental, not only for a
   blocking one.
3. MUST still produce a `QA-REPORT-SHAPE` finding for a report that is genuinely
   malformed — a row with a non-terminal status, and a table whose rows a reader
   cannot parse.
4. MUST NOT construct the report by hand; the point is that the writer's own
   output is readable by the reader.

## Subtasks

- [ ] Replay the refused-then-rerun sequence through the real writer and stage.
- [ ] Cover the environmental refusal.
- [ ] Keep the genuinely-malformed cases refused.

## Acceptance Criteria

- [ ] A report written by a precondition-refused gate produces no
      `QA-REPORT-SHAPE` finding when the mechanical stage reads it.
- [ ] The same holds for an environmentally-blocked refusal.
- [ ] A report with a non-terminal row status is still refused.
- [ ] The report under test came from the writer, not from a fixture string.

## Verification

- `go test -count=1 ./internal/daemon ./internal/speccheck -run 'TestRefusedReportDoesNotBlockItsSuccessor' -v > /tmp/0113-t02.log 2>&1; s=$?; grep -q '^--- PASS: TestRefusedReportDoesNotBlockItsSuccessor' /tmp/0113-t02.log || { cat /tmp/0113-t02.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing sequence test; fails today, where no such test exists.
- `test -s /tmp/0113-t02.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestRefusedReportDoesNotBlockItsSuccessor' /tmp/0113-t02.log || { echo 'the sequence test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t02.log > /tmp/0113-t02-n.txt; test "$(cat /tmp/0113-t02-n.txt)" -ge 4 || { echo "expected the blocking refusal, the environmental refusal, and the two still-refused shapes, got $(cat /tmp/0113-t02-n.txt)"; cat /tmp/0113-t02.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving both directions are covered.
- `f=$(grep -rl 'TestRefusedReportDoesNotBlockItsSuccessor' internal/ | head -1); test -n "$f" || { echo 'the sequence test does not exist'; exit 1; }; grep -q 'writeMechanicalQAReport\|WriteMechanicalResult' "$f" || { echo "FAIL: $f builds its report by hand instead of through the writer"; exit 1; }` — expected: exits 0, proving the sequence runs against the writer's real output. Fails today.

## Context

- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → Build Order 2; Testing Approach, the successor. `_prd.md` →
Core Feature 2; Goal 1; Success Metrics. ADR-0132, ADR-0096.

## Result

### Implementation

- The daemon integration test now replays four report-writer-to-mechanical-stage
  sequences under `TestRefusedReportDoesNotBlockItsSuccessor`.
- The blocking and environmental cases pass the real writer's unchanged output
  to the mechanical stage and assert that it produces no `QA-REPORT-SHAPE`
  finding.
- The negative controls start from the same real writer output, then introduce a
  non-terminal row status or replace the `Status` table header so the reader
  cannot parse any row. Both assert a specific `QA-REPORT-SHAPE` finding.

### Focused checks

- Before the edit, repository search found no test named
  `TestRefusedReportDoesNotBlockItsSuccessor`; the existing integration seam
  covered only the blocking sequence under another name.
- `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test -v
  ./internal/daemon -run '^TestRefusedReportDoesNotBlockItsSuccessor$'
  -count=1` passed the parent test and all four named subtests.
- `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test
  ./internal/daemon -run
  '^(TestRefusedReportDoesNotBlockItsSuccessor|TestWriteMechanicalQAReportRecordsTheRefusal)$'
  -count=1` passed the sequence and adjacent writer tests.
- `rtk git diff --check` passed.
- The first focused invocation without the task-local `GOCACHE` did not compile
  because the sandbox denied Go's host cache path; rerunning with the writable
  `/tmp` cache above passed.
- The Task's declared Verification commands were not run; the Daemon owns that
  gate.

### Evidence per acceptance criterion

1. `blocking precondition refusal is readable` writes a blocking precondition
   result through `writeMechanicalQAReport`, confirms the refusal cause in the
   emitted report, and observes no `QA-REPORT-SHAPE` finding on the successor
   mechanical-stage read.
2. `environmental precondition refusal is readable` writes a result whose skipped
   external-evidence detector names missing environment access, confirms that
   cause in the emitted report, and observes no `QA-REPORT-SHAPE` finding.
3. `non-terminal row status is refused` changes the writer-emitted terminal
   status to `running` and observes a `QA-REPORT-SHAPE` finding naming the
   non-terminal status. `unreadable results table is refused` changes the emitted
   status header to `Outcome` and observes the empty-parsed-table shape finding.
4. Every subtest calls `writeMechanicalQAReport`. No positive report is assembled
   from a fixture string; only the two negative controls alter the emitted bytes
   to introduce their named malformed condition.
