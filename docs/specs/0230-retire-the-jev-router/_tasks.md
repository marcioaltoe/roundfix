---
schema: spec-tasks/v1
spec: 0230-retire-the-jev-router
qa: task_04
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: [task_03]
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_02]
---

# Tasks — Retire the Jev Router

| id      | title                                                                                                          | type    | complexity | needs            |
| ------- | -------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Agent sessions no longer reach OpenRouter, and the runner refuses the models the subscription rule reserves   | backend | high       | task_03          |
| task_02 | Configuration refuses the models the subscription rule reserves, and the router's floor key degrades to a warning | backend | medium     | task_01          |
| task_03 | The guides, the Roundfix Skill and the Project Config comment state the subscription rule instead of the router | docs    | medium     | —                |
| task_04 | Run the final QA gate                                                                                          | qa      | high       | task_02          |

Waves: 1 → task_03 · 2 → task_01 · 3 → task_02 · 4 → task_04

task_03 writes the guide and the Roundfix Skill reference first, so the CLI
engine change of task_01 ships with its guide (the skill-sync rule), as Spec
0229 ordered its documentation Task. task_02 needs task_01 because both edit
the config package's profile code and task_02's configuration rule calls the
predicate task_01 adds. The gate follows the single leaf, task_02.
