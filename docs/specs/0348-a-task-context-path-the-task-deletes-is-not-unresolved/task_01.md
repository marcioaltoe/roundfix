---
task: task_01
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
status: completed
type: backend
complexity: medium
---

# Task 01: A Task Context path tracked at HEAD or the Delivery Base resolves and the fix stops saying Create

## Overview

The Spec Consistency Check reports `SC-REF-UNRESOLVED` for an `interface:`
or `instruction:` Task Context path that the checked tree no longer holds,
so a Task that removes a file it lists fails its own settlement and the
finding's fix ("Create <path> ...") leads the repair turn to swap in another
file. This Task makes such a path resolve when `HEAD` or the Delivery Base
tracks it, changes the fix to suggest removing the entry or declaring the
deletion, and proves both through `speccheck.Check` and the Daemon's Task
cycle. It answers the Fiscus report of 2026-11-04 (ai-memory message
`01a1a5c5`, adopted as
[references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md](references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md)).

## Requirements

1. MUST characterize first: add `TestTaskContextPathTrackedAtTheBaseResolves`
   in `internal/speccheck/task_context_base_test.go` with the twelve rows of
   `_techspec.md` → Testing Approach and today's counts (every row 1 except
   `present`), run it green on the unchanged check, and record that run in
   the Task Result; then set the counts after the change and implement.
2. MUST add `taskContextTrackedAtBase` and `taskContextUnresolvedFix` in the
   new `internal/speccheck/task_context_base.go`, unexported, exactly as
   `_techspec.md` → API Contracts 1 and 2 state: the top-level guard
   (`git rev-parse --show-toplevel` and `repoRoot` naming the same directory
   after `filepath.EvalSymlinks`), `HEAD` and the Delivery Base resolved with
   `preflight.ResolveDefaultBranchHead` and `git merge-base`, `git cat-file
   -e <revision>:<path>`, `-c core.fsmonitor=false` on every `git` call, no
   fetch, no ref write, and false on any Git failure.
3. MUST call `taskContextTrackedAtBase` in `detectTaskContextReferences`
   (`internal/speccheck/citations.go`) only after the three existing
   resolutions fail, and take the finding's fix from
   `taskContextUnresolvedFix`; code, severity, summary and locations keep
   their bytes, and the reference index and adopted source findings keep
   theirs.
4. MUST make every reporting row of the table assert the code, the
   severity, the exact summary `<task> declares unresolved Task Context path <path>`,
   both locations and the exact fix
   `Remove the Task Context entry for <path> from <task>, or declare it under deletes: when the Task removes it; replace it only with a path the Task itself reads or changes.`
   The fixtures write `refs/remotes/origin/main` with `git update-ref` and
   never contact a remote; the `nested below another work tree` row commits
   the path at an outer repository's top level and checks a plain
   subdirectory of it.
5. MUST change `TestTaskContextResolvesThroughTheArchiveRecord`
   (`internal/speccheck/archive_record_test.go`) only as `_techspec.md` →
   Declared breaks states: its fixture commits the archived Spec folder's
   removal; its assertions keep their expectations.
6. MUST add to `internal/daemon/task_deletes_settlement_test.go`
   `TestAGatedTaskWhoseInterfacePathItDeletesSettles` (modes `disk-only` and
   `staged`, the Context entry labeled `interface:`, `TaskCycle` completes
   one Task) and `TestSettlementStillRefusesAContextPathNoRevisionTracks`
   (the turn adds an `interface:` entry for a path no revision tracks; the
   Task fails and the spec consistency diagnostics carry
   `SC-REF-UNRESOLVED`, the path and the new fix), reusing
   `retiredHistoryDeletionFixture`. Before the change the first test fails
   with `SC-REF-UNRESOLVED`; record that in the Task Result.
7. MUST start every new top-level test with `t.Parallel()`, keep every test
   hermetic (temporary repositories, `gittest`, no network, no real
   `~/.roundfix`), and start no process that outlives its test.
8. MUST NOT change `SC-TOOLING-UNDECLARED`, the mechanical authorization
   audit, the scope audit, the `creates:` and `deletes:` handling, the
   `deletes:` settlement refusal, any skill or guide, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize today's check with the case table.
- [ ] Add the base reader and the fix, and call them from the Task Context check.
- [ ] Commit the removal in the archive record fixture.
- [ ] Add the two settlement tests.

## Acceptance Criteria

- [ ] A missing `interface:` or `instruction:` path that `HEAD` or the Delivery Base tracks reports no `SC-REF-UNRESOLVED`; every other missing path reports one with today's summary.
- [ ] The fix of the Task Context finding is the API Contract 2 text.
- [ ] A gated Task that deletes its `interface:` path settles; a Context path no revision tracks still fails settlement with the new fix.

## Context

