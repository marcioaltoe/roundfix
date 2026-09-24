---
spec: 0161-deliver-ready-for-real-repositories
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Deliver, ready for real repositories

## Executive Summary

Make `roundfix deliver` wait for unreported checks, create untracked,
per-delivery item branches recorded before creation, restore a clean checkout
on park, accept an archive commit after a crash, and release a stale owner.

## Project Constraints

- Identifier strategy: applicable — per-delivery branch suffix. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — `gh` and Git carry their own
  credentials; Roundfix handles none. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Checks

`CurrentHeadChecks` in `internal/delivery/github.go` maps exit 1 with empty
stdout and a `no checks reported` stderr to an empty report. The wait loop in
`internal/delivery/engine.go` treats an empty report as pending and retries a
read error until the bounded deadline, parking only on failure, cancellation or
timeout.

## Branches

`CreateItemBranch` in `internal/cli/deliver_workflow.go` derives
`roundfix/deliver-<slug>-<suffix>`, records it on the item through
`internal/store/delivery.go` before any Git command, and creates it with
`git switch --no-track -c`. When the recorded branch already exists, a resume
switches to it instead of creating it.

## Park and resume

On park, the workflow discards uncommitted changes the item made and returns
the checkout to the branch it started from. On resume at the archiving stage, a
HEAD whose parent is the reviewed head and whose diff is exactly the move of the
Spec folder into the archive root is accepted. `resume` in
`internal/cli/deliver.go` compares the live identity of the recorded PID with
the recorded identity and releases the owner on a mismatch.

## API Contracts

1. `roundfix deliver` waits for checks that have not reported.
2. Item branches are untracked and unique per delivery.
3. `deliver start` on a parked Spec, and `deliver resume` after a crash, succeed.

## Coverage Map

- Goal 1 → Checks, Branches; API Contracts 1-2.
- Goal 2 → Branches, Park and resume; API Contract 3.
- Core Feature 1 → Checks.
- Core Features 2 and 3 → Branches.
- Core Features 4, 5 and 6 → Park and resume.
- Success Metric 1 → Testing Approach 1.
- Success Metrics 2 and 3 → Testing Approach 2.
- Success Metrics 4 and 5 → Testing Approach 3.
- API Contracts 1-3 → Checks, Branches, Park and resume.

## Integration Points

- **Spec 0156.** Owns the stages this Spec corrects.

## Testing Approach

1. **Checks.** The exact gh output for no reported checks becomes an empty
   report; the engine waits and merges once checks pass; a read error inside
   the wait is retried.
2. **Branches.** With real Git, the item branch has no upstream, two deliveries
   of one slug get different branches, and a resume after a crash between
   creation and the stage write reuses the recorded branch.
3. **Park and resume.** With real Git, a non-exact archive park leaves a clean
   checkout and the next item starts; a crash after the archive commit resumes
   past archiving; a reused PID with a different identity is released.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Checks (depends on: none).
2. Branches (depends on: 1).
3. Park and resume (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **Discarding work on park.** Only changes the item itself made are discarded;
  commits on the item branch stay on it.
