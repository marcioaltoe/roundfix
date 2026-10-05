---
task: task_06
spec: 0228-queue-items-that-stay-current-with-main
status: completed
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

## Result

Implemented the assigned corrective slice for Daemon Verification:

- `resolveDerivedLineHunks` now collects both sides of each hunk, refuses
  unequal line counts, and compares corresponding lines. Identical lines are
  preserved; a differing line must match the declaration on both sides. The
  output keeps the default branch's side. The regeneration bound in
  `derivedLineChangesAllowed` is unchanged.
- Added real local Git fixtures with canonical and mirrored skill frontmatter,
  the repository's derived-record and line-scoped declarations, a local bare
  default-branch remote, and shell regeneration. An invocation log outside
  the worktree proves the command ran exactly once or never ran, even when a
  merge is aborted.
- Updated Delivery conflict recovery in the configuration guide to describe
  equal line counts and differences confined to matching lines, including
  "identical lines between them are kept" and the default-side rule.

Acceptance evidence from focused implementation checks:

1. Real skill collision merges and regenerates one patch above main:
   `TestARealSkillVersionHunkIsMergedAndRegeneratedOnePatchAboveMain` first
   proves, through its own two-sided Git merge followed by abort, exactly one
   canonical-file conflict hunk with both versions and the author/source lines
   between them. It then asserts both version fields in both files become
   `0.1.33`, both branches' body edits survive, the record is regenerated,
   exactly one command ran, the merge parents are item then main, the
   `Roundfix-Delivery: derived-merge` trailer is present, no source paths are
   returned, and the worktree is clean.
2. Non-matching differences remain source conflicts:
   `TestAVersionHunkWhoseInterveningLineDiffersAbortsTheMerge` asserts that a
   differing author line returns both `SKILL.md` paths, runs no regeneration,
   preserves the item head and clean worktree, and leaves no `MERGE_HEAD`.
   Every existing test in `deliver_conflict_test.go` and
   `deliver_derived_lines_test.go` passed unedited in the final focused run,
   including the regeneration-bound and whole-file-precedence tests.

Red and mutation evidence (Requirement 6):

- Before the production fix,
  `GOCACHE=/tmp/roundfix-task06-gocache rtk proxy go test ./internal/cli -count=1 -run '^TestARealSkillVersionHunk' -v`
  exited 1 at the resolution assertion:
  `Head:` empty,
  `SourcePaths:[.agents/skills/example/SKILL.md skills/example/SKILL.md]`,
  `Regenerated:[]`. The fixture's one-hunk proof had already succeeded.
- With a temporary mutation that exempted differing `  author: ` lines from
  the pattern check,
  `GOCACHE=/tmp/roundfix-task06-gocache rtk proxy go test ./internal/cli -count=1 -run '^TestAVersionHunkWhoseInterveningLine' -v`
  exited 1: the refusal assertion saw a merged head, `SourcePaths:[]`, and
  a regeneration command. The mutation was removed before the final run.
- With the intended fix, the focused two-test selection
  `GOCACHE=/tmp/roundfix-task06-gocache rtk proxy go test ./internal/cli -count=1 -run 'Test(ARealSkillVersionHunk|AVersionHunkWhoseInterveningLine)' -v`
  exited 0 with both tests passing.
- Final focused regression check after removing the mutation:
  `GOCACHE=/tmp/roundfix-task06-gocache rtk proxy go test ./internal/cli -count=1 -run 'Test.*(Conflict|Regeneration|DerivedPaths|VersionHunk)' -v`
  exited 0 (`ok roundfix/internal/cli`, 9.563s). Both new tests and every
  existing test in the two protected test files passed.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0, and
  `rtk proxy git -c core.fsmonitor=false diff --exit-code -- internal/cli/deliver_conflict_test.go internal/cli/deliver_derived_lines_test.go`
  exited 0, confirming both existing files are unchanged.

The initial worktree already contained the Daemon's `pending` to
`in_progress` status change. That status is preserved. No Task Graph, other
Task, owned skill, or tooling configuration was edited. No declared
Verification command or repository settlement gate was run; these remain
Daemon-owned. No commit, push, or Pull Request was created. No follow-up was
needed for this slice.

## Carry-forward provenance

- Source Run: `run_20261005T112526Z_f4195dc79490f939`
- Source commit: `b957f7cfe233c19e865dfdf80b10465b77bdac75`
