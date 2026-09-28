---
schema: spec-tasks/v1
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
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
| task_01 | backend | An unselected provider receives no request |
| task_02 | backend | Every finding gets a recorded disposition |
| task_03 | backend | A findings verdict stands for its head |
| task_04 | backend | A blocking review after archive parks publication |
| task_05 | qa | Run the final QA gate |
