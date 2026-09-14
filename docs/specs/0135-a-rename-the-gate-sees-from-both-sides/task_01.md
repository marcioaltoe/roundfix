---
task: task_01
spec: 0135-a-rename-the-gate-sees-from-both-sides
status: completed
type: backend
complexity: medium
---

# Task 01: Carry a rename's source through the changed-path reader

## Overview

The changed-path reader appends a porcelain record's destination and then skips
the field carrying its source, so renaming a Governed Path to an ungoverned
name produces a changed set naming only the ungoverned destination. No
comparison of snapshot pairs can recover a path that never enters either
snapshot. This slice preserves both sides where the record is read, and proves
it against porcelain output produced by a real rename.

## Requirements

1. MUST emit both the destination and the source of a rename or copy record, so
   a rename of a Governed Path presents its governed side to the classifier.
2. MUST prove the reader's behavior against porcelain output produced by a real
   rename in a repository the Task itself creates and removes, not against a
   hand-written snapshot pair. An injected pair is what allowed the previous
   rename test to pass while this path stayed broken.
3. MUST prove that a governed rename to an ungoverned name is refused when the
   changed set comes from the real reader.
4. MUST keep the existing real-repository staging behavior staging exactly the
   Agent's own changes.
5. MUST keep an ordinary removal settling with no operation required.
6. MUST NOT change the governed-mutation classifier's comparison, and MUST NOT
   weaken or delete an existing assertion.

## Subtasks

- [ ] Emit the source alongside the destination for rename and copy records.
- [ ] Prove the reader against a real rename.
- [ ] Prove the refusal with the changed set the real reader returns.
- [ ] Confirm staging and the ordinary path are unchanged.

## Acceptance Criteria

- [ ] A rename performed in a real repository yields both paths from the
      reader; today it yields only the destination.
- [ ] A governed rename to an ungoverned name is refused without the required
      operation, with the changed set taken from the real reader.
- [ ] An ordinary removal still settles with no operation required.
- [ ] The existing real-repository staging assertion still passes unchanged.

## Context

- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/task_engine.go`

## Verification

- `grep -q 'func TestSnapshotCarriesRenameSource' internal/daemon/daemon_test.go && go test -count=1 ./internal/daemon -run '^TestSnapshotCarriesRenameSource$'` — a real rename yields both paths; this fails today.
- `grep -q 'func TestGovernedRenameRefusesFromRealSnapshot' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestGovernedRenameRefusesFromRealSnapshot$'` — the refusal holds with the changed set the real reader returns.
- `grep -q 'func TestSnapshotCarriesRenameSource' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo$'` — real-repository staging is unchanged by the wider changed set.
- `grep -q 'func TestSnapshotCarriesRenameSource' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestOrdinaryRemovalDoesNotRequireOperation$'` — an ordinary removal still needs no operation.
- `grep -q 'func TestSnapshotCarriesRenameSource' internal/daemon/daemon_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestGovernedMutationDetectionUsesTheUnfilteredSnapshot$'` — detection still reads the unfiltered snapshot.

## References

- `_prd.md` → Goal 1; Core Features 1-2, 4; Declared intentional breaks 1.
- `_techspec.md` → Implementation Design: A rename names both of its paths;
  Testing Approach observations 1-3; Build Order 1.

## Result

The changed-path reader now appends a rename or copy record's source after its
destination. The governed-mutation classifier remains unchanged.

- `TestSnapshotCarriesRenameSource` creates a temporary Git repository,
  commits `Makefile`, performs `git mv Makefile notes.txt`, and reads both paths
  through `GitWorktreeSnapshotter`.
- `TestGovernedRenameRefusesFromRealSnapshot` renames a committed `Makefile`
  inside the Task-cycle repository and uses `GitWorktreeSnapshotter` to prove
  the existing missing-`implement` refusal prevents a commit.
- Before the reader change,
  `rtk sh -c 'GOCACHE=/tmp/roundfix-go-build-0135-task-01 go test -count=1 ./internal/daemon -run "^(TestSnapshotCarriesRenameSource|TestGovernedRenameRefusesFromRealSnapshot)$"'`
  failed because the real rename snapshot was `[notes.txt]`.
- After the reader change,
  `rtk sh -c 'GOCACHE=/tmp/roundfix-go-build-0135-task-01 go test -count=1 ./internal/daemon -run "^(TestSnapshotCarriesRenameSource|TestGovernedRenameRefusesFromRealSnapshot|TestSnapshotDiffCommitStagesOnlyAgentChangesInRealRepo|TestOrdinaryRemovalDoesNotRequireOperation|TestGovernedMutationDetectionUsesTheUnfilteredSnapshot)$"'`
  passed.
- `rtk sh -c 'GOCACHE=/tmp/roundfix-go-build-0135-task-01 go test -count=1 ./internal/daemon'`
  passed on the one unchanged retry after the sandbox denied GPG lock access
  during the first attempt.
- `rtk gofmt -d internal/daemon/daemon.go internal/daemon/daemon_test.go internal/daemon/task_engine_test.go`
  produced no output, and `rtk git diff --check` passed.
- `rtk sh -c 'GOCACHE=/tmp/roundfix-go-build-0135-task-01 make verify-incremental'`
  stopped in `fmt-check` because the unchanged files
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. This Task did
  not change those out-of-scope files.

Acceptance evidence:

- The real-repository snapshot test observes exactly `Makefile` and
  `notes.txt`, covering both sides of the rename.
- The real-snapshot Task-cycle test observes one refused Task, the existing
  `implement` refusal reason, and no commit.
- The unchanged ordinary-removal regression passed in the focused command.
- The unchanged real-repository staging regression passed in the focused
  command, and the unchanged unfiltered-snapshot regression confirms the
  existing classifier boundary still receives the reader's full result.
