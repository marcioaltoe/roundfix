---
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
prd: _prd.md
created: 2026-11-05
---

# A queue item shows its turn, and the pre-PR review gets its own capacity — Technical Spec

## Executive Summary

Four slices. First, the documents and ADR-0365, which state the behavior the
code slices then implement. Second, `delivery.review_concurrency` (User
Config, default `2`) and the `Limits:` line that prints it. Third, the
Delivery Engine runs the
pre-PR review of items in `reviewing` in review slots owned by one `Run`
call, up to the Review Capacity, while its loop keeps advancing one item at a
time through every other stage; a queued item waits while an earlier item is
in review, a poller hands retried items to the slots every 30 seconds, and a
reviewed item continues without its worktree being prepared again. Fourth,
the engine records each item's Queue Turn through a recorder the command
writes as the Turn Record beside the owner's console log, and `deliver
status` prints one `Turn:` line per item for an owner it proves live. The trade-off this design accepts is that a
reviewed item still waits for the owner's turn for its archive, gate,
publication and merge, and that the next queued item still waits for every
earlier item to merge, so a fresh queue gains nothing in the normal flow;
in exchange the queue keeps starting each item from a default branch that
holds every earlier merged item, and no schema, stage or column changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Turn
  Record sits in the directory named by the existing Queue Owner id
  (`delivery-<12 hex of the Git root>-<pid>`), and items keep their Spec
  slugs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the owner and `deliver status` read the Run Database,
  the configuration and a local file. Source: `docs/agents/go.md`,
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0304: "A new `waiting` stage or an item column was rejected: a schema change makes
  older item binaries refuse the shared Run Database."; the turn lives in a
  file, and no stage, column or schema changes. The prerequisite check of a queued item runs first,
  as today, and the turn names its wait. ADR-0279: "Status stays read-only and never reclaims the owner record"; status
  reads the Turn Record only for an owner it proves live, resolves the
  artifact directory without creating it and writes nothing. ADR-0199:
  "Tokens outside Runs, such as the pre-PR review, do not count."; the
  ceiling check stays where it is. ADR-0225: "Before each of those three children it runs the
  build in the item worktree"; each review builds its item binary in its own
  item worktree and copies it to its own per-Spec path, as today. ADR-0311: "a User
  Config key, 60 minutes by default, that the operator changes without a
  release."; `delivery.review_concurrency` follows the same operator-only
  form, and the Check Wait is unchanged. ADR-0237: "A Delivery Retry records a merge made outside the queue"
  carries no obligation: the retry's merge observation does not change. ADR-0365,
  written by task_01, records the decision below. Source:
  `docs/agents/domain.md`, `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 in the authorization record; bounded files
  are the Roundfix skill files task_01 changes. No Makefile, workflow,
  `go.mod`, `go.sum`, lint, formatter or Baseline module change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `ad77a797`.

1. `runDeliveryOwner` (`internal/cli/deliver.go`) claims the queue and calls
   `Engine.Run` in a loop until `ReleaseIdleDeliveryQueueOwner` finds no item
   outside `merged` and `parked`.
2. `Engine.Run` (`internal/delivery/engine.go`) reads the queue once and, for
   each item in position order, calls `advanceItem`, which loops
   `runCandidate`, `reviewCandidate`, `archiveCandidate`, `gateCandidate`,
   `publishCandidate`, `checkCandidate` and `mergeCandidate` until the item is
   `merged` or `parked`. The pre-PR review runs inline, through
   `commandDeliveryWorkflow.Review`, which runs `roundfix review` in the item
   worktree; the review session is read-only (`agent.SessionAccessReadOnly`).
3. `Engine.Retry` of a `review-findings`, `corrective-spec-required`,
   `pull-request-conflict` or accepted archived park, and of a park whose
   Tasks all completed, stores the stage `reviewing`, and `deliver retry`
   prints `Handed <slug> to Delivery Queue owner PID <pid>.`. The owner
   reaches the item only when its loop comes back to its position, after the
   item it is advancing, and every later item, merged or parked. This is what
   Fiscus saw: five items at `reviewing` while one Implement Run held the
   owner.
