---
task: task_07
spec: 0056-profiles-configure-merge-semantics
status: pending
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
