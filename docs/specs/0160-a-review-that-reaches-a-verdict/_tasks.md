---
schema: spec-tasks/v1
spec: 0160-a-review-that-reaches-a-verdict
qa: task_05
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
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Classify the verdict by substance and keep the answer |
| task_02 | backend | Run the claude provider |
| task_03 | backend | Review a Spec's delivery against its decisions |
| task_04 | docs | Keep the shipped skill and the guide true |
| task_05 | qa | Run the final QA gate |
