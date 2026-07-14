---
task: task_02
spec: 0027-review-loop-integrity
status: pending
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

- [ ] Add the field to the frontmatter struct and the public issue struct
- [ ] Extend the status-setting function signature with the reason parameter
- [ ] Update existing callers (agent settle paths and any others) with empty reasons
- [ ] Table-test the round trip: set status with reason, re-parse, assert field and neighbors intact

## Acceptance Criteria

- [ ] Setting a status with a non-empty reason persists `terminal_reason` in the artifact and re-parsing returns it
- [ ] Setting a status with an empty reason leaves no misleading stale reason behind
- [ ] Existing artifacts without the field still parse
- [ ] The full test suite passes

## Context

- interface: `internal/rounds/rounds.go`
- interface: `internal/agent/agent.go`

## Verification

- `rg -q "terminal_reason" internal/rounds/rounds.go` — expected: exit 0 (field exists)
- `go test ./internal/rounds/... ./internal/agent/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 5, User Story 7, Core Feature 8; `_techspec.md` → Build Order 2, Interfaces (SetIssueStatus), Data Models (Issue artifact frontmatter).
