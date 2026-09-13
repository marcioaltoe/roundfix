---
schema: spec-tasks/v1
spec: 0133-a-fixture-that-does-not-spawn-per-test
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

# Tasks — A fixture that does not spawn per test

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Stop paying Git setup per Task-cycle test | test | medium | — |
| task_02 | Run the final QA gate | qa | medium | task_01 |

Waves: 1 → task_01 · 2 → task_02.
