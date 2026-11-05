---
task: task_04
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: completed
type: backend
complexity: high
---

# Task 04: The owner records each item's Queue Turn, and deliver status prints the Turn lines

## Overview

`deliver status` prints each item's stored stage, which names the step the
item is at, and nothing tells whether that step runs or waits ("O `deliver
status` mostra `reviewing`, mas não diz que o item espera a vez", reported by
fiscus on 2026-11-05,
[references](references/2026-11-05-fiscus-delivery-queue-takes-days.md)).
After this Task the Delivery Engine records the Queue Turn of every item it
has not finished through a recorder, the command writes it as the Turn
Record beside the owner's console log, and `deliver status` prints one
`Turn:` line per item for an owner it proves live, in the six forms of
`_techspec.md` → API Contract 6.

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestDeliverStatus.*Owner.*|TestDeliverStatusPrintsTheItemWorktree|TestDeliverStatusPrintsAParkLinePerParkedItem)$' ./internal/cli`
   and record in the Result that it passes; these pin the status lines around
   the new ones.
2. MUST create `internal/delivery/turn.go` with `TurnState` and its five
   values, `ItemTurn`, `Turn`, `TurnRecorder` and `TurnLines` exactly as
   `_techspec.md` → Interfaces and API Contract 6 state, and add the
   `TurnRecorder` field to `EngineDependencies` and `Engine` (API Contract 3).
3. MUST record the turn from `internal/delivery/engine.go`, through
   `Engine.reviewSlotsChanged()` and the pass, and from
   `internal/delivery/turn.go`, as API Contract 5 states, without changing
   `internal/delivery/review_slots.go`: at each event
   it names, from a fresh queue read, serialized by one mutex, with the
   advancing item, the slots' in-flight and waiting items and the
   prerequisite waits; a failed read or record logs exactly
   `roundfix: warning: Delivery Queue turn record: <error>` and delivery goes
   on. A nil recorder records nothing and adds no log line.
4. MUST create `internal/cli/deliver_turn.go` with `turnRecordPath` and
   `ownerTurnRecorder` as API Contract 7 states, writing the JSON of
   `_techspec.md` → Data Models, and set
   `TurnRecorder: &ownerTurnRecorder{loaded: loaded}` in
   `commandDeliveryEngineDependencies` (`internal/cli/deliver_workflow.go`).
5. MUST change `runDeliverStatus` and `deliveryOwnerLine` in
   `internal/cli/deliver.go` as API Contract 8 states: the `Turn:` lines come
   after the `Warning:` lines and before the `Park:` lines, only for an owner
   proven live whose Turn Record names its PID and identity; the artifact
   directory is resolved with `roundconfig.ResolveArtifactDirectory`; a record
   that cannot be read or decoded, or has another `schema`, writes exactly
   `roundfix: Delivery Queue turn record <path> could not be read: <error>`
   to stderr; status exits `0` and writes no file.
6. MUST create `internal/delivery/turn_test.go`, opening with the `Suite`,
   `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, each
   top-level test calling `t.Parallel()` first, with:
   - `TestTurnLinesNameEachItemsTurn`: a queue of the items `a` (`queued`),
     `b` (`running`), `c` (`gating`), `d` (`queued`), `e`, `f`, `g`
     (`reviewing`), `h` (`queued`), `i` (`archiving`) and `j` (`parked`) at
     positions 0 to 9, and a `Turn` of capacity 2 with `a`, `c`, `d`
     `turn-wait`, `b` `advancing`, `e` `reviewing`, `f` and `g`
     `review-slot-wait`, `h` `prerequisite-wait` on `0090-x` and `0091-y`, and
     `i` recorded as `reviewing` at stage `reviewing`, each other entry at its
     item's stage. `TurnLines` returns exactly, in this order:
     `Turn: a waiting for the owner's turn (3 of 3 waiting)`,
     `Turn: b advancing at running`,
     `Turn: c waiting for the owner's turn (1 of 3 waiting)`,
     `Turn: d waiting for the owner's turn (2 of 3 waiting)`,
     `Turn: e reviewing; 1 of 2 review slot(s) held`,
     `Turn: f waiting for a review slot (1 of 2 waiting); 1 of 2 review slot(s) held`,
     `Turn: g waiting for a review slot (2 of 2 waiting); 1 of 2 review slot(s) held`,
     `Turn: h waiting for prerequisites: 0090-x, 0091-y`,
     `Turn: i waiting for the owner to read it`.
     A subtest `nothing advancing` drops `b` from the `Turn` and expects `a`,
     `c` and `d` at places 1, 2 and 3 of 3 and `b` waiting for the owner to
     read it.
   - `TestTheOwnerRecordsEachItemsTurn`: `Engine.Run` with
     `ReviewCapacity: 1`, `TurnPollInterval` of `10 * time.Millisecond`, a
     recorder that keeps every `Turn` under a mutex, and a queue of
     `implementing` (`running`), `reviewed` and `waiting` (`reviewing`, with
     provisioned worktrees under `t.TempDir()`) and `later` (`queued`); the
     fake `RunSpec` of `implementing` blocks, through `testwait.Until`, until
     a recorded `Turn` has capacity 1 and `implementing` `advancing`,
     `reviewed` `reviewing`, `waiting` `review-slot-wait` and `later`
     `turn-wait`; then every item reaches `merged`. Subtest
     `prerequisite wait`: `blocker` (`reviewing`, its fake `Review` held) and
     `dependent` (`queued`, whose fake prerequisite reader reports `blocker`
     unmet until `blocker` merges) record `dependent` as `prerequisite-wait`
     with `Prerequisites` `[blocker]` before the review is released. Subtest
     `recorder failure`: a recorder returning `errors.New("disk full")`
     leaves every item `merged` and the console log carries
     `roundfix: warning: Delivery Queue turn record: disk full`.
7. MUST create `internal/cli/deliver_turn_test.go`, with the four comment
   lines and `t.Parallel()` first in each top-level test, seeding queues with
   `store.CreateDeliveryQueue`, `UpdateDeliveryQueueItem` and
   `ClaimDeliveryQueueOwner` and proving the owner with
   `fakeDeliveryOwnerProcesses` as `deliver_owner_status_test.go` does:
   - `TestDeliverStatusReproducesSurfaceTranscriptQueueTurns`: Surface
     Transcript 1 of `_techspec.md`, with the owner PID `os.Getpid()`, the
     identity `recorded-owner-identity` and the Turn Record it describes
     written at `turnRecordPath`, asserting every stdout line, every stderr
     line and the exit code.
   - `TestDeliverStatusIgnoresATurnRecordOfAnotherOwner`, with subtests
     `other identity` and `other pid` (no `Turn:` line), `owner not running`
     (PID `1 << 30`: the `Owner:` line and no `Turn:` line), `no record` (no
     `Turn:` line, and the resolved artifact directory still does not exist
     when it did not before), `unreadable record` (the file holds `{`: no
     `Turn:` line, stderr is exactly one line starting
     `roundfix: Delivery Queue turn record ` and containing
     ` could not be read: `, exit `0`) and `wrong schema` (the same stderr
     line).
   - `TestTheTurnRecorderWritesBesideTheConsoleLog`: `RecordTurn` of an
     `ownerTurnRecorder` over a test `roundconfig.Loaded` writes
     `<artifact dir>/delivery/<deliveryOwnerID(gitRoot, os.Getpid())>/turn.json`
     with `schema` `roundfix/delivery-turn/v1`, `owner_pid` `os.Getpid()`,
     `owner_identity` equal to `store.OwnerProcessIdentity` of that PID, the
     capacity and the items in order, no `prerequisites` key on an item
     without prerequisites, mode `0600`; a second call replaces the file and
     leaves no other file in the directory.
   - `TestTheDeliveryEngineGetsTheTurnRecorder`:
     `commandDeliveryEngineDependencies` returns a `TurnRecorder` that is an
     `*ownerTurnRecorder`.
8. MUST keep every existing test of `./internal/delivery` and every `deliver`
   test of `./internal/cli` passing unchanged.
9. MUST run sabotages 6 and 7 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails and on which input.
10. MUST NOT change the item line of `deliver status`, the `Owner:` and
    `Park:` lines, any blocker text, Delivery Stage, Park Class, the Run
    Database, the review slots' scheduling, the configuration, any document
    or skill, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the status lines around the new ones.
- [ ] Add the Queue Turn types and `TurnLines`.
- [ ] Record the turn from the pass, the slots, the prerequisite wait and the poller.
- [ ] Write the Turn Record beside the console log and wire the recorder.
- [ ] Print the `Turn:` lines for a proven live owner, and add the tests and the two sabotages.

## Acceptance Criteria

- [ ] `TurnLines` prints the six forms with their places in line exactly as API Contract 6 states.
- [ ] The owner records each item's turn, and a failed record never stops delivery.
- [ ] `deliver status` prints Surface Transcript 1 and prints no `Turn:` line for an owner not proven live, a record of another owner or no record.

## Context

- creates: `internal/delivery/turn.go`
- creates: `internal/delivery/turn_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/testwait/testwait.go`
- creates: `internal/cli/deliver_turn.go`
- creates: `internal/cli/deliver_turn_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_owner_status_test.go`
- interface: `internal/config/config.go`
- interface: `internal/store/process.go`

## Consumers

- `delivery.Turn`, `TurnRecorder` and `TurnLines` — consumers: `./internal/delivery`, `./internal/cli`
- The `deliver status` output and the Turn Record file — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTurnLinesNameEachItemsTurn TestTheOwnerRecordsEachItemsTurn TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 30m -v -run 'Deliver|Delivery|Queue|Turn' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverStatusReproducesSurfaceTranscriptQueueTurns TestDeliverStatusIgnoresATurnRecordOfAnotherOwner TestTheTurnRecorderWritesBesideTheConsoleLog TestTheDeliveryEngineGetsTheTurnRecorder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/delivery ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/delivery ./internal/cli` — expected: exit 0. Before this Task none of the six new tests exists, so the first pass line is missing; afterwards the turn lines, the recording, the Turn Record and the status reader behave as API Contracts 5 to 8 and Surface Transcript 1 state, every existing engine and `deliver` test passes, and both platforms build.

## References

- `_prd.md` → Goals; User Story 1; User Story 5; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 4
- `_techspec.md` → Current behavior 4; Data Models; API Contract 3; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Surface Transcript 1; Invariant 1; Invariant 7; Invariant 8; Testing Approach; Build Order 3; ADR-0279; ADR-0304


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-concurrency, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging.

Implemented this Task's turn recording and status slice. The pre-existing
`status: in_progress` is preserved for Daemon settlement. No existing test,
review-slot implementation, Task Graph, other Task, configuration, guide,
skill, database schema, stage, blocker or Park Class was changed.

Acceptance evidence:

- `TurnLines` renders all six forms in queue order, excludes merged and
  parked items, ignores stale non-advancing entries when counting slots and
  waits, and ranks the owner's waiting items after the advancing position
  before wrapping. `TestTurnLinesNameEachItemsTurn` passes the required
  ten-item input, `nothing advancing`, and an advancing entry whose stored
  stage changed.
- The engine records fresh queue reads at pass start/stop, slot callbacks,
  prerequisite waits and polls, with the advancing item and prerequisite
  state protected by the same mutex as the read and recorder call. The
  existing review slots remain unchanged. `TestTheOwnerRecordsEachItemsTurn`
  observes capacity 1 with `implementing` advancing, `reviewed` reviewing,
  `waiting` waiting for a review slot, and `later` waiting for the owner;
  every item then merges. Its `prerequisite wait` subtest observes
  `dependent` waiting on `[blocker]` before releasing the review, then
  preserves the existing next-Run prerequisite deferral. Its `recorder
  failure` subtest merges both items despite `disk full` and observes the
  exact required warning. A nil recorder returns before any recording read
  or warning.
- The command injects `*ownerTurnRecorder`, which caches the current PID,
  identity and path and atomically replaces `turn.json` beside the console
  log. `TestTheTurnRecorderWritesBesideTheConsoleLog` checks the v1 schema,
  literal JSON field names, identity, capacity, ordered items, omitted empty
  prerequisites, file mode `0600`, directory mode `0700`, and replacement
  with no temporary file left. `TestTheDeliveryEngineGetsTheTurnRecorder`
  checks the production dependency builder.
- `deliver status` reads turns only after proving the owner live and matching
  both recorded PID and identity. `TestDeliverStatusReproducesSurfaceTranscriptQueueTurns`
  checks every line of Surface Transcript 1, empty stderr and exit 0.
  `TestDeliverStatusIgnoresATurnRecordOfAnotherOwner` covers another identity,
  another PID, an absent owner, no record, malformed JSON, a wrong schema,
  unproven identity and unreadable owner identity. It checks that status
  preserves the queue, creates no artifact directory for a missing record,
  and emits exactly one diagnostic for an unreadable or unsupported record
  while exiting 0. `TestDeliverStatusOrdersWarningsTurnsAndParks` checks the
  exact ordering and unchanged adjacent lines.

Focused checks (all exit 0 unless explicitly identified as the red signal):

- Before changing code, the required characterization passed:
  `go test -count=1 -timeout 20m -v -run '^(TestDeliverStatus.*Owner.*|TestDeliverStatusPrintsTheItemWorktree|TestDeliverStatusPrintsAParkLinePerParkedItem)$' ./internal/cli`.
- Red signal after adding the engine tests, before implementation:
  `go test -count=1 -run '^TestTurnLinesNameEachItemsTurn$' ./internal/delivery`
  exited 1 with the absent `TurnState`, `Turn` and turn constants.
- `go test -count=1 -timeout 3m -v -run '^(TestTurnLinesNameEachItemsTurn|TestTheOwnerRecordsEachItemsTurn)$' ./internal/delivery`.
- `go test -count=1 -timeout 3m -run '^TestDeliverStatusPrintsNoOwnerLineForAProvenOwner$' ./internal/cli`.
- `go test -count=1 -timeout 3m -v -run '^(TestDeliverStatusReproducesSurfaceTranscriptQueueTurns|TestDeliverStatusIgnoresATurnRecordOfAnotherOwner|TestTheTurnRecorderWritesBesideTheConsoleLog|TestTheDeliveryEngineGetsTheTurnRecorder)$' ./internal/cli`.
- `go test -race -count=1 -timeout 5m -v -run '^(TestTurnLinesNameEachItemsTurn|TestTheOwnerRecordsEachItemsTurn|TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements|TestReviewCapacityBoundsConcurrentReviews|TestAQueuedItemWaitsForAnEarlierItemInReview|TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing)$' ./internal/delivery`.
- `go test -race -count=1 -timeout 5m -v -run '^(TestDeliverStatus.*Owner.*|TestDeliverStatusReproducesSurfaceTranscriptQueueTurns|TestDeliverStatusOrdersWarningsTurnsAndParks|TestDeliverStatusPrintsTheItemWorktree|TestDeliverStatusPrintsAParkLinePerParkedItem|TestTheTurnRecorderWritesBesideTheConsoleLog|TestTheDeliveryEngineGetsTheTurnRecorder)$' ./internal/cli`.
- After strengthening the independent JSON-field and directory-mode assertions:
  `go test -race -count=1 -timeout 3m -v -run '^TestTheTurnRecorderWritesBesideTheConsoleLog$' ./internal/cli`.
- `git diff --check` reports no whitespace errors. The changed-file scope is
  this Task file and the seven declared Go implementation/test files only.

Required sabotage evidence, each on the isolated copy at
`/private/tmp/roundfix-0355-turn-sabotage-vsdxuu76/repo`, with each mutation
restored after its probe:

- Sabotage 6 disabled the rotation after the advancing position, leaving
  position-only ranking. `go test -count=1 -timeout 3m -v -run '^TestTurnLinesNameEachItemsTurn$' ./internal/delivery`
  exited 1 in `TestTurnLinesNameEachItemsTurn`: the input with `b` advancing
  at position 1 and `a`, `c`, `d` waiting produced places 1, 2, 3 instead of
  3, 1, 2. Exact transcript:
  `/private/tmp/roundfix-0355-turn-sabotage-vsdxuu76/sabotage-6.log`.
- Sabotage 7 removed only the Turn Record owner-identity comparison, retaining
  the PID comparison. `go test -count=1 -timeout 3m -v -run '^TestDeliverStatusIgnoresATurnRecordOfAnotherOwner$/^other_identity$' ./internal/cli`
  exited 1 in `TestDeliverStatusIgnoresATurnRecordOfAnotherOwner/other_identity`:
  a record with the live owner's PID but `another-owner-identity` incorrectly
  printed `Turn: 0300-example advancing at checking`. Exact transcript:
  `/private/tmp/roundfix-0355-turn-sabotage-vsdxuu76/sabotage-7.log`.

The declared Verification, cross-platform builds and repository gates were
not run; they remain Daemon-owned. No commit, push or pull request was made.
