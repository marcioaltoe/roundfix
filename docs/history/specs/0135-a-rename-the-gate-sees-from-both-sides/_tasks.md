---
schema: spec-tasks/v1
spec: 0135-a-rename-the-gate-sees-from-both-sides
qa: task_03
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
---

# Tasks — A rename the gate sees from both sides

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Carry a rename's source through the changed-path reader | backend | medium | — |
| task_02 | Report an unresolvable revision as an unavailable revision | backend | low | — |
| task_03 | Run the final QA gate | qa | medium | task_01, task_02 |

Waves: 1 → task_01, task_02 · 2 → task_03.
