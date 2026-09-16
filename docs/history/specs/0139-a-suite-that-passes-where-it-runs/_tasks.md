---
schema: spec-tasks/v1
spec: 0139-a-suite-that-passes-where-it-runs
qa: task_05
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_03, task_04]
---

# Tasks — A suite that passes where it runs

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Stop ACPX fixtures churning links to the test binary | test | medium | — |
| task_02 | Bind Task-cycle waits to the test deadline | test | medium | — |
| task_03 | Prove the historical audit with a controlled fixture | test | medium | — |
| task_04 | Run the repository Verification in the QA gate step | backend | medium | task_02 |
| task_05 | Run the final QA gate | qa | medium | task_01, task_03, task_04 |

Waves: 1 → task_01, task_02, task_03 · 2 → task_04 · 3 → task_05.
