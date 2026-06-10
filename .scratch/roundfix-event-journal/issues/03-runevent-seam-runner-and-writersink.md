---
title: "Run Event seam: runevent package, runner publication, WriterSink"
type: AFK
category: enhancement
state: completed
labels:
  - enhancement
  - completed
user_stories:
  - 1
  - 2
  - 3
  - 4
  - 5
  - 6
  - 7
  - 9
  - 10
  - 12
  - 13
  - 16
  - 17
blocked_by:
  - 01-bump-acp-sdk-and-bubbletea.md
---

# Run Event seam: runevent package, runner publication, WriterSink

## Parent

.scratch/roundfix-event-journal/01-run-event-seam-prd.md

## What to build

Core of the Run Event seam PRD: the new leaf `runevent` module (RunEvent type, namespaced kinds, `Sink` interface,
critical/best-effort fanout, summary-bounding helper), and the full replacement of the runner output path. The ACP
runner converts stream updates to Run Events at publish time, stamping Run ID and Batch from the execute request,
with the raw ACP `session/update` JSON as payload per ADR 0008. The exec runner emits `agent.raw` events. Runner
cancellation publishes an `agent.status` stopped event before returning. Non-TTY/headless output becomes a
WriterSink adapter preserving today's text contract; the type-asserting publish path and the old single-consumer
sink interface are deleted.

End-to-end demo: a headless resolve Run with a fake Agent prints the same text as today, while a fake sink observes
Run Events with byte-exact raw payloads.

## Acceptance criteria

- [x] Leaf module with RunEvent, Source/Kind vocabularies, `Publish(ctx, event) error` sink, fanout with critical vs best-effort policy
- [x] Runner publishes Run Events with Run ID/Batch stamped; raw ACP payload byte-equal to what the runtime sent
- [x] Kind mapping covers message, thought, tool_started, tool_updated, plan, status, raw
- [x] Stopped status event published on cancellation
- [x] WriterSink preserves the existing non-TTY output contract; formatting tests migrate to it
- [x] Old type-assert publish path and single-consumer sink interface removed
- [x] Fanout race tests; critical error propagates, best-effort error swallowed and never blocks

## Blocked by

- 01-bump-acp-sdk-and-bubbletea.md

## Comments

**2026-06-10 (agent):** Built `internal/runevent` as a stdlib-only leaf module: `RunEvent` (Run ID, Batch, source,
kind, Review Issue ref, tool id/state, bounded summary, producer time, raw JSON payload), source/kind vocabularies
(`agent.message` … `agent.raw`), `Sink` with `Publish(ctx, event) error`, `BoundSummary` (2048-byte rune-safe), and
`Fanout` (critical sinks synchronous with joined error propagation; best-effort sinks behind per-sink bounded
queues + one delivery goroutine each, drops and failures counted, `Close` drains). Raw payload fidelity: the ACP
runner now tees the agent's newline-delimited JSON-RPC stdout through an interceptor that publishes each
`session/update` with its raw params bytes verbatim (ADR 0008) before forwarding the line to the SDK — proven
byte-equal in tests; the SDK `SessionUpdate` callback is a no-op. Runner-generated events (`agent.status`,
`agent.raw`) carry small JSON payloads; both runners publish a stopped status event on cancellation (using
`context.WithoutCancel` so durable sinks still receive it). Runner interface is now
`Run(ctx, req, sink runevent.Sink)`; `publishStreamUpdate` and `StreamUpdateSink` are deleted. Console text
formatting moved to `agent.ConsoleText` (shared by `agent.WriterSink` and the TUI), so the non-TTY contract is
byte-identical; the TUI formatting test migrated to `TestWriterSinkRendersConsoleTextContract`. CLI wires
TTY → interim live-view sink (replaced by the bounded TUI sink in issue 04), non-TTY → `WriterSink`. Verification:
`rtk go vet ./...` clean, `rtk go test ./...` and `rtk go test -race ./...` 161 passed in 15 packages,
`rtk go run ./cmd/roundfix --help` green.
