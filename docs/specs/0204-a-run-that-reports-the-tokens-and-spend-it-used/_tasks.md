---
schema: spec-tasks/v1
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A prompt's result carries the usage its adapter reported |
| task_02 | data | Every prompt of a Run leaves a usage record and a `usage` event |
| task_03 | backend | `runs show`, the Implement Run summary and `deliver status` print the usage |
| task_04 | backend | A Delivery Queue stops starting items at its token ceiling |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
