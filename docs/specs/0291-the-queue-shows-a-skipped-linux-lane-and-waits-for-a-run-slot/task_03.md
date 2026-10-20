---
task: task_03
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: pending
type: backend
complexity: high
---

# Task 03: The engine waits for a Run slot up to runs.slot_wait

## Overview

The queue parked Spec 0289 as `delivery-error` only because other
repositories' Runs held every slot of the Active Implement Run ceiling, and it
waited for a manual retry (`_prd.md` → Sources 2). This Task adds
`runs.slot_wait` to User Config and the Run Slot Wait to the Delivery Engine:
on a ceiling-only refusal the running item records the wait, the engine polls
the Run slots, starts the Run when one frees and parks at the earlier deadline
exactly as today. It is verifiable on its own with the engine fakes, a fake
clock, a fake sleeper, a fake slot reader and a real Run Database.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md`](references/2026-10-20-the-queue-parks-an-item-at-the-run-ceiling.md)
   (`_prd.md` → Sources 2), by implementing API Contracts 5 and 6 of
   `_techspec.md` in `internal/config/config.go` and
   `internal/delivery/engine.go`, with `Runs.SlotWait`, `WaitRunSlot`,
   `RunSlotHolder`, `RunSlots`, `RunSlotRefusal`, `RunSlotReader`,
   `RunSlotWaitText` and the three `EngineDependencies` fields as
   `_techspec.md` → Interfaces sketches them.
2. MUST read `runs.slot_wait` as a Go duration with the existing duration
   parser, default `2h`, apply it only outside Project Config, remove a
   Project Config value with the existing ignored-setting warning, refuse a
   negative value with `runs.slot_wait must be greater than or equal to 0`,
   and leave the init templates unchanged.
3. MUST state the hazard in the code's behavior, as API Contract 6 does. The
   data sources are the Implement step's refusal, the queue's recorded limits
   and the slot reader; the destination is the item's blocker text in the Run
   Database, which `roundfix deliver status` prints, and the Queue Owner's
   console log. The wait never starts a Run itself: it only calls `RunSpec`
   again, so the store's `CreateRun` transaction stays the only authority on
   the ceiling. It never waits past `runs.slot_wait` or the queue deadline,
   and a cancelled owner clears the wait text through a context without
   cancellation.
4. MUST keep today's behavior byte for byte when `RunSlots` is unset,
   `RunSlotWait` is zero, or `RunSpec` fails with anything but a
   `*RunSlotRefusal`; at the wait's own deadline the park reason equals the
   no-wait reason of the same refusal. Every existing delivery and
   configuration test keeps its name and assertions; this Task declares no
   break.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, and declaring no real-time budget under ten seconds
   (waits use the fake clock and sleeper; the race test waits on channel
   events through the `internal/testwait` helpers, bounded only by the test
   deadline, the harness wait budget):
   - in `internal/config/run_slot_wait_test.go`,
     `TestRunSlotWaitIsReadFromUserConfig`,
     `TestProjectConfigCannotSetTheRunSlotWait` and
     `TestRunSlotWaitRefusesANegativeOrInvalidValue`;
   - in `internal/delivery/run_slot_wait_test.go`,
     `TestARunCeilingRefusalWaitsAndStartsWhenASlotFrees`,
     `TestARunSlotWaitParksAsTodayAtItsDeadline`,
     `TestTheQueueDeadlineEndsARunSlotWait`, `TestAZeroRunSlotWaitParksAtOnce`,
     `TestOtherRunSpecFailuresDoNotWait`,
     `TestAStoppedOwnerClearsTheRunSlotWait` and
     `TestTwoQueuesWaitingForOneSlotStartOneRun`, with the cases
     `_techspec.md` → Testing Approach lists. The race test runs two engines
     over two Delivery Queues in one Run Database whose fake runners create
     their Run through `store.CreateRun` with `MaxActiveImplementRuns: 1`, and
     asserts that exactly one Run is created after the holder completes and
     that Active Implement Runs never exceed one.
6. MUST document `runs.slot_wait` in `docs/user-guide/configuration.md`: in
   the Active Implement Run ceiling section, the key, its default `2h`, `0`
   parking at once, User Config only with the Project Config warning, that
   only the Delivery Queue waits (`roundfix implement` itself still refuses),
   and that an older Roundfix binary refuses a User Config that sets it, with
   that section's YAML example gaining the line `slot_wait: 2h` under `runs:`
   (the Full example stays as it is); and a
   `runs.slot_wait` row in the settings table beside `runs.max_active`.
7. MUST re-record the Behavior Surface Record, because this Task edits a user
   guide page: `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; no record is hand-edited.

## Subtasks

- [ ] Read, apply and validate `runs.slot_wait`.
- [ ] Add the Run slot types, the wait text and the engine dependencies.
- [ ] Run the wait in `runCandidate` with both deadlines and the cancellation path.
- [ ] Add the configuration, engine and race tests.
- [ ] Document the key and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A ceiling refusal keeps the item `running` with the wait text, and the Run starts after a slot frees.
- [ ] The wait parks with today's exact reason at `runs.slot_wait`, as `queue-deadline` when the queue's deadline is earlier, and at once with `0` or another failure.
- [ ] Two queues waiting for one slot create exactly one Run.
- [ ] `runs.slot_wait` defaults to `2h`, is ignored with a warning in Project Config, and is documented.

## Consumers

- `delivery.RunSlotRefusal`, `delivery.RunSlots`, `delivery.RunSlotReader`, `delivery.RunSlotWaitText` and the new `EngineDependencies` fields — consumers: `./internal/delivery`, `./internal/cli`
- `config.Runs.SlotWait` — consumers: `./internal/config`, `./internal/cli`

## Context

- instruction: `docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md`
- instruction: `docs/adr/0199-a-queue-token-ceiling-stops-new-starts-and-never-a-running-task.md`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/delivery/token_ceiling_test.go`
- instruction: `internal/config/config_test.go`
- instruction: `internal/store/store.go`
- interface: `internal/config/config.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/config/run_slot_wait_test.go`
- creates: `internal/delivery/run_slot_wait_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunSlotWaitIsReadFromUserConfig TestProjectConfigCannotSetTheRunSlotWait TestRunSlotWaitRefusesANegativeOrInvalidValue TestProjectConfigCannotRaiseTheRunCeiling; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARunCeilingRefusalWaitsAndStartsWhenASlotFrees TestARunSlotWaitParksAsTodayAtItsDeadline TestTheQueueDeadlineEndsARunSlotWait TestAZeroRunSlotWaitParksAtOnce TestOtherRunSpecFailuresDoNotWait TestAStoppedOwnerClearsTheRunSlotWait TestTwoQueuesWaitingForOneSlotStartOneRun TestAQueuedItemParksAtTheTokenCeilingWithoutAWorktree; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'runs.slot_wait' docs/user-guide/configuration.md || { echo 'missing runs.slot_wait in the configuration guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'slot_wait: 2h' || { echo 'missing the slot_wait example' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/config ./internal/delivery ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the configuration and engine tests pass with the existing ceiling, token-ceiling tests and the whole delivery package, the configuration guide documents the key with its example, the packages build for Windows, and the Behavior Surface Record is current.

## References

- `_prd.md` → Sources; Goals; Core Feature 3; Core Feature 5; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 5; API Contract 6; API Contract 9; Invariants 4-7; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 3
- ADR-0304; ADR-0199; ADR-0279; ADR-0259; ADR-0271; ADR-0285; ADR-0280
