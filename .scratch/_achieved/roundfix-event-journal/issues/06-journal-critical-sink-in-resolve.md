---
title: "Journal as critical sink in Resolve Runs"
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

- [x] Resolve Runs append all Agent Run Events to the journal with expected source, kind, Batch, tool metadata, summary, and raw payload
- [x] Existing console and log output unchanged
- [x] Journal append error surfaces as a Run failure, not swallowed
- [x] Canceled Agent execution journals a stopped event before returning
- [x] Tests use fake Agent runners and temporary Run Databases; no real ACP runtimes

## Blocked by

- 04-tui-implements-run-event-sink.md
- 05-run-event-journal-store.md

## Comments

**2026-06-10 (agent):** Added `store.JournalSink` (Publish → `AppendRunEvent`) and wired the resolve cycle through
`newAgentEventFanout` in the CLI: the journal is a critical sink in every resolve path (resolve command and the
watch resolver, which share `executeResolveCycle`); the Live Run View joins as best-effort when live; the non-TTY
writer sink stays critical so console ordering and write-error behavior are byte-identical to before — all 65
pre-existing CLI tests pass unchanged. The fanout closes before the console stream so queued events drain. Journal
failure path: a critical-sink append error propagates from `Fanout.Publish` through the runner's publish-error
recording into the engine, failing the Run (`TestResolveCycleFailsRunWhenCriticalJournalSinkFails`). Stop path: the
runner's stopped status event flows through the critical journal sink synchronously before return
(`TestStoppedResolveJournalsStoppedEventBeforeReturning`, end-to-end at the CLI with a temp Run Database). Durable
demo: `TestResolveJournalsAgentRunEventsDurably` runs a full Resolve Run with a fake Agent, then opens a separate
read-only Run Database connection and finds the queryable event timeline (Run ID, Batch, agent source/kind,
bounded summary, raw payload). Agent log files keep being written by the runner, untouched. Verification:
`rtk go vet ./...` clean, `rtk go test ./...` 184 passed in 15 packages, `rtk go test -race` on cli/daemon/store
96 passed, `rtk go run ./cmd/roundfix --help` green.
