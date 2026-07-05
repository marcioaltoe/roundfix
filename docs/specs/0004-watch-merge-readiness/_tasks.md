---
schema: spec-tasks/v1
spec: 0004-watch-merge-readiness
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_02, task_03, task_04]
---

# Tasks — Watch Merge Readiness

| id      | title                                                          | type    | complexity | needs                     |
| ------- | -------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | Poll-first ordering and pre-settled quiet skip                 | backend | medium     | —                         |
| task_02 | Deterministic stdout reports for watch and resolve             | backend | medium     | —                         |
| task_03 | Agent-console suppression flag and filter sink                 | backend | low        | —                         |
| task_04 | Merge-readiness confirm phase: CheckFunc, gh adapter, outcomes | backend | high       | task_01                   |
| task_05 | Docs and skill sync                                            | docs    | low        | task_02, task_03, task_04 |

Waves: 1 → task_01, task_02, task_03 · 2 → task_04 · 3 → task_05
