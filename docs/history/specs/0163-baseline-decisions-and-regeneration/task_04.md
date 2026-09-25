---
task: task_04
spec: 0163-baseline-decisions-and-regeneration
status: completed
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

## Result

### Implementation

- `ReconcileSkillsLock` now loads a built-in Profile, requires an exact 40-hex
  commit, and classifies the selected source's ordered lock entries as
  `present`, `moved`, `obsolete`, or `required-removed`; entries from other
  sources are `unrelated`.
- Restore and reconciliation share one immutable acquisition path for the Git
  fetch and `FETCH_HEAD^{commit}` identity check. Failed acquisition returns no
  classifications or plan.
- Reconciliation plans only `remove-lock-entry` changes for `obsolete`
  entries. It reports every installed tree as `retained`, preserves ordered JSON
  and unknown surviving fields, and blocks the complete plan when a required
  Profile skill is absent upstream.
- Preview returns a content-bound plan digest. Apply requires that exact digest
  and writes only `skills-lock.json` through the restore transaction and its
  preimage revalidation; no installed skill tree enters the transaction.

### Focused checks

- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run '^TestReconcileSkillsLockRejectsMutableRevisionBeforeAcquisition$' ./internal/baseline`
  failed to compile because the reconciliation request, payload, dispositions,
  and engine did not exist.
- After implementation,
  `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run '^TestReconcileSkillsLock' ./internal/baseline`
  passed all seven reconciliation cases, including the six acceptance cases and
  the immutable-revision negative case.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache make verify-changed` passed
  after process-table access was granted. The first sandboxed attempt reached
  the existing force-stop integration tests but could not enumerate the owner
  process tree (`operation not permitted`); the permission-enabled rerun and
  the final post-edit rerun both exited 0.

### Acceptance evidence

1. `TestReconcileSkillsLockRemovesOnlyAbsentUnrequiredEntries` uses the exported
   API and a local bare source repository; preview plans only the absent,
   unrequired entry, and confirmed apply keeps the present entry.
2. `TestReconcileSkillsLockBlocksARequiredRemovedSkill` resolves the real
   `go-cli-tui` Profile, classifies its absent `coding-guidelines` skill as
   `required-removed`, emits no applicable plan or digest, and preserves the
   repository bytes.
3. `TestReconcileSkillsLockUnreachableSourceRemovesNothing` proves a failed Git
   fetch returns no classifications, changes, or digest and leaves the visible
   tree unchanged.
4. `TestReconcileSkillsLockKeepsMovedAndUnrelatedEntries` proves a same-named
   `SKILL.md` at a new path is `moved`, another source is `unrelated`, and the
   empty reconciliation writes nothing.
5. `TestReconcileSkillsLockPreservesInstalledTreesAndUnknownFields` confirms
   apply retains the obsolete skill's installed tree, unknown root and retained
   entry fields, and the relative order of every surviving lock entry.
6. `TestReconcileSkillsLockStalePreimageWritesNothing` changes the lock after
   transaction staging; preimage revalidation returns
   `plan.confirmation.stale` and preserves the concurrent bytes.

`TestReconcileSkillsLockRejectsMutableRevisionBeforeAcquisition` additionally
proves a mutable ref is rejected before the deliberately unreachable source
path is inspected. The Task's declared `## Verification` command was not run;
Verification remains Daemon-owned.

## Carry-forward provenance

- Source Run: `run_20260925T135958Z_cb4c08f055bb883c`
- Source commit: `fef89f30576c678a66e7a03fd30dbd70b605ef25`
