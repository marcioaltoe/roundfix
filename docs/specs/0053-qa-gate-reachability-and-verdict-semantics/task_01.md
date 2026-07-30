---
task: task_01
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: completed
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

## Result

Implemented typed `rows_blocked_environment` and `rows_blocked_finding`
frontmatter parsing in the QA verdict reader. Absent counts resolve to zero;
present counts must be non-negative YAML integer scalars. The reader returns
the existing `QAReportError` for malformed counts and for a `pass` report with
finding-blocked rows. Environment-blocked `pass` reports remain readable.
Unknown frontmatter keys remain tolerated.

Focused implementation evidence:

- Before the implementation, the new focused contract checks produced the
  expected red signal: the Spec matrix had 4 passing and 6 failing subtests,
  and the archive matrix had 4 passing and 2 failing subtests.
- `GOCACHE=$PWD/.gocache rtk go test -count=1 -run 'Test(QAVerdict|NewestQAReport)' ./internal/spec`
  passed 33 subtests after the final edit.
- `GOCACHE=$PWD/.gocache rtk go test -count=1 -run 'TestRunArchive' ./internal/cli`
  passed 9 subtests after the final edit.
- `rtk git diff --check` passed after the final edit.
- The Task's two declared Daemon Verification commands were not run.

Acceptance evidence:

- Finding-blocked `pass` refusal: the focused run passed
  `TestQAVerdictValidatesBlockedCounts/finding-blocked_pass_is_unreadable`,
  which asserts `QAReportError`, and
  `TestRunArchiveRefusesMissingOrNonPassingQA/finding-blocked_pass_QA_Report`,
  which asserts archive exit 2 and the unreadable-report diagnostic.
- Environment-blocked `pass`: the focused run passed
  `TestQAVerdictValidatesBlockedCounts/environment-blocked_pass_is_readable`;
  `TestRunArchiveMovesCompletedSpecAndStampsMetadata` now archives a report
  carrying `rows_blocked_environment: 3`.
- Backward compatibility: the focused run passed
  `TestQAVerdictValidatesBlockedCounts/absent_counts_default_to_zero` and the
  existing supported-verdict cases. The fixture's existing unknown `surfaces`
  key remained accepted.
- Malformed counts: the focused run passed negative and non-integer cases for
  both blocked-count fields, each requiring `QAReportError`.
- Unchanged verdict and archive semantics: focused cases keep finding-blocked
  `partial` and `fail` readable, while the archive matrix now exercises
  archivable `pass`, refused `partial`, and refused `fail`.
