---
task: task_13
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

Apply results now carry the transaction's verified History Relocation ledger.
Both `baseline apply` and the apply phase of `baseline update` render every
performed move with its source, destination, and content identity; an idempotent
reapply reports no newly performed move. The existing plan History Relocation
block is unchanged.

Layout discovery now records each retained orphan Review Artifact as a Baseline
finding. Live and undecidable reviews use distinct finding codes, and each
finding names the Review Artifact path, the liveness answer, and the reason
returned by the existing local-Git classifier. These findings flow through the
Plan and apply result warning collections without influencing relocation or
retention decisions.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-task-13-gocache go test -count=1 ./internal/baseline -run '^TestHistoryMoveApplyReport$'`
  and the corresponding `'^TestRetainedReviewReport$'` run both failed to build:
  `Result.VerifiedHistoryMoves` did not exist and `planHistoryMoves` returned no
  report collection.
- After implementation, each focused command above exited 0 independently.
- `rtk env GOCACHE=/tmp/roundfix-task-13-gocache go test -count=1 ./internal/baseline`
  exited 0 (`ok roundfix/internal/baseline`, 115.641s).
- `rtk env GOCACHE=/tmp/roundfix-task-13-gocache go test -count=1 ./internal/cli -run '^TestBaseline(ApplyTextReportsHistoryMoves|UpdateTextReportsHistoryMoves)$' -v`
  exited 0 and named both passing renderer tests.
- `rtk env GOCACHE=/tmp/roundfix-task-13-gocache go test -count=1 ./internal/cli`
  exited 0 (`ok roundfix/internal/cli`, 139.463s).
- `rtk git diff --check` exited 0. The task's declared `## Verification`
  commands remain pending Daemon execution and were not run in this Agent turn.

Acceptance evidence:

1. `TestHistoryMoveApplyReport` fixes two planned moves, applies them, and proves
   the result's `verifiedHistoryMoves` contains each source, destination, and
   content identity. Its negative control reapplies the same Plan and proves no
   move is reported as newly performed.
2. `TestRetainedReviewReport` fixes a reachable non-ancestor review, proves it
   remains unmoved, and verifies the `baseline.history.review.live` finding
   includes the classifier's recorded-head reason.
3. The same test fixes a no-head review, proves it remains unmoved, and verifies
   the distinct `baseline.history.review.undecidable` finding includes the newest
   Round reason. The production switch names `spec.ReviewUndecidable` explicitly.
4. The retention test first captures the public relocation and collision sets,
   then obtains the added report and proves both retained states still produce no
   moves or collisions. `DiscoverHistoryLayout` continues to project only the
   unchanged relocation and collision decisions from the internal report.
