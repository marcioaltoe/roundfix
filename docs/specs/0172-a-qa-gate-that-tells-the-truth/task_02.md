---
task: task_02
spec: 0172-a-qa-gate-that-tells-the-truth
status: pending
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
