---
task: task_02
spec: 0242-an-archive-that-leaves-an-archive-record
status: completed
type: backend
complexity: high
---

# Task 02: Reconcile, the pre-PR review and the queue's other readers work from the Archive Record

## Overview

After task_01, an archived Spec may be a record with no folder. The Delivery
Queue's item inspection and prerequisites, the review correction proof,
reconcile's merge evidence and Task completion, its leftover and QA
supersession checks, the pre-PR review's archived-Spec context and override
convention, and Run causes all open the folder today. Each one would
silently stop seeing the Spec. This Task moves them to `ReadArchivedSpec`.
They read the pre-archive tree from Git at `source_revision` when the
repository at hand holds it, and decide conservatively when it does not.
Existing tests keep proving the legacy folder form unchanged. It answers the
Backlog Entry "History keeps only what the Secondbrain needs" of 2026-10-06.

## Requirements

1. MUST make the Delivery Queue's `InspectItem` read `Archived` and
   `QAOverride` through `ReadArchivedSpec`. MUST make `UnmetPrerequisites`
   and `validateDeliveryPrerequisites` count a prerequisite whose record
   exists at the inspected ref or on disk. MUST make `ProveReviewCorrection`
   allow the exact record path where it allows the archived folder today.
2. MUST make reconcile use the record at the default-branch head:
   - `specArchivedAtMergedHead`, `chooseMergedHead` and
     `provenDeliveryEvidence` take the record as the archived `_prd.md`. The
     delivery commit is the first-parent commit that added the record
     (ADR-0232).
   - Task completion follows Invariant 10.
   - `dirtyPathsInArchivedSpec` and the QA supersession in
     `internal/worktree/worktree.go` read Task Context paths and QA Reports
     at `source_revision`, per Invariant 9. Without that commit they keep
     the worktree and report why.
3. MUST make the pre-PR review find an archived Spec by a record added under
   the archive root. Its PRD Decisions and TechSpec are read at
   `source_revision`, per Invariant 11, so ADR-0165's corrective park still
   fires. `reviewSpecQAOverride` reads the override from the record. The
   review diff omits the deletions under the record's `source` and keeps the
   record.
4. MUST make Run causes read an archived Task Graph and Task files at
   `source_revision`. When that commit is absent, it reports the Spec as
   archived with its record's disposition instead of `SpecsNotFound`.
5. MUST keep every existing test in `internal/cli`, `internal/worktree` and
   `internal/delivery` passing without editing it. Those tests build legacy
   folders and prove that the legacy form still reads.
6. MUST add the tests named in Verification. Each reader gets a record
   fixture whose Spec reached `main` through a squash merge, so `main` never
   held the folder, and also a fixture where `source_revision` is absent from
   the repository. The tests are:
   - in `internal/cli/archive_record_readers_test.go`: item inspection with
     an override record, prerequisites, the review correction path, the
     review's archived-Spec context with a blocking finding that parks
     `corrective-spec-required`, the override convention, the review diff
     without the removed folder, and Run causes with and without the
     revision;
   - in `internal/worktree/archive_record_test.go`: merge evidence after a
     squash, Task completion with and without a QA override Task status,
     leftovers read at `source_revision`, the worktree kept when that
     revision is absent, and QA supersession.
7. MUST NOT change `spec.Archive`, the Archive Command, the exact-retirement
   proof, the Spec check, the Spec audit or the suite guard, and MUST NOT
   touch any file under `docs/history`.

## Subtasks

- [ ] Move the queue's inspection, prerequisites and review correction to the resolver.
- [ ] Move reconcile's merge evidence, Task completion, leftovers and QA supersession.
- [ ] Move the pre-PR review's archived-Spec context, override convention and diff.
- [ ] Move Run causes.
- [ ] Add the record and missing-revision tests.

## Acceptance Criteria

- [ ] Every listed reader sees a record-only archived Spec after a squash merge.
- [ ] A missing `source_revision` never releases a worktree or drops a park.
- [ ] Every existing reader test passes unchanged.

## Context

- creates: `internal/cli/archive_record_readers_test.go`
- creates: `internal/worktree/archive_record_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_review_correction.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/review_conventions.go`
- interface: `internal/cli/runs_causes.go`
- interface: `internal/runcause/report.go`
- interface: `internal/worktree/merged_head.go`
- interface: `internal/worktree/worktree.go`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`
- instruction: `docs/adr/0165-a-blocking-review-after-archive-parks-publication-for-a-corrective-spec.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestInspectItemReadsTheArchiveRecord|TestPrerequisitesCountAnArchiveRecord|TestReviewCorrectionAllowsTheArchiveRecordPath|TestReviewParksForACorrectiveSpecOnAnArchiveRecord|TestReviewOverrideConventionReadsTheArchiveRecord|TestReviewDiffOmitsTheRemovedSpecFolder|TestRunCausesReadAnArchivedTaskGraphFromGit|TestRunCausesReportAnArchiveRecordWithoutItsRevision)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestInspectItemReadsTheArchiveRecord TestPrerequisitesCountAnArchiveRecord TestReviewCorrectionAllowsTheArchiveRecordPath TestReviewParksForACorrectiveSpecOnAnArchiveRecord TestReviewOverrideConventionReadsTheArchiveRecord TestReviewDiffOmitsTheRemovedSpecFolder TestRunCausesReadAnArchivedTaskGraphFromGit TestRunCausesReportAnArchiveRecordWithoutItsRevision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eight tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestMergeEvidenceFromTheArchiveRecordAfterASquash|TestTaskCompletionFollowsTheArchiveRecord|TestLeftoversReadTaskContextAtTheSourceRevision|TestLeftoversKeepTheWorktreeWithoutTheSourceRevision|TestQASupersessionReadsTheArchiveRecord)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMergeEvidenceFromTheArchiveRecordAfterASquash TestTaskCompletionFollowsTheArchiveRecord TestLeftoversReadTaskContextAtTheSourceRevision TestLeftoversKeepTheWorktreeWithoutTheSourceRevision TestQASupersessionReadsTheArchiveRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.

