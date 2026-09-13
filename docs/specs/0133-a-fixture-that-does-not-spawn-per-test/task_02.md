---
task: task_02
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: pending
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
