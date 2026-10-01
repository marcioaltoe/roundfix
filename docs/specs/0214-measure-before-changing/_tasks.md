---
schema: spec-tasks/v1
spec: 0214-measure-before-changing
qa: task_05
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | `roundfix runs causes` says why Verification failed and why corrective Tasks were added |
| task_02 | test | A flag-gated harness re-scores the Task acceptance question against the attempt count |
| task_03 | docs | Measure who reads archived evidence and propose, removing nothing |
| task_04 | docs | Record why corrective Tasks happen and what the Task acceptance question is worth |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04 · 4 → task_05
