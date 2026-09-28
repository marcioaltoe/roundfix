---
spec: 0175-cleanup-after-a-squash-merge
status: active
created: 2026-09-28
surfaces: [backend, cli, docs]
---

# Cleanup after a squash merge

## Executive Summary

Prove a terminal Run of a merged Spec commit by commit against the merged head,
where the merged head is the pull request head the Delivery Queue recorded or,
failing that, the default branch that carries the archived Spec. Release those
Runs automatically after a delivery merge and through `roundfix reconcile
--apply` otherwise. Key legacy Runs from their Run Worktree, own and sweep
carry-forward staging worktrees, refuse to delete a nested registered worktree,
and route every production `git worktree` add, remove, prune and move through
the Spec 0171 administration lock.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git, local files and the Run
  Database only; no credential and no network call. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0004, ADR-0014, ADR-0023, ADR-0044,
  ADR-0052, ADR-0053, ADR-0057, ADR-0115, ADR-0158 and ADR-0163 hold; ADR-0020,
  ADR-0022, ADR-0038, ADR-0056, ADR-0127, ADR-0141, ADR-0159 and ADR-0160 do
  not apply; this Spec
  adds ADR-0161; the gate is bound by ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Decision: when cleanup runs, and what a squash merge proves

Cleanup runs in two places, both through the same proof and the same
revalidated release (`ApplyTerminalRun`):

- **Automatically**, inside the delivery owner, right after a merged item's
  merge receipt and before `RemoveItemBranch`. The merge is recorded, the
  Spec's Runs are still in the Run Database, and the merged pull request head
  is still local.
- **Explicitly**, through `roundfix reconcile --apply`, for merges made outside
  the delivery loop and for merged items whose automatic release left a Run
  behind. Dry-run stays read-only.

A squash merge rewrites every commit, so ancestry against the default branch
proves nothing, and whole-tree equality fails as soon as `main` moves or the
Spec folder is archived. What survives is per-commit: a Daemon Task commit
names its Spec and Task in the `Roundfix-Spec` and `Roundfix-Task` trailers,
and the merged tree records whether that Task completed. The proof therefore
walks the commits the Run Branch holds and the merged head lacks, and settles
each one on its own. Tree equality survives only for the files a
non-Daemon commit changed. ADR-0161 records this decision.

## The merged-head proof

`internal/worktree/merged_head.go` adds:

```go
// MergedHead records the head a Spec's pull request merged at.
type MergedHead struct {
	SpecSlug     string
	TargetBranch string
	Head         string
	MergeCommit  string
	PullRequest  string
}

func InspectTerminalRunMerged(ctx context.Context, run store.Run, merged []MergedHead) (RunWorktreeReconciliation, error)
```

`InspectTerminalRun(ctx, run)` keeps its signature and behaves as
`InspectTerminalRunMerged(ctx, run, nil)`. Every existing precondition and
precedence in `inspectTerminalRun` stays: metadata, registration, `dirty` and
`unknown` are decided first, and a present target keeps its ancestry check and
its `supersedingQAReport` check.

**Choosing the merged head H.** Among `merged`, the entries whose `SpecSlug`
equals the Run's Spec slug are candidates. Exactly one distinct `Head` among
them is the record source (source label `pull request #<N> merged head <H12>`,
or `merged head <H12>` when `PullRequest` is empty). The record source is
unusable when there is no candidate, when the candidates name different heads,
or when `git cat-file -e <H>^{commit}` fails. The Run then falls back to the
default-branch source: H is the default branch head from
`resolveDefaultBranchHead` (label `default branch "<name>"`), usable only when
the Spec is archived at H (the existing
`newestQAReportAtHeadInDirectories(..., archivedQAReportDirectory(slug))`
check).

**When the proof runs.** For a present target, the proof runs only after
ancestry and `supersedingQAReport` both fail, and only when a source is usable.
Without a usable source, the Run keeps today's `unintegrated` result and
reason. For a deleted target, `inspectDeletedTargetRunByContent` runs the
proof with the record source when usable, else with the default-branch source.
The existing order stays: representation first, then the archived-Spec
requirement of the default-branch source, so a represented Run of an
unarchived Spec keeps the "Spec is not archived" reason.

