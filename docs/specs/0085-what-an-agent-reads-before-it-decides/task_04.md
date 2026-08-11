---
task: task_04
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: docs
complexity: high
---

# Task 04: Relocate the retired artifacts under one root

## Overview

The move Task 03 made cheap. Archived Specs and findings relocate under the
single archive root, the resolver's answers change once, and the corpus golden is
re-recorded deliberately with the reason stated rather than regenerated on
reflex.

## Requirements

1. MUST relocate archived Specs and archived findings under the single archive
   root, preserving history for every moved path.
2. MUST change the resolver's answers in one place, so no consumer needs editing.
3. MUST re-record the characterization corpus golden in the same commit, with the
   reason recorded, and MUST NOT regenerate it silently.
4. MUST break the corpus-golden case Task 01 declared, and update it to the new
   layout in the same commit.
5. MUST leave every archived artifact byte-identical in content; this Task moves
   files and changes nothing inside them.

## Subtasks

- [ ] Relocate the archived Specs and findings.
- [ ] Change the resolver's answers.
- [ ] Re-record the corpus golden with its reason.

## Acceptance Criteria

- [ ] Every archived Spec and finding resolves under the single root.
- [ ] The resolver is the only file that changed to make that true.
- [ ] The corpus golden is re-recorded with the reason stated.
- [ ] No archived artifact's content changed.

## Bounded scope

This Task may create or modify only:

- `internal/spec/archive.go`
- `internal/spec/archive_layout_characterization_test.go`
- `internal/speccheck/testdata/corpus-golden.json`
- `docs/specs/_archived/**`
- `docs/findings/_archived/**`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveDirAnswersEveryRetiredKind' -count=1 -v 2>&1 | grep -q '^--- PASS'` — expected: exits 0.
- `grep -q 'Re-recorded because' internal/speccheck/testdata/corpus-golden.json` — expected: exits 0. The `|| grep ... task_04.md` fallback this command used to carry made it self-satisfying: the string occurs in this file because the command writes it, so the check passed before any re-recording happened.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 4.
