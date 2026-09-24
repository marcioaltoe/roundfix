---
task: task_05
spec: 0163-baseline-decisions-and-regeneration
status: pending
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
