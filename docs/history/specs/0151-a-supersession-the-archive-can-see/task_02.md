---
task: task_02
spec: 0151-a-supersession-the-archive-can-see
status: completed
type: backend
complexity: medium
---

# Task 02: Archive's second accepted proof

## Overview

Archive proves a Spec is finished by reading completed Tasks and a passing gate.
A Spec whose content another Spec delivered was never decomposed and cannot
acquire either. This Task gives archive a second proof for exactly that Spec,
without loosening the first.

## Requirements

1. MUST accept a recorded supersession in place of completed Tasks and a passing
   gate, for a Spec that has no Task Graph.
2. MUST keep every other archive precondition in force on that path.
3. MUST leave archive's behavior unchanged for a Spec that has a Task Graph,
   whether or not a supersession is present.
4. MUST keep refusing a Spec that has neither a Task Graph nor a supersession,
   with the reason it gives today.
5. MUST move the Spec folder to history intact, deleting and rewriting nothing.

## Subtasks

- [ ] Read the supersession record in archive's preflight.
- [ ] Accept it as proof only where no Task Graph exists.
- [ ] Add tests for all four combinations of Task Graph and supersession.

## Acceptance Criteria

- [ ] A fixture with no Task Graph and a supersession archives.
- [ ] The same fixture without a supersession refuses with the
      missing-Task-Graph reason.
- [ ] A fixture with a Task Graph behaves exactly as before, with and without a
      supersession.
- [ ] The archived folder's files are byte-identical to what they were.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestArchive" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestArchiveAcceptsARecordedSupersession"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — What archive changes

## Result

Implemented archive's second proof only for an active Spec whose Task Graph is
missing. Archive now validates that Spec's `_supersession.md`, retains the
existing missing-Task-Graph refusal when no record exists, and never consults a
supersession when a Task Graph loads. The supersession path applies the existing
destination checks and moves the directory without stamping or rewriting its
PRD. Archive help now describes both proof paths and the byte-preserving move.

Focused checks:

- Before the production change,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1 -run '^TestArchiveAcceptsARecordedSupersession$' ./internal/cli`
  failed because the recorded-supersession fixture still returned the existing
  missing-Task-Graph refusal.
- After the final edit,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1 -run '^(TestArchiveAcceptsARecordedSupersession|TestRunArchiveMovesCompletedSpecAndStampsMetadata|TestRunArchiveRefusesIncompleteTask|TestRunArchiveRefusesMissingOrNonPassingQA|TestRunArchiveHelp)$' ./internal/cli`
  exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -count=1 ./internal/spec`
  exited 0.
- `rtk git diff --check` exited 0.

Acceptance evidence:

- `TestArchiveAcceptsARecordedSupersession/no_Task_Graph_with_supersession_archives_byte-identically`
  exercised a valid record with no Task Graph and observed archive exit 0.
- `TestArchiveAcceptsARecordedSupersession/no_Task_Graph_without_supersession_keeps_the_missing_manifest_refusal`
  observed exit 2 with the unchanged `file does not exist; run the write-tasks
  workflow to create the Task Graph` reason and no filesystem mutation.
- The two `Task_Graph_keeps_its_incomplete-Task_refusal` cases observed the
  same incomplete-Task refusal with and without `_supersession.md`; the focused
  regression set also exercised the existing completed-Task and QA paths.
- The accepted supersession case snapshots every file before the command and
  compares the archived directory byte-for-byte afterward. The destination
  collision case separately proves the supersession path keeps that precondition
  and leaves both source and destination content unchanged on refusal.

The Task's declared `## Verification` command was not run; the Daemon owns that
check and settlement.
