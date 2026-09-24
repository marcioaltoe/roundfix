---
task: task_04
spec: 0163-baseline-decisions-and-regeneration
status: pending
type: backend
complexity: high
---

# Task 04: Reconcile the lock on proven absence

## Overview

`baseline skills restore` can restore a `skills-lock.json` entry but never retire one upstream removed, so an obsolete entry is left in place or removed by a hand edit on unrecorded evidence.

## Requirements

1. MUST add `baseline.ReconcileSkillsLock` in `internal/baseline/skills_reconcile.go`, taking a repository, a built-in Profile, a source repository, a 40-hex commit, an optional source directory and an optional confirmation digest.
2. MUST acquire the commit through `acquireRestoreGroup`'s fetch and commit check, and read the lock with the ordered JSON of `skills_restore_lock.go`.
3. MUST classify each entry of the selected source as `present`, `moved`, `obsolete` or `required-removed` as the TechSpec defines, and every other entry as `unrelated`.
4. MUST plan lock removals only for `obsolete` entries, report each one's installed tree as retained, and keep unknown lock fields and entry order.
5. MUST block the whole plan on any `required-removed` entry, and write nothing on a failed fetch, a commit mismatch or a mutable ref.
6. MUST return a plan digest from the preview, and apply only with that digest, through the restore transaction and its preimage check.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Against a local bare repository, only an entry absent at the commit and not required by the Profile is removed.
- [ ] A required skill absent at the commit blocks and nothing is written.
- [ ] An unreachable source removes nothing.
- [ ] Moved and unrelated entries are kept.
- [ ] Installed trees and unknown lock fields survive apply.
- [ ] A stale preimage writes nothing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_restore_git.go`
- interface: `internal/baseline/skills_restore_lock.go`
- creates: `internal/baseline/skills_reconcile.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries|TestReconcileSkillsLockBlocksARequiredRemovedSkill|TestReconcileSkillsLockUnreachableSourceRemovesNothing|TestReconcileSkillsLockKeepsMovedAndUnrelatedEntries|TestReconcileSkillsLockPreservesInstalledTreesAndUnknownFields|TestReconcileSkillsLockStalePreimageWritesNothing)$" ./internal/baseline 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries TestReconcileSkillsLockBlocksARequiredRemovedSkill TestReconcileSkillsLockUnreachableSourceRemovesNothing TestReconcileSkillsLockKeepsMovedAndUnrelatedEntries TestReconcileSkillsLockPreservesInstalledTreesAndUnknownFields TestReconcileSkillsLockStalePreimageWritesNothing; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Lock reconciliation
