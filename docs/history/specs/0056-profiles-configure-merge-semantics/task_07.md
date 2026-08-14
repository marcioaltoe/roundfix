---
task: task_07
spec: 0056-profiles-configure-merge-semantics
status: completed
type: backend
complexity: medium
---

# Task 07: Remove a category without touching its neighbours' spacing

## Overview

Deleting a category removes its key/value pair and also swallows the blank
line that separated `profiles` from the next top-level key. The data is
correct, but the promise this Spec makes is that a write touches only the
lines the change requires — a removal that reflows unrelated spacing is the
same class of noise the original defect hid behind, just smaller.

## Requirements

1. MUST leave blank lines that do not belong to the removed category exactly
   where they were, including the separator between the `profiles` mapping and
   any following top-level key.
2. MUST leave blank lines and comments belonging to surviving categories
   untouched when a neighbouring category is removed.
3. MUST remove the blank lines and comments that belong to the removed
   category itself, so a removal leaves no orphaned fragment behind.
4. MUST keep every guarantee earlier slices established: byte-identical
   untouched categories, atomic replacement, the alias-anchor rejection, and
   the empty-file and absent-section paths.
5. MUST extend the characterization corpus with the removal-adjacent spacing
   cases rather than relaxing an existing golden to accept the churn.

## Subtasks

- [ ] Reproduce the swallowed separator in a failing corpus case.
- [ ] Attribute blank lines and comments to the category they belong to.
- [ ] Remove only the removed category's own trivia.
- [ ] Confirm replacement and addition paths did not change.

## Acceptance Criteria

- [ ] Removing a category from a file whose `profiles` mapping is followed by
      another top-level key leaves the separating blank line in place.
- [ ] Removing a middle category leaves the blank lines and comments around
      its surviving neighbours unchanged.
- [ ] Removing a category removes its own leading comment and trailing blank
      line, leaving no orphaned trivia.
- [ ] A removal's diff touches only the removed category's lines.
- [ ] Every earlier corpus case still passes, and the new cases were added
      rather than existing goldens loosened.
- [ ] `git status --porcelain` shows no path outside `internal/config/` and
      this task file.

## Context

- interface: `internal/config/profile_config.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/config -run TestProfilesConfigureRemovalPreservesSpacing -count=1`
  — expected: exit 0; the separator and neighbouring trivia survive a removal.
- `go test ./internal/config -run TestProfilesConfigureMergePreservesOtherCategories -count=1`
  — expected: exit 0; the guarantees from task 03 still hold.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0; the corpus passes with the removal cases added.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 3; Core Features 4; Non-Goals (no general YAML
  formatting pass).
- `_techspec.md` → Implementation Design: Interfaces; Risks.
- `qa/qa-report-2026-08-01.md` → F-003.

## Result

### Implementation

- The removal path now attributes source-only blank lines to the surviving
  category or following top-level key that owns them. Unique internal comment
  markers carry those exact bytes through `yaml.v3` encoding and are removed
  before validation, while the removed category's node comments and trailing
  blank lines remain absent.
- `TestProfilesConfigureRemovalPreservesSpacing` covers the top-level
  separator, a middle category's owned comment and trailing blank, a surviving
  neighbour's blank line, and dangling-alias rejection without persistence.
- The characterization corpus gained two new input/golden pairs for removal
  before a top-level separator and removal beside category trivia. No existing
  golden was changed.

### Focused checks

- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/config -run 'TestProfilesConfigureRemovalPreservesSpacing/.+$' -count=1 -v`
  — exit 0; all four removal spacing and alias-rejection subtests passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/config -run 'TestProfilesConfigWriterCharacterization/.+$' -count=1 -v`
  — exit 0; all 11 cases passed, including the nine prior cases and two new
  removal cases.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/config -run 'TestProfilesConfigureMergePreservesOtherCategories/(replaces_one_category_without_touching_the_other_four|replaces_the_whole_category_object|appends_an_added_category|adds_profiles_to_empty_and_sectionless_configs|rejects_a_replacement_that_strands_an_alias)$' -count=1 -v`
  — exit 0; replacement, atomic-category replacement, addition,
  empty/sectionless, and alias canaries passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/config -run 'TestProfileConfigAtomic(WritesUserAndProjectProfiles|DryRunAndFailuresLeaveBytesUnchanged)$' -count=1 -v`
  — exit 0; atomic persistence and no-mutation failure canaries passed.
- `rtk git diff --check` — exit 0.
- `rtk git status --porcelain` — exit 0 with the worktree's fsmonitor warning;
  every changed path is under `internal/config/` or is this task file.

### Acceptance evidence

1. `keeps_separator_before_following_top-level_key` compares the complete
   written file byte-for-byte and retains the blank line before `watch`.
2. `keeps_surviving_category_trailing_blank_line` retains the surviving
   category's blank line and the following category's comment when the middle
   category is removed.
3. `removes_only_middle_category_trivia` proves the removed category's leading
   comment and trailing blank leave no orphaned fragment.
4. Every removal spacing case compares the complete output with an exact
   expected file, so any diff outside the removed category fails the test.
5. The 11-case corpus check passed; Git status lists only the four newly added
   removal corpus files and no modified prior golden.
6. The changed-path status check lists only `internal/config/` and this task
   file.

The commands under `## Verification` were not run; Daemon Verification owns
those checks and the terminal Task status.
