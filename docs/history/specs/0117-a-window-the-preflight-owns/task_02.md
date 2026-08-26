---
status: completed
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

Implemented the repository-scoped `roundfix window set|show|clear` command on
the schema-13 Run Window store. The command loads configuration and resolves
the Git root through the same `preflight.InspectGit` path as `implement`, so a
call from a nested working-tree directory addresses the repository's one
window.

Acceptance evidence:

- Command surface and Git-root resolution: root and command help advertise
  `window`; `TestWindowSetResolvesNextOccurrenceFromNestedWorktreePath` passed
  through the real CLI runner, a nested Git directory, and the real Run
  Database.
- Cutoff resolution and refusal: the focused suite passed the 23:00 to next-day
  07:00 case, the same-day-while-ahead case, and a future absolute local
  instant. `TestWindowSetRejectsPastAbsoluteCutoffWithoutStoring` exited `2`,
  named the literal cutoff, and proved no window row existed. Malformed input,
  missing input, extra arguments, and unknown subcommands also exited `2`.
- Idempotent replacement: `TestWindowSetPreservesExistingWindowUnlessForced`
  proved a second set without `--force` reports the standing cutoff while both
  stored instants remain identical; the same flow with `--force` replaced the
  cutoff.
- Read and clear states: `TestWindowShowReportsSetAndAbsentStates` proved both
  exit `0`, with the set state printing cutoff, current time, and remaining
  duration. `TestWindowClearReportsWhetherWindowExisted` proved the first clear
  reports removal and the second reports absence.
- Bound explanation: `TestWindowHelpExplainsStartAndRunBounds` proved command
  help and show output say the Run Window bounds when a Run may start and
  `budget.max_run_duration` bounds how long one may run.

Focused checks:

- Pre-change signal:
  `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test ./internal/cli -run '^TestWindowSetResolvesNextOccurrenceFromNestedWorktreePath$'`
  failed to build because the command clock dependency and Run Window command
  output did not exist.
- After implementation:
  `rtk env GOCACHE=/tmp/roundfix-task02-go-cache go test ./internal/cli -run '^TestWindow' -shuffle=on`
  passed in 1.153s after the final implementation edits.
- Required incremental gate: `rtk make verify-incremental` passed every
  reported package except `internal/cli`; its only failures were the existing
  force-stop integration tests, where the sandbox denied process-table
  enumeration with `operation not permitted`. The Run Window tests passed in
  the same package run.

The Daemon-owned `## Verification` command was not run. Task 05 owns the
already-planned `CONTEXT.md` term and user-guide updates; they remain outside
this Task's slice.
