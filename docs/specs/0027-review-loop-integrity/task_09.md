---
task: task_09
spec: 0027-review-loop-integrity
status: completed
type: backend
complexity: high
---

# Task 09: Propagate per-issue outcomes with Outcome Comments

## Overview

Make the pull request self-auditable: at each Batch settlement — the earliest point compatible with the Verification gate — every settled Review Issue is propagated to the Review Source individually, with an Outcome Comment for every non-resolved outcome. Failed threads stay open with an explanation; invalid and duplicated threads are resolved only after their comment publishes; issues still unresolved at Run end get a closing comment.

## Requirements

1. MUST replace the resolved/invalid-only batch propagation with per-issue propagation at Batch settlement, after Verification has passed for the batch (or the batch has terminally failed): resolved → resolve the thread; invalid → publish an Outcome Comment with the triage reason, then resolve; duplicated → publish an Outcome Comment naming the canonical thread, then resolve; failed → publish an Outcome Comment naming the failed step and the needed action (from the terminal reason), leaving the thread open.
2. MUST NOT resolve any thread before the batch's Verification confirmed the result.
3. MUST publish, at Run end, an Outcome Comment on each still-unresolved Review Issue stating why it remains open and that a later Round retries it.
4. MUST make every comment idempotent via the marker convention — a retried Run or Batch never posts a duplicate.
5. MUST journal each propagation as a Run Event carrying the Review Issue reference and the action taken.
6. MUST keep a single propagation failure from aborting the batch: report it, journal it, continue with the remaining issues.

## Subtasks

- [ ] Extend the batch source-resolution step to iterate settled issues individually with per-status actions
- [ ] Build Outcome Comment bodies from artifact status, terminal reason, and duplicate-of reference, with the idempotency marker
- [ ] Add the Run-end pass for still-unresolved issues
- [ ] Journal propagation events with issue references
- [ ] Engine tests with a fake source: per-status action matrix, comment-before-resolve ordering, no-resolve-before-verify, idempotent retry, propagation-failure continuation

## Acceptance Criteria

- [ ] A settled batch with resolved, invalid, duplicated, and failed issues produces exactly the per-status actions above, in comment-before-resolve order
- [ ] Re-running propagation posts no duplicate comments
- [ ] A failed thread remains unresolved on the Review Source and carries the explanatory comment
- [ ] Run-end leaves no still-unresolved issue without a comment
- [ ] The full test suite passes

## Context

- interface: `internal/daemon/engine.go`
- interface: `internal/reviewsource/reviewsource.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/rounds/rounds.go`

## Verification

- `go test ./internal/daemon/... ./internal/reviewsource/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, User Story 6, Core Features 6–7, Decisions (propagation granularity, resolve-after-comment); `_techspec.md` → Build Order 7, Testing Approach (engine tests), Decisions (batch-settlement propagation point).

## Result

- Implemented per-issue Review Source propagation at Batch settlement in `internal/daemon/engine.go`: resolved issues resolve the thread; invalid and duplicated issues comment before resolving; failed issues comment and remain open; source propagation failures are reported, journaled, and do not abort remaining issue propagation.
- Added Outcome Comment marker/body support and issue-level Review Source operations in `internal/reviewsource/reviewsource.go` and `internal/reviewsource/coderabbit/coderabbit.go`; CodeRabbit now skips duplicate comments when the marker already exists on the thread.
- Added the Run-end unresolved propagation pass and daemon Run Event journaling with issue path, source reference, status, action, and failure details.
- Acceptance evidence: `TestResolveCyclePropagatesSettledIssueOutcomesIndividually`, `TestResolveCycleOutcomeCommentsAreIdempotent`, `TestResolveCycleSourcePropagationFailureContinues`, and `TestResolveCycleRunEndLeavesUnresolvedIssuesCommented` cover the per-status matrix, comment-before-resolve ordering, idempotency, failed-thread open state, run-end comments, and continuation on propagation failure.
- Verification passed: `go test ./internal/daemon/... ./internal/reviewsource/...` (117 tests), `go test ./internal/daemon/... ./internal/reviewsource/... ./internal/cli/...` (532 tests), `go test ./...` (1192 tests), `go build -buildvcs=false ./...`, and `make verify`.
- Verification blocker: exact task command `go build ./...` fails before compilation with `error obtaining VCS status: exit status 128`; `go build -x ./cmd/roundfix` shows Go running `git status --porcelain` from `/Users/marcio`, where `/Users/marcio/.git` is an invalid parent marker. Because the task's exact build verification does not pass, this task is settled as failed despite the implementation and repo gate passing.
