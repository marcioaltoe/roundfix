---
schema: spec-tasks/v1
spec: 0144-a-run-stops-when-its-budget-is-spent
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
    - id: task_06
      file: task_06.md
      needs: [task_04]
    - id: task_05
      file: task_05.md
      needs: [task_06]
---

# Tasks — A Run stops when its budget is spent

| id      | title                                          | type    | complexity | needs   |
| ------- | ---------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Bound the Run by its configured maximum         | backend | medium     | —       |
| task_02 | Settle the bounded Run as BudgetExceeded        | backend | medium     | task_01 |
| task_03 | Accept the outcome for carry-forward            | backend | medium     | task_02 |
| task_04 | Tell the skill what a spent budget does         | docs    | low        | task_03 |
| task_06 | Bound the whole Run, and say so consistently     | backend | medium     | task_04 |
| task_05 | Run the final QA gate                           | qa      | high       | task_06 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_06 · 6 → task_05
