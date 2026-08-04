---
task: task_02
spec: 0070-declared-unreachable-acceptance
status: pending
type: backend
complexity: low
---

# Task 02: Count declared-unreachable rows as their own cause

## Overview

ADR-0080 types blocked rows by cause and counts each separately. This slice adds
the third count, `rows_blocked_declared`, read from QA report frontmatter beside
the two that exist. It is a new cause, never a reclassification: folding causes
together to make a report look cleaner is the dishonesty this Spec removes.

## Requirements

1. MUST read `rows_blocked_declared` from QA report frontmatter with the same
   rules the existing counts follow: absent means zero, present must be a
   non-negative integer scalar, anything else is a typed report error.
2. MUST expose the count to callers alongside the existing two, so the archive
   boundary can reason about all three.
3. MUST NOT change what any existing verdict means. `pass` still requires
   `rows_blocked_finding` to be zero, and no count is folded into another.
4. MUST leave a report with no `rows_blocked_declared` field behaving exactly as
   it does today, so every archived report stays readable.

## Subtasks

- [ ] Read and validate the third count.
- [ ] Expose all three counts to callers.
- [ ] Add fixtures: absent, zero, positive, and invalid.

## Acceptance Criteria

- [ ] A report with `rows_blocked_declared: 3` exposes three.
- [ ] A report without the field exposes zero and no error.
- [ ] A report with a negative or non-integer value returns a typed error
      naming the field.
- [ ] A `pass` verdict with a positive `rows_blocked_finding` is still
      rejected, unchanged.
- [ ] Every existing QA report in the archived corpus still reads without
      error, asserted over the corpus rather than one fixture.

## Context

- interface: `internal/spec/qa.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'QA|Verdict|Blocked' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the verdict and count tests ran and passed.
- `go test ./internal/spec -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 3 and 8.
- `_techspec.md` → Data Models; Build Order 2.
- ADR-0080.
