---
task: task_01
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
status: pending
type: backend
complexity: high
---

# Task 01: The checking stage makes a Base Update when GitHub reports the Pull Request BEHIND

## Overview

The Delivery Engine's checking stage learns the merge state `BEHIND`. Read at a poll, it stops the Check Wait and the Queue Owner merges the refreshed default branch into the item branch through the queue's existing merge of the default branch, records the merged head as the newest candidate and sends the item back to `gating`. Three Base Updates per item and Engine bound it; a fourth `BEHIND` parks `pull-request-behind`, of Park Class `environment`. A command-engine test with a strict local forge proves the item that parks `delivery-error` today now merges. This answers the Backlog Entry of 2026-11-16 reported by quaestor, [`references/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md`](references/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md).

## Requirements

1. MUST add to `internal/delivery/engine.go` the exported constants `BlockerPullRequestBehind = "pull-request-behind"` and `BaseUpdateLimit = 3`, and the Engine field `baseUpdates map[string]int` keyed by Spec slug, as `_techspec.md` → Interfaces gives them.
2. MUST change `checkCandidate` as `_techspec.md` → Checking stage and invariants 1 to 4 say, right after the `CONFLICTING` branch and before the checks are classified, for a report whose `MergeState` is `BEHIND`:
   - with no `ConflictResolver`, park `pull-request-behind: pull request #<n> is behind its base branch`;
   - when the item already has `BaseUpdateLimit` Base Updates from this Engine, park `pull-request-behind: pull request #<n> is still behind its base branch after 3 Base Updates` without calling the resolver;
   - otherwise call `ResolveConflict(ctx, workDir, item.SpecSlug, head)` with the same arguments the conflict path uses; return its error wrapped as `bring pull request up to date: %w`; park `pull-request-conflict: <path>, …` for source paths, joined by `, ` as the conflict path joins them;
   - for a returned head that is not empty and differs from the candidate: add one to the item's count, write the owner log line of `_techspec.md` → API Contract 1 exactly, `roundfix: base update: Delivery Queue item <slug>: pull request #<n> was behind its base branch; merged the default branch at <head>`, append the head to `CandidateCommits` and set the stage `gating`;
   - for an empty head or the candidate itself, count nothing and treat the poll as pending, so the Check Wait goes on.
   `<n>` is the item's recorded `PullRequestNumber`.
