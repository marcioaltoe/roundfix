---
task: task_04
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: pending
type: backend
complexity: medium
---

# Task 04: Keep superseded branches out of automatic integration

## Overview

Branch Integrity Preflight fast-forwards pending Run Branches onto the user's
head branch. A branch holding only a failing QA report gets merged silently,
dragging that report into the target. Exclude QA-report-only branches from
automatic integration and point the user at the explicit release instead.

## Requirements

1. MUST exclude a pending Run Branch that the shared probe proves QA-report-only
   from automatic fast-forward integration.
2. MUST list those branches separately in the refusal, with
   `superseded QA report — release with: roundfix reconcile --apply` rather than
   a `git merge --ff-only` command.
3. MUST keep genuinely pending task-work branches on their current listing and
   integration behavior.
4. MUST reuse the probe from task_03 rather than reimplementing the proof.
5. MUST treat an unprovable probe as task work, preserving today's behavior.

## Subtasks

- [ ] Consume the probe in the preflight and split the two branch classes.
- [ ] Render the supersession guidance in the refusal.
- [ ] Cover both classes and the unprovable case in the preflight table.

## Acceptance Criteria

- [ ] A pending QA-report-only branch is not fast-forwarded and appears with the
      reconcile guidance.
- [ ] A pending branch with task work is fast-forwarded exactly as today.
- [ ] A branch whose probe cannot be proven is treated as task work.
- [ ] A run with both classes lists each under its own guidance.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/worktree/worktree.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/` — expected: pass, including the mixed-class
  preflight case.

## References

`_prd.md` → Goal 3 Story 4, Feature 7; `_techspec.md` → Build Order 4, API
Contracts (Branch Integrity Preflight), Risks (auto-integration behavior
change).
