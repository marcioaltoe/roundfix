---
schema: spec-tasks/v1
spec: 0162-a-durable-repository-key-per-run
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
      needs: [task_01, task_02, task_04, task_05]
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Record the repository key |
| task_02 | backend | Every lookup uses the key |
| task_04 | backend | Every lookup, and only the recorded key |
| task_05 | backend | Removed worktrees and bare layouts |
| task_03 | qa | Run the final QA gate |
