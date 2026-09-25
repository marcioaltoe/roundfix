---
schema: spec-tasks/v1
spec: 0170-authoring-that-fails-before-dispatch
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
      needs: []
    - id: task_04
      file: task_04.md
      needs: []
    - id: task_05
      file: task_05.md
      needs: [task_01, task_02, task_03, task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Refuse an undeclared Governed Path at authoring |
| task_02 | backend | Report a CLI change that ships without its guide |
| task_03 | backend | Close three checker honesty gaps and pin the closure cases |
| task_04 | docs | Make the authoring templates and skills produce what the checker accepts |
| task_05 | qa | Run the final QA gate |
