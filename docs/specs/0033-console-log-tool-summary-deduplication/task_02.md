---
task: task_02
spec: 0033-console-log-tool-summary-deduplication
status: completed
type: backend
complexity: high
---

# Task 02: Preserve lossless Run evidence across deduplicated output

## Overview

Complete the display lifecycle and prove that Console Log cleanup cannot alter durable Run evidence or replay surfaces. The sink must keep state bounded and retry-safe under terminal events, writer failures, and concurrent publication while the Run Event Journal remains lossless.

## Requirements

1. MUST release one tool call's remembered summary only after processing its terminal `completed`, `failed`, or `stopped` event.
2. MUST clear all remembered tool summaries after a session terminal event while leaving unknown tool states active.
3. MUST update remembered state only after a complete successful write; writer errors and short writes MUST remain observable and MUST NOT cause a later retry to be suppressed.
4. MUST serialize state mutation and writer access so concurrent publications preserve comparison and write ordering without data races.
5. MUST retain every original Run Event, raw payload byte, event kind, tool identifier, tool state, and journal cursor regardless of the display decision.
6. MUST keep `--no-agent-console`, Attach, and the Live Run View on their existing journal and display contracts.
7. MUST prove the sanitized dogfood lifecycle pair renders once in plain text while both events remain independently replayable from the Run Event Journal.

## Subtasks

- [x] Bound per-tool and session-level deduplication state at terminal lifecycle boundaries.
- [x] Preserve retry behavior across writer failures and short writes.
- [x] Exercise concurrent publication under the race detector.
- [x] Prove journal fanout retains both duplicate lifecycle events byte-for-byte and cursor-by-cursor.
- [x] Cover Agent-console suppression, Attach, Live Run View, and sanitized dogfood replay boundaries.
- [x] Run the repository verification and full race gates.

## Acceptance Criteria

- [x] A terminal tool event is compared and rendered or suppressed before its identifier is removed, and later identifier reuse begins with empty display state.
- [x] A session terminal event clears every active identifier without changing its own existing console behavior.
- [x] A failed or short write leaves comparison state unchanged and returns an error that the critical fanout can propagate.
- [x] Concurrent publication completes without data races and preserves serialized writer output.
- [x] The duplicate lifecycle pair produces two byte-identical journal payloads with distinct cursors and one Console Log summary.
- [x] Distinct tool calls, `--no-agent-console`, Attach, and Live Run View retain their existing observable behavior.
- [x] The reported dogfood event shape passes as a sanitized bounded regression fixture.
- [x] The full verification and race gates pass without retries or relaxed assertions.

## Context

- interface: `internal/agent/event.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/cli/runui.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/runevent/fanout.go`
- interface: `internal/store/journal.go`

## Verification

- `rtk go test ./internal/agent -run 'TestConsoleDisplaySink' -count=100` — expected: lifecycle cleanup, retry safety, and comparison behavior remain stable across repeated execution.
- `rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1` — expected: journal preservation, suppression, replay boundaries, and the sanitized dogfood pair pass.
- `rtk go test -race ./internal/agent ./internal/cli -run 'TestConsoleDisplaySink|TestAgentConsoleDisplaySink' -count=1` — expected: mutable display state and CLI fanout integration report no data races.
- `rtk go test -race ./...` — expected: all repository packages pass under the race detector.
- `rtk make verify` — expected: formatting, Go tests, context-driven setup tests and assets, skill synchronization, and build all pass.

## References

- `_prd.md` → Goal 2; User Story 2; Core Features 4-7; Success Metrics; Non-Goals / Out of Scope; Decisions.
- `_techspec.md` → System Architecture; Data Models; API Contracts 5-6; Coverage Map; Integration Points; Testing Approach; Build Order 4; Risks & Considerations; Decisions.
- `docs/adr/0008-run-event-payload-stores-raw-producer-json.md` → lossless Agent payload preservation.
- `docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md` → Attach and Live Run View replay boundary.
- `docs/adr/0030-agent-run-logs-are-opt-in.md` → unconditional Detached Run Console Log boundary.

## Result

Implemented the remaining display lifecycle safeguards and regression coverage for lossless Run evidence:

- `ConsoleDisplaySink` now clears remembered summaries after any successfully processed session-terminal status path, preserving writer-error retry behavior.
- Agent tests cover terminal tool cleanup for `completed`, `failed`, and `stopped`, session-terminal cleanup across multiple active identifiers, writer errors, short writes, critical fanout propagation, and serialized concurrent publication.
- CLI tests cover distinct tool-call visibility, `--no-agent-console` suppression, byte-identical duplicate lifecycle payloads with distinct journal cursors, one Console Log summary, and replay through the Attach/Live Run View timeline.

Evidence:

- Terminal tool cleanup: `TestConsoleDisplaySinkReleasesTerminalToolStateAfterProcessing` passed in `rtk go test ./internal/agent -run 'TestConsoleDisplaySink' -count=100`.
- Session terminal cleanup: `TestConsoleDisplaySinkClearsSessionStateAfterTerminalStatus` passed in `rtk go test ./internal/agent -run 'TestConsoleDisplaySink' -count=100`.
- Retry safety and critical fanout propagation: `TestConsoleDisplaySinkDoesNotAdvanceStateWhenWriteFails` passed in `rtk go test ./internal/agent -run 'TestConsoleDisplaySink' -count=100`.
- Concurrent publication: `TestConsoleDisplaySinkSerializesConcurrentPublish` passed in `rtk go test -race ./internal/agent ./internal/cli -run 'TestConsoleDisplaySink|TestAgentConsoleDisplaySink' -count=1`.
- Journal losslessness and sanitized dogfood replay: `TestAgentConsoleDisplaySinkUsesStatefulSinkForNonTTYAndDetachedLogWriter` passed in `rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1`.
- Distinct tool calls and Agent-console suppression: `TestAgentConsoleDisplaySinkKeepsDistinctToolCallsVisible` and `TestAgentConsoleDisplaySinkKeepsNoAgentConsoleSuppression` passed in `rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1`.
- Full gates passed:
  - `GOCACHE="/private/tmp/roundfix-go-cache-task-02" rtk go test ./internal/agent -run 'TestConsoleDisplaySink' -count=100`
  - `GOCACHE="/private/tmp/roundfix-go-cache-task-02" rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1`
  - `GOCACHE="/private/tmp/roundfix-go-cache-task-02" rtk go test -race ./internal/agent ./internal/cli -run 'TestConsoleDisplaySink|TestAgentConsoleDisplaySink' -count=1`
  - `GOCACHE="/private/tmp/roundfix-go-cache-task-02" rtk go test -race ./...`
  - `GOCACHE="/private/tmp/roundfix-go-cache-task-02" rtk make verify`

Follow-ups: none.
