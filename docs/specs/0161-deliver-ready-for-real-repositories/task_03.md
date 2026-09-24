---
task: task_03
spec: 0161-deliver-ready-for-real-repositories
status: pending
type: backend
complexity: medium
---

# Task 03: Clean parks, crash-safe archive, stale owners

## Overview

A non-exact archive park leaves the checkout dirty so every later item refuses to start; a crash after the archive commit resumes as `review-stale`; a reused owner PID makes both `resume` and `stop` refuse.

## Requirements

1. MUST, on park, discard uncommitted changes the item made and return the checkout to the branch it started from.
2. MUST accept on resume a HEAD whose parent is the reviewed head and whose change is exactly the Spec move into the archive root.
3. MUST release, on `resume`, an owner whose live process identity does not match the recorded identity.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, the next item starts after a non-exact archive park.
- [ ] A crash after the archive commit resumes past archiving.
- [ ] A reused PID with another identity is released.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/delivery/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAParkLeavesACleanCheckout|TestResumeAcceptsTheArchiveCommit|TestResumeReleasesAStaleOwner)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAParkLeavesACleanCheckout TestResumeAcceptsTheArchiveCommit TestResumeReleasesAStaleOwner; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Park and resume