**Representation.** Let U be `git rev-list <RunHead> ^<H>`. An empty U means
the Run Branch is contained in H. Otherwise each commit in U is exactly one of
the following:

1. A **Task commit**: its `Roundfix-Spec` trailer equals the Run's Spec slug,
   its `Roundfix-Task` trailer names a Task ID, and that Task's file at H, found
   through the Spec's `_tasks.md` at H under `docs/specs/<slug>/` or
   `docs/history/specs/<slug>/`, reads `completed` through
   `spec.CarryForwardStatus`.
2. A **QA Report commit**: its message matches `matchesQAReportCommitMessage`,
   it changes only paths under the Spec's QA directories, and
   `supersedingQAReport(H, <commit>)` would prove the report it adds
   superseded. A same-named report needs an identical blob; otherwise H's
   newest report must be newer.
3. **Other**: any remaining commit. The files every other commit changed are
   compared with `git diff --name-only -z --no-renames` between RunHead and H,
   restricted by pathspec to those files. A file present at RunHead and absent
   at H counts as Run-only, a file present at both with different content as
   differing shared, and a file the Run deleted that H still holds as a
   retained Run deletion. All three counts must be zero. `git diff` stays the
   comparison command, so a failed diff still preserves the Run with "could not
   prove integration".

A Task commit whose Task is not completed at H, or whose trailer names another
Spec, and a QA Report commit that H does not supersede, are never re-examined
as "other". They preserve the Run.

**Outcomes and reasons**, each bounded by `boundedReconciliationReason`:

| Outcome | State | Reason |
| --- | --- | --- |
| U empty, record source | `safe` | `Run Branch is contained in <source>` |
| No Task or QA commit represented, content equal | `safe` | `Run Branch content is fully represented on <source>` (unchanged text for the default branch) |
| At least one Task or QA commit represented | `superseded` | `Run work is superseded at <source>: <t> Task commit(s) completed, <q> QA Report commit(s) superseded` |
| A Task commit not completed | `unintegrated` | `Run commit <C12> is not represented at <source>: Task <id> is not completed` |
| A Task commit of another Spec | `unintegrated` | `Run commit <C12> is not represented at <source>: it belongs to Spec <slug>` |
| A QA Report commit not superseded | `unintegrated` | `Run commit <C12> is not represented at <source>: its QA Report is not superseded` |
| Content differs | `unintegrated` | today's `Run Branch content is not fully represented: <counts> against <source>` |

The `superseded` and `safe` results carry evidence, so `ApplyTerminalRun`
accepts them. `terminalRunReconciliationEvidence` also keeps `merged`, and
`revalidateTerminalRunApply` re-inspects with the same records, so a changed
record or a moved head makes the apply stale and preserves the Run.
`CONTEXT.md` updates the Run Worktree Reconciliation entry so that `safe`
includes containment in, or content represented at, the merged head, and
`superseded` includes the merged-head representation.

## Repository keys for legacy Runs

`applyRepositoryRootMigration` in `internal/store/store.go` splits its backfill
into a function that `migrate` also runs when the schema is already current.
It runs inside the same write transaction and only on write-mode `Open`, never
on `OpenReader`. For each Run with `repository_root = ''`, it keeps today's
derivation from `git_root` when `git_root/.git` exists. Otherwise, when
`work_dir` is set and `work_dir/.git` exists, it derives the key with
`roundconfig.RepositoryRoot(work_dir)` and accepts the result only when it
differs from the cleaned `work_dir`. A Run Worktree whose administrative entry
is gone resolves to itself and stays unkeyed. The schema version stays 20, and
the backfill updates only rows that are still unkeyed, so repeated opens are
idempotent.

`PruneTerminalReport` in `internal/worktree/worktree.go` stops comparing
`run.GitRoot` with `userRoot` by path. It compares repository keys:
`run.RepositoryRoot`, or `roundconfig.RepositoryRoot(run.GitRoot)` when it is
empty, against `roundconfig.RepositoryRoot(userRoot)`. It inspects a matching
Run with `GitRoot` set to `userRoot`, as `loadReconcileRuns` already does. A
Run keyed to another repository is skipped.

## Reconcile reads the merge record

