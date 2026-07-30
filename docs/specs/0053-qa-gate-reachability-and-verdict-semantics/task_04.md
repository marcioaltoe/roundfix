---
task: task_04
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: completed
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

## Result

### Implementation

- Branch Integrity Preflight now attributes pending Run Branches to their
  recorded Implement Runs and calls task_03's shared
  `worktree.QAReportOnlyBranch` probe with the current target head, Run Branch,
  and Spec slug.
- A proven QA-report-only branch is excluded from the automatic integration
  plan and causes a refusal whose branch entry directs the operator to
  `roundfix reconcile --apply`. Task-work branches retain the existing
  fast-forward integration and `git merge --ff-only` refusal guidance.
- A false or errored probe remains task work. The preflight does not convert
  missing or unreadable proof into supersession evidence.

### Focused checks

- Red signal before the production change:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260730T131049Z_a02854137f3dd85c/.gocache go test ./internal/cli -run '^TestBranchIntegrityPreflight(ClassifiesPendingRunBranches|ListsTaskWorkAndSupersededQAReportSeparately)$'`
  — failed because the QA-report-only branch was integrated and the mixed
  refusal omitted its reconcile guidance.
- After implementation:
  `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260730T131049Z_a02854137f3dd85c/.gocache go test ./internal/cli -run '^(TestBranchIntegrityPreflight|TestBranchIntegrityIntegrationPlan)'`
  — passed, covering the new classification table, mixed refusal, and the
  existing refusal, automatic integration, migration, and Active Run cases.
- An initial focused-check invocation with relative `GOCACHE=.gocache` did not
  start the tests because Go requires an absolute cache path; the commands
  above use the writable absolute repository cache.
- `rtk git diff --check` — passed after the final implementation and Result
  edits.
- The commands in `## Verification` were not run; Daemon Verification owns
  them.

### Acceptance evidence

- QA-report-only branch: `TestBranchIntegrityPreflightClassifiesPendingRunBranches/QA_report_only`
  proves the CLI returns the preflight refusal, performs no integration, emits
  the reconcile guidance, and omits a merge command for that branch.
- Task work: the same table's `task_work` case proves a false probe preserves
  the existing successful automatic fast-forward integration path.
- Unprovable probe: the table's `unprovable_probe` case proves a probe error
  also preserves the task-work integration path.
- Mixed classes:
  `TestBranchIntegrityPreflightListsTaskWorkAndSupersededQAReportSeparately`
  proves task work retains its `git merge --ff-only` entry while the
  QA-report-only branch receives only the superseded-report reconcile
  guidance, with no integration attempt during the refusal.
