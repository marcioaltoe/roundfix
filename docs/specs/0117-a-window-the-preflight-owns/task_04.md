---
status: pending
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
A Run that may outlast the window is created and says so; a Run already under
way is never interrupted by it.
