---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
kind: rollup
absorbed_by: 0125-repository-identity-and-run-branch-policy
members:
  - 2026-09-08-run-branches-ignore-the-selected-prefix.md
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

The Run findings show one lifecycle spread across process trees, Run Branches, Task and Run Worktrees, refs, artifacts, notifications, and database storage. Failures recur when one surface decides terminal state without disposing or classifying the resources created on another.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-rollup-run-lifecycle-and-branch-integrity.md`.
