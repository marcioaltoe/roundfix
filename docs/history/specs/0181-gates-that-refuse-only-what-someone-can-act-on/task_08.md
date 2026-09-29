---
task: task_08
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: completed
type: backend
complexity: medium
---

# Task 08: A Task commit stages what it records, and the record never eats authored text

## Overview

The pre-PR review of the 0181 candidate found two defects in the recorded-paths slice.

- `prepareTaskCommit` records the files that `expandUntrackedCommitDirectories` lists under an untracked directory, using `git ls-files --others --exclude-standard`. The commit, however, still stages the raw directory entry through `git add -f -- <paths>`. The force flag adds ignored files inside that directory too, so an ignored file, possibly a secret, can be committed without appearing in `## Recorded paths`.
- `recordedPathsSectionOffset` in `internal/spec/recorded_paths.go` finds the section by raw substring. A `## Recorded paths` line inside authored prose or a fenced example therefore matches, and replacing the section truncates every authored line after it.

## Requirements

1. MUST make the Daemon's Task commit stage exactly the expanded file list, the same list the record is computed from, instead of an untracked directory entry. A file Git ignores inside a new directory MUST never be staged by a Task commit. Tracked removals and ordinary files keep staging as today.
2. MUST make `recordedPathsSectionOffset` recognize the section only as a line exactly equal to `## Recorded paths`, outside any fenced code block, that starts the file's last `## ` section. A matching line inside a fence, or one followed by another `## ` heading, MUST NOT be treated as the Daemon's section. Writing, replacing and reading the section MUST all use this rule, and every authored byte MUST survive a replacement.
3. MUST change no exported function signature and rename or remove no top-level test. It MUST update only the existing tests this change invalidates and name each in its Result.
4. MUST put the new tests in `internal/daemon/task_commit_ignored_test.go`, driving the existing Task cycle fixture over a `gittest` repository, and in `internal/spec/recorded_paths_section_test.go`.
5. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`, because this Spec's QA gate still runs on the v0.20.0 auditor.

## Subtasks

- [ ] Stage the expanded file list in the Task commit.
- [ ] Anchor the section to an exact, unfenced, final heading.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A Task that creates a new directory holding one ordinary file and one file its `.gitignore` excludes commits and records only the ordinary file.
- [ ] A Task file whose authored text contains `## Recorded paths` inside a fenced block keeps every byte of that block when the Daemon writes its section.
- [ ] A Task file whose authored text has a `## Recorded paths` heading followed by later `## ` sections keeps those sections, and the Daemon's section is appended at the end.
- [ ] Replacing an existing Daemon section at the end of the file still replaces it rather than appending a second one.

## Context

- instruction: `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/spec/recorded_paths.go`
- creates: `internal/daemon/task_commit_ignored_test.go`
- creates: `internal/spec/recorded_paths_section_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory|TestRecordedPathsKeepAFencedHeadingInAuthoredText|TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections|TestRecordedPathsReplaceTheTrailingDaemonSection|TestTaskCommitRecordsTheFilesOfANewUntrackedPackage|TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte)$" ./internal/daemon ./internal/spec 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory TestRecordedPathsKeepAFencedHeadingInAuthoredText TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections TestRecordedPathsReplaceTheTrailingDaemonSection TestTaskCommitRecordsTheFilesOfANewUntrackedPackage TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four new named tests exists, so the command fails.

## References

- `_prd.md` → Core Feature 1
- `_techspec.md` → Recorded paths
- ADR-0166

## Result

The Task commit now stages the expanded file list used to compute recorded
paths. An ignored file under a newly created directory therefore stays out of
both the commit and the Task file record. Recorded-path section discovery now
accepts only the exact unfenced heading that starts the final level-two
section, so authored fenced examples and earlier authored sections remain
byte-identical when the Daemon appends its section.

New tests:

- `TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory`
- `TestRecordedPathsKeepAFencedHeadingInAuthoredText`
- `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections`
- `TestRecordedPathsReplaceTheTrailingDaemonSection`

Existing top-level tests changed: none. The behavior change invalidated no
existing assertion. `TestTaskCommitRecordsTheFilesOfANewUntrackedPackage` and
`TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte` passed
unchanged in the focused checks.

Focused implementation evidence:

- Before the production edit,
  `TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory` failed because the
  commit contained both `internal/newpkg/ignored.secret` and
  `internal/newpkg/new.go`.
- Before the production edit, the section tests showed that
  `RecordedTaskPaths` read a fenced authored bullet and that writing a record
  truncated an authored `## Notes` section.
- `rtk env GOCACHE=/tmp/roundfix-task08-gocache go test -count=1 -run '^(TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory|TestTaskCommitRecordsTheFilesOfANewUntrackedPackage)$' ./internal/daemon`: passed.
- `rtk env GOCACHE=/tmp/roundfix-task08-gocache go test -count=1 -run '^(TestRecordedPathsKeepAFencedHeadingInAuthoredText|TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections|TestRecordedPathsReplaceTheTrailingDaemonSection|TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte)$' ./internal/spec`: passed.
- `rtk env GOCACHE=/tmp/roundfix-task08-gocache go test -count=1 ./internal/daemon ./internal/spec`: passed.
- `rtk make verify-incremental`: passed after the sandbox-blocked first attempt
  was rerun with access to its test-time network dependency.
- `rtk git diff --check`: passed.

Acceptance evidence:

- `TestTaskCommitNeverStagesAnIgnoredFileInANewDirectory` drives the real Git
  Task-cycle fixture and proves the commit and record contain only the ordinary
  file.
- `TestRecordedPathsKeepAFencedHeadingInAuthoredText` proves the reader ignores
  the fenced heading and the writer preserves every authored byte before
  appending the Daemon section.
- `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections` proves an earlier
  authored heading is not owned by the Daemon and all later authored sections
  survive.
- `TestRecordedPathsReplaceTheTrailingDaemonSection` proves a trailing Daemon
  section is replaced once, with its preceding authored bytes preserved.

The Daemon must run the declared Verification and settle the Task; Task status
was left unchanged.
