---
schema: spec-tasks/v1
spec: 0007-command-lifecycle
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
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_01, task_02, task_03, task_04, task_05]
---

# Tasks — Command Lifecycle

| id      | title                                                        | type    | complexity | needs                                        |
| ------- | ------------------------------------------------------------ | ------- | ---------- | -------------------------------------------- |
| task_01 | Graceful Stop Requests through the Run Database              | backend | high       | —                                            |
| task_02 | Setup Command: environment bootstrap                         | backend | high       | —                                            |
| task_03 | Upgrade Command and version freshness check                  | backend | high       | —                                            |
| task_04 | Force stop with cooperative Agent Session cancel             | backend | medium     | task_01                                      |
| task_05 | Spec-Run push at Clean via Project Config                    | backend | medium     | —                                            |
| task_06 | Docs and skill sync                                          | docs    | low        | task_01, task_02, task_03, task_04, task_05  |

Waves: 1 → task_01, task_02, task_03, task_05 · 2 → task_04 · 3 → task_06
