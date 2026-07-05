---
task: task_01
spec: 0003-dogfood-polish
status: completed
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

- [x] Lowercase derivation in the Task commit subject
- [x] Unscoped QA commit subject
- [x] Table tests covering letter, non-letter, and unicode first runes

## Acceptance Criteria

- [x] `TaskCommitMessage` table tests show `feat: build the acpx invocation
      core`-style subjects, including a digit-first and a unicode-first title.
- [x] `QACommitMessage` yields `docs: qa report for <slug> (<verdict>)`.
- [x] Both generated subjects pass `cog verify` under this repository's
      configuration (asserted in a test or recorded verbatim in the Result).
- [x] The full suite passes with only the deliberate assertion updates.

## Verification

- `rtk go test ./internal/daemon/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Decisions. `_techspec.md` →
Interfaces, Build Order 1. Dogfood findings 1 and 9.

## Result

Implemented Daemon commit-subject normalization for this slice:

- `TaskCommitMessage` now lowercases only the first rune of the derived commit
  subject title. The task file title remains unchanged, digit-first titles pass
  through unchanged, unicode uppercase first runes lowercase, and the existing
  docs/test/chore/feat type mapping plus `Roundfix-Spec` and `Roundfix-Task`
  trailers remain intact.
- `QACommitMessage` now emits `docs: qa report for <slug> (<verdict>)` with the
  existing `Roundfix-Spec` trailer.
- Existing Daemon and CLI assertions that deliberately encoded the old subjects
  were updated to the new subjects.

Acceptance evidence:

- `TestTaskCommitMessageDerivesSubjectAndTrailers` covers
  `feat: build the acpx invocation core`, `feat: 2FA setup`, and
  `feat: über tracing`, plus docs/test/chore type pass-through and trailer
  preservation.
- `TestQACommitMessageDerivesUnscopedSubjectAndTrailer` covers
  `docs: qa report for 0003-dogfood-polish (pass)`.
- `rtk cog verify "feat: build the acpx invocation core"` passed with
  `Type: feat`, `Scope: none`.
- `rtk cog verify "docs: qa report for 0003-dogfood-polish (pass)"` passed with
  `Type: docs`, `Scope: none`.

Verification:

- `rtk go test ./internal/daemon/` passed: 45 tests.
- `rtk go test ./...` passed: 453 tests in 16 packages.
- `rtk make verify` passed: full Go suite, `roundfix skills check`, and build.

Follow-ups: none.
