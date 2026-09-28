---
spec: 0173-a-delivery-queue-that-recovers
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# A delivery queue that recovers

## Executive Summary

Stage carry-forward candidates in the order the Run integrated them, commit
Project Config only on a named grant and report every exclusion, retry a
timed-out profile proof once, record the Run that parks an item, and add
`roundfix deliver retry <slug>`: it carries the item's settled Tasks back to
the item branch, derives the re-entry stage from recorded evidence, moves the
item out of `parked` in one guarded transaction and hands it to a live owner or
starts one. No Run Database schema change: the `run_id` and
`candidate_commits` columns already exist.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git, the Run Database
  and ACP Runtime subprocesses only; no credential and no new network call.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0002, ADR-0014, ADR-0026, ADR-0044,
  ADR-0050, ADR-0052, ADR-0053, ADR-0057, ADR-0080, ADR-0091, ADR-0093,
  ADR-0096, ADR-0104, ADR-0107, ADR-0114, ADR-0117, ADR-0130, ADR-0138,
  ADR-0139, ADR-0140, ADR-0153, ADR-0155, ADR-0156, ADR-0157, ADR-0158 and
  ADR-0160 hold; ADR-0020, ADR-0038, ADR-0056, ADR-0069, ADR-0097, ADR-0127
  and ADR-0159 do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Carry-forward in integration order

`carryForwardTaskCommits` in `internal/cli/carryforward.go` already walks
`git rev-list --reverse <run.HeadSHA>..HEAD` in the Run Worktree, which is the
order ADR-0026 integrated the settled Tasks. It keeps that order: alongside the
per-Task commit map it returns the Task IDs in the position of their first
settlement commit. `inspectCarryForwardsForRun` iterates the settled-completed
Tasks in that order instead of `graph.Tasks` order; a settled-completed Task
with no settlement commit keeps its refusal and follows the ordered Tasks in
Task Graph order. `unstagedCarryForwards` lists the Tasks after a conflicting
one in the same integration order. `applyCarryForwards` already stages the
candidates in the order it receives them, so the carried history follows the
Run's history.

Every proof is unchanged: each candidate's declared inputs are still compared
with the accumulating staged state, which is now the checkout plus the Tasks
the Run integrated before it — the state its settlement parent had when the
checkout has not moved. An input changed on the checkout still refuses the
whole set. The `carryForwards` array of `roundfix-reconcile/v1` and the Tasks
`implement` names follow the same order. The reconcile section of
`docs/user-guide/commands.md` says the candidates are proved in the order the
Run integrated them.

## Project Config in Daemon commits

`diffSnapshots` in `internal/daemon/engine.go` becomes a pure snapshot diff: it
no longer skips `.roundfixrc.yml`. Each caller decides instead, and every
exclusion becomes a `DroppedStagePath` published through
`publishDroppedStagePath`:

- `prepareTaskCommit` in `internal/daemon/task_engine.go` keeps
  `.roundfixrc.yml` in the Task's changed paths when the plan's frozen
  authorization (`plan.Authorization`) has outcome
  `spec.AuthorizationGranted` and its `Record.Paths` lists `.roundfixrc.yml`;
  the existing governed-mutation checks then apply. Otherwise it drops the path
  with reason `Project Config outside the Spec's authorization` and
  `Lost: true`, so `lostStagePathsFailureReason` settles the Task failed with
  `Task commit lost output: .roundfixrc.yml (Project Config outside the Spec's
  authorization)`, ahead of the governed-mutation refusals as every lost output
  is today.
- `commitBatch` in `internal/daemon/engine.go` and `commitQAReport` in
  `internal/daemon/task_engine.go` drop it with reason `Project Config is never
  committed by a Batch or QA Report commit` and `Lost: false`, so the commit
  proceeds without it.

`publishDroppedStagePath` gains a branch for both reasons: progress prints
`roundfix: Project Config .roundfixrc.yml omitted from the commit: <reason>`
and the event summary reads `Project Config .roundfixrc.yml omitted from the
commit: <reason>.`; the payload keeps `decision: dropped`, `path` and
`reason`. A `.roundfixrc.yml` already dirty in the before snapshot stays out of
every diff, as any pre-existing change does. The Project Config entry of
`CONTEXT.md` and the Daemon commit paragraph of the Roundfix skill state the
rule with the phrase `Project Config outside the Spec's authorization`.

## Retry in the engine and store

`runCandidate` in `internal/delivery/engine.go` records
`strings.TrimSpace(result.RunID)` on the item before it parks it
`run-unresolved`.

`internal/store/delivery.go` adds two transactions:

