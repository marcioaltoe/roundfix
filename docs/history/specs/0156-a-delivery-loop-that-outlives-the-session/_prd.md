---
spec: 0156-a-delivery-loop-that-outlives-the-session
status: archived
created: 2026-09-24
surfaces: [backend, cli, docs]
archived: "2026-09-24"
source_slug: 0156-a-delivery-loop-that-outlives-the-session
---


# A delivery loop that outlives the session

The Daemon survives a closed terminal. The delivery does not.

`roundfix implement --detach` runs a Spec's Task Graph as a detached process, so
a Run finishes whether anyone is watching. Everything after the Run — the
pre-PR review, archiving, the pull request, waiting for checks, the merge, and
starting the next Spec — is done by whoever drives the session. When the
session stops, the chain stops.

Measured on Spec 0154 on 2026-09-24: about 25 hours elapsed, 2 hours 15 minutes
of actual Run time, so roughly 91 percent of the elapsed time was idle. The
longest gaps were 7 hours between two Runs and about 15 hours overnight, when the
session closed and took its background jobs with it. The machine does not sleep
and the Runs survive; the orchestration is what lives in the chat.

There is no code in Roundfix today that opens, follows or merges a pull request.
Those steps happen only because an operator types them.

The maintainer restructured the queue on 2026-09-24. Jev ranked this delivery
first with probability 1.0; the original portfolio had scheduled it last, behind
eight prerequisite Specs.

## Project Constraints

- Identifier strategy: applicable — each queue item is keyed by its Spec slug and
  records the Run identity, the candidate commits, the pull request number and
  the merge commit it produced, so a resumed queue addresses the same external
  objects. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — pull request actions use the GitHub CLI
  already authenticated on the host; this Spec stores no credential and adds no
  token handling. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0153 makes pre-PR review an explicit provider
  policy, ADR-0155 makes the `qa` Task declare the matrix and ADR-0156 makes a
  declared promise name a consuming Task. This Spec's gate is bound by ADR-0080,
  ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — a new command changes the public CLI surface
  the shipped skill documents. Express maintainer authorization: the standing
  grants of 2026-09-18 and 2026-09-21, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/cli/cli_test.go`. Sanctioned regeneration: `make skills-sync`. The
  bounded set is the intersection of this Spec's changed paths with
  `internal/speccheck/governed.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- A queue of Specs advances from Run to merge without an operator in the loop.
- Losing the session loses nothing: the queue resumes from what actually
  happened.
- Nothing is published without the evidence the workflow requires, and nothing
  is published twice.

## Core Features

