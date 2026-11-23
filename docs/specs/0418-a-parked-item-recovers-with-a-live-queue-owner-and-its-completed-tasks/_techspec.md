---
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
prd: _prd.md
created: 2026-11-23
---

# A parked item recovers with a live Queue Owner and its completed Tasks — Technical Spec

## Executive Summary

Four changes to code that exists, and no new package. The Queue Owner writes on the item the Run its Implement executor reports as soon as the executor returns, whatever it returned. A Delivery Retry accepts a QA-Only Failed Run as one more source of Task Carry-Forward. `roundfix deliver retry` holds back its `Retried` line until a live Queue Owner is proven, and otherwise parks the item again as it was and exits `2`. `roundfix deliver status` gains an `Owner:` line for a queue nobody holds and one `Run:` line per item at `running`. No table, column, stage, command, flag or exit code is added. The measurements behind the design are in [`references/2026-11-23-characterization.md`](references/2026-11-23-characterization.md).

The primary trade-off is when the item's Run is recorded: at the Run's end, not at its start. A child process creates the Implement Run, and the owner reads that child's output only when it exits, so it could learn the Run earlier only by reading the Run Database on a timer beside its child. The design gives up the item's record during an Active Run. Status reads the Active Run from the Run Database instead, and a retry refuses an item that is not parked, so no reader depends on the record while a Run is active.

The second trade-off is that the retry undoes instead of reordering. An owner started before the stage is stored finds nothing to advance and releases the queue, so the stage is stored first, the owner is proven second, and a failed proof reverts the item in one guarded transaction.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The design reads Run identifiers, process identifiers and start-time identities that already exist, and the decision record takes the number 0450. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added. Every changed command reads the local Run Database, local Git objects and the local process table. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands. ADR-0279 and ADR-0395 are the decisions this design amends in part, by the maintainer's decision the PRD quotes, and ADR-0450 is the record a Task writes for it. ADR-0028, ADR-0158, ADR-0170, ADR-0256, ADR-0284, ADR-0304, ADR-0349, ADR-0365, ADR-0373 and ADR-0398 are the ones it obeys or uses. ADR-0030, ADR-0034, ADR-0056, ADR-0098, ADR-0148, ADR-0164, ADR-0311, ADR-0315, ADR-0346, ADR-0384, ADR-0385, ADR-0406, ADR-0412 and ADR-0435 are listed only to record why they do not apply. This stage adds no ADR to the PRD's list, and `## ADR Alternatives` relates each one. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's two decisions of 2026-11-23, quoted with their provenance in the authorization record, which grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. The design changes no Makefile, workflow, Go module, lint or formatter file. The governed set was measured again with `speccheck.GovernedPath` over every file the Build Order names (`references/2026-11-23-characterization.md`, section 7): the three granted paths are the only governed ones. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Every component extends something that exists. The new files are one Go source file, the new test files and one ADR.

### The item's latest Run

In `internal/cli/deliver_workflow.go`, `deliveryRunResult` maps the exit of `roundfix implement` to a result, and its failure branch returns an empty one, so a Run the child created and ended Failed never reaches the engine. The branch now returns that Run with the failure. In `internal/delivery/engine.go`, `runCandidate` and the Run Slot Wait call a new `recordRun` first thing after the executor returns.

### The QA-Only Failed Run as a source

`commandDeliveryWorkflow.CarryForward` walks every terminal Implement Run of the item's Spec on the Item Branch and skips an outcome outside `runworktree.CarryForwardAcceptedOutcomes`. It now classifies a skipped Run whose outcome is Failed, and a QA-Only Failed Run goes through the same inspection, proofs and application as an accepted outcome. `inspectCarryForwards` learns of the classification from one new field of `reconcileRunSelection`. The accepted-outcome list does not change, so the Reconcile Command and the Implement Command's Preflight, which never set the field, behave as before.

### The retry's proof

`runDeliverRetry` in `internal/cli/deliver.go` keeps its hand-off to a recorded owner proven alive. Otherwise it starts the owner with the owner's report buffered, reads the queue again and proves the owner it records with the process controller it already uses. A failed proof calls one new store operation and prints one failure report.

### The status lines

`deliveryOwnerLine` gains the case of no recorded owner. A new file, `internal/cli/deliver_run_line.go`, derives the `Run:` lines from two reads the store already offers, `ListRuns` and `RunEventsOfKinds`, on the read-only connection status already opens.

### The decision record, the glossary, the guide and the skill

