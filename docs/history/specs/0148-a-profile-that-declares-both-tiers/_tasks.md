---
schema: spec-tasks/v1
spec: 0148-a-profile-that-declares-both-tiers
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02]
---

# Tasks — A profile that declares both tiers

| id      | title                                        | type    | complexity | needs   |
| ------- | -------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Declare the incremental tier in the Profile   | backend | medium     | —       |
| task_02 | Regenerate the derived catalog and plan       | chore   | low        | task_01 |
| task_04 | Run the final QA gate                         | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_04
