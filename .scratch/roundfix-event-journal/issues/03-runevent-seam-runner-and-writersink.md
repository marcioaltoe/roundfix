---
title: "Run Event seam: runevent package, runner publication, WriterSink"
type: AFK
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
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

- [ ] Leaf module with RunEvent, Source/Kind vocabularies, `Publish(ctx, event) error` sink, fanout with critical vs best-effort policy
- [ ] Runner publishes Run Events with Run ID/Batch stamped; raw ACP payload byte-equal to what the runtime sent
- [ ] Kind mapping covers message, thought, tool_started, tool_updated, plan, status, raw
- [ ] Stopped status event published on cancellation
- [ ] WriterSink preserves the existing non-TTY output contract; formatting tests migrate to it
- [ ] Old type-assert publish path and single-consumer sink interface removed
- [ ] Fanout race tests; critical error propagates, best-effort error swallowed and never blocks

## Blocked by

- 01-bump-acp-sdk-and-bubbletea.md
