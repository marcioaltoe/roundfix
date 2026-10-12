---
task: task_03
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
type: backend
complexity: medium
---

# Task 03: The Archive Command records the pull request and warns about a squash-dropped revision

## Overview

The Archive Command runs on the delivery branch, before the pull request
exists, and records that branch's commit as `source_revision`. This Task adds
`roundfix archive <slug> --pull-request <n>`. It stores the number when
archiving, and also stores it in a record the current branch already added.
It also adds a standard-error warning when `source_revision` is not an
ancestor of the default branch. It is verifiable on its own through the five
Surface Transcripts in temporary repositories whose `origin` is a local bare
repository.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12
   [an Archive Record points at a branch commit a squash merge drops](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)
   on the archive side, implementing `_techspec.md` → API Contract 1 and API
   Contract 2 exactly, including every refusal of API Contract 1.
2. MUST add `spec.RecordArchivePullRequest` as in `_techspec.md` →
   Interfaces, outside `internal/spec/archive.go`. Every record field other
   than `pull_request` MUST round-trip unchanged, except an outcome paragraph
   that `RenderArchiveRecord` shortens to keep the 2,048-byte target. A record
   that already stores the same number MUST stay byte-identical.
3. MUST decide the record-only form by the absence of the active Spec folder
   and the presence of `<archive-root>/<slug>.md` in the working tree. It
   MUST refuse when `<default-ref>:<record path>` exists. A record already on
   the default branch stays byte-identical.
4. MUST write the warning only to standard error and keep exit 0 and the
   existing standard output line. With `--pull-request`, standard output
   gains the suffix `; pull request <n>`.
5. MUST assert each of `_techspec.md` → Surface Transcripts 1-5 whole, in one
   named test each: every line of standard output and standard error in
   order, and the exit code, with only the transcript's declared variations.
   The tests are `TestArchiveOnABranchWarnsAboutASquashDroppedRevision` (1),
   `TestArchiveWithAPullRequestNamesItsHeadRef` (2),
   `TestArchiveRecordsAPullRequestInTheBranchRecord` (3),
   `TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch` (4) and
   `TestArchiveHelpNamesThePullRequestOption` (5). Also add
   `TestArchiveOnTheDefaultBranchDoesNotWarn`, `TestArchivePullRequestRefusals`
   (non-integer, different stored number, unknown default branch, `--plan`,
   record-only with `--promote` or `--qa-override`) and, in `internal/spec`,
   `TestRecordArchivePullRequestChangesOnlyThatField`. Each top-level test
   calls `t.Parallel()` first or carries a `// Sequential: <reason>` comment.
   No test opens a network connection.
6. MUST keep `TestArchiveCommandLeavesTheArchiveRecord` passing unchanged.
7. MUST re-record the Behavior Surface Record, because the archive help
   changes, with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   never by hand.
8. MUST prove each new gate can fail. The Result MUST record two sabotages,
   restored afterwards: the warning suppressed, and the default-branch refusal
   removed. For each it names the test that failed.
9. MUST NOT change `internal/spec/archive.go`, the Archive Record schema, the
   archive cut, the QA override or any skill text.

## Subtasks

- [ ] Add the record field writer beside the record code.
- [ ] Add the flag, its record-only form and its refusals.
- [ ] Add the ancestry warning.
- [ ] Assert the five transcripts and the refusals; re-record the Behavior Surface Record.
- [ ] Record the sabotages.

## Acceptance Criteria

