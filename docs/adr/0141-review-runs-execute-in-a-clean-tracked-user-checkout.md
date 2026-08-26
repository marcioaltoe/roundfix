---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Review Runs execute in a clean tracked user checkout

Review Runs (fetch, resolve, watch) create no Run Worktree: they operate directly on the checked-out PR Head Branch — a review fix is by definition a delta over the pull request's published HEAD, and the worktree indirection produced Final Pushes that silently omitted stranded work. The guard is a deterministic Branch Integrity Preflight requiring zero unintegrated Run Branch commits, no other Run bound to the branch, and a clean tracked working tree (each dirty path named; untracked files stay allowed because batch commits stage only paths changed since the batch snapshot). Consequence: after a failed batch, everything dirty in the checkout is Agent work by construction. ADR-0023's worktree rule keeps applying to Spec Runs, where Task concurrency needs isolation.

Consolidates ADR-0042 and ADR-0045 (2026-08-26); both are archived under docs/history/adr/.
