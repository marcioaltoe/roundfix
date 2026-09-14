---
schema: spec-tasks/v1
spec: 0137-a-fixture-that-carries-its-committer
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

# Tasks — A fixture that carries its committer

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
| task_01 | Write the committer identity into the shared fixture seed | test | low | — |
| task_02 | Run the final QA gate | qa | low | task_01 |

Waves: 1 → task_01 · 2 → task_02.