- `RetryDeliveryQueueItem(ctx, gitRoot string, item DeliveryQueueItem,
  parkedBlocker string) (ownerPID int, ownerIdentity string, err error)`
  refuses a target stage other than `running`, `reviewing`, `gating` or
  `checking`, and refuses unless the stored item is `parked` with exactly
  `parkedBlocker` at the item's stored position; a refusal names the stored
  stage and blocker and changes nothing. Otherwise, in the same write
  transaction, it persists the target stage, an empty blocker, the item's Run
  ID and candidate commits, and reads the queue's recorded owner PID and
  identity, which it returns.
- `ReleaseIdleDeliveryQueueOwner(ctx, gitRoot string, pid int, identity string)
  (bool, error)` clears the owner in one write transaction only when the caller
  is the recorded owner and no item is in a stage other than `merged` or
  `parked`; with an advanceable item it returns `false` and keeps the claim,
  and when the caller is not the recorded owner it returns an error naming the
  recorded owner.

`internal/delivery/engine.go` adds the item recovery boundary and `Retry`:

```go
type ItemState struct {
	Archived        bool
	UnfinishedTasks []string
	Head            string
}

type CarryForwardResult struct {
	RunID   string
	Carried []string
}

type ItemRecovery interface {
	InspectItem(ctx context.Context, workDir, specSlug string) (ItemState, error)
	CarryForward(ctx context.Context, workDir, specSlug, branch, runID string) (CarryForwardResult, error)
}

type RetryResult struct {
	SpecSlug      string
	Blocker       string
	Stage         store.DeliveryStage
	CarriedFrom   CarryForwardResult
	OwnerPID      int
	OwnerIdentity string
}

func (engine *Engine) Retry(ctx context.Context, gitRoot, specSlug string) (RetryResult, error)
```

`EngineDependencies` gains `Recovery ItemRecovery`; `Run` does not require it
and `Retry` does. `Retry` reads the queue, finds the item and refuses an item
that is absent or not `parked`. It calls `UseItemBranch` with the recorded
branch and worktree, which recreates a missing worktree; `ErrItemWorktreeMissing`
refuses the retry because the item branch is gone. It then inspects the item:

| Recorded evidence | Re-entry stage | Candidate commits |
| --- | --- | --- |
| Spec active, a Task unfinished after carry-forward | `running` | unchanged |
| Spec active, every Task completed after carry-forward | `reviewing` | item head appended when it differs from the last one |
| Spec archived, item head equals the candidate head, no pull request | `gating` | unchanged |
| Spec archived, item head equals the candidate head, pull request recorded | `checking` | unchanged |
| Spec archived, item head differs from the candidate head | refused | unchanged |

For an active Spec, `Retry` calls `CarryForward` with the item worktree, slug,
item branch and recorded Run ID before it inspects the item again, and records
the returned Run ID on the item. A carry-forward error refuses the retry and is
wrapped with `%w`, so its next action survives. `Retry` finally calls
`RetryDeliveryQueueItem` with the item's blocker and returns the owner it read.
Every refusal happens before that call, so a refused retry leaves the stored
item unchanged; the only work done before it is the carry-forward itself, which
is the operator's explicit request.

## Item recovery in the delivery workflow

`commandDeliveryWorkflow` in `internal/cli/deliver_workflow.go` implements
`delivery.ItemRecovery`, and `newCommandDeliveryEngine` passes it as
`Recovery`.

`InspectItem` resolves the Specs Root for the item worktree with
`roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)`. When the Spec folder
exists there it loads the Task Graph and lists, in Task Graph order, every Task
whose status is not `completed`; when it does not, it reports `Archived` when
the archive destination `archivePaths` names exists, and fails otherwise. It
reads the head with `git rev-parse HEAD`.

`CarryForward` resolves the Run: a recorded Run ID must name an Implement Run of
this Spec; with no recorded Run ID it takes the newest Implement Run of this
Spec whose recorded local branch is the item branch, and with none it carries
nothing. It carries nothing when the Run's outcome is not one
`runworktree.CarryForwardAcceptedOutcomes` returns, when
`loadReconcileTaskCoverage` finds no Task the Run settled completed, or when
every such Task is already `completed` in the item worktree — the case after an
earlier retry or an operator's own `reconcile --carry-forward`. Otherwise it
refuses when the Run Worktree is gone or the Specs Root is external, runs
`inspectCarryForwards` with the item worktree as the repository (resolved
through `filepath.EvalSymlinks`, as reconcile does), refuses with
`carryForwardRefusalReason` when any candidate refuses, and applies the set
with `applyCarryForwards`. A refusal is an error whose `NextAction` names
`roundfix reconcile <run-id> --carry-forward` in the item worktree and
`roundfix deliver retry <slug>`. The result names the Run and the carried Task
IDs in integration order.

