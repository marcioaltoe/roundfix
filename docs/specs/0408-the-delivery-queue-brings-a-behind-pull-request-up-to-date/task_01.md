---
task: task_01
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
status: completed
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


## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate, golang-structs-interfaces.

The checking stage now handles `BEHIND` after the unchanged `CONFLICTING` path. It uses the existing resolver, records only a different nonempty head, writes the exact Base Update log, and returns that head to `gating`. The Engine counts Base Updates by Spec slug, bounds them with `BaseUpdateLimit = 3`, and drops the count before every park. `BEHIND` remains pending when nothing merges and appears in its timeout blocker. `pull-request-behind` has the environment Park Class and the specified retry action. The policy-refusal path, Head Check Set evaluation and stale/outside-change recovery code are unchanged.

Scope: production changes are confined to `internal/delivery/engine.go` and `internal/delivery/park_class.go`. Added the two declared test files and the single class-table row. The only pre-existing changed path was this Task file, whose Daemon-owned `status: in_progress` is preserved. No other Task, manifest, guide, skill, ADR or Governed Path was changed.

### Focused checks

- Before changing production code, `rtk proxy go test ./internal/delivery -run '^TestABehindPullRequestReturnsToTheGateWithTheMergedHead$' -count=1` exited 1 at the new stored-item assertion: stage `merging`, candidates `[candidate]`, resolver calls `[]`, one poll. This reproduced the missing Base Update.
- Final positive check, exit 0 for both packages:

  ```sh
  rtk proxy go test ./internal/delivery ./internal/cli -run 'Test(ABehind|ABaseUpdate|AFourthBehind|TheBaseUpdate|AnOwnerWithoutABase|APolicyRefusalThatComesBack|AConflictingAndBehind|CheckingWaits|APolicyRefused|ASecondPolicy|ABlockedMerge|AnEmptyMerge|AConflictingPull|AnUnknownMergeable|ADerivedConflict|ASourceConflict|ALaggingPull|RetryResumesAResolved|ParkClasses|EveryBlocker|CommandDeliveryEngine(BringsABehindPullRequestUpToDate|RecordsThePullRequest))' -count=1 -v > /private/tmp/roundfix-0408-task01-focused.log
  ```

  The log includes every new test's pass line, both required no-op subtests, `TestParkClassesAndNextActions/pull-request-behind`, `TestEveryBlockerHasAParkClass`, the selected existing conflict/merge-state/policy regressions, both command-engine scenarios, and the Base Update owner line.
- `rtk proxy git diff --check` exited 0. `rtk proxy gofmt -l` over all five changed Go files exited 0 with no output.

### Evidence by acceptance criterion

| Criterion | Implementation and focused evidence |
| --- | --- |
| A behind head is never merged; the updated head goes through gating and merges | `TestABehindPullRequestReturnsToTheGateWithTheMergedHead` asserts persisted stage/candidate history, resolver arguments and exact owner line. `TestABehindPullRequestIsMergedOnlyAtItsUpdatedHead` reports readiness by head, finishes persisted `merged`, and asserts one merge at `updated`, one gate/publication, and no review/Run/archive replay. |
| Source conflicts park; no-op updates wait and time out with BEHIND | `TestABaseUpdateWithASourceConflictParksWithItsPaths` asserts persisted `pull-request-conflict: source.go` and conflict class. `TestABaseUpdateWithNothingToMergeKeepsWaiting/merges_once_clean` covers both empty and unchanged returned heads and two polls before `merging`; `times_out_with_the_merge_state` asserts the exact 3s timeout blocker and four fake-clock polls without a candidate append. |
| Fourth BEHIND parks with the exact blocker/class/action; any park resets the count | `TestAFourthBehindPullRequestParksPullRequestBehind` asserts three resolver calls, persisted park and exact action. `TestTheBaseUpdateCountStartsOverAfterAPark` resumes the same Engine at checking and observes the fourth resolver call and persisted `gating`. Additional `TestABaseUpdateErrorKeepsItsCauseAndAParkResetsTheBound` proves `%w` preserves the resolver error and that the Run's `delivery-error` park resets prior updates. The missing-resolver, policy-refusal recovery and CONFLICTING-precedence tests also pass. |
| Strict local forge merges with the moved base file | `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate` passes through the real constructor and child process, bare origin, separate clone moving main, gate command `true`, record publication and strict local gh. It asserts persisted `merged`, a two-parent newest candidate with moved main as second parent, the `Roundfix-Delivery: derived-merge` trailer, `moved.txt` in the final squash, the child's Base Update line and no refusal marker. A temporary overlay restoring the full pre-change BEHIND behavior made this test fail at persisted `parked` with `delivery-error: merge pull request: merge pull request: X Pull request #7 is not mergeable: the base branch policy prohibits the merge.` |
| Existing delivery tests remain unchanged except the class-table row | The focused conflict, lagging-head, conflict-retry, BLOCKED/UNKNOWN/empty-state, policy-refusal, class-table and blocker-sweep tests pass. All existing delivery test files remain byte-identical except the declared class-table addition. The whole-package run belongs to the Daemon's declared Verification and was not run here. |
| Removing each new guard makes a named test fail | Five temporary Go overlays independently removed the BEHIND path, counter increment, park reset, awaited-state pending guard and timeout suffix. Each exited 1 at its expected behavioral assertion as detailed below. |

