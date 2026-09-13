---
schema: spec-tasks/v1
spec: 0133-a-fixture-that-does-not-spawn-per-test
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
      needs: [task_01, task_02]
---

# Tasks — The regressions this branch must not ship

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | State the two authorities apart | backend | high | — |
| task_02 | Stop paying Git setup per Task-cycle test | test | medium | task_01 |
| task_03 | Run the final QA gate | qa | medium | task_01, task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03.
