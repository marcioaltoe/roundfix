---
task: task_05
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: high
---

# Task 05: Enforce the Branch Integrity Preflight on fetch, resolve, and watch

## Overview

Wire the Branch Integrity Preflight into all three review commands: before any Run work, pending Run Branch work must be zero (auto-integrated when fast-forward resolves it, refused otherwise) and no Active Run may be bound to the target. A single explicit bypass flag skips both guardrails, but only after publishing an audit comment on the pull request. This delivers the deterministic guardrail end to end; review Runs still execute in their worktree until the next task removes it.

## Requirements

1. MUST run the Branch Integrity Preflight in fetch, resolve, and watch before any Review Source fetch, Agent Session, comment, or code change: enumerate pending Run Branch work against the PR Head Branch, fast-forward-integrate what resolves cleanly (reporting each integration), and fail Preflight Validation otherwise, naming each pending branch, its worktree path, its ahead-commit count, and the exact integration command.
2. MUST fail Preflight Validation while any Active Run is bound to the target, naming the run id and the stop command, including the force form for the case where the owning process may be dead.
3. MUST add a `--skip-branch-integrity` flag to all three commands that skips both guardrails; when used, Roundfix publishes a pull request comment — carrying the idempotency marker — recording who and when, the run id, which guardrails were skipped, and the ignored state (pending branches and active runs enumerated). The comment publishes after Run creation but before any fetch, Agent Session, or code change.
4. MUST fail the command when the bypass audit comment cannot be published — the audit trail is part of the contract.
5. MUST surface refusals through the existing preflight-failure convention: reason, "No side effects" statement, next action, exit code 2.
6. MUST journal auto-integrations and bypasses as Run Events once a Run exists.

## Subtasks

- [ ] Build the branch-integrity inspection composing the pending-work enumeration and the Active-Run store lookups
- [ ] Implement fast-forward auto-integration with per-branch reporting during preflight
- [ ] Add the flag to the shared command parsing and the request struct for all three commands
- [ ] Implement the bypass audit comment (content, marker, publish-or-fail rule)
- [ ] Write the refusal messages following the preflight-failure convention
- [ ] Integration-test refusal, auto-integration, active-run block, bypass success, and bypass publish-failure via buffer-captured CLI runs with fakes

## Acceptance Criteria

- [ ] With a diverged Run Branch present, fetch, resolve, and watch each refuse with exit 2, naming the branch and the integration command, before creating any Run artifacts
- [ ] With a fast-forwardable Run Branch present, the command integrates it, reports the integration, and proceeds
- [ ] With an Active Run on the target, the command refuses naming the run id and stop command
- [ ] With the bypass flag, the audit comment body contains the marker, run id, skipped guardrails, and ignored state; a failing publish fails the command
- [ ] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/preflight/preflight.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/store/store.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`

## Verification

- `rg -q "skip-branch-integrity" internal/cli` — expected: exit 0 (flag exists)
- `go test ./internal/cli/... ./internal/worktree/...` — expected: all tests pass
- `go build ./...` — expected: clean build
- `go run -buildvcs=false ./cmd/roundfix watch --help 2>&1 | rg -q "skip-branch-integrity"` — expected: exit 0 (help documents the flag)

## References

`_prd.md` → Goal 1, User Stories 2, 3, 4, Core Features 2, 3, 4; `_techspec.md` → Build Order 4, Interfaces (BranchIntegrityReport), API Contracts (new flag, preflight failure text); ADR-0042.
