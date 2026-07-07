---
schema: spec-tasks/v1
spec: 0019-run-outcome-notifications
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — Run Outcome Notifications

| id      | title                                                          | type    | complexity | needs   |
| ------- | -------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | notify config section with enabled and command                 | backend | low        | —       |
| task_02 | notify package: payload, command and native notifiers          | backend | medium     | task_01 |
| task_03 | Terminal-outcome wiring with best-effort warning reporting     | backend | medium     | task_02 |
| task_04 | Docs and skill sync for Run Outcome Notifications              | docs    | low        | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
