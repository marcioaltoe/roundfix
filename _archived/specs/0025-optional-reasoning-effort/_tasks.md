---
schema: spec-tasks/v1
spec: 0025-optional-reasoning-effort
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
      needs: [task_01]
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
---

# Tasks — Optional Reasoning Effort

| id      | title                                                | type    | complexity | needs            |
| ------- | ---------------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Select the Agent without a reasoning option          | backend | medium     | —                |
| task_02 | Surface model-managed reasoning across CLI commands  | backend | medium     | task_01          |
| task_03 | Accept model-managed reasoning in Interactive Input  | backend | low        | task_01          |
| task_04 | Ship gpt-5.6-sol config and optional-effort guidance | docs    | low        | task_02, task_03 |

Waves: 1 → task_01 · 2 → task_02, task_03 · 3 → task_04
