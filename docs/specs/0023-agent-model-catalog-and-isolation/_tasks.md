---
schema: spec-tasks/v1
spec: 0023-agent-model-catalog-and-isolation
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
      needs: [task_01, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks - Agent Model Catalog and Isolation

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Resolve per-runtime Agent selections | backend | high | - |
| task_02 | Apply explicit selection to Agent Sessions | backend | high | task_01 |
| task_03 | Reject unavailable selections during Preflight Validation | infra | high | task_02 |
| task_04 | Expose per-Run selection controls | frontend | high | task_01, task_03 |
| task_05 | Persist effective selection for Run inspection | data | high | task_04 |
| task_06 | Ship Agent selection guidance | docs | medium | task_05 |

Waves: 1 -> task_01 · 2 -> task_02 · 3 -> task_03 · 4 -> task_04 · 5 -> task_05 · 6 -> task_06
