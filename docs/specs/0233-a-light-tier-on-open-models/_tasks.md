---
schema: spec-tasks/v1
spec: 0233-a-light-tier-on-open-models
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
---

# Tasks — A light tier on open models

| id      | title                                                                                                                         | type    | complexity | needs            |
| ------- | ----------------------------------------------------------------------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | The guides, the Roundfix Skill and the write-tasks skill describe the light tier and the judge's tasks stage                   | docs    | medium     | —                |
| task_02 | User Config carries the light models and the light ceiling, and a tier rule and a Light Spend Log decide and record light Tasks | backend | medium     | —                |
| task_03 | Light Tasks dispatch on an open model with the implementation key, record their spend and escalate once to their profile        | backend | high       | task_01, task_02 |
| task_04 | The advisory judge suggests a model tier for each Task at the tasks stage                                                       | backend | medium     | task_01          |
| task_05 | Run the final QA gate                                                                                                         | qa      | high       | task_03, task_04 |

Waves: 1 → task_01, task_02 · 2 → task_03, task_04 · 3 → task_05

task_01 writes the guides and skills first, so the CLI changes ship with
their record; task_04 follows it because both describe the judge's `tasks`
stage, and task_03 follows it so the engine change ships with the guide it
cites. task_03 also needs the tier rule, the keys, the key helper and the Light
Spend Log that task_02 adds. No two Tasks declare the same file. No Task is
`low`, so none of this Spec's own Tasks would run on the light tier. The gate
follows both leaves.
