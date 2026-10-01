---
task: task_04
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: pending
type: backend
complexity: high
---

# Task 04: Reconcile releases the Runs a merged Spec superseded

## Overview

Reconcile keeps every Run of a Spec merged by squash: commits that changed
only the Spec's directory never match the archived copy, a dirty Run Worktree
is classified before any proof runs, and a Run Branch left without a target
branch has no merged-head fallback. This Task applies ADR-0212: when the
merged head holds the Spec archived, Spec-directory work is represented and
leftovers inside the Spec's declared or recorded scope are superseded, with
forced removal only on a re-proved leftover set; and it describes the rule in
the Roundfix Skill's `reconcile` reference and the `reconcile` command guide.

## Requirements

1. MUST make `inspectRunAtMergedHead` treat paths of non-Task, non-QA commits
   under the Spec's active or archive directory as represented when
   `<archive root>/<slug>/_prd.md` exists at the merged head, and MUST count
   them in the superseded reason of "The merged Spec's Runs".
2. MUST make `inspectTerminalRunMerged` record a dirty worktree's paths, run
   the merged-head proof, and classify `superseded` with the appended reason
   only when the proof yields `safe` or `superseded`, the Spec is archived at
   the merged head, and every dirty path is under the Spec's directories or
   declared (`interface`, `creates`, `deletes`) or recorded by a Task file of
   the archived Spec at that head; otherwise `dirty` with today's reason.
3. MUST make `cleanupTerminalRun` use `git worktree remove --force` only for a
   worktree whose evidence holds superseded dirty paths and whose fresh
   revalidation yields the same state, Run head and dirty path set; every
   other removal MUST stay as today.
4. MUST make `classifyRunBranchSet`, for an absent target branch, apply the
   merged-head proof before preserving a terminal Run, and consider a Run
   whose working root is a linked worktree of the repository.
5. MUST keep preserving any Run with an unrepresented Task commit or a dirty
   path outside the Spec's scope, and MUST NOT read GitHub.
6. MUST add the tests named in Verification to the new file
   `internal/worktree/merged_spec_leftovers_test.go`, over
   `newMergedHeadTestFixture` and `terminalRunFixture`, and MUST leave the
   existing merged-head and inspection tests unedited and passing.
7. MUST describe the rule and its reason words in
   `.agents/skills/roundfix/references/reconcile.md` and
   `docs/user-guide/commands/reconcile.md`, raise the Roundfix Skill's
   version by one patch level from the value on this Task's base in both
   front-matter fields, run `make skills-sync`, and re-record the version.

## Subtasks

- [ ] Represent Spec-directory paths at an archived merged head.
- [ ] Classify declared leftovers as superseded.
- [ ] Force removal only on a re-proved leftover set.
- [ ] Give an absent target branch the merged-head proof.
- [ ] Add the tests.
- [ ] Describe the rule in the skill reference and the command guide, and raise the version.

## Acceptance Criteria

- [ ] A Run whose non-Task commits changed the Spec directory before an
      archive move is `superseded` with the Spec-directory count.
- [ ] A dirty worktree whose changes are a declared path and a Spec-directory
      file is `superseded`, and apply removes it.
- [ ] A dirty undeclared file outside the Spec keeps the Run `dirty`, and an
      unrepresented Task commit keeps it `unintegrated`.
- [ ] Apply refuses as stale when the dirty path set changed after inspection.
- [ ] A Run Branch without its target branch is released when the merged
      head supersedes it.
- [ ] The existing merged-head and inspection tests pass unedited.
- [ ] The reconcile reference and guide carry the reason words, the mirrors
      equal their canonical files, and the raised version is recorded.

## Context

- interface: `internal/worktree/merged_head.go`
- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/merged_spec_leftovers_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/reconcile.md`
- instruction: `internal/worktree/merged_head_test.go`
- instruction: `internal/worktree/worktree_test.go`
- instruction: `internal/spec/recorded_paths.go`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`
- instruction: `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md`
- instruction: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- instruction: `docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/references/2026-10-01-reconcile-releases-the-runs-of-a-squash-merged-spec.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded|TestMergedSpecRunWithDeclaredLeftoversIsSuperseded|TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty|TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated|TestApplyRefusesWhenTheLeftoverSetChanged|TestApplyRemovesASupersededLeftoverWorktree|TestAbsentTargetBranchIsReleasedOnTheMergedHead)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded TestMergedSpecRunWithDeclaredLeftoversIsSuperseded TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated TestApplyRefusesWhenTheLeftoverSetChanged TestApplyRemovesASupersededLeftoverWorktree TestAbsentTargetBranchIsReleasedOnTheMergedHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved|TestMergedHeadRefusesAnUnrepresentedCommit|TestInspectTerminalRunRequiresArchivedEvidence|TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergedHeadDefaultBranchReleasesAnArchivedSpecRunAfterMainMoved TestMergedHeadRefusesAnUnrepresentedCommit TestInspectTerminalRunRequiresArchivedEvidence TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing merged-head and inspection tests run unedited beside the new one, which does not exist before this Task.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/reconcile.md | grep -qF -- "uncommitted path(s) superseded by the archived Spec" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/reconcile.md "uncommitted path(s) superseded by the archived Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/reconcile.md | grep -qF -- "Spec-directory path(s) archived" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/reconcile.md "Spec-directory path(s) archived" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- "uncommitted path(s) superseded by the archived Spec" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/reconcile.md "uncommitted path(s) superseded by the archived Spec" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- "Spec-directory path(s) archived" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/reconcile.md "Spec-directory path(s) archived" >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable|TestTaskAuthoringGuidanceNamesDeclarations)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable TestTaskAuthoringGuidanceNamesDeclarations; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the reconcile reference and guide carry none of these words, so the command fails.

## References

- `_prd.md` → User Story 5; Core Features 7-8; Success Metric 5; Acceptance evidence
- `_techspec.md` → The merged Spec's Runs; API Contract 5; Testing Approach 4; Build Order 4
- ADR-0212; ADR-0053; ADR-0161; ADR-0115; ADR-0052
