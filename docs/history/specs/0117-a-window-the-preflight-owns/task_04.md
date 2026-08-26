---
status: completed
type: backend
---

# Task: Report a Run that may cross the window

The window traps the start, never the finish. A Run created inside the window
runs to its own terminal outcome — ending one on a clock would be a Stop
Request on a timer, which ADR-0022 already owns and this Spec declines.

## Work

- When the window is open but `now + budget.MaxRunDuration` falls past the
  cutoff, create the Run and print one line alongside the existing startup
  report, naming the cutoff, the time remaining, and the configured maximum.
- Report, never refuse: refusing here would deny the last Run of a window on a
  prediction.
- This is the implement path's first read of `budget.MaxRunDuration`, which is
  consumed only by `internal/watch` today. It is read for reporting and bounds
  nothing here; the watch package's use of it must stay byte-identical.
- Nothing reads the window after Run creation.

## References

- User Story 3: A running Run finishes
- Core Feature 4: The bound governs starting, never finishing

## Verification
- `grep -q "MaxRunDuration" internal/cli/implement.go && go test -count=1 ./internal/cli -run 'TestImplementRunWindowCrossing' 2>&1 | grep -q "^ok" && go test -count=1 ./internal/watch 2>&1 | grep -q "^ok"`

## Result

Implementation:

- The implement preflight prepares a crossing report when the already-read Run
  Window is open and `now + budget.MaxRunDuration` is after its cutoff. The
  report is emitted only after Run creation and names the cutoff, remaining
  time, and configured maximum.
- A maximum duration ending exactly at the cutoff does not report a crossing.
  The comparison never refuses, stops, or shortens the Run.
- The implement path performs no Run Window read after Run creation, and
  `internal/watch` remains byte-identical.

Focused-check evidence:

- Red signal: `rtk go test ./internal/cli -run
  '^TestImplementRunWindowCrossing$' -count=1` failed because the crossing Run
  emitted zero reports while both cases still created their Runs.
- After implementation, the same focused command passed all 3 reported tests.
- `rtk go test ./internal/cli -run
  '^TestImplementRunWindow($|Crossing$)' -count=1` passed all 7 reported tests,
  covering closed, open, absent, crossing, and cutoff-equality behavior.
- `rtk git -c core.fsmonitor=false diff --check` exited 0, and `rtk git -c
  core.fsmonitor=false diff -- internal/watch` produced no output.
- The Daemon-owned `## Verification` command was not run in this Agent turn.
