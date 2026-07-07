---
schema: spec-tasks/v1
spec: 0017-run-discovery
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
      needs: [task_02, task_03]
---

# Tasks — Run discovery

| id      | title                                                            | type    | complexity | needs            |
| ------- | ---------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Store listing query with repository scope and active filters     | backend | low        | —                |
| task_02 | runs list command with stable columns and empty-result contract  | backend | medium     | task_01          |
| task_03 | Attach picker: no-argument Interactive Input and no-input error  | backend | medium     | task_01          |
| task_04 | Docs and skill sync for the Run discovery surface                | docs    | low        | task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04
