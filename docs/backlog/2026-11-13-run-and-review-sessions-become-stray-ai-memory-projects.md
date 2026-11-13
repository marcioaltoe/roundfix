---
type: fix
status: open
created: 2026-11-13
spec: null
reason: null
---

# Run and review sessions become stray ai-memory projects

## Symptom

Observed by the operator on 2026-11-13 on the mini, with Roundfix 0.142.2. Between 02:35 and 03:55 BRT, while the Delivery Queue delivered Spec 0386, `ai-memory list-projects` gained two projects with no pages:

- `roundfix-deliver-0386-a-message-that-asks-for-a-completion-notice-receives-one-69c397275560e6ae-15487c51-55161c1d`, the basename of a directory under `/Volumes/Brain-SSD/cache/.roundfix/worktrees/` that holds only a Run Worktree;
- `run_20261113T045011Z_cf0fd882ba9fb7d9`, the basename of that Run Worktree, whose root holds this repository's committed `.ai-memory.toml` (`project = "marcioaltoe-roundfix"`).

Older projects of the same shape were already listed: `run_20261112T164648Z_add31f5e8ad8c351.task_01`, `run_20261112T230635Z_aed4833ad749b5f1` and `roundfix-deliver-0085-…-fdf98040-c384e759` (a Quaestor item). The pre-PR review of that delivery recorded `validator: not-needed`, so the review validator that Spec 0385 sealed with the reviewed checkout's `.ai-memory.toml` did not run; these projects come from other Agent sessions.

## Where

The Agent sessions of a Run Worktree or Task Worktree and the directory that holds them, and how ai-memory resolves their project. The cause is not established: the Run Worktree carries the repository marker, so either a session starts in the parent directory, or ai-memory resolves the project from the basename before it reads the marker.

## Expected

Every Agent session a Run or review starts in this repository's worktrees records under `marcioaltoe-roundfix`, and no Delivery Queue item adds a project to `ai-memory list-projects`.

## Evidence

`ai-memory list-projects` before the 0386 delivery (26 projects, snapshot of 2026-11-12) and after it (28); project ids `01a1d35b-27d3-75b0-9be6-447ea01653bf` and `01a1d343-1b58-7e72-8bee-ef248195a0c2`. Not reproduced in isolation. If the cause lies in ai-memory, the outcome is a Cross-project Message to its owner, not a Roundfix change.
