---
task: task_03
spec: 0011-storage-lifecycle
status: completed
type: backend
complexity: medium
---

# Task 03: Archive Command with all-completed and QA-passed precondition

## Overview

Add `roundfix archive <slug>`, the support command that retires a completed Spec:
it verifies every Task is completed and QA passed, stamps archive metadata, and
moves the Spec folder under the archived root. It enforces the precondition
deterministically so a half-finished Spec can never be archived by mistake.

## Requirements

1. MUST add a non-interactive `archive` support command taking a Spec slug.
2. MUST verify, before touching the filesystem, that every task in the Spec's
   Task Graph has `status: completed` and that the newest QA report carries a
   passing verdict; MUST refuse with an error naming the first unmet condition
   otherwise.
3. On pass, MUST stamp archive metadata (at least the archive date and source
   slug) and move `docs/specs/<slug>/` to `docs/specs/_archived/<slug>/`.
4. MUST create no Run and never push; MUST write requested output to stdout and
   diagnostics to stderr, returning a stable non-zero exit code on refusal.

## Subtasks

- [x] `archive` command parsing and dispatch
- [x] Precondition check: all tasks completed + passing QA verdict
- [x] Metadata stamp and folder move to `_archived/`
- [x] Tests: happy path moves and stamps; incomplete task refused; missing/failing QA refused

## Acceptance Criteria

- [x] Archiving a Spec whose Tasks are all completed and QA passed moves its folder to `docs/specs/_archived/<slug>/` with archive metadata stamped.
- [x] Archiving a Spec with any non-completed Task is refused with that condition named and the folder left in place.
- [x] Archiving a Spec with no passing QA verdict is refused with that condition named.
- [x] The command creates no Run, never pushes, and returns a stable non-zero exit code on refusal.

## Verification

- `rtk go test ./internal/cli/` — expected: the archive command tests pass.
- `rtk go run ./cmd/roundfix archive --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 3; Decisions. `_techspec.md` → Archive
Command, Build Order 3, Interfaces: `ArchiveRequest`. CONTEXT.md → Archive
Command. Work-plan Spec 0011 archiving.

## Result

- Added the non-interactive `roundfix archive <slug>` support command and help text. Success writes only `archived <slug> -> docs/specs/_archived/<slug>` to stdout; refusals use the existing Preflight failure path on stderr with exit code 2.
- Added `spec.Archive`, which loads the Task Graph, checks the first non-completed Task before checking QA, requires the newest QA Report verdict to be `pass`, stamps `_prd.md` with `status: archived`, `archived`, and `source_slug`, then moves the folder to `docs/specs/_archived/<slug>/`.
- Evidence for the happy path: `TestRunArchiveMovesCompletedSpecAndStampsMetadata` asserts the active Spec folder is removed, the archived folder exists, task files move with it, `_prd.md` contains archive metadata, stdout is deterministic, stderr is empty, no Run Database is created, and Run engine collaborators are never constructed.
- Evidence for incomplete Task refusal: `TestRunArchiveRefusesIncompleteTask` asserts exit code 2, empty stdout, stderr names `Task "task_02" is "pending"` and the completed-task requirement, the active folder stays in place, `_archived/<slug>` is absent, `_prd.md` is not stamped, and no Run Database is created.
- Evidence for QA refusal: `TestRunArchiveRefusesMissingOrFailingQA` covers missing QA Reports and a failing latest QA verdict; each case exits 2, names `no passing QA verdict`, leaves the active folder in place, leaves `_archived/<slug>` absent, and creates no Run Database.
- Verification: `rtk go test ./internal/cli/` passed (`274 passed in 1 packages`); `rtk go run ./cmd/roundfix archive --help` passed and printed concise archive usage; `rtk go test ./internal/spec/` passed (`43 passed in 1 packages`); `rtk go test ./...` passed (`761 passed in 17 packages`); `rtk make verify` passed (full tests, skill check, build).
