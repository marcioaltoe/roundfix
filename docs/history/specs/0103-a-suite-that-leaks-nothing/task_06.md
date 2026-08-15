---
task: task_06
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

### Implementation

- The Doctor Command now reads every stored Run through the read-only Run
  Database connection, excludes every non-terminal Run before process
  inspection, and reports one line per process in a terminal Run's proven spawn
  lineage.
- The existing owner-process controller now returns start time and CPU time for
  proven lineage members. Unix reads only those PIDs through `ps`; Windows uses
  process accounting APIs. A process that disappears during inspection is
  omitted, while every other unreadable entry remains a partial diagnostic.
- Doctor renders `found`, explicit `ok (no process residue found)`, and `partial`
  cases. Neither `found` nor `partial` changes the command's exit status, and
  the diagnostic never opens the Run Database writer.

### Focused checks

- `GOCACHE=/tmp/roundfix-task06-go-cache rtk go test ./internal/store -run '^TestParsePSOwnedProcess$'` passed 5 cases covering ordinary, multi-day,
  malformed, and overflowing process accounting input.
- `GOCACHE=/tmp/roundfix-task06-go-cache rtk go test ./internal/cli -run '^TestDoctorReportsProcessResidue(Reported|ExcludesLiveRun|Empty|Unreadable|DoesNotWriteRunRecord)$'` passed 5 independently runnable cases.
- `GOCACHE=/tmp/roundfix-task06-go-cache rtk go test ./internal/cli -run 'Doctor'` passed 43 Doctor-focused cases after the output-contract update.
- `GOCACHE=/tmp/roundfix-task06-go-cache rtk go test ./internal/store` passed
  252 package cases, and `GOCACHE=/tmp/roundfix-task06-go-cache rtk go test ./internal/cli` passed 1,052 package cases.
- `GOOS=windows GOCACHE=/tmp/roundfix-task06-go-cache rtk proxy go build ./internal/store ./internal/cli` passed the production portability check.
- `GOCACHE=/tmp/roundfix-task06-go-cache rtk make verify-incremental` passed:
  formatting, the full Go suite, skill checks, and the production build all
  exited successfully.
- `rtk git diff --check` passed.

### Acceptance evidence

- Reported process: `TestDoctorReportsProcessResidueReported` observes PID 5252
  with age `2h3m4s`, CPU `7m8s`, and originating Run `run-residue`.
- Live Run exclusion: `TestDoctorReportsProcessResidueExcludesLiveRun` proves a
  non-terminal Run causes zero process-table reads and no reported PID.
- Empty inventory: `TestDoctorReportsProcessResidueEmpty` observes the words
  `no process residue found`.
- Partial inventory: `TestDoctorReportsProcessResidueUnreadable` preserves the
  readable PID and reports the Run-scoped process-table error instead of an
  empty inventory.
- Exit status: the reported, live-Run, empty, and unreadable tests all observe
  `exitOK`; residue status is not part of Doctor's failure aggregation.
- No Run writes: `TestDoctorReportsProcessResidueDoesNotWriteRunRecord` invokes
  the real default read-only reader and compares the complete Run listing before
  and after.

### Not run

- The commands under `## Verification` were not run; the Daemon owns those
  commands and Task settlement.

### Follow-up

- `CONTEXT.md` does not yet define the TechSpec's coined Process Residue term.
  Task 09 already owns the Spec's glossary check, so this diff does not widen
  into that Task.
