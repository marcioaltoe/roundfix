---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-26
kind: rollup
members:
  - 2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md
  - 2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md
  - 2026-07-16-vortex-pr87-detached-watch-notification.md
  - 2026-07-17-global-run-storage-sanitation-and-compaction.md
  - 2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md
  - 2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md
  - 2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md
  - 2026-07-30-run-termination-does-not-reach-the-acpx-child.md
  - 2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md
  - 2026-08-04-branch-integrity-preflight-prescribes-a-remedy-that-reintroduces-superseded-work.md
  - 2026-08-04-watch-derives-a-review-head-it-never-checks-is-reachable.md
  - 2026-08-05-preflight-prescribes-integrating-a-superseded-run-branch.md
  - 2026-08-06-six-parallel-runs-on-one-machine-show-the-seams.md
---

# Run lifecycle and branch integrity — every created resource needs one terminal disposition (2026-08-06)

The Run findings show one lifecycle spread across process trees, Run Branches,
Task and Run Worktrees, refs, artifacts, notifications, and database storage.
Failures recur when one surface decides terminal state without disposing or
classifying the resources created on another.

## Consolidated learning

- Owner identity and termination must cover the real process tree without
  depending on a new process at the moment the host is already constrained.
- Failed and superseded Run Branches need explicit classifications;
  preflight must not prescribe integrating work that reconciliation can prove
  superseded.
- Branch and review-head checks must prove reachability against the intended
  Pull Request, not infer it from names or the global refs namespace.
- Reconciliation, GC, and notifications must report actionable terminal state
  across repositories while preserving evidence needed to audit cleanup.

## Live edge

Specs 0039, 0055, 0059, 0066, and 0068 closed important parts of the lifecycle.
The rollup remains `pending` because parallel Runs still expose cross-surface
ownership and classification seams that no single terminal audit covers.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
