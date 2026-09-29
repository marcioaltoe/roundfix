---
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: active
created: 2026-09-29
surfaces: [backend, cli, docs]
---

# Delivery that reviews and retries from where the item stands

## Executive Summary

Resolve the merge base of the review head and the selected base ref once,
right after the base ref resolves, and use that one commit everywhere the
review reads a base. Record the resolved tip beside it for provenance only.
Make the shared Task Carry-Forward inspection record a Task already
`completed` on its target as nothing to carry, so reconcile, the implement
Preflight and a Delivery Retry stop refusing a Run over work that is already
home. Make the retry walk every terminal Implement Run of the item's Spec on the
item branch, newest first. The trade-off accepted: a retry that carries one Run
and then meets a refusing older Run leaves the first Run's proved Tasks on the
item branch instead of rolling them back. Each Run stays whole-or-nothing, and
no Run Database schema changes.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git, the Run Database
  and the configured reviewer runtime only; no credential and no new network
  call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0169 and ADR-0170 (this Spec) are
  implemented here; ADR-0014, ADR-0026, ADR-0044, ADR-0052, ADR-0053,
  ADR-0057, ADR-0080, ADR-0090, ADR-0091, ADR-0093, ADR-0094, ADR-0096,
  ADR-0104, ADR-0117, ADR-0139, ADR-0151, ADR-0153, ADR-0155, ADR-0156,
  ADR-0158 and ADR-0165 hold as the PRD states. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The review diffs from the merge base

`runReviewCommand` in `internal/cli/review.go` keeps
`resolveReviewBaseCommit` as the resolver of the base ref's commit, the tip.
Immediately after it, a new `resolveReviewMergeBase` runs
`git merge-base <tip> <head>` through the same `preflight.GitRunner` with the
head from `preflight.InspectGit`. A Git failure or an empty answer fails the
command through `printReviewCommandFailure` with exit `2`, before
`newReviewRecord`, the reuse lookup, the provider switch or any readiness
probe. The message names both commits and says that the head shares no history
with the base, ending with `pass --base <ref>`.

```go
func resolveReviewMergeBase(
	ctx context.Context, gitRoot, tip, head string, runner preflight.GitRunner,
) (string, error)
```

The merge base then replaces the tip in every consumer:

- `newReviewRecord` receives it as `baseCommit`, and the record's new
  `BaseTipCommit string` field (`json:"baseTipCommit,omitempty"`) receives the
  tip. `validateReviewRecord` does not require the new field, so older
  records stay readable.
- `reusableReviewRecord` keeps comparing `BaseCommit`, now the merge base, and
  never compares `BaseTipCommit`.
- `runConfiguredReviewSession` passes it to `reviewCandidateDiff`,
  `reviewCandidateSpecContexts` and `buildReviewPrompt`, so the reviewer's
  diff, the changed-Spec list with `specs`, `skippedSpecs` and
  `archivedSpecs`, and the prompt's `Base commit:` line all start at the merge
  base. `runReviewSession` callers pass the merge base too.

The prompt text itself does not change. `internal/delivery/engine.go` needs no
change: it parks `corrective-spec-required` from the record's `archivedSpecs`,
which now names only archived Specs the candidate changed.

The Pre-PR review section of the Roundfix skill and the review section of
`docs/user-guide/commands.md` say that the review diffs from the merge base of
the current head and the selected base, and that the record names that merge
base as `baseCommit` and the base tip as `baseTipCommit`. Both documents drop
the phrase that the diff runs "from the selected base to the current head".
They say that a findings verdict stands while the base branch moves, and that
a head sharing no history with the base exits `2`. The `CONTEXT.md` entry
**Delivery Base** gains the sentence
`The pre-PR review diffs the candidate from it.`

## A completed Task is nothing to carry

`internal/cli/carryforward.go` adds the action
`carryForwardCompletedAction = "already completed; nothing to carry"`.

`inspectCarryForwards` reads, per selected Run, which Tasks are completed on
the target. When `<resolvedSpecsRoot.Path>/<run.SpecSlug>` exists, it loads
that Task Graph with `spec.Load` and keeps the IDs whose status is
`spec.StatusCompleted`. When the folder does not exist, no Task counts as
completed. Any other load error fails the inspection. It passes that set to
`inspectCarryForwardsForRun`. There, a Task in the set produces a candidate
with the completed action, the Run ID, the Task file, its settlement commit
when it has exactly one, and no refusal reason. It comes before every proof and
is never staged. A completed Task therefore never appears among the Tasks
`unstagedCarryForwards` reports after a conflict. The remaining Tasks keep
integration order and every existing proof.

