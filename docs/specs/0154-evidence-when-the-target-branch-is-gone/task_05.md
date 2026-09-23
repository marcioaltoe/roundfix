---
task: task_05
spec: 0154-evidence-when-the-target-branch-is-gone
status: pending
type: backend
complexity: medium
---

# Task 05: Archive-state evidence for an absent target

## Overview

Pre-PR review found that the new fallback accepts evidence this Spec's own Core
Feature 3 says must preserve. Confirmed in source: `qaReportDirectories`
returns both `docs/specs/<slug>/qa` and the archived path, so
`supersedingQAReport` proves delivery from a QA Report under the *active* Spec
tree.

For a present target branch that is the accepted, pre-existing behavior and this
Spec must not change it. For the absent-target fallback it means `--apply` can
remove a Run's retained surfaces while its Spec is still active — which Core
Feature 3 forbids: "an unarchived Spec ... leaves it preserved".

The distinction matters because the two cases carry different risk. With the
target present, the branch itself is the evidence trail. With the target gone,
the archive is the only durable record that the work concluded.

## Requirements

1. MUST require archive-state evidence in the absent-target fallback: a QA
   Report under the archived Spec path proves delivery, one under the active
   Spec path does not.
2. MUST preserve the Run, naming the missing proof, when only active-path
   evidence exists.
3. MUST leave the present-target path accepting both paths, exactly as it does
   today.
4. MUST keep every other behavior Tasks 01 through 03 delivered.

## Subtasks

- [ ] Restrict the fallback's evidence to the archived Spec path.
- [ ] Add the preserved reason for active-only evidence.
- [ ] Add tests for archived evidence, active-only evidence, and the unchanged
      present-target path.

## Acceptance Criteria

- [ ] An absent target with archived-path evidence releases.
- [ ] An absent target with only active-path evidence preserves, and the reason
      says the Spec is not archived.
- [ ] A present target releases on either path, unchanged.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReconcileFallbackRequiresArchivedEvidence" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcileFallbackRequiresArchivedEvidence"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReconcilePresentTargetAcceptsEitherSpecPath" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReconcilePresentTargetAcceptsEitherSpecPath"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Core Features