One Task writes ADR-0450 with its notices and its export line. `CONTEXT.md` changes in the Task that changes the behavior an entry describes. The last implementation Task rewrites the guide and the skill reference and regenerates the mirror, the version fields and the Behavior Surface Record.

## Implementation Design

### Interfaces

```go
// internal/delivery/engine.go
// recordRun writes the Run the Implement executor reported on the item, in
// memory and in the Run Database, under a context the owner's stop does not cancel.
func (engine *Engine) recordRun(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem, runID string) error

// internal/cli/deliver_workflow.go
// qaOnlyFailedRun names the QA gate's Task of a QA-Only Failed Run.
func (workflow *commandDeliveryWorkflow) qaOnlyFailedRun(ctx context.Context, repository string, specsRoot roundconfig.SpecsRoot, run store.Run) (qaTask string, ok bool, err error)

// internal/cli/reconcile.go
type reconcileRunSelection struct {
	// existing fields stay
	qaOnlyFailed map[string]bool // Run ids a Delivery Retry classified
}
```

```go
// internal/store/delivery.go
// RevertDeliveryQueueItemRetry parks a retried item again as it was before the retry.
func (store *Store) RevertDeliveryQueueItemRetry(ctx context.Context, gitRoot string, before DeliveryQueueItem, retried DeliveryStage) (bool, error)

// internal/cli/deliver.go
type retryOwnerProof struct {
	Proven   bool   // a live Queue Owner holds the queue, or no item can advance
	Reason   string // why not, as the failure report prints it
	Reparked bool   // the item is parked again as before the retry
}

// internal/cli/deliver_run_line.go
func deliveryRunLines(ctx context.Context, runStore *store.Store, gitRoot string, queue store.DeliveryQueue) ([]string, error)
```

### Data Models

No schema change. `delivery_queue_items.run_id` keeps its column and its meaning and is written at one more moment. The revert writes `stage`, `blocker`, `retry_count` and `candidate_commits` of one row. Status reads `runs` and `run_events` and writes nothing.

### API Contracts

1. API Contract: the item's recorded Run. The executor returns the identifier of a Run its child created for every exit of the child, a Run that ended Failed included. When that identifier differs from the item's, the owner writes it to the item before it reads the error or the outcome, with a context its own stop does not cancel. A store error of that write is returned joined with the executor's error. The park that follows keeps its blocker and its Park Class.
2. API Contract: the QA-Only Failed Run. A terminal Implement Run is one when its outcome is Failed, its Run Branch `roundfix/run-<run-id>` exists, the Spec's Task Graph loads from the commit at that branch's tip, the graph has a QA gate's Task, that Task is not `completed` there, and every other Task is. The read goes through local Git objects with the loader the Implement Command uses for a committed graph. A Run whose branch is gone or whose graph does not load there is not a source and is skipped as today. A Delivery Retry treats a QA-Only Failed Run as it treats an accepted outcome, with one difference: the evidence of the QA gate's Task is removed before inspection, so that Task is never a candidate. A refusal is whole and is reported as today's refusal.
3. API Contract: the end of `roundfix deliver retry`, after the engine stored the item's re-entry stage.
   - A recorded owner: unchanged. Alive and proven, the item is handed to it and the command exits `0`. Proven gone, its record is reclaimed and the command goes on. Alive with an identity that cannot be checked, the command fails with today's text.
   - With no recorded owner, or after reclaiming one, the command starts an owner, buffering what the start writes to standard output, and reads the queue again. When no item can advance, the owner it started already left every item merged or parked, and the command succeeds. When an item can advance, the queue must record an owner whose PID is alive and whose start-time identity the process controller proves.
   - Success prints, in today's order, the `Merged outside the queue` line if any, the `Carried forward from Run` lines, the `Retried` line, the `Review:` line if any, then the buffered owner report, and exits `0`.
   - Failure prints nothing on standard output. When the owner did not start, when the record is of an owner proven gone, or when no owner is recorded, the command releases a record proven gone and calls API Contract 4. When the identity of a live PID cannot be checked, it changes nothing more, because that owner counts as present. It then prints the failure report on standard error and exits `2`.
   - The report's `Reason:` is one of: `the Queue Owner did not start`; `Queue Owner PID <pid> is not running`; `Queue Owner PID <pid> has a different process identity`; `no Queue Owner holds the Delivery Queue after the start`; `Queue Owner PID <pid> could not be checked: <error>`. Its `Item:` line prints the stage and the blocker the Run Database holds after the command, with `-` for no blocker. Its `Carried forward:` block has one line per Run this command carried from and is absent when it carried none. Its `Next:` line names `roundfix deliver retry <slug>` when the item is parked again and `roundfix deliver status` otherwise.