- interface: `internal/speccheck/citations.go`
- creates: `internal/speccheck/task_context_base.go`
- creates: `internal/speccheck/task_context_base_test.go`
- interface: `internal/speccheck/archive_record_test.go`
- interface: `internal/daemon/task_deletes_settlement_test.go`

## Consumers

- The Task Context `SC-REF-UNRESOLVED` finding read by `speccheck.Check` — consumers: `./internal/speccheck`, `./internal/daemon`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTaskContextPathTrackedAtTheBaseResolves$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTaskContextPathTrackedAtTheBaseResolves ' || { echo 'missing pass: TestTaskContextPathTrackedAtTheBaseResolves' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/speccheck` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes with the whole Spec Check suite, the archive record fixture included.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAGatedTaskWhoseInterfacePathItDeletesSettles|TestSettlementStillRefusesAContextPathNoRevisionTracks)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAGatedTaskWhoseInterfacePathItDeletesSettles TestSettlementStillRefusesAContextPathNoRevisionTracks; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/daemon` — expected: exit 0. Before this Task neither test exists; afterwards both pass with the whole Daemon suite.

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 1-4; API Contracts 1-2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, systematic-debugging, no-workarounds, golang-context, golang-testing, testing-boss, tech-writer, evidence-gate.

Pre-change characterization: `go test -count=1 -timeout 5m -v -run '^TestTaskContextPathTrackedAtTheBaseResolves$' ./internal/speccheck` exited 0 with all twelve rows on the unchanged production check. Counts were 0 for `present` and 1 for every other row; reporting rows pinned the existing code, severity, summary, both locations, and old fix. An initial fixture-only failure (missing PRD frontmatter) was corrected before this run.

Pre-change settlement reproduction: `go test -count=1 -timeout 5m -v -run '^TestAGatedTaskWhoseInterfacePathItDeletesSettles$' ./internal/daemon` exited 1. Both `disk-only` and `staged` produced `Completed:0 Failed:1` with `settlement check: spec consistency: SC-REF-UNRESOLVED` naming `packages/backend/test/gone.test.ts`.

Implemented the two unexported helpers and invoked the base reader only after all three existing Task Context resolutions fail. The reader verifies the symlink-resolved Git top level, checks `HEAD`, then the Delivery Base from `preflight.ResolveDefaultBranchHead` and `git merge-base`. All Git calls use `preflight.ExecGitRunner`, which disables fsmonitor; reads are local, synchronous, and write no refs. Unavailable revisions leave the finding standing. Only the Task Context finding's fix changes; its code, severity, summary, locations, and other reference findings retain their existing construction.

Focused acceptance evidence:

- Missing tracked `interface:` and `instruction:` paths: the twelve-row table now passes with the five specified removal rows changed from 1 to 0. The six unresolved rows retain count 1 and assert the complete finding, including both locations and the exact API Contract 2 fix. Temporary fixtures use `gittest`, local `update-ref` calls, and an outer-root tracked source for the nested-root refusal.
- Existing archive resolution: the active-spelling archive-record fixture commits its source-folder removal; all original assertions and the archived-spelling fixture retain their expectations.
- Settlement: both `disk-only` and `staged` interface deletions now complete one Task through `TaskCycle`. The never-tracked Context entry still fails settlement without a commit, and its spec-consistency diagnostic log contains `SC-REF-UNRESOLVED`, the path, and the complete new fix. Existing declared-deletion settlement, undeleted-source refusal, broken-link refusal, and unauthorized governed-deletion checks also pass.

Post-change focused commands:

- `go test -count=1 -timeout 5m -v -run '^(TestTaskContextPathTrackedAtTheBaseResolves|TestTaskContextResolvesThroughTheArchiveRecord|TestTaskContextResolvesThroughTheArchive)$' ./internal/speccheck` — exit 0; all twelve table rows and both archive routes passed.
- `go test -count=1 -timeout 5m -v -run '^(TestAGatedTaskWhoseInterfacePathItDeletesSettles|TestSettlementStillRefusesAContextPathNoRevisionTracks|TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles|TestSettlementStillReportsALinkToADeletedFile|TestSettlementRefusesAnUnauthorizedGovernedDeletion|TestATaskThatLeavesItsDeletesPathDoesNotSettle|TestATaskThatRemovesItsDeletesPathSettles)$' ./internal/daemon` — exit 0; seven top-level tests passed, including both deletion modes.
- `git diff --check` — exit 0. Changed-path review found only the five declared implementation/test paths and this Task file; the Task file was already modified at entry for Daemon-owned `status: in_progress`, which was preserved.

Every new top-level test starts with `t.Parallel()`; the new subtests are parallel too. The focused runs' suite guard reported no repository mutation or helper residue. No skill, guide, tooling configuration, Task Graph, or other Task was edited. Declared Verification and repository Verification were not run; the Daemon owns them and terminal settlement. No commits, pushes, or pull requests were made.
