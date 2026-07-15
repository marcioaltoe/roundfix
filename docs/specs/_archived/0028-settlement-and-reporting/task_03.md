---
task: task_03
spec: 0028-settlement-and-reporting
status: completed
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

- [x] Add the normalization function with its documented map
- [x] Apply it in task parsing/reload before validation
- [x] Daemon reload rewrites normalized values to canonical in the task file
- [x] Table tests: every synonym, every canonical value, rejected garbage; engine test proving a `done` task file settles completed and ends up rewritten canonical

## Acceptance Criteria

- [x] A task file with `status: done` reloads as completed and the file afterward reads `status: completed`
- [x] A task file with `status: finished` still fails with the unsupported-status diagnostics
- [x] Canonical statuses round-trip untouched
- [x] The full test suite passes

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

## Result

- Added `NormalizeStatus` with the documented exact synonym map: `done` → `completed`, `in-progress`/`in progress` → `in_progress`; pre-change signal was `rtk grep -q "NormalizeStatus" internal/spec/spec.go` exiting 1.
- Applied normalization before task-status validation and carried a `StatusNormalized` marker so daemon post-agent reloads rewrite only synonym-authored frontmatter through `spec.SetStatus`.
- Evidence: `TestTaskCycleRewritesNormalizedStatusAfterAgentReload` proves `status: done` settles completed and the file reads `status: completed`.
- Evidence: `TestReloadTaskReportsBrokenAgentEdits` and the invalid-load fixture prove `status: finished` still fails with the unsupported-status diagnostic naming `pending, in_progress, completed, failed`.
- Evidence: `TestNormalizeStatus` and `TestReloadTaskNormalizesStatusValues` cover every canonical status, every synonym, and garbage passthrough.
- Verification passed: `rtk grep -q "NormalizeStatus" internal/spec/spec.go`; `rtk go test ./internal/spec/... ./internal/daemon/...` (157 passed); `rtk go build -buildvcs=false ./...`; `rtk make verify` (1222 passed, skills check passed, build clean).
