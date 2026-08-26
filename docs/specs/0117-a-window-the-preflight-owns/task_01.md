---
status: pending
type: data
---

# Task: Store the Run Window at schema 13

The Run Window is durable state in the central Run Database, keyed by
repository. There is no per-session state file in this repository today and
ADR-0004 already decided that Run state is central, so a file under `.git/`
would be a second store to keep honest.

## Work

- Raise `schemaVersion` from 12 to 13 and add the additive migration `case 12`
  creating `run_windows(git_root TEXT PRIMARY KEY, cutoff_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL)`. No existing row changes.
- Add `SetRunWindow(ctx, gitRoot, cutoff, replace)` returning the effective
  window and whether it wrote, `RunWindowFor(ctx, gitRoot)` returning the
  window and whether one exists, and `ClearRunWindow(ctx, gitRoot)` returning
  whether one was removed.
- `SetRunWindow` with `replace` false against an existing row writes nothing and
  returns the standing window: re-opening a closed window by accident is the
  failure this asymmetry exists to prevent.
- Scope is `git_root`, the same key `ActiveRunInGitRoot` already uses. Absence
  of a window is a state, never an error.

## References

- User Story 1: Declare when a session stops taking on new work
- Core Feature 1: A stored bound on Run creation
- Core Feature 3: Setting a cutoff does not silently move one

## Verification
- `grep -q "run_windows" internal/store/store.go && grep -q "schemaVersion = 13" internal/store/store.go && grep -q "func (store \*Store) SetRunWindow" internal/store/store.go && go test -count=1 ./internal/store 2>&1 | grep -q "^ok"`

## Result
The Run Window is durable, repository-scoped state that a restarted process
finds unchanged.
