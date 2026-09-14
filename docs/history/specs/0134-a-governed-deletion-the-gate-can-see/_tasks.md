---
schema: spec-tasks/v1
spec: 0134-a-governed-deletion-the-gate-can-see
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_02, task_03]
---

# Tasks — A governed deletion the gate can see

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Characterize classification for addition, removal and rename | test | low | — |
| task_02 | Classify a removed Governed Path as a mutation | backend | medium | task_01 |
| task_03 | Let an enumerated regeneration list be the whole list | backend | medium | — |
| task_04 | Run the final QA gate | qa | medium | task_02, task_03 |

Waves: 1 → task_01, task_03 · 2 → task_02 · 3 → task_04.
