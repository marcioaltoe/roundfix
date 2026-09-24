---
task: task_02
spec: 0156-a-delivery-loop-that-outlives-the-session
status: pending
type: backend
complexity: medium
---

# Task 02: The pull request boundary

## Overview

No code in Roundfix opens, follows or merges a pull request today. Add that boundary behind an interface, backed by the GitHub CLI already authenticated on the host, with a fake for tests.

## Requirements

1. MUST push a branch and report the remote head.
2. MUST find an open pull request for a head branch before creating one, and create one only when none exists.
3. MUST report the current-head checks of a pull request.
4. MUST report an existing merge instead of merging again.
5. MUST provide a fake implementation for tests, and MUST NOT call GitHub from any test.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An existing pull request for the head is returned instead of a new one.
- [ ] An already-merged pull request is reported, not merged again.
- [ ] No test reaches the network.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement.go`
- creates: `internal/delivery/github.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestPullRequestBoundaryReusesAnOpenPullRequest$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPullRequestBoundaryReusesAnOpenPullRequest"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestPullRequestBoundaryDoesNotMergeTwice$" ./internal/delivery 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPullRequestBoundaryDoesNotMergeTwice"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Components
