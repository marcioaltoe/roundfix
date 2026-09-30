---
schema: spec-tasks/v1
spec: 0187-a-queue-that-recovers-without-a-supervisor
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
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Post-merge cleanup refreshes the default branch before it reads it |
| task_02 | backend | A retry refused by a Spec amendment prints how to recover it |
| task_03 | backend | An item started by an older queue owner records a warning |
| task_04 | backend | A Task commit is authorized by the grant it ran under |
| task_05 | qa | Run the final QA gate |
