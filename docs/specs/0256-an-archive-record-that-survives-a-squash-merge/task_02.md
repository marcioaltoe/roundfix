---
task: task_02
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: pending
type: backend
complexity: high
---

# Task 02: Archive-root readers read only real records and recover a squash-dropped source revision

## Overview

Two reader defects share one set of readers. A review whose diff changes a
prose `README.md` in the archive root stops before calling its provider,
because every `*.md` there is read as an Archive Record. Every reader that
opens a record's folder at `source_revision` also gives up when a squash
merge has dropped that commit, although the forge still holds it under the
pull request head ref. This Task adds the record predicate and Source Revision
Recovery in `internal/spec` and wires both into each reader of
`_techspec.md` → Reader inventory. It is verifiable on its own with temporary
repositories and a local bare repository standing in for the forge.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12
   [review reads any Markdown in the archive root](references/2026-10-12-review-reads-any-markdown-in-the-archive-root-as-an-archive-record.md).
   It adds `spec.IsArchiveRecordName` and `spec.IsArchiveRecord` as in
   `_techspec.md` → Interfaces. The review scope, the review Spec context and
   `spec.ArchivedSpecSlugs` MUST skip every file directly under an archive
   root that is not an Archive Record. A file that passes the predicate but
   fails `spec.ParseArchiveRecord` MUST still fail the reader as today.
2. MUST characterize the defect before fixing it.
   `TestReviewIgnoresAProseReadmeInTheArchiveRoot` drives a review, through
   its public entry point with a fake provider, of a diff that changes a prose
   `docs/history/specs/README.md`. The Result MUST record that on the starting
   tree it failed with
   `read review archive record: missing YAML frontmatter opening marker`.
3. MUST answer the Backlog Entry of 2026-10-12
   [an Archive Record points at a branch commit a squash merge drops](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)
   on the reader side. It adds `spec.EnsureSourceRevision`,
   `spec.PullRequestFromSubject`, `spec.SourceRecoveryFetchTimeout` and
   `spec.ErrSourceRevisionUnavailable` as in `_techspec.md` → Interfaces and
   API Contract 3. The data comes from the configured remote. The fetch MUST
   run only when the commit is missing and a number matching
   `^[1-9][0-9]{0,9}$` is known. It MUST run with `GIT_TERMINAL_PROMPT=0`,
   with no standard input and under a context bounded by
   `SourceRecoveryFetchTimeout` that cancels (kills) the process, never
   abandons it. It MUST write no ref and no `FETCH_HEAD`, and it MUST NOT push.
4. MUST call `spec.EnsureSourceRevision` before the existing availability
   check in the review Spec context, the Run causes archived-graph reader,
   `worktree` `archivedSourceAtHead` and the Spec check's
   `archivedTaskContextExists`. The remote is `watch.push_remote` (default
   `origin`) where the caller holds the Project Config, otherwise `origin`.
   When recovery fails, each reader's output, JSON and exit code MUST stay
   exactly as today.
5. MUST keep `TestRunCausesReportAnArchiveRecordWithoutItsRevision` and
   `TestLeftoversKeepTheWorktreeWithoutTheSourceRevision` passing unchanged.
6. MUST add these tests, each top-level test calling `t.Parallel()` first or
   carrying a `// Sequential: <reason>` comment:
   - in `internal/spec`: `TestOnlyANamedSchemaFileIsAnArchiveRecord`,
     `TestArchivedSpecSlugsIgnoreNonRecordMarkdown`,
     `TestSourceRevisionRecoveryFetchesThePullRequestHead`,
     `TestSourceRevisionRecoveryReadsTheNumberFromTheAddingCommit`,
     `TestSourceRevisionRecoveryWritesNoRef`,
     `TestSourceRevisionRecoveryLeavesAnUnreachableRevisionUnavailable` and
     `TestSourceRevisionRecoveryRefusesAMalformedNumber`;
   - in `internal/cli`: `TestReviewIgnoresAProseReadmeInTheArchiveRoot`,
     `TestReviewStillRefusesAMalformedArchiveRecord`,
     `TestReviewContextRecoversASquashDroppedRevision` and
     `TestRunCausesRecoverASquashDroppedRevision`;
   - in `internal/worktree`: `TestLeftoversRecoverASquashDroppedRevision`;
   - in `internal/speccheck`: `TestArchivedTaskContextRecoversASquashDroppedRevision`.

   The forge is a bare repository on disk. `refs/pull/7/head` points at a
   commit no other ref of the clone holds. The clone is made so that it lacks
   that commit, and every test checks that `git cat-file -e` fails before the
   reader runs. The no-ref test compares `git for-each-ref` output before and
   after and requires no `FETCH_HEAD`.
