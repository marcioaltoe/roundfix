---
schema: spec-tasks/v1
spec: 0031-decision-driven-setup-generation
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
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Decision-driven setup generation

| id      | title                                                         | type    | complexity | needs            |
| ------- | ------------------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Introduce validated Decision Plan contracts                    | backend | high       | —                |
| task_02 | Expose conditional setup preview                               | backend | high       | task_01          |
| task_03 | Make Spec scaffolding and external triage decision-driven       | backend | medium     | task_02          |
| task_04 | Make autonomous work and Secondbrain decision-driven            | backend | high       | task_02          |
| task_05 | Render decision values and enforce semantic audit               | backend | high       | task_03, task_04 |
| task_06 | Migrate existing manifests and synchronize the workflow         | backend | high       | task_05          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03, task_04 · 4 → task_05 · 5 → task_06
