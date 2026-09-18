---
schema: spec-tasks/v1
spec: 0146-a-gate-that-runs-the-analyzer
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

# Tasks — A gate that runs the analyzer

| id      | title                                        | type    | complexity | needs   |
| ------- | -------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Prove the analyzer can fail                   | test    | medium     | —       |
| task_02 | Compose the analyzer into the gate            | chore   | low        | task_01 |
| task_03 | Run the final QA gate                         | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03