7. MUST prove each new gate can fail. The Result MUST record two sabotages,
   restored afterwards: the predicate bypassed in the review scope, and the
   fetch removed from `EnsureSourceRevision`. For each it names the test that
   failed.
8. MUST NOT change the Archive Record schema, `internal/spec/archive.go`, any
   archived record or any reader that takes a known slug. It MUST NOT add a
   dependency.

## Subtasks

- [ ] Characterize the README failure through the review entry point.
- [ ] Add the record predicate and use it in the three enumerating readers.
- [ ] Add Source Revision Recovery with its hazard bounds.
- [ ] Wire recovery into the four readers that open the folder.
- [ ] Add the tests and record the sabotages.

## Acceptance Criteria

- [ ] A review of a diff that changes a prose README in the archive root reaches its provider.
- [ ] A file that declares the record schema but is malformed still blocks the review.
- [ ] Each of the four readers reads the folder of a squash-dropped revision through a bare-repository forge, with the number in the record or in the adding commit's subject.
- [ ] No ref and no `FETCH_HEAD` is written, and an unreachable forge leaves today's outcome.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`
- instruction: `internal/cli/archive_record_readers_test.go`
- instruction: `internal/cli/review_scope_test.go`
- instruction: `internal/worktree/archive_record_test.go`
- instruction: `internal/spec/archive_reader.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/runs_causes.go`
- interface: `internal/runcause/report.go`
- interface: `internal/worktree/archive_reader.go`
- interface: `internal/speccheck/citations.go`
- creates: `internal/spec/source_revision.go`
- creates: `internal/spec/source_revision_test.go`
- creates: `internal/spec/archive_record_identity_test.go`
- creates: `internal/cli/review_archive_root_markdown_test.go`
- creates: `internal/cli/source_revision_recovery_test.go`
- creates: `internal/worktree/archive_source_recovery_test.go`
- creates: `internal/speccheck/archived_context_recovery_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestOnlyANamedSchemaFileIsAnArchiveRecord|TestArchivedSpecSlugsIgnoreNonRecordMarkdown|TestSourceRevisionRecoveryFetchesThePullRequestHead|TestSourceRevisionRecoveryReadsTheNumberFromTheAddingCommit|TestSourceRevisionRecoveryWritesNoRef|TestSourceRevisionRecoveryLeavesAnUnreachableRevisionUnavailable|TestSourceRevisionRecoveryRefusesAMalformedNumber)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestOnlyANamedSchemaFileIsAnArchiveRecord TestArchivedSpecSlugsIgnoreNonRecordMarkdown TestSourceRevisionRecoveryFetchesThePullRequestHead TestSourceRevisionRecoveryReadsTheNumberFromTheAddingCommit TestSourceRevisionRecoveryWritesNoRef TestSourceRevisionRecoveryLeavesAnUnreachableRevisionUnavailable TestSourceRevisionRecoveryRefusesAMalformedNumber; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of these tests exists, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestReviewIgnoresAProseReadmeInTheArchiveRoot|TestReviewStillRefusesAMalformedArchiveRecord|TestReviewContextRecoversASquashDroppedRevision|TestRunCausesRecoverASquashDroppedRevision|TestRunCausesReportAnArchiveRecordWithoutItsRevision)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewIgnoresAProseReadmeInTheArchiveRoot TestReviewStillRefusesAMalformedArchiveRecord TestReviewContextRecoversASquashDroppedRevision TestRunCausesRecoverASquashDroppedRevision TestRunCausesReportAnArchiveRecordWithoutItsRevision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestLeftoversRecoverASquashDroppedRevision|TestLeftoversKeepTheWorktreeWithoutTheSourceRevision)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLeftoversRecoverASquashDroppedRevision TestLeftoversKeepTheWorktreeWithoutTheSourceRevision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestArchivedTaskContextRecoversASquashDroppedRevision$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestArchivedTaskContextRecoversASquashDroppedRevision" || { printf 'missing pass: TestArchivedTaskContextRecoversASquashDroppedRevision\n' >&2; exit 1; }` — expected: exit 0; before this Task the worktree and speccheck recovery tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; Goal 4; User Stories 1, 4 and 5; Core Feature 3; Core Feature 5; Success Metric 2; Success Metric 4
- [_techspec.md](_techspec.md) — Reader inventory; Interfaces; API Contract 3; Testing Approach; Build Order 2
- ADR-0261; ADR-0247; ADR-0165; ADR-0169; ADR-0259
