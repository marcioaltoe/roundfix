---
task: task_01
spec: 0137-a-fixture-that-carries-its-committer
status: completed
type: test
complexity: low
---

# Task 01: Write the committer identity into the shared fixture seed

## Overview

The shared Task-cycle fixture seed initialises a repository and commits into it
through the test helper, which supplies the committer identity as a
per-invocation override. Overrides are not written to disk and repository
hardening appends only maintenance keys, so a copy of the seed carries no
identity and the Daemon's own committer depends on host discovery. The setup
Spec 0133 replaced wrote the identity into the repository configuration; this
slice restores that in the seed.

## Requirements

1. MUST write a committer identity into the shared fixture seed's repository
   configuration, so every repository copied from the seed inherits it.
2. MUST assert the identity by reading the copied repository's resolved
   configuration with host, global and system configuration excluded from the
   answer. An assertion that only commits successfully passes on a machine
   whose Git can compose an identity from the local user and hostname, which is
   the blind spot that let this reach a Pull Request gate.
3. MUST keep the seed created once, so the per-test cost Spec 0133 removed
   stays removed.
4. MUST keep the Daemon's real-repository Task-cycle behavior unchanged: one
   commit per Task, pre-existing dirt excluded.
5. MUST NOT change production code, and MUST NOT weaken or delete an existing
   assertion.

## Subtasks

- [ ] Persist the identity in the seed's repository configuration.
- [ ] Assert it from a copy with host configuration excluded.
- [ ] Confirm the seed is still created once and the journey is unchanged.

## Acceptance Criteria

- [ ] A repository copied from the seed resolves a committer name and email
      from its own configuration with host, global and system configuration
      excluded; today it resolves none.
- [ ] The real-repository Task-cycle journey still commits once per Task and
      still excludes pre-existing dirt.
- [ ] The seed is still created once.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/gittest/gittest.go`

## Verification

- `grep -q 'func TestTaskCycleFixtureSeedCarriesACommitterIdentity' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestTaskCycleFixtureSeedCarriesACommitterIdentity$'` — a copy resolves the identity from its own configuration; this fails today.
- `grep -q 'func TestTaskCycleFixtureSeedCarriesACommitterIdentity' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt$'` — the real-repository journey is unchanged.
- `grep -q 'func TestTaskCycleFixtureSeedCarriesACommitterIdentity' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestTaskCycleFixtureSeedIsCreatedOnce$'` — the seed is still created once.

## References

- `_prd.md` → Goal 1; Core Features 1-3; Regression locks.
- `_techspec.md` → Implementation Design: Write the identity where a copy can
  inherit it; Testing Approach observations 1-3; Build Order 1.

## Result

- Implementation: the once-only Task-cycle fixture seed now persists the
  existing deterministic test identity before snapshotting the repository.
  Copies therefore inherit the identity through their own `.git/config`.
- Pre-change signal: `rtk proxy env
  GOCACHE=/private/tmp/roundfix-go-build-0137-task-01 rtk go test -count=1
  ./internal/daemon -run
  '^(TestTaskCycleFixtureSeedCarriesACommitterIdentity|TestTaskCycleFixtureSeedIsCreatedOnce)$'`
  reported the seed-created-once test passing and both copied-repository
  identity reads failing with exit status 1.
- Focused check: `rtk proxy env
  GOCACHE=/private/tmp/roundfix-go-build-0137-task-01 rtk go test -count=1
  ./internal/daemon -run
  '^(TestTaskCycleFixtureSeedCarriesACommitterIdentity|TestTaskCycleFixtureSeedIsCreatedOnce|TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt)$'`
  exited 0 and reported five passing tests.
- Acceptance criterion 1 evidence: the new fixture-seed test copies the shared
  seed, invokes plain Git without the test helper's command-line identity, and
  uses the isolated Git environment that excludes inherited global and system
  configuration. Both `user.name` and `user.email` resolved as non-empty.
- Acceptance criterion 2 evidence: the focused run exercised the existing real
  repository journey without changing it; its assertions still cover two Task
  commits and keep `user-wip.txt` uncommitted.
- Acceptance criterion 3 evidence: the focused run exercised the existing
  pointer, creation-count, copy-isolation, graph-isolation, and provenance
  assertions in `TestTaskCycleFixtureSeedIsCreatedOnce`; the identity write
  remains inside the existing `sync.Once` body.
- Repository incremental check: `rtk make verify` exited 2 at `fmt-check`
  because the unchanged files `internal/cli/baseline_skills_restore_test.go`
  and `internal/cli/baseline_assets_sync_test.go` need formatting. They are
  outside this Task's slice and were not edited.
- Daemon Verification was not run in this Agent turn.
