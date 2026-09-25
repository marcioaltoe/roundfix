---
spec: 0168-deliver-one-worktree-per-item
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# Deliver, one worktree per item

## Executive Summary

Give every `roundfix deliver` item its own linked worktree, created from the
refreshed default branch, and run every per-item action there. Park then
touches no checkout, resume finds or recreates the worktree, and merge removes
it. The restore machinery Spec 0161 added to park is deleted.

## Project Constraints

- Identifier strategy: applicable — the item worktree path derives from
  `worktree.location`, the repository and the per-delivery item branch. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — `gh` and Git carry their own
  credentials; Roundfix handles none. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the queue's checkout behaviour and
  `deliver status` output are public CLI surface the shipped skill documents.
  Express maintainer authorization: recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## The item record

`DeliveryQueueItem` in `internal/store/delivery.go` gains `Worktree`, persisted
in a new `worktree` column by a migration from the current schema version;
earlier databases open with an empty value. The item branch and the worktree
path are recorded in one write, before any Git command creates either, and the
first recorded pair is never overwritten. `StartingBranch` stops being read or
written.

## The item worktree

`CreateItemBranch` in `internal/cli/deliver_workflow.go` becomes a preparation
step that returns the item branch and worktree. It derives
`roundfix/deliver-<slug>-<suffix>` as today and a worktree path under
`worktree.location` through the path derivation in `internal/worktree`, records
both, refreshes the default branch, and runs
`git worktree add --no-track -b <branch> <path> <remote>/<default>` from the
user's repository. The worktree is provisioned like a Run Worktree: the
configured `worktree.copy` files are copied from the user's checkout and the
configured bootstrap runs, through helpers exported from `internal/worktree`.
No step inspects, requires clean or changes the user's checkout.

The `ItemWorkspace` interface in `internal/delivery/engine.go` returns the
worktree with the branch, and the engine passes that path as the working
directory to the runner, reviewer, archiver, gate, authorization reader,
publication planner and push. Queue and action records stay keyed by the user's
repository root. Archive paths resolve relative to the item worktree. Push in
`internal/delivery/github.go` runs with the item worktree as its working
directory. `deliver status` in `internal/cli/deliver.go` prints a fourth
tab-separated field, the item worktree or `-`.

## Park, resume and merge

Park persists the blocker and changes no checkout; the item worktree stays at
its recorded path. On resume, a recorded worktree that Git lists with the
recorded branch checked out is used as is. A recorded worktree that is missing
is recreated with `git worktree add <path> <branch>` from the recorded branch.
When the recorded branch is missing too, the item parks with the blocker
`item-worktree-missing` instead of replaying its stage. After the merged stage
is persisted, the item worktree is removed and the local item branch is
deleted; a failed removal leaves the item merged and is retried on the next
resume.

## Retired

- `StartingBranch` on the item and the starting-branch argument of
  `RecordDeliveryQueueItemBranch`; the column stays so earlier databases open.
- `ParkItem`'s reset, literal-pathspec clean and switch, its untracked-path
  tracking, and the clean-checkout checks in `CreateItemBranch` and
  `UseItemBranch`.
- The tests that pin that behaviour — `TestAParkLeavesACleanCheckout`,
  `TestParkKeepsUntrackedFilesTheItemDidNotCreate`,
  `TestParkNeverTouchesACheckoutTheItemRefused` and
  `TestParkRestoresTheRecordedStartingBranchAfterACrash` — are replaced by this
  Spec's tests of an untouched checkout and a kept item worktree.

These guarantees are preserved and keep their tests: an item branch has no
upstream, each delivery gets its own branch, a resume reuses the recorded
branch, and a real archive commit is accepted on resume while one with extra
changes is refused.

## API Contracts

1. `deliver start` and `deliver resume` never change the user's checkout and do
   not require it clean.
2. `deliver status` prints, per item, the Spec slug, stage, blocker and item
   worktree, tab-separated, with `-` for an empty value.
3. Item worktrees live under `worktree.location`; a parked item keeps its
   worktree, and a merged item's worktree and local item branch are removed.

## Coverage Map

- Goal 1 → The item worktree; API Contract 1.
- Goal 2 → Park, resume and merge; API Contracts 2-3.
- Core Features 1 and 2 → The item record, The item worktree.
- Core Features 3 and 4 → The item worktree, Park, resume and merge.
- Core Features 5 and 6 → Park, resume and merge.
- Core Feature 7 → Retired.
- Success Metrics 1 and 2 → Testing Approach 2.
- Success Metrics 3 and 4 → Testing Approach 3.
- API Contracts 1-3 → The item worktree, Park, resume and merge.

## Integration Points

- **Spec 0161.** Owns the branch, archive and owner behaviour this Spec keeps,
  and the park behaviour it retires.
- **Specs 0157 and 0162.** A Run started in an item worktree shares the
  repository identity, and stays listed after the worktree is removed.
- **`internal/worktree`.** Path derivation and provisioning are reused, not
  duplicated.

## Testing Approach

1. **Item record.** The branch and worktree are recorded together and the first
   pair is kept; a database at the previous schema migrates with an empty
   worktree.
2. **Isolation.** With real Git, a user checkout that is dirty, on a
   non-default branch, with an ignored file and an untracked nested repository,
   is byte-for-byte unchanged after one item parks and another merges; each
   item's actions receive its own worktree, created from the refreshed default
   branch; `deliver status` prints it.
3. **Lifecycle.** With real Git, a resume uses a recorded worktree, recreates a
   removed one from the recorded branch, and parks with `item-worktree-missing`
   when the branch is gone; a merged item leaves no worktree and no local item
   branch.
4. **Documentation.** The skill, its mirror and the user guide describe the
   worktree behaviour, and the skill check passes.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The item record (depends on: none).
2. The item worktree (depends on: 1).
3. Park, resume and merge; retire the restore machinery (depends on: 2).
4. The shipped skill and the user guide (depends on: 3).
5. Terminal QA (depends on: 1, 2, 3, 4, 6).
6. Provision every item worktree, and let migrated merged items rest (depends on: 4).

## Risks & Considerations

- **Worktree residue.** Parked items keep their worktrees by design; the skill
  and guide name them and `deliver status` prints them, so they are never
  hidden.
- **Provisioning cost.** Bootstrap now runs once per item worktree; that is the
  cost `roundfix implement` already pays per Run Worktree.