## The retry command and the owner hand-off

`internal/cli/deliver.go` adds `retry` to `deliverUsage` (`roundfix deliver
retry <slug>` and a Commands line) and to `runDeliverCommand`, and the
`deliveryEngine` interface gains `Retry`. `runDeliverRetry` accepts exactly one
slug through the deliver flag set, opens the Run Database, builds the engine
with `newDeliveryEngine` and calls `Retry`. It prints to stdout, in order:
`Carried forward from Run <run-id>: <task>, <task>` when Tasks were carried,
then `Retried <slug>: <blocker> -> <stage>`. With a recorded owner that is
alive and proven through `ownerProcesses.ProveOwner`, it prints `Handed <slug>
to Delivery Queue owner PID <pid>.` and exits `0`. With a recorded owner that
is dead or whose identity is unproven, it releases that record with
`ReleaseDeliveryQueueOwner`, writes the same stderr notice `deliver resume`
writes, and starts an owner; with no owner it starts one through
`startDeliveryOwner`. Every refusal and failure goes through
`printDeliverFailure` and exits `2`. The top-level usage line in
`internal/cli/cli.go` becomes `roundfix deliver <start|status|resume|retry|stop>
[<slug> ...]`, and `TestRunCommandHelp` in `internal/cli/cli_test.go` expects
`roundfix deliver retry <slug>`.

`runDeliveryOwner` loops: it runs the engine, then calls
`ReleaseIdleDeliveryQueueOwner`; it exits when that releases the queue and runs
the engine again when it returns `false`. `Engine.Run` leaves every item it
read `merged` or `parked`, so a further pass happens only for an item a retry
moved after the pass read the queue. The deferred release stays for the error
paths.

The deliver section of `docs/user-guide/commands.md` and the Delivery queue
section of the Roundfix skill document the command, the re-entry table, the
carry-forward step and the hand-off. `CONTEXT.md` gains **Delivery Queue** and
**Delivery Retry** entries, and its Task Carry-Forward entry says the act is
reached through the Reconcile Command's `--carry-forward` switch or through a
Delivery Retry, and is never automatic.

## Proof timeouts

`proveProfileSelectionsWithOptions` in `internal/cli/profiles_validate.go`
proves each tuple through a helper that retries once: when
`proveProfileSelection` returns an error for which `errors.Is(err,
context.DeadlineExceeded)` holds while `ctx.Err()` is nil, the proof's own
setup deadline expired, and the helper proves the tuple again with a fresh
disposable session. A second such timeout returns a `profileProofTimeoutError`
wrapping it. `profileProofClassification` returns `temporary` for that error
before it treats a deadline as unclassified, and `profileProofNextAction`
returns advice to rerun the command when load drops, stating that the
configured profile was not shown to be wrong, without the `roundfix profiles
configure` action. A cancelled or expired command context, a runtime
construction error, an access-policy error and every other proof error are
returned after one attempt, as today. `profiles validate`, the Doctor
Command's profile check and every operational Preflight share the helper. The
profiles section of `docs/user-guide/commands.md` and the profile Preflight
paragraph of the Roundfix skill say that a proof whose setup times out is
retried once and that a second timeout is classified `temporary`.

## API Contracts

1. `roundfix deliver retry <slug>` takes exactly one Spec slug; it exits `0`
   after printing `Retried <slug>: <blocker> -> <stage>` and either `Handed
   <slug> to Delivery Queue owner PID <pid>.` or the detached owner report, and
   exits `2` for every refusal with the item unchanged.
2. The re-entry stage follows the table in "Retry in the engine and store".
3. `Store.RetryDeliveryQueueItem` moves only a `parked` item with the expected
   blocker to `running`, `reviewing`, `gating` or `checking` and returns the
   recorded owner; `Store.ReleaseIdleDeliveryQueueOwner` releases only an idle
   queue owned by the caller.
4. An item parked `run-unresolved` records the Run ID the Implement executor
   returned.
5. Task Carry-Forward stages, proves and reports candidates in the order the
   Run integrated their settlement commits.
6. A Task commit carries `.roundfixrc.yml` only when the frozen authorization
   bounds it and otherwise fails with reason `Task commit lost output:
   .roundfixrc.yml (Project Config outside the Spec's authorization)`; Batch
   and QA Report commits drop it with a published reason.
7. A `roundfix/profiles-validate/v1` proof may carry classification
   `temporary`, after exactly two attempts that each timed out.

## Coverage Map

- Goal 1 → Retry in the engine and store; The retry command and the owner
  hand-off; API Contracts 1-3.
