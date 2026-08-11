---
spec: 0033-console-log-tool-summary-deduplication
prd: _prd.md
created: 2026-07-16
---

# Console Log tool-summary deduplication — Technical Spec

## Executive Summary

The fix adds a stateful decorator only to the non-TTY Agent console branch of the existing Run Event fanout. It keys previously written text by ACP tool call identifier and suppresses an event only when the current `ConsoleText` bytes equal the prior bytes for that same call. The Run Event Journal remains the first critical sink and receives every event unchanged under ADR-0008. The design accepts a small amount of per-active-tool memory to avoid changing the producer, journal, `ConsoleText`, Attach, or Live Run View contracts.

## System Architecture

The existing flow remains:

```text
Agent stream → Run Event fanout → Run Event Journal
                              └→ plain-text Agent console → stderr / Console Log
```

- `internal/agent` gains a stateful plain-text console sink that composes the existing event decoder and `ConsoleText` renderer.
- `internal/cli` constructs that sink in `agentConsoleDisplaySink`; `--no-agent-console` keeps wrapping it with the existing Agent source filter.
- `internal/runevent` and `internal/store` do not change. The journal receives the event before the display decision through the existing critical fanout.
- `internal/tui` does not change. ADR-0009 keeps the Live Run View and Attach on cursor-ordered journal replay.
- The stateless `WriterSink` and `ConsoleText` functions keep their byte contracts for direct callers and TUI rendering.

No new package, event kind, configuration key, flag, schema, or dependency is introduced.

## Implementation Design

### Interfaces

```go
type ConsoleDisplaySink struct {
    Writer io.Writer
    mu sync.Mutex
    lastByTool map[string]string
}

func NewConsoleDisplaySink(io.Writer) *ConsoleDisplaySink
func (*ConsoleDisplaySink) Publish(context.Context, runevent.RunEvent) error
```

The constructor owns map initialization. `Publish` is safe for concurrent fanout calls even though current publication is serialized; the mutex protects both state and write ordering.

### Data Models

`lastByTool` maps a non-empty `RunEvent.ToolID` to the last successfully written `ConsoleText` bytes for that active tool call. It stores no raw payload, path list, diff body, terminal body, or user content beyond the already bounded visible summary.

Terminal states use the existing ACP vocabulary `completed`, `failed`, and `stopped`. Unknown states remain active so a future update can still be compared; a session-level terminal event clears the complete map as a defensive lifecycle boundary.

### API Contracts

For every Agent-source event:

1. Decode it through `StreamUpdateFromEvent`. Unknown or undecodable events keep the current no-output behavior.
2. Render it through `ConsoleText`. Empty output keeps the current no-output behavior.
3. If the event is not a tool lifecycle event or has an empty tool call identifier, write the rendered bytes unchanged.
4. For a tool event with an identifier, compare the exact rendered bytes with `lastByTool[id]`.
5. If equal, suppress the writer call. If different or absent, write once and update the map only after a successful complete write.
6. After processing a terminal tool state, delete that identifier. After a session terminal event, clear all identifiers.

The sink returns writer failures unchanged. A failed or short write never advances deduplication state, so a later retry cannot be suppressed as if it had been displayed. Suppressed events return success because the journal already owns durability and the display behavior intentionally performed no write.

The equality key is `(tool call identifier, exact rendered bytes)`. Run ID, Batch, path, tool title, and line counts remain part of the rendered/event context but cannot substitute for the producer-issued identifier. This prevents identical edits in separate tool calls from collapsing.

The journal fanout remains unfiltered. Both `agent.tool_started` and `agent.tool_updated` events retain their original kinds, summaries, tool state, tool identifier, raw payload, and cursor even when the second event produces no Console Log write.

## Coverage Map

- Goal 1 and Story 1 → `ConsoleDisplaySink` same-ID exact-byte comparison.
- Goal 2 and Story 3 → identifier-scoped state and distinct-ID tests.
- Goal 3 and Story 2 → unchanged JournalSink fanout and payload/cursor assertions.
- Goal 4 and Story 4 → non-tool passthrough, changed-summary rendering, and CLI display-sink wiring.
- Core Feature 4 → terminal-state and session-lifecycle state cleanup tests.

## Integration Points

- ACP tool lifecycle metadata already parsed into `RunEvent.ToolID` and `RunEvent.ToolState`.
- `agent.WriterSink`, `StreamUpdateFromEvent`, and `ConsoleText` remain the source of current visible bytes.
- `startRunUI` remains the only wiring point for non-TTY plain-text output and Detached Run Console Logs.
- `runevent.Fanout` continues treating both the journal and display sink as critical; a display writer error remains a command error rather than silently dropping output.

## Cross-Spec dependencies

This Spec has no functional dependency on Spec 0032, but it must follow Spec 0032 so its race and full verification gates use the deterministic Agent Session lifecycle tests. It must land before Spec 0035 because both change `internal/agent`, `internal/cli`, and Agent Run Event presentation; the order minimizes merge conflicts without making Console Log deduplication a prerequisite for profile behavior.

## Testing Approach

Extend the existing Agent console suite rather than creating a parallel renderer suite. Table tests publish real raw ACP lifecycle payloads and assert writer bytes for same-ID duplicates, distinct IDs, changed summaries, missing IDs, terminal cleanup, session cleanup, unknown events, and writer failures. Tests use the raw started/updated payload shape from the recorded finding so the regression protects the actual producer boundary.

A CLI integration test builds the existing journal-plus-display fanout, publishes the duplicate lifecycle pair, and proves two journal entries with byte-identical original payloads coexist with one plain-text line. Companion cases prove distinct tool calls remain visible twice and `--no-agent-console` still hides all Agent output without affecting the journal.

The full Run Event, Attach, and TUI suites run unchanged. A saved dogfood fixture or bounded replay test uses cursors 6923/6924's event shape without copying sensitive raw content and asserts the same outcome. `go test -race ./...` is mandatory because the new sink owns mutable state.

## Build Order

1. Failing regression fixture for the same-ID started/updated duplicate and distinct-ID negative companion.
2. Stateful `ConsoleDisplaySink` with exact-byte comparison, terminal cleanup, and writer-error semantics (depends on: 1).
3. Non-TTY and Detached Run display wiring through `agentConsoleDisplaySink` (depends on: 2).
4. Journal-preservation, Agent-suppression, Attach, TUI, and race verification (depends on: 2, 3).

## Risks & Considerations

- ACP Runtimes might reuse a tool identifier incorrectly. Terminal cleanup bounds the effect, and changed visible text still renders; distinct identifiers remain the required identity contract.
- A tool call that never reports a terminal state retains one bounded string until the session ends. Session cleanup bounds that leak without introducing timers or eviction heuristics.
- Deduplicating in `ConsoleText` would also alter TUI replay and erase the identifier context needed for safe state. The decorator stays at the display sink.
- Deduplicating in the Run Event Journal would violate ADR-0008 and prevent future renderers from interpreting the lifecycle pair differently. Journal writes remain untouched.
- Exact bytes intentionally do not normalize whitespace or paths. This may leave near-duplicates visible, but it avoids hiding content that is not proven identical.

## Decisions

- ADR-0008 remains unchanged: all raw Agent payloads are journaled.
- ADR-0009 remains unchanged: Attach and the Live Run View consume the journal, not this sink.
- ADR-0030 remains unchanged: the Detached Run Console Log stays unconditional and becomes less repetitive only at render time.
- Deduplication uses the first visible summary, same tool call identifier, and exact rendered bytes; it never buffers until terminal completion.
