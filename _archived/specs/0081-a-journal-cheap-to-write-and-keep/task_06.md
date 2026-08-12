---
task: task_06
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: frontend
complexity: medium
---

# Task 06: Make cockpit cost track new events

## Overview

The cockpit rescans a Run's entire journal from cursor zero, payloads
included, on every data-version change — which is essentially after every
append. On the 42,000-event Run from the measurement that is a full re-read
per poll. The batch-clock refresh in the same file already does the right
thing with a forward-only cursor; the task journal refresh simply never
learned it.

## Requirements

1. MUST advance a forward-only cursor for the task journal refresh, in the
   pattern the batch-clock refresh already uses in the same model.
2. MUST use the header projection for the fields it renders from headers, and
   the full read only where a payload field is actually parsed.
3. MUST keep every rendered line byte-identical to today's output for the same
   event stream, since ADR-0009 makes the cockpit the live view and a
   rendering change is a user-visible change.
4. MUST keep the existing summary fallback behaviour when a payload field is
   missing.
5. MUST NOT change the stream contract, the poll trigger, or any keybinding.

The declared Verification names `TestCockpitRefreshCostTracksNewEvents`, which does not exist yet, so it can
fail before the work. Create it to assert that refresh cost tracks new events rather than total events. A broad pattern over
this package matches cases that already pass and would approve the Task
before it starts.

## Subtasks

- [x] Give the task journal refresh a forward-only cursor.
- [x] Use the header projection where no payload field is read.
- [x] Prove rendering identical against the recorded snapshots.

## Acceptance Criteria

- [ ] Refresh cost tracks new events rather than total events.
- [ ] Rendered output is unchanged against the existing snapshot fixtures.
- [ ] The summary fallback still applies when payload fields are absent.

## Context

- interface: internal/tui/cockpit.go
- interface: internal/tui/timeline.go

## Verification

- `output="$(go test -count=1 ./internal/tui -run '^TestCockpitRefreshCostTracksNewEvents$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the cockpit tests are selected and pass.
- `output="$(go test -count=1 ./internal/tui -run 'ForwardCursor' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; a named test proves the refresh advances rather than
  rescanning.
  — expected: exit 0; the snapshot fixtures confirm rendering is unchanged.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 4; Goal 3; User Story 4.
- `_techspec.md` → System Architecture (the cockpit); Build Order 6.
- ADR-0009.

## Result

The Task journal refresh now folds forward from a cursor the way
`refreshBatchClocks` already did: one poll costs the events that arrived since
the last poll, kinds come from the payload-free header projection, and only
the two daemon kinds whose payload fields the fold parses are read whole.

### What changed (behavior)

- **`internal/tui/cockpit.go` — `cockpitModel.taskJournalCursor`**: a new
  forward-only cursor beside the existing `batchTimeCursor`.
  `refreshTaskJournalEvents` no longer wipes `taskJournalStates` and no longer
  restarts at cursor zero. The per-Task phase fold is a left fold over events
  in cursor order, so folding each event exactly once lands on the same state
  a full rescan produced — including the sticky terminal phases that ignore a
  late duplicate.
- **Header projection where no payload is parsed**: the refresh pages
  `RunEventHeadersAfter` (already on `CockpitSource` from Task 05) and reads
  the whole event only when `foldsIntoTaskJournal(header.Kind)` is true —
  `daemon.task` and `daemon.verification`, the two kinds whose payload fields
  (`task`, `work_item`, `phase`, `status`) and `review_issue` column the fold
  actually parses. Agent output, the bulk of a large journal, is now never
  loaded with its payload by this path. `readTaskJournalEvent` fetches that one
  event through the unchanged `RunEventsAfter(ctx, runID, cursor-1, 1)`.
- **Error path**: a failed whole read returns with the cursor still behind the
  unread event, so the next poll retries it instead of folding a gap. The
  earlier code reset the phases to empty at the top of every refresh, so a
  mid-scan read error briefly blanked every Task row; phases now survive it.
- **Cursor rewind**: when the Task row count changes, `refreshTasks`
  reallocates the state slices and rewinds `taskJournalCursor` to zero, so the
  new rows are rebuilt from the journal's start.
