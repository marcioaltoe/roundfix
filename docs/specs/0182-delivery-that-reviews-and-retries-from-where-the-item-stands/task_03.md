---
task: task_03
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: completed
type: backend
complexity: high
---

# Task 03: A Delivery Retry carries forward from every Run of its item, newest first

## Overview

A Delivery Retry carries forward only from the Run recorded on the item. That record goes stale whenever a later Implement Run fails in a way the delivery workflow does not record, so on 2026-09-28 a retry of Spec 0175 re-carried the older Run and refused, while the newer Run's proved Tasks stayed behind. This Task makes the retry walk every terminal Implement Run of the item's Spec on the item branch, together with the recorded Run, newest first. It carries each Run's remaining proved Tasks through the Task Carry-Forward inspection that task_02 made skip Tasks already completed, rereads the item's Task Graph between Runs, records the newest Run on the item and prints one line per Run that moved Tasks. The Runs come from the Run Database, the Tasks and statuses from each Run Branch and the item worktree, and the carried commits land on the item branch that the next Implement Run and the pre-PR review read.

## Requirements

1. MUST replace `CarryForwardResult{RunID, Carried}` in `internal/delivery/engine.go` with `CarryForwardResult{RunID string; Runs []CarriedRun}` and `CarriedRun{RunID string; Carried []string}`, keep `Engine.Retry` recording `RunID` on the item and returning the result as `CarriedFrom`, and make `printDeliverRetryResult` in `internal/cli/deliver.go` print one `Carried forward from Run <run-id>: <task>, <task>` line per entry of `Runs`, in order, before `Retried <slug>: <blocker> -> <stage>`.
2. MUST replace `deliveryCarryForwardRun` in `internal/cli/deliver_workflow.go` with a selection that keeps the refusal for a recorded Run ID that is not an Implement Run of the Spec, and returns that Run plus every Implement Run of the Spec whose recorded local branch is the item branch, deduplicated and ordered newest first; `RunID` MUST name the first of them, or be empty when there is none.
3. MUST make `CarryForward` walk that list and, for each Run: skip it when its outcome is not accepted or it has no settled-completed Task; reread the item's Task Graph and skip it when every settled-completed Task is `completed` there, without requiring its Run Worktree; otherwise refuse when its Run Worktree is gone or the Specs Root is external; otherwise run `inspectCarryForwards` in the item worktree, refuse on a refusal reason, apply the ready candidates and append a `CarriedRun` when any Task moved.
4. MUST make a refusal name the Runs already carried in the same retry before its reason, keep its `NextAction` naming `roundfix reconcile <refusing-run-id> --carry-forward` in the item worktree and `roundfix deliver retry <slug>`, and leave the item branch advanced by exactly the Runs carried before it.
5. MUST update the existing tests in `internal/cli/deliver_recovery_test.go` and `internal/delivery/retry_test.go` to the new result shape without renaming any, and put the new workflow and output tests in `internal/cli/deliver_retry_runs_test.go` against real Git repositories and the Run Database fixtures of package `cli`, and the new engine test in `internal/delivery/retry_test.go`.
6. MUST update the retry paragraph of the Delivery queue section of `.agents/skills/roundfix/SKILL.md` and the `deliver retry` part of `docs/user-guide/commands.md` to contain `every terminal Implement Run of the item's Spec on the item branch, newest first` and to say one line is printed per Run carried from, run `make skills-sync`, and make the **Delivery Retry** entry of `CONTEXT.md` say it performs Task Carry-Forward `from every Run of the item's Spec, newest first`.

## Subtasks

- [ ] Change the recovery result and print one line per carried Run.
- [ ] Select every Run of the item's Spec on the item branch, newest first.
- [ ] Carry each Run's remaining proved Tasks, rereading the item graph between Runs.
- [ ] Name the carried Runs in a refusal.
- [ ] Update the guide, the skill and its mirror, and the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] With two BudgetExceeded Runs where the item records the older one and its Tasks are already on the item branch, the retry carries the newer Run's Tasks and `RunID` names the newer Run.
- [ ] Two Runs holding disjoint proved Tasks both carry, newest first, and the result lists both.
- [ ] A Run whose Tasks are all completed on the item is skipped even when its Run Worktree is gone.
- [ ] An older Run refusing after a newer Run carried names both Runs, and the item head holds only the newer Run's Tasks.
- [ ] `Engine.Retry` records the newest Run on the item, and the retry output prints one `Carried forward from Run` line per carried Run.
- [ ] The guide, the skill and its mirror, and the glossary state the rule, and `make skills-sync-check` passes.

## Context

- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_recovery_test.go`
- interface: `internal/delivery/retry_test.go`
- creates: `internal/cli/deliver_retry_runs_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0170-a-task-already-completed-on-its-target-is-nothing-to-carry.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne|TestItemRecoveryCarriesEveryRunNewestFirst|TestItemRecoverySkipsACompletedRunWithAGoneWorktree|TestItemRecoveryRefusalNamesTheRunsAlreadyCarried|TestDeliverRetryPrintsOneLinePerCarriedRun|TestItemRecoveryCarriesNothingForARunAlreadyCarried|TestRetriedItemRunsWithItsCompletedTasksCarried)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne TestItemRecoveryCarriesEveryRunNewestFirst TestItemRecoverySkipsACompletedRunWithAGoneWorktree TestItemRecoveryRefusalNamesTheRunsAlreadyCarried TestDeliverRetryPrintsOneLinePerCarriedRun TestItemRecoveryCarriesNothingForARunAlreadyCarried TestRetriedItemRunsWithItsCompletedTasksCarried; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the five new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestRetryRecordsTheNewestRunItCarriedFrom|TestRetryCarriesForwardBeforeReenteringTheRun|TestRetryRefusesWhenCarryForwardRefuses)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryRecordsTheNewestRunItCarriedFrom TestRetryCarriesForwardBeforeReenteringTheRun TestRetryRefusesWhenCarryForwardRefuses; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task `TestRetryRecordsTheNewestRunItCarriedFrom` does not exist, so the command fails.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "every terminal Implement Run of the item's Spec on the item branch, newest first" || { printf 'missing phrase in %s: %s\n' "$file" "every terminal Implement Run of the item's Spec on the item branch, newest first" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "from every Run of the item's Spec, newest first" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "from every Run of the item's Spec, newest first" >&2; exit 1; }; make skills-sync-check` — expected: exit 0; before this Task neither phrase exists.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — A Delivery Retry carries from every Run of its item; API Contract 4; Testing Approach 3; Build Order 3
- ADR-0170; ADR-0090; ADR-0158; ADR-0052; ADR-0044

## Result

Implemented a newest-first, deduplicated Delivery Retry carry-forward across the recorded Run and every terminal Implement Run of the Spec on the item branch. The retry rereads the item Task Graph per eligible Run, skips settled Tasks already completed on the item before requiring the Run Worktree, applies each accepted Run independently, preserves earlier carried commits when an older Run refuses, records the newest selected Run, and reports carried Tasks grouped by Run.

Acceptance evidence:

- Stale recorded Run: `TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne` uses two real BudgetExceeded Runs and asserts the selector returns exactly the newer and older Runs once each, the newer Run becomes `RunID`, and only its remaining Task moves.
- Multiple Runs: `TestItemRecoveryCarriesEveryRunNewestFirst` proves two disjoint Run task sets both move and appear newest first in `Runs`.
- Gone Worktree: `TestItemRecoverySkipsACompletedRunWithAGoneWorktree` removes the real Run Worktree after completing its Task on the item and proves the retry skips it without refusal.
- Partial progress and refusal: `TestItemRecoveryRefusalNamesTheRunsAlreadyCarried` proves the newer Run remains on the item, the refusing older Run does not move, the diagnostic names both Runs before its reason, and the next action still names reconcile in the item worktree followed by deliver retry.
- Engine and output: `TestRetryRecordsTheNewestRunItCarriedFrom` proves `Engine.Retry` persists the newest Run ID and returns the grouped result; `TestDeliverRetryPrintsOneLinePerCarriedRun` proves ordered per-Run lines precede the retried line.
- Documentation and mirror: the guide and canonical skill contain `every terminal Implement Run of the item's Spec on the item branch, newest first` and per-Run output wording; `CONTEXT.md` contains `from every Run of the item's Spec, newest first`; `rtk make skills-sync` exited 0 and `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited 0.

Focused checks:

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne$' ./internal/cli` failed to compile before implementation because `delivery.CarriedRun` and `CarryForwardResult.Runs` did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^(TestItemRecoveryCarriesTheNewestRunWhenTheItemRecordsAnOlderOne|TestItemRecoveryCarriesEveryRunNewestFirst|TestItemRecoverySkipsACompletedRunWithAGoneWorktree|TestItemRecoveryRefusalNamesTheRunsAlreadyCarried|TestDeliverRetryPrintsOneLinePerCarriedRun)$' ./internal/cli` — passed after the final test edit.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -run '^TestRetryRecordsTheNewestRunItCarriedFrom$' ./internal/delivery` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/delivery` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 ./internal/cli` — passed with process-table permission; the sandboxed attempt reached two unrelated Force Stop integration tests and failed only because macOS process-table access returned `operation not permitted`.
- `rtk git diff --check` — passed.

Daemon Verification was not run and Task status was not edited, per the Daemon settlement contract.

## Carry-forward provenance

- Source Run: `run_20260929T174719Z_24238e6aebd23928`
- Source commit: `a5da146dc709d8cc75149f21c8bc6a5d08f52617`
