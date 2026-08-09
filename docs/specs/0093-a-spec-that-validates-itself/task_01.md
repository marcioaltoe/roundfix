---
task: task_01
spec: 0093-a-spec-that-validates-itself
status: pending
type: test
complexity: medium
---

# Task 01: Record that a false citation passes every check today

## Overview

The corpus, with an unusually good fixture available: the exact text Spec 0090's
PRD carried when its QA gate failed it, claiming ADR-0083 makes `make verify`
the authoritative gate while ADR-0083 is about adopted sources moving to their
owning Spec. This Task records that the current checker passes that text, and
that no detector reads a cited record's body at all.

## Requirements

1. MUST record that an artifact claiming a decision record establishes something
   the record does not mention passes every current check.
2. MUST use Spec 0090's original wording as the fixture, so the corpus carries
   the real defect rather than an invented one.
3. MUST record that no current detector reads the body of a cited ADR.
4. MUST record how the existing citation checks differ from this one: they
   verify a record was listed or accounted for, never that a claim about it
   holds.
5. MUST declare the break against Task 02.
6. MUST NOT change any production behaviour.

## Subtasks

- [ ] Capture the false-citation fixture from Spec 0090's original text.
- [ ] Assert it passes today.
- [ ] Record that no detector reads a cited record's body.

## Acceptance Criteria

- [ ] A test proves the false citation produces no finding today.
- [ ] A test proves the existing citation checks pass on it for their own
      reasons, which are about listing rather than support.
- [ ] The declared break names Task 02.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/citation_characterization_test.go`
- `internal/speccheck/testdata/citation/**`
- `docs/specs/0093-a-spec-that-validates-itself/task_01.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitationCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationCharacterizationFalseCitationPassesEveryCheck'` — expected: exits 0. A `-run` pattern selecting no cases exits 0, so this asserts the named case ran.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitationCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationCharacterizationNoDetectorReadsACitedRecordBody'` — expected: exits 0.
- `grep -q 'Declared break: task_02' internal/speccheck/citation_characterization_test.go` — expected: exits 0.

## References

- `_prd.md` → Goals 1 and 4.
- `_techspec.md` → Build Order 1; Testing Approach.
