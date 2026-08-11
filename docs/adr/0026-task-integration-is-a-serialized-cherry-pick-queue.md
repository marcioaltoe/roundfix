---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Task integration is a serialized cherry-pick queue

Concurrent Tasks branch from the same Run Branch tip, so only the first finisher can fast-forward; the rest integrate through a single serialized queue that cherry-picks each settled Task's commit onto the Run Branch in completion order. A cherry-pick conflict settles that Task `failed` with its Task Worktree kept for inspection — never an in-place merge resolution — because Tasks that write-tasks declared independent are file-disjoint by contract, making a conflict a real graph defect to surface, not noise to auto-resolve. The Run Branch → user branch integration keeps the ADR-0024 porcelain-only protocol unchanged.