4. API Contract: `RevertDeliveryQueueItemRetry`. In one write transaction, and only when the queue records no owner and the stored item is at the stage the retry stored with an empty blocker, it sets the stage to `parked` and restores the blocker, the retry count and the candidate commits of the item as read before the retry. It leaves the recorded Run, the Item Branch, the item worktree and the queue's Run links. It returns `false` and changes nothing when either condition fails.
5. API Contract: the `Owner:` line. When at least one item is neither `merged` nor `parked` and the queue records no owner, status prints `Owner: Delivery Queue has no recorded owner and <n> item(s) can still advance; run roundfix deliver resume`, where `<n>` counts those items, in the place of today's `Owner:` line. The line of a recorded owner proven gone, the diagnostic of an owner that cannot be checked, the silence for a proven live owner and the silence when no item can advance keep their bytes.
6. API Contract: the `Run:` line. For each item at `running`, in queue order, after the `Turn:` lines and before the `Park:` lines, status prints one line. The item's Active Run is the newest non-terminal Implement Run of the repository whose Spec is the item's and whose local branch is the Item Branch. With none: `Run: <slug> no Active Run`. With one, its Tasks in progress are the Tasks whose latest `daemon.task` event of phase `started` has no later event of phase `settled` or `skipped`, in the order they started: `Run: <slug> <run-id>; Tasks in progress: <task>, <task>`, or `Run: <slug> <run-id>; no Task in progress`. A read error of the Run Database fails status as its other reads do.
7. API Contract: the decision record. `docs/adr/0450-a-delivery-retry-ends-with-a-live-queue-owner-and-status-reports-a-queue-nobody-advances.md` is accepted, carries a true Knowledge Mark and has the title `A Delivery Retry ends with a live Queue Owner, and status reports a queue nobody advances`. It states API Contracts 1 to 6 as rules, the two options this design gives up, a record at the Run's start and a restart of the owner by Roundfix, and that it amends ADR-0279 and ADR-0395 in part. Each of those two gains `amended_by: [ADR-0450]` and a paragraph that begins `**Amended by [ADR-0450](<file>) (<date>).**`, with the date of ADR-0450's `created_at`. In ADR-0279 the paragraph follows the sentence that a queue with no recorded owner adds no line; in ADR-0395 it follows the paragraph that an owner that stops with an item left keeps its record. `.secondbrain-export` gains the ADR's path once, in order, inside the Standard Export List.
8. API Contract: the glossary. `CONTEXT.md` gains QA-Only Failed Run with the definition of API Contract 2. Task Carry-Forward says that a Delivery Retry also accepts a QA-Only Failed Run and that the Reconcile Command and the Preflight do not. Delivery Retry says that it ends with a Queue Owner proven alive, or fails and parks the item again as it was. Queue Owner says that status also reports a queue with no recorded owner while an item can advance, and the Run of each item at `running`. Each revised entry cites ADR-0450.
9. API Contract: the guide and the skill. `docs/user-guide/commands/deliver.md` and `.agents/skills/roundfix/references/deliver.md` describe the line of API Contract 5, the three forms of API Contract 6 with their place, the failure report of API Contract 3 with its exit code, and the QA-Only Failed Run. Neither keeps the words "no recorded owner" in a sentence about states that add no `Owner:` line, and the guide's sentence that a retry with no recorded owner "starts one directly" says what the retry then proves.

### Surface Transcripts

1. Surface Transcript: status of a queue with one item at `running`, no recorded owner and no Active Run, the state quaestor measured.

   ```transcript
   $ roundfix deliver status
   stdout:
   <slug>	running	-	<item worktree>
   Owner: Delivery Queue has no recorded owner and 1 item(s) can still advance; run roundfix deliver resume
   Run: <slug> no Active Run
   Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
   Usage: <usage>
   stderr:
   exit: 0
   ```

2. Surface Transcript: status under a proven live owner whose Active Run has two Tasks in progress.

   ```transcript
   $ roundfix deliver status
   stdout:
   <slug>	running	-	<item worktree>
   Turn: <slug> advancing at running
   Run: <slug> <run-id>; Tasks in progress: task_03, task_04
   Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
   Usage: <usage>
   stderr:
   exit: 0
   ```