- Unchanged: `RunEventsAfter`/`RunEventHeadersAfter` themselves, the
  `CockpitSource` interface, the `poll()` data-version trigger, every
  keybinding, `timeline.go`, and the `events` stream contract.

### Evidence per Acceptance Criterion

- **Refresh cost tracks new events rather than total events.**
  `TestCockpitRefreshCostTracksNewEvents`
  (`internal/tui/cockpit_forward_cursor_test.go`) builds the same Run twice —
  once over a 20-event backlog, once over 400 — folds the opening replay, then
  appends the identical two events (one agent line, one `daemon.task`
  settlement) and measures the next refresh through a recording source. Both
  refreshes cost `{headerRows:2 fullRows:1 payloadBytes:…}`: two header rows
  for two new events and one whole read for the one folded kind, identical
  across a 20× difference in journal size. The opening fold is asserted to
  page the backlog exactly once. `go test ./internal/tui -run
  'TestCockpitRefreshCostTracksNewEvents|ForwardCursor' -v` → 3 passed.
  `TestCockpitTaskJournalRefreshUsesForwardCursorAndHeaderProjection` pins the
  mechanism: the refresh issues one header read from the tail cursor (32, not
  0), reads whole only the `daemon.verification` event and none of the three
  agent lines beside it, advances the cursor to the new tail, and leaves the
  render untouched on an idle poll that reads nothing.
- **Rendered output is unchanged against the existing snapshot fixtures.**
  The ten golden files under `internal/tui/testdata/cockpit_snapshots/` are
  untouched (`git status` shows only `cockpit.go`, the new test file, and this
  task file) and `TestCockpitRenderSnapshots` passes against them.
  `TestCockpitSpecRunInterleavedTaskReplayMatchesLivePolling` — the existing
  proof that incremental polling renders byte-identically to a full replay —
  passes, as does `TestCockpitSpecRunTaskSettlementResistsStaleAndReplayedEvents`
  (terminal-phase stickiness across polls) and the whole `internal/tui` suite.
- **The summary fallback still applies when payload fields are absent.**
  `TestCockpitTaskJournalForwardCursorKeepsSummaryFallback` folds
  payload-less `daemon.task` events ("Task task_01 settled completed.", "Task
  task_02 skipped.") with both Task files left `pending` on disk, so the
  summary is the only possible source, and asserts the rows render
  `[done] Completed` and `[skip] Skipped`.

### Focused checks

- `go test -count=1 ./internal/tui -run
  'TestCockpitRefreshCostTracksNewEvents|ForwardCursor' -v` → 3 passed.
- `go test -count=1 ./internal/tui ./internal/store` → 443 passed.
- `go build -buildvcs=false ./...` and `go vet ./internal/tui/` → clean.
  (`-buildvcs=false` only because the Task Worktree's git dir is not the build
  VCS root; without it `go build` errors before compiling anything.)

### Gate proven fail-closed

Each new assertion was proven to fail against the behavior it guards, by
sabotaging the implementation and re-running before restoring it:

1. Pre-change fold restored (rescan from zero over `RunEventsAfter`) → both
   new tests FAIL (`expected the opening fold to page 22 headers, read 0`).
2. Header projection kept but the cursor rewound to 0 each refresh → cost test
   FAILS on the cost equality itself (`headerRows:24` vs `404`), forward-cursor
   test FAILS with `expected one header read from the forward cursor 32, got
   [0]`.
3. Forward cursor kept but every kind read whole → FAILS `expected exactly one
   whole read, for the folded event at cursor …`.
4. `applyTaskJournalSummary` removed from `parseTaskJournalEvent` → the
   fallback test FAILS (`expected task_01 row marker "[done] Completed"`).

`internal/tui/cockpit.go` was restored from a pre-sabotage copy after each
run, and the final tree is the implementation described above.

### Notes for later Tasks

- The whole read is one query per folded daemon event. That is bounded by the
  daemon events in one poll (typically zero or one), not by journal size, but
  a burst of many `daemon.verification` events in a single poll costs one
  small indexed query each. If that ever shows up, the fix is a header field
  or a bounded range read, not a return to the rescan.
