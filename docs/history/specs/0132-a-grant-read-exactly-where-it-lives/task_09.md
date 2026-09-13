---
task: task_09
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
type: test
complexity: low
---

# Task 09: Give the Daemon's external Spec Root fixtures real Git provenance

## Overview

Deriving the record path from the resolved Spec Root made the authorization
reader ask that root for committed bytes. The Daemon's external-root fixtures
build a plain temporary directory, so three existing journeys now refuse
Implement or commit authority with `not a git repository`. The production
external-root helper creates a committed repository; the test data contract
diverged from it. This slice closes that divergence, and it is verifiable alone:
the three journeys reach their settlement assertions again.

`internal/daemon/task_engine_test.go` is not a Governed Path, so this Task needs
no tooling grant.

## Requirements

1. MUST give the external Spec Root fixture real Git provenance: initialize a
   repository at that root and commit the seeded Spec, so the authorization
   reader can read committed bytes from it.
2. MUST give the symlinked Spec Root fixture the same provenance, since it
   reaches the same reader through a different path shape.
3. MUST restore the three affected journeys to asserting settlement and staging
   behavior, not authorization refusal. A journey that now passes by expecting
   the refusal would record the defect as the contract.
4. MUST keep every journey's existing assertion about settlement, staging and
   commit content unchanged; this Task repairs the fixture's provenance, not
   what the journeys prove.
5. MUST NOT weaken the authorization reader, add a bypass for test data, or make
   the reader accept a non-Git Spec Root. The fixture moves to meet production,
   never the reverse.
6. MUST reuse the repository's existing Git test helper rather than adding a
   second way to create a fixture repository.

## Subtasks

- [ ] Initialize and commit the external Spec Root fixture.
- [ ] Do the same for the symlinked Spec Root fixture.
- [ ] Confirm the three journeys assert settlement rather than refusal.
- [ ] Confirm no reader or production path changed.

## Acceptance Criteria

- [ ] The external Spec Root fixture is a Git repository whose seeded Spec is
      committed, created through the existing Git test helper.
- [ ] `TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged`
      and `TestTaskCycleQAReportExternalProceedsWithoutStaging` reach their
      settlement and staging assertions; both fail today with
      `not a git repository`.
- [ ] `TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths`
      reaches its commit-content assertion through the symlinked root.
- [ ] No file outside `internal/daemon/task_engine_test.go` and this Task file
      changed, proven by inspecting the changed paths rather than by assertion.
- [ ] The authorization reader still refuses a Spec Root that genuinely carries
      no Git provenance.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/gittest/gittest.go`

## Verification

- `grep -q 'gittest' internal/daemon/task_engine_test.go || exit 1; body="$(sed -n '/^func (fixture \*taskCycleFixture) useExternalSpecRoot/,/^}$/p' internal/daemon/task_engine_test.go)"; test -n "$body" || exit 1; printf '%s' "$body" | grep -q 'gittest'` — the external root fixture creates its repository through the shared Git helper; today that function body contains no Git call at all.
- `go test -count=1 ./internal/daemon -run '^(TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged|TestTaskCycleQAReportExternalProceedsWithoutStaging)$'` — both external-root journeys pass; both exit 1 today with `not a git repository`.
- `go test -count=1 ./internal/daemon -run '^TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths$'` — the symlinked root journey passes.

## References

- `_prd.md` → Goals 2; Core Features 2.
- `_techspec.md` → Implementation Design: Two roots, named separately.
- `qa/qa-report-2026-09-11.md` → F-001.

## Result

Implemented the fixture-side provenance repair. `useExternalSpecRoot` now
initializes its Spec Root with `gittest.InitRepo`, seeds the Spec, and commits
those bytes through `gittest.Run` before loading it. The symlink-crossing
journey initializes and commits its target Spec Root the same way. The existing
settlement, staging, event, warning, and commit-content assertions are
unchanged, and no authorization reader or production path changed.

Focused evidence by acceptance criterion:

1. Before the edit,
   `GOCACHE=/private/tmp/roundfix-task09-go-cache rtk proxy go test ./internal/daemon -run TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged -count=1 -v`
   exited 1 at Implement dispatch with `fatal: not a git repository`. The source
   diff now shows `gittest.InitRepo`, `git add -A`, and a fixture commit before
   `spec.Load` in `useExternalSpecRoot`.
2. `GOCACHE=/private/tmp/roundfix-task09-go-cache rtk go test ./internal/daemon -run 'TestTask(CycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged|CycleQAReportExternalProceedsWithoutStaging|CommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths)' -count=1 -v`
   passed all three selected journeys, including both external settlement and
   staging cases.
3. The same focused three-journey command passed
   `TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths`; its
   existing repository-path commit-content and dropped symlink-path assertions
   are byte-unchanged.
4. The preflight `rtk git status --short` showed only this Task file's
   Daemon-owned status edit. Final `rtk git diff --name-only` inspection listed
   exactly this Task file and `internal/daemon/task_engine_test.go`; the Go diff
   contains six fixture-setup additions and no assertion or production change.
5. `GOCACHE=/private/tmp/roundfix-task09-go-cache rtk go test ./internal/spec -run TestOperationAuthorityReportsUnresolvableSpecRoot -count=1 -v`
   passed the existing negative control: a genuine non-Git Spec Root remains
   unresolved with the `unreadable_record` reason on `spec_root`, so it grants
   no operation.

The Daemon-owned commands under `## Verification` were not run.
