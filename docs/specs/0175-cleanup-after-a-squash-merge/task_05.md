---
task: task_05
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: high
---

# Task 05: Own and sweep carry-forward staging worktrees

## Overview

`createCarryForwardStaging` in `internal/cli/carryforward.go` runs `git
worktree add` and `git worktree remove --force` directly through
`reconcileGitRaw`, outside the worktree administration lock Spec 0171 added.
It runs on every reconcile dry-run and implement preflight that inspects
carry-forward. A process killed during the add left
`$TMPDIR/roundfix-carry-forward-2309636395` `locked initializing` since
2026-09-24 14:57. One `--force` cannot remove a locked worktree, and nothing
sweeps `roundfix-carry-forward-*`. This Task moves staging into
`internal/worktree` with an owner record and sweeps stale staging worktrees.
The owner record is written by the process that creates the staging and read
by any later Roundfix process in the same repository. The sweep removes a
directory tree and a worktree registration. A concurrent carry-forward that
is still live must therefore never be swept, and staleness needs proof, never
age.

## Requirements

1. MUST add `internal/worktree/staging.go` with `AddCarryForwardStaging`,
   `CarryForwardStaging.Remove`, `InspectCarryForwardStaging` and
   `ReleaseCarryForwardStaging` as the TechSpec's "Staging worktrees" section
   defines:
   - `owner.json` holds the creator's PID and `store.OwnerProcessIdentity`.
   - Add and remove go through `runWorktreeCommand`.
   - A candidate is a registered worktree named `worktree` whose parent's name
     starts with `roundfix-carry-forward-`, wherever it lives.
2. MUST decide staleness only on proof. Either `owner.json` names a PID whose
   `store.OwnerProcessIdentity` fails or differs from the record, or there is
   no owner record and the registration is `locked` with reason
   `initializing`. Every other candidate is kept with a reason naming why.
3. MUST add `withWorktreeAdminLock` to `internal/worktree/adminlock.go`, and
   make `ReleaseCarryForwardStaging` hold the lock once while it re-proves
   staleness and runs `git worktree remove --force --force` directly through
   the runner. It MUST NOT call `runWorktreeCommand` while holding the lock,
   because that would wait on the lock it already holds. It then removes only
   the `roundfix-carry-forward-*` root.
4. MUST make `createCarryForwardStaging` take a `parent string` argument,
   delegate to `AddCarryForwardStaging`, and stop running `git worktree`
   itself. Production callers in `carryforward.go` and `reconcile.go` pass
   `""`, and `internal/cli/carryforward_test.go` passes `t.TempDir()`.
5. MUST add the `stagingCandidates` array (kind `staging`, worktree, owner
   PID, proof, action, refusal reason) and the `debrisSummary` fields
   `stagingCandidates` and `stagingApplied` to the `roundfix-reconcile/v1`
   report. Kept candidates go to `preservedCandidates`. The text report prints
   one `Staging candidate:` block per candidate and adds
   `staging-candidates=<n> staging-applied=<n>` to its debris summary line.
   Every existing field keeps its name and meaning.
6. MUST report stale staging in dry-run without removing it, release it under
   `--apply`, and release it under `--carry-forward` before that mode creates
   its own staging.
7. MUST add `stagingCandidates` to the expected top-level field list of
   `TestRunReconcileJSONMatchesTextFields` in `internal/cli/cli_test.go`, and
   change nothing else in that file.
8. MUST describe the staging sweep, its proof and the new report fields in
   `docs/user-guide/commands.md` (reconcile) and in the Run Worktree
   reconciliation section of `.agents/skills/roundfix/SKILL.md`, using the
   phrases `stagingCandidates` and `locked initializing`. Then regenerate
   `skills/roundfix/SKILL.md` with `make skills-sync`.
9. MUST put new tests in `internal/worktree/staging_test.go` and
   `internal/cli/reconcile_staging_test.go`, and MUST NOT remove or rename an
   existing test.

## Subtasks

