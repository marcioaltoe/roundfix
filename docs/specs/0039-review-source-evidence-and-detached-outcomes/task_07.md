---
task: task_07
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
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

- [ ] Return exact commit identity from artifact publication.
- [ ] Prove parent, current head, and review root.
- [ ] Validate non-empty artifact-only diffs.
- [ ] Require verified parent Evidence and no unresolved threads.
- [ ] Publish inherited Evidence with both heads.
- [ ] Add real Git positive and refusal fixtures.

## Acceptance Criteria

- [ ] An exact Daemon artifact-only descendant inherits verified parent
      Evidence without another Roundfix request or wait.
- [ ] Mixed-path, user-authored, wrong-parent, and empty commits do not inherit.
- [ ] A path outside the resolved review root refuses inheritance.
- [ ] Unresolved CodeRabbit threads refuse inheritance.
- [ ] Inherited Evidence names `artifact_only_descendant` and both heads.
- [ ] Every refusal falls back to normal current-head polling.
- [ ] Existing separate review-artifact commit policy remains intact.

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
