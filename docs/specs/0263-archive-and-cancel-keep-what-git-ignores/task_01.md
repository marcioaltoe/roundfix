---
task: task_01
spec: 0263-archive-and-cancel-keep-what-git-ignores
status: completed
type: backend
complexity: medium
---

# Task 01: Archive and cancel refuse Ignored Spec Files

## Overview

Both commands guard the Spec folder with a `git status` call that omits files
Git ignores, so they remove such files and claim Git keeps them
(`_techspec.md` → Current behavior). This Task replaces both calls with one
shared helper that also reads ignored files, refuses while any exists, lists
each one and writes nothing. It is verifiable on its own on temporary
repositories: the refusal transcripts match and the tree stays byte-identical.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [archive and cancel remove ignored files in the Spec folder](references/2026-10-12-archive-and-cancel-remove-ignored-files-in-the-spec-folder.md),
   for its refusal half.
2. MUST add `readSpecFolderStatus` and `refuseSpecFolderStatus` following
   `_techspec.md` → Interfaces and Data Models, in the new
   `internal/cli/spec_folder_status.go`, reading the raw NUL-separated output
   of one `git status --porcelain=v1 -z --no-renames --untracked-files=all
   --ignored=traditional -- <folder>` run in the repository that holds the
   Specs Root, never the trimmed text `gitOutput` returns.
3. MUST call the helper from the Archive Command and the Cancel Command in
   place of their current `git status` calls, before any write, following API
   Contract 1: uncommitted changes refuse first with today's reason byte for
   byte; otherwise each Ignored Spec File is listed, repository-relative and
   sorted, one per file even inside an ignored directory, and the command
   exits 2 through Preflight Validation having changed no byte.
4. MUST add the refusal lines of API Contract 4 to both help texts: the archive
   line `Refuses, writing nothing, while the Spec folder has uncommitted
   changes or Git ignores a file under it.` and the cancel refusal sentence
   naming `while Git ignores a file under the Spec folder`. The exit-code lines
   belong to task_02 and MUST NOT change here.
5. MUST NOT change the symbolic-link destination parent checks of
   `internal/spec` (ADR-0269).
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. `TestArchiveRefusesIgnoredSpecFiles` asserts every line of Surface
   Transcript 1 and `TestCancelRefusesIgnoredSpecFiles` every line of Surface
   Transcript 2, each with only the transcript's declared variations, and both
   prove the Spec folder, the three ignored files and
   `git status --porcelain --untracked-files=all --ignored` are unchanged and
   no record exists. `TestCancelTranscriptPrintsHelp` and `TestRunArchiveHelp`
   take the new help lines; every other existing archive and cancel test MUST
   pass unchanged.
7. MUST describe the refusal in `docs/user-guide/commands/archive.md` and
   `docs/user-guide/commands/cancel.md`, each containing the phrase
   `remove each file or commit it with git add -f`, add
   `internal/cli/spec_folder_status.go` to the sources of the `command archive`
   and `command cancel` entries of `docs/references/skill-coverage.json`, and
   re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the shared helper and its raw `git status` parsing.
- [ ] Call it from both commands and keep the uncommitted refusal first.
- [ ] Add the help lines and update the help tests.
- [ ] Add the refusal tests, including the external Specs Root.
- [ ] Update both user guides, the coverage map and the Behavior Surface Record.

## Acceptance Criteria

- [ ] With three ignored files, archive and cancel exit 2, list all three, and leave every byte and the Git status unchanged.
- [ ] With an uncommitted change and an ignored file, both report today's uncommitted reason.
- [ ] An ignore rule held only by an external Specs Root's repository is honored.
- [ ] Both help texts and both user guides describe the refusal.

## Consumers

- Archive Command and Cancel Command refusal output — consumers: `./internal/cli`
- Behavior Surface Record of the `archive` and `cancel` help and guides — consumers: `./internal/docscontract`

## Context