4. A stage is the step an item is at, written when the step is chosen. No
   record says whether that step is running; `deliver status` prints the
   stored stage, blocker and worktree, `Warning:` lines, an `Owner:` line for
   an owner proven gone (`deliveryOwnerLine`), `Park:` lines and
   `Limits: deadline <d>, retries per item <r>, concurrency 1, tokens <t>`
   from `printDeliveryLimits`, whose `concurrency 1` is a literal.
5. `advanceItem` of a queued item checks prerequisites first: unmet ones
   still ahead in the queue log `roundfix: prerequisite wait: …` and return
   with no change.
6. The Run Database schema is version 23; `delivery_queue_items` has no
   column for a wait. `delivery.check_timeout`, `runs.max_active` and
   `runs.slot_wait` are operator-only keys: `parseConfig` removes them from a
   Project Config with `warnIgnoredProjectSetting`.
7. A prototype on a scratch copy of `ad77a797` implemented API Contracts 3
   and 4 below; all 363 tests of `./internal/delivery` and the 452 tests of
   `./internal/cli` matched by `Deliver|Delivery|Queue` passed unchanged.

## System Architecture

- `internal/config`: the `delivery.review_concurrency` key.
- `internal/delivery`: the review slots (`review_slots.go`), the queued-start
  rule and the slot hand-off in `engine.go`, and the Queue Turn types, the
  recording and `TurnLines` (`turn.go`).
- `internal/cli`: the `Limits:` line, the engine wiring, the Turn Record
  writer and reader (`deliver_turn.go`) and the `Turn:` lines of `deliver
  status`.
- Documents: the Roundfix skill's `deliver` reference, the `deliver` and
  configuration user guides, `CONTEXT.md` and ADR-0365.

## Implementation Design

### Interfaces

```go
// internal/config
const DeliveryReviewConcurrencyKey = "delivery.review_concurrency"
const DefaultDeliveryReviewConcurrency = 2
type Delivery struct { /* existing fields */; ReviewConcurrency int }

// internal/delivery
type EngineDependencies struct { /* existing */; ReviewCapacity int; TurnPollInterval time.Duration; TurnRecorder TurnRecorder }
type TurnState string // "advancing", "reviewing", "review-slot-wait", "prerequisite-wait", "turn-wait"
type ItemTurn struct { SpecSlug string; Stage store.DeliveryStage; State TurnState; Prerequisites []string }
type Turn struct { ReviewCapacity int; Items []ItemTurn }
type TurnRecorder interface { RecordTurn(ctx context.Context, turn Turn) error }
func TurnLines(queue store.DeliveryQueue, turn Turn) []string

// internal/cli
func printDeliveryLimits(output io.Writer, limits store.DeliveryQueueLimits, reviewConcurrency int)
type ownerTurnRecorder struct { loaded roundconfig.Loaded /* pid, identity, path resolved once */ }
func turnRecordPath(artifactDir, gitRoot string, pid int) string
// internal/delivery, engine.go: called by the review slots on each take, start, end and poll
func (engine *Engine) reviewSlotsChanged()
```

### Data Models

No Run Database change. The Turn Record is a JSON file:

```json
{"schema":"roundfix/delivery-turn/v1","owner_pid":4242,"owner_identity":"<identity>","review_capacity":2,
 "items":[{"spec":"0066-alpha","stage":"reviewing","turn":"reviewing"},
          {"spec":"0080-echo","stage":"queued","turn":"prerequisite-wait","prerequisites":["0079-foxtrot"]}]}
```

`prerequisites` is omitted when empty; `items` lists, in queue order, every
item that is neither `merged` nor `parked` when the record is written.

### API Contracts