Every reader of the candidates follows the new action:

- `specCarryForward.wouldCarry` holds when at least one candidate is ready and
  none carries a refusal reason. `carriable` still counts the ready ones.
- `carryForwardRefusalReason` is unchanged. A completed candidate has no
  reason, and zero candidates still mean "Run has no Tasks with completed
  settlement evidence".
- `applyCarryForwards` in `internal/cli/reconcile.go` stages only ready
  candidates. With none, it returns before touching the checkout.
- After a successful apply, reconcile marks only the previously ready
  candidates `carried forward`. Completed candidates keep their action in the
  text and `roundfix-reconcile/v1` reports.
- `reportImplementNonCarriableCarryForwards` in `internal/cli/implement.go`
  notes only a Run whose candidates carry a refusal, so a Run whose Tasks are
  all completed on the checkout is silent. `selectImplementCarryForward` names
  only ready Tasks.

The reconcile and implement Preflight paragraphs of
`docs/user-guide/commands.md` and of the Roundfix skill say that a Task already
completed on the checkout is reported as `already completed; nothing to carry`
and never refuses the set. They drop the statement that a caller "effectively
gets one carry-forward for overlapping work". The **Task Carry-Forward** entry
of `CONTEXT.md` replaces "a carried Task's own file becomes a moved input
afterwards, so the act does not repeat itself" with
`A Task already completed on the target is nothing to carry, so the act does not repeat itself.`

## A Delivery Retry carries from every Run of its item

`internal/delivery/engine.go` changes the recovery result:

```go
type CarriedRun struct {
	RunID   string
	Carried []string // Task IDs in the Run's integration order
}

type CarryForwardResult struct {
	RunID string       // newest Implement Run of the Spec on the item branch
	Runs  []CarriedRun // Runs that moved at least one Task, newest first
}
```

`Engine.Retry` still records `carried.RunID` on the item and returns the
result as `CarriedFrom`. `printDeliverRetryResult` in `internal/cli/deliver.go`
prints one `Carried forward from Run <run-id>: <task>, <task>` line per entry
of `Runs`, in order, before `Retried <slug>: <blocker> -> <stage>`.

In `internal/cli/deliver_workflow.go`, `deliveryCarryForwardRun` becomes
`deliveryCarryForwardRuns`. A recorded Run ID must still name an Implement Run
of the Spec, or the retry fails naming it. The result is that Run plus every
Implement Run of the Spec whose recorded local branch is the item branch, from
`ListRuns` over the repository, which is already newest first. Duplicates are
removed and the list is ordered by creation time, newest first. `RunID` is the
first of them. `CarryForward` walks the list:

1. It skips a Run whose outcome `runworktree.CarryForwardAcceptedOutcomes`
   does not accept, and a Run with no Task that `loadReconcileTaskCoverage`
   settled completed.
2. It rereads the item's Task Graph and skips a Run whose settled-completed
   Tasks are all `completed` there. That check needs no Run Worktree.
3. It refuses when the Run Worktree is gone or the Specs Root is external, as
   today.
4. It runs `inspectCarryForwards` for that Run in the item worktree. That
   inspection now records completed Tasks as nothing to carry. A refusal
   reason refuses the retry. Otherwise it applies the ready candidates and
   appends a `CarriedRun` when any moved.

`deliveryCarryForwardRefusal` gains the Runs already carried in this retry,
and its message names them before the refusal. Its `NextAction` still names
`roundfix reconcile <refusing-run-id> --carry-forward` in the item worktree and
`roundfix deliver retry <slug>`.

The retry paragraph of the Delivery queue section of the Roundfix skill and the
`deliver retry` part of `docs/user-guide/commands.md` say that the retry
considers every terminal Implement Run of the item's Spec on the item branch,
newest first, and prints one line per Run it carried from. The
**Delivery Retry** entry of `CONTEXT.md` says it performs Task Carry-Forward
`from every Run of the item's Spec, newest first`.

## API Contracts

1. `roundfix review [--base <ref>]` diffs, lists changed Specs and prompts
   from `git merge-base <base-tip> HEAD`. Its record's `baseCommit` is that
   merge base, and `baseTipCommit` is the tip. A findings record is reused by
   repository, `baseCommit`, `headCommit` and provider. A head with no merge
   base exits `2` before any reviewer call.
2. Task Carry-Forward records a Task already `completed` on its target with
   action `already completed; nothing to carry`. It never stages, proves or
   refuses that Task. `roundfix reconcile <run-id> --carry-forward` exits `0`
   when every candidate is completed or carried.
