---
spec: 0161-deliver-ready-for-real-repositories
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Deliver, ready for real repositories

Spec 0156 shipped `roundfix deliver` with six recorded limits after its
corrective ceiling was spent. It also stated that no release may ship the
command before this Spec merges. Four of the limits make the loop park or
publish wrongly on an ordinary repository; two make it unrecoverable after a
crash or a reboot.

**Checks that have not reported yet park the item.** Right after
`gh pr create`, `gh pr checks` exits 1 with `no checks reported`, and the item
parks as `delivery-error` instead of waiting.

**The item branch tracks the default branch.** The branch is created with
`git switch -c <branch> origin/<default>`, so its upstream is `origin/main`.
With `implement.auto_push: true` a Clean Run pushes straight to the default
branch before review.

**A parked Spec cannot be delivered again.** The branch name is fixed per slug
and created with `git switch -c`, so a second delivery, or a resume after a
crash between branch creation and the stage write, parks with "branch already
exists".

**A non-exact archive leaves the checkout dirty,** and every later item then
refuses to start.

**A crash after the archive commit resumes as `review-stale`,** because HEAD's
parent is the reviewed head.

**A reused owner PID locks the queue:** `resume` checks only that the PID is
alive, and `stop` cannot prove the recorded identity.

## Project Constraints

- Identifier strategy: applicable — an item branch name gains a per-delivery
  suffix so a new delivery never collides with an earlier one. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — `gh` and Git carry their own
  credentials; Roundfix handles none. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A delivery waits for checks and never publishes to the default branch.
- A parked or crashed delivery can always be resumed or delivered again.

## Core Features

1. **Wait for checks that have not reported.** A `no checks reported` answer is
   an empty report, awaited within the bounded timeout; a transient read error
   inside the wait is retried until the deadline.
2. **An item branch that tracks nothing.** The item branch is created without an
   upstream, so no command can push it to the default branch implicitly.
3. **Re-delivery and crash-safe branch creation.** The item branch name carries
   a per-delivery suffix and is recorded on the item before it is created; a
   resume reuses a recorded branch that exists.
4. **A park leaves a clean checkout.** Parking restores the checkout the item
   started from, so the next item can start.
5. **An archive commit survives a crash.** On resume, a HEAD whose parent is the
   reviewed head and whose change is exactly the Spec move is accepted as the
   archive result.
6. **A stale owner is released.** `resume` compares the live process identity
   with the recorded one and releases an owner whose identity does not match.

## Non-Goals / Out of Scope

- New stages or new external actions.
- Changing what the authority check requires.

## Success Metrics

1. An item whose pull request has no reported checks waits and merges once they
   pass.
2. An item branch has no upstream.
3. A parked Spec delivers again, and a crash between branch creation and the
   stage write resumes.
4. A park after a non-exact archive leaves the next item able to start.
5. A crash after the archive commit resumes without `review-stale`, and a reused
   owner PID does not lock the queue.

## Recorded limits

The corrective ceiling of two Tasks was spent on the defects the first pre-PR
review of 2026-09-24 found. The second review found that park still mutates the
user's checkout unsafely. The cause is structural — the loop runs every item in
the user's own checkout — so the fix is carried to Spec 0168, which gives each
item its own worktree. No release may ship `roundfix deliver` until Spec 0168
merges.

- Park can delete files that were ignored when the item started, when the item
  branch ignores less than the starting branch. Reproduction: commit a local
  `.gitignore` entry for `local.env` on `main` that `origin/main` lacks, create
  `local.env`, deliver one Spec and let it park.
- After an item merges, the next item records the merged item branch as its
  starting branch, so its park returns the user to a stale delivery branch.
- A starting branch deleted after it was recorded makes park fail after the
  reset and clean, so the item never parks and every resume replays its stage.
- An untracked nested repository created by the item survives park, leaving the
  checkout dirty and the item unparked.

## Decisions

- **Fix in place, not behind a flag.** Every limit is a defect of the shipped
  behaviour, so the command changes rather than growing an option.
- **Suffix over reuse.** A fresh branch per delivery avoids force-pushing over a
  branch an earlier delivery published.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph, using real Git where the defect lives in Git behaviour. The
negative cases carry the weight: a branch with an upstream, or a park that
leaves the checkout dirty, passes a fake-backed happy path.

## Research basis

The six limits and their reproductions are recorded in the archived PRD of Spec
0156, found by two independent pre-PR review rounds on 2026-09-24. The
`no checks reported` behaviour was confirmed in the installed gh 2.101.0.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
