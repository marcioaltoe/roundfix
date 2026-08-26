---
status: pending
type: backend
---

# Task: The window command

A Supervisor needs one surface to declare, read, and remove the bound.

## Work

- Add `roundfix window <set|show|clear>`, resolving the git root the way
  `implement` does.
- `set <HH:MM|YYYY-MM-DDTHH:MM> [--force]` resolves `HH:MM` to the **next
  occurrence** in local time — today when still ahead, tomorrow otherwise. The
  naive same-day comparison reports a closed window at 23:00 for an 07:00
  cutoff, killing a night session at the moment it starts. Resolution happens
  once; the stored value is an instant, so nothing re-derives it later.
- A cutoff resolving into the past is refused at set time, exit `2`, and is not
  stored.
- Without `--force`, an existing window is reported and left byte-identical,
  exit `0`. With `--force`, it is replaced.
- `show` prints the cutoff, the current time, and the remaining duration, or
  states that no window is set — exit `0` either way, because absence is a
  state. `clear` reports whether a window was there.
- The help text states that the window bounds when a Run may **start** and
  names `budget.max_run_duration` as the bound on how long one may **run**.

## References

- User Story 1: Declare when a session stops taking on new work
- Core Feature 2: The cutoff is the next occurrence
- Core Feature 3: Setting a cutoff does not silently move one

## Verification
- `grep -q "\"window\"" internal/cli/cli.go && grep -q "RunWindow" internal/cli/cli.go && go test -count=1 ./internal/cli -run 'TestWindow' 2>&1 | grep -q "^ok"`

## Result
A Supervisor sets, reads, and clears the Run Window, and a night-start cutoff
resolves forward instead of reporting a closed window.
