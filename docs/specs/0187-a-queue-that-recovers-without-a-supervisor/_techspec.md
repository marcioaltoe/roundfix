---
spec: 0187-a-queue-that-recovers-without-a-supervisor
prd: _prd.md
created: 2026-09-30
---

# A queue that recovers without a supervisor — Technical Spec

## Executive Summary

Each of the four defects comes from Roundfix reading a fact from the wrong
place:

- The post-merge cleanup reads unrefreshed local refs.
- The retry refusal reads only the Tasks' moved inputs, not the commit that
  moved them.
- The queue never compares its owner binary with its items' main.
- The authorization audit reads the grant at the item's fork point, not the
  grant the Task ran under.

Each fix is local to the reader that erred and adds no command, flag or schema.
The trade-off this design accepts is in the retry hint. Roundfix prints a
destructive recovery (a branch reset) for the operator to run instead of
performing it. The alternative, an automatic reset and cherry-pick, would put a
history rewrite inside an unattended owner. The investigation also found the
audit's real rule: the grant is read at `merge-base <delivery target>
<commit>^1`. ADR-0178 widens that rule.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the warning and the
  hint name existing commit and Run IDs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the Run
  Database only; the one `git fetch` targets the delivery remote the queue
  already fetches when it creates an item. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0178 (this Spec), ADR-0161,
  ADR-0170, ADR-0053, ADR-0090, ADR-0158, ADR-0160, ADR-0081, ADR-0149,
  ADR-0166, ADR-0167, ADR-0169 and ADR-0176 hold as the PRD states. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155
  and ADR-0156, and by ADR-0093 and ADR-0094. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix skill files ride the standing
  grant of 2026-09-18, recorded in [_authorization.md](_authorization.md) under
  the maintainer's 2026-09-30 request to release every remaining fix; bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No new package or seam. Each fix extends the existing reader:

| Defect | Reader | File |
| --- | --- | --- |
| Cleanup before fetch | `resolveMergedReleaseEvidence` | `internal/cli/deliver_workflow.go` |
| Retry hint | `deliveryCarryForwardRefusal` and a new helper beside `carryForwardRefusalReason` | `internal/cli/deliver_workflow.go`, `internal/cli/carryforward.go` |
| Owner warning | `commandDeliveryWorkflow.Revalidate` and the engine's item start | `internal/cli/deliver_revalidate.go`, `internal/delivery/engine.go` |
| Grant ran under | `detectMechanicalAuthPaths` and `mechanicalAuthorizingRevision` | `internal/speccheck/mechanical.go` |

## Implementation Design

### Interfaces

```go
// internal/cli/carryforward.go
// carryForwardAmendments returns, oldest first, the non-Task commits in
// run.HeadSHA..HEAD of repository that changed the moved inputs, when every
// refused candidate was refused only for moved inputs and no Task commit in
// that range touched them. ok is false otherwise.
func carryForwardAmendments(ctx context.Context, repository string, run store.Run,
	candidates []spec.CarryForward) (amendments []string, ok bool, err error)

