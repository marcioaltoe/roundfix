---
task: task_04
spec: 0024-context-efficient-runs
status: completed
type: backend
complexity: medium
---

# Task 04: Compact Agent reads and edits

## Overview

Render repetitive ACP file reads and edits as bounded one-line Console Log
summaries while preserving the original Agent Run Event payload. The slice is
verifiable with measured-shape fixtures containing 31 edits and 330 reads and
with byte-for-byte journal payload assertions.

## Requirements

1. MUST retain ACP Tool Kind, locations, and diff old/new line information in the structured stream model.
2. MUST render reads as `read <path> (N lines)` and edits as `edit <path> (+N/-M)`.
3. MUST omit raw input, raw output, file bodies, ACP JSON, and unified diffs from Console Log and Live Run View text.
4. MUST degrade incomplete structured tool metadata to a bounded marker rather than raw payload output.
5. MUST keep original Agent Run Event payload bytes lossless and must not enable acpx read suppression.
6. MUST use one compact formatter for Detached Run output, journal summaries, and the Live Run View timeline.

## Subtasks

- [x] Extend ACP parsing with structured read/edit metadata.
- [x] Calculate deterministic read and diff line counts.
- [x] Replace raw tool detail rendering with compact summaries.
- [x] Add metadata-only fallback behavior.
- [x] Preserve raw Run Event payload capture.
- [x] Add measured-shape and Live Run View regression fixtures.

## Acceptance Criteria

- [x] A 31-edit/330-read fixture produces exactly 31 edit lines and 330 read lines.
- [x] No compact output contains a file body, old/new diff text, raw tool output, or serialized ACP object.
- [x] Edit counts reflect added and removed lines from structured diff content.
- [x] A structured read reports the correct path and line count.
- [x] Missing optional metadata produces a bounded tool marker without exposing raw content.
- [x] The corresponding journal events retain their original payload bytes exactly.
- [x] Console Log and Live Run View use the same summary wording.

## Verification

- `rtk go test ./internal/agent ./internal/tui` - expected: ACP parsing, compact formatting, measured fixture, timeline, and lossless-payload tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/agent/acp_stream.go`
- interface: `internal/agent/stream.go`
- interface: `internal/agent/event.go`
- interface: `internal/tui/timeline.go`

## References

`_prd.md` -> User Story 4; Core Features 8-9; Success Metrics. `_techspec.md` -> API Contracts: Console rendering; Integration Points: acpx/ACP; Build Order 4. ADR-0008.

## Result

- ACP tool updates now retain tool kind, locations, read line counts, and edit old/new line counts in the stream model; covered by `TestStreamUpdateFromACPReadEditMetadata`.
- Reads render as `read <path> (N lines)` and edits render as `edit <path> (+N/-M)`; covered by `TestConsoleTextCompactsMeasuredReadEditFixtureAndPreservesPayload` and `TestWriterSinkRendersConsoleTextContract`.
- The measured fixture generated exactly 330 read lines and 31 edit lines with no raw file body, old/new diff text, raw tool output, serialized ACP object fields, or unified diff text; covered by `TestConsoleTextCompactsMeasuredReadEditFixtureAndPreservesPayload`.
- Incomplete structured read metadata falls back to `[TOOL] read <path> · completed` without exposing raw output; covered by `TestConsoleTextFallsBackToBoundedToolMarkerForIncompleteMetadata`.
- Journal payload bytes remain byte-identical when `newAgentRunEvent` records read/edit updates; covered by `TestConsoleTextCompactsMeasuredReadEditFixtureAndPreservesPayload`.
- Detached Run output, journal summaries, and Live Run View timeline use `agent.ConsoleText`; covered by `TestWriterSinkRendersConsoleTextContract` and `TestRunTimelineUsesConsoleTextForCompactReadEditSummaries`.
- acpx prompt args remain exact and do not include read suppression; existing `acpxPromptArgs` tests passed in the focused and full gates.
- While running the full gate, `TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks` exposed a deterministic fake-worktree copy race. The fix serializes only the fake test helper's filesystem copies so the existing scheduler test can prove independent-task isolation reliably.

Verification evidence:

- `rtk go test ./internal/agent ./internal/tui` passed: 301 tests passed in 2 packages.
- `rtk go test ./internal/daemon -run TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks -count=1` passed after the fake worktree helper fix.
- `rtk make verify` passed: `rtk go test ./...` reported 1082 tests passed in 19 packages, `roundfix skills check` passed, and `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` passed.
