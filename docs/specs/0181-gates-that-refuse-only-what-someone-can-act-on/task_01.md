---
task: task_01
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: backend
complexity: high
---

# Task 01: Every Task commit records the paths the Task changed without declaring them

## Overview

A Task's `## Context` names the paths it plans to edit. An Agent that needs a file it did not plan, usually a new test file, cannot declare it: it may edit only its Task's `## Result`. The authored QA gate then refuses a completed, verified Task on that path. This Task makes both writers of a Task's standard commit record those paths themselves. Those writers are the Daemon's `commitTask` in `internal/daemon/task_engine.go` and the Settle Command's `settleTaskAndCommit` in `internal/cli/settle.go`. The record is a Daemon-owned `## Recorded paths` section at the end of the Task file, committed with the work (ADR-0166), and the Daemon's Task commit event lists the same paths. The record is read later by the QA gate and by reviewers, in the Task file on the Run and item branches. It is never read by Task Carry-Forward, path reservation or the consistency check.

## Requirements

1. MUST add `internal/spec/recorded_paths.go` with `RecordedPathsHeading`, `UndeclaredTaskPaths`, `RecordTaskPaths` and `RecordedTaskPaths`, as the TechSpec's "Recorded paths" section states:
   - `UndeclaredTaskPaths` returns, sorted and unique, every committed path that is not the Task file, not declared as `interface:` or `creates:` in the Task's Context, and not reported true by the `governed` function the caller passes. An edited `instruction:` path is recorded.
   - `RecordTaskPaths` writes the section at the end of the file, replacing an existing one and preserving every other byte. It writes nothing for an empty list, except to remove an existing section. It refuses, with an error and without writing, a path containing a backtick, a carriage return or a line feed.
   - `RecordedTaskPaths` reads the paths back.
2. MUST keep the Task parser, `CarryForwardInputs`, the fifty-entry Context limit and `RecordCarryForward` unaware of the section. A Task file carrying it MUST parse to the same Context and yield the same `CarryForwardInputs` as without it. The authored `## Context` MUST stay byte-identical.
3. MUST make `prepareTaskCommit` compute `recorded` for every Task whose type is not `qa`. It uses `spec.UndeclaredTaskPaths` with `speccheck.GovernedPath` over the staged paths, after expanding each untracked directory entry into the files `git ls-files --others --exclude-standard -z -- <dir>` lists under it.
4. MUST make `commitTask`, which runs only for a completed Task, write the section with `spec.RecordTaskPaths` before `Committer.Commit` when `recorded` is not empty, so the section rides in the same Task commit. The `daemon.commit` event payload MUST carry `recorded_paths` with the same paths only when it is not empty. A failed Task, a QA Task and a Task whose every changed path is declared MUST record nothing and carry no `recorded_paths` key.
5. MUST make `settleTaskAndCommit` compute the same record for a non-QA Task from the paths `stagedSettlePaths` reports, deleted paths included. It writes the section with `spec.RecordTaskPaths` and stages the Task file again before `Committer.Commit`. Its stdout, stderr and exit codes MUST not change otherwise.
6. MUST never record a Governed Path. A Governed Path stays under its authorization exactly as before.
7. MUST name the section in `docs/user-guide/commands.md`, both where the Task commit an Implement Run creates is described and in the `roundfix settle` section, including that a recorded path is disclosed and reserves nothing.
8. MUST change no exported function signature and rename or remove no top-level test. It MUST update only the existing tests this change invalidates and name each in the Result.
9. MUST put the new tests in `internal/spec/recorded_paths_test.go`, `internal/daemon/recorded_paths_test.go` and `internal/cli/settle_recorded_paths_test.go`. The daemon tests MUST drive the existing Task cycle fixture (`taskFakeRunner` over a `gittest` repository), not a hand-built engine. One daemon test MUST reproduce Spec 0180's `task_02` shape: the Task declares `interface: internal/store/delivery.go` and also changes `internal/store/delivery_test.go`.

## Subtasks

- [ ] Add the spec-package record, writer and reader.
- [ ] Record in the Daemon's Task commit and its commit event.
- [ ] Record in the Settle Command's Task commit.
- [ ] Describe the section in the commands guide.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A completed Task that declares `interface: internal/store/delivery.go` and changes `internal/store/delivery_test.go` commits a Task file whose `## Recorded paths` names exactly `internal/store/delivery_test.go`, and its `daemon.commit` event carries the same `recorded_paths`.
- [ ] A Task that adds a new untracked package records each file in it, not the directory.
- [ ] A Task whose every changed path is declared, a failed Task and a QA Task leave no section and no `recorded_paths` key.
- [ ] An authorized Governed Path a Task changes is never recorded.
- [ ] Replacing a section keeps every other byte of the Task file, an empty list removes the section, and an unrecordable path is refused without a write.
- [ ] A Task file with the section parses to the same Context and yields the same `CarryForwardInputs`.
- [ ] `roundfix settle` records the undeclared paths of the Task commit it creates, and records nothing when every path is declared.

## Context

- instruction: `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- instruction: `internal/spec/spec.go`
- creates: `internal/spec/recorded_paths.go`
- creates: `internal/spec/recorded_paths_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/recorded_paths_test.go`
- interface: `internal/cli/settle.go`
- creates: `internal/cli/settle_recorded_paths_test.go`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestUndeclaredTaskPathsExcludesDeclaredGovernedAndTheTaskFile|TestUndeclaredTaskPathsRecordsAnEditedInstructionPath|TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte|TestRecordTaskPathsWithNoPathsRemovesTheSection|TestRecordTaskPathsRefusesAnUnrecordablePath|TestARecordedSectionLeavesContextAndCarryForwardInputsUnchanged|TestTaskCommitRecordsAnUndeclaredTestFile|TestTaskCommitRecordsTheFilesOfANewUntrackedPackage|TestTaskCommitRecordsNothingWhenEveryPathIsDeclared|TestTaskCommitNeverRecordsAGovernedPath|TestAFailedTaskRecordsNoPaths|TestAQATaskRecordsNoPaths|TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit|TestSettleRecordsNothingWhenEveryPathIsDeclared|TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage)$" ./internal/spec ./internal/daemon ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUndeclaredTaskPathsExcludesDeclaredGovernedAndTheTaskFile TestUndeclaredTaskPathsRecordsAnEditedInstructionPath TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte TestRecordTaskPathsWithNoPathsRemovesTheSection TestRecordTaskPathsRefusesAnUnrecordablePath TestARecordedSectionLeavesContextAndCarryForwardInputsUnchanged TestTaskCommitRecordsAnUndeclaredTestFile TestTaskCommitRecordsTheFilesOfANewUntrackedPackage TestTaskCommitRecordsNothingWhenEveryPathIsDeclared TestTaskCommitNeverRecordsAGovernedPath TestAFailedTaskRecordsNoPaths TestAQATaskRecordsNoPaths TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit TestSettleRecordsNothingWhenEveryPathIsDeclared TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "## Recorded paths" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands.md "## Recorded paths" >&2; exit 1; }` — expected: exit 0; before this Task none of the fourteen new named tests exists and the guide does not name the section, so the command fails.

## References

- [_techspec.md](_techspec.md) — Recorded paths
- `_prd.md` → Goal 1; Goal 4; Core Feature 1; Success Metric 1
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Testing Approach 1
- [references/2026-09-29-a-file-a-task-creates-fails-the-qa-scope-audit.md](references/2026-09-29-a-file-a-task-creates-fails-the-qa-scope-audit.md)

## Result