1. **A durable queue.** `roundfix deliver start <slug>...` records an ordered
   queue and starts a detached owner, like the Daemon's. `deliver status` shows
   each item's stage and any blocker; `deliver stop` ends the owner; `deliver
   resume` restarts it.
2. **The delivery order the maintainer set.** Each item goes: Spec Run to Clean,
   pre-PR review, archive on the branch, repository gate, push, pull request,
   current-head checks, squash merge, next item. With pre-PR review set to
   `none`, archive follows the QA gate directly.
3. **Review binds the published head.** After the review, the only commit
   allowed before publication is the archive commit, and it must contain exactly
   the move of the Spec's folder. Any other change invalidates the review.
4. **Intent before action, receipt after.** Before a push, a pull request or a
   merge, the queue records the intent; after it, the receipt. On resume, an
   intent without a receipt is reconciled against the observed state — the
   branch on the remote, an open pull request for the head, a merged pull
   request — before anything is retried. The queue never opens a second pull
   request for the same head and never merges twice.
5. **Authority per Spec.** Publication proceeds only when the Spec's
   authorization record grants `push`, `pull_request` and `merge`.
6. **A blocker parks, it does not improvise.** A Run that ends Unresolved, a
   review with findings or a blocked review, a failing gate or failing checks
   parks the item with its reason, and the queue continues with the next item.

## Non-Goals / Out of Scope

- Inventorying the backlog to build a queue (Spec 0127 Core Feature 1); the
  maintainer supplies the ordered list.
- Modelling dependencies between queued Specs or revalidating later Specs
  (Spec 0127 Core Feature 6).
- Queue-wide time or spending limits (Spec 0127 Core Feature 5); each Run keeps
  its existing budget.
- Delegation from the `implement-spec` skill (Spec 0127 Core Feature 8).
- Releases.
- Resolving review threads on the open pull request, which `roundfix watch`
  already owns.

## Success Metrics

1. A fixture queue of two Specs, with the GitHub boundary faked, reaches merged
   for both with no operator action between them.
2. Killing the owner after a recorded push intent and restarting it yields
   exactly one pull request and one merge.
3. A commit other than the Spec's folder move between review and publication
   parks the item instead of publishing it.
4. A Spec whose authorization record lacks `merge` is parked before any merge
   is attempted.

## Recorded limits

The maintainer's corrective ceiling of two Tasks was spent on the migration
ladder and on the publication-safety defects the pre-PR review of 2026-09-24
found. Two lower-severity defects from that review are carried:

- A crash after the archive commit but before the stage advances resumes as
  `review-stale`, because the new HEAD's parent is the reviewed head.
  Reproduction: kill the owner between the archive commit and the gating stage,
  then `deliver resume`. Carried fix: accept a HEAD whose parent is the reviewed
  head and whose change is exactly the Spec move.
- `resume` checks only that the recorded owner PID is alive, so a reused PID
  after a reboot makes both `resume` and `stop` refuse. Reproduction: record an
  owner, reuse its PID with an unrelated process, then `deliver resume`. Carried
  fix: compare the recorded process identity and release a stale owner.

The second pre-PR review round of 2026-09-24 found four more defects, carried
to Spec 0161 as a blocking follow-up. No release may ship `roundfix deliver`
until Spec 0161 is merged.

- Checks not yet reported: `gh pr checks` exits 1 with `no checks reported`
  right after the pull request is created, which parks the item as
  `delivery-error` instead of waiting. Reproduction: deliver one Spec to a
  repository whose CI has not attached a check yet.
- Item branch upstream: the item branch is created tracking
  `origin/<default>`, so with `implement.auto_push: true` a Clean Run pushes to
  the default branch before review. Reproduction: set `implement.auto_push:
  true` and deliver one Spec; `git rev-parse @{u}` on the item branch prints
  `origin/main`.
- Re-delivery: the item branch name is fixed per slug and created with
  `git switch -c`, so delivering a parked Spec again, or resuming after a crash
  between branch creation and the stage write, parks with "branch already
  exists".
- Dirty checkout after a non-exact archive: the archive changes stay
  uncommitted, and every later item refuses to start on a dirty checkout.

## Decisions

- **Review before archive.** The maintainer set the order on 2026-09-22: archive
  at the end of the pre-PR review, or after the QA gate when review is
  disabled. It supersedes the archive-then-review order that Spec 0127 Core
  Features 3 and 9 and Spec 0126 Core Feature 8 proposed.
- **Park and continue.** One blocked Spec stops only itself. Stopping the queue
  on the first blocker would recreate the idle gaps this Spec exists to remove.
- **Reuse, do not duplicate.** The Spec Run is the existing Implement executor,
  the review is `roundfix review`, archive is `roundfix archive`; this Spec adds
  the chain between them and the pull request boundary that does not exist.
- **The GitHub CLI as the boundary.** It is already authenticated on the host
  and already what the maintainer uses; a wrapper behind an interface keeps the
  queue testable with a fake.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a queue that retried a merge
after a lost receipt, or published a head the review never saw, would move fast
and leave the repository in a state nobody approved.

## Research basis

Measured on this repository: 91 percent idle on Spec 0154; no code path creates,
follows or merges a pull request; per-Run duration is already enforced through
`MaxRunDuration`. External research on 2026-09-24 found durable agent
orchestration converging on persisted state and resume from observed results
rather than always-on processes. The restructuring and its scoring are recorded
in the secondbrain at `inbox/roundfix/2026-09-24-reestruturacao-da-fila-com-jev.md`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