3. Surface Transcript: the same queue between two Tasks.

   ```transcript
   $ roundfix deliver status
   stdout:
   <slug>	running	-	<item worktree>
   Turn: <slug> advancing at running
   Run: <slug> <run-id>; no Task in progress
   Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
   Usage: <usage>
   stderr:
   exit: 0
   ```

4. Surface Transcript: a retry after a QA-Only Failed Run with five completed Tasks, whose new owner is proven alive.

   ```transcript
   $ roundfix deliver retry <slug>
   stdout:
   Carried forward from Run <run-id>: task_01, task_02, task_03, task_04, task_05
   Retried <slug>: <blocker> -> running
   Delivery Owner: <owner id>
   Console Log: <path>
   Status: roundfix deliver status
   Stop: roundfix deliver stop
   stderr:
   exit: 0
   ```

5. Surface Transcript: a retry whose owner does not start, with nothing carried. The first lines of standard error are the start's own diagnostic.

   ```transcript
   $ roundfix deliver retry <slug>
   stdout:
   stderr:
   ...
   Retry failed

   Reason:
     the Queue Owner did not start

   Item:
     stage: parked; blocker: <blocker>

   Next:
     roundfix deliver retry <slug>
   exit: 2
   ```

6. Surface Transcript: a retry that carried five Tasks and whose owner ended before the proof.

   ```transcript
   $ roundfix deliver retry <slug>
   stdout:
   stderr:
   Retry failed

   Reason:
     Queue Owner PID <pid> is not running

   Item:
     stage: parked; blocker: <blocker>

   Carried forward:
     Run <run-id>: task_01, task_02, task_03, task_04, task_05

   Next:
     roundfix deliver retry <slug>
   exit: 2
   ```

### Invariants

```text
1. Every line status prints on the Delivery Base keeps its bytes and its order relative to the others.
2. The item row keeps four tab-separated fields.
3. The accepted-outcome list of Task Carry-Forward keeps three outcomes.
4. The QA gate's Task is never carried.
5. A retry that fails spends no retry: the count after it equals the count before it.
6. Status opens the Run Database for reading only and writes no file.
7. A test finds ADR-0450 by its title, never by its number.
8. The help of every deliver command keeps its bytes.
```

## Coverage Map

- Goal 1 → The retry's proof (API Contracts 3 and 4)
- Goal 2 → The QA-Only Failed Run as a source (API Contract 2)
- Goal 3 → The item's latest Run (API Contract 1)
- Goal 4 → The status lines (API Contract 5)
- Goal 5 → The status lines (API Contract 6)
- Story 1 → `qaOnlyFailedRun`, `CarryForward` (API Contract 2, Surface Transcript 4)
- Story 2 → `runDeliverRetry`, `RevertDeliveryQueueItemRetry` (API Contracts 3 and 4, Surface Transcripts 5 and 6)
- Story 3 → `deliveryOwnerLine` (API Contract 5, Surface Transcript 1)
- Story 4 → `deliveryRunLines` (API Contract 6, Surface Transcripts 2 and 3)
- Story 5 → `recordRun`, `deliveryRunResult` (API Contract 1)
- Core Feature 1 → `recordRun`, `deliveryRunResult` (API Contract 1)
- Core Feature 2 → `qaOnlyFailedRun`, `CarryForward`, `inspectCarryForwards` (API Contract 2, invariant 4)
- Core Feature 3 → `inspectCarryForwards` and the unchanged accepted-outcome list (invariant 3)
- Core Feature 4 → `runDeliverRetry`, `retryOwnerProof`, `RevertDeliveryQueueItemRetry` (API Contracts 3 and 4, invariant 5)
- Core Feature 5 → `deliveryOwnerLine` (API Contract 5, invariants 1 and 6)
- Core Feature 6 → `deliveryRunLines` (API Contract 6, invariants 1 and 2)
- Core Feature 7 → `references/2026-11-23-characterization.md`, the tests of Testing Approach, Build Order 1 to 3
- Core Feature 8 → the decision record and the glossary (API Contracts 7 and 8, invariant 7)
- Core Feature 9 → the guide and the skill (API Contract 9, invariant 8)
- Success Metric 1 → `TestTheItemRecordsTheRunOfAFailedImplementExecutor`, `TestAFailedImplementChildNamesTheRunItCreated`
- Success Metric 2 → `TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun`, with its controls
- Success Metric 3 → `TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof`, `TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun`
- Success Metric 4 → `TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner`, with the existing owner status tests
- Success Metric 5 → `TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress`, `TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress`
- Success Metric 6 → `TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules`, `TestTheDeliverGuideAndReferenceDescribeTheRecoveryLines`, with the existing documentation, skills and glossary contracts
- Success Metric 7 → `TestTheBuiltBinaryReportsTheQueueQuaestorMeasured`, run by the QA gate

