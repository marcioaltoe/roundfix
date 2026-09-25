---
task: task_05
spec: 0163-baseline-decisions-and-regeneration
status: completed
type: backend
complexity: medium
---

# Task 05: The lock reconciliation command

## Overview

`ReconcileSkillsLock` needs a public, confirmation-gated command, and `baseline skills` today accepts only `restore`.

## Requirements

1. MUST add `roundfix baseline skills reconcile --profile <id> --source <owner/repo> --revision <commit> [--source-dir <path>] [--confirm-plan <digest>] [--repo <path>] [--format <text|json>]` in `internal/cli/baseline_skills_reconcile.go`, dispatched from `internal/cli/baseline_profile.go`.
2. MUST use the exit categories and output shape of `baseline skills restore`: a non-empty preview exits 3 with its plan digest, an empty reconciliation exits 0, and a revision that is not a 40-hex commit is refused as invalid before any fetch.
3. MUST add the command's usage to `commandUsage` and document it in `docs/user-guide/commands.md`.
4. MUST leave Doctor unchanged: it reports an obsolete lock entry without editing the lock.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A preview exits 3 with a digest; confirming that digest removes the obsolete entry.
- [ ] A mutable revision is refused before any fetch.
- [ ] The help names the confirmation contract.
- [ ] Doctor leaves an obsolete lock entry untouched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/baseline_skills_restore.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/cli/doctor.go`
- creates: `internal/cli/baseline_skills_reconcile.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBaselineSkillsReconcilePreviewThenConfirm|TestBaselineSkillsReconcileRejectsAMutableRevision|TestBaselineSkillsReconcileHelpNamesTheConfirmationContract|TestDoctorLeavesAnObsoleteLockEntryUntouched)$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestBaselineSkillsReconcilePreviewThenConfirm TestBaselineSkillsReconcileRejectsAMutableRevision TestBaselineSkillsReconcileHelpNamesTheConfirmationContract TestDoctorLeavesAnObsoleteLockEntryUntouched; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Lock reconciliation

## Result

### Implementation

- `roundfix baseline skills reconcile` now accepts the selected built-in
  Profile, source repository, exact revision, optional offline source,
  confirmation digest, repository and text/JSON format. `baseline skills`
  dispatches both `restore` and `reconcile`.
- The command delegates classification and mutation to `ReconcileSkillsLock`,
  uses the restoration command's exit categories, and renders the same
  blocked/applied/no-changes, finding, Plan Digest and planned-change shape.
  A changed preview exits `3`; confirmed apply and an empty reconciliation exit
  `0`.
- Root and Baseline help plus `docs/user-guide/commands.md` document the command,
  its immutable 40-hex revision requirement and its exact Plan Digest
  confirmation contract.
- Doctor production code is unchanged. A regression test runs its real
  Repository Skill Set check with an extra, no-longer-required lock entry and
  proves the lock bytes remain identical.

### Focused checks

- Before implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestBaselineSkillsReconcilePreviewThenConfirm$' ./internal/cli`
  reached public dispatch and failed with exit `2` because `baseline skills`
  accepted only `restore`. The fixture first needed repository-local
  `commit.gpgsign=false` so inherited workstation signing did not mask that
  product signal.
- After implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestBaselineSkillsReconcile' ./internal/cli`
  passed the preview/confirm/no-op, mutable-revision and help cases.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 -run '^TestDoctorLeavesAnObsoleteLockEntryUntouched$' ./internal/cli`
  passed with the real repository skill checker.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -count=1 ./internal/cli`
  passed after the existing GitHub-backed test boundary received network
  access; the first sandboxed attempt was blocked from `api.github.com`.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache make verify-incremental`
  passed, including formatting, vet, the repository test packages, skill mirror
  checks and the build. `rtk git diff --check` also passed before this Result
  update.

### Acceptance evidence

1. `TestBaselineSkillsReconcilePreviewThenConfirm` drives the public command
   through a real local Git source: preview exits `3` with one digest-bound lock
   removal, confirmation removes the obsolete entry, and a second unconfirmed
   run exits `0` with no changes.
2. `TestBaselineSkillsReconcileRejectsAMutableRevision` passes `main` together
   with a deliberately missing offline source and receives
   `reconcile.commit-invalid`, proving refusal occurs before source acquisition.
3. `TestBaselineSkillsReconcileHelpNamesTheConfirmationContract` proves public
   help names the full command, exit-`3` preview and exact `--confirm-plan`
   contract.
4. `TestDoctorLeavesAnObsoleteLockEntryUntouched` proves Doctor reports
   Repository Skill Set readiness and leaves the lock, including the obsolete
   entry, byte-identical.

The Task's declared `## Verification` command was not run; Verification remains
Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20260925T135958Z_cb4c08f055bb883c`
- Source commit: `22b642d8798cc6799f315b716823927f1fad67ad`
