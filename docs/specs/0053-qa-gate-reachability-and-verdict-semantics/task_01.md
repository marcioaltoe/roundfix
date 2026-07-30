---
task: task_01
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: pending
type: backend
complexity: medium
---

# Task 01: Distinguish environment-blocked rows in the verdict contract

## Overview

A QA report today carries one `verdict` scalar, so a row the environment made
unreachable is indistinguishable from a row a real finding blocked. The gate
therefore cannot legitimately reach `pass` whenever any journey needs something
the environment cannot supply. Add typed blocked-cause counts to the report
frontmatter and enforce their consistency with the verdict, so an
environment-blocked `pass` is readable and a finding-blocked `pass` is refused.

## Requirements

1. MUST read two optional frontmatter scalars, `rows_blocked_environment` and
   `rows_blocked_finding`, each defaulting to `0` when absent.
2. MUST refuse a report as unreadable when `verdict: pass` is accompanied by
   `rows_blocked_finding` greater than zero, using the existing QA report error
   type.
3. MUST accept `verdict: pass` with `rows_blocked_environment` greater than
   zero; Roundfix validates count consistency only, never evidence quality
   (ADR-0080).
4. MUST refuse a negative or non-integer count as unreadable rather than
   coercing it to zero.
5. MUST leave `partial` and `fail` handling, the verdict vocabulary, the
   Daemon's settlement mapping, and the archive `pass` gate unchanged.
6. MUST keep tolerating unknown frontmatter keys.

## Subtasks

- [ ] Parse and validate the two counts in the QA verdict reader.
- [ ] Cover the consistency matrix: finding-blocked `pass` refused,
      environment-blocked `pass` accepted, malformed counts refused, absent
      counts defaulting to zero.
- [ ] Add the missing `partial` case to the archive gate matrix alongside a
      `pass`-with-environment-counts case.

## Acceptance Criteria

- [ ] A report with `verdict: pass` and `rows_blocked_finding: 1` is reported
      unreadable, and `roundfix archive` refuses it.
- [ ] A report with `verdict: pass` and `rows_blocked_environment: 3` is
      readable and archivable.
- [ ] A report with neither count behaves exactly as it does today.
- [ ] A report with `rows_blocked_environment: -1` or a non-integer value is
      reported unreadable.
- [ ] The archive gate matrix covers `pass`, `partial`, and `fail`.

## Context

- interface: `internal/spec/qa.go`
- interface: `internal/spec/qa_test.go`
- interface: `internal/cli/archive_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/spec/ ./internal/cli/` — expected: pass,
  including the consistency matrix and the `partial` archive case.

## References

`_prd.md` → Goal 1, Stories 1–2, 5; `_techspec.md` → Build Order 1,
Interfaces; ADR-0080.
