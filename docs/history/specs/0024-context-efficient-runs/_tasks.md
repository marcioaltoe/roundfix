---
schema: spec-tasks/v1
spec: 0024-context-efficient-runs
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
      needs: []
    - id: task_05
      file: task_05.md
      needs: []
    - id: task_06
      file: task_06.md
      needs: [task_02, task_03, task_04, task_05]
---

# Tasks - Context-Efficient Runs

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Capture failed Verification diagnostics | backend | high | - |
| task_02 | Allow one same-session Verification repair | backend | high | task_01 |
| task_03 | Emit the Supervisor Run Event Stream | backend | high | task_01 |
| task_04 | Compact Agent reads and edits | backend | medium | - |
| task_05 | Build bounded Spec Context Bundles | backend | high | - |
| task_06 | Ship context-efficient Run guidance | docs | medium | task_02, task_03, task_04, task_05 |

Waves: 1 -> task_01, task_04, task_05 · 2 -> task_02, task_03 · 3 -> task_06
