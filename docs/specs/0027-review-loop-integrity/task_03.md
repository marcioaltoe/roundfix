---
task: task_03
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: medium
---

# Task 03: Add thread-reply and PR-comment mutations to the Review Source client

## Overview

The Review Source client can only resolve threads today. This task widens the GitHub boundary with two write operations — reply to a review thread and comment on a pull request — plus the idempotency-marker convention that later propagation and bypass-audit work will rely on. Fully testable against the existing fake GitHub client seam.

## Requirements

1. MUST add a reply-to-review-thread operation (thread id + body) to the GitHub client interface and its concrete implementation, using the same external `gh` runner pattern as the existing resolve mutation.
2. MUST add a comment-on-pull-request operation (PR number + body) the same way.
3. MUST define an idempotency-marker convention: every Roundfix-authored comment body carries a stable machine-checkable marker line, and a helper reports whether a marker already appears in already-fetched thread or PR comment data.
4. MUST surface failures with the operation name and the underlying error, per the CLI error contract.
5. MUST NOT invoke the new operations from any command in this task — wiring happens in later work.

## Subtasks

- [ ] Extend the GitHub client interface and the `gh`-backed implementation with both mutations
- [ ] Implement the marker convention and the marker-detection helper
- [ ] Extend the existing fakes so downstream tests can assert calls and bodies
- [ ] Unit-test both mutations against a scripted `gh` runner, including failure paths and marker detection

## Acceptance Criteria

- [ ] A reply and a PR comment each produce exactly one correctly-shaped `gh` invocation in tests
- [ ] The marker helper detects an existing marker and reports absence correctly
- [ ] Failures name the operation and wrap the cause
- [ ] The full test suite passes

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/reviewsource.go`

## Verification

- `rg -q "ReplyToReviewThread" internal/reviewsource/coderabbit/coderabbit.go` — expected: exit 0
- `rg -q "CommentOnPullRequest" internal/reviewsource/coderabbit/coderabbit.go` — expected: exit 0
- `go test ./internal/reviewsource/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 3, User Stories 4 and 6, Core Features 4 and 7; `_techspec.md` → Build Order 3, Interfaces (GitHubClient), Integration Points (GitHub via gh).
