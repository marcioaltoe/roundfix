---
schema: spec-tasks/v1
spec: 0136-a-rename-the-committer-can-stage
qa: task_04
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Tasks — A rename the committer can stage

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Stage only paths Git can match | backend | medium | — |
| task_02 | Derive the final push's authority from the changed paths | backend | low | — |
| task_03 | Carry a resolvable record path into the unresolved result | backend | low | — |
| task_04 | Run the final QA gate | qa | medium | task_01, task_02, task_03 |

Waves: 1 → task_01, task_02, task_03 · 2 → task_04.
