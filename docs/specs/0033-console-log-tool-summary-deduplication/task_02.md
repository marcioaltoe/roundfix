---
task: task_02
spec: 0033-console-log-tool-summary-deduplication
status: pending
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

- [ ] Bound per-tool and session-level deduplication state at terminal lifecycle boundaries.
- [ ] Preserve retry behavior across writer failures and short writes.
- [ ] Exercise concurrent publication under the race detector.
- [ ] Prove journal fanout retains both duplicate lifecycle events byte-for-byte and cursor-by-cursor.
- [ ] Cover Agent-console suppression, Attach, Live Run View, and sanitized dogfood replay boundaries.
- [ ] Run the repository verification and full race gates.

## Acceptance Criteria

- [ ] A terminal tool event is compared and rendered or suppressed before its identifier is removed, and later identifier reuse begins with empty display state.
- [ ] A session terminal event clears every active identifier without changing its own existing console behavior.
- [ ] A failed or short write leaves comparison state unchanged and returns an error that the critical fanout can propagate.
- [ ] Concurrent publication completes without data races and preserves serialized writer output.
- [ ] The duplicate lifecycle pair produces two byte-identical journal payloads with distinct cursors and one Console Log summary.
- [ ] Distinct tool calls, `--no-agent-console`, Attach, and Live Run View retain their existing observable behavior.
- [ ] The reported dogfood event shape passes as a sanitized bounded regression fixture.
- [ ] The full verification and race gates pass without retries or relaxed assertions.

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
