---
task: task_03
spec: 0028-settlement-and-reporting
status: pending
type: backend
complexity: low
---

# Task 03: Normalize task status synonyms before validation

## Overview

Stop voiding finished work over wording: task-file parsing normalizes a small documented set of unambiguous status synonyms to the canonical vocabulary before the allowed-status gate, and the Daemon rewrites the frontmatter to canonical form when a reload finds a normalized value. Everything outside the synonym set still fails exactly as today.

## Requirements

1. MUST normalize exactly this documented set: `done` → `completed`, plus hyphen and space variants of the canonical statuses (for example `in-progress` and `in progress` → `in_progress`); the mapping is exact-match after trimming, no fuzzy matching.
2. MUST apply normalization in the task parse/reload path ahead of the existing allowed-status rejection, so a synonym never fails a Task.
3. MUST rewrite the task file to the canonical status through the existing byte-preserving status rewrite helper when the Daemon's post-agent reload normalized a synonym.
4. MUST keep rejecting statuses outside the canonical and synonym sets with the existing diagnostics naming the allowed values.
5. MUST document the synonym set where the status vocabulary is defined.

## Subtasks

- [ ] Add the normalization function with its documented map
- [ ] Apply it in task parsing/reload before validation
- [ ] Daemon reload rewrites normalized values to canonical in the task file
- [ ] Table tests: every synonym, every canonical value, rejected garbage; engine test proving a `done` task file settles completed and ends up rewritten canonical

## Acceptance Criteria

- [ ] A task file with `status: done` reloads as completed and the file afterward reads `status: completed`
- [ ] A task file with `status: finished` still fails with the unsupported-status diagnostics
- [ ] Canonical statuses round-trip untouched
- [ ] The full test suite passes

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/spec/task.go`
- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q "NormalizeStatus" internal/spec/spec.go` — expected: exit 0
- `go test ./internal/spec/... ./internal/daemon/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 2, User Story 2, Core Feature 2, Decisions; `_techspec.md` → Build Order 3, Interfaces (NormalizeStatus), Decisions (synonym set).
