---
spec: 0156-a-delivery-loop-that-outlives-the-session
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# A delivery loop that outlives the session

## Executive Summary

A detached queue owner advances each queued Spec through the existing Implement
executor, `roundfix review`, `roundfix archive`, a repository gate and a new
pull request boundary, recording an intent before every external action and a
receipt after it, and reconciling any unmatched intent against observed state on
resume.

## Project Constraints

- Identifier strategy: applicable — items are keyed by Spec slug and record the
  Run, candidate commits, pull request number and merge commit. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the GitHub CLI already authenticated on
  the host carries every pull request action; no credential is stored. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0153, ADR-0155 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — a new command family is public CLI surface.
  Express maintainer authorization: recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/cli/cli_test.go`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Components

| Component | Location | Responsibility |
| --- | --- | --- |
| Queue store | `internal/store` | Persist the queue, ordered items, each item's stage and blocker, and action intents with their receipts. |
| Pull request boundary | `internal/delivery` | Push a branch; find an open pull request for a head or create one; read current-head checks; merge or observe an existing merge. Behind an interface, backed by the GitHub CLI, faked in tests. |
| Delivery engine | `internal/delivery` | The per-item stage machine, the review-binds-head rule, parking, and reconciliation on resume. |
| Command family | `internal/cli/deliver.go` | `deliver start`, `status`, `resume`, `stop`; starts the owner detached with the same mechanism `implement --detach` uses. |

## The stages

`queued` → `running` (Implement executor) → `reviewing` (`roundfix review`, or a
recorded omission under `none`) → `archiving` (`roundfix archive` on the branch)
→ `gating` (the repository gate on the archived head) → `publishing` (push, then
find-or-create the pull request) → `checking` (current-head checks) → `merging`
(squash merge) → `merged`.

Any stage may end in `parked`, with the reason. The owner then moves to the next
queued item.

## Review binds the head

The review record from Spec 0153 names the head it examined. Between that head
and publication the engine permits exactly one commit: the archive commit, whose
diff must be precisely the rename of `docs/specs/<slug>/` to the archive root.
Any other difference parks the item as `review-stale`.

## Intents and receipts

Before `push`, `create-pull-request` and `merge`, the engine writes an intent
row; after the action it writes the receipt. On resume, for each intent without
a receipt:

- push: is the branch on the remote at the recorded head?
- create-pull-request: is there an open pull request for the head branch?
- merge: is the pull request merged, and at which commit?

A positive observation becomes the receipt; a negative one allows the retry.
This is what makes a lost acknowledgement harmless.

## Authority

Before `publishing`, the engine reads the Spec's authorization record and
requires `push`, `pull_request` and `merge` among its operations. A missing
operation parks the item as `unauthorized` before any external action.

## API Contracts

1. `deliver start <slug>...` records the queue and starts a detached owner;
   `deliver status` prints each item's stage and blocker; `deliver stop` ends the
   owner; `deliver resume` restarts it from the persisted queue.
2. A resumed owner never creates a second pull request for a head and never
   merges a merged pull request.
3. A parked item does not stop the items after it.

## Coverage Map

- Goal 1 → The stages; Components.
- Goal 2 → Intents and receipts; API Contract 2.
- Goal 3 → Review binds the head; Authority; API Contract 2.
- Core Feature 1 → Components; API Contract 1.
- Core Feature 2 → The stages.
- Core Feature 3 → Review binds the head.
- Core Feature 4 → Intents and receipts; API Contract 2.
- Core Feature 5 → Authority.
- Core Feature 6 → The stages; API Contract 3.
- Success Metric 1 → Testing Approach 3.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 5.
- Success Metric 4 → Testing Approach 6.
- API Contracts 1-3 → Components, Intents and receipts, The stages.

## Integration Points

- **Implement executor.** Invoked as-is for each item's Run.
- **`roundfix review` (Spec 0153).** Its record decides the review stage.
- **`roundfix archive`.** Invoked on the branch after a clean review.
- **Repository gate.** The command in `defaults.verification`.

## Testing Approach

1. **Store.** Queue, items, stages, blockers, intents and receipts round-trip.
2. **Boundary.** The fake pull request boundary finds an existing pull request
   instead of creating one, and reports an existing merge instead of merging.
3. **End to end.** A two-Spec queue with fake Runs, a fake review and the fake
   boundary reaches `merged` for both with no operator action.
4. **Crash safety.** The owner is stopped after a recorded push intent and after
   a recorded merge intent; each resume yields exactly one pull request and one
   merge.
5. **Stale review.** A commit other than the folder move between review and
   publication parks the item as `review-stale`.
6. **Authority.** A Spec without `merge` in its authorization parks as
   `unauthorized` with no external action attempted.
7. **The skill is true.** The shipped skill and its mirror describe the command
   family, and the mirror matches the canonical file.
8. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The queue store (depends on: none).
2. The pull request boundary (depends on: none).
3. The delivery engine (depends on: 1, 2).
4. The command family (depends on: 3).
5. The shipped skill and the user guide (depends on: 4).
6. Terminal QA (depends on: 1, 2, 3, 4, 5, 7).
7. A migration ladder that applies from every earlier version (depends on: 5).

## Risks & Considerations

- **Publishing what nobody reviewed.** The review-binds-head rule is the
  control; a lenient comparison would let an unreviewed commit ride the archive
  commit into `main`.
- **Double effects after a crash.** Intent-before-action with observation on
  resume is the control; the crash-safety tests exercise both push and merge.
- **A queue that looks busy while stuck.** Parking with a named reason and
  `deliver status` make every stop visible; nothing waits silently.
