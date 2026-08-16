---
task: task_04
spec: 0096-a-failure-the-agent-can-read
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 04: Let the vacuity refusal name its offenders

## Overview

The event that reports a Verification refused as vacuous summarises a count and
carries the offending commands under a key named `commands`, which reads as the
commands the tool ran. A reader takes the list for the wrong thing and looks for
a defect in commands that were fine.

## Requirements

1. MUST name the offending commands in the event summary rather than only their
   number.
2. MUST carry every probed command with its own verdict, under a key whose name
   says what the list holds.
3. MUST point at the probe log that settles what ran.
4. MUST leave the event's classification and phase unchanged, so existing readers
   keep working.

## Subtasks

- [ ] Name the offenders in the summary.
- [ ] Carry every command with its verdict under a self-describing key.
- [ ] Include the probe log path.

## Acceptance Criteria

- [ ] The summary names the offending commands.
- [ ] The payload lists every probed command with its verdict.
- [ ] The payload names the probe log path.
- [ ] The classification and phase values are unchanged.

## Verification

- `go test -count=1 ./internal/daemon -run 'TestVacuousPreWorkEvent' -v > /tmp/0096-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestVacuousPreWorkEvent' /tmp/0096-t04.log || { cat /tmp/0096-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0096-t04.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t04.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0096-t04.log > /tmp/0096-t04-n.txt; test "$(cat /tmp/0096-t04-n.txt)" -ge 3 || { echo "expected the summary, the per-command verdicts, and the probe log cases, got $(cat /tmp/0096-t04-n.txt)"; cat /tmp/0096-t04.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each assertion runs on its own.
- `grep -n '"commands":' internal/daemon/task_engine.go && { echo 'the payload still uses the ambiguous commands key'; exit 1; }; grep -q 'VerificationClassificationVacuous' internal/daemon/task_engine.go || { echo 'the vacuous classification was removed rather than kept'; exit 1; }; exit 0` — expected: exits 0, proving the ambiguous key is gone and the classification survives. It prints the offending line on failure. Fails today.

## Context

- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 4; System Architecture, the vacuity event.
`_prd.md` → Core Feature 5. ADR-0111.
