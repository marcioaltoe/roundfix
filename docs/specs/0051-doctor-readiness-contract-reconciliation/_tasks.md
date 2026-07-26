---
schema: spec-tasks/v1
spec: 0051-doctor-readiness-contract-reconciliation
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
      needs: [task_04]
---

# Tasks — Doctor Readiness Contract Reconciliation

| id      | title                                            | type    | complexity | needs   |
| ------- | ------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Tidy authorized Go module metadata               | chore   | low        | —       |
| task_02 | Make external skill hash ordering total          | backend | medium     | task_01 |
| task_03 | Make Repository Skill Set inspection cancellable | backend | high       | task_02 |
| task_04 | Reconcile Doctor readiness contracts             | backend | medium     | task_03 |
| task_05 | Synchronize Roundfix Doctor guidance             | docs    | low        | task_04 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05