`internal/cli/merged_release.go` adds `loadDeliveryMergedHeads(ctx, reader
*store.Store, repository string) ([]runworktree.MergedHead, error)`. It reads
`reader.DeliveryQueue` for the resolved repository path and, when it differs,
for the checkout's Git root as loaded. It maps every item at stage `merged`
with a non-empty merge commit and candidate head to a `MergedHead`, whose
`Head` is the last candidate commit, the head `mergeCandidate` verified against
the merged pull request. A missing queue yields no records, not an error.

`runReconcileCommand` loads the records once, and `inspectReconcileRuns` calls
`runworktree.InspectTerminalRunMerged` for every selected Run.
`applyReconcileBranchDisposition` no longer overwrites the classification or
action of a Run the merged-head proof already classified `safe` or
`superseded`, except under `--discard-superseded`, whose behavior is
unchanged. The JSON and text shapes are unchanged by this section, and the
evidence field carries the proof. `docs/user-guide/commands.md` (reconcile) and
the Roundfix skill's Run Worktree reconciliation section describe the
merged-head proof and its two sources.

## Automatic release after a delivery merge

`delivery.ItemWorkspace` gains:

```go
ReleaseMergedRuns(ctx context.Context, gitRoot string, item store.DeliveryQueueItem) error
```

In `Engine.Run`, a merged item first calls `ReleaseMergedRuns`, then
`RemoveItemBranch`, even when the release failed. The two errors are joined
into one `warning: cleanup failed: ...` blocker. A later clean pass clears the
warning, exactly as today.

`commandDeliveryWorkflow.ReleaseMergedRuns` in `internal/cli/deliver_workflow.go`
builds the item's `MergedHead`: Spec slug, item branch, last candidate commit,
merge commit and pull request number. It then calls
`releaseMergedSpecRuns(ctx, runStore, gitRoot, merged)` in the new
`internal/cli/deliver_release.go`, which does the following:

- It lists Runs with `ListRuns` for the repository and keeps implement Runs of
  that Spec slug.
- It leaves Active Runs alone and names them in the returned error.
- For each terminal Run, with `GitRoot` set to the repository, it calls
  `InspectTerminalRunMerged`. It applies `ApplyTerminalRun` to `safe` and
  `superseded` results, counts `released` as done, and collects every other
  result with its reason.
- It returns nil only when nothing was preserved or failed. Otherwise it
  returns one error that names each Run and its reason.

`docs/user-guide/commands.md` (deliver) and the Roundfix skill's Delivery queue
section state that a merge releases the Spec's provable Runs and that a kept
Run appears in the item's cleanup warning.

## Staging worktrees

`internal/worktree/staging.go` owns carry-forward staging:

- `AddCarryForwardStaging(ctx, repository, head, parent string)`
  (`CarryForwardStaging, error`). It creates `<parent or os.TempDir()>/roundfix-carry-forward-*`,
  writes `owner.json` there with the current PID and
  `store.OwnerProcessIdentity`, then runs `worktree add --detach <root>/worktree
  <head>` through `runWorktreeCommand`. `CarryForwardStaging.Remove(ctx)` runs
  `worktree remove --force` through the lock and removes the root.
- `InspectCarryForwardStaging(ctx, repository)` (`[]StagingCandidate, error`)
  lists registered worktrees whose basename is `worktree` and whose parent's
  basename starts with `roundfix-carry-forward-`, wherever they live. It reads
  their lock reason from `git worktree list --porcelain`. A candidate is
  stale when `owner.json` names a PID whose `store.OwnerProcessIdentity` fails
  or differs, or when there is no owner record and the entry is `locked` with
  reason `initializing`. Every other candidate is kept with a reason.
- `ReleaseCarryForwardStaging(ctx, repository, candidate)` holds the
  administration lock once for the whole act, through a new
  `withWorktreeAdminLock` helper in `adminlock.go`. It re-proves staleness,
  runs `git worktree remove --force --force <path>` directly through the
  runner, and then removes the `roundfix-carry-forward-*` root. Inside the
  helper it must never call `runWorktreeCommand`, which would wait on the lock
  it already holds.

