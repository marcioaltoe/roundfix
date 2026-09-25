---
type: feat
status: open
created: 2026-09-25
spec: null
reason: null
---

# Release a delivered Spec's Runs, branches and worktrees after its squash merge

## Opportunity

After every squash merge the Spec's Run Worktrees, `roundfix/run-*` branches and carry-forward staging worktrees stay behind; on 2026-09-25 this repository held 16 run branches and 16 Run worktrees for Specs already merged. `roundfix deliver` records `merged` and stops, and nothing else releases them.

## Value

A repository that stays clean without manual `git worktree remove` against the skill's guidance; retained Runs stop blocking the Branch Integrity Preflight; storage stops growing with every delivery.

## Shape

A release step run by `roundfix deliver` right after the merge receipt, and an explicit `roundfix reconcile --release-merged <pr|spec>` for manual merges, releasing every Run of the merged Spec whose work is proven in the merge commit — including superseded and carried-forward Unresolved Runs — and sweeping stale carry-forward staging worktrees. Depends on B03, B05 and B16 so the proofs can pass.

Evidence: `docs/findings/2026-09-25-squash-merge-leaves-runs-and-branches-behind.md`.
