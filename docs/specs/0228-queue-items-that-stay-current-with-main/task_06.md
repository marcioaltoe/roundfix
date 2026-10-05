---
task: task_06
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: backend
complexity: medium
---

# Task 06: The real skill version hunk is merged and regenerated

## Overview

Corrective Task for QA finding F1 of the 2026-10-05 QA Report (row Q03,
Requirement 3) of Run `run_20261005T101255Z_f599a5cc8ea7f4c2`. The Roundfix
Skill keeps its metadata version field, then the `author` and `source` lines,
then the top-level version field. When two branches raise both fields, Git
merges the two changes into one conflict hunk whose sides also hold the
identical `author` and `source` lines. `resolveDerivedLineHunks` in
`internal/cli/deliver_workflow.go` refuses any non-matching line inside a
hunk, so the real collision parks as a source conflict before the
regeneration runs. The fixtures of task_02 put the two version fields far
apart, so each lands in its own hunk and the tests pass.

The intended rule (Requirement 3, ADR-0233): a hunk whose two sides differ
only on lines matching the declared pattern takes the default branch's side;
identical lines between them are allowed. A hunk where any non-matching line
differs between the sides is still refused.

## Requirements

1. MUST change `resolveDerivedLineHunks` so that a conflict hunk is resolved
   when its two sides have the same number of lines and every line that
   differs between them matches `match` on both sides. Lines identical on
   both sides are kept. The resolver writes the default branch's side of each
   resolved hunk, as today.
2. MUST keep refusing, as a source conflict that parks
   `pull-request-conflict: <path>`, a hunk whose sides differ on any
   non-matching line or have different line counts. The regeneration bound
   of task_02, which allows changes only on matching lines with the line count
   kept, MUST stay unchanged.
3. MUST add `internal/cli/deliver_derived_skill_layout_test.go`, over real
   Git repositories with a local default-branch remote and a shell
   regeneration command, with:
   - `TestARealSkillVersionHunkIsMergedAndRegeneratedOnePatchAboveMain`. A
     canonical `.agents/skills/<name>/SKILL.md` and its mirror
     `skills/<name>/SKILL.md` reproduce the Roundfix Skill frontmatter: a
     `metadata:` block with `  version:`, `  author:` and `  source:`, then a
     top-level `version:`. The declaration is this repository's, with
     line-scoped paths `.agents/skills/*/SKILL.md` and `skills/*/SKILL.md`
     and the match `^ *version: `. The item raises both files to `0.1.31` and
     edits a body line, and the default branch raises them to `0.1.32` and
     edits another body line. Each side also changes the derived record.
     Before resolving, the test proves with its own
     `merge.conflictStyle=merge` merge, then aborted, that the canonical file
     has exactly one conflict hunk holding both version fields and the
     `source` line between them. The regeneration command fails unless both
     fields of both files read `0.1.32`, and then raises them to `0.1.33` and
     rewrites the record. The test asserts a merged head with two parents,
     the item head first, the `Roundfix-Delivery: derived-merge` trailer, no
     source paths, one regeneration, `version: 0.1.33` exactly twice in each
     file, both body edits kept and a clean worktree.
   - `TestAVersionHunkWhoseInterveningLineDiffersAbortsTheMerge`: the same
     layout, but the default branch also changes the `author` line inside
     that hunk. The test asserts that both `SKILL.md` paths are source paths,
     that no regeneration ran, and that the item head, a clean worktree and
     no `MERGE_HEAD` remain.
4. MUST update the "Delivery conflict recovery" section of
   `docs/user-guide/configuration.md` so it states that a hunk's sides may
   differ only on matching lines, with the phrase "identical lines between
   them are kept", and keep the phrase "takes the default branch's side of
   each conflict hunk".
5. MUST leave every existing test unedited and passing, in
   `internal/cli/deliver_conflict_test.go` and
   `internal/cli/deliver_derived_lines_test.go`.
6. MUST prove that the first new test fails before the fix with both
   `SKILL.md` paths as source paths, the behavior QA observed, and that the
   second fails when the fix accepts a differing non-matching line. Record
   both in the Result.

## Subtasks

- [ ] Compare the two sides of each hunk line by line and allow identical non-matching lines.
- [ ] Add the real-layout merge test and the differing-line refusal test.
- [ ] Update the configuration guide.

## Acceptance Criteria

- [ ] The real Roundfix Skill version collision merges and regenerates one
      patch above the default branch.
- [ ] A hunk whose sides differ on a non-matching line still parks, and every
      existing derived-merge test passes unedited.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_derived_skill_layout_test.go`
- interface: `docs/user-guide/configuration.md`
- instruction: `internal/cli/deliver_derived_lines_test.go`
- instruction: `internal/cli/deliver_conflict_test.go`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "identical lines between them are kept" || { printf 'missing phrase in configuration guide\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestARealSkillVersionHunkIsMergedAndRegeneratedOnePatchAboveMain|TestAVersionHunkWhoseInterveningLineDiffersAbortsTheMerge|TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated|TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge|TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge|TestWholeFileDerivedPathsTakePrecedenceOverDeclaredLines|TestADerivedConflictIsMergedAndRegenerated|TestASourceConflictAbortsTheMerge|TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge|TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch|TestConflictRecoveryRunsMatchedDeclarationsInOrder)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestARealSkillVersionHunkIsMergedAndRegeneratedOnePatchAboveMain TestAVersionHunkWhoseInterveningLineDiffersAbortsTheMerge TestAConflictConfinedToDeclaredLinesIsMergedAndRegenerated TestAConflictHunkOutsideDeclaredLinesAbortsTheMerge TestARegenerationThatChangesAnUndeclaredLineAbortsTheMerge TestWholeFileDerivedPathsTakePrecedenceOverDeclaredLines TestADerivedConflictIsMergedAndRegenerated TestASourceConflictAbortsTheMerge TestARegenerationThatWritesAnUndeclaredPathAbortsTheMerge TestConflictRecoveryReadsDerivedPathsFromTheDefaultBranch TestConflictRecoveryRunsMatchedDeclarationsInOrder; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide lacks the phrase and the two new tests do not exist, so the command fails.

## References

- QA Report 2026-10-05 of Run `run_20261005T101255Z_f599a5cc8ea7f4c2` → finding F1; row Q03
- task_02 → Requirement 2
- `_prd.md` → Core Feature 2; Success Metric 2
- `_techspec.md` → Line-scoped derived paths; Testing Approach 2
- ADR-0233; ADR-0192
