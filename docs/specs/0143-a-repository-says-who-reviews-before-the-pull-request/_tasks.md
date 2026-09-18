---
schema: spec-tasks/v1
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
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

# Tasks — A repository says who reviews before the Pull Request

| id      | title                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Carry the pre-Pull-Request reviewer in configuration | backend | medium     | —       |
| task_02 | Report the resolved policy and its source            | backend | medium     | task_01 |
| task_03 | Document the key beside the Review Source            | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04
