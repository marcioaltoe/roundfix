---
task: task_01
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: completed
type: backend
complexity: low
---

# Task 01: A blocked reconcile honors the output contract

## Overview

A Profile-required skill can be absent at the revision selected for
`roundfix baseline skills reconcile`. In that case `buildSkillsReconcilePlan`
in `internal/baseline/skills_reconcile.go` returns a payload whose
`PlannedChanges` was never set, so `--format json` prints
`"plannedChanges": null`. The command exits `3`, which the help in
`internal/cli/cli.go` documents only for plan confirmation. An agent reading
the JSON or the help then looks for a Plan Digest that does not exist.

The payload is read by any agent or script driving the command, and the help,
the user guide and the Roundfix skill are read by the agents that follow them.

## Requirements

1. MUST build the required-removed payload in `buildSkillsReconcilePlan` with
   `PlannedChanges: []RestorePlannedChange{}`, so JSON prints `[]`. It MUST keep
   `PlanDigest` `nil` (JSON `null`), every `LockEdit` `nil`, `OK` false, and
   the `reconcile.required-removed` error in category `action_required`, so
   `baselineSkillsRestoreExit` still maps it to exit `3`.
2. MUST change the exit `3` line of the `baseline skills reconcile` help in
   `internal/cli/cli.go` to the text in the TechSpec, which contains the phrase
   `Profile-required skill is absent at the selected revision` and names
   `reconcile.required-removed`. The other exit lines stay unchanged.
3. MUST add `TestReconcileSkillsLockRequiredRemovedPayloadMarshalsAnEmptyPlan`
   to `internal/baseline/skills_reconcile_test.go`. It marshals the blocked
   payload and asserts `"plannedChanges":[]` and `"planDigest":null` in the
   bytes. `TestReconcileSkillsLockBlocksARequiredRemovedSkill` stays green and
   unchanged.
4. MUST add `TestBaselineSkillsReconcileRequiredRemovedPrintsAnEmptyPlan` to
   `internal/cli/baseline_skills_reconcile_test.go`. It drives the command
   through `RunContext` with `--format json` against an offline source that
   lacks a Profile-required skill, and asserts all of these:
   - exit `3`;
   - `plannedChanges` decodes as an empty JSON array, not `null`;
   - `planDigest` is `null`, `ok` is `false`, and `finding.code` is
     `reconcile.required-removed`;
   - `skills-lock.json` is byte-identical afterwards.
5. MUST add `TestBaselineSkillsReconcileHelpDocumentsTheRequiredRemovedExit` to
   the same file, asserting that the help contains `Profile-required skill is
   absent at the selected revision`.
6. MUST keep the non-blocked preview returning a Plan Digest with exit `3`:
   `TestBaselineSkillsReconcilePreviewThenConfirm` stays green and unchanged.
7. MUST state in the `baseline skills reconcile` paragraph of
   `docs/user-guide/commands.md`, and in the reconcile paragraph of
   `.agents/skills/roundfix/SKILL.md`, what a required skill absent at the
   revision does: it blocks with exit `3` and finding
   `reconcile.required-removed`, prints `plannedChanges: []` and
   `planDigest: null`, and writes nothing. Both MUST name
   `reconcile.required-removed`. MUST regenerate `skills/roundfix/SKILL.md`
   with `make skills-sync`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A blocked reconcile prints `plannedChanges: []` and `planDigest: null`,
      exits `3`, and writes nothing.
- [ ] The help, the user guide and the Roundfix skill document the
      required-removed exit.
- [ ] A non-blocked preview still returns its Plan Digest.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/skills_reconcile.go`
- interface: `internal/baseline/skills_reconcile_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/baseline_skills_reconcile_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileSkillsLockRequiredRemovedPayloadMarshalsAnEmptyPlan|TestReconcileSkillsLockBlocksARequiredRemovedSkill|TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries|TestBaselineSkillsReconcileRequiredRemovedPrintsAnEmptyPlan|TestBaselineSkillsReconcileHelpDocumentsTheRequiredRemovedExit|TestBaselineSkillsReconcilePreviewThenConfirm|TestBaselineSkillsReconcileHelpNamesTheConfirmationContract)$" ./internal/baseline ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileSkillsLockRequiredRemovedPayloadMarshalsAnEmptyPlan TestReconcileSkillsLockBlocksARequiredRemovedSkill TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries TestBaselineSkillsReconcileRequiredRemovedPrintsAnEmptyPlan TestBaselineSkillsReconcileHelpDocumentsTheRequiredRemovedExit TestBaselineSkillsReconcilePreviewThenConfirm TestBaselineSkillsReconcileHelpNamesTheConfirmationContract; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "reconcile.required-removed" docs/user-guide/commands.md && grep -q "reconcile.required-removed" .agents/skills/roundfix/SKILL.md && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task neither new test exists and neither guide names `reconcile.required-removed`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The blocked reconcile result