Deviations from the PRD:

- Core Feature 1 asks for the record as soon as the owner can read the Run's identifier. The owner records it when the Implement executor returns, so during an Active Run the item still names the Run that ended before it, and status names the Active Run. The reason is in the Executive Summary.
- Core Feature 4 asks for exit `0` only with a Queue Owner proven alive at the command's end. When the owner the retry started has already left every item merged or parked, the command exits `0`: nothing can advance, and parking the item again would overwrite the blocker that owner just recorded.
- Core Feature 4 asks for the item parked again when the owner "cannot be identified". The command exits `2` with the reason and leaves the item and the owner record as they are, because an accepted rule counts such an owner as present (ADR-0373: "are not proof, and the owner counts as present").
- Core Feature 6 names three states. While a Run is in its QA gate the line reads `no Task in progress`, because the Run Database holds no start event for the gate's Task.
- Core Feature 7 asks for named tests that reproduce the four defects on the unchanged tree. This stage's probes reproduce them and are recorded in `references/2026-11-23-characterization.md`. Each Task's named tests then fail on the Delivery Base with the defect, which the Verification probe proves before the Run. A test that pins a defect would have to be inverted by the Task that fixes it.

## Absence Evidence

1. A write of the Run on the item for a failed executor — command: `rg -n 'item\.RunID = ' internal/delivery/engine.go`; result: two matches, the retry's write at line 836 and the write at line 1351, which follows the return of the executor's error.
2. The QA-Only Failed Run in the carry-forward code, the glossary, the ADRs, the guide or the skill — command: `rg -n -i 'qa-only|qaonly|qa_only' internal/cli/deliver_workflow.go internal/cli/carryforward.go internal/cli/reconcile.go CONTEXT.md docs/adr docs/user-guide .agents/skills/roundfix`; result: no match. The same words over `internal` match fourteen lines in five other files, each about a Run or a branch that holds only the QA step.
3. A store operation that reverts a retry — command: `rg -n -i 'revert|repark|restore' internal/store/delivery.go`; result: no match.
4. A proof of the owner after the retry starts it — command: `rg -n 'startDeliveryOwner\(' internal/cli/deliver.go`; result: three matches, at lines 233, 386 and 597, each returning the start's exit code directly.
5. A status line for a queue with no recorded owner — command: `rg -n 'no recorded owner' internal cmd`; result: no match.
6. A status line that names a Run or its Tasks — command: `rg -n 'Run: ' internal/cli/deliver.go internal/cli/deliver_turn.go internal/delivery/turn.go`; result: no match, and `rg -n 'Tasks in progress|no Task in progress' internal cmd docs/user-guide .agents/skills/roundfix` finds none.
7. The decision record — command: `ls docs/adr | rg '^045'`; result: no match, and `rg -l 'nobody advances' docs/adr CONTEXT.md` finds no file.
8. `internal/cli/deliver_run_line.go` and a reader of Task events for status — command: `rg -n 'RunEventsOfKinds\(' internal --glob '!*_test.go'`; result: six matches in four files, the store's definition, the Run cause report, the corrections report and the unresolved-Run reader of `deliver_workflow.go`; none in `deliver.go`, and `ls internal/cli/deliver_run_line.go` finds no file.

## ADR Alternatives