- instruction: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0262-a-cancelled-spec-leaves-an-archive-record-and-returns-its-sources.md`
- creates: `internal/cli/spec_folder_status.go`
- creates: `internal/cli/spec_folder_status_test.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/cancel.go`
- interface: `internal/cli/archive_test.go`
- interface: `internal/cli/cancel_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/cancel.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for f in docs/user-guide/commands/archive.md docs/user-guide/commands/cancel.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'remove each file or commit it with git add -f' || { printf 'missing refusal phrase in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- 'internal/cli/spec_folder_status.go' docs/references/skill-coverage.json || { printf 'coverage map lacks the helper source\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task neither guide names the refusal and the map lacks the helper, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestArchiveRefusesIgnoredSpecFiles|TestCancelRefusesIgnoredSpecFiles|TestSpecFolderStatusRefusesUncommittedChangesFirst|TestSpecFolderStatusReadsTheExternalSpecsRootRepository|TestCancelTranscriptPrintsHelp|TestRunArchiveHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveRefusesIgnoredSpecFiles TestCancelRefusesIgnoredSpecFiles TestSpecFolderStatusRefusesUncommittedChangesFirst TestSpecFolderStatusReadsTheExternalSpecsRootRepository; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Archive|Cancel' ./internal/cli || exit 1` — expected: exit 0. Before this Task the four new tests do not exist, so the PASS check fails; afterwards they, the two help tests and every archive, cancel and Delivery Queue archive test pass.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Core Feature 3; Success Metric 1; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 4; Surface Transcript 1; Surface Transcript 2; Interfaces; Data Models; Invariants 1, 2, 5 and 6; Testing Approach; Build Order 1
- ADR-0269; ADR-0247; ADR-0262

## Result

Implemented the refusal half of the 2026-10-12 Backlog Entry. Archive and
cancel now share `readSpecFolderStatus` and `refuseSpecFolderStatus`: one raw
NUL-separated porcelain-v1 status call, with renames disabled and every
untracked and ignored file enumerated. Non-ignored changes keep the previous
reason and precedence; otherwise ignored paths are sorted, listed and refused
through Preflight Validation before any write. Ignored-folder display resolves
symbolic-link path aliases against Git's physical repository root; the
existing uncommitted-folder display is preserved.

Both help texts and user guides describe the refusal and how to remove or
commit the files. The coverage map names the helper for both commands, and
the sanctioned record command regenerated only the archive and cancel command
and guide fingerprints. Help exit-code lines and `internal/spec` destination
parent checks were not changed; folder removal recovery and skill/glossary
changes remain the later Tasks' slices.

### Focused implementation evidence

Before the guard change,
`go test -count=1 ./internal/cli -run '^(TestArchiveRefusesIgnoredSpecFiles|TestCancelRefusesIgnoredSpecFiles)$'`
exited 1: both commands exited 0, removed the three ignored files and claimed
their bytes were kept in Git. This reproduced the recorded defect in disposable
repositories.

After the final implementation edits:

- `go test -count=1 ./internal/cli -run '^(TestArchiveRefusesIgnoredSpecFiles|TestCancelRefusesIgnoredSpecFiles|TestSpecFolderStatus.*)$'`
  exited 0. The four required new tests and the additional raw-path test ran.
- `go test -count=1 ./internal/cli -run '^(TestRunArchiveHelp|TestCancelTranscriptPrintsHelp|TestArchiveCommandRefusesUncommittedSpecChanges|TestCancelRefusalsLeaveTheTreeUnchanged|TestRunArchiveUsesConfiguredExternalSpecRoot|TestArchiveCommandLeavesTheArchiveRecord|TestCancelTranscriptCancelsASpecWithSources)$'`
  exited 0. These selected existing tests cover the help changes, prior
  uncommitted refusals, external-root success, and ordinary archive and cancel
  writes including Source Return.
- Required regeneration:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 after the final source and guide edits.
- `git diff --check` exited 0 before this Result was appended.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Three ignored files refuse without changes | `TestArchiveRefusesIgnoredSpecFiles` and `TestCancelRefusesIgnoredSpecFiles` assert each full Surface Transcript, empty stdout and exit 2. Each compares all repository file bytes, the Spec folder hash, and `git status --porcelain --untracked-files=all --ignored`, checks all three ignored files remain, and checks no Archive Record or Run Database exists. The fixture ignores both logs and the `out/` directory, proving individual files are enumerated inside an ignored directory. |
| Uncommitted changes refuse first | `TestSpecFolderStatusRefusesUncommittedChangesFirst` checks both commands against the existing reason byte for byte with modified, staged, deleted, untracked and renamed files alongside ignored files. It also checks file bytes and Git status remain unchanged. |
| External Specs Root honors its own ignore rules | `TestSpecFolderStatusReadsTheExternalSpecsRootRepository` checks both commands with the ignore rules committed only in the external repository and no project `.gitignore`. It asserts the same repository-relative refusal, unchanged files and Git status in both repositories, and no record under the external archive root. |
| Help and guides describe the refusal | All three help tests pass (including the additional full archive help transcript checked in the feedback repair below). Both guides contain `remove each file or commit it with git add -f`; both coverage entries name the helper. The sanctioned generator passed and changed only the two commands' and two guides' recorded fingerprints. |

`TestSpecFolderStatusPreservesRawPaths` additionally proves the first modified
entry's leading status space is not trimmed and paths with a leading space,
tab or newline stay unquoted and intact.

The Task's declared Verification commands, the selected repository
Verification and the full archive/cancel suite were not run by this Agent; the Daemon owns
those checks and settlement. Task status and the Task Graph were not edited.
No commit, push or pull request was made. The assigned work is handed back for
Daemon Verification, without a terminal Task verdict.


### Verification feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261013T163453Z_215d9dc7373d8dae/verification/batch-001-attempt-1.log`.
It identified `TestArchiveHelpNamesThePullRequestOption`: its full archive
help expectation omitted this Task's newly required refusal sentence.

Added exactly that sentence to the expected transcript in
`internal/cli/archive_pull_request_test.go`. This additional changed path owns
another full transcript of the same public help surface; the test continues
to compare all stdout, stderr and exit code. Production behavior and help
exit-code lines were not changed in this repair.

Focused check:
`go test -count=1 ./internal/cli -run '^(TestArchiveHelpNamesThePullRequestOption|TestRunArchiveHelp|TestCancelTranscriptPrintsHelp)$'`
exited 0. All three help tests pass with the authored refusal line. The
previously recorded guard and generator evidence concerns unchanged sources.

Declared Verification was not rerun by this Agent. Task status remains
Daemon-owned, and no other Task or the Task Graph was edited. The repair is
handed back for the Daemon's next Verification attempt, without a terminal
Task verdict. No commit, push or pull request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/archive_pull_request_test.go`