`createCarryForwardStaging` in `internal/cli/carryforward.go` gains a `parent
string` argument. Production passes `""`, and `carryforward_test.go` passes
`t.TempDir()`. It delegates to `AddCarryForwardStaging` and stops running `git
worktree` itself. The reconcile report adds `stagingCandidates`
(`reconcileDebrisResult` with kind `staging`, worktree, owner PID, proof,
action, refusal reason). Kept candidates go to `preservedCandidates`, and
`debrisSummary` adds `stagingCandidates` and `stagingApplied`. Dry-run reports
without removing. `--apply` releases every stale candidate, and
`--carry-forward` releases them before creating its own staging. The text
report prints one `Staging candidate:` block per candidate and the two counts
on the debris summary line.

## Item cleanup and the lock for every caller

`removeUnregisteredItemWorktree` lists the registered worktrees of
`ref.UserRoot` before `os.RemoveAll`. It refuses, with an error naming the
nested path, when any registered worktree's canonical path lies under the item
path. `CleanupItem` returns that error, so delivery records it as the cleanup
warning.

`internal/worktree/disposal.go` adds `RemoveRegisteredWorktree(ctx, gitRoot,
path string) error` (`worktree remove <path>`, no force) and
`AddDetachedWorktree(ctx, repository, path, revision string) error` and
`RemoveWorktreeForce(ctx, repository, path string) error`, all through
`runWorktreeCommand`. `discardSupersededBranch` in `internal/daemon/reconcile.go`
removes the Run Worktree through `RemoveRegisteredWorktree`, and
`internal/speccheck/checkout.go` creates and removes its disposable checkout
through the other two. A guard test then fails when any non-test Go file under
`internal/` or `cmd/`, outside `internal/worktree`, passes the string literal
`"worktree"` followed by `"add"`, `"remove"`, `"prune"` or `"move"` as
consecutive arguments of one call.

## API Contracts

1. `worktree.InspectTerminalRunMerged(ctx, run, merged []worktree.MergedHead)`
   classifies a terminal Run with the merged-head proof, and
   `worktree.InspectTerminalRun(ctx, run)` equals it with no records.
2. A merged-head result is `safe` with reason `Run Branch is contained in
   <source>` or `Run Branch content is fully represented on <source>`,
   `superseded` with reason `Run work is superseded at <source>: ...`, or
   `unintegrated` with a reason naming the first unrepresented commit or the
   content counts.
3. `delivery.ItemWorkspace.ReleaseMergedRuns(ctx, gitRoot, item)` runs for
   every merged item before `RemoveItemBranch`, and its failure becomes the
   item's `warning: cleanup failed:` blocker without skipping the item-branch
   removal.
4. `roundfix-reconcile/v1` adds the top-level `stagingCandidates` array and the
   `debrisSummary` fields `stagingCandidates` and `stagingApplied`; every
   existing field keeps its name and meaning.
5. A write-mode `store.Open` keys an unkeyed Run from its Run Worktree when its
   recorded checkout is gone; the schema version stays 20.
6. `worktree.AddCarryForwardStaging`, `worktree.InspectCarryForwardStaging`,
   `worktree.ReleaseCarryForwardStaging`, `worktree.RemoveRegisteredWorktree`,
   `worktree.AddDetachedWorktree` and `worktree.RemoveWorktreeForce` are the
   only production paths outside `internal/worktree` that add, remove or prune
   a Git worktree.

## Coverage Map

- Goal 1 → The merged-head proof; API Contracts 1-2.
- Goal 2 → Automatic release after a delivery merge; API Contract 3.
- Goal 3 → Reconcile reads the merge record.
- Goal 4 → Repository keys for legacy Runs; API Contract 5.
- Goal 5 → Staging worktrees; Item cleanup and the lock for every caller; API
  Contracts 4 and 6.
- Core Feature 1 → The merged-head proof.
- Core Feature 2 → Automatic release after a delivery merge.
- Core Feature 3 → Reconcile reads the merge record.
- Core Feature 4 → Repository keys for legacy Runs.
- Core Feature 5 → Staging worktrees; Item cleanup and the lock for every
  caller.
- Core Feature 6 → Item cleanup and the lock for every caller.
- Success Metric 1 → Testing Approach 1 and 3.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 4.
- Success Metric 4 → Testing Approach 2.
- Success Metric 5 → Testing Approach 5 and 6.

