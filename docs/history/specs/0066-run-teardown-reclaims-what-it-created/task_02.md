---
task: task_02
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
type: backend
complexity: high
---

# Task 02: Classify a target's Run Branches as a set

## Overview

Spec 0053 shipped a `superseded` classification proving one Run Branch holds
nothing but an obsolete QA report. The case its design does not cover is
several such branches at once, none passing — four QA cycles left four Run
Branches that neither `reconcile --apply` nor Branch Integrity Preflight could
clear, and the accumulation then refused the very `watch` those cycles were
trying to make Clean.

This slice classifies the set: the newest carries current evidence, older
superseded-only branches become releasable, and anything ambiguous is
preserved.

## Requirements

1. MUST classify the Run Branches of one target together, identifying the one
   holding current evidence and the older superseded-only branches.
2. MUST reuse the existing superseded proof rather than introducing a second
   notion of obsolescence.
3. MUST preserve any branch it cannot prove superseded, with a reason, and MUST
   NOT release a branch carrying unintegrated implementation work.
4. MUST never classify a branch belonging to an Active Run as releasable, per
   ADR-0052.
5. MUST leave the existing single-branch classification behaving as it does
   today.

## Subtasks

- [ ] Enumerate the Run Branches of one target.
- [ ] Identify the current-evidence branch and the superseded-only older ones.
- [ ] Preserve the ambiguous with reasons.
- [ ] Add the four-branch fixture and the Active Run guard.

## Acceptance Criteria

- [ ] Four failed-cycle Run Branches on one target classify as one current and
      three releasable.
- [ ] A branch carrying unintegrated implementation work is never releasable.
- [ ] A branch that cannot be proven superseded is preserved with a reason.
- [ ] An Active Run's branch is never releasable, proven by an injected active
      fixture.
- [ ] The existing single-branch superseded classification is unchanged,
      asserted by the existing tests passing.

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/gittest`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/worktree -count=1 -run 'BranchSet|Superseded|Classify|Active' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the set-classification tests ran and passed.
- `go test ./internal/worktree -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Success Metric 2.
- `_techspec.md` → Interfaces; Build Order 2.
- ADR-0052, ADR-0053.

## Result

### Implementation

- Added target-scoped Run Branch set classification in `internal/worktree`.
  The classifier uses Run metadata to attribute branches and protect Active
  Runs, selects the branch with the newest QA Report as current evidence, and
  reuses `supersedingQAReport` for every releasable decision.
- Added a real-Git four-cycle fixture plus negative fixtures for another
  target, unintegrated implementation work, malformed supersession evidence,
  and an Active Run branch. Preserved branches carry bounded reasons.

### Focused checks

- Pre-change signal: the focused branch-set tests failed to compile because
  `ClassifyRunBranchSet` and `BranchSetClassification` did not exist.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-spec0066-go-cache go test ./internal/worktree -run '^(TestClassifyRunBranchSet.*|TestInspectTerminalRunClassifiesSupersededQAReport)$' -count=1`
  — passed (`ok roundfix/internal/worktree`).
- The commands under `## Verification` were not run; the Daemon owns them.

### Acceptance criterion evidence

- Four failed cycles: `TestClassifyRunBranchSetFourFailedCycles` observes the
  newest branch as `Current` and the other three as `Releasable`.
- Unintegrated implementation work:
  `TestClassifyRunBranchSetPreservesUnintegratedImplementationWork` observes
  the valuable-work branch only in `Preserved`, with a reason.
- Unprovable supersession:
  `TestClassifyRunBranchSetPreservesBranchWithoutSupersededProof` observes a
  QA-path commit that violates the Daemon QA-commit contract preserved with a
  reason.
- Active Run guard: `TestClassifyRunBranchSetPreservesActiveRunBranch` injects
  `store.StateActive` into an otherwise releasable older branch and observes
  it preserved, never releasable.
- Existing single-branch behavior:
  `TestInspectTerminalRunClassifiesSupersededQAReport` passed unchanged in the
  focused check. The Daemon-owned full `internal/worktree` suite remains the
  terminal regression check.
