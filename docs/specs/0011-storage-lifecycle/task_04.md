---
task: task_04
spec: 0011-storage-lifecycle
status: pending
type: backend
complexity: low
---

# Task 04: Review Issue title hygiene, status-poll dedup, merge-readiness note

## Overview

Three small output-shaping fixes bundled as one slice: strip markup and emoji
from CodeRabbit table-fragment Review Issue titles, print the watch status-poll
line only when it changes, and add the documentation expectation to the
merge-readiness `missing` path. Each is verifiable on its own with no config
surface.

## Requirements

1. MUST strip Review Source markdown (table pipes, heading/emoji markup) and
   surrounding whitespace from CodeRabbit Review Issue titles so titles store as
   plain text.
2. MUST dedup the watch status-poll stderr line so an unchanged status prints
   once, not once per interval.
3. MUST append the documentation expectation that explains the merge-readiness
   `missing` state to its stderr note, pointing at the next useful action.
4. MUST keep all other output byte-stable and MUST NOT alter stdout for the
   requested command output.

## Subtasks

- [ ] Title derivation strips markup/emoji in the CodeRabbit source
- [ ] Status-poll writer prints on change only
- [ ] Merge-readiness `missing` note names the docs expectation
- [ ] Tests: table-fragment inputs → plain titles; repeated poll → one line; missing-note text

## Acceptance Criteria

- [ ] A CodeRabbit table-fragment issue yields a plain-text title with no pipes, heading markers, or emoji.
- [ ] A stable watch poll emits one status line across many intervals, not one per interval.
- [ ] The merge-readiness `missing` output names the documentation expectation and the next action.
- [ ] stdout is unchanged for requested command output; the changes are stderr/artifact-only.

## Verification

- `rtk go test ./internal/reviewsource/... ./internal/cli/` — expected: title, poll, and merge-readiness tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 7-8; Core Features 4-6. `_techspec.md` → Title hygiene
and poll dedup, Build Order 4. Work-plan findings R3-2, R3-3, R3-4.
