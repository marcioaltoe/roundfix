---
task: task_02
spec: 0027-review-loop-integrity
status: completed
type: backend
complexity: low
---

# Task 02: Persist a terminal reason on Review Issue artifacts

## Overview

Review Issue artifacts currently flip status with no explanation, leaving `Decision: resolved` and `status: failed` side by side with nothing bridging them. This task adds an optional terminal-reason field to the issue artifact frontmatter and threads it through the status-setting API, so later engine work can persist why an issue ended the way it did.

## Requirements

1. MUST add an optional `terminal_reason` string field to the Review Issue artifact frontmatter and its parsed representation, empty by default and omitted or empty for resolved issues.
2. MUST extend the issue status-setting function to accept and rewrite the terminal reason alongside status and duplicate-of.
3. MUST update all existing callers to compile against the new signature, passing an empty reason where no reason is known yet.
4. MUST preserve the field through a parse → rewrite round trip without disturbing other frontmatter fields.

## Subtasks

- [x] Add the field to the frontmatter struct and the public issue struct
- [x] Extend the status-setting function signature with the reason parameter
- [x] Update existing callers (agent settle paths and any others) with empty reasons
- [x] Table-test the round trip: set status with reason, re-parse, assert field and neighbors intact

## Acceptance Criteria

- [x] Setting a status with a non-empty reason persists `terminal_reason` in the artifact and re-parsing returns it
- [x] Setting a status with an empty reason leaves no misleading stale reason behind
- [x] Existing artifacts without the field still parse
- [x] The full test suite passes

## Context

- interface: `internal/rounds/rounds.go`
- interface: `internal/agent/agent.go`

## Verification

- `grep -q "terminal_reason" internal/rounds/rounds.go` — expected: exit 0 (field exists)
- `go test ./internal/rounds/... ./internal/agent/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 5, User Story 7, Core Feature 8; `_techspec.md` → Build Order 2, Interfaces (SetIssueStatus), Data Models (Issue artifact frontmatter).

## Result

Implemented optional `terminal_reason` persistence for Review Issue artifacts. `rounds.Issue` and the artifact frontmatter now carry the field, `SetIssueStatus` rewrites it with status and `duplicate_of`, and existing callers pass an empty reason until later engine tasks can provide real terminal reasons.

Evidence:

- Pre-change signal: `grep -q "terminal_reason" internal/rounds/rounds.go` exited 1 before implementation.
- Red test: `go test ./internal/rounds -run TestSetIssueStatusRoundTripsTerminalReason` failed before implementation because `Issue.TerminalReason` and the four-argument `SetIssueStatus` signature did not exist.
- `go test ./internal/rounds -run TestSetIssueStatusRoundTripsTerminalReason`: passed, 4 tests.
- `go test ./internal/rounds/... ./internal/agent/...`: passed, 147 tests across 2 packages.
- `grep -q "terminal_reason" internal/rounds/rounds.go`: passed.
- `go build -buildvcs=false ./...`: passed.
- `make verify`: passed; `go test ./...` reported 1150 tests across 19 packages, `roundfix skills check` passed, and the Makefile build completed with `-buildvcs=false`.

Acceptance evidence:

- Non-empty reason persistence: `TestSetIssueStatusRoundTripsTerminalReason` covers failed and duplicated status rewrites with non-empty reasons and verifies re-parsed `TerminalReason`.
- Empty reason clearing: the same test seeds a stale reason, rewrites with an empty reason, and asserts the field and stale text are absent.
- Backward parse compatibility: the same test parses an artifact without `terminal_reason` before any rewrite and asserts an empty reason.
- Full suite: `make verify` passed after the implementation.

Verification note: `rg` is not installed in this execution environment, so the task-local field check uses `grep`. The build command uses `-buildvcs=false`, matching the repository Makefile, because bare `go build ./...` fails in this Roundfix worktree when Go VCS stamping probes the invalid parent `/Users/marcio/.git`.
