---
task: task_04
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: completed
type: backend
complexity: medium
---

# Task 04: The Turn Record lists the checks an item waits for, and deliver status prints them

## Overview

Fluxus asked that `deliver status`, not only the console log, say what the
checking stage waits for (reported by fluxus on 2026-11-05,
[references](references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md)).
After this Task the Queue Owner records the pending members of the Head Check
Set on the Turn Record of the item it advances in `checking`, and `deliver
status` prints them on that item's `Turn:` line, as `_techspec.md` → Data
Models, API Contract 6 and Surface Transcript 1 state.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestDeliverStatusReproducesSurfaceTranscriptQueueTurns|TestDeliverStatusOrdersWarningsTurnsAndParks)$' ./internal/cli`
   and record in the Result that it passes.
2. MUST add `Checks []string` with the JSON key `checks`, omitted when
   empty, to `ItemTurn`, the per-slug pending list under `turnMu`,
   `Engine.checkTurn` and the checking form of `TurnLines` in
   `internal/delivery/turn.go`, and set `ItemTurn.Checks` for the advancing
   item in `recordTurnLocked`, exactly as API Contract 6 states; the line is
   exactly `Turn: <slug> advancing at checking; waiting for <n> check(s): <member>; <member>`.
3. MUST call `checkTurn` from `checkCandidate` in
   `internal/delivery/engine.go` each time it writes the waiting line, with
   the pending list, and with `nil` on every return; a nil recorder records
   nothing and adds no log line.
4. MUST add to `internal/delivery/turn_test.go`:
   - `TestTurnLinesNameTheChecksAnItemWaitsFor`: a queue with `0149-example`
     at `checking` and `0150-other` at `queued`, and a `Turn` with
     `0149-example` `advancing` and `Checks` `[CI | Backend / Lint, Type
     Check, Test & Build, .github/workflows/ci-frontends.yml (no run yet)]`
     and `0150-other` `turn-wait`; `TurnLines` returns exactly
     `Turn: 0149-example advancing at checking; waiting for 2 check(s): CI | Backend / Lint, Type Check, Test & Build; .github/workflows/ci-frontends.yml (no run yet)`
     and `Turn: 0150-other waiting for the owner's turn (1 of 1 waiting)`.
     Subtest `no checks`: without `Checks` the first line is
     `Turn: 0149-example advancing at checking`.
   - `TestTheOwnerRecordsTheChecksAnItemWaitsFor`: an engine over
     `newCheckWaitFixture` with a recorder that keeps every `Turn` under a
     mutex, `advancingTurn` called for the item, and `checkCandidate` over
     two polls of `recoveryReport("pending")` and then
     `recoveryReport("pass")`; some recorded `Turn` lists the item
     `advancing` with `Checks` `[CI / verify]`, the last recorded `Turn` has
     no `Checks` for it, and the item reaches `merging`.
5. MUST add to `internal/cli/deliver_turn_test.go`
   `TestDeliverStatusPrintsTheChecksAnItemWaitsFor`, reproducing Surface
   Transcript 1 of `_techspec.md`: a queue of `0149-example` at `checking`
   with worktree `/worktrees/0149-example`, the owner PID `os.Getpid()`
   with identity `recorded-owner-identity` proven by
   `fakeDeliveryOwnerProcesses`, and the Turn Record the transcript
   describes written at `turnRecordPath`; it asserts every stdout line, the
   empty stderr and exit `0`.
6. MUST keep every other existing test of `./internal/delivery` and every
   `deliver` test of `./internal/cli` passing unchanged.
7. MUST run sabotage 7 of `_techspec.md` → Testing Approach on a temporary
   copy and record which named test fails.
8. MUST NOT change any other `Turn:` form, the item line of `deliver
   status`, any blocker or log text, `internal/delivery/github.go`,
   `internal/delivery/head_check_set.go`, the Run Database, any document,
   the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the status transcript tests.
- [ ] Record the pending members on the Turn Record.
- [ ] Print the checking `Turn:` form.
- [ ] Add the three tests and run the sabotage.

## Acceptance Criteria

- [ ] The Turn Record lists the pending members of the item the owner advances in `checking`, and drops them when the stage ends.
- [ ] `deliver status` prints Surface Transcript 1 exactly.

## Context

- interface: `internal/delivery/turn.go`
- interface: `internal/delivery/turn_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/check_wait_test.go`
- interface: `internal/cli/deliver_turn.go`
- interface: `internal/cli/deliver_turn_test.go`
- interface: `internal/cli/deliver.go`

## Consumers

