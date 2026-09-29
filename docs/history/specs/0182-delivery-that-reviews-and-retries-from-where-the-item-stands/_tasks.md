---
schema: spec-tasks/v1
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
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
      needs: [task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | The pre-PR review diffs the candidate from its merge base |
| task_02 | backend | Task Carry-Forward treats a Task completed on its target as nothing to carry |
| task_03 | backend | A Delivery Retry carries forward from every Run of its item, newest first |
| task_04 | qa | Run the final QA gate |
