---
task: task_02
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: pending
type: backend
complexity: medium
---

# Task 02: The GitHub boundary reads the workflow runs and the Expected Workflows of the head

## Overview

The checking stage reads only `gh pr checks`, which on Fluxus pull requests
#466, #469 and #477 listed one passing check while two workflow runs of the
head were already in progress (reported by fluxus on 2026-11-05,
[references](references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md);
timing in [references](references/2026-11-06-fluxus-check-timing.md)). After
this Task `GitHubCLI.CurrentHeadChecks` also returns the GitHub Actions
workflow runs of the head commit and the head's Expected Workflows in the
`CheckReport`, as `_techspec.md` → API Contracts 1 and 2 state. The engine
does not read the new fields yet (task_03), so delivery behaves as today.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestPullRequestBoundaryReportsCurrentHeadChecks|TestNoChecksReportedIsAnEmptyReport|TestPullRequestBoundaryReportsFailedCurrentHeadChecks|TestPullRequestBoundaryRejectsChecksFromAMovedHead|TestPullRequestReadCarriesTheMergeState)$' ./internal/delivery`
   and record in the Result that it passes.
2. MUST create `internal/delivery/head_check_set.go` with `WorkflowRun` and
   `ExpectedWorkflows` exactly as `_techspec.md` → Interfaces and API
   Contract 2 state, parsing with `gopkg.in/yaml.v3` nodes so that the key
   `on` is read as written, never as a boolean.
3. MUST change `internal/delivery/github.go`: add `baseRefName` to
   `pullRequestJSONFields` (giving exactly
   `number,url,state,baseRefName,headRefName,headRefOid,mergedAt,mergeCommit,mergeable,mergeStateStatus`),
   `PullRequest.BaseBranch`, the `WorkflowRuns` and `ExpectedWorkflows`
   fields of `CheckReport`, and the two reads of API Contract 1 in
   `CurrentHeadChecks`, in that order, with the error texts it states.
4. MUST create `internal/delivery/head_check_set_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, each
   top-level test calling `t.Parallel()` first, with
   `TestExpectedWorkflowsFollowTheTriggerRule`: one call of
   `ExpectedWorkflows` over these files with base `main`, returning exactly
   `[a.yml c.yml e.yml f.yml g.yml h.yml k.yml]`:
   - `a.yml` `on: {pull_request: {branches: [main], types: [opened, reopened, edited, synchronize]}}`;
   - `b.yml` `pull_request` with `branches: [main]` and `paths: ['x/**']`;
   - `c.yml` `on: [push, pull_request]`;
   - `d.yml` `push` with `branches: [main]`;
   - `e.yml` the key written `"on"` with an empty `pull_request`;
   - `f.yml` `pull_request` with `types: [opened, reopened, synchronize]` and `branches: [main]`, and `push` with `branches: [main]`;
   - `g.yml` `on: pull_request`;
   - `h.yml` an empty `push`;
   - `i.yml` `pull_request` with `branches: ['releases/**']`;
   - `j.yml` `pull_request` with `types: [opened]`;
   - `k.yml` `pull_request_target` with `branches: main`;
   - `l.yml` `pull_request` with `branches-ignore: [docs]`;
   - `m.yml` `pull_request` with `branches: [main, '!main-old']`;
   - `n.yml` `on: workflow_dispatch`;
   - `o.yml` the text `on: [` (not YAML);
   - `p.yml` `pull_request` with `paths-ignore: ['docs/**']`.
5. MUST add to `internal/delivery/github_test.go`
   `TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead`: a scripted runner
   with `gh pr view 18` (head `head-two`, base `main`), `gh pr checks 18`
   (one passing `Validate PR title` of `CI | Conventions` linked to run 1),
   `gh api repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100`
   answering `total_count` 4 and the runs, in this order: id 9 `CI | Backend`
   `.github/workflows/ci-backend.yml` `pull_request` `completed` `cancelled`
   run_number 9; id 10 of the same file and event `in_progress` run_number 10
   attempt 1; id 11 `CI | Backend` with `head_sha` `head-one`; id 1
   `CI | Conventions` `.github/workflows/ci-conventions.yml` `pull_request`
   `completed` `success`; `git ls-tree -z --name-only head-two -- .github/workflows/`
   listing `.github/workflows/ci-backend.yml`, `.github/workflows/ci-conventions.yml`
   and `.github/workflows/README.md`; `git show` of the two YAML files (the
   backend one with a `paths` filter, the conventions one with `branches:
   [main]` and the four Fluxus types); and `gh pr view 18` again. The report
   carries exactly the runs `10` (`in_progress`, attempt 1, its `html_url`)
   and `1`, in that order, and the Expected Workflows
   `[.github/workflows/ci-conventions.yml]`. Subtests: `incomplete page`
   (`total_count` 101 with one run) returns an error containing `more than
   one page`; `runs read fails` (exit 1, stderr `HTTP 403`) returns an error
   starting `read workflow runs:`; `ls-tree fails` returns an error starting
   `list workflow files:`.
6. MUST update every test that scripts the commands of
   `GitHubCLI.CurrentHeadChecks` on the default branch when the Task runs, so
   that after the `gh pr checks` step it gains the `gh api` step answering
   `{"total_count":0,"workflow_runs":[]}` and the `git ls-tree` step with an
   empty listing, and nothing else changes (`_techspec.md` → Declared
   breaks), and make the fake `gh` of
   `TestCommandDeliveryEngineRecordsThePullRequest`
   (`internal/cli/deliver_record_pull_request_test.go`) answer `api` with
   `{"total_count":0,"workflow_runs":[]}`.
7. MUST run sabotages 4 and 5 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails and on which input.
8. MUST NOT change `internal/delivery/engine.go`, `internal/delivery/turn.go`,
   `internal/delivery/park_class.go`, any blocker or log text, any document,
   the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the scripted boundary tests.
- [ ] Add `WorkflowRun` and `ExpectedWorkflows` with their trigger rule.
- [ ] Read the workflow runs and the workflow files in `CurrentHeadChecks`.
- [ ] Add the tests, update the scripted ones and run the two sabotages.

## Acceptance Criteria

- [ ] `ExpectedWorkflows` returns exactly the filter-free triggers of API Contract 2.
- [ ] `CurrentHeadChecks` returns the latest run per workflow file and event of the head and the Expected Workflows, and fails with the stated errors.
- [ ] Every existing engine and `deliver` test passes unchanged apart from the scripted steps Declared breaks names.

## Context

- creates: `internal/delivery/head_check_set.go`
- creates: `internal/delivery/head_check_set_test.go`
- interface: `internal/delivery/github.go`
- interface: `internal/delivery/github_test.go`
- interface: `internal/delivery/merge_state_test.go`
- interface: `internal/delivery/check_rerun.go`
- interface: `internal/cli/deliver_record_pull_request_test.go`

## Consumers

- `CheckReport`, `PullRequest` and the `gh pr view` field list — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestExpectedWorkflowsFollowTheTriggerRule TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run '^TestCommandDeliveryEngineRecordsThePullRequest$' ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/delivery && GOOS=linux go build -buildvcs=false ./internal/delivery` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards the trigger rule and the reader behave as API Contracts 1 and 2 state, every existing test of the package passes, the command test answers the new read, and both platforms build.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1; Current behavior 4; Interfaces; API Contract 1; API Contract 2; Declared breaks; Testing Approach; Build Order 2; ADR-0369