// internal/delivery/engine.go
type Revalidation struct {
	Findings        []string
	ChangedPremises []string
	ChangedBy       []string
	OwnerWarning    string // "" or "owner-older-than-main: …"
}
```

### Data Models

No schema change. `DeliveryQueueItem.Warning` keeps one string. When both a
`premise-changed` warning and an owner warning exist, they are joined with
`; `, premise first. A retry keeps the recorded warning, as today.

### Cleanup refreshes the default branch first

`resolveMergedReleaseEvidence` resolves the delivery remote exactly as
`CreateItemBranch` does: `Watch.PushRemote`, or `origin` when that is empty or
the workflow has no loaded config. The remote exists when `git remote get-url
<remote>` succeeds.

- **With the remote.** It runs `git fetch <remote> <default>` before resolving
  anything. A fetch failure returns `refresh default branch "<name>": <err>`,
  which the engine records as today's cleanup warning and retries on the next
  owner pass. The default head is `refs/remotes/<remote>/<default>`, and the
  merge commit must be its ancestor.
- **Without the remote.** Today's local resolution is unchanged.

The merge-commit and candidate-head checks are otherwise unchanged.

### Retry hint

`CarryForward` calls `carryForwardAmendments` when
`carryForwardRefusalReason(candidates)` refuses.

- The helper returns ok only when:
  - every candidate with a refusal has `InputsMoved` and at least one moved
    input;
  - `run.HeadSHA` resolves;
  - `git log --format=%H%x1f%(trailers:key=Roundfix-Task,valueonly,unfold)
    <run.HeadSHA>..HEAD -- <moved inputs>` lists at least one commit, and none
    of them carries a `Roundfix-Task` trailer.
- `deliveryCarryForwardRefusal` gains `amendments []string`. With amendments,
  `Error()` appends `; amended by <sha>, …`, and `NextAction()` returns the
  ordered commands:

```text
git -C <worktree> branch roundfix-amended-<run-id> HEAD
git -C <worktree> reset --hard <first-amendment>^
(cd <worktree> && roundfix reconcile <run-id> --carry-forward)
git -C <worktree> cherry-pick <amendment> …
roundfix deliver retry <slug>
```

Every printed argument that comes from state (the worktree path, the Run ID,
commit IDs and the slug) is POSIX single-quoted through one helper, so a path
with spaces or shell metacharacters is copied as one literal argument. Roundfix
only prints these commands and never runs them.

Without amendments, both texts are unchanged.

### Owner warning

The owner check runs once, at item start, in the first `Revalidate` before the
item's first Run. There the item worktree's `HEAD` is exactly the starting
main the item branch was created from. That commit is the `<starting main>`
below, and the result is recorded in the item's `Warning`. Later
revalidations (retries) never recompute it, so a later merge of main into the
item cannot move the comparison. After the strict findings it calls
`spec.ResolveAuditorEvidence(ctx, workDir, <starting main>, app.Auditor())`.

- It records a warning only when `SelfAudit` is true and `Ancestry` is
  `app.AncestryOlder`, and `git diff --name-only <build> <starting main> -- cmd internal
  go.mod go.sum` in `workDir` lists at least one path.
- It then sets `OwnerWarning` to `owner-older-than-main: owner build
  <build[:12]> predates starting main <starting main[:12]>`.
- A Git error in the diff leaves the warning empty and never fails
  revalidation.

`advanceItem` joins `OwnerWarning` into `item.Warning` and logs it like
`premise-changed`. It never parks for it.

### Grant ran under

In `detectMechanicalAuthPaths`, the fork-point read stays first. When it is
not `granted`, or when a changed governed path is outside its bounded paths,
the detector tries the parent's grant. It reads the blob of
`<authorizationPath>` at `<commit>^1` and compares it with the blob at the
same path at the delivery target's current tip (`<targetRevision>`) only.
An identical record anywhere else in history never matches, so a grant that
was narrowed or revoked on the default branch can never be revived.

- On a match, it reads the grant at `<targetRevision>` through
  `readMechanicalAuthorization`. The read's `Revision` is that commit, and the
  bounded-path check reruns against it.
- With no match, or a missing parent blob, the fork-point result stands.
- The self-approval check still runs first on the commit's changed paths.
- `!sameRepository` keeps today's behavior.

### API Contracts

1. API Contract: `roundfix deliver status` — a merged item whose merge commit
   exists only on the delivery remote shows no `cleanup failed` warning after
   the next owner pass. An item started by an owner older than its starting
   main shows `Warning: <slug> owner-older-than-main: owner build <commit>
   predates starting main <commit>`.
2. API Contract: `roundfix deliver retry <slug>` — a refusal caused only by
   amending commits exits `2`. It names them in the reason (`amended by
   <sha>, …`) and prints the five ordered commands as its next action. Every
   other refusal keeps today's text.
3. API Contract: QA report "Authorization audit inputs" — a Task commit
   authorized by the grant it ran under shows outcome `granted`, and its
   `Revision` is the delivery-target commit that holds the matching record.

## Coverage Map

- Goal 1 → Cleanup refreshes the default branch first; API Contract 1.
- Goal 2 → Retry hint; API Contract 2.
- Goal 3 → Owner warning; API Contract 1.
- Goal 4 → Grant ran under; API Contract 3.
- Core Feature 1 → Cleanup refreshes the default branch first.
- Core Feature 2 → Retry hint.
- Core Feature 3 → Owner warning.
- Core Feature 4 → Grant ran under.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.

## Integration Points

- **Delivery remote.** The fetch uses the remote the item creation already
  fetches. No credential is added.
- **Auditor Staleness.** The owner warning reuses
  `spec.ResolveAuditorEvidence` unchanged.
- **QA report writer.** `WriteMechanicalResult` already prints the read's
  `Revision`, so it needs no change.

## Testing Approach

1. **Cleanup.** A new `internal/cli/deliver_release_refresh_test.go` uses a
   bare `origin` and two clones. The item is merged and pushed from one clone,
   and `ReleaseMergedRuns` is run in the other without a prior fetch; it
   releases. A clone whose `origin` is unreachable returns a
   `refresh default branch` error. The existing
   `internal/cli/deliver_release_evidence_test.go` tests, which have no
   remote, pass unchanged.
2. **Retry hint.** A new `internal/cli/deliver_retry_amendment_test.go` has a
   Run with a settled task, then a non-Task commit on the item branch that
   edits `_prd.md`. `CarryForward` refuses, and the refusal's `Error()` and
   `NextAction()` carry the amending SHA and the five commands in order. A
   refusal caused by a Task commit touching the input keeps today's texts.
3. **Owner warning.** A new `internal/cli/deliver_owner_staleness_test.go`
   sets `app.BuildCommit` to an older commit of a temporary repository whose
   HEAD changed `internal/x.go`, and `Revalidate` returns the warning. A HEAD
   that changed only `docs/` returns none. A build commit absent from the
   repository returns none. A new `internal/delivery/owner_staleness_test.go`
   proves that `advanceItem` joins the warning, logs it and does not park.
4. **Grant ran under.** A new
   `internal/speccheck/mechanical_grant_ran_under_test.go` builds a temporary
   repository:
   - main widens the grant after the fork;
   - the item branch cherry-picks the widened record, then a Task commit
     changes the newly bounded path. It is granted, with `Revision` equal to
     main's widening commit.
   - A parent grant that differs from every main version refuses.
   - A Task commit that edits the grant itself is refused as self-approval.
5. **Docs.** Each Task phrase-checks its own behavior in
   `docs/user-guide/commands.md` and the Roundfix skill, plus
   `make skills-sync-check`.

## Build Order

1. Cleanup refreshes the default branch first, with its guide and skill
   text, task_01 (depends on: none).
2. Retry hint, with its guide and skill text, task_02 (depends on: 1). Both
   edit `internal/cli/deliver_workflow.go`, the guide and the skill.
3. Owner warning, with its guide and skill text, task_03 (depends on: 2),
   serialized because every Task edits the guide and the skill.
4. Grant ran under, with its guide and skill text, task_04 (depends on: 3).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A network fetch inside cleanup.** Cleanup already runs right after a
  merge. The same fetch runs at item creation, and a failure stays a warning
  that is retried.
- **A printed reset.** The commands save the branch first, and nothing runs
  them automatically.
- **Warning noise.** Only source paths count, and only a build commit present
  in the repository.
- **Grant forgery.** A matching record must exist on the default branch, where
  it passed review and the maintainer's grant rule. The item branch alone can
  never authorize.

## Decisions

- The delivery remote is resolved the way item creation resolves it.
- The hint is printed, never executed.
- The owner warning is based on source ancestry.
- The grant a Task ran under is accepted when the delivery target holds it.
  See ADR-0178.
