---
task: task_01
spec: 0033-console-log-tool-summary-deduplication
status: completed
type: backend
complexity: high
---

# Task 01: Render one summary per tool call in plain-text output

## Overview

Deliver same-tool summary deduplication through the existing non-TTY plain-text path. One tool call must render its first bounded summary immediately and suppress only later byte-identical summaries, while separate calls and changed summaries remain visible.

## Requirements

1. MUST provide the stateful plain-text display sink defined by the TechSpec without changing the stateless `WriterSink` or `ConsoleText` contracts.
2. MUST suppress output only when a tool lifecycle event has a non-empty tool call identifier and its rendered bytes exactly match the last successfully displayed bytes for that identifier.
3. MUST render changed bytes for the same identifier and byte-identical summaries from different identifiers in publication order.
4. MUST preserve existing output for tool events without an identifier, non-tool events, empty console text, and unknown or undecodable future events.
5. MUST use the stateful display sink for non-TTY command output and Detached Run Console Logs while retaining the existing Agent-source suppression wrapper.
6. MUST store only the bounded rendered console text required for comparison; raw tool bodies and payloads MUST NOT enter display state or comparison logic.

## Subtasks

- [x] Add the stateful console display sink around the existing event decoder and console renderer.
- [x] Cover the reported same-identifier started/updated duplicate with a sanitized ACP fixture.
- [x] Cover changed summaries, distinct identifiers, missing identifiers, and non-tool passthrough.
- [x] Wire the sink into the shared non-TTY display construction path.
- [x] Preserve the existing Agent-console suppression and stateless rendering contracts.

## Acceptance Criteria

- [x] A started and updated lifecycle pair with one tool call identifier and byte-identical compact edit text writes exactly one console line.
- [x] Two tool call identifiers with byte-identical compact edit text write two console lines.
- [x] Two different rendered summaries from one tool call both appear in order.
- [x] Tool events without identifiers and non-tool events retain their previous writer bytes.
- [x] Unknown, undecodable, and empty-output events remain silent without changing deduplication state.
- [x] Non-TTY command output and the Detached Run Console Log use the same stateful display behavior.
- [x] The existing `WriterSink`, `ConsoleText`, and `--no-agent-console` observable contracts remain unchanged.

## Context

- interface: `internal/agent/event.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/cli/runui.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `rtk go test ./internal/agent -run 'TestConsoleDisplaySink|TestWriterSinkRendersConsoleTextContract' -count=1` — expected: same-call duplicates collapse, distinct or changed summaries remain visible, and stateless rendering stays byte-compatible.
- `rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1` — expected: non-TTY display wiring and Agent-source suppression use the stateful sink without changing unrelated output.

## References

- `_prd.md` → Goals 1, 3, and 4; User Stories 1, 3, and 4; Core Features 1-3, 5, and 7; User Experience; Non-Goals / Out of Scope; Decisions.
- `_techspec.md` → System Architecture; Interfaces; Data Models; API Contracts 1-5; Integration Points; Testing Approach; Build Order 1-3; Risks & Considerations.
- `docs/adr/0030-agent-run-logs-are-opt-in.md` → unconditional Detached Run Console Log boundary.

## Result

Implemented the stateful plain-text Agent console display sink and wired the non-TTY display factory to use it. The sink decodes each Run Event through `StreamUpdateFromEvent`, renders with `ConsoleText`, stores only the last successfully written bounded console text per non-empty `RunEvent.ToolID`, suppresses exact same-tool byte repeats, and releases tool state on terminal tool states.

Pre-change signal:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run 'TestConsoleDisplaySink|TestWriterSinkRendersConsoleTextContract' -count=1` failed because `NewConsoleDisplaySink` did not exist.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1` failed because the non-TTY display path still printed both duplicate summaries.

Verification:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run 'TestConsoleDisplaySink|TestWriterSinkRendersConsoleTextContract' -count=1` passed: 9 tests.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run 'TestAgentConsoleDisplaySink' -count=1` passed: 3 tests.
- `GOCACHE="/private/tmp/roundfix-go-cache-task-01" rtk make verify` passed: Go formatting check, `go test ./...` with 1290 tests, setup-context checks, `roundfix skills check`, and build.

Acceptance evidence:

- Same-ID started/updated duplicate: `TestConsoleDisplaySinkDeduplicatesToolSummaries/same_identifier_exact_bytes_collapse` renders one compact edit line.
- Distinct IDs: `TestConsoleDisplaySinkDeduplicatesToolSummaries/distinct_identifiers_keep_byte-identical_summaries` renders two identical compact edit lines.
- Changed same-ID summaries: `TestConsoleDisplaySinkDeduplicatesToolSummaries/same_identifier_renders_changed_summaries_in_order` renders both summaries in order.
- Missing IDs and non-tool passthrough: `TestConsoleDisplaySinkDeduplicatesToolSummaries/missing_identifiers_and_non-tool_events_match_writer_sink_bytes` matches `WriterSink` bytes.
- Silent events: `TestConsoleDisplaySinkDeduplicatesToolSummaries/silent_events_do_not_change_deduplication_state` keeps unknown, undecodable, and empty-output events silent without advancing state.
- Non-TTY and Detached Run Console Log writer: `TestAgentConsoleDisplaySinkUsesStatefulSinkForNonTTYAndDetachedLogWriter` publishes through `startRunUI` to a file writer and journals both lifecycle events while writing one console line.
- Existing contracts: `TestWriterSinkRendersConsoleTextContract`, `TestAgentConsoleDisplaySinkKeepsWriterBytesByDefault`, and `TestAgentConsoleDisplaySinkKeepsNoAgentConsoleSuppression` passed.

Follow-ups: none.
