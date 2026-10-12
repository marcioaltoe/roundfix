---
task: task_03
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: pending
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
