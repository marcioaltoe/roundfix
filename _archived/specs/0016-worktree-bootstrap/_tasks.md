---
schema: spec-tasks/v1
spec: 0016-worktree-bootstrap
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
---

# Tasks — Worktree Bootstrap

| id      | title                                                            | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Bootstrap config and Run Worktree bootstrap with failure mapping | backend | medium     | —                |
| task_02 | Task Worktree bootstrap for concurrent Runs                      | backend | medium     | task_01          |
| task_03 | Docs and skill sync (bootstrap + env-file recipe)               | docs    | low        | task_01, task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03
