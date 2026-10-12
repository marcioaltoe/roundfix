---
task: task_04
spec: 0261-a-cancelled-spec-leaves-a-record
status: completed
type: backend
complexity: medium
---

# Task 04: A cancelled Spec meets no prerequisite, supersedes nothing and is never retried

## Overview

Three readers treat an archived Spec as delivered work: the Supersede
Command's deliverer check, `deliver start`'s prerequisite check and
`deliver retry`. This Task makes each refuse a cancelled Spec, and proves
that the review, the Run causes report, the Spec check and the History
Sanitize Command read a repository holding a cancelled record without an
error and without rewriting it (`_techspec.md` → Reader inventory, ADR-0262).
It answers
[the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)
of 2026-10-12 on the reader side. Its fixtures write cancelled records with
task_01's renderer, so it is verifiable without the Cancel Command.

## Requirements

1. MUST refuse `roundfix supersede --by <slug>` for a Spec whose Archive
   Record is cancelled, with the message of `_techspec.md` → API Contract 4
   and exit 2, and keep accepting every other archived deliverer.
2. MUST refuse `roundfix deliver start` when a queued Spec's `requires` names
   a Spec whose Archive Record is cancelled, with the message of API
   Contract 4 and exit 2, before any queue is recorded.
3. MUST refuse `roundfix deliver retry <slug>` when the checkout's Archive
   Record for the Spec is cancelled, with the message of API Contract 4 and
   exit 2, before the engine reads or writes the queue. A Spec without a
   record keeps today's retry path.
4. MUST NOT change `UnmetPrerequisites`, the review, the Run causes report,
   the Spec check or the History Sanitize Command; each is characterized
   instead.
5. MUST create `internal/cli/cancelled_record_readers_test.go` with
   `TestSupersedeRefusesACancelledDeliverer`,
   `TestDeliverStartRefusesARequiredCancelledSpec`,
   `TestDeliverRetryRefusesACancelledSpec`,
   `TestReviewReadsACancelledArchiveRecord`,
   `TestRunCausesReportACancelledRecord` and
   `TestHistorySanitizeLeavesACancelledRecord`, and
   `internal/speccheck/cancelled_record_test.go` with
   `TestSpecCheckReadsACancelledRecordAndItsReturnedEntry`. The review test
   MUST show a diff that removes a Spec folder and adds its cancelled record
   reaching the fake provider with the folder omitted as an archived Spec
   source. The Run causes test MUST show the disposition `cancelled`. The
   sanitize test MUST show a plan that lists no unit for the record, and the
   record byte-identical after `--apply` of a batch over another unit. The
   Spec check test MUST show no finding for a returned `open` Backlog Entry
   beside the cancelled record. The retry test MUST show the queue rows
   unchanged. Every test MUST use temporary repositories and homes, and every
   new top-level test MUST call `t.Parallel()` first or carry a
   `// Sequential: <reason>` comment.
6. MUST keep `TestSupersedeAcceptsAnArchiveRecord`,
   `TestPrerequisitesCountAnArchiveRecord` and
   `TestReviewDiffOmitsTheRemovedSpecFolder` passing unchanged.
7. MUST NOT change any command's help, the Cancel Command, `internal/spec` or
   any skill or guide.

## Subtasks

- [ ] Refuse a cancelled deliverer in the Supersede Command.
- [ ] Refuse a cancelled prerequisite in `deliver start`.
- [ ] Refuse a cancelled Spec in `deliver retry` before the engine.
- [ ] Characterize review, Run causes, Spec check and History Sanitize.

## Acceptance Criteria

- [ ] Each of the three refusals exits 2 with its message and writes nothing.
- [ ] Every other reader reads a cancelled record as it reads any record.
- [ ] The kept tests pass unchanged.

## Context

- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- instruction: `docs/adr/0193-a-spec-names-its-prerequisites-and-the-owner-waits-for-them.md`
- instruction: `internal/cli/archive_record_readers_test.go`
- instruction: `internal/cli/archive_record_test.go`
- instruction: `internal/cli/review_archived_spec_test.go`
- instruction: `internal/cli/history_test.go`
- instruction: `internal/speccheck/backlog_test.go`
- instruction: `internal/runcause/report.go`
- interface: `internal/cli/supersede.go`
- interface: `internal/cli/deliver.go`
- creates: `internal/cli/cancelled_record_readers_test.go`
- creates: `internal/speccheck/cancelled_record_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSupersedeRefusesACancelledDeliverer|TestDeliverStartRefusesARequiredCancelledSpec|TestDeliverRetryRefusesACancelledSpec|TestReviewReadsACancelledArchiveRecord|TestRunCausesReportACancelledRecord|TestHistorySanitizeLeavesACancelledRecord|TestSupersedeAcceptsAnArchiveRecord|TestPrerequisitesCountAnArchiveRecord|TestReviewDiffOmitsTheRemovedSpecFolder)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSupersedeRefusesACancelledDeliverer TestDeliverStartRefusesARequiredCancelledSpec TestDeliverRetryRefusesACancelledSpec TestReviewReadsACancelledArchiveRecord TestRunCausesReportACancelledRecord TestHistorySanitizeLeavesACancelledRecord TestSupersedeAcceptsAnArchiveRecord TestPrerequisitesCountAnArchiveRecord TestReviewDiffOmitsTheRemovedSpecFolder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^TestSpecCheckReadsACancelledRecordAndItsReturnedEntry$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestSpecCheckReadsACancelledRecordAndItsReturnedEntry" || { printf 'missing pass: TestSpecCheckReadsACancelledRecordAndItsReturnedEntry\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; User Story 5; Core Feature 5; Success Metric 3; Success Metric 4
- [_techspec.md](_techspec.md) — Reader inventory; API Contract 4; Testing Approach; Build Order 4
- ADR-0262; ADR-0193; ADR-0223; ADR-0229; ADR-0237; ADR-0165; ADR-0169; ADR-0248; ADR-0251; ADR-0259

## Result

Implemented the reader slice for daemon Verification; Task status and the
Subtasks/Acceptance Criteria checkboxes remain daemon-owned.

- The Supersede Command checks the archived deliverer's disposition and
  refuses `cancelled` with API Contract 4's message and exit 2.
- `deliver start` rejects a cancelled prerequisite during prerequisite
  validation, before recording a queue or opening the Run Database.
- `deliver retry` resolves the checkout's Archive Record before opening the
  Run Database or constructing the Delivery Engine. A cancelled record
  refuses with API Contract 4's message and exit 2. An absent record continues
  through the existing retry path.
- Added the six required CLI tests and the required Spec-check test. Fixtures
  render cancelled records with `spec.RenderArchiveRecord` and use disposable
  repositories and homes. CLI top-level tests call `t.Parallel()` first; the
  Spec-check test declares why isolating process-wide HOME/XDG_CONFIG_HOME
  requires sequential execution.

Acceptance evidence:

1. Each refusal exits 2, emits its exact reason with empty stdout, and leaves
   repository files, home files and Git status unchanged. Start creates no
   database. Retry preserves every queue row, does not construct the engine,
   and returns the same cancellation refusal with an absent or malformed
   database, demonstrating that it does not open the queue store.
2. Review sends the added cancelled record and recovered Spec context to the
   fake provider, omits both deleted Spec source files as
   `archived-spec-source`, and leaves record bytes unchanged. Run causes
   reports disposition `cancelled` in JSON with source history present or
   absent; its existing text output names the disposition when source history
   is unavailable. History sanitize plans only the other Legacy Archive
   Folder and leaves the cancelled record byte-identical after applying that
   one-unit batch. Spec check emits no finding for the returned `open` Backlog
   Entry or cancelled record and resolves a retired Finding's archive license
   through that record; both record and entry bytes remain unchanged.
3. The three required compatibility tests passed unchanged. Additional
   coverage accepts all six other Archive Record dispositions as superseding
   deliverers and prerequisites. The existing no-record retry test passed.

Focused checks (all final commands exited 0):

- `rtk proxy go test -count=1 -run '^Test(SupersedeRefusesACancelledDeliverer|DeliverStartRefusesARequiredCancelledSpec|DeliverRetryRefusesACancelledSpec|SupersedeAcceptsOtherArchiveRecordDispositions)$' ./internal/cli`
- `rtk proxy go test -count=1 -run '^Test(ReviewReadsACancelledArchiveRecord|RunCausesReportACancelledRecord|HistorySanitizeLeavesACancelledRecord)$' ./internal/cli`
- `rtk proxy go test -count=1 -run 'CancelledRecordAndItsReturnedEntry' ./internal/speccheck`
- `rtk proxy go test -count=1 -run '^Test(SupersedeAcceptsAnArchiveRecord|PrerequisitesCountAnArchiveRecord|ReviewDiffOmitsTheRemovedSpecFolder|DeliverRetryStartsAnOwnerWhenNoneIsRunning)$' ./internal/cli`
- `rtk proxy git diff --check`

The initial focused refusal check failed as expected before the guards:
Supersede accepted the cancelled deliverer, start reached authorization
validation, and retry reached the engine or malformed database. Test-fixture
corrections followed the existing reader contracts: Spec-check diagnostics
use `Where`, Run causes emits its archive text line only for unavailable
source history, and only QA override fixtures carry override fields.

The sandbox rejected command startup because the assigned worktree root
contains a symlink; approved elevated commands operated on this same checkout.
This Agent ran no Task Verification command, selected repository Verification,
commit, push or pull request. `UnmetPrerequisites`, the other readers, command help,
Cancel Command, `internal/spec`, skills, guides, Task Graph and other Task files
were not edited. No follow-up work was added to this slice.


### Verification Feedback repair

The Daemon's attempt 1 of `make verify-changed` reported regressions in
`TestDeliverRetryRecordsAPullRequestMergedByHandAfterTheItemBranchIsGone` and
`TestDeliverRetryRecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest`.
Diagnostic provenance:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T171408Z_a9139d7510b0db63/verification/batch-004-attempt-1.log`.
This is Daemon-provided evidence, not an Agent rerun of that command.

