---
schema: spec-tasks/v1
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
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

# Tasks — A corrective Task after a completed gate passes the Wave collision check

| id      | title                                                                      | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Every Wave collision message names the edge that resolves it               | backend | medium     | —       |
| task_02 | The spec guide and the Roundfix reference describe the Wave collision finding | docs | low        | task_01 |
| task_03 | Run the final QA gate                                                      | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

task_01 builds the collision texts from one helper pair and pins oraculum's
corrective shape; task_02 states the finding it writes in the `spec` guide
and the Roundfix reference, so it follows task_01. The gate follows the last
of them. No Spec is listed in `requires`.
