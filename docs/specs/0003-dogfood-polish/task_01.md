---
task: task_01
spec: 0003-dogfood-polish
status: pending
type: backend
complexity: low
---

# Task 01: Normalize generated commit subjects

## Overview

Make every Daemon-generated commit read like the repository's own style: Task
commit subjects lowercase their first rune, and the QA Report commit drops its
scope. Trailers and everything else stay byte-identical. Verifiable through
the existing commit-message unit tests.

## Requirements

1. MUST lowercase the first rune of the derived Task commit subject
   (`<type>: <title>`) while leaving the task file's own title untouched;
   non-letter first runes pass through unchanged.
2. MUST change the QA Report commit subject to `docs: qa report for <slug>
   (<verdict>)` — unscoped; the `Roundfix-Spec` trailer stays.
3. MUST keep the type mapping (docs/test/chore pass through, else feat) and
   both trailers exactly as shipped.
4. MUST update every test asserting the old subjects deliberately — no other
   assertion changes.

## Subtasks

- [ ] Lowercase derivation in the Task commit subject
- [ ] Unscoped QA commit subject
- [ ] Table tests covering letter, non-letter, and unicode first runes

## Acceptance Criteria

- [ ] `TaskCommitMessage` table tests show `feat: build the acpx invocation
      core`-style subjects, including a digit-first and a unicode-first title.
- [ ] `QACommitMessage` yields `docs: qa report for <slug> (<verdict>)`.
- [ ] Both generated subjects pass `cog verify` under this repository's
      configuration (asserted in a test or recorded verbatim in the Result).
- [ ] The full suite passes with only the deliberate assertion updates.

## Verification

- `rtk go test ./internal/daemon/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Decisions. `_techspec.md` →
Interfaces, Build Order 1. Dogfood findings 1 and 9.
