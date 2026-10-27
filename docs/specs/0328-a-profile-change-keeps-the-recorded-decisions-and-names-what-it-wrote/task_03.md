---
task: task_03
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
status: pending
type: backend
complexity: medium
---

# Task 03: The apply names each Root Backup and reconcile says what it compared

## Overview

A plan that backs up a root instruction carrier labels that file change as a
Root Backup in its text, and its apply carries one warning per Root Backup in
text and JSON, so the adopter knows to commit it. A `baseline skills
reconcile` with nothing to remove explains what it removes and names the
Setup Snapshot revision the installed skills follow. This answers frictions
(c) and (d) of the adopted Backlog Entry (Quaestor, 2026-10-26).

## Requirements

1. MUST implement API Contract 4: `printBaselinePlanText` appends `; root
   backup of <carrier>` inside the parentheses of each file change whose
   managed entries hold a `backup:<carrier>` entry, and `ApplyPlan` appends,
   after the plan's warnings, one `baseline.preservation.backup.written`
   warning per `backup:` managed entry with the message the contract fixes.
   The plan document and the result schema keep their fields.
2. MUST keep the managed refresh backup-free and warning-free (Invariant 3):
   a `baseline update` that applies prints no
   `baseline.preservation.backup.written` warning.
3. MUST implement API Contract 5: export `SetupSnapshotRevisions` from
   `internal/baseline/skills_reconcile.go`, returning the sorted distinct
   Setup Snapshot refs of the profile's external skills from the given
   source, and print the three no-change lines only when the payload is OK,
   not applied, and has no planned change and no finding. Reconcile's JSON
   and its planning stay as they are (Invariant 4).
4. MUST add the tests `_techspec.md` → Testing Approach names for task_03:
   `TestBaselineApplyNamesTheRootBackupItWrote` (adopt with a committed
   `AGENTS.md` under `preservation.mode=greenfield` through `baseline plan`
   and `baseline apply`, text and JSON, asserting the backup path equals the
   created file and the carrier is `AGENTS.md`),
   `TestBaselinePlanTextLabelsTheRootBackup`,
   `TestBaselineUpdateWritesNoRootBackupWarning` and
   `TestBaselineSkillsReconcileNoChangesNamesTheSetupSnapshotRevision`,
   whose expected ref is read from the embedded catalog through the
   exported function or the snapshot asset, never written as a literal.
   Each fails when its behavior is removed.

## Subtasks

- [ ] Label backup file changes in the plan text.
- [ ] Add the backup warnings to the apply result.
- [ ] Add the snapshot revisions and the reconcile no-change lines.
- [ ] Add the named tests.

## Acceptance Criteria

- [ ] An apply that writes `AGENTS.<digest>.md` names it in a `Warning:` line and in JSON `warnings`.
- [ ] A managed refresh apply carries no backup warning.
- [ ] A no-change reconcile prints the Setup Snapshot revision the embedded catalog pins for the profile's skills from that source, and its JSON is unchanged.

## Context

- instruction: `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- instruction: `docs/user-guide/commands/baseline.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`
- instruction: `internal/baseline/preservation.go`
- instruction: `internal/cli/baseline_apply.go`
- interface: `internal/baseline/apply.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/baseline/skills_reconcile.go`
- interface: `internal/cli/baseline_skills_reconcile.go`
- creates: `internal/cli/baseline_root_backup_test.go`
- interface: `internal/cli/baseline_skills_reconcile_test.go`
- interface: `internal/cli/baseline_apply_test.go`
- interface: `internal/baseline/apply_test.go`

## Consumers

- The apply `Result` warnings and the plan text — consumers: `./internal/baseline`, `./internal/cli`
- The reconcile text — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestBaselineApply.*|TestBaselinePlanTextLabelsTheRootBackup|TestBaselineUpdateWritesNoRootBackupWarning|TestBaselineSkillsReconcile.*|TestReconcileSkillsLock.*|TestApply.*)$' ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineApplyNamesTheRootBackupItWrote TestBaselinePlanTextLabelsTheRootBackup TestBaselineUpdateWritesNoRootBackupWarning TestBaselineSkillsReconcileNoChangesNamesTheSetupSnapshotRevision TestBaselineApplyCommand TestBaselineSkillsReconcilePreviewThenConfirm; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the four new tests do not exist, so the first pass check fails. Afterwards they pass with the existing apply and reconcile tests that consume the changed output.

## References

- `_prd.md` → Goals; Core Feature 4; Core Feature 5; User Experience; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior 4-5; Interfaces; API Contract 4; API Contract 5; Invariants 3-4; Testing Approach; Build Order 2
- ADR-0264; ADR-0341
