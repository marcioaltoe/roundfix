---
schema: spec-tasks/v1
spec: 0166-docs-changes-run-the-tests-that-read-them
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

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | No path selects nothing |
| task_02 | qa | Run the final QA gate |
