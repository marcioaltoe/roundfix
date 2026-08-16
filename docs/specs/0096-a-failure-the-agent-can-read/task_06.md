---
task: task_06
spec: 0096-a-failure-the-agent-can-read
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 06: Name the surface a Task file was read from

## Overview

Recovery output names the Task and not where its file was read from, so a fix is
applied once and then discovered by trial to be needed in a second place. Naming
the surface on the same line costs nothing and removes the trial.

## Requirements

1. MUST name, in recovery output that identifies a Task file, the surface it was
   read from.
2. MUST keep the Task identifier on the same line, so one line answers both
   questions.
3. MUST leave every other field of that output unchanged.

## Subtasks

- [ ] Carry the surface to the recovery output.
- [ ] Render it beside the Task.

## Acceptance Criteria

- [ ] Recovery output naming a Task file also names its surface.
- [ ] Both appear on one line.
- [ ] No other field changed.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestRecoveryNamesTheSurface' -v > /tmp/0096-t06.log 2>&1; s=$?; grep -q '^--- PASS: TestRecoveryNamesTheSurface' /tmp/0096-t06.log || { cat /tmp/0096-t06.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0096-t06.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t06.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0096-t06.log > /tmp/0096-t06-n.txt; test "$(cat /tmp/0096-t06-n.txt)" -ge 2 || { echo "expected the named-surface case and the unchanged-fields case, got $(cat /tmp/0096-t06-n.txt)"; cat /tmp/0096-t06.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving both assertions run.
- `go test -count=1 ./internal/daemon > /tmp/0096-t06-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0096-t06-regress.log && { echo 'the Daemon regressed:'; grep -B 3 -A 8 'FAIL' /tmp/0096-t06-regress.log | head -30; exit 1; }; grep -rq 'TestRecoveryNamesTheSurface' internal/daemon || { echo 'the package passes, but the case does not exist'; exit 1; }; exit $s` — expected: exits 0, proving the package still passes with the change in place, anchored to the case this Task adds.

## Context

- interface: `internal/daemon/agent_session_owner.go`

## References

`_techspec.md` → Build Order 6. `_prd.md` → Core Feature 4.
