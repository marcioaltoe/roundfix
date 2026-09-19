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
    - id: task_03
      file: task_03.md
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
---

# Tasks — A profile that declares both tiers

| id      | title                                        | type    | complexity | needs   |
| ------- | -------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Declare the incremental tier in the Profile   | backend | medium     | —       |
| task_02 | Publish both tiers in the generated guidance  | backend | medium     | task_01 |
| task_03 | Regenerate the derived guides and pins        | chore   | low        | task_02 |
| task_04 | Run the final QA gate                         | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
