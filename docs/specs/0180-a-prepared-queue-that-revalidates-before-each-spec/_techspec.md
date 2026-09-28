---
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# A prepared queue that revalidates before each Spec

## Executive Summary

Add Delivery Revalidation to the Delivery Engine: after an item's worktree is
created, and before its first Run, the item passes the strict Spec Consistency
Check there and a premise check against the merge commits of earlier items.
Record queue limits in one new Run Database schema version and enforce them in
the engine and the retry transaction. Add `roundfix deliver plan`, a read-only
report of which Specs are approved to run, and make `deliver start` refuse a
Spec without delivery authority. Print the limits and one Pending Question in
`deliver status`. Rewrite the owned `implement-spec` skill so it hands
implementation to Roundfix. No exported function changes its signature.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; no credential and no new network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0044, ADR-0052, ADR-0057,
  ADR-0080, ADR-0091, ADR-0093, ADR-0094, ADR-0096, ADR-0104, ADR-0117,
  ADR-0130, ADR-0137, ADR-0139, ADR-0153, ADR-0155, ADR-0156, ADR-0158 and
  ADR-0160 hold; ADR-0020, ADR-0038, ADR-0053, ADR-0056, ADR-0097, ADR-0127
  and ADR-0159 do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/implement-spec/SKILL.md`,
  `skills/implement-spec/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Queue limits and the item warning in the store and the engine

`internal/store/store.go` raises `schemaVersion` by exactly one from its value
on the Task's starting main. It adds `delivery_queues.deadline_unix INTEGER NOT
NULL DEFAULT 0`, `delivery_queues.max_retries INTEGER NOT NULL DEFAULT 0`,
`delivery_queue_items.retry_count INTEGER NOT NULL DEFAULT 0` and
`delivery_queue_items.warning TEXT NOT NULL DEFAULT ''`. A fresh database
creates the columns. Every older schema that already has the delivery tables
gains them through idempotent column-existence checks, as
`deliveryWorktreeMigrationStatements` does, so the existing downgrade fixtures
in `internal/store/delivery_test.go` keep migrating.

`internal/store/delivery.go` adds:

```go
type DeliveryQueueLimits struct {
	Deadline   time.Time // zero means none; stored as UTC Unix seconds
	MaxRetries int       // zero means none
}