- [ ] Archiving on a branch prints the warning and exits 0; on the default branch it does not warn.
- [ ] `--pull-request` stores the number when archiving and in the branch's own record.
- [ ] A record already on the default branch is refused byte-identical.
- [ ] The help names the new form and option.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `internal/cli/archive_record_test.go`
- instruction: `internal/preflight/default_branch.go`
- instruction: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/spec/archive_record_pull_request.go`
- creates: `internal/spec/archive_record_pull_request_test.go`
- creates: `internal/cli/archive_pull_request_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveOnABranchWarnsAboutASquashDroppedRevision|TestArchiveWithAPullRequestNamesItsHeadRef|TestArchiveRecordsAPullRequestInTheBranchRecord|TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch|TestArchiveHelpNamesThePullRequestOption|TestArchiveOnTheDefaultBranchDoesNotWarn|TestArchivePullRequestRefusals|TestArchiveCommandLeavesTheArchiveRecord)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveOnABranchWarnsAboutASquashDroppedRevision TestArchiveWithAPullRequestNamesItsHeadRef TestArchiveRecordsAPullRequestInTheBranchRecord TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch TestArchiveHelpNamesThePullRequestOption TestArchiveOnTheDefaultBranchDoesNotWarn TestArchivePullRequestRefusals TestArchiveCommandLeavesTheArchiveRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^TestRecordArchivePullRequestChangesOnlyThatField$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestRecordArchivePullRequestChangesOnlyThatField" || exit 1; help="$(go run -buildvcs=false ./cmd/roundfix archive --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -qF -- "--pull-request Store the delivery pull request number in the Archive Record" || { printf 'archive help lacks the option\n' >&2; exit 1; }; out="$(go test -count=1 -tags docscontract -v -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheSkillCoverageMapIsCurrent" || exit 1` — expected: exit 0; before this Task the record writer test does not exist and the help lacks the option, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 3; User Story 2; Core Feature 1; Core Feature 4; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Data Models; API Contract 1; API Contract 2; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Testing Approach; Build Order 3
- ADR-0261; ADR-0247; ADR-0184; ADR-0256

## Result

Implemented the archive-side contracts of ADR-0261 and the adopted squash-merge
Backlog Entry. `--pull-request` records a validated decimal number during a
normal archive or in an existing branch record, with byte-identical replay of
the same number. Record-only detection uses the missing active folder and the
working-tree record; it checks the record path at the resolved default tip
before writing. Full archive prints the existing confirmation with the optional
pull request suffix, and checks ancestry for the standard-error warning without
changing its successful exit code.

`RecordArchivePullRequest` lives beside the record code, outside `archive.go`.
It uses the existing parser and renderer, refuses conflicting stored numbers,
and preserves every other field, including override, supersession, relink,
regeneration and delivery metadata. Only the renderer's outcome shortening is
permitted. No Archive Record schema, archive cut, QA override implementation,
skill text, Task Graph or other Task file changed.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Branch warning, exit 0; no warning on the default branch | `TestArchiveOnABranchWarnsAboutASquashDroppedRevision` asserts Surface Transcript 1 whole. `TestArchiveWithAPullRequestNamesItsHeadRef` asserts Transcript 2 whole. `TestArchiveOnTheDefaultBranchDoesNotWarn` covers both origin/HEAD and the local main fallback. `TestArchiveWithAnUnknownDefaultBranchWarns` asserts the exact unknown-default warning. |
| Store the number during archive and in the branch record | Transcript 2 checks the stored field. `TestArchiveRecordsAPullRequestInTheBranchRecord` asserts Transcript 3 whole, checks that other files are unchanged, and asserts the idempotent transcript and byte identity. `TestRecordArchivePullRequestChangesOnlyThatField` covers every disposition, all supported metadata, short and oversized UTF-8 outcomes, noncanonical replay, invalid values, conflicting values, and malformed or missing records. |
| Refuse a default-branch record byte-identical | `TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch` asserts Transcript 4 whole and checks unchanged bytes. `TestArchivePullRequestRefusesAnExternalDefaultBranchRecord` also covers an external Spec root nested below its Git repository root. |
| Help names the new form and option | `TestArchiveHelpNamesThePullRequestOption` asserts Transcript 5 as the complete existing help plus the new usage form and option, empty stderr, and exit 0. The Behavior Surface Record was regenerated; only the command archive fingerprint changed. |

`TestArchivePullRequestRefusals` asserts full preflight output, exit 2 and no
file changes for non-integer, empty, zero, negative, leading-zero and overlong
numbers; a different stored number; unknown default branch with either form;
`--plan`; and record-only combinations with `--promote`, `--qa-override`,
`--approval` or `--reason`. Every new top-level test calls `t.Parallel()` first.
Tests use disposable repositories and local bare origins; none opens a network
connection.

The existing `TestArchiveCommandLeavesTheArchiveRecord` body is byte-identical
to HEAD and passes in the focused run. Its shared `commitArchiveFixture` now
places the disposable archive fixture on main, making its existing empty-stderr
expectation exercise the default-branch case under the new warning contract.
This also preserves the existing Archive Relink transcripts. Those four helper
lines are the only change in `internal/cli/archive_record_test.go`.

### Focused checks

- Initial red check: `go test ./internal/cli -run '^TestArchiveHelpNamesThePullRequestOption$' -count=1`
  exited 1 because the new tests referenced the then-missing
  `spec.RecordArchivePullRequest` entry point.
- Final check after implementation, fixture adjustments and sabotage restoration:
  `GOPROXY=off GOSUMDB=off go test ./internal/cli ./internal/spec -run '^(TestArchive|TestRecordArchivePullRequest)' -count=1`
  exited 0; both packages reported `ok`. This includes all new tests and the
  unchanged legacy archive test, plus the existing archive regressions.
- Required generation:
  `GOPROXY=off GOSUMDB=off go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and generated `docs/references/behavior-surfaces.json`; the record
  was not edited by hand.
- `git diff --check` exited 0. An inspection compared the legacy test body with
  HEAD and confirmed identity, and confirmed both sabotaged gates were restored.

### Required sabotages

Each sabotage changed only `internal/cli/archive.go` temporarily and restored
its exact pre-sabotage bytes in a `finally` block.

1. Suppressed the call to `warnArchiveSourceRevision`. Running
   `go test ./internal/cli -run '^TestArchiveOnABranchWarnsAboutASquashDroppedRevision$' -count=1`
   exited 1 and reported that named test as failing: exit and stdout still
   matched, but stderr was empty instead of the warning. Restored the call.
2. Disabled the `entry != ""` default-branch refusal. Running
   `go test ./internal/cli -run '^TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch$' -count=1`
   exited 1 and reported that named test as failing: the command wrote the
   number, printed the recording confirmation and exited 0 instead of the
   full preflight refusal and exit 2. Restored the refusal.

### Execution boundary

The baseline contained only the Daemon's existing `status: in_progress` change
in this Task file. Task status remains untouched. The Task's declared
Verification commands, selected repository Verification and incremental
Verification were not run; the Daemon owns them and settlement. No commit,
push or pull request was created in the work repository. Local pushes occurred
only between disposable test repositories and their local bare origins.

The executor initially refused all sandboxed commands because the configured
writable root contains a symlink. Approved sandbox bypass allowed the scoped
work and checks. Source Revision Recovery and Delivery Queue recording remain
with task_02 and task_04 respectively; no additional follow-up was found.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/archive_record_test.go`

## Carry-forward provenance

- Source Run: `run_20261012T140955Z_23b372ef5c35c71c`
- Source commit: `219a9b9ed461a37ef11e239f0d1f6583d826dba5`
