---
task: task_02
spec: 0134-a-governed-deletion-the-gate-can-see
status: completed
type: backend
complexity: medium
---

# Task 02: Classify a removed Governed Path as a mutation

## Overview

The classifier builds a set from the paths before the Agent turn and reports a
mutation only for a path in the after list that is absent from it, so a removal
is invisible and a rename reads as an ungoverned addition. This slice compares
both directions and shares one classifier with the public Settle path. It is
verifiable alone: a Task removing a Governed Path without the required operation
refuses, and every refusal that works today keeps its token.

## Requirements

1. MUST classify a Governed Path present before the turn and absent after it as
   a governed mutation, so the required operation is asked for.
2. MUST classify a rename of a Governed Path from its source as well as its
   destination, which follows from comparing both directions.
3. MUST compare path snapshots rather than parse porcelain rename records: the
   snapshot pair already carries both sides, and a second source of truth is
   how this defect returns.
4. MUST have the Daemon path and the public Settle path consume one classifier,
   so the two cannot reach different verdicts for the same change.
5. MUST keep every refusal that works today working with its existing token,
   including a governed change outside the bounded set, a grant edited in the
   commit that consumes it, and a Spec that declares no authorization.
6. MUST update only the characterization rows this Task intentionally moves.

## Subtasks

- [ ] Compare both directions in the classifier.
- [ ] Share the classifier with the Settle path.
- [ ] Prove a removal and a rename refuse without the operation.
- [ ] Move the two recorded characterization rows and no others.

## Acceptance Criteria

- [ ] A Task removing a Governed Path without the `commit` operation refuses;
      it settles silently today.
- [ ] A Task renaming a Governed Path to an ungoverned name refuses, classified
      from its source.
- [ ] The public Settle path reaches the same verdict as the Daemon path for the
      same removal, exercised through the public command.
- [ ] A Task removing an ordinary source path still settles, so the change adds
      cases rather than refusing everything.
- [ ] Every existing refusal keeps its token and its condition.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/settle.go`

## Verification

- `grep -q 'func TestGovernedRemovalRequiresOperation' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestGovernedRemovalRequiresOperation$'` — a removed Governed Path refuses without the operation; this fails today.
- `grep -q 'func TestGovernedRenameClassifiesFromSource' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestGovernedRenameClassifiesFromSource$'` — a rename is classified from its source.
- `grep -q 'func TestSettleClassifiesGovernedRemoval' internal/cli/settle_test.go && go test -count=1 ./internal/cli -run '^TestSettleClassifiesGovernedRemoval$'` — the public Settle path reaches the same verdict through its own call site.
- `grep -q 'func TestGovernedRemovalRequiresOperation' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon -run '^TestGovernedMutationClassificationCharacterization$'` — the recorded answers hold with exactly the two intended rows moved.
- `grep -q 'func TestGovernedRemovalRequiresOperation' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/daemon ./internal/cli` — both packages pass with one shared classifier.

## References

- `_prd.md` → Goals 1; Core Features 1-2, 4; Declared intentional breaks 1.
- `_techspec.md` → Implementation Design: Symmetric classification; Build Order 2.

## Result

- Implemented one exported snapshot classifier that compares both path sets.
  The Daemon's implementation and QA commit paths call it with their before and
  after snapshots, and Settle calls it with the current changed-path snapshot.
  The classifier does not inspect porcelain rename records.
- Added Task-cycle regressions for a governed removal and a rename from a
  governed source with `implement` authority but no `commit` authority, plus an
  ordinary-source removal companion. Added a public Settle regression that
  deletes a tracked `Makefile` and observes the preflight refusal before
  Verification or mutation.
- Pre-change signal: `rtk go test -count=1 ./internal/daemon ./internal/cli -run '^(TestGovernedRemovalRequiresOperation|TestGovernedRenameClassifiesFromSource|TestOrdinaryRemovalDoesNotRequireOperation|TestSettleClassifiesGovernedRemoval)$'`
  was initially blocked by `operation not permitted` in the host Go build
  cache. One unchanged retry with cache access ran the tests: the ordinary and
  Settle cases passed, while the governed removal and rename each failed
  because the Task settled with `Completed:1 Failed:0`.
- Focused check after the implementation: the same four-test command exited 0
  with 4 tests passing across `internal/daemon` and `internal/cli`.
- Acceptance criterion 1: `TestGovernedRemovalRequiresOperation` now observes
  one failed Task, the existing `commit` operation refusal text and the
  `_authorization.md` path, with no commit written.
- Acceptance criterion 2: `TestGovernedRenameClassifiesFromSource` supplies
  `Makefile` only in the before snapshot and `internal/ordinary.go` only in the
  after snapshot; it observes the same `commit` refusal with no commit written.
- Acceptance criterion 3: `TestSettleClassifiesGovernedRemoval` exercises the
  public command against a removed tracked `Makefile`. It observes exit 2 and
  the existing `commit` refusal while stdout, Task content, HEAD and the dirty
  deletion remain unchanged, and its Verification command does not execute.
- Acceptance criterion 4: `TestOrdinaryRemovalDoesNotRequireOperation`
  supplies `internal/ordinary.go` only in the before snapshot with no
  authorization record; it observes one settled Task and one Task commit.
- Acceptance criterion 5: `rtk go test -count=1 ./internal/daemon ./internal/cli ./internal/speccheck -run '^(TestGovernedMutationClassificationCharacterization|TestGovernedMutationDetectionUsesTheUnfilteredSnapshot|TestGovernedMutationRefusesMissingImplementAuthority|TestGovernedRemovalRequiresOperation|TestGovernedRenameClassifiesFromSource|TestOrdinaryRemovalDoesNotRequireOperation|TestSettleClassifiesGovernedRemoval|TestSettleRefusesMissingCommitAuthority|TestSettleCommitsOrdinaryWorkWithoutRecord|TestGovernedChangeStillRefusesWithoutRecord|TestAuditRefusesOutOfGrantUnderBothCitationForms|TestAuditRefusesSelfApprovalAndRetroactiveGrants)$'`
  exited 0 with 21 tests and subtests passing across three packages. This keeps
  the `QA-AUTH-PATHS` refusals for an out-of-grant path and a same-commit grant
  edit, and the `authorization operation "implement" is not permitted`
  refusal for a Spec without an authorization record.
- The characterization's addition row remains unchanged. Only the removal and
  governed-source rename rows moved from `false` to `true` and dropped their
  Task 02 transition notes.
- The Daemon-owned `## Verification` commands were not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20260913T221055Z_c70caa90b7bbad59`
- Source commit: `c589e900dcc020edb62634db2aa5029aebde4e35`
