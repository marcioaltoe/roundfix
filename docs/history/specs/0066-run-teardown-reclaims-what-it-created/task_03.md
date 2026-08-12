---
task: task_03
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
type: backend
complexity: medium
---

# Task 03: Stop blocking review on superseded Run Branch work

## Overview

Branch Integrity Preflight refuses fetch, resolve, and watch while
unintegrated Run Branch commits remain bound to the PR Head Branch. That
refusal is correct for work that would be lost, and wrong for accumulated
superseded QA reports — which is how four failed cycles came to refuse the
`watch` that would have resolved them.

This slice teaches the Preflight the difference. It is this Spec's **declared
break**: every other refusal it makes today it still makes.

## Requirements

1. MUST let a review Run proceed when the only unintegrated Run Branch work is
   proven superseded by task_02's classification.
2. MUST keep refusing for every other reason it refuses today, including
   unintegrated implementation work, another Run bound to the branch, and any
   branch it cannot classify.
3. MUST name, in its diagnostic, which branches it disregarded as superseded
   and the proof, so the relaxation is visible rather than silent.
4. MUST NOT delete, move, or modify any branch. Preflight decides whether to
   proceed; reclamation is task_04's surface.
5. MUST NOT weaken the explicit bypass or its audit comment.

## Subtasks

- [ ] Consume the set classification in the Preflight decision.
- [ ] Proceed when the only obstruction is proven superseded.
- [ ] Name the disregarded branches and their proof in the diagnostic.
- [ ] Assert every other refusal unchanged.

## Acceptance Criteria

- [ ] A repository with four failed-cycle Run Branches on one target lets
      `watch` proceed, with no ref hand-deleted.
- [ ] The diagnostic names each disregarded branch and its proof.
- [ ] Unintegrated implementation work still refuses.
- [ ] An unclassifiable branch still refuses.
- [ ] Another Run bound to the branch still refuses.
- [ ] No branch is modified by the Preflight, proven by Git state identical
      before and after.
- [ ] The explicit bypass and its audit comment are unchanged.

## Context

- interface: `internal/preflight`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/preflight ./internal/cli -count=1 -run 'Preflight|BranchIntegrity' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the Preflight tests ran and passed.
- `go test ./internal/preflight ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 4; Decisions (declared break); Success Metric 2.
- `_techspec.md` → API Contracts; Build Order 3.
- ADR-0052.

## Result

### Implementation

- Branch Integrity Preflight now groups target-bound Implement Run Branches by
  Spec and consumes `ClassifyRunBranchSet`. It disregards only a proven failed
  QA set: the current-evidence branch when older branches are releasable, plus
  each branch classified `Releasable`.
- Each disregarded branch emits a stderr diagnostic with its branch, worktree,
  ahead count, and the current or target QA Report that proves the decision.
  The decision path performs no Git mutation; branch reclamation remains owned
  by task_04.
- Classification errors, `Preserved` branches, unintegrated implementation
  work, and Active spec Run Branches remain blocking even when Git reports
  them fast-forwardable. The existing Active review Run guard and explicit
  bypass/audit path are unchanged.

### Focused checks

- Pre-change signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task03-cache go test ./internal/cli -run '^TestBranchIntegrityPreflightWatchDisregardsFourFailedQACycles$' -count=1`
  — failed because all four failed-cycle branches were listed as pending Run
  Branch work and `watch` exited 2.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task03-cache go test ./internal/cli -run 'BranchIntegrity' -count=1`
  — passed (`ok roundfix/internal/cli`).
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-task03-cache go test ./internal/worktree -run '^TestClassifyRunBranchSet' -count=1`
  — passed (`ok roundfix/internal/worktree`).
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance criterion evidence

- Four failed cycles and no hand-deleted ref:
  `TestBranchIntegrityPreflightWatchDisregardsFourFailedQACycles` creates four
  real Run Branches and terminal failed QA Runs, then observes `watch` proceed
  without calling the integration seam.
- Visible proof per branch: the same test matches one diagnostic line per Run
  Branch. The newest names its current QA Report; each older branch names that
  report as its superseding proof.
- Unintegrated implementation work:
  `TestBranchIntegrityPreflightClassifiesPendingRunBranches/task_work`
  supplies a fast-forwardable `Preserved` implementation branch and observes
  exit 2 with no integration.
- Unclassifiable work:
  `TestBranchIntegrityPreflightClassifiesPendingRunBranches/unclassifiable_branch`
  makes set classification fail and observes the branch preserved as a
  blocking integration candidate.
- Another Run bound to the branch:
  `TestBranchIntegrityPreflightRefusesActiveImplementRunBranch`,
  `TestBranchIntegrityPreflightRejectsActiveRunForReviewCommands`, and
  `TestBranchIntegrityRejectsLocklessBypassRun` passed in the focused suite.
- No branch modification: the four-cycle test compares the complete
  `for-each-ref` name/object snapshot before and after `watch`; the snapshots
  are identical.
- Explicit bypass and audit comment:
  `TestBranchIntegrityBypassPublishesAuditBeforeFetch` and
  `TestBranchIntegrityBypassAuditsActiveRunAndProceeds` passed unchanged; the
  bypass and audit-rendering functions were not edited.

### Follow-up

- Applying reclamation to the now-disregarded branches remains task_04's
  surface and is intentionally absent from this diff.
