---
title: "Journal as critical sink in Resolve Runs"
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
  - 8
  - 9
  - 10
  - 11
  - 12
  - 13
  - 15
  - 16
blocked_by:
  - 04-tui-implements-run-event-sink.md
  - 05-run-event-journal-store.md
---

# Journal as critical sink in Resolve Runs

## Parent

.scratch/roundfix-event-journal/03-journaled-agent-stream-prd.md

## What to build

The journaled ACP Agent stream PRD: register the Run Event Journal as a critical sink in the resolve path fanout, so
every Agent Run Event is durably appended while console/log adapters keep working. Agent log files remain a parallel
artifact. Journal append failure after Run start fails the Run. The stop event is journaled before later daemon
actions are skipped.

End-to-end demo: a Resolve Run with a fake Agent leaves a queryable event timeline in the Run Database after the
process exits.

## Acceptance criteria

- [ ] Resolve Runs append all Agent Run Events to the journal with expected source, kind, Batch, tool metadata, summary, and raw payload
- [ ] Existing console and log output unchanged
- [ ] Journal append error surfaces as a Run failure, not swallowed
- [ ] Canceled Agent execution journals a stopped event before returning
- [ ] Tests use fake Agent runners and temporary Run Databases; no real ACP runtimes

## Blocked by

- 04-tui-implements-run-event-sink.md
- 05-run-event-journal-store.md