## Integration Points

- **Spec 0171.** Added the worktree administration lock in
  `internal/worktree/adminlock.go`. Every release in this Spec goes through it.
- **Spec 0168.** Gave each delivery item its own worktree and recorded
  `half-removed item` cleanup, whose nested-worktree limit Core Feature 6
  closes.
- **Specs 0157 and 0162.** Introduced the durable repository key, whose
  backfill Core Feature 4 completes.

## Testing Approach

1. **Merged-head proof.** Real Git fixtures in `internal/worktree`: the Spec
   0172 shape (Task 01–04 commits in H, a unique Task 05 commit and a unique
   failed QA Report commit) is `superseded` against a record and against the
   default branch carrying the archived Spec. The Spec 0164 shape (archived
   Spec, and a default branch that later changed other files and a file the Run
   touched) is released. Each refusal is its own test: a plain unrepresented
   commit, a Task not completed at H, a Task commit of another Spec, a QA
   Report H does not supersede, a record for another Spec, and a dirty Run
   Worktree. The existing deleted-target tests stay green unchanged, and apply
   revalidation refuses a changed record.
2. **Repository keys.** A write-mode open keys an unkeyed Run from its Run
   Worktree, and the mirror cases stay unkeyed: a gone worktree, a worktree
   without a common directory, and a worktree of another repository keyed to
   that repository. Repeated opens are idempotent. The prune releases a keyed
   Run whose checkout was removed and skips one keyed elsewhere. `roundfix
   reconcile` lists the keyed Run. `TestJournalConsumerCorpusReplaysEveryConsumer`
   stays green.
3. **Reconcile.** `roundfix reconcile` with a seeded Delivery Queue merge
   record classifies and, with `--apply`, releases the Spec 0172 shape. The
   mirror cases each keep the Run: a dry-run, a record for another Spec, an
   unrepresented commit. A Branch Disposition does not replace the apply action.
4. **Delivery.** The engine calls the release before the item-branch removal,
   records a kept Run as the cleanup warning while still removing the item
   branch, clears the warning on a clean retry, and releases nothing for an
   unmerged item. The command workflow releases every provable Run of the
   merged Spec, keeps an unrepresented one, and never touches an Active Run or
   another Spec's Run.
5. **Staging.** A dead owner's staging and a legacy `locked initializing` entry
   are stale. A live owner's staging and an unlocked legacy entry are kept.
   Adding staging waits on a held administration lock. Dry-run reports without
   removing, `--apply` and `--carry-forward` release, and the JSON keeps every
   existing field.
6. **Item cleanup and the lock.** A half-removed item that holds a nested
   registered worktree is refused and the nested worktree keeps its files. One
   without a nested worktree is still removed. The exported removal waits on a
   held lock. The guard test fails on today's tree and passes after the routing.
7. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The merged-head proof (depends on: none).
2. Repository keys for legacy Runs (depends on: 1).
3. Reconcile reads the merge record (depends on: 2).
4. Automatic release after a delivery merge (depends on: 3).
5. Staging worktrees (depends on: 4).
6. Item cleanup and the lock for every caller (depends on: 5).
7. Terminal QA (depends on: 1, 2, 3, 4, 5, 6).

## Risks & Considerations

- **Shared files force a chain.** Tasks 1, 2 and 6 edit
  `internal/worktree/worktree.go`. Tasks 3 and 5 edit `internal/cli/reconcile.go`.
  Tasks 3, 4 and 5 edit `docs/user-guide/commands.md` and the Roundfix skill,
  and every Task reads the same implement-task instruction, so the graph is a
  chain. Each Task puts its new tests in files of its own.
- **Discarding superseded work.** The Task-commit rule releases a Run whose
  Task was redone differently. That is the Branch Disposition's existing
  meaning of superseded, and the merged head is the delivered version. A
  commit without Daemon trailers is never discarded on that ground.
- **Parallel Specs.** Onda 2 Specs authored in parallel may claim ADR-0161 or
  touch `internal/cli/reconcile.go`. The ordinal check reports an ADR claim
  once both Specs are on `main`.
