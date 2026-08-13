---
task: task_13
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 13: Report what the migration moved and what it kept

## Overview

The QA gate found the migration silent: a correct Review Artifact relocation
happened with nothing said about it, and an orphan review retained as live was
retained without its reason reaching the operator. Core Features 5 and 8 already
require the opposite. This slice makes the public surfaces say what happened,
without changing what happens.

## Requirements

1. MUST report every performed relocation in the apply output, at the detail the
   Baseline already gives the repository bytes it manages.
2. MUST report an orphan Review Artifact that was retained rather than relocated,
   naming the liveness answer and its reason.
3. MUST report an undecidable review distinctly from a live one, since ADR-0123
   resolves both to retained for different reasons.
4. MUST keep the plan output's existing relocation reporting unchanged; this
   slice covers apply and retention, which are what the gate found silent.
5. MUST NOT change which files relocate, which reviews are retained, or any
   liveness answer; this slice changes reporting only.

## Subtasks

- [ ] Report performed relocations in apply output.
- [ ] Report retained orphan reviews with their liveness answer and reason.
- [ ] Distinguish an undecidable review from a live one in that report.
- [ ] Cover each report shape with a test.

## Acceptance Criteria

- [ ] Applying a plan with relocations names each performed move in its output.
- [ ] An orphan review retained as live is reported with its reason.
- [ ] An orphan review retained as undecidable is reported distinctly from a live
      one.
- [ ] The set of files relocated and reviews retained is unchanged by this slice,
      proven by a test that fixes the decisions and varies only the report.

## Verification

- `go test -count=1 ./internal/baseline -run 'TestHistoryMoveApplyReport|TestRetainedReviewReport' -v > /tmp/0094-task-13.log 2>&1; s=$?; grep -q '^--- PASS: TestHistoryMoveApplyReport' /tmp/0094-task-13.log && grep -q '^--- PASS: TestRetainedReviewReport' /tmp/0094-task-13.log || { cat /tmp/0094-task-13.log; exit 1; }; exit $s` — expected: exits 0 and the log names both passing tests; fails today, where neither exists.
- `! grep -qi 'no tests to run' /tmp/0094-task-13.log` — expected: exits 0, refusing a vacuous run.
- `grep -q 'ReviewUndecidable' internal/baseline/*.go` — expected: exits 0, proving the report distinguishes the third liveness answer rather than folding it into live. Fails today, where that answer is read in `internal/spec` and never reported by the Baseline.

## Context

- interface: `internal/baseline/apply.go`
- interface: `internal/baseline/history_layout.go`
- interface: `internal/spec/review_liveness.go`

## References

`_prd.md` → Core Features 5 and 8; User Experience, the retained-review sentence.
`_techspec.md` → API Contracts. QA report `qa/qa-report-2026-08-13.md` → F-002,
the silent-relocation half. ADR-0123.