3. The Implement Command's carry-forward Preflight refuses only for ready
   Tasks, and notes only Runs whose candidates refuse.
4. `roundfix deliver retry <slug>` carries from every terminal Implement Run of
   the item's Spec on the item branch, newest first. It records the newest Run
   on the item and prints one `Carried forward from Run <run-id>: …` line per
   Run that moved Tasks.

## Coverage Map

- Goal 1 → The review diffs from the merge base; API Contract 1.
- Goal 2 → The review diffs from the merge base; API Contract 1.
- Goal 3 → A completed Task is nothing to carry; API Contracts 2–3.
- Goal 4 → A Delivery Retry carries from every Run of its item; API Contract 4.
- Core Feature 1 → The review diffs from the merge base.
- Core Feature 2 → The review diffs from the merge base.
- Core Feature 3 → A completed Task is nothing to carry.
- Core Feature 4 → A Delivery Retry carries from every Run of its item.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 3.
- API Contract 1 → The review diffs from the merge base.
- API Contracts 2–3 → A completed Task is nothing to carry.
- API Contract 4 → A Delivery Retry carries from every Run of its item.

## Integration Points

- **Spec 0179's review record.** The dispositions ledger and the reuse rule
  keep working; only the base they are keyed on moves to the merge base.
- **Spec 0173's Delivery Retry and Task Carry-Forward.** The retry keeps its
  re-entry table, guarded transition and owner hand-off. Task Carry-Forward
  keeps integration order and every proof for the Tasks that remain.
- **The delivery engine's review stage.** It consumes `archivedSpecs` from the
  record unchanged.

## Testing Approach

1. **Review base.** Through `runReviewCommand` with the existing fake reviewer
   runner, in a real repository where the default branch advances after the
   candidate branch is cut:
   - the prompt's diff carries only the candidate's change;
   - `baseCommit` is the merge base and `baseTipCommit` is the tip;
   - a main-side edit of an archived Spec leaves `archivedSpecs` empty, while
     the candidate's own edit of an archived Spec is listed;
   - a findings record from before main moved is reused with no reviewer call;
   - an orphan base ref exits `2` with no reviewer call.

   The existing default-base and prompt tests stay green.
2. **Completed on the target.** Through the public `reconcile` command and the
   implement Preflight inspection, against real Git:
   - a Run whose `task_01` is completed on the checkout carries only
     `task_02`, and JSON reports `task_01` as `already completed; nothing to
     carry`;
   - a repeated carry exits `0` with `HEAD` unchanged;
   - a moved input of the remaining `task_02` still refuses with `HEAD`
     unchanged;
   - the Preflight names only the ready Task, and stays silent for a Run whose
     Tasks are all completed.

   `TestCarryForwardRefusesRatherThanCarryingASubset` stays green unchanged.
3. **Retry over several Runs.** Through `commandDeliveryWorkflow.CarryForward`
   against real Git and the Run Database, with two BudgetExceeded Runs:
   - with the item recording the older Run, whose Tasks are already on the
     item branch, the newer Run's Tasks are carried and `RunID` names the
     newer Run;
   - two Runs holding disjoint proved Tasks both carry, newest first;
   - a fully completed Run with a gone Run Worktree is skipped;
   - an older Run refusing after a newer one carried names both Runs, and the
     item head holds only the newer Run's Tasks.

   Through `Engine.Retry` with fake recovery, the newest Run is recorded on the
   item. Through `printDeliverRetryResult`, one line is printed per Run.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The review diffs from the merge base (depends on: none).
2. A completed Task is nothing to carry (depends on: 1 — both edit
   `docs/user-guide/commands.md`, the Roundfix skill and `CONTEXT.md`).
3. A Delivery Retry carries from every Run of its item (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **Shared guides make a chain.** All three slices edit the same three guides,
  so the graph runs them in sequence. None of their code overlaps.
- **A reused review after a merge of main.** Merging main into the candidate
  changes both the head and the merge base, so the key changes and the
  candidate gets a fresh review. Nothing stale is reused.
- **Trusting the target's status.** A completed status is written by the Daemon
  at settlement or by carry-forward. A hand edit is outside every proof in the
  Task Graph, and the PRD records it as a limit.
- **A partly advanced item branch after a refusal.** The refusal names the
  Runs already carried. Their Tasks are proved and complete, and the next
  retry treats them as completed.

## Decisions

- Merge base resolved once at the base resolution. See ADR-0169.
- Reuse keyed on the merge base; tip kept for provenance. See ADR-0169.
- Completed-on-target filter in the shared inspection. See ADR-0170.
- Newest Run first, one Run at a time, each whole-or-nothing. See ADR-0170.