- [ ] Move staging into `internal/worktree` with its owner record and lock.
- [ ] Report and sweep stale staging in the Reconcile Command.
- [ ] Update the reconcile guide and the Roundfix skill.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Staging records its owner, and adding it waits while another holder has the administration lock.
- [ ] A dead owner's staging and a legacy `locked initializing` registration are released, while a live owner's staging and an unlocked legacy registration are kept with a reason.
- [ ] Dry-run lists stale staging and removes nothing. `--apply` and `--carry-forward` release it, and the JSON report keeps every existing field.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `internal/worktree/staging.go`
- creates: `internal/worktree/staging_test.go`
- interface: `internal/worktree/adminlock.go`
- interface: `internal/cli/carryforward.go`
- interface: `internal/cli/carryforward_test.go`
- interface: `internal/cli/reconcile.go`
- creates: `internal/cli/reconcile_staging_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCarryForwardStagingRecordsItsOwner|TestCarryForwardStagingWaitsForTheAdminLock|TestCarryForwardStagingSweepReleasesADeadOwnersWorktree|TestCarryForwardStagingSweepKeepsALiveOwnersWorktree|TestCarryForwardStagingSweepReleasesALockedInitializingLegacyWorktree|TestCarryForwardStagingSweepKeepsAnUnlockedLegacyWorktree|TestReconcileDryRunReportsStaleStagingWithoutRemovingIt|TestReconcileApplyReleasesStaleStaging|TestReconcileCarryForwardReleasesStaleStagingFirst|TestReconcileKeepsALiveStagingAsPreserved|TestRunReconcileJSONMatchesTextFields|TestInspectSpecCarryForwards|TestCarryForwardStagingFailureIsNotClassifiedAsAConflict|TestCarryForwardOperationalCherryPickFailureIsNotAConflict|TestCarryForwardAcceptsAnUnresolvedRun|TestWorktreeAdministrationIsSerializedAcrossProcesses)$" ./internal/worktree ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCarryForwardStagingRecordsItsOwner TestCarryForwardStagingWaitsForTheAdminLock TestCarryForwardStagingSweepReleasesADeadOwnersWorktree TestCarryForwardStagingSweepKeepsALiveOwnersWorktree TestCarryForwardStagingSweepReleasesALockedInitializingLegacyWorktree TestCarryForwardStagingSweepKeepsAnUnlockedLegacyWorktree TestReconcileDryRunReportsStaleStagingWithoutRemovingIt TestReconcileApplyReleasesStaleStaging TestReconcileCarryForwardReleasesStaleStagingFirst TestReconcileKeepsALiveStagingAsPreserved TestRunReconcileJSONMatchesTextFields TestInspectSpecCarryForwards TestCarryForwardStagingFailureIsNotClassifiedAsAConflict TestCarryForwardOperationalCherryPickFailureIsNotAConflict TestCarryForwardAcceptsAnUnresolvedRun TestWorktreeAdministrationIsSerializedAcrossProcesses; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "stagingCandidates" docs/user-guide/commands.md && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "locked initializing" && grep -q "stagingCandidates" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide names `stagingCandidates`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Staging worktrees
- `_prd.md` → Core Feature 5; Success Metric 5
- `_techspec.md` → API Contracts 4 and 6

## Result

Implemented owned carry-forward staging in `internal/worktree`: creation writes
the current PID and process identity to `owner.json`, add/remove share the Git
worktree administration lock, inspection classifies only positive stale proof,
and release re-proves under one lock before running double-force removal and
deleting only the staging root. Carry-forward now delegates staging lifecycle
to that package and accepts a caller-selected parent for isolated tests.

Reconcile now reports `stagingCandidates`, includes staging counts in
`debrisSummary`, places kept staging in `preservedCandidates`, releases stale
staging under `--apply`, and releases it before `--carry-forward` creates its
own staging. The text report, user guide, canonical Roundfix skill, and synced
skill copy describe the same proof and report fields.

Focused checks:

- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^TestCarryForwardStaging' ./internal/worktree` — passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^TestWorktreeAdministration' ./internal/worktree` — passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^Test(ReconcileDryRunReportsStaleStagingWithoutRemovingIt|ReconcileApplyReleasesStaleStaging|ReconcileCarryForwardReleasesStaleStagingFirst|ReconcileKeepsALiveStagingAsPreserved|RunReconcileJSONMatchesTextFields|InspectSpecCarryForwards|CarryForwardStagingFailureIsNotClassifiedAsAConflict|CarryForwardOperationalCherryPickFailureIsNotAConflict|CarryForwardAcceptsAnUnresolvedRun)$' ./internal/cli` — passed.
- `make skills-sync-check` and `git diff --check` — passed.
- `make skills-sync` regenerated the distributed skill; `make baseline-digests`
  passed and reported no derived-artifact changes.

Acceptance evidence:

- Owner recording and lock waiting: `TestCarryForwardStagingRecordsItsOwner`
  exercised a real Git worktree and matched the record to
  `store.OwnerProcessIdentity`; `TestCarryForwardStagingWaitsForTheAdminLock`
  proved add did not reach Git while another holder owned the repository lock.
- Stale proof and negative cases:
  `TestCarryForwardStagingSweepReleasesADeadOwnersWorktree` and
  `TestCarryForwardStagingSweepReleasesALockedInitializingLegacyWorktree`
  released real registrations, while the separate live-owner and unlocked
  legacy tests preserved their registrations with refusal reasons.
- Reconcile behavior and schema: the four staging reconcile tests proved
  dry-run preservation, apply release, carry-forward release-before-use, and
  live-owner preservation. `TestRunReconcileJSONMatchesTextFields` passed with
  the additive `stagingCandidates` top-level field while retaining every
  existing field.

## Carry-forward provenance

- Source Run: `run_20260928T174529Z_c0cad0da5734a85f`
- Source commit: `4313558f97ea10aa1183f919276f0dff7b123ef9`