1. API Contract: the key. `Delivery.ReviewConcurrency` defaults to
   `DefaultDeliveryReviewConcurrency` (`2`). `parseConfig` removes
   `delivery.review_concurrency` from a Project Config, valid or not, with
   `warnIgnoredProjectSetting(DeliveryReviewConcurrencyKey)`, before strict
   decoding, so it prints
   `config: delivery.review_concurrency in Project Config is ignored; set delivery.review_concurrency in User Config`.
   `deliveryOverlay` gains `ReviewConcurrency *int` with YAML key
   `review_concurrency`, applied only when the source is not Project Config.
   `Validate` refuses a value below `1` with
   `delivery.review_concurrency must be greater than 0`; a non-integer fails
   strict decoding. The User and Project templates are unchanged.
2. API Contract: the `Limits:` line. `printDeliveryLimits` prints
   `Limits: deadline <d>, retries per item <r>, concurrency 1, review concurrency <n>, tokens <t>`
   with `<n>` the `reviewConcurrency` argument; `deliver start` and `deliver
   status` pass `loaded.Config.Delivery.ReviewConcurrency`.
3. API Contract: the wiring. `EngineDependencies.ReviewCapacity` below `1`
   means `1`; `TurnPollInterval` at or below `0` means `30s`; a nil
   `TurnRecorder` records nothing. `commandDeliveryEngineDependencies` sets
   `ReviewCapacity: loaded.Config.Delivery.ReviewConcurrency` and (API
   Contract 7) `TurnRecorder: &ownerTurnRecorder{loaded: loaded}`.
4. API Contract: the review slots.
   - `Run` starts the review slots and a poller that belong to that call, and
     returns only after both have ended; on any return it stops the poller,
     drops waiting items, cancels the reviews' context and waits for every
     review to end.
   - Each pass reads the queue and goes through it in position order as
     today. It first continues an item the slots returned (below), whatever
     stage the pass read for it. It skips an item at `reviewing` that the
     slots hold (waiting, in flight or claimed); otherwise it claims it, runs the existing resume
     prologue of `advanceItem` and hands it to the slots, releasing the claim.
     When advancing reaches `reviewing`, the item is handed to the slots and
     the pass goes on to the next item; `reviewCandidate` is no longer called
     by the pass.
   - The slots start waiting items in position order while fewer than
     `ReviewCapacity` reviews run. A review calls `reviewCandidate` unchanged
     with its item. An error while the context is live parks the item as
     `delivery-error: <error>`, the blocker the pass writes today; an error
     after cancellation leaves the item at `reviewing`. An item handed to the
     slots while every slot is held logs once
     `roundfix: Delivery Queue item <slug>: waiting for a review slot; <held> of <capacity> review slot(s) held`.
     A review that ends with the item at `archiving` returns the item to the
     pass, which continues it with the stage loop alone (no
     `UseItemBranch`).
   - The slots and the poller live in `review_slots.go`. They report each
     take, start and end of an item and each poll to the engine through one
     method in `engine.go`, `Engine.reviewSlotsChanged()`, which wakes a
     pass waiting on the slots; the slots answer `state(slug)` with `none`,
     `claimed`, `waiting` or `in-flight` for the pass and the engine.
   - Every `TurnPollInterval` the poller reads the queue and hands to the
     slots each item at `reviewing` with `WorktreeProvisioned`, a non-empty
     `Worktree` whose directory exists, and no hold or claim; any other item
     waits for a pass.
   - A queued item, after the prerequisite check and before
     `CreateItemBranch`, re-reads the queue and returns with no change and no
     log line when an item at a lower position is `running`, `reviewing`,
     `archiving`, `gating`, `publishing`, `checking` or `merging`, or when
     this rule already held an earlier queued item in the same pass, so a
     later item never starts ahead of an earlier one whose blocker ended in
     between. A queued item that waits for its prerequisites holds no later
     item, as today. (Without the second clause, a prototype run under the
     race detector started item 5 of `TestDeliveryEngineParksDeclaredBlockersAndContinues`
     ahead of item 4.)
   - After a pass, while the slots hold an item and none was returned, `Run`
     waits for a review to end, the poll interval or cancellation, then
     passes again; it passes again at once when an item was returned, and it
     returns, with the items of a fresh queue read, after a pass during which
     the slots stayed empty and nothing was handed to them or returned. A
     pass after the first skips items that were `merged` when it began, so
     merged cleanup still runs once per `Run` call.
