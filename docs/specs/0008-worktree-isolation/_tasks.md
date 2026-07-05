---
schema: spec-tasks/v1
spec: 0008-worktree-isolation
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_03, task_04, task_05]
---

# Tasks — Worktree Isolation

| id      | title                                                            | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Worktree package: lifecycle and porcelain integration            | backend | high       | —                         |
| task_02 | Schema v6 work_dir, worktree.copy config, artifact default move  | data    | medium     | —                         |
| task_03 | Implement path on the Run Worktree                               | backend | high       | task_01, task_02          |
| task_04 | Resolve and watch on the Run Worktree                            | backend | high       | task_01, task_02          |
| task_05 | Live Run View and Attach read the execution workspace            | frontend| medium     | task_03                   |
| task_06 | Docs and skill sync                                              | docs    | low        | task_03, task_04, task_05 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05 · 4 → task_06
