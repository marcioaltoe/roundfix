---
task: task_06
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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

- [ ] Distinguish an assigned repair from an observation.
- [ ] Perform, verify, and record it.
- [ ] Keep unassigned findings reported.
- [ ] Cover the performed, the skipped, and the unassigned cases.

## Acceptance Criteria

- [ ] A gate whose Task names a repair performs it, and the report records what
      was performed.
- [ ] A gate that leaves a named repair unmade fails.
- [ ] A finding the Task did not assign is reported, not performed.
- [ ] The gate writes nothing outside the paths its Task names.

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
