---
task: task_02
spec: 0151-a-supersession-the-archive-can-see
status: pending
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
