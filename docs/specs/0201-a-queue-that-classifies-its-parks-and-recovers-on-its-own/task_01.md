---
task: task_01
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: pending
type: backend
complexity: high
---

# Task 01: Every park names its class and next command, and a check that failed outside the change is re-run once

## Overview

A parked Delivery Queue item carries only a blocker string, and the Pending Question answers two blockers specifically and the rest generically. This Task adds the Park Class table, prints one `Park:` line per parked item in `deliver status`, and answers the Pending Question from the same table. It also lets the owner re-run once a failed GitHub Actions check whose failing Go packages all lie outside the item's change, and park a second failure as `flaky-check`.

## Requirements

1. MUST add `ClassifyPark(queue, item)` and the `ParkClass` values `dependency`, `conflict`, `environment`, `flaky-check`, `finding`, `budget`, `review`, `authorization` and `unclassified` in a new `internal/delivery/park_class.go`, mapping each blocker as the TechSpec's Park Classes table states. A blocker matches exactly or by its `<blocker>:` prefix. Blockers later Tasks add are classified by those Tasks.
2. MUST make `PendingQuestionFor` answer with `ClassifyPark(...).Next`, and keep every answer that exists today byte-identical, including the three cases of `TestPendingQuestionAnswersEachBlockerClass`.
3. MUST make `roundfix deliver status` print `Park: <slug> <class>: <next>` for each parked item in queue order, after the Warning lines and before the Limits line, and leave the output of a queue without a parked item byte-identical. The test that pins a parked item's full status output, `TestDeliverStatusPrintsTheItemWorktree`, changes on purpose to expect the `Park:` line.
4. MUST add a test that parses the Go files of `internal/delivery` and fails when an exported `Blocker…` constant maps to `unclassified`, so a blocker added without a class fails.
5. MUST add the optional `CheckRecovery` dependency and `CheckFailure` type, and implement them on `GitHubCLI` in a new `internal/delivery/check_rerun.go`, as the TechSpec's Failed check re-run section states: the run ID from the check `link`, `gh run view <id> --json attempt`, `gh run view <id> --log-failed`, Go package attribution from `FAIL\t<import path>\t<n>s` lines through the worktree's `go.mod` module line, and the item's changed paths from `git diff --name-only` against the merge base with the refreshed remote default branch. A `[build failed]` or `[setup failed]` line, or a log without a Go package, makes the failure unattributable.
6. MUST make `checkCandidate` re-run a failed check once (`gh run rerun <id> --failed`) only when `OutsideChange` holds and `Attempt` is `1`, restart the check timeout once, log the re-run, and keep polling. A re-run check that later passes MUST join `flaky-check: <check> passed on re-run` to the item's Warning; a second failure MUST park `flaky-check: <package>, …` (class `flaky-check`). Every other failure MUST park `checks-failed` as today, and a nil `CheckRecovery` MUST keep today's behavior.
7. MUST wire `delivery.NewGitHubCLI(loaded.GitRoot)` as `Checks` in `newCommandDeliveryEngine`.
8. MUST describe the Park Classes, the `Park:` line and the one re-run in `docs/user-guide/commands/deliver.md`, the delivery command file Spec 0194 creates, with the phrases `Park:` and `flaky-check`.
9. MUST NOT reach GitHub or run `gh` in any test, change a command's flags or help text, or rename or remove an existing top-level test.

## Subtasks

- [ ] Add the Park Class table and route the Pending Question through it.
- [ ] Print the `Park:` lines in `deliver status`.
- [ ] Attribute a failed check to Go packages and re-run it once when it failed outside the change.
- [ ] Describe both behaviors in the delivery command guide.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] Every exported blocker constant maps to a class other than `unclassified`, and a new unclassified constant fails the sweep test.
- [ ] `deliver status` with two parked items prints two `Park:` lines in queue order, reproducing the shape of Surface Transcript 1; a queue without a parked item prints today's output.
- [ ] A failed check whose log names only packages the item did not change is re-run exactly once; a pass then records the warning, and a second failure parks `flaky-check`.
- [ ] A failure in a changed package, a build failure, a log without a Go package, or a run already past its first attempt parks `checks-failed` with no re-run.

## Context

- interface: `internal/delivery/question.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/github.go`
- creates: `internal/delivery/park_class.go`
- creates: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/check_rerun.go`
- creates: `internal/delivery/check_rerun_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_test.go`
- creates: `internal/cli/deliver_park_status_test.go`
- creates: `docs/user-guide/commands/deliver.md`
- instruction: `internal/delivery/question_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEveryBlockerHasAParkClass|TestAFailedCheckOutsideTheChangeIsRerunOnce|TestARerunCheckThatPassesRecordsAWarning|TestARerunCheckThatFailsAgainParksAsFlakyCheck|TestAFailedCheckInAChangedPackageParksWithoutARerun|TestAnUnattributableCheckFailureParksWithoutARerun|TestACheckPastItsFirstAttemptIsNotRerun|TestGitHubCLIAttributesAFailedCheckToGoPackages|TestPendingQuestionAnswersEachBlockerClass)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryBlockerHasAParkClass TestAFailedCheckOutsideTheChangeIsRerunOnce TestARerunCheckThatPassesRecordsAWarning TestARerunCheckThatFailsAgainParksAsFlakyCheck TestAFailedCheckInAChangedPackageParksWithoutARerun TestAnUnattributableCheckFailureParksWithoutARerun TestACheckPastItsFirstAttemptIsNotRerun TestGitHubCLIAttributesAFailedCheckToGoPackages TestPendingQuestionAnswersEachBlockerClass; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the new named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDeliverStatusPrintsAParkLinePerParkedItem|TestDeliverStatusPrintsTheItemWorktree|TestDeliverStatusPrintsNoWarningLineWithoutAWarning)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverStatusPrintsAParkLinePerParkedItem TestDeliverStatusPrintsTheItemWorktree TestDeliverStatusPrintsNoWarningLineWithoutAWarning; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `TestDeliverStatusPrintsAParkLinePerParkedItem` does not exist, so the command fails.
- `for pair in "docs/user-guide/commands/deliver.md|Park:" "docs/user-guide/commands/deliver.md|flaky-check"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task neither phrase is in the delivery command guide, so the command fails.

## References

- `_prd.md` → Goal 4; User Stories 5-6; Core Features 4-5; Success Metrics 4-5; Declared breaks
- `_techspec.md` → Park Classes; Failed check re-run; API Contract 1; Surface Transcript 1; Testing Approach 1; Build Order 1
- ADR-0153; ADR-0165
