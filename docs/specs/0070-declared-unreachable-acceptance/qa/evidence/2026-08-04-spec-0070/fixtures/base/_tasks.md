---
schema: spec-tasks/v1
spec: qa-case
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — QA case

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Prepare acceptance | backend | low | — |
| task_02 | Run QA | qa | low | task_01 |
