---
schema: spec-tasks/v1
spec: 0022-cleanup-robustness
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — Cleanup Robustness

| id      | title                                                       | type    | complexity | needs   |
| ------- | ----------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Forced worktree removal and warn-and-continue clean path    | backend | medium     | —       |
| task_02 | Docs and skill rules: settle recovery and authoring gates   | docs    | low        | task_01 |

Waves: 1 → task_01 · 2 → task_02
