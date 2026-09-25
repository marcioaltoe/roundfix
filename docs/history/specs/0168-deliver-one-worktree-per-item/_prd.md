---
spec: 0168-deliver-one-worktree-per-item
status: archived
created: 2026-09-24
surfaces: [backend, cli, docs]
archived: "2026-09-25"
source_slug: 0168-deliver-one-worktree-per-item
---


# Deliver, one worktree per item

Spec 0161 shipped with four recorded limits, all in park, and stated that no
release may ship `roundfix deliver` until this Spec merges. This Spec inherits
that gate: the command stays unreleased until it merges.

The four limits share one cause. `roundfix deliver` runs every queue item in the
user's own checkout: it switches branches there, commits there and, on park,
resets and cleans there to put the checkout back. Every attempt to make that
restore safe added state that guesses what the user had:

- **Park can delete ignored files.** When the item branch ignores less than the
  starting branch, a file the user's branch ignored is untracked on the item
  branch and park's clean removes it.
- **Park restores a stale branch.** After an item merges, the next item records
  the merged item branch as its starting branch, so its park returns the user to
  a delivery branch that is already done.
- **A deleted starting branch parks nothing.** Park fails after the reset and
  clean, so the item never parks and every resume replays its stage.
- **A nested repository survives park.** An untracked nested repository the item
  created stays behind, the checkout stays dirty and the item stays unparked.

## Project Constraints

- Identifier strategy: applicable — each item worktree lives at a path derived
  under `worktree.location` from the repository and the item's per-delivery
  branch, so two deliveries of one Spec never share a worktree. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — `gh` and Git carry their own
  credentials; Roundfix handles none. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — `roundfix deliver` stops changing the user's
  checkout and `deliver status` prints each item's worktree, which the shipped
  skill documents. Express maintainer authorization: the standing grant of
  2026-09-18, recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.
  Sanctioned regeneration: `make skills-sync`. The bounded set is the
  intersection of this Spec's changed paths with `GovernedPath`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- `roundfix deliver` never switches, resets or cleans the user's checkout.
- A parked item stays resumable and inspectable, and a merged item leaves
  nothing behind.

## Core Features

1. **One worktree per item.** Each item gets its own linked Git worktree,
   created from the refreshed default branch on the item's untracked,
   per-delivery branch, under `worktree.location`.
2. **Every per-item action runs there.** `roundfix implement`, pre-PR review,
   archive, the repository gate, the authorization read and push run with the
   item worktree as their working directory.
3. **The user's checkout is never touched.** The queue does not switch, reset
   or clean it, and does not need it clean to start an item.
4. **Park keeps the worktree.** Parking changes no checkout; the item worktree
   stays at the path recorded on the item, and `deliver status` prints it.
5. **Resume finds or recreates the worktree.** A resume uses the recorded
   worktree, recreates it from the recorded branch when it is gone, and parks
   the item when the branch is gone too, rather than replaying its stage.
6. **Merge cleans up.** After an item merges, its worktree and its local item
   branch are removed.
7. **The restore machinery is gone.** The starting branch, the untracked-file
   baseline and park's reset, clean and switch that Spec 0161 added are removed.

## Non-Goals / Out of Scope

- New stages, new external actions or new configuration keys.
- Running items in parallel; the queue stays sequential.
- Removing the worktrees of parked items automatically.
- Changing `roundfix implement`, `roundfix review` or `roundfix archive`.

## Success Metrics

1. With the user's checkout dirty, on a non-default branch, holding an ignored
   file and an untracked nested repository, a queue of one item that parks and
   one that merges leaves the checkout's branch, HEAD, status and those files
   unchanged.
2. Each item's implement, review, archive, gate and push run in that item's own
   worktree, created from the refreshed default branch.
3. A parked item's worktree remains at its recorded path; a merged item leaves
   no worktree and no local item branch.
4. A resume whose recorded worktree is gone recreates it from the recorded
   branch; with the branch gone too, the item parks with a named blocker.

## Decisions

- **A worktree per item, not a safer park.** Four park defects in two review
  rounds all came from restoring a checkout the loop does not own. Each patch —
  a recorded starting branch, an untracked baseline, literal pathspecs — added
  state to guess what the user had, and ignored files, a moved or deleted
  branch and nested repositories remain state no reset or clean can classify
  safely. A worktree the loop creates is entirely its own, so park needs no
  restore and the user's checkout needs no protection. It also lets the user
  keep working while the queue runs, and matches how `roundfix implement`
  already isolates each Run in a Run Worktree.
- **Park keeps, merge removes.** A parked worktree holds exactly what a human
  needs to inspect or finish the item; after merge the work is on the default
  branch and the worktree is residue.
- **Recreate, then park.** The branch is the durable record and the worktree is
  disposable, so a missing worktree is rebuilt from the branch; only a missing
  branch, which holds the work, parks the item.
- **Reuse `worktree.location`.** It is already validated to lie outside the
  repository tree, so no new setting is needed.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph, using real Git where the defect lives in Git behaviour. The negative
cases carry the weight: a user checkout whose ignored file, nested repository,
branch or uncommitted change is altered by a delivery passes every fake-backed
happy path.

## Research basis

The four limits and their reproductions are recorded in the archived PRD of
Spec 0161, found by its second pre-PR review on 2026-09-24. This Spec relies on
the repository-identity work: Spec 0157 made a Run started from a linked
worktree share the repository's identity and artifact directory, and Spec 0162
stores a durable repository key on each Run so it stays listed after its
worktree is removed.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
