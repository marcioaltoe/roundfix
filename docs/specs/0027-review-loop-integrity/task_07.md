---
task: task_07
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: medium
---

# Task 07: Record terminal reasons from the engine's settle paths

## Overview

Fill the terminal-reason field from the Daemon side: whenever the engine settles a Review Issue as failed — agent error, issue left unsettled after a batch, or Verification failure including the repair attempt — the artifact records which step failed, the command, its exit status, and where the diagnostics live. Agent-provided reasons for invalid and duplicated triage are preserved when present.

## Requirements

1. MUST pass a non-empty terminal reason whenever the Daemon settles an issue as failed, naming the failed step and, when the failure came from Verification, the command, exit status, and diagnostics location.
2. MUST NOT overwrite a non-empty agent-provided terminal reason with an empty one when settling.
3. MUST keep the reason a single concise line suitable for report suffixes and Outcome Comments.
4. SHOULD source the Verification detail from the batch outcome the engine already tracks in memory.

## Subtasks

- [ ] Thread reasons through the batch-failure settle paths (agent error, post-agent unsettled sweep, verify failure, repair failure)
- [ ] Preserve pre-existing non-empty reasons on settle
- [ ] Engine tests: after a failed Verification, the assigned issues' artifacts carry the command, exit status, and diagnostics location; after an agent error, the reason names the failed step

## Acceptance Criteria

- [ ] A batch failing Verification leaves every assigned issue artifact with a terminal reason naming the command and exit status
- [ ] An issue the agent triaged with its own reason keeps that reason through Daemon settlement
- [ ] Resolved issues carry no terminal reason
- [ ] The full test suite passes

## Context

- interface: `internal/daemon/engine.go`
- interface: `internal/agent/agent.go`
- interface: `internal/rounds/rounds.go`

## Verification

- `go test ./internal/daemon/... ./internal/agent/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 5, User Story 7, Core Feature 8; `_techspec.md` → Build Order 8, Data Models (Issue artifact frontmatter), Risks (comment content sourcing).