- ADR-0450: the record a Task of this Spec writes, absent until then; it rejects nothing yet. API Contract 7 fixes what it states and the two options it records as given up.
- ADR-0279: amends it in part and takes one option it rejected, by the maintainer's decision of 2026-11-23 that the PRD's Core Features 5 and 6 name. It kept a queue without an owner silent (ADR-0279: "A queue with no recorded owner adds no line.") and froze the rest of the output (ADR-0279: "keeps its exit code, output order and every other line"). The design reports that queue and adds lines between the frozen ones. Its other rejections stand: a PID alone is never proof, and status never reclaims a record (ADR-0279: "Reclaiming the owner record from status would make a read command write.").
- ADR-0395: amends it in part and takes the option it rejected for ADR-0279's reason, by the same decision (ADR-0395: "Keeping the record of an owner that died names exactly the defect."). Its own rule stands, and the retry's proof relies on it (ADR-0395: "An owner that stops with an item left keeps its record."). The panic recovery it rejected is not added.
- ADR-0170: obeys it. The item's recorded Run stays provenance (ADR-0170: "The Run recorded on the item is provenance, not the selector."), and a set is carried whole or refused whole. Tasks carried before a failure stay and are named (ADR-0170: "A later Run's refusal leaves earlier Runs' proved Tasks on the item branch and names them.").
- ADR-0158: extends what it accepted, for a Delivery Retry only, "under the proof requirements it already applies". The relabelling it rejected is not taken (ADR-0158: "Settling these Runs as `Stopped` was rejected."): a QA-Only Failed Run keeps the outcome Failed.
- ADR-0284: obeys it. The Binding Text proof decides which planning amendments a carried Task survives, and the option it kept rejected stays out (ADR-0284: "Carrying the unaffected part of a refused set stays rejected").
- ADR-0304: obeys it (ADR-0304: "No stage, column or schema is added."). An item in a Run Slot Wait prints `no Active Run` beside its unchanged blocker text, and no exit code is added, which it rejected for its own refusal (ADR-0304: "A distinct exit code for the ceiling refusal was rejected").
- ADR-0365: obeys it. It rejected deriving the turn from the Run Database (ADR-0365: "nothing there says which item the owner is advancing") and a column for it. The `Run:` line derives the Run and its Tasks, never the turn, and the Turn Record is read as before (ADR-0365: "Status reads it only for the owner it proves live").
- ADR-0373: obeys it. The retry's proof reads an absent PID and a mismatch as gone and an unreadable identity as present, and does not take the option it rejected (ADR-0373: "Treating an unreadable identity as gone too").
- ADR-0398: obeys it (ADR-0398: "A parked item always has a way back through a Delivery Retry"). A retry that fails leaves the item parked and retryable with its count restored, since no retry took place. It rejected pushing from the retry (ADR-0398: "the retry would reach the network"); the proof reads only the local Run Database and process table.
- ADR-0028: obeys it (ADR-0028: "a Run's lifetime must not depend on the caller's"). The owner stays detached; the retry waits for the handshake it already waits for and adds one read. It records no rejected option.
- ADR-0256: obeys it. The guide's `deliver` page is a Behavior Surface whose fingerprint moves, and the covering skill reference changes in the same range (ADR-0256: "is a Lagging Surface unless a covering skill file changed in the same range"). The record is regenerated by its command, never judged by a model, which it rejected.
- ADR-0349: uses it. The authorization record carries a `patch` Release Ceiling and the Spec is a `fix`. No Task performs the release.
- ADR-0030: not applicable — no Agent log is read or written.
- ADR-0034: not applicable — the Worktree Bootstrap is not changed.
- ADR-0056: not applicable — no capacity is read or changed.
- ADR-0098: not applicable — the design reads Run Events and appends none.
- ADR-0148: not applicable — the Verification probe is used as it is.
- ADR-0164: not applicable — the Run Budget and its renewal are not changed.
- ADR-0311: not applicable — the Check Wait is not changed.
- ADR-0315: not applicable — no rule of the Spec Consistency Check changes.
- ADR-0346: not applicable — the Task prompt and the skills an Agent loads are not changed.
- ADR-0384: not applicable — the Machine Verification Capacity and its status lines keep their bytes.
- ADR-0385: not applicable — no Verification environment changes.
- ADR-0406: not applicable — the authoring probe's checkout is used as it is.
- ADR-0412: not applicable — the release plan's reading of a commit is not changed.
- ADR-0435: not applicable — a Pull Request behind its base is handled as before.

## Integration Points

- The Run Database: read by status through the read-only connection, written by the owner and the retry through the store's write transactions. No migration.
- Local Git: the QA-Only Failed Run is classified from the committed Task Graph at the Run Branch's tip. No remote is read.
- The process table: the retry's proof uses the injected `OwnerProcessController` and `store.ProcessAlive`, as the hand-off does.
- The detached start: unchanged. The retry passes it a buffer instead of its own standard output.
- Routines that read status, the Praetor's sweep among them: they receive two new line prefixes between lines they know. The operator tells their owners after the release; no Task sends a message.
- The Secondbrain: reads `.secondbrain-export` after the merge and receives ADR-0450.

## Testing Approach

No new seam. The tests use what the delivery tests already inject: the fake runner of `internal/delivery`, and in `internal/cli` the Delivery Engine, the owner start and the process controller of `commandDependencies`, over a real temporary Run Database and temporary repositories.

