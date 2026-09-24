---
schema: spec-tasks/v1
spec: 0165-a-measured-run-ceiling
qa: task_03
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
      needs: [task_01, task_02]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | A machine-wide Active Run ceiling |
| task_02 | test | A budget test that does not race the header |
| task_03 | qa | Run the final QA gate |
