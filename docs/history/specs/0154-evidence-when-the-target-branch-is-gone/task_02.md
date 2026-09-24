---
task: task_02
spec: 0154-evidence-when-the-target-branch-is-gone
status: completed
type: backend
complexity: low
---

# Task 02: Preserved reasons that name the missing proof

## Overview

`target branch "<name>" is absent` is true and unhelpful: it says what was not
found, never what was sought. Once Task 01 makes the absent target send the
check elsewhere, that string stops being a terminal answer and a reader needs to
know which proof failed.

## Requirements

1. MUST distinguish, in the preserved reason, an absent target whose default
   branch carries no evidence for this Spec from one whose default branch could
   not be resolved.
2. MUST name the Spec the evidence was sought for.
3. MUST keep the existing reasons for a present target branch byte-identical.
4. MUST stay within the reason's existing length bound.

## Subtasks

- [ ] Add the reasons for each missing-proof case.
- [ ] Add tests asserting each reason.

## Acceptance Criteria

- [ ] Each preserved case produces its own reason, naming the proof sought.
- [ ] Present-target reasons are unchanged.
- [ ] No reason exceeds the existing maximum length.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReconcileReasonNamesTheMissingProof" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcileReasonNamesTheMissingProof"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — What still preserves

## Result

### Implementation

- The absent-target classifier now chooses one preserved reason when the
  default branch resolves without superseding QA Report evidence and another
  when the default branch cannot be resolved.
- Both reasons put the Run's Spec slug before the missing-proof detail and pass
  through the existing 160-byte reason bound.
- The present-target classification block and its reason strings were not
  changed.
- The Task 01 test entry point remains in place and shares the same fixture body
  with the new Task 02 reason test.

### Acceptance evidence

- `TestReconcileReasonNamesTheMissingProof/default_branch_has_no_superseding_evidence`
  asserted the exact preserved reason, including the Spec slug and the missing
  superseding QA Report on the default branch.
- `TestReconcileReasonNamesTheMissingProof/default_branch_cannot_be_resolved`
  asserted a distinct exact reason, including the same Spec slug and that the
  superseding QA Report could not be sought.
- Both cases asserted that their reason was no longer than
  `reconciliationReasonMaxBytes`; the production helpers also retain
  `boundedReconciliationReason` as the final constructor.
- The production diff from Task 01 changes only the absent-target block and its
  former absent-target reason helper. `rtk make verify-incremental` also passed
  the existing present-target package tests.

### Focused checks

- Initial regression:
  `rtk env GOCACHE="$PWD/.cache/go-build" go test ./internal/worktree -run '^TestReconcileReasonNamesTheMissingProof$/default_branch_has_no_superseding_evidence$'`
  — failed because the preserved reason was still `target branch
  "ma/deleted-target" is absent`.
- After implementation:
  `rtk env GOCACHE=/private/tmp/roundfix-task02-focused.EYAzBg go test ./internal/worktree -run 'Test(Reconcile(PreservesWithoutDefaultBranchEvidence|ReasonNamesTheMissingProof)|ClassifyRunBranchSetPreservesAbsentTarget)$'`
  — passed both Task test entry points, both missing-proof cases and the
  existing absent-target compatibility test.
- Repository incremental gate:
  `rtk env GOCACHE=/private/tmp/roundfix-task02-focused.EYAzBg make verify-incremental`
  — passed vet, all Go tests, skill checks and the build outside the sandbox.
- An earlier incremental attempt with `GOCACHE` under the repository failed the
  suite guard because the cache itself changed repository paths; the cache was
  removed and the passing rerun used `/private/tmp`.
- `rtk git diff --check` passed before this Result was appended.
- The command under `## Verification` was not run; the Daemon owns it.
