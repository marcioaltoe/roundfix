---
task: task_03
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: pending
type: backend
complexity: medium
---

# Task 03: A failed QA without an override becomes the failed-qa disposition

## Overview

Gives a Legacy Archive Folder that was archived with a failing QA Report and no
override stamp a valid Archive Record. The record keeps the `fail` verdict and
the report name, and it is never a pass or an override. This Task answers
report 4 of the Backlog Entry "History sanitize refuses older archived Specs
and aborts the whole plan" of 2026-10-07. It is verifiable on its own through
the spec tests below.

## Requirements

1. MUST add `ArchiveFailedQA ArchiveDisposition = "failed-qa"`. In
   `BuildArchiveRecord`, a newest QA Report with verdict `fail` and no override
   gives `failed-qa`, keeping `QAReport` and `QAVerdict` (`_techspec.md` →
   Invariant 5). An override still gives `qa-override`, and supersession and
   `no-qa` keep today's rules.
2. MUST make `ParseArchiveRecord` accept `failed-qa` only with `qa_verdict:
   fail`, a non-empty `qa_report` and no `qa_override` key. It returns the
   three messages of `_techspec.md` → Invariant 6 otherwise. Every other
   disposition keeps today's checks and messages.
3. MUST make `ArchivedTaskCompleted` return false for the QA Task of a
   `failed-qa` record and keep every other answer (`_techspec.md` →
   Invariant 7).
4. MUST leave unchanged every other reader in `_techspec.md` → "Readers that
   branch on disposition": reconcile through `ArchivedTaskCompleted`, review,
   delivery, `SC-ARCHIVE-LICENSE`, cause reports, the sanitize plan line,
   supersede and the Spec check. MUST NOT add a disposition branch to any of
   them.
5. MUST keep the Archive Command's refusal of an active Spec with a failing QA
   and no override, with its message (`_techspec.md` → Invariant 8).
6. MUST add `internal/spec/archive_failed_qa_test.go` with synthetic folders
   only:
   - `TestLegacyConversionOfAFailedQAIsFailedQA`: a `verdict: fail` report
     with no `qa_override` gives `failed-qa`, `qa_verdict: fail` and the
     report name, with no `qa_override` field.
   - `TestFailedQARecordRoundTrips`.
   - `TestFailedQARecordRefusesAnotherVerdictOrAnOverride`: one case for each
     Invariant 6 message.
   - `TestArchivedTaskCompletedExcludesTheQATaskOfAFailedQA`: the QA Task is
     false and another Task is true.
   - `TestArchiveStillRefusesAFailingQAWithoutAnOverride`: the Archive Command
     on an active Spec refuses and changes no file.
7. MUST NOT change the lenient reading from task_02, the CLI, or any existing
   test expectation.

## Subtasks

- [ ] Add the disposition and map a failed QA without an override to it.
- [ ] Guard it in the parser.
- [ ] Exclude its QA Task from completion.
- [ ] Write the spec tests.

## Acceptance Criteria

- [ ] A failed QA without an override converts to a parseable `failed-qa`
      record that keeps the verdict and the report.
- [ ] A `failed-qa` record with another verdict, without a report or with an
      override is refused.
- [ ] The Archive Command and every existing archive record test behave as
      before.

## Context

- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/archive_reader.go`
- creates: `internal/spec/archive_failed_qa_test.go`
- instruction: `internal/spec/archive_record_test.go`
- instruction: `internal/spec/archive.go`
- instruction: `internal/worktree/merged_head.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `internal/cli/review_conventions.go`
- instruction: `internal/speccheck/citations.go`
- instruction: `internal/runcause/report.go`
- instruction: `docs/adr/0251-a-legacy-archive-folder-is-read-leniently-and-a-failed-qa-keeps-its-verdict.md`
- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacyConversionOfAFailedQAIsFailedQA|TestFailedQARecordRoundTrips|TestFailedQARecordRefusesAnotherVerdictOrAnOverride|TestArchivedTaskCompletedExcludesTheQATaskOfAFailedQA|TestArchiveStillRefusesAFailingQAWithoutAnOverride)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestLegacyConversionOfAFailedQAIsFailedQA TestFailedQARecordRoundTrips TestFailedQARecordRefusesAnotherVerdictOrAnOverride TestArchivedTaskCompletedExcludesTheQATaskOfAFailedQA TestArchiveStillRefusesAFailingQAWithoutAnOverride; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the five tests exists. After it, a failed QA converts to `failed-qa`, the parser guards it, and the Archive Command refuses as before.

## References

- `_prd.md` → Core Feature 2; Goals; Success Metric 2; Success Metric 4
- `_techspec.md` → API Contract 2; Invariant 5; Invariant 6; Invariant 7; Invariant 8; Build Order 3
- ADR-0251
- ADR-0154
