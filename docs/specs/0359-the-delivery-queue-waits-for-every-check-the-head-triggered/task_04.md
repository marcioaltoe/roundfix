---
task: task_04
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: pending
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
