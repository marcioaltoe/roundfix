---
task: task_01
spec: 0028-settlement-and-reporting
status: pending
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

- [ ] Schema migration adding the owner-pid column, following the existing column-add precedent
- [ ] Thread the owner pid through the Run creation request and its call sites
- [ ] Implement the liveness helper with unix signal-0 semantics and a no-proof fallback for other platforms
- [ ] Tests: liveness against the current process (alive) and a reaped child process (dead); migration from the prior schema version; created Runs carry the pid

## Acceptance Criteria

- [ ] A newly created Run row carries the creating process's pid
- [ ] The liveness helper reports the test process alive and a spawned-then-reaped child dead
- [ ] A database at the previous schema version opens, migrates, and keeps old rows pid-less
- [ ] The full test suite passes

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
