---
task: task_02
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
status: completed
type: backend
complexity: medium
---

# Task 02: The archive subject and the pull request title keep the slug while it fits

## Overview

The Delivery Queue's archive commit subject `docs: archive <slug>` and pull
request title `feat: deliver <slug>` carry the whole slug, so a long slug
passes 100 characters: a commitlint `commit-msg` hook refuses the archive
commit, and the squash commit lands with an over-long subject. This Task keeps
both messages byte for byte while they fit and otherwise names the Spec by its
Spec Number, the title followed by the Spec title cut at a word, and proves
that every reader of the two messages still works.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md`](references/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-6 of
   `_techspec.md` in `internal/cli/deliver_workflow.go`.
2. MUST extract `specNumber(specSlug string) (string, bool)` from
   `pullRequestRecordSubject` (API Contract 1) and make
   `pullRequestRecordSubject` use it, returning exactly today's bytes for
   every slug.
3. MUST add `queueHeaderLength`, counting UTF-16 code units with
   `unicode/utf16`, the limit 100 and the squash reserve ` (#99999)` (API
   Contract 2), and `archiveCommitSubject(specSlug string) string` (API
   Contract 3), and make `Archive` commit with it through the existing
   `commitQueueChange`, so a hook refusal still quotes the subject it refused.
4. MUST add `pullRequestTitle` (API Contract 4), whose slug form is chosen
   without calling the Spec title reader, whose number form is `feat: deliver
   <spec number>: <text>` and whose form without a Spec Number is `feat:
   deliver <text>`, with the whitespace collapse, the word cut, the trailing
   `.`, `,`, `;`, `:` and space removal, and the empty-text forms exactly as
   API Contract 4 states; and make `Publication` use it (API Contract 5),
   reading the Spec title from `spec.ReadArchivedSpec` at the archive root
   `archivePaths` resolves, only when the slug form does not fit, and failing
   planning with `plan publication: read Spec title: <error>` when that read
   fails. The body stays `Delivers Spec <slug>.`; the Linux Lane paragraph
   the Delivery Engine appends to it is untouched.
5. MUST state the hazard in the code's behavior: the Spec title is
   repository-authored text from the PRD's first heading, already stored in
   the Archive Record; it goes, cut and on one line, into the pull request
   title on the forge and, through a squash merge, into the default branch's
   history. The title never carries text other than the Spec Number, the
   fixed words and a prefix of that title.
6. MUST keep API Contract 6: `reconcileArchiveCommit`,
   `archiveCommitIsExact`, `commitQueueChange`, the derived merge message and
   every file under `internal/delivery`, `internal/spec` and
   `internal/releaseplan` keep their bytes. No existing test, golden or
   fixture changes.
