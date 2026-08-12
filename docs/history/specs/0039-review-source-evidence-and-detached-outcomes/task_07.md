---
task: task_07
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: backend
complexity: high
---

# Task 07: Inherit evidence only for Daemon artifact commits

## Overview

Prove when the current head is exactly the Daemon-created review-artifact-only
descendant of a verified parent and inherit Evidence only in that narrow case.
Real Git fixtures reject mixed paths, user-authored commits, wrong parents,
empty diffs, unresolved threads, and unrecognized artifact roots.

## Requirements

1. MUST return created commit SHA, parent SHA, and resolved review root from
   Daemon review-artifact commit creation.
2. MUST require accepted verified Evidence for the exact parent.
3. MUST require the current head to equal the exact Daemon-created artifact
   commit.
4. MUST require a non-empty parent-to-current diff wholly under the resolved
   review-artifact root.
5. MUST require zero unresolved CodeRabbit threads.
6. MUST emit `artifact_only_descendant` Evidence carrying both heads.
7. MUST fall back to normal current-head Evidence polling on any mismatch and
   MUST NOT publish another Roundfix review request for an inherited head.

## Subtasks

- [x] Return exact commit identity from artifact publication.
- [x] Prove parent, current head, and review root.
- [x] Validate non-empty artifact-only diffs.
- [x] Require verified parent Evidence and no unresolved threads.
- [x] Publish inherited Evidence with both heads.
- [x] Add real Git positive and refusal fixtures.

## Acceptance Criteria

- [x] An exact Daemon artifact-only descendant inherits verified parent
      Evidence without another Roundfix request or wait.
- [x] Mixed-path, user-authored, wrong-parent, and empty commits do not inherit.
- [x] A path outside the resolved review root refuses inheritance.
- [x] Unresolved CodeRabbit threads refuse inheritance.
- [x] Inherited Evidence names `artifact_only_descendant` and both heads.
- [x] Every refusal falls back to normal current-head polling.
- [x] Existing separate review-artifact commit policy remains intact.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/config/config.go`
- interface: `internal/config/config_test.go`
- interface: `internal/preflight/preflight.go`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/daemon_test.go`

## Verification

- `rtk go test ./internal/cli -run 'Test.*Artifact.*Evidence.*(Inherited|Mixed|Parent|Empty|User|Thread|Root)' -count=1`
  — expected: real Git fixtures accept only the exact Daemon artifact
  descendant.
- `rtk go test ./internal/config ./internal/daemon -run 'Test.*(ReviewRoot|ReviewArtifactsCommit)' -count=1`
  — expected: review-root resolution and Daemon commit identity remain stable.
- `rtk go test -race ./internal/cli -run 'Test.*Artifact.*Evidence' -count=1`
  — expected: artifact identity and Evidence inheritance are race-free.

## References

- `_prd.md` → Goal 5; User Story 8; Core Features 3–4; Success Metrics.
- `_techspec.md` → API Contracts: Artifact-only descendant; Artifact
  persistence; Build Order 7.
- `../../adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md` →
  separate
  artifact commit.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  narrowly inherited Evidence.

## Result

Watch now verifies the code head before creating its separate review-artifact
commit. Artifact publication returns the created commit SHA, its exact parent,
and the resolved review root. Inheritance then proves the current head, single
parent, Daemon commit subject, non-empty diff, and review-root-only paths
against real Git state.

The parent Evidence must be verified and bound to the exact parent before
publication. A fresh parent observation must remain verified, which confirms
that zero unresolved CodeRabbit threads remain. Successful inheritance emits
`artifact_only_descendant` Evidence with the artifact and parent heads. Any
refusal enters the existing current-head Evidence polling path, and the
artifact publisher runs at most once.

### Verification

- `rtk env GOCACHE=/tmp/roundfix-task07-gocache go test ./internal/cli -run 'Test.*Artifact.*Evidence.*(Inherited|Mixed|Parent|Empty|User|Thread|Root)' -count=1`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task07-gocache go test ./internal/config ./internal/daemon -run 'Test.*(ReviewRoot|ReviewArtifactsCommit)' -count=1`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task07-gocache go test -race ./internal/cli -run 'Test.*Artifact.*Evidence' -count=1`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task07-gocache go test ./internal/watch ./internal/cli ./internal/reviewsource ./internal/runevent ./internal/config ./internal/daemon -count=1`
  — passed.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

The initial compile-only probe used the host Go cache and was blocked by the
managed sandbox. Re-running with the task-local cache above exercised the same
packages successfully. The Daemon remains authoritative for the task's
declared Verification after this turn.

### Acceptance evidence

- `TestRunWatchArtifactEvidenceInheritedWithoutCurrentHeadPolling` creates a
  real Daemon review-artifact commit only after the exact parent is verified,
  refreshes the parent once for unresolved-thread proof, and records no
  current-artifact-head poll or wait.
- `TestReviewArtifactEvidenceMixedParentEmptyUserRootRefused` uses real Git
  repositories to reject mixed paths, a wrong recorded parent, an empty diff,
  a user-authored commit subject, and a path outside the resolved review root.
- `TestRunWatchArtifactEvidenceThreadRefusesAndFallsBack` makes the fresh
  parent observation non-verified, refuses inheritance, and records the
  subsequent current-head Evidence poll.
- `TestWatchArtifactEvidenceMixedFallsBackToCurrentHeadPolling` proves every
  non-inherited artifact publication shares one normal current-head polling
  branch and does not invoke artifact publication again.
- The inherited `daemon.review_status` payload contains
  `artifact_only_descendant`, the artifact SHA as expected and observed head,
  and the verified parent SHA; the terminal outcome carries the artifact and
  verified parent heads separately.
- The review-root and Daemon commit-message contract tests pass, and the
  affected package suites preserve the separate ADR-0036 artifact commit
  policy.
