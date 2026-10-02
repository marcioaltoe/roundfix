---
task: task_03
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: completed
type: backend
complexity: high
---

# Task 03: The queue merges only on GitHub's mergeable state and resumes a delivery-error retry

## Overview

The checking stage moves to merge when every check `gh pr checks` lists
passes, so a re-run that has not registered yet lets the queue ask GitHub to
merge a Pull Request its branch policy still blocks, and GitHub's refusal
parks `delivery-error`. This Task reads GitHub's merge state beside the
checks and waits while it is blocked, returns a base-branch policy refusal to
checking once per head, pins that a retry of a `delivery-error` park resumes
checking, and prints a refused retry as a refusal instead of a usage error.

## Requirements

1. MUST add `mergeStateStatus` to the Pull Request read and carry it as
   `PullRequest.MergeState` and `CheckReport.MergeState`, as "The merge state"
   and the TechSpec's Interfaces state.
2. MUST keep an otherwise mergeable item in `checking` while the merge state
   is `BLOCKED` or `UNKNOWN`, under the existing deadline, whose expiry parks
   `checks-timeout`; an empty merge state MUST keep today's decision (API
   Contract 2). The stage MUST log the merge state once per change.
3. MUST return `MergePolicyRefusalError` from the merge when `gh`'s standard
   error contains `base branch policy prohibits the merge`, case-insensitively,
   and MUST make the merging stage return the first such refusal per
   `<slug>@<head>` per owner process to `checking`, with the log line of "The
   merge refusal", and park a second as `delivery-error` (API Contract 3).
4. MUST pin, by a test that parks an item `delivery-error` through a refused
   merge and retries it, that the retry of a `delivery-error` park on an
   archived candidate at an unchanged head with a recorded Pull Request
   resumes `checking` and merges exactly once, and MUST remove any refusal
   that reproduction exposes on that path ("The retry of a `delivery-error`
   park").
5. MUST print an `Engine.Retry` refusal of `roundfix deliver retry` as Surface
   Transcript 1, with no `Usage:` block and exit `2`, and MUST keep the
   `Preflight failed` form with its usage hint for argument, configuration and
   store errors (API Contract 1).
6. MUST add the tests named in Verification: the `internal/delivery` ones in
   the new file `internal/delivery/merge_state_test.go`, over the scripted
   command runner and the fake Pull Request boundary with the fake clock, and
   the two retry-output tests in `internal/cli/deliver_retry_test.go`. Every
   existing `internal/delivery` test MUST pass unedited.

## Subtasks

- [ ] Read and carry the merge state.
- [ ] Wait in checking on a blocked or unknown merge state.
- [ ] Type a policy refusal and return it to checking once per head.
- [ ] Pin the `delivery-error` retry resumption.
- [ ] Print a refused retry as a refusal.

## Acceptance Criteria

- [ ] With every listed check passing and the merge state `BLOCKED`, no merge
      is issued until the state is `CLEAN`, and then exactly one.
- [ ] `UNKNOWN` waits like `BLOCKED`, a state that stays `BLOCKED` parks
      `checks-timeout`, and an empty state merges as before.
- [ ] A merge refused with `the base branch policy prohibits the merge`
      returns the item to `checking` and the next clean read merges with no
      park; a second refusal at the same head parks `delivery-error`.
- [ ] A `delivery-error` park at an unchanged archived head with a Pull
      Request, retried, resumes `checking` and merges once.
- [ ] A refused retry prints Surface Transcript 1's standard error with no
      `Usage:` line and exits `2`; an unexpected argument still prints
      `Preflight failed` and `Usage:`.
- [ ] The existing check re-run and retry tests pass unedited.

## Context

- interface: `internal/delivery/github.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/merge_state_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_retry_test.go`
- instruction: `internal/delivery/github_test.go`
- instruction: `internal/delivery/check_rerun_test.go`
- instruction: `internal/delivery/retry_test.go`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/cli/cli.go`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `.agents/skills/roundfix/references/deliver.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPullRequestReadCarriesTheMergeState|TestCheckingWaitsWhileTheMergeStateIsBlocked|TestCheckingWaitsWhileTheMergeStateIsUnknown|TestABlockedMergeStateParksAtTheChecksTimeout|TestAnEmptyMergeStateKeepsTheListedChecksDecision|TestAPolicyRefusalIsTypedFromGhStderr|TestAPolicyRefusedMergeReturnsToChecking|TestASecondPolicyRefusalAtTheSameHeadParks|TestADeliveryErrorRetryAtAnUnchangedArchivedHeadResumesChecking)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPullRequestReadCarriesTheMergeState TestCheckingWaitsWhileTheMergeStateIsBlocked TestCheckingWaitsWhileTheMergeStateIsUnknown TestABlockedMergeStateParksAtTheChecksTimeout TestAnEmptyMergeStateKeepsTheListedChecksDecision TestAPolicyRefusalIsTypedFromGhStderr TestAPolicyRefusedMergeReturnsToChecking TestASecondPolicyRefusalAtTheSameHeadParks TestADeliveryErrorRetryAtAnUnchangedArchivedHeadResumesChecking; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the nine tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDeliverRetryPrintsARefusalWithoutUsage|TestDeliverRetryKeepsUsageForAnArgumentError)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverRetryPrintsARefusalWithoutUsage TestDeliverRetryKeepsUsageForAnArgumentError; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestAFailedCheckOutsideTheChangeIsRerunOnce|TestCheckRerunRestartsTheTimeout|TestRetryAfterArchiveReentersCheckingWithAPullRequest|TestRetriedArchivedItemPublishesAndMergesOnce|TestRetryRefusesAnArchivedItemWhoseHeadMoved|TestPullRequestReadCarriesTheMergeState)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAFailedCheckOutsideTheChangeIsRerunOnce TestCheckRerunRestartsTheTimeout TestRetryAfterArchiveReentersCheckingWithAPullRequest TestRetriedArchivedItemPublishesAndMergesOnce TestRetryRefusesAnArchivedItemWhoseHeadMoved TestPullRequestReadCarriesTheMergeState; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing check re-run and retry tests run unedited beside the merge-state read test, which does not exist before this Task.

