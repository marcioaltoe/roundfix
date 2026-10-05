---
schema: spec-tasks/v1
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
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

# Tasks — A queue that sees a Pull Request merged by hand

| id      | title                                                                                                    | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A Delivery Retry records a parked item as merged when a Merge Observer reports its merge                 | backend | medium     | —       |
| task_02 | The queue observes a merge through the recorded Pull Request or the Spec's merge evidence, and says so   | backend | medium     | task_01 |
| task_03 | Run the final QA gate                                                                                    | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_02 implements the observer that task_01's retry asks, and both change
`docs/user-guide/commands/deliver.md`, so they run in series. The gate follows
the only leaf.
