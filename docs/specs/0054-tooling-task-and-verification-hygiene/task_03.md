---
task: task_03
spec: 0054-tooling-task-and-verification-hygiene
status: pending
type: backend
complexity: medium
---

# Task 03: Refuse executable files in Daemon commits

## Overview

A build artifact is never legitimate Work Item output, but the Daemon's
snapshot-diff staging has no file-mode filter, so a stray binary in the
worktree reached a Task commit once already. Drop executable files at the
same seam that already drops repository-external and symlink-crossing paths,
so every commit surface — Batch, Task, and QA — is covered by one guard.

## Requirements

1. MUST drop a changed path from staging when it is a regular file carrying
   any execute permission, reporting the refusal with the path and its mode.
2. MUST apply to every commit the Daemon creates, by attaching to the shared
   stageable-path filter rather than to one call site.
3. MUST keep the existing drop classes and their reporting behavior
   unchanged, and MUST keep committing the remaining paths exactly as today.
4. MUST preserve the contract that an explicitly selected tracked path is
   still staged even when an ignore rule matches it.
5. SHOULD word the refusal so an operator can tell a build artifact from a
   deliberately executable repository file.

## Subtasks

- [ ] Add the executable-mode drop to the shared stageable-path filter.
- [ ] Report the refused path and mode through the existing dropped-path
      channel.
- [ ] Cover the guard for Task, Batch, and QA commit paths.

## Acceptance Criteria

- [ ] A Task whose worktree gained an executable file commits its other
      changes and reports the executable as refused with its path and mode.
- [ ] The repository-external and symlink-crossing drops keep their current
      behavior and wording.
- [ ] An explicitly selected tracked path matched by an ignore rule is still
      staged.
- [ ] A Task with no executable change produces a byte-identical commit to
      today.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/daemon_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/daemon/` — expected: pass, including the new executable-drop coverage and the unchanged ignore-override contract.

## References

`_prd.md` → User Story 6, Core Feature 6; `_techspec.md` → Build Order 4,
Interfaces: executable staging drop; ADR-0045.