3. MUST treat `BEHIND` as an awaited merge state everywhere `BLOCKED` and `UNKNOWN` are awaited: the end-of-poll pending test, and `checksTimeoutBlocker`, which then ends `; merge state BEHIND` (invariant 6).
4. MUST make `Engine.park` drop the item's count before it writes the park, so every park, including the step-error park `delivery-error`, starts the count over (invariant 5).
5. MUST NOT change the `CONFLICTING` path, the policy refusal handling of `mergeCandidate` (invariant 7), the Head Check Set evaluation, the stale and outside-change re-run rules, or the behavior for any merge state other than `BEHIND`.
6. MUST give `pull-request-behind` the Park Class `environment` in `ClassifyPark` (`internal/delivery/park_class.go`), with the next action of `_techspec.md` → API Contract 2 exactly: `run roundfix deliver retry <slug> once the base branch stops moving; the retry resumes at checking, which merges the default branch into the item branch again`, where `<slug>` is the item's Spec slug.
7. MUST add `internal/delivery/base_update_test.go`. Its tests call `t.Parallel()` first, use the existing fakes (`fakePullRequestBoundary`, `fakeDeliveryClock`, `newFakeDeliveryWorkflow`, `openDeliveryEngineStore`, `recordDeliveryItemWorkspace`) and a resolver fake of their own that returns a sequence of heads or source paths and counts its calls; they MUST NOT reuse the name `fakeConflictResolver`. Each test asserts the stored item, not only the returned error:
   - `TestABehindPullRequestReturnsToTheGateWithTheMergedHead`: one `BEHIND` poll; the item is at `gating` with `CandidateCommits` ending in the merged head, the resolver ran once with the item worktree and the candidate, and the log holds the exact line of requirement 2;
   - `TestABehindPullRequestIsMergedOnlyAtItsUpdatedHead`: through `Engine.Run` from `checking`, the forge fake reports `BEHIND` for the first head and `CLEAN` for the merged head, and returns a policy refusal for any merge at the first head; the item ends `merged`, exactly one merge was attempted and it named the merged head;
   - `TestABaseUpdateWithASourceConflictParksWithItsPaths`: the resolver returns `source.go`; the item parks `pull-request-conflict: source.go` of Park Class `conflict`;
   - `TestABaseUpdateWithNothingToMergeKeepsWaiting` with the subtests `merges_once_clean`, where `BEHIND` then `CLEAN` with an empty resolver head ends at `merging` after two polls, and `times_out_with_the_merge_state`, where `BEHIND` on every poll parks `checks-timeout: current-head checks not settled within delivery.check_timeout 3s; merge state BEHIND`;
   - `TestAFourthBehindPullRequestParksPullRequestBehind`: one Engine, four checking turns that each read `BEHIND`, a resolver that returns a new head each call; the resolver ran three times and the item parks with the exact bound blocker of requirement 2, Park Class `environment` and the exact next action of requirement 6;
   - `TestTheBaseUpdateCountStartsOverAfterAPark`: after that park, the same Engine returns the item to `checking` as a Delivery Retry does, and the next `BEHIND` calls the resolver a fourth time and moves the item to `gating`;
   - `TestAnOwnerWithoutABaseUpdateParksPullRequestBehind`: an Engine without a `ConflictResolver` parks `pull-request-behind: pull request #1 is behind its base branch` at the first `BEHIND`;
   - `TestAPolicyRefusalThatComesBackBehindGetsABaseUpdate`: the first poll is `CLEAN`, the merge is refused by the branch policy, the next poll is `BEHIND`; the resolver runs and the item reaches `gating` without a second merge attempt at that head;
   - `TestAConflictingAndBehindPullRequestTakesTheConflictPath`: a report with `Mergeable` `CONFLICTING` and `MergeState` `BEHIND` takes the conflict path unchanged and counts no Base Update.
8. MUST add `{BlockerPullRequestBehind, ParkClassEnvironment}` to the table of `TestParkClassesAndNextActions` in `internal/delivery/park_class_test.go`; `TestEveryBlockerHasAParkClass` then holds the new constant.
9. MUST add `internal/cli/deliver_base_update_test.go` with `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate`, built on the child-process pattern of `TestCommandDeliveryEngineRecordsThePullRequest` (`internal/cli/deliver_record_pull_request_test.go`): a temporary bare origin, a local `gh` script on the child's `PATH` through an environment variable of its own, `setCommandEnvForTest(t, cliTestHelperEnv, "1")` before the engine starts, `Defaults.Verification` set to `true` and `Defaults.ArtifactDir` set to a temporary directory. The item starts at `publishing` with its archive commit, as in that test. After the item branch is cut, the test pushes a commit adding `moved.txt` to `main` of the origin from a temporary clone. The `gh` script is strict:
   - `gh pr view` prints `mergeStateStatus` `BEHIND` while `main` of the origin is not an ancestor of the Pull Request head, `CLEAN` otherwise;
   - `gh pr merge` while behind writes a refusal marker into its state directory, prints `X Pull request #7 is not mergeable: the base branch policy prohibits the merge.` to standard error and exits 1; otherwise it builds the squash commit as the record test's script does.
   The test asserts that the item ends `merged`; that its newest candidate is a commit with two parents whose second parent is the moved `main` and whose message carries `Roundfix-Delivery: derived-merge`; that the merge commit on the origin holds `moved.txt`; that the child's output holds `roundfix: base update: Delivery Queue item `; and that no refusal marker exists.
10. MUST NOT change any file under `internal/cli` other than the new test, any guide, skill, glossary or ADR, or any Governed Path.

## Subtasks

- [ ] Add the `BEHIND` path, the count and its reset, and the awaited merge state to the checking stage.
- [ ] Give `pull-request-behind` its Park Class and next action.
- [ ] Add the Delivery Engine tests and the class-table case.
- [ ] Add the command-engine test with the strict local forge.

