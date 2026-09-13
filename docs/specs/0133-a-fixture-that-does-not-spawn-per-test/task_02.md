---
task: task_02
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: completed
type: test
complexity: medium
---

# Task 02: Stop paying Git setup per Task-cycle test

## Overview

Giving the Task-cycle fixtures real Git provenance put repository setup in the
base fixture that every Task-cycle test builds. Measured on 2026-09-13: the
Daemon package runs 3.2s at the delivery target and 6.8s here, and the base
fixture went from zero Git invocations to three plus a revision read and an
authorization resolution — roughly five process spawns per test across about
forty tests. This suite is spawn-bound, so that setup is the whole delta. The
provenance stays; paying for it once per test is the defect.

`internal/daemon/task_engine_test.go` is not a Governed Path, so this Task needs
no tooling grant.

## Requirements

1. MUST stop performing Git repository setup inside the per-test base fixture.
   Create the committed seed once for the package and let each test reuse it, or
   obtain the same provenance without spawning a process per test.
2. MUST keep every journey's provenance real: the authorization reader must
   still read committed bytes, and no test may pass because the reader was
   weakened, stubbed, or given a non-Git shortcut.
3. MUST keep the three external and symlinked Spec Root journeys passing, since
   they exercise a different root shape and were the reason provenance was added.
4. MUST bring the Daemon package's wall clock back into the delivery target's
   range, measured as one concurrent package run rather than in isolation.
5. MUST NOT raise a time budget, lengthen a deadline, reduce parallelism, or
   skip a case. The cost goes, not the signal.
6. MUST NOT change production code. The defect is in test setup; production
   already resolves the authorization once per Implement and once per Settle.

## Subtasks

- [ ] Create the committed seed once for the package instead of per test.
- [ ] Point the base fixture at that shared seed.
- [ ] Confirm the external and symlinked journeys still read real provenance.
- [ ] Measure the package against the delivery target.

## Acceptance Criteria

- [ ] The per-test base fixture performs no Git repository initialization; a
      search of its body finds no such call, where it finds three today.
- [ ] A named test asserts the shared seed is created once and reused, so a
      later change reintroducing per-test setup fails here.
- [ ] The three external and symlinked Spec Root journeys pass, and the
      authorization reader still refuses a Spec Root with no Git provenance.
- [ ] The Daemon package passes as one concurrent run and its wall clock is
      within the delivery target's range rather than roughly double it.
- [ ] No production file changed, and no time budget, deadline, parallelism
      setting or skip was introduced.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/gittest/gittest.go`

## Verification

- `body="$(sed -n '/^func newTaskCycleFixture/,/^}$/p' internal/daemon/task_engine_test.go)"; test -n "$body" || exit 1; printf '%s' "$body" | grep -q 'InitRepo' && exit 1; exit 0` — the per-test base fixture no longer initializes a repository; it does today, so this fails before the work.
- `grep -q 'func TestTaskCycleFixtureSeedIsCreatedOnce' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestTaskCycleFixtureSeedIsCreatedOnce$'` — the shared seed is asserted to be created once and reused.
- `grep -q 'func TestTaskCycleFixtureSeedIsCreatedOnce' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^(TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged|TestTaskCycleQAReportExternalProceedsWithoutStaging|TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths)$'` — the three provenance journeys still pass once the shared seed is in place. The guard reads the new test, so this preservation check cannot pass before the work.
- `grep -q 'func TestTaskCycleFixtureSeedIsCreatedOnce' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon` — the package passes as one concurrent run with the shared seed in place.
- `grep -q 'func TestTaskCycleFixtureSeedIsCreatedOnce' internal/daemon/task_engine_test.go || exit 1; widened="$(grep -n 'time.After(' internal/daemon/task_engine_test.go | grep -v '2 \* time.Second' || true)"; test -z "$widened" || { printf '%s\n' "$widened"; exit 1; }; changed="$(git diff --name-only HEAD -- internal/ | grep -v '^internal/daemon/task_engine_test.go' || true)"; test -z "$changed" || { printf '%s\n' "$changed"; exit 1; }` — every deadline still reads two seconds and no file outside the fixture moved. The guard reads the new test so this cannot pass before the work.

## References

- `_prd.md` → User Stories 3-4; Core Features 4-7; Decisions: Regression locks.
- `_techspec.md` → Risks & Considerations.
- Spec 0132 qa/qa-report-2026-09-11-01.md → the passing gate this Task must not regress.

## Result

### Implementation

- The package now initializes and commits one Task-cycle repository seed behind
  `sync.Once`, snapshots its regular files into an immutable in-memory file
  system, and copies that history into each fixture without a Git process.
- `newTaskCycleFixture` overlays each journey's Task source on its private copy.
  Tests can still write and commit independently, while the real authorization
  reader resolves the authorization record from the copied commit.
- `TestTaskCycleFixtureSeedIsCreatedOnce` creates two fixtures with different
  Task graphs and asserts one seed creation, shared seed identity, distinct Git
  roots, isolated writes, isolated graphs, and committed authorization
  provenance.

### Focused checks

- Red signal: `rtk go test ./internal/daemon -run 'TaskCycleFixtureSeedIsCreatedOnce' -count=1`
  reached the expected build failure on the missing seed seam after one
  unchanged retry with Go cache access.
- `rtk go test ./internal/daemon -run 'TaskCycleFixtureSeedIsCreatedOnce' -count=1`
  passed after the implementation (1 test).
- `rtk go test ./internal/daemon ./internal/spec -run 'TaskCycleFixtureSeedIsCreatedOnce|External|Symlink|OperationAuthorityReportsUnresolvableSpecRoot|Test(Task|Governed|Commit|Verification)' -count=1`
  passed after the last code edit (144 tests in 2 packages). This selection
  includes the three named external/symlink journeys and the existing refusal
  for a Spec Root without Git provenance.
- `rtk git diff --check` passed.

### Acceptance evidence

1. Source inspection after the edit shows `newTaskCycleFixture` calls only the
   shared seed accessor and filesystem copy before writing its journey data; it
   contains no `InitRepo`, `git add`, or `git commit` call.
2. The named regression test observes two fixture constructions and requires
   `taskCycleFixtureSeedCreations == 1` plus identical seed pointers.
3. The 144-test focused run exercised the external Task settlement, external QA
   Report, symlink-crossing commit, and unresolvable Spec Root cases through
   their real Git boundaries.
4. The focused Task-cycle run completed in 4.80 seconds including 144 selected
   tests across two packages. The full concurrent Daemon-package timing remains
   for Daemon-owned authored Verification.
5. The changed-path inspection under `internal/` names only
   `internal/daemon/task_engine_test.go`; the diff adds no deadline, time budget,
   parallelism change, or skip. The only other changed path is this assigned
   Task file, whose pre-existing status edit remains Daemon-owned.
