---
task: task_06
spec: 0028-settlement-and-reporting
status: completed
type: backend
complexity: medium
---

# Task 06: List settle commit paths and warn on shared worktrees

## Overview

Make the Settle Command honest about what it sweeps: before the settled line, settle prints one `commit <path>` line per path included in its commit, in sorted order; and when other Tasks in the Spec are failed at settle time, one stderr warning names them and states their work may be swept into this commit. Field-proven need: a settle in this repo silently carried a sibling task's whole implementation.

## Requirements

1. MUST print one `commit <path>` stdout line per committed path, sorted, between the verification lines and the existing `settled <task> completed — <sha>` line; the existing lines are unchanged.
2. MUST emit one stderr warning naming the other failed Tasks of the Spec when any exist at settle time, stating their work may be included in this commit; no warning when none exist.
3. MUST print nothing extra when the settle creates no commit (nothing stageable).
4. MUST source the sibling-failed check from the Task Graph the command already loads.

## Subtasks

- [ ] Surface the stageable path set in the settle report before committing
- [ ] Sibling-failed warning from the loaded Task Graph
- [ ] CLI tests (buffer-captured): path lines sorted and complete in a temp repo; warning present with a second failed Task and absent without; report byte shape for the unchanged lines

## Acceptance Criteria

- [ ] A settle committing three paths prints exactly three sorted `commit <path>` lines before the settled line
- [ ] With another failed Task in the Spec, stderr names it in the warning; with none, stderr has no warning
- [ ] Existing verification and settled lines are byte-identical to today
- [ ] The full test suite passes

## Context

- interface: `internal/cli/settle.go`

## Verification

- `go test ./internal/cli/...` — expected: all tests pass, including the settle path-listing coverage
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 4, User Story 4, Core Feature 4, Decisions (transparency over pathspec restriction); `_techspec.md` → Build Order 6, API Contracts (settle report), Risks (report consumers / skill sync).

## Result

- Added settle commit reporting in `internal/cli/settle.go`: successful commits now print sorted `commit <path>` lines before the existing settled line.
- Added a shared-worktree warning sourced from the loaded Task Graph when other failed Tasks exist and a settle commit is created.
- Acceptance evidence:
  - `TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage` asserts exactly three sorted commit path lines before the settled line.
  - `TestRunSettleWarnsWhenOtherSpecTasksAreFailed` asserts stderr names `task_02` and states work may be included; existing no-sibling paths assert empty stderr.
  - `TestRunSettleNoCommitPrintsNoCommitPathsOrSharedWarning` asserts no commit lines or shared warning when nothing is stageable.
  - Exact stdout assertions keep existing verification and settled lines byte-identical around the inserted commit lines.
- Verification:
  - `rtk go test ./internal/cli/...` — passed, 423 tests.
  - `rtk go build -buildvcs=false ./...` — passed.
  - `rtk make verify` — passed: `go test ./...` 1227 tests, skill check passed, build passed.
