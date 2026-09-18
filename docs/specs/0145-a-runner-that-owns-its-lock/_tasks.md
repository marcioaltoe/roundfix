---
schema: spec-tasks/v1
spec: 0145-a-runner-that-owns-its-lock
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — A runner that owns its lock

| id      | title                                   | type    | complexity | needs   |
| ------- | --------------------------------------- | ------- | ---------- | ------- |
| task_01 | Give the runner's state one owner        | backend | medium     | —       |
| task_02 | Run the final QA gate                    | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02