## References

- `_prd.md` → User Stories 2-4; Core Features 2-5; Success Metrics 2-4; Acceptance evidence
- `_techspec.md` → The merge state; The merge refusal; The retry of a `delivery-error` park; The refusal output; API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Testing Approach 2; Testing Approach 3; Build Order 3
- ADR-0161; ADR-0192; ADR-0199

## Result

Implemented the Task 03 slice for Daemon Verification. The daemon-owned
`status: in_progress` is preserved; no Task Graph, other Task, commit, push or
Pull Request was changed or created.

- GitHub Pull Request reads request `mergeStateStatus`, normalize it to an
  upper-case trimmed `PullRequest.MergeState`, and carry the after-read value
  in `CheckReport.MergeState`. The conflict short path carries its observed
  value too.
- Checking keeps successful listed checks pending for `BLOCKED` and `UNKNOWN`,
  under the existing checks deadline. Other states, including an omitted
  state, preserve the existing decision. State changes print the specified
  checks log line.
- A non-zero merge whose stderr contains the branch-policy refusal,
  case-insensitively, returns `MergePolicyRefusalError`. The Engine remembers
  `<slug>@<head>` refusals in memory: the first logs the specified recovery
  line and returns to checking with the merge intent unmatched; the second
  takes the existing `delivery-error` park path.
- The archived retry path needed no production change. The new regression
  parks through two policy-refused merges, retries the same archived candidate
  with its recorded Pull Request, observes `checking`, and merges once. The
  successful merge also closes the original unmatched action intent.
