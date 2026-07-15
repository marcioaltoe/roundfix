---
task: task_01
spec: 0028-settlement-and-reporting
status: completed
type: backend
complexity: medium
---

# Task 01: Record the Run owner PID and prove process liveness

## Overview

Groundwork for orphaned-lock reclamation: every Run row records the process id of its owning process at creation, and a platform-guarded liveness helper answers "does this pid provably exist". No consumer acts on the data yet — this task makes the identity and the proof primitive exist, migrated and tested.

## Requirements

1. MUST add a nullable owner-pid column to the Run rows through the store's existing versioned-migration pattern; pre-migration rows have no pid.
2. MUST record the calling process id on the Run creation request at every Run creation call site (foreground and detached both create the Run from the owning process).
3. MUST provide a process-liveness helper where only a definitive "no such process" result reports dead; permission errors and non-unix platforms report alive (no proof), per the refined ADR-0044 decision.
4. MUST prove the migration: opening a database created at the previous schema version upgrades it and leaves old rows pid-less.

## Subtasks

- [x] Schema migration adding the owner-pid column, following the existing column-add precedent
- [x] Thread the owner pid through the Run creation request and its call sites
- [x] Implement the liveness helper with unix signal-0 semantics and a no-proof fallback for other platforms
- [x] Tests: liveness against the current process (alive) and a reaped child process (dead); migration from the prior schema version; created Runs carry the pid

## Acceptance Criteria

- [x] A newly created Run row carries the creating process's pid
- [x] The liveness helper reports the test process alive and a spawned-then-reaped child dead
- [x] A database at the previous schema version opens, migrates, and keeps old rows pid-less
- [x] The full test suite passes

## Context

- interface: `internal/store/store.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/detach.go`

## Verification

- `grep -q "owner_pid" internal/store/store.go` — expected: exit 0 (column exists)
- `go test ./internal/store/... ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 1, Core Feature 1; `_techspec.md` → Build Order 1, Interfaces (CreateRunRequest, ProcessAlive), Data Models (runs table); ADR-0044.

## Result

Implemented owner PID groundwork for Runs:

- Added schema version 8 with nullable `runs.owner_pid`; migration v7→v8 preserves legacy rows with `NULL` owner PID.
- Added `CreateRunRequest.OwnerPID` and persisted it on Run creation; review and implement production call sites pass `os.Getpid()`, including detached children.
- Added `store.ProcessAlive(pid)` with Unix signal-0 semantics and a non-Unix no-proof fallback that reports alive.

Evidence:

- Pre-change signal: `rtk grep "owner_pid" internal/store/store.go` exited 1 before implementation.
- Owner PID persistence: `TestCreateRunPersistsOwnerPIDAcrossRunQueries` passed via `rtk go test ./internal/store -run 'Test(ProcessAliveReportsCurrentProcessAlive|ProcessAliveReportsReapedChildDead|CreateRunPersistsOwnerPIDAcrossRunQueries|OpenMigratesV7RunDatabaseAddingOwnerPID)$'`.
- Liveness helper: `TestProcessAliveReportsCurrentProcessAlive` and `TestProcessAliveReportsReapedChildDead` passed in the focused store test command.
- Previous-schema migration: `TestOpenMigratesV7RunDatabaseAddingOwnerPID` passed in the focused store test command and asserted legacy `owner_pid` remains `NULL`.
- Task verification: `rtk grep -q "owner_pid" internal/store/store.go && rtk go test ./internal/store/... ./internal/cli/... && rtk go build -buildvcs=false ./...` exited 0.
- Full gate: `rtk make verify` exited 0 (`go test ./...`: 1197 passed in 19 packages; `roundfix skills check`: passed; build: passed).
