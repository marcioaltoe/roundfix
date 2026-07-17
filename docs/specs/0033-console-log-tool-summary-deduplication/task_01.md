---
task: task_01
spec: 0033-console-log-tool-summary-deduplication
status: pending
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

- [ ] Add the stateful console display sink around the existing event decoder and console renderer.
- [ ] Cover the reported same-identifier started/updated duplicate with a sanitized ACP fixture.
- [ ] Cover changed summaries, distinct identifiers, missing identifiers, and non-tool passthrough.
- [ ] Wire the sink into the shared non-TTY display construction path.
- [ ] Preserve the existing Agent-console suppression and stateless rendering contracts.

## Acceptance Criteria

- [ ] A started and updated lifecycle pair with one tool call identifier and byte-identical compact edit text writes exactly one console line.
- [ ] Two tool call identifiers with byte-identical compact edit text write two console lines.
- [ ] Two different rendered summaries from one tool call both appear in order.
- [ ] Tool events without identifiers and non-tool events retain their previous writer bytes.
- [ ] Unknown, undecodable, and empty-output events remain silent without changing deduplication state.
- [ ] Non-TTY command output and the Detached Run Console Log use the same stateful display behavior.
- [ ] The existing `WriterSink`, `ConsoleText`, and `--no-agent-console` observable contracts remain unchanged.

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