- `Engine.Retry` errors print `Retry refused`, the refusal reason, the recorded
  stage/blocker when readable, and the specified no-side-effects text, with
  exit `2` and no usage block. Parse, configuration and store-opening errors
  still use the preflight path. Existing CLI refusal expectations were updated
  in `deliver_retry_test.go` and `deliver_token_ceiling_test.go` because this
  Task intentionally changes that public output; their no-owner-start and
  unchanged-queue assertions remain. The dependency Task's command and Skill
  references already describe this output and needed no edit here.

### Acceptance evidence

| Acceptance criterion | Implementation and focused-check evidence |
| --- | --- |
| Passing checks with `BLOCKED` wait until `CLEAN`, then merge exactly once | `TestCheckingWaitsWhileTheMergeStateIsBlocked` observes three checking polls, two fake-clock sleeps, one merge, one log per observed state, and no unmatched merge intent. |
| `UNKNOWN` waits; persistent `BLOCKED` times out; empty state keeps the prior decision | `TestCheckingWaitsWhileTheMergeStateIsUnknown`, `TestABlockedMergeStateParksAtTheChecksTimeout`, and `TestAnEmptyMergeStateKeepsTheListedChecksDecision` exercise these outcomes through `Engine.Run`; the timeout test observes exactly three seconds of fake time and no merge attempt. |
| First policy refusal returns to checking; second at the same head parks | `TestAPolicyRefusalIsTypedFromGhStderr` checks lower/upper-case policy stderr and an unrelated generic failure. `TestAPolicyRefusedMergeReturnsToChecking` observes another checks read, one successful merge and the recovery log. `TestASecondPolicyRefusalAtTheSameHeadParks` observes two refused attempts, no successful merge and `delivery-error`. |
| An unchanged archived `delivery-error` item with a Pull Request retries into checking and merges once | `TestADeliveryErrorRetryAtAnUnchangedArchivedHeadResumesChecking` reproduces the park through refused merges, invokes `Engine.Retry`, observes the recorded checking stage and incremented retry count, then runs the Engine and observes exactly one successful merge. No retry refusal was exposed on this path. |
| Retry refusal has the transcript shape without usage; an argument error retains preflight usage | `TestDeliverRetryPrintsARefusalWithoutUsage` compares stderr exactly, checks empty stdout and exit `2`, verifies no owner starts and confirms the item remains unchanged. `TestDeliverRetryKeepsUsageForAnArgumentError` asserts `Preflight failed`, `Usage:` and the unexpected-argument reason. `TestDeliverRetryRefusalOmitsAnUnreadableItem` covers omission of the Item block when its read fails. |
| Existing delivery check re-run and retry tests pass unedited | The complete `internal/delivery` package test run exited `0`; no existing file under `internal/delivery/*_test.go` was edited. `TestPullRequestReadCarriesTheMergeState` also exercises the scripted GitHub CLI read, normalization and the after-read state. |

### Focused checks

- Red starting evidence:
  `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/delivery -run TestCheckingWaitsWhileTheMergeStateIsBlocked -count=1`
  exited `1` because the merge-state fields and policy-refusal type did not
  exist.
- Red refusal reproduction:
  `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run TestDeliverRetryPrintsARefusalWithoutUsage -count=1`
  exited `1`, showing the old `Preflight failed` heading and `Usage:` block.
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/delivery -count=1`
  exited `0` (`ok roundfix/internal/delivery`, 26.626s).
- The initial broader CLI retry selection exposed two assertions for the old
  preflight heading. After updating those assertions to the authored output
  contract, `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestDeliver' -count=1`
  exited `0` (`ok roundfix/internal/cli`, 4.184s).
- `rtk proxy git -c core.fsmonitor=false diff --check` exited `0`.

The declared `## Verification` commands and repository-wide Verification were
not run in this child turn; the daemon owns those checks and Task settlement.
All GitHub behavior was checked through the specified scripted runner and fake
boundary; no live GitHub operation was performed. No follow-up work was found
within this slice.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/deliver_token_ceiling_test.go`

## Carry-forward provenance

- Source Run: `run_20261002T081705Z_a516b29ed11eed22`
- Source commit: `88c0100660eb1ea755f091fcce4a4740fde5cfa4`