### Negative probes

Each probe used `go test -overlay <temporary-overlay.json> <package> -run '^<named-test>$' -count=1` through a Python driver. The overlays replaced only the compiler's view of `engine.go`; repository source was not changed. The driver inspected both exit status and the named failing test. Temporary probe files were removed after execution.

| Removed behavior | Named test and observed failure |
| --- | --- |
| BEHIND resolution path | `TestABehindPullRequestReturnsToTheGateWithTheMergedHead`: timeout park with unchanged candidate and zero resolver calls instead of gating. |
| Counter increment | `TestAFourthBehindPullRequestParksPullRequestBehind`: fourth resolver call and fourth head recorded at gating instead of the bound park. |
| Count reset in park | `TestTheBaseUpdateCountStartsOverAfterAPark`: retry parks immediately without a fourth resolver call. |
| BEHIND in end-of-poll awaited states | `TestABaseUpdateWithNothingToMergeKeepsWaiting`: no-op cases reach merging after one poll; timeout case also reaches merging. |
| BEHIND in timeout-state formatting | `TestABaseUpdateWithNothingToMergeKeepsWaiting/times_out_with_the_merge_state`: stored timeout omits `; merge state BEHIND`. |
| Full pre-change BEHIND behavior, for command-engine reproduction | `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate`: two policy refusals and the exact original delivery-error park quoted above. |

An exploratory command-engine probe initially removed only the resolution path while retaining the new pending-state guard. It kept waiting and was explicitly stopped; its termination is not acceptance evidence. The corrected command probe removed the whole new BEHIND behavior and produced the behavioral failure recorded above. A final positive focused run followed the probes.

The Task's declared Verification commands, the complete delivery package suite and repository Verification were not run; the Daemon owns them and Task settlement. No terminal Task verdict is claimed. Documentation/glossary/skill changes remain assigned to task_02. No additional follow-up was identified.


### Verification feedback repair — attempt 1

Inspected the supplied Daemon artifact at `/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261119T035531Z_b63aa3cc448dbba2/verification/batch-001-attempt-1.log` and the Real-Time Budget detector. The supplied run passed `internal/delivery` but its overall `make verify-changed` failed because the new test file declared an inline short CheckTimeout without an existing file allowance. The fake clock prevented a wall-clock wait, but the detector also counts that constant field declaration.

Reworked only `newBaseUpdateFixture` in the declared `internal/delivery/base_update_test.go`: it now composes the existing `newMergeStateFixture` engine, stored item, fake workflow and fake clock, then supplies its own resolver, forge callbacks and log. The shared fixture already owns the same three-second Check Wait and one-second polling interval. A checked interface assertion obtains its workflow for the unchanged gate/publication assertions; its fixed recovery head is replaced by this fixture's direct checking setup. Removed the duplicate Engine/timing construction and unused time import. No detector, ceiling, configuration, production code or behavioral assertion changed in this repair.

Fresh focused evidence after the fixture change:

- `rtk proxy go test ./internal/testfixture -run '^Test(EveryRealTimeBudgetIsWithinItsFileCeiling|RealTimeBudgetRuleRejectsSeededViolations)$' -count=1` exited 0, covering the reported contract and its seeded negative cases.
- `rtk proxy go test ./internal/delivery -run 'Test(ABehind|ABaseUpdate|AFourthBehind|TheBaseUpdate|AnOwnerWithoutABase|APolicyRefusalThatComesBack|AConflictingAndBehind|CheckingWaits|APolicyRefused|ASecondPolicy|ABlockedMerge|AnEmptyMerge)' -count=1 -v` passed. The output at `/private/tmp/roundfix-0408-task01-repair-focused.log` contains all new Base Update tests, both required no-op subtests and the existing shared-fixture merge-state/policy tests, including the exact three-second BEHIND timeout assertion.
- Repeated the five temporary-overlay negative probes against the revised fixture. Removing the BEHIND path, counter increment, park reset, awaited-state pending guard or timeout suffix still failed its respective named behavioral assertion. The driver checked the nonzero status, named failing test and expected persisted-state diagnostic for each; no source mutation remained.
- `rtk proxy git diff --check` exited 0, and `rtk proxy gofmt -l internal/delivery/base_update_test.go` produced no output.

This repair changes only the new delivery test fixture and this Result. The Task remains `in_progress`. The Agent did not rerun `make verify-changed`, either declared Verification command or repository Verification; the Daemon owns the next complete Verification sequence and settlement.
