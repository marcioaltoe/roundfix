# Runs execute in per-Run git worktrees

**Partially superseded by ADR-0042 (2026-07-15): review Runs (fetch, resolve, watch) now execute in the user's checkout; this decision keeps governing spec Runs (implement).**

Every operational Run (resolve, watch, implement) executes in its own git worktree — created at Run start on a named Run Branch (`roundfix/run-<id>`) from the Run's head commit, living under Roundfix Home, recorded on the Run row — instead of the user's checkout. One worktree serves the whole Run sequentially, preserving the failed-work and before-snapshot semantics of ADRs 0010/0013/0014 unchanged inside it; a Clean, integrated Run removes its worktree and branch, while every other outcome keeps them as the single inspection and settle surface. Worktree-per-Task was rejected for now: it changes failed-work semantics, multiplies orphan states, and buys nothing before parallel execution exists.
