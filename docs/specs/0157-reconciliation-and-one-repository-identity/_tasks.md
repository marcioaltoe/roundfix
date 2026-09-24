---
schema: spec-tasks/v1
spec: 0157-reconciliation-and-one-repository-identity
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
      needs: [task_01, task_02, task_03, task_05, task_06]
    - id: task_05
      file: task_05.md
      needs: [task_03]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Refuse a candidate without evidence |
| task_02 | backend | Report identity survives archiving |
| task_03 | backend | One identity for every worktree |
| task_04 | qa | Run the final QA gate |
| task_05 | backend | Keep each Run's checkout, share only the identity |
| task_06 | backend | Same report means same content |