7. MUST add, in the new file `internal/cli/deliver_queue_subject_test.go`,
   the tests `_techspec.md` → Testing Approach lists, each top-level test
   calling `t.Parallel()` first, compiling for Windows, using temporary
   repositories, homes and Git configuration, writing hooks with
   `writeScriptFixture`, and waiting on no real time; the tests that make real
   Git commits re-run themselves under an isolated home with
   `queueCommitTestHome`:
   - `TestQueueArchiveSubjectKeepsTheSlugWhileItFits` and
     `TestQueuePullRequestTitleKeepsTheSlugWhileItFits`, with the cases listed
     there, including the 86/87-unit archive boundary, the 77/78-unit title
     boundary with a reader that fails the test if it is called for the slug
     form, Spec 0278's slug and title giving `feat: deliver 0278: A fixed
     round-2 finding gets one Correction Review, and the plan names` (90
     units), and a title holding an astral rune counted as two units; each
     archive subject is at most 100 units and each title plus ` (#99999)` at
     most 100 units.
   - `TestQueueMessagesStayReadableBySquashReaders`: each title plus ` (#575)`
     gives `575` to `spec.PullRequestFromSubject` and is classified by
     `releaseplan.ClassifyCommit` with type `feat`, a minor impact and no
     breaking change; each archive subject form is classified `docs`.
   - `TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook`: built like the
     archive subtest of `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery`,
     with a 100-character synthetic slug starting `0296-` and a `commit-msg`
     hook that exits 1 with `header must not be longer than 100 characters` on
     standard error when the message's first line exceeds 100 characters;
     `Archive` returns an exact head whose subject is `docs: archive the 0296
     Spec`, and the worktree is clean.
   - `TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject`: the same
     slug with a refusing `pre-commit` hook; after the refusal the hook is made
     to pass, the staged archive is committed by hand under `docs: archive
     <slug>`, the old long subject, and a second `Archive` from the reviewed
     head returns that head as an exact Spec move with no further commit.
   - `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`: an archived
     item fixture built like `queueRecordHookFixture` with a 100-character
     slug and a known Archive Record title; `Publication` returns the
     number-form title and the body `Delivers Spec <slug>.`; with the record
     removed from the worktree it fails with an error that starts `plan
     publication: read Spec title:`.

## Subtasks

- [ ] Extract the Spec Number helper and keep the record subject's bytes.
- [ ] Build the archive subject and the title from the slug while it fits, else from the Spec Number.
- [ ] Read the Spec title from the Archive Record only for the number form.
- [ ] Add the subject, title, reader, hook, adoption and publication tests.

## Acceptance Criteria

- [ ] Behind a 100-character header hook, the archive commit of a 100-character slug succeeds with `docs: archive the 0296 Spec`.
- [ ] Every archive subject is at most 100 units and every title plus ` (#99999)` at most 100 units, and both keep the slug form byte for byte while it fits.
- [ ] An archive commit made by hand under the old long subject is adopted.
- [ ] Squash subjects built on every title form still yield their number and the type `feat`.
- [ ] The record subject test and every existing hook and publication test pass unchanged.

## Consumers

- The archive subject (`Archive`) and the pull request title (`Publication`), and the extracted Spec Number helper behind `pullRequestRecordSubject` — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- instruction: `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/delivery/engine.go`
- instruction: `internal/delivery/github.go`
- instruction: `internal/spec/source_revision.go`
- instruction: `internal/spec/history_sanitize.go`
- instruction: `internal/spec/archive_record.go`
- instruction: `internal/releaseplan/classifier.go`
- instruction: `internal/cli/deliver_commit_hook_test.go`
- instruction: `internal/cli/deliver_publication_scope_test.go`
- instruction: `internal/cli/deliver_record_pull_request_test.go`
- instruction: `internal/cli/deliver_archive_record_test.go`
- instruction: `internal/cli/archive_record_test.go`
- instruction: `internal/cli/archive_test.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/cli/script_fixture_test.go`
- creates: `internal/cli/deliver_queue_subject_test.go`
- interface: `internal/cli/deliver_workflow.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestQueueArchiveSubjectKeepsTheSlugWhileItFits|TestQueuePullRequestTitleKeepsTheSlugWhileItFits|TestQueueMessagesStayReadableBySquashReaders|TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook|TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject|TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug|TestPullRequestRecordSubjectNamesTheSpecNumber|TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery|TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook|TestDeliveryPublicationRequestsNoCodeRabbitReview)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQueueArchiveSubjectKeepsTheSlugWhileItFits TestQueuePullRequestTitleKeepsTheSlugWhileItFits TestQueueMessagesStayReadableBySquashReaders TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug TestPullRequestRecordSubjectNamesTheSpecNumber TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery TestDeliveryPublicationRequestsNoCodeRabbitReview; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the six new tests do not exist, so the command fails at its first pass check. Afterwards the archive subject, title, reader, hook, adoption and publication tests pass with the unchanged record subject, hook refusal and publication tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 2
- ADR-0302
- ADR-0309


## Result

Implemented API Contracts 1–6 in `internal/cli/deliver_workflow.go`.
`specNumber` now supplies the record subject without changing its output.
Archive subjects retain fitting slugs and otherwise name the Spec Number.
Publication titles reserve ` (#99999)`, count UTF-16 units, and lazily read
the Archive Record title only for an over-budget slug. Whitespace collapse,
word/rune cuts, trailing punctuation removal, empty-title forms and wrapped
reader errors follow the TechSpec. A code comment names the repository-authored
PRD heading's path through the Archive Record to the forge and squash history.
The publication body keeps the complete slug.

Added `internal/cli/deliver_queue_subject_test.go`; every top-level test starts
with `t.Parallel()`. Git scenarios use `queueCommitTestHome`, disposable
repositories and configuration, and `writeScriptFixture` hooks. No real-time
waits or existing test/fixture changes were introduced.

### Acceptance evidence

- Header hook: `TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook`
  passed for a 100-character slug, returning an exact archive head with
  `docs: archive the 0296 Spec` and a clean worktree.
- Bounds and byte compatibility: the archive/title table tests passed at
  86/87-unit and 77/78-unit boundaries, with long slugs, ten/eleven-digit
  prefixes, absent/non-ASCII numbers, whitespace, punctuation, empty text,
  word cuts and astral runes. Fitting titles never call the reader. Every
  subject/title is checked against its UTF-16 budget; Spec 0278 produces the
  specified 90-unit title.
- Adoption: `TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject`
  passed after a real hook refusal and hand commit under the old long
  subject; retry returned that head, exactly one commit after the reviewed
  head, with a clean worktree.
- Readers: `TestQueueMessagesStayReadableBySquashReaders` passed for every
  successful title table case with ` (#575)`, yielding number `575`, type
  `feat`, minor impact and no breaking change; archive forms classify `docs`.
- Existing contracts: unchanged record-subject, hook-limit, hook-refusal,
  record-adoption, record-writing and publication tests passed in the
  focused regression command below. Long-slug publication also passed with
  both the built-in and a configured archive root; deleting the record
  produced `plan publication: read Spec title:`.

### Focused checks

- Before implementation:
  `go test ./internal/cli -run '^TestQueueArchiveSubjectKeepsTheSlugWhileItFits$' -count=1`
  failed compilation because the three new message helpers were absent.
- `go test ./internal/cli -run '^TestQueue(ArchiveSubject|PullRequestTitle|Messages)' -count=1`
  — exit 0. An initial run exposed off-by-one expectations in four new
  cases: the numbered prefix occupies 20 units, leaving 71; corrected those
  expectations before this passing run.
- `go test ./internal/cli -run '^TestDelivery(WorkflowArchiveCommitPassesAHeaderLimitHook|WorkflowAdoptsAnArchiveCommitWhateverItsSubject|PublicationTitleNamesTheSpecNumberForALongSlug)$' -count=1 -v`
  — exit 0; all three scenarios passed.
- `go test ./internal/cli -run '^Test(PullRequestRecordSubject|DeliveryQueueCommit|DeliveryWorkflowRecordCommit|DeliveryPublicationRequests|DeliveryWorkflowRecordsThePullRequest|DeliveryWorkflowAdoptsARecord)' -count=1 -v`
  — exit 0; eight existing top-level regression tests passed.
- After adding configured-root coverage:
  `go test ./internal/cli -run '^TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug$' -count=1 -v`
  — exit 0.
- `GOOS=windows go test -c ./internal/cli -o /tmp/roundfix-task0296-02-cli.test.exe`
  — exit 0 after the final test edit; compilation only, no Windows execution.
- `git diff --check` — exit 0. A byte comparison against HEAD confirmed
  `reconcileArchiveCommit`, `archiveCommitIsExact`, `commitQueueChange` and
  `ResolveConflict` unchanged; `git diff --exit-code HEAD -- internal/delivery internal/spec internal/releaseplan`
  also exited 0.

Environment: the configured writable root contains a symlink, so commands
required escalation to operate in the assigned resolved Run checkout.
Only the workflow, the new test file and this Result were changed by this
turn; the Task's pre-existing `in_progress` status was preserved. No follow-up
outside this slice was identified. Declared Task Verification and repository
Verification were not run and remain Daemon-owned. No commit, push or pull
request was made.

## Carry-forward provenance

- Source Run: `run_20261022T014206Z_8044b8cc31aec528`
- Source commit: `9fc4ffb020169ddb92cd06643de0fc38518c773c`
