---
task: task_06
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
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

## Subtasks

- [ ] Give the task journal refresh a forward-only cursor.
- [ ] Use the header projection where no payload field is read.
- [ ] Prove rendering identical against the recorded snapshots.

## Acceptance Criteria

- [ ] Refresh cost tracks new events rather than total events.
- [ ] Rendered output is unchanged against the existing snapshot fixtures.
- [ ] The summary fallback still applies when payload fields are absent.

## Context

- interface: internal/tui/cockpit.go
- interface: internal/tui/timeline.go

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test -count=1 ./internal/tui -run 'Cockpit|Journal|Cursor' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the cockpit tests are selected and pass.
- `output="$(go test -count=1 ./internal/tui -run 'ForwardCursor' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; a named test proves the refresh advances rather than
  rescanning.
- `go test -count=1 ./internal/tui/...`
  — expected: exit 0; the snapshot fixtures confirm rendering is unchanged.

## References

- `_prd.md` → Core Feature 4; Goal 3; User Story 4.
- `_techspec.md` → System Architecture (the cockpit); Build Order 6.
- ADR-0009.
