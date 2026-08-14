---
task: task_06
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 06: Report what Roundfix left running

## Overview

A tool whose central promise is detached execution has no command that answers
what it detached that is still running. `runs list` correctly reports nothing,
because a process that outlived its Run has no Run record. The readiness
diagnostic is where a fact about the machine belongs.

## Requirements

1. MUST report processes Roundfix started that no live Run record owns, each with
   its age, its CPU time, and its originating Run when the Run Database still
   knows it.
2. MUST say plainly when there is nothing to report, rather than printing an
   empty table.
3. MUST report what it could not read rather than claiming an empty inventory it
   did not establish.
4. MUST NOT inspect a process outside Roundfix's own spawn lineage.
5. MUST NOT change the diagnostic's exit status; residue is surfaced, not
   refused.
6. MUST NOT write any Run record or settle any Task status.

## Subtasks

- [ ] Read the process table for Roundfix's own lineage.
- [ ] Join against live Run records.
- [ ] Render the reported, empty, and partial cases.

## Acceptance Criteria

- [ ] A process with no live Run record is reported with age and CPU time.
- [ ] A process whose Run is still live is not reported as residue.
- [ ] With nothing to report, the diagnostic says so in words.
- [ ] An unreadable process table reports what it could not read.
- [ ] The diagnostic's exit status is unchanged in every case.
- [ ] No Run record is written.

## Verification

- `go test -count=1 ./internal/cli -run 'TestDoctorReportsProcessResidue' -v > /tmp/0103-t06.log 2>&1; s=$?; grep -q '^--- PASS: TestDoctorReportsProcessResidue' /tmp/0103-t06.log || { cat /tmp/0103-t06.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t06.log` — expected: exits 0, refusing a vacuous run.
- `grep -c '^--- PASS' /tmp/0103-t06.log > /tmp/0103-t06-n.txt; test "$(cat /tmp/0103-t06-n.txt)" -ge 4 || { echo "expected the reported, live-Run, empty, and unreadable cases as their own cases, got $(cat /tmp/0103-t06-n.txt)"; cat /tmp/0103-t06.log; exit 1; }` — expected: exits 0, proving each case runs separately rather than as one combined assertion.
- `go build -buildvcs=false -o /tmp/0103-t06-roundfix ./cmd/roundfix && /tmp/0103-t06-roundfix doctor > /tmp/0103-t06-doctor.log 2>&1; grep -qi 'residue' /tmp/0103-t06-doctor.log || { echo 'the built diagnostic does not report residue:'; cat /tmp/0103-t06-doctor.log; exit 1; }` — expected: exits 0, proving the check reaches the built command rather than only its tests. Fails today.

## Context

- interface: `internal/cli/doctor.go`

## References

`_techspec.md` → Build Order 6; Implementation Design, Interfaces; Risks &
Considerations, the partial process table. `_prd.md` → Core Feature 5; Goal 4;
User Story 4; Success Metrics. ADR-0127, ADR-0014.