func (store *Store) CreateDeliveryQueueWithLimits(
	ctx context.Context, gitRoot string, specSlugs []string, limits DeliveryQueueLimits,
) (DeliveryQueue, error)
```

`CreateDeliveryQueue` keeps its signature and delegates with zero limits. A
negative `MaxRetries` is refused. `DeliveryQueue` gains `Limits
DeliveryQueueLimits` and `DeliveryQueueItem` gains `RetryCount int` and
`Warning string`, all read by `DeliveryQueue`. `UpdateDeliveryQueueItem`
persists `Warning` with the other item fields, so the item carries it through
every later stage; the Delivery Revalidation is its only writer. `RetryDeliveryQueueItem` keeps its signature. In its write
transaction it reads the queue's `max_retries` and the item's `retry_count`. It
refuses with an error wrapping the new sentinel `ErrDeliveryRetryLimit` when
the limit is non-zero and reached, changing nothing. Otherwise it increments
`retry_count` in the same `UPDATE` that moves the item.

`internal/delivery/engine.go` adds `BlockerQueueDeadline = "queue-deadline"`:

- `Engine.Run` parks an item that is still `queued` as `queue-deadline` when the
  queue has a deadline and `engine.clock.Now()` is not before it. The engine
  never calls `CreateItemBranch` for that item. An item in any later stage
  advances as before, so a started item runs to its own outcome.
- `Engine.Retry` refuses a `queue-deadline` item before any work, naming the
  deadline and `roundfix deliver start` as the way to continue. It refuses an
  item whose `RetryCount` has reached a non-zero `MaxRetries` before
  `UseItemBranch` or carry-forward, wrapping `store.ErrDeliveryRetryLimit`.
  The store's in-transaction check stays the guard against a concurrent retry.

## The Delivery Plan and the start refusal

A new file `internal/cli/deliver_plan.go` holds the two helpers the plan and
the Delivery Revalidation share:

```go
func strictSpecFindings(specsRoot, repoRoot, specSlug string) ([]speccheck.Finding, error)
func productionPremises(graph *spec.Graph) []string
```

`strictSpecFindings` runs `speccheck.Check(specsRoot, repoRoot, specSlug)`,
applies `speccheck.PromoteGaps`, and returns
`speccheck.GatePrecondition(result).Findings`. That is the verdict `roundfix
spec check <slug> --strict` and the Daemon's `qaGatePrecondition` reach, so the
three can never disagree. `productionPremises` returns, sorted and unique, the
`interface:` paths (`spec.ContextKindInterface`) of every non-`qa` Task that
end in `.go` and not in `_test.go`.

The same file implements `roundfix deliver plan [--json] [<slug>...]`,
dispatched from `runDeliverCommand`. It resolves the
configuration and Specs Root with `loadDeliveryCommand`. With no slug, it takes
every active Spec from `spec.ListActiveDetailed` and reports skipped
directories through `printSkippedSpecDiagnostics`. A slug `spec.Load` refuses
exits `2`. For each Spec, in the given order, it computes:

- the number of Tasks and the number whose status is not `completed`;
- the authorization reasons from `deliveryAuthorizationReasons(ctx, loaded,
  specsRoot, slug) []string`, a helper in the same file.

  The helper reads the Spec's authorization with
  `spec.ReadSpecAuthorization(ctx, loaded.GitRoot, specsRoot.Path, slug, "")`.
  A resolution whose outcome is not granted yields `authorization <outcome>:
  <reason code>`. A granted one that does not `Permit` every operation of
  `implement`, `commit`, `push`, `pull_request` and `merge` yields
  `authorization lacks <op>, <op>`.
- `spec check: <code>, <code>` from `strictSpecFindings(specsRoot.Path,
  loaded.GitRoot, slug)`;
- the shared premises: for each earlier Spec in the order, the paths both
  Specs' `productionPremises` name. A shared row is information: it tells the
  operator which later item will carry a `premise-changed` warning once the
  earlier one merges, and it never blocks a Spec or a start.

A Spec with no reason is `approved`; any other is `blocked`. The plan then
lists repository intent that is not approved to run:

- every `docs/backlog/*.md` with its front-matter `status`;
- every `docs/findings/*.md` with its front-matter `status`;
- every regular file under `docs/_inbox/`, with status `-`.

Each list is sorted by path, and a missing directory contributes nothing. Text
output is tab-separated, one row per fact, with the row kind first:

```text
spec	<slug>	<approved|blocked>	<unfinished>/<tasks>	<reason>; <reason>|-
shared	<slug>	<earlier-slug>	<path>, <path>
backlog	<path>	<status>
finding	<path>	<status>
inbox	<path>	-
```

`--json` prints one `roundfix-deliver-plan/v1` document carrying the same
facts:

- `schema`;
- `specs`, whose entries hold `slug`, `verdict`, `tasks`, `unfinishedTasks`,
  `reasons` and `sharedPremises` entries of `with` and `paths`;
- `intent`, whose entries hold `kind`, `path` and `status`.

The plan exits `0` when every reported Spec is approved and `1` when any is
blocked; a usage or preflight error exits `2`. It never opens the Run Database,
creates a worktree or writes a file.

`runDeliverStart` in `internal/cli/deliver.go` calls
`deliveryAuthorizationReasons` for every slug after `spec.Load` and before
`store.Open`. When any slug has a reason, it fails through `printDeliverFailure`
with exit `2`, names every such slug with its reasons and `roundfix deliver
plan`, and records no queue. Strict findings never refuse a start.

`deliverUsage` gains `roundfix deliver plan [--json] [<slug>...]` and a
Commands line. The top-level usage line in `internal/cli/cli.go` becomes
`roundfix deliver <plan|start|status|resume|retry|stop> [<slug> ...]`. The
`deliver` row of `TestRunCommandHelp` in `internal/cli/cli_test.go` expects the
plan line. The expectation of `TestTopLevelUsageNamesDeliverRetry` in
`internal/cli/deliver_retry_test.go` follows the new top-level line, keeping the
test's name. The tests that start a queue in `internal/cli/deliver_test.go`
first grant all five operations with `setImplementFixtureAuthorizationOperations`.

## Strict consistency and Delivery Revalidation

`internal/delivery/engine.go` adds:

```go
const BlockerRevalidationFailed = "revalidation-failed"

const WarningPremiseChanged = "premise-changed"

type Revalidation struct {
	Findings        []string // unique, sorted finding codes
	ChangedPremises []string // unique, sorted repository-relative paths
	ChangedBy       []string // the prior merge commits that changed at least one of them, in queue order
}

type ItemRevalidator interface {
	Revalidate(ctx context.Context, workDir, specSlug string, priorMerges []string) (Revalidation, error)
}
```

`EngineDependencies` gains `Revalidator ItemRevalidator` and `Log io.Writer`.
Both `validate` and `validateRetry` refuse an engine without a revalidator, so
no path can skip the check; a nil `Log` discards. `Engine.Run` passes each item
the merge commits of the items before it that are `merged` with a non-empty
merge commit, in queue position order. In the `queued` branch of
`advanceItem`, after the deadline check and after `CreateItemBranch` has
recorded the branch and worktree, and before the stage becomes `running`, the
engine calls `Revalidate` with the item worktree:

| Revalidation | Next step |
| --- | --- |
| at least one finding | park `revalidation-failed: <code>, <code>` |
| no finding | stage `running`, as today |
| an error | returned; `Engine.Run` parks the item as `delivery-error: <error>`, as today |

Whatever the next step, when `ChangedPremises` is non-empty the engine sets the
item's `Warning` to `premise-changed: <path>, <path> (merge <sha>, <sha>)` and
persists it with that step. It also writes `roundfix: warning: Delivery Queue
item <slug>: <warning>` to `Log`. A changed premise never parks an item and never
delays its Run. An item with no changed premise keeps an empty `Warning`, and
nothing is written to `Log`. The Runner is never called for a parked item,
which keeps its branch and worktree so a Delivery Retry can reach it.

`Engine.Retry` gains one step for an active Spec whose item records no Run ID.
This covers an item parked by revalidation and one whose revalidation errored.
Before carry-forward, it calls `Revalidate` in the item worktree with no prior
merges and refuses while `Findings` is non-empty. The refusal names the codes
and leaves the item unchanged. A retry never changes the recorded `Warning`.

`commandDeliveryWorkflow` implements `delivery.ItemRevalidator` in a new file
`internal/cli/deliver_revalidate.go`. `newCommandDeliveryEngine` in
`internal/cli/deliver_workflow.go` passes it as `Revalidator` and passes
`os.Stderr` as `Log`, which is the delivery console log of a detached owner.
`Revalidate` works as follows:

1. Resolves the Specs Root for the worktree with
   `roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)`.
2. Reports the codes of `strictSpecFindings(specsRoot.Path, workDir, specSlug)`.
3. Loads the graph with `spec.Load` and takes its `productionPremises`.
4. For each prior merge, reads `git diff --name-only <merge>^ <merge>` in the
   worktree through `workflow.git`. It reports the premises that diff names,
   and the merge whenever it names at least one.

A merge commit that cannot be read is an error, never an empty change set.

`runDeliverStatus` in `internal/cli/deliver.go` prints the item rows unchanged,
then one `Warning: <slug> <warning>` line for each item with a non-empty
`Warning`, in position order. The deliver section of
`docs/user-guide/commands.md` and the Delivery queue section of the Roundfix
skill describe the revalidation, the `revalidation-failed` blocker, the
`premise-changed` warning and the retry rule.

## Limit flags, status and the Pending Question

`parseDeliverStart` accepts `--max-duration <duration>`, a positive Go
duration, and `--max-retries <n>`, an integer of at least `1`, both declared as
value flags to `hoistCommandFlags`. A zero, negative or malformed value exits
`2`. `runDeliverStart` sets the deadline to `time.Now().UTC()` plus the
duration, truncated to whole seconds. It records the queue with
`CreateDeliveryQueueWithLimits` and prints to stdout, before the detached owner
report, one line:

```text
Limits: deadline <RFC 3339 UTC|none>, retries per item <n|none>, concurrency 1, spend not measured
```

`deliverUsage` keeps the line `roundfix deliver start <slug>...` and gains a
Flags block naming both flags. `runDeliverStatus` prints the item rows
and the warning lines unchanged, then the same `Limits:` line. When any item is parked, it prints:

```text
Pending question: <slug> parked <blocker>
Answer: <answer>
Waiting behind it: <n> parked item(s)
```

It prints the last line only when `n` is greater than zero.

A new file `internal/delivery/question.go` derives the question:

```go
type PendingQuestion struct {
	SpecSlug string
	Blocker  string
	Answer   string
	Waiting  int
}

func PendingQuestionFor(queue store.DeliveryQueue) (PendingQuestion, bool)
```

The question is the parked item with the lowest position. `Waiting` counts the
other parked items. The answer depends on the blocker:

| Blocker | Answer |
| --- | --- |
| `revalidation-failed…` | `amend the Spec on its item branch in <worktree>, then run roundfix deliver retry <slug>` |
| `queue-deadline` | `record a new queue for the remaining Specs with roundfix deliver start` |
| any other | `resolve the blocker, then run roundfix deliver retry <slug>` |

Nothing in the engine or the owner changes a parked item. Only a Delivery Retry
or a new queue does.

`CONTEXT.md` gains **Delivery Plan**, **Delivery Revalidation**, **Delivery
Queue Limit** and **Pending Question**, and its **Delivery Queue** entry names
the limits. The deliver section of `docs/user-guide/commands.md` and the
Delivery queue section of the Roundfix skill document the flags, the `Limits:`
line and the Pending Question.

## The implement-spec entry point

`.agents/skills/implement-spec/SKILL.md` is rewritten, and
`skills/implement-spec/SKILL.md` is regenerated with `make skills-sync`. The new
skill has the Supervisor:

1. Run `roundfix deliver plan <slug>...` and stop on a blocked Spec, reporting
   its reasons.
2. Hand a single Spec on the current branch to `roundfix implement --spec
   <slug>`, detached with `--detach` when the session may end. Hand a
   merge-through sequence to `roundfix deliver start [--max-duration
   <duration>] [--max-retries <n>] <slug>...`.
3. Monitor through `roundfix deliver status` or the Run's events, and ask the
   maintainer only the Pending Question.

The QA gate is the Spec's terminal Task, which the Daemon runs. The Supervisor
never writes code or tests, never runs a Task or the implement-task cycle, and
never runs the QA gate itself. The skill defers to the Roundfix skill for
command details. It keeps `disable-model-invocation: true`, and both of its
version declarations move to `0.1.0`. `make baseline-digests` runs after the
edit.

## API Contracts

1. `roundfix deliver plan [--json] [<slug>...]` prints one `spec` row per Spec
   and the `shared`, `backlog`, `finding` and `inbox` rows, or the
   `roundfix-deliver-plan/v1` document. It exits `0` when every Spec is
   approved, `1` when any is blocked and `2` on a usage or preflight error, and
   it writes nothing.
2. `roundfix deliver start` exits `2` and records no queue when any slug's
   authorization does not grant `implement`, `commit`, `push`, `pull_request`
   and `merge`.
3. A queued item parks as `revalidation-failed: <codes>` before its first Run
   when its strict check fails. An item whose declared production Go file an
   earlier item's merge changed records the warning `premise-changed: <paths>
   (merge <shas>)`, which `deliver status` prints as `Warning: <slug>
   <warning>` and the console log carries, and continues to its Run.
4. `Engine.Retry` of an active item with no recorded Run refuses while the
   strict check reports findings, and never changes the recorded warning.
5. `roundfix deliver start --max-duration <duration> --max-retries <n>` records
   a deadline and a per-item retry limit, and prints the `Limits:` line; each
   omitted limit is `none`.
6. A `queued` item at or after the deadline parks as `queue-deadline` without a
   worktree, and a retry of it is refused. A retry beyond the limit is refused
   with `store.ErrDeliveryRetryLimit` and changes nothing.
7. `roundfix deliver status` prints the `Limits:` line and at most one Pending
   Question.
8. `Store.CreateDeliveryQueue`, `Store.RetryDeliveryQueueItem` and every other
   exported signature are unchanged.

## Coverage Map

- Goal 1 → The Delivery Plan and the start refusal; API Contract 1.
- Goal 2 → The Delivery Plan and the start refusal; API Contract 2.
- Goal 3 → Strict consistency and Delivery Revalidation; API Contracts 3-4.
- Goal 4 → Queue limits and the item warning in the store and the engine; Limit flags, status and
  the Pending Question; API Contracts 5-6.
- Goal 5 → Limit flags, status and the Pending Question; API Contract 7.
- Goal 6 → The implement-spec entry point.
- Core Feature 1 → The Delivery Plan and the start refusal.
- Core Feature 2 → The Delivery Plan and the start refusal.
- Core Feature 3 → Strict consistency and Delivery Revalidation.
- Core Feature 4 → Queue limits and the item warning in the store and the engine; Limit flags,
  status and the Pending Question.
- Core Feature 5 → Limit flags, status and the Pending Question.
- Core Feature 6 → The implement-spec entry point.
- Success Metric 1 → Testing Approach 3.
- Success Metric 2 → Testing Approach 1, Testing Approach 6.
- Success Metric 3 → Testing Approach 2, Testing Approach 4.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 5.
- API Contracts 1-2 → The Delivery Plan and the start refusal.
- API Contracts 3-4 → Strict consistency and Delivery Revalidation.
- API Contract 5 → Limit flags, status and the Pending Question.
- API Contract 6 → Queue limits and the item warning in the store and the
  engine.
- API Contract 7 → Limit flags, status and the Pending Question.
- API Contract 8 → Queue limits and the item warning in the store and the
  engine.

## Integration Points

- **Delivery Retry.** The revalidation re-check and the limit refusals run
  inside `Engine.Retry` before its existing refusals change anything, so a
  refused retry still leaves the item unchanged.
- **Spec Consistency Check.** `strictSpecFindings` composes the exported
  `speccheck` functions the CLI and the Daemon already compose; no detector
  changes.
- **Authorization.** The plan and the start refusal read the committed record
  through `spec.ReadSpecAuthorization` at the checkout's `HEAD`, the same
  reader `commandDeliveryWorkflow.Authorization` uses before publication.
- **Run Database.** One schema version is added. An older binary refuses the
  migrated database and names `roundfix migrate`, as Spec 0174 made it do.

## Testing Approach

1. **Revalidation in the engine.** With a fake revalidator:
   - a finding parks `revalidation-failed` before any Run;
   - a changed premise records the `premise-changed` warning, logs it and
     continues to the Run, while an item with no overlap records and logs no
     warning;
   - a clean item runs;
   - an error parks `delivery-error`;
   - prior merges arrive in position order, excluding unmerged items;
   - a retry of an item with no Run refuses while findings remain, proceeds when
     they clear, and keeps the recorded warning;
   - an engine without a revalidator is refused.
2. **Limits in the store and the engine.** Covers the limits round trip and
   `none`; a refused negative retry limit; the retry count increment; a refusal
   at the limit that leaves the item unchanged; migration of the previous
   schema version with rows kept; a queued item parked at the deadline without
   `CreateItemBranch`; a started item continuing past it; and retries refused
   for a `queue-deadline` item and at the limit, before carry-forward.
3. **Plan and start.** Through the public CLI in a disposable repository:
   - an approved and a blocked Spec and their exit codes;
   - JSON parity;
   - shared premises;
   - intent rows;
   - no Run Database and an unchanged checkout;
   - an unknown slug exiting `2`;
   - a start refused for a Spec lacking `merge` with no queue recorded, and a
     start accepted once all five operations are granted;
   - help and top-level usage.
4. **Flags, status and the question.** Covers limit recording and printing;
   omitted limits printed as `none`; refused flag values; one Pending Question
   for two parked items and none without a parked item; each blocker's answer;
   the question unchanged after an owner pass with the clock advanced; and
   help naming the flags.
5. **Skill.** Phrase checks on the canonical skill for the three hand-off
   commands, absence checks for the implement-task loop, and
   `make skills-sync-check`.
6. **Workflow revalidation against real Git.** A declared file that an earlier
   merge removed yields `SC-REF-UNRESOLVED` only after that merge. A declared
   production Go file that an earlier merge changed is named. A changed test,
   guide or undeclared file is not. An unreadable merge commit is an error.
7. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Queue limits and the item warning in the store and the engine, task_02
   (depends on: none).
2. The Delivery Plan and the start refusal, task_03 (depends on: none).
3. Strict consistency and Delivery Revalidation, task_01 (depends on: 1, 2).
4. Limit flags, status and the Pending Question, task_04 (depends on: 1, 2, 3).
5. The implement-spec entry point, task_05 (depends on: 2, 4).
6. Terminal QA, task_06 (depends on: 1, 2, 3, 4, 5).

## Risks & Considerations

- **Shared files make a chain.** task_02 and task_01 both edit
  `internal/delivery/engine.go`, and task_01 writes the item warning column
  task_02 adds. task_01 reuses the helpers task_03 adds, and task_01, task_03
  and task_04 edit `internal/cli/deliver.go`, `internal/cli/deliver_test.go`,
  `docs/user-guide/commands.md` and the Roundfix skill. task_02 and task_03
  share no file and run together; the Task ids keep their authored numbers, so
  the graph, not the numbering, gives the order.
- **A warning nobody reads.** Two Specs that edit the same production Go file
  are common, so the premise check records a warning instead of stopping the
  queue. The warning is durable on the item, printed by `deliver status` and
  written to the console log, and the plan's `shared` rows predict it before
  the queue starts. The Task's own Verification stays the detector for a real
  break.
- **Schema versions that collide.** Another Spec that raises `schemaVersion`
  and merges first moves the starting value. The requirement is relative to
  the Task's starting main, and the idempotent column checks keep either order
  migrating.
- **Earlier queue items that change the same files.** Spec 0175 changes
  `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go` and
  `internal/cli/deliver_test.go` and adds a method to `delivery.ItemWorkspace`.
  This Spec adds a separate interface and dependency field instead of changing
  an existing one, and its Tasks name every construction site to update.
