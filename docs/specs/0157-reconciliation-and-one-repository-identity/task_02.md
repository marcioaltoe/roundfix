---
task: task_02
spec: 0157-reconciliation-and-one-repository-identity
status: completed
type: backend
complexity: low
---

# Task 02: Report identity survives archiving

## Overview

Archive moves a QA Report between the active and archived Spec roots without renaming it, and the supersession comparison breaks ties on the full path, so an archived copy is not recognised as the same report.

## Requirements

1. MUST treat a QA Report and its same-named archived copy as the same report when judging supersession.
2. MUST keep genuinely newer reports winning by date and sequence as today.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A Run-side report and its archived copy are the same report.
- [ ] A newer report still supersedes an older one.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestSupersedingQAReportRecognisesAnArchivedCopy$" ./internal/worktree 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSupersedingQAReportRecognisesAnArchivedCopy"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md)

## Result

### Implementation

- Supersession now compares the Run-side and target-side QA Report filenames,
  independent of whether each report lives below the active or archived Spec
  root, while retaining the target's full path as reconciliation evidence.
- The existing report recency rules still select distinct reports by embedded
  date and numeric sequence.

### Focused checks

- Before the production change,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 -run 'TestSupersedingQAReport(RecognisesAnArchivedCopy|StillPrefersNewerReport)' ./internal/worktree`
  exited 1 because the archived-copy case returned `proven = false`; the later
  date and later sequence cases passed.
- After the production change, the same focused command exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/worktree`
  exited 0 after the final implementation edit.
- The first `make verify-incremental` run reached the test suite but the sandbox
  blocked its GitHub-dependent check. The approved rerun,
  `rtk env GOCACHE=/private/tmp/roundfix-task02-gocache make verify-incremental`,
  exited 0 after vet, all Go package tests, skill checks, and the build.
- `rtk git diff --check` exited 0.
- The Task's authored `## Verification` command was not run; the Daemon owns
  that check and Task settlement.

### Acceptance-criterion evidence

- **A Run-side report and its archived copy are the same report:**
  `TestSupersedingQAReportRecognisesAnArchivedCopy` creates the same filename
  below the Run's active Spec root and the target's archived Spec root, then
  proves supersession with the archived path returned as evidence.
- **A newer report still supersedes an older one:**
  `TestSupersedingQAReportStillPrefersNewerReport` covers both a later date and
  a later same-day numeric sequence, with each target report selected as the
  superseding evidence.
