---
task: task_02
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: completed
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


## Result

Skills loaded: implement-task, coding-guidelines, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, golang-cli, evidence-gate, systematic-debugging, no-workarounds, tech-writer.

The GitHub boundary now returns the head's workflow runs and Expected
Workflows alongside the existing checks. It reads runs before listing the
head's workflow files, reads only `.yml` and `.yaml` entries in listing order,
and rechecks the pull request head after those reads. The pull request field
list includes `baseRefName`, and `PullRequest.BaseBranch` trims that value.
Run selection discards other heads, groups by workflow path and event,
replaces a superseded run in its original position, and preserves the
attempt and URL while lower-casing status and conclusion.

`ExpectedWorkflows` parses YAML nodes, preserving quoted and unquoted `on`,
and sorts the paths admitted by API Contract 2. Filtered, malformed and
unsupported trigger declarations create no expectation. The engine, re-run
implementation, turn records, park classification and blocker/log text were
not changed. Existing scripted reads gained only the two declared empty
steps, and the command test's fake forge now answers `api` with empty runs.

Focused evidence:

- Before any edits, the required characterization command passed (exit 0):
  `go test -count=1 -timeout 20m -v -run '^(TestPullRequestBoundaryReportsCurrentHeadChecks|TestNoChecksReportedIsAnEmptyReport|TestPullRequestBoundaryReportsFailedCurrentHeadChecks|TestPullRequestBoundaryRejectsChecksFromAMovedHead|TestPullRequestReadCarriesTheMergeState)$' ./internal/delivery`.
  All five named tests reported PASS.
- The new tests first failed to compile on the absent `WorkflowRun`,
  `ExpectedWorkflows`, `CheckReport` fields and `PullRequest.BaseBranch`.
- Acceptance criterion 1: after implementation,
  `TestExpectedWorkflowsFollowTheTriggerRule` returned exactly
  `[a.yml c.yml e.yml f.yml g.yml h.yml k.yml]` from the one required
  sixteen-file input. `TestExpectedWorkflowsSkipUnsupportedTriggers` also
  passed, covering empty mappings, scalar events, missing/non-mapping YAML
  roots, invalid branch/type values and unsupported keys.
- Acceptance criterion 2: `TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead`
  passed with exactly runs 10 and 1 and only
  `.github/workflows/ci-conventions.yml` expected. Its incomplete-page,
  runs-read, listing, file-read, malformed JSON and runner-cancellation
  subtests passed; cancellation remains discoverable with `errors.Is`.
  `TestCurrentHeadChecksKeepsWorkflowEventsSeparate` passed, covering
  independent events and files with identical names, descending run numbers,
  foreign heads, lowercase normalization, `.yaml` reads and a trimmed base.
- The final focused boundary command passed (exit 0):
  `go test -count=1 -timeout 20m -v -run '^(TestExpectedWorkflows.*|TestCurrentHeadChecks.*|TestPullRequestReadTrimsTheBaseBranch|TestPullRequestBoundary.*|TestNoChecksReportedIsAnEmptyReport|TestCheckCommandFailuresRemainErrors|TestPullRequestReadCarriesTheMergeState|TestGitHubCLIReadsTheMergeableState)$' ./internal/delivery`.
- Acceptance criterion 3, focused regression evidence: the following
  existing engine/recovery checks passed without assertion changes (exit 0):
  `go test -count=1 -timeout 20m -run '^(TestGitHubCLI.*|TestDeliveryEngine.*|TestDeliveryWaitsThroughUnreportedChecks|TestCheckReadErrorsAreRetriedUntilTheDeadline|TestPersistentCheckReadErrorsParkAtTheDeadline|TestChecksTimeoutParksTheItem|TestAConflictingPullRequestStopsTheCheckWaitAtOnce|TestAnUnknownMergeableStateKeepsWaiting|TestABlockedMergeStateParksAtTheChecksTimeout|TestTheQueueWaitsThroughABlockedMergeState|TestARetryAfterTheDefaultBranchMovedRerunsTheStaleCheckAndMerges|TestAStaleCheckKeepsPollingWhileGitHubReportsTheAttemptItReran|TestAStaleRerunThatFailsOnTheCurrentTipGetsTheOutsideChangeRerun|TestAStaleCheckWhoseRerunIsRefusedParksChecksFailed|TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce|TestCancelledChecksAreNotRerun|TestCheckRerunRestartsTheTimeout|TestAnOldCheckAttemptCannotKeepResettingTheTimeout)$' ./internal/delivery`.
- Focused CLI archive-record regressions passed (exit 0):
  `go test -count=1 -timeout 20m -run '^(TestDeliveryWorkflow.*|TestDeliveryRecordProofAllowsOnlyTheRenderedOutcomeShortening)$' ./internal/cli`.
  The declared command-engine test and full-package/platform checks remain
  for Daemon Verification; these focused checks do not claim their outcome.

Sabotage evidence (temporary repositories made from `git archive HEAD` with
the six implementation/test files overlaid; the Run checkout was not
sabotaged):

- Sabotage 4 accepted `paths` in the pull request trigger. In
  `/private/tmp/roundfix-0359-task02-sabotage-ucqi6xn6/4`,
  `go test -count=1 -timeout 20m -v -run '^TestExpectedWorkflowsFollowTheTriggerRule$' ./internal/delivery`
  exited 1 and failed `TestExpectedWorkflowsFollowTheTriggerRule` on
  `b.yml`, `branches: [main], paths: ['x/**']`: the actual list incorrectly
  included `b.yml`. Full transcript:
  `/private/tmp/roundfix-0359-task02-sabotage-ucqi6xn6/sabotage-4.log`.
- Sabotage 5 kept the first run for an already-seen workflow path/event.
  In `/private/tmp/roundfix-0359-task02-sabotage-ucqi6xn6/5`,
  `go test -count=1 -timeout 20m -v -run '^TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead$' ./internal/delivery`
  exited 1 and failed
  `TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead/reads_head_runs`:
  the cancelled backend run 9 preceding run 10 was retained instead of the
  in-progress run 10 with attempt 1 and its URL. Full transcript:
  `/private/tmp/roundfix-0359-task02-sabotage-ucqi6xn6/sabotage-5.log`.

`git diff --check` passed. No Task status, Task Graph, other Task file,
protected tooling, dependency manifest or unrelated document was edited by
this Agent. No commit, push or pull request was made. The declared
Verification commands and repository gates were not run.
