---
task: task_02
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
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


## Result

Implemented this Task's reader slice for Daemon Verification. The incoming
`in_progress` status was left untouched. The initial worktree had only the
Daemon's status change in this Task file.

### Implementation

- Added the named-schema predicate to the review scope, review Spec context
  and archived-slug enumeration. YAML schema identity is separate from full
  record validation, so schema-declaring malformed records still fail.
- Added Source Revision Recovery with the record number or adding-commit
  subject, strict positive-number validation, local availability first and
  an availability check after fetching. Fetch uses the configured remote,
  no tags, no destination ref, no `FETCH_HEAD`, an empty refmap and no automatic
  maintenance. It has no stdin, forces `GIT_TERMINAL_PROMPT=0`, and runs through
  `exec.CommandContext` under the 60-second `SourceRecoveryFetchTimeout`.
- Wired recovery into all four folder readers. Review and Run causes use
  `watch.push_remote`; the other two use `origin`. Existing unavailable
  diagnostics, skipped context, report metadata and unresolved citations
  retain their previous behavior.
- Added the required tests and a reusable local-forge/main-only-clone fixture
  in `internal/gittest/pull_request.go`. Every recovery fixture proves
  `git cat-file -e <source_revision>^{commit}` fails before its reader runs.
  Every new top-level test starts with `t.Parallel()`; the process-boundary
  probe re-executes the test binary with a private environment and uses
  `testfixture.FixtureBinary` for its compiled Git wrapper (ADR-0125).
- Updated only the enumeration assertion of
  `TestReadArchivedSpecReadsRecordAndLegacyFolder`: its unnumbered `demo.md`
  is intentionally excluded by the new predicate. All its known-slug read
  assertions remain. No schema, archived record, archive cut, dependency,
  Task Graph or other Task file was changed.

### Characterization and sabotages

Before the production change:

```sh
go test ./internal/cli -run '^TestReviewIgnoresAProseReadmeInTheArchiveRoot$' -count=1
```

Exit 1: `TestReviewIgnoresAProseReadmeInTheArchiveRoot` reported provider calls
0, review exit 2 and
`read review archive record: missing YAML frontmatter opening marker`.

Two temporary sabotages were applied individually and restored in `finally`
blocks before subsequent checks:

1. Bypassed `IsArchiveRecord` in the review scope. The characterization
   command above exited 1 with the same README error and zero provider calls.
2. Removed the fetch invocation from `EnsureSourceRevision`.
   `go test ./internal/spec -run '^TestSourceRevisionRecoveryFetchesThePullRequestHead$' -count=1`
   exited 1 in `TestSourceRevisionRecoveryFetchesThePullRequestHead` with
   `check fetched source revision: source_revision is unavailable`.

### Focused-check evidence

The final restored implementation passed this focused regression selection
(exit 0 in all four packages):

```sh
go test ./internal/spec ./internal/cli ./internal/worktree ./internal/speccheck -run '(ArchiveRecord|ArchivedSpec|ArchivedTask|SourceRevisionRecovery|PullRequestFromSubject|ReviewIgnoresAProse|ReviewStillRefuses|Recovers?ASquashDroppedRevision|LeftoversReadTaskContext|LeftoversKeepTheWorktree|RunCausesReadAnArchivedTaskGraph)' -count=1
```

| Acceptance criterion | Focused evidence |
| --- | --- |
| A prose archive-root README reaches the provider | `TestReviewIgnoresAProseReadmeInTheArchiveRoot` asserts review exit 0, one provider call, prose in the prompt and no archived/skipped Spec context. |
| A malformed declared record blocks review | `TestReviewStillRefusesAMalformedArchiveRecord` asserts refusal before any provider call. The identity and enumeration tests also reject malformed declared records without treating prose or another schema as records. |
| All four readers recover the folder | `TestReviewContextRecoversASquashDroppedRevision`, `TestRunCausesRecoverASquashDroppedRevision`, `TestLeftoversRecoverASquashDroppedRevision` and `TestArchivedTaskContextRecoversASquashDroppedRevision` each cover the record number and adding subject through a bare forge. Review and causes use a remote named `delivery` with an unreachable `origin`; Spec check covers built-in and custom roots. Leftovers also reads the archived QA report. |
| No ref or FETCH_HEAD; unavailable behavior preserved | `TestSourceRevisionRecoveryWritesNoRef` compares all refs before/after, includes a configured pull-ref mapping and requires absent `FETCH_HEAD`. The unreachable-revision test covers an absent forge and a head lacking the revision. Each reader test has an unreachable-forge case preserving its prior outcome. The unchanged `TestRunCausesReportAnArchiveRecordWithoutItsRevision` and `TestLeftoversKeepTheWorktreeWithoutTheSourceRevision` pass. |

Additional evidence: malformed numbers fetch no objects; an available local
commit succeeds with an unreachable remote; a later record edit carrying
`(#8)` still recovers from the adding commit's `(#7)`.
`TestSourceRevisionRecoveryCancelsTheFetchProcess` observes prompt setting 0,
EOF on stdin, cancellation by a tighter parent deadline, and a dead/reaped
fetch PID on return.

The broader selection first exposed the obsolete `demo.md` enumeration
expectation above; it passed after updating that intentionally changed
contract. Initial CLI recovery fixtures also exposed macOS `/var` versus
`/private/var` repository keys; canonicalizing the clone paths corrected the
fixture artifacts and disposable Run-store keys.

`git diff --check` passed. The Agent did not run the Task's declared
Verification commands, repository Verification or incremental Verification;
the Daemon owns those commands. No commit, push or pull request was made.
No follow-up outside this Task was implemented.


### Verification Feedback repair — attempt 1

The Daemon's `make verify-changed` attempt reported an executable-fixture
inventory failure in the cancellation probe. Diagnostic source:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T140955Z_23b372ef5c35c71c/verification/batch-002-attempt-1.log`.
The probe wrote an executable shell wrapper, contrary to ADR-0125's compiled
fixture rule.

Replaced that wrapper with a Go binary built through the existing
`testfixture.FixtureBinary` helper. The compiled fixture delegates local
reads to real Git and blocks the fetch inside its own process. The probe
retains its private child environment, prompt, EOF and dead-PID assertions.
This repair changes only this Task's test and Result; the executable
inventory and production implementation remain unchanged.

Focused checks after the repair:

- `go test ./internal/spec -run '^TestSourceRevisionRecoveryCancelsTheFetchProcess$' -count=1`
  — exit 0; the compiled-fixture probe verifies the non-interactive fetch
  boundary and cancellation that kills and reaps its process.
- `go test ./internal/testfixture -run '^TestNoTestWritesAnExecutableOutsideTheResidue$' -count=1`
  — exit 0; the repository inventory accepts the repaired fixture.
- `git diff --check` — exit 0.

Status remains `in_progress`, owned by the Daemon. The Agent did not rerun
`make verify-changed` or any declared Verification command. The configured
Verification rerun and settlement remain with the Daemon.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/gittest/pull_request.go`
- `internal/spec/archive_record_test.go`

## Carry-forward provenance

- Source Run: `run_20261012T140955Z_23b372ef5c35c71c`
- Source commit: `003f2a8769e196cae168af15aa9920bc0db0f632`
