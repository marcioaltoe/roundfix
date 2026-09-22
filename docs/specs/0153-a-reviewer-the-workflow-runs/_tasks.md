---
schema: spec-tasks/v1
spec: 0153-a-reviewer-the-workflow-runs
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
| task_01 | backend | The review record and its writer |
| task_02 | backend | The command, the Codex path and every refusal |
| task_03 | docs | Describe the command in the shipped skill and the guide |
| task_04 | qa | Run the final QA gate |