- Goal 2 → Item recovery in the delivery workflow; API Contract 4.
- Goal 3 → Carry-forward in integration order; API Contract 5.
- Goal 4 → Project Config in Daemon commits; API Contract 6.
- Goal 5 → Proof timeouts; API Contract 7.
- Core Feature 1 → Retry in the engine and store; The retry command and the
  owner hand-off.
- Core Feature 2 → Retry in the engine and store; Item recovery in the delivery
  workflow.
- Core Feature 3 → The retry command and the owner hand-off.
- Core Feature 4 → Carry-forward in integration order.
- Core Feature 5 → Project Config in Daemon commits.
- Core Feature 6 → Proof timeouts.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 3, Testing Approach 4.
- Success Metric 3 → Testing Approach 5.
- Success Metric 4 → Testing Approach 2.
- Success Metric 5 → Testing Approach 6.
- API Contract 1 → The retry command and the owner hand-off.
- API Contracts 2-4 → Retry in the engine and store.
- API Contract 5 → Carry-forward in integration order.
- API Contract 6 → Project Config in Daemon commits.
- API Contract 7 → Proof timeouts.

## Integration Points

- **Spec 0168.** Gave each Delivery Queue item its own branch and worktree; the
  retry reuses `UseItemBranch` to recreate a missing worktree.
- **Task Carry-Forward.** The retry's carry step reuses `inspectCarryForwards`,
  `carryForwardRefusalReason` and `applyCarryForwards` unchanged apart from the
  order, so `reconcile --carry-forward` and the retry prove the same thing.
- **Implement Preflight.** After a retry carries a Run forward, `implement`'s
  own carry-forward inspection finds that Run's Task files moved and proceeds,
  as the Task Carry-Forward entry of `CONTEXT.md` already describes.

## Testing Approach

1. **Integration order.** A Run whose Run Branch holds `task_02`'s settlement
   before `task_01`'s, where `task_01` changes an input `task_02` declares,
   carries both through `reconcile --carry-forward`; a checkout change to that
   input still refuses the whole set with HEAD unchanged; the Tasks left
   unevaluated after a conflict follow integration order; the existing serial
   and drift cases stay green.
2. **Project Config.** A bounded Task commits `.roundfixrc.yml`; an unbounded
   Task settles failed with the lost-output reason and a dropped-path event; a
   Batch and a QA Report commit omit it with their reason; a pre-existing dirty
   Project Config stays out.
3. **Engine and store.** Store transitions and their refusals, the idle release
   and its refusals, the recorded Run ID, and every row of the re-entry table
   through `Engine.Retry` with fake recovery, including the refusals that leave
   the item unchanged and a re-entered archived item that pushes and merges
   exactly once.
4. **Workflow recovery.** Against real Git: an unresolved Run's completed Tasks
   reach the item branch with provenance while its failed Task stays pending; a
   Run found by item branch when none is recorded; a second carry that carries
   nothing; a moved input that refuses with HEAD unchanged; an archived and an
   active Spec inspected.
5. **Command and hand-off.** Through the public CLI: a retry that starts an
   owner, one handed to a live owner, one that reclaims a stale owner, argument
   and refusal exits; the owner makes a second pass for an item retried during
   its first and releases an idle queue after one; help and top-level usage.
6. **Proof timeouts.** A timeout then a pass succeeds after two attempts; two
   timeouts report `temporary` with rerun advice; a rejection and a cancelled
   command are attempted once.
7. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Carry-forward in integration order (depends on: none).
2. Project Config in Daemon commits (depends on: none).
3. Retry in the engine and store (depends on: none).
4. Item recovery in the delivery workflow (depends on: 1, 3).
5. The retry command and the owner hand-off (depends on: 1, 2, 4).
6. Proof timeouts (depends on: 5).
7. Terminal QA (depends on: 1, 2, 3, 4, 5, 6).

## Risks & Considerations

- **Shared guides make a chain.** Tasks 1, 5 and 6 edit
  `docs/user-guide/commands.md`, and Tasks 2, 5 and 6 edit the Roundfix skill
  and, with Task 5, `CONTEXT.md`, so Tasks 5 and 6 wait on them; Tasks 1, 2 and
  3 share no file and run together.
- **A carry the operator did not expect.** The retry carries Tasks forward on
  the item branch before it re-stages the item. That is the operator's explicit
  request, it carries only proved Tasks, and a refusal carries none, so the
  item branch is either unchanged or advanced by exactly the proved set.
- **An owner that never exits.** The owner loops only while an advanceable item
  remains after a pass; `Engine.Run` parks every item it cannot advance, so the
  loop ends unless an operator keeps retrying.
- **A longer Preflight.** One retry per timed-out tuple can double that tuple's
  proof time under load; the alternative, a refused dispatch, costs a full
  relaunch.