## Acceptance Criteria

- [ ] A Pull Request reported `BEHIND` never reaches a merge at that head; the merged head goes back through `gating`, and a run merges at it.
- [ ] A source conflict parks `pull-request-conflict: <path>`; nothing to merge keeps polling and times out with `; merge state BEHIND`.
- [ ] A fourth `BEHIND` parks `pull-request-behind` with its exact blocker, class and next action, and the count starts over after any park.
- [ ] The strict local forge, which parks the item `delivery-error` on the tree before this Task, ends with the item merged and the base branch's file in the merge.
- [ ] Every existing test of `internal/delivery` passes unchanged except the class-table addition.
- [ ] Removing the `BEHIND` path, the count, its reset in `park` or `BEHIND` from the awaited states makes a named test fail.

## Context

- instruction: `docs/adr/0435-a-delivery-pull-request-behind-its-base-branch-is-brought-up-to-date-by-a-merge-of-the-default-branch.md`
- instruction: `docs/adr/0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md`
- instruction: `internal/delivery/merge_state_test.go`
- instruction: `internal/delivery/conflict_test.go`
- instruction: `internal/cli/deliver_record_pull_request_test.go`
- instruction: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/base_update_test.go`
- creates: `internal/cli/deliver_base_update_test.go`

## Consumers

- The blocker `pull-request-behind`, its Park Class and next action from `delivery.ClassifyPark`, and the checking stage's transitions — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestABehindPullRequestReturnsToTheGateWithTheMergedHead|TestABehindPullRequestIsMergedOnlyAtItsUpdatedHead|TestABaseUpdateWithASourceConflictParksWithItsPaths|TestABaseUpdateWithNothingToMergeKeepsWaiting|TestAFourthBehindPullRequestParksPullRequestBehind|TestTheBaseUpdateCountStartsOverAfterAPark|TestAnOwnerWithoutABaseUpdateParksPullRequestBehind|TestAPolicyRefusalThatComesBackBehindGetsABaseUpdate|TestAConflictingAndBehindPullRequestTakesTheConflictPath|TestParkClassesAndNextActions|TestEveryBlockerHasAParkClass)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestABehindPullRequestReturnsToTheGateWithTheMergedHead TestABehindPullRequestIsMergedOnlyAtItsUpdatedHead TestABaseUpdateWithASourceConflictParksWithItsPaths TestABaseUpdateWithNothingToMergeKeepsWaiting/merges_once_clean TestABaseUpdateWithNothingToMergeKeepsWaiting/times_out_with_the_merge_state TestAFourthBehindPullRequestParksPullRequestBehind TestTheBaseUpdateCountStartsOverAfterAPark TestAnOwnerWithoutABaseUpdateParksPullRequestBehind TestAPolicyRefusalThatComesBackBehindGetsABaseUpdate TestAConflictingAndBehindPullRequestTakesTheConflictPath TestParkClassesAndNextActions/pull-request-behind TestEveryBlockerHasAParkClass; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/delivery` — expected: exit 0. Before this Task the new tests and the class-table case do not exist, so their pass lines are missing; afterwards each passes and the whole Delivery Engine package passes.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestCommandDeliveryEngineBringsABehindPullRequestUpToDate$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestCommandDeliveryEngineBringsABehindPullRequestUpToDate ' || { echo 'missing pass: TestCommandDeliveryEngineBringsABehindPullRequestUpToDate' >&2; exit 1; }; printf '%s\n' "$out" | grep -q 'roundfix: base update: Delivery Queue item ' || { echo 'missing Base Update log line' >&2; exit 1; }` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards the strict local forge merges the item after a Base Update.

## References

- `_prd.md` → Goals (first, second, fourth, fifth and sixth bullets); Core Features 1 and 2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Decisions
- `_techspec.md` → Checking stage; The count; Park Class; What the command engine already provides; Interfaces and invariants 1 to 7; API Contract 1; API Contract 2; Testing Approach (Task 1); Build Order 1
- [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md)
- ADR-0435, ADR-0192, ADR-0369, ADR-0311
