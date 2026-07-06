---
task: task_03
spec: 0011-storage-lifecycle
status: pending
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

- [ ] `archive` command parsing and dispatch
- [ ] Precondition check: all tasks completed + passing QA verdict
- [ ] Metadata stamp and folder move to `_archived/`
- [ ] Tests: happy path moves and stamps; incomplete task refused; missing/failing QA refused

## Acceptance Criteria

- [ ] Archiving a Spec whose Tasks are all completed and QA passed moves its folder to `docs/specs/_archived/<slug>/` with archive metadata stamped.
- [ ] Archiving a Spec with any non-completed Task is refused with that condition named and the folder left in place.
- [ ] Archiving a Spec with no passing QA verdict is refused with that condition named.
- [ ] The command creates no Run, never pushes, and returns a stable non-zero exit code on refusal.

## Verification

- `rtk go test ./internal/cli/` — expected: the archive command tests pass.
- `rtk go run ./cmd/roundfix archive --help` — expected: concise, truthful help.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 3; Decisions. `_techspec.md` → Archive
Command, Build Order 3, Interfaces: `ArchiveRequest`. CONTEXT.md → Archive
Command. Work-plan Spec 0011 archiving.
