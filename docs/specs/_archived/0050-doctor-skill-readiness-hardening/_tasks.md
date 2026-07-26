---
schema: spec-tasks/v1
spec: 0050-doctor-skill-readiness-hardening
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
---

# Tasks — Doctor Skill Readiness Hardening

| id      | title                                            | type    | complexity | needs   |
| ------- | ------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | Add the authorized Unicode collation dependency | chore   | low        | —       |
| task_02 | Centralize external skill hash compatibility    | backend | medium     | task_01 |
| task_03 | Anchor Repository Skill Set filesystem reads    | backend | high       | task_02 |
| task_04 | Harden Doctor coordination and evidence         | backend | medium     | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

