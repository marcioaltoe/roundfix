---
task: task_05
spec: 0154-evidence-when-the-target-branch-is-gone
status: completed
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

## Result

### Implementation

- The absent-target branch-set fallback now releases a Run Branch only when
  the shared supersession proof identifies a QA Report under the archived Spec
  path.
- Active-path-only evidence now preserves the Run Branch and names both the
  missing archived-path proof and that the Spec is not archived.
- Present-target reconciliation still uses the unchanged shared supersession
  rule, so active and archived Spec paths remain accepted.

### Focused checks

- Before the production change,
  `rtk go test -run '^TestReconcileFallbackRequiresArchivedEvidence/active_Spec_evidence_preserves$' ./internal/worktree`
  exited 1 because active-path evidence incorrectly left the Run Branch in
  `Releasable`.
- `rtk go test -run '^TestReconcile(FallbackRequiresArchivedEvidence|PresentTargetAcceptsEitherSpecPath|FallsBackToTheDefaultBranch)$' ./internal/worktree`
  exited 0 with 7 tests passing.
- `rtk go test ./internal/worktree` exited 0 with 123 tests passing.
- `rtk make verify-incremental` exited 0 after formatting, vet, all Go package
  tests, skill checks, and the build.
- The Task's authored `## Verification` commands were not run; the Daemon owns
  those checks and Task settlement.

### Acceptance-criterion evidence

- **Absent target with archived-path evidence releases:**
  `TestReconcileFallbackRequiresArchivedEvidence/archived_Spec_evidence_releases`
  asserts the archived report path is the release proof.
- **Absent target with active-only evidence preserves:**
  `TestReconcileFallbackRequiresArchivedEvidence/active_Spec_evidence_preserves`
  asserts no release and a reason containing `archived Spec path` and
  `Spec is not archived`.
- **Present target accepts either path:**
  `TestReconcilePresentTargetAcceptsEitherSpecPath` covers active and archived
  report paths and asserts the same release classification for both.
