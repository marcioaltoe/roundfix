---
schema: spec-tasks/v1
spec: 0014-run-store-retention
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — Run Store Retention

| id      | title                                                       | type    | complexity | needs                     |
| ------- | ----------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Journal retention config and terminal-only prune in the store| backend | medium     | —                         |
| task_02 | GC Command with dry-run and artifact cleanup                | backend | medium     | task_01                   |
| task_03 | Best-effort retention prune in the preflight sweep          | backend | low        | task_01                   |
| task_04 | Docs and skill sync                                         | docs    | low        | task_01, task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04
