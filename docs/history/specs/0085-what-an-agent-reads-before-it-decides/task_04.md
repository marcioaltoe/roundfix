---
task: task_04
spec: 0085-what-an-agent-reads-before-it-decides
status: completed
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

- `_archived/specs/**`
- `_archived/findings/**`

- `internal/spec/archive.go`
- `internal/spec/archive_layout_characterization_test.go`
- `internal/docscontract/testdata/corpus-golden.json`
- `docs/specs/_archived/**`
- `docs/findings/_archived/**`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_04.md`

Corrected on 2026-08-11. The list first named
`internal/speccheck/testdata/corpus-golden.json`, which does not exist — the
live golden is under `internal/docscontract/` — and omitted the destinations
`_archived/specs/**` and `_archived/findings/**` that Requirements 1 and 4
oblige. The Agent refused rather than move files outside its boundary and named
both gaps. Enumerating a boundary from guessed paths is what cost this Run, for
the second time in this Spec.

## Verification

- `grep -q '"_archived/specs"' internal/spec/archive.go` — expected: exits 0. The resolver answers `docs/specs/_archived` today; this Task is what moves its answer to the single root, so asserting the resolver test passes proves nothing — Task 02 created that test and it already passes.
- `test -d _archived/specs && test -d _archived/findings` — expected: exits 0, proving the artifacts moved rather than only the resolver's answer.
- `grep -q 'Re-recorded because' internal/docscontract/testdata/corpus-golden.json` — expected: exits 0. The `|| grep ... task_04.md` fallback this command used to carry made it self-satisfying: the string occurs in this file because the command writes it, so the check passed before any re-recording happened.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 4.

## Result

Relocated the existing archived Spec and finding trees to
`_archived/specs/` and `_archived/findings/`. The archive resolver now owns the
four `_archived/<kind>` answers, and the default Archive Command destination
uses that resolver while a configured non-default Spec Root retains its
colocated archive behavior. Re-recorded the active-corpus golden with the
reason for its intentional update and updated Task 01's characterization at
the declared break.

- Criterion 1: the updated archive-layout characterization passed and observed
  the relocated Spec and finding directories. Raw destination counts are 2,142
  Spec files and 65 finding files.
- Criterion 2: a production-source search for the four new literal paths found
  them only in `internal/spec/archive.go`; the consumer packages continue to
  ask `ArchiveDir` for the layout.
- Criterion 3: changing only the golden first made
  `TestArchiveLayoutCharacterizationPinsCorpusGoldenBeforeRelocation` fail on
  the old `update` explanation. After the characterization update, the focused
  characterization and docs-contract golden tests passed. The golden now
  starts its explanation with `Re-recorded because` and names this relocation.
- Criterion 4: the pre-move Git trees were
  `1781d09ecf5f525fb3d4231fe1d83435ea14c34b` for Specs and
  `d506a2f8d845fb3fbabcc0317100ef44b00be3ce` for findings. After extracting
  those `HEAD` trees to
  `/private/tmp/roundfix-spec0085-task04.SXZOOB`, `rtk diff -qr` against each
  relocated destination exited 0 with no output, proving every moved file is
  byte-identical at the same relative path.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-spec0085-task04-gocache go test
  ./internal/spec -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenBeforeRelocation$'
  -count=1` failed before the characterization update with the intentional
  old-versus-new golden explanation mismatch.
- `rtk env GOCACHE=/private/tmp/roundfix-spec0085-task04-gocache go test
  ./internal/spec -run '^TestArchiveLayoutCharacterization' -count=1` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-spec0085-task04-gocache go test -tags
  docscontract ./internal/docscontract -run '^TestCheckCorpusGolden$' -count=1`
  passed.
- `rtk git diff --check` reported no whitespace errors.
- A changed-path inspection reported 4,404 old/new move entries and bounded
  file edits, with zero paths outside this Task's scope.

Follow-up boundary note: predecessor unit and fixture tests outside this Task's
bounded paths still contain the old archive layout. They were not edited here;
the Run-level regression gate owns surfacing and routing any required
follow-up.

The Daemon-owned `## Verification` commands were not run.
