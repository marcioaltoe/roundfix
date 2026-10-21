---
schema: spec-tasks/v1
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
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

# Tasks — The archive rewrites the satisfied-by links it records as unproven

| id      | title                                                                           | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The archive rebases each unproven entry it records                              | backend | medium     | —       |
| task_02 | The Record Link Repair and the strict check point a promoted target at its copy | backend | medium     | task_01 |
| task_03 | The Roundfix skill references state both rewrites                              | docs    | low        | task_02 |
| task_04 | Run the final QA gate                                                           | qa      | high       | task_03 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04

task_01 changes the record builder and proves it through the archive and the
strict check; task_02 changes the repair and the check's fix line and proves
them through `roundfix history sanitize` and `roundfix baseline update`. Both
edit `CONTEXT.md` and the Behavior Surface Record, so task_02 follows task_01.
task_03 repeats the guide wording of both in the Roundfix skill and raises its
version once, so it follows task_02. The gate follows the last of them. No
Spec is listed in `requires`.
