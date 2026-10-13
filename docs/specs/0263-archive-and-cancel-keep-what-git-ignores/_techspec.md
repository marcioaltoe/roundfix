---
spec: 0263-archive-and-cancel-keep-what-git-ignores
prd: _prd.md
created: 2026-10-12
---

# Archive and cancel keep what Git ignores — Technical Spec

## Executive Summary

One helper in `internal/cli` replaces the two `git status` calls that guard
the Archive Command and the Cancel Command. It adds `--ignored`, so both
commands refuse an Ignored Spec File, list every one and write nothing. In
`internal/spec`, `Archive` and `Cancel` return a typed folder removal error
together with their filled result when the final `os.RemoveAll` fails. The
commands turn it into Folder Removal Recovery: everything written stays, the
exit code is 1, and standard error names the record and the exact finishing
command. Help, user guides, the Roundfix skill's archive reference, the
archive-spec skill and the glossary follow.

The trade-off accepted is ADR-0269's refusal without an escape hatch: a stray
`.DS_Store` or an ignored log now blocks archive and cancel until the operator
removes or commits it, including inside a Delivery Queue item worktree. In
exchange, no command ever removes a byte Git does not hold, and the recovery
never has to undo a partial `os.RemoveAll`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or Archive
  Record field changes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; tests run on temporary repositories and homes. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0269 governs every decision below.
  ADR-0247: "so the recorded revision holds exactly what was removed".
  ADR-0248: "Existing history is sanitized in batches after a history-full
  tag"; History Sanitize is unchanged.
  ADR-0251: "A Legacy Archive Folder is read leniently" and ADR-0254:
  "Baseline update sanitizes pending history in the change it plans"; neither
  reader nor sanitizer changes.
  ADR-0262: "refuses while a Run or a Delivery Queue item is active for the
  Spec". ADR-0261: "An operator records the number with" the archive's
  `--pull-request`. ADR-0260: "An archive relinks live citations". ADR-0184:
  "A TechSpec states a command surface as a transcript". ADR-0256: "A release
  waits for the skills that describe a changed surface". ADR-0244: "A Spec
  declares the domain terms it introduces". ADR-0257: "Inside a Run the Daemon
  is the only full gate". ADR-0258: "A Spec is authored against the QA rerun
  classes". ADR-0259: "A test runs in parallel unless it names why it cannot". ADR-0271: "A test bounds by real time only what it asserts"; the new tests give no child, goroutine or Git a real-time allowance they do not assert.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AF grant of 2026-10-12 ("Sim, como no AE") and the 2026-09-30 skills grant
  ("considere autorizado a ajustar todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0263-archive-and-cancel-keep-what-git-ignores/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`, `internal/spec/archive.go`,
  `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-12 with `bin/roundfix` built at `94fbf85b`, in a
temporary repository and home whose `.gitignore` holds `*.log`, with a Spec
`0001-demo` of one Task:

- `git status --porcelain --untracked-files=all -- docs/specs/0001-demo`
  printed nothing for an ignored `run.log`; with `--ignored` it printed
  `!! docs/specs/0001-demo/run.log`, and with an ignored directory holding two
  logs it printed one `!!` line per file.
- `roundfix cancel 0001-demo --reason characterize` exited 0 and printed
  `removed 4 file(s) (376 bytes) kept in Git at 81add6598229`: three tracked
  files and `run.log`. `git log --all -- docs/specs/0001-demo/run.log` lists no
  commit.
- With the Task completed and a passing QA Report, `roundfix archive 0001-demo`
  exited 0 and printed `removed 6 file(s) (489 bytes) kept in Git at
  35a1bcc2b46f`: five tracked files and `run.log`.
- With `references/` made read-only (`chmod 0555`), `roundfix cancel` exited 1
  with `roundfix: cancel write failed: remove cancelled Spec; leftover folder
  "<folder>": unlinkat <folder>/references/note.md: permission denied`, after
  writing the record and removing the three other files.
- In the same state `roundfix archive` printed `Preflight failed`, the
  leftover-folder reason and `No side effects`, and exited 2, after writing
  `docs/history/specs/0001-demo.md` and removing four of five files.
- A second `roundfix archive` or `roundfix cancel` is refused because the
  record exists, so the printed message is the operator's only guide.

## System Architecture

Both guards live in the command layer today, `internal/cli/archive.go` and
`internal/cli/cancel.go`, so the shared helper goes there too, in one new file
`internal/cli/spec_folder_status.go`; the finishing command formatter and the
recovery report go in `internal/cli/folder_removal_recovery.go`.
`internal/spec` keeps the writes; it gains one error type in a new file
`internal/spec/folder_removal.go`, and `Archive` and `Cancel` return it. The
Delivery Queue runs `roundfix archive` as a child process and parks on any
non-zero exit, so it needs no change. No new package, flag, output format or
record field.

## Implementation Design

### Interfaces

```go
// internal/cli/spec_folder_status.go
type specFolderStatus struct {
	Changes []string // tracked changes and untracked, non-ignored paths
	Ignored []string // Ignored Spec Files, repository-relative, sorted
}

func readSpecFolderStatus(ctx context.Context, repo, folder string) (specFolderStatus, error)
func refuseSpecFolderStatus(command, slug, folderRel string, status specFolderStatus) error

// internal/cli/folder_removal_recovery.go
func folderRemovalFinishCommand(folder string) string

// internal/spec/folder_removal.go
type FolderRemovalError struct {
	Folder string // absolute Spec folder
	Err    error
}

func (e *FolderRemovalError) Error() string
func (e *FolderRemovalError) Unwrap() error
```

### Data Models

`readSpecFolderStatus` runs `git -C <repo> status --porcelain=v1 -z
--no-renames --untracked-files=all --ignored=traditional -- <folder>` and
parses the raw output. It must not read the trimmed text `gitOutput` returns,
because trimming removes the leading space of a first entry such as
` M path`. An entry whose status is `!!` is an Ignored Spec File; every other
entry is a change. `repo` is the repository that holds the Specs Root: the Git
root, or the external Specs Root's own repository.

### API Contracts

1. API Contract: the Archive Command and the Cancel Command call
   `readSpecFolderStatus` where they call `git status` today, before any
   write. Changes refuse first with today's reason, byte for byte. Otherwise
   one or more Ignored Spec Files refuse through Preflight Validation, exit 2,
   with the reason `Spec "<slug>" cannot <command> while Git ignores <n>
   file(s) under <folder> that no commit holds: <path>, <path>; remove each
   file or commit it with git add -f, then retry`, where `<command>` is
   `archive` or `cancel`, `<folder>` and each `<path>` are relative to that
   repository's root, and the paths are sorted and joined with `, `. The
   refusal writes nothing: the Spec folder, the ignored files, Git status and
   every other file stay byte-identical.
2. API Contract: when the final removal of the Spec folder fails, `Archive`
   and `Cancel` return their filled result with a `*FolderRemovalError` that
   wraps the `os.RemoveAll` error. Nothing they wrote is rolled back.
3. API Contract: on that error the Archive Command first records
   `--pull-request` when given, then both commands print nothing on standard
   output and exactly two lines on standard error, and exit 1:
   `roundfix: <command> wrote <record> but could not remove the Spec folder
   <folder>: <cause>` and `roundfix: finish the removal with: <finish>`.
   `<record>` is relative to the Git root as in the success line, `<folder>`
   is relative to the repository root of the Specs Root, and `<cause>` is the
   wrapped error. `<finish>` is `folderRemovalFinishCommand(<absolute
   folder>)`: `chmod -R u+w '<folder>' && rm -rf '<folder>'`, with the folder
   single-quoted and each embedded `'` written as `'\''`. Running it exactly
   as printed leaves the tree a successful run would leave.
4. API Contract: help text. The archive help gains, after `archive creates no
   Run and never pushes.`, the line `Refuses, writing nothing, while the Spec
   folder has uncommitted changes or Git ignores a file under it.`, and its
   exit codes gain `  1  archive could not finish after writing; stderr says
   how to finish` between `0` and `2`. The cancel help's refusal sentence
   becomes `Refuses while an Active Run or an unfinished Delivery Queue item
   holds the Spec, while an active Spec requires it, while Git ignores a file
   under the Spec folder, and for a Spec that carries a supersession.`, and its
   exit code `1` reads `cancel write or removal failed`. Write failures before
   the removal keep their current messages and exit codes.

### Surface Transcripts

Each fixture holds a Spec `<slug>` under `docs/specs/`, a committed
`.gitignore` with `*.log`, and three ignored files: `run.log` at the folder
root and `qa/evidence/out/a.log` and `qa/evidence/out/b.log`. The recovery
fixtures instead hold a committed `references/note.md` under a read-only
`references/`.

1. Surface Transcript: archive refuses Ignored Spec Files.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "<slug>" cannot archive while Git ignores 3 file(s) under docs/specs/<slug> that no commit holds: docs/specs/<slug>/qa/evidence/out/a.log, docs/specs/<slug>/qa/evidence/out/b.log, docs/specs/<slug>/run.log; remove each file or commit it with git add -f, then retry

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix archive --help' for usage.
   exit: 2
   ```

2. Surface Transcript: cancel refuses Ignored Spec Files.

   ```transcript
   $ roundfix cancel <slug> --reason <text>
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "<slug>" cannot cancel while Git ignores 3 file(s) under docs/specs/<slug> that no commit holds: docs/specs/<slug>/qa/evidence/out/a.log, docs/specs/<slug>/qa/evidence/out/b.log, docs/specs/<slug>/run.log; remove each file or commit it with git add -f, then retry

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix cancel --help' for usage.
   exit: 2
   ```

3. Surface Transcript: archive recovers forward from a failed removal.

   ```transcript
   $ roundfix archive <slug>
   stdout:
   stderr:
   roundfix: archive wrote docs/history/specs/<slug>.md but could not remove the Spec folder docs/specs/<slug>: <cause>
   roundfix: finish the removal with: chmod -R u+w '<folder>' && rm -rf '<folder>'
   exit: 1
   ```

4. Surface Transcript: cancel recovers forward from a failed removal.

   ```transcript
   $ roundfix cancel <slug> --reason <text>
   stdout:
   stderr:
   roundfix: cancel wrote docs/history/specs/<slug>.md but could not remove the Spec folder docs/specs/<slug>: <cause>
   roundfix: finish the removal with: chmod -R u+w '<folder>' && rm -rf '<folder>'
   exit: 1
   ```

### Invariants

```text
1. Neither command removes a file that git status reports as ignored or untracked under the Spec folder.
2. Every refusal happens before the first write and changes no byte.
3. A failed folder removal never rolls back; the record and every other write stay.
4. The finishing command quotes the absolute folder for a POSIX shell.
5. Symbolic-link destination parents are accepted as before.
6. Uncommitted changes keep their refusal text and precedence.
```

## Coverage Map

Deviations from the PRD: None.

- User Story 1 → `readSpecFolderStatus`, `refuseSpecFolderStatus`, Surface Transcripts 1-2
- User Story 2 → `FolderRemovalError`, `folderRemovalFinishCommand`, Surface Transcripts 3-4
- Core Feature 1 → `readSpecFolderStatus`, `refuseSpecFolderStatus`, API Contract 1
- Core Feature 2 → `FolderRemovalError`, API Contracts 2-3
- Core Feature 3 → API Contract 4, user guides, Roundfix skill archive reference, archive-spec skill, `CONTEXT.md`
- Success Metric 1 → refusal tests in `internal/cli/spec_folder_status_test.go`
- Success Metric 2 → recovery tests in `internal/cli/folder_removal_recovery_test.go` and `internal/spec/folder_removal_test.go`
- Success Metric 3 → transcript tests and the QA gate

## Absence Evidence

1. Ignored file guard — command: `grep -rn -- '--ignored' internal/cli/archive.go internal/cli/cancel.go internal/spec/archive.go internal/spec/cancel.go`; result: no match at `94fbf85b`; both guards call `status --porcelain --untracked-files=all` only.
2. Shared Spec folder status helper — command: `grep -rln 'untracked-files=all' internal/cli`; result: `archive.go` and `cancel.go` each run their own call; the other production matches (`history.go`, `reconcile.go`, `spec_check.go`, `baseline_update_history.go`) read other trees for other commands, and the rest are tests; none is shared by archive and cancel.
3. Folder removal error type — command: `grep -rn 'leftover folder' internal`; result: two plain `fmt.Errorf` strings, `internal/spec/archive.go:420` and `internal/spec/cancel.go:303`, and no type a caller can match.
4. Finishing command — command: `grep -rln 'rm -rf' internal/cli internal/spec --include=*.go`; result: only `internal/cli/cli_test.go` at `94fbf85b`, a test fixture; no production code prints a removal command.
5. Glossary terms — command: `grep -nE '^[*][*](Ignored Spec File|Folder Removal Recovery)[*][*]' CONTEXT.md`; result: no match at `94fbf85b`.

## ADR Alternatives

- ADR-0269: this Spec's own decision; it rejects copying, committing or warning about ignored files, a rollback, and a bare `rm -rf`.
- ADR-0247: no rejected alternative is taken; the record stays small and the folder stays in Git, now for every removed byte.
- ADR-0248: no rejected alternative is taken; History Sanitize and its batches are untouched.
- ADR-0251: no rejected alternative is taken; Lenient Legacy Reading is untouched.
- ADR-0254: no rejected alternative is taken; a Baseline update's Pending History is untouched.
- ADR-0262: no rejected alternative is taken; the Cancel Command still records and returns sources, with one more refusal.
- ADR-0261: no rejected alternative is taken; `--pull-request` still writes only the record, also after a failed removal.
- ADR-0260: no rejected alternative is taken; relinked documents stay written, and the relink plan is unchanged.
- ADR-0184: no rejected alternative is taken; the four outputs are transcripts asserted whole.
- ADR-0256: no rejected alternative is taken; the skills Task updates the Roundfix skill with the user guides.
- ADR-0244: no rejected alternative is taken; the two terms follow its declaration contract.
- ADR-0257: no rejected alternative is taken; agents run focused tests and the Daemon runs the declared Verification.
- ADR-0258: no rejected alternative is taken; the Spec declares every file and break it touches.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; the new tests bound by real time only what they assert, and no test file gains a Real-Time Budget.

## Integration Points

None outside the repository. Git is read through `git status`, as today. The
Delivery Queue's archive step treats any non-zero exit as a failed archive and
parks the item, so a refusal or a failed removal parks it as before.

## Testing Approach

All new tests call `t.Parallel()` first and run on temporary repositories and
homes built by the existing fixtures (`newImplementWorkspace`,
`cancelFixture`, `writeArchiveQAReport`).

- `internal/cli/spec_folder_status_test.go`: `TestArchiveRefusesIgnoredSpecFiles`
  asserts Surface Transcript 1 whole and that the tree, the ignored files and
  `git status --ignored` are unchanged; `TestCancelRefusesIgnoredSpecFiles`
  does the same for Surface Transcript 2;
  `TestSpecFolderStatusRefusesUncommittedChangesFirst` keeps today's reason
  when both exist; `TestSpecFolderStatusReadsTheExternalSpecsRootRepository`
  uses an ignore rule that only the external root's repository holds.
- `internal/cli/folder_removal_recovery_test.go`:
  `TestArchiveFolderRemovalFailureRecoversForward` and
  `TestCancelFolderRemovalFailureRecoversForward` assert Surface Transcripts 3
  and 4 whole, that the record and other writes exist, then run the printed
  command with `sh -c` and assert the tree a successful run leaves;
  `TestArchiveFolderRemovalFailureRecordsThePullRequest` and
  `TestFolderRemovalFinishCommandQuotesTheFolder`, with a folder whose name
  holds a space and a `'`.
- `internal/spec/folder_removal_test.go`: `TestArchiveReturnsAFolderRemovalError`
  and `TestCancelReturnsAFolderRemovalError` match the error with `errors.As`
  and check the filled result and the written record.

The removal failure comes from a committed subdirectory made read-only with
`chmod 0555` and restored in `t.Cleanup`, as `TestCancelRollsBackAFailedWrite`
already does. Declared breaks: `TestCancelTranscriptPrintsHelp` takes the new
cancel help, `TestRunArchiveHelp` gains the new archive lines, and
`TestCancelWriteFailureUsesExitOne` keeps its write-failure prefix. The
Behavior Surface Record is re-recorded with
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
after each change to help or a user guide, and the Skill Coverage Map lists
`internal/cli/spec_folder_status.go` and `internal/cli/folder_removal_recovery.go`
among the sources of `command archive` and `command cancel`.

## Build Order

1. Shared Spec folder status helper, the Ignored Spec File refusal in both
   commands, help lines, user guides, coverage map and record, and tests.
2. `FolderRemovalError`, Folder Removal Recovery in both commands, help exit
   lines, user guides, record and tests (depends on: 1).
3. Roundfix skill archive reference, archive-spec skill, mirrors, versions and
   glossary (depends on: 1, 2).
4. QA gate (depends on: 3).

No other Spec is a prerequisite. If Spec 0258 merges first, this Spec already
carries its sections; rebase and run the strict check again.

## Risks & Considerations

- An ignored `.DS_Store` or log in a Spec folder now blocks archive and
  cancel; the reason names each file and the two ways out.
- A Delivery Queue item whose worktree holds an ignored Spec file parks at
  archive; the operator removes or commits it and retries, as for any failed
  archive.
- The read-only fixture cannot fail as root; the suite runs as an ordinary
  user locally and in CI, as `TestCancelRollsBackAFailedWrite` already
  assumes.
- `git status --ignored` walks ignored directories with
  `--untracked-files=all`; a Spec folder is small, so the cost is negligible.

## Glossary

None.

## Decisions

- One helper in `internal/cli` reads tracked changes, untracked files and
  Ignored Spec Files in one `git status` call. See ADR-0269.
- `internal/spec` returns a typed folder removal error with the filled result,
  and the command layer prints the recovery. See ADR-0269.
- The finishing command restores write permission before removing. See
  ADR-0269.