The new retry preflight had called the general archive reader even when no
Archive Record existed. That reader also validates Legacy Archive Folders,
which imposed a new refusal on the existing merge-recovery path. The repair
checks for the Archive Record file first and uses the reader only when that
record exists. An absent record leaves legacy-folder handling to the existing
Delivery Engine path. Existing record errors still refuse, and a cancelled
record still refuses before opening the Run Database or constructing the
engine. Only `internal/cli/deliver.go` and this Result changed in the repair
turn; the existing regression tests remain unchanged.

Focused repair evidence:

- `rtk proxy go test -count=1 -run '^TestDeliverRetry(RecordsAPullRequestMergedByHandAfterTheItemBranchIsGone|RecordsASpecArchivedOnTheDefaultBranchWithoutARecordedPullRequest|RefusesACancelledSpec|StartsAnOwnerWhenNoneIsRunning|WithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore|RefusesAPullRequestClosedWithoutMerging)$' ./internal/cli`
  exited 0. Both reported regressions passed, as did normal no-record retry,
  open/closed Pull Request handling and cancelled-record refusal with parked,
  absent and malformed queue databases. The cancellation test still proves
  unchanged queue rows and repository/home bytes.
- `rtk proxy git diff --check` exited 0.

Task status remains `in_progress`. This Agent did not rerun authored Task
Verification or the configured repository Verification; the Daemon owns the
next full Verification sequence and settlement. No commit, push or pull
request was created.
