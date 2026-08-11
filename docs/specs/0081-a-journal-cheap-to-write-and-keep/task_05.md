---
task: task_05
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: data
complexity: medium
---

# Task 05: Read only what the reader uses

## Overview

Every read of the journal selects the payload column, including the reads that
discard it. A payload-free projection lands beside the existing full read
rather than replacing it, so no consumer changes shape and each one keeps
exactly the columns it actually uses.

The asymmetry matters and must be honoured rather than flattened: the
`events` stream is daemon-only and hard-fails on an empty daemon payload,
while agent payloads serve human timeline rendering. A projection that treated
all payloads alike would break the stream contract supervisors depend on.

## Requirements

1. MUST add a header-only read path projecting cursor, batch, source, kind,
   summary, and creation time, leaving the existing full read untouched.
2. MUST move to the header path only those consumers that provably discard the
   payload, naming each one in the Task Result.
3. MUST keep every consumer that reads payload fields on the full read: the
   `events` stream's daemon payload fields, timeline rendering of agent
   payloads, the cockpit's task and verification parsing, the attach capacity
   replay, and the reconcile replay probe that matches on payload equality.
4. MUST assert consumer non-regression against a journal recorded before the
   change, so the proof is a corpus replay rather than an argument.
5. MUST NOT change the stream schema, its categories, or any command's output
   shape.

## Subtasks

- [ ] Add the header projection beside the full read.
- [ ] Move only the provably payload-free consumers.
- [ ] Replay the recorded corpus through every consumer.

## Acceptance Criteria

- [ ] The header path exists and the full path is unchanged.
- [ ] Every consumer that needs payload still uses the full read.
- [ ] A pre-change journal replays identically through `events`, the timeline,
      the cockpit parsers, attach replay, reconcile, and `gc`.
- [ ] The stream schema and command outputs are byte-identical.

## Context

- interface: internal/store/journal.go
- interface: internal/runevent/stream.go

## Verification

- `output="$(go test -count=1 ./internal/store -run 'HeaderProjection|EventHeaders' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the projection tests are selected and pass.
- `output="$(go test -count=1 ./internal/... -run 'ConsumerCorpus|ReplayCorpus' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the recorded-corpus replay exists and passes.
  — expected: exit 0; every reader package stays green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 4; Goal 3; User Story 5.
- `_techspec.md` → Implementation Design (header projection); Integration
  Points; Build Order 5.
- ADR-0009, ADR-0008.
