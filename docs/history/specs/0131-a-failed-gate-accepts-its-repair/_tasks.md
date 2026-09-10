---
schema: spec-tasks/v1
spec: 0131-a-failed-gate-accepts-its-repair
qa: task_03
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
---

# Tasks — A failed gate accepts its repair

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Characterize the loader's answer for each gate verdict | test | low | — |
| task_02 | Accept a failed gate above incomplete dependencies | backend | low | task_01 |
| task_03 | Run the final QA gate | qa | medium | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03.
