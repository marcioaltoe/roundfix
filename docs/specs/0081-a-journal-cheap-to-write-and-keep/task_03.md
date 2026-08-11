---
task: task_03
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: data
complexity: medium
---

# Task 03: Stop scanning the journal at every Run start

## Overview

The retention sweep runs at the start of every `implement`, `watch`, and
`resolve`. Its eligibility query joins the entire event table and groups by
Run with no time predicate in SQL — the cutoff is applied in Go, after the
scan — and it runs inside the write transaction. So every Run start takes the
machine-wide write lock and holds it across a full scan of the largest table,
even when nothing is eligible, which the measurement showed is the steady
state: 302 eligible Runs, zero prunable rows.

The event count that scan produces is used for reporting only.

## Requirements

1. MUST express the retention cutoff as a SQL predicate so eligibility work is
   bounded by the candidate set rather than by the table.
2. MUST move the eligibility scan out of the write transaction; only the prune
   itself needs the writer, and only when rows are actually eligible.
3. MUST either derive the reported event count cheaply from the candidate set
   or drop it from the hot path, and MUST say which in the Task Result.
4. MUST preserve the retention policy's meaning exactly: age-based,
   terminal-only, `runs` rows and Active Run locks never touched, and zero
   meaning keep everything.
5. MUST keep `roundfix gc` output shapes and exit behaviour unchanged for
   every case that has rows, and MUST state any output change for the
   no-eligible-rows case.
6. MUST NOT change the schema or vacuum anything — compaction stays the
   separate, explicitly fenced command it already is.

The declared Verification names `TestRetentionScanIsBoundedByCandidates`, which does not exist yet, so it can
fail before the work. Create it to assert that eligibility work is bounded by the candidate set rather than the event table, with no aggregate over that table inside a write transaction. A broad pattern over
this package matches cases that already pass and would approve the Task
before it starts.

## Subtasks

- [ ] Move the cutoff into SQL and the scan out of the write transaction.
- [ ] Resolve the reported event count cheaply or drop it.
- [ ] Prove policy meaning and `gc` behaviour unchanged.

## Acceptance Criteria

- [ ] Eligibility work is bounded by the candidate set, not the event table.
- [ ] No aggregate over the event table runs inside a write transaction.
- [ ] The retention policy's meaning is unchanged for every case.
- [ ] `gc` still prunes exactly the Runs it pruned before, and never touches
      `runs` rows or locks.

## Context

- interface: internal/store/journal.go
- interface: internal/cli/gc.go

## Verification

- `output="$(go test -count=1 ./internal/store -run '^TestRetentionScanIsBoundedByCandidates$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the retention tests are selected and pass.
- `output="$(go test -count=1 ./internal/store -run 'RetentionScanOutsideWriteTransaction' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; a named test proves the scan left the write transaction,
  rather than the property being asserted in prose.
  — expected: exit 0; the store and the `gc` command stay green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 2; Goal 1; User Story 1.
- `_techspec.md` → System Architecture (the retention path); Build Order 3.
- ADR-0033.
- `references/2026-08-06-event-journal-payload-economics.md`.