- `TestTheItemRecordsTheRunOfAFailedImplementExecutor` (`internal/delivery`): a runner returns a Run with an error, and the item parks with that Run and the unchanged blocker. A second case cancels the context inside the runner and finds the Run on the item though nothing parked it.
- `TestAFailedImplementChildNamesTheRunItCreated` (`internal/cli`): `deliveryRunResult` for exit `1` with a new Run in state Failed returns the Run with the unchanged failure text, and none when no Run was created.
- `TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun` (`internal/cli`): the fixture `createDeliveryRetryRun` with the outcome Failed and a QA gate's Task in the graph. The production `CarryForward` carries the settled Task, stamps it and leaves the gate's Task pending. Controls in the same file: a Failed Run with another Task not completed carries nothing; a Failed Run whose Run Branch is gone carries nothing; the same Run Unresolved carries as before; `inspectSpecCarryForwards`, which the Preflight uses, reports no candidate for the QA-Only Failed Run.
- `TestRevertDeliveryQueueItemRetryRestoresTheParkedItem` (`internal/store`): the revert after a retry, and its two refusals.
- `TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun`, `TestDeliverRetryReproducesSurfaceTranscriptOwnerDidNotStart` and `TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof` (`internal/cli`): Surface Transcripts 4, 5 and 6, each whole. The engine fake stores the stage as the production engine does; the start fake claims the queue with the test's own PID, returns `2`, or claims it with a PID that does not exist. The last two also assert the stored item: `parked`, the earlier blocker, the earlier retry count, no recorded owner. Further cases, not transcripts: an owner with a different identity; an identity that cannot be checked, which leaves the item and the record; an owner that already parked the item, which exits `0`.
- `TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner`, `TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress` and `TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress` (`internal/cli`): Surface Transcripts 1, 2 and 3, each whole, over Runs and `daemon.task` events appended to the temporary Run Database.
- `TestTheBuiltBinaryReportsTheQueueQuaestorMeasured` (`internal/cli`): builds the binary as `buildRoundfixBinaryForMacro` does, rebuilds in a temporary home the state of `references/sources.md`, section 2, and asserts Surface Transcript 1 on the binary's output. A second case makes `<home>/.roundfix/artifacts` a regular file and asserts Surface Transcript 5 and the stored item on the binary's `deliver retry`.
- `TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules` (`internal/docscontract`, with the package's build tag and marker): API Contract 7 by invariant 7, as `TestTheIndexusDecisionAmendsTheConsultationRule` does for its ADR.
- `TestTheDeliverGuideAndReferenceDescribeTheRecoveryLines` (`internal/docscontract`): the phrases of API Contract 9 in the guide, the canonical reference and its mirror, and the absence of the replaced sentence.
- Existing tests that move, listed with what each gains in `references/2026-11-23-characterization.md`, sections 5 and 6. Eight status tests gain the `Owner:` line, because their queue has an item that can advance and no recorded owner; `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner` among them is renamed for what it now asserts. `TestDeliverStatusReproducesSurfaceTranscriptQueueTurns` gains one `Run:` line and `TestDeliverStatusShowsARunSlotWaitTranscript` gains both. The owner start of `newOutsideMergeFixture` claims the queue as a live owner does.
- Existing tests that keep passing unchanged: `TestCarryForwardAcceptedOutcomes`, the item recovery tests, the other owner status tests and every other retry test.

The QA gate runs these by name, with the two cases of `TestTheBuiltBinaryReportsTheQueueQuaestorMeasured` as its built-product rows, and no network.

## Build Order

1. The item's latest Run and the QA-Only Failed Run: API Contracts 1 and 2 and the glossary entries QA-Only Failed Run and Task Carry-Forward, with the first three tests of Testing Approach. Files: `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go`, `internal/cli/carryforward.go`, `internal/cli/reconcile.go`, `CONTEXT.md` and two new test files.
2. The retry's proof: API Contracts 3 and 4 and the glossary entry Delivery Retry, with the store test and the three retry transcript tests (depends on: 1, because both change `CONTEXT.md` and Surface Transcript 4 carries from a QA-Only Failed Run). Files: `internal/cli/deliver.go`, `internal/store/delivery.go`, `internal/cli/deliver_merged_outside_test.go`, `CONTEXT.md` and two new test files.
3. The status lines and the decision record: API Contracts 5, 6 and 7 and the glossary entry Queue Owner, with the three status transcript tests, the built-binary test and the decision record test (depends on: 2, because both change `internal/cli/deliver.go` and `CONTEXT.md` and the built-binary test asserts the retry's failure). Files: `internal/cli/deliver.go`, the new `internal/cli/deliver_run_line.go`, the nine moved test files, the three ADR files, `.secondbrain-export`, `CONTEXT.md` and three new test files. The two `amended_by` entries and the export line are written by hand or by the documentation part of `roundfix baseline update`, which writes those three files and no other (`references/2026-11-23-characterization.md`, section 8).
4. The guide and the skill: API Contract 9, with its test (depends on: 3, whose lines it describes). Edit the guide and the canonical reference, then run `make skills-sync`, the skill record command and the Behavior Surface record command of the authorization record. Files: `docs/user-guide/commands/deliver.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/references/deliver.md`, both `SKILL.md` files, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json` and one new test file.
5. The QA gate (depends on: 1, 2, 3 and 4).

Steps 1 to 3 change no help text and no page of the user guide, so no Behavior Surface moves before step 4 (invariant 8). `docs/references/skill-coverage.json` does not change: the map names no source file this design adds.

Files shared with other Specs: the Roundfix skill's version fields, the owned-skill version record and the Behavior Surface Record. Step 4 regenerates them on the tree it runs on and writes the next version, never a number. No Verification of this Spec depends on another Spec, and the Task Graph needs no `requires`.

## Risks & Considerations

- A refused set now stops a retry that the unchanged tree would have let through. Before, a Failed Run was skipped and the Spec ran again; now a QA-Only Failed Run whose Task moved an input refuses the retry, as an Unresolved Run does, and the reason names the Task and what moved. An older QA-Only Failed Run of the same Item Branch is considered too.
- The `Run:` line reports what the Run Database records. An Active Run whose process died stays non-terminal until it is reclaimed, and its line keeps naming it; the `Owner:` line beside it is what says nobody is advancing the queue.
- The queue's `Usage:` line and its token ceiling now count a Failed Run, because the Run is linked when the item records it. The line's format does not change.
- The retry returns later: it waits for the start handshake, as today, and adds two reads. The detached start's timeouts bound the wait.
- An owner killed without a chance to write leaves the earlier Run on the item until its next Run ends. No reader depends on it.
- The revert is refused when another command claimed the queue between the start and the revert. The report then prints the item as it is and names `roundfix deliver status`.
- The proof is not absolute: a start-time identity is compared, never a handle on the process, and the process can end the instant after the proof. Status reports that case.
- ADR-0450's number can be taken while the Pull Request waits. The operator then renames the file, its links in the two notices and its export line; invariant 7 keeps the test independent of the number.

## Glossary

None.

## Decisions

- The PRD's first open question is settled by measurement: a QA gate's Task settles `failed`, and the Run ends Failed when the QA Report commit or its format step fails after that settlement (`references/2026-11-23-characterization.md`, section 1). That quaestor's Run took this path is inferred.
- Its second is settled by measurement: the retry stores the stage and prints `Retried` before it starts the owner, and a start that does not succeed leaves the state quaestor measured, through the built binary as well (section 3). The Turn Record quaestor found is the last write of the owner that ran the Failed Run (section 2). That the start is what failed on 2026-11-22 is inferred.
- Its fourth is settled by a source: the Run Database holds a start and a settlement event for each Task of a Run (section 4). Its fifth is settled by the Surface Transcripts.
- Delegated Decision of 2026-11-23: When does the Queue Owner record on the item the Run it started?
- Delegated Decision of 2026-11-23: What does the Run line of status show while a Run is in its QA gate?
- Delegated Decision of 2026-11-23: What does a Delivery Retry do when the owner it started has already left every item merged or parked?
- The retry leaves the item and the owner record alone when a live owner's identity cannot be checked. A source settles it, ADR-0373, which rejected the opposite (ADR-0373: "Treating an unreadable identity as gone too").
- The QA-Only Failed Run is classified from the committed graph at the Run Branch's tip, not from the Run Worktree's files, so an uncommitted status never counts as completed.
- A retry's failure has one report for every reason. The hand-off's own failure keeps today's text.
- The export line and the two `amended_by` entries may be written by hand: the Secondbrain guide allows it for an ADR's path (`docs/agents/secondbrain.md`: "the only line added by hand is the path of an ADR, in order") and the docs-layout guide for the header (`docs/agents/docs-layout.md`: "In the change that adds the later ADR, its header gains").
- The PRD's `## Glossary` declares the term this Spec adds and the three it changes; this stage names no other.
- The Secondbrain pages and the published sources consulted on 2026-11-23, with what each changed, are recorded in `references/2026-11-23-characterization.md`, section 9.
