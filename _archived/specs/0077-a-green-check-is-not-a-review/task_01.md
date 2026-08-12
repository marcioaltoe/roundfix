---
task: task_01
spec: 0077-a-green-check-is-not-a-review
status: completed
type: backend
complexity: high
---

# Task 01: Make verified require proof a review ran

## Overview

`verified` is currently reached by falling through: a check on the expected head
that is neither pending nor the one recognised skip title lands on
`settledEvidenceState`, which returns `verified` when no thread is unresolved.
That is how `Review rate limited` — green by design — merged 125 files
unreviewed on Pull Request #107.

This slice inverts the default. `verified` becomes reachable only through a
recognised review-completed shape; everything else resolves `pending`.

It leads the graph because it closes the gate **on its own**, including for
refusal shapes the next Task has not learned yet.

## Requirements

1. MUST make `settledEvidenceState` reachable only when the signal is a
   recognised review-completed shape.
2. MUST resolve any other signal on the expected head as `pending`, including a
   check with a success conclusion whose name or output Roundfix does not
   recognise.
3. MUST NOT infer that a review ran from a conclusion alone. The conclusion is
   green by design when the source refuses.
4. MUST keep every recorded payload that resolves `verified` today resolving
   `verified`, asserted over the existing corpus unchanged. Tightening this gate
   goes wrong by making real reviews stop counting.
5. MUST keep evidence bound to the head it was observed on: a signal recorded
   against an earlier commit settles nothing for the current one.
6. MUST NOT add refusal-specific recognition here. That is task_02, and this
   Task must be provable without it.

## Subtasks

- [ ] Add the recognised review-completed predicate.
- [ ] Route `settledEvidenceState` behind it; default the rest to `pending`.
- [ ] Add the default-deny test for an unrecognised green check.
- [ ] Assert the recorded corpus still resolves as it does today.

## Acceptance Criteria

- [ ] An unrecognised check name with a success conclusion on the expected head
      resolves `pending`, not `verified`.
- [ ] The `Review rate limited` payload from Pull Request #107 resolves
      `pending` after this Task alone, proving the default closes the gate
      without refusal recognition.
- [ ] A recognised completed review with no unresolved thread still resolves
      `verified`.
- [ ] A recognised completed review with unresolved threads still resolves
      `reviewed`.
- [ ] Every payload in the recorded corpus resolves exactly as it does today,
      asserted by the existing tests passing unchanged.
- [ ] A completed review recorded against an earlier commit does not verify the
      current head.

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- instruction: `docs/adr/0054-review-source-evidence-determines-review-outcomes.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/reviewsource/... -count=1 -run 'Evidence|Classify|Verified|Pending' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the classification tests ran and passed.
- `go test ./internal/reviewsource/... ./internal/watch -count=1` — expected:
  exit 0; nothing in the review path regressed.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 2 and 5; Goals; Success Metrics 3 and 5.
- `_techspec.md` → Interfaces; Build Order 1; Risks & Considerations.
- ADR-0054.

## Result

Implementation:

- Retained the CodeRabbit commit-status description so classification can
  distinguish the positive `Review completed` shape from other green statuses.
- Added positive completion predicates for recognised check-run and commit-status
  shapes. Only those predicates and a current-head CodeRabbit approval can reach
  `settledEvidenceState`.
- Routed every other current-head CodeRabbit signal to `pending` after preserving
  explicit skip, reviewing, approval, and failure precedence.
- Added a recorded-status corpus for a completed review and Pull Request #107's
  `Review rate limited` status, plus default-deny check-run cases.

Focused checks:

- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/reviewsource/coderabbit -count=1`
  — exit 0; the complete CodeRabbit package suite passed.
- `rtk env GOCACHE=<worktree>/.gocache go test ./internal/reviewsource/coderabbit -count=1 -run '^(TestEvidenceHierarchyPrecedence|TestEvidenceRecordedCommitStatusCorpus|TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals)$'`
  — exit 0; the focused classification, recorded-payload, and stale-head tests
  passed after the last implementation edit.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

Acceptance evidence:

1. `TestEvidenceHierarchyPrecedence/unrecognised successful check stays pending`
   and `/unrecognised successful check output stays pending` assert that neither a
   green unknown name nor unknown output can verify the expected head.
2. `TestEvidenceRecordedCommitStatusCorpus/pull request 107 rate limit stays pending`
   maps the recorded `CodeRabbit` / `success` / `Review rate limited` payload and
   asserts `pending`; the hierarchy table asserts the same state and evidence kind.
3. The unchanged `successful check verifies with no unresolved threads` case and
   the recorded `Review completed` corpus case both assert `verified`.
4. The unchanged `successful check is reviewed with unresolved thread` case
   asserts `reviewed`.
5. The complete package suite passed, including the existing CodeRabbit evidence
   corpus and its verified outcomes.
6. `TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals` passed with a recognised
   completed check and approval recorded against the earlier head; the current head
   remained `pending`.

Declared `## Verification` commands were not run; the Daemon owns them.

Follow-ups: task_02 owns refusal-specific recognition and task_03 owns operator
diagnostics. Neither is included in this diff.
