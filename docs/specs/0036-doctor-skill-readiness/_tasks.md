---
schema: spec-tasks/v1
spec: 0036-doctor-skill-readiness
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — Doctor Skill Readiness

| id      | title                                            | type    | complexity | needs   |
| ------- | ------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Diagnose Repository Skill Set readiness          | backend | high       | —       |
| task_02 | Synchronize Doctor skill-readiness guidance      | docs    | medium     | task_01 |

Waves: 1 → task_01 · 2 → task_02