## References

- `_prd.md` → Goal 4; User Story 3; Core Feature 4; Success Metric 2; Success Metric 5
- `_techspec.md` → Readers; Invariants 7, 9, 10 and 11; API Contract 5; Build Order 2
- ADR-0247; ADR-0232; ADR-0165; ADR-0229

## Result

Implemented this Task's runtime-reader slice. Task status remains Daemon-owned;
the declared Verification commands were not run, and no workspace commit, push
or Pull Request was created.

### Implementation

- Queue inspection resolves Archive Records through `ReadArchivedSpec`, reads
  their override metadata, and works when the last active Spec has left the
  Specs Root. Legacy folders retain their prior namespace-based inspection,
  including unstamped PRDs. Prerequisites count records on disk and at the
  refreshed default ref; review correction permits only the exact named
  record path alongside the legacy folder paths.
- Reconciliation recognizes the default-head record and finds its delivery
  commit through first-parent history. Task completion follows Invariant 10.
  Task Context and recorded paths come from the source revision. A missing
  source revision keeps clean and dirty worktrees, with the missing-proof
  reason placed before the slug so reason truncation preserves it.
- QA supersession reads reports at the source revision, including identical
  same-path reports whose blobs must match across the Run and source trees.
- Review discovers record-only archives, loads PRD Decisions and TechSpec
  from Git, reads override conventions from the record and omits source
  deletions while keeping the record. Missing source context stays explicitly
  skipped without losing archived identity or the corrective-Spec park.
  The original C5 sentence remains verbatim, followed by record guidance.
- Run causes reads the historical graph and regular Task blobs from the
  archive's owning repository. Missing history produces archived disposition
  metadata (`archived_specs`) instead of `specs_not_found`; corrective
  classification stays unknown. The shared manifest parser preserves legacy
  validation. Command references document the review and causes behavior.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Every listed reader sees a record-only archived Spec after a squash merge. | The eight named CLI tests and five named worktree tests were added in the required files. Their fixtures squash the record onto `main` and assert that `main` never held the Spec folder. Focused tests exercise inspection, prerequisites, correction paths, review context and parking, conventions, diff scope, causes, merge evidence, completion, leftovers and QA supersession. |
| A missing `source_revision` never releases a worktree or drops a park. | Missing-revision variants retain archive/override identity and prerequisite satisfaction, preserve the `corrective-spec-required` park, and report causes with the record disposition. `TestMergeEvidenceFromTheArchiveRecordAfterASquash` preserves clean worktrees; `TestLeftoversKeepTheWorktreeWithoutTheSourceRevision` preserves dirty worktrees and refuses apply, including leftovers inside the former Spec directory. QA supersession refuses missing-source proof. |
| Every existing reader test passes unchanged. | The final incremental check exited 0 with the complete existing `internal/cli`, `internal/worktree` and `internal/delivery` suites. Existing tests were not edited. Focused checks also re-ran legacy item inspection, retry, cause reporting and the verbatim C5 prompt contract. |

### Checks run

- `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk proxy go test ./internal/cli ./internal/worktree ./internal/runcause -run 'ArchiveRecord|SourceRevision|AfterASquash|ArchivedTaskGraphFromGit|RemovedSpecFolder|TestInspectItemReportsAnArchivedSpec|TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore|TestBuildReportsAttemptsTasksAndCorrectives|TestReviewPromptCarriesConventionC5' -count=1`
  — exit 0 on the final implementation; CLI 4.691 s, worktree 9.828 s,
  causes 0.507 s. This is a focused implementation check, not either declared
  Verification command.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk make verify-incremental`
  — exit 0 on the stable final implementation, with process-table access for
  existing owner-process tests. Formatting, vet, all Go packages, skill checks
  and the CLI build passed; CLI 284.528 s and worktree 46.764 s.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0 before the
  Result append; repeated for the final handoff diff.

The first build could not access the default Go cache; the task-scoped cache
resolved that environment restriction. An initial sandboxed broad run also
could not enumerate owner processes, and overlapping edits invalidated its
suiteguard snapshot. Those attempts are not passing evidence. Stable checks
ran after fixing legacy namespace inspection, absent legacy graphs, Git
root-relative tree reads and the C5 prompt compatibility regression.

Only this Task's code, tests, command references and Result were changed. The
pre-existing `pending` to `in_progress` diff belongs to the Daemon. No other
Task, Task Graph manifest, Archive Command, exact-retirement proof, Spec check,
Spec audit, suite guard or file under `docs/history` was edited. Task 04
continues to own the broader skill, Baseline and glossary updates.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/user-guide/commands/review.md`
- `docs/user-guide/commands/runs.md`
- `internal/cli/archive_reader.go`
- `internal/spec/archive_reader.go`
- `internal/spec/cause_graph.go`
- `internal/spec/spec.go`
- `internal/worktree/archive_reader.go`

## Carry-forward provenance

- Source Run: `run_20261006T232220Z_7c2e4a372c1dd72e`
- Source commit: `b2330535f0b1faba68e7949ab375ce15cb07bf72`
