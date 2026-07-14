---
task: task_03
spec: 0027-review-loop-integrity
status: completed
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

- [x] Extend the GitHub client interface and the `gh`-backed implementation with both mutations
- [x] Implement the marker convention and the marker-detection helper
- [x] Extend the existing fakes so downstream tests can assert calls and bodies
- [x] Unit-test both mutations against a scripted `gh` runner, including failure paths and marker detection

## Acceptance Criteria

- [x] A reply and a PR comment each produce exactly one correctly-shaped `gh` invocation in tests
- [x] The marker helper detects an existing marker and reports absence correctly
- [x] Failures name the operation and wrap the cause
- [x] The full test suite passes

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`
- interface: `internal/reviewsource/reviewsource.go`

## Verification

- `grep -q "ReplyToReviewThread" internal/reviewsource/coderabbit/coderabbit.go` — expected: exit 0
- `grep -q "CommentOnPullRequest" internal/reviewsource/coderabbit/coderabbit.go` — expected: exit 0
- `go test ./internal/reviewsource/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, User Stories 4 and 6, Core Features 4 and 7; `_techspec.md` → Build Order 3, Interfaces (GitHubClient), Integration Points (GitHub via gh).

## Result

Implemented the Review Source write surface without wiring it into any command. `coderabbit.GitHubClient` and `GHClient` now support replying to a review thread and commenting on a pull request through the existing `gh` runner pattern, and failures wrap the underlying `gh` error with the operation name. Roundfix-authored comment bodies now have a stable marker-line helper plus marker detection for already-fetched comment bodies.

Evidence:

- Pre-change signal: `grep -q "ReplyToReviewThread" internal/reviewsource/coderabbit/coderabbit.go` and `grep -q "CommentOnPullRequest" internal/reviewsource/coderabbit/coderabbit.go` exited 1 before implementation.
- Red test: `go test ./internal/reviewsource/coderabbit -run 'TestGHClientWriteMutationsInvokeGHOnce|TestGHClientWriteMutationFailuresWrapCause|TestRoundfixCommentMarkerHelpers'` failed before implementation because the new methods, marker helpers, mutation constant, and testable `runGH` seam did not exist.
- `go test ./internal/reviewsource/...`: passed, 27 tests across 2 packages.
- `grep -q "ReplyToReviewThread" internal/reviewsource/coderabbit/coderabbit.go`: passed.
- `grep -q "CommentOnPullRequest" internal/reviewsource/coderabbit/coderabbit.go`: passed.
- `go build -buildvcs=false ./...`: passed.
- `make verify`: passed; `go test ./...` reported 1157 tests across 19 packages, `roundfix skills check` passed, and the Makefile build completed with `-buildvcs=false`.

Acceptance evidence:

- `TestGHClientWriteMutationsInvokeGHOnce` asserts one `gh` invocation for `ReplyToReviewThread` and one for `CommentOnPullRequest`, including the GraphQL mutation name and REST endpoint shape.
- `TestRoundfixCommentMarkerHelpers` asserts marker construction, marker appending, marker detection, and absence when the marker is not a standalone line.
- `TestGHClientWriteMutationFailuresWrapCause` asserts both new write failures name their operation and wrap the underlying cause with `errors.Is`.
- Full suite: `make verify` passed after the implementation.

Verification note: `rg` is not installed in this execution environment, so the task-local field checks use `grep`. The build command uses `-buildvcs=false`, matching the repository Makefile, because bare `go build ./...` fails in this Roundfix worktree when Go VCS stamping probes the invalid parent `/Users/marcio/.git`.
