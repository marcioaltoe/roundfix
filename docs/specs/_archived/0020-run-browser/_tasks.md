---
schema: spec-tasks/v1
spec: 0020-run-browser
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
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Run Browser

| id      | title                                                           | type     | complexity | needs            |
| ------- | --------------------------------------------------------------- | -------- | ---------- | ---------------- |
| task_01 | Run state filter and limit on the store listing query           | backend  | low        | —                |
| task_02 | runs list enrichment: columns, state and limit flags, notes     | backend  | medium     | task_01          |
| task_03 | RunBrowser TUI model with row formatting                        | frontend | high       | task_01          |
| task_04 | Browser entry points and the attach loop                        | backend  | medium     | task_02, task_03 |
| task_05 | Docs and skill sync for the Run Browser surface                 | docs     | low        | task_04          |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04 · 4 → task_05
