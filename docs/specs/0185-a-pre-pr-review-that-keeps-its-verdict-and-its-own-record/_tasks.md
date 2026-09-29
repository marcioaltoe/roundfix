---
schema: spec-tasks/v1
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
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
| task_01 | backend | Agent messages keep their boundaries and the review reads the final message |
| task_02 | backend | Each checkout keeps its own pre-PR review record |
| task_03 | qa | Run the final QA gate |
