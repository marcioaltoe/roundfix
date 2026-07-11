---
task: task_03
spec: 0024-context-efficient-runs
status: completed
type: backend
complexity: high
---

# Task 03: Emit the Supervisor Run Event Stream

## Overview

Add a stable read-only JSONL projection that lets a Supervisor replay or follow
one explicit Run without parsing Console Log text. The slice is verifiable as a
standalone CLI contract over stored Run Event Journal fixtures and Active Run
follow transitions.

## Requirements

1. MUST implement `roundfix events <run-id>` with replay and optional `--follow` for one explicit Run.
2. MUST emit only `roundfix-events/v1` JSONL records for task-status, batch, verification, and outcome categories by default.
3. MUST accept a comma-separated `--filter` subset of those categories and reject empty or unknown values with exit `2`.
4. MUST project structured Daemon payload fields without exposing raw Agent payloads, command bodies, diagnostic paths, or internal Run Event kind vocabulary.
5. MUST keep stdout records-only, route diagnostics to stderr, drain terminal Runs, and exit `130` on follow cancellation.
6. MUST share cursor paging and the data-version follow mechanism with attach without changing attach behavior.
7. MUST fail malformed relevant Daemon payloads instead of inferring state from human summaries.

## Subtasks

- [x] Define the stable category and record projection contract.
- [x] Refactor journal replay/follow into a command-neutral component.
- [x] Add events argument and filter parsing.
- [x] Implement JSONL replay, follow, terminal drain, and cancellation.
- [x] Add malformed-payload and stdout/stderr error handling.
- [x] Cover the complete public command contract with journal fixtures.

## Acceptance Criteria

- [x] Default replay emits only the four stable categories in journal cursor order with one valid JSON object per line.
- [x] Filtered replay emits only requested categories and never emits a raw Agent payload.
- [x] `--follow` emits no duplicate at the replay boundary and exits after the terminal event is drained.
- [x] A terminal Run replays and exits immediately with status `0`.
- [x] Missing Run ID, unknown Run, and invalid filter errors emit no stdout records and use the specified exit codes.
- [x] SIGINT/SIGTERM during follow exits `130` without a stdout trailer.
- [x] Existing attach replay/follow tests remain byte-compatible.

## Verification

- `rtk go test ./internal/runevent ./internal/store ./internal/cli` - expected: projection, JSONL, filter, cursor, follow, cancellation, and attach regression tests pass.
- `rtk go run -buildvcs=false ./cmd/roundfix events --help` - expected: concise truthful help renders successfully.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/cli.go`
- interface: `internal/store/journal.go`
- interface: `internal/runevent/event.go`

## References

`_prd.md` -> User Story 3; Core Features 5-7; User Experience; Success Metrics. `_techspec.md` -> Interfaces: Supervisor projection; API Contracts: events; Build Order 3. ADR-0008.

## Result

- Default replay emits only the stable `task-status`, `batch`, `verification`, and `outcome` categories in journal cursor order as JSONL; covered by `TestEventsReplayDefaultAndFilterJSONLRecordsOnly` and `TestProjectStreamEventCoversStableCategoriesAndRedactsPayload`.
- Filtered replay emits only requested categories and skips raw Agent payloads; covered by `TestEventsReplayDefaultAndFilterJSONLRecordsOnly`.
- `--follow` shares the attach cursor/data-version follower, emits no duplicate replay-boundary record, drains the terminal outcome event, and exits `0`; covered by `TestEventsFollowDrainsTerminalWithoutDuplicateBoundary`.
- Terminal Runs replay and exit immediately with status `0`; covered by `TestEventsTerminalRunReplaysAndExitsImmediately`.
- Missing Run ID, unknown Run, and invalid filter values emit no stdout records and exit `2`; covered by `TestEventsValidationErrorsEmitNoStdout`.
- Follow cancellation exits `130` without a stdout trailer; covered by `TestEventsFollowCancellationExits130WithoutTrailer`.
- Attach replay/follow behavior remains covered by the existing attach tests in `internal/cli`, which passed with the shared follower.
- Malformed selected Daemon payloads fail instead of inferring from summaries; covered by `TestEventsMalformedRelevantPayloadFailsNoStdout` and `TestProjectStreamEventRejectsMalformedRelevantDaemonPayload`.

Verification evidence:

- `rtk go test ./internal/runevent ./internal/store ./internal/cli` passed: 452 tests passed in 3 packages.
- `rtk go run -buildvcs=false ./cmd/roundfix events --help` passed and rendered concise command help.
- `rtk make verify` passed: `rtk go test ./...` reported 1078 tests passed in 19 packages, `roundfix skills check` passed, and `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` passed.
