---
task: task_02
spec: 0172-a-qa-gate-that-tells-the-truth
status: completed
type: backend
complexity: medium
---

# Task 02: The allocator writes where the reader looks

## Overview

`writeMechanicalQAReport` in `internal/daemon/task_engine.go` creates the first free name of the day, so with `qa-report-D.md` and `qa-report-D-02.md` present it writes `qa-report-D-01.md`. `NewestQAReport` in `internal/spec/qa.go` picks the highest suffix, so settlement reads the older `-02` report and the gate settles on a stale verdict. The names come from the Spec's `qa/` directory, where earlier gates, deleted reports and hand-run qa-gate sessions leave gaps, and the chosen path is read by settlement and by the QA Agent's prompt.

## Requirements

1. MUST allocate the new report of date D one above the highest numeric sequence among existing `qa-report-D*.md` names, counting an unsuffixed report as zero: unsuffixed when no report of D exists, `-01` after an unsuffixed report alone, one above the highest otherwise; an earlier gap is never filled.
2. MUST keep the exclusive create and, on a collision, re-read the directory and allocate again.
3. MUST refuse, with an error naming the report and writing nothing, when a report dated after D already exists, because `NewestQAReport` would never select the new report.
4. MUST replace the qa-gate skill's naming sentence in `.agents/skills/qa-gate/SKILL.md` so it says the suffix is one above the highest existing suffix of that date (the phrase `one above the highest`), removing `next unused numeric`, and regenerate `skills/qa-gate/SKILL.md` with `make skills-sync`.
5. MUST put the new tests in `internal/daemon/qa_report_sequence_test.go`, each asserting that `NewestQAReport` returns the path the allocator wrote.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] With `qa-report-D.md` and `qa-report-D-02.md` present, the next report is `qa-report-D-03.md` and `NewestQAReport` returns it.
- [ ] A fresh date starts unsuffixed, and an unsuffixed report alone is followed by `-01`.
- [ ] A report dated after D refuses the allocation and no file is written.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/spec/qa.go`
- interface: `.agents/skills/qa-gate/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMechanicalQAReportAllocatesAboveTheHighestSequence|TestMechanicalQAReportStartsUnsuffixedOnAFreshDate|TestMechanicalQAReportFollowsTheUnsuffixedReportWithSequenceOne|TestMechanicalQAReportRefusesBehindALaterDatedReport)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMechanicalQAReportAllocatesAboveTheHighestSequence TestMechanicalQAReportStartsUnsuffixedOnAFreshDate TestMechanicalQAReportFollowsTheUnsuffixedReportWithSequenceOne TestMechanicalQAReportRefusesBehindALaterDatedReport; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "one above the highest" .agents/skills/qa-gate/SKILL.md && ! grep -q "next unused numeric" .agents/skills/qa-gate/SKILL.md && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the named cases exists and the skill still says `next unused numeric`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The allocator

## Result

Implementation:

- The mechanical QA Report allocator now re-reads the report directory before
  each exclusive-create attempt, allocates one above the highest numeric suffix
  for the current date, and never fills an earlier gap. An `os.ErrExist`
  collision restarts allocation from a fresh directory snapshot.
- Allocation now refuses when any valid later-dated QA Report name is present;
  the error names that report and no current-date report is created.
- The qa-gate naming rule now says `one above the highest existing suffix of
  that date`; `make skills-sync` regenerated the shipped skill mirror.
- The allocation coverage lives in
  `internal/daemon/qa_report_sequence_test.go`, and every successful allocation
  asserts that `NewestQAReport` selects the path just written.

Focused checks:

- Red reproduction before the production change:
  `GOCACHE=/private/tmp/roundfix-0172-task02-gocache go test -count=1 -run
  '^TestMechanicalQAReport(AllocatesAboveTheHighestSequence|StartsUnsuffixedOnAFreshDate|FollowsTheUnsuffixedReportWithSequenceOne|RefusesBehindALaterDatedReport)$'
  ./internal/daemon` failed because the gap case wrote `-01` instead of `-03`
  and the later-date case wrote a current-date report instead of refusing.
- `GOCACHE=/private/tmp/roundfix-0172-task02-gocache go test -count=1 -run
  '^(TestMechanicalQAReport.*|TestWriteMechanicalQAReportPreservesSameDayNamingAndPriorReport)$'
  ./internal/daemon` exited 0 after the implementation change.
- `GOCACHE=/private/tmp/roundfix-0172-task02-gocache go test -count=1 -run
  '^TestNewestQAReport' ./internal/spec` exited 0.
- `make skills-sync-check` exited 0. A focused search found `one above the
  highest` in both QA skill copies and no `next unused numeric` occurrence.

Acceptance evidence:

- Highest sequence with a gap: `TestMechanicalQAReportAllocatesAboveTheHighestSequence`
  wrote `qa-report-D-03.md` beside the unsuffixed and `-02` reports, then
  observed `NewestQAReport` return `-03`.
- Fresh and second same-day allocations:
  `TestMechanicalQAReportStartsUnsuffixedOnAFreshDate` observed the unsuffixed
  name, and `TestMechanicalQAReportFollowsTheUnsuffixedReportWithSequenceOne`
  observed `-01`; both observed `NewestQAReport` return the allocated path.
- Later-date refusal: `TestMechanicalQAReportRefusesBehindALaterDatedReport`
  observed an error containing the later report's name, an empty returned path,
  and an unchanged QA directory containing only that later report.

The Task's authored `## Verification` command was not run; the Daemon owns that
command and settlement.