- `ItemTurn`, `TurnLines` and the Turn Record file — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run 'Turn|CheckWait|Check' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTurnLinesNameTheChecksAnItemWaitsFor TestTheOwnerRecordsTheChecksAnItemWaitsFor; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 30m -v -run 'Deliver|Delivery|Queue|Turn' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestDeliverStatusPrintsTheChecksAnItemWaitsFor " || { echo "missing pass: TestDeliverStatusPrintsTheChecksAnItemWaitsFor" >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/delivery ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/delivery ./internal/cli` — expected: exit 0. Before this Task none of the three new tests exists, so the first pass line is missing; afterwards the Turn Record and the `Turn:` line behave as API Contract 6 and Surface Transcript 1 state, the engine and `deliver` tests pass, and both platforms build.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 4
- `_techspec.md` → Current behavior 5; Data Models; API Contract 6; Surface Transcript 1; Testing Approach; Build Order 4; ADR-0365; ADR-0369

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-concurrency, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, evidence-gate, tech-writer.

The advancing item's Turn Record now carries optional `checks`, with the
pending list copied per slug under `turnMu`. Each waiting log write calls
`checkTurn` with that list; a deferred call clears it on every return from
`checkCandidate`, including early returns. Without a Turn Recorder,
`checkTurn` returns before changing state, reading storage or writing a log.
The checking `Turn:` line names the pending members with the specified count
and separators. Existing Turn forms, status item lines, blockers and log
texts are unchanged.

Acceptance evidence:

- Pending members and stage exit: `TestTheOwnerRecordsTheChecksAnItemWaitsFor`
  passed with two pending polls followed by one passing poll. Its recorded
  turns include the advancing item with `[CI / verify]`; the last turn has
  no checks and its item and the stored item are at `merging`.
- Exact status transcript: `TestDeliverStatusPrintsTheChecksAnItemWaitsFor`
  passed through public `deliver status` dispatch with the live owner proof
  and a serialized Turn Record. It asserts all stdout lines, empty stderr
  and exit 0. `TestTurnLinesNameTheChecksAnItemWaitsFor` passed for both the
  two-member checking line and its `no checks` subtest, including the other
  item's unchanged owner-turn wait line.

Focused checks (Go commands invoked through `rtk proxy`):

- Before changing code,
  `go test -count=1 -timeout 20m -v -run '^(TestDeliverStatusReproducesSurfaceTranscriptQueueTurns|TestDeliverStatusOrdersWarningsTurnsAndParks)$' ./internal/cli`
  exited 0; both characterization tests passed.
- After adding the tests and before implementation,
  `go test -count=1 -timeout 2m -run '^TestTurnLinesNameTheChecksAnItemWaitsFor$' ./internal/delivery`
  exited 1 because `ItemTurn.Checks` did not exist.
- After implementation,
  `go test -count=1 -timeout 3m -v -run '^(TestTurnLinesNameTheChecksAnItemWaitsFor|TestTheOwnerRecordsTheChecksAnItemWaitsFor|TestTurnLinesNameEachItemsTurn|TestTheOwnerRecordsEachItemsTurn)$' ./internal/delivery`
  exited 0; all selected tests and subtests passed.
- `go test -count=1 -timeout 3m -v -run '^(TestDeliverStatusPrintsTheChecksAnItemWaitsFor|TestDeliverStatusReproducesSurfaceTranscriptQueueTurns|TestDeliverStatusOrdersWarningsTurnsAndParks|TestDeliverStatusIgnoresATurnRecordOfAnotherOwner|TestTheTurnRecorderWritesBesideTheConsoleLog)$' ./internal/cli`
  exited 0; all selected tests and subtests passed unchanged alongside the
  new transcript test.
- `go test -race -count=1 -timeout 3m -run '^(TestTurnLinesNameTheChecksAnItemWaitsFor|TestTheOwnerRecordsTheChecksAnItemWaitsFor|TestTheOwnerRecordsEachItemsTurn)$' ./internal/delivery`
  exited 0, with no race report.
- `go test -count=1 -timeout 3m -v -run '^(TestTheCheckWaitDefaultsToSixtyMinutes|TestTheCheckWaitUsesTheConfiguredDuration|TestEachCheckWaitStartsItsOwnDeadline|TestNilCheckRecoveryKeepsFailedCheckBehavior)$' ./internal/delivery`
  exited 0; the selected existing check-wait and failed-check tests passed.
- `git diff --check` exited 0.

Sabotage 7 ran in a temporary copy at
`/private/tmp/roundfix-task04-sabotage7-6o78h3j3`, with independent Git
metadata. Removing only the checking suffix from `TurnLines`, then running
`go test -count=1 -timeout 3m -v -run '^TestDeliverStatusPrintsTheChecksAnItemWaitsFor$' ./internal/cli`
exited 1: `TestDeliverStatusPrintsTheChecksAnItemWaitsFor` failed because
actual stdout contained `Turn: 0149-example advancing at checking`, while
the expected line named the two pending checks. Exit 0 and empty stderr
still matched, so the failure was the missing status text. The working
tree retained the suffix.

Only `internal/delivery/engine.go`, `internal/delivery/turn.go`,
`internal/delivery/turn_test.go`, `internal/cli/deliver_turn_test.go` and
this Result were edited. The initial dirty path was this Task's
Daemon-owned status change. Task status and the Task Graph were left alone;
the declared Verification, repository gates, full suites and cross-platform
builds were not run, for the Daemon to execute after handoff. No commit,
push or pull request was made.