5. API Contract: recording the turn. With a recorder, the engine writes a
   `Turn`, from `engine.go` and `turn.go` alone (the slot events arrive
   through `Engine.reviewSlotsChanged`), whenever the pass starts or stops advancing an item, the slots take,
   start or end an item, a queued item waits for prerequisites, and on every
   poll. It builds the `Turn` from a fresh queue read, serialized by one
   mutex: for each item neither `merged` nor `parked`, in queue order, its
   stored stage and `advancing` for the item the pass advances, `reviewing`
   for an item in flight, `review-slot-wait` for an item waiting in the slots,
   `prerequisite-wait` with the unmet slugs for a queued item whose last
   check found prerequisites still ahead, and `turn-wait` otherwise;
   `ReviewCapacity` is the effective capacity. A failed read or record logs
   `roundfix: warning: Delivery Queue turn record: <error>` and delivery goes
   on.
6. API Contract: `TurnLines`. For each item of `queue` neither `merged` nor
   `parked`, in queue order, one line:
   - no entry for its slug, or an entry other than `advancing` whose stage
     differs from the item's: `Turn: <slug> waiting for the owner to read it`;
   - `advancing`: `Turn: <slug> advancing at <item stage>`;
   - `reviewing`: `Turn: <slug> reviewing; <held> of <capacity> review slot(s) held`;
   - `review-slot-wait`:
     `Turn: <slug> waiting for a review slot (<n> of <m> waiting); <held> of <capacity> review slot(s) held`;
   - `prerequisite-wait`: `Turn: <slug> waiting for prerequisites: <a>, <b>`;
   - `turn-wait`: `Turn: <slug> waiting for the owner's turn (<n> of <m> waiting)`.
   `<held>` counts the `reviewing` entries whose stage matches; `<capacity>`
   is `turn.ReviewCapacity`. For a review slot, `<n>` is the place by
   position among the matching `review-slot-wait` entries and `<m>` their
   count. For the owner's turn, `<m>` counts the matching `turn-wait` entries
   and `<n>` is the place in the order the pass will visit them: the items
   after the `advancing` item's position first, then the items before it,
   each by position; with no `advancing` entry, by position.
7. API Contract: the Turn Record. `turnRecordPath` is
   `<artifact dir>/delivery/<deliveryOwnerID(gitRoot, pid)>/turn.json`, the
   console log's directory. `ownerTurnRecorder.RecordTurn` resolves the path
   for `os.Getpid()`, reads `store.OwnerProcessIdentity` once, creates the
   directory (`0o700`) when absent, writes the JSON of Data Models to a
   temporary file in it (`0o600`) and renames it over `turn.json`.
8. API Contract: status. `deliveryOwnerLine` also reports whether the owner
   is proven live. When it is, status resolves the artifact directory with
   `roundconfig.ResolveArtifactDirectory` (creating nothing) and reads the
   Turn Record for the queue's owner PID: a missing file prints nothing; a
   file that cannot be read or decoded, or whose `schema` is not
   `roundfix/delivery-turn/v1`, writes
   `roundfix: Delivery Queue turn record <path> could not be read: <error>`
   to stderr; a record whose `owner_pid` or `owner_identity` differs from the
   queue's prints nothing; otherwise status prints `TurnLines` after the
   `Warning:` lines and before the `Park:` lines. Status exits `0` in every
   case and writes no file.

### Surface Transcripts

