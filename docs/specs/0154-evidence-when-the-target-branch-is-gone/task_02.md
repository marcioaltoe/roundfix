---
task: task_02
spec: 0154-evidence-when-the-target-branch-is-gone
status: pending
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
