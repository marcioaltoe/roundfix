---
task: task_06
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 06: Perform the repairs the Task names

## Overview

A gate given two repairs by its own Task file on 2026-08-15 found both, wrote them
as findings, and failed. Reporting a repair the contract assigned leaves the work
to whoever reads the report. This slice makes an assigned repair work the gate
performs and then verifies, bounded to what the Task names.

## Requirements

1. MUST treat a repair the gate's Task file names as work to perform, then verify,
   rather than as a finding to report.
2. MUST bound what the gate may write to the paths its Task names; a repair is not
   licence to edit the Spec at large.
3. MUST fail when a named repair was neither performed nor verifiable, so a gate
   that skips one does not pass.
4. MUST keep reporting, not performing, anything the Task did not assign.
5. MUST record in the report what it performed, so the change is auditable rather
   than silent.

## Subtasks

- [x] Distinguish an assigned repair from an observation.
- [x] Perform, verify, and record it.
- [x] Keep unassigned findings reported.
- [x] Cover the performed, the skipped, and the unassigned cases.

## Acceptance Criteria

- [x] A gate whose Task names a repair performs it, and the report records what
      was performed.
- [x] A gate that leaves a named repair unmade fails.
- [x] A finding the Task did not assign is reported, not performed.
- [x] The gate writes nothing outside the paths its Task names.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestGatePerformsAssignedRepairs' -v > /tmp/0113-t06.log 2>&1; s=$?; grep -q '^--- PASS: TestGatePerformsAssignedRepairs' /tmp/0113-t06.log || { cat /tmp/0113-t06.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t06.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestGatePerformsAssignedRepairs' /tmp/0113-t06.log || { echo 'the repair test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t06.log > /tmp/0113-t06-n.txt; test "$(cat /tmp/0113-t06-n.txt)" -ge 4 || { echo "expected the performed, skipped, unassigned, and out-of-bounds cases, got $(cat /tmp/0113-t06-n.txt)"; cat /tmp/0113-t06.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving the negative cases run — a gate that reports instead of acting passes nothing, so the skipped case is the one that matters.
- `go test -count=1 ./internal/speccheck ./internal/daemon > /tmp/0113-t06-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0113-t06-regress.log && { echo 'the gate or its Daemon caller regressed:'; grep -B 3 -A 8 'FAIL' /tmp/0113-t06-regress.log | head -30; exit 1; }; grep -rq 'TestGatePerformsAssignedRepairs' internal/speccheck || { echo 'both packages pass, but the repair contract does not exist'; exit 1; }; exit $s` — expected: exits 0, proving the gate and the Daemon that runs it agree, anchored to the case this Task adds.

## Context

- interface: `internal/speccheck/mechanical.go`
- instruction: `.agents/skills/qa-gate/SKILL.md`

## References

`_techspec.md` → Build Order 6; Risks & Considerations, a gate that writes.
`_prd.md` → Core Feature 7; Goal 5; User Story 5; Non-Goals. ADR-0134.

## Result

Implemented a declarative assigned-repair phase in the pre-QA mechanical stage.
Each repair names an exact Task repair path and one unambiguous before/after
replacement. The stage preflights the complete repair batch before writing,
preserves file modes, reads each changed file back, and records only verified
writes under `## Performed repairs`. Work it cannot safely perform or verify is
kept separate from observations as an assigned repair failure and blocks the
mechanical result. The Daemon includes that failure in the terminal refusal row.

Unassigned detector output remains a `MechanicalFinding` and is never used as a
repair instruction. Non-canonical paths, paths absent from the Task's exact
repair-path set, symlinks and other non-regular targets, missing before/after
text, duplicate identifiers, and ambiguous replacements fail before the repair
batch writes any file.

Focused evidence:

- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-0113-task06-gocache go test ./internal/speccheck -run '^TestGatePerformsAssignedRepairs$/assigned_repair_is_performed_verified_and_recorded$'`
  failed to compile because the repair request, result, and report contract did
  not exist.
- `rtk env GOCACHE=/tmp/roundfix-0113-task06-gocache go test -count=1 ./internal/speccheck -run '^TestGatePerformsAssignedRepairs$/'`
  passed the performed-and-recorded, skipped-and-blocking, unassigned-finding,
  and out-of-bounds-no-write cases.
- `rtk env GOCACHE=/tmp/roundfix-0113-task06-gocache go test -count=1 ./internal/daemon -run '^TestWriteMechanicalQAReportRecordsTheRefusal$/'`
  passed, including the assigned-repair failure's terminal refusal provenance.
- `rtk env GOCACHE=/tmp/roundfix-0113-task06-gocache go test -count=1 ./internal/speccheck -run '^(TestMaterializeMechanicalResult|TestMechanicalReportShape)$'`
  passed, covering the changed report carrier beside the existing shape rules.
- `rtk git diff --check` exited 0.

Acceptance evidence:

- Performed and recorded: the positive subtest observes the changed file, the
  `PerformedRepair` record, and its `verified after write` report row.
- Unmade repair fails: the skipped subtest observes an unchanged file, one
  `RepairFailure`, no `MechanicalFinding`, and `Blocking: true`; the Daemon
  subtest observes the matching failed terminal row.
- Unassigned observation stays reported: the malformed-report subtest receives
  `QA-REPORT-SHAPE`, leaves the report byte-identical, and records no performed
  repair.
- Exact write bound: the out-of-bounds subtest assigns only the PRD path,
  attempts a repair to `CONTEXT.md`, observes a blocking repair failure, and
  confirms `CONTEXT.md` remains byte-identical.

The Task's declared `## Verification` commands were not run; the Daemon owns
those commands and settlement.