1. Surface Transcript: a live owner's turns, for a queue whose owner is the
   test process proven live, with a Turn Record of capacity 2 listing
   `0066-alpha` (`reviewing`, reviewing), `0068-bravo` (`running`,
   advancing), `0070-charlie` (`reviewing`, reviewing), `0074-delta`
   (`reviewing`, review-slot-wait), `0080-echo` (`queued`, prerequisite-wait
   on `0079-foxtrot`) and `0081-golf` (`queued`, turn-wait), and the item
   `0083-hotel` at `reviewing` absent from the record; no User Config.

   ```transcript
   $ roundfix deliver status
   stdout:
   0066-alpha	reviewing	-	/worktrees/0066-alpha
   0068-bravo	running	-	/worktrees/0068-bravo
   0070-charlie	reviewing	-	/worktrees/0070-charlie
   0074-delta	reviewing	-	/worktrees/0074-delta
   0080-echo	queued	-	-
   0081-golf	queued	-	-
   0083-hotel	reviewing	-	/worktrees/0083-hotel
   Turn: 0066-alpha reviewing; 2 of 2 review slot(s) held
   Turn: 0068-bravo advancing at running
   Turn: 0070-charlie reviewing; 2 of 2 review slot(s) held
   Turn: 0074-delta waiting for a review slot (1 of 1 waiting); 2 of 2 review slot(s) held
   Turn: 0080-echo waiting for prerequisites: 0079-foxtrot
   Turn: 0081-golf waiting for the owner's turn (1 of 1 waiting)
   Turn: 0083-hotel waiting for the owner to read it
   Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
   Usage: no Runs recorded
   stderr:
   exit: 0
   ```

## Invariants

1. No Run Database schema, Delivery Stage, item column, item line of `deliver
   status`, blocker text, Park Class or exit code changes.
2. One Queue Owner runs at most `ReviewCapacity` reviews at once, and never
   two reviews of one item.
3. The pass never calls the workspace for an item the slots hold.
4. A queued item never creates its item branch while an item at a lower
   position is between `running` and `merging`, nor before an earlier queued
   item that does not wait for its prerequisites.
5. Every review outcome parks or advances its item exactly as
   `reviewCandidate` does today; a cancelled review leaves `reviewing`.
6. `Run` returns only after every review it started and its poller ended.
7. `deliver status` prints `Turn:` lines only for an owner it proves live
   whose Turn Record names the same PID and identity, and writes nothing.
8. A failed turn record never parks, stops or delays an item.

## Declared breaks

- task_02 changes the `Limits:` line, so the expected text of these tests
  gains `review concurrency 2, ` before `tokens`, and nothing else changes:
  `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`
  (`deliver_lane_skip_visible_test.go`), `TestDeliverStartRecordsAndPrintsItsLimits`
  and `TestDeliverStartRecordsNoneForOmittedLimits` (`deliver_limits_test.go`),
  `TestDeliverStatusPrintsOwnerAfterWarningsAndBeforeParks` and the helper
  `assertDeliveryOwnerStatus` (`deliver_owner_status_test.go`),
  `TestDeliverStatusPrintsAParkLinePerParkedItem` and
  `TestDeliverStatusReproducesSurfaceTranscriptOne` (`deliver_park_status_test.go`),
  `TestDeliverStatusPrintsAnItemWarning` and
  `TestDeliverStatusPrintsNoWarningLineWithoutAWarning` (`deliver_revalidate_test.go`),
  `TestDeliverStatusShowsARunSlotWaitTranscript` (`deliver_run_slot_wait_test.go`),
  `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark`
  (`deliver_run_unresolved_park_test.go`), `TestDeliverStatusPrintsTheItemWorktree`
  (`deliver_test.go`), `TestDeliverStartRecordsAndPrintsTheTokenCeiling` and
  `TestDeliverStatusPrintsTheTokenCeilingQuestion` (`deliver_token_ceiling_test.go`),
  and `TestDeliverStartPrintsTokensNone`, `TestDeliverStatusPrintsNoRunsRecorded`
  and `TestDeliverStatusPrintsTheUsageLine` (`deliver_usage_test.go`), all in
  `internal/cli`.

No other test is expected to change; the prototype of Current behavior 7
proves it for the review slots, and the `Turn:` lines need a Turn Record no
existing test writes.

## ADR-0365

task_01 creates `docs/adr/0365-the-pre-pr-review-gets-its-own-capacity-and-status-shows-each-items-turn.md`
with exactly this content, apart from line wrapping:

