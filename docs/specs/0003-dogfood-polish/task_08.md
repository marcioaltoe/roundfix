---
task: task_08
spec: 0003-dogfood-polish
status: completed
type: backend
complexity: low
---

# Task 08: Surface skipped spec folders in discovery

## Overview

Active-Spec discovery silently skips folders whose `_prd.md` is missing,
unreadable, or inactive — right for picker robustness, invisible for
debugging a typo. Discovery gains a detailed variant reporting skipped
folders with reasons, and the picker prints them as stderr diagnostics.
Verifiable through spec-package and collector tests.

## Requirements

1. MUST add a detailed discovery variant returning active Specs plus skipped
   entries `{Dir, Reason}` (missing `_prd.md`, unreadable frontmatter,
   `status` not active — each with the offending value); `_archived/` is not
   a skip, it is out of scope by definition.
2. MUST keep the existing simple discovery behavior unchanged for callers
   that want only active Specs.
3. MUST print one stderr line per skipped folder when the implement picker
   lists Specs (`skipped docs/specs/<dir>: <reason>`), leaving stdout and the
   picker rendering untouched.
4. MUST keep the no-active-Specs failure message as shipped.

## Subtasks

- [x] Detailed discovery with typed skip reasons
- [x] Picker wiring for stderr diagnostics
- [x] Table tests over broken fixture folders

## Acceptance Criteria

- [x] Fixtures for missing PRD, broken frontmatter, and archived-status
      folders produce the exact skip reasons; active Specs list unchanged.
- [x] Interactive implement over a repo with one broken folder shows the
      diagnostic on stderr and the picker still works.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/spec/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Feature 7. `_techspec.md` → Interfaces
(ListActiveDetailed), Build Order 8. Dogfood finding 6.

## Result

- Added `spec.ListActiveDetailed`, returning active Specs plus typed
  `SkippedSpec` entries. `TestListActiveDetailedReportsSkippedSpecFolders`
  covers missing `_prd.md`, broken frontmatter, `status: archived`, keeps
  `_archived/` out of skipped results, and confirms `ListActive` still returns
  the unchanged active list.
- Wired the implement Interactive Input picker to print skipped-folder
  diagnostics to stderr only, in the form `skipped docs/specs/<dir>: <reason>`.
  The collector test covers a broken spec folder, confirms the diagnostic
  appears on stderr, and confirms the picker rendering and stdout stay
  unchanged.
- Verification:
  - `rtk go test ./internal/spec/ ./internal/cli/`: passed (`Go test: 200
    passed in 2 packages`).
  - `rtk go test ./...`: passed (`Go test: 479 passed in 16 packages`).
  - `rtk make verify`: passed (`go test`, `skills check`, and `go build` all
    exited 0).
