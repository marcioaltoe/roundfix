---
task: task_05
spec: 0027-review-loop-integrity
status: completed
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

- [x] Build the branch-integrity inspection composing the pending-work enumeration and the Active-Run store lookups
- [x] Implement fast-forward auto-integration with per-branch reporting during preflight
- [x] Add the flag to the shared command parsing and the request struct for all three commands
- [x] Implement the bypass audit comment (content, marker, publish-or-fail rule)
- [x] Write the refusal messages following the preflight-failure convention
- [x] Integration-test refusal, auto-integration, active-run block, bypass success, and bypass publish-failure via buffer-captured CLI runs with fakes

## Acceptance Criteria

- [x] With a diverged Run Branch present, fetch, resolve, and watch each refuse with exit 2, naming the branch and the integration command, before creating any Run artifacts
- [x] With a fast-forwardable Run Branch present, the command integrates it, reports the integration, and proceeds
- [x] With an Active Run on the target, the command refuses naming the run id and stop command
- [x] With the bypass flag, the audit comment body contains the marker, run id, skipped guardrails, and ignored state; a failing publish fails the command
- [x] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/preflight/preflight.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/store/store.go`
- interface: `internal/reviewsource/coderabbit/coderabbit.go`

## Verification

- `grep -R -q "skip-branch-integrity" internal/cli` — expected: exit 0 (flag exists; `rg` is unavailable in this execution environment)
- `go test ./internal/cli/... ./internal/worktree/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build (matches the repository Makefile in this worktree)
- `go run -buildvcs=false ./cmd/roundfix watch --help 2>&1 | grep -q "skip-branch-integrity"` — expected: exit 0 (help documents the flag; `rg` is unavailable in this execution environment)

## References

`_prd.md` → Goal 1, User Stories 2, 3, 4, Core Features 2, 3, 4; `_techspec.md` → Build Order 4, Interfaces (BranchIntegrityReport), API Contracts (new flag, preflight failure text); ADR-0042.

## Result

Implemented the Branch Integrity Preflight wiring for fetch, resolve, and watch. Each command now inspects pending Run Branch work and Active Runs after standard PR preflight and before Review Source fetches, Agent Sessions, comments, or code changes. Fast-forwardable pending work is integrated with a stderr report and journaled after Run creation; non-fast-forward pending work refuses with exit 2, branch/worktree/ahead count, and `git merge --ff-only <branch>`. Active Run refusals name the run id plus `roundfix stop <id>` and `roundfix stop --force <id>`.

Added `--skip-branch-integrity` to fetch, resolve, and watch. The bypass records the ignored pending branches and Active Run state in a PR audit comment using the Roundfix idempotency marker, then journals the bypass before any fetch, Agent Session, or code change. If an Active Run is ignored, the bypass-created review Run is created without acquiring the target lock so the existing Active Run remains the lock owner. If comment publishing fails, the command exits 2, marks the created Run failed, and does not fetch, start an Agent, commit, or push.

Evidence:

- Pre-change signal: `rtk grep -R "skip-branch-integrity" internal/cli` exited 1 before implementation.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache go test ./internal/cli -run 'TestBranchIntegrity|TestReviewCommandHelpDocumentsSkipBranchIntegrity'`: passed.
- `rtk grep -R -q "skip-branch-integrity" internal/cli`: passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache go test ./internal/cli/... ./internal/worktree/... ./internal/store/...`: passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache go build ./...`: failed before compilation because this worktree cannot provide Go VCS stamping metadata (`error obtaining VCS status: exit status 128`).
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache go build -buildvcs=false ./...`: passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache go run -buildvcs=false ./cmd/roundfix watch --help 2>&1 | rtk grep -q "skip-branch-integrity"`: passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache make verify`: passed; reported 1176 Go tests across 19 packages, `roundfix skills check` passed, and the Makefile build completed with `-buildvcs=false`.

Acceptance evidence:

- `TestBranchIntegrityPreflightRejectsPendingRunBranchForReviewCommands` covers fetch, resolve, and watch refusing exit 2 before Run creation/fetch when a diverged Run Branch is present, including branch name, ahead count, and integration command.
- `TestBranchIntegrityPreflightIntegratesFastForwardRunBranchAndJournals` covers fast-forward auto-integration, stderr reporting, command continuation, and the Run Event payload.
- `TestBranchIntegrityPreflightRejectsActiveRunForReviewCommands` covers fetch, resolve, and watch refusing an Active Run with the run id and both stop commands.
- `TestBranchIntegrityBypassPublishesAuditBeforeFetch` covers marker, run id, skipped guardrails, ignored pending state, and bypass Run Event journaling; `TestBranchIntegrityBypassAuditsActiveRunAndProceeds` covers ignored Active Run state and continuation; `TestBranchIntegrityBypassFailsWhenAuditCommentPublishFails` covers publish failure exiting 2 before fetch.
- Full verification: `make verify` passed.