```markdown
---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# The pre-PR review gets its own capacity, and status shows each item's turn

The Queue Owner runs the pre-PR review of items in `reviewing` beside the one
item it advances through every other stage, up to a Review Capacity set by
`delivery.review_concurrency` in User Config, `2` by default. It records each
item's Queue Turn in a Turn Record beside its console log, and `roundfix
deliver status` prints one `Turn:` line per item for an owner it proves live.

## Context

On 2026-11-05 Fiscus reported a queue of about twenty-one Specs running for
two days, with five items in `reviewing` and no review running while one
Implement Run held the owner, and `deliver status` saying nothing about the
wait. The owner advanced one item at a time through every stage, the review
inline; a Delivery Retry stored `reviewing` and the owner reached the item
only when its loop came back to it. The maintainer chose: "O deliver status
mostra 'aguardando capacidade' e a posição. O review pré-PR, que só lê, ganha
uma concorrência própria, separada da do implement."

## Decision

**A stage is the next step.** The stored stage keeps its name and meaning;
whether the step runs or waits is the item's Queue Turn.

**Reviews run in slots.** Each `Run` of the owner holds review slots and a
poller. The owner's loop hands an item that reaches `reviewing` to the slots
and goes on; the slots review waiting items in position order, up to the
Review Capacity. The poller hands every retried item at `reviewing` whose
item worktree exists to the slots every 30 seconds, so a review never waits
for another item's Implement Run. A review reads only its own item worktree
in a read-only session, builds its item binary into its own path, and the
owner never prepares the worktree of an item the slots hold. A reviewed item
continues in the owner's loop without preparing its worktree again.

**Order is kept.** A queued item starts only when no earlier item is between
its Implement Run and its merge, and never ahead of an earlier queued item,
so it still starts from a default branch holding every earlier merged item. Archive, gate, publication, check wait
and merge stay one item at a time.

**The knob is the operator's.** `delivery.review_concurrency` is a User
Config integer of at least `1`; Project Config is ignored with the warning of
the other operator-only keys. The owner reads it when it starts. `deliver
start` and `deliver status` print it on the `Limits:` line as
`review concurrency <n>`.

**The turn is a file.** The owner writes the Queue Turn of every item it has
not finished, its PID, start-time identity and Review Capacity to
`turn.json` beside its console log whenever a turn changes and on every poll.
Status reads it only for the owner it proves live and only when the PID and
identity match, and writes nothing.

## Considered Options

- A `waiting` stage or an item column: rejected, as in ADR-0304, because a
  schema change makes older item binaries refuse the shared Run Database.
- Deriving the turn in status from the Run Database alone: rejected, because
  nothing there says which item the owner is advancing.
- Starting the next queued item while an earlier one is in review: rejected,
  because the next item would start from a default branch without the
  earlier item, and Fiscus already pays for Specs that start from an older
  main.
- A Project Config key or a `deliver start` flag: rejected, because the
  reviews spend the operator's provider sessions on the operator's machine,
  and a flag recorded per queue needs a schema change.
- A default of `1`: rejected, because one slot leaves a batch of retried
  items in a line; two reviews beside one Implement Run of Task Capacity 2
  stay within four Agent sessions, fewer than the machine runs under the
  default ceiling of three Active Implement Runs.

## Consequences

A retried item is reviewed while another item implements, and several waiting
reviews run together. A reviewed item still waits for the owner's turn for its
archive and merge, and a fresh queue reviews one item at a time as before.
Review tokens stay outside the Queue Token Ceiling (ADR-0199), and reviews
count neither against `runs.max_active` nor a Run's Verification Capacity. An
older owner writes no Turn Record, so its status shows no `Turn:` line until
it is restarted on a binary that writes one. An older Roundfix binary,
including an item binary built from a branch older than this decision,
refuses a User Config that sets `delivery.review_concurrency` as an unknown
key; the default needs no setting.
```

## Coverage Map

- Goal 1, Story 1, Story 5, Core Features 1–3, Success Metric 4 → `Turn`,
  `TurnLines`, the engine recording, `ownerTurnRecorder`, `turnRecordPath`
  and the status reader (task_04).
