---
task: task_06
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: completed
type: backend
complexity: medium
---

# Task 06: Apply releases a merged Run whose target branch is gone

## Overview

Corrective Task for QA finding F1 of the 2026-10-02 QA Report (row 5, Requirement 5). After a squash merge deletes the target branch, `roundfix reconcile` classifies the Run superseded through the merged-head proof, but `reconcile --apply` does not release it. `applyTerminalRunWithStore` in `internal/worktree/worktree.go` passes the empty `TargetHead` of the absent branch into the Store's `IntegrationReconciliation`, and `validateIntegrationReconciliation` in `internal/store/store.go` refuses it with `target head is required`. The proof head that justified the classification never reaches the Store.

## Requirements

1. MUST, when a terminal Run is classified superseded by a merged-head proof while its target branch is absent, send the Store the proof it was classified on: the merged head as `TargetHead`, and as `TargetBranch` the branch or ref the merged head was read from. A Run whose target branch is present keeps its current request unchanged.
2. MUST NOT relax `validateIntegrationReconciliation` for any other request: an empty target head without a merged-head proof is still refused.
3. MUST add the tests named in Verification to the new file `internal/worktree/merged_spec_absent_target_test.go`, over real Git repositories and the real Store, as `internal/worktree/worktree_test.go` already does:
   - `TestApplyReleasesAMergedRunWhoseTargetBranchIsGone`: dry-run reports superseded, and apply removes the Run Worktree and the Run Branch and records the reconciliation with the merged head.
   - `TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone`: a Run with an undeclared leftover stays preserved, and nothing is removed.
   - `TestApplyStillRefusesAnEmptyTargetHeadWithoutProof`: an empty target head without a merged-head proof is still refused.
4. MUST leave every existing `internal/worktree` and `internal/store` test unedited and passing.
5. MUST prove each new test fails before the fix and passes after, and record that in the Result.

## Subtasks

- [ ] Carry the merged-head proof into the Store request for an absent target.
- [ ] Add the three tests over real repositories and the real Store.

## Acceptance Criteria

- [ ] `reconcile --apply` releases a superseded Run whose target branch is gone.
- [ ] An undeclared leftover is still preserved, and an empty target head without proof is still refused.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/merged_head.go`
- interface: `internal/store/store.go`
- creates: `internal/worktree/merged_spec_absent_target_test.go`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/worktree/merged_spec_leftovers_test.go`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestApplyReleasesAMergedRunWhoseTargetBranchIsGone|TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone|TestApplyStillRefusesAnEmptyTargetHeadWithoutProof)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestApplyReleasesAMergedRunWhoseTargetBranchIsGone TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone TestApplyStillRefusesAnEmptyTargetHeadWithoutProof; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0
- `out="$(go test -count=1 ./internal/worktree ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; test -f internal/worktree/merged_spec_absent_target_test.go` — expected: exit 0

## References

- QA Report 2026-10-02 → F1; row 5
- task_04 → Requirement 5
- ADR-0212; ADR-0053; ADR-0161

## Result

Implemented the Task 06 slice for QA finding F1 (2026-10-02, row 5).
Merged-head evidence now retains its source ref alongside its head and
revalidates both before apply. When the recorded target branch is absent,
apply sends that proof head and source ref to the Store. The default-branch
fallback records its branch name; an explicit merge record records the
immutable commit ref used to read the proof.

The Store request carries the original recorded target separately for its
existing Run metadata check. The stored Run's target is unchanged. A present
target keeps the original request branch and head. No changes were made to
`validateIntegrationReconciliation`, existing tests, or the database schema.

Acceptance evidence from the new real-Git/real-Store tests:

- Release a superseded Run whose target is gone:
  `TestApplyReleasesAMergedRunWhoseTargetBranchIsGone` checks read-only
  inspection reports `superseded`, then apply removes both the Worktree and
  Run Branch and persists the proof head/ref in the reconciliation event.
  Covers clean Runs and declared leftovers, plus a present-target control
  that checks the original request branch/head remain in the event.
- Preserve undeclared leftovers:
  `TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone` checks the Run stays
  dirty, apply refuses, the file contents and Git resources remain, and no
  reconciliation event or state transition occurs. Removing only the
  undeclared file permits the remaining declared leftovers to be released.
- Refuse an empty target head without proof:
  `TestApplyStillRefusesAnEmptyTargetHeadWithoutProof` checks the real Store
  returns `target head is required` without recording an event or changing
  the Run. Its positive control applies an inspected explicit merged-head
  proof and records the immutable commit ref.

Focused red/green check:

`GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test -count=1 -v -run '^TestApply(ReleasesAMergedRun|KeepsAnUndeclaredLeftover|StillRefusesAnEmptyTargetHead)' ./internal/worktree`

Before the production fix, exit 1 with all three named tests failing:
the release test's clean case recorded `main` rather than proof source
`trunk`, and its declared-leftover case refused `target head is required`;
the preservation test's release control and the empty-head test's inspected
proof control also refused `target head is required`. Their negative checks
already held. After the fix, exit 0 with all three named tests and the
present-target control passing.

`GOCACHE=/private/tmp/roundfix-task06-gocache rtk make verify-incremental`:
the first run passed `internal/worktree` and `internal/store` but exited 2
because two unrelated CLI force-stop tests could not read the process table
inside the sandbox (`operation not permitted`). Rerunning with process-table
access exited 0: formatting, vet, the package suites, skill checks, and build
passed; successful package results were reused from the first run's cache.

Task status and authored Verification remain Daemon-owned. No declared
Verification command was run, and no commit, push, or Pull Request was made.
No follow-up work was added to this diff.

### Verification feedback — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T114737Z_50dd817fd855fe86/verification/batch-001-attempt-1.log`
and the related CLI detach test fixtures. The recorded `make verify-changed`
attempt passed `internal/worktree` and `internal/store`; its exit 2 came from
`TestImplementDetachChildEndsWhenItsTestBinaryDies` and
`TestRunImplementDetachSurvivesCallerProcessGroupKill`, which reported
`operation not permitted` while probing or terminating fixture process
groups. Those fixtures use `syscall.Kill` on their recorded process groups;
they are outside the Task 06 reconciliation slice.

Focused diagnostic check, with elevated process-control access:

`GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test -count=1 -v -run '^(TestImplementDetachChildEndsWhenItsTestBinaryDies|TestRunImplementDetachSurvivesCallerProcessGroupKill)$' ./internal/cli`

Exit 1: the survival test passed, while the teardown test still reported
`operation not permitted` for its fixture process-group probe and cleanup.
Elevated access did not resolve that failure. Its underlying process-control
cause remains unresolved; no CLI code, tests, or Verification configuration
were changed to bypass it.

Reran the focused Task 06 check recorded above: exit 0, all three named tests
and the present-target control passed. No Task 06 implementation repair was
indicated by these diagnostics. Only this Result addendum was changed in the
feedback turn. Follow-up outside this slice: investigate the CLI fixture
process-group refusal if it recurs in the Daemon's configured Verification.
Task status remains unchanged, and the configured Verification rerun remains
with the Daemon; no declared Verification command was rerun here.
