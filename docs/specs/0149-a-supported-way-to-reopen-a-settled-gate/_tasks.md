---
schema: spec-tasks/v1
spec: 0149-a-supported-way-to-reopen-a-settled-gate
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
| task_01 | backend | Expose the QA Task behind a stale-gate refusal |
| task_02 | backend | The reopen command and its refusals |
| task_03 | docs | Keep the shipped skill and the guide true |
| task_04 | qa | Run the final QA gate |