- Goal 2, Story 2, Story 3, Core Feature 4, Success Metrics 1–2 → the review
  slots, the poller and the continuation of a reviewed item (task_03).
- Goal 3, Story 4, Core Features 5 and 7, Success Metric 3 → the queued-start
  rule and the unchanged review outcomes (task_03).
- Goal 4, Core Feature 6, Success Metric 5 → `delivery.review_concurrency` and
  the `Limits:` line (task_02).
- Goal 5, Core Feature 8, Success Metric 6 → the `deliver` guide paragraphs
  on the stage meaning and the `Turn:` lines, the configuration guide's
  `delivery.review_concurrency` row and section, the Roundfix `deliver`
  reference paragraph that tells Agents to read the `Turn:` lines and its
  mirror, the four `CONTEXT.md` entries, and ADR-0365 (task_01).

Deviations from the PRD: None.

## Absence Evidence

1. `delivery.review_concurrency` — command: `grep -rn 'review_concurrency\|ReviewConcurrency' internal docs/user-guide`; result: no match.
2. The review slots and the poller — command: `grep -rn 'go func\|sync\.' internal/delivery/engine.go`; result: no match; `reviewCandidate` is called only from the stage loop of `advanceItem`.
3. The queued-start rule — command: `grep -n 'DeliveryStageQueued' internal/delivery/engine.go`; result: the queued branch of `advanceItem` checks prerequisites, the archived proof and revalidation, and reads no earlier item's stage.
4. The Queue Turn and `TurnLines` — command: `grep -rn 'TurnLines\|TurnRecorder\|turn.json\|Turn:' internal`; result: no match.
5. The Turn Record reader in status — command: `grep -n 'ResolveArtifactDirectory\|delivery-turn' internal/cli/deliver.go`; result: no match.
6. ADR-0365 and the glossary terms — command: `ls docs/adr | grep '^0365'; grep -n 'Review Capacity\|Queue Turn\|Turn Record' CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0304: it rejected a distinct exit code, retrying the child on every
  poll, a `waiting` stage or an item column, a Project Config key and an
  unbounded wait; this design takes none of them: the wait is a file, the
  poller reads the Run Database instead of starting children, and the key is
  User Config only.
- ADR-0279: it records no rejected alternative; status stays read-only and
  adds the `Turn:` lines only where it adds no `Owner:` line.
- ADR-0199: it records no rejected alternative; the ceiling still bounds new
  starts and retries, and review tokens still do not count.
- ADR-0225: it rejected accepting an unknown key in an older binary and
  re-executing the whole owner as the item binary; this design takes neither,
  and each review keeps building and copying its own item binary.
- ADR-0311: it rejected keeping a constant and a Project Config key; this
  design keeps the Check Wait and makes the new key User Config only too.
- ADR-0237: it records no rejected alternative this design touches; the
  retry's merge observation is unchanged.
- ADR-0365: written here; it takes none of its own rejected options.

## Integration Points

No network, provider or other repository. The reviews run `roundfix review`
as today, one child per item worktree; `git worktree` administration stays
behind the existing worktree admin lock, and the pass alone calls the
workspace.

## Testing Approach

Engine tests drive `Engine.Run` over the real Run Database with fakes that
guard their state with a mutex and signal through channels, and wait with
`testwait.Until`; no test waits on a short real-time timer, and the poll
interval is set through `EngineDependencies.TurnPollInterval`, which the
Real-Time Budget rule does not count. Seeded items at `reviewing` carry a
`Worktree` under `t.TempDir()` and `WorktreeProvisioned`. Command tests seed
queues with `store.CreateDeliveryQueue`, `UpdateDeliveryQueueItem` and
`ClaimDeliveryQueueOwner`, prove the owner with `fakeDeliveryOwnerProcesses`
as `deliver_owner_status_test.go` does, and write the Turn Record at
`turnRecordPath`. The Project Config test marks its work directory with
`gittest.MarkRepository` from Spec 0352.

Sabotages, each on a temporary copy, with the test and the input that sees
it:

1. Call `reviewCandidate` inline instead of handing the item to the slots:
   `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements` fails (run
   with `-timeout 3m`): the earlier item's Implement Run never sees the review
   start.
2. Start every waiting item regardless of capacity:
   `TestReviewCapacityBoundsConcurrentReviews` fails on its three items with
   capacity 2.
3. Drop the queued-start rule: `TestAQueuedItemWaitsForAnEarlierItemInReview`
   fails: the later item creates its branch while the earlier review runs.
   Drop only its same-pass clause: its subtest `an earlier blocker ends
   during the pass` fails: `third` creates its branch before `second`.
4. Let the pass prepare an item the slots hold:
   `TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements` fails on its
   count of `UseItemBranch` calls for the reviewed item.
5. Park a review that ends by cancellation:
   `TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing` fails.
6. Rank `turn-wait` items by position alone: `TestTurnLinesNameEachItemsTurn`
   fails on its input where the advancing item sits between two waiting ones.
7. Print the turn lines without comparing the owner identity:
   `TestDeliverStatusIgnoresATurnRecordOfAnotherOwner` fails.
8. Keep `concurrency 1` as the whole line: `TestDeliverLimitsNameTheReviewConcurrency`
   fails on its User Config of `3`.
9. Apply a Project Config value: `TestProjectConfigCannotSetTheDeliveryReviewConcurrency`
   fails.

## Build Order

1. ADR-0365, `CONTEXT.md`, the Roundfix `deliver` reference and the `deliver`
   and configuration guides (task_01), written first so every later Task that
   changes the `deliver` command has its guide among its dependencies and
   implements the behavior those documents state.
2. `delivery.review_concurrency` and the `Limits:` line (task_02; depends on:
   1, for the guides of the `deliver` command).
3. The review slots, the poller, the queued-start rule and the engine wiring
   of the capacity (task_03; depends on: 2, for the configuration field).
4. The Queue Turn, its recording, the Turn Record and the `Turn:` lines of
   `deliver status` (task_04; depends on: 3, for the slot events it records
   through `Engine.reviewSlotsChanged` and the shared `engine.go` and
   `deliver_workflow.go`).
5. QA gate (depends on: 4).

Spec 0352 must merge first (`requires`): task_02's Project Config test uses
`gittest.MarkRepository`, and 0352 changes `internal/config/config.go`.
Spec 0356, delivered before this Spec, may change `Engine.Retry`,
`internal/cli/deliver.go` (the retry output) and the Roundfix `deliver` and
`review` references; Spec 0357, delivered after, may change the queued branch
of `advanceItem` where it revalidates. Each Task here rebases on whatever has
merged and keeps its edits to those functions narrow; the record commands
recompute versions and fingerprints.

## Risks & Considerations

- A reviewed item waits for the owner's turn behind a long Implement Run of
  another item that started before it was retried; the turn line says so.
  Running archive and publication beside the owner is a larger change left
  out of scope.
- Two reviews build two item binaries at once; the builds share the Go build
  cache, which is safe for concurrent builds, and write to per-Spec paths.
- Concurrent writes to the Run Database go through its existing write
  transactions; each review writes only its own item.
- The console log may interleave lines from the slots and the pass; each line
  is one write.
- A Turn Record can be up to 30 seconds old for an item retried while the
  owner is busy; the line `waiting for the owner to read it` covers it.
- An owner started on an older binary writes no record; status shows no turn
  line for it.

## Glossary

- adds: **Review Capacity**
- adds: **Queue Turn**
- adds: **Turn Record**
- changes: **Queue Owner**

## Decisions

- The review slots belong to one `Run` call, so the goroutines have an owner
  that waits for them, and the owner loop of `runDeliveryOwner` is unchanged.
- Only the pass calls the workspace; the poller hands over only items whose
  worktree a Delivery Retry already proved, and a reviewed item continues
  without the prologue. See ADR-0365.
- The queued-start rule logs nothing, so a normal queue's console log is
  byte-identical to today's.
- The Turn Record is a file beside the console log, read by status only for
  a proven live owner. See ADR-0365.
- The Review Capacity is User Config only, `2` by default. See ADR-0365.
