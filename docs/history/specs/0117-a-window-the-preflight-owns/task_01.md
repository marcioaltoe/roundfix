---
status: completed
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

Implemented schema 13 with an additive, idempotent v12 migration for the
`run_windows` table. Fresh databases create the table directly, and historical
schema fixtures migrate their copied working databases before read-only
consumer replay; the recorded corpus remains unchanged. The durable-table
lifecycle policy now records that only explicit force-set and clear operations
replace or remove Run Windows.

Added `SetRunWindow`, `RunWindowFor`, and `ClearRunWindow` with context-aware
SQLite operations under the existing serialized write transaction. The stored
row and returned `RunWindow` use Unix-second instants. A non-replacing set reads
and returns the standing row without issuing an insert or update.

Focused evidence:

- Schema migration and existing-row preservation:
  `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test ./internal/store -run 'Test(OpenMigratesV11RunDatabaseAddingOwnerIdentityUnproven|OpenMigratesV12RunDatabaseAddingRunWindows|RunWindowPersistsAndClearsByGitRoot|SetRunWindowPreservesExistingWindowWithoutReplace|DurableTableLifecyclePolicyCoversEveryTable|JournalConsumerCorpusReplaysEveryConsumer)$'`
  passed. `TestOpenMigratesV12RunDatabaseAddingRunWindows` proves the existing
  Run survives and the new table starts empty.
- Durable repository scope, absence, and clear semantics: the same command
  passed `TestRunWindowPersistsAndClearsByGitRoot`, which closes and reopens the
  store, isolates two Git roots, treats absence as non-error state, and proves
  clear's removed/not-removed result.
- Set asymmetry: the same command passed
  `TestSetRunWindowPreservesExistingWindowWithoutReplace`; it returns the
  standing window with `written=false`, and SQLite `total_changes()` remains
  unchanged. Its forced-set case replaces the cutoff with `written=true`.
- Existing outdated-database caller:
  `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test ./internal/cli -run '^TestBranchIntegrityPreflightMigratesOutdatedRunDatabase$'`
  passed.
- Required incremental gate: `rtk make verify-incremental` ran after the final
  edits. `internal/store` and every package except `internal/cli` passed;
  `internal/cli` failed only two force-stop integration tests because the
  sandbox denied process-table enumeration with `operation not permitted`.
  The Daemon-owned Task Verification command was not run.
